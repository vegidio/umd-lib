package redgifs

// SourceType is the interface that all source types implement.
// The isSourceType method is unexported to restrict implementation to the same package.
type SourceType interface {
	isSourceType()
}

// SourceVideo represents a video source type.
type SourceVideo struct {
	Id string
}

// isSourceType implements the SourceType interface for Video.
func (s SourceVideo) isSourceType() {}

// SourceUser represents a user source type.
type SourceUser struct {
	Name string
}

// isSourceType implements the SourceType interface for User.
func (s SourceUser) isSourceType() {}
