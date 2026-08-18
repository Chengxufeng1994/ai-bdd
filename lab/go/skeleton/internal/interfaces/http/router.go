package http

import (
	"fmt"

	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/gin-gonic/gin"

	"skeleton/internal/interfaces/http/apigen"
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

	apigen.RegisterHandlers(engine, apigen.NewStrictHandler(si, nil))

	return engine, nil
}
