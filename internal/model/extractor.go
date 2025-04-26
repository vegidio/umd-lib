package model

type External interface {
	ExpandMedia(media []Media, ignoreHost string, metadata *Metadata, parallel int) []Media
}

// Extractor defines the interface for extractors.
type Extractor interface {
	// Type returns the name of the extractor.
	Type() ExtractorType

	// SourceType determines the source type of the extractor.
	SourceType() (SourceType, error)

	// QueryMedia queries media from the given URL with specified limit and extensions.
	QueryMedia(limit int, extensions []string, deep bool) (*Response, error)
}
