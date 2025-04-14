package umd

import (
	"fmt"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/internal/extractors/coomer"
	"github.com/vegidio/umd-lib/internal/extractors/fapello"
	"github.com/vegidio/umd-lib/internal/extractors/imaglr"
	"github.com/vegidio/umd-lib/internal/extractors/reddit"
	"github.com/vegidio/umd-lib/internal/extractors/redgifs"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"reflect"
)

// Umd represents a Universal Media Downloader instance.
type Umd struct {
	metadata  model.Metadata
	callback  func(event event.Event)
	extractor model.Extractor
	source    model.SourceType
}

// New creates a new instance of Umd.
//
// # Parameters:
//   - metadata: A map containing metadata information.
//   - callback: A function to handle events of type event.Event.
//
// # Returns:
//   - Umd: A new instance of Umd.
func New(metadata model.Metadata, callback func(event event.Event)) Umd {
	if metadata == nil {
		metadata = make(model.Metadata)
	}

	return Umd{
		metadata: metadata,
		callback: callback,
	}
}

// FindExtractor attempts to find a suitable extractor for the given URL.
//
// # Parameters:
//   - url: The URL for which an extractor is to be found.
//
// # Returns:
//   - model.Extractor: The extractor instance if found.
//   - error: An error if no suitable extractor is found.
func (u *Umd) FindExtractor(url string) (model.Extractor, error) {
	var extractor model.Extractor
	extractors := []func(string, model.Metadata, func(event.Event), model.External) model.Extractor{
		coomer.New, fapello.New, imaglr.New, reddit.New, redgifs.New,
	}

	for _, newExtractor := range extractors {
		if e := newExtractor(url, u.metadata, u.callback, External{}); e != nil {
			extractor = e
			break
		}
	}

	if extractor == nil {
		return nil, fmt.Errorf("no extractor found for URL: %s", url)
	}

	if u.callback != nil {
		name := utils.LastRightOf(reflect.TypeOf(extractor).String(), ".")
		u.callback(event.OnExtractorFound{Name: name})
	}

	return extractor, nil
}
