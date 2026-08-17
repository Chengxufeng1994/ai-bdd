package acceptance_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"

	apihttp "fitness-tracker/internal/interfaces/http"
	"fitness-tracker/internal/interfaces/http/apigen"
)

// registerSteps binds Gherkin step text to Go functions.
//
// A step definition may optionally take context.Context as its first argument,
// and may return nothing, an error, a context.Context, or both. Take the context
// whenever the step touches scenario state — with Concurrency above 1 it is the
// only safe place to keep it.
//
// Use Given/When/Then rather than Step where the keyword is meaningful: it stops
// a Then-shaped assertion from being reused as a Given-shaped setup, which is
// how step definitions quietly turn into a second, undocumented API.
//
// Before adding a step, search the existing ones for something equivalent.
// Duplicate steps that differ only in wording are the start of an unmaintainable
// glue layer. To see what is currently registered:
//
//	go test ./test/acceptance/ -godog.definitions
func registerSteps(sc *godog.ScenarioContext) {
	sc.Given(`^the service is running at version "([^"]*)"$`, theServiceIsRunningAtVersion)
	sc.When(`^a client asks for the service version$`, aClientAsksForTheServiceVersion)
	sc.Then(`^the service reports version "([^"]*)"$`, theServiceReportsVersion)
}

func theServiceIsRunningAtVersion(ctx context.Context, version string) error {
	w, err := worldFrom(ctx)
	if err != nil {
		return err
	}

	w.router = apihttp.NewRouter(apihttp.NewServer(version))

	return nil
}

// The request is issued against the router in-process via httptest rather than
// over a real socket. That keeps scenarios fast and free of port allocation,
// while still exercising the generated routing, the strict handler wrapper and
// JSON encoding — everything except the network itself.
func aClientAsksForTheServiceVersion(ctx context.Context) error {
	w, err := worldFrom(ctx)
	if err != nil {
		return err
	}
	if w.router == nil {
		return fmt.Errorf("no service is running: a Given step must start one first")
	}

	w.resp = httptest.NewRecorder()
	w.router.ServeHTTP(w.resp, httptest.NewRequest(http.MethodGet, "/version", nil))

	return nil
}

func theServiceReportsVersion(ctx context.Context, expected string) error {
	w, err := worldFrom(ctx)
	if err != nil {
		return err
	}
	if w.resp == nil {
		return fmt.Errorf("no response recorded: a When step must make a request first")
	}

	if w.resp.Code != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d: %s", w.resp.Code, w.resp.Body.String())
	}

	// Decoding into the generated type rather than a map is deliberate: if the
	// contract renames a field, this stops compiling instead of silently
	// asserting against a zero value.
	var body apigen.Version
	if err := json.Unmarshal(w.resp.Body.Bytes(), &body); err != nil {
		return fmt.Errorf("response body is not a Version document: %w (body: %s)", err, w.resp.Body.String())
	}

	if body.Version != expected {
		return fmt.Errorf("expected version %q, got %q", expected, body.Version)
	}

	return nil
}
