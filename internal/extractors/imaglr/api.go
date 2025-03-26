package imaglr

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/umd-lib/fetch"
	"strings"
	"time"
)

const BaseUrl = "https://imaglr.com/"

var f = fetch.New(nil, 0)

func getPost(id string) (*Post, error) {
	url := BaseUrl + fmt.Sprintf("post/%s", id)
	html, err := f.GetText(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	author, _ := doc.Find("meta[name='author']").Attr("content")
	mediaType, _ := doc.Find("meta[property='og:type']").Attr("content")
	image, _ := doc.Find("meta[property='og:image']").Attr("content")
	video, _ := doc.Find("meta[property='og:video']").Attr("content")

	var result map[string]interface{}
	jsonString, _ := doc.Find("div#app").Attr("data-page")
	err = json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}

	timestamp := result["props"].(map[string]interface{})["post"].(map[string]interface{})["data"].(map[string]interface{})["created_at_timestamp"].(float64)
	createdAt := time.Unix(int64(timestamp), 0)

	post := &Post{
		Id:        id,
		Author:    author,
		Type:      mediaType,
		Image:     image,
		Video:     video,
		Timestamp: createdAt,
	}

	return post, nil
}
