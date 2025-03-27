package fapello

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/vegidio/umd-lib/fetch"
	"regexp"
	"strconv"
	"strings"
)

const BaseUrl = "https://fapello.com/"

var f = fetch.New(nil, 0)

func getPages(name string) (int, error) {
	url := BaseUrl + name
	html, err := f.GetText(url)
	if err != nil {
		return 0, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return 0, err
	}

	pages, _ := doc.Find("div#showmore").Attr("data-max")
	numPages, err := strconv.Atoi(pages)
	if err != nil {
		return 0, err
	}

	return numPages, nil
}

func getPosts(name string, page int) ([]Post, error) {
	posts := make([]Post, 0)

	pageUrl := fmt.Sprintf("%s/ajax/model/%s/page-%d/", BaseUrl, name, page)
	html, err := f.GetText(pageUrl)
	if err != nil {
		return posts, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return posts, err
	}

	doc.Find("img.object-cover").Each(func(i int, s *goquery.Selection) {
		parentLink, _ := s.ParentsFiltered("a").Attr("href")
		matches := regexp.MustCompile(`/(\d+)/?$`).FindStringSubmatch(parentLink)
		id, _ := strconv.Atoi(matches[1])

		link, _ := s.Attr("src")
		url := strings.Replace(link, "_300px.", ".", 1)

		posts = append(posts, Post{
			Id:   id,
			Name: name,
			Url:  url,
		})
	})

	return posts, nil
}
