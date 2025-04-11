package model

// Result is a generic struct that represents the result of an operation.
//
// Parameters:
//   - Data is a data of type T.
//   - Err is an error that indicates if the operation failed.
type Result[T any] struct {
	Data T
	Err  error
}
