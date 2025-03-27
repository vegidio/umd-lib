package fapello

// SourceType is the interface that all source types implement.
// The isSourceType method is unexported to restrict implementation to the same package.
type SourceType interface {
	isSourceType()
}

// SourceModel represents a model source type.
type SourceModel struct {
	Name string
}

// isSourceType implements the SourceType interface for Post.
func (s SourceModel) isSourceType() {}
