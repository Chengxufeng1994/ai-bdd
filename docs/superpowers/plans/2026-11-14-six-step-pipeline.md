# Six-Step Pipeline Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** discovery writes only `discovery.md` (原文 + 粗切 story + 問答); `bdd-spec` writes `spec.md` from it and runs readiness before `.feature`; `design.md` disappears into `plan.md`.

**Architecture:** Scripts first (TDD), because every skill's wording depends on which file holds which section. Then the skills in pipeline order — discovery, spec, plan/implement — each task fixing the cross-skill links it breaks so every commit validates. Docs and the on-disk parking run last.

**Tech Stack:** Markdown skills, Python 3 stdlib scripts + `unittest`, `claude plugin validate`, `skills/skill-rules/scripts/audit_skill.py`.

**Spec:** `docs/superpowers/specs/2026-11-14-six-step-pipeline-design.md`

## Global Constraints

- Prose 繁中; commits conventional title + `WHAT:` / `WHY:` / `HOW:` body, ending with `Claude-Session: https://claude.ai/code/session_01Qkb1mBWbYa6q2KffkJjgeT`.
- Moves use `git mv`.
- Every commit passes `claude plugin validate .`, `python3 -m unittest discover -s tests`, and `audit_skill.py` on every skill it touched.
- Source markers stay exactly `← PRD §n` and `← Q-<n>`; `PRD §n` now resolves in `discovery.md` `## 原文`, `Q-<n>` in `discovery.md` `## Open Questions`.
- Status words: `discovery.md` uses `澄清中` / `已收斂`; `spec.md` uses `待對焦` / `已對焦`. Neither file uses the other's words.
- Not built: `bdd-review`, `bdd-verify`. Not edited: `docs/superpowers/**` history, `docs/bdd.md`, `docs/tdd.md`, `docs/sdd.md`, `step-grammar.md`, `technical-probes.md` content, the sixteen checks' logic.

## Review Focus

1. A `spec.md` whose `← Q-7` has no row in `discovery.md`, or points at a `待答` row — check_spec must fail (Task 1 tests).
2. A feature directory with `discovery.md` but no `spec.md` yet — check_spec prints 「還沒寫 spec，先跑 bdd-spec」 and is not an error (Task 1 test).
3. The 11-13 layout (`spec.md` holding `## Open Questions` / `## 原文`, no `discovery.md`) — both scripts exit 1 with the split hint, never parse it as the new `spec.md` (Task 1 tests).
4. A resumed `bdd-discovery` after a hand-back — it must read `spec.md` `## 回 DISCOVERY 補問` and `plan.md` `## 缺口` and open them as new `Q-<n>`, not re-ask answered ones (Task 2 grep + skill text).
5. A second `bdd-spec` run after a round trip — FR/AC numbers unchanged, new ones appended (Task 3 skill text; verified by reviewer, no script test).

---

### Task 1: Scripts read `discovery.md` + `spec.md`

**Files:**
- Modify: `skills/bdd-discovery/scripts/status.py`, `skills/bdd-formulation/scripts/check_spec.py`
- Test: `tests/test_status.py`, `tests/test_check_spec.py`

**Interfaces:**
- Produces (`check_spec.py`): `cited_qs(text: str) -> set[str]` — the `<n>` of every `← Q-<n>` on a FR/AC line under `## User Stories`. Existing `has_verbatim_section`, `missing_cited_paragraphs`, `qs_in_table` keep their signatures and are now called on `discovery.md` text (for `missing_cited_paragraphs`, pass `spec_text + "\n" + discovery_text` so citations come from `spec.md` and `## 原文` from `discovery.md`).
- Produces (both): `legacy_hint(specs_dir: Path) -> str | None` — returns the refusal message for the first old-layout directory, else `None`. Two layouts: `prd.md` present (10-02; keep today's hint); `spec.md` containing `## Open Questions` or `## 原文` with no `discovery.md` beside it (11-13) → 「`<dir>` 是 11-13 佈局：把 `spec.md` 的 `## Open Questions` 與 `## 原文` 搬到 `discovery.md`，`design.md` 的技術節併進 `plan.md`」. Duplicated in both scripts (they are standalone, like today's guard).

