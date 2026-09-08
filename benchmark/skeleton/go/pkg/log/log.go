// Package log defines the logging contract and builds an slog-backed
// implementation of it.
//
// It sits in pkg/ rather than internal/ because logging is the one concern with
// no layer of its own: the composition root, adapters and infrastructure all
// need it, and none of them owns it.
//
// This package depends on nothing else in the module — in particular not on
// pkg/config. Configuration decides *what* to build; this decides *how*. Wiring
// the two together is the composition root's job, which is what keeps either of
// them replaceable on its own.
package log

import (
	"io"
	"log/slog"
	"os"
)

// Logger is what the rest of the codebase depends on.
//
// The usual Go advice — let the consumer declare the interface it needs — is
// right when there is one consumer. Logging has every layer as a consumer, and
// each declaring its own three-method interface would duplicate one contract
// many times over and let the copies drift.
//
// Returning an interface from a constructor is normally a smell; here the
// interface *is* the package's contract and the concrete type has nothing extra
// to offer, so there is nothing to lose by hiding it.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)

	// With returns a logger that carries the given key/value pairs on every
	// subsequent record — a request id, a workout id, whatever identifies the
	// unit of work.
	With(args ...any) Logger
}

// Format selects the output encoding.
//
// A named type rather than a string: a string field can hold "xml", which then
// needs validating, defaulting and documenting in two places. This cannot hold
// an invalid value, so the package needs no validation at all.
type Format int

const (
	// FormatText emits human-readable key=value lines. This is the zero value,
	// because the common case for this skeleton is a terminal.
	FormatText Format = iota

	// FormatJSON emits structured JSON, one object per line.
	FormatJSON
)

// Options is the resolved configuration of a Logger.
//
// It is exported so the full set of knobs is visible in one place rather than
// scattered across the With… functions. Build it through those functions rather
// than by hand — they are what keeps New's signature stable when a field is
// added.
type Options struct {
	// Level is the minimum level that will be emitted. Zero value: info.
	Level slog.Level

	// Format selects text or JSON. Zero value: text.
	Format Format

	// Writer is where records go. Zero value: os.Stdout.
	Writer io.Writer
}

// Option configures a Logger.
type Option func(*Options)

// WithLevel sets the minimum level that will be emitted.
func WithLevel(level slog.Level) Option {
	return func(o *Options) { o.Level = level }
}

// WithFormat selects the output encoding.
func WithFormat(format Format) Option {
	return func(o *Options) { o.Format = format }
}

// WithWriter sends output somewhere other than stdout.
func WithWriter(w io.Writer) Option {
	return func(o *Options) { o.Writer = w }
}

// New returns a Logger. With no options it writes text at info level to stdout.
func New(opts ...Option) Logger {
	options := Options{
		Level:  slog.LevelInfo,
		Format: FormatText,
		Writer: os.Stdout,
	}

	for _, opt := range opts {
		opt(&options)
	}

	handlerOpts := &slog.HandlerOptions{Level: options.Level}

	var handler slog.Handler
	switch options.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(options.Writer, handlerOpts)
	default:
		handler = slog.NewTextHandler(options.Writer, handlerOpts)
	}

	return slogLogger{slog.New(handler)}
}

// Discard returns a Logger that writes nothing.
//
// Tests exercising code that logs should pass this rather than letting output
// leak into the test report — a failing test is harder to read when it is buried
// in log lines that had nothing to do with the failure.
func Discard() Logger {
	return slogLogger{slog.New(slog.DiscardHandler)}
}

// slogLogger adapts *slog.Logger to Logger.
//
// Embedding gives Debug, Info, Warn and Error for free; only With needs
// overriding, because slog's returns *slog.Logger rather than the interface.
type slogLogger struct{ *slog.Logger }

// With returns a logger carrying the given key/value pairs, satisfying Logger
// where slog's own With cannot: it returns *slog.Logger, not the interface.
func (l slogLogger) With(args ...any) Logger {
	return slogLogger{l.Logger.With(args...)}
}
