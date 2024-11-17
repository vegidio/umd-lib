package utils

import (
	"github.com/thoas/go-funk"
	"github.com/vegidio/umd-lib/internal/model"
)

func MergeMedia(media []model.Media, newMedia []model.Media) ([]model.Media, int) {
	amountBefore := len(media)

	unique := funk.UniqBy(append(media, newMedia...), func(m model.Media) string {
		return m.Url
	}).([]model.Media)

	amountQueried := len(unique) - amountBefore

	return unique, amountQueried
}
