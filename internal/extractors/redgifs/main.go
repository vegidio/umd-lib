package redgifs

import (
	"fmt"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
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

func New(url string, metadata model.Metadata, callback func(event event.Event)) model.Extractor {
	switch {
	case utils.HasHost(url, "redgifs.com"):
		return &Redgifs{Metadata: metadata, Callback: callback}
	}

	return nil
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
	return fetch.New(nil, 0)
}

func (r *Redgifs) SetExternal(external model.External) {
	r.external = external
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Redgifs)(nil)

// region - Private methods

func (r *Redgifs) getSourceType(url string) (SourceType, error) {
	regexVideo := regexp.MustCompile(`/(ifr|watch)/([^/\n?]+)`)

	var source SourceType
	var id string

	switch {
	case regexVideo.MatchString(url):
		matches := regexVideo.FindStringSubmatch(url)
		id = matches[2]
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

func (r *Redgifs) getNewOrSavedToken() (string, error) {
	token, exists := r.Metadata[model.RedGifs]["token"].(string)

	if !exists {
		log.Debug("Issuing new RedGifs token")

		auth, err := getToken()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Failed to issue RedGifs token")

			return "", err
		}

		token = auth.Token

		if r.responseMetadata[model.RedGifs] == nil {
			r.responseMetadata[model.RedGifs] = make(map[string]interface{})
		}

		// Save the token to be reused in the future
		r.responseMetadata[model.RedGifs]["token"] = token
	} else {
		log.WithFields(log.Fields{
			"token": token,
		}).Debug("Reusing RedGifs token")
	}

	return token, nil
}

func (r *Redgifs) fetchMedia(source SourceType, limit int, extensions []string, deep bool) ([]model.Media, error) {
	media := make([]model.Media, 0)
	videos := make([]Video, 0)
	amountQueried := 0
	var err error

	token, err := r.getNewOrSavedToken()
	if err != nil {
		return nil, err
	}

	switch s := source.(type) {
	case SourceVideo:
		videos, err = r.fetchVideo(s, token)
	}

	if err != nil {
		return media, err
	}

	sourceName := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
	newMedia := videosToMedia(videos, sourceName)
	media, amountQueried = utils.MergeMedia(media, newMedia)

	if r.Callback != nil {
		r.Callback(event.OnMediaQueried{Amount: amountQueried})
	}
	
	return media, nil
}

func (r *Redgifs) fetchVideo(source SourceVideo, token string) ([]Video, error) {
	video, err := getVideo(
		fmt.Sprintf("Bearer %s", token),
		fmt.Sprintf("https://www.redgifs.com/watch/%s", source.Id),
		source.Id,
	)

	if err != nil {
		return make([]Video, 0), err
	}

	return []Video{*video}, nil
}

// endregion

// region - Private functions

func videosToMedia(videos []Video, sourceName string) []model.Media {
	return lo.Map(videos, func(video Video, _ int) model.Media {
		return model.NewMedia(video.Gif.Url.Hd, model.RedGifs, map[string]interface{}{
			"name":    video.Gif.Username,
			"source":  strings.ToLower(sourceName),
			"created": video.Gif.Created,
			"id":      video.Gif.Id,
		})
	})
}

// endregion
