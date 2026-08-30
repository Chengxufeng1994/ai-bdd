package http

import (
	"context"

	"skeleton/internal/interfaces/http/apigen"
	"skeleton/internal/interfaces/http/errmap"
	"skeleton/internal/interfaces/http/mapper"
	"skeleton/internal/interfaces/http/presenter"
)

// GetVersion reports the version the running binary was built from.
//
// The handler does four things and decides nothing: map the request, call the
// port, present the result, translate any failure. A decision here would be a
// business rule sitting in the wrong layer.
//
// A failure is returned as a typed 500 response with a nil error, not as a
// returned error. Returning the error would hand the response shape to
// oapi-codegen's default handler, and the body would stop being the one
// api/openapi.yaml promises.
func (s *Server) GetVersion(ctx context.Context, request apigen.GetVersionRequestObject) (apigen.GetVersionResponseObject, error) {
	v, err := s.svc.GetVersion.Handle(ctx, mapper.ToGetVersion(request))
	if err != nil {
		return apigen.GetVersion500ApplicationProblemPlusJSONResponse{
			InternalServerErrorApplicationProblemPlusJSONResponse: errmap.ToInternalServerError(err),
		}, nil
	}

	return presenter.ToGetVersionResponse(v), nil
}
