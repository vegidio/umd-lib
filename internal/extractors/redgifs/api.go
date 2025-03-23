package redgifs

import (
	"fmt"
	"github.com/vegidio/umd-lib/fetch"
)

const BaseUrl = "https://api.redgifs.com/"

var f = fetch.New(nil, 0)

func getToken() (*Auth, error) {
	var auth *Auth
	url := BaseUrl + "v2/auth/temporary"
	headers := map[string]string{
		"Content-Type": "application/json",
		"Origin":       "https://www.redgifs.com",
		"Referer":      "https://www.redgifs.com/",
	}

	resp, err := f.GetResult(url, headers, &auth)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching authorization token: %s", resp.Status())
	}

	return auth, nil
}

func getVideo(token string, videoUrl string, videoId string) (*Video, error) {
	var video *Video
	url := BaseUrl + fmt.Sprintf("v2/gifs/%s?views=yes", videoId)
	headers := map[string]string{
		"Authorization":  token,
		"X-CustomHeader": videoUrl,
	}

	resp, err := f.GetResult(url, headers, &video)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching video ID '%s': %s", videoId, resp.Status())
	}

	return video, nil
}
