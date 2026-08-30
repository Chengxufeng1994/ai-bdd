// Package query holds the inputs to read-side use cases.
//
// Query types carry basic types and no json tags: they are semantic input to the
// application, not a wire format. A query naming an apigen type would make the
// application depend on HTTP and stop every other protocol from reusing it.
package query

// GetVersion asks for the version the running service was built from.
//
// It has no fields. The type exists anyway so that every use case has the same
// shape, which is what lets one decorator wrap them all.
type GetVersion struct{}

// GetVersionResult is what the get-version use case returns.
type GetVersionResult struct {
	// Value is the version the running service was built from.
	Value string
}
