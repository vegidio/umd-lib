package extractors

import (
	"github.com/vegidio/kmd-lib/internal/models"
	"github.com/vegidio/kmd-lib/pkg"
)

// Extractor defines the interface for extractors.
type Extractor interface {
	// IsMatch checks if the given URL matches the criteria for this extractor.
	IsMatch(url string) bool

	// QueryMedia queries media from the given URL with specified limit and extensions.
	QueryMedia(url string, limit int, extensions []string) models.Response

	// GetFetch returns the Fetch instance used by this extractor.
	GetFetch() pkg.Fetch
}
