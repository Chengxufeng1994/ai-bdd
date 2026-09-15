# 待辦清理 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把計畫一與計畫二累積的十一項 follow-up 清乾淨，其中最重要的一項是讓 `bdd-spec` 交件前跑得到那八項只讀 `spec.md` 的檢查。

**Architecture:** 四個任務按「動什麼」分：腳本與流程（唯一會改到程式碼的）、`spec-format.md` 的內文、fixture、以及散落在兩個 skill 與兩份文件的小項。每個任務各自獨立，順序無關緊要，但腳本那個排第一——它是唯一有行為改變的，早跑早暴露。

**Tech Stack:** Python 3 腳本（`check_spec.py`、`status.py`、`audit_skill.py`）、Markdown skill 文件、`claude plugin validate`。沒有 CI，全部手跑。

**Spec:** 無單一 spec。本計畫的依據是計畫一與計畫二的 final review 與 ledger 所記錄的十一項 follow-up，逐項在下方複述並標明出處與現況。每一項都已對 `main` 的現況驗證過（2026-09-15，`20b126d`）。

## Global Constraints

- **`check_spec.py` 的十三項檢查，行為一項都不准改。** 本計畫只改**它們在什麼時候跑**，不改任何一項在檢查什麼、怎麼判、印什麼訊息。
- **接縫不變：** `FR-<n>` ↔ `@rule-<n>`、`AC-<n>.<m>` ↔ `@example-<n>.<m>`。
- **`spec.md` 的十六節：標題文字一字不改，順序不動，不增不減。** 本計畫不加節也不刪節。
- **AC 的層級不動。**
- `docs/bdd.md`、`docs/tdd.md`、`docs/sdd.md` 轉錄外部來源，**不得**寫入專案意見；本計畫完全不碰。
- 散文用繁體中文；commit message 用英文，conventional-commit 標題加 `WHAT:` / `WHY:` / `HOW:` 正文，結尾兩行：
  `Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>`
  `Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN`
- skill 目錄只准有 `rules/`、`references/`、`examples/`、`scripts/`（S3，MUST）；子目錄檔案要被同 skill 指到（S5，MUST）；`SKILL.md` 本文超過 500 行拿 S10（SHOULD）。
- 提交前清 `__pycache__`：`find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null`

### 基線（動手前已量過）

| 量什麼 | 現值 |
| --- | --- |
| `audit_skill.py` 跑 8 個 skill | 全部 exit 0，沒有任何 S10 |
| `claude plugin validate .` | 通過，**恰好 1 個警告**（plugin root 的 `CLAUDE.md`） |
| `status.py` 對 fixture | `合計  2  1  1`，exit 0 |
| `check_spec.py` 對 fixture（空 `features/`） | `4 個問題`，exit 1 |
| `check_spec.py` 對 fixture（**沒有** `features/`） | 印 `找不到任何 .feature`，exit 1，**八項 spec.md 檢查一項都沒跑** |
| `spec-format.md` | 736 行 |
| `examples/minimal-spec.md` | 182 行，16 個 `^## ` 標題 |
| `bdd-spec/SKILL.md` | 292 行 |

**兩個 fixture**，本計畫每個任務都會用到：

```bash
# A：有 features/（空的）—— 現行迴歸基準
FA="$(mktemp -d)/fx"; mkdir -p "$FA/specs/2026-09-09-parking" "$FA/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FA/specs/2026-09-09-parking/spec.md"

# B：完全沒有 features/ —— 這是 bdd-spec 交件當下的樣子
FB="$(mktemp -d)/fx"; mkdir -p "$FB/specs/2026-09-09-parking"
cp skills/bdd-spec/examples/minimal-spec.md "$FB/specs/2026-09-09-parking/spec.md"
```

---

## 十一項 follow-up 與它們的出處

