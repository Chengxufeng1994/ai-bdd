package errors

// Kind classifies why a use case failed, independently of any transport.
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
//
// KindUnclassified is the zero value on purpose. An Error a caller built
// without setting Kind — or one some future refactor forgets to set — falls
// through to an adapter's internal-server-error branch instead of silently
// becoming, say, a KindNotFound. Default-deny starts at the type's zero
// value, not at a review comment asking every caller to remember.
type Kind int

// The kinds a use case may classify a failure as. See Kind's doc comment for
// why these six exist and why the list stops here.
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
