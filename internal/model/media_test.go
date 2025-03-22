package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMedia_UrlWithValidExtension(t *testing.T) {
	media := NewMedia("http://example.com/file.mp4", Reddit, nil)
	assert.Equal(t, "mp4", media.Extension)
}

func TestMedia_UrlWithValidExtensionAndQueryString(t *testing.T) {
	media := NewMedia("http://example.com/file.mp4?query=1", Reddit, nil)
	assert.Equal(t, "mp4", media.Extension)
}

func TestMedia_UrlWithNoExtension(t *testing.T) {
	media := NewMedia("http://example.com/file", Reddit, nil)
	assert.Equal(t, "", media.Extension)
}

func TestMedia_EmptyUrl(t *testing.T) {
	media := NewMedia("", Reddit, nil)
	assert.Equal(t, "", media.Extension)
}

func TestMedia_InvalidUrl(t *testing.T) {
	assert.Panics(t, func() {
		_ = NewMedia("://invalid-url", Reddit, nil)
	})
}

func TestMedia_TypeImage(t *testing.T) {
	media := NewMedia("http://example.com/image.jpg", Reddit, nil)
	assert.Equal(t, Image, media.Type)

	media = NewMedia("http://example.com/image.avif", Reddit, nil)
	assert.Equal(t, Image, media.Type)

	media = NewMedia("http://example.com/image.png", Reddit, nil)
	assert.Equal(t, Image, media.Type)
}

func TestMedia_TypeVideo(t *testing.T) {
	media := NewMedia("http://example.com/video.mp4", Reddit, nil)
	assert.Equal(t, Video, media.Type)

	media = NewMedia("http://example.com/video.mkv", Reddit, nil)
	assert.Equal(t, Video, media.Type)

	media = NewMedia("http://example.com/video.webm", Reddit, nil)
	assert.Equal(t, Video, media.Type)
}

func TestMedia_TypeUnknown(t *testing.T) {
	media := NewMedia("http://example.com/file.unknown", Reddit, nil)
	assert.Equal(t, Unknown, media.Type)
}
