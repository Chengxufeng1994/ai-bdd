// Command server serves the skeleton's HTTP API.
//
// This file is the composition root: it is the only place allowed to know about
// every layer at once. Keep it to wiring — any decision made here is a decision
// that cannot be tested without starting a process.
//
// Dependencies are wired by hand rather than with a container. The graph is
// small enough to read top to bottom, and a reader can answer "what does this
// depend on" without learning a second tool.
package main

import (
	"fmt"
	"net/http"
	"os"

	apihttp "skeleton/internal/interfaces/http"
	"skeleton/pkg/config"
	"skeleton/pkg/log"
)

// version is overridden at build time:
//
//	go build -ldflags "-X main.version=$(git describe --tags)" ./cmd/server
//
// The default is deliberately not a plausible version number: seeing "dev" in a
// deployed environment should look wrong immediately.
var version = "dev"

func main() {
	// Configuration is validated before the logger exists, so a bad value has
	// to be reported on stderr rather than through slog. That is the correct
	// order: a logger built from invalid settings is not trustworthy.
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Mapping configuration onto logger options happens here, in the only place
	// that is allowed to know about both packages.
	logger := log.New(
		log.WithLevel(cfg.Log.Level),
		log.WithFormat(logFormat(cfg.Log.Format)),
	)

	router, err := apihttp.NewRouter(apihttp.NewServer(version))
	if err != nil {
		logger.Error("build router", "error", err)
		os.Exit(1)
	}

	logger.Info("listening", "addr", cfg.Addr, "version", version)

	if err := http.ListenAndServe(cfg.Addr, router); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

// logFormat maps a validated configuration value onto a logger option.
//
// This mapping lives here rather than in either package because it is the one
// place allowed to know both. Putting it in config would make config import log;
// putting it in log would make log import config. Either way the two stop being
// replaceable on their own, which is the property this whole arrangement exists
// to protect.
func logFormat(f config.LogFormat) log.Format {
	if f == config.LogFormatJSON {
		return log.FormatJSON
	}
	return log.FormatText
}
