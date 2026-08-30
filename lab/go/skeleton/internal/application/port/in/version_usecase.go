package in

import (
	"skeleton/internal/application/dto"
	"skeleton/internal/application/query"
)

// GetVersionUseCase is the driving port for reporting the running version.
//
// It is an alias rather than a new interface type so that the generic shape
// survives: a decorator declared as WithLogging[Q, R](QueryHandler[Q, R]) still
// wraps it, and the service bundle's field type stays the contract itself. The
// cost is that any handler with the same signature satisfies it, which review
// has to catch rather than the compiler.
type GetVersionUseCase = QueryHandler[query.GetVersion, dto.Version]

// VersionUseCase bundles the use cases that report build information.
//
// It lives here, beside the port it bundles, rather than in a separate
// package: it is the application's published API, so an adapter holding it
// depends on contracts and nothing else. It is a struct with fields and no
// methods because every field is already a driving port — there is no
// behaviour left to abstract, and an interface here would only add a
// GetVersion() method and an indirection between the field and the call.
//
// Adding a use case is one field and one line at the composition root. The
// forwarding-method alternative costs two edits in two files every time, all
// of it mechanical — and mechanical duplication is what drifts.
type VersionUseCase struct {
	// GetVersion reports the version the running service was built from.
	GetVersion GetVersionUseCase
}
