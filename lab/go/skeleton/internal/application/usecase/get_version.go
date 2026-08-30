// Package usecase implements the driving ports declared in port/in.
//
// A use case orchestrates: it calls the ports it needs and returns a dto. It
// holds no business rules — those live in domain — and it never returns a domain
// aggregate, because an adapter holding one could mutate it past its invariants.
package usecase

import (
	"context"
	"fmt"

	"skeleton/internal/application/port/out"
	"skeleton/internal/application/query"
)

// GetVersion reports the version the running service was built from.
type GetVersion struct {
	versions out.VersionProvider
}

// Compile-time check that GetVersion still satisfies the generic shape it is
// held by. Without it, a signature drift would surface as a confusing error at
// the composition root instead of here.
var _ QueryHandler[query.GetVersion, query.GetVersionResult] = GetVersion{}

// NewGetVersion returns a GetVersion reading from the given provider.
func NewGetVersion(versions out.VersionProvider) GetVersion {
	return GetVersion{versions: versions}
}

// Handle returns the running version, or the provider's failure with context
// added.
func (h GetVersion) Handle(ctx context.Context, _ query.GetVersion) (query.GetVersionResult, error) {
	v, err := h.versions.Version(ctx)
	if err != nil {
		return query.GetVersionResult{}, fmt.Errorf("read build version: %w", err)
	}

	return query.GetVersionResult{Value: v}, nil
}
