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

	// MessageKey is the translation key, and the identity a client branches
	// on once an adapter renders it as an RFC 9457 type.
	MessageKey string

	// Kind is the classification an adapter maps to its own vocabulary.
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
// This is the only way Code, MessageKey and Kind are set. Filling them
// individually is what the Descriptor exists to prevent, so no setter for them
// is offered — a caller who needs a new combination adds a Descriptor, where
// the three stay together and a reviewer sees them at once.
func New(d Descriptor, where Where, err error) Error {
	return Error{
		Kind:       d.Kind,
		Code:       d.Code,
		MessageKey: d.MessageKey,
		Where:      where,
		Err:        err,
	}
}

// WithParams returns a copy of e carrying the values its message interpolates.
//
// It returns a copy rather than mutating so that a Descriptor-built Error can
// be shared and specialised without one caller's params reaching another's.
func (e Error) WithParams(params map[string]any) Error {
	e.Params = params
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
