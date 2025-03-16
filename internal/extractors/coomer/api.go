package coomer

import (
	"fmt"
	"github.com/vegidio/umd-lib/fetch"
)

var f = fetch.New(nil, 10)
var baseUrl string

func getUserPosts(service string, user string) ([]Post, error) {
	posts := make([]Post, 0)
	offset := 0

	for {
		var newPosts []Post
		url := fmt.Sprintf(baseUrl+"/api/v1/%s/user/%s?o=%d", service, user, offset)
		resp, err := f.GetResult(url, &newPosts)

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
	var response *Response
	url := fmt.Sprintf(baseUrl+"/api/v1/%s/user/%s/post/%s", service, user, id)
	resp, err := f.GetResult(url, &response)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching user '%s' posts: %s", user, resp.Status())
	}

	return response.Post, nil
}
