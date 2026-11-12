# Drop CAPTURE Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove the `bdd-capture` skill and the `prd.md` artifact; DISCOVERY builds and fills `specs/<date>-<feature>/spec.md`, SPEC writes `design.md`.

**Architecture:** Three renames and one merge. (1) The technical `spec.md` becomes `design.md` everywhere. (2) The business `prd.md` becomes `spec.md` in the two scripts, their tests, and every skill. (3) `bdd-capture`'s opening work (verbatim 原文, header sections, ≤4 questions, rough story cut, originator sign-off) becomes step 0 of `bdd-discovery`'s opening; its references/examples `git mv` into `bdd-discovery`. The rename order matters: `spec.md → design.md` first, `prd.md → spec.md` second, or the two strings collide.

**Tech Stack:** Markdown skills, Python 3 stdlib scripts + `unittest`, `claude plugin validate`, `skills/skill-rules/scripts/audit_skill.py`.

**Spec:** `docs/superpowers/specs/2026-11-13-drop-capture-design.md`

## Global Constraints

- Prose in Traditional Chinese; commit titles conventional-commit, body `WHAT:` / `WHY:` / `HOW:` (English).
- Moves use `git mv`; `git log --follow` must chain.
- Every commit passes `claude plugin validate .` and `python3 -m unittest discover -s tests`.
- Every touched skill passes `python3 skills/skill-rules/scripts/audit_skill.py skills/<name>`.
- Source markers stay exactly `← PRD §n` and `← Q-<n>`; `check_spec.py` check 14 unchanged.
- `docs/superpowers/specs/*`, `docs/superpowers/plans/*`, `docs/bdd.md`, `docs/tdd.md`, `docs/sdd.md`, `runs/`, and history passages in `PLAN.md` / `benchmark/cases/*/grading.md` are not edited. `benchmark/cases/personal-memo/input-prd.md` keeps its name.
- `bdd-discovery/SKILL.md` is 472 lines; if step 0 pushes it past the `skill-rules` length guidance, put step 0 in `references/opening.md` and link it.

## Review Focus

1. A repo in the old layout (`prd.md` + technical `spec.md` side by side) — both scripts must exit 1 naming the `git mv` fix, not parse the technical file as the business spec. (Task 2 tests.)
2. A `specs/<x>/` with only `design.md` — scripts say "沒有 spec.md —— 先跑 bdd-discovery", never mention `bdd-capture`. (Task 2 tests.)
3. The shipped example `bdd-discovery/examples/minimal-spec.md` still passes the eleven spec-only checks after the move and the in-text rename. (Task 3 test.)
4. `bdd-discovery` invoked with only a pasted idea and no `spec.md` — the skill text must route to step 0 (build the file), not refuse with "回 bdd-capture". (Task 3 grep.)
5. `benchmark/skeleton/go/README.md` lines 228–258 use `spec.md` in the pre-10-02 sense (FR/AC live there); the Task 1 rename must **not** touch them. (Task 1 verification step.)

---

### Task 1: Technical `spec.md` → `design.md`

**Files:**
- Rename: `skills/bdd-spec/references/spec-format.md` → `design-format.md`; `skills/bdd-spec/examples/minimal-spec.md` → `minimal-design.md`
- Modify: `skills/bdd-spec/SKILL.md`, `skills/bdd-plan/SKILL.md`, `skills/bdd-plan/references/plan-format.md`, `skills/bdd-implement/SKILL.md`, `skills/bdd-formulation/SKILL.md`, `skills/bdd-formulation/references/*.md`, `skills/bdd-formulation/examples/anti-patterns.md`, `skills/example-mapping/SKILL.md`, `skills/clarify/SKILL.md`, `skills/story-splitting/SKILL.md`, `skills/bdd-discovery/SKILL.md`, `skills/skill-rules/examples/anti-patterns.md`, `skills/bdd-formulation/scripts/check_spec.py` (comment near line 494 only), `tests/test_status.py:71`, `tests/test_check_spec.py:158`
- Modify: `PLAN.md` (open 未決事項 rows that mean the technical file, e.g. "`spec.md` 649 行"; file tree line 85; `spec.md` section header line 113–117), `CLAUDE.md:103,126-127,166`, `benchmark/README.md`, `benchmark/skeleton/go/docs/ARCHITECTURE.md`

**Interfaces:**
- Produces: the string `spec.md` no longer appears anywhere under `skills/` or `tests/` except as the script filename `check_spec.py`; `design.md` is the technical artifact's name downstream tasks use.

- [ ] **Step 1: Move the two bdd-spec files**

