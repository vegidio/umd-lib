package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib"
	"time"
)

func main() {
	//log.SetOutput(io.Discard)
	//log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	extractor, _ := umd.New(nil).
		FindExtractor("https://www.reddit.com/user/atomicbrunette18/")

	resp, stop := extractor.QueryMedia(99_999, nil, true)

	go func() {
		time.Sleep(10 * time.Second)
		stop()
	}()

	err := resp.Track(func(queried, total int) {
		log.Info("Queried: ", queried, " - Size: ", total)
	})

	if err != nil {
		log.Error(err)
	}

	log.Info("Done")
}
