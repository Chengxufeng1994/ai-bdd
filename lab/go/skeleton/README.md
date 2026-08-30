# lab/go/skeleton

A Go + godog project used to dogfood the `ai-bdd` skills end to end:
**CLARIFY → SPEC → PLAN → IMPLEMENT → VERIFY → REVIEW**.

It is not a product. Its job is to make the skills fail visibly when they are
weak.

## Why this domain

A training log looks like CRUD, so the interesting part is deliberately **total
training volume**, not record-keeping. Volume is where the rules live, and the
rules are genuinely ambiguous:

- Pull-ups carry no external weight — is the volume zero, or bodyweight × reps?
- Assisted pull-ups at −20 kg — negative volume, or (bodyweight − 20) × reps?
- A single-arm row at 20 kg × 8, done per side — 160 or 320?
- Does a 20 kg empty-bar warm-up count toward volume?
- A set targeting 10 reps that reached 7 — recorded as 7, or as "10, failed"?
- A dropset of 60 kg × 8 straight into 40 kg × 6 — one set or two?
- kg and lb mixed in one session
- Smith machine 60 kg versus barbell 60 kg — can they be summed?
- A 60-second plank has neither reps nor load — what is its volume?
- "Total volume" over what — one session, one exercise, one week, one muscle group?

Most people think of three or four of these unprompted.

**This is the bar for `bdd-clarify`:** a run that does not surface the
bodyweight and per-side questions is a run that failed, regardless of how
polished its output looks.

## One scenario at a time

This is a *skeleton*, not a fitness app. It hosts whichever scenario is being
used to evaluate the skills; `prompts/` names them:

```
prompts/1-fitness-tracker-clarify.md      a casual brief
prompts/1-fitness-tracker-clarify-prd.md  the same scenario as an approved PRD
docs/bdd/                                 that run's example maps and questions
features/                               the scenarios it produced
```

Only CLARIFY, SPEC and PLAN are language-agnostic enough to be swapped freely.
Once IMPLEMENT fills `internal/`, the code belongs to one scenario — so archive
a completed run before starting another rather than trying to keep two live.

Other languages get their own skeleton beside this one: `lab/python/skeleton`,
`lab/typescript/skeleton`. Everything upstream of IMPLEMENT is shared; only the
harness differs.

## Ground rule

**No business code until CLARIFY has run.** Every type in `domain/` must trace
back to a rule in an example map, and every scenario in `features/` back to a
concrete example. Writing the model first defeats the entire point of the
testbed — it would prove the skills work by handing them the answer.

Current state: `domain/` is empty. `application/` and `infrastructure/` hold the
`/version` read slice and nothing else — it encodes no business rule, and exists
to keep the walking skeleton's chain honest through every layer. The
only thing implemented is the walking skeleton below, which carries no domain
meaning and therefore does not pre-empt CLARIFY.

## The walking skeleton

`features/version.feature` and the `/version` endpoint exist to keep one thin
path through the whole chain alive from day one:

```
api/openapi.yaml → go generate → api.gen.go → server.go → version.feature → godog
```

A break anywhere in that chain then shows up as a failing test, rather than
being discovered while trying to write the first real feature — when it would be
tangled up with domain questions and much harder to diagnose.

### Regenerating from the contract

```bash
make gen-api      # or: go generate -run oapi-codegen ./internal/...
```

`make help` lists every target.

Every `//go:generate` directive lives in the build-tagged `generate.go` at the
repository root, and the Makefile sits on top of them. Two reasons: all
directives then run from the same directory, so every path — spec, config,
output — is root-relative and reads alike; and repo-scoped generators like
mockery, which expect to run from the root and belong to no package, have
somewhere to live.

`go generate -run <tool>` selects one generator at a time, which is what makes
per-tool Makefile targets possible without splitting the file up.

The generator is pinned as a `tool` dependency in `go.mod`, so it stays out of
the module's real requirements while everyone still regenerates with the same
version. `api.gen.go` is committed: a plain `go build` then needs no generator,
and a review diff shows exactly how a contract change altered the server
interface.

