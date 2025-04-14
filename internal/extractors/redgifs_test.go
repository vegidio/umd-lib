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

	u := umd.New(nil, nil)
	extractor, _ := u.FindExtractor("https://www.redgifs.com/watch/sturdycuddlyicefish")
	resp, _ := extractor.QueryMedia(99999, nil, true)

	media := resp.Media[0]
	request, _ := grab.NewRequest("video.mp4", media.Url)
	f := fetch.New(nil, 0)
	downloadResponse := f.DownloadFile(request)

	assert.NoError(t, downloadResponse.Err())
	assert.Equal(t, int64(15_212_770), downloadResponse.Size())
	assert.Equal(t, "sturdycuddlyicefish", media.Metadata["id"])
	assert.Equal(t, "sonya_18yo", media.Metadata["name"])
}

func TestRedGifs_FetchUser(t *testing.T) {
	u := umd.New(nil, nil)
	extractor, _ := u.FindExtractor("https://www.redgifs.com/users/atomicbrunette18")
	resp, err := extractor.QueryMedia(180, nil, true)

	media := resp.Media[0]

	assert.NoError(t, err)
	assert.Equal(t, 180, len(resp.Media))
	assert.Equal(t, "user", media.Metadata["source"])
	assert.Equal(t, "atomicbrunette18", media.Metadata["name"])
}

func TestRedGifs_ReuseToken(t *testing.T) {
	// Create a buffer and redirect global log output to it
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.DebugLevel)

	// First query
	u := umd.New(nil, nil)
	extractor, _ := u.FindExtractor("https://www.redgifs.com/watch/sturdycuddlyicefish")
	resp, _ := extractor.QueryMedia(99999, nil, true)

	// Second query
	u = umd.New(resp.Metadata, nil)
	extractor, _ = u.FindExtractor("https://www.redgifs.com/watch/ecstaticthickasiansmallclawedotter")
	_, _ = extractor.QueryMedia(99999, nil, true)

	// Check the log output
	output := buf.String()
	assert.Equal(t, 1, strings.Count(output, "Issuing new RedGifs token"))
	assert.Equal(t, 1, strings.Count(output, "Reusing RedGifs token"))
}
