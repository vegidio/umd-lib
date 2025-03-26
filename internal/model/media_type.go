package model

type MediaType struct {
	value string
}

func (t MediaType) String() string {
	return t.value
}

var (
	Image   = MediaType{"Image"}
	Video   = MediaType{"Video"}
	Unknown = MediaType{"Unknown"}
)
