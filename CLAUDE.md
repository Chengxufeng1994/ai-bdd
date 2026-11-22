# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

A Claude Code **plugin** (`benny-skills`) whose product is the `skills/` directory; its skills
currently form one set, `ai-bdd`, plus five addyosmani skills kept verbatim outside the pipeline
(`incremental-implementation`, `test-driven-development`, `context-engineering`,
`source-driven-development`, `doubt-driven-development`; manual `/<name>` only). Everything
else either explains why those skills look the way they do (`docs/`, `PLAN.md`) or exists
to make them fail visibly when they are weak (`benchmark/`).

The plugin's subject is a six-step pipeline:

**DISCOVERY → SPEC → PLAN → IMPLEMENT → REVIEW → VERIFY**

(brainstorming → to spec → to plan → implement → review → verify.) Only the first four
steps exist as skills; REVIEW and VERIFY are planned. PLAN is served by one skill, `plan`
(following addyosmani's planning flow: collect; design — the seam, then the `.feature` files delegated to
`formulation` (business language only, before any design), then Domain Model, Database Schema and the API
contract side by side, the contract delegated to `api-and-interface-design`, plus the technical
decisions with the seam confirmed by the user; dependency graph; vertical slices; write the tasks;
order and checkpoints; `plan.md` is written during the steps, and only `tasks.md` waits for the user's nod
on the slicing — one section per ticket, so tickets can run in parallel — and STOPs for
`/implement <ticket>`). IMPLEMENT is served by
`implement` (one ticket end to end: test layer red first, gate, merge-back) and `behavior-driven-development`
(green and refactor, delegated by `implement`; `disable-model-invocation`, so only
that delegation or a manual `/behavior-driven-development` runs it). When all of a feature's tickets are merged,
`archive` writes what it delivered into the project's current truth under `docs/` and moves
its directory from `specs/ongoing/` to `specs/sleeping/`. SPEC is served by `spec`, `feature-mapping`, `example-mapping` and
`spec` (writes `spec.md` from the recorded answers — what to build and what not to: context,
personas, goals and non-goals, prioritised stories with FR → AC, domain concepts, edge cases;
it stops once `已對焦` and writes no technical decisions) and `spec-by-example` (writes
business-language acceptance journeys under `acceptance/`, not executed). `formulation`, BDD's own
Formulation practice, is delegated by `plan`. `docs/ai-sdlc.md`
explains which steps map onto BDD's practices and which are this project's own additions.

## Commands

### Plugin

```bash
claude plugin validate .              # the gate
python3 skills/skill-rules/scripts/audit_skill.py skills/<name>   # mechanical skill audit
```

There is no CI configuration — `validate` and the testbed's `make verify` are run by hand.

**`--strict` cannot pass here, and that is a deliberate trade.** It fails on one
warning: a `CLAUDE.md` at the plugin root is not loaded when the plugin is *installed*,
so the validator treats shipping one as a mistake. This repo is both the plugin and the
place the plugin is developed, and the file earns its place in the second role. Use plain
`validate` as the gate; run `--strict` when you want to see everything, expecting that
one warning and nothing else.

### Testbed (`benchmark/skeleton/go/`)

```bash
make help              # authoritative target list
make verify            # everything CI should run: fmt, generated, stamp, lint, vet, test -race
make test              # unit + acceptance, no external deps
make test-acceptance   # .feature files only
make test-wip          # only scenarios tagged @wip
make gen               # every generator (oapi-codegen, mockery)
make lint / lint-fix
make steps             # list every registered godog step definition
```

Single Go test: `go test ./internal/interfaces/http/ -run TestVersion`.

`generate.go` holds every `//go:generate` directive for the module and carries
`//go:build generate`, so it is excluded from all builds — `go vet`/lint never see it.
To check it: `go vet -tags generate .`.

### Skill scripts

Both take a project root and default to `.`; both exit 1 and name what is missing rather
than inferring anything from zero input.

`status.py` belongs to `discovery`:

```bash
python3 skills/discovery/scripts/status.py [root]   # red cards and probe coverage per feature, read from interview-log.md's ## 問答
```

`check_spec.py` belongs to `formulation`:

