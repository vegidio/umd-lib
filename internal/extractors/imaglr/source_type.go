package imaglr

// SourceType is the interface that all source types implement.
// The isSourceType method is unexported to restrict implementation to the same package.
type SourceType interface {
	isSourceType()
}

// SourcePost represents a post source type.
type SourcePost struct {
	Id string
}

// isSourceType implements the SourceType interface for Post.
func (s SourcePost) isSourceType() {}
