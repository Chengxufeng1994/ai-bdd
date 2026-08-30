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
