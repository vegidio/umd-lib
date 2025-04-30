package fetch

import (
	"fmt"
	"github.com/dromara/dongle"
	log "github.com/sirupsen/logrus"
	"github.com/zeebo/blake3"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// DownloadFile downloads a single file based on the provided request.
//
// Parameters:
//   - request: a Request object containing the details of the file to download.
//
// Returns:
//   - A Response object that contains the status and details of the download process.
func (f Fetch) DownloadFile(request *Request) *Response {
	response := &Response{
		Request: request,
		Done:    make(chan struct{}, 1),
	}

	go func() {
		defer close(response.Done)

		_, err := f.createHttpRequest(request)
		if err != nil {
			response.err = fmt.Errorf("failed to create request: %w", err)
			return
		}

		response.Request = request

		// File Writer
		file, err := os.Create(request.FilePath)
		if err != nil {
			response.err = fmt.Errorf("failed to create file: %w", err)
			return
		}

		defer file.Close()

		f.sendAndRetry(response, file)
	}()

	return response
}

// DownloadFiles downloads multiple files concurrently.
//
// Parameters:
//   - requests: a slice of *Request objects representing the files to download.
//   - parallel: the maximum number of concurrent downloads.
//
// Returns:
//   - A channel of *Response objects, where each response corresponds to a file download.
func (f Fetch) DownloadFiles(requests []*Request, parallel int) <-chan *Response {
	result := make(chan *Response)

	go func() {
		defer close(result)

		var wg sync.WaitGroup
		sem := make(chan struct{}, parallel)

		for _, request := range requests {
			wg.Add(1)

			go func(r *Request) {
				defer func() {
					<-sem
					wg.Done()
				}()

				sem <- struct{}{}
				resp := f.DownloadFile(request)
				result <- resp

				// Waiting for the download the complete before continuing
				_ = resp.Error()
			}(request)
		}

		wg.Wait()
		close(sem)
	}()

	return result
}

// TrackProgress monitors the progress of a file download and invokes a callback function.
//
// Parameters:
//   - resp: a *Response representing the download response to track.
//   - callback: a function that takes three parameters (completed, total, progress) and is called
//     whenever the download progress updates.
//
// Returns:
//   - error: an error if the download fails or completes with an error.
func TrackProgress(
	resp *Response,
	callback func(completed, total int64, progress float64),
) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	oldValue := int64(-1)

	for {
		select {
		case <-ticker.C:
			downloaded := resp.Downloaded
			if downloaded != oldValue {
				oldValue = downloaded
				callback(downloaded, resp.Size, resp.Progress)
			}

		case <-resp.Done:
			return resp.Error()
		}
	}
}

// region - Private functions

func (f Fetch) createHttpRequest(request *Request) (*http.Request, error) {
	req, err := http.NewRequest("GET", request.Url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add the user agent to the request
	req.Header.Set("User-Agent", userAgent)

	// Add any additional headers to the request
	for key, value := range f.headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (f Fetch) sendAndRetry(response *Response, file *os.File) {
	client := &http.Client{}
	var resp *http.Response
	var err error

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	// Hash Writer
	hasher := blake3.New()

	for attempt := 0; attempt <= f.retries; attempt++ {
		// It means we already sent the request at least once, so the subsequent requests should have some sort of
		// backoff strategy before sending the request again.
		if attempt > 0 {
			sleep := time.Duration(fibonacci(attempt+1)) * time.Second

			log.WithFields(log.Fields{
				"attempt": attempt,
				"error":   err,
				"url":     response.Request.Url,
			}).Warn("failed to download file; retrying in ", sleep)

			time.Sleep(sleep)
		}

		req, _ := f.createHttpRequest(response.Request)

		resp, err = client.Do(req)
		if err != nil {
			err = fmt.Errorf("failed to send request: %w", err)
			continue
		}

		response.StatusCode = resp.StatusCode
		response.Size = resp.ContentLength

		pw := &progressWriter{
			file:  file,
			total: resp.ContentLength,
			callback: func(downloaded int64, progress float64) {
				response.Downloaded = downloaded
				response.Progress = progress
			},
		}

		// Multi Writer
		mw := io.MultiWriter(pw, hasher)

		_, err = io.Copy(mw, resp.Body)
		if err != nil {
			err = fmt.Errorf("failed to copy response body: %w", err)
			continue
		}

		sum := hasher.Sum(nil)
		response.Hash = dongle.Encode.FromBytes(sum).ByBase91().ToString()

		break
	}

	response.err = err
	response.Done <- struct{}{}
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// endregion
