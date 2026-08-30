package errmap_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"skeleton/internal/interfaces/http/errmap"
)

func TestToInternalServerErrorReportsStatus500(t *testing.T) {
	got := errmap.ToInternalServerError(errors.New("anything"))

	if got.Status != http.StatusInternalServerError {
		t.Errorf("Status: want 500, got %d", got.Status)
	}
	if got.Title == "" {
		t.Error("Title: want a human-readable summary, got empty")
	}
	if got.Type == "" {
		t.Error("Type: want a stable identifier, got empty")
	}
}

// ARCHITECTURE.md §7 requires that a 500 body stay generic: stack traces, SQL
// fragments and internal hostnames belong in logs, not in a response that
// reaches whoever made the request. Until now that rule lived only in prose.
//
// The assertion is on the rendered JSON rather than on Detail alone, so that
// adding a new field cannot open a second leak without failing this test.
func TestToInternalServerErrorDoesNotLeakTheUnderlyingError(t *testing.T) {
	leaky := errors.New(`pq: password authentication failed for user "admin" at db-prod-3.internal`)

	got := errmap.ToInternalServerError(leaky)

	body, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal problem document: %v", err)
	}

	for _, secret := range []string{"password", "admin", "db-prod-3"} {
		if strings.Contains(string(body), secret) {
			t.Errorf("500 body leaks %q from the underlying error: %s", secret, body)
		}
	}
}
