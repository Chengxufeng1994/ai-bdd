package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"skeleton/pkg/log"
)

// hookCause reads like the errors these hooks are actually handed: a binding
// failure quoting the value it choked on, and a transport failure naming the
// socket. Both carry things a client must never be shown — a bearer token, an
// internal host — which is the whole reason the hooks are replaced.
const hookCause = `parsing "sk-live-DEADBEEF" from db-prod-3.internal at 10.4.2.9:5432: invalid UUID`

// Every generated error path is replaced, and each keeps the status its own
// path means. Two obligations pull against each other here: the client must
// learn nothing about the failure, and an operator must learn everything.
//
// The status half is not decoration. Collapsing a binding failure into a 500
// tells the client the server broke when the client sent a bad value, and it
// pages whoever owns the service for someone else's typo. Asserting only "the
// body is generic" would be satisfied by a handler that answered 500 to all of
// them.
//
// The hooks are taken from strictOptions and serverOptions rather than built
// fresh here, so the status each path answers with is the wired one and not one
// the test chose. They are called directly rather than driven through NewRouter
// because two of the four have no generated call site while /version is the
// only operation — it declares no parameters and no request body. Waiting for
// an operation that reaches them would mean shipping them untested for exactly
// as long as they are unwatched.
func TestEveryGeneratedErrorPathKeepsItsStatusAndLeaksNothing(t *testing.T) {
	tests := []struct {
		name string
		want int
		call func(log.Logger, *gin.Context, error)
	}{
		{
			// RequestErrorHandlerFunc: the spec accepted the request but the
			// generated glue could not bind it. The client's mistake.
			name: "request could not be bound",
			want: http.StatusBadRequest,
			call: func(l log.Logger, c *gin.Context, err error) {
				strictOptions(l).RequestErrorHandlerFunc(c, err)
			},
		},
		{
			// HandlerErrorFunc: a handler returned an error instead of a typed
			// response. The server's mistake.
			name: "handler returned an error",
			want: http.StatusInternalServerError,
			call: func(l log.Logger, c *gin.Context, err error) {
				strictOptions(l).HandlerErrorFunc(c, err)
			},
		},
		{
			// ResponseErrorHandlerFunc: the answer would not serialise, or the
			// client hung up mid-write.
			name: "response could not be delivered",
			want: http.StatusInternalServerError,
			call: func(l log.Logger, c *gin.Context, err error) {
				strictOptions(l).ResponseErrorHandlerFunc(c, err)
			},
		},
		{
			// GinServerOptions.ErrorHandler: parameter binding. This one is
			// handed its status by the generated wrapper, so the test asserts
			// the hook forwards it rather than choosing one.
			name: "parameter would not parse",
			want: http.StatusBadRequest,
			call: func(l log.Logger, c *gin.Context, err error) {
				serverOptions(l).ErrorHandler(c, err, http.StatusBadRequest)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logged bytes.Buffer
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tt.call(log.New(log.WithWriter(&logged)), ctx, errors.New(hookCause))

			if rec.Code != tt.want {
				t.Errorf("status: want %d, got %d", tt.want, rec.Code)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/problem+json" {
				t.Errorf("Content-Type: want application/problem+json, got %q", got)
			}

			// One list, two opposite assertions. That is the two-channel
			// contract stated as a test: nothing here may reach the client,
			// and all of it must reach the operator. Asserting either half
			// alone passes for a hook that silently drops the failure.
			secrets := []string{"sk-live-DEADBEEF", "db-prod-3", "10.4.2.9"}

			body := rec.Body.String()
			for _, secret := range secrets {
				if strings.Contains(body, secret) {
					t.Errorf("body leaks %q from the underlying error: %s", secret, body)
				}
			}

			// The document must still identify itself honestly: a status that
			// disagrees with the status line is the RFC 9457 §3.1.4 violation
			// this codebase has already shipped once.
			var problem struct {
				Status int    `json:"status"`
				Title  string `json:"title"`
				Type   string `json:"type"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &problem); err != nil {
				t.Fatalf("body is not a Problem document: %v: %s", err, body)
			}
			if problem.Status != tt.want {
				t.Errorf("body status: want %d, got %d", tt.want, problem.Status)
			}
			if problem.Title != http.StatusText(tt.want) {
				t.Errorf("title: want %q, got %q", http.StatusText(tt.want), problem.Title)
			}
			if problem.Type != "about:blank" {
				t.Errorf("type: want about:blank, got %q", problem.Type)
			}

			// A failure the body refuses to explain is one only the log can.
			// Checked per fragment rather than against hookCause whole,
			// because slog escapes the quotes inside the value.
			for _, needed := range secrets {
				if !strings.Contains(logged.String(), needed) {
					t.Errorf("no record carries %q, got %q", needed, logged.String())
				}
			}
		})
	}
}

// The hooks being correct is worth nothing if one is not wired, and that is the
// failure that actually happened: three of the four were set, the fourth kept
// oapi-codegen's leaking default, and every test passed.
//
// Reflection rather than four named assertions, because the regression this
// guards against is a hook nobody knows about yet. Regenerating apigen against
// a new oapi-codegen can add a fifth; a test naming four would stay green while
// the new one leaked. Any func field left nil fails here, which turns "a hook
// was added" from silent into a build-time conversation.
func TestEveryGeneratedErrorHookIsWired(t *testing.T) {
	logger := log.Discard()

	for name, opts := range map[string]any{
		"StrictGinServerOptions": strictOptions(logger),
		"GinServerOptions":       serverOptions(logger),
	} {
		v := reflect.ValueOf(opts)
		for i := range v.NumField() {
			field := v.Type().Field(i)
			if field.Type.Kind() != reflect.Func {
				continue
			}
			if v.Field(i).IsNil() {
				t.Errorf("%s.%s is nil: oapi-codegen fills every unset hook with a default that writes err.Error() into the response body", name, field.Name)
			}
		}
	}
}
