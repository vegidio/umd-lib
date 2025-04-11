package coomer

import (
	"fmt"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/model"
)

var f = fetch.New(nil, 10)
var baseUrl string

func getUser(service string, user string) <-chan model.Result[Response] {
	result := make(chan model.Result[Response])

	go func() {
		offset := 0
		defer close(result)

		for {
			var posts []Post
			url := fmt.Sprintf(baseUrl+"/api/v1/%s/user/%s?o=%d", service, user, offset)
			resp, err := f.GetResult(url, nil, &posts)

			if err != nil {
				result <- model.Result[Response]{Err: err}
			} else if resp.IsError() {
				result <- model.Result[Response]{Err: fmt.Errorf("error fetching user '%s' posts: %s", user,
					resp.Status())}
			}

			if len(posts) == 0 {
				break
			}

			for _, post := range posts {
				response, postErr := getPost(post.Service, post.User, post.Id)
				if postErr != nil {
					result <- model.Result[Response]{Err: postErr}
				}

				result <- model.Result[Response]{Data: *response}
			}

			offset += 50
		}
	}()

	return result
}

func getPost(service string, user string, id string) (*Response, error) {
	var response *Response
	url := fmt.Sprintf(baseUrl+"/api/v1/%s/user/%s/post/%s", service, user, id)
	resp, err := f.GetResult(url, nil, &response)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching user '%s' posts: %s", user, resp.Status())
	}

	return response, nil
}
