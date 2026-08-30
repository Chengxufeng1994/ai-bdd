// Package out declares what the application needs from the outside world.
//
// Every interface here is a requirement the application states and that
// infrastructure complies with. An interface declared in infrastructure and
// imported here would have the dependency backwards, and would drag storage and
// transport concepts into use cases.
package out

import "context"

// VersionProvider reports the version the running service was built from.
//
// It is a port rather than a plain string because the source lies outside the
// application: a build-time stamp today, possibly a file or a configuration
// service later. Those can fail, which is why Version returns an error even
// though the current implementation never does.
type VersionProvider interface {
	Version(ctx context.Context) (string, error)
}
