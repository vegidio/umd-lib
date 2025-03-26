package extractors

import (
	"bytes"
	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/fetch"
	"os"
	"strings"
	"testing"
)

func TestRedGifs_DownloadVideo(t *testing.T) {
	// Delete any previous file before continuing
	const FilePath = "video.mp4"
	_ = os.Remove(FilePath)
	
	u, _ := umd.New("https://www.redgifs.com/watch/sturdycuddlyicefish", nil, nil)
	resp, _ := u.QueryMedia(99999, nil, true)

	media := resp.Media[0]
	request, _ := grab.NewRequest("video.mp4", media.Url)
	f := fetch.New(nil, 0)
	downloadResponse := f.DownloadFile(request)

	assert.NoError(t, downloadResponse.Err())
	assert.Equal(t, int64(15_212_770), downloadResponse.Size())
}

func TestRedGifs_ReuseToken(t *testing.T) {
	// Create a buffer and redirect global log output to it
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.DebugLevel)

	// First query
	u, _ := umd.New("https://www.redgifs.com/watch/sturdycuddlyicefish", nil, nil)
	resp, _ := u.QueryMedia(99999, nil, true)

	// Second query
	u, _ = umd.New("https://www.redgifs.com/watch/ecstaticthickasiansmallclawedotter", resp.Metadata, nil)
	_, _ = u.QueryMedia(99999, nil, true)

	// Check the log output
	output := buf.String()
	assert.Equal(t, 1, strings.Count(output, "Issuing new RedGifs token"))
	assert.Equal(t, 1, strings.Count(output, "Reusing RedGifs token"))
}
