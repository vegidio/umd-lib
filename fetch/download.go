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

// NewRequest creates a new download request with the specified URL and file path.
//
// Parameters:
//   - url: The URL to download the file from.
//   - filePath: The path where the downloaded file will be saved.
//
// Returns:
//   - A Request object containing the URL and file path.
//   - An error if the request creation fails.
func (f Fetch) NewRequest(url string, filePath string) (*Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range f.headers {
		req.Header.Set(key, value)
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", userAgent)

	// Create a new Request object
	return &Request{
		Url:      url,
		FilePath: filePath,
		httpReq:  req,
	}, nil
}

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

		// File Writer
		file, err := os.Create(request.FilePath)
		if err != nil {
			response.err = fmt.Errorf("failed to create file: %w", err)
			return
		}

		defer file.Close()

		// Progress Writer
		pw := &progressWriter{
			file: file,
			callback: func(downloaded int64) {
				response.Downloaded += downloaded
				if response.Size < response.Downloaded {
					response.Size = response.Downloaded
				}
				if response.Size > 0 {
					response.Progress = float64(response.Downloaded) / float64(response.Size)
				}
			},
		}

		// Hash Writer
		hasher := blake3.New()

		// Multi Writer
		mw := io.MultiWriter(pw, hasher)

		f.sendAndRetry(response, mw)

		sum := hasher.Sum(nil)
		response.Hash = dongle.Encode.FromBytes(sum).ByBase91().ToString()
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

// region - Private functions

func (f Fetch) sendAndRetry(response *Response, mw io.Writer) {
	var resp *http.Response
	var err error

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

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

		request, _ := f.NewRequest(response.Request.Url, response.Request.FilePath)
		response.Request = request

		resp, err = f.httpClient.Do(response.Request.httpReq)
		if err != nil {
			err = fmt.Errorf("failed to send request: %w", err)
			continue
		}

		response.StatusCode = resp.StatusCode
		response.Size = resp.ContentLength

		_, err = io.Copy(mw, resp.Body)
		if err != nil {
			err = fmt.Errorf("failed to copy response body: %w", err)
			continue
		}

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
