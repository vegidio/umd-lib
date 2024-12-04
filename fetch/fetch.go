package fetch

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
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
						}).Debug("Too many requests - Retrying in ", sleep)

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
		log.Error(err)
		return "", fmt.Errorf("fetch error - GetText - %v", err)
	}

	if resp.IsError() {
		log.Error(err)
		return "", fmt.Errorf("fetch error - GetText - %d", resp.StatusCode())
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
func (f Fetch) GetHtml(browser *rod.Browser, url string, element string) (string, error) {
	var e proto.NetworkResponseReceived

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer page.Close()
	context := browser.GetContext()

	for attempt := 1; attempt <= f.retries; attempt++ {
		wait := page.Context(context).WaitEvent(&e)
		err = page.Timeout(10 * time.Second).Navigate(url)
		if err != nil {
			log.Error(err)
			return "", err
		}

		wait()

		if e.Response.Status == http.StatusTooManyRequests {
			sleep := time.Duration(fibonacci(attempt+1)) * time.Second

			log.WithFields(log.Fields{
				"attempt": attempt,
				"url":     url,
			}).Debug("Too many requests - Retrying in ", sleep)

			time.Sleep(sleep)
		} else {
			log.WithFields(log.Fields{
				"status": e.Response.Status,
				"url":    url,
			}).Debug("Successful response")

			break
		}
	}

	el, err := page.Element(element)
	if err != nil {
		log.Error(err)
		return "", err
	}

	err = el.Timeout(10 * time.Second).WaitVisible()
	if err != nil {
		log.Error(err)
		return "", err
	}

	html, err := page.HTML()
	if err != nil {
		log.Error(err)
		return "", err
	}

	return html, nil
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
		log.Error(err)
		return 0, fmt.Errorf("fetch error - DownloadFile - %v", err)
	}

	if resp.IsError() {
		err = fmt.Errorf("fetch error - DownloadFile - %d", resp.StatusCode())
		log.Error(err)
		return 0, err
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

// endregion
