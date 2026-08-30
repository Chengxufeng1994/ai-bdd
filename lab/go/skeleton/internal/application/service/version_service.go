// Package service implements the application's driving ports by delegating to
// use cases.
//
// A service holds each use case by interface rather than by concrete type, so
// that a decorated use case — one wrapped with logging or tracing — substitutes
// cleanly for a plain one without the service noticing. Each method is thin
// delegation and nothing more: anything that needs a decision belongs in the
// use case, not here, because a decision here would be a business rule with no
// scenario able to reach it.
package service

import (
	"context"

	"skeleton/internal/application/port/in"
	"skeleton/internal/application/query"
	"skeleton/internal/application/usecase"
)

// VersionService implements in.VersionService by delegating to use cases.
type VersionService struct {
	getVersion usecase.QueryHandler[query.GetVersion, query.GetVersionResult]
}

// Compile-time check that VersionService still satisfies the driving port.
// Without it, a signature drift would surface as a confusing error at the
// composition root instead of here.
var _ in.VersionService = VersionService{}

// NewVersionService returns a VersionService delegating to the given use case.
func NewVersionService(getVersion usecase.QueryHandler[query.GetVersion, query.GetVersionResult]) VersionService {
	return VersionService{getVersion: getVersion}
}

// GetVersion reports the version the running service was built from.
func (s VersionService) GetVersion(ctx context.Context, q query.GetVersion) (query.GetVersionResult, error) {
	return s.getVersion.Handle(ctx, q)
}