- [ ] **Step 1: Split the fixtures and write the failing tests**

In both test files split today's fixture into `SPEC` (Document Overview with `**狀態**：待對焦`, `## User Stories` unchanged) and `DISCOVERY` (`## Document Overview` `**狀態**：澄清中`, `## Open Questions` table unchanged, `## 原文`). Helper writes both into `specs/2026-10-02-parking/`.

`tests/test_status.py`:
- `test_counts_rows_from_discovery` — writes only `discovery.md`; same assertions as today's `test_counts_rows_from_spec`.
- `test_missing_discovery_exits_1_and_names_discovery` — only `plan.md` present; `"沒有任何 discovery.md"` and `"bdd-discovery"` in out.
- `test_1113_layout_exits_1_with_split_hint` — `spec.md` = today's combined fixture, no `discovery.md`; rc 1, `"11-13 佈局"` and `"discovery.md"` in out.
- keep both `prd.md` tests.

`tests/test_check_spec.py`:
- existing `CheckReadsSpec` tests write both files; assertions unchanged.
- `test_q_marker_missing_from_discovery_fails` — SPEC with `← Q-9`; rc 1, `"Q-9"` and `"discovery.md"` in out.
- `test_q_marker_pointing_at_unanswered_fails` — AC `← Q-2` (Q-2 is `待答`); rc 1, `"Q-2"` and `"已答"` in out.
- `test_discovery_only_dir_is_a_notice` — only `discovery.md`; rc 0, `"先跑 bdd-spec"` in out.
- `test_1113_layout_exits_1_with_split_hint` — as above for `check`.
- `test_missing_verbatim_section_reported` now drops `## 原文` from DISCOVERY.
- `test_shipped_example_spec_passes_spec_only_checks` — loads `skills/bdd-spec/examples/minimal-spec.md` and `skills/bdd-discovery/examples/minimal-discovery.md` (created in Tasks 2–3); mark `@unittest.skip("examples land in Task 3")` here and remove the skip in Task 3.

- [ ] **Step 2: Run to verify RED**

Run: `python3 -m unittest discover -s tests 2>&1 | grep -E "^(FAIL|ERROR|Ran|OK|FAILED)"`
Expected: the new and split tests FAIL (scripts still read `## Open Questions` / `## 原文` from `spec.md`).

- [ ] **Step 3: Implement**

`status.py.main`: `legacy_hint` → glob `*/discovery.md` → message 「`{specs}` 底下沒有任何 discovery.md —— 先跑 bdd-discovery」 → parse as today.

`check_spec.check`: move the header print below `legacy_hint` (fixes the deferred minor from 11-13). Iterate feature dirs: dir with `discovery.md` and no `spec.md` → print notice, skip, no problem. `spec.md` without `discovery.md` (and not 11-13) → ✗ 「沒有 discovery.md，`← Q-<n>` 與 `← PRD §n` 無處可查 —— 先跑 bdd-discovery」. Otherwise read both; checks 10/11 use `qs_in_table(dtext)`; 15/16 use `dtext`; after check 14 add the cross-file check from `cited_qs` (missing → ✗ 「`← Q-<n>` 在 discovery.md 沒有這一題」; not `已答` → ✗ 「`← Q-<n>` 指向還不是已答的題」). Rename every message string that sends the user to DISCOVERY for a FR-shaped problem (「開場的半成品」「DISCOVERY 開場還沒切 story」) to `bdd-spec`. Update the docstring's numbered list for items 10, 11, 14, 15, 16 to name `discovery.md`.

- [ ] **Step 4: Verify GREEN**

Run: `python3 -m unittest discover -s tests -v 2>&1 | tail -3`
Expected: OK (one skipped).

- [ ] **Step 5: Commit**

```bash
git add skills/bdd-discovery/scripts/status.py skills/bdd-formulation/scripts/check_spec.py tests/
git commit -m "feat(scripts): read questions and 原文 from discovery.md, rules from spec.md"
```

---

### Task 2: `bdd-discovery` writes only `discovery.md`

