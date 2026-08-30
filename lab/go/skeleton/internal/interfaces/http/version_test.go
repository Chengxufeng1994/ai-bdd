package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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
