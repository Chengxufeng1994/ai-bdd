package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apihttp "skeleton/internal/interfaces/http"
	"skeleton/internal/interfaces/http/apigen"
)

// leakyServer implements apigen.StrictServerInterface by returning an error
// instead of a typed response — the mistake interfaces/doc.go's MUST forbids a
// handler from making. Every real handler in this codebase honours the rule,
// so this is the only way to drive NewRouter's configured HandlerErrorFunc.
type leakyServer struct{}

// GetVersion always fails, simulating a handler that returned its error
// instead of a typed 500 response.
func (leakyServer) GetVersion(context.Context, apigen.GetVersionRequestObject) (apigen.GetVersionResponseObject, error) {
	return nil, errors.New(`pq: password authentication failed for user "admin" at db-prod-3.internal`)
}

// ARCHITECTURE.md §7 requires a 500 body stay generic, and interfaces/doc.go's
// MUST assumes every handler cooperates by never returning an error. This
// covers the backstop for when one does not: NewRouter must not fall back to
// oapi-codegen's generated default, which would write err.Error() straight
// into the body (see apigen.NewStrictHandler's HandlerErrorFunc).
//
// The assertion scans the whole rendered body rather than one field, mirroring
// errmap_test.go, so a new field cannot open a second leak quietly.
func TestRouterHandlerErrorFuncDoesNotLeakTheUnderlyingError(t *testing.T) {
	engine, err := apihttp.NewRouter(leakyServer{})
	if err != nil {
		t.Fatalf("build router: %v", err)
	}

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("Code: want 500, got %d: %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Errorf("Content-Type: want application/problem+json, got %q", got)
	}

	body := rec.Body.String()
	for _, secret := range []string{"password", "admin", "db-prod-3"} {
		if strings.Contains(body, secret) {
			t.Errorf("500 body leaks %q from the underlying error: %s", secret, body)
		}
	}
}
