package redgifs

import (
	"fmt"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"math"
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
	regexUser := regexp.MustCompile(`/users/([^/\n?]+)`)

	var source SourceType
	var name string

	switch {
	case regexVideo.MatchString(url):
		matches := regexVideo.FindStringSubmatch(url)
		name = matches[2]
		source = SourceVideo{Id: name}
	case regexUser.MatchString(url):
		matches := regexUser.FindStringSubmatch(url)
		name = matches[1]
		source = SourceUser{Name: name}
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
	gifs := make([]Gif, 0)
	amountQueried := 0
	var err error

	token, err := r.getNewOrSavedToken()
	if err != nil {
		return nil, err
	}

	switch s := source.(type) {
	case SourceVideo:
		gifs, err = r.fetchGif(s, token)
	case SourceUser:
		gifs, err = r.fetchUser(s, token, limit)
	}

	if err != nil {
		return media, err
	}

	sourceName := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
	newMedia := videosToMedia(gifs, sourceName)
	media, amountQueried = utils.MergeMedia(media, newMedia)

	if r.Callback != nil {
		r.Callback(event.OnMediaQueried{Amount: amountQueried})
	}

	// Limiting the number of results
	if len(media) > limit {
		media = media[:limit]
	}

	return media, nil
}

func (r *Redgifs) fetchGif(source SourceVideo, token string) ([]Gif, error) {
	response, err := getGif(
		fmt.Sprintf("Bearer %s", token),
		fmt.Sprintf("https://www.redgifs.com/watch/%s", source.Id),
		source.Id,
	)

	if err != nil {
		return make([]Gif, 0), err
	}

	return []Gif{response.Gif}, nil
}

func (r *Redgifs) fetchUser(source SourceUser, token string, limit int) ([]Gif, error) {
	gifs := make([]Gif, 0)

	bearer := fmt.Sprintf("Bearer %s", token)
	url := fmt.Sprintf("https://www.redgifs.com/users/%s", source.Name)
	response, err := getUser(bearer, url, source.Name, 1)

	if err != nil {
		return gifs, err
	}

	gifs = append(gifs, response.Gifs...)
	maxPages := math.Ceil(float64(limit) / 100)
	numPages := int(math.Min(float64(response.Pages), maxPages))

	for i := 2; i <= numPages; i++ {
		response, err = getUser(bearer, url, source.Name, i)
		if err != nil {
			return gifs, err
		}

		gifs = append(gifs, response.Gifs...)
	}

	return gifs, nil
}

// endregion

// region - Private functions

func videosToMedia(gifs []Gif, sourceName string) []model.Media {
	return lo.Map(gifs, func(gif Gif, _ int) model.Media {
		url := gif.Url.Hd
		if url == "" {
			url = gif.Url.Sd
		}

		return model.NewMedia(url, model.RedGifs, map[string]interface{}{
			"name":    gif.Username,
			"source":  strings.ToLower(sourceName),
			"created": gif.Created,
			"id":      gif.Id,
		})
	})
}

// endregion
