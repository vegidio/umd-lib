package model

import (
	"fmt"
	"net/url"
	"path"
)

// Media represents a media object.
type Media struct {
	// Url is the URL of the media.
	Url string

	// Extension is extension of the media file, derived from the URL.
	Extension string

	// Type is the type of media, determined based on the file extension.
	Type MediaType

	// Extractor is the extractor used to fetch the media.
	Extractor ExtractorType

	// Metadata contains metadata about the media. Default is an empty map.
	Metadata map[string]interface{}
}

func (m Media) String() string {
	return fmt.Sprintf("{Url: %s, Extension: %s, Type: %s, Extractor: %s, Metadata: %v}",
		m.Url, m.Extension, m.Type, m.Extractor, m.Metadata)
}

func NewMedia(url string, extractor ExtractorType, metadata map[string]interface{}) Media {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	extension := getExtension(url)
	return Media{Url: url, Extension: extension, Type: getType(extension), Extractor: extractor, Metadata: metadata}
}

// region - Private functions

func getExtension(urStr string) string {
	u, err := url.Parse(urStr)
	if err != nil {
		panic(err)
	}

	ext := path.Ext(u.Path)
	if ext == "" {
		return ""
	}

	return ext[1:]
}

func getType(extension string) MediaType {
	switch extension {
	case "jpg", "jpeg", "png", "gif", "avif":
		return Image
	case "gifv", "mp4", "m4v", "webm", "mkv":
		return Video
	default:
		return Unknown
	}
}

// endregion
