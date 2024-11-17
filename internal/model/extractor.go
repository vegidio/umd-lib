package model

import (
	"github.com/vegidio/umd-lib/fetch"
)

type External interface {
	ExpandMedia(media []Media, metadata *Metadata, parallel int) []Media
}

// Extractor defines the interface for extractors.
type Extractor interface {
	// QueryMedia queries media from the given URL with specified limit and extensions.
	QueryMedia(url string, limit int, extensions []string, deep bool) (*Response, error)

	// GetFetch returns the Fetch instance used by this extractor.
	GetFetch() fetch.Fetch

	// SetExternal sets the external utilities.
	SetExternal(external External)
}
