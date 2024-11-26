package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib/fetch"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	f := fetch.New(nil, 10)
	html, _ := f.GetText("https://httpbin.org/status/429")
	fmt.Println(html)
}
