package umd

import (
	"fmt"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/extractors/reddit"
	"github.com/vegidio/umd-lib/internal/model"
	"reflect"
)

type Umd struct {
	url       string
	metadata  map[string]interface{}
	extractor model.Extractor
}

func New(url string, metadata map[string]interface{}, callback func(event event.Event)) (*Umd, error) {
	extractor := findExtractor(url, callback)

	// throw an error if no extractor was found
	if extractor == nil {
		return nil, fmt.Errorf("no extractor found for URL: %s", url)
	}

	return &Umd{
		url:       url,
		metadata:  metadata,
		extractor: extractor,
	}, nil
}

func (u Umd) QueryMedia(limit int, extensions []string) (*Response, error) {
	return u.extractor.QueryMedia(u.url, limit, extensions)
}

func (u Umd) GetFetch() fetch.Fetch {
	return u.extractor.GetFetch()
}

// region - Private functions

func findExtractor(url string, callback func(event event.Event)) model.Extractor {
	var extractor model.Extractor

	if reddit.IsMatch(url) {
		extractor = reddit.Reddit{Callback: callback}
	}

	if callback != nil && extractor != nil {
		name := reflect.TypeOf(extractor).Name()
		callback(event.OnExtractorFound{Name: name})
	}

	return extractor
}

// endregion
