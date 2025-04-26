package extractors

import (
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd-lib"
	"testing"
)

func TestFapello_QueryPosts(t *testing.T) {
	const NumberOfPosts = 91

	extractor, _ := umd.New(nil).FindExtractor("https://fapello.com/darja-sobakinskaja/")
	resp := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "model", resp.Media[0].Metadata["source"])
	assert.Equal(t, "darja-sobakinskaja", resp.Media[0].Metadata["name"])
}
