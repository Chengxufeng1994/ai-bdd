package http_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"skeleton/internal/application/dto"
	"skeleton/internal/application/query"
	"skeleton/internal/application/service"
	apihttp "skeleton/internal/interfaces/http"
	"skeleton/internal/interfaces/http/apigen"
)

// stubUseCase stands in for in.GetVersionUseCase.
//
// It is hand-written for the same reason stubProvider is: the port has one
// method. Note that satisfying the alias needs nothing more than the matching
// Handle signature — that is the known cost of D1 in the design.
type stubUseCase struct {
	result dto.Version
	err    error
}

// Handle satisfies in.GetVersionUseCase, returning whatever the test set up.
func (s stubUseCase) Handle(context.Context, query.GetVersion) (dto.Version, error) {
	return s.result, s.err
}

func TestGetVersionReturnsThe200Response(t *testing.T) {
	s := apihttp.NewServer(service.VersionService{
		GetVersion: stubUseCase{result: dto.Version{Value: "1.2.3"}},
	})

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
	s := apihttp.NewServer(service.VersionService{
		GetVersion: stubUseCase{err: errors.New("provider unavailable")},
	})

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
