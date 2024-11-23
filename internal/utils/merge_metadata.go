package utils

import "github.com/vegidio/umd-lib/internal/model"

func MergeMetadata(originalMedia model.Media, expandedMedia model.Media) model.Media {
	for k, v := range originalMedia.Metadata {
		expandedMedia.Metadata[k] = v
	}

	return expandedMedia
}
