package extractors

import (
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

	extractor, _ := umd.New(nil).FindExtractor("https://imaglr.com/post/5778297")
	resp, _ := extractor.QueryMedia(99999, nil, true)
	<-resp.Done

	media := resp.Media[0]
	f := fetch.New(nil, 0)
	request, _ := f.NewRequest(media.Url, "video.mp4")
	downloadResponse := f.DownloadFile(request)

	assert.NoError(t, downloadResponse.Error())
	assert.Equal(t, int64(75_520_497), downloadResponse.Size)
}
