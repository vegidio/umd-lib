package umd

import (
	"fmt"
	"github.com/vegidio/umd-lib/event"
	"github.com/vegidio/umd-lib/fetch"
	"github.com/vegidio/umd-lib/internal/extractors/reddit"
	"github.com/vegidio/umd-lib/internal/extractors/redgifs"
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"reflect"
	"sync"
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
	extractor, err := findExtractor(url, metadata, callback)

	if err != nil {
		return nil, err
	}

	if metadata == nil {
		metadata = make(model.Metadata)
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
func (u Umd) QueryMedia(limit int, extensions []string, deep bool) (*Response, error) {
	return u.extractor.QueryMedia(u.url, limit, extensions, deep)
}

// GetFetch returns the fetch.Fetch instance associated with the extractor.
//
// # Returns:
//   - fetch.Fetch - The fetch instance used by the extractor.
func (u Umd) GetFetch() fetch.Fetch {
	return u.extractor.GetFetch()
}

// region - Private functions

func findExtractor(url string, metadata model.Metadata, callback func(event event.Event)) (model.Extractor, error) {
	var extractor model.Extractor

	switch {
	case reddit.IsMatch(url):
		extractor = &reddit.Reddit{Metadata: metadata, Callback: callback}
	case redgifs.IsMatch(url):
		extractor = &redgifs.Redgifs{Metadata: metadata, Callback: callback}
	}

	if extractor == nil {
		return nil, fmt.Errorf("no extractor found for URL: %s", url)
	}

	extractor.SetExternal(External{})

	if callback != nil {
		name := utils.LastRightOf(reflect.TypeOf(extractor).String(), ".")
		callback(event.OnExtractorFound{Name: name})
	}

	return extractor, nil
}

type External struct{}

func (External) ExpandMedia(media []model.Media, metadata *model.Metadata, parallel int) []model.Media {
	result := make([]model.Media, 0)

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, parallel)

	for _, m := range media {
		wg.Add(1)

		go func(current Media) {
			defer func() {
				wg.Done()
				<-sem
			}()

			sem <- struct{}{}

			if current.Type == model.Unknown {
				uObj, err := New(current.Url, *metadata, nil)
				if err != nil {
					result = append(result, current)
					return
				}

				resp, err := uObj.QueryMedia(1, make([]string, 0), false)
				if err != nil {
					result = append(result, current)
					return
				}

				_, exists := (*metadata)[resp.Extractor]
				if !exists {
					mu.Lock()
					(*metadata)[resp.Extractor] = resp.Metadata[resp.Extractor]
					mu.Unlock()
				}

				if len(resp.Media) > 0 {
					resp.Media[0] = utils.MergeMetadata(m, resp.Media[0])
					result = append(result, resp.Media[0])
				}
			} else {
				result = append(result, current)
			}
		}(m)
	}

	wg.Wait()
	close(sem)

	return result
}

// endregion
