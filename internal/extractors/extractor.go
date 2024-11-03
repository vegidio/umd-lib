package extractors

import (
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/pkg"
)

// Extractor defines the interface for extractors.
type Extractor interface {
	// QueryMedia queries media from the given URL with specified limit and extensions.
	QueryMedia(url string, limit int, extensions []string) model.Response

	// GetFetch returns the Fetch instance used by this extractor.
	GetFetch() pkg.Fetch
}
