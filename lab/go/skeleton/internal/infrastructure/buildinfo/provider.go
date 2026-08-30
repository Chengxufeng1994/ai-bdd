// Package buildinfo adapts build-time facts onto the ports the application
// declares.
package buildinfo

import (
	"context"

	"skeleton/internal/application/port/out"
)

// Provider reports a version fixed at construction.
//
// The version is a field rather than a package variable so that two providers at
// different versions can coexist in one process, which the acceptance suite
// relies on since scenarios run concurrently.
type Provider struct {
	version string
}

// Compile-time check that Provider still satisfies the port it exists to
// implement.
var _ out.VersionProvider = Provider{}

// NewProvider returns a Provider reporting the given version.
func NewProvider(version string) Provider {
	return Provider{version: version}
}

// Version returns the version this Provider was constructed with.
//
// It never fails. The port declares an error because the source may later be a
// file or a configuration service; this implementation returns a value it
// already holds.
func (p Provider) Version(_ context.Context) (string, error) {
	return p.version, nil
}
