package fetch

import (
	"fmt"
	"io"
	"os"
)

// Request

type Request struct {
	Url      string
	FilePath string
}

// Response

type Response struct {
	Request    *Request
	StatusCode int
	Size       int64
	Downloaded int64
	Progress   float64
	Hash       string
	Done       chan struct{}
	err        error
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

// Bytes reads the file specified in the Request's FilePath and returns its content as a byte slice.
// It returns an error if the file cannot be read.
func (r *Response) Bytes() ([]byte, error) {
	data, err := os.ReadFile(r.Request.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// ProgressWriter

type progressWriter struct {
	file       io.Writer
	downloaded int64
	total      int64
	progress   float64
	callback   func(downloaded int64, progress float64)
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.file.Write(p)
	if err != nil {
		return n, err
	}

	pw.downloaded += int64(n)
	if pw.total > 0 {
		pw.progress = float64(pw.downloaded) / float64(pw.total)
	}

	if pw.callback != nil {
		pw.callback(pw.downloaded, pw.progress)
	}

	return n, nil
}
