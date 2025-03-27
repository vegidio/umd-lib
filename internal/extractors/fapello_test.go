package extractors

import (
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd-lib"
	"testing"
)

func TestFapello_QueryPosts(t *testing.T) {
	const NumberOfPosts = 49

	u, _ := umd.New("https://fapello.com/caylinlive-33/", nil, nil)
	resp, err := u.QueryMedia(99999, nil, true)
	media := resp.Media

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(media))
	assert.Equal(t, "model", media[0].Metadata["source"])
	assert.Equal(t, "caylinlive-33", media[0].Metadata["name"])
}
