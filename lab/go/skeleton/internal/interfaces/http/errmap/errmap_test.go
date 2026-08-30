package errmap_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	apperrors "skeleton/internal/application/errors"
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

// TestStatusFor covers every kind this table classifies, plus the three ways
// StatusFor must fall back to 500: an apperrors.Error explicitly carrying the
// zero-value KindUnclassified, a plain error wrapped with %w that never
// became an apperrors.Error, and an unrecognised error guaranteed non-nil.
// The default-deny order this table enforces — recognise what we classify,
// 500 for everything else — is what stops an unrecognised error's message
// from ever reaching a client; see StatusFor's doc comment.
func TestStatusFor(t *testing.T) {
	tests := map[string]struct {
		err  error
		want int
	}{
		"not found": {
			err:  apperrors.Error{Kind: apperrors.KindNotFound, Err: errors.New("workout 42 does not exist")},
			want: http.StatusNotFound,
		},
		"invalid": {
			err:  apperrors.Error{Kind: apperrors.KindInvalid, Err: errors.New("weight must be >= 0")},
			want: http.StatusUnprocessableEntity,
		},
		"conflict": {
			err:  apperrors.Error{Kind: apperrors.KindConflict, Err: errors.New("workout already recorded")},
			want: http.StatusConflict,
		},
		"unauthorized": {
			err:  apperrors.Error{Kind: apperrors.KindUnauthorized, Err: errors.New("no credentials supplied")},
			want: http.StatusUnauthorized,
		},
		"forbidden": {
			err:  apperrors.Error{Kind: apperrors.KindForbidden, Err: errors.New("not allowed to delete this workout")},
			want: http.StatusForbidden,
		},
		"unavailable": {
			err:  apperrors.Error{Kind: apperrors.KindUnavailable, Err: errors.New("build stamp unreadable")},
			want: http.StatusServiceUnavailable,
		},
		"unclassified apperrors.Error falls back to 500": {
			err:  apperrors.Error{Err: errors.New("nobody set Kind")},
			want: http.StatusInternalServerError,
		},
		"plain error wrapped with %w, never classified, falls back to 500": {
			err:  fmt.Errorf("read build version: %w", errors.New("disk unavailable")),
			want: http.StatusInternalServerError,
		},
		"unrecognised non-nil error falls back to 500": {
			err:  errors.New("anything"),
			want: http.StatusInternalServerError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := errmap.StatusFor(tc.err)
			if got != tc.want {
				t.Errorf("StatusFor: want %d, got %d", tc.want, got)
			}
		})
	}
}
