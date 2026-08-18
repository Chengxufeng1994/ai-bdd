# Architecture Overview

A living document for agents and humans arriving at this codebase. Update it as
the codebase evolves.

> **Read this first.** This is a *testbed*, not a product. It exists to run the
> `ai-bdd` skills end to end and make them fail visibly when they are weak. Most
> of it is deliberately empty — see §9.

## Related documents

This document gives orientation. It deliberately does not repeat what these say,
because two copies of a rule drift apart and then nobody knows which is current.

| Document | Read it when | Kept separate because |
| --- | --- | --- |
| [DATAFLOW.md](./DATAFLOW.md) | About to write a handler, mapper, assembler or presenter | It is consulted repeatedly while coding; this document is read once on arrival |
| [../README.md](../README.md) | You want to know why this testbed exists and what "done" means for it | It is the testbed's charter, not its architecture |
| `internal/*/doc.go` | Working inside one layer and unsure what it may import | Rules belong next to the code they constrain, where a compiler error sends you |
| [../prompts/1-fitness-tracker-clarify.md](../prompts/1-fitness-tracker-clarify.md) | Evaluating the skills | It is a test fixture, not documentation |
| `make help` | Looking for a command | The Makefile is the authoritative list; a copy here would go stale |

The `doc.go` files are authoritative for their own layer. Where this document and
a `doc.go` disagree, the `doc.go` is right and this one needs updating.

---

## 1. Project Structure

```
skeleton/
├── api/
│   ├── openapi.yaml          the HTTP contract; source of truth for §3.2
│   └── cfg.yaml              oapi-codegen configuration
├── cmd/server/               composition root — wiring only, no decisions
├── internal/                 Go enforces this boundary at compile time
│   ├── domain/               business rules and entities; imports nothing
│   ├── application/          use-case orchestration; declares ports
│   ├── infrastructure/       port implementations (SQL, clients, clocks)
│   └── interfaces/http/      HTTP adapter
│       ├── apigen/           generated from openapi.yaml — never edited
│       ├── server.go         handlers; no gin types in their signatures
│       └── router.go         engine, spec validation, route registration
├── pkg/
│   ├── config/               the only place that reads os.Getenv
│   └── log/                  the Logger interface, backed by slog
├── features/                 .feature files — specifications, not test code
├── test/acceptance/          the godog harness that executes them
├── prompts/                  fixed inputs for evaluating the skills
├── docs/
│   ├── ARCHITECTURE.md       this document
│   └── DATAFLOW.md           the full request-to-response chain
├── generate.go               every //go:generate directive, build-tagged
├── Makefile                  the entry point: make help
└── go.mod                    deps plus pinned `tool` dependencies
```

Each layer's `doc.go` states its own rules; those files are the authoritative
version and this document does not repeat them.

---

## 2. High-Level System Diagram

Runtime — one process, one protocol so far:

```
[HTTP client] ──▶ [gin router] ──▶ [generated strict handler]
                                          │
                                          ▼
                                   [application]  ◀── ports ──▶ [infrastructure]
                                          │                       (none yet)
                                          ▼
                                     [domain]
```

Development — the chain the testbed exists to exercise:

```
prompts/ ──▶ CLARIFY ──▶ docs/bdd/*/example-mapping.md
                              │
                              ▼
                            SPEC ──▶ features/*.feature
                                          │
             api/openapi.yaml ──▶ generate ──▶ apigen/
                                          │
                                          ▼
                        IMPLEMENT ──▶ internal/** ──▶ VERIFY ──▶ go test
```

Dependencies point inward. `infrastructure` looks inverted but is not: it
imports application's port interfaces, and application never imports it.

### The read and write paths cost different amounts, on purpose

```
write   apigen request ─▶ command ─▶ aggregate ─▶ dto ─▶ apigen response   3 conversions
read    apigen params  ─▶ query   ──────────────▶ dto ─▶ apigen response   2 conversions
                                    ↑
                          reads never load the aggregate
```

That asymmetry is the whole return on CQRS: the read side is served by its own
port returning a view shaped for what the caller displays, so it can later come
from a denormalised table or a cache without touching the domain. A read path
that costs the same as a write path means the split is decoration.

`docs/DATAFLOW.md` has the full chain: what each conversion is called, where it
lives, which of them may fail, and the three levels at which input is validated.

