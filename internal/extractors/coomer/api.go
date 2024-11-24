package coomer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/model"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var f = fetch.New(nil, 6)

func countPages(url string) (int, error) {
	html, err := f.GetText(url)
	if err != nil {
		return 0, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return 0, err
	}

	result := doc.Find("div#paginator-top small").Text()
	matches := regexp.MustCompile(`of (\d+)`).FindStringSubmatch(result)

	num, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, err
	}

	pages := int(math.Ceil(num / 50))
	return pages, nil
}

func getPostUrls(url string) ([]string, error) {
	urls := make([]string, 0)

	html, err := f.GetText(url)
	if err != nil {
		return urls, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return urls, err
	}

	return doc.Find("article").
		Map(func(i int, s *goquery.Selection) string {
			service, _ := s.Attr("data-service")
			user, _ := s.Attr("data-user")
			id, _ := s.Attr("data-id")

			return fmt.Sprintf("https://coomer.su/%s/user/%s/post/%s", service, user, id)
		}), nil
}

func getPostMedia(url string, service string, user string) ([]model.Media, error) {
	media := make([]model.Media, 0)

	html, err := f.GetText(url)
	if err != nil {
		return media, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return media, err
	}

	postId, exists := doc.Find("meta[name='id']").Attr("content")
	if !exists {
		return media, err
	}

	result := doc.Find("div.post__published").Text()
	matches := regexp.MustCompile(`Published: (.+)`).FindStringSubmatch(result)
	dateTime := matches[1]

	parsedTime, err := time.Parse("2006-01-02 15:04:05", dateTime)
	if err != nil {
		parsedTime = time.Now()
	}

	images := doc.Find("a.fileThumb").
		Map(func(i int, s *goquery.Selection) string {
			link, _ := s.Attr("href")
			return link
		})

	videos := doc.Find("a.post__attachment-link").
		Map(func(i int, s *goquery.Selection) string {
			link, _ := s.Attr("href")
			return link
		})

	links := append(images, videos...)

	media = lo.Map(links, func(link string, _ int) model.Media {
		return model.NewMedia(link, model.Coomer, map[string]interface{}{
			"source":  service,
			"name":    user,
			"id":      postId,
			"created": parsedTime,
		})
	})

	return media, nil
}