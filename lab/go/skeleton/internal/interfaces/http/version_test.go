package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	apperrors "skeleton/internal/application/errors"
	"skeleton/internal/application/usecase/query"
	apihttp "skeleton/internal/interfaces/http"
	"skeleton/internal/interfaces/http/apigen"
	"skeleton/pkg/i18n"
	"skeleton/pkg/log"
)

// stubVersionService stands in for in.VersionService.
//
// It is hand-written for the same reason stubProvider is: the port has one
// method.
type stubVersionService struct {
	result query.GetVersionResult
	err    error
}

// GetVersion satisfies in.VersionService, returning whatever the test set up.
func (s stubVersionService) GetVersion(context.Context, query.GetVersion) (query.GetVersionResult, error) {
	return s.result, s.err
}

func TestGetVersionReturnsThe200Response(t *testing.T) {
	s := apihttp.NewServer(stubVersionService{result: query.GetVersionResult{Value: "1.2.3"}}, log.Discard(), i18n.NewBundle(nil))

	resp, err := s.GetVersion(context.Background(), apigen.GetVersionRequestObject{})
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}

	got, ok := resp.(apigen.GetVersion200JSONResponse)
	if !ok {
		t.Fatalf("want a 200 response, got %T", resp)
	}
	if got.Version != "1.2.3" {
		t.Errorf("Version: want 1.2.3, got %q", got.Version)
	}
}

// The contract declares a 500 for this operation, and the generator emitted the
// type for it. Nothing returned that type until now, so this is what turns a
// dead branch into a live one. The acceptance suite cannot reach it: reading a
// build stamp does not fail in production, and inventing a scenario for it would
// be inventing a business behaviour.
func TestGetVersionReturnsThe500ProblemWhenTheUseCaseFails(t *testing.T) {
	s := apihttp.NewServer(stubVersionService{err: errors.New("provider unavailable")}, log.Discard(), i18n.NewBundle(nil))

	resp, err := s.GetVersion(context.Background(), apigen.GetVersionRequestObject{})
	if err != nil {
		t.Fatalf("the failure must travel in the response, not as a returned error: %v", err)
	}

	got, ok := resp.(apigen.GetVersion500ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Fatalf("want the generated 500 response, got %T", resp)
	}
	if got.Status != http.StatusInternalServerError {
		t.Errorf("Status: want 500, got %d", got.Status)
	}
}

// ToProblem sets Problem.Status from errmap.StatusFor's classification — 503
// for KindUnavailable — but /version's contract has only one failure
// response, so VisitGetVersionResponse always writes a 500 status line. RFC
// 9457 §3.1.4 requires the body's status to agree with the line it travels
// under; an uncorrected body would say "status": 503 underneath a 500,
// contradicting itself. The test above cannot catch this: it stubs a plain
// error, which takes ToProblem's default branch where Status is already 500,
// so the two numbers can never disagree there. Only a classified error
// reaches the branch where they can.
func TestGetVersionCoercesAClassifiedBodyStatusToMatchTheStatusLine(t *testing.T) {
	classified := apperrors.Error{
		Kind:       apperrors.KindUnavailable,
		MessageKey: "version.unavailable",
	}.Wrap(errors.New("provider unavailable"))

	s := apihttp.NewServer(stubVersionService{err: classified}, log.Discard(), i18n.NewBundle(nil))

	resp, err := s.GetVersion(context.Background(), apigen.GetVersionRequestObject{})
	if err != nil {
		t.Fatalf("the failure must travel in the response, not as a returned error: %v", err)
	}

	got, ok := resp.(apigen.GetVersion500ApplicationProblemPlusJSONResponse)
	if !ok {
		t.Fatalf("want the generated 500 response, got %T", resp)
	}

	body, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal problem: %v", err)
	}
	var decoded struct {
		Status int `json:"status"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("unmarshal problem: %v", err)
	}
	if decoded.Status != http.StatusInternalServerError {
		t.Errorf("body status: want %d to match the 500 status line, got %d: %s", http.StatusInternalServerError, decoded.Status, body)
	}
}

// A 500 tells its recipient nothing, on purpose: errmap keeps Where,
// DetailedError and the cause out of every body, so the operator woken up by
// that 500 has exactly one place left to look. This is the test for that place.
//
// Two failures it must catch, and neither is about a particular string:
//
//   - The handler stops logging. The suite would otherwise stay green, because
//     every other test here passes log.Discard() and a discarded record is
//     indistinguishable from a record never written. A failure that produces a
//     generic body and no record is an outage with no way in.
//
//   - Error() stops rendering Where. The chain "usecase.GetVersion" is the only
//     thing in the record that says *which* operation failed; without it the
//     operator has a cause and no location, and in a service with more than one
//     use case that is not enough to act on. Three tests already prove Where
//     never reaches a client; this is the one proving it reaches an operator.
//
// The final assertion pins the two channels together. What the operator must
// see and what the client must not see are the same values, so a change that
// widens one had better fail the test guarding the other.
func TestGetVersionLogsWhereTheFailureHappenedAndWhatCausedIt(t *testing.T) {
	const (
		where = "usecase.GetVersion"
		cause = "build stamp unreadable"
	)

	// Built the way usecase/query.GetVersion builds it, so that what is
	// asserted here is the record a real failure would produce.
	classified := apperrors.Error{
		Kind:       apperrors.KindUnavailable,
		MessageKey: "version.unavailable",
		Where:      where,
	}.Wrap(fmt.Errorf("read build version: %w", errors.New(cause)))

	logger := newLogRecorder()
	s := apihttp.NewServer(stubVersionService{err: classified}, logger, i18n.NewBundle(nil))

	resp, err := s.GetVersion(context.Background(), apigen.GetVersionRequestObject{})
	if err != nil {
		t.Fatalf("the failure must travel in the response, not as a returned error: %v", err)
	}

	records := logger.errorRecords()
	if len(records) == 0 {
		t.Fatal("the failure was rendered into a deliberately empty body and logged nowhere: nothing anywhere records that it happened")
	}

	var located, explained bool
	for _, r := range records {
		if strings.Contains(r.Text, where) {
			located = true
		}
		if strings.Contains(r.Text, cause) {
			explained = true
		}
	}
	if !located {
		t.Errorf("no record names the operation that failed (%q), so an operator cannot locate it: %v", where, records)
	}
	if !explained {
		t.Errorf("no record carries the cause (%q), so an operator cannot explain it: %v", cause, records)
	}

	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}
	for _, internal := range []string{where, cause} {
		if strings.Contains(string(body), internal) {
			t.Errorf("the response body carries %q, which belongs to the log channel alone: %s", internal, body)
		}
	}
}
