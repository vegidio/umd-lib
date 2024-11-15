package reddit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSubmission(t *testing.T) {
	id := "1bxsmnr"
	submissions, _ := getSubmission(id)

	assert.NotNil(t, submissions)
	assert.Greater(t, len(submissions.Data.Children), 0)
}

func TestGetUserSubmissions(t *testing.T) {
	user := "atomicbrunette18"
	after := ""
	limit := 100
	submission, _ := getUserSubmissions(user, after, limit)

	assert.NotNil(t, submission)
	assert.Greater(t, len(submission.Data.Children), 0)
}

func TestGetSubredditSubmissions(t *testing.T) {
	subreddit := "nsfw"
	after := ""
	limit := 100
	submission, _ := getSubredditSubmissions(subreddit, after, limit)

	assert.NotNil(t, submission)
	assert.Greater(t, len(submission.Data.Children), 0)
}
