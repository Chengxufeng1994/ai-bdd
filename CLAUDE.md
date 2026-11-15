# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

A Claude Code **plugin** (`ai-bdd`) whose product is the `skills/` directory. Everything
else either explains why those skills look the way they do (`docs/`, `PLAN.md`) or exists
to make them fail visibly when they are weak (`benchmark/`).

The plugin's subject is a six-step pipeline:

**DISCOVERY → SPEC → PLAN → IMPLEMENT → REVIEW → VERIFY**

(brainstorming → to spec → to plan → implement → review → verify.) Only the first four
steps exist as skills; REVIEW and VERIFY are planned. PLAN is served by two skills —
`bdd-plan` (technical design, `plan.md`) and `bdd-tasks` (slices that design into one
file per ticket under `tasks/`, so tickets can run in parallel). IMPLEMENT is served by
`bdd-implement` (one ticket end to end: test layer red first, gate, merge-back) and `bdd`
(green and refactor, delegated by `bdd-implement`; `disable-model-invocation`, so only
that delegation or a manual `/bdd` runs it). SPEC is served by three skills —
`bdd-spec` (writes `spec.md` from the recorded answers), `bdd-spec-by-example` (writes
business-language acceptance journeys under `acceptance/`, not executed) and
`bdd-formulation`, BDD's own Formulation practice. `docs/ai-sdlc.md`
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

`status.py` belongs to `bdd-discovery`:

```bash
python3 skills/bdd-discovery/scripts/status.py [root]   # red cards and probe coverage per feature, read from discovery.md
```

`check_spec.py` belongs to `bdd-formulation`:

```bash
python3 skills/bdd-formulation/scripts/check_spec.py [root]   # spec.md shape + source markers resolved against discovery.md; checks acceptance journeys' @covers- tags; adds .feature ↔ spec.md coverage once features exist
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
| `docs/tdd.md` | Canon TDD, the three laws, the two schools — BDD's parent, and the inner loop of `bdd.md`'s double-loop figure | ″ |
| `docs/sdd.md` | Spec-Driven Development per Spec Kit / Kiro / Tessl, and why its "spec" is not BDD's | ″ |
| `docs/ai-sdlc.md` | This project's position: why six steps, which parts are not BDD, the AI-specific constraints | Evolves with real runs |

The first four transcribe outside sources and say so in their own opening paragraph. Each
figure has exactly one owner: `tdd.md` points at `bdd.md` for the double-loop diagram
rather than redrawing it.

`PLAN.md` is the roadmap and is explicitly labelled a hypothesis, not a decision.

`RESOURCES.md` lists every outside source: the articles `docs/` transcribes and the skill
projects individual skills borrowed from. Add a row there when a new reference is used.

### Artifact format is the interface

Each step reads what the previous step actually emitted: DISCOVERY writes
`discovery.md` (原文, rough stories, questions and answers — no rules), SPEC writes
`spec.md` from it and the `.feature` files once `已對焦`, PLAN writes `plan.md` (technical
design) from `spec.md` and the `.feature` files, then `tasks/` (one file per ticket plus
`index.md`) from `plan.md`, REVIEW audits the whole chain.

```
需求方的話（discovery.md ## 原文）→ Story → FR → AC → Scenario → Ticket → Step → 程式碼
```

Every arrow is a place intent can be lost. Changing an artifact's format breaks the chain,
so formats are settled before downstream work starts.

### The three modes, three disciplines

Each of the first three steps holds because it runs in one main working mode:

- **DISCOVERY** is brainstorming: it opens in scribe mode (the originator speaks, the
  agent records `## 原文` verbatim and rough stories, asks at most four analyst questions)
  and then interrogates. Every answer is recorded on the spot in `discovery.md` with the
  originator's words and who answered. **It writes no FR/AC.**
- **SPEC** is writing mode: `bdd-spec` turns an `已收斂` `discovery.md` into `spec.md`,
  `bdd-spec-by-example` writes acceptance journeys the originator reads, `example-mapping`
  takes it to `已對焦`, then `bdd-formulation` writes `.feature` files —
  **none of them interviews**; holes go into `spec.md`'s `## 回 DISCOVERY 補問` and send
  `discovery.md` back to `澄清中`. **An FR or AC without `← PRD §n` or `← Q-<n>` is a
  defect**, and a `← Q-<n>` must name an answered row in `discovery.md` — `check_spec.py`
  check 14 rejects both.
- **PLAN** owns the technical design (seam, modules, contracts, schema, test levels) —
  `bdd-plan` **does not change requirements**; holes go to `## 缺口`. `bdd-tasks` then
  slices tracer bullets and **does not change the design**; its holes go to the same
  `## 缺口`. `bdd-implement` writes only its own ticket file, never `plan.md` or
  `tasks/index.md` — that is what lets tickets run in parallel. A ticket done in a worktree
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
  `specs/<date>-<feature>/` (`discovery.md`, `spec.md`, `splitting.md`, `acceptance/*.feature`, `plan.md`, `tasks/`), one directory that can be deleted whole. Two
  exceptions: executable `.feature` files stay where the runner expects them (`features/` for
  Cucumber-family tools); `docs/CONTEXT.md` outlives any single feature, so SPEC may create
  it or append entries but must not rewrite its existing sections or write to anything else
  under `docs/`.
- Renames use `git mv` so `git log --follow` keeps working; the testbed's history depends on it.
