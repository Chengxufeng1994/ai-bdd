package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/gin-gonic/gin"

	"skeleton/internal/interfaces/http/apigen"
	"skeleton/internal/interfaces/http/errmap"
	"skeleton/pkg/log"
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
//
// logger is not for request logging — it is what the error hooks below write
// to. Every one of them answers with a body that deliberately says nothing, so
// without a logger here a backstop 500 would leave the process with no trace
// of why it happened.
func NewRouter(si apigen.StrictServerInterface, logger log.Logger) (*gin.Engine, error) {
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

	// All three hooks, not just the one a handler can trip. Every field left
	// nil here is filled by NewStrictHandlerWithOptions with a default that
	// writes gin.H{"msg": err.Error()} into the body, so setting two of three
	// leaves the third leaking — and the third is reached by failures nobody
	// writes code for: a response that fails to serialise, a connection the
	// client dropped mid-write. Those carry socket addresses and internal
	// hostnames.
	//
	// RequestErrorHandlerFunc's default answers 400 rather than 500, and
	// collapsing it into the generic 500 loses that distinction. It costs
	// nothing real: OapiRequestValidator above rejects a request that does not
	// match api/openapi.yaml before the strict handler is reached, so a
	// genuine "I cannot read this" is already a 400 from the validator, and
	// what remains here is a decode failure on a request the spec accepted —
	// a server-side surprise, not a client mistake.
	onError := errorFunc(logger)
	strict := apigen.NewStrictHandlerWithOptions(si, nil, apigen.StrictGinServerOptions{
		RequestErrorHandlerFunc:  onError,
		HandlerErrorFunc:         onError,
		ResponseErrorHandlerFunc: onError,
	})
	apigen.RegisterHandlers(engine, strict)

	return engine, nil
}

// errorFunc builds the replacement for every error hook apigen installs by
// default. Each default writes gin.H{"msg": err.Error()} straight into the
// response body — exactly the leak interfaces/doc.go's MUST forbids.
//
// This is the backstop, not the normal path. Every handler still returns a
// typed response and a nil error, so the contract stays compiler-checked; this
// only fires when one forgets, or when the generated glue fails somewhere no
// handler can see. It reuses errmap.ToInternalServerError for the body so there
// is still exactly one place that owns the Problem shape.
//
// The body says nothing on purpose, so the record is not optional: it is the
// only account of a failure that reached the boundary uncategorised, and it is
// written here because this is the last place the error still exists.
func errorFunc(logger log.Logger) func(*gin.Context, error) {
	return func(ctx *gin.Context, err error) {
		logger.Error("unhandled error at the http boundary", "error", err)

		body, marshalErr := json.Marshal(errmap.ToInternalServerError(err))
		if marshalErr != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.Data(http.StatusInternalServerError, "application/problem+json", body)
	}
}
