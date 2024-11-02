package models

// Event is the interface that all event types implement.
// The unexported isEvent method restricts implementation to the same package.
type Event interface {
	isEvent()
}

// OnExtractorFound represents an event when an extractor is found.
type OnExtractorFound struct {
	// Name of the extractor found.
	Name string
}

// isEvent implements the Event interface for OnExtractorFound.
func (OnExtractorFound) isEvent() {}

// OnExtractorTypeFound represents an event when an extractor type is found.
type OnExtractorTypeFound struct {
	// Type is the type of the extractor found.
	Type string
	// Name is the name of the extractor found.
	Name string
}

// isEvent implements the Event interface for OnExtractorTypeFound.
func (OnExtractorTypeFound) isEvent() {}

// OnMediaQueried represents an event when media is queried.
type OnMediaQueried struct {
	// Amount is the number of media items queried.
	Amount int
}

// isEvent implements the Event interface for OnMediaQueried.
func (OnMediaQueried) isEvent() {}

// OnQueryCompleted represents an event when a query is completed.
type OnQueryCompleted struct {
	// Total is the total number of items queried.
	Total int
}

// isEvent implements the Event interface for OnQueryCompleted.
func (OnQueryCompleted) isEvent() {}
