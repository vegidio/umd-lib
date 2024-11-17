package redgifs

import (
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"reflect"
	"regexp"
	"strings"
)

type Redgifs struct {
	Metadata model.Metadata
	Callback func(event event.Event)

	responseMetadata model.Metadata
	external         model.External
}

func IsMatch(url string) bool {
	return utils.HasHost(url, "redgifs.com")
}

func (r *Redgifs) QueryMedia(url string, limit int, extensions []string, deep bool) (*model.Response, error) {
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
		Extractor: model.RedGifs,
		Metadata:  r.responseMetadata,
	}, nil
}

func (r *Redgifs) GetFetch() fetch.Fetch {
	return fetch.New(map[string]string{
		"User-Agent": "UMD",
	}, 0)
}

func (r *Redgifs) SetExternal(external model.External) {
	r.external = external
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Redgifs)(nil)

// region - Private methods

func (r *Redgifs) getSourceType(url string) (SourceType, error) {
	regexVideo := regexp.MustCompile(`/watch/([^/\n?]+)`)

	var source SourceType
	var id string

	switch {
	case regexVideo.MatchString(url):
		matches := regexVideo.FindStringSubmatch(url)
		id = matches[1]
		source = SourceVideo{Id: id}
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

func (r *Redgifs) fetchMedia(source SourceType, limit int, extensions []string, _ bool) ([]model.Media, error) {
	media := make([]model.Media, 0)
	videos := make([]Video, 0)
	amountQueried := 0
	var err error

	sourceName := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
	id := reflect.ValueOf(source).FieldByName("Id").String()

	switch s := source.(type) {
	case SourceVideo:
		videos, err = r.fetchVideo(s)
	}

	if err != nil {
		return media, err
	}

	newMedia := videosToMedia(videos, sourceName, id)
	media, amountQueried = utils.MergeMedia(media, newMedia)

	if r.Callback != nil {
		r.Callback(event.OnMediaQueried{Amount: amountQueried})
	}

	media = append(media, newMedia...)
	return media, nil
}

func (r *Redgifs) fetchVideo(source SourceVideo) ([]Video, error) {
	video, err := getVideo(source.Id)

	if err != nil {
		return make([]Video, 0), err
	}

	return []Video{*video}, nil
}

// endregion

// region - Private functions

func videosToMedia(videos []Video, sourceName string, id string) []model.Media {
	return funk.Map(videos, func(video Video) model.Media {
		return model.NewMedia(video.Url, model.RedGifs, map[string]interface{}{
			"source": sourceName,
			"id":     id,
		})
	}).([]model.Media)
}

// endregion
