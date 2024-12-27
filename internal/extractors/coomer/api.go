package coomer

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func setBaseUrl(baseUrl string) {
	client.SetBaseURL(baseUrl)
}

func getUserPosts(service string, user string) ([]Post, error) {
	posts := make([]Post, 0)
	offset := 0

	for {
		url := fmt.Sprintf("/api/v1/%s/user/%s?o=%d", service, user, offset)
		var newPosts []Post

		resp, err := client.R().
			SetResult(&newPosts).
			Get(url)

		if err != nil {
			return nil, err
		} else if resp.IsError() {
			return nil, fmt.Errorf("error fetching user '%s' posts: %s", user, resp.Status())
		}

		if len(newPosts) == 0 {
			break
		}

		posts = append(posts, newPosts...)
		offset += 50
	}

	return posts, nil
}

func getPost(service string, user string, id string) (*Post, error) {
	url := fmt.Sprintf("/api/v1/%s/user/%s/post/%s", service, user, id)
	var response *Response

	resp, err := client.R().
		SetResult(&response).
		Get(url)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching user '%s' posts: %s", user, resp.Status())
	}

	return response.Post, nil
}
