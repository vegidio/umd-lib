package reddit

import (
	"github.com/thoas/go-funk"
	"github.com/vegidio/umd-lib/model"
	"github.com/vegidio/umd-lib/pkg"
	"reflect"
	"regexp"
	"strings"
)

type Reddit struct {
	Callback func(event model.Event)
}

func IsMatch(url string) bool {
	return pkg.HasHost(url, "reddit.com")
}

func (r Reddit) QueryMedia(url string, limit int, extensions []string) model.Response {
	source := r.getSourceType(url)
	submissions := r.fetchSubmissions(source, limit, extensions)

	sourceName := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
	media := submissionsToMedia(submissions, sourceName, url)

	if r.Callback != nil {
		r.Callback(model.OnQueryCompleted{Total: len(media)})
	}

	return model.Response{
		Url:       url,
		Media:     media,
		Extractor: model.Reddit,
		Metadata:  map[string]interface{}{},
	}
}

func (r Reddit) GetFetch() pkg.Fetch {
	return pkg.Fetch{}
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ model.Extractor = (*Reddit)(nil)

// region - Private methods

func (r Reddit) getSourceType(url string) SourceType {
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
	}

	if regexUser.MatchString(url) {
		matches := regexUser.FindStringSubmatch(url)
		name = matches[1]
		source = SourceUser{Name: name}
	}

	if regexSubreddit.MatchString(url) {
		matches := regexSubreddit.FindStringSubmatch(url)
		name = matches[1]
		source = SourceSubreddit{Name: name}
	}

	if r.Callback != nil {
		sourceType := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
		r.Callback(model.OnExtractorTypeFound{Type: sourceType, Name: name})
	}

	return source
}

func (r Reddit) fetchSubmissions(source SourceType, limit int, extensions []string) []Child {
	submissions := make([]Child, 0)
	after := ""

	for {
		var submission Submission

		switch s := source.(type) {
		case SourceSubmission:
			submission = getSubmission(s.ID)[0]
		case SourceUser:
			submission = getUserSubmissions(s.Name, after, limit)
		case SourceSubreddit:
			submission = getSubredditSubmissions(s.Name, after, limit)
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
			r.Callback(model.OnMediaQueried{Amount: queried})
		}

		if len(submission.Data.Children) == 0 || len(submissions) >= limit || after == "" {
			break
		}
	}

	// Limiting the number of results
	if len(submissions) > limit {
		submissions = submissions[:limit]
	}

	return submissions
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
			"created": submission.Data.Created,
		})
	}).([]model.Media)
}

// endregion
