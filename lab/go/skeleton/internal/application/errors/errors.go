// Package errors is the application layer's error vocabulary.
//
// A use case that fails classifies the failure by Kind, not by a transport
// code: the application does not know whether it is being called over HTTP,
// gRPC or a CLI, so it cannot pick a status. Each adapter's errmap package
// owns the table from Kind to its own protocol's vocabulary — see
// ../../interfaces/doc.go for HTTP's copy of it.
//
// The package is split by how often each part changes. kind.go holds Kind and
// its constants: adding a classification is a significant decision, so that
// file should almost never change. This file holds Error itself — its
// fields, its methods, and Wrap — which changes only when the type's shape
// evolves.
//
// # There is no catalog
//
// MessageKey is API surface: an adapter renders it into the RFC 9457 `type`
// a client branches on, so changing one is a breaking change that looks, at
// the call site, like renaming a local string constant. A catalog — one file
// listing every Kind/Code/MessageKey triple a service can produce — is what
// would make that surface reviewable as a single list instead of scattered
// literals.
//
// There isn't one here. The one failure this codebase currently produces
// sets its Kind, Code and MessageKey inline at the call site that raises it
// (see usecase/query.GetVersion), the same way this project has already left
// assembler/, command/, readmodel/ and a Violations field unbuilt: a
// container built before there is anything to put in it. That is affordable
// with one key and no consumer outside this service; it stops being
// affordable at the second literal MessageKey, or the first client outside
// this service that needs the list to stay stable.
package errors

import (
	// Aliased because this package is itself named errors, and Kind, Error
	// and the four With* methods below stay unqualified.
	stderrors "errors"
	"strings"
)

// Error is an application-layer failure classified by Kind.
//
// Two channels leave this type and they must not be confused. Params is
// rendered into the response a client receives; Where, Details and Err go to
// logs and never to a client. The split is what keeps a 500's body generic
// while an operator still gets everything.
//
// Code, MessageKey and Where are all plain string. Nothing but the parameter
// name stops WithWhere(e.MessageKey) from compiling: the three are adjacent
// identifier strings with no distinct type keeping them apart. That is
// acceptable today because no variable flows between them — Code and
// MessageKey are set together in one literal at whichever call site raises
// the failure, and Where only ever arrives as WithWhere's argument. Revisit
// this the day a function needs two of them as parameters at once, where
// argument order becomes an invisible trap.
type Error struct {
	// Kind classifies the failure. The zero value, KindUnclassified, is
	// treated by every adapter as an unrecognised error.
	Kind Kind

	// Code is the stable identifier support tooling quotes. It is set
	// together with MessageKey and Kind in one literal at the call site
	// that raises the failure.
	Code string

	// MessageKey identifies the message this failure renders to. An adapter
	// uses it both as the translation key and, rendered, as the RFC 9457
	// type a client branches on — so it is API surface, not an internal
	// string. It is set together with Code and Kind.
	MessageKey string

	// Params are the values a message template interpolates. They reach the
	// client, so they carry domain values — an identifier, a count, a unit —
	// and never the text of an underlying error.
	Params map[string]any

	// Where names the operation this failure was produced in. It names
	// internal structure and must never reach a client.
	Where string

	// Details is context the caller knows that Err does not say: how many
	// retries, over what timeout, against which upstream. It goes to logs
	// only. If it would merely restate Err, leave it empty.
	Details string

	// Err is the underlying failure this Error classifies.
	Err error
}

// Error renders e for logs: where it happened, which code it carries, what the
// caller observed, and what caused it. It is not what a client sees — see the
// HTTP adapter's errmap for that — so it deliberately says more.
func (e Error) Error() string {
	var b strings.Builder

	for _, part := range []string{e.Where, e.Code, e.Details} {
		if part != "" {
			b.WriteString(part)
			b.WriteString(": ")
		}
	}

	if e.Err != nil {
		b.WriteString(e.Err.Error())
	} else {
		b.WriteString(e.MessageKey)
	}

	return b.String()
}

// Unwrap returns the failure Error classifies, so that errors.Is and
// errors.As can see through it to whatever a driven port returned.
func (e Error) Unwrap() error {
	return e.Err
}

// WithWhere returns a copy of e naming the operation it happened in.
//
// Where composes: each layer that wraps an Error calls WithWhere again with
// its own name, so unwrapping the chain yields a path — see Where's doc
// comment above.
func (e Error) WithWhere(where string) Error {
	e.Where = where
	return e
}

// WithErr returns a copy of e wrapping the given failure.
//
// A caller names what generally went wrong; WithErr attaches the specific
// occurrence, so the result wraps err and errors.As and errors.Unwrap can
// still reach it.
func (e Error) WithErr(err error) Error {
	e.Err = err
	return e
}

// WithParams returns a copy of e carrying the values its message interpolates.
//
// The map is copied because an Error outlives the call that built it: it is
// wrapped, returned, and logged after the caller has moved on. Storing the
// caller's map by reference would let a later mutation there change what a
// past failure says it was.
func (e Error) WithParams(params map[string]any) Error {
	copied := make(map[string]any, len(params))
	for k, v := range params {
		copied[k] = v
	}
	e.Params = copied
	return e
}

// WithDetails returns a copy of e carrying operator context.
//
// Details goes to logs only. Say what the caller knows that Err does not; if it
// would restate Err, leave it empty.
func (e Error) WithDetails(details string) Error {
	e.Details = details
	return e
}

// Wrap returns err annotated with where, preserving the classification beneath it.
//
// Identity travels outward and location does not: the wrapper carries the same
// Kind, Code, MessageKey and Params as the error it wraps, because it is the same
// failure seen from one layer further out. Where and Details stay per-layer —
// each names what that layer knows.
//
// Copying the identity is what keeps errors.As correct. As stops at the first
// Error in the chain, so a wrapper that left Kind at its zero value would make an
// adapter classify a NotFound as an unclassified 500 — silently, since the inner
// Error is still there, just never reached.
func Wrap(err error, where string) error {
	if err == nil {
		return nil
	}

	var e Error
	if stderrors.As(err, &e) {
		return Error{
			Kind:       e.Kind,
			Code:       e.Code,
			MessageKey: e.MessageKey,
			Params:     e.Params,
			Where:      where,
			Err:        err,
		}
	}

	return Error{Where: where, Err: err}
}
