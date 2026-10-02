# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

A Claude Code **plugin** (`ai-bdd`) whose product is the `skills/` directory. Everything
else either explains why those skills look the way they do (`docs/`, `PLAN.md`) or exists
to make them fail visibly when they are weak (`benchmark/`).

The plugin's subject is a seven-step pipeline:

**CAPTURE → DISCOVERY → SPEC → PLAN → IMPLEMENT → VERIFY → REVIEW**

Only the first four steps exist as skills, and SPEC is served by two of them —
`bdd-spec` and `bdd-formulation`, BDD's own Formulation practice. `docs/ai-sdlc.md`
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
python3 skills/bdd-discovery/scripts/status.py [root]   # red cards and probe coverage per feature, read from prd.md
```

`check_spec.py` belongs to `bdd-formulation`:

```bash
python3 skills/bdd-formulation/scripts/check_spec.py [root]   # prd.md shape incl. source markers; adds .feature ↔ prd.md coverage once features exist
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
| `docs/tdd.md` | Canon TDD, the three laws, the two schools — BDD's parent, and the inner loop of `bdd.md`'s double-loop figure | ″ |
| `docs/sdd.md` | Spec-Driven Development per Spec Kit / Kiro / Tessl, and why its "spec" is not BDD's | ″ |
| `docs/ai-sdlc.md` | This project's position: why six steps, which parts are not BDD, the AI-specific constraints | Evolves with real runs |

The first three transcribe outside sources and say so in their own opening paragraph. Each
figure has exactly one owner: `tdd.md` points at `bdd.md` for the double-loop diagram
rather than redrawing it.

`PLAN.md` is the roadmap and is explicitly labelled a hypothesis, not a decision.

### Artifact format is the interface

Each step reads what the previous step actually emitted: DISCOVERY reads CAPTURE's
draft `prd.md`, SPEC reads the `已對焦` `prd.md`, PLAN reads `spec.md` and the `.feature`
files, REVIEW audits the whole chain.

```
需求方的話（prd.md ## 原文）→ Story → FR → AC → Scenario → Ticket → Step → 程式碼
```

Every arrow is a place intent can be lost. Changing an artifact's format breaks the chain,
so formats are settled before downstream work starts.

### The three modes, four disciplines

Each of the first four steps holds because every skill under it runs in exactly one
working mode. SPEC has two skills, so there are four disciplines:

- **CAPTURE** is scribe mode: the originator speaks, the agent writes `prd.md` as a
  draft. It asks at most four analyst questions and **never fills in rules or examples**.
- **DISCOVERY** is interrogation mode: the agent asks, the originator answers, and each
  answer is written on the spot as an FR/AC under its story with a `← Q-<n>` source.
  **An FR or AC without `← PRD §n` or `← Q-<n>` is a defect** — `check_spec.py` check 14
  rejects it. This replaces the older "CLARIFY does not write" rule, which could not be
  verified.
- **SPEC** is writing mode: `bdd-formulation` turns an `已對焦` `prd.md` into `.feature`
  files; `bdd-spec` writes the technical `spec.md` — **neither interviews**, and holes go
  into `spec.md`'s `## 回 DISCOVERY 補問`.
- **PLAN** slices tracer bullets — it does not design.

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
trace back to a rule in `prd.md`. `internal/domain/` currently holds only `doc.go` —
the `/version` walking skeleton is the sole implemented slice and carries no domain meaning.

## Conventions

- **Prose documents are written in Traditional Chinese; commit messages and PR text are
  English.** Commits follow conventional-commit titles with a `WHAT:` / `WHY:` / `HOW:` body.
  Scopes seen in history: `lab`, `bdd-spec`, `bdd-clarify`.
- Documentation artifacts the skills produce go to the consuming repo's
  `specs/<date>-<feature>/` (`prd.md`, `spec.md`, `plan.md`), one directory that can be deleted whole. Two
  exceptions: `.feature` files stay where the runner expects them (`features/` for
  Cucumber-family tools); `docs/CONTEXT.md` outlives any single feature, so SPEC may create
  it or append entries but must not rewrite its existing sections or write to anything else
  under `docs/`.
- Renames use `git mv` so `git log --follow` keeps working; the testbed's history depends on it.
