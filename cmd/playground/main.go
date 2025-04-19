package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/event"
)

func main() {
	//log.SetOutput(io.Discard)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	umdObj := umd.New(nil, func(ev event.Event) {
		switch e := ev.(type) {
		case event.OnMediaQueried:
			log.Info("Found ", e.Amount, " media")
		}
	})

	extractor, err := umdObj.FindExtractor("https://coomer.su/onlyfans/user/melindalondon")
	if err != nil {
		log.Error(err)
		return
	}

	resp, err := extractor.QueryMedia(100, nil, true)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Media found: ", resp.Media)
}
