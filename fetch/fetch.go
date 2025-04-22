package fetch

import (
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"github.com/go-resty/resty/v2"
	"github.com/go-rod/rod"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Fetch struct {
	restClient *resty.Client
	grabClient *grab.Client
	retries    int
}

var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
	"Chrome/133.0.0.0 Safari/537.36"

// New creates a new Fetch instance with specified headers and retry settings.
//
// Parameters:
//   - headers: a map of headers to be set on each request.
//   - retries: the number of retry attempts for failed requests.
func New(headers map[string]string, retries int) Fetch {
	logger := log.New()

	f := resty.New()
	f.SetHeader("User-Agent", userAgent)
	g := grab.NewClient()
	g.UserAgent = userAgent

	return Fetch{
		restClient: f.
			SetLogger(logger).
			SetHeaders(headers).
			SetRetryCount(retries).
			SetRetryWaitTime(0).
			AddRetryCondition(
				func(r *resty.Response, err error) bool {
					if (err != nil || r.IsError()) && r.Request.Attempt <= retries {
						sleep := time.Duration(fibonacci(r.Request.Attempt+1)) * time.Second

						log.WithFields(log.Fields{
							"attempt": r.Request.Attempt,
							"error":   r.Error(),
							"status":  r.StatusCode(),
							"url":     r.Request.URL,
						}).Warn("failed to get data; retrying in ", sleep)

						time.Sleep(sleep)
						return true
					}

					return false
				},
			),

		grabClient: g,
		retries:    retries,
	}
}

// GetText performs a GET request to the specified URL and returns the response body as a string.
//
// Parameters:
//   - url: the URL to send the GET request to.
//
// Returns the response body as a string and an error if the request fails.
func (f Fetch) GetText(url string) (string, error) {
	resp, err := f.restClient.R().
		Get(url)

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"status": resp.StatusCode(),
			"url":    url,
		}).Error("Error getting text")

		return "", err
	}

	if resp.IsError() {
		log.WithFields(log.Fields{
			"error":  err,
			"status": resp.StatusCode(),
			"url":    url,
		}).Error("Error getting text")

		return "", fmt.Errorf(resp.Status())
	}

	return resp.String(), nil
}

// GetResult performs a GET request to the specified URL and unmarshals the response body into the provided result
// interface.
//
// Parameters:
//   - url: the URL to send the GET request to.
//   - result: a pointer to the variable where the response body will be unmarshalled.
//
// Returns:
//   - *resty.Response: the response from the GET request.
//   - error: an error if the request fails or the response indicates an error.
func (f Fetch) GetResult(url string, headers map[string]string, result interface{}) (*resty.Response, error) {
	resp, err := f.restClient.R().
		SetHeaders(headers).
		SetResult(result).
		Get(url)

	if err != nil {
		log.WithFields(log.Fields{
			"error":  resp.Error(),
			"status": resp.StatusCode(),
			"url":    url,
		}).Error("error getting result")

		return resp, err
	}

	if resp.IsError() {
		log.WithFields(log.Fields{
			"error":  err,
			"status": resp.StatusCode(),
			"url":    url,
		}).Error("Error getting result")

		return resp, fmt.Errorf(resp.Status())
	}

	return resp, nil
}

// GetHtml uses the browser to perform a GET request to the specified URL and returns the response body as a string.
//
// Parameters:
//   - browser: the browser instance to use for the request.
//   - url: the URL to send the GET request to.
//   - element: the selector for the element to wait for before returning the response body.
//
// Returns the response body as a string and an error if the request fails.
func (f Fetch) GetHtml(page *rod.Page, url string, element string) (string, error) {
	var html string
	var err error

	for attempt := 1; attempt <= f.retries; attempt++ {
		var el *rod.Element

		// Navigate to the URL
		err = page.Timeout(5 * time.Second).Navigate(url)
		err = replaceError(err, fmt.Errorf("failed to navigate to URL"))

		if err == nil {
			el, err = page.Timeout(5 * time.Second).Element(element)
			err = replaceError(err, fmt.Errorf("failed to get element"))
		}

		if err == nil {
			err = el.Timeout(5 * time.Second).WaitVisible()
			err = replaceError(err, fmt.Errorf("failed waiting for element to be visible"))
		}

		if err == nil {
			html, err = page.Timeout(5 * time.Second).HTML()
			err = replaceError(err, fmt.Errorf("failed to get page HTML"))
		}

		if err == nil {
			break
		} else {
			sleep := time.Duration(fibonacci(attempt+1)) * time.Second

			log.WithFields(log.Fields{
				"attempt": attempt,
				"element": element,
				"error":   err,
				"url":     url,
			}).Warn("Error getting data; retrying in ", sleep)

			time.Sleep(sleep)
		}
	}

	return html, err
}

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
	var wg sync.WaitGroup
	result := make(chan *grab.Response)
	pool, _ := ants.NewPool(parallel)

	go func() {
		defer pool.Release()
		defer close(result)

		for _, request := range requests {
			wg.Add(1)

			_ = pool.Submit(func() {
				defer wg.Done()
				resp := f.DownloadFile(request)
				result <- resp
			})
		}

		wg.Wait()
	}()

	return result
}

// region - Private functions

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func replaceError(err error, fallback error) error {
	if err != nil {
		return fallback
	}
	return err
}

// endregion
