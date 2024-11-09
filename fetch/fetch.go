package fetch

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type Fetch struct {
	client *resty.Client
}

// region

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
			SetRetryWaitTime(1 * time.Second).
			SetRetryMaxWaitTime(60 * time.Second).
			AddRetryCondition(
				func(r *resty.Response, err error) bool {
					return err != nil || r.StatusCode() == http.StatusTooManyRequests
				},
			),
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
		return "", fmt.Errorf("fetch error - GetString - %v", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("fetch error - GetString - %d", resp.StatusCode())
	}

	return resp.String(), nil
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
		return 0, fmt.Errorf("fetch error - DownloadFile - %v", err)
	}

	if resp.IsError() {
		return 0, fmt.Errorf("fetch error - DownloadFile - %d", resp.StatusCode())
	}

	return resp.Size(), nil
}

// endregion
