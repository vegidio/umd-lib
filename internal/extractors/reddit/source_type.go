package reddit

// SourceSubmission represents a submission source type.
type SourceSubmission struct {
	Id string

	name string
}

func (s SourceSubmission) GetName() string {
	return s.name
}

// SourceUser represents a user source type.
type SourceUser struct {
	name string
}

func (s SourceUser) GetName() string {
	return s.name
}

// SourceSubreddit represents a subreddit source type.
type SourceSubreddit struct {
	name string
}

func (s SourceSubreddit) GetName() string {
	return s.name
}
