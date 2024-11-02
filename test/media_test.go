package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "github.com/vegidio/umd-lib/internal/models"
	"testing"
)

func TestUrlWithValidExtension(t *testing.T) {
	media := Media{Url: "http://example.com/file.mp4"}
	assert.Equal(t, "mp4", media.Extension())
}

func TestUrlWithValidExtensionAndQueryString(t *testing.T) {
	media := Media{Url: "http://example.com/file.mp4?query=1"}
	assert.Equal(t, "mp4", media.Extension())
}

func TestUrlWithNoExtension(t *testing.T) {
	media := Media{Url: "http://example.com/file"}
	assert.Equal(t, "", media.Extension())
}

func TestEmptyUrl(t *testing.T) {
	media := Media{Url: ""}
	assert.Equal(t, "", media.Extension())
}

func TestInvalidUrl(t *testing.T) {
	media := Media{Url: "://invalid-url"}
	require.Panics(t, func() {
		media.Extension()
	}, "expected panic for invalid URL")
}

func TestTypeImage(t *testing.T) {
	media := Media{Url: "http://example.com/image.jpg"}
	assert.Equal(t, Image, media.Type())

	media = Media{Url: "http://example.com/image.avif"}
	assert.Equal(t, Image, media.Type())

	media = Media{Url: "http://example.com/image.png"}
	assert.Equal(t, Image, media.Type())
}

func TestTypeVideo(t *testing.T) {
	media := Media{Url: "http://example.com/video.mp4"}
	assert.Equal(t, Video, media.Type())

	media = Media{Url: "http://example.com/video.mkv"}
	assert.Equal(t, Video, media.Type())

	media = Media{Url: "http://example.com/video.webm"}
	assert.Equal(t, Video, media.Type())
}

func TestTypeUnknown(t *testing.T) {
	media := Media{Url: "http://example.com/file.unknown"}
	assert.Equal(t, Unknown, media.Type())

	media = Media{Url: "http://example.com/file.unknown"}
	assert.Equal(t, Unknown, media.Type())
}
