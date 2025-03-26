package model

type ExtractorType struct {
	value string
}

func (t ExtractorType) String() string {
	return t.value
}

var (
	Coomer  = ExtractorType{"Coomer"}
	Imaglr  = ExtractorType{"Imaglr"}
	Reddit  = ExtractorType{"Reddit"}
	RedGifs = ExtractorType{"RedGifs"}
	Kemono  = ExtractorType{"Kemono"}
)
