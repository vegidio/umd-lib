package imaglr

// SourcePost represents a post source type.
type SourcePost struct {
	name string
}

func (s SourcePost) GetName() string {
	return s.name
}