---

## 3. Core Components

### 3.1. Frontend

**None.** No UI exists and none is planned — the skills under test operate on
requirements and specifications, and a UI would add cost without adding a single
new question for them to surface.

### 3.2. Backend Services

#### 3.2.1. fitness

- **Description**: A single HTTP service. Currently exposes exactly one
  endpoint, `GET /version`, which reports the version the binary was built from.
  It carries no domain meaning; it exists to keep one thin path through the whole
  generation-and-test chain alive so a break anywhere shows up as a failing test.
- **Technologies**: Go 1.26, gin, oapi-codegen (contract-first), godog
- **Deployment**: none — see §6
- **Entry point**: `cmd/fitness`, configurable via `APP_ADDR` (default `:8080`)

Business endpoints are absent on purpose. See §9.

---

## 4. Data Stores

**None yet.** Nothing is persisted; the service holds no state between requests.
PostgreSQL with pgx, sqlc and golang-migrate is chosen for when a scenario needs
it — see §3.2.

`internal/infrastructure/` is where a repository would live, and
`internal/application/port/out/` is where its interface would be declared — the
application states what it needs and infrastructure complies, never the reverse.

When a store does arrive, its row-shaped model stays unexported inside the
repository. A public DTO that mirrors an aggregate field-for-field recreates the
coupling the port was introduced to break; see
`internal/infrastructure/doc.go`.

---

## 5. External Integrations / APIs

**None.**

---

## 6. Deployment & Infrastructure

**Not deployed.** No cloud provider, no CI pipeline, no monitoring.

`make build` produces `bin/fitness` with the version stamped in via ldflags:

```bash
make build   # -X main.version=$(git describe --tags --always --dirty)
```

`make verify` is what a CI job would run today: formatting, generated-code
staleness, lint, `go vet`, and `go test -race`. Every one of those gates has
been confirmed to fail when it should.

---

## 7. Security Considerations

**No authentication or authorisation exists.** `GET /version` is unauthenticated
and exposes only a build string.

Requests are validated against `api/openapi.yaml` by middleware before reaching a
handler, so anything off-contract is rejected at the edge — verified: a `POST` to
`/version` returns 400 rather than reaching the router, and an undeclared path
returns 404.

Two decisions already made that will matter once endpoints arrive:

- Errors follow **RFC 9457 Problem Details**. The `detail` field on a 500 must
  stay generic — stack traces, SQL fragments and internal hostnames belong in
  logs, not in a response body that reaches whoever made the request.
- The reusable error responses in `api/openapi.yaml` include 401 and 403, but
  they are **not** attached to `/version`. Declaring an auth failure on an
  endpoint that never authenticates makes the contract lie.

---

## 8. Development & Testing Environment

```bash
make help              # every target
make gen               # regenerate from api/openapi.yaml
make test              # unit + acceptance, with -race
make test-integration  # adds tests behind the `integration` build tag
make verify            # what CI would run
```

### Configuration

Every setting arrives through the environment, and `pkg/config` is the only
package that reads it. Anything reading `os.Getenv` elsewhere cannot be varied in
a test without setting global state, and its default stops being visible.

| Variable | Default | Accepts |
| --- | --- | --- |
| `APP_ADDR` | `:8080` | `host:port` |
| `APP_LOG_LEVEL` | `info` | `debug` `info` `warn` `error` |
| `APP_LOG_FORMAT` | `text` | `text` `json` |

An unrecognised value refuses to start and reports **every** offending variable
at once. A silent fallback is worse than a crash here: a mistyped
`APP_LOG_LEVEL` that quietly resolves to `info` costs an hour of wondering why
the debug lines never appear.

`pkg/config` and `pkg/log` do not import each other. Configuration decides
*what* to build, the logger decides *how*; the mapping between them lives in
`cmd/server`, which is the only place allowed to know both.

### Three levels of test

Which level a scenario belongs at is decided during PLAN, not guessed. Pushing
every scenario to the acceptance level is the most common way a Cucumber suite
becomes slow, flaky, and eventually switched off.

| Level | Where | Command |
| --- | --- | --- |
| Unit | `*_test.go` beside the code | `make test` |
| Integration | `*_test.go` with `//go:build integration` | `make test-integration` |
| Acceptance | `features/*.feature` + `test/acceptance/` | `make test-acceptance` |

