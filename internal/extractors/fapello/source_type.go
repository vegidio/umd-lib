package fapello

// SourceModel represents a model source type.
type SourceModel struct {
	name string
}

func (s SourceModel) Type() string {
	return "Model"
}

func (s SourceModel) Name() string {
	return s.name
}
