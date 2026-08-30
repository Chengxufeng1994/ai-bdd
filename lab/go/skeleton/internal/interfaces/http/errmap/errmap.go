// Package errmap turns errors into HTTP's vocabulary: a status and an RFC 9457
// Problem document.
//
// Application errors carry a kind, not a transport code, because the
// application does not know which protocol is serving it. The kinds are
// declared in ../../../application/errors; the table from kind to status
// belongs to each adapter, and this is HTTP's copy of it.
//
// StatusFor classifies an error into the status this adapter would answer
// with; ToInternalServerError renders the one body every unclassified error
// gets. ToProblem is what a handler calls: it combines the two, translating a
// classified error's message into a Problem document whose Status is
// StatusFor's classification, and calling ToInternalServerError for the generic
// body of anything it does not recognise. It calls it rather than rebuilding
// the same literal: two documents assembled from the same two constants agree
// only until a field is added to one of them. An operation whose
// contract declares more than one failure response switches on the status
// ToProblem returns to pick among its generated response types; one whose
// contract declares only 500, like /version, coerces the body's Status back
// to 500 instead of forwarding it — see http/version.go's GetVersion for that
// concrete case.
package errmap

import (
	"errors"
	"net/http"
	"strings"

	apperrors "skeleton/internal/application/errors"
	"skeleton/internal/interfaces/http/apigen"
	"skeleton/pkg/i18n"
)

// unclassifiedProblemType is the RFC 9457 type for a failure with no more
// specific identity. RFC 9457 reserves "about:blank" for exactly this.
const unclassifiedProblemType = "about:blank"

// unclassifiedProblemTitle is the human-readable summary for the same.
const unclassifiedProblemTitle = "Internal Server Error"

// ToInternalServerError renders err as the shared 500 Problem document.
//
// The error is deliberately not copied into Detail. Stack traces, SQL fragments
// and internal hostnames belong in logs, not in a response body that reaches
// whoever made the request — and "just this once, to help debugging" is how they
// get there.
//
// It is not logged here, and that is a division of labour rather than the gap
// it once was. Every caller holds err at a point where it still carries the
// context worth recording, and writes that record itself: router.go's error
// hooks log before handing this body to the client, and version.go's GetVersion
// logs before rendering its failure at all. Logging again here would emit each
// of those failures twice, the second time under a message naming no operation.
//
// The parameter is unused and stays anyway. It is what makes a call site read
// as "render this error" rather than "render the 500", which is the difference
// between a caller that still has the error in hand and one that has already
// dropped it — and it keeps a later decision to lift something safe out of err,
// a correlation id say, a change to this function rather than to every
// caller's signature.
func ToInternalServerError(_ error) apigen.InternalServerErrorApplicationProblemPlusJSONResponse {
	return apigen.InternalServerErrorApplicationProblemPlusJSONResponse{
		Type:   unclassifiedProblemType,
		Title:  unclassifiedProblemTitle,
		Status: http.StatusInternalServerError,
	}
}

// StatusFor reports the HTTP status err classifies to.
//
// The order is a security property, not a style choice: recognise what we
// classify, and let everything else — a kind this table has no entry for, an
// apperrors.Error with the zero-value KindUnclassified, or an error that
// carries no classification at all — fall to http.StatusInternalServerError.
// A default that instead passed an unrecognised error through unmapped would
// leak err.Error() to a client the first time an unclassified error reached
// this function; recognise-then-deny is what rules that out.
func StatusFor(err error) int {
	var e apperrors.Error
	if !errors.As(err, &e) {
		return http.StatusInternalServerError
	}

	switch e.Kind {
	case apperrors.KindNotFound:
		return http.StatusNotFound
	case apperrors.KindInvalid:
		return http.StatusUnprocessableEntity
	case apperrors.KindConflict:
		return http.StatusConflict
	case apperrors.KindUnauthorized:
		return http.StatusUnauthorized
	case apperrors.KindForbidden:
		return http.StatusForbidden
	case apperrors.KindUnavailable:
		return http.StatusServiceUnavailable
	default:
		// Also reached by KindUnclassified, the zero value: an Error nobody
		// classified is exactly the case this default must not miss.
		return http.StatusInternalServerError
	}
}

// problemTypeBase prefixes a message key to form an RFC 9457 type.
const problemTypeBase = "https://errors.skeleton.local"

// ToProblem renders err as the status and Problem document this adapter
// answers with, translating its message into locale.
//
// It is the top half of recognise-then-default: an error carrying a
// classification gets that classification's status and its own message; an
// error carrying none — including one whose Kind is the zero value,
// KindUnclassified — gets the same generic 500 ToInternalServerError
// produces. Reversing that order would leak err.Error() to a client the
// first time an unclassified error arrived.
//
// Carrying a classification means both halves of one. Kind selects the status,
// but MessageKey is what supplies the two fields a client reads: the type it
// branches on and the title it shows. An Error with a Kind and no MessageKey
// would render a document that satisfies the schema and identifies nothing —
// a type that is the base URL and a trailing slash, and an empty title the
// schema marks required. That is the same "somebody forgot" KindUnclassified
// stands for reaching this function through the other field, so it takes the
// same default path.
//
// Only Params crosses into the body. Where, DetailedError and the wrapped
// cause are the log channel and are deliberately absent from what this
// returns; the caller is expected to log err itself.
func ToProblem(err error, locale string, tr i18n.Translator) (int, apigen.Problem) {
	var e apperrors.Error
	if !errors.As(err, &e) || e.Kind == apperrors.KindUnclassified || e.MessageKey == "" {
		// The generic body comes from ToInternalServerError rather than from a
		// second literal here. The generated response type is a defined type
		// whose underlying type is Problem, so the conversion is free — and it
		// makes "the body an unrecognised error gets" one definition instead of
		// two that happen to read the same two constants.
		return http.StatusInternalServerError, apigen.Problem(ToInternalServerError(err))
	}

	status := StatusFor(e)

	return status, apigen.Problem{
		Type:   problemTypeBase + "/" + strings.ReplaceAll(e.MessageKey, ".", "-"),
		Title:  tr.Translate(locale, e.MessageKey, e.Params),
		Status: int32(status),
	}
}
