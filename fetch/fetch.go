package fetch

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/go-rod/rod"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Fetch struct {
	client  *resty.Client
	retries int
}

// New creates a new Fetch instance with specified headers and retry settings.
//
// Parameters:
//   - headers: a map of headers to be set on each request.
//   - retries: the number of retry attempts for failed requests.
func New(headers map[string]string, retries int) Fetch {
	return Fetch{
		client: resty.New().
			SetHeaders(headers).
			SetRetryCount(retries).
			SetRetryWaitTime(0).
			AddRetryCondition(
				func(r *resty.Response, err error) bool {
					if r.StatusCode() == http.StatusTooManyRequests {
						sleep := time.Duration(fibonacci(r.Request.Attempt+1)) * time.Second

						log.WithFields(log.Fields{
							"attempt": r.Request.Attempt,
							"url":     r.Request.URL,
						}).Debug("Too many requests; retrying in ", sleep)

						time.Sleep(sleep)
						return true
					}

					return false
				},
			),
		retries: retries,
	}
}

// GetText performs a GET request to the specified URL and returns the response body as a string.
//
// Parameters:
//   - url: the URL to send the GET request to.
//
// Returns the response body as a string and an error if the request fails.
func (f Fetch) GetText(url string) (string, error) {
	resp, err := f.client.R().
		Get(url)

	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Error("Error getting text: ", err)

		return "", err
	}

	if resp.IsError() {
		log.WithFields(log.Fields{
			"status": resp.StatusCode(),
		}).Error("Error getting text: ", resp.Status())

		return "", fmt.Errorf(resp.Status())
	}

	return resp.String(), nil
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
			}).Warn("Error getting HTML; retrying in ", sleep)

			time.Sleep(sleep)
		}
	}

	return html, err
}

// DownloadFile performs a GET request to the specified URL and saves the response body to the specified file path.
//
// Parameters:
//   - url: the URL to send the GET request to.
//   - filePath: the path where the response body should be saved.
//
// Returns the size of the downloaded file and an error if the request fails.
func (f Fetch) DownloadFile(url string, filePath string) (int64, error) {
	resp, err := f.client.R().
		SetOutput(filePath).
		Get(url)

	if err != nil {
		log.WithFields(log.Fields{
			"filePath": filePath,
			"url":      url,
		}).Debug("Error downloading file: ", err)

		return 0, err
	}

	if resp.IsError() {
		log.WithFields(log.Fields{
			"status": resp.StatusCode(),
			"url":    url,
		}).Debug("Error downloading file: ", resp.Status())

		return 0, fmt.Errorf(resp.Status())
	}

	return resp.Size(), nil
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
