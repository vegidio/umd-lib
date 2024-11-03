package reddit

import (
	"github.com/vegidio/umd-lib/pkg"
)

type Submission struct {
	Data SubmissionData `json:"data"`
}

type SubmissionData struct {
	After    string  `json:"after"`
	Children []Child `json:"children"`
}

type Child struct {
	Data ChildData `json:"data"`
}

type ChildData struct {
	Author        string                 `json:"author"`
	Url           string                 `json:"url"`
	Created       pkg.EpochTime          `json:"created"`
	IsGallery     bool                   `json:"is_gallery"`
	MediaMetadata map[string]interface{} `json:"media_metadata"`
	SecureMedia   SecureMedia            `json:"secure_media"`
}

type MediaMetadata struct {
	Status string `json:"status"`
	S      S      `json:"s"`
}

type S struct {
	Image string `json:"u"`
	Gif   string `json:"gif"`
}

type SecureMedia struct {
	RedditVideo RedditVideo `json:"reddit_video"`
}

type RedditVideo struct {
	FallbackUrl string `json:"fallback_url"`
}
