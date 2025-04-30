package fetch

import (
	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// DownloadFile attempts to download a file.
//
// Parameters:
//   - request: a *grab.Request containing the file URL and other settings.
//
// Returns:
//   - *grab.Response: the response from the download.
func (f Fetch) DownloadFile(request *grab.Request) *grab.Response {
	return f.grabClient.Do(request)
}

// DownloadFiles downloads multiple files concurrently using a worker pool.
//
// Parameters:
//   - requests: a slice of *grab.Request objects representing the files to be downloaded.
//   - parallel: the number of concurrent workers to use for downloading.
//
// Returns:
//   - <-chan *grab.Response: a channel through which the download responses are sent.
func (f Fetch) DownloadFiles(requests []*grab.Request, parallel int) <-chan *grab.Response {
	result := make(chan *grab.Response)

	go func() {
		defer close(result)

		var wg sync.WaitGroup
		sem := make(chan struct{}, parallel)

		for _, request := range requests {
			wg.Add(1)

			go func(r *grab.Request) {
				defer func() {
					<-sem
					wg.Done()
				}()

				sem <- struct{}{}
				resp := f.DownloadFile(r)
				result <- resp

				// Waiting for the download the complete before continuing
				_ = resp.Err()
			}(request)
		}

		wg.Wait()
		close(sem)
	}()

	return result
}

// RetryDownload retries a failed file download using an exponential backoff strategy.
//
// Parameters:
//   - response: a *grab.Response representing the initial failed download attempt.
//
// Returns:
//   - <-chan *grab.Response: a channel that emits the updated *grab.Response for each retry attempt.
func (f Fetch) RetryDownload(response *grab.Response) <-chan *grab.Response {
	result := make(chan *grab.Response)

	go func() {
		defer close(result)

		// If the download is not completed yet we wait until it's.
		// Also, if there's no error after it completes then we don't need to retry.
		if !response.IsComplete() {
			result <- response
			if response.Err() == nil {
				return
			}
		}

		for attempt := 1; attempt <= f.retries; attempt++ {
			request, _ := grab.NewRequest(response.Filename, response.Request.URL().String())
			sleep := time.Duration(fibonacci(attempt+1)) * time.Second

			log.WithFields(log.Fields{
				"attempt": attempt,
				"error":   response.Err(),
				"url":     request.URL(),
			}).Warn("failed to download file; retrying in ", sleep)

			time.Sleep(sleep)

			response = f.grabClient.Do(request)
			result <- response

			if response.Err() == nil {
				break
			}
		}
	}()

	return result
}

// TrackProgress monitors the progress of a file download and invokes a callback function.
//
// Parameters:
//   - resp: a *grab.Response representing the download response to track.
//   - callback: a function that takes three parameters (completed, total, progress) and is called
//     whenever the download progress updates.
//
// Returns:
//   - error: an error if the download fails or completes with an error.
func TrackProgress(
	resp *grab.Response,
	callback func(completed, total int64, progress float64),
) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	oldValue := int64(-1)

	for {
		select {
		case <-ticker.C:
			completed := resp.BytesComplete()
			if completed != oldValue {
				oldValue = completed
				callback(completed, resp.Size(), resp.Progress())
			}

		case <-resp.Done:
			return resp.Err()
		}
	}
}

// region - Private functions

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// endregion