**Files:**
- Create: `skills/bdd-discovery/references/discovery-format.md`, `skills/bdd-discovery/examples/minimal-discovery.md`
- Rename: `skills/bdd-discovery/references/spec-format.md` → `skills/bdd-spec/references/spec-format.md`; `skills/bdd-discovery/examples/minimal-spec.md` → `skills/bdd-spec/examples/minimal-spec.md`
- Modify: `skills/bdd-discovery/SKILL.md`, `skills/bdd-discovery/references/opening.md`, `skills/bdd-discovery/references/persona-definition.md`, `skills/clarify/SKILL.md`, `skills/story-splitting/SKILL.md`, every file linking `../bdd-discovery/references/spec-format.md` (retarget to `../bdd-spec/references/spec-format.md`)

**Interfaces:**
- Consumes: Task 1's file split.
- Produces: `discovery-format.md` defines D2 exactly — sections `## Document Overview` / `## 粗切 story` / `## Open Questions` / `## 原文`; status `澄清中` → `已收斂`; `擋住：<story slug 或 §n>`; `歸屬` values = `spec.md` section names or `plan.md` technical section names (`Tech Preferences & Constraints` / `Known Boundaries` / `Implementation Decisions` / `Testing Decisions` / `Risks`). Later tasks link to it as `../bdd-discovery/references/discovery-format.md`.

- [ ] **Step 1: Move the spec format and example out**

```bash
mkdir -p skills/bdd-spec/examples
git mv skills/bdd-discovery/references/spec-format.md skills/bdd-spec/references/spec-format.md
git mv skills/bdd-discovery/examples/minimal-spec.md skills/bdd-spec/examples/minimal-spec.md
```
Retarget every link found by `grep -rn "bdd-discovery/references/spec-format\|references/spec-format.md" skills/`.

- [ ] **Step 2: Write `discovery-format.md` and `minimal-discovery.md`**

`discovery-format.md`: the D2 table, the Open Questions format section copied from `spec-format.md` (table + 小節 + `面向` 值域), and the `## 原文` rules. `minimal-discovery.md`: the `## Open Questions` and `## 原文` sections cut from `bdd-spec/examples/minimal-spec.md`, plus a `## Document Overview` (`已收斂`, a 「需求方確認建檔」 revision row) and a `## 粗切 story` section listing its stories as one-liners. Remove those two sections from `minimal-spec.md`; set its status to `已對焦`.

- [ ] **Step 3: Rewrite `bdd-discovery/SKILL.md` and `opening.md`**

- description / title / intro: produces `discovery.md`; no 「長成 FR／AC」; 收尾 `已收斂`, next `bdd-spec`.
- opening step 0 builds `discovery.md`: 原文, 粗切 story section, four questions as `Q-<n>` rows with `面向` `—`, sign-off row. Remove the instructions that write Background / Goal / Scope / Personas / Assumptions sections (to spec does that) and the User Stories template with empty FR/NFR/Q groups — replace with the one-line `## 粗切 story` form.
- 「記錄：三件事一起做」 table: drop row 2 (write FR/AC) and row 3 (remove `#### Q` pointer); keep 1, 4, 5. Drop the 「FR 與 AC 的行尾只准」 MUST and the 「第 2 項不可以留到最後一次做」 paragraph.
- 續跑: add — read `spec.md` `## 回 DISCOVERY 補問` and `plan.md` `## 缺口` in the same directory; each entry not yet a `Q-<n>` becomes a new `待答` row whose 小節 names its source file.
- 結束: remove the steps moved to `bdd-spec` (persona-used-by-FR, command/query split, check 14 at zero, 「藍卡太多」); remove the `example-mapping` 「逼問中」 call and the 「流程內」 hand-off; keep frontier-empty, 待答 zero, 稽核分類, sign-off. 完成後 收斂 → 「`已收斂`，下一步 `bdd-spec`」.
- 產物 block: `discovery.md  DISCOVERY ★`. MUST NOT write `spec.md`, `plan.md`, `.feature`.
- Skill Boundaries: 「要把答案寫成 FR／AC → `bdd-spec`」; 「要判就緒 → `bdd-spec`（它呼叫 `example-mapping`）」.

