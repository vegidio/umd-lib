package model

type ExtractorType int

const (
	// Coomer represents the Coomer (coomer.su) extractor type.
	Coomer ExtractorType = iota
	// Imaglr represents the Imaglr (imaglr.com) extractor type.
	Imaglr
	// Reddit represents the Reddit (reddit.com) extractor type.
	Reddit
	// RedGifs the RedGifs (redgifs.com) extractor type.
	RedGifs
	// Kemono the Kemono (kemono.su) extractor type.
	Kemono
)

func (e ExtractorType) String() string {
	switch e {
	case Coomer:
		return "Coomer"
	case Imaglr:
		return "Imaglr"
	case Kemono:
		return "Kemono"
	case Reddit:
		return "Reddit"
	case RedGifs:
		return "RedGifs"
	}

	return "Unknown"
}
