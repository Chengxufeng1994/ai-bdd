# bdd-spec-by-example Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 新 skill `bdd-spec-by-example` 在就緒判定之前把 `spec.md` 寫成 `specs/<date>-<feature>/acceptance/*.feature` 的業務語意旅程，`check_spec.py` 加第 17 項擋掉引用不存在 AC 的旅程。

**Architecture:** 旅程是不執行的規格。`check_spec.py` 在每份 `spec.md` 的迴圈裡讀同目錄 `acceptance/*.feature`，比對 `@covers-` 與 `frs_in_spec()` 的 AC；`find_features` 的 rglob 退路排除 `specs/`。其餘是 Markdown skill 與文件的改動。

**Tech Stack:** Python 3 stdlib（`unittest`）、Markdown skills、`claude plugin validate .`

**Spec:** `docs/superpowers/specs/2026-11-14-spec-by-example-design.md`

## Global Constraints

- 文件繁中；commit 英文，conventional title ＋ `WHAT:`／`WHY:`／`HOW:` body，結尾 `Claude-Session: https://claude.ai/code/session_01Qkb1mBWbYa6q2KffkJjgeT`
- 搬檔用 `git mv`；不用 `--no-verify`
- 新增或改過的 skill 先讀 `skills/skill-rules/`，交件前 `python3 skills/skill-rules/scripts/audit_skill.py skills/<name>` 過
- 旅程：關鍵字 `Feature`／`Rule`／`Example` 英文、內容中文；不受封閉步驟文法約束
- `@covers-<n>.<m>` 一個 tag 一個 AC；`Example:` 上方至少一個
- 不執行旅程、不改可執行 `.feature` 的格式與 tag、不要求每條 AC 都被涵蓋
- 實作在新分支 `spec-by-example` 上做（從 `main` 切）

## Review Focus

1. `Feature:` 上方的 tag 含 `@covers-` 時，不得被算成底下某個 `Example:` 的錨點（Task 1 `test_feature_level_covers_do_not_anchor_examples`）
2. 一行多個 tag（`@covers-1.1 @covers-1.2`）兩個都要讀到（Task 1 `test_multiple_covers_on_one_line`）
3. 形式錯的 tag（`@covers-1`、`@covers-AC-1.1`）不能安靜略過，要報成不存在的 AC（Task 1 `test_malformed_covers_reported_as_unknown`）
4. 有 `features/` 又有 `specs/x/acceptance/` 時，旅程不得出現在孤兒 `.feature` 清單（Task 1 `test_acceptance_not_listed_as_orphan_when_features_dir_exists`）
5. `specs/` 以外的巢狀 `.feature`（如 `app/features/`）仍要被 rglob 找到——排除只針對 `specs/`（Task 1 `test_rglob_still_finds_features_outside_specs`）

---

### Task 1: `check_spec.py` 第 17 項與 `find_features` 排除 `specs/`

**Files:**
- Modify: `skills/bdd-formulation/scripts/check_spec.py`（`find_features` :63-69；新函式放 `cited_qs` 之後；`check()` 在 AC 編號對帳 `mismatched` 那段之後呼叫；模組 docstring 「十六件事」→「十七件事」並加第 17 項；spec-only 訊息「十一項」→「十二項」、第 6-16 →第 6-17）
- Test: `tests/test_check_spec.py`（新 class `AcceptanceJourneys`，放在 `if __name__` 區塊**之前**）

**Interfaces:**
- Produces: `journey_covers(text: str) -> tuple[set[str], list[str]]` —— 回 (所有 `@covers-` 後面的值，原樣字串；沒有任何 `@covers-` 的 `Example:`／`Scenario:` 標題列表)
- Produces: `find_features(root)` 的 rglob 退路略過路徑含 `root / "specs"` 的檔

`journey_covers` 的掃法：逐行，`@` 開頭的行把 `@covers-(\S+)` 全部收進 pending；空行與 `#` 註解行不清 pending；`Example:`／`Scenario:`（含 `Scenario Outline:`）行消耗 pending——pending 沒有 covers 就記這個標題；其他任何行（含 `Feature:`、`Rule:`、步驟）清掉 pending。

- [ ] **Step 1: 寫失敗的測試**

fixture：`JOURNEY` 常數，`Feature: 訪客出場` 底下一個 `Rule:`、一個 `@covers-1.1` 的 `Example:`，配既有的 `GOOD`／`DISCOVERY`（`GOOD` 只有 AC-1.1）。helper `write_journey(root, name, text)` 寫到 `specs/2026-10-02-parking/acceptance/<name>`。

