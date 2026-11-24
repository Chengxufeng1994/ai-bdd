# Skill Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Every ai-bdd skill reads as numbered steps that each state their output, with a `## 產物` table matching the duty card.

**Architecture:** Markdown skills plus one Python audit script. Rules S17/S18 land in `skill-rules` with mechanical checks first; `implement` is the pilot the user approves; then one skill per task, each checked by a fresh reviewer for lost rules.

**Tech Stack:** Markdown; `skills/skill-rules/scripts/audit_skill.py`; stdlib `unittest`; `claude plugin validate`.

**Spec:** `docs/superpowers/specs/2026-11-24-skill-cleanup-design.md`

## Global Constraints

- Prose Traditional Chinese; commits English, conventional title, `WHAT:`／`WHY:`／`HOW:` body, last line `Claude-Session: https://claude.ai/code/session_01Qkb1m77tBhqzve8Hn4TbBP`; work on `main`, push after each task.
- No rule's content changes — only its wording and position. A rule that moves to a reference file stays pointed to from SKILL.md (S5).
- Step shape: `### N. <動詞>` then a first line `產出：…`; rules one line each, `——` then a one-line why; no multi-paragraph rationale.
- `## 產物` is a table `| 檔 | 哪一步寫 | 格式在哪 |`; every row's file also appears in the duty card's `寫`, and every path in `寫` has a row.
- Untouched: `api-and-interface-design`, `incremental-implementation`, `test-driven-development`, `context-engineering`, `source-driven-development`, `doubt-driven-development`, `archive`; nothing under `runs/`.
- References are edited only to delete passages that duplicate SKILL.md.
- Section order stays as `rules/structure.md`〈章節順序〉; rule numbers S1–S16 are not renumbered.

## Review Focus

- **Rule dropped in the rewrite:** a `MUST`／`NEVER` or a Skill Boundaries row in the old SKILL.md absent from the new SKILL.md and its references — expected: reviewer lists it, it is restored before commit. Every skill task, step "lost-rule review".
- **Step with no file output** (asks the user, reads only): expected `產出：給使用者的…` or `產出：無，只讀`, not a missing line. Task 1 test.
- **`寫` lists a path the 產物 table lacks:** expected S18 warning. Task 1 test.
- **Un-numbered subsections under 流程** (`### 續跑`, `### 不越界`): expected exempt from S17. Task 1 test.
- **Cross-skill pointers** (`implement 的 §6`, `plan 的 Step 4`) break when steps are renumbered — expected: `grep -rn` for the old step names in other skills after each task, fixed in the same commit.

---

### Task 1: S17 and S18 in skill-rules

**Files:**
- Modify: `skills/skill-rules/rules/structure.md` (S17, S18 rows after S16 in the SKILL.md 章節 table; a `## S17／S18` section with the D1 example)
- Modify: `skills/skill-rules/SKILL.md` (`## 產物` row becomes 必要 for skills with `## 流程`; mention S17)
- Modify: `skills/skill-rules/scripts/audit_skill.py`
- Test: `tests/test_audit_skill.py`

**Interfaces:**
- Produces: audit emits `S17` (SHOULD) per numbered step under `## 流程` whose first non-blank line does not start with `產出：`, naming the step; `S18` (SHOULD) when `## 流程` exists and `## 產物` lacks the three-column table, or a path in a table row's first cell (backticked) is not in the card's `寫` cell, or a backticked path in `寫` has no row. Both skip skills without `## 流程`.

- [ ] **Step 1: Write failing tests** in a new `class StepOutputs(unittest.TestCase)` using `rules_hit_any`:
  - `test_numbered_step_without_output_is_flagged`: `## 流程\n### 1. 讀\n- x\n` → `"S17" in rules`
  - `test_step_with_output_passes`: `### 1. 讀\n產出：無，只讀\n` → `"S17" not in`
  - `test_unnumbered_subsection_is_exempt`: `### 續跑\n- x\n` → `"S17" not in`
  - `test_no_flow_section_skips_both`: body without `## 流程` → neither S17 nor S18
  - `test_product_row_missing_from_card_is_flagged`: card `寫` = `` `a.md` ``, table rows `` `a.md` `` and `` `b.md` `` → `"S18" in`
  - `test_card_path_missing_from_table_is_flagged`: card `寫` = `` `a.md`、`b.md` ``, table row only `` `a.md` `` → `"S18" in`
  - `test_matching_card_and_table_pass` → `"S18" not in`
- [ ] **Step 2:** `python3 -m unittest tests.test_audit_skill -v` → the new tests FAIL
- [ ] **Step 3:** Implement the two checks in `audit()`, after S16, reusing `preamble` for the card; paths = backticked tokens.
- [ ] **Step 4:** `python3 -m unittest discover -s tests` → OK
- [ ] **Step 5:** Edit structure.md and skill-rules SKILL.md as listed; `for d in skills/*/; do python3 skills/skill-rules/scripts/audit_skill.py $d | tail -1; done` → still `✔ 通過` on every skill (S17/S18 are warnings; the ones they list are the rollout's to-do)
- [ ] **Step 6:** Commit `feat(skill-rules): require step outputs and a product table (S17, S18)`; push.

### Task 2: Pilot — `implement`

**Files:**
- Modify: `skills/implement/SKILL.md`

- [ ] **Step 1:** Save the old version: `git show HEAD:skills/implement/SKILL.md > <scratchpad>/old-implement.md`.
- [ ] **Step 2:** Rewrite per Global Constraints: each `### N.` gets `產出：`; merge `### 不越界` into the steps or the card's `不做`; collapse rationale to one line per rule; `## 產物` becomes the table.
- [ ] **Step 3:** `audit_skill.py skills/implement` → no S17／S18; `wc -l` reported before → after.
- [ ] **Step 4: Lost-rule review.** Dispatch a fresh `general-purpose` subagent (model sonnet) with only the two file paths: "List every rule, boundary, or required action in OLD that NEW (plus the reference files it points to) does not state. Quote OLD's line. Output an empty list if none." Restore each hit.
- [ ] **Step 5:** `grep -rn "implement.*§\|implement 的" skills` → fix any pointer to a renamed step.
- [ ] **Step 6:** Full gate (unittest, audit all, `claude plugin validate .`); commit `refactor(implement): one output per step, product table`; push.
- [ ] **Step 7: STOP.** Show the user before → after of one step and the 產物 table; continue only on approval. Any change to the shape goes back into Task 1's rules first.

### Tasks 3–12: one skill each, same steps as Task 2 (minus Step 7)

In this order (pipeline tail first, so pointers into already-cleaned skills stay stable):

| Task | Skill |
| --- | --- |
| 3 | `behavior-driven-development` |
| 4 | `plan` |
| 5 | `formulation` |
| 6 | `spec` |
| 7 | `discovery` |
| 8 | `spec-by-example` |
| 9 | `example-mapping` |
| 10 | `feature-mapping` |
| 11 | `story-splitting` |
| 12 | `skill-rules` |

Each commits `refactor(<skill>): one output per step, product table` and pushes. After Task 12: the audit loop shows no S17／S18 on any of the eleven, and `README.md`'s product table is re-checked against the cards (S16 says they must agree).
