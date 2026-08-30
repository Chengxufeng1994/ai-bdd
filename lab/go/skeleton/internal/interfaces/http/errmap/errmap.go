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
// gets. An operation whose contract declares more than a 500 combines the
// two: switch on StatusFor to pick among its generated response types, and
// fall back to ToInternalServerError for whatever StatusFor could not place.
// See http/version.go's GetVersion for today's one-operation case, where that
// combination has no production caller yet.
package errmap

import (
	"errors"
	"net/http"

	apperrors "skeleton/internal/application/errors"
	"skeleton/internal/interfaces/http/apigen"
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
