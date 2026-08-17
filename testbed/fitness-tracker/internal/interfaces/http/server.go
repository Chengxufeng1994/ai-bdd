package http

import (
	"fitness-tracker/internal/interfaces/http/apigen"
)

// Server implements the generated StrictServerInterface.
//
// The strict interface is deliberate: its methods take a context.Context and a
// typed request, and return a typed response — no *gin.Context anywhere. That
// keeps the web framework at the edge, in router.go and the generated package,
// where swapping it would touch two files and nothing deeper.
type Server struct {
	version string
}

// Compile-time check that Server still satisfies the contract. Without it, a
// change to api/openapi.yaml that adds an operation would fail somewhere in
// router.go with a less obvious message.
var _ apigen.StrictServerInterface = (*Server)(nil)

// NewServer takes the version rather than reading it from a package variable so
// that a test can run two servers at different versions in the same process —
// which the acceptance suite does, since scenarios run concurrently.
func NewServer(version string) *Server {
	return &Server{version: version}
}
