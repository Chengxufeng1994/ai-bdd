// Package bootstrap assembles the application.
//
// It is the composition root: the only place that knows every layer at once.
// cmd/server and the acceptance suite both call it, so "how this system is
// wired" has exactly one description — and the chain a scenario exercises is the
// chain the binary runs. Two separate wirings would drift, and the drift would
// be silent: the suite could stay green while the binary was assembled wrong.
//
// Keep this to wiring. Any decision made here is a decision that cannot be
// tested without assembling the whole application.
package bootstrap

import (
	"net/http"

	"skeleton/internal/application/port/in"
	"skeleton/internal/application/usecase"
	"skeleton/internal/infrastructure/buildinfo"
	apihttp "skeleton/internal/interfaces/http"
)

// NewHandler assembles the HTTP surface reporting the given version.
//
// The version is a parameter and must stay one. Reaching for pkg/version from
// inside would make every assembled handler share one value, and the acceptance
// suite runs scenarios concurrently at different versions — they would pollute
// each other intermittently, which is the hardest kind of failure to trace.
//
// It returns http.Handler rather than *gin.Engine so that router.go remains the
// only hand-written file naming gin.
func NewHandler(version string) (http.Handler, error) {
	svc := in.VersionService{
		GetVersion: usecase.NewGetVersion(buildinfo.NewProvider(version)),
	}

	engine, err := apihttp.NewRouter(apihttp.NewServer(svc))
	if err != nil {
		return nil, err
	}
	return engine, nil
}
