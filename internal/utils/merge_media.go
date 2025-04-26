package utils

import (
	"github.com/samber/lo"
	"github.com/vegidio/umd-lib/internal/model"
)

func MergeMedia(media *[]model.Media, newMedia []model.Media) int {
	*media = lo.UniqBy(append(*media, newMedia...), func(m model.Media) string {
		return m.Url
	})

	return len(*media)
}
