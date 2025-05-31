package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/fetch"
)

func main() {
	//log.SetOutput(io.Discard)
	//log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	downloadFiles()
}

func query() {
	extractor, _ := umd.New(nil).
		FindExtractor("https://www.redgifs.com/watch/suddenthinjohndory")

	resp, _ := extractor.QueryMedia(99_999, nil, true)

	err := resp.Track(func(queried, total int) {
		log.Info("Queried: ", queried, " - Size: ", total)
	})

	if err != nil {
		log.Error(err)
	}

	log.Info("Done")
}

func download() {
	f := fetch.New(nil, 10)
	request, _ := f.NewRequest("https://www.redgifs.com/watch/suddenthinjohndory", "test.mp4")
	resp := f.DownloadFile(request)

	err := resp.Track(func(completed, total int64, progress float64) {
		log.Info("Progress: ", progress)
	})

	if err != nil {
		log.Error(err)
	}

	log.Info("Done")
}

func downloadFiles() {
	f := fetch.New(nil, 10)
	request, _ := f.NewRequest("https://www.redgifs.com/watch/suddenthinjohndory", "test.mp4")
	requests := []*fetch.Request{request}

	result, _ := f.DownloadFiles(requests, 5)

	for file := range result {
		log.Info("Downloading ", file.Request.Url)
	}

	log.Info("Done")
}
