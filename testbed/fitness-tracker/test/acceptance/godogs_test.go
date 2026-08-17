package acceptance_test

// Acceptance suite: runs the .feature files in ../../features through godog as
// Go subtests.
//
// Sample commands (note the package path comes before any -godog.* flag, or
// `go test` falls back to building `.` and fails with "no Go files"):
//
//	go test ./test/acceptance/
//	go test -test.v ./test/acceptance/ -test.run "^TestFeatures/"
//	go test ./test/acceptance/ -godog.tags=@wip
//	go test ./test/acceptance/ -godog.paths=../../features/volume.feature
//	go test -test.v ./test/acceptance/ -test.run "^TestFeatures/A_single_barbell_set$"
//	go test ./test/acceptance/ -godog.format=cucumber:report.json,pretty
//	go test ./test/acceptance/ -godog.definitions   # list registered steps
//	go test ./test/acceptance/ -godog.help
//
// Running one scenario by name via -test.run is often more precise than tagging
// it @wip, because it needs no edit to the feature file — which means no risk of
// committing a stray tag.

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

// opts carries the defaults. BindFlags turns every option godog knows about
// into a `-godog.*` flag using these values as the flag defaults, so the whole
// option set is reachable from the command line without hand-rolling flags.
var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
	Paths:  []string{"../../features"},

	// Without Strict, godog treats undefined and pending steps as TODOs and
	// still passes. That would let a scenario with no implementation report
	// green, which breaks the outer loop of outside-in development: the
	// scenario has to be able to go red.
	Strict: true,

	// Scenarios must be independent, so run them in parallel and let any hidden
	// coupling surface as a failure rather than as an occasional mystery. This
	// is only safe because per-scenario state lives in the context rather than
	// in package variables — see world_test.go.
	Concurrency: 4,

	// -1 asks for a fresh seed each run; the seed is printed so a failing order
	// can be reproduced with -godog.random=<seed>.
	Randomize: -1,
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestFeatures(t *testing.T) {
	o := opts // copy, so TestingT does not leak into the package-level defaults
	o.TestingT = t

	status := godog.TestSuite{
		Name:                 "fitness-tracker",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &o,
	}.Run()

	// Status 2 is godog's "option error". It covers two very different things:
	// informational modes that print and stop (-godog.definitions), and genuine
	// misconfiguration (an unknown formatter, an unwritable output file).
	//
	// Upstream's example skips on 2 so the informational flags do not fail the
	// suite. Skipping on both is wrong here: `go test` prints a plain "ok" for a
	// skip, so a mistyped flag would look exactly like a passing run. Ask which
	// case it is instead of guessing.
	if status == 2 {
		if o.ShowStepDefinitions {
			t.Skip("listed step definitions; no scenarios were run")
		}
		t.Fatalf("godog rejected its options (status 2) — see the error above")
	}

	if status != 0 {
		t.Fatalf("expected status 0, got %d", status)
	}
}

// InitializeTestSuite holds setup shared by the whole run — migrations, a test
// container, fixtures loaded once. Anything per-scenario belongs in
// InitializeScenario instead, so scenarios stay independent.
func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {})
	ctx.AfterSuite(func() {})
}

// InitializeScenario gives each scenario a fresh world and registers the step
// definitions.
//
// Every scenario must start from the same state. With Concurrency above 1 this
// stops being a matter of hygiene and becomes a correctness requirement: two
// scenarios sharing a package-level variable would race.
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return withWorld(ctx, newWorld()), nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, scenarioErr error) (context.Context, error) {
		// Release per-scenario resources here. scenarioErr is non-nil when the
		// scenario failed, which is the place to dump state for diagnosis.
		return ctx, nil
	})

	registerSteps(ctx)
}
