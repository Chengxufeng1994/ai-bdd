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

// A parameter value that itself contains a brace must not be mistaken for an
// unresolved placeholder. Scanning the rendered result for "{" would discard a
// fully-supplied translation and wrongly return the key.
func TestParamValueContainingBraceRendersCorrectly(t *testing.T) {
	b := i18n.NewBundle(map[string]map[string]string{
		"en": {"greet": "Hello {name}!"},
	})

	got := b.Translate("en", "greet", map[string]any{"name": "{bob}"})

	if got != "Hello {bob}!" {
		t.Errorf("want the brace in the value preserved, got %q", got)
	}
}

// Two params where one's value looks like the other's placeholder must render
// deterministically regardless of map iteration order. Substituting into the
// rendered result rather than the original template makes the outcome depend
// on which param is applied first, so this runs several times in one test.
func TestParamValueResemblingAnotherPlaceholderIsDeterministic(t *testing.T) {
	b := i18n.NewBundle(map[string]map[string]string{
		"en": {"pair": "{a} and {b}"},
	})

	for i := 0; i < 20; i++ {
		got := b.Translate("en", "pair", map[string]any{"a": "{b}", "b": "x"})

		if got != "{b} and x" {
			t.Errorf("run %d: want deterministic rendering, got %q", i, got)
		}
	}
}

// An unterminated brace is a malformed template, which is a gap in the message
// table rather than at the call site — but it is still a gap, and this
// package's answer to every gap is the key.
func TestUnterminatedBraceRendersAsTheKey(t *testing.T) {
	b := i18n.NewBundle(map[string]map[string]string{
		"en": {"broken": "Hello {name"},
	})

	got := b.Translate("en", "broken", map[string]any{"name": "world"})

	if got != "broken" {
		t.Errorf("want the key for an unterminated brace, got %q", got)
	}
}
