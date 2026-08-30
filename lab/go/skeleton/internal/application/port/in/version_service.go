package in

import (
	"context"

	"skeleton/internal/application/usecase/query"
)

// VersionService is the application's published API for reporting build
// information.
//
// It is an interface, not a struct of handlers, because a driving port is a
// capability the application offers outward: an adapter should depend on what
// the application can do, not on how its use cases are assembled underneath.
// A struct of fields only lets an adapter fill in handlers one by one; an
// interface lets a whole implementation — including a decorated one — be
// substituted for another.
type VersionService interface {
	// GetVersion reports the version the running service was built from.
	GetVersion(ctx context.Context, q query.GetVersion) (query.GetVersionResult, error)
}
