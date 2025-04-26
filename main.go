package umd

import (
	"fmt"
	"github.com/vegidio/umd-lib/internal/extractors/coomer"
	"github.com/vegidio/umd-lib/internal/extractors/fapello"
	"github.com/vegidio/umd-lib/internal/extractors/imaglr"
	"github.com/vegidio/umd-lib/internal/extractors/reddit"
	"github.com/vegidio/umd-lib/internal/extractors/redgifs"
	"github.com/vegidio/umd-lib/internal/model"
)

// Umd represents a Universal Media Downloader instance.
type Umd struct {
	metadata model.Metadata
}

// New creates a new instance of Umd.
//
// # Parameters:
//   - metadata: A map containing metadata information.
//
// # Returns:
//   - Umd: A new instance of Umd.
func New(metadata model.Metadata) Umd {
	if metadata == nil {
		metadata = make(model.Metadata)
	}

	return Umd{metadata: metadata}
}

// FindExtractor attempts to find a suitable extractor for the given URL.
//
// # Parameters:
//   - url: The URL for which an extractor is to be found.
//
// # Returns:
//   - model.Extractor: The extractor instance if found.
//   - error: An error if no suitable extractor is found.
func (u Umd) FindExtractor(url string) (model.Extractor, error) {
	var extractor model.Extractor
	extractors := []func(string, model.Metadata, model.External) model.Extractor{
		coomer.New, fapello.New, imaglr.New, reddit.New, redgifs.New,
	}

	for _, newExtractor := range extractors {
		if e := newExtractor(url, u.metadata, External{}); e != nil {
			extractor = e
			break
		}
	}

	if extractor == nil {
		return nil, fmt.Errorf("no extractor found for URL: %s", url)
	}

	return extractor, nil
}
