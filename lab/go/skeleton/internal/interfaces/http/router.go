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

	// All four generated error paths, not just the one a handler can trip.
	// Three are hooks on the strict handler; the fourth is
	// GinServerOptions.ErrorHandler, which the parameter-binding wrappers call
	// and which RegisterHandlers leaves unset by passing an empty options
	// struct. Each one left unset keeps a default that writes
	// gin.H{"msg": err.Error()} into the body, so closing three is not "the
	// leak is closed" — and the ones easiest to forget are exactly the ones
	// reached by failures nobody writes code for: a parameter that will not
	// parse, a response that fails to serialise, a connection the client
	// dropped mid-write. Those errors carry socket addresses and internal
	// hostnames.
	//
	// Two of the four have no generated call site while /version is the only
	// operation, because it declares no parameters and no request body. They
	// are wired anyway: what makes them live is an edit to api/openapi.yaml,
	// and nothing there would tell whoever makes it that a leak switched on.
	//
	// Each keeps the status its own path means rather than collapsing to 500.
	// A request the glue could not bind is the client's mistake and stays a
	// 400. OapiRequestValidator above catches most of those first, but not all
	// — it does not enforce every string format the spec can declare, so a
	// malformed value can still reach the binder. Answering 500 there would
	// blame the server for a client's typo, and page somebody for it.
	strict := apigen.NewStrictHandlerWithOptions(si, nil, strictOptions(logger))
	apigen.RegisterHandlersWithOptions(engine, strict, serverOptions(logger))

	return engine, nil
}

// strictOptions and serverOptions exist as named functions rather than as
// literals inline above so that a test can assert every hook is set. Wiring
// them inline would leave the one thing worth checking — that none was
// forgotten — visible only by reading.
func strictOptions(logger log.Logger) apigen.StrictGinServerOptions {
	onServerError := boundaryErrorFunc(logger, http.StatusInternalServerError)

	return apigen.StrictGinServerOptions{
		RequestErrorHandlerFunc:  boundaryErrorFunc(logger, http.StatusBadRequest),
		HandlerErrorFunc:         onServerError,
		ResponseErrorHandlerFunc: onServerError,
	}
}

// serverOptions carries the fourth generated error path, the one the
// parameter-binding wrappers call.
func serverOptions(logger log.Logger) apigen.GinServerOptions {
	return apigen.GinServerOptions{ErrorHandler: bindErrorFunc(logger)}
}

// boundaryErrorFunc builds the replacement for an error hook that carries no
// status of its own, so the caller names the one that path means.
//
// Each generated default writes gin.H{"msg": err.Error()} straight into the
// response body — exactly the leak interfaces/doc.go's MUST forbids.
//
// This is the backstop, not the normal path. Every handler still returns a
// typed response and a nil error, so the contract stays compiler-checked; this
// only fires when one forgets, or when the generated glue fails somewhere no
// handler can see.
func boundaryErrorFunc(logger log.Logger, status int) func(*gin.Context, error) {
	return func(ctx *gin.Context, err error) {
		writeGenericProblem(ctx, logger, err, status)
	}
}

// bindErrorFunc is the same replacement for GinServerOptions.ErrorHandler, the
// one generated error path that is handed the status it should answer with.
//
// Taking that status rather than choosing one keeps the generated wrappers'
// judgement about what a binding failure means, while replacing only their
// judgement about what to put in the body.
func bindErrorFunc(logger log.Logger) func(*gin.Context, error, int) {
	return func(ctx *gin.Context, err error, status int) {
		writeGenericProblem(ctx, logger, err, status)
	}
}

// writeGenericProblem is the body every replaced hook shares: log the error
// where it still exists, then answer with a document that explains nothing.
//
// The two halves are one decision. The body says nothing on purpose, which is
// what makes the record not optional — it is the only account of a failure that
// reached the boundary uncategorised, and this is the last place the error is
// still in hand.
func writeGenericProblem(ctx *gin.Context, logger log.Logger, err error, status int) {
	logger.Error("unhandled error at the http boundary", "error", err, "status", status)

	body, marshalErr := json.Marshal(errmap.ToGenericProblem(err, status))
	if marshalErr != nil {
		ctx.Status(status)
		return
	}
	ctx.Data(status, "application/problem+json", body)
}
