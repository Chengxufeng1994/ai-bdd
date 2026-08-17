package http

import (
	"github.com/gin-gonic/gin"

	"fitness-tracker/internal/interfaces/http/apigen"
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
func NewRouter(si apigen.StrictServerInterface) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	apigen.RegisterHandlers(engine, apigen.NewStrictHandler(si, nil))

	return engine
}
