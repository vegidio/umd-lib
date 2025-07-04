package main

import (
	"fmt"
	"github.com/samber/lo"
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

	fileCookies()
}

func query() {
	extractor, _ := umd.New(nil).
		FindExtractor("https://www.reddit.com/user/LustChloe/")

	resp, _ := extractor.QueryMedia(500, []string{"jpg"}, true)

	err := resp.Track(func(queried, total int) {
		log.Info("Queried: ", queried, " - Size: ", total)
	})

	if err != nil {
		log.Error(err)
	}

	log.Info("Done")
}

func queryDownload() {
	extractor, _ := umd.New(nil).
		FindExtractor("https://coomer.su/onlyfans/user/belledelphine")

	resp, _ := extractor.QueryMedia(50, nil, true)

	err := resp.Track(func(queried, total int) {
		log.Info("Queried: ", queried, " - Size: ", total)
	})

	if err != nil {
		log.Error(err)
	}

	co, _ := fetch.GetBrowserCookies("https://coomer.su/onlyfans/user/belledelphine", "header.user-header")
	header := fetch.CookiesToHeader(co)
	headers := map[string]string{
		"Cookie": header,
	}

	f := fetch.New(headers, 10)

	requests := lo.Map(resp.Media, func(media umd.Media, index int) *fetch.Request {
		req, _ := f.NewRequest(media.Url, fmt.Sprintf("media%d.blah", index))
		return req
	})

	result, _ := f.DownloadFiles(requests, 5)

	for file := range result {
		log.Info("Downloading ", file.Request.Url)
	}

	log.Info("Done")
}

func download() {
	f := fetch.New(nil, 10)
	request, _ := f.NewRequest("https://www.reddit.com/r/u_SecretSlutAdventures/comments/yeuy8j/14k_followers_time_forama/", "test.txt")
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

func browserCookies() {
	co, err := fetch.GetBrowserCookies("https://coomer.su/onlyfans/user/belledelphine", "header.user-header")

	if err != nil {
		log.Error(err)
		return
	}

	header := fetch.CookiesToHeader(co)
	fmt.Println(header)
}

func fileCookies() {
	co, err := fetch.GetFileCookies("/Users/vegidio/Desktop/cookies.txt")

	if err != nil {
		log.Error(err)
		return
	}

	header := fetch.CookiesToHeader(co)
	fmt.Println(header)
}
