package imaglr

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"reflect"
	"regexp"
	"strings"
)

type Imaglr struct {
	Metadata model.Metadata
	Callback func(event event.Event)

	url              string
	source           model.SourceType
	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, callback func(event event.Event), external model.External) model.Extractor {
	switch {
	case utils.HasHost(url, "imaglr.com"):
		return &Imaglr{Metadata: metadata, Callback: callback, url: url, external: external}
	}

	return nil
}

func (i *Imaglr) Type() model.ExtractorType {
	return model.Imaglr
}

func (i *Imaglr) SourceType() (model.SourceType, error) {
	regexPost := regexp.MustCompile(`/post/([^/\n?]+)`)

	var source model.SourceType
	var id string

	switch {
	case regexPost.MatchString(i.url):
		matches := regexPost.FindStringSubmatch(i.url)
		id = matches[1]
		source = SourcePost{name: id}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", i.url)
	}

	if i.Callback != nil {
		sourceType := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
		i.Callback(event.OnExtractorTypeFound{Type: sourceType, Name: id})
	}

	i.source = source
	return source, nil
}

func (i *Imaglr) QueryMedia(limit int, extensions []string, deep bool) (*model.Response, error) {
	var err error

	if i.responseMetadata == nil {
		i.responseMetadata = make(model.Metadata)
	}

	if i.source == nil {
		i.source, err = i.SourceType()
		if err != nil {
			return nil, err
		}
	}

	media, err := i.fetchMedia(i.source, limit, extensions, deep)
	if err != nil {
		return nil, err
	}

	if i.Callback != nil {
		i.Callback(event.OnQueryCompleted{Total: len(media)})
	}

	return &model.Response{
		Url:       i.url,
		Media:     media,
		Extractor: model.Imaglr,
		Metadata:  i.responseMetadata,
	}, nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Imaglr)(nil)

// region - Private methods

func (i *Imaglr) fetchMedia(source model.SourceType, limit int, extensions []string, deep bool) ([]model.Media, error) {
	media := make([]model.Media, 0)
	posts := make([]Post, 0)
	amountQueried := 0
	var err error

	switch s := source.(type) {
	case SourcePost:
		posts, err = i.fetchPost(s)
	}

	if err != nil {
		return media, err
	}

	sourceName := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
	newMedia := postsToMedia(posts, sourceName)
	media, amountQueried = utils.MergeMedia(media, newMedia)

	if i.Callback != nil {
		i.Callback(event.OnMediaQueried{Amount: amountQueried})
	}

	return media, nil
}

func (i *Imaglr) fetchPost(source SourcePost) ([]Post, error) {
	post, err := getPost(source.name)

	if err != nil {
		return make([]Post, 0), err
	}

	return []Post{*post}, nil
}

// endregion

// region - Private functions

func postsToMedia(posts []Post, sourceName string) []model.Media {
	return lo.Map(posts, func(post Post, _ int) model.Media {
		var url string
		if post.Type == "video" {
			url = post.Video
		} else {
			url = post.Image
		}

		return model.NewMedia(url, model.Imaglr, map[string]interface{}{
			"id":      post.Id,
			"name":    post.Author,
			"source":  strings.ToLower(sourceName),
			"created": post.Timestamp,
		})
	})
}

// endregion
