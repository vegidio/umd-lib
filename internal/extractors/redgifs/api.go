package redgifs

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/umd-lib/fetch"
	"strings"
	"time"
)

var f = fetch.New(nil, 0)

func getVideo(videoId string) (*Video, error) {
	html, err := f.GetText(fmt.Sprintf("https://www.redgifs.com/watch/%s", videoId))
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	jsonString := doc.Find("script[type='application/ld+json']").Text()
	err = json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}

	video := result["video"].(map[string]interface{})
	file := strings.ReplaceAll(lastRightOf(video["contentUrl"].(string), "/"), "-silent", "")

	return &Video{
		Author:  video["author"].(string),
		Url:     fmt.Sprintf("https://files.redgifs.com/%s", file),
		Created: parseCustomDateTime(video["uploadDate"].(string)),
	}, nil
}

func lastRightOf(s string, substring string) string {
	lastSlashIndex := strings.LastIndex(s, substring)
	if lastSlashIndex == -1 {
		return s
	}

	return s[lastSlashIndex+1:]
}

func parseCustomDateTime(input string) time.Time {
	var year, month, day, hour int

	_, err := fmt.Sscanf(input, "%%%d-%%%d-%%%dUTC%%%d:", &year, &month, &day, &hour)
	if err != nil {
		return time.Now()
	}

	return time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC)
}
