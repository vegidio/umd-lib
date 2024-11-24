package utils

import (
	"github.com/samber/lo"
	"github.com/vegidio/umd-lib/internal/model"
)

func MergeMedia(media []model.Media, newMedia []model.Media) ([]model.Media, int) {
	amountBefore := len(media)

	unique := lo.UniqBy(append(media, newMedia...), func(m model.Media) string {
		return m.Url
	})

	amountQueried := len(unique) - amountBefore

	return unique, amountQueried
}
