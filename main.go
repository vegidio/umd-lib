package umd

import (
	"fmt"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/extractors/reddit"
	"github.com/vegidio/umd-lib/internal/extractors/redgifs"
	"github.com/vegidio/umd-lib/internal/model"
	"reflect"
)

// Umd represents a Universal Media Downloader instance.
type Umd struct {
	url       string
	metadata  model.Metadata
	extractor model.Extractor
}

// New creates a new instance of Umd with the provided URL, metadata, and callback function.
// It finds the appropriate extractor for the given URL and initializes the Umd instance with it.
// If no extractor is found, it returns an error.
//
// # Parameters:
//   - url - The URL from which media will be extracted
//   - metadata - A map containing metadata for different types of extractors.
//   - callback - An optional function to be called with events during the extraction process.
//
// # Returns:
//   - *Umd - A pointer to the newly created Umd instance.
//   - error - An error if no extractor is found for the given URL.
func New(url string, metadata model.Metadata, callback func(event event.Event)) (*Umd, error) {
	extractor := findExtractor(url, metadata, callback)

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

// QueryMedia queries the media found in the URL.
//
// # Parameters:
//   - limit - The maximum number of files to query.
//   - extensions - A slice of file extensions to be queried.
//
// # Returns:
//   - *Response - A pointer to the response containing the queried media.
//   - error - An error if the query fails.
func (u Umd) QueryMedia(limit int, extensions []string) (*Response, error) {
	return u.extractor.QueryMedia(u.url, limit, extensions)
}

// GetFetch returns the fetch.Fetch instance associated with the extractor.
//
// # Returns:
//   - fetch.Fetch - The fetch instance used by the extractor.
func (u Umd) GetFetch() fetch.Fetch {
	return u.extractor.GetFetch()
}

// region - Private functions

func findExtractor(url string, metadata model.Metadata, callback func(event event.Event)) model.Extractor {
	var extractor model.Extractor

	switch {
	case reddit.IsMatch(url):
		extractor = reddit.Reddit{Callback: callback}
	case redgifs.IsMatch(url):
		data := getMetadata(metadata, model.RedGifs)
		extractor = &redgifs.Redgifs{Metadata: data, Callback: callback}
	}

	if callback != nil && extractor != nil {
		name := reflect.TypeOf(extractor).Name()
		callback(event.OnExtractorFound{Name: name})
	}

	return extractor
}

func getMetadata(metadata model.Metadata, extractorType model.ExtractorType) map[string]interface{} {
	data, exists := metadata[extractorType]
	if !exists {
		data = make(map[string]interface{})
	}

	return data
}

// endregion
