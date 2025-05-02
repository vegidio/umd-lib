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
	regexPost := regexp.MustCompile(`com/([a-zA-Z0-9-_.]+)/(\d+)`)
	regexModel := regexp.MustCompile(`com/([a-zA-Z0-9-_.]+)/?`)

	var source model.SourceType

	switch {
	case regexPost.MatchString(f.url):
		matches := regexPost.FindStringSubmatch(f.url)
		name := matches[1]
		id := matches[2]
		source = SourcePost{Id: id, name: name}
	case regexModel.MatchString(f.url):
		matches := regexModel.FindStringSubmatch(f.url)
		name := matches[1]
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
	out := make(chan model.Result[[]model.Media])

	go func() {
		defer close(out)
		var posts <-chan model.Result[Post]

		switch s := source.(type) {
		case SourcePost:
			posts = f.fetchPost(s)
		case SourceModel:
			posts = f.fetchModel(s, limit)
		}

		for post := range posts {
			if post.Err != nil {
				out <- model.Result[[]model.Media]{Err: post.Err}
				return
			}

			media := postsToMedia(post.Data, source.Type())
			out <- model.Result[[]model.Media]{Data: media}
		}
	}()

	return out
}

func (f *Fapello) fetchPost(source SourcePost) <-chan model.Result[Post] {
	result := make(chan model.Result[Post])

	go func() {
		defer close(result)

		link := fmt.Sprintf("https://fapello.com/%s/%s", source.name, source.Id)

		post, err := getPost(link, source.name)
		if err != nil {
			result <- model.Result[Post]{Err: err}
		}

		result <- model.Result[Post]{Data: *post}
	}()

	return result
}

func (f *Fapello) fetchModel(source SourceModel, limit int) <-chan model.Result[Post] {
	result := make(chan model.Result[Post])

	go func() {
		defer close(result)

		links, err := getLinks(source.name, limit)
		if err != nil {
			result <- model.Result[Post]{Err: err}
			return
		}

		for _, link := range links {
			post, postErr := getPost(link, source.name)
			if postErr != nil {
				result <- model.Result[Post]{Err: postErr}
			}

			result <- model.Result[Post]{Data: *post}
		}
	}()

	return result
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
