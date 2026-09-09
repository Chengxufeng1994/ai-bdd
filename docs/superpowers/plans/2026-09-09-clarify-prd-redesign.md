# CLARIFY 產出單一 PRD Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** CLARIFY 改成一個 feature 產出一份自足的 `specs/<date>-<feature>/prd.md` 且不切分；切 story 搬到 SPEC；詞彙表由 SPEC 寫進 `docs/CONTEXT.md`。

**Architecture:** 先定契約再改生產者。Task 1 寫出新格式與一份可解析的 fixture，Task 2–3 用它把兩支腳本 TDD 到綠，Task 4–6 才改 skill 與下游文件。這個順序就是這個 repo 自己的「產物格式就是介面」——介面先固定，兩端才動得起來。

**Tech Stack:** Markdown（skill 與 reference）、Python 3 標準庫（兩支腳本，無外部依賴）

**Spec:** [`docs/superpowers/specs/2026-09-09-clarify-prd-redesign.md`](../specs/2026-09-09-clarify-prd-redesign.md)

## Global Constraints

- **不引進測試框架。** repo 沒有任何測試基礎設施，`pytest` 也沒安裝。驗證一律是「建 fixture → 跑腳本 → 驗 exit code 與 stdout」，與 Makefile 既有的 `test -z "$(gofmt -l .)"` 同風格。
- **兩支腳本只用標準庫。** 現行 `status.py` 與 `check_spec.py` 都是 stdlib-only，維持。
- **fail loud 不可退化。** 找不到輸入時必須非零離開並指名缺什麼，**絕不可從零輸入推導出結論**。這是改資料來源時最容易弄壞的一條。
- **散文用繁體中文，commit 訊息用英文**，WHAT/WHY/HOW 三段。
- **`docs/CONTEXT.md` 的例外寫得很窄**：SPEC 得建立該檔或在其中追加條目，不得改寫既有段落，不得寫入 `docs/` 底下其他檔案。
- 每個 task 結束時工作區乾淨、`python3 skills/skill-rules/scripts/audit_skill.py` 對動過的 skill 為 exit 0。

## File Structure

| 檔案 | 責任 | Task |
| --- | --- | --- |
| `skills/bdd-clarify/references/prd-format.md` | **新建。** `prd.md` 的欄位、章節、編號規則、Open Questions 表格式 | 1 |
| `skills/bdd-clarify/examples/minimal-prd.md` | **新建。** 最小可解析 fixture，兩支腳本的測試對象 | 1 |
| `skills/bdd-clarify/scripts/status.py` | 改讀 `prd.md` 的 Open Questions 表 | 2 |
| `skills/bdd-spec/scripts/check_spec.py` | 改讀 `prd.md` 的 FR/EX ＋ `spec.md` 的 story 歸屬；新增 FR-without-EX 檢查 | 3 |
| `skills/bdd-clarify/SKILL.md` | 拿掉切 story 與分批；三個 Pass 寫同一份 `prd.md` | 4 |
| `skills/bdd-clarify/references/{clarify-format,map-format,prd-breakdown}.md` | 刪除，內容併入 `prd-format.md` | 4 |
| `skills/bdd-spec/SKILL.md` | 新增切 story 與 `docs/CONTEXT.md` | 5 |
| `skills/clarify-loop/SKILL.md` | 找待答題的方式 | 6 |
| `PLAN.md`／`CLAUDE.md`／`docs/ai-sdlc.md` | 下游更新 | 6 |

---

### Task 1: 定義新契約 —— `prd-format.md` 與 fixture

這是介面。後面五個 task 全部對著它做。

**Files:**
- Create: `skills/bdd-clarify/references/prd-format.md`
- Create: `skills/bdd-clarify/examples/minimal-prd.md`

**Interfaces:**
- Produces: `prd.md` 的章節名與編號格式，Task 2–6 全部依賴：
  - `## Functional Requirements` → `### FR-<n>` → `#### AC`（`AC-<n>.<m>`）／`#### Examples`（`EX-<n>.<m>`）
  - `## Open Questions` 的表格欄位：`| Q | 問題 | 面向 | 狀態 | 答案 |`，狀態值為 `已答`／`待答`／`n/a`
  - `## Non-Functional Requirements`、`## Scope`、`## Actors`、`## Assumptions / Constraints`

- [ ] **Step 1: 寫 `references/prd-format.md`**

內容涵蓋六件事，每一件都要有範例而不只是規則：

1. **章節順序與各節該裝什麼**（照 spec 的〈產物契約〉）
2. **編號規則**：`FR-<n>` 平鋪不分組；`AC-<n>.<m>` 與 `EX-<n>.<m>` 掛在對應 FR 底下；**編號在此定版**，`.feature` 的 `@example-<n>.<m>` 回指 `EX`
3. **EARS 句式**：只說「用 EARS 寫，六個 pattern 見 [`docs/sdd.md`](../../../docs/sdd.md) 的 `## EARS`」——**不重述那六個 pattern**
4. **AC 與 EX 的分野**：AC 給 PM（使用者看得到什麼），EX 給 SPEC（具體到含數字）
5. **`## Open Questions` 的雙層格式**：表是機械可解析的索引，每題底下的 `### Q<n> <標題>` 小節裝追問串
6. **角色留在 `prd.md`、詞義決策住 Open Questions**（spec 講死的兩條界線）

