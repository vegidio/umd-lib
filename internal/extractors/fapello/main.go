package fapello

import (
	"fmt"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Fapello struct {
	Metadata model.Metadata
	Callback func(event event.Event)

	url              string
	source           model.SourceType
	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, callback func(event event.Event), external model.External) model.Extractor {
	switch {
	case utils.HasHost(url, "fapello.com"):
		return &Fapello{Metadata: metadata, Callback: callback, url: url, external: external}
	}

	return nil
}

func (f *Fapello) Type() model.ExtractorType {
	return model.Fapello
}

func (f *Fapello) SourceType() (model.SourceType, error) {
	regexPost := regexp.MustCompile(`/([a-zA-Z0-9-_.]+)/?$`)

	var source model.SourceType
	var name string

	switch {
	case regexPost.MatchString(f.url):
		matches := regexPost.FindStringSubmatch(f.url)
		name = matches[1]
		source = SourceModel{name: name}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", f.url)
	}

	if f.Callback != nil {
		sourceType := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
		f.Callback(event.OnExtractorTypeFound{Type: sourceType, Name: name})
	}

	f.source = source
	return source, nil
}

func (f *Fapello) QueryMedia(limit int, extensions []string, deep bool) (*model.Response, error) {
	var err error

	if f.responseMetadata == nil {
		f.responseMetadata = make(model.Metadata)
	}

	if f.source == nil {
		f.source, err = f.SourceType()
		if err != nil {
			return nil, err
		}
	}

	media, err := f.fetchMedia(f.source, limit, extensions, deep)
	if err != nil {
		return nil, err
	}

	if f.Callback != nil {
		f.Callback(event.OnQueryCompleted{Total: len(media)})
	}

	return &model.Response{
		Url:       f.url,
		Media:     media,
		Extractor: model.Fapello,
		Metadata:  f.responseMetadata,
	}, nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Fapello)(nil)

// region - Private methods

func (f *Fapello) fetchMedia(source model.SourceType, limit int, extensions []string, deep bool) ([]model.Media, error) {
	media := make([]model.Media, 0)
	var err error

	switch s := source.(type) {
	case SourceModel:
		media, err = f.fetchModel(s, limit)
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

func (f *Fapello) fetchModel(source SourceModel, limit int) ([]model.Media, error) {
	media := make([]model.Media, 0)
	amountQueried := 0
	var err error

	links, err := getLinks(source.name, limit)
	if err != nil {
		return media, err
	}

	for _, link := range links {
		post, postErr := getPost(link, source.name)
		if postErr != nil {
			return media, err
		}

		newMedia := postsToMedia(*post, "model")
		media, amountQueried = utils.MergeMedia(media, newMedia)

		if f.Callback != nil {
			f.Callback(event.OnMediaQueried{Amount: amountQueried})
		}

		if len(media) >= limit {
			break
		}
	}

	return media, nil
}

// endregion

// region - Private functions

func postsToMedia(post Post, sourceName string) []model.Media {
	now := time.Date(1980, time.October, 6, 17, 7, 0, 0, time.UTC)

	return []model.Media{model.NewMedia(post.Url, model.Fapello, map[string]interface{}{
		"id":      post.Id,
		"name":    post.Name,
		"source":  strings.ToLower(sourceName),
		"created": now.Add(time.Duration(post.Id*24) * time.Hour),
	})}
}

// endregion
