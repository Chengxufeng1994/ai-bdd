# Plan as Commander Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `plan` the system-analysis step — inventory interfaces, design domain model, schema and API contract (the last via a new `api-and-interface-design` skill copied from addyosmani), then plan dependencies, risks, checkpoints and parallel work before `task-breakdown`.

**Architecture:** Markdown skills in a Claude Code plugin. One new skill directory; `plan` rewritten; format references split so the three designs have their own format file; downstream skills retargeted. No script logic changes; the gates are the repo's audit script, the unit tests and `claude plugin validate`.

**Tech Stack:** Markdown skills, Python 3 stdlib scripts (`audit_skill.py`, `check_spec.py`, `status.py`), `claude plugin validate`.

**Spec:** `docs/superpowers/specs/2026-11-20-plan-as-commander-design.md`

## Global Constraints

- Prose in Traditional Chinese; commit messages English, conventional title with `WHAT:`／`WHY:`／`HOW:` body, ending with the `Claude-Session:` line.
- Every `SKILL.md` keeps `disable-model-invocation: true`, the 職責卡, `## 使用時機`, `## Skill Boundaries`, `## Verification` with `- [ ]` items (house style, `skills/skill-rules/`).
- Every file in a skill's subdirectories is referenced from its `SKILL.md` or a referenced file (S5).
- `skills/discovery/SKILL.md` must not grow (it is 513 lines, over the 500 signal).
- Examples use the library-lending domain (`skills/spec/examples/minimal-spec.md`); never the parking or personal-memo benchmark domains.
- `plan.md` stays self-contained: the five sections Tech Stack, Commands, Project Structure, Code Style, Boundaries are written out in full.
- `api-and-interface-design/SKILL.md` body is addyosmani's text verbatim (MIT); only D1's additions and Chinese trigger words in the description.
- Renames with `git mv`; nothing under `runs/` is touched.
- Old formats are not accepted: a `spec.md` with `### 聚合與不變條件` or a `plan.md` without `## Domain Model` makes `task-breakdown`／`implement`／`archive` stop and ask the user to re-run `plan`.

## Review Focus

- **Contract drift:** a ticket copies operations from `plan.md`'s `## API 契約` into `api/openapi.yaml` and changes their shape — expected: the ticket's acceptance item catches it, and changes go back through `plan`. Pinned in Task 5.
- **No HTTP surface:** the seam is an application service and there is no `api/openapi.yaml` — expected: `## API 契約` holds the code-level interface contract and the copy-into-openapi step and its acceptance item are skipped. Pinned in Tasks 2 and 5.
- **Designs disagree:** an aggregate boundary splits what one table holds, or an operation writes two aggregates — expected: plan's cross-check names the mismatch and records it in `## 缺口`. Pinned in Task 3.
- **Interface with no source:** an interface listed "because systems usually have a database" — expected: rejected; every interface cites an FR／AC or `Q-<n>`. Pinned in Task 3.
- **Old-format feature:** `spec.md` still has `### 聚合與不變條件` — expected: downstream stops with a re-run-`plan` hint instead of silently reading the wrong file. Pinned in Task 5.

---

### Task 1: Add the `api-and-interface-design` skill

**Files:**
- Create: `skills/api-and-interface-design/SKILL.md`
- Create: `skills/api-and-interface-design/NOTICE`
- Modify: `RESOURCES.md` (skill-projects table)

**Interfaces:**
- Produces: a skill named `api-and-interface-design` that `plan` DELEGATEs to. Its 職責卡: **讀** `spec.md` 的 FR／AC、涵蓋的 `.feature`、專案既有的 `api/openapi.yaml`；**寫** 不寫檔——把 OpenAPI 片段（或沒有 HTTP 時的程式碼層介面契約）交回 `plan`；**上游** `plan` 第 3 步；**下游** `plan`。

- [ ] **Step 1: Fetch the original verbatim**

Run: `gh api repos/addyosmani/agent-skills/contents/skills/api-and-interface-design/SKILL.md --jq .content | base64 -d > /tmp/api-original.md`
Expected: 369 lines, frontmatter `name: api-and-interface-design`.

- [ ] **Step 2: Write `skills/api-and-interface-design/SKILL.md`**