| # | 是什麼 | 出處 | Task |
| --- | --- | --- | --- |
| 1 | `## Implementation Decisions` 的「每一條要指得出它來自哪個**已答的問題或哪條規則**」不接受 `推論`，而那一節的六個 bullet 現在帶 `← <標記>`（接受 `推論`），fixture 的條目也是 `← 推論` | 計畫二 final review；**計畫二讓它惡化** | 2 |
| 2 | `## Risks` 的欄位 fixture 與骨架不一致：骨架 `風險｜從哪條規則長出來｜影響什麼`，fixture `風險｜影響｜目前怎麼辦` | 計畫二 final review #8 | 3 |
| 3 | `## Testing Decisions` 的表格把 `← 推論` 貼在格子值後面，產生 `\| ✗ ← 推論 \|` | 計畫二 scoped re-review | 2 |
| 4 | fixture 的 `## Testing Decisions` 沒有骨架要求的 Seam 表、每個 Example 的測試層級表、Prior art | 計畫二 final review #7 | 3 |
| 5 | fixture 沒有任何 story 寫「為什麼是這個切法」，而 `SKILL.md:83` 把它列成 MUST | 計畫二 final review #9 | 3 |
| 6 | `spec-format.md:508` 「空著的分類保留標題並寫為什麼空」沒有範例 | 計畫二 task review | 2 |
| 7 | `spec-format.md:648` `## Further Notes` 的「沒有就寫「（無）」」沒有範例 | 計畫二 final review | 2 |
| 8 | `bdd-spec/SKILL.md:222` 重述了一個比 `spec-format.md:98` 窄的 該問誰 規則，沒有指路說哪個是規範 | 計畫二 final review | 4 |
| 9 | `bdd-spec-review` 被兩個 skill 的 Skill Boundaries 指名、`PLAN.md:204` 列為未實作，但 `skills/bdd-spec-review/` 不存在 | 計畫一 final review | 4 |
| 10 | `README.md` 的 skill 樹列了 7 個，漏掉 `example-mapping/` | 計畫一 final review | 4 |
| 11 | `check_spec.py` 十三項裡有**八項只讀 `spec.md`**，但沒有 `.feature` 時腳本在跑到它們之前就退出；而 `bdd-spec` 的流程裡沒有任何一步跑任何腳本 | 計畫一 park，計畫二 final review 重申 | 1 |

**使用者已裁定（2026-09-15）**：第 11 項用「讓 `check_spec.py` 降級跑」，不拆成兩支腳本。理由是零重複邏輯——同一件事兩份實作是這個 repo 最大的風險——而跨 skill 引用腳本有先例（`bdd-discovery/SKILL.md:403` 指著 `../bdd-spec/scripts/status.py`）。

---

## Task 1: `check_spec.py` 降級跑 ＋ `bdd-spec` 加自檢步驟

本計畫唯一動到程式碼的任務。

**Files:**
- Modify: `skills/bdd-formulation/scripts/check_spec.py`（`check()` 一個函式，其餘不動）
- Modify: `skills/bdd-spec/SKILL.md`（流程加一步）
- Modify: `skills/bdd-formulation/SKILL.md`（稽核那一步說明它跑得到什麼）
- Modify: `CLAUDE.md`（腳本那兩行的註解）

**Interfaces:**
- Consumes: 無
- Produces: `check_spec.py` 在沒有 `.feature` 時的新行為，Task 3 的 fixture 驗證會用到

- [ ] **Step 1: 先讓它紅燈——證明八項現在跑不到**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FB="$(mktemp -d)/fx"; mkdir -p "$FB/specs/2026-09-09-parking"
cp skills/bdd-spec/examples/minimal-spec.md "$FB/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-formulation/scripts/check_spec.py "$FB"; echo "exit=$?"
echo "$FB" > /tmp/fb.txt
```

Expected: 只印 `找不到任何 .feature`，`exit=1`。**沒有任何一行提到 story、FR、AC 或 Q**——那就是這一步要修好的事。

- [ ] **Step 2: 造一份會踩中那八項的壞 `spec.md`，證明現在抓不到**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FB=$(cat /tmp/fb.txt)
python3 - "$FB/specs/2026-09-09-parking/spec.md" <<'PY'
import pathlib, sys, re
p = pathlib.Path(sys.argv[1]); t = p.read_text(encoding="utf-8")
# 壞法一：AC 編號掛錯 FR（第 12 項）
t = t.replace("- **AC-1.2**", "- **AC-9.2**", 1)
p.write_text(t, encoding="utf-8")
print("已把一條 AC 的編號改成掛不上它的 FR")
PY
python3 skills/bdd-formulation/scripts/check_spec.py "$FB"; echo "exit=$?"
```

