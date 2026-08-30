// Package version carries the version stamped into the binary at build time.
//
// It sits in pkg/ rather than internal/ for the same reason pkg/config does: a
// build-time fact belongs to no layer. The practical reason matters more — the
// linker's -X flag targets a symbol path, so putting this under
// internal/infrastructure would break the stamp whenever that tree is
// reorganised, and break it silently.
package version

// value is overridden at build time:
//
//	go build -ldflags "-X skeleton/pkg/version.value=$(git describe --tags)"
//
// The default is deliberately not a plausible version number: seeing "dev" in a
// deployed environment should look wrong immediately.
var value = "dev"

// Build returns the version stamped into this binary.
//
// Only the composition root may call it. Every layer below receives the version
// by injection, because the acceptance suite runs concurrent scenarios at
// different versions inside one process.
func Build() string { return value }