Start from the original. Additions only:
- frontmatter: `disable-model-invocation: true`; append to `description` a line `觸發詞：「設計 API」「API 契約」「介面設計」「openapi 怎麼寫」「錯誤格式」「冪等」` (keep the English text).
- right under the `# API and Interface Design` title: one line `> 原文：addyosmani/agent-skills skills/api-and-interface-design（MIT，見同目錄 NOTICE）。本文原樣保留，以下兩節與職責卡是本 plugin 的增補。`, then the 職責卡 above.
- before `## Core Principles`: `## 使用時機`（`plan` 第 3 步設計 API 契約時；單獨叫用時只給設計建議，不寫檔；原文的 When to Use 見上）and `## Skill Boundaries`（| 要寫進 `plan.md` | `plan` | ／ | 要把契約搬進 `api/openapi.yaml` | `implement` | ／ | 要定 seam、schema、聚合 | `plan` |）.
- in the original `## Verification`, append `- [ ] 每個 operation 指得出它服務哪幾條 AC（\`AC-<n>.<m>\`）`.

- [ ] **Step 3: Write `NOTICE`**

The MIT license text with `Copyright (c) Addy Osmani` and the source URL `https://github.com/addyosmani/agent-skills`.

- [ ] **Step 4: Verify the body is unchanged**

Run: `diff <(sed '1,/^---$/d' /tmp/api-original.md | sed '1,/^---$/d') <(sed '1,/^---$/d' skills/api-and-interface-design/SKILL.md | sed '1,/^---$/d')`
Expected: only added lines (`>` hunks) for the source line, 職責卡, `## 使用時機`, `## Skill Boundaries`, and the extra Verification item; no `<` lines.

- [ ] **Step 5: Audit**

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/api-and-interface-design`
Expected: `✔ 通過`

- [ ] **Step 6: Add the `RESOURCES.md` row**

Row in the skill-projects table: source addyosmani `api-and-interface-design` (link), used by `skills/api-and-interface-design`, taken: 原樣保留，增補職責卡、使用時機、Skill Boundaries 與一項 Verification.

- [ ] **Step 7: Commit**

```bash
git add skills/api-and-interface-design RESOURCES.md
git commit   # feat(api-and-interface-design): adopt addyosmani's API design skill
```

### Task 2: Formats for the three designs and the new `plan.md` skeleton

**Files:**
- Create: `skills/plan/references/design-format.md`
- Modify: `skills/plan/references/plan-format.md` (skeleton and example)
- Modify: `skills/plan/references/tech-format.md:227-246` (Implementation Decisions without API contract and schema)

**Interfaces:**
- Produces: section names later tasks use verbatim — `## 系統介面`, `## Domain Model`, `## Database Schema`, `## API 契約`, `## 依賴圖`, `## 順序與 Checkpoint` — in this order in `plan.md`: Overview, 系統介面, Tech Stack, Commands, Project Structure, Code Style, Boundaries, Implementation Decisions, Testing Decisions, Domain Model, Database Schema, API 契約, 依賴圖, Risks, 順序與 Checkpoint, Parallel Work, 缺口, Further Notes.

- [ ] **Step 1: Write `design-format.md`**

One section per design, each with its skeleton and a library-lending example:
- `## 系統介面`: per interface `端點類型`／`主要介面`／`需求原文依據`（FR／AC 或 `Q-<n>`）; then `### Wave <n>` with `平行分析介面` and `安排理由`; then the 對應 table (API 契約、事件、請求回應、錯誤碼 → `api-and-interface-design`; 實體、欄位、生命週期、儲存 → `plan`). Credit aixbdd `system-analysis` (Apache-2.0) in one line.
- `## Domain Model`: the `../../archive/references/domain-model-format.md` skeleton with headings one level down (`### <聚合名>`), plus 「原本是什麼、為什麼改」 for changed existing aggregates; 被誰保證 names the aggregate only.
- `## Database Schema`: same table shape as `docs/DATABASE-SCHEMA.md` in `../../archive/references/docs-format.md`; one extra column 守住哪條不變條件.
- `## API 契約`: a fenced ```yaml OpenAPI fragment (new and changed `paths` and `components` only) when there is an HTTP surface; otherwise the code-level interface contract (port methods, inputs, outputs, errors). Each operation cites its ACs.
- `## 依賴圖`: nodes are modules and contracts (not tickets); `task-breakdown` treats it as a constraint.

- [ ] **Step 2: Update `plan-format.md`**

Replace the skeleton with the Interfaces order above; point each new section at `design-format.md`; keep the Risks／Parallel Work／缺口 text; rewrite the example's Overview so it names the three designs.

- [ ] **Step 3: Update `tech-format.md` Implementation Decisions**

