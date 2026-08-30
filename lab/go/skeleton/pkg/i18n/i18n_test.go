package i18n_test

import (
	"strings"
	"testing"

	"skeleton/pkg/i18n"
)

func bundle() i18n.Bundle {
	return i18n.NewBundle(map[string]map[string]string{
		"en": {
			"version.unavailable": "The service version is unavailable.",
			"workout.not_found":   "Workout {id} does not exist.",
		},
		"zh-TW": {
			"version.unavailable": "無法取得服務版本。",
		},
	})
}

func TestTranslateRendersTheRequestedLocale(t *testing.T) {
	got := bundle().Translate("zh-TW", "version.unavailable", nil)

	if got != "無法取得服務版本。" {
		t.Errorf("want the zh-TW message, got %q", got)
	}
}

func TestTranslateInterpolatesParams(t *testing.T) {
	got := bundle().Translate("en", "workout.not_found", map[string]any{"id": "01H8X"})

	if got != "Workout 01H8X does not exist." {
		t.Errorf("want the id interpolated, got %q", got)
	}
}

// A missing translation renders as the key itself rather than falling back to a
// default locale. A fallback would look like a finished message and hide the
// gap; "workout.not_found" on screen never does.
func TestMissingKeyRendersAsTheKey(t *testing.T) {
	for name, tc := range map[string]struct{ locale, key string }{
		"unknown key":    {"en", "nope.missing"},
		"unknown locale": {"fr", "version.unavailable"},
	} {
		t.Run(name, func(t *testing.T) {
			if got := bundle().Translate(tc.locale, tc.key, nil); got != tc.key {
				t.Errorf("want the key %q returned, got %q", tc.key, got)
			}
		})
	}
}

// A template needing a param nobody supplied renders as the key too. Emitting
// "Workout {id} does not exist." would look like a bug in the product rather
// than a gap in the call site, which is the harder thing to diagnose.
func TestMissingParamRendersAsTheKey(t *testing.T) {
	got := bundle().Translate("en", "workout.not_found", nil)

	if got != "workout.not_found" {
		t.Errorf("want the key when a param is missing, got %q", got)
	}
	if strings.Contains(got, "{") {
		t.Errorf("an unresolved placeholder escaped: %q", got)
	}
}
