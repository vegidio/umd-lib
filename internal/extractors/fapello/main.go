package fapello

import (
	"fmt"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/fetch"
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

	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, callback func(event event.Event)) model.Extractor {
	switch {
	case utils.HasHost(url, "fapello.com"):
		return &Fapello{Metadata: metadata, Callback: callback}
	}

	return nil
}

func (r *Fapello) QueryMedia(url string, limit int, extensions []string, deep bool) (*model.Response, error) {
	if r.responseMetadata == nil {
		r.responseMetadata = make(model.Metadata)
	}

	source, err := r.getSourceType(url)
	if err != nil {
		return nil, err
	}

	media, err := r.fetchMedia(source, limit, extensions, deep)
	if err != nil {
		return nil, err
	}

	if r.Callback != nil {
		r.Callback(event.OnQueryCompleted{Total: len(media)})
	}

	return &model.Response{
		Url:       url,
		Media:     media,
		Extractor: model.Fapello,
		Metadata:  r.responseMetadata,
	}, nil
}

func (r *Fapello) GetFetch() fetch.Fetch {
	return fetch.New(nil, 0)
}

func (r *Fapello) SetExternal(external model.External) {
	r.external = external
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Fapello)(nil)

// region - Private methods

func (r *Fapello) getSourceType(url string) (SourceType, error) {
	regexPost := regexp.MustCompile(`/([a-zA-Z0-9-_.]+)/?$`)

	var source SourceType
	var name string

	switch {
	case regexPost.MatchString(url):
		matches := regexPost.FindStringSubmatch(url)
		name = matches[1]
		source = SourceModel{Name: name}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", url)
	}

	if r.Callback != nil {
		sourceType := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
		r.Callback(event.OnExtractorTypeFound{Type: sourceType, Name: name})
	}

	return source, nil
}

func (r *Fapello) fetchMedia(source SourceType, limit int, extensions []string, deep bool) ([]model.Media, error) {
	media := make([]model.Media, 0)
	var err error

	switch s := source.(type) {
	case SourceModel:
		media, err = r.fetchModel(s, limit)
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

func (r *Fapello) fetchModel(source SourceModel, limit int) ([]model.Media, error) {
	media := make([]model.Media, 0)
	amountQueried := 0
	var err error

	links, err := getLinks(source.Name, limit)
	if err != nil {
		return media, err
	}

	for _, link := range links {
		post, postErr := getPost(link, source.Name)
		if postErr != nil {
			return media, err
		}

		newMedia := postsToMedia(*post, "model")
		media, amountQueried = utils.MergeMedia(media, newMedia)

		if r.Callback != nil {
			r.Callback(event.OnMediaQueried{Amount: amountQueried})
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
