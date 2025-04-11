package extractors

import (
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/internal/model"
	"testing"
)

func TestCoomer_QueryUser(t *testing.T) {
	const NumberOfPosts = 50

	u, _ := umd.New("https://coomer.su/onlyfans/user/melindalondon", nil, nil)
	resp, err := u.QueryMedia(NumberOfPosts, nil, true)
	media := resp.Media

	assert.NoError(t, err)
	assert.Equal(t, NumberOfPosts, len(media))
	assert.Equal(t, "onlyfans", media[0].Metadata["source"])
	assert.Equal(t, "melindalondon", media[0].Metadata["name"])
}

func TestCoomer_QueryPost(t *testing.T) {
	u, _ := umd.New("https://coomer.su/onlyfans/user/melindalondon/post/357160243", nil, nil)
	resp, err := u.QueryMedia(99999, nil, true)
	media := resp.Media

	assert.NoError(t, err)
	assert.Equal(t, 1, len(media))
	assert.Equal(t, model.Image, media[0].Type)
	assert.Equal(t, "onlyfans", media[0].Metadata["source"])
	assert.Equal(t, "melindalondon", media[0].Metadata["name"])
}
