package coomer

// SourceType is the interface that all source types implement.
// The isSourceType method is unexported to restrict implementation to the same package.
type SourceType interface {
	isSourceType()
}

// SourceUser represents a user source type.
type SourceUser struct {
	Service string
	User    string
}

// isSourceType implements the SourceType interface for User.
func (SourceUser) isSourceType() {}

// SourcePost represents a post source type.
type SourcePost struct {
	Service string
	User    string
	Id      string
}

// isSourceType implements the SourceType interface for Post.
func (SourcePost) isSourceType() {}
