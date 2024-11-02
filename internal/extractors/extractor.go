package extractors

import . "github.com/vegidio/umd-lib/internal/models"

type Extractor interface {
	QueryMedia(url string, limit int, extensions []string) Response
}
