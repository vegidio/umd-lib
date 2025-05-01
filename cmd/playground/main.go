package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib/fetch"
)

func main() {
	//log.SetOutput(io.Discard)
	//log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	f := fetch.New(nil, 1)

	request, _ := f.NewRequest("https://httpbingo.org/json", "test.json")
	resp := f.DownloadFile(request)

	err := resp.Track(func(completed, total int64, progress float64) {
		log.Info("Downloaded: ", completed, "; Total: ", total, "; Progress: ", progress)
	})

	if err != nil {
		log.Error("Failed to download")
		return
	}

	log.Info("Download successful; Hash: ", resp.Hash)
}