Generated code sits in its own `apigen` package so that "do not edit" is
structural rather than a naming convention — nothing hand-written shares a
namespace with it, so a type declared in `server.go` can never collide with one
the generator emits. The directory is therefore disposable, which doubles as the
cheapest check that the pipeline still works:

```bash
rm -rf internal/interfaces/http/apigen && go generate ./... && go test ./...
```

`api/cfg.yaml` sits beside the spec it configures. Its `output` path is resolved
relative to the directory the generator runs in — **not** relative to the config
file, verified by running it from elsewhere and watching the output land in the
wrong place. Because every directive runs from the root, every path in that file
is root-relative, which is the only arrangement where it reads correctly from
wherever you happen to open it.

### Response conventions

Success responses are shaped by the domain and differ per operation. Failures
share one shape everywhere: **RFC 9457 Problem Details**, served as
`application/problem+json`. Clients then need one error path instead of one per
endpoint.

`type` is the machine-readable identity of a failure and is what clients branch
on — HTTP status alone is too coarse, since one 409 can mean several different
things. `title` is for humans and may be reworded without notice.

`components/responses/` holds the reusable failures (400, 401, 403, 404, 409,
422, 429, 500, 503). **They are shared, not pasted onto every operation.**
Declaring 401 on an endpoint that never authenticates makes the contract lie,
and generated clients carry dead branches for it. oapi-codegen only emits code
for responses a path actually references, so an unused component costs nothing
until an endpoint declares it.

400 versus 422 is worth keeping straight: 400 means "I cannot read this",
422 means "I read it and it is not allowed". Conflating them robs clients of the
ability to tell a client bug from a user mistake.

## Layout

```
skeleton/
├── api/
│   ├── openapi.yaml         HTTP contract; paths derive from scenarios
│   └── cfg.yaml             oapi-codegen configuration
├── generate.go              every //go:generate directive, build-tagged
├── cmd/fitness/             composition root; wiring only
├── internal/
│   ├── domain/              business rules and entities; depends on nothing
│   ├── application/         use-case orchestration; depends on domain; declares ports
│   ├── infrastructure/      port implementations; depends on domain + application
│   └── interfaces/http/     HTTP adapter
│       ├── apigen/          generated — never edited, fully rebuildable
│       │   └── api.gen.go
│       ├── server.go        handlers; no gin types in sight
│       └── router.go        engine, middleware, route registration
├── docs/
│   ├── ARCHITECTURE.md      start here for the whole picture
│   └── DATAFLOW.md          every conversion on the path, and who owns it
├── prompts/                 fixed inputs for evaluating the skills
│   ├── 1-fitness-tracker-clarify.md       casual brief
│   └── 1-fitness-tracker-clarify-prd.md   the same scenario as an approved PRD
├── features/                .feature files — specifications, not test code
│   └── version.feature      the walking skeleton
└── test/acceptance/         the godog harness that executes them
    ├── godogs_test.go       runner, suite and scenario hooks
    ├── world_test.go        per-scenario state, carried in context.Context
    └── steps_test.go        step definitions
```

Dependencies point inward: everything may import `domain`, and `domain` imports
nothing. Each layer's `doc.go` states its own rule.

The layers sit under `internal/` because Go enforces that boundary at compile
time — no other module can import them. What remains importable from outside is
exactly the three contracts: the `cmd` binary, `api/openapi.yaml`, and
`features/`. That makes "nobody reaches past `interfaces/` straight into
`domain/`" a compiler guarantee rather than a code-review convention.

`features/` sits at the top level next to `api/openapi.yaml`, not inside
`test/`. Both are contracts meant to be read by people who do not read Go.
Filing feature files under `test/` quietly reframes them as test scripts, which
is the first step toward "BDD means writing Cucumber tests" — the failure mode
this whole project exists to avoid.

## Three levels of test

Which level a given scenario belongs at is decided in PLAN, not guessed here.
Pushing every scenario to the acceptance level is the most common way a
Cucumber suite becomes slow, flaky, and eventually switched off.

