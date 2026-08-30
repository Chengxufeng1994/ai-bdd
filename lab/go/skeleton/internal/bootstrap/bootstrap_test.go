package bootstrap_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"skeleton/internal/bootstrap"
	"skeleton/internal/interfaces/http/apigen"
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

			router, err := bootstrap.NewHandler(version)
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
