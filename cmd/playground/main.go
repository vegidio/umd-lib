package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib/fetch"
	"time"
)

func main() {
	//log.SetOutput(io.Discard)
	//log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	f := fetch.New(nil, 1)

	resp := f.DownloadFile(&fetch.Request{
		Url:      "https://httpbingo.org/json",
		FilePath: "test.json",
	})

	if err := queryUpdates(resp); err != nil {
		log.Error("Failed to download")
		return
	}

	log.Info("Download successful")
}

func queryUpdates(resp *fetch.Response) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	oldValue := int64(0)

	for {
		select {
		case <-ticker.C:
			if oldValue != resp.Downloaded {
				oldValue = resp.Downloaded
				log.Info("Downloaded: ", resp.Downloaded, "; Progress: ", resp.Progress)
			}

		case <-resp.Done:
			log.Info("Downloaded: ", resp.Downloaded, "; Total: ", resp.Size)
			return resp.Error()
		}
	}
}
