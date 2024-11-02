package reddit

import (
	"github.com/thoas/go-funk"
	"github.com/vegidio/kmd-lib/internal/extractors"
	"github.com/vegidio/kmd-lib/internal/models"
	"github.com/vegidio/kmd-lib/pkg"
	"reflect"
	"regexp"
	"strings"
)

type Reddit struct {
	Callback func(models.Event)
}

func (r Reddit) IsMatch(url string) bool {
	found := pkg.HasHost(url, "reddit.com")
	if found && r.Callback != nil {
		r.Callback(models.OnExtractorFound{Name: "Reddit"})
	}

	return found
}

func (r Reddit) QueryMedia(url string, limit int, extensions []string) models.Response {
	source := r.getSourceType(url)
	submissions := r.fetchSubmissions(source, limit, extensions)

	sourceName := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
	media := submissionsToMedia(submissions, sourceName, url)

	if r.Callback != nil {
		r.Callback(models.OnQueryCompleted{Total: len(media)})
	}

	return models.Response{}
}

func (r Reddit) GetFetch() pkg.Fetch {
	return pkg.Fetch{}
}

// Compile-time assertion to ensure the extractor implements the Extractor interface
var _ extractors.Extractor = (*Reddit)(nil)

// region - Private methods

func (r Reddit) getSourceType(url string) SourceType {
	regexSubmission := regexp.MustCompile(`/(?:r|u|user)/([^/?]+)/comments/([^/\n?]+)`)
	regexUser := regexp.MustCompile(`/(?:u|user)/([^/\n?]+)`)
	regexSubreddit := regexp.MustCompile(`/r/([^/\n]+)`)

	var source SourceType

	if regexSubmission.MatchString(url) {
		matches := regexSubmission.FindStringSubmatch(url)
		name := matches[1]
		id := matches[2]
		source = SourceSubmission{Name: name, ID: id}
	}

	if regexUser.MatchString(url) {
		matches := regexUser.FindStringSubmatch(url)
		name := matches[1]
		source = SourceUser{Name: name}
	}

	if regexSubreddit.MatchString(url) {
		matches := regexSubreddit.FindStringSubmatch(url)
		name := matches[1]
		source = SourceSubreddit{Name: name}
	}

	if r.Callback != nil {
		sourceName := strings.TrimPrefix(reflect.TypeOf(source).Name(), "Source")
		r.Callback(models.OnExtractorTypeFound{Name: sourceName})
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
		submissions = append(submissions, filteredSubmissions...)

		if r.Callback != nil {
			queried := len(submissions) - amountBefore
			r.Callback(models.OnMediaQueried{Amount: queried})
		}

		if len(submission.Data.Children) == 0 || len(submissions) >= limit || after == "" {
			break
		}
	}

	return submissions[:limit]
}

// endregion

// region - Private functions

func submissionsToMedia(submissions []Child, sourceName string, name string) []models.Media {
	return funk.Map(submissions, func(submission Child) models.Media {
		url := submission.Data.SecureMedia.RedditVideo.FallbackUrl
		if url == "" {
			url = submission.Data.Url
		}

		return models.Media{
			Url: url,
			Metadata: map[string]interface{}{
				"source":  sourceName,
				"name":    name,
				"created": submission.Data.Created,
			},
		}
	}).([]models.Media)
}

// endregion