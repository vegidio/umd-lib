package imaglr

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"regexp"
	"slices"
	"strings"
)

type Imaglr struct {
	Metadata model.Metadata

	url              string
	source           model.SourceType
	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, external model.External) model.Extractor {
	switch {
	case utils.HasHost(url, "imaglr.com"):
		return &Imaglr{Metadata: metadata, url: url, external: external}
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

	i.source = source
	return source, nil
}

func (i *Imaglr) QueryMedia(limit int, extensions []string, deep bool) (*model.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if i.responseMetadata == nil {
		i.responseMetadata = make(model.Metadata)
	}

	response := &model.Response{
		Url:       i.url,
		Media:     make([]model.Media, 0),
		Extractor: model.Imaglr,
		Metadata:  i.responseMetadata,
		Done:      make(chan error),
	}

	go func() {
		defer close(response.Done)

		if i.source == nil {
			i.source, err = i.SourceType()
			if err != nil {
				response.Done <- err
				return
			}
		}

		mediaCh := i.fetchMedia(i.source, extensions, deep)

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
var _ model.Extractor = (*Imaglr)(nil)

// region - Private methods

func (i *Imaglr) fetchMedia(
	source model.SourceType,
	extensions []string,
	_ bool,
) <-chan model.Result[[]model.Media] {
	out := make(chan model.Result[[]model.Media])

	go func() {
		defer close(out)

		posts := make([]Post, 0)
		var err error

		switch s := source.(type) {
		case SourcePost:
			posts, err = i.fetchPost(s)
		}

		if err != nil {
			out <- model.Result[[]model.Media]{Err: err}
			return
		}

		media := postsToMedia(posts, source.Name())

		// Filter files with certain extensions
		if len(extensions) > 0 {
			media = lo.Filter(media, func(m model.Media, _ int) bool {
				return slices.Contains(extensions, m.Extension)
			})
		}

		out <- model.Result[[]model.Media]{Data: media}
	}()

	return out
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
