package extractors

import (
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/internal/model"
	"testing"
)

func TestCoomer_QueryUser(t *testing.T) {
	const NumberOfPosts = 50

	u := umd.New(nil, nil)
	extractor, _ := u.FindExtractor("https://coomer.su/onlyfans/user/melindalondon")
	resp, err := extractor.QueryMedia(NumberOfPosts, nil, true)
	media := resp.Media

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(media))
	assert.Equal(t, "onlyfans", media[0].Metadata["source"])
	assert.Equal(t, "melindalondon", media[0].Metadata["name"])
}

func TestCoomer_QueryPost(t *testing.T) {
	u := umd.New(nil, nil)
	extractor, _ := u.FindExtractor("https://coomer.su/onlyfans/user/melindalondon/post/357160243")
	resp, err := extractor.QueryMedia(99999, nil, true)
	media := resp.Media

	assert.NoError(t, err)
	assert.Equal(t, 1, len(media))
	assert.Equal(t, model.Image, media[0].Type)
	assert.Equal(t, "onlyfans", media[0].Metadata["source"])
	assert.Equal(t, "melindalondon", media[0].Metadata["name"])
}
