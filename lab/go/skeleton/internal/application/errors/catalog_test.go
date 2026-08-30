package errors_test

import (
	"errors"
	"fmt"
	"testing"

	"strings"

	apperrors "skeleton/internal/application/errors"
)

// New 是 Code、MessageKey、Kind 唯一的設定路徑。這個測試存在是為了讓「個別賦值」
// 這條路成為編譯期就不存在的選項——若有人加了個別的 setter，這裡不會失敗，但
// 目錄的意義會消失，所以規則同時寫在 doc comment 裡。
func TestNewTakesAllThreeFromTheDescriptor(t *testing.T) {
	got := apperrors.New(apperrors.VersionUnavailable, "usecase.GetVersion", nil)

	if got.Code != apperrors.VersionUnavailable.Code {
		t.Errorf("Code: want %q, got %q", apperrors.VersionUnavailable.Code, got.Code)
	}
	if got.MessageKey != apperrors.VersionUnavailable.MessageKey {
		t.Errorf("MessageKey: want %q, got %q", apperrors.VersionUnavailable.MessageKey, got.MessageKey)
	}
	if got.Kind != apperrors.VersionUnavailable.Kind {
		t.Errorf("Kind: want %v, got %v", apperrors.VersionUnavailable.Kind, got.Kind)
	}
}

func TestWithParamsAndWithDetailsDoNotMutateTheOriginal(t *testing.T) {
	base := apperrors.New(apperrors.VersionUnavailable, "usecase.GetVersion", nil)

	derived := base.WithParams(map[string]any{"id": "01H8X"}).WithDetails("3 retries over 5s")

	if base.Params != nil {
		t.Errorf("Params: want the original untouched, got %v", base.Params)
	}
	if base.Details != "" {
		t.Errorf("Details: want the original untouched, got %q", base.Details)
	}
	if derived.Params["id"] != "01H8X" {
		t.Errorf("derived Params: want id 01H8X, got %v", derived.Params)
	}
	if derived.Details != "3 retries over 5s" {
		t.Errorf("derived Details: got %q", derived.Details)
	}
}

// Where 疊加是「在哪裡發生」這個需求的實作方式：每一層包裝加上自己的，解開就得到
// 一條路徑。這比 stack trace 便宜，且命名的是邏輯操作而非 Go 呼叫框。
func TestWhereComposesThroughWrapping(t *testing.T) {
	inner := apperrors.New(apperrors.VersionUnavailable, "usecase.GetVersion",
		fmt.Errorf("read build version: %w", errors.New("connection refused")))

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
