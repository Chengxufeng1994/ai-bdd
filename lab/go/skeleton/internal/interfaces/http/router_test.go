package http_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"skeleton/internal/application/usecase/query"
	apihttp "skeleton/internal/interfaces/http"
	"skeleton/internal/interfaces/http/apigen"
	"skeleton/pkg/i18n"
)

// leakyCause is what the failing handler below returns. It reads like the
// errors that actually reach this path in production — a driver message
// naming a user and an internal host — because the point of the test is that
// none of it reaches a client and all of it reaches an operator.
const leakyCause = `pq: password authentication failed for user "admin" at db-prod-3.internal`

// leakyServer implements apigen.StrictServerInterface by returning an error
// instead of a typed response — the mistake interfaces/doc.go's MUST forbids a
// handler from making. Every real handler in this codebase honours the rule,
// so this is the only way to drive the error hooks NewRouter configures.
type leakyServer struct{}

// GetVersion always fails, simulating a handler that returned its error
// instead of a typed 500 response.
func (leakyServer) GetVersion(context.Context, apigen.GetVersionRequestObject) (apigen.GetVersionResponseObject, error) {
	return nil, errors.New(leakyCause)
}

// The backstop has two obligations and they pull in opposite directions: the
// client must learn nothing, and an operator must learn everything. Asserting
// only the first would be satisfied by a router that dropped the failure on
// the floor — which is what this one did before it was given a logger, and
// which no test could tell apart from correct redaction.
//
// ARCHITECTURE.md §7 requires a 500 body stay generic, and interfaces/doc.go's
// MUST assumes every handler cooperates by never returning an error. This
// covers the backstop for when one does not: NewRouter must not fall back to
// oapi-codegen's generated defaults, which write err.Error() straight into the
// body (see apigen.NewStrictHandlerWithOptions, which fills every hook left
// nil with exactly that).
//
// The body assertion scans the whole rendered document rather than one field,
// mirroring errmap_test.go, so a new field cannot open a second leak quietly.
func TestRouterBackstopKeepsTheBodyGenericAndStillLeavesATrace(t *testing.T) {
	logger := newLogRecorder()

	engine, err := apihttp.NewRouter(leakyServer{}, logger)
	if err != nil {
		t.Fatalf("build router: %v", err)
	}

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("Code: want 500, got %d: %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Errorf("Content-Type: want application/problem+json, got %q", got)
	}

	body := rec.Body.String()
	for _, secret := range []string{"password", "admin", "db-prod-3"} {
		if strings.Contains(body, secret) {
			t.Errorf("500 body leaks %q from the underlying error: %s", secret, body)
		}
	}

	// The other half. A 500 nobody can explain is an outage nobody can end,
	// so the cause the body just refused to carry has to be somewhere.
	records := logger.errorRecords()
	if len(records) == 0 {
		t.Fatal("the backstop 500 left no error record: the cause reaches the log channel or it reaches nothing at all")
	}

	var found bool
	for _, r := range records {
		if strings.Contains(r.Text, leakyCause) {
			found = true
		}
	}
	if !found {
		t.Errorf("no error record carries the cause %q, got %v", leakyCause, records)
	}
}

// droppedCause is what a write to a client that went away actually says. The
// address pair is the point: it names a host inside the cluster and the port
// it serves on, which is why this error must never be rendered into a body.
const droppedCause = `write tcp 10.4.2.9:8080->203.0.113.7:52344: broken pipe`

// droppedConnection is an http.ResponseWriter that fails every write while
// recording what was handed to it, standing in for a client that hung up
// mid-response.
//
// Recording the bytes is what makes the assertion possible at all: the
// interesting question is not what the client received — it received nothing —
// but what the server tried to send it. A leak here is a leak whether or not
// the socket survived to carry it.
type droppedConnection struct {
	header http.Header

	// sent accumulates every byte the server attempted to write.
	sent bytes.Buffer
}

// Header satisfies http.ResponseWriter.
func (w *droppedConnection) Header() http.Header {
	if w.header == nil {
		w.header = http.Header{}
	}
	return w.header
}

// Write records p and then fails, the way a write to a closed socket does.
func (w *droppedConnection) Write(p []byte) (int, error) {
	w.sent.Write(p)
	return len(p), errors.New(droppedCause)
}

// WriteHeader satisfies http.ResponseWriter.
func (w *droppedConnection) WriteHeader(int) {}

// ResponseErrorHandlerFunc is the hook nobody writes a handler for, and it is
// the one that gets handed transport errors: the write failed, the client
// vanished, the response would not serialise. Those errors carry socket
// addresses and internal hostnames, so oapi-codegen's default for this hook —
// gin.H{"msg": err.Error()}, identical to the HandlerErrorFunc default — puts
// an internal IP and port into the response stream.
//
// NewStrictHandlerWithOptions fills every hook left nil with that default, so
// configuring HandlerErrorFunc alone is not "the leak is fixed", it is "one of
// three leaks is fixed". This drives the second one.
func TestRouterResponseErrorHookDoesNotLeakTheTransportError(t *testing.T) {
	logger := newLogRecorder()

	// A working handler, not a leaky one: this path is reached after the
	// handler has already done everything right and the failure is in
	// delivering its answer.
	engine, err := apihttp.NewRouter(apihttp.NewServer(stubVersionService{result: query.GetVersionResult{Value: "1.2.3"}}, logger, i18n.NewBundle(nil)), logger)
	if err != nil {
		t.Fatalf("build router: %v", err)
	}

	w := &droppedConnection{}
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/version", nil))

	if sent := w.sent.String(); strings.Contains(sent, "10.4.2.9") || strings.Contains(sent, "203.0.113.7") {
		t.Errorf("the response stream leaks the transport error's addresses: %s", sent)
	}

	records := logger.errorRecords()
	if len(records) == 0 {
		t.Fatal("a failed write left no error record: the cause reaches the log channel or it reaches nothing at all")
	}
	var found bool
	for _, r := range records {
		if strings.Contains(r.Text, droppedCause) {
			found = true
		}
	}
	if !found {
		t.Errorf("no error record carries the transport failure %q, got %v", droppedCause, records)
	}
}
