// Package query holds the read side's use cases: their inputs, their results
// and the handlers that connect them.
//
// Query types carry basic types and no json tags: they are semantic input to the
// application, not a wire format. A query naming an apigen type would make the
// application depend on HTTP and stop every other protocol from reusing it.
package query

import (
	"context"
	"fmt"

	"skeleton/internal/application/port/out"
	"skeleton/internal/application/usecase"
)

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

// getVersionHandler answers GetVersion by reading through out.VersionProvider
// and wrapping a provider failure with context.
//
// It is unexported: an adapter must import this package to name GetVersion and
// GetVersionResult, and would otherwise also see this concrete handler —
// internal/interfaces/doc.go forbids that. Returning an interface from a
// constructor is normally a Go smell; here it is the same case pkg/log's Logger
// documents — the interface is the package's contract and the concrete type has
// nothing extra to offer, so there is nothing to lose by hiding it.
type getVersionHandler struct {
	versions out.VersionProvider
}

// NewGetVersion returns a QueryHandler that reads the version from the given
// provider.
func NewGetVersion(versions out.VersionProvider) usecase.QueryHandler[GetVersion, GetVersionResult] {
	return getVersionHandler{versions: versions}
}

// Handle returns the running version, or the provider's failure with context
// added.
func (h getVersionHandler) Handle(ctx context.Context, _ GetVersion) (GetVersionResult, error) {
	v, err := h.versions.Version(ctx)
	if err != nil {
		return GetVersionResult{}, fmt.Errorf("read build version: %w", err)
	}

	return GetVersionResult{Value: v}, nil
}
