package acceptance_test

import (
	"context"
	"errors"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// world is the state one scenario builds up as its steps run: what was set up
// in Given, what happened in When, what has to be asserted in Then.
//
// It is carried in the scenario's context.Context rather than in a package
// variable, so scenarios cannot leak state into each other. With Concurrency
// above 1 that is not hygiene, it is correctness.
//
// Fields are added as steps need them. Keep it to test state — application
// state belongs in the layers under test.
type world struct {
	// router is the assembled HTTP surface for this scenario. Each scenario
	// builds its own, so two scenarios can run servers configured differently
	// at the same time.
	router *gin.Engine

	// resp holds the last response, so a Then step can assert on it instead of
	// the When step having to know what will be checked.
	resp *httptest.ResponseRecorder
}

type worldKey struct{}

func newWorld() *world { return &world{} }

func withWorld(ctx context.Context, w *world) context.Context {
	return context.WithValue(ctx, worldKey{}, w)
}

// worldFrom retrieves the current scenario's world.
//
// It returns an error rather than panicking so a misconfigured suite reports a
// readable failure instead of a stack trace.
func worldFrom(ctx context.Context) (*world, error) {
	w, ok := ctx.Value(worldKey{}).(*world)
	if !ok {
		return nil, errors.New("no world in context: InitializeScenario must install one in a Before hook")
	}
	return w, nil
}

func init() {
	// Keep gin from printing its startup banner and debug warnings for every
	// scenario; a failing acceptance run should be readable.
	gin.SetMode(gin.TestMode)
}