`## Open Questions` 那一節逐字寫成：

````markdown
## Open Questions

表是索引，給腳本讀；小節是決策史，給人讀。

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q1 | 月租費帳務在不在範圍 | 範圍 | 已答 | 名單與期限在，錢不在 |
| Q2 | 消防法規對閘門的要求 | 法規 | 待答 | — |
| Q3 | 機車費率 | 範圍 | n/a | 本期不做機車 |

**`面向` 欄不可省。** `status.py` 用它算追問覆蓋——十一個業務面向掃過幾個。
省掉它不會報錯，只會讓覆蓋率永遠顯示 0，而那是靜默地錯。

### Q2 消防法規對閘門的要求

**狀態**：待答　**輪次**：3　**該問誰**：消防設備師、柵欄廠商

2026-09-09 問管委會 → 「我知道有規定但背不出來，不敢猜」

**擋住**：FR-11（緊急車輛進場）
````

- [ ] **Step 2: 寫 fixture `examples/minimal-prd.md`**

最小但涵蓋所有腳本要解析的結構。逐字：

````markdown
# PRD：社區地下停車場計費系統

## Problem / Goal / Success Metrics

訪客臨停目前無法計費。目標是讓訪客離場前完成付款。

## Actors

| 角色 | 是誰 | 怎麼取得 |
| --- | --- | --- |
| 訪客 | 非月租名單上的車 | 入口抽紙票 |
| 月租戶 | 名單上且未到期 | 感應卡 |

## Scope — In / Out

**In**：訪客計費、出場放行
**Out**：月租費線上收款

## Functional Requirements

### FR-1  When 訪客車出場, the system shall 依停留時長計費，前 30 分鐘免費。

#### AC
- AC-1.1  訪客在繳費機看得到金額與停留時長
- AC-1.2  付款完成前柵欄不開

#### Examples
- EX-1.1  停 29 分 → 收 0 元
- EX-1.2  停剛好 30 分 → 收 0 元
- EX-1.3  停 31 分 → 收 30 元

### FR-2  When 月租戶刷卡出場, the system shall 直接開啟柵欄，不計費。

#### AC
- AC-2.1  月租戶出場不需經過繳費機

#### Examples
- EX-2.1  名單內且未到期 → 開柵欄，金額 0

## Non-Functional Requirements

- NFR-1  出場柵欄自刷卡到開啟不超過 2 秒

## Open Questions

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q1 | 月租費帳務在不在範圍 | 範圍 | 已答 | 名單與期限在，錢不在 |
| Q2 | 消防法規對閘門的要求 | 法規 | 待答 | — |
| Q3 | 機車費率 | 範圍 | n/a | 本期不做機車 |

### Q2 消防法規對閘門的要求

**狀態**：待答　**輪次**：3　**該問誰**：消防設備師、柵欄廠商

2026-09-09 問管委會 → 「我知道有規定但背不出來，不敢猜」

**擋住**：FR-11（緊急車輛進場）

## Assumptions / Constraints

- 沒有車牌辨識，去年區權會否決
````

這份 fixture **刻意包含一條待答題**——Task 2 要驗「有未答問題時就緒判定擋得住」。

- [ ] **Step 3: 驗證 fixture 沒有孤兒引用**

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify`
Expected: exit 0。若 S5（孤兒檔案）報 `examples/minimal-prd.md`，在 `prd-format.md` 裡加一句指向它的文字引用即可——那是 S5 的已知偽陽性處理方式。

- [ ] **Step 4: Commit**

```bash
git add skills/bdd-clarify/references/prd-format.md skills/bdd-clarify/examples/minimal-prd.md
git commit -F - <<'EOF'
docs(bdd-clarify): define the prd.md contract and a parsable fixture

WHAT: Add references/prd-format.md describing the single self-contained
      PRD, and examples/minimal-prd.md as the fixture both scripts parse
WHY: The format is the interface between all four steps, so it has to
     exist and be testable before either producer or consumer moves.
     Writing it first is what lets the scripts be driven from a fixture
     rather than from a skill that has not been rewritten yet.
HOW: The Open Questions table carries a 面向 column that the design
     document omitted — status.py computes probe coverage from it, so
     dropping it would have left coverage silently reporting zero. EARS
     is referenced to docs/sdd.md rather than restated, keeping one owner
     per definition.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 2: `status.py` 改讀 `prd.md`

**Files:**
- Modify: `skills/bdd-clarify/scripts/status.py`

**Interfaces:**
- Consumes: Task 1 的 `examples/minimal-prd.md` 與 Open Questions 表格式
- Produces: `parse_open_questions(text) -> list[dict]`，每筆含 `q`／`問題`／`面向`／`狀態`／`答案`

