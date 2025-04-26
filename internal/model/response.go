package model

import "fmt"

// Response represents a response from a service.
type Response struct {
	// Url is the URL from which the response was obtained.
	Url string

	// Media is a list of Media objects associated with the response.
	Media []Media

	// Extractor is the type of extractor used to obtain the response.
	Extractor ExtractorType

	// Metadata contains additional metadata about the response.
	Metadata Metadata

	// Done is a channel used to signal when the media query is complete.
	Done chan error
}

// Error waits for the query to finish and returns any error that occurred during the process.
func (r *Response) Error() error {
	err := <-r.Done
	return err
}

func (r *Response) String() string {
	return fmt.Sprintf("{Url: %s, Media: %v, Extractor: %s, Metadata: %v}",
		r.Url, r.Media, r.Extractor, r.Metadata)
}