```python
class AcceptanceJourneys(unittest.TestCase):
    def test_journey_citing_known_ac_passes(self):        # write_spec(GOOD) + JOURNEY → code 0
    def test_journey_citing_unknown_ac_fails(self):       # JOURNEY 改 @covers-9.9 → code 1, "9.9" 與檔名 in out
    def test_example_without_covers_fails(self):          # 拿掉 tag → code 1, Example 標題 in out
    def test_no_acceptance_dir_is_silent(self):           # 只 write_spec(GOOD) → code 0, "acceptance" not in out
    def test_feature_level_covers_do_not_anchor_examples(self):
        # "@covers-1.1\nFeature: x\n  Example: y\n" → journey_covers 回 ({"1.1"}, ["y"])
    def test_multiple_covers_on_one_line(self):
        # "@covers-1.1 @covers-1.2\nExample: y\n" → ({"1.1","1.2"}, [])
    def test_malformed_covers_reported_as_unknown(self):  # @covers-1 → code 1, "@covers-1" in out

class FindFeatures(unittest.TestCase):
    def test_acceptance_only_returns_none(self):          # 只有 specs/x/acceptance/a.feature → find_features is None
    def test_rglob_still_finds_features_outside_specs(self):  # app/features/a.feature + specs/x/acceptance/b.feature → root/"app"/"features"
    def test_acceptance_not_listed_as_orphan_when_features_dir_exists(self):
        # write_spec(GOOD) + features/version.feature + JOURNEY → "訪客出場" 檔名 not in 孤兒那一行
```

- [ ] **Step 2: 跑測試確認失敗**

Run: `python3 -m unittest discover -s tests -k Acceptance -k FindFeatures 2>&1 | tail -20`
Expected: `AttributeError: ... journey_covers`、`test_acceptance_only_returns_none` FAIL（回的是 acceptance 目錄）

- [ ] **Step 3: 實作**

`find_features`：rglob 結果濾掉 `specs_dir in p.parents`。`check()`：`acceptance = spec_path.parent / "acceptance"`；不存在就不印任何東西；存在就對每個 `*.feature` 跑 `journey_covers`，`unknown = covers - {ac for acs in frs.values() for ac in acs}`，兩種問題各印一行 `✗ <dir> 的 acceptance/<file> ...` 並 `problems += 1`。訊息措辭自訂，但要帶檔名與 tag（`@covers-9.9`）或 Example 標題。

- [ ] **Step 4: 跑全部測試**

Run: `python3 -m unittest discover -s tests`
Expected: `OK`，舊的 27 個加新的 10 個都過（若舊測試斷言「十一項」字串，跟著改成「十二項」並 ledger 一條 ruling）

- [ ] **Step 5: Commit**

```bash
git add skills/bdd-formulation/scripts/check_spec.py tests/test_check_spec.py
git commit  # feat(bdd-formulation): check acceptance journeys cite existing ACs
```

---

### Task 2: 新 skill `bdd-spec-by-example`

**Files:**
- Create: `skills/bdd-spec-by-example/SKILL.md`
- Create: `skills/bdd-spec-by-example/references/journey-format.md`
- Create: `skills/bdd-spec-by-example/examples/訪客出場結帳.feature`
- Test: `tests/test_check_spec.py`（`AcceptanceJourneys` 加一個測試）

**Interfaces:**
- Consumes: Task 1 的第 17 項
- Produces: skill 名 `bdd-spec-by-example`；範例路徑 `skills/bdd-spec-by-example/examples/訪客出場結帳.feature`

內容來源：
- `SKILL.md`：照 `skills/skill-rules/` 的結構（frontmatter、使用時機、Skill Boundaries、參考檔案、流程）。輸入 `待對焦` 的 `spec.md`、補問清單「（無）」；輸出 `acceptance/<旅程>.feature`；有 `# [need clarification]` → 同一題寫進 `spec.md` `## 回 DISCOVERY 補問`（`**來自**` 標旅程檔名），交回 `bdd-spec` 第 7 步；沒有 → 交回 `example-mapping` 提就緒問句。不訪談、不改 FR／AC、不寫 `已對焦`。再跑規則照設計 D3。反目標照設計〈反目標〉兩條
- `references/journey-format.md`：設計 D2 的六條規則原文 ＋ 檔名與語言那一段 ＋ 一段骨架
- 範例：覆蓋 `bdd-spec/examples/minimal-spec.md` 的 AC-1.1、1.2、1.3（訪客停 29／30／31 分的一條長旅程可拆兩個 `Example:`），只用業務語言，每個 `Example:` 有 `@covers-`，每個 `Then` 對得到它 cover 的 AC。不放 `# [need clarification]`（範例是已收斂的）

