// Package in declares the application's published API: what the outside world
// may ask of it.
//
// The usual Go advice — let the consumer declare the interface it needs — is
// right when there is one consumer. This one has several: HTTP today, gRPC,
// GraphQL or a CLI later. Each declaring its own copy would duplicate one
// contract many times over and let the copies drift.
package in

import "context"

// CommandHandler handles a use case that changes state.
//
// Commands return as little as possible, an identifier rather than a view, so
// that reads and writes can be routed differently later. Nothing implements this
// yet; it is declared beside QueryHandler because the pair is the contract, and
// adding it later invites a second, ad-hoc shape.
type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

// QueryHandler handles a use case that reads.
//
// The read side does not go through the aggregate: a query has its own view in
// dto and its own reader in port/out, shaped for what the caller displays rather
// than for what the business rules need.
type QueryHandler[Q, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
