package imaglr

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"reflect"
	"regexp"
	"strings"
)

type Imaglr struct {
	Metadata model.Metadata
	Callback func(event event.Event)

	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, callback func(event event.Event)) model.Extractor {
	switch {
	case utils.HasHost(url, "imaglr.com"):
		return &Imaglr{Metadata: metadata, Callback: callback}
	}

	return nil
}

func (r *Imaglr) QueryMedia(url string, limit int, extensions []string, deep bool) (*model.Response, error) {
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
		Extractor: model.Imaglr,
		Metadata:  r.responseMetadata,
	}, nil
}

func (r *Imaglr) GetFetch() fetch.Fetch {
	return fetch.New(nil, 0)
}

func (r *Imaglr) SetExternal(external model.External) {
	r.external = external
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Imaglr)(nil)

// region - Private methods

func (r *Imaglr) getSourceType(url string) (SourceType, error) {
	regexPost := regexp.MustCompile(`/post/([^/\n?]+)`)

	var source SourceType
	var id string

	switch {
	case regexPost.MatchString(url):
		matches := regexPost.FindStringSubmatch(url)
		id = matches[1]
		source = SourcePost{Id: id}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", url)
	}

	if r.Callback != nil {
		sourceType := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
		r.Callback(event.OnExtractorTypeFound{Type: sourceType, Name: id})
	}

	return source, nil
}

func (r *Imaglr) fetchMedia(source SourceType, limit int, extensions []string, deep bool) ([]model.Media, error) {
	media := make([]model.Media, 0)
	posts := make([]Post, 0)
	amountQueried := 0
	var err error

	switch s := source.(type) {
	case SourcePost:
		posts, err = r.fetchPost(s)
	}

	if err != nil {
		return media, err
	}

	sourceName := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
	newMedia := postsToMedia(posts, sourceName)
	media, amountQueried = utils.MergeMedia(media, newMedia)

	if r.Callback != nil {
		r.Callback(event.OnMediaQueried{Amount: amountQueried})
	}

	return media, nil
}

func (r *Imaglr) fetchPost(source SourcePost) ([]Post, error) {
	post, err := getPost(source.Id)

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