```bash
python3 skills/formulation/scripts/check_spec.py [root]   # spec.md shape + source markers resolved against its own ## Open Questions and ## 原文; checks acceptance journeys' @covers- tags; adds .feature ↔ spec.md coverage once features exist
```

Both have stdlib unit tests:

```bash
python3 -m unittest discover -s tests
```

## Architecture

### The document contract

`docs/` holds two kinds of file with deliberately different change rules. Respect the
split — it exists because an earlier version blended them and nobody went back to check
the source.

| File | Contains | Changes when |
| --- | --- | --- |
| `docs/bdd.md` | Cucumber's own definitions of Discovery / Formulation / Automation | A source changed the definition, or we misread it. **Never** to record a project opinion. |
| `docs/example-mapping.md` | Matt Wynne's Example Mapping article — four cards, the table as feedback, the 25-minute box, rough examples | ″ |
| `docs/feature-mapping.md` | John Ferguson Smart's Feature Mapping article — example → steps → consequences rows, grouped by rule, breadth-first | ″ |
| `docs/tdd.md` | Canon TDD, the three laws, the two schools — BDD's parent, and the inner loop of `bdd.md`'s double-loop figure | ″ |
| `docs/sdd.md` | Spec-Driven Development per Spec Kit / Kiro / Tessl, and why its "spec" is not BDD's | ″ |
| `docs/ai-sdlc.md` | This project's position: why six steps, which parts are not BDD, the AI-specific constraints | Evolves with real runs |

The first five transcribe outside sources and say so in their own opening paragraph. Each
figure has exactly one owner: `tdd.md` points at `bdd.md` for the double-loop diagram
rather than redrawing it.

`PLAN.md` is the roadmap and is explicitly labelled a hypothesis, not a decision.

`RESOURCES.md` lists every outside source: the articles `docs/` transcribes and the skill
projects individual skills borrowed from. Add a row there when a new reference is used.

### Artifact format is the interface

Each step reads what the previous step actually emitted: DISCOVERY writes only an
append-only interview log `interview-log.md` (原文, sessions with stakeholders and rough stories,
one `### Q-<n>` per question, confirmed decisions, open items — no rules, no `spec.md`), SPEC is the
only writer of `spec.md` (copies 原文 verbatim, curates the log's questions into `## Open Questions`
under the same numbers, writes what to build up to `已對焦`) and the `.feature` files, PLAN writes
`plan.md` (system interfaces, technical decisions, the three designs, then the execution plan) from `spec.md` and the `.feature`
files, then `tasks.md` (one section per ticket, in dependency order) from both, REVIEW audits the
whole chain.

```
需求方的話（interview-log.md → spec.md ## 原文）→ Story → FR → AC → Scenario → Ticket → Step → 程式碼
```

Every arrow is a place intent can be lost. Changing an artifact's format breaks the chain,
so formats are settled before downstream work starts.

### The three modes, three disciplines

Each of the first three steps holds because it runs in one main working mode:

- **DISCOVERY** is brainstorming: it opens in scribe mode (the originator speaks, the
  agent records `## 原文` verbatim and rough stories, asks at most four analyst questions)
  and then interrogates one question at a time, each with context and a guessed answer,
  until it can predict the originator's reaction and the coverage check passes; it ends
  with an intent restatement the originator must answer yes to. Every answer is recorded on
  the spot in `interview-log.md` with the originator's words and who answered.
  **It writes no `spec.md` and no FR/AC.** It is the only skill that asks the originator questions.
- **SPEC** is writing mode: `spec` reads the log, copies its 原文 and questions into `spec.md`
  (the log is evidence and never edited; on a mismatch its wording wins) and turns the answered ones into FR/AC,
  `feature-mapping` (whole feature, breadth) and `example-mapping` (one story, depth)
  diagnose it, `spec-by-example` writes acceptance journeys the originator reads, `example-mapping`
  takes it to `已對焦` —
  **none of them interviews**; a hole is opened as a new `待答` `Q-<n>` in `## Open Questions`,
  the status goes back to `澄清中`, and the user runs `/discovery`, which reads those questions, appends them to
  the log under the same number and asks them. **An FR or AC without `← PRD §n` or `← Q-<n>` is a
  defect**, and a `← Q-<n>` must name an answered row in `## Open Questions` — `check_spec.py`
  check 14 rejects both. No skill calls the next step: every step is a user-typed `/<skill>`; a skill
  may only DELEGATE helpers inside its own step.
