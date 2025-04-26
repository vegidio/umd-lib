package fapello

import (
	"fmt"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"regexp"
	"strings"
	"time"
)

type Fapello struct {
	Metadata model.Metadata

	url              string
	source           model.SourceType
	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, external model.External) model.Extractor {
	switch {
	case utils.HasHost(url, "fapello.com"):
		return &Fapello{Metadata: metadata, url: url, external: external}
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

	f.source = source
	return source, nil
}

func (f *Fapello) QueryMedia(limit int, extensions []string, deep bool) *model.Response {
	var err error

	if f.responseMetadata == nil {
		f.responseMetadata = make(model.Metadata)
	}

	response := &model.Response{
		Url:       f.url,
		Media:     make([]model.Media, 0),
		Extractor: model.Fapello,
		Metadata:  f.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if f.source == nil {
			f.source, err = f.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		for result := range f.fetchMedia(f.source, limit, extensions, deep) {
			if result.Err != nil {
				response.Done <- result.Err
				return
			}

			// Limiting the number of results
			if utils.MergeMedia(&response.Media, result.Data) >= limit {
				response.Media = response.Media[:limit]
				break
			}
		}

		response.Done <- nil
	}()

	return response
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Fapello)(nil)

// region - Private methods

func (f *Fapello) fetchMedia(
	source model.SourceType,
	limit int,
	extensions []string,
	_ bool,
) <-chan model.Result[[]model.Media] {
	result := make(chan model.Result[[]model.Media])

	go func() {
		defer close(result)

		media := make([]model.Media, 0)
		var err error

		switch s := source.(type) {
		case SourceModel:
			media, err = f.fetchModel(s, limit)
		}

		if err != nil {
			result <- model.Result[[]model.Media]{Err: err}
			return
		}

		result <- model.Result[[]model.Media]{Data: media}
	}()

	return result
}

func (f *Fapello) fetchModel(source SourceModel, limit int) ([]model.Media, error) {
	media := make([]model.Media, 0)
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
		utils.MergeMedia(&media, newMedia)

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
