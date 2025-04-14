package model

type External interface {
	ExpandMedia(media []Media, ignoreHost string, metadata *Metadata, parallel int) []Media
}

// Extractor defines the interface for extractors.
type Extractor interface {
	// GetSourceType determines the source type of the extractor.
	GetSourceType() (SourceType, error)

	// QueryMedia queries media from the given URL with specified limit and extensions.
	QueryMedia(limit int, extensions []string, deep bool) (*Response, error)
}
