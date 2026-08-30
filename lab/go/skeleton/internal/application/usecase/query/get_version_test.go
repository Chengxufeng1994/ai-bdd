package query_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	apperrors "skeleton/internal/application/errors"
	"skeleton/internal/application/usecase/query"
)

// stubProvider stands in for out.VersionProvider.
//
// The port has one method, so a hand-written stub is enough — see
// internal/application/doc.go on keeping ports small enough to fake without a
// mocking library. The generated mock exists for ports that outgrow this.
type stubProvider struct {
	version string
	err     error
}

// Version satisfies out.VersionProvider, returning whatever the test set up.
//
// The comment is not decoration: .golangci.yml enables revive's `exported` rule
// with checkPrivateReceivers, so an exported method on a private receiver needs
// one — see pkg/log's slogLogger.With for the same pattern.
func (s stubProvider) Version(context.Context) (string, error) {
	return s.version, s.err
}

func TestGetVersionReportsWhatTheProviderReturns(t *testing.T) {
	h := query.NewGetVersion(stubProvider{version: "1.2.3"})

	got, err := h.Handle(context.Background(), query.GetVersion{})
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
	if got.Value != "1.2.3" {
		t.Errorf("Value: want 1.2.3, got %q", got.Value)
	}
}

// The provider is the only thing in this slice that can fail, so this is the
// only place the error branch can be reached at all. Without it the branch is
// written but never executed, which is indistinguishable from not writing it.
//
// It also gives application/errors' vocabulary its first producer: KindUnavailable
// is "precisely what a failing out.VersionProvider — or any other driven port —
// is" per errors.go's package doc, so this asserts the classification, not just
// the wrapping.
func TestGetVersionPropagatesProviderFailure(t *testing.T) {
	want := errors.New("build stamp unreadable")

	h := query.NewGetVersion(stubProvider{err: want})

	got, err := h.Handle(context.Background(), query.GetVersion{})
	if !errors.Is(err, want) {
		t.Fatalf("want the provider's error wrapped, got %v", err)
	}
	if got.Value != "" {
		t.Errorf("Value: want empty on failure, got %q", got.Value)
	}

	var classified apperrors.Error
	if !errors.As(err, &classified) {
		t.Fatalf("want a classified apperrors.Error, got %T: %v", err, err)
	}
	if classified.Kind != apperrors.KindUnavailable {
		t.Errorf("Kind: want KindUnavailable, got %v", classified.Kind)
	}
	if !strings.Contains(err.Error(), "read build version: ") {
		t.Errorf("want the %q context preserved, got %q", "read build version: ", err.Error())
	}
}
