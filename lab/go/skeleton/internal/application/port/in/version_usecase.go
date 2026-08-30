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