- [ ] **Step 4: Update `clarify`, `story-splitting`, `persona-definition.md`**

`clarify`: its file is `discovery.md`; the standalone three-section file is a `discovery.md` missing `## 粗切 story`. `story-splitting`: two call sites per D4 table (discovery edits `discovery.md` `## 粗切 story`; `bdd-spec` edits `spec.md` `## User Stories` before `已對焦` only). `persona-definition.md`: answers live in `discovery.md`, `bdd-spec` writes `## Personas`; point its 來源標記 link at `../../bdd-spec/references/spec-format.md` (closes the 11-13 deferred minor).

- [ ] **Step 5: Verify**

Run: `grep -n "FR-<n>\*\*\|長成.*FR\|已對焦\|待對焦\|spec\.md" skills/bdd-discovery/SKILL.md skills/bdd-discovery/references/opening.md`
Expected: only lines that name `spec.md` as something discovery **reads** (回頭清單) or **does not write**; no `已對焦` / `待對焦`.

Run: `for s in bdd-discovery clarify story-splitting bdd-spec example-mapping bdd-formulation bdd-plan bdd-implement; do python3 skills/skill-rules/scripts/audit_skill.py skills/$s >/dev/null 2>&1 || echo FAIL $s; done; claude plugin validate . 2>&1 | tail -1; python3 -m unittest discover -s tests 2>&1 | tail -1`
Expected: no FAIL, validate passes, OK.

- [ ] **Step 6: Commit**

```bash
git add -A skills
git commit -m "feat(bdd-discovery): write only discovery.md and leave rules to bdd-spec"
```

---

### Task 3: `bdd-spec` writes `spec.md` from `discovery.md`

**Files:**
- Modify: `skills/bdd-spec/SKILL.md` (rewrite), `skills/bdd-spec/references/spec-format.md`, `skills/bdd-spec/references/design-format.md` (shrink to the `domain-model.md` skeleton; rename to `domain-model-format.md` with `git mv`), `skills/example-mapping/SKILL.md`, `skills/bdd-formulation/SKILL.md`, `skills/bdd-formulation/references/state-tags.md`, `tests/test_check_spec.py` (remove the Task 1 skip)

**Interfaces:**
- Consumes: `discovery-format.md` (Task 2), `spec-format.md` now under `bdd-spec`.
- Produces: `bdd-spec` flow D4a; the technical sections of `design-format.md` are **cut out here and pasted verbatim into `skills/bdd-plan/references/plan-format.md` under a temporary heading `## 技術節（待 Task 4 整理）`** so no content is lost between commits.

- [ ] **Step 1: Remove the skip and run the shipped-example test**

Run: `python3 -m unittest tests.test_check_spec -k shipped_example -v`
Expected: PASS (both example files exist since Task 2). If it fails, fix the examples, not the test.

- [ ] **Step 2: Move the technical sections out of `design-format.md`**

Cut `## Tech Preferences & Constraints`, `## Known Boundaries`, `## Implementation Decisions`, `## Testing Decisions`, `## Risks`, the 「來源標記：三選一」 section and `## Document Overview` / 核心關係人 rules into `plan-format.md` as above. Move `## 回 DISCOVERY 補問` into `spec-format.md`. What remains is the `domain-model.md` skeleton → `git mv` to `domain-model-format.md`. Delete `skills/bdd-spec/examples/minimal-design.md` with `git rm` after pasting its technical sections under the same temporary heading in `skills/bdd-plan/references/plan-format.md`.

- [ ] **Step 3: Update `spec-format.md`**

Remove `## Open Questions` and `## 原文` sections (now in `discovery-format.md`; link to it). Status table: `待對焦`（`bdd-spec` 寫完）→ `已對焦`（`example-mapping`）. 「誰第一次寫」 column: every row `bdd-spec`. Source markers: `PRD §n` → `discovery.md` `## 原文`; `Q-<n>` → `discovery.md` `## Open Questions`, must be `已答`. Add `## 回 DISCOVERY 補問` to the section order.

