// Package presenter converts application dto into HTTP responses.
//
// It is mechanical translation, outbound. A presenter may drop fields — each
// protocol exposes as little as it likes — but it never reaches past the dto
// into the domain.
package presenter

import (
	"skeleton/internal/application/dto"
	"skeleton/internal/interfaces/http/apigen"
)

// ToGetVersionResponse renders the running version as the 200 response.
func ToGetVersionResponse(d dto.Version) apigen.GetVersion200JSONResponse {
	return apigen.GetVersion200JSONResponse{Version: d.Value}
}
