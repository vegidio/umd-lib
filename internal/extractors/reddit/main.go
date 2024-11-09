package reddit

import (
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal"
	"github.com/vegidio/umd-lib/internal/model"
	"reflect"
	"regexp"
	"strings"
)

type Reddit struct {
	Callback func(event event.Event)
}

func IsMatch(url string) bool {
	return internal.HasHost(url, "reddit.com")
}

func (r Reddit) QueryMedia(url string, limit int, extensions []string) (*model.Response, error) {
	source, err := r.getSourceType(url)
	if err != nil {
		return nil, err
	}

	submissions, err := r.fetchSubmissions(source, limit, extensions)
	if err != nil {
		return nil, err
	}

	sourceName := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
	media := submissionsToMedia(submissions, sourceName, url)

	if r.Callback != nil {
		r.Callback(event.OnQueryCompleted{Total: len(media)})
	}

	return &model.Response{
		Url:       url,
		Media:     media,
		Extractor: model.Reddit,
		Metadata:  map[string]interface{}{},
	}, nil
}

func (r Reddit) GetFetch() fetch.Fetch {
	return fetch.New(make(map[string]string), 0)
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Reddit)(nil)

// region - Private methods

func (r Reddit) getSourceType(url string) (SourceType, error) {
	regexSubmission := regexp.MustCompile(`/(?:r|u|user)/([^/?]+)/comments/([^/\n?]+)`)
	regexUser := regexp.MustCompile(`/(?:u|user)/([^/\n?]+)`)
	regexSubreddit := regexp.MustCompile(`/r/([^/\n]+)`)

	var source SourceType
	var name string

	if regexSubmission.MatchString(url) {
		matches := regexSubmission.FindStringSubmatch(url)
		name = matches[1]
		id := matches[2]
		source = SourceSubmission{Name: name, ID: id}

	} else if regexUser.MatchString(url) {
		matches := regexUser.FindStringSubmatch(url)
		name = matches[1]
		source = SourceUser{Name: name}

	} else if regexSubreddit.MatchString(url) {
		matches := regexSubreddit.FindStringSubmatch(url)
		name = matches[1]
		source = SourceSubreddit{Name: name}
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

func (r Reddit) fetchSubmissions(source SourceType, limit int, extensions []string) ([]Child, error) {
	submissions := make([]Child, 0)
	after := ""

	for {
		var submission *Submission
		var err error

		switch s := source.(type) {
		case SourceSubmission:
			submission, err = getSubmission(s.ID)
		case SourceUser:
			submission, err = getUserSubmissions(s.Name, after, limit)
		case SourceSubreddit:
			submission, err = getSubredditSubmissions(s.Name, after, limit)
		}

		if err != nil {
			return make([]Child, 0), err
		}

		filteredSubmissions := submission.Data.Children
		after = submission.Data.After
		amountBefore := len(submissions)

		// Append the arrays together, but removing duplicates
		submissions = funk.UniqBy(append(submissions, filteredSubmissions...), func(c Child) string {
			url := c.Data.SecureMedia.RedditVideo.FallbackUrl
			if url == "" {
				url = c.Data.Url
			}

			return url
		}).([]Child)

		if r.Callback != nil {
			queried := len(submissions) - amountBefore
			r.Callback(event.OnMediaQueried{Amount: queried})
		}

		if len(submission.Data.Children) == 0 || len(submissions) >= limit || after == "" {
			break
		}
	}

	// Limiting the number of results
	if len(submissions) > limit {
		submissions = submissions[:limit]
	}

	return submissions, nil
}

// endregion

// region - Private functions

func submissionsToMedia(submissions []Child, sourceName string, name string) []model.Media {
	return funk.Map(submissions, func(submission Child) model.Media {
		url := submission.Data.SecureMedia.RedditVideo.FallbackUrl
		if url == "" {
			url = submission.Data.Url
		}

		return model.NewMedia(url, map[string]interface{}{
			"source":  sourceName,
			"name":    name,
			"created": submission.Data.Created.Time,
		})
	}).([]model.Media)
}

// endregion
