package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMedia_UrlWithValidExtension(t *testing.T) {
	media := Media{Url: "http://example.com/file.mp4"}
	assert.Equal(t, "mp4", media.Extension)
}

func TestMedia_UrlWithValidExtensionAndQueryString(t *testing.T) {
	media := Media{Url: "http://example.com/file.mp4?query=1"}
	assert.Equal(t, "mp4", media.Extension)
}

func TestMedia_UrlWithNoExtension(t *testing.T) {
	media := Media{Url: "http://example.com/file"}
	assert.Equal(t, "", media.Extension)
}

func TestMedia_EmptyUrl(t *testing.T) {
	media := Media{Url: ""}
	assert.Equal(t, "", media.Extension)
}

func TestMedia_InvalidUrl(t *testing.T) {
	media := Media{Url: "://invalid-url"}
	assert.Equal(t, "", media.Extension)
}

func TestMedia_TypeImage(t *testing.T) {
	media := Media{Url: "http://example.com/image.jpg"}
	assert.Equal(t, Image, media.Type)

	media = Media{Url: "http://example.com/image.avif"}
	assert.Equal(t, Image, media.Type)

	media = Media{Url: "http://example.com/image.png"}
	assert.Equal(t, Image, media.Type)
}

func TestMedia_TypeVideo(t *testing.T) {
	media := Media{Url: "http://example.com/video.mp4"}
	assert.Equal(t, Video, media.Type)

	media = Media{Url: "http://example.com/video.mkv"}
	assert.Equal(t, Video, media.Type)

	media = Media{Url: "http://example.com/video.webm"}
	assert.Equal(t, Video, media.Type)
}

func TestMedia_TypeUnknown(t *testing.T) {
	media := Media{Url: "http://example.com/file.unknown"}
	assert.Equal(t, Unknown, media.Type)

	media = Media{Url: "http://example.com/file.unknown"}
	assert.Equal(t, Unknown, media.Type)
}
