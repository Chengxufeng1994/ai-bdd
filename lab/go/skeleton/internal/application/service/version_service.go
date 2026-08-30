// Package service bundles handlers so that an adapter depends on one value
// instead of on a constructor parameter per use case.
//
// A service has fields and no methods. With no methods there is nowhere for
// logic to accumulate; anything needing a decision is a use case and belongs in
// handler. Unlike "keep it thin", that rule is visible in a diff.
//
// Adding a use case is one field and one line at the composition root. The
// forwarding-method alternative costs two edits in two files every time, all of
// it mechanical — and mechanical duplication is what drifts.
package service

import "skeleton/internal/application/port/in"

// VersionService bundles the use cases that report build information.
type VersionService struct {
	// GetVersion reports the version the running service was built from.
	GetVersion in.GetVersionUseCase
}
