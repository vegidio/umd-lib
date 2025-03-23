package reddit

import (
	"fmt"
	"github.com/vegidio/umd-lib/fetch"
)

const BaseUrl = "https://www.reddit.com/"

var f = fetch.New(nil, 10)

// getSubmission retrieves a list of submissions for a given Reddit post ID.
//
// Example: https://www.reddit.com/comments/1bxsmnr.json?raw_json=1, where <1bxsmnr> is the ID.
//
// # Parameters:
//   - id: The ID of the Reddit post.
//
// # Returns:
//   - A slice of Submission structs containing the details of the submissions.
func getSubmission(id string) (*Submission, error) {
	submissions := make([]Submission, 0)
	url := fmt.Sprintf(BaseUrl+"comments/%s.json?raw_json=1", id)
	resp, err := f.GetResult(url, nil, &submissions)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching post id '%s' submissions: %s", id, resp.Status())
	}

	return &submissions[0], nil
}

// getUserSubmissions retrieves a list of submissions for a given Reddit user.
//
// Example: https://www.reddit.com/user/atomicbrunette18/submitted.json?sort=new&raw_json=1&after=&limit=100, where
// <atomicbrunette18> is the username.
//
// # Parameters:
//   - user: The username of the Reddit user.
//   - after: The ID of the last submission to start after (for pagination).
//   - limit: The maximum number of submissions to retrieve.
//
// # Returns:
//   - A Submission struct containing the details of the submissions.
func getUserSubmissions(user string, after string, limit int) (*Submission, error) {
	var submission *Submission
	url := fmt.Sprintf(BaseUrl+"user/%s/submitted.json?sort=new&raw_json=1&after=%s&limit=%d", user, after, limit)
	resp, err := f.GetResult(url, nil, &submission)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching user '%s' submissions: %s", user, resp.Status())
	}

	return submission, nil
}

// getSubredditSubmissions retrieves a list of submissions for a given subreddit.
//
// Example: https://www.reddit.com/r/nsfw/hot.json?raw_json=1&after=&limit=100, where <nsfw> is the subreddit name.
//
// # Parameters:
//   - subreddit: The name of the subreddit.
//   - after: The ID of the last submission to start after (for pagination).
//   - limit: The maximum number of submissions to retrieve.
//
// # Returns:
//   - A Submission struct containing the details of the submissions.
func getSubredditSubmissions(subreddit string, after string, limit int) (*Submission, error) {
	var submission *Submission
	url := fmt.Sprintf(BaseUrl+"r/%s/hot.json?raw_json=1&after=%s&limit=%d", subreddit, after, limit)
	resp, err := f.GetResult(url, nil, &submission)

	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return nil, fmt.Errorf("error fetching subreddit '%s' submissions: %s", subreddit, resp.Status())
	}

	return submission, nil
}