```bash
git mv skills/bdd-spec/references/spec-format.md skills/bdd-spec/references/design-format.md
git mv skills/bdd-spec/examples/minimal-spec.md skills/bdd-spec/examples/minimal-design.md
```

- [ ] **Step 2: Rewrite every technical-sense `spec.md` to `design.md` in the files listed above**

Read each hit of `grep -rn "spec.md" skills/ tests/ PLAN.md CLAUDE.md benchmark/README.md benchmark/skeleton/go/docs/ARCHITECTURE.md` before editing. In `bdd-spec/SKILL.md` also change the description trigger 「寫 spec.md」→「寫 design.md」 and the `references/`/`examples/` paths. In the two tests the fixture `spec.md` written as `"# old"` becomes `design.md` (the assertion strings change in Task 2, not here). Skip `benchmark/skeleton/go/README.md` entirely (Review Focus 5) and the history rows in `PLAN.md` (定版點被改過三次, 2026-10-02 進度 bullets).

- [ ] **Step 3: Verify no technical-sense `spec.md` remains**

Run: `grep -rn "spec\.md" skills/ tests/ | grep -v "check_spec.py"`
Expected: no output.

Run: `grep -n "spec\.md" benchmark/skeleton/go/README.md | wc -l`
Expected: unchanged from `git stash; <same grep>; git stash pop` (the file was not edited).

- [ ] **Step 4: Audit and validate**

Run: `for s in bdd-spec bdd-plan bdd-implement bdd-formulation example-mapping clarify story-splitting bdd-discovery; do python3 skills/skill-rules/scripts/audit_skill.py skills/$s || echo "FAIL $s"; done; claude plugin validate . && python3 -m unittest discover -s tests`
Expected: no `FAIL`, validate passes, tests pass.

- [ ] **Step 5: Commit**

```bash
git add -A skills tests PLAN.md CLAUDE.md benchmark/README.md benchmark/skeleton/go/docs/ARCHITECTURE.md
git commit -m "refactor(bdd-spec): rename the technical spec.md to design.md"
```

---

### Task 2: Scripts and tests: `prd.md` → `spec.md`, old-layout guard

**Files:**
- Modify: `skills/bdd-discovery/scripts/status.py:109-111`, `skills/bdd-formulation/scripts/check_spec.py` (lines 2–53 docstring, 365–382, 390–391, 425–432, 490; every `prd.md` and `bdd-capture`/`CAPTURE` in strings and comments)
- Test: `tests/test_status.py`, `tests/test_check_spec.py`

**Interfaces:**
- Consumes: Task 1's `design.md` fixtures in both tests.
- Produces: both `main(root)` glob `specs/*/spec.md`; messages 「沒有 spec.md —— 先跑 bdd-discovery」; a new guard in each `main` that returns 1 with the message 「舊版佈局：先 `git mv spec.md design.md`，再 `git mv prd.md spec.md`」 when any `specs/<x>/prd.md` exists.

- [ ] **Step 1: Write the failing tests**

In `tests/test_status.py` rename `test_missing_prd_exits_1_and_names_capture` → `test_missing_spec_exits_1_and_names_discovery`; fixture writes `design.md`; assertions become `assertIn("沒有任何 spec.md", out)`, `assertIn("bdd-discovery", out)`, `assertNotIn("bdd-capture", out)`. Rename `test_counts_rows_from_prd` → `test_counts_rows_from_spec` and write the fixture to `spec.md`. Add:

```python
def test_old_layout_prd_beside_spec_exits_1_and_names_git_mv(self):
    # arrange: specs/x/prd.md + specs/x/spec.md (old technical file)
    # act: main(root)
    # assert: rc == 1; "舊版佈局" in out; "git mv prd.md spec.md" in out
```

In `tests/test_check_spec.py`: helper at line 51 writes `spec.md`; `test_spec_only_dir_exits_1` fixture `design.md`, assertions `"沒有 spec.md"` / `"bdd-discovery"` / `assertNotIn("bdd-capture", out)`; `test_capture_fresh_prd_points_at_discovery` → `test_fresh_spec_points_at_discovery` with `assertNotIn("bdd-capture", out)` kept; `test_good_prd_passes_in_spec_only_mode` → `test_good_spec_passes_in_spec_only_mode`; add the same `test_old_layout_prd_beside_spec_exits_1_and_names_git_mv` for `check_spec.main`. Leave `test_shipped_example_prd_passes_prd_only_checks` pointing at `skills/bdd-capture/examples/minimal-prd.md` for now (Task 3 moves it).

- [ ] **Step 2: Run tests to verify they fail**

