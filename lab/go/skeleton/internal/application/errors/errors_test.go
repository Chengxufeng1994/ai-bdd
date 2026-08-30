package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "skeleton/internal/application/errors"
)

// TestErrorsAsRecoversThroughWrapping is what makes the whole design work: a
// use case wraps an apperrors.Error with %w on its way up, possibly more than
// once, and an adapter must still find it with errors.As.
func TestErrorsAsRecoversThroughWrapping(t *testing.T) {
	original := apperrors.Error{Kind: apperrors.KindNotFound, Err: errors.New("workout 42 does not exist")}
	wrapped := fmt.Errorf("get workout: %w", original)

	var got apperrors.Error
	if !errors.As(wrapped, &got) {
		t.Fatalf("errors.As did not recover an apperrors.Error from %v", wrapped)
	}
	if got.Kind != apperrors.KindNotFound {
		t.Errorf("Kind: want KindNotFound, got %v", got.Kind)
	}
}

// TestKindZeroValueIsUnclassified guards the default-deny property the
// package doc describes: an Error nobody set Kind on must not be mistaken for
// any specific classification.
func TestKindZeroValueIsUnclassified(t *testing.T) {
	var got apperrors.Kind

	if got != apperrors.KindUnclassified {
		t.Errorf("zero value: want KindUnclassified, got %v", got)
	}
}

// TestUnwrapReturnsTheWrappedError checks Unwrap directly, rather than only
// through errors.As, so a future change that breaks Unwrap while leaving
// errors.As passing (for instance by matching on Error() text) still fails.
func TestUnwrapReturnsTheWrappedError(t *testing.T) {
	want := errors.New("build stamp unreadable")
	e := apperrors.Error{Kind: apperrors.KindUnavailable, Err: want}

	got := e.Unwrap()
	if !errors.Is(got, want) {
		t.Errorf("Unwrap: want %v, got %v", want, got)
	}
}

// TestWithParamsCopiesTheCallersMap guards the aliasing risk that actually
// matters: an Error is wrapped, returned and logged well after the call site
// that built it has moved on, so WithParams must not keep the caller's map by
// reference — a later mutation there must not silently change what a past
// failure says it was.
func TestWithParamsCopiesTheCallersMap(t *testing.T) {
	params := map[string]any{"id": "01H8X"}
	base := apperrors.Error{Kind: apperrors.KindUnavailable}

	err := base.WithParams(params)
	params["id"] = "mutated"

	if err.Params["id"] != "01H8X" {
		t.Errorf("Params: want the value at the time of the call, got %v", err.Params["id"])
	}
}

// TestWithDetailsSetsTheField checks that WithDetails attaches the operator
// context it is given.
func TestWithDetailsSetsTheField(t *testing.T) {
	base := apperrors.Error{Kind: apperrors.KindUnavailable}
	got := base.WithDetails("3 retries over 5s")

	if got.Details != "3 retries over 5s" {
		t.Errorf("Details: got %q", got.Details)
	}
}
