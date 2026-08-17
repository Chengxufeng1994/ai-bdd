//go:build generate

// Package fitnesstracker exists only to hold the code-generation directives for
// the whole module.
//
// Every generator is listed here rather than beside the package it writes into,
// for three reasons:
//
//   - They all run from the same directory, so every path — spec, config,
//     output — is root-relative and reads the same way. Directives scattered
//     across packages each carry their own relative frame, and a config file's
//     `output` then only makes sense next to the directive that runs it.
//   - Repo-scoped generators have nowhere else to go. mockery, buf and sqlc read
//     a config describing the whole repository and expect to run from its root;
//     there is no package they belong to.
//   - One file is the inventory. "What generates this?" is answered by reading
//     it, not by grepping.
//
// Directives execute in the order they appear, which is how a dependency
// between generators is expressed — mocks come after the interfaces they stand
// in for. Where a generator is not a Go tool at all (a TypeScript client, say),
// it belongs in the Makefile instead; `make gen` is the entry point that covers
// both kinds.
//
// The build tag keeps this file out of every build, so the module root stays a
// package with no compiled contents. `go generate` ignores build constraints, so
// the directives still run under a plain `go generate ./...`.
package fitnesstracker

//go:generate go tool oapi-codegen -config api/cfg.yaml api/openapi.yaml

//go:generate go tool mockery
