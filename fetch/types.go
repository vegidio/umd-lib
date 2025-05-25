package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Request

type Request struct {
	Url      string
	FilePath string

	httpReq *http.Request
}

// Response

type Response struct {
	Request    *Request
	StatusCode int
	Size       int64
	Downloaded int64
	Progress   float64
	Done       chan struct{} `json:"-"`

	cancel context.CancelFunc
	err    error
}

// Error waits for the download to complete and returns any error that occurred during the process.
func (r *Response) Error() error {
	<-r.Done
	return r.err
}

// IsComplete checks if the download process is complete.
func (r *Response) IsComplete() bool {
	select {
	case <-r.Done:
		return true
	default:
		return false
	}
}

// Cancel stops the download process.
func (r *Response) Cancel() {
	r.cancel()
}

// Bytes read the file specified in the Request's FilePath and return its content as a byte slice.
// It returns an error if the file cannot be read.
func (r *Response) Bytes() ([]byte, error) {
	data, err := os.ReadFile(r.Request.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// Track monitors the progress of the download and invokes the provided callback function with the current progress
// details. The callback receives the number of bytes downloaded, the total size of the file, and the progress
// percentage.
//
// # Parameters:
//   - callback: A function that takes three arguments: completed bytes (int64), total bytes (int64),
//     and progress percentage (float64).
//
// # Returns:
//   - An error if one occurred during the download process.
func (r *Response) Track(callback func(completed, total int64, progress float64)) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	oldValue := int64(-1)

	for {
		select {
		case <-ticker.C:
			if r.Downloaded != oldValue {
				oldValue = r.Downloaded
				callback(r.Downloaded, r.Size, r.Progress)
			}

		case <-r.Done:
			if r.Downloaded != oldValue {
				oldValue = r.Downloaded
				callback(r.Downloaded, r.Size, r.Progress)
			}
			return r.Error()
		}
	}
}

// ProgressWriter

type progressWriter struct {
	file     io.Writer
	callback func(downloaded int64)
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.file.Write(p)
	if err != nil {
		return n, err
	}

	if pw.callback != nil {
		pw.callback(int64(n))
	}

	return n, nil
}
