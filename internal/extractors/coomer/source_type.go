package coomer

// SourceUser represents a user source type.
type SourceUser struct {
	Service string

	name string
}

func (s SourceUser) GetName() string {
	return s.name
}

// SourcePost represents a post source type.
type SourcePost struct {
	Service string
	Id      string

	name string
}

func (s SourcePost) GetName() string {
	return s.name
}
