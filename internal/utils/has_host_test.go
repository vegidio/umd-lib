package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasHost_ValidURLAndMatchingSuffix(t *testing.T) {
	urlStr := "http://example.com"
	suffix := "example.com"
	assert.True(t, HasHost(urlStr, suffix), "Expected true for URL %s with suffix %s", urlStr, suffix)
}

func TestHasHost_ValidURLAndNonMatchingSuffix(t *testing.T) {
	urlStr := "http://example.com"
	suffix := "test.com"
	assert.False(t, HasHost(urlStr, suffix), "Expected false for URL %s with suffix %s", urlStr, suffix)
}

func TestHasHost_InvalidURL(t *testing.T) {
	urlStr := "://invalid-url"
	suffix := "example.com"
	assert.False(t, HasHost(urlStr, suffix), "Expected false for invalid URL %s", urlStr)
}

func TestHasHost_URLContainingPort(t *testing.T) {
	urlStr := "http://example.com:8080"
	suffix := "example.com"
	assert.True(t, HasHost(urlStr, suffix), "Expected true for URL %s with suffix %s", urlStr, suffix)
}

func TestHasHost_EmptyURL(t *testing.T) {
	urlStr := ""
	suffix := "example.com"
	assert.False(t, HasHost(urlStr, suffix), "Expected false for empty URL")
}

func TestHasHost_EmptySuffix(t *testing.T) {
	urlStr := "http://example.com"
	suffix := ""
	assert.True(t, HasHost(urlStr, suffix), "Expected true for URL %s with empty suffix", urlStr)
}
