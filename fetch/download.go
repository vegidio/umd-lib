package fetch

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
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
func (f *Fetch) NewRequest(url string, filePath string) (*Request, error) {
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
func (f *Fetch) DownloadFile(request *Request) *Response {
	ctx, cancel := context.WithCancel(context.Background())

	// Tie the http.Request to the context.
	// It means that if the context is canceled, the request will be canceled too.
	request.httpReq = request.httpReq.WithContext(ctx)

	response := &Response{
		Request: request,
		Done:    make(chan struct{}, 1),
		cancel:  cancel,
	}

	go func() {
		defer close(response.Done)

		// How many bytes are already on the disk?
		var offset int64
		if info, err := os.Stat(request.FilePath); err == nil {
			offset = info.Size()
		}

		// Open (or create) a file for appending
		file, err := os.OpenFile(request.FilePath, os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			response.err = fmt.Errorf("could not open file: %w", err)
			return
		}

		defer file.Close()

		// Seek to the end of existing data
		if _, fErr := file.Seek(offset, io.SeekStart); fErr != nil {
			response.err = fmt.Errorf("could not seek: %w", fErr)
			return
		}

		// Set up the progress callback
		pw := &progressWriter{
			file: file,
			callback: func(downloaded int64) {
				response.Downloaded += downloaded
				if response.Size > 0 {
					response.Progress = float64(response.Downloaded) / float64(response.Size)
				}
			},
		}

		// Perform the download (with resume & retries)
		f.downloadWithRetries(response, offset, file, pw, ctx)
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
//   - A function that can be called to cancel all downloads.
func (f *Fetch) DownloadFiles(requests []*Request, parallel int) (<-chan *Response, func()) {
	result := make(chan *Response)
	done := make(chan struct{})

	var (
		wg      sync.WaitGroup
		sem     = make(chan struct{}, parallel)
		mu      sync.Mutex
		cancels []func()
	)

	// cancelAll cancels all ongoing downloads
	cancelAll := func() {
		select {
		case <-done:
			// already canceled
		default:
			close(done)
		}

		mu.Lock()
		defer mu.Unlock()
		for _, cancelFn := range cancels {
			cancelFn()
		}
	}

	go func() {
		defer close(result)

		for _, req := range requests {
			wg.Add(1)

			go func(r *Request) {
				defer wg.Done()

				// Either grab a slot or exit if canceled
				select {
				case sem <- struct{}{}:
					// slot acquired
				case <-done:
					return
				}

				defer func() { <-sem }()

				// Start the download
				resp := f.DownloadFile(r)

				// Capture the Cancel() function
				mu.Lock()
				cancels = append(cancels, resp.cancel)
				mu.Unlock()

				result <- resp

				// Waiting for the download the complete before continuing
				_ = resp.Error()
			}(req)
		}

		wg.Wait()
	}()

	return result, cancelAll
}

// region - Private functions

func (f *Fetch) downloadWithRetries(
	response *Response,
	offset int64,
	file *os.File,
	writer io.Writer,
	ctx context.Context,
) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= f.retries; attempt++ {
		// Before each attempt, see if we've been canceled
		select {
		case <-ctx.Done():
			response.err = ctx.Err()
			return
		default:
		}

		if attempt > 0 {
			backoff := time.Duration(fibonacci(attempt+1)) * time.Second

			log.WithFields(log.Fields{
				"attempt": attempt,
				"error":   err,
				"url":     response.Request.Url,
			}).Warn("failed to download file; retrying in ", backoff)

			time.Sleep(backoff)
		}

		isRangeReq := offset > 0
		if isRangeReq {
			response.Request.httpReq.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
		} else {
			response.Request.httpReq.Header.Del("Range")
		}

		// Send it
		resp, err = f.httpClient.Do(response.Request.httpReq)
		if err != nil {
			response.err = fmt.Errorf("request error: %w", err)
			continue
		}

		// Fallback if server doesn’t support Range
		if isRangeReq && resp.StatusCode == http.StatusOK {
			// Truncate file and reset offset
			if tErr := file.Truncate(0); tErr != nil {
				response.StatusCode = resp.StatusCode
				response.Size = 0
				response.err = fmt.Errorf("truncate failed: %w", tErr)
			}

			offset = 0

			if _, sErr := file.Seek(0, io.SeekStart); sErr != nil {
				response.StatusCode = resp.StatusCode
				response.Size = 0
				response.err = fmt.Errorf("seek after truncate failed: %w", sErr)
			}

			attempt-- // retry same attempt count with fresh download
			resp.Body.Close()
			continue
		}

		response.Downloaded = offset

		// Handle '416 Range Not Satisfiable' (already complete)
		if resp.StatusCode == http.StatusRequestedRangeNotSatisfiable {
			response.StatusCode = resp.StatusCode
			response.Size = offset
			response.Progress = float64(offset) / float64(response.Size)

			resp.Body.Close()
			break
		}

		// If we get an error (anything that is not 2xx), then we abort this loop and go to the next attempt.
		// We don't do that for HTTP 404 and 410, because those are cases where we know the file is not there.
		if resp.StatusCode != 404 && resp.StatusCode != 410 && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
			response.err = fmt.Errorf("unexpected status: %d", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		// Compute total size from Content-Range or Content-Length
		if cr := resp.Header.Get("Content-Range"); cr != "" {
			// e.g. "bytes 500-999/1234"
			var start, end, total int64
			if _, scanErr := fmt.Sscanf(cr, "bytes %d-%d/%d", &start, &end, &total); scanErr == nil {
				response.Size = total
			} else {
				response.Size = offset + resp.ContentLength
			}
		} else {
			response.Size = offset + resp.ContentLength
		}

		// Track where this attempt started
		startOffset := offset

		// Actually copy data
		_, err = io.Copy(writer, resp.Body)
		if err != nil {
			if ctx.Err() != nil {
				response.err = ctx.Err()
				resp.Body.Close()
				break
			}

			// figure out how many bytes actually made it to the disk
			newOffset, seekErr := file.Seek(0, io.SeekEnd)
			if seekErr != nil {
				response.err = fmt.Errorf("seek after partial download failed: %w", seekErr)
				resp.Body.Close()
				break
			}

			// bump offset to resume after what we have
			offset = newOffset
			response.err = fmt.Errorf("download interrupted (wrote %d bytes), will resume: %w",
				newOffset-startOffset, err)
			resp.Body.Close()
			continue
		}

		// Success
		if response.Size == -1 {
			response.Size = response.Downloaded
			response.Progress = 1
		}

		response.StatusCode = resp.StatusCode
		response.err = nil
		resp.Body.Close()
		break
	}
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// endregion
