package redgifs

// SourceVideo represents a video source type.
type SourceVideo struct {
	name string
}

func (s SourceVideo) GetName() string {
	return s.name
}

// SourceUser represents a user source type.
type SourceUser struct {
	name string
}

func (s SourceUser) GetName() string {
	return s.name
}