Run: `python3 -m unittest discover -s tests -v 2>&1 | tail -20`
Expected: the renamed/new tests FAIL on the glob (`沒有任何 prd.md`) and on the missing 舊版佈局 message.

- [ ] **Step 3: Implement in both scripts**

In each `main(root: Path) -> int`, after the `specs.is_dir()` check and before the glob:

```python
legacy = sorted(specs.glob("*/prd.md"))
if legacy:
    print(f"{legacy[0].parent} 是舊版佈局（prd.md）：先 `git mv spec.md design.md`，再 `git mv prd.md spec.md`")
    return 1
```

Then glob `*/spec.md`; messages 「{specs} 底下沒有任何 spec.md —— 先跑 bdd-discovery」 (status) and 「{specs_dir} 底下沒有 spec.md —— 先跑 bdd-discovery」 (check_spec). Replace every remaining `prd.md` with `spec.md` and every `CAPTURE`/`bdd-capture` in strings and comments with DISCOVERY 開場／`bdd-discovery` (e.g. line 391 「CAPTURE 的產物不完整」→「DISCOVERY 開場沒建完」; 425–432 both branches now route to `bdd-discovery`; 490 「CAPTURE 還沒切 story」→「DISCOVERY 開場還沒切 story」). Do not touch the sixteen checks' logic.

- [ ] **Step 4: Run tests to verify they pass**

Run: `python3 -m unittest discover -s tests -v`
Expected: all PASS.

Run: `grep -n "prd.md\|capture\|CAPTURE" skills/bdd-discovery/scripts/status.py skills/bdd-formulation/scripts/check_spec.py`
Expected: only the 舊版佈局 guard lines.

- [ ] **Step 5: Commit**

```bash
git add skills/bdd-discovery/scripts/status.py skills/bdd-formulation/scripts/check_spec.py tests/
git commit -m "feat(scripts): read spec.md instead of prd.md and refuse the old layout"
```

---

### Task 3: Merge `bdd-capture` into `bdd-discovery`

**Files:**
- Rename: `skills/bdd-capture/references/prd-format.md` → `skills/bdd-discovery/references/spec-format.md`; `skills/bdd-capture/references/persona-definition.md` → `skills/bdd-discovery/references/persona-definition.md`; `skills/bdd-capture/examples/minimal-prd.md` → `skills/bdd-discovery/examples/minimal-spec.md`
- Delete: `skills/bdd-capture/SKILL.md` (then the empty directory)
- Create (if length requires): `skills/bdd-discovery/references/opening.md`
- Modify: `skills/bdd-discovery/SKILL.md`, `skills/bdd-discovery/references/technical-probes.md:14`, the three moved files (in-text `prd.md` → `spec.md`, CAPTURE → DISCOVERY 開場), `tests/test_check_spec.py:145`

**Interfaces:**
- Consumes: Task 2's messages (`先跑 bdd-discovery`).
- Produces: `bdd-discovery` is the only entry point; its description carries capture's trigger phrases; the reference path other skills will cite in Task 4 is `../bdd-discovery/references/spec-format.md`.

- [ ] **Step 1: Move the three files and retarget the test**

```bash
git mv skills/bdd-capture/references/prd-format.md skills/bdd-discovery/references/spec-format.md
git mv skills/bdd-capture/references/persona-definition.md skills/bdd-discovery/references/persona-definition.md
git mv skills/bdd-capture/examples/minimal-prd.md skills/bdd-discovery/examples/minimal-spec.md
```

Edit `tests/test_check_spec.py:144-145`: rename the test to `test_shipped_example_spec_passes_spec_only_checks` and load `skills/bdd-discovery/examples/minimal-spec.md`.

- [ ] **Step 2: Retarget the two cross-skill links so this commit stays green**

`skills/clarify/SKILL.md:89` and `skills/bdd-spec/SKILL.md:40`: `../bdd-capture/references/prd-format.md` → `../bdd-discovery/references/spec-format.md`. Nothing else in those files changes here (Task 4 does the rest).

- [ ] **Step 3: Rename in-text inside the three moved files**

`spec-format.md`: title `# \`spec.md\` 的骨架`; the opening paragraph lists `bdd-discovery` 建檔與每輪寫答案、`bdd-formulation`、`example-mapping`; the「誰第一次寫」column: every `CAPTURE` → `DISCOVERY 開場`; path `specs/<date>-<feature>/spec.md`; the「接縫」section's technical-file mentions already say `design.md` after Task 1 — verify. `minimal-spec.md`: path/title lines only; FR/AC content untouched. `persona-definition.md`: `prd.md` → `spec.md`.

