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

	u, _ := umd.New("https://www.reddit.com/r/bigtiddygothgf/comments/1gz2vxn/how_it_started_vs_how_it_ended/", nil, func(ev event.Event) {
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

	log.Info(resp.Media)
}
