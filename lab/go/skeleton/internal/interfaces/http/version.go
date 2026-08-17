package http

import (
	"context"

	"skeleton/internal/interfaces/http/apigen"
)

// GetVersion reports the version the running binary was built from.
//
// It reads the value captured at construction rather than a package variable so
// that two servers at different versions can coexist in one process — which the
// acceptance suite relies on, since scenarios run concurrently.
func (s *Server) GetVersion(_ context.Context, _ apigen.GetVersionRequestObject) (apigen.GetVersionResponseObject, error) {
	return apigen.GetVersion200JSONResponse{Version: s.version}, nil
}
