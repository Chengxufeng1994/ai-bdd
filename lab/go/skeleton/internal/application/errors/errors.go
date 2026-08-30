// Package errors is the application layer's error vocabulary.
//
// A use case that fails classifies the failure by Kind, not by a transport
// code: the application does not know whether it is being called over HTTP,
// gRPC or a CLI, so it cannot pick a status. Each adapter's errmap package
// owns the table from Kind to its own protocol's vocabulary — see
// ../../interfaces/doc.go for HTTP's copy of it.
//
// # The kinds are transcribed, not invented
//
// Six kinds exist because six failure shapes are documented elsewhere in this
// codebase. Adding a seventh without a matching citation would be a guess
// dressed up as a classification.
//
//	KindNotFound      the addressed resource does not exist
//	KindInvalid       the request parsed but violates a domain rule (a 422,
//	                  not a 400 — see below)
//	KindConflict      the request conflicts with the current state
//	KindUnauthorized  no credentials were supplied, or they were not valid
//	KindForbidden     authenticated, but not allowed to perform this operation
//	KindUnavailable   a dependency the service needs is down
//
// NotFound, Invalid, Conflict and Unauthorized are the four the HTTP kind
// table in ../../interfaces/doc.go already lists.
//
// Forbidden and Unavailable are not new decisions; they close gaps in
// documents that already assumed them. DATAFLOW.md's three-layer validation
// table charges the application layer with checking "identifier format,
// permission, resource existence" and lists 401 and 403 side by side as its
// failure modes — the HTTP kind table carries 401 but had never been given
// 403 to carry. And api/openapi.yaml's ServiceUnavailable response is
// documented as "a dependency the service needs is down", which is precisely
// what a failing out.VersionProvider — or any other driven port — is; today
// it is the only kind of failure this codebase can actually produce.
//
// # Why there is no KindBadRequest
//
// 400 is deliberately absent. DATAFLOW.md's validation table assigns shape
// checking — malformed JSON, a missing required field, a value of the wrong
// type — to the interfaces layer, where OapiRequestValidator rejects the
// request against api/openapi.yaml before any handler runs. A use case never
// sees a request that shallow, so it can never produce "I cannot read this";
// a kind for it would sit in this type with no code able to ever set it.
//
// Success is absent for the same reason from the other direction: 200, 201
// and 204 are chosen by an operation's contract and its presenter, never by
// this package. Kind classifies a non-nil error and nothing else.
package errors

import "strings"

// Kind classifies why a use case failed, independently of any transport.
//
// KindUnclassified is the zero value on purpose. An Error a caller built
// without setting Kind — or one some future refactor forgets to set — falls
// through to an adapter's internal-server-error branch instead of silently
// becoming, say, a KindNotFound. Default-deny starts at the type's zero
// value, not at a review comment asking every caller to remember.
type Kind int

// The kinds a use case may classify a failure as. See the package doc for why
// these six exist and why the list stops here.
const (
	// KindUnclassified is the zero value: a failure nobody has classified.
	KindUnclassified Kind = iota

	// KindNotFound reports that the addressed resource does not exist.
	KindNotFound

	// KindInvalid reports that the request parsed but violates a domain rule.
	KindInvalid

	// KindConflict reports that the request conflicts with the current state.
	KindConflict

	// KindUnauthorized reports that no credentials were supplied, or they
	// were not valid.
	KindUnauthorized

	// KindForbidden reports that the caller is authenticated but not allowed
	// to perform this operation.
	KindForbidden

	// KindUnavailable reports that a dependency the service needs is down.
	KindUnavailable
)

// Where names the logical operation an Error was produced in, such as
// "usecase.GetVersion".
//
// It composes: each layer that wraps an Error adds its own, so unwrapping the
// chain yields a path — service.VersionService.GetVersion: usecase.GetVersion:
// read build version: connection refused. That is cheaper than a stack trace,
// which costs at capture time, and it survives refactoring, because it names
// what the code was doing rather than which frames it happened to be in.
type Where string

// Error is an application-layer failure classified by Kind.
//
// Two channels leave this type and they must not be confused. Params is
// rendered into the response a client receives; Where, Details and Err go to
// logs and never to a client. The split is what keeps a 500's body generic
// while an operator still gets everything.
type Error struct {
	// Kind classifies the failure. The zero value, KindUnclassified, is
	// treated by every adapter as an unrecognised error.
	Kind Kind

	// Code is the stable identifier support tooling quotes. A catalog entry
	// in catalog.go sets it together with MessageKey and Kind in one literal.
	Code string

	// MessageKey identifies the message this failure renders to. An adapter
	// uses it both as the translation key and, rendered, as the RFC 9457
	// type a client branches on — so it is API surface, not an internal
	// string. A catalog entry in catalog.go sets it together with Code and
	// Kind.
	MessageKey string

	// Params are the values a message template interpolates. They reach the
	// client, so they carry domain values — an identifier, a count, a unit —
	// and never the text of an underlying error.
	Params map[string]any

	// Where names the operation this failure was produced in. It names
	// internal structure and must never reach a client.
	Where Where

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

	for _, part := range []string{string(e.Where), e.Code, e.Details} {
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
