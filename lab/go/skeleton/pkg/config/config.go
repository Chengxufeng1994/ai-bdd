// Package config reads the application's settings from the environment.
//
// It is the only place that touches os.Getenv. Any package that reads the
// environment itself cannot be configured differently in a test without setting
// global state, and its defaults become invisible — spread across the codebase
// instead of listed in one struct.
//
// Environment variables rather than a config library: the settings fit on one
// screen, and a loader supporting files, flags and remote providers would be
// more code than the thing it configures. When that stops being true, this file
// is the only one that changes.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Config is the fully resolved configuration.
//
// It holds standard-library and primitive types only, and imports nothing else
// from this module. A config package that referenced pkg/log would make the two
// impossible to change independently; mapping these values onto log options is
// the composition root's job.
type Config struct {
	// Addr is the address the HTTP server listens on. APP_ADDR, default ":8080".
	Addr string

	// Log holds the already-parsed logging settings.
	Log LogConfig
}

// LogConfig is the validated logging configuration.
type LogConfig struct {
	// Level is parsed from APP_LOG_LEVEL. Default: info.
	Level slog.Level

	// Format is parsed from APP_LOG_FORMAT. Default: text.
	Format LogFormat
}

// LogFormat is a validated APP_LOG_FORMAT value.
//
// A named string rather than a bool: a bool encodes "two choices" into the type
// itself, so a third format would make the field name wrong and every reader of
// it wrong with it. It also stays close to what it represents — the environment
// hands over strings, and this is that string after validation.
//
// This deliberately does not reuse pkg/log's Format. The two describe different
// things: one is what the environment asked for, the other is how a handler
// encodes. They happen to share two members today; coupling them would mean
// pkg/config importing pkg/log, which is what keeps them independent.
type LogFormat string

// The formats APP_LOG_FORMAT accepts.
const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

// Load reads the environment and validates it.
//
// It reports every problem at once rather than stopping at the first. Failing on
// one variable at a time means one restart per typo, and the person fixing them
// has no way to see how many are left.
func Load() (Config, error) {
	var problems []error

	cfg := Config{Addr: getenv("APP_ADDR", ":8080")}

	if !strings.Contains(cfg.Addr, ":") {
		problems = append(problems, fmt.Errorf(
			"APP_ADDR=%q: want host:port, for example :8080", cfg.Addr))
	}

	rawLevel := getenv("APP_LOG_LEVEL", "info")
	if err := cfg.Log.Level.UnmarshalText([]byte(rawLevel)); err != nil {
		problems = append(problems, fmt.Errorf(
			"APP_LOG_LEVEL=%q: want debug, info, warn or error", rawLevel))
	}

	switch rawFormat := LogFormat(strings.ToLower(getenv("APP_LOG_FORMAT", "text"))); rawFormat {
	case LogFormatText, LogFormatJSON:
		cfg.Log.Format = rawFormat
	default:
		problems = append(problems, fmt.Errorf(
			"APP_LOG_FORMAT=%q: want %s or %s", rawFormat, LogFormatText, LogFormatJSON))
	}

	// An invalid value is louder than a silent fallback. A mistyped
	// APP_LOG_LEVEL that quietly resolves to info costs an hour of wondering why
	// the debug lines never appear; refusing to start costs a second.
	if len(problems) > 0 {
		return Config{}, fmt.Errorf("invalid configuration: %w", errors.Join(problems...))
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
