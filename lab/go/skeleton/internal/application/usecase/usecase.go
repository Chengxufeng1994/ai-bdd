// Package usecase holds the shape every use case takes.
//
// CommandHandler and QueryHandler are declared here, once, generically, for
// every use case there will ever be. This package holds only that shape; it
// implements no ports itself. Its subpackages hold the use cases themselves,
// split by CQRS side: usecase/query for reads, usecase/command for writes.
//
// A use case orchestrates: it calls the ports it needs and returns a result. It
// holds no business rules — those live in domain — and it never returns a domain
// aggregate, because an adapter holding one could mutate it past its invariants.
package usecase

import "context"

// CommandHandler is the shape every use case that changes state takes.
//
// Commands return as little as possible, an identifier rather than a view, so
// that reads and writes can be routed differently later. Nothing implements this
// yet; it is declared beside QueryHandler because the pair is the contract, and
// adding it later invites a second, ad-hoc shape.
type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

// QueryHandler is the shape every use case that reads takes.
//
// The read side does not go through the aggregate: a query has its own result
// type declared beside it and its own reader in port/out, shaped for what the
// caller displays rather than for what the business rules need.
//
// This pair is also what lets service hold a use case by interface rather than
// by concrete type: a field typed QueryHandler[Q, R] accepts a decorated use
// case — one wrapped with logging or tracing — exactly as it accepts a plain
// one, because both satisfy the same generic shape.
type QueryHandler[Q, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