Expected: **仍然只印 `找不到任何 .feature`**，`exit=1`。壞掉的編號沒有被提到。這就是第 11 項的具體樣子：`spec.md` 壞了，而看得出來的只有「還沒寫 .feature」。

- [ ] **Step 3: 改 `check()` 的早退為旗標**

`skills/bdd-formulation/scripts/check_spec.py` 的 `check()` 函式，目前開頭是：

```python
    specs_dir = root / "specs"
    feat_dir = find_features(root)
    if not specs_dir.is_dir():
        print(f"找不到 {specs_dir} —— 沒有 SPEC 的產物可以比對")
        return 1
    if feat_dir is None:
        print("找不到任何 .feature")
        return 1

    problems = 0
    print(f"specs: {specs_dir}    feature: {feat_dir}\n")
```

改成：

```python
    specs_dir = root / "specs"
    feat_dir = find_features(root)
    if not specs_dir.is_dir():
        print(f"找不到 {specs_dir} —— 沒有 SPEC 的產物可以比對")
        return 1

    # 十三項裡有八項只讀 spec.md（第 6-13 項）。SPEC 交件到 FORMULATION
    # 開跑之間還沒有任何 .feature，而那正是 spec.md 最需要被檢查的時刻——
    # 早退會讓那八項在整段空窗期一項都跑不到，壞掉的 spec.md 於是要等到
    # 下一個 skill 才被發現。所以缺 .feature 不再是「不能跑」，是「跑得到
    # 的先跑，跑不到的明講跳過」。
    spec_only = feat_dir is None

    problems = 0
    if spec_only:
        print(f"specs: {specs_dir}    feature: 找不到任何 .feature\n"
              f"只跑 spec.md 自己的八項檢查；覆蓋比對、狀態 tag、方言陷阱、"
              f"缺口註解、樣板重用率五項需要 .feature，跳過。\n")
    else:
        print(f"specs: {specs_dir}    feature: {feat_dir}\n")
```

- [ ] **Step 4: 守住三段用到 `.feature` 的程式碼**

`check()` 裡有三段需要 `feat_dir`，全部在那八項檢查**之後**。逐段加守衛，**不要改動段落內部的任何一行**：

1. **逐 story 比對的內層迴圈**，開頭是 `for slug, fr_ids in sorted(stories.items()):`。整段包進 `if not spec_only:`。但 `covered.add(slug)` 這一行要留在外面——它記錄的是「spec.md 提到哪些 story」，跟 `.feature` 無關，而第 3 段會用到。所以改成：

```python
        for slug in sorted(stories):
            covered.add(slug)

        if not spec_only:
            for slug, fr_ids in sorted(stories.items()):
                ...原本的內容，整段往內縮一層，一個字不改...
```

2. **樣板重用率那一段**，開頭是 `written = [p.read_text(...) for p in feat_dir.glob("*.feature")`。整段包進 `if not spec_only:`。
3. **孤兒 `.feature` 那一段**，開頭是 `orphans = [p.name for p in feat_dir.glob("*.feature")`。整段包進 `if not spec_only:`。

- [ ] **Step 5: 兩個 fixture 都要對**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FB=$(cat /tmp/fb.txt)
echo "--- B：沒有 features/，spec.md 有一條 AC 編號壞掉 ---"
python3 skills/bdd-formulation/scripts/check_spec.py "$FB"; echo "exit=$?"
```

Expected: 印出 `只跑 spec.md 自己的八項檢查…` 的說明，**並且指名那條編號對不上的 AC**，`exit=1`。

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FB2="$(mktemp -d)/fx"; mkdir -p "$FB2/specs/2026-09-09-parking"
cp skills/bdd-spec/examples/minimal-spec.md "$FB2/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-formulation/scripts/check_spec.py "$FB2"; echo "exit=$?"
```

Expected: 印出說明，然後 `全部通過`，**`exit=0`**。乾淨的 `spec.md` 在還沒有 `.feature` 時應該是綠的——這是這次改動最重要的一個新行為。

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FA="$(mktemp -d)/fx"; mkdir -p "$FA/specs/2026-09-09-parking" "$FA/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FA/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-formulation/scripts/check_spec.py "$FA" | sed "s|$FA|FX|g"; echo "exit=${PIPESTATUS[0]}"
```

Expected: **與基線逐字節相同**——`4 個問題`，`exit=1`，開頭仍是 `specs: FX/specs    feature: FX/features`。有 `.feature` 時什麼都沒變，這是不准動十三項行為的那條約束。

- [ ] **Step 6: `bdd-spec/SKILL.md` 加自檢步驟**

現有流程是五步（`### 1. 挑 story` 到 `### 5. 增修 specs/domain-model.md`）。在 `### 5.` **之後**新增：

