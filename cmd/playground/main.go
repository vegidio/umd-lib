package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/event"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	u, _ := umd.New("https://coomer.su/onlyfans/user/belledelphine", nil, func(ev event.Event) {
		switch e := ev.(type) {
		case event.OnMediaQueried:
			log.Info("Found ", e.Amount, " media")
		}
	})

	resp, err := u.QueryMedia(99999, nil, true)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(len(resp.Media))
}
