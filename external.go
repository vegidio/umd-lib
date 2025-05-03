package umd

import (
	"github.com/vegidio/umd-lib/internal/model"
	"github.com/vegidio/umd-lib/internal/utils"
	"sync"
)

type External struct{}

func (External) ExpandMedia(media []model.Media, ignoreHost string, metadata *model.Metadata, parallel int) []model.Media {
	result := make([]model.Media, 0)

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, parallel)

	for _, m := range media {
		wg.Add(1)

		go func(current Media) {
			defer func() {
				<-sem
				wg.Done()
			}()

			sem <- struct{}{}

			if current.Type == model.Unknown && !utils.HasHost(current.Url, ignoreHost) {
				extractor, err := New(*metadata).FindExtractor(current.Url)
				if err != nil {
					mu.Lock()
					result = append(result, current)
					mu.Unlock()

					return
				}

				resp, _ := extractor.QueryMedia(1, make([]string, 0), false)
				if resp.Error() != nil {
					mu.Lock()
					result = append(result, current)
					mu.Unlock()
					return
				}

				_, exists := (*metadata)[resp.Extractor]
				if !exists {
					mu.Lock()
					(*metadata)[resp.Extractor] = resp.Metadata[resp.Extractor]
					mu.Unlock()
				}

				if len(resp.Media) > 0 {
					mu.Lock()
					resp.Media[0] = utils.MergeMetadata(m, resp.Media[0])
					result = append(result, resp.Media[0])
					mu.Unlock()
				}
			} else {
				mu.Lock()
				result = append(result, current)
				mu.Unlock()
			}
		}(m)
	}

	wg.Wait()
	close(sem)

	return result
}

// endregion
