package bootstrap_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"skeleton/internal/bootstrap"
	"skeleton/internal/interfaces/http/apigen"
	"skeleton/pkg/i18n"
	"skeleton/pkg/log"
)

// features/version.feature runs one scenario, so one version, at a time — it
// cannot catch NewHandler reading a future package-level variable written by
// pkg/version instead of the version parameter it is handed. This builds two
// handlers at different versions from concurrent goroutines and checks each
// still answers its own, under -race, which is what would actually catch it.
func TestNewHandlerIsIndependentAcrossConcurrentVersions(t *testing.T) {
	versions := []string{"1.0.0", "2.0.0"}

	var wg sync.WaitGroup
	for _, version := range versions {
		wg.Add(1)
		go func(version string) {
			defer wg.Done()

			router, err := bootstrap.NewHandler(bootstrap.Deps{Version: version, Logger: log.Discard(), Translator: i18n.NewBundle(nil)})
			if err != nil {
				t.Errorf("NewHandler(%q): %v", version, err)
				return
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))

			var body apigen.Version
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Errorf("decode response for %q: %v", version, err)
				return
			}
			if body.Version != version {
				t.Errorf("Version: want %q, got %q", version, body.Version)
			}
		}(version)
	}
	wg.Wait()
}

// A Deps missing its Logger or its Translator assembles without complaint and
// serves the happy path correctly. The gap only opens on the failure path,
// where GetVersion logs the error and errmap translates it — so the first
// request that fails nil-panics, gin.Recovery() catches it, and the client gets
// a bodiless 500 while the process survives.
//
// That is worse than a crash, because of what it replaces: the structured
// record naming the operation and its cause — the thing that branch exists to
// produce — becomes an unstructured panic trace, at the one moment anybody
// needed it. And nothing before that moment gives any sign: assembly returned
// a nil error, the health of the service looked fine, tests of the 200 path
// passed.
//
// Incomplete wiring is a startup failure. NewHandler already returns an error,
// which is the whole reason it can say so.
func TestNewHandlerRejectsDepsItCannotServeAFailureWith(t *testing.T) {
	tests := map[string]bootstrap.Deps{
		"no logger":     {Version: "1.0.0", Translator: i18n.NewBundle(nil)},
		"no translator": {Version: "1.0.0", Logger: log.Discard()},
	}

	for name, deps := range tests {
		t.Run(name, func(t *testing.T) {
			handler, err := bootstrap.NewHandler(deps)
			if err == nil {
				t.Fatal("want assembly to fail, got a handler that would panic on its first failing request")
			}
			if handler != nil {
				t.Errorf("want no handler alongside the error, got %T", handler)
			}
		})
	}
}
