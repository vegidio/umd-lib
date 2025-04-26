package main

import (
	"fmt"
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

	extractor, err := umd.New(nil).
		FindExtractor("https://www.reddit.com/user/atomicbrunette18/")

	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Extractor..: ", extractor.Type().String())

	source, _ := extractor.SourceType()
	log.Info("Source Type: ", source.Type())
	log.Info("Source Name: ", source.Name())

	resp := extractor.QueryMedia(99_999, nil, true)
	if err = queryUpdates(resp); err != nil {
		log.Error(err)
	}

	log.Info("Amount.....: ", len(resp.Media))
}

func queryUpdates(resp *umd.Response) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Size:", len(resp.Media))

		case <-resp.Done:
			fmt.Println("Size:", len(resp.Media))
			return resp.Error()
		}
	}
}
