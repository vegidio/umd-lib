package extractors

import (
	"github.com/cavaliergopher/grab/v3"
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/fetch"
	"os"
	"testing"
)

func TestImaglr_DownloadVideo(t *testing.T) {
	// Delete any previous file before continuing
	const FilePath = "video.mp4"
	_ = os.Remove(FilePath)

	u, _ := umd.New("https://imaglr.com/post/5778297", nil, nil)
	resp, _ := u.QueryMedia(99999, nil, true)

	media := resp.Media[0]
	request, _ := grab.NewRequest("video.mp4", media.Url)
	f := fetch.New(nil, 0)
	downloadResponse := f.DownloadFile(request)

	assert.NoError(t, downloadResponse.Err())
	assert.Equal(t, int64(75_520_497), downloadResponse.Size())
}
