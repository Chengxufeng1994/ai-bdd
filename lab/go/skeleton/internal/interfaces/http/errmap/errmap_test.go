package errmap_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	apperrors "skeleton/internal/application/errors"
	"skeleton/internal/interfaces/http/apigen"
	"skeleton/internal/interfaces/http/errmap"
	"skeleton/pkg/i18n"
)

// assertGenericProblem checks that problem is the body this adapter gives a
// failure it refuses to describe, identity included.
//
// The status alone proves nothing here, and that is the whole reason this
// helper exists. Both of ToProblem's branches can produce a 500 — the generic
// one always does, and the classified one does whenever a kind maps there — so
// a test asserting only the number passes no matter which branch ran, and
// cannot fail when the guard that separates them is deleted. What actually
// distinguishes them is what the client can identify the failure as: the
// generic branch answers with RFC 9457's reserved "about:blank" and a fixed
// summary, the classified branch with a type derived from the message key and
// a translated title.
//
// The literals are spelled out rather than read from errmap's constants on
// purpose. They are the wire contract a client branches on; a test that read
// the same constants the code does would follow them wherever they went.
func assertGenericProblem(t *testing.T, problem apigen.Problem) {
	t.Helper()

	if problem.Type != "about:blank" {
		t.Errorf("Type: want RFC 9457's reserved about:blank for a failure with no specific identity, got %q", problem.Type)
	}
	if problem.Title != "Internal Server Error" {
		t.Errorf("Title: want the generic summary, got %q", problem.Title)
	}
	if problem.Status != http.StatusInternalServerError {
		t.Errorf("body status: want 500, got %d", problem.Status)
	}
}

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

// The two channels must not cross. Params is rendered into the body; Where,
// DetailedError and the wrapped cause are for logs. This asserts on the
// marshalled document rather than on one field, so a future field cannot open
// a second leak silently — the same shape as
// TestToInternalServerErrorDoesNotLeakTheUnderlyingError above.
func TestToProblemRendersParamsButNeverTheLogChannel(t *testing.T) {
	tr := i18n.NewBundle(map[string]map[string]string{
		"en": {"version.unavailable": "Version {id} is unavailable."},
	})

	appErr := apperrors.Error{
		Kind:          apperrors.KindUnavailable,
		MessageKey:    "version.unavailable",
		Params:        map[string]any{"id": "01H8X"},
		Where:         "usecase.GetVersion",
		DetailedError: "3 retries over 5s",
	}.Wrap(errors.New(`pq: password authentication failed for user "admin"`))

	status, problem := errmap.ToProblem(appErr, "en", tr)

	if status != http.StatusServiceUnavailable {
		t.Errorf("status: want 503, got %d", status)
	}
	if problem.Title != "Version 01H8X is unavailable." {
		t.Errorf("Title: want the rendered message, got %q", problem.Title)
	}

	body, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("marshal problem: %v", err)
	}
	for _, secret := range []string{"usecase.GetVersion", "3 retries", "password", "admin"} {
		if strings.Contains(string(body), secret) {
			t.Errorf("the body leaks %q from the log channel: %s", secret, body)
		}
	}
}

// An unrecognised error takes the same path an unclassified one does: 500
// with a generic body. Recognise-then-default is a security property, so it
// is asserted rather than assumed — and asserted on the body's identity, not
// only its status, so that the test can tell "took the generic branch" from
// "took the classified branch and landed on 500 anyway".
func TestToProblemDefaultsUnrecognisedErrorsTo500(t *testing.T) {
	status, problem := errmap.ToProblem(errors.New("something nobody classified"), "en", i18n.NewBundle(nil))

	if status != http.StatusInternalServerError {
		t.Errorf("status: want 500, got %d", status)
	}
	assertGenericProblem(t, problem)

	body, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("marshal problem: %v", err)
	}
	if strings.Contains(string(body), "something nobody classified") {
		t.Errorf("the body leaks the unrecognised error: %s", body)
	}
}

// An apperrors.Error whose Kind was never set — the zero value,
// KindUnclassified — takes the same default path an unrecognised error does,
// rather than being treated as some accidental classification.
//
// The identity assertion is what gives this test teeth. Delete ToProblem's
// KindUnclassified guard and the error falls through to StatusFor, whose
// default is also 500 — so a status-only assertion stays green while the
// client starts receiving a type and a title derived from a message key
// nobody classified, which is precisely the accidental classification this
// test claims to forbid.
func TestToProblemDefaultsZeroKindTo500(t *testing.T) {
	appErr := apperrors.Error{
		MessageKey: "whatever",
		Where:      "usecase.GetVersion",
	}.Wrap(errors.New("nobody set Kind"))

	status, problem := errmap.ToProblem(appErr, "en", i18n.NewBundle(nil))

	if status != http.StatusInternalServerError {
		t.Errorf("status: want 500, got %d", status)
	}
	assertGenericProblem(t, problem)

	body, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("marshal problem: %v", err)
	}
	if strings.Contains(string(body), "nobody set Kind") {
		t.Errorf("the body leaks the underlying error: %s", body)
	}
}

// A missing translation renders the key. The key is also what the type is
// derived from, so both stay consistent when a locale is incomplete.
func TestToProblemRendersTheKeyWhenUntranslated(t *testing.T) {
	appErr := apperrors.Error{
		Kind:       apperrors.KindUnavailable,
		MessageKey: "version.unavailable",
	}.Wrap(errors.New("boom"))

	_, problem := errmap.ToProblem(appErr, "fr", i18n.NewBundle(nil))

	if problem.Title != "version.unavailable" {
		t.Errorf("Title: want the key, got %q", problem.Title)
	}
}

// A Kind is not enough to describe a failure with. Classification picks the
// status; MessageKey is what supplies the other two fields a client actually
// reads — the RFC 9457 type it branches on and the title it shows. An Error
// carrying a Kind and no MessageKey therefore renders a document that satisfies
// the schema and says nothing: a type that is the base URL with a trailing
// slash and nothing after it, and a title that is the empty string, which the
// schema marks required.
//
// So it belongs on the same default path an unclassified error takes. Half a
// classification is not a classification; it is the same "somebody forgot"
// KindUnclassified already stands for, arriving through the other field.
func TestToProblemDefaultsAClassifiedErrorWithNoMessageKeyTo500(t *testing.T) {
	appErr := apperrors.Error{
		Kind:  apperrors.KindUnavailable,
		Where: "usecase.GetVersion",
	}.Wrap(errors.New("nobody set MessageKey"))

	status, problem := errmap.ToProblem(appErr, "en", i18n.NewBundle(nil))

	if status != http.StatusInternalServerError {
		t.Errorf("status: want 500, got %d", status)
	}
	assertGenericProblem(t, problem)
}
