package http

import (
	"skeleton/internal/application/port/in"
	"skeleton/internal/interfaces/http/apigen"
)

// Server implements the generated StrictServerInterface.
//
// The strict interface is deliberate: its methods take a context.Context and a
// typed request, and return a typed response — no *gin.Context anywhere. That
// keeps the web framework at the edge, in router.go and the generated package,
// where swapping it would touch two files and nothing deeper.
type Server struct {
	svc in.VersionService
}

// Compile-time check that Server still satisfies the contract. Without it, a
// change to api/openapi.yaml that adds an operation would fail somewhere in
// router.go with a less obvious message.
var _ apigen.StrictServerInterface = (*Server)(nil)

// NewServer takes the driving port rather than reading anything from a
// package variable, so that a test can run two servers configured differently
// in the same process — which the acceptance suite does, since scenarios run
// concurrently.
//
// It depends on the port rather than on a concrete service: in.VersionService
// names only capabilities, so holding it grants access to those capabilities
// and nothing more. A concrete service would bypass that contract and lose the
// ability to substitute a decorated implementation.
func NewServer(svc in.VersionService) *Server {
	return &Server{svc: svc}
}
