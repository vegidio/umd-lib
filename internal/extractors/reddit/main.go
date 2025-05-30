package reddit

import (
	"context"
	"fmt"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"regexp"
	"strings"
)

const Host = "reddit.com"

type Reddit struct {
	Metadata model.Metadata

	url              string
	source           model.SourceType
	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, external model.External) model.Extractor {
	switch {
	case utils.HasHost(url, Host):
		return &Reddit{Metadata: metadata, url: url, external: external}
	}

	return nil
}

func (r *Reddit) Type() model.ExtractorType {
	return model.Reddit
}

func (r *Reddit) SourceType() (model.SourceType, error) {
	regexSubmission := regexp.MustCompile(`/(?:r|u|user)/([^/?]+)/comments/([^/\n?]+)`)
	regexUser := regexp.MustCompile(`/(?:u|user)/([^/\n?]+)`)
	regexSubreddit := regexp.MustCompile(`/r/([^/\n]+)`)

	var source model.SourceType
	var name string

	switch {
	case regexSubmission.MatchString(r.url):
		matches := regexSubmission.FindStringSubmatch(r.url)
		name = matches[1]
		id := matches[2]
		source = SourceSubmission{Id: id, name: name}

	case regexUser.MatchString(r.url):
		matches := regexUser.FindStringSubmatch(r.url)
		name = matches[1]
		source = SourceUser{name: name}

	case regexSubreddit.MatchString(r.url):
		matches := regexSubreddit.FindStringSubmatch(r.url)
		name = matches[1]
		source = SourceSubreddit{name: name}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", r.url)
	}

	r.source = source
	return source, nil
}

func (r *Reddit) QueryMedia(limit int, extensions []string, deep bool) (*model.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if r.responseMetadata == nil {
		r.responseMetadata = make(model.Metadata)
	}

	response := &model.Response{
		Url:       r.url,
		Media:     make([]model.Media, 0),
		Extractor: model.Reddit,
		Metadata:  r.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if r.source == nil {
			r.source, err = r.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := r.fetchMedia(r.source, extensions, deep)

		for {
			select {
			case <-ctx.Done():
				return

			case result, ok := <-mediaCh:
				if !ok {
					return
				}

				if result.Err != nil {
					response.Done <- result.Err
					return
				}

				// Limiting the number of results
				if utils.MergeMedia(&response.Media, result.Data) >= limit {
					response.Media = response.Media[:limit]
					return
				}
			}
		}
	}()

	return response, stop
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Reddit)(nil)

// region - Private methods

func (r *Reddit) fetchMedia(
	source model.SourceType,
	extensions []string,
	deep bool,
) <-chan model.Result[[]model.Media] {
	out := make(chan model.Result[[]model.Media])

	go func() {
		defer close(out)
		var children <-chan model.Result[ChildData]

		switch s := source.(type) {
		case SourceSubmission:
			children = getSubmission(s.Id)
		case SourceUser:
			children = getUserSubmissions(s.name)
		case SourceSubreddit:
			children = getSubredditSubmissions(s.name)
		}

		for child := range children {
			if child.Err != nil {
				out <- model.Result[[]model.Media]{Err: child.Err}
				return
			}

			media := r.childToMedia(child.Data, source.Type(), source.Name())
			if deep {
				media = r.external.ExpandMedia(media, Host, &r.responseMetadata, 5)
			}

			out <- model.Result[[]model.Media]{Data: media}
		}
	}()

	return out
}

func (r *Reddit) childToMedia(child ChildData, sourceName string, name string) []model.Media {
	url := child.SecureMedia.RedditVideo.FallbackUrl
	if url == "" {
		url = child.Url
	}

	newMedia := model.NewMedia(url, model.Reddit, map[string]interface{}{
		"source":  strings.ToLower(sourceName),
		"name":    name,
		"created": child.Created.Time,
	})

	return []model.Media{newMedia}
}

// endregion