- [ ] **Step 4: Write step 0 of the opening in `bdd-discovery/SKILL.md`**

Replace the body of `bdd-capture/SKILL.md`'s 流程 1–4 and 收尾 into a new subsection `#### 0. 建檔（只在 spec.md 不存在時）` placed before `#### 掃原文找隱含假設`, or into `references/opening.md` linked from that spot (Global Constraints). It must keep, verbatim in meaning: 原文逐字進 `## 原文` 並編 `§<n>`；從原文抽開頭幾節，行尾 `← PRD §n`，抽不出來寫「原文沒提」；最多四題（表照抄），MUST NOT 問第五題；粗切 story 標「暫定」，三個分組空著，MUST NOT 填 FR／AC；收尾三條（Open Questions 紅卡、Document Overview 狀態 `澄清中`、**給需求方看一遍改到他說「這是我的意思」**，這是自己的一個回合，點頭後才進掃描）；`clarify` 三節檔的處理。

Then edit the rest of `SKILL.md`:
- description: 「把一個想法或一份需求文件建成 `spec.md` 並問到完整——原文逐字存檔、粗切 story，再從 `## 原文` 掃隱含假設…」; append capture's trigger phrases (「我有一個想法」「記一下這個需求」「幫我把需求寫下來」「這是我想做的」「需求進來了」「把這份 PRD 收進來」, English `capture this idea, write down this requirement, intake this requirement document`). Check total description ≤ 1024 chars.
- 前置確認: 「**唯一必要的輸入是需求本身**——一段描述、一份文件的內容或路徑、或已存在的 `specs/<date>-<feature>/spec.md`。」 Replace `NEVER: prd.md 不存在時憑一段描述開始問…回 bdd-capture` with `NEVER: 還沒把原文逐字寫進 ## 原文 就開始問`. Add: Ask user the feature slug when building (from capture).
- Skill Boundaries: delete 「還沒有 `prd.md` → 先跑 `bdd-capture`」; add 「要把一則太大的 story 切開 → 粗切在開場，沿規則細切用 `story-splitting`」.
- Lines 19, 33, 84, 104, 110, 125, 129, 146, 161, 296, 432 and every other `prd.md` / `bdd-capture` / `CAPTURE`: `spec.md`, `../bdd-discovery/references/…` becomes `references/…`, 「CAPTURE 抄進來的」→「開場抄進來的」, 「開場做過沒有」 detection text now also says 第 0 件事只寫 `面向` 為 `—` 的列.
- Flow diagram line 84: `需求（想法／PRD 檔）→ 第 0 件事建 spec.md（澄清中）` as the left entry.
- 產物 block line 429: `spec.md  DISCOVERY ★ 建檔，從 澄清中 寫到 待對焦`.

- [ ] **Step 5: Delete bdd-capture**

```bash
git rm skills/bdd-capture/SKILL.md
rmdir skills/bdd-capture
```

- [ ] **Step 6: Verify**

Run: `grep -rn "prd.md\|bdd-capture\|CAPTURE" skills/bdd-discovery/ tests/`
Expected: no output.

