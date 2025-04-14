package extractors

import (
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd-lib"
	"testing"
)

func TestFapello_QueryPosts(t *testing.T) {
	const NumberOfPosts = 91

	u := umd.New(nil, nil)
	extractor, _ := u.FindExtractor("https://fapello.com/darja-sobakinskaja/")
	resp, err := extractor.QueryMedia(99999, nil, true)
	media := resp.Media

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(media))
	assert.Equal(t, "model", media[0].Metadata["source"])
	assert.Equal(t, "darja-sobakinskaja", media[0].Metadata["name"])
}