| Level | Where | Command |
| --- | --- | --- |
| **Unit** | `*_test.go` beside the code, no external dependencies | `go test ./...` |
| **Integration** | `*_test.go` with `//go:build integration`, real adapters | `go test -tags=integration ./...` |
| **Acceptance (BDD)** | `features/*.feature` + `steps_test.go` | `go test -count=1 .` |

```bash
go test ./...                                        # unit + acceptance
go test -tags=integration ./...                      # adds integration
go test -count=1 ./test/acceptance/ -godog.tags=@wip # only the current scenario
```

Note the argument order in the last one: `go test` needs the package path
*before* flags it does not recognise, or it falls back to building `.` and
fails with "no Go files".

## Organising features and step definitions

Step definitions are matched by regular expression against **every** feature
file. They are a single global namespace, not scoped to the feature they were
written for. Two consequences follow, and most Cucumber suites rot because they
are ignored.

**Do not mirror features with step files.** One `steps_x_test.go` per
`x.feature` looks tidy and guarantees duplication, because a phrase like
`a workout with 3 sets` is naturally shared. The two hierarchies are organised
on different axes on purpose:

| | Grouped by | Example |
| --- | --- | --- |
| `features/` | capability | `features/volume/total_volume.feature` |
| step definitions | domain noun | `steps_workout_test.go`, `steps_volume_test.go` |

**Search before writing a step.** Duplicate definitions that differ only in
wording are the beginning of a glue layer nobody can safely change.

### Keeping the chain machine-checkable

Gherkin's `Rule:` keyword maps onto Example Mapping's cards, so the specification
can be checked against the example map mechanically rather than by reading:

| Example map | Gherkin | What can be checked |
| --- | --- | --- |
| Story | `Feature:` | — |
| Rule Rule 1 (blue card) | `Rule:` | every rule has a `Rule:` block |
| Example Example 1.1 (green card) | `@Example1.1` on a `Scenario` | every example has a scenario; no orphan scenarios |

```gherkin
Feature: Total training volume

  Rule: Volume is the sum of load times reps across working sets

    @Example1.1
    Scenario: A single barbell set
      ...
```

This is why example-map numbering must never be renumbered — the tags are
references, and renumbering silently repoints them.

### Tags

| Tag | Meaning |
| --- | --- |
| `@Example1.1` | traces to an example map entry |
| `@wip` | being worked on right now |
| `@smoke` | must pass before anything else is trusted |
| `@slow` | excluded from the fast loop (`-godog.tags="~@slow"`) |

### Health signals

- Step-definition count growing as fast as scenario count means steps are
  written imperatively; declarative steps get reused and the curve flattens.
- A step definition matched by no scenario is an orphan — delete it.
- A `Scenario Outline` with a large Examples table is usually a unit test that
  escaped into the acceptance level. It belongs one layer down.

`gopls` reports "No packages found" for integration files — that is the build
tag doing its job, not an error.

### Runner behaviour worth knowing

`godog.BindFlags("godog.", flag.CommandLine, &opts)` exposes godog's whole
option set as `-godog.*` flags, using the values in `opts` as their defaults.
Anything below can therefore be overridden per run without editing code:

- **`Strict: true`.** By default godog treats undefined and pending steps as
  TODOs and still passes, which would let a scenario with no implementation
  report green and break the outer loop of outside-in development. Strict makes
  an unimplemented scenario go red, as it must.
- **`Concurrency: 4`.** Scenarios must be independent, so they run in parallel
  and any hidden coupling surfaces as a failure. This is only safe because
  per-scenario state lives in the context rather than in package variables —
  see `world_test.go`. Output is buffered until each scenario finishes, so
  parallel runs do not interleave.
- **`Randomize: -1`.** A fresh seed each run; the seed is printed, so a failing
  order can be reproduced with `-godog.random=<seed>`.

```bash
go test ./test/acceptance/ -godog.concurrency=1 -godog.format=progress
```

With no feature files present godog reports `No scenarios` and passes — a green
run here means the harness works, not that anything is implemented.
