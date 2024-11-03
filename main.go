package umd

import (
	"fmt"
	"github.com/vegidio/umd-lib/internal/extractors/reddit"
	"github.com/vegidio/umd-lib/model"
	"reflect"
)

type Umd struct {
	url       string
	metadata  map[string]interface{}
	extractor model.Extractor
}

func New(url string, metadata map[string]interface{}, callback func(event model.Event)) (*Umd, error) {
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

func (u Umd) QueryMedia(limit int, extensions []string) model.Response {
	return u.extractor.QueryMedia(u.url, limit, extensions)
}

// region - Private functions

func findExtractor(url string, callback func(event model.Event)) model.Extractor {
	var extractor model.Extractor

	if reddit.IsMatch(url) {
		extractor = reddit.Reddit{Callback: callback}
	}

	if callback != nil && extractor != nil {
		name := reflect.TypeOf(extractor).Name()
		callback(model.OnExtractorFound{Name: name})
	}

	return extractor
}

// endregion