- **PLAN** owns the `.feature` files (delegated to `formulation` right after the seam, before any
  design, so scenarios stay in business language), the technical decisions (seam, modules, test levels), the three designs
  (Domain Model with aggregates and invariants, Database Schema, API contract as an OpenAPI
  fragment — aggregates live in `plan.md`, not `spec.md`, whose `## Domain Concepts` holds only
  the glossary) and the execution plan (risks, parallel work, join points) — `plan` **does not change
  `spec.md`, and touches `.feature` files only through `formulation`**; the seam is the one thing confirmed with the user,
  technical lines may be `← 推論` or `← 專案 <path>`, and holes go to `## 缺口` and back to
  `spec` (a scenario not observable at the seam is first re-delegated to `formulation` for a rewording). `plan` also
  slices tracer bullets into `tasks.md`; a rerun redoes only tickets not yet started. `implement` writes only its own ticket's section of `tasks.md`, never `plan.md` or
  another ticket's section (it copies its operations from `plan.md`'s API contract into `api/openapi.yaml` unchanged) — that is what lets tickets run in parallel. A ticket done in a worktree
  is merged back by the session that did it: generated files are regenerated, append-only
  conflicts are unioned, and the gate re-runs on the merged code.

This matters more for an agent than for a person: a human who finds a hole while writing a
spec stops and asks; an agent supplies a plausible answer and continues.

### Skills layout

Skill directories are flat under `skills/`. Skill names are global and flat at runtime, and
`claude plugin validate` does not recurse into nested directories, so the step grouping
lives in `PLAN.md`'s table rather than in the filesystem.

`skills/skill-rules/` is the house style for authoring and auditing skills — read it before
adding or editing any skill.

### The testbed

`benchmark/skeleton/go/` is a Go + godog project used to dogfood the pipeline: layered
architecture (cmd / domain / application / infrastructure / interfaces), three test levels,
and `api/openapi.yaml`. `benchmark/` is not a plugin convention directory, so nothing in it is
loaded and it cannot affect plugin behaviour.

Two things about it that look like breakage but are not:

- **`go test ./...` fails on purpose.** The acceptance suite is red while scenarios are
  specified but unimplemented. `make verify` is the green gate, not `go test ./...`.
- **The scripts exit 1 against it** when they find no artifacts. That is fail-loud
  behaviour: a dashboard that quietly reports "all fine" gets believed.

Its ground rule: **no business code until DISCOVERY has run.** Every type in `domain/` must
trace back to a rule in `spec.md`. `internal/domain/` currently holds only `doc.go` —
the `/version` walking skeleton is the sole implemented slice and carries no domain meaning.

## Conventions

- **Prose documents are written in Traditional Chinese; commit messages and PR text are
  English.** Commits follow conventional-commit titles with a `WHAT:` / `WHY:` / `HOW:` body.
  Scopes seen in history: `lab`, `bdd-spec`, `bdd-clarify`.
- Documentation artifacts the skills produce go to the consuming repo's
  `specs/ongoing/<date>-<feature>/` (`spec.md`, `splitting.md`, `acceptance/*.feature`, `plan.md`, `tasks.md`), one directory that can be deleted whole. When all its
  tickets are merged, `archive` `git mv`s that directory to `specs/sleeping/`; `discovery`
  moves it back to `ongoing/` (only when the user says so) to correct it. Exceptions:
  executable `.feature` files stay where the runner expects them (`features/` for
  Cucumber-family tools); the project's current truth lives in `docs/` and only `archive`
  writes it, once a feature is done — it may create `docs/CONTEXT.md` or append entries but
  must not rewrite its existing sections, and may create `docs/DOMAIN-MODEL.md`,
  `docs/DATABASE-SCHEMA.md` and `docs/TECH.md` or change only the `<!-- ai-bdd:begin … -->`
  blocks inside them. Nothing else under `docs/` is written.
- Renames use `git mv` so `git log --follow` keeps working; the testbed's history depends on it.