- [ ] **Step 1: 建測試 fixture 樹並跑現行腳本，確認它對新結構是壞的**

```bash
rm -rf /tmp/prdtest && mkdir -p /tmp/prdtest/specs/2026-09-09-parking-billing
cp skills/bdd-clarify/examples/minimal-prd.md /tmp/prdtest/specs/2026-09-09-parking-billing/prd.md
python3 skills/bdd-clarify/scripts/status.py /tmp/prdtest; echo "exit=$?"
```

Expected: exit 1，訊息說 `specs` 底下沒有任何 slice 帶 `questions/` 目錄。**這是對的**——現行腳本本來就找不到新結構，這一步是確認起點。

- [ ] **Step 2: 寫新的解析函式**

在 `status.py` 加入，取代 `parse_status` 與 `scan_questions`：

```python
def parse_open_questions(text: str) -> list[dict[str, str]]:
    """prd.md 的 `## Open Questions` 表 -> 每題一個欄位字典。

    只讀那一節的表格列，不掃整份檔案：prd.md 的其他章節也有表格
    （Actors、NFR），對整份掃會把它們算成問題。

    表頭固定五欄 `| Q | 問題 | 面向 | 狀態 | 答案 |`；欄數不符的列
    直接丟進 unparseable，不猜。
    """
    section = re.search(
        r"^## Open Questions\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return []
    rows = []
    for line in section.group(1).splitlines():
        line = line.strip()
        if not line.startswith("|"):
            continue
        cells = [c.strip() for c in line.strip("|").split("|")]
        if len(cells) != 5 or cells[0] in ("Q", "---"):
            continue
        if set(cells[0]) <= {"-", " "}:
            continue
        rows.append(dict(zip(("q", "問題", "面向", "狀態", "答案"), cells)))
    return rows
```

- [ ] **Step 3: 改 `main()` 的來源**

把 `qdirs = sorted(p for p in specs.glob("*/questions") if p.is_dir())` 換成：

```python
    prds = sorted(specs.glob("*/prd.md"))
    if not prds:
        print(f"{specs} 底下沒有任何 prd.md —— 先跑 bdd-clarify")
        return 1
```

每份 `prd.md` 的統計改用 `parse_open_questions(prd.read_text(encoding="utf-8"))`：
`狀態` 為 `已答`／`待答`／`n/a` 各自計數，`面向` 收集成集合餵給既有的 `coverage()`。

**`coverage()` 不動。** 它吃的是面向集合，來源換了但介面沒換。

- [ ] **Step 4: 跑測試，驗證數字正確**

```bash
python3 skills/bdd-clarify/scripts/status.py /tmp/prdtest; echo "exit=$?"
```

Expected: exit 0；輸出顯示 `2026-09-09-parking-billing` 有 **已答 1／待答 1／n/a 1**，面向涵蓋 `範圍`、`法規` 兩項。

- [ ] **Step 5: 驗證 fail loud 沒有退化**

```bash
python3 skills/bdd-clarify/scripts/status.py /tmp/nonexistent; echo "exit=$?"
rm -rf /tmp/prdempty && mkdir -p /tmp/prdempty/specs
python3 skills/bdd-clarify/scripts/status.py /tmp/prdempty; echo "exit=$?"
```

Expected: 兩次都 `exit=1`，且都印出**缺什麼**而不是印一份看起來正常的空報表。

這一步不可省。改資料來源最容易產生的 bug 就是「找不到就當成零」，而那正是這個 repo 反覆說的靜默地錯。

- [ ] **Step 6: Commit**

```bash
git add skills/bdd-clarify/scripts/status.py
git commit -F - <<'EOF'
feat(bdd-clarify): read clarification progress from prd.md

WHAT: Replace status.py's per-file question scan with a parser for the
      Open Questions table in specs/<date>-<feature>/prd.md
WHY: CLARIFY now writes one self-contained PRD per feature instead of a
     directory of question files, so the old glob finds nothing. The
     script's own claim is unchanged — progress is computed, never
     stored — only where it counts from moves.
HOW: The parser reads only the Open Questions section rather than the
     whole file, because Actors and NFR are tables too and a whole-file
     scan would count them as questions. Rows whose column count does not
     match are collected as unparseable rather than guessed at. Verified
     that both a missing directory and an empty specs/ still exit 1 and
     name what is absent: reading from a new source is exactly where a
     "not found means zero" regression would hide.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 3: `check_spec.py` 改讀 `prd.md` 與 `spec.md`

**Files:**
- Modify: `skills/bdd-spec/scripts/check_spec.py`

**Interfaces:**
- Consumes: Task 1 的 `FR-<n>` / `EX-<n>.<m>` 格式
- Produces: `frs_in_prd(text) -> dict[str, set[str]]`（FR 編號 → 該 FR 的 EX 編號集合）、`stories_in_spec(text) -> dict[str, set[str]]`（story-slug → 它涵蓋的 FR 編號集合）

**這個 task 有一件 spec 沒講死的事，在這裡定案：** 現行 `stories_in_clarify()` 從 `clarify.md` 讀 story→例子。切 story 搬到 SPEC 之後，**FR 與 EX 的目錄在 `prd.md`，story 由哪些 FR 組成則在 `spec.md`**。所以覆蓋檢查要讀兩個檔：

```
prd.md   FR-1 → {EX-1.1, EX-1.2, EX-1.3}      例子的目錄
spec.md  visitor-billing → {FR-1}             切分的結果
.feature @example-1.1 …                        實際寫出來的場景
```

- [ ] **Step 1: 跑現行腳本，確認它對新結構是壞的**

```bash
rm -rf /tmp/spectest && mkdir -p /tmp/spectest/specs/2026-09-09-parking-billing /tmp/spectest/features
cp skills/bdd-clarify/examples/minimal-prd.md /tmp/spectest/specs/2026-09-09-parking-billing/prd.md
python3 skills/bdd-spec/scripts/check_spec.py /tmp/spectest; echo "exit=$?"
```

Expected: exit 1，訊息說 `specs` 底下沒有 `clarify.md`。確認起點。

- [ ] **Step 2: 補齊 fixture —— 寫 `spec.md` 與一個 `.feature`**

```bash
cat > /tmp/spectest/specs/2026-09-09-parking-billing/spec.md <<'EOF'
# SPEC：社區地下停車場計費系統

## Stories

### visitor-billing
涵蓋 FR-1。沿規則切：訪客計費是一條獨立成立的約束，自己可驗收。

### monthly-pass-exit
涵蓋 FR-2。
EOF

cat > /tmp/spectest/features/visitor-billing.feature <<'EOF'
Feature: Visitor billing

  Rule: 前 30 分鐘免費

    @example-1.1
    Scenario: 停 29 分
      Given 一台訪客車停了 29 分鐘
      When 它在繳費機結帳
      Then 應收金額是 0 元

    @example-1.2
    Scenario: 停剛好 30 分
      Given 一台訪客車停了 30 分鐘
      When 它在繳費機結帳
      Then 應收金額是 0 元
EOF
```

注意 `EX-1.3`（停 31 分 → 30 元）**故意沒有對應場景**——Step 4 要驗雙向覆蓋抓得到漏做的例子。

- [ ] **Step 3: 寫兩個新解析函式，取代 `stories_in_clarify`**

```python
def frs_in_prd(text: str) -> dict[str, set[str]]:
    """prd.md 的 `## Functional Requirements` 段：FR 編號 -> EX 編號集合。

    來源是 prd.md 而不是 spec.md：例子的定版編號誕生在 CLARIFY，
    spec.md 只記哪些 FR 湊成一則 story，不重述例子。
    """
    section = re.search(
        r"^## Functional Requirements\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return {}
    out: dict[str, set[str]] = {}
    current = None
    for line in section.group(1).splitlines():
        fr = re.match(r"^### FR-(\d+)\b", line)
        if fr:
            current = fr.group(1)
            out.setdefault(current, set())
            continue
        if current:
            for ex in re.findall(r"\bEX-(\d+\.\d+)\b", line):
                out[current].add(ex)
    return out


def stories_in_spec(text: str) -> dict[str, set[str]]:
    """spec.md 的 `## Stories` 段：story-slug -> 它涵蓋的 FR 編號集合。"""
    section = re.search(r"^## Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return {}
    out: dict[str, set[str]] = {}
    current = None
    for line in section.group(1).splitlines():
        st = re.match(r"^### (\S+)\s*$", line)
        if st:
            current = st.group(1)
            out.setdefault(current, set())
            continue
        if current:
            for fr in re.findall(r"\bFR-(\d+)\b", line):
                out[current].add(fr)
    return out
```

- [ ] **Step 4: 改 `check()` 的來源與覆蓋比對**

`clarifies = sorted(specs_dir.glob("*/clarify.md"))` 換成 `prds = sorted(specs_dir.glob("*/prd.md"))`，找不到時訊息改成 `底下沒有 prd.md —— 先跑 bdd-clarify`。

每個 feature 的比對邏輯：

```python
    frs = frs_in_prd(prd_path.read_text(encoding="utf-8"))
    spec_path = prd_path.parent / "spec.md"
    if not spec_path.exists():
        print(f"✗ {prd_path.parent.name} 有 prd.md 但沒有 spec.md —— 先跑 bdd-spec")
        return 1
    stories = stories_in_spec(spec_path.read_text(encoding="utf-8"))

    # 逐 story 比，不可把所有 story 的例子聯集起來跟單一 .feature 比：
    # 一則 story 一個 .feature，聯集會讓 A 的例子出現在 B 檔裡也算通過。
    for slug, fr_ids in sorted(stories.items()):
        expected = {ex for fr in fr_ids for ex in frs.get(fr, set())}
        fpath = feat_dir / f"{slug}.feature"
        tags = (set(re.findall(r"@example-(\d+\.\d+)",
                               fpath.read_text(encoding="utf-8")))
                if fpath.exists() else set())
        missing = expected - tags      # 漏做的例子
        invented = tags - expected     # 發明出來的驗收條件
```

- [ ] **Step 5: 新增「FR 沒有任何例子」檢查**

這是 spec 的 D3 兩層結構唯一的漂移防線：

```python
    barren = sorted(fr for fr, exs in frs.items() if not exs)
    if barren:
        print(f"✗ 這些 FR 沒有任何 EX，SPEC 無從寫出場景：{', '.join('FR-' + b for b in barren)}")
```

有 `barren` 時整體回傳非零。

- [ ] **Step 6: 跑測試，驗證三件事**

```bash
python3 skills/bdd-spec/scripts/check_spec.py /tmp/spectest; echo "exit=$?"
```

Expected: exit 非零，且輸出同時包含——
1. `EX-1.3` 被列為漏做（`.feature` 沒有 `@example-1.3`）
2. `monthly-pass-exit` 的 `EX-2.1` 被列為漏做（那個 `.feature` 根本還沒寫）
3. **沒有** `barren` 訊息（fixture 裡兩條 FR 都有 EX）

- [ ] **Step 7: 驗證 barren 檢查真的會觸發**

```bash
printf '\n### FR-3  When 系統故障, the system shall 開啟柵欄。\n\n#### AC\n- AC-3.1  故障時不困住車輛\n' \
  >> /tmp/spectest/specs/2026-09-09-parking-billing/prd.md
python3 skills/bdd-spec/scripts/check_spec.py /tmp/spectest 2>&1 | grep "FR-3"
```

Expected: 輸出含 `FR-3`，說它沒有任何 EX。

- [ ] **Step 8: 驗證 fail loud**

```bash
python3 skills/bdd-spec/scripts/check_spec.py /tmp/nonexistent; echo "exit=$?"
```

Expected: `exit=1` 並指名缺什麼。

- [ ] **Step 9: Commit**

```bash
git add skills/bdd-spec/scripts/check_spec.py
git commit -F - <<'EOF'
feat(bdd-spec): check coverage against prd.md and spec.md

WHAT: Replace the clarify.md story scan with two parsers — FR to EX from
      prd.md, story to FR from spec.md — and add a check for FRs that
      carry no examples at all
WHY: Story slicing moved from CLARIFY to SPEC, which splits what used to
     be one lookup into two: the catalogue of examples is decided in
     CLARIFY and numbered there, while which requirements make up a story
     is now SPEC's own output. Reading both keeps the chain checkable end
     to end rather than trusting either file alone.
HOW: The barren-FR check is the only mechanical defence for the design's
     two-layer split between acceptance criteria and examples: it cannot
     catch an AC edited without its examples, but it does catch a
     requirement with no examples at all, which is the case where SPEC
     would otherwise have to invent them — something it is forbidden to
     do. Verified against a fixture that deliberately omits one example
     and one whole feature file, so both directions of the coverage check
     are exercised rather than assumed.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 4: `bdd-clarify` —— 拿掉切分，三個 Pass 寫同一份 PRD

**Files:**
- Modify: `skills/bdd-clarify/SKILL.md`
- Delete: `skills/bdd-clarify/references/clarify-format.md`、`map-format.md`、`prd-breakdown.md`
- Modify: `skills/bdd-clarify/references/actor-definition.md`（改指 `prd.md` 的 `## Actors`）

**Interfaces:**
- Consumes: Task 1 的 `prd-format.md`
- Produces: CLARIFY 寫出 `specs/<date>-<feature>/prd.md`，Task 5 的 SPEC 讀它

- [ ] **Step 1: 改流程圖**

把 `## 流程` 的 code block 換成：

```
PRD／一段敘述（新的）          或      繼續澄清（bdd-clarify [feature]）
        ↓                                          ↓
Pass 1 · 廣度 —— 對整個 feature                    （跳過 Pass 1）
  1. 拆成目標／成功指標／範圍邊界＋隱含假設         ↓
  2. 識別角色                                       ↓
  3. 只問「會改變範圍」的題目                       ↓
        ↓                                          ↓
Pass 2 · 深度 —— 對整個 feature ←──────────────────┘
  4. 澄清循環（每輪）
       ├─ 逐題問：一題一列進 ## Open Questions，狀態記已答／待答
       ├─ 抽出本輪能抽的 FR 與 EX
       └─ 還有未知？ → 下一輪
  5. 就緒判定（提案）
  6. 收尾：標記完成
        ↓
Pass 3 · 技術 —— 逐項掃過 technical-probes.md
  量得出來的 → ## Non-Functional Requirements
  不可協商的外部限制 → ## Constraints
        ↓
一份 specs/<date>-<feature>/prd.md
```

- [ ] **Step 2: 刪掉 `### 4. 切 story` 與 `### 5. 分批` 兩整節**

從 `### 4. 切 story —— 依據來自上面三步` 到 `## Pass 2 · 深度` 之前全部刪除，並把 `## Pass 2 · 深度 —— 每一則 story` 改成 `## Pass 2 · 深度 —— 對整個 feature`。

Pass 2 底下的步驟編號從 5、6、7 改成 4、5、6。

- [ ] **Step 3: 改 `## 產物` 一節**

整節換成：

````markdown
## 產物

```
specs/<date>-<feature-name>/
└── prd.md          CLARIFY  ★ 唯一產物，完整自足

（spec.md 與 plan.md 是 SPEC 與 PLAN 的產物，本 skill 不得寫入）
features/*.feature           SPEC
docs/CONTEXT.md              SPEC
```

一個 feature 一份 `prd.md`，不切分——**切 story 是 SPEC 的事**，它需要領域知識
而且需要範圍先穩定。格式 → `references/prd-format.md`。

規則與例子的**定版編號在這裡誕生**：`FR-<n>`、`EX-<n>.<m>`。`.feature` 的
`@example-<n>.<m>` 回指它們，`spec.md` 不重述。

MUST NOT: 另存一份進度儀表板。進度是**算出來的**——`scripts/status.py` 從
`## Open Questions` 的表直接數。
````

- [ ] **Step 4: 改 Skill Boundaries 與續跑段**

- `- 產出是 brief.md／actor.md／glossary.md／questions/ 的問答` → `- 產出是一份 prd.md，**不是**規格、程式碼或資料模型`
- 續跑段的 `brief.md 不存在時才回頭跑 Pass 1` → `prd.md 不存在時才回頭跑 Pass 1`
- `還開著的：grep -L '狀態.*已答' questions/*.md` → `還開著的：python3 scripts/status.py`
- 「切不動、或切完仍然太大時 → `story-splitting`」整句刪除（那是 SPEC 的事了）

- [ ] **Step 5: 刪掉三份過時的 reference，改寫第四份**

```bash
git rm skills/bdd-clarify/references/clarify-format.md \
       skills/bdd-clarify/references/map-format.md \
       skills/bdd-clarify/references/prd-breakdown.md
```

`actor-definition.md` 裡所有指向 `actor.md` 與 `glossary.md` 的句子，改成指
`prd.md` 的 `## Actors`，並說明**詞彙不在這裡**——它由 SPEC 寫進 `docs/CONTEXT.md`。

- [ ] **Step 6: 稽核與殘留檢查**

```bash
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify; echo "exit=$?"
grep -rn "clarify\.md\|brief\.md\|actor\.md\|glossary\.md\|<slice-slug>\|分批" skills/bdd-clarify/
```

Expected: audit exit 0；grep 只剩刻意保留的說明（例如「詞彙由 SPEC 寫」那句提到 `CONTEXT.md`）。

- [ ] **Step 7: Commit**

```bash
git add -A skills/bdd-clarify
git commit -F - <<'EOF'
refactor(bdd-clarify): produce one self-contained PRD, stop slicing

WHAT: Drop Pass 1's story-splitting and batching steps, collapse the four
      root-level artifacts into a single specs/<date>-<feature>/prd.md,
      and delete the three reference files that described the old shape
WHY: Slicing ran in Pass 1, before the basis for slicing was stable. A
     live parking-billing run made that concrete: across fifteen scope
     questions the stakeholder revised three times and each revision
     changed the shape, so any slice drawn after the first batch would
     have been wrong. The cause was structural — the main document lived
     at specs/<slice-slug>/, so a slug had to exist before it could be
     written, which forced slicing to come first.
HOW: With batches gone, the cross-batch accumulators have nothing left to
     accumulate across, so brief/actor/glossary/questions collapse into
     the PRD rather than being relocated. Pass 3's output is routed
     explicitly — measurable things to NFR, non-negotiable external
     limits to Constraints — because a technical probe with no named
     destination is one that quietly lands nowhere. The reference files
     are deleted rather than emptied so a stale format cannot be found
     and followed.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 5: `bdd-spec` —— 接下切 story 與詞彙表

**Files:**
- Modify: `skills/bdd-spec/SKILL.md`

**Interfaces:**
- Consumes: Task 4 的 `prd.md`
- Produces: `spec.md` 的 `## Stories` 段（`### <story-slug>` 底下列 `FR-<n>`），Task 3 的 `stories_in_spec()` 解析它

- [ ] **Step 1: 改輸入契約**

開頭那句 `一份就緒的 example map 進來，出去的是 .feature` 改成：

```
一份就緒的 `prd.md` 進來，出去的是 `.feature` 與 `spec.md`：**先切 story——
決定哪幾條 FR 湊成一則可獨立驗收的 story——再把 FR 變成 `Rule:` 區塊、
EX 變成 `Example:`、編號變成 tag**，步驟套一組封閉文法，外加一張覆蓋表。
```

- [ ] **Step 2: 新增 `## 切 story` 一節**

放在現行「這份 `.feature` 是寫給誰讀的」之前：

````markdown
## 切 story

`prd.md` 的 FR 是平鋪的。第一件事是決定哪幾條湊成一則 story——**一則 story
＝ 一個 `.feature` 檔 ＝ 一次可獨立驗收的交付**。

沿規則切，不沿使用者旅程或畫面切：後者容易切出「做一半」的 story，前半段
交付了但沒有任何一條 FR 被完整滿足。

每一則要通過 INVEST（**Small 除外**，那是切分本身要解決的）。切不動、或切完
仍然太大時 → `story-splitting`，那裡有九種切分模式。

切分結果寫進 `specs/<date>-<feature>/spec.md`：

```markdown
## Stories

### visitor-billing
涵蓋 FR-1、FR-4。沿規則切：訪客計費是一條獨立成立的約束，自己可驗收。

### monthly-pass-exit
涵蓋 FR-2。
```

**`### <slug>` 與它底下提到的 `FR-<n>` 是機械可解析的**——`scripts/check_spec.py`
用它算「這則 story 該有哪些 `@example` tag」。slug 就是 `.feature` 的檔名。

MUST: 寫下**為什麼是這個切法**。半年後有人要加一則 story 時，第一個該讀的
就是這一段。

MUST NOT: 切出 `prd.md` 的 `## Scope — Out` 明確排除的東西。範圍要擴張就回
`bdd-clarify` 公開改範圍——悄悄擴張比公開改糟，因為沒有人有機會反對。
````

- [ ] **Step 3: 新增 `## 詞彙表` 一節**

````markdown
## 詞彙表 —— `docs/CONTEXT.md`

寫 `.feature` 時措辭必須一致，所以詞彙在這一步才變成承重的東西。把 `prd.md`
裡**已經有答案**的詞抬進 `docs/CONTEXT.md`。

**只轉錄，不裁決。** 詞義有爭議時不要自己定一個——那是訪談，SPEC 不訪談。
回 `bdd-clarify` 把它變成一題。

MUST: 只得建立 `docs/CONTEXT.md`，或在其中**追加**條目。
MUST NOT: 改寫該檔既有的任何段落；寫入 `docs/` 底下其他任何檔案。

這是「文件型產物一律寫入 `specs/`」的第二個例外（第一個是 `.feature`，位置由
runner 決定）。理由是**壽命**：刪掉某個 feature 的 spec 目錄之後，「訪客」的
定義應該還在。
````

- [ ] **Step 4: 改 Skill Boundaries 與產物清單**

- `例子還不夠、還有紅卡 → 回 bdd-clarify` 保留
- `產出是 .feature ＋ specs/<slice-slug>/spec.md` → `產出是 .feature ＋ specs/<date>-<feature>/spec.md ＋ docs/CONTEXT.md 的追加`
- 所有 `clarify.md` 的引用改成 `prd.md`

- [ ] **Step 5: 稽核與殘留檢查**

```bash
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
grep -rn "clarify\.md\|example map\|<slice-slug>" skills/bdd-spec/
```

Expected: audit exit 0；grep 只剩刻意保留的歷史說明。

- [ ] **Step 6: Commit**

```bash
git add skills/bdd-spec
git commit -F - <<'EOF'
feat(bdd-spec): take over story slicing and own the glossary

WHAT: Add a story-slicing section that groups prd.md's flat FRs into
      stories recorded in spec.md, and a glossary section that appends
      settled terms to docs/CONTEXT.md
WHY: Slicing stories is domain judgement and slicing tickets is not, so
     they cannot share a home. ai-sdlc.md already names PLAN as the only
     step needing no domain knowledge, which rules it out; SPEC was
     already doing the work unnamed, since features/<story-slug>.feature
     means it decides which scenarios go in which file. The glossary
     lands here for the same reason — wording only becomes load-bearing
     when step definitions start reusing it.
HOW: spec.md's Stories section is written to be machine-readable, because
     check_spec.py now derives each story's expected @example tags from
     it; prose alone would have made the chain uncheckable at exactly the
     point where slicing moved. The CONTEXT.md permission is deliberately
     narrow — create or append only, nothing else under docs/ — so that
     "do not modify the project's existing files" survives having an
     exception at all.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 6: `clarify-loop` 與下游文件

**Files:**
- Modify: `skills/bdd-plan/SKILL.md`（一行引用更新）
- Modify: `skills/clarify-loop/SKILL.md`
- Modify: `PLAN.md`、`CLAUDE.md`、`docs/ai-sdlc.md`

- [ ] **Step 0: 改 `bdd-plan` 對 `clarify.md` 的引用**

**spec 說「`bdd-plan` 無變動」，那是錯的。** `skills/bdd-plan/SKILL.md:63` 寫著
`其餘依 clarify.md 的 Business Rules 切`。改成：

```
- 其餘依 `clarify.md` 的 Business Rules 切，一張票一組相關的 `@example-N.M`。
+ 其餘依 `prd.md` 的 Functional Requirements 切，一張票一組相關的 `@example-N.M`。
```

這是純引用更新，PLAN 的職責與讀的檔（`spec.md` ＋ `.feature`）都沒變。

- [ ] **Step 1: 改 `clarify-loop` 找待答題的方式**

- `grep -L '狀態.*已答' specs/questions/*.md specs/*/questions/*.md` → `python3 skills/bdd-clarify/scripts/status.py`
- 「`questions/` 裝的是這個 feature 問過的每一題」整段改成描述 `prd.md` 的 `## Open Questions` 表
- `specs/<topic-slug>/clarify-log.md` 保留（那是 clarify-loop 單獨使用時的產物，不受影響）

- [ ] **Step 2: 改 `PLAN.md` 未決事項 #1**

```markdown
1. ~~**產物存放路徑。**~~ 已定：文件型產物一律寫入使用端 repo 的 `specs/`，
   一個目錄裝完，刪掉即乾淨。**兩個例外**，各有具體理由：

   | 例外 | 理由 |
   | --- | --- |
   | `features/*.feature` | 位置由測試 runner 的慣例決定 |
   | `docs/CONTEXT.md` | 壽命長於任何 feature——刪掉某個 feature 的 spec 目錄之後，詞彙定義應該還在 |

   例外的範圍很窄：SPEC 得建立 `docs/CONTEXT.md` 或在其中追加條目，**不得改寫
   該檔既有段落，不得寫入 `docs/` 底下其他檔案**。
```

原未決事項 #2（`.feature` 的位置）併進上表，該條標記為已併入。

- [ ] **Step 3: 改 `CLAUDE.md`**

- 「三條紀律」那段的 **CLARIFY** 一列補上「不切 story——那是 SPEC 的事」
- 「產物」一段：`specs/` 的說明改成 `specs/<date>-<feature>/prd.md`，並補上 `docs/CONTEXT.md` 這個例外

- [ ] **Step 4: 改 `docs/ai-sdlc.md`**

「六步，以及哪幾段是 BDD」的表格，CLARIFY 與 SPEC 兩列的「存在的理由」欄更新：

- **CLARIFY**：`兩件事：**對齊業務目標、產出一份對焦用的 PRD**（不是 BDD）＋把需求逼成規則與例子（BDD 的 Discovery）`
- **SPEC**：`三件事：切 story（不是 BDD）＋把例子寫成可執行規格（BDD 的 Formulation）＋把技術決定記下來（不是 BDD——它讓 .feature 以外的產物一起構成 live document）`

「四塊不屬於 BDD 的東西」那句改成**五塊**，加入「SPEC 的切 story」。

- [ ] **Step 5: 全域殘留檢查與連結檢查**

```bash
grep -rn "clarify\.md\|brief\.md\|<slice-slug>\|questions/" skills/ PLAN.md CLAUDE.md docs/ai-sdlc.md --exclude-dir=superpowers
python3 - <<'EOF'
import re, pathlib, sys
bad = []
for md in pathlib.Path('.').rglob('*.md'):
    if '.git' in md.parts: continue
    text = md.read_text(encoding='utf-8')
    text = re.sub(r'^(```|````).*?^\1', '', text, flags=re.S | re.M)
    for link in re.findall(r'\]\((\.[^)]+)\)', text):
        if not (md.parent / link.split('#')[0]).exists():
            bad.append(f"{md}: {link}")
print('\n'.join(bad) or 'all links resolve'); sys.exit(1 if bad else 0)
EOF
```

Expected: grep 只剩 `docs/superpowers/` 底下的 spec 與計畫（它們記述變更本身）；連結全部解得到。**特別注意 Task 4 刪掉了三份 reference，指向它們的連結會在這裡現形。**

- [ ] **Step 6: plugin 驗證**

```bash
claude plugin validate .; echo "exit=$?"
```

Expected: 通過（`--strict` 仍有 CLAUDE.md 那則已知警告，見 CLAUDE.md 的說明）。

- [ ] **Step 7: Commit**

```bash
git add skills/bdd-plan skills/clarify-loop PLAN.md CLAUDE.md docs/ai-sdlc.md
git commit -F - <<'EOF'
docs: bring the loop skill and the pipeline docs onto the PRD contract

WHAT: Point clarify-loop at status.py instead of globbing question files,
      record docs/CONTEXT.md as the second exception to the specs/-only
      rule, and update the six-step table now that SPEC slices stories
WHY: The artifact format is the interface between steps, so a document
     still describing the old one is not merely stale — it is an
     instruction to produce something no downstream step reads. The
     exception in particular has to be written down with its reason,
     because an undocumented exception is indistinguishable from a rule
     nobody enforces.
HOW: ai-sdlc.md's count of things that are not BDD goes from four to
     five: story slicing is domain judgement this pipeline adds, not
     something Cucumber's three practices cover. PLAN.md's two open items
     about artifact location merge into one entry with a table, since
     they were always the same decision with two exceptions rather than
     two decisions. Verified with a link checker after three reference
     files were deleted in an earlier task — deletions are where dangling
     links appear, and they exit zero everywhere.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```
