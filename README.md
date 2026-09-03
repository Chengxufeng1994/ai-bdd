# ai-bdd

A Claude Code plugin for AI-assisted BDD (Behavior-Driven Development): turning
requirements into Gherkin feature files, generating step definitions, running BDD
suites, and keeping scenarios in sync with the code they describe.

> **Status: early.** CLARIFY, SPEC and PLAN are implemented; IMPLEMENT,
> VERIFY and REVIEW are not. The division of labour between the first three
> was re-cut on 2026-09-02 and the skills have not caught up yet — see
> [PLAN.md](PLAN.md).

The workflow is six steps, with BDD as the method inside each:

**CLARIFY → SPEC → PLAN → IMPLEMENT → VERIFY → REVIEW**

## Layout

```
ai-bdd/
├── .claude-plugin/
│   └── plugin.json    # manifest: name, description, version, author
├── PLAN.md            # what to build, in what order
├── docs/
│   ├── bdd-workflow.md # what BDD itself is — external reference
│   └── ai-sdlc.md      # why this pipeline looks the way it does
├── skills/
│   ├── bdd-clarify/
│   ├── bdd-spec/
│   ├── bdd-plan/
│   ├── clarify-loop/
│   ├── story-splitting/
│   └── skill-rules/
├── lab/
│   └── go/
│       └── skeleton/     # Go + godog dogfooding ground
└── README.md
```

Skill directories are flat; the step grouping lives in the roadmap's table,
not in the filesystem — skill names are global and flat at runtime, and
`claude plugin validate` does not recurse into nested directories.

Capabilities live in conventional directories at the repo root, and are picked up
automatically once they exist:

| Directory   | Holds                                                    |
| ----------- | -------------------------------------------------------- |
| `skills/`   | `<skill-name>/SKILL.md` — model-invoked workflows         |
| `commands/` | `<name>.md` — user-invoked `/slash` commands              |
| `agents/`   | `<name>.md` — specialised subagents                       |
| `hooks/`    | `hooks.json` — lifecycle hooks                            |

A skill's directory name is what identifies it, and its frontmatter `description`
is what decides when it gets triggered — both are worth getting right.

None of these are created yet; add one when there is a first real capability to
put in it.

## Validate

```bash
claude plugin validate . --strict
```

`--strict` fails on unrecognised fields and missing metadata that the runtime
would otherwise tolerate, which makes it suitable as a CI gate.

## Install

This repo is a single plugin, not a marketplace, so it is installed by reference
from a marketplace entry:

```json
{
  "name": "ai-bdd",
  "source": { "source": "github", "repo": "<owner>/ai-bdd" }
}
```

Then:

```bash
claude plugin install ai-bdd@<marketplace-name>
```