Remove the `**API 契約**` and `**Schema**` bullets; add one line: API 契約、schema、聚合 are their own sections, see `design-format.md`.

- [ ] **Step 4: Verify**

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/plan`
Expected: `✘ [S5]` for `design-format.md` (not yet referenced from `SKILL.md`) — fixed in Task 3; no other item.
Run: `grep -nE "停車|訪客|月租|待辦" skills/plan/references/design-format.md`
Expected: no output.

- [ ] **Step 5: Commit** (together with Task 3 — the audit is only green once SKILL.md references the file; stage now, commit at the end of Task 3).

### Task 3: Rewrite `plan/SKILL.md` as the commander

**Files:**
- Modify: `skills/plan/SKILL.md` (whole 流程, 職責卡, description, Verification, 完成後)

**Interfaces:**
- Consumes: Task 1's skill name; Task 2's section names and `design-format.md`.
- Produces: `plan`'s four steps that Task 5's downstream text refers to as 「`plan` 第 2 步（盤點）」「第 3 步（設計）」「第 4 步（排）」.

- [ ] **Step 1: Rewrite the 流程 in verb-led steps** (`READ`／`WRITE`／`CHECK`／`DELEGATE`／`ASK`／`STOP`, as `skills/spec/SKILL.md` does), each step ending with **做完的樣子**:
1. 收集 — READ the files listed in design D2 step 1.
2. 盤點 — WRITE `## 系統介面`; CHECK every interface cites an FR／AC or `Q-<n>` (an interface "because systems usually have one" is rejected); CHECK each Wave only holds interfaces that do not wait for each other.
3. 設計 — ASK seam (keep today's seam rules and table); WRITE `## Domain Model`, `## Database Schema`; DELEGATE `../api-and-interface-design/SKILL.md` for `## API 契約`, in Wave order; CHECK the three side by side — every invariant has an aggregate and a constraint, every table belongs to one aggregate, every operation writes at most one aggregate unless a cross-aggregate rule says otherwise; mismatches go to `## 缺口`. Then the seven technical sections (unchanged rules, `tech-format.md`).
4. 排 — WRITE `## 依賴圖`, `## Risks`, `## 順序與 Checkpoint` (addyosmani Step 2 and 5; credit in 參考檔案), `## Parallel Work` (ASK, unchanged; add 「先定契約，再並行」); CHECK every `.feature` Then is observable at the seam (unchanged step); STOP — tell the user `/task-breakdown`.

- [ ] **Step 2: Update 職責卡, description, 參考檔案, Verification, 完成後**

Description gains the three designs and the DELEGATE; 參考檔案 lists `references/design-format.md` and `../api-and-interface-design/SKILL.md`; Verification adds: 每個系統介面有需求依據；三份設計並排對過，對不上的進了 `## 缺口`；API 契約每個 operation 指得出 AC；有 HTTP 時是 OpenAPI 片段、沒有時是程式碼層介面契約。

- [ ] **Step 3: Verify**

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/plan`
Expected: `✔ 通過`
Run: `wc -l skills/plan/SKILL.md`
Expected: under 500.

- [ ] **Step 4: Commit Tasks 2 and 3**

```bash
git add skills/plan
git commit   # refactor(plan)!: design interfaces, domain model, schema and API before planning
```

### Task 4: Domain model leaves `spec`

**Files:**
- Modify: `skills/spec/references/spec-format.md` (`## Domain Concepts` section and order table)
- Modify: `skills/spec/SKILL.md` (step 4, 產物, Verification, 完成後, 再跑 step 4)
- Modify: `skills/spec/examples/minimal-spec.md` (`## Domain Concepts`)

**Interfaces:**
- Produces: `spec.md` `## Domain Concepts` holds only `### 詞彙`; `check_spec.py` item 18 unchanged (still requires the section once 已對焦).

- [ ] **Step 1: Edit `spec-format.md`** — `## Domain Concepts` keeps the `### 詞彙` table and its rules; remove `### 聚合與不變條件` and its 判準; add one line: 聚合與不變條件 are designed in `plan.md` `## Domain Model`. Order-table row and the archive line (CONTEXT.md only) follow.

- [ ] **Step 2: Edit `spec/SKILL.md`** — step 4 keeps READ `docs/CONTEXT.md`, WRITE `### 詞彙`, CHECK conflicting definitions; drop READ `docs/DOMAIN-MODEL.md`, the aggregate WRITE and CHECK. 產物 and Verification lose aggregates.

- [ ] **Step 3: Edit `minimal-spec.md`** — drop its `### 聚合與不變條件` subsection (it says 「（無）」 today).

- [ ] **Step 4: Verify**

Run: `python3 -m unittest discover -s tests`
Expected: `OK` (the shipped-example tests still pass).
Run: `grep -rn "聚合與不變條件" skills/spec`
Expected: only the one pointer line to `plan.md`.
Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/spec`
Expected: `✔ 通過`

- [ ] **Step 5: Commit** — `refactor(spec)!: leave aggregates and invariants to plan`

### Task 5: Downstream reads the three designs

**Files:**
- Modify: `skills/task-breakdown/SKILL.md`, `skills/task-breakdown/references/tasks-format.md`
- Modify: `skills/implement/SKILL.md`
- Modify: `skills/archive/SKILL.md`, `skills/archive/references/docs-format.md`, `skills/archive/references/domain-model-format.md`

**Interfaces:**
- Consumes: Task 2's section names; Task 3's step names.

- [ ] **Step 1: `task-breakdown`** — 讀 adds `## 系統介面`, `## Domain Model`, `## Database Schema`, `## API 契約`, `## 依賴圖`; CHECK no ticket is ordered before a node it depends on in `## 依賴圖`; each 票節 gains `**契約**:` (the operations it covers) and `**資料表**:` lines in `tasks-format.md`; when `## API 契約` is an OpenAPI fragment, each ticket's 驗收條件 gains `- [ ] 搬進 api/openapi.yaml 的 operation 與 plan.md 的片段一致`. Precondition: `plan.md` without `## Domain Model` → stop, re-run `plan`.

- [ ] **Step 2: `implement`** — when the ticket has `**契約**:` and `## API 契約` is OpenAPI: copy those operations and the schemas they use from `plan.md` into `api/openapi.yaml` unchanged (append-type join point, existing merge rules); a needed change goes back to `plan`, not edited in place. When the ticket has `**資料表**:`: write the migration from `## Database Schema`. No OpenAPI fragment → skip the copy. Precondition as Step 1.

- [ ] **Step 3: `archive`** — `DOMAIN-MODEL.md` from `plan.md` `## Domain Model`; `DATABASE-SCHEMA.md` from `## Database Schema`; `TECH.md` from the five sections; `CONTEXT.md` from `spec.md` `### 詞彙` (unchanged). `domain-model-format.md`'s line about `spec.md` points at `plan.md`. Old format (`spec.md` with `### 聚合與不變條件`) → stop, re-run `plan`.

- [ ] **Step 4: Verify**

Run: `for s in task-breakdown implement archive; do python3 skills/skill-rules/scripts/audit_skill.py skills/$s; done`
Expected: three `✔ 通過`.
Run: `grep -rn "Implementation Decisions.\{0,10\}Schema\|spec.md.\{0,10\}Domain Model\|聚合與不變條件" skills/task-breakdown skills/implement skills/archive`
Expected: only old-format detection lines.

- [ ] **Step 5: Commit** — `refactor(skills)!: build tickets and docs from plan's three designs`

### Task 6: Top-level docs and final gate

**Files:**
- Modify: `README.md` (artifact table, `plan` section, new `api-and-interface-design` section, directory tree), `CLAUDE.md` (PLAN paragraph, skills list), `PLAN.md` (step table, skill table row), `docs/ai-sdlc.md` (PLAN rationale), `RESOURCES.md` (aixbdd `system-analysis` row)

- [ ] **Step 1: Edit the five files** to describe `plan` as the system-analysis step (four steps), the three designs in `plan.md`, aggregates no longer in `spec.md`, and the new skill. `docs/ai-sdlc.md` gains the reason: 「aixbdd 的 system-analysis 就等於 addyosmani 的 plan」 (需求方) and why the three designs stay in one skill (they must be read side by side).

- [ ] **Step 2: Full gate**

Run: `python3 -m unittest discover -s tests`
Expected: `OK`
Run: `for d in skills/*/; do python3 skills/skill-rules/scripts/audit_skill.py $d | tail -1; done`
Expected: 14 lines, all `✔ 通過`
Run: `claude plugin validate .`
Expected: `✔ Validation passed with warnings`
Run: `grep -rnE "停車|月租|柵欄|繳費" skills`
Expected: no output.

- [ ] **Step 3: Commit and push** — `docs: describe plan as the system-analysis step`, then `git push`.
