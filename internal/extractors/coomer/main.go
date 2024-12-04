package coomer

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"reflect"
	"regexp"
	"strings"
)

type Coomer struct {
	Metadata model.Metadata
	Callback func(event event.Event)

	responseMetadata model.Metadata
	external         model.External
	browser          *rod.Browser
}

func IsMatch(url string) bool {
	return utils.HasHost(url, "coomer.su")
}

func (c *Coomer) QueryMedia(url string, limit int, extensions []string, deep bool) (*model.Response, error) {
	c.browser = rod.New().MustConnect()
	defer c.browser.MustClose()

	if c.responseMetadata == nil {
		c.responseMetadata = make(model.Metadata)
	}

	source, err := c.getSourceType(url)
	if err != nil {
		return nil, err
	}

	media, err := c.fetchMedia(source, limit, extensions, deep)
	if err != nil {
		return nil, err
	}

	if c.Callback != nil {
		c.Callback(event.OnQueryCompleted{Total: len(media)})
	}

	return &model.Response{
		Url:       url,
		Media:     media,
		Extractor: model.Coomer,
		Metadata:  c.responseMetadata,
	}, nil
}

func (c *Coomer) GetFetch() fetch.Fetch {
	return fetch.New(make(map[string]string), 10)
}

func (c *Coomer) SetExternal(external model.External) {
	c.external = external
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Coomer)(nil)

// region - Private methods

func (c *Coomer) getSourceType(url string) (SourceType, error) {
	regexPost := regexp.MustCompile(`(onlyfans|fansly|candfans)/user/([^/]+)/post/([^/\n?]+)`)
	regexUser := regexp.MustCompile(`(onlyfans|fansly|candfans)/user/([^/\n?]+)`)

	var source SourceType
	var user string

	switch {
	case regexPost.MatchString(url):
		matches := regexPost.FindStringSubmatch(url)
		service := matches[1]
		user = matches[2]
		id := matches[3]
		source = SourcePost{Service: service, User: user, Id: id}

	case regexUser.MatchString(url):
		matches := regexUser.FindStringSubmatch(url)
		service := matches[1]
		user = matches[2]
		source = SourceUser{Service: service, User: user}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", url)
	}

	if c.Callback != nil {
		sourceType := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
		c.Callback(event.OnExtractorTypeFound{Type: sourceType, Name: user})
	}

	return source, nil
}

func (c *Coomer) fetchMedia(source SourceType, limit int, extensions []string, deep bool) ([]model.Media, error) {
	media := make([]model.Media, 0)
	var err error

	switch s := source.(type) {
	case SourceUser:
		media, err = c.fetchUserMedia(s, limit, extensions)
	case SourcePost:
		url := fmt.Sprintf("https://coomer.su/%s/user/%s/post/%s", s.Service, s.User, s.Id)
		media, err = getPostMedia(c.browser, url, s.Service, s.User)
	}

	if err != nil {
		return media, err
	}

	// Limiting the number of results
	if len(media) > limit {
		media = media[:limit]
	}

	return media, nil
}

func (c *Coomer) fetchUserMedia(source SourceUser, limit int, extensions []string) ([]model.Media, error) {
	media := make([]model.Media, 0)
	amountQueried := 0

	url := fmt.Sprintf("https://coomer.su/%s/user/%s", source.Service, source.User)
	numPages, err := countPages(c.browser, url)
	if err != nil {
		return media, err
	}

outerLoop:
	for i := 0; i < numPages; i++ {
		url = fmt.Sprintf("https://coomer.su/%s/user/%s?o=%d", source.Service, source.User, i*50)
		postUrls, err1 := getPostUrls(c.browser, url)
		if err1 != nil {
			return media, err1
		}

		for _, postUrl := range postUrls {
			postMedia, err2 := getPostMedia(c.browser, postUrl, source.Service, source.User)
			if err2 != nil {
				return media, err2
			}

			media, amountQueried = utils.MergeMedia(media, postMedia)

			if c.Callback != nil && amountQueried > 0 {
				c.Callback(event.OnMediaQueried{Amount: amountQueried})
			}

			if len(media) >= limit {
				break outerLoop
			}
		}
	}

	return media, nil
}

// endregion