Run: `grep -rn "bdd-capture/" skills/`
Expected: no output (no skill links into the deleted directory).

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-discovery && claude plugin validate . && python3 -m unittest discover -s tests`
Expected: all pass.

Run: `git log --follow --oneline skills/bdd-discovery/references/spec-format.md | tail -1`
Expected: a commit older than today (history chained).

- [ ] **Step 7: Commit**

```bash
git add -A skills/bdd-capture skills/bdd-discovery skills/clarify/SKILL.md skills/bdd-spec/SKILL.md tests/test_check_spec.py
git commit -m "feat(bdd-discovery): absorb CAPTURE as the opening step and drop bdd-capture"
```

---

### Task 4: Every other skill: `prd.md` → `spec.md`, no more `bdd-capture`

**Files:**
- Modify: `skills/bdd-spec/SKILL.md:40` (reference path), `skills/bdd-spec/references/design-format.md`, `skills/bdd-spec/examples/minimal-design.md`, `skills/bdd-formulation/SKILL.md:36,63,110,249,405` and all other hits, `skills/bdd-formulation/references/*.md`, `skills/bdd-formulation/examples/anti-patterns.md`, `skills/bdd-plan/SKILL.md:30` and hits, `skills/bdd-plan/references/plan-format.md`, `skills/bdd-implement/SKILL.md`, `skills/example-mapping/SKILL.md:137,159` and hits, `skills/clarify/SKILL.md:89,204,229` and hits, `skills/story-splitting/SKILL.md`, `skills/skill-rules/examples/anti-patterns.md`

**Interfaces:**
- Consumes: Task 3's path `../bdd-discovery/references/spec-format.md`.
- Produces: no `prd.md` / `bdd-capture` / `CAPTURE` string under `skills/`.

- [ ] **Step 1: Rewrite the hits**

List them: `grep -rn "prd.md\|bdd-capture\|CAPTURE" skills/`. Rules: `prd.md` → `spec.md`; `../bdd-capture/references/prd-format.md` → `../bdd-discovery/references/spec-format.md`; 「沒有 `prd.md` → `bdd-capture`」 → 「沒有 `spec.md` → `bdd-discovery`，它會建」; 「CAPTURE 與 DISCOVERY 的產物」→「DISCOVERY 的產物」; `example-mapping` status table row `澄清中 | bdd-capture | 建檔時` → `bdd-discovery | 開場建檔時`; `clarify` line 229's 「之後跑 `bdd-capture`」 → 「之後 `bdd-discovery` 第 0 件事從它補開頭幾節」.

- [ ] **Step 2: Verify**

Run: `grep -rn "prd.md\|bdd-capture\|CAPTURE" skills/`
Expected: no output.

Run: `for s in $(ls skills); do python3 skills/skill-rules/scripts/audit_skill.py skills/$s || echo "FAIL $s"; done; claude plugin validate . && python3 -m unittest discover -s tests`
Expected: no `FAIL`, validate passes (one expected `--strict` warning only if run with `--strict`), tests pass.

- [ ] **Step 3: Commit**

```bash
git add skills
git commit -m "refactor(skills): read spec.md and route missing input to bdd-discovery"
```

---

### Task 5: Documents: six steps

**Files:**
- Modify: `docs/ai-sdlc.md` (lines 31–57 seven-step figure and table; 61–75 三種模式 table → three rows; 85–90 「為什麼是五個 skill」→「四個」 with a short 「CAPTURE 曾經獨立」 paragraph; 127 chain), `PLAN.md` (lines 14, 16, 38, 43, 58–70, 83–111; add a migration note to the 「重跑 parking 某張票」 item: that run needs `git mv` to the new layout; leave 定版點 and 2026-10-02 bullets), `CLAUDE.md` (lines 13, 64–70, 102–107, 118–127, 157, 166), `benchmark/README.md:31,44,61`, `benchmark/skeleton/go/README.md:4,12`, `benchmark/skeleton/go/docs/ARCHITECTURE.md:90,366`

- [ ] **Step 1: Edit the documents**

`ai-sdlc.md` figure:

```
DISCOVERY  →  SPEC  →  PLAN  →  IMPLEMENT  →  VERIFY  →  REVIEW
└──── BDD ────┘       ✗       └─ BDD ─┘        ✗          ✗
```

DISCOVERY row gains CAPTURE's 「存在的理由」 (原文逐字、粗切 story 是開場的事); the 「為什麼是四個 skill」 paragraph says why CAPTURE was split in 09-09～09-11 and why it folded back on 2026-11-13 (防線從 skill 邊界換成 `check_spec.py` 第 14 項). `CLAUDE.md` 「three modes, four disciplines」 → three rows; CAPTURE bullet removed, DISCOVERY bullet gains 「opens by writing `## 原文` verbatim and cutting stories, asks at most four analyst questions before the first round」. `PLAN.md` 七步流程 → 六步流程, `prd.md` section header → `spec.md`, format-contract path → `skills/bdd-discovery/references/spec-format.md`.

- [ ] **Step 2: Verify the only remaining hits are history**

Run: `grep -rn "prd.md\|bdd-capture\|CAPTURE" --include=*.md --include=*.py . | grep -v "docs/superpowers\|^./runs\|^./.git"`
Expected: every line is one of — `docs/ai-sdlc.md` 「CAPTURE 曾經獨立」paragraph; `PLAN.md` 定版點 / 2026-10-02 進度 / migration note; `benchmark/cases/*/grading.md` 10-02 passages; `benchmark/cases/personal-memo/input-prd.md` filename references; the two scripts' 舊版佈局 guard. Paste the list into the commit body.

Run: `claude plugin validate . && python3 -m unittest discover -s tests`
Expected: pass.

- [ ] **Step 3: Commit**

```bash
git add docs/ai-sdlc.md PLAN.md CLAUDE.md benchmark/README.md benchmark/skeleton/go/README.md benchmark/skeleton/go/docs/ARCHITECTURE.md
git commit -m "docs: describe the six-step pipeline without CAPTURE"
```
