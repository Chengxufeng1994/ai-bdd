package http_test

import (
	"fmt"
	"strings"
	"sync"

	"skeleton/pkg/log"
)

// logRecord is one thing a logRecorder was told to log.
type logRecord struct {
	// Level is the severity the call site chose.
	Level string

	// Text is the record rendered the way a handler would emit it: the
	// message followed by every argument. Rendering with %v is what puts an
	// error's Error() — the method that assembles Where, DetailedError and
	// the cause into one line — inside what a test can assert on.
	Text string
}

// logRecorder is a log.Logger that keeps its records instead of writing them.
//
// Every other test in this package passes log.Discard(), which makes the log
// channel unobservable: a test holding one cannot tell "logged the cause" from
// "logged nothing at all". That matters here more than it would elsewhere,
// because half of what an application error carries — Where, DetailedError and
// the wrapped cause — is deliberately kept out of every response body, so the
// log channel is the *only* place an operator can ever see it. Asserting a
// generic body without asserting the record proves the redaction and leaves
// open the possibility that the failure vanished silently.
//
// It is safe for concurrent use because the handlers under test may be driven
// from more than one goroutine, and because the suite runs under -race.
type logRecorder struct {
	mu      *sync.Mutex
	records *[]logRecord
	base    []any
}

// Compile-time check that logRecorder still satisfies the contract it fakes.
var _ log.Logger = (*logRecorder)(nil)

// newLogRecorder returns an empty logRecorder.
func newLogRecorder() *logRecorder {
	return &logRecorder{mu: &sync.Mutex{}, records: &[]logRecord{}}
}

// Debug records a debug-level call.
//
// The comment is not decoration: .golangci.yml enables revive's `exported`
// rule with checkPrivateReceivers — see query_test's stubProvider.Version.
func (l *logRecorder) Debug(msg string, args ...any) { l.record("DEBUG", msg, args) }

// Info records an info-level call.
func (l *logRecorder) Info(msg string, args ...any) { l.record("INFO", msg, args) }

// Warn records a warn-level call.
func (l *logRecorder) Warn(msg string, args ...any) { l.record("WARN", msg, args) }

// Error records an error-level call.
func (l *logRecorder) Error(msg string, args ...any) { l.record("ERROR", msg, args) }

// With returns a recorder that carries args on every subsequent record and
// shares this one's storage, so a caller that decorates its logger is still
// observed by the test that built it.
func (l *logRecorder) With(args ...any) log.Logger {
	return &logRecorder{mu: l.mu, records: l.records, base: append(l.base[:len(l.base):len(l.base)], args...)}
}

func (l *logRecorder) record(level, msg string, args []any) {
	var b strings.Builder
	b.WriteString(msg)
	for _, arg := range append(l.base[:len(l.base):len(l.base)], args...) {
		fmt.Fprintf(&b, " %v", arg)
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	*l.records = append(*l.records, logRecord{Level: level, Text: b.String()})
}

// errorRecords returns every error-level record, which is the only severity
// a failing request is expected to produce.
func (l *logRecorder) errorRecords() []logRecord {
	l.mu.Lock()
	defer l.mu.Unlock()

	var out []logRecord
	for _, r := range *l.records {
		if r.Level == "ERROR" {
			out = append(out, r)
		}
	}
	return out
}
