package coomer

import (
	"fmt"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/model"
)

var f = fetch.New(nil, 10)
var baseUrl string

func getUser(service string, user string) <-chan model.Result[Response] {
	out := make(chan model.Result[Response])

	go func() {
		offset := 0
		defer close(out)

		for {
			var posts []Post
			url := fmt.Sprintf(baseUrl+"/api/v1/%s/user/%s?o=%d", service, user, offset)
			resp, err := f.GetResult(url, nil, &posts)

			if err != nil {
				out <- model.Result[Response]{Err: err}
			} else if resp.IsError() {
				out <- model.Result[Response]{Err: fmt.Errorf("error fetching user '%s' posts: %s", user,
					resp.Status())}
			}

			if len(posts) == 0 {
				break
			}

			for _, post := range posts {
				result := <-getPost(post.Service, post.User, post.Id)
				if result.Err != nil {
					out <- model.Result[Response]{Err: result.Err}
					continue
				}

				out <- model.Result[Response]{Data: result.Data}
			}

			offset += 50
		}
	}()

	return out
}

func getPost(service string, user string, id string) <-chan model.Result[Response] {
	out := make(chan model.Result[Response])

	go func() {
		defer close(out)

		var response Response
		url := fmt.Sprintf(baseUrl+"/api/v1/%s/user/%s/post/%s", service, user, id)
		resp, err := f.GetResult(url, nil, &response)

		if err != nil {
			out <- model.Result[Response]{Err: err}
			return
		} else if resp.IsError() {
			out <- model.Result[Response]{Err: fmt.Errorf("error fetching user '%s' posts: %s",
				user, resp.Status())}
			return
		}

		out <- model.Result[Response]{Data: response}
	}()

	return out
}