- [ ] **Step 1: 寫失敗的測試**

```python
    def test_shipped_journey_example_passes(self):
        # minimal-spec.md + minimal-discovery.md + examples/訪客出場結帳.feature 放進 acceptance/
        # → code 0
```

- [ ] **Step 2: 跑測試確認失敗**

Run: `python3 -m unittest discover -s tests -k shipped_journey`
Expected: FAIL，`FileNotFoundError` 指向範例檔

- [ ] **Step 3: 寫三個檔**

- [ ] **Step 4: 驗證**

Run: `python3 -m unittest discover -s tests && python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec-by-example && claude plugin validate .`
Expected: `OK`；audit 無錯；validate 通過

- [ ] **Step 5: Commit**

```bash
git add skills/bdd-spec-by-example tests/test_check_spec.py
git commit  # feat(bdd-spec-by-example): write business-language acceptance journeys before readiness
```

---

### Task 3: 接進流程——`bdd-spec`、`example-mapping`、`bdd-formulation`

**Files:**
- Modify: `skills/bdd-spec/SKILL.md`（〈流程〉圖 :69-80、〈Skill Boundaries〉、第 5 步 :172-176、第 7 步「三種情況」→加「旅程有 `# [need clarification]`」）
- Modify: `skills/example-mapping/SKILL.md`（〈就緒判定〉的問句範例加一句「驗收旅程在 `specs/<dir>/acceptance/`」）
- Modify: `skills/bdd-formulation/SKILL.md`（第 1 步「讀三樣東西」表加一列 `acceptance/*.feature`＝**用語的參考**，以及設計 D2a 的 MUST NOT 與「跟 `docs/CONTEXT.md` 打架以 CONTEXT 為準、回報、不改旅程」）

**Interfaces:**
- Consumes: Task 2 的 skill 名

第 5 步改成：補問清單「（無）」→ `example-mapping` 攤地圖、跑四條診斷 → `bdd-spec-by-example` 寫旅程、跑 `check_spec.py` → 補問清單仍是「（無）」才交回 `example-mapping` 提就緒問句。`已對焦` 的寫入者、喊停、代演規則不動。

- [ ] **Step 1: 改三個 SKILL.md**

- [ ] **Step 2: 驗證**

Run: `grep -n "bdd-spec-by-example" skills/bdd-spec/SKILL.md && grep -n "acceptance" skills/example-mapping/SKILL.md skills/bdd-formulation/SKILL.md && for s in bdd-spec example-mapping bdd-formulation; do python3 skills/skill-rules/scripts/audit_skill.py skills/$s; done && claude plugin validate .`
Expected: 三個 grep 都有命中；第 5 步出現的順序是 診斷 → `bdd-spec-by-example` → 就緒問句；三個 audit 無錯；validate 通過

- [ ] **Step 3: Commit**

```bash
git add skills/bdd-spec/SKILL.md skills/example-mapping/SKILL.md skills/bdd-formulation/SKILL.md
git commit  # feat(bdd-spec): run acceptance journeys between diagnostics and readiness
```

---

### Task 4: 文件

**Files:**
- Modify: `docs/ai-sdlc.md`（:48 與 :70 的 SPEC 列加驗收旅程 `bdd-spec-by-example`）
- Modify: `PLAN.md`（〈步驟 → skills 對照〉SPEC 底下加一列，狀態「**已實作**」；〈產物佈局〉加 `acceptance/`；〈進度〉加一條 `- [x] bdd-spec-by-example`）
- Modify: `CLAUDE.md`（SPEC 的描述加旅程；Conventions 的產物清單加 `acceptance/*.feature`；`check_spec.py` 那行說明加「and acceptance journeys' `@covers-` once they exist」）

- [ ] **Step 1: 改三份文件**

- [ ] **Step 2: 驗證**

Run: `grep -n "acceptance" docs/ai-sdlc.md PLAN.md CLAUDE.md && python3 -m unittest discover -s tests && claude plugin validate .`
Expected: 三份都有命中；`OK`；validate 通過

- [ ] **Step 3: Commit**

```bash
git add docs/ai-sdlc.md PLAN.md CLAUDE.md docs/superpowers/plans/2026-11-14-spec-by-example.md
git commit  # docs: place acceptance journeys in the SPEC step
```
