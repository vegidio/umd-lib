package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	u, _ := umd.New("https://coomer.su/onlyfans/user/missbella", nil, nil)
	resp, _ := u.QueryMedia(99999, nil, false)

	log.Info("Found ", len(resp.Media), " media")
}
