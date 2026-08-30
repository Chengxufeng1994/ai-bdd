// Package errmap turns errors into HTTP's vocabulary: a status and an RFC 9457
// Problem document.
//
// Domain errors carry a kind, not a transport code, because the domain does not
// know which protocol is serving it. The table from kind to status belongs to
// each adapter, and this is HTTP's copy of it.
//
// # Growth point
//
// There is one branch today. The kind table arrives with the first domain error
// kind; until CLARIFY has produced one there is nothing to classify, and a table
// with invented entries would be a guess nobody asked for. When kinds exist,
// this package gains ToNotFound, ToUnprocessableEntity and so on — each
// returning the matching shared response component from apigen.
package errmap

import (
	"net/http"

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
