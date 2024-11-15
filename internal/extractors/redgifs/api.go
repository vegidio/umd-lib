package redgifs

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

var client = resty.New().
	SetBaseURL("https://api.redgifs.com/")

func getToken() (*Auth, error) {
	var auth *Auth
	url := "v2/auth/temporary"
	headers := map[string]string{
		"Content-Type": "application/json",
		"Origin":       "https://www.redgifs.com",
		"Referer":      "https://www.redgifs.com/",
		"User-Agent":   "UMD",
	}

	resp, err := client.R().
		SetHeaders(headers).
		SetResult(&auth).
		Get(url)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching authorization token: %s", resp.Status())
	}

	return auth, nil
}

func getVideo(token string, videoUrl string, videoId string) (*Video, error) {
	var watch *Video
	url := fmt.Sprintf("v2/gifs/%s?views=yes", videoId)
	headers := map[string]string{
		"Authorization":  token,
		"X-CustomHeader": videoUrl,
		"User-Agent":     "UMD",
	}

	resp, err := client.R().
		SetHeaders(headers).
		SetResult(&watch).
		Get(url)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching video ID '%s': %s", videoId, resp.Status())
	}

	return watch, nil
}
