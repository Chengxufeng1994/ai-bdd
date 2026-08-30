// Package usecase implements the driving ports declared in port/in.
//
// A use case orchestrates: it calls the ports it needs and returns a dto. It
// holds no business rules — those live in domain — and it never returns a domain
// aggregate, because an adapter holding one could mutate it past its invariants.
package usecase

import (
	"context"
	"fmt"

	"skeleton/internal/application/dto"
	"skeleton/internal/application/port/in"
	"skeleton/internal/application/port/out"
	"skeleton/internal/application/query"
)

// GetVersion reports the version the running service was built from.
type GetVersion struct {
	versions out.VersionProvider
}

// Compile-time check that GetVersion still satisfies the driving port. Without
// it, a signature drift would surface as a confusing error at the composition
// root instead of here.
var _ in.GetVersionUseCase = GetVersion{}

// NewGetVersion returns a GetVersion reading from the given provider.
func NewGetVersion(versions out.VersionProvider) GetVersion {
	return GetVersion{versions: versions}
}

// Handle returns the running version, or the provider's failure with context
// added.
func (h GetVersion) Handle(ctx context.Context, _ query.GetVersion) (dto.Version, error) {
	v, err := h.versions.Version(ctx)
	if err != nil {
		return dto.Version{}, fmt.Errorf("read build version: %w", err)
	}

	return dto.Version{Value: v}, nil
}
