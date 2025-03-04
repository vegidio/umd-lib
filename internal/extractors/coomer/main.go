package coomer

import (
	"fmt"
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

	baseUrl          string
	extractor        model.ExtractorType
	services         string
	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, callback func(event event.Event)) model.Extractor {
	switch {
	case utils.HasHost(url, "coomer.su") || utils.HasHost(url, "coomer.party"):
		return &Coomer{
			Metadata:  metadata,
			Callback:  callback,
			baseUrl:   "https://coomer.su",
			extractor: model.Coomer,
			services:  "onlyfans|fansly|candfans",
		}
	case utils.HasHost(url, "kemono.su") || utils.HasHost(url, "kemono.party"):
		return &Coomer{
			Metadata:  metadata,
			Callback:  callback,
			baseUrl:   "https://kemono.su",
			extractor: model.Kemono,
			services:  "patreon|fanbox|discord|fantia|afdian|boosty|gumroad|subscribestar|dlsite",
		}
	}

	return nil
}

func (c *Coomer) QueryMedia(url string, limit int, extensions []string, deep bool) (*model.Response, error) {
	var err error
	setBaseUrl(c.baseUrl)

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
		Extractor: c.extractor,
		Metadata:  c.responseMetadata,
	}, nil
}

func (c *Coomer) GetFetch() fetch.Fetch {
	return fetch.New(nil, 10)
}

func (c *Coomer) SetExternal(external model.External) {
	c.external = external
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Coomer)(nil)

// region - Private methods

func (c *Coomer) getSourceType(url string) (SourceType, error) {
	regexPost := regexp.MustCompile(`(` + c.services + `)/user/([^/]+)/post/([^/\n?]+)`)
	regexUser := regexp.MustCompile(`(` + c.services + `)/user/([^/\n?]+)`)

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

func (c *Coomer) fetchMedia(source SourceType, limit int, extensions []string, _ bool) ([]model.Media, error) {
	media := make([]model.Media, 0)
	var err error

	switch s := source.(type) {
	case SourceUser:
		media, err = c.fetchUserMedia(s, limit, extensions)
	case SourcePost:
		media, err = c.fetchPostMedia(s, limit, extensions)
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

	posts, err := getUserPosts(source.Service, source.User)
	if err != nil {
		return media, err
	}

	for _, post := range posts {
		newMedia := c.postToMedia(post)

		media, amountQueried = utils.MergeMedia(media, newMedia)
		if c.Callback != nil && amountQueried > 0 {
			c.Callback(event.OnMediaQueried{Amount: amountQueried})
		}

		// Limiting the number of results
		if len(media) > limit {
			break
		}
	}

	return media, nil
}

func (c *Coomer) fetchPostMedia(source SourcePost, limit int, extensions []string) ([]model.Media, error) {
	media := make([]model.Media, 0)
	amountQueried := 0

	post, err := getPost(source.Service, source.User, source.Id)
	if err != nil {
		return media, err
	}

	newMedia := c.postToMedia(*post)

	media, amountQueried = utils.MergeMedia(media, newMedia)
	if c.Callback != nil && amountQueried > 0 {
		c.Callback(event.OnMediaQueried{Amount: amountQueried})
	}

	return media, nil
}

func (c *Coomer) postToMedia(post Post) []model.Media {
	media := make([]model.Media, 0)

	if post.File.Path != "" {
		url := c.baseUrl + post.File.Path
		newMedia := model.NewMedia(url, c.extractor, map[string]interface{}{
			"source":  post.Service,
			"name":    post.User,
			"created": post.Published.Time,
		})

		media = append(media, newMedia)
	}

	for _, attachment := range post.Attachments {
		if attachment.Path != "" {
			url := c.baseUrl + attachment.Path
			newMedia := model.NewMedia(url, c.extractor, map[string]interface{}{
				"source":  post.Service,
				"name":    post.User,
				"created": post.Published.Time,
			})

			media = append(media, newMedia)
		}
	}

	return media
}

// endregion
