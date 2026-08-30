package errors

// VersionUnavailable reports that the build version could not be read.
//
// Catalog entries are prototypes: they carry the identity a failure always has
// — its kind, its code and its message key, written together in one literal so
// they cannot be given values that disagree — and callers specialise a copy
// with what varies per occurrence. That agreement holds only here, in this one
// declaration: nothing in the type stops a caller writing
// Error{Code: "X", Kind: k} directly elsewhere and skipping the catalog
// entirely. What the catalog buys is a single place where the three are
// deliberately kept together, not a compiler-enforced one.
//
// It is the only entry because it is the only failure this codebase can
// currently produce: a driven port that did not answer. Entries are added when
// a use case can actually raise them, never to make the catalog look complete.
//
// MUST NOT: build a MessageKey by interpolation. A key assembled as
// "workout." + id + ".not_found" would be unbounded, and it would put user data
// into the RFC 9457 type — the field clients log and group by. Keys are
// declared here as constants; variable data belongs in Params.
var VersionUnavailable = Error{
	Kind:       KindUnavailable,
	Code:       "E0001",
	MessageKey: "version.unavailable",
}

// WithWhere returns a copy of e naming the operation it happened in.
//
// Where composes: each layer that wraps an Error calls WithWhere again with
// its own name, so unwrapping the chain yields a path — see Where's doc
// comment in errors.go.
func (e Error) WithWhere(where Where) Error {
	e.Where = where
	return e
}

// WithErr returns a copy of e wrapping the given failure.
//
// A catalog entry names what generally went wrong; WithErr attaches the
// specific occurrence, so the result wraps err and errors.As and errors.Unwrap
// can still reach it.
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
