package redgifs

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"math"
	"regexp"
	"strings"
)

type Redgifs struct {
	Metadata model.Metadata

	url              string
	source           model.SourceType
	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, external model.External) model.Extractor {
	switch {
	case utils.HasHost(url, "redgifs.com"):
		return &Redgifs{Metadata: metadata, url: url, external: external}
	}

	return nil
}

func (r *Redgifs) Type() model.ExtractorType {
	return model.RedGifs
}

func (r *Redgifs) SourceType() (model.SourceType, error) {
	regexVideo := regexp.MustCompile(`/(ifr|watch)/([^/\n?]+)`)
	regexUser := regexp.MustCompile(`/users/([^/\n?]+)`)

	var source model.SourceType
	var name string

	switch {
	case regexVideo.MatchString(r.url):
		matches := regexVideo.FindStringSubmatch(r.url)
		name = matches[2]
		source = SourceVideo{name: name}
	case regexUser.MatchString(r.url):
		matches := regexUser.FindStringSubmatch(r.url)
		name = matches[1]
		source = SourceUser{name: name}
	}

	if source == nil {
		return nil, fmt.Errorf("source type not found for URL: %s", r.url)
	}

	r.source = source
	return source, nil
}

func (r *Redgifs) QueryMedia(limit int, extensions []string, deep bool) (*model.Response, func()) {
	var err error
	ctx, stop := context.WithCancel(context.Background())

	if r.responseMetadata == nil {
		r.responseMetadata = make(model.Metadata)
	}

	response := &model.Response{
		Url:       r.url,
		Media:     make([]model.Media, 0),
		Extractor: model.RedGifs,
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

		mediaCh := r.fetchMedia(r.source, limit, extensions, deep)

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
var _ model.Extractor = (*Redgifs)(nil)

// region - Private methods

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

func (r *Redgifs) fetchMedia(
	source model.SourceType,
	limit int,
	extensions []string,
	_ bool,
) <-chan model.Result[[]model.Media] {
	out := make(chan model.Result[[]model.Media])

	go func() {
		defer close(out)
		var gifs <-chan model.Result[[]Gif]

		token, err := r.getNewOrSavedToken()
		if err != nil {
			out <- model.Result[[]model.Media]{Err: err}
			return
		}

		switch s := source.(type) {
		case SourceVideo:
			gifs = r.fetchGif(s, token)
		case SourceUser:
			gifs = r.fetchUser(s, token, limit)
		}

		for gif := range gifs {
			if gif.Err != nil {
				out <- model.Result[[]model.Media]{Err: gif.Err}
				return
			}

			media := videosToMedia(gif.Data, source.Type())
			out <- model.Result[[]model.Media]{Data: media}
		}
	}()

	return out
}

func (r *Redgifs) fetchGif(source SourceVideo, token string) <-chan model.Result[[]Gif] {
	result := make(chan model.Result[[]Gif])

	go func() {
		defer close(result)

		response, err := getGif(
			fmt.Sprintf("Bearer %s", token),
			fmt.Sprintf("https://www.redgifs.com/watch/%s", source.name),
			source.name,
		)

		if err != nil {
			result <- model.Result[[]Gif]{Err: err}
			return
		}

		result <- model.Result[[]Gif]{Data: []Gif{response.Gif}}
	}()

	return result
}

func (r *Redgifs) fetchUser(source SourceUser, token string, limit int) <-chan model.Result[[]Gif] {
	result := make(chan model.Result[[]Gif])

	go func() {
		defer close(result)

		bearer := fmt.Sprintf("Bearer %s", token)
		url := fmt.Sprintf("https://www.redgifs.com/users/%s", source.name)
		response, err := getUser(bearer, url, source.name, 1)

		if err != nil {
			result <- model.Result[[]Gif]{Err: err}
			return
		}

		result <- model.Result[[]Gif]{Data: response.Gifs}
		maxPages := math.Ceil(float64(limit) / 100)
		numPages := int(math.Min(float64(response.Pages), maxPages))

		for i := 2; i <= numPages; i++ {
			response, err = getUser(bearer, url, source.name, i)
			if err != nil {
				result <- model.Result[[]Gif]{Err: err}
				return
			}

			result <- model.Result[[]Gif]{Data: response.Gifs}
		}
	}()

	return result
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
			"created": gif.Created.Time,
			"id":      gif.Id,
		})
	})
}

// endregion
