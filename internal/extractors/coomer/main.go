package coomer

import (
	"fmt"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"reflect"
	"regexp"
	"strings"
)

type Coomer struct {
	Metadata model.Metadata
	Callback func(event event.Event)

	url              string
	extractor        model.ExtractorType
	source           model.SourceType
	services         string
	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, callback func(event event.Event), external model.External) model.Extractor {
	switch {
	case utils.HasHost(url, "coomer.su") || utils.HasHost(url, "coomer.party"):
		baseUrl = "https://coomer.su"

		return &Coomer{
			Metadata: metadata,
			Callback: callback,

			url:       url,
			extractor: model.Coomer,
			services:  "onlyfans|fansly|candfans",
			external:  external,
		}
	case utils.HasHost(url, "kemono.su") || utils.HasHost(url, "kemono.party"):
		baseUrl = "https://kemono.su"

		return &Coomer{
			Metadata: metadata,
			Callback: callback,

			url:       url,
			extractor: model.Kemono,
			services:  "patreon|fanbox|discord|fantia|afdian|boosty|gumroad|subscribestar|dlsite",
			external:  external,
		}
	}

	return nil
}

func (c *Coomer) Type() model.ExtractorType {
	return c.extractor
}

func (c *Coomer) SourceType() (model.SourceType, error) {
	regexPost := regexp.MustCompile(`(` + c.services + `)/user/([^/]+)/post/([^/\n?]+)`)
	regexUser := regexp.MustCompile(`(` + c.services + `)/user/([^/\n?]+)`)

	var source model.SourceType
	var user string

	switch {
	case regexPost.MatchString(c.url):
		matches := regexPost.FindStringSubmatch(c.url)
		service := matches[1]
		user = matches[2]
		id := matches[3]
		source = SourcePost{Service: service, Id: id, name: user}

	case regexUser.MatchString(c.url):
		matches := regexUser.FindStringSubmatch(c.url)
		service := matches[1]
		user = matches[2]
		source = SourceUser{Service: service, name: user}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", c.url)
	}

	if c.Callback != nil {
		sourceType := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
		c.Callback(event.OnExtractorTypeFound{Type: sourceType, Name: user})
	}

	c.source = source
	return source, nil
}

func (c *Coomer) QueryMedia(limit int, extensions []string, deep bool) (*model.Response, error) {
	var err error

	if c.responseMetadata == nil {
		c.responseMetadata = make(model.Metadata)
	}

	if c.source == nil {
		c.source, err = c.SourceType()
		if err != nil {
			return nil, err
		}
	}

	media, err := c.fetchMedia(c.source, limit, extensions, deep)
	if err != nil {
		return nil, err
	}

	if c.Callback != nil {
		c.Callback(event.OnQueryCompleted{Total: len(media)})
	}

	return &model.Response{
		Url:       c.url,
		Media:     media,
		Extractor: c.extractor,
		Metadata:  c.responseMetadata,
	}, nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Coomer)(nil)

// region - Private methods

func (c *Coomer) fetchMedia(source model.SourceType, limit int, extensions []string, _ bool) ([]model.Media, error) {
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
	results := getUser(source.Service, source.name)

	for result := range results {
		if result.Err != nil {
			return media, result.Err
		}

		newMedia := c.postToMedia(result.Data)

		media, amountQueried = utils.MergeMedia(media, newMedia)
		if c.Callback != nil && amountQueried > 0 {
			c.Callback(event.OnMediaQueried{Amount: amountQueried})
		}

		// Limiting the number of results
		if len(media) >= limit {
			break
		}
	}

	return media, nil
}

func (c *Coomer) fetchPostMedia(source SourcePost, limit int, extensions []string) ([]model.Media, error) {
	media := make([]model.Media, 0)
	amountQueried := 0

	response, err := getPost(source.Service, source.name, source.Id)
	if err != nil {
		return media, err
	}

	newMedia := c.postToMedia(*response)

	media, amountQueried = utils.MergeMedia(media, newMedia)
	if c.Callback != nil && amountQueried > 0 {
		c.Callback(event.OnMediaQueried{Amount: amountQueried})
	}

	return media, nil
}

func (c *Coomer) postToMedia(response Response) []model.Media {
	media := make([]model.Media, 0)

	for _, image := range response.Images {
		if image.Path != "" {
			url := image.Server + "/data" + image.Path
			newMedia := model.NewMedia(url, c.extractor, map[string]interface{}{
				"source":  response.Post.Service,
				"name":    response.Post.User,
				"created": response.Post.Published.Time,
			})

			media = append(media, newMedia)
		}
	}

	for _, video := range response.Videos {
		if video.Path != "" {
			url := video.Server + "/data" + video.Path
			newMedia := model.NewMedia(url, c.extractor, map[string]interface{}{
				"source":  response.Post.Service,
				"name":    response.Post.User,
				"created": response.Post.Published.Time,
			})

			media = append(media, newMedia)
		}
	}

	return media
}

// endregion
