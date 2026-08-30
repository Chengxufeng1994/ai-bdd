package errors_test

import (
	"errors"
	"fmt"
	"testing"

	"strings"

	apperrors "skeleton/internal/application/errors"
)

// TestCatalogEntrySpecialisedWithAtAndWrapping checks that a catalog prototype
// keeps its identity — Kind, Code, MessageKey — through specialisation, and
// that At and Wrapping attach what varies per occurrence.
func TestCatalogEntrySpecialisedWithAtAndWrapping(t *testing.T) {
	cause := errors.New("connection refused")
	got := apperrors.VersionUnavailable.At("usecase.GetVersion").Wrapping(cause)

	if got.Code != apperrors.VersionUnavailable.Code {
		t.Errorf("Code: want %q, got %q", apperrors.VersionUnavailable.Code, got.Code)
	}
	if got.MessageKey != apperrors.VersionUnavailable.MessageKey {
		t.Errorf("MessageKey: want %q, got %q", apperrors.VersionUnavailable.MessageKey, got.MessageKey)
	}
	if got.Kind != apperrors.VersionUnavailable.Kind {
		t.Errorf("Kind: want %v, got %v", apperrors.VersionUnavailable.Kind, got.Kind)
	}
	if got.Where != "usecase.GetVersion" {
		t.Errorf("Where: want %q, got %q", "usecase.GetVersion", got.Where)
	}
	if !errors.Is(got.Err, cause) {
		t.Errorf("Err: want %v, got %v", cause, got.Err)
	}
}

// TestWithParamsCopiesTheCallersMap guards the aliasing risk that actually
// matters: an Error is wrapped, returned and logged well after the call site
// that built it has moved on, so WithParams must not keep the caller's map by
// reference — a later mutation there must not silently change what a past
// failure says it was.
func TestWithParamsCopiesTheCallersMap(t *testing.T) {
	params := map[string]any{"id": "01H8X"}

	err := apperrors.VersionUnavailable.WithParams(params)
	params["id"] = "mutated"

	if err.Params["id"] != "01H8X" {
		t.Errorf("Params: want the value at the time of the call, got %v", err.Params["id"])
	}
}

func TestWithDetailsSetsTheField(t *testing.T) {
	got := apperrors.VersionUnavailable.WithDetails("3 retries over 5s")

	if got.Details != "3 retries over 5s" {
		t.Errorf("Details: got %q", got.Details)
	}
}

// Where 疊加是「在哪裡發生」這個需求的實作方式：每一層包裝加上自己的，解開就得到
// 一條路徑。這比 stack trace 便宜，且命名的是邏輯操作而非 Go 呼叫框。
func TestWhereComposesThroughWrapping(t *testing.T) {
	inner := apperrors.VersionUnavailable.At("usecase.GetVersion").
		Wrapping(fmt.Errorf("read build version: %w", errors.New("connection refused")))

	outer := fmt.Errorf("%w", apperrors.Error{
		Where: "service.VersionService.GetVersion",
		Err:   inner,
	})

	msg := outer.Error()
	for _, want := range []string{
		"service.VersionService.GetVersion",
		"usecase.GetVersion",
		"connection refused",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("want %q in the rendered chain, got %q", want, msg)
		}
	}
}
