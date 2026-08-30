package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/gin-gonic/gin"

	"skeleton/internal/interfaces/http/apigen"
	"skeleton/internal/interfaces/http/errmap"
)

// NewRouter wires the generated handlers onto a gin engine.
//
// This is the only hand-written file that names gin. Keeping routing separate
// from the handlers in server.go means middleware, engine configuration and
// route registration can change without touching request handling, and vice
// versa — the two churn for entirely different reasons.
//
// gin.New rather than gin.Default: Default installs the request logger, which
// turns an acceptance run into a wall of noise. Recovery is kept, because a
// panicking handler should become a 500 that a scenario can assert on rather
// than take the whole test process down.
func NewRouter(si apigen.StrictServerInterface) (*gin.Engine, error) {
	spec, err := apigen.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load embedded OpenAPI spec: %w", err)
	}

	// The spec carries server URLs that would otherwise be matched against the
	// request host, rejecting anything served from a different address.
	// Validation here is about the request shape, not about where it arrived.
	spec.Servers = nil

	engine := gin.New()
	engine.Use(gin.Recovery())

	// Shape validation comes from the contract, not from struct tags.
	//
	// This is the boundary layer's whole job: is the JSON well-formed, are the
	// required fields present, are the types right. A request that fails here
	// never reaches a handler, and the rule that rejected it exists in exactly
	// one place — api/openapi.yaml. Duplicating those rules as validator tags
	// would create a second source of truth that drifts from the spec the
	// clients were generated from.
	//
	// What this deliberately does not check is meaning. "Weight must not be
	// negative" is a domain invariant and belongs in a domain constructor; if it
	// also lived here it would have two homes and they would disagree.
	engine.Use(middleware.OapiRequestValidator(spec))

	strict := apigen.NewStrictHandlerWithOptions(si, nil, apigen.StrictGinServerOptions{
		HandlerErrorFunc: handlerErrorFunc,
	})
	apigen.RegisterHandlers(engine, strict)

	return engine, nil
}

// handlerErrorFunc replaces the HandlerErrorFunc that apigen.NewStrictHandler
// installs by default, which writes gin.H{"msg": err.Error()} straight into
// the response body — exactly the leak interfaces/doc.go's MUST forbids.
//
// This is the backstop, not the normal path. Every handler still returns a
// typed response and a nil error, so the contract stays compiler-checked; this
// only fires when one forgets, and turns that mistake into a generic 500
// rather than a leak. It reuses errmap.ToInternalServerError for the body so
// there is still exactly one place that owns the Problem shape.
func handlerErrorFunc(ctx *gin.Context, err error) {
	body, marshalErr := json.Marshal(errmap.ToInternalServerError(err))
	if marshalErr != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Data(http.StatusInternalServerError, "application/problem+json", body)
}
