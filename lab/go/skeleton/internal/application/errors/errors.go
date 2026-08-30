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
// file should almost never change. This file holds Error itself — its fields
// and its methods — which changes only when the type's shape evolves.
//
// # There is no catalog
//
// MessageKey is API surface: an adapter renders it into the RFC 9457 `type`
// a client branches on, so changing one is a breaking change that looks, at
// the call site, like renaming a local string constant. A catalog — one file
// listing every Kind/MessageKey pair a service can produce — is what would
// make that surface reviewable as a single list instead of scattered
// literals.
//
// There isn't one here. The one failure this codebase currently produces
// sets its Kind and MessageKey inline at the call site that raises it
// (see usecase/query.GetVersion), the same way this project has already left
// assembler/, command/, readmodel/ and a Violations field unbuilt: a
// container built before there is anything to put in it. That is affordable
// with one key and no consumer outside this service; it stops being
// affordable at the second literal MessageKey, or the first client outside
// this service that needs the list to stay stable.
package errors

import "strings"

// Error is an application-layer failure classified by Kind.
//
// Two channels leave this type and they must not be confused. Params is
// rendered into the response a client receives; Where, DetailedError and Err
// go to logs and never to a client. The split is what keeps a 500's body
// generic while an operator still gets everything.
//
// MessageKey and Where are both plain string. Nothing but the field name in
// a literal stops one from being given the other's value — they are
// adjacent identifier strings with no distinct type keeping them apart. That
// is acceptable today because every call site fills this struct with one
// literal naming each field explicitly, so a mismatch would be a copy-paste
// mistake visible in the literal, not an invisible argument-order trap.
// Revisit this if a helper ever takes both as parameters.
//
// Do not nest an Error inside another Error's Err field. errors.As stops at
// the first Error it finds in a chain, so a layer that writes
// Error{Where: "service.X", Err: innerAppErr} — leaving its own Kind at the
// zero value — silently downgrades whatever classification innerAppErr
// carried to KindUnclassified; nothing in the type or the compiler catches
// it. A layer that wants to add context to a failure it did not itself
// classify should wrap with fmt.Errorf("...: %w", err) instead: that
// wrapper has Unwrap but is not an Error, so errors.As passes through it to
// the classification underneath.
type Error struct {
	// Kind classifies the failure. The zero value, KindUnclassified, is
	// treated by every adapter as an unrecognised error.
	Kind Kind

	// MessageKey identifies the message this failure renders to. An adapter
	// uses it both as the translation key and, rendered, as the RFC 9457
	// type a client branches on — the machine-readable identity of the
	// failure and what clients branch on, so it is API surface, not an
	// internal string. It is set together with Kind in one literal at the
	// call site that raises the failure.
	//
	// Never build one by interpolation. fmt.Sprintf("workout.%s.not_found",
	// id) reads like a more precise key and is a different thing entirely:
	// the RFC 9457 type is what clients log, group by and alert on, so a key
	// carrying a request's own data gives every occurrence a distinct
	// identity — nothing groups, no alert can be written against it, and the
	// value is republished into a field no log pipeline knows to redact.
	// Values that vary per occurrence belong in Params, which exists for
	// exactly that and is rendered into the message rather than its identity.
	MessageKey string

	// Params are the values a message template interpolates. They reach the
	// client, so they carry domain values — an identifier, a count, a unit —
	// and never the text of an underlying error.
	//
	// Set by struct literal, this field is stored by reference: writing
	// Params: m aliases whatever map m names, so a later mutation of m would
	// silently change what a past failure says it was. That is safe today
	// because every call site writes an inline map[string]any{...} literal
	// with no other holder; revisit this the first time a call site builds
	// Params in a variable before assigning it.
	Params map[string]any

	// Where names the logical operation this failure was produced in, such
	// as "usecase.GetVersion". It names a logical operation rather than a
	// stack frame: that costs nothing at run time and survives refactoring,
	// where a file:line does not.
	//
	// The json:"-" tag is a guard, not an active rule — this type is not
	// marshalled today, because what reaches a client is the Problem
	// document errmap renders, never the error itself. The tag is here so
	// that the day someone does marshal an Error — for structured logging,
	// or to pass a failure between services — Where is already excluded,
	// rather than relying on whoever writes that code remembering that it
	// names internal structure a caller must not see.
	Where string `json:"-"`

	// DetailedError is internal context the caller knows that Err does not
	// say: how many retries, over what timeout, against which upstream. The
	// name is borrowed from Mattermost's AppError.DetailedError, whose own
	// comment calls it "internal error string to help the developer" —
	// exactly this field's job. The contract is not borrowed with it: their
	// field is tagged to reach the client, with an opt-out to clear it
	// selectively. Here ARCHITECTURE.md §7 requires a 500 body to stay
	// generic, and errmap_test.go enforces that by scanning the whole
	// rendered document — so DetailedError stays log-only, unconditionally,
	// with no tag or opt-out to later "restore." If it would merely restate
	// Err, leave it empty.
	DetailedError string

	// Err is the underlying failure this Error classifies.
	Err error
}

// Error renders e for logs: where it happened, what detail the caller
// attached, and what caused it. It is not what a client sees — see the HTTP
// adapter's errmap for that — so it deliberately says more.
func (e Error) Error() string {
	var b strings.Builder

	for _, part := range []string{e.Where, e.DetailedError} {
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

// Wrap returns a copy of e with the given failure as its cause.
//
// Wrap pairs with Unwrap: it sets the cause this Error classifies, and
// Unwrap is what exposes that cause to errors.Is and errors.As.
func (e Error) Wrap(err error) Error {
	e.Err = err
	return e
}
