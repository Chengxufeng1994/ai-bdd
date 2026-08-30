// Package domain holds the business rules and entities of the fitness tracker.
//
// It depends on nothing. Every other layer may import domain; domain imports no
// other layer and no framework. If a rule cannot be expressed without a
// database or an HTTP request, it does not belong here.
//
// # Structure
//
// Packages are split by aggregate, never by pattern:
//
//	domain/
//	├── workout/           one aggregate, everything about it
//	│   ├── workout.go     the aggregate root
//	│   ├── set.go         entities and value objects it owns
//	│   ├── volume.go      a domain service, when a rule belongs to no single entity
//	│   ├── errors.go      the failures this aggregate can produce
//	│   └── events.go      what it announces when it changes
//	├── exercise/
//	└── shared/            value objects genuinely used by several aggregates
//
// The tempting alternative — domain/entities, domain/valueobjects,
// domain/services — scatters one business concept across four directories, so
// changing a single rule means editing four packages and nothing tells a reader
// which files belong together. The test for a correct split: deleting one
// directory should delete exactly one business concept.
//
// # Aggregate boundaries
//
// An aggregate is a consistency boundary: everything inside it is saved or
// rejected together, and everything outside it is referenced by identity, never
// by pointer. Two aggregates that always change together are one aggregate; an
// aggregate that is only ever partly loaded is probably two.
//
// # Errors
//
// A failure originating here is classified with application/errors' Kind, the
// same vocabulary every other layer uses — see ../application/errors/errors.go.
// There is no domain-local kind: one vocabulary beats two that would drift
// apart. Domain still has no idea it is being served over HTTP, gRPC or a CLI;
// the adapter's kind → protocol table, not domain, is the only mapping to a
// transport code.
//
// # Not yet populated
//
// No aggregate exists here yet. Which aggregates the fitness tracker has is a
// question CLARIFY answers — inventing them first would be guessing at exactly
// the thing the testbed is meant to test. See ../../README.md.
package domain
