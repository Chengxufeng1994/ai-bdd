# Plan Follows addyosmani Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rewrite `plan` on addyosmani's planning-and-task-breakdown structure, with a design step before Step 2, and fold `task-breakdown` into it.

**Architecture:** Markdown skills. `plan/SKILL.md` is rewritten; `task-breakdown`'s rules move into it and its `tasks-format.md` moves under `plan/references/`; `task-breakdown/` is deleted; every reference to it is retargeted.

**Tech Stack:** Markdown skills; `audit_skill.py`, `check_spec.py`, `status.py`; `claude plugin validate`.

**Spec:** `docs/superpowers/specs/2026-11-20-plan-follows-addyosmani-design.md`

## Global Constraints

- Prose Traditional Chinese; commits English, conventional title, `WHAT:`／`WHY:`／`HOW:` body, last line `Claude-Session: https://claude.ai/code/session_01Qkb1m77tBhqzve8Hn4TbBP`; work on `main`, push after each task.
- Every `SKILL.md`: `disable-model-invocation: true`, 職責卡, `## 使用時機`, `## Skill Boundaries`, `## Verification` with `- [ ]`; every subdirectory file referenced (S5); `plan/SKILL.md` under 500 lines.
- Library-lending examples only; never parking or personal-memo.
- Renames and deletions with `git mv`／`git rm`; nothing under `runs/` touched.
- `tasks.md` format unchanged; `implement`'s one-ticket-per-session, merge-back and join-point rules unchanged.
- Kept from today's `plan`: seam ASK and its rules, the three designs and `design-format.md` (minus `## 系統介面`), the five technical sections in full, Implementation／Testing Decisions, the "`.feature` 在 seam 上驗得到嗎" check, Parallel Work ASK, `## 缺口`, `## Further Notes`, design-side mismatches fixed in the same run, old-format stop rules.

## Review Focus

- **Rerun with half-done tickets:** `tasks.md` already has checked tickets — expected: checked ones untouched, unchecked re-cut (addyosmani's "never overwrite an incomplete plan"). Task 1.
- **User rejects the slicing:** the ASK before writing `tasks.md` gets "too coarse" — expected: re-slice and ask again, nothing written until approved. Task 1.
- **No HTTP surface:** no OpenAPI fragment — expected: tickets carry no openapi acceptance item. Task 1.
- **Stale pointer:** a skill or doc still says "run `/task-breakdown`" — expected: none remain except history. Task 2.
- **Old `tasks/` directory layout:** expected: the fold-into-`tasks.md` rule from `task-breakdown` survives in `plan`. Task 1.

---

### Task 1: Rewrite `plan` and absorb `task-breakdown`

**Files:**
- Modify: `skills/plan/SKILL.md` (whole)
- Move: `skills/task-breakdown/references/tasks-format.md` → `skills/plan/references/tasks-format.md` (`git mv`)
- Modify: `skills/plan/references/design-format.md` (drop `## 系統介面`), `skills/plan/references/plan-format.md` (skeleton without `## 系統介面`; `## 缺口` described as addyosmani's Open Questions), `skills/plan/references/tasks-format.md` (self-references from task-breakdown → plan)
- Delete: `skills/task-breakdown/` (`git rm -r`)

**Interfaces:**
- Produces: `plan`'s steps named exactly `Step 1 收集`, `設計`, `Step 2 依賴圖`, `Step 3 垂直切片`, `Step 4 寫任務`, `Step 5 排序與 Checkpoint`, then `ASK 切法` and `WRITE`; downstream (Task 2) says 「`plan` 的 Step 4」 where it used to say `task-breakdown`.

- [ ] **Step 1:** Fetch addyosmani's text: `gh api repos/addyosmani/agent-skills/contents/skills/planning-and-task-breakdown/SKILL.md --jq .content | base64 -d`; and mattpocock `to-tickets` (`gh api repos/mattpocock/skills/contents/skills/engineering/to-tickets/SKILL.md --jq .content | base64 -d`) for the Quiz step.
- [ ] **Step 2:** Rewrite `skills/plan/SKILL.md` in the verb-led style of `skills/spec/SKILL.md`, flow exactly as design D1. Move into it from `skills/task-breakdown/SKILL.md`: the slicing rules (tracer bullet, fits a fresh session, demoable), Blocked by, join points per ticket, wide refactor expand–contract, the rerun rule (checked tickets untouched), folding an old `tasks/` directory into `tasks.md`, the precondition and old-format stops. Add addyosmani's Task Sizing (XS–XL, L+ must split) and its Common Rationalizations／Red Flags items that apply. Credit addyosmani and mattpocock in 參考檔案. Keep under 500 lines; push detail into references.
- [ ] **Step 3:** `git mv` tasks-format.md; edit it and design-format.md／plan-format.md as listed; `git rm -r skills/task-breakdown`.
- [ ] **Step 4:** Verify: `python3 skills/skill-rules/scripts/audit_skill.py skills/plan` → `✔ 通過`; `wc -l skills/plan/SKILL.md` → < 500; `grep -n "系統介面\|Wave" skills/plan -r` → no output.
- [ ] **Step 5:** Commit `refactor(plan)!: follow addyosmani's planning flow and fold task-breakdown in` with `BREAKING CHANGE: /task-breakdown no longer exists; /plan writes tasks.md`; push.

### Task 2: Retarget every reference and the docs

**Files:**
- Modify: every file from `grep -rln "task-breakdown" skills README.md CLAUDE.md PLAN.md docs/ai-sdlc.md RESOURCES.md` (implement, behavior-driven-development, spec, spec-format, wake-up, story-splitting, discovery, formulation, archive if listed, tech-format, top-level docs)
- Modify: `RESOURCES.md` — rows for addyosmani `planning-and-task-breakdown` and mattpocock `to-tickets` (Quiz step)

**Interfaces:**
- Consumes: Task 1's step names.

- [ ] **Step 1:** Replace each `task-breakdown` reference with the right `plan` step (e.g. 「`plan` 的 Step 4」 for ticket-writing, 「`/plan`」 for the next command); historical sentences in PLAN.md／ai-sdlc.md stay as history, reworded to past tense if needed. `discovery/SKILL.md` must not grow.
- [ ] **Step 2:** Full gate: `python3 -m unittest discover -s tests` → OK; `for d in skills/*/; do python3 skills/skill-rules/scripts/audit_skill.py $d | tail -1; done` → 13 lines, all `✔ 通過` except `api-and-interface-design`'s accepted F6 warning; `claude plugin validate .` → passes with warnings; `grep -rn "task-breakdown" skills` → no output.
- [ ] **Step 3:** Commit `docs: point every task-breakdown reference at plan`; push.
