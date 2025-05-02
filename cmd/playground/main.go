package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib"
	"time"
)

func main() {
	//log.SetOutput(io.Discard)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	extractor, _ := umd.New(nil).
		FindExtractor("https://coomer.su/onlyfans/user/corinnakopf")

	resp := extractor.QueryMedia(100, nil, true)
	if err := queryUpdates(resp); err != nil {
		log.Error(err)
	}

	log.Info("Done")
}

func queryUpdates(resp *umd.Response) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	oldValue := -1

	for {
		select {
		case <-ticker.C:
			size := len(resp.Media)
			if size != oldValue {
				oldValue = size
				log.Info("Size: ", size)
			}

		case <-resp.Done:
			log.Info("Size: ", len(resp.Media))
			return resp.Error()
		}
	}
}
