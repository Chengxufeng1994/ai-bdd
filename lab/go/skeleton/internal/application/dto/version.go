// Package dto holds what the application hands back to its adapters.
//
// It imports nothing, which is what lets port/in, port/out and query all depend
// on it without a cycle. Each type carries what one use case needs and nothing
// more: a dto that mirrors an aggregate field for field re-creates the coupling
// it was introduced to break.
package dto

// Version is what the get-version use case returns.
type Version struct {
	// Value is the version the running service was built from.
	Value string
}
