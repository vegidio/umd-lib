package extractors

import (
	"bytes"
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

	extractor, _ := umd.New(nil).FindExtractor("https://www.redgifs.com/watch/sturdycuddlyicefish")
	resp := extractor.QueryMedia(99999, nil, true)
	<-resp.Done

	media := resp.Media[0]
	request := &fetch.Request{Url: media.Url, FilePath: "video.mp4"}
	f := fetch.New(nil, 0)
	downloadResponse := f.DownloadFile(request)

	assert.NoError(t, downloadResponse.Error())
	assert.Equal(t, int64(15_212_770), downloadResponse.Size)
	assert.Equal(t, "sturdycuddlyicefish", media.Metadata["id"])
	assert.Equal(t, "sonya_18yo", media.Metadata["name"])
}

func TestRedGifs_FetchUser(t *testing.T) {
	extractor, _ := umd.New(nil).FindExtractor("https://www.redgifs.com/users/atomicbrunette18")
	resp := extractor.QueryMedia(180, nil, true)
	err := resp.Error()

	assert.NoError(t, err)
	assert.Equal(t, 180, len(resp.Media))
	assert.Equal(t, "user", resp.Media[0].Metadata["source"])
	assert.Equal(t, "atomicbrunette18", resp.Media[0].Metadata["name"])
}

func TestRedGifs_ReuseToken(t *testing.T) {
	// Create a buffer and redirect global log output to it
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.DebugLevel)

	// First query
	u := umd.New(nil)
	extractor, _ := u.FindExtractor("https://www.redgifs.com/watch/sturdycuddlyicefish")
	r1 := extractor.QueryMedia(99999, nil, true)
	<-r1.Done

	// Second query
	u = umd.New(r1.Metadata)
	extractor, _ = u.FindExtractor("https://www.redgifs.com/watch/ecstaticthickasiansmallclawedotter")
	r2 := extractor.QueryMedia(99999, nil, true)
	<-r2.Done

	// Check the log output
	output := buf.String()
	assert.Equal(t, 1, strings.Count(output, "Issuing new RedGifs token"))
	assert.Equal(t, 1, strings.Count(output, "Reusing RedGifs token"))
}
