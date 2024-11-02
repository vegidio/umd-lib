package models

// Response represents a response from a service.
type Response struct {
	// Url is the URL from which the response was obtained.
	Url string

	// Media is a list of Media objects associated with the response.
	Media []Media

	// Extractor is the type of extractor used to obtain the response.
	Extractor ExtractorType

	// Metadata contains additional metadata about the response.
	Metadata map[string]interface{}
}