- [ ] **Step 4: Rewrite `bdd-spec/SKILL.md`**

Description: 「讀已收斂的 `discovery.md`，寫 `spec.md`——story、FR／AC（行尾 `← Q-<n>`／`← PRD §n`）、角色、範圍——再判就緒、產 `.feature`；並把詞彙抬進 `docs/CONTEXT.md`、聚合增修進 `specs/domain-model.md`。不訪談。BDD 六步流程的 to spec 步驟。」 Body sections, in order:
1. 前置確認: accepts `discovery.md` `已收斂` or 結束方式 `喊停` (D4d); refuses `阻塞`／`零進展`; NEVER without `## 原文`.
2. 寫 `spec.md`: from 原文 and 已答 rows; every FR/AC line ends with a marker pointing at a 已答 row or a §n; 喊停 leaves 待答 rows as `#### Q` pointers.
3. 收尾檢查 (D4b list), with `story-splitting` allowed before `已對焦`.
4. `docs/CONTEXT.md` glossary and `specs/domain-model.md` — keep today's text.
5. `example-mapping` → 就緒問句 → `已對焦` → `bdd-formulation`.
6. 回頭 (D4c): set `discovery.md` to `澄清中` + revision row 「回頭：<原因>」; tell the user to run `bdd-discovery`.
7. 再跑時修補 (D4c): never renumber FR/AC; append new numbers; revision row per change; state back to `待對焦`.
8. 自檢: `check_spec.py`.
Remove every seam / Implementation / Testing instruction (Task 4 owns them).

- [ ] **Step 5: `example-mapping`, `bdd-formulation`, `state-tags.md`**

`example-mapping`: caller is `bdd-spec`; delete the 「流程內·逼問中」 row; status table rows `待對焦`（`bdd-spec`）/ `已對焦`（本 skill）; 「紅卡」 now means `discovery.md` 待答 rows reached through `#### Q`. `bdd-formulation`: `spec.md` 由 `bdd-spec` 維護; it is called by `bdd-spec` after `已對焦`; its own gate unchanged. `state-tags.md`: `@ready` row unchanged; `@done` 「`bdd-review` 通過後」 → 「`bdd-verify` 通過後」.

- [ ] **Step 6: Verify**

Run: `grep -rn "design\.md\|design-format\|minimal-design" skills/ | grep -v "11-13 佈局"`
Expected: only `bdd-plan` / `bdd-implement` hits (Task 4 removes them).

Run the audit / validate / unittest line from Task 2 Step 5.
Expected: no FAIL, validate passes, OK (no skips).

- [ ] **Step 7: Commit**

```bash
git add -A skills tests
git commit -m "feat(bdd-spec): write spec.md from discovery.md and run readiness before formulation"
```

---

### Task 4: Technical design lives in `plan.md`

**Files:**
- Modify: `skills/bdd-plan/SKILL.md`, `skills/bdd-plan/references/plan-format.md`, `skills/bdd-implement/SKILL.md`

**Interfaces:**
- Consumes: the temporary `## 技術節（待 Task 4 整理）` block in `plan-format.md` (Task 3).
- Produces: `plan.md` section order per D5; nothing under `skills/` mentions `design.md` except the 11-13 hint.

- [ ] **Step 1: Restructure `plan-format.md`**

Turn the temporary block into the D5 order: header (`**來源**`, `**seam**`), 核心關係人 table (copied from `discovery.md`), `## Tech Preferences & Constraints`, `## Known Boundaries`, `## Implementation Decisions`, `## Testing Decisions`, `## Risks`, then the existing `## 前置整理`, `## 票`, `## 缺口`. Source marker rule: `PRD §x` / `Q-<n>` / `推論`; every `推論` also listed in `## 缺口`. State that `plan.md`'s markers are not script-checked yet (known gap, `bdd-review`).

- [ ] **Step 2: `bdd-plan/SKILL.md`**

