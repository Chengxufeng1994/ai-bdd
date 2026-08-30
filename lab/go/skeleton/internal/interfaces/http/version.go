package http

import (
	"context"
	"net/http"

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
	v, err := s.svc.GetVersion(ctx, mapper.ToGetVersion(request))
	if err != nil {
		// The error is logged here and nowhere else: errmap renders only what
		// a client may see, so this is the one place Where, DetailedError and
		// the cause are still in hand.
		s.logger.Error("get version", "error", err)

		// api/openapi.yaml declares only 200 and 500 for this operation, so
		// oapi-codegen never generated a 404/422/409/... response type for it
		// to return here — there is nothing errmap.StatusFor's classification
		// could select among, which is why the status is discarded rather than
		// switched on. An operation whose contract declares more than one
		// failure response switches on the status ToProblem returns to choose
		// among its generated response types instead.
		_, problem := errmap.ToProblem(err, localeOf(ctx), s.tr)

		// Unconditional: every failure here renders as 500 whatever its kind,
		// and RFC 9457 §3.1.4 requires the body to say the same number as the
		// status line it travels under. The classification is not lost — Type
		// still carries the specific identity, which is what a client branches
		// on.
		problem.Status = http.StatusInternalServerError

		return apigen.GetVersion500ApplicationProblemPlusJSONResponse{
			InternalServerErrorApplicationProblemPlusJSONResponse: apigen.InternalServerErrorApplicationProblemPlusJSONResponse(problem),
		}, nil
	}

	return presenter.ToGetVersionResponse(v), nil
}

// localeOf reports the locale a request asked for.
//
// The strict handler hands operations a typed request rather than the HTTP
// one, so there is no Accept-Language header here to read. Returning a fixed
// locale keeps the rendering path honest — it still exercises translation —
// until a middleware puts the negotiated locale in the context.
func localeOf(_ context.Context) string { return "en" }
