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
// StatusFor's classification, and falling back to ToInternalServerError's
// generic body for anything it does not recognise. An operation whose
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
// It is not logged either, and that is a gap rather than a decision. Nothing in
// this adapter holds a logger: Server has no field for one, router.go installs
// only gin.Recovery(), and the handler that calls this discards err. So a 500
// leaves the process with no trace at all today, and the context a use case
// attached on its way out — "read build version: %w" — is written nowhere. The
// parameter stays in the signature so that closing the gap is a change to this
// function and its callers, not a change to this package's shape.
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
// Only Params crosses into the body. Where, DetailedError and the wrapped
// cause are the log channel and are deliberately absent from what this
// returns; the caller is expected to log err itself.
func ToProblem(err error, locale string, tr i18n.Translator) (int, apigen.Problem) {
	var e apperrors.Error
	if !errors.As(err, &e) || e.Kind == apperrors.KindUnclassified {
		return http.StatusInternalServerError, apigen.Problem{
			Type:   unclassifiedProblemType,
			Title:  unclassifiedProblemTitle,
			Status: http.StatusInternalServerError,
		}
	}

	status := StatusFor(e)

	return status, apigen.Problem{
		Type:   problemTypeBase + "/" + strings.ReplaceAll(e.MessageKey, ".", "-"),
		Title:  tr.Translate(locale, e.MessageKey, e.Params),
		Status: int32(status),
	}
}
