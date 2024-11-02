package extractors

import (
	. "github.com/vegidio/kmd-lib/internal/models"
	. "github.com/vegidio/kmd-lib/pkg"
)

type Extractor interface {
	QueryMedia(url string, limit int, extensions []string) Response
	GetFetch() Fetch
}
