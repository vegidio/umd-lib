package extractors

import (
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/internal/model"
	"testing"
)

func TestCoomer_QueryUser(t *testing.T) {
	const NumberOfPosts = 50

	extractor, _ := umd.New(nil).FindExtractor("https://coomer.su/onlyfans/user/melindalondon")
	resp, _ := extractor.QueryMedia(NumberOfPosts, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(resp.Media))
	assert.Equal(t, "onlyfans", resp.Media[0].Metadata["source"])
	assert.Equal(t, "melindalondon", resp.Media[0].Metadata["name"])
}

func TestCoomer_QueryPost(t *testing.T) {
	extractor, _ := umd.New(nil).FindExtractor("https://coomer.su/onlyfans/user/melindalondon/post/357160243")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Media))
	assert.Equal(t, model.Image, resp.Media[0].Type)
	assert.Equal(t, "onlyfans", resp.Media[0].Metadata["source"])
	assert.Equal(t, "melindalondon", resp.Media[0].Metadata["name"])
}
