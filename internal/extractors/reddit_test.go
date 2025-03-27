package extractors

import (
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/internal/model"
	"os"
	"testing"
)

func TestReddit_QuerySubreddit(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("This test doesn't work when executed from GitHub Actions")
	}

	const NumberOfPosts = 50

	u, _ := umd.New("https://www.reddit.com/r/PristineGirls/", nil, nil)
	resp, err := u.QueryMedia(NumberOfPosts, nil, true)
	media := resp.Media

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(media))
	assert.Equal(t, "subreddit", media[0].Metadata["source"])
	assert.Equal(t, "PristineGirls", media[0].Metadata["name"])
}

func TestReddit_QuerySubmissions(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("This test doesn't work when executed from GitHub Actions")
	}

	const NumberOfPosts = 50

	u, _ := umd.New("https://www.reddit.com/user/atomicbrunette18/", nil, nil)
	resp, err := u.QueryMedia(NumberOfPosts, nil, true)
	media := resp.Media

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(media))
	assert.Equal(t, "user", media[0].Metadata["source"])
	assert.Equal(t, "atomicbrunette18", media[0].Metadata["name"])
}

func TestReddit_QuerySingleSubmission(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("This test doesn't work when executed from GitHub Actions")
	}
	
	u, _ := umd.New("https://www.reddit.com/r/needysluts/comments/1aenk3e/if_im_wearing_this_for_our_date_you_have_bo/", nil, nil)
	resp, err := u.QueryMedia(99999, nil, true)
	media := resp.Media

	assert.NoError(t, err)
	assert.Equal(t, 1, len(media))
	assert.Equal(t, model.Video, media[0].Type)
	assert.Equal(t, "submission", media[0].Metadata["source"])
	assert.Equal(t, "needysluts", media[0].Metadata["name"])
}
