package reddit

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

var client = resty.New().
	SetBaseURL("https://www.reddit.com/")

// GetSubmission retrieves a list of submissions for a given Reddit post ID.
//
// Example: https://www.reddit.com/comments/1bxsmnr.json?raw_json=1, where <1bxsmnr> is the ID.
//
// Parameters:
//   - id: The ID of the Reddit post.
//
// Returns:
//   - A slice of Submission structs containing the details of the submissions.
func getSubmission(id string) []Submission {
	submissions := make([]Submission, 0)
	url := fmt.Sprintf("comments/%s.json?raw_json=1", id)

	_, _ = client.R().
		SetResult(&submissions).
		Get(url)

	return submissions
}

// GetUserSubmissions retrieves a list of submissions for a given Reddit user.
//
// Example: https://www.reddit.com/user/atomicbrunette18/submitted.json?sort=new&raw_json=1&after=&limit=100, where
// <atomicbrunette18> is the username.
//
// Parameters:
//   - user: The username of the Reddit user.
//   - after: The ID of the last submission to start after (for pagination).
//   - limit: The maximum number of submissions to retrieve.
//
// Returns:
//   - A Submission struct containing the details of the submissions.
func getUserSubmissions(user string, after string, limit int) Submission {
	var submission Submission
	url := fmt.Sprintf("user/%s/submitted.json?sort=new&raw_json=1&after=%s&limit=%d", user, after, limit)

	_, _ = client.R().
		SetResult(&submission).
		Get(url)

	return submission
}
