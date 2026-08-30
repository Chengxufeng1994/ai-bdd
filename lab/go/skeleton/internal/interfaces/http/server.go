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
	uc in.VersionUseCase
}

// Compile-time check that Server still satisfies the contract. Without it, a
// change to api/openapi.yaml that adds an operation would fail somewhere in
// router.go with a less obvious message.
var _ apigen.StrictServerInterface = (*Server)(nil)

// NewServer takes the use case bundle rather than reading anything from a
// package variable, so that a test can run two servers configured differently
// in the same process — which the acceptance suite does, since scenarios run
// concurrently.
//
// It depends on the bundle rather than on a concrete handler: the bundle lives
// in port/in, so its fields are driving ports and holding it grants access to
// contracts and nothing more. A concrete handler would bypass those contracts
// and lose the ability to decorate.
func NewServer(uc in.VersionUseCase) *Server {
	return &Server{uc: uc}
}