Inputs: `spec.md`（`已對焦`）, `.feature`, `discovery.md` (technical-dimension 已答 rows by `歸屬`, 核心關係人). Move in from the old `bdd-spec` text: 「決定 seam」 four-rule table, the one Ask-user to confirm seam, 「seam 一旦寫進就往下傳」. Discipline line: 「不設計」 → 「不改需求」 (may make technical decisions; never change FR/AC/story; holes → `## 缺口`, and if it is a requirements hole set `discovery.md` back to `澄清中` per D4c). Remove 「還沒有 `design.md` → `bdd-spec`」.

- [ ] **Step 3: `bdd-implement/SKILL.md`**

Every `design.md` read point → `plan.md` (`## Implementation Decisions` / `## Testing Decisions`); 「還沒有 design.md → bdd-spec」 → 「`plan.md` 沒有技術節 → `bdd-plan`」; holes go to `plan.md` `## 缺口` only.

- [ ] **Step 4: Verify**

Run: `grep -rn "design\.md\|design-format\|bdd-spec\` 的事\|回 \`bdd-spec\`.*seam" skills/ | grep -v "11-13 佈局"`
Expected: no output.

Run the audit / validate / unittest line from Task 2 Step 5.
Expected: no FAIL, validate passes, OK.

- [ ] **Step 5: Commit**

```bash
git add skills
git commit -m "feat(bdd-plan): own the technical design that used to live in design.md"
```

---

### Task 5: Docs and the parking run

**Files:**
- Modify: `docs/ai-sdlc.md`, `PLAN.md`, `CLAUDE.md`, `benchmark/README.md`, `benchmark/skeleton/go/README.md`, `benchmark/skeleton/go/docs/ARCHITECTURE.md`
- On disk only (gitignored): `runs/go-parking-billing/specs/2026-10-03-parking-billing/`

- [ ] **Step 1: Edit the documents**

Pipeline line everywhere: `DISCOVERY → SPEC → PLAN → IMPLEMENT → REVIEW → VERIFY` (REVIEW and VERIFY marked 未實作). `ai-sdlc.md`: step table rows for DISCOVERY (只記問答, `discovery.md`), SPEC (`spec.md` + readiness + `.feature`), PLAN (技術設計 + 票, 紀律「不改需求」); modes table; add a short 「為什麼 discovery 只記問答」 paragraph (B1 vs B2 from the spec); keep the CAPTURE history paragraph. `CLAUDE.md`: artifact list `discovery.md`, `spec.md`, `plan.md`; modes bullets; script descriptions (`status.py` reads `discovery.md`; `check_spec.py` reads both). `PLAN.md`: 六步流程 table, file tree, skill table (`bdd-spec` row rewritten; `design.md` gone), 進度 bullet for 2026-11-14. Leave history passages as written.

- [ ] **Step 2: Split the parking run**

In `runs/go-parking-billing/specs/2026-10-03-parking-billing/`: create `discovery.md` from `spec.md`'s `## Document Overview` (status `已收斂`), its story headlines as `## 粗切 story` one-liners, `## Open Questions`, `## 原文`; remove those two sections from `spec.md` and set its status `已對焦`; paste `design.md`'s technical sections at the top of `plan.md` in D5 order; delete `design.md`. Plain `mv`/editing — this directory is not under git.

- [ ] **Step 3: Verify**

Run: `python3 skills/bdd-discovery/scripts/status.py runs/go-parking-billing; echo rc=$?; python3 skills/bdd-formulation/scripts/check_spec.py runs/go-parking-billing | tail -2`
Expected: `rc=0`; `全部通過`.

Run: `grep -rn "design\.md" --include=*.md --include=*.py . | grep -v "docs/superpowers\|^./runs\|^./.git"`
Expected: only history passages and the 11-13 hint; list them in the commit body.

Run: `claude plugin validate . && python3 -m unittest discover -s tests`
Expected: pass.

- [ ] **Step 4: Commit**

```bash
git add docs/ai-sdlc.md PLAN.md CLAUDE.md benchmark/README.md benchmark/skeleton/go/README.md benchmark/skeleton/go/docs/ARCHITECTURE.md
git commit -m "docs: describe the six-step pipeline with discovery.md, spec.md and plan.md"
```