### Two godog defaults are deliberately inverted

- **`Strict: true`** — by default godog treats undefined steps as TODOs and
  still passes, which would let an unimplemented scenario report green and break
  the outer loop of outside-in development.
- **`Concurrency: 4` with `-race`** — scenarios must be independent, so they run
  in parallel and any coupling surfaces as a failure. Per-scenario state lives in
  `context.Context`, which is what makes this safe.

### Code generation

Every `//go:generate` directive lives in the build-tagged `generate.go` at the
root, so all paths share one frame of reference and repo-scoped generators
(mockery today, buf or sqlc later) have somewhere to belong. Generators are
pinned as `tool` dependencies in `go.mod`:

```bash
go tool                 # list them
go get tool             # upgrade all
go build -o bin/ tool   # compile all, for CI caching
```

`apigen/` is fully disposable, which doubles as the cheapest check that the
pipeline still works:

```bash
rm -rf internal/interfaces/http/apigen && make gen && make test
```

### Code quality

`golangci-lint` v2, pinned as a tool dependency so CI and every developer run
the same version. Config in `.golangci.yml`; generated code and mocks are
excluded because findings there cannot be acted on without editing a file marked
DO NOT EDIT.

---

## 9. Future Considerations / Roadmap

**The empty layers are the roadmap, and they are empty by rule.**

`domain/`, `application/` and `infrastructure/` contain only `doc.go` files
stating their dependency rules. No business code may be written until a CLARIFY
run has produced an example map: every type in `domain/` must trace back to a
rule, and every scenario in `features/` back to a concrete example.

Writing the model first would defeat the point of the testbed — it would prove
the skills work by handing them the answer.

Planned, in order:

1. Run `bdd-clarify` against `prompts/1-fitness-tracker-clarify.md`; produce the first example maps
2. SPEC turns those examples into `.feature` files
3. PLAN assigns each scenario a test level
4. IMPLEMENT fills `domain/` and `application/` outside-in

Known gaps that are decisions, not oversights: no persistence, no auth, no
second protocol. Each arrives when a scenario requires it.

---

## 10. Project Identification

| | |
| --- | --- |
| Project name | skeleton — the Go lab for the `ai-bdd` plugin |
| Repository | none — lives inside the `ai-bdd` repo at `lab/go/skeleton` |
| Module path | `skeleton` |
| Primary contact | Benny.XF.Cheng |
| Last updated | 2026-08-17 |

---

## 11. Glossary

### BDD artefacts

| Term | Meaning |
| --- | --- |
| **Story** | One independently deliverable behaviour. One story, one example map |
| **Rule** | A constraint that holds on its own. Rules are the natural fault lines for splitting a story |
| **Example** | A concrete case illustrating a rule, with real data. `Example 1.1` is Rule 1's first |
| **Question** | Something nobody could answer during discovery — the evidence that a story is not ready |
| **Shared** | A rule that constrains several stories; lives in `docs/bdd/glossary.md`, referenced not copied |
| **Example map** | The output of CLARIFY: story, rules, examples, questions, readiness |

Numbers in these identifiers are **references**. Downstream Gherkin scenarios tag
themselves `@Example1.1`, so renumbering silently repoints them.

### Architectural terms

| Term | Meaning |
| --- | --- |
| **Aggregate** | A consistency boundary in `domain/`. Saved or rejected as a whole |
| **Driving port** (`port/in`) | The application's published API. Adapters depend on it |
| **Driven port** (`port/out`) | What the application needs from outside. Infrastructure implements it |
| **Command / Query** | Write and read use cases. Reads do not go through the aggregate |
| **Mapper** | Protocol request → command/query. Mechanical, in the adapter |
| **Assembler** | Domain → dto. In application, and rare — see `docs/DATAFLOW.md` |
| **Presenter** | Dto → protocol response. Mechanical, in the adapter |
| **Walking skeleton** | `/version` — the thinnest path exercising the whole chain |

### Acronyms

| | |
| --- | --- |
| **BDD** | Behaviour-Driven Development |
| **CQRS** | Command Query Responsibility Segregation |
| **DTO** | Data Transfer Object |
| **RFC 9457** | Problem Details for HTTP APIs, the error format used throughout |