```markdown
### 6. 自檢，然後才算完成

```bash
python3 <ai-bdd>/skills/bdd-formulation/scripts/check_spec.py <專案根>
```

`.feature` 還不存在，所以它只跑得到八項——全部是 `spec.md` 自己的形狀：FR 有沒有掛
AC、`## User Stories` 底下有沒有 v2 的舊標題、story 標題形式對不對、`#### Q` 的形式、
紅卡指向不存在的題號、紅卡列了已答的題、`AC-<n>.<m>` 的 `<n>` 對不對得上它的 FR、
有沒有 FR 不屬於任何一則 story。它會明講另外五項跳過了。

**這八項全是「安靜地錯」的形狀**：每一行都合規、檔案照樣 parse、讀起來通順，
但下游的解析器看不到那條 FR，於是它的 AC 從覆蓋率裡消失，而**沒有任何東西會報錯**。
交件前跑一次是這條鏈上唯一一個能在寫的人還記得為什麼這樣寫的時候抓到它們的位置。

MUST: `exit` 非 0 就修完再交件。修不掉的（例如缺的答案要回頭問）寫進
`## Scope — In / Out` 的「回 CLARIFY 補問」，不要留給下一步去撞。
```

同時把 `## 完成後` 的回報清單加一條，講自檢跑出什麼。

- [ ] **Step 7: `bdd-formulation/SKILL.md` 的稽核那一步要講清楚它跑得到什麼**

那一節現在寫「檢查十三件事」。`bdd-spec` 現在也會跑同一支腳本，所以要說明兩邊看到的不一樣：

在該節的 `python3 <skill>/scripts/check_spec.py <專案根>` 那個區塊底下加一段：

```markdown
**`bdd-spec` 也會跑同一支腳本，但它看到的是八項。** 沒有 `.feature` 時，需要它的
五項（雙向覆蓋、狀態 tag、方言陷阱、缺口註解、樣板重用率）會明講跳過。
跑到這一步才是十三項全開——**這五項只有這裡驗得到**，所以這一步不能因為
「`bdd-spec` 已經跑過了」而省略。
```

- [ ] **Step 8: `CLAUDE.md` 的腳本註解**

`CLAUDE.md` 的 Skill scripts 一節，那兩行現在是：

```bash
python3 skills/bdd-spec/scripts/status.py [root]                # clarification progress per feature
python3 skills/bdd-formulation/scripts/check_spec.py [root]     # .feature ↔ spec.md consistency
```

第二行的註解已經不完整了——它現在也在沒有 `.feature` 時檢查 `spec.md` 自己。改成：

```bash
python3 skills/bdd-spec/scripts/status.py [root]                # clarification progress per feature
python3 skills/bdd-formulation/scripts/check_spec.py [root]     # spec.md shape; adds .feature ↔ spec.md coverage once features exist
```

同一節開頭那句「Both take a project root and default to `.`; both exit 1 and name what is missing rather than inferring anything from zero input.」**要重讀一遍再決定改不改**——`check_spec.py` 現在對一份乾淨的、還沒有 `.feature` 的專案會 exit 0，那句話是否還成立，說出你的判斷。

