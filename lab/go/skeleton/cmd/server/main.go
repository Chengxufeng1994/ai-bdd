// Command server serves the skeleton's HTTP API.
//
// This is the entry point, not the composition root — see bootstrap, which it
// calls, for that. Keep it to wiring — any decision made here is a decision
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

	"skeleton/internal/bootstrap"
	"skeleton/pkg/config"
	"skeleton/pkg/i18n"
	"skeleton/pkg/log"
	"skeleton/pkg/version"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run assembles and serves, returning the first failure instead of exiting.
//
// Keeping main to three lines is what makes everything below it reachable from
// a test: an os.Exit inside the wiring would take the test process with it.
func run() error {
	// Configuration is validated before the logger exists, so a bad value has
	// to be reported on stderr rather than through slog. That is the correct
	// order: a logger built from invalid settings is not trustworthy.
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := log.New(
		log.WithLevel(cfg.Log.Level),
		log.WithFormat(bootstrap.LogFormat(cfg.Log.Format)),
	)

	// Read once here, in the only place allowed to touch the package variable.
	v := version.Build()

	router, err := bootstrap.NewHandler(bootstrap.Deps{
		Version:    v,
		Logger:     logger,
		Translator: i18n.NewBundle(messages),
	})
	if err != nil {
		return fmt.Errorf("assemble the application: %w", err)
	}

	logger.Info("listening", "addr", cfg.Addr, "version", v)

	return http.ListenAndServe(cfg.Addr, router)
}

// messages is the message table this binary ships with.
//
// It lives here rather than in pkg/i18n because which locales a deployment
// carries is a property of the deployment, not of the translator.
var messages = map[string]map[string]string{
	"en": {
		"version.unavailable": "The service version is unavailable.",
	},
}
