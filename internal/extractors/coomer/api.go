package coomer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/model"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var f = fetch.New(nil, 10)

func countPages(page *rod.Page, url string) (int, error) {
	element := "div#paginator-top small"
	html, err := f.GetHtml(page, url, element)
	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Error("countPages: Failed to get the HTML: ", err)

		return 0, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Error("countPages: Failed to parse the HTML: ", err)

		return 0, err
	}

	result := doc.Find("div#paginator-top small").Text()
	matches := regexp.MustCompile(`of (\d+)`).FindStringSubmatch(result)

	num, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Error("countPages: Failed to convert string to number: ", err)

		return 0, err
	}

	pages := int(math.Ceil(num / 50))
	return pages, nil
}

func getPostUrls(page *rod.Page, url string) ([]string, error) {
	urls := make([]string, 0)

	element := "article.post-card"
	html, err := f.GetHtml(page, url, element)
	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Error("getPostUrls: Failed to get the HTML: ", err)

		return urls, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Error("getPostUrls: Failed to parse the HTML: ", err)

		return urls, err
	}

	return doc.Find(element).
		Map(func(i int, s *goquery.Selection) string {
			service, _ := s.Attr("data-service")
			user, _ := s.Attr("data-user")
			id, _ := s.Attr("data-id")

			return fmt.Sprintf("https://coomer.su/%s/user/%s/post/%s", service, user, id)
		}), nil
}

func getPostMedia(page *rod.Page, url string, service string, user string) ([]model.Media, error) {
	media := make([]model.Media, 0)

	// Get the post ID
	index := strings.LastIndex(url, "/")
	postId := url[index+1:]

	element := "div.post__published"
	html, err := f.GetHtml(page, url, element)
	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Error("getPostMedia: Failed to get the HTML: ", err)

		return media, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Error("getPostMedia: Failed to parse the HTML: ", err)

		return media, err
	}

	result := doc.Find(element).Text()
	matches := regexp.MustCompile(`Published: (.+)`).FindStringSubmatch(result)
	dateTime := matches[1]

	parsedTime, err := time.Parse("2006-01-02T15:04:05", dateTime)
	if err != nil {
		log.WithFields(log.Fields{
			"dateTime": dateTime,
			"url":      url,
		}).Error("getPostMedia: Failed to convert string to time: ", err)

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