- [ ] **Step 9: 全套驗收與提交**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for d in skills/*/; do printf "%-18s" "$(basename $d)"; \
  python3 skills/skill-rules/scripts/audit_skill.py "$d" >/dev/null 2>&1 && echo "exit 0" || echo "exit $?"; done
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
claude plugin validate . 2>&1 | tail -3
```

Expected: 8 個 skill 全 `exit 0`；`validate` 恰好 1 個警告。

提交拆兩筆：腳本行為一筆，四份文件一筆。腳本那筆的訊息要說清楚十三項的**行為**一項未改，改的只是**時機**。

---

## Task 2: `spec-format.md` 的四項內文

**Files:**
- Modify: `skills/bdd-spec/references/spec-format.md`

**Interfaces:**
- Consumes: 無
- Produces: 無

- [ ] **Step 1: 第 1 項——`## Implementation Decisions` 的來源句**

那一節開頭是：

> 已經決定好的技術事項。每一條要指得出它來自哪個**已答的問題或哪條規則**。

問題有三層，要一次解決：
- 這一節的六個 bullet 現在帶 `← <標記>`，而 `<標記>` 的三個值包含 `推論`；
- fixture 的兩條也是 `← 推論`；
- `推論` 的定義是「兩者都不是，是從別的地方推導出來的」，字面上不在「已答的問題或哪條規則」裡。

**改法：把這句話改成講它真正要講的事。** 它想要的不是「來源只能是這兩種」，是
**「每一條都說得出理由，不准憑空」**——那跟來源標記是兩個不同的要求，一個管
可追溯性，一個管標記法。重寫成兩句：一句要求可追溯（並明講 `推論` 合法，
但推論本身要指得出從什麼推來的），一句指向 `## 來源標記：三選一`。

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "已答的問題或哪條規則" skills/bdd-spec/references/spec-format.md
```

改完要零命中。

- [ ] **Step 2: 第 3 項——表格裡的 `← 推論`**

`## Testing Decisions` 一節現在有三列把標記貼在格子值後面：

```
| `application.WorkoutService` | 既有 | 規則全在 domain 與 application，HTTP 只是轉發 ← 推論 |
| 2.1 進度不可倒退 | 單調遞增判斷 | ✗ ← 推論 |
| 6.1 查詢回傳完整內容 | — | ✓ 要讀得到前一次寫入 ← 推論 |
```

兩個毛病：`:279` 說「行尾加 `← <標記>`」，而表格裡它在 `|` 之前不是行尾；`✗ ← 推論`
把「這格的值」和「這個值的來源」黏成一團，讓一個符號欄位的值域變成開放的。

**改法（兩半）：**

1. **Seam 表加第四欄 `來源`**，把標記移進去。這是這個 repo 既有的做法——
   `## Risks` 的表就是用 `從哪條規則長出來` 當欄位，沒有用 `←`。
2. **每個 Example 的測試層級那張表，兩列的標記直接拿掉**，並在
   `## 來源標記：三選一` 的列舉裡，把 `## Testing Decisions` 那一項限縮成
   **只有 Seam 的決定**。理由寫進去：那張表的單位是**場景**不是**決定**，
   而且它的標記永遠只可能是 `推論`——沒有任何 PRD 或澄清答案說得出
   「場景 2.1 不需要資料庫」。**一個不會變的欄位沒有鑑別力**，而鑑別力
   正是 `:285` 那條 MUST 說這個標記存在的理由。

- [ ] **Step 3: 第 6、7 項——兩條沒有範例的規則**

兩條規則各補一個最小範例，就地放在規則底下：

- `:508` 「空著的分類保留標題並寫為什麼空」（`## Known Boundaries`）
- `:648` 「沒有就寫「（無）」，不要刪掉標題」（`## Further Notes`）

補的範例要**示範空的樣子**，不是再示範一次有內容的樣子——這兩條規則講的就是
「沒東西的時候長什麼樣」，而那正是目前整份文件裡沒有任何地方示範過的。

- [ ] **Step 4: 規則 vs 範例全檔自檢**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -c "^MUST\|^MUST NOT\|^SHOULD\|判準" skills/bdd-spec/references/spec-format.md
```

這一步不是數數字。**逐條走過這份檔案裡的每一條 MUST／MUST NOT／SHOULD／判準，
找出它對應的範例，確認範例遵守它。** 這個分支家族至今抓到的缺陷有三個是這一族，
包括本任務第 1 與第 3 項。把你檢查過的清單寫進報告——沒有清單的「檢查過了」
不算檢查過。

- [ ] **Step 5: 驗收與提交**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FA="$(mktemp -d)/fx"; mkdir -p "$FA/specs/2026-09-09-parking" "$FA/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FA/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-spec/scripts/status.py "$FA" | grep "^合計"
python3 skills/bdd-formulation/scripts/check_spec.py "$FA" 2>&1 | tail -1
diff <(grep -o '^| `## [^`]*`' skills/bdd-spec/references/spec-format.md | sed 's/^| `//;s/`$//') \
     <(grep '^## ' skills/bdd-spec/examples/minimal-spec.md) && echo "16/16 對帳 ✓"
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
```

Expected: `合計  2  1  1`；`4 個問題`；`16/16 對帳 ✓`；稽核 exit 0。本任務不碰 fixture，
所以這三個都必須紋風不動。

---

## Task 3: fixture 的三項

`examples/minimal-spec.md` 是 canonical 解析對象，也是下一個人照抄的東西。

**Files:**
- Modify: `skills/bdd-spec/examples/minimal-spec.md`

**Interfaces:**
- Consumes: Task 1 改好的 `check_spec.py`（Step 4 會用它驗）
- Produces: 無

- [ ] **Step 1: 記下基線**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FA="$(mktemp -d)/fx"; mkdir -p "$FA/specs/2026-09-09-parking" "$FA/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FA/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-spec/scripts/status.py "$FA" > /tmp/base-status2.txt 2>&1
python3 skills/bdd-formulation/scripts/check_spec.py "$FA" | sed "s|$FA|FX|g" > /tmp/base-check2.txt 2>&1
grep "^合計" /tmp/base-status2.txt; tail -1 /tmp/base-check2.txt
```

- [ ] **Step 2: 第 2 項——`## Risks` 的欄位**

骨架（`spec-format.md:565`）是 `| 風險 | 從哪條規則長出來 | 影響什麼 |`；
fixture 是 `| 風險 | 影響 | 目前怎麼辦 |`。

**以骨架為準改 fixture。** 骨架那一欄是有理由的：`:561` 要求「每一條指出從哪條規則
長出來」，而 fixture 現在的三欄裡沒有任何一欄裝得下它。現有那一列的內容要重新分配
到三個新欄位，不是硬塞——`消防法規可能要求警報時柵欄常開` 這條風險長自哪條規則
（FR-2），影響什麼（FR-2 的行為要改），都寫得出來。「目前怎麼辦」（Q-2 待答，先不
實作）沒有欄位可放，**判斷它該去哪裡並說明**：可能是 `## Open Questions` 的 Q-2
小節，也可能骨架該有第四欄——如果你認為是後者，**不要自己加欄**，寫進報告讓我裁決。

- [ ] **Step 3: 第 4 項——`## Testing Decisions` 缺 Seam 表**

fixture 那一節現在只有兩個 bullet。骨架要求三塊：**Seam** 表、**每個 Example 的
測試層級**表、**Prior art**。

這是教學用的虛構例子，所以替它寫一個 seam 不是「發明沒人同意過的東西」——
整份 fixture 都是虛構的。但**要寫得像是真的做過那個判斷**：停車場計費的 seam
選在哪一層、是既有的還是新的、為什麼是那一層，三條判準（優先既有、取最高、
理想數量一個）要看得出來被套用過。

`Prior art` 這一塊在虛構專案裡沒有真實路徑可寫——**照骨架自己的規則處理**：
「沒有內容的節保留標題並寫為什麼空」。

Task 2 Step 2 會決定 Seam 表有沒有第四欄 `來源`。**兩個任務都會動到 Seam 表的形狀，
先做完 Task 2 再做這一步**，並照 Task 2 定下的形狀寫。

- [ ] **Step 4: 第 5 項——story 的切法理由**

`SKILL.md:83` 是 MUST：「寫下**為什麼是這個切法**」。fixture 的兩則 story 都只有
「身為⋯我想⋯以便⋯」那一句，沒有切法說明。

兩則各補一句。`spec-format.md` 的 story 範例（約 `:171-191`）**也沒有**——
一併補上，兩份要示範同一個形狀。骨架裡 `SKILL.md:66-78` 有現成的範例可參考
（「涵蓋 FR-1、FR-4。沿規則切：訪客計費是一條獨立成立的約束，自己可驗收。」）。

- [ ] **Step 5: 迴歸——改完解析不准變**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FA="$(mktemp -d)/fx"; mkdir -p "$FA/specs/2026-09-09-parking" "$FA/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FA/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-spec/scripts/status.py "$FA" > /tmp/new-status2.txt 2>&1
python3 skills/bdd-formulation/scripts/check_spec.py "$FA" | sed "s|$FA|FX|g" > /tmp/new-check2.txt 2>&1
diff /tmp/base-status2.txt /tmp/new-status2.txt && echo "status.py 零變化 ✓"
diff /tmp/base-check2.txt  /tmp/new-check2.txt  && echo "check_spec.py 零變化 ✓"
```

Expected: 兩個 `diff` 都無輸出。有輸出就停下來——本任務不該動到任何被解析的東西。

- [ ] **Step 6: 用 Task 1 的新行為驗一次**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FB="$(mktemp -d)/fx"; mkdir -p "$FB/specs/2026-09-09-parking"
cp skills/bdd-spec/examples/minimal-spec.md "$FB/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-formulation/scripts/check_spec.py "$FB"; echo "exit=$?"
```

Expected: `全部通過`，`exit=0`。**canonical fixture 在 `bdd-spec` 交件的那一刻必須是綠的**
——它是示範「一份合格的 `spec.md` 長什麼樣」的那一份。

- [ ] **Step 7: 16/16 對帳與提交**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
diff <(grep -o '^| `## [^`]*`' skills/bdd-spec/references/spec-format.md | sed 's/^| `//;s/`$//') \
     <(grep '^## ' skills/bdd-spec/examples/minimal-spec.md) && echo "16/16 ✓"
grep -c '^## ' skills/bdd-spec/examples/minimal-spec.md
```

Expected: `diff` 無輸出，`16`。

---

## Task 4: 兩個 skill ＋ 兩份文件的小項

三項獨立，**各自一筆提交**。

**Files:**
- Modify: `skills/bdd-spec/SKILL.md`（第 8 項）
- Modify: `skills/bdd-formulation/SKILL.md`、`PLAN.md`（第 9 項）
- Modify: `README.md`（第 10 項）

**Interfaces:**
- Consumes: 無
- Produces: 無

- [ ] **Step 1: 第 8 項——`SKILL.md:222` 的窄版重述**

那一段（`### 4. 寫 spec.md` 的第一條 MUST）結尾寫著：

> ⋯`## Open Questions` 每一題的「該問誰」都要指向這張表裡的人，兩邊一起編的話那條規則會輕鬆通過而毫無意義。

`spec-format.md:98` 已經是通則了——「文件裡**任何一個**「該問誰」都必須指向這張表裡的
某個關係人，**不限** `## Open Questions`」。這裡的窄版仍然為真，但它是規範的**子集**
而沒有指路，而且它是寫步驟 4 的 agent 先讀到的那一份。

**改法：保留這一段的論證**（不要兩邊一起編，否則規則會空轉），**把規則本身換成指路**。

- [ ] **Step 2: 第 9 項——`bdd-spec-review` 指向不存在的 skill**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
ls -d skills/bdd-spec-review 2>&1 | head -1
grep -rn "bdd-spec-review" skills/ PLAN.md
```

兩個 skill 的 Skill Boundaries 都寫「要稽核既有 `.feature` 寫得好不好 → 改用 `bdd-spec-review`」，
而它不存在；`PLAN.md:204` 老實列它為未實作。

**不要新建那個 skill**——本計畫是清理，不是開新功能。**改成誠實的形式**：
兩行各自改寫，說明這件事目前沒有 skill 接，以及在那之前該怎麼辦
（`bdd-formulation` 的 `examples/anti-patterns.md` 是現成的自檢清單，
`check_spec.py` 的十三項裡有五項就是在稽核 `.feature`）。`PLAN.md:204` 那一列不動——
它本來就誠實。

- [ ] **Step 3: 第 10 項——`README.md` 的 skill 樹**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -n '28,38p' README.md; ls skills/
```

樹裡列了 7 個，`skills/` 有 8 個，漏掉 `example-mapping/`。補上，位置與縮排照相鄰行。

**`README.md` 的 skill 樹是照 pipeline 順序列的**（bdd-discovery → bdd-spec → bdd-formulation → bdd-plan，然後才是支援性的那幾個），**不是字母序**。`example-mapping` 要放進它在流程裡的位置，不要為了讓比對指令好寫就把整棵樹重排。下面 Step 4 的比對是**集合比對**（兩邊都 `sort` 過），刻意不管順序。

**順手核對整棵樹**——漏一個就可能漏第二個，其餘每一行都要對得上 `ls skills/` 的結果。

- [ ] **Step 4: 驗收**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for d in skills/*/; do printf "%-18s" "$(basename $d)"; \
  python3 skills/skill-rules/scripts/audit_skill.py "$d" >/dev/null 2>&1 && echo "exit 0" || echo "exit $?"; done
