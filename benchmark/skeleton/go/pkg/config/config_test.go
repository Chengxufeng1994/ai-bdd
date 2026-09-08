package config_test

import (
	"log/slog"
	"strings"
	"testing"

	"skeleton/pkg/config"
)

// These are the skeleton's unit tests: no I/O, no container, no HTTP. They run
// under `make test` alongside the acceptance suite.

func TestLoadUsesDefaultsWhenNothingIsSet(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}

	if cfg.Addr != ":8080" {
		t.Errorf("Addr: want :8080, got %q", cfg.Addr)
	}
	if cfg.Log.Level != slog.LevelInfo {
		t.Errorf("Log.Level: want info, got %v", cfg.Log.Level)
	}
	if cfg.Log.Format != config.LogFormatText {
		t.Errorf("Log.Format: want text, got %q", cfg.Log.Format)
	}
}

func TestLoadReadsEnvironment(t *testing.T) {
	t.Setenv("APP_ADDR", "127.0.0.1:9000")
	t.Setenv("APP_LOG_LEVEL", "debug")
	t.Setenv("APP_LOG_FORMAT", "json")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}

	if cfg.Addr != "127.0.0.1:9000" {
		t.Errorf("Addr: want 127.0.0.1:9000, got %q", cfg.Addr)
	}
	if cfg.Log.Level != slog.LevelDebug {
		t.Errorf("Log.Level: want debug, got %v", cfg.Log.Level)
	}
	if cfg.Log.Format != config.LogFormatJSON {
		t.Errorf("Log.Format: want json, got %q", cfg.Log.Format)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	tests := map[string]struct {
		env  map[string]string
		want string
	}{
		"address without a port": {
			env:  map[string]string{"APP_ADDR": "localhost"},
			want: "APP_ADDR",
		},
		"log level that looks like a typo": {
			env:  map[string]string{"APP_LOG_LEVEL": "debgu"},
			want: "APP_LOG_LEVEL",
		},
		"unsupported log format": {
			env:  map[string]string{"APP_LOG_FORMAT": "xml"},
			want: "APP_LOG_FORMAT",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			_, err := config.Load()
			if err == nil {
				t.Fatal("want an error, got nil")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error should name %s, got: %v", tc.want, err)
			}
		})
	}
}

// A mistyped variable is cheap to fix once you can see all of them. Reporting
// one problem per restart is what this test exists to prevent regressing.
func TestLoadReportsEveryProblemAtOnce(t *testing.T) {
	t.Setenv("APP_ADDR", "localhost")
	t.Setenv("APP_LOG_LEVEL", "debgu")
	t.Setenv("APP_LOG_FORMAT", "xml")

	_, err := config.Load()
	if err == nil {
		t.Fatal("want an error, got nil")
	}

	for _, want := range []string{"APP_ADDR", "APP_LOG_LEVEL", "APP_LOG_FORMAT"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error should mention %s, got: %v", want, err)
		}
	}
}
