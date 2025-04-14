package fapello

// SourceModel represents a model source type.
type SourceModel struct {
	name string
}

func (s SourceModel) GetName() string {
	return s.name
}
