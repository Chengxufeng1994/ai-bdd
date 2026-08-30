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

	"skeleton/internal/application/service"
	"skeleton/internal/application/usecase/query"
	"skeleton/internal/infrastructure/buildinfo"
	apihttp "skeleton/internal/interfaces/http"
	"skeleton/pkg/config"
	"skeleton/pkg/i18n"
	"skeleton/pkg/log"
)

// Deps are the constructed objects the application is assembled from.
//
// It is deliberately not pkg/config.Config: that holds settings read from the
// environment, this holds objects already built from them. Merging the two
// would make "was this read or constructed" unanswerable at a glance, and the
// answer decides who is allowed to change a field.
type Deps struct {
	// Version is the build version this instance reports. It is a parameter
	// and must stay one: the acceptance suite runs scenarios concurrently at
	// different versions in one process, so a package-level read here would
	// make them pollute each other intermittently.
	Version string

	// Logger receives what never reaches a client — an error's Where, Details
	// and cause.
	Logger log.Logger

	// Translator renders message keys for the locale a request asked for.
	Translator i18n.Translator
}

// NewHandler assembles the HTTP surface from deps.
//
// It returns http.Handler rather than *gin.Engine so that router.go remains the
// only hand-written file naming gin.
func NewHandler(deps Deps) (http.Handler, error) {
	svc := service.NewVersionService(query.NewGetVersion(buildinfo.NewProvider(deps.Version)))

	engine, err := apihttp.NewRouter(apihttp.NewServer(svc, deps.Logger, deps.Translator))
	if err != nil {
		return nil, err
	}
	return engine, nil
}

// LogFormat maps a validated configuration value onto a logger option.
//
// It lives here rather than in either package because the composition root is
// the one place allowed to know both. Putting it in config would make config
// import log; putting it in log would make log import config. Either way the
// two stop being replaceable on their own. It is here rather than in main
// because this package is testable and main is not.
func LogFormat(f config.LogFormat) log.Format {
	if f == config.LogFormatJSON {
		return log.FormatJSON
	}
	return log.FormatText
}
