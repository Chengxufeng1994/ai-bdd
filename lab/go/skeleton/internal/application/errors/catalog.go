package errors

// Descriptor is one failure's identity: the code support quotes, the key its
// message renders from, and the kind an adapter maps to a status.
//
// The three live together because they describe one failure. Set individually
// they would drift — a code updated without its key, a kind that no longer
// matches either — and nothing would report the disagreement.
type Descriptor struct {
	// Code is the stable identifier support tooling quotes.
	Code string

	// MessageKey identifies the message this failure renders to. An adapter
	// uses it both as the translation key and, rendered, as the RFC 9457
	// type a client branches on — so it is API surface, not an internal
	// string.
	MessageKey string

	// Kind is the classification an adapter maps to its own vocabulary. The
	// zero value, KindUnclassified, is treated by every adapter as an
	// unrecognised error.
	Kind Kind
}

// VersionUnavailable reports that the build version could not be read.
//
// It is the only entry here because it is the only failure this codebase can
// currently produce: a driven port that did not answer. Entries are added when
// a use case can actually raise them, never to make the catalog look complete.
//
// MUST NOT: build a MessageKey by interpolation. A key assembled as
// "workout." + id + ".not_found" would be unbounded, and it would put user
// data into the RFC 9457 type — the field clients log and group by. Keys are
// declared here as constants; variable data belongs in Params.
var VersionUnavailable = Descriptor{
	Code:       "E0001",
	MessageKey: "version.unavailable",
	Kind:       KindUnavailable,
}

// New builds an Error from a Descriptor.
//
// Embedding Descriptor in Error is what makes Code, MessageKey and Kind
// inseparable: there is no path that sets one without the other two. What
// this does not do is stop a caller writing Error{Descriptor: Descriptor{...}}
// inline instead of naming a catalog entry — that residue stays visible in
// review, which is the point; New exists so the common path never has a
// reason to reach for the struct literal at all.
func New(d Descriptor, where Where, err error) Error {
	return Error{Descriptor: d, Where: where, Err: err}
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
