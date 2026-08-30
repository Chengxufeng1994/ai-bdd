// Package mapper converts HTTP requests into application commands and queries.
//
// It is mechanical translation, inbound. Its counterpart is presenter, going the
// other way; the two are named apart so the direction is visible at a glance
// rather than inferred from a signature.
package mapper

import (
	"skeleton/internal/application/query"
	"skeleton/internal/interfaces/http/apigen"
)

// ToGetVersion builds the query behind GET /version.
//
// The request carries no parameters, so there is nothing to copy. The function
// exists anyway: it is where a parameter would land, and without it the inbound
// half of the conversion has no home and drifts into the handler — where the
// next person copying this slice would reasonably conclude it belongs.
func ToGetVersion(_ apigen.GetVersionRequestObject) query.GetVersion {
	return query.GetVersion{}
}
