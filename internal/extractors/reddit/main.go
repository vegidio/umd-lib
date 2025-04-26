package reddit

import (
	"encoding/json"
	"fmt"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"reflect"
	"regexp"
	"strings"
)

const Host = "reddit.com"

type Reddit struct {
	Metadata model.Metadata
	Callback func(event event.Event)

	url              string
	source           model.SourceType
	responseMetadata model.Metadata
	external         model.External
}

func New(url string, metadata model.Metadata, callback func(event event.Event), external model.External) model.Extractor {
	switch {
	case utils.HasHost(url, Host):
		return &Reddit{Metadata: metadata, Callback: callback, url: url, external: external}
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

	if r.Callback != nil {
		sourceType := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
		r.Callback(event.OnExtractorTypeFound{Type: sourceType, Name: name})
	}

	r.source = source
	return source, nil
}

func (r *Reddit) QueryMedia(limit int, extensions []string, deep bool) (*model.Response, error) {
	var err error

	if r.responseMetadata == nil {
		r.responseMetadata = make(model.Metadata)
	}

	if r.source == nil {
		r.source, err = r.SourceType()
		if err != nil {
			return nil, err
		}
	}

	media, err := r.fetchMedia(r.source, limit, extensions, deep)
	if err != nil {
		return nil, err
	}

	if r.Callback != nil {
		r.Callback(event.OnQueryCompleted{Total: len(media)})
	}

	return &model.Response{
		Url:       r.url,
		Media:     media,
		Extractor: model.Reddit,
		Metadata:  r.responseMetadata,
	}, nil
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Reddit)(nil)

// region - Private methods

func (r *Reddit) fetchMedia(source model.SourceType, limit int, extensions []string, deep bool) ([]model.Media, error) {
	media := make([]model.Media, 0)
	amountQueried := 0
	after := ""

	sourceName := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
	name := source.Name()

	for {
		var submission *Submission
		var err error

		switch s := source.(type) {
		case SourceSubmission:
			submission, err = getSubmission(s.Id)
		case SourceUser:
			submission, err = getUserSubmissions(s.name, after, limit)
		case SourceSubreddit:
			submission, err = getSubredditSubmissions(s.name, after, limit)
		}

		if err != nil {
			return media, err
		}

		newMedia := submissionsToMedia(submission.Data.Children, sourceName, name)

		if deep {
			newMedia = r.external.ExpandMedia(newMedia, Host, &r.responseMetadata, 5)
		}

		media, amountQueried = utils.MergeMedia(media, newMedia)

		if r.Callback != nil {
			r.Callback(event.OnMediaQueried{Amount: amountQueried})
		}

		after = submission.Data.After
		if len(newMedia) == 0 || len(media) >= limit || after == "" {
			break
		}
	}

	// Limiting the number of results
	if len(media) > limit {
		media = media[:limit]
	}

	return media, nil
}

// endregion

// region - Private functions

func submissionsToMedia(submissions []Child, sourceName string, name string) []model.Media {
	media := make([]model.Media, 0)

	for _, submission := range submissions {
		if submission.Data.IsGallery {
			newMedia := getGalleryMedia(submission, sourceName, name)
			media = append(media, newMedia...)
		} else {
			url := submission.Data.SecureMedia.RedditVideo.FallbackUrl
			if url == "" {
				url = submission.Data.Url
			}

			newMedia := model.NewMedia(url, model.Reddit, map[string]interface{}{
				"source":  strings.ToLower(sourceName),
				"name":    name,
				"created": submission.Data.Created.Time,
			})

			media = append(media, newMedia)
		}
	}

	return media
}

func getGalleryMedia(submission Child, sourceName string, name string) []model.Media {
	media := make([]model.Media, 0)

	for _, value := range submission.Data.MediaMetadata {
		var metadata MediaMetadata
		jsonData, _ := json.Marshal(value)
		json.Unmarshal(jsonData, &metadata)

		if metadata.Status == "valid" {
			url := metadata.S.Image
			if url == "" {
				url = metadata.S.Gif
			}

			newMedia := model.NewMedia(url, model.Reddit, map[string]interface{}{
				"source":  strings.ToLower(sourceName),
				"name":    name,
				"created": submission.Data.Created.Time,
			})

			newMedia.Url = url
			media = append(media, newMedia)
		}
	}

	return media
}

// endregion
