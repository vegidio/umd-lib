package reddit

// SourceType is the interface that all source types implement.
// The isSourceType method is unexported to restrict implementation to the same package.
type SourceType interface {
	isSourceType()
}

// SourceSubmission represents a submission source type.
type SourceSubmission struct {
	Name string
	Id   string
}

// isSourceType implements the SourceType interface for Submission.
func (SourceSubmission) isSourceType() {}

// SourceUser represents a user source type.
type SourceUser struct {
	Name string
}

// isSourceType implements the SourceType interface for User.
func (SourceUser) isSourceType() {}

// SourceSubreddit represents a subreddit source type.
type SourceSubreddit struct {
	Name string
}

// isSourceType implements the SourceType interface for Subreddit.
func (SourceSubreddit) isSourceType() {}
