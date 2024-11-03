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

	// Metadata contains metadata about the media. Default is an empty map.
	Metadata map[string]interface{}
}

// region - Methods

// Extension is the extension of the media file, derived from the URL.
//
// It parses the URL and extracts the file extension from the path. If the URL is invalid or the path has no extension,
// it returns an empty string.
func (m Media) Extension() string {
	u, err := url.Parse(m.Url)
	if err != nil {
		panic(err)
	}

	ext := path.Ext(u.Path)
	if ext == "" {
		return ""
	}

	return ext[1:]
}

// Type is the type of media, determined based on the file extension.
//
// It uses the Extension method to get the file extension and maps it to a MediaType. If the extension is not
// recognized, it returns Unknown.
func (m Media) Type() MediaType {
	ext := m.Extension()

	switch ext {
	case "jpg", "jpeg", "png", "gif", "avif":
		return Image
	case "gifv", "mp4", "m4v", "webm", "mkv":
		return Video
	default:
		return Unknown
	}
}

func (m Media) String() string {
	return fmt.Sprintf("{Url: %s, Metadata: %v, Extension: %s, Type: %s}",
		m.Url, m.Metadata, m.Extension(), m.Type())
}

// endregion