echo "--- skill 樹對不對得上 ---"
diff <(awk '/├── skills\//{f=1;next} f&&/^├── /{exit} f' README.md | grep -oE '[a-z-]+/' | sort) \
     <(ls -d skills/*/ | xargs -n1 basename | sed 's|$|/|' | sort) \
  && echo "README 的 skill 樹與 skills/ 是同一組 ✓"
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
claude plugin validate . 2>&1 | tail -3
```

Expected: 8 個 skill 全 `exit 0`；`diff` 無輸出；`validate` 恰好 1 個警告。

---

## Self-Review

**1. 覆蓋**

十一項全部有歸屬：Task 1 收第 11 項，Task 2 收第 1、3、6、7 項，Task 3 收第 2、4、5 項，
Task 4 收第 8、9、10 項。

**2. 佔位字掃描**

沒有 TBD／TODO。Task 1 的程式碼改動給了改前與改後的完整片段。Task 2、3、4 有幾處
給的是**要達到的狀態與判準**而非逐字稿（重寫來源句、寫 Seam 表、補切法理由）——
那幾處都需要配合上下文的語氣，逐字指定會產出接不上的段落；每一處都附了驗證指令，
而且有兩處明講「不要自己決定，寫進報告讓我裁決」。

**3. 任務之間的相依**

只有一處：**Task 2 Step 2 決定 Seam 表有沒有第四欄，Task 3 Step 3 要照那個形狀寫。**
已在 Task 3 Step 3 明寫「先做完 Task 2 再做這一步」。其餘任務彼此無關，順序隨意。

**4. 我在這份計畫裡最可能犯的錯**

前兩份計畫我犯過三次同一種：**驗收預期與自己 mandate 的施工對不上**（預測兩行實際六行、
叫人改不存在的行號、把 reviewer 舉的四個例子當成完整的九項）。這一份的對策：

- 每一個數字都是剛才實際量過的（736 行、182 行、16 節、292 行、八項／五項）。

**而我在寫完之後的自我複查裡，又犯了同一族兩次。** 第一次：我把 `spec-format.md` 寫成 660 行、
`minimal-spec.md` 寫成 179 行、該問誰通則的行號寫成 `:97`——全是憑印象寫的，實際是 736、182、`:98`，
因為計畫二又往那兩個檔加了東西。第二次更典型：Task 4 Step 4 的 README 比對指令，我寫的正則
`^│   [├└]── [a-z-]+/` 會把 `benchmark/` 底下的 `cases/` 與 `skeleton/` 一起掃進來，產生兩個假陽性；
改成只掃 `skills/` 那一段之後，我又要求「順序也一致」——但 README 的 skill 樹是照 pipeline 順序排的，
不是字母序，那個要求會逼執行者把整棵樹重排。三次才寫對。**驗收指令本身也要跑過一次才算數。**
- Task 1 Step 1 和 Step 2 先**證明缺陷存在**再修——先紅燈，而且紅燈的內容要具體
  （「壞掉的 AC 編號沒有被提到」），不是「exit 1」。
- Task 3 Step 1 的基線由執行者**自己量**，不是我寫死。
- Task 2 Step 4 與 Task 4 Step 3 都明講「不是數數字」「漏一個就可能漏第二個」，
  並要求把檢查過的清單寫進報告。

**5. 一個我沒有解決的問題**

`check_spec.py` 降級之後，一份乾淨的、還沒有 `.feature` 的專案會 exit 0。`CLAUDE.md`
說這兩支腳本「exit 1 and name what is missing rather than inferring anything from zero
input」，那句話的字面在新行為下不再完全成立。Task 1 Step 8 要求執行者重讀那句並
**說出判斷**，而不是直接改掉它——因為那句話背後的原則（fail-loud，不要安靜地說一切正常）
在新行為下仍然成立，改的只是「什麼算 zero input」。這一處是這份計畫裡我最不確定的地方。
