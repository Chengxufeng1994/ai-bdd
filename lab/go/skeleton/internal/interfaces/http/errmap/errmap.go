// Package errmap turns errors into HTTP's vocabulary: a status and an RFC 9457
// Problem document.
//
// Domain errors carry a kind, not a transport code, because the domain does not
// know which protocol is serving it. The table from kind to status belongs to
// each adapter, and this is HTTP's copy of it.
//
// StatusFor classifies an error into the status this adapter would answer
// with; ToInternalServerError renders the one body every unclassified error
// gets. An operation whose contract declares more than a 500 combines the
// two: switch on StatusFor to pick among its generated response types, and
// fall back to ToInternalServerError for whatever StatusFor could not place.
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
// get there. The caller is expected to log err itself.
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
