# `spec.md` 加四節 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 讓 `spec.md` 自帶「這個 feature 為誰做」「技術上准用什麼、禁用什麼」「哪些事自己做、哪些要先問、哪些絕不做」，並把 `bdd-spec` 的 seam 步驟補齊成 D6 要求的形狀。

**Architecture:** `references/spec-format.md`（骨架）與 `examples/minimal-spec.md`（canonical 解析對象）是**一個原子單位**——格式加了一節而 fixture 沒有，等於留下一份示範違反自己規格的範例，所以兩者同一筆提交。`bdd-spec/SKILL.md` 隨後補上 seam 的兩條守則與節數。最後清掉計畫一留下的四項 follow-up。

**Tech Stack:** Markdown skill 文件、Python 3 稽核腳本（`audit_skill.py`、`check_spec.py`、`status.py`）、`claude plugin validate`。沒有 CI，全部手跑。

**Spec:** `docs/superpowers/specs/2026-09-14-grilling-and-formulation-design.md`（D5 與 D6 是本計畫的依據；「計畫二」那一節是它的範圍聲明）

## Global Constraints

- **AC 的層級不動。** 這一份只加章節，不改既有結構——`#### FR` 底下掛 `- **AC-<n>.<m>**` ＋ Given/When/Then 三行，維持現狀。
- **接縫不變：** `FR-<n>` ↔ `@rule-<n>`、`AC-<n>.<m>` ↔ `@example-<n>.<m>`。
- **既有十三節的標題文字一字不改，順序不動。** 新節插進去，舊節不改名、不重排、不合併。改名會動到 `spec-format.md` 內文每一處引用該節名的地方，以及 fixture、`bdd-spec/SKILL.md`、`bdd-formulation/SKILL.md`——那是另一份計畫的規模。
- **腳本的解析行為必須零變化。** `check_spec.py` 與 `status.py` 用 `(?=^## |\Z)` 當節邊界，只讀 `## User Stories` 與 `## Open Questions`。加節之後對 fixture 的輸出要**逐字節相同**。這是本計畫最重要的迴歸護欄。
- 散文用繁體中文；commit message 用英文，conventional-commit 標題加 `WHAT:` / `WHY:` / `HOW:` 正文，結尾兩行：
  `Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>`
  `Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN`
- 每個 skill 目錄底下只准有 `rules/`、`references/`、`examples/`、`scripts/`（S3，MUST）；每個子目錄檔案都要被同一個 skill 指到（S5，MUST）；`SKILL.md` 本文超過 500 行拿 S10（SHOULD）。
- 提交前清 `__pycache__`：`find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null`

### 基線（動手前已量過，收工要對得上）

| 量什麼 | 現值 |
| --- | --- |
| `audit_skill.py` 跑 8 個 skill | 全部 exit 0，沒有任何 S10 |
| `claude plugin validate .` | 通過，**恰好 1 個警告**（plugin root 的 `CLAUDE.md`，刻意的取捨） |
| `status.py` 對 fixture | `合計  2  1  1`，exit 0 |
| `check_spec.py` 對 fixture | `4 個問題`，exit 1 |
| `spec-format.md` | 633 行；「章節順序」表 **13** 列 |
| `examples/minimal-spec.md` | 152 行；`grep -c "^## "` ＝ **13** |
| `bdd-spec/SKILL.md` | 280 行；「十三節」出現 **2** 次（第 4 行 frontmatter、第 243 行產物格式圖） |

**fixture** 的建法固定，每個任務都用它：

```bash
FX="$(mktemp -d)/fx"
mkdir -p "$FX/specs/2026-09-09-parking" "$FX/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FX/specs/2026-09-09-parking/spec.md"
```

`features/` 是空的，所以 `check_spec.py` 必然報 4 個問題。**那 4 個問題的文字在本計畫前後必須一字不差**——它證明加節沒有改變任何解析。

---

## 已決定的事（不要重新討論）

**16 節，不是 15。** 「技術偏好與限制」自己一節。使用者在 2026-09-15 明確選了這個。理由：他當初列的四件事裡 `Tech stack preferences and constraints` 就是跟 `Known boundaries` 並列的一項；D5 的抱怨本來就是「散在兩邊」，不給它一個家等於沒解決。

**四件新東西的落點：**

| D5 的項目 | 落在哪 |
| --- | --- |
| 目標使用者 | **不是新節**——`## Personas` 開頭加一行「這個 feature 為誰做」 |
| 技術偏好與限制 | 新節 `## Tech Preferences & Constraints`，排在 `## Assumptions / Constraints` 之後 |
| Known Boundaries | 新節 `## Known Boundaries`，排在上者之後 |
| Further Notes | 新節 `## Further Notes`，排在**最後**，在 `## Open Questions` 之後 |

**改完的章節順序（16 節）：**

```
 1 ## Document Overview
 2 ## Background
 3 ## Goal
 4 ## Scope — In / Out
 5 ## Personas                        ← 開頭加一行
 6 ## User Stories
 7 ## Non-Functional Requirements
 8 ## Assumptions / Constraints
 9 ## Tech Preferences & Constraints  ★ 新
10 ## Known Boundaries                ★ 新
11 ## Success Metrics
12 ## Implementation Decisions
13 ## Testing Decisions
14 ## Risks
15 ## Open Questions
16 ## Further Notes                   ★ 新
```

**`## Further Notes` 排在 `## Open Questions` 之後，是刻意的。** 其餘十五節順序完全不動，新的兩節插在第 8 與第 11 之間。`## Open Questions` 是機械可解析的索引，`status.py` 用 `(?=^## |\Z)` 找它的結尾——在它後面加一節等於給它一個明確的結尾標記，比讓它靠檔尾 `\Z` 結束更穩。

---

## File Structure

| 檔案 | 這份計畫對它做什麼 |
| --- | --- |
| `skills/bdd-spec/references/spec-format.md` | 章節順序表 13→16 列；Personas 一節加一行說明；新增三節的格式規範；新增一節講 `Tech Preferences & Constraints` 與 `Assumptions / Constraints` 怎麼分邊 |
| `skills/bdd-spec/examples/minimal-spec.md` | 跟著加三節 ＋ Personas 那一行。**與上者同一筆提交** |
| `skills/bdd-spec/SKILL.md` | seam 步驟補 D6 的兩條；「十三節」→「十六節」兩處；步驟 4 指向新節 |
| `skills/bdd-spec/references/persona-definition.md` | 計畫一的 follow-up：`前置（狀態）` 這個詞與 `bdd-spec` 這個歸屬都過期了 |
| `PLAN.md` | 計畫一的 follow-up：未決事項 #13 的行號引用偏移兩行 |
| `docs/ai-sdlc.md` | 計畫一的 follow-up：`:82` 起的邊界論證少了 `bdd-spec`↔`bdd-formulation` 這一條 |
| `skills/bdd-formulation/scripts/check_spec.py` | 計畫一的 follow-up：docstring 開頭「十三件事」、結尾「上面十二項」。**只改註解，不改一行程式碼** |

---

## Task 1: `spec-format.md` ＋ `minimal-spec.md` 加三節

這兩個檔是一個原子單位，**同一筆提交**。格式加了一節而 canonical 範例沒有，等於附上一份示範違反自己規格的範例——而那份範例正是下一個人照抄的東西。

**Files:**
- Modify: `skills/bdd-spec/references/spec-format.md`
- Modify: `skills/bdd-spec/examples/minimal-spec.md`

**Interfaces:**
- Consumes: 無（第一個任務）
- Produces: 三個新節名 `## Tech Preferences & Constraints`、`## Known Boundaries`、`## Further Notes`，Task 2 的 `SKILL.md` 要指向它們；16 這個數字，Task 2 要寫進 frontmatter 與產物格式圖

- [ ] **Step 1: 記下基線，等一下要證明它沒變**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FX="$(mktemp -d)/fx"
mkdir -p "$FX/specs/2026-09-09-parking" "$FX/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FX/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-spec/scripts/status.py "$FX" > /tmp/base-status.txt 2>&1; echo "status exit=$?"
python3 skills/bdd-formulation/scripts/check_spec.py "$FX" | sed "s|$FX|FX|g" > /tmp/base-check.txt 2>&1
echo "check exit=${PIPESTATUS[0]}"
grep "^合計" /tmp/base-status.txt; tail -2 /tmp/base-check.txt
echo "$FX" > /tmp/fx-path.txt
```

Expected: `status exit=0`，`合計` 那行是 `2     1     1`；`check exit=1`，最後是 `4 個問題`。

- [ ] **Step 2: 改「章節順序」表——13 列變 16 列**

`skills/bdd-spec/references/spec-format.md` 第 17–30 行是那張表。在 `` | `## Assumptions / Constraints` | … | `` 那一列**之後**插入兩列：

```markdown
| `## Tech Preferences & Constraints` | 語言、框架、不准引入的依賴——**我們自己**的技術選型偏好與禁令 |
| `## Known Boundaries` | `總是要做`／`要先問`／`絕對不做` 三張清單 |
```

並在表的**最後一列** `` | `## Open Questions` | … | `` 之後插入：

```markdown
| `## Further Notes` | 放不進其他節、但下一個接手的人會想知道的 |
```

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -n '15,34p' skills/bdd-spec/references/spec-format.md
grep -c '^| `##' skills/bdd-spec/references/spec-format.md
```

Expected: 表變成 16 列，`grep -c` 回 **16**；`## Assumptions / Constraints` 之後是 `## Tech Preferences & Constraints` 再 `## Known Boundaries` 再 `## Success Metrics`；`## Open Questions` 之後是 `## Further Notes`。

- [ ] **Step 3: `## Personas` 一節加「這個 feature 為誰做」**

`spec-format.md` 的 `## Personas：在三欄上加兩欄` 一節（約第 121 行起）。在那段 markdown 範例**之前**插入：

```markdown
一節開頭先寫一行，說明**這個 feature 是為誰做的**——不是「誰在流程裡出現」：

```markdown
**這個 feature 為誰做**：P-1（訪客）。P-2 出現在流程裡但不是這次的受益者 ← 推論
```

兩件事分開寫的理由：**出現在流程裡的角色，不一定是這次要讓誰過得更好。**
月租戶會刷卡出場，所以他在流程裡；但這個 feature 要解的是訪客算不出錢的問題。
不分開的話，story 會開始為「流程裡出現過的每一個人」各切一則，而其中幾則
沒有人在等它。
```

**注意**：上面那段文字本身含一個三引號區塊，寫進檔案時**內外都要保留**——外層是說明文字的一部分，不是本計畫的圍欄。

- [ ] **Step 4: 新增 `## Tech Preferences & Constraints` 的格式規範**

在 `spec-format.md` 的 `## Non-Functional Requirements 與 Constraints：怎麼分邊` 一節**之後**（約第 419 行、`## Implementation Decisions` 之前）插入一整節：

```markdown
## Tech Preferences & Constraints

三個具名區塊，每條照樣標來源：

```markdown
## Tech Preferences & Constraints

**語言與框架**
- Go 1.23，驗收層用 godog ← PRD §1

**不准引入的依賴**
- 不引入需要商業授權的套件 ← Q-6

**既有資產要沿用的**
- HTTP 層沿用既有的 chi router，不另起一套 ← 推論
```

**跟 `## Assumptions / Constraints` 怎麼分邊：限制的來源在不在我們的控制範圍內。**

| 這一條 | 進哪一節 |
| --- | --- |
| 區權會否決車牌辨識 | `## Assumptions / Constraints`——外部決定的，我們改不了 |
| 已發包的閘門硬體只吃某種訊號 | `## Assumptions / Constraints`——外部既成事實 |
| 不准引入需要授權費的套件 | `## Tech Preferences & Constraints`——我們自己的規矩 |
| 沿用既有的 router 不另起 | `## Tech Preferences & Constraints`——我們自己的選擇 |

判準是**誰有權改它**，不是它有多硬。外部限制改不了，只能記下來；自己的偏好
改得了，所以要寫清楚「現在的決定是什麼」，讓下一個人知道破例要先問誰。

MUST NOT: 在這一節寫「已經決定好的技術事項」——那是 `## Implementation Decisions`。
判準：這一條是在說**這個 feature 要怎麼做**（決定），還是在說**任何做法都不准
越過的線**（偏好與禁令）？
```

**注意**：這一節的文字含一個三引號區塊，寫進檔案時內外都要保留。

- [ ] **Step 5: 新增 `## Known Boundaries` 的格式規範**

緊接在 Step 4 那一節之後插入：

```markdown
## Known Boundaries

三張清單。**這一節等於每個 feature 自帶一份 agent 約束**——它是這一版最值得加的東西。

```markdown
## Known Boundaries

**總是要做**
- 金額一律用整數分計算，不用浮點 ← Q-7

**要先問**
- 動到既有的出場閘門控制流程 ← 推論（那是已發包的設備）

**絕對不做**
- 寫入車牌辨識相關的任何欄位 ← PRD §5（區權會否決）
```

**三個分類的判準是誰承擔後果：** 自己做、問了再做、做了要負責的事不做。

| 分類 | 意思 | 判準 |
| --- | --- | --- |
| `總是要做` | 不用問就照做 | 做錯的後果落在這個 feature 內，而且改得回來 |
| `要先問` | 可以做，但動手前要有人點頭 | 做錯會影響這個 feature 以外的東西 |
| `絕對不做` | 不論理由 | 已經有人否決過，或後果不可逆 |

MUST: 每一條照樣標來源。`要先問` 的條目要寫**該問誰**——寫「要先問」而不寫問誰，
等於把一件事掛起來而沒有人接。

MUST NOT: 把 `## Scope — In / Out` 的「明確排除的能力」搬進 `絕對不做`。
兩者不同：Out of Scope 是**這次不做**，`絕對不做` 是**做了要負責**。
下一期可能會做的東西進 Out of Scope；下一期也不准做的才進這裡。

**空著的分類保留標題並寫為什麼空。** 三張清單都空的 feature 是存在的（純查詢、
沒有副作用），但那要寫出來——空清單與沒想過在檔案上長得一樣。
```

**注意**：這一節的文字含一個三引號區塊，寫進檔案時內外都要保留。

- [ ] **Step 6: 新增 `## Further Notes` 的格式規範**

在 `spec-format.md` 的 `## Open Questions：表是索引，小節是決策史` 與 `## Open Questions` 兩節**之後**、`## 角色與詞義：兩條界線` 之前插入：

```markdown
## Further Notes

放不進其他節、但下一個接手的人會想知道的。沒有就寫「（無）」，不要刪掉標題。

典型內容：試過但放棄的做法與放棄的理由、這個 feature 跟另一個進行中的工作
會撞在哪、某個決定當時的氣氛（誰有疑慮、後來為什麼還是這樣定）。

MUST NOT: 把應該進其他節的東西丟進來。這一節是**最後手段**，不是收納箱——
一條規則進了這裡，下游沒有任何東西會去讀它。判準：說得出它該進哪一節嗎？
說得出就進那一節。
```

- [ ] **Step 7: fixture 跟著改——`minimal-spec.md`**

`skills/bdd-spec/examples/minimal-spec.md`：

**(a)** `## Personas` 標題的下一行（現在是空行接 `### P-1  訪客`）之間插入：

```markdown
**這個 feature 為誰做**：P-1（訪客）。P-2 出現在流程裡但不是這次的受益者 ← 推論
```

**(b)** `## Assumptions / Constraints` 那一節之後、`## Success Metrics` 之前插入：

```markdown
## Tech Preferences & Constraints

**語言與框架**
- （無資料）——需求方沒有指定，技術選型由實作端決定 ← 推論

**不准引入的依賴**
- 不引入需要連外網的服務，停車場網路是封閉的 ← PRD §5

**既有資產要沿用的**
- 出場閘門的既有控制介面沿用，不另起 ← PRD §5

## Known Boundaries

**總是要做**
- 金額一律用整數元計算，不用浮點 ← 推論

**要先問**
- 動到既有的出場閘門控制流程 —— 問管委會主委 ← PRD §5（已發包的設備）

**絕對不做**
- 寫入車牌辨識相關的任何欄位 ← PRD §5（去年區權會否決）
```

**(c)** 檔案最後、`## Open Questions` 的表格之後插入：

```markdown

## Further Notes

免費時段的長度（30 分）是沿用隔壁社區的做法，不是算出來的。之後若要調整，
Q-4 的答案（30 分整仍算免費）是跟著這個數字定的，兩者要一起改。 ← 推論
```

- [ ] **Step 8: 節數對帳**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
echo "spec-format 章節順序表: $(grep -c '^| `##' skills/bdd-spec/references/spec-format.md)"
echo "minimal-spec 標題數:    $(grep -c '^## ' skills/bdd-spec/examples/minimal-spec.md)"
diff <(grep -o '^| `## [^`]*`' skills/bdd-spec/references/spec-format.md | sed 's/^| `//;s/`$//') \
     <(grep '^## ' skills/bdd-spec/examples/minimal-spec.md) && echo "順序與名稱完全一致"
```

Expected: 兩個數字都是 **16**，`diff` 無輸出並印出 `順序與名稱完全一致`。
這一條比數數字強：它同時驗了**名稱拼法**與**順序**——章節順序表宣告的東西，
fixture 要一字不差地照著長。

- [ ] **Step 9: 迴歸——腳本的輸出必須一字不差**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FX="$(mktemp -d)/fx"
mkdir -p "$FX/specs/2026-09-09-parking" "$FX/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FX/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-spec/scripts/status.py "$FX" > /tmp/new-status.txt 2>&1; S=$?
python3 skills/bdd-formulation/scripts/check_spec.py "$FX" | sed "s|$FX|FX|g" > /tmp/new-check.txt 2>&1
C=${PIPESTATUS[0]}
diff /tmp/base-status.txt /tmp/new-status.txt && echo "status.py 輸出零變化 (exit $S，基線 0)"
diff /tmp/base-check.txt  /tmp/new-check.txt  && echo "check_spec.py 輸出零變化 (exit $C，基線 1)"
```

Expected: 兩個 `diff` 都無輸出，兩行確認訊息都印出來，`exit 0` 與 `exit 1`。

**有輸出就停下來。** 那代表加節改變了解析——最可能的原因是新節插在 `## User Stories`
或 `## Open Questions` 中間，把它們的 `(?=^## |\Z)` 邊界提前了。回 Step 2 檢查落點。

- [ ] **Step 10: 稽核**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
```

Expected: `✔ 通過`，`exit=0`。`spec-format.md` 變長不影響 S10——S10 只量 `SKILL.md` 本文。

- [ ] **Step 11: 提交**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
git add skills/bdd-spec/references/spec-format.md skills/bdd-spec/examples/minimal-spec.md
git commit -F - <<'MSG'
feat(bdd-spec): spec.md gains tech constraints, known boundaries and further notes

WHAT: Add three sections to spec.md's skeleton -- Tech Preferences &
Constraints, Known Boundaries, Further Notes -- and one line at the top of
Personas naming who the feature is actually for. Thirteen sections become
sixteen. The canonical fixture gains all four in the same commit.

WHY: Known Boundaries is the one that earns its place: 總是要做 / 要先問 /
絕對不做 is a per-feature agent contract, and an agent that has one stops
guessing which changes it may make unasked. The other two close gaps the
document already had -- technology preferences were scattered across
Implementation Decisions and Assumptions/Constraints with no single home, and
anything that fit no section was simply lost. Personas listed everyone in the
flow, which quietly produced a story per participant, some of which nobody was
waiting for.

HOW: Format and fixture in one commit, because a skeleton whose own example
disobeys it ships a violation as the thing people copy. The three new sections
sit between Assumptions/Constraints and Success Metrics, and after Open
Questions, so no existing section moves and the two the scripts parse keep
their boundaries; verified by diffing status.py and check_spec.py output
against the pre-change run, byte for byte. Section names and order are
reconciled between the skeleton table and the fixture by diff, not by counting.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
MSG
git log --oneline -1
```

---

## Task 2: `bdd-spec/SKILL.md` —— seam 步驟補齊 ＋ 節數 ＋ 指向新節

**Files:**
- Modify: `skills/bdd-spec/SKILL.md`

**Interfaces:**
- Consumes: Task 1 定下的三個新節名與 16 這個數字
- Produces: 無下游任務依賴

- [ ] **Step 1: 「十三節」→「十六節」，兩處**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "十三節" skills/bdd-spec/SKILL.md
sed -i 's/十三節/十六節/g' skills/bdd-spec/SKILL.md
grep -n "十六節" skills/bdd-spec/SKILL.md
grep -c "十三節" skills/bdd-spec/SKILL.md || echo "十三節 零殘留"
```

Expected: 改前 2 行（第 4 行 frontmatter、第 243 行產物格式圖），改後同樣 2 行變成「十六節」，`十三節` 零命中。

- [ ] **Step 2: seam 步驟補 D6 的兩條**

`### 3. 決定 seam —— 驗收測試打在哪一層` 一節。現有的三條規則表與 `Ask user` 段落**保留不動**。

**(a)** 三條規則的表格底下加第四條：

```markdown
| 需要新 seam 就在**你能做到的最高點**提出 | 往下挪一層通常是為了好寫，代價是驗到的東西離使用者更遠 |
```

**(b)** 在 `IMPORTANT: seam 一旦寫進 `spec.md` 就往下傳` 那一段**之前**插入：

```markdown
MUST NOT: 把 seam 的討論變成需求訪談。這一次徵詢問的是**技術判斷**——
「打這一層對嗎」——不是「這個規則該怎麼定」。談的時候發現需求有洞，
寫進 `## Scope — In / Out` 的「回 CLARIFY 補問」，**不要順手問掉**。

理由是這一步沒有 `clarify-log.md` 的記錄機制：訪談的問答會被寫進 log、
被後面每一步讀到；在 seam 討論裡順口問到的答案只活在這一次對話裡，
下一個人打開 `spec.md` 只會看到一條沒有出處的規則。
```

- [ ] **Step 3: 步驟 4 指向三個新節**

`### 4. 寫 `spec.md`` 一節。在 `MUST NOT: 在這一步做新決定。` 那一段**之前**插入：

```markdown
MUST: `## Known Boundaries` 的三張清單要填，空的也要寫為什麼空。
這一節是下游每個 agent 讀到的約束——`總是要做` 省了它一次確認，
`要先問` 擋下一次它不該自己決定的動作，`絕對不做` 是唯一寫得下
「這件事有人否決過」的地方。三者都空而不解釋，等於宣告這個 feature
沒有邊界，而那幾乎不會是真的。

MUST: `要先問` 的每一條要寫**該問誰**，指向 `## Document Overview` 的
`核心關係人` 表裡的人。寫「要先問」而不寫問誰，等於把一件事掛起來而沒有人接。
```

- [ ] **Step 4: 負向檢查 ＋ 稽核**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
echo "--- 三個新節名有沒有在 SKILL.md 被提到 ---"
grep -n "Known Boundaries\|Tech Preferences\|Further Notes" skills/bdd-spec/SKILL.md
echo "--- seam 兩條新規則 ---"
grep -n "最高點\|變成需求訪談" skills/bdd-spec/SKILL.md
echo "--- 節數 ---"; grep -n "十六節" skills/bdd-spec/SKILL.md
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
wc -l skills/bdd-spec/SKILL.md
```

Expected: `Known Boundaries` 至少 2 次（步驟 4 的兩條 MUST）；`最高點` 與 `變成需求訪談` 各 1 次；`十六節` 2 次；稽核 `✔ 通過` exit 0，**沒有 S10**（280 行加十幾行仍遠低於 500）。

- [ ] **Step 5: 提交**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
git add skills/bdd-spec/SKILL.md
git commit -F - <<'MSG'
feat(bdd-spec): fence the seam conversation and require Known Boundaries

WHAT: Add a fourth seam rule (propose a new seam at the highest point you
can reach), a MUST NOT against letting the seam conversation turn into
requirements interviewing, and two MUSTs in step 4 requiring Known Boundaries
to be filled and every 要先問 entry to name who to ask. 十三節 becomes 十六節
in the two places that state it.

WHY: The seam check is this skill's only conversation with a human, which
makes it the one place requirements can enter without passing through
CLARIFY. Nothing recorded them if they did: an answer given in passing during
a seam discussion lives only in that conversation, and the next reader finds a
rule with no source. Known Boundaries is worth a MUST because an unfilled one
is indistinguishable from a feature that genuinely has no boundaries, and that
is almost never true.

HOW: The existing three seam rules and the Ask user paragraph are untouched;
the additions sit around them. The 要先問 rule points at Document Overview's
核心關係人 table, which is the only place in the document that says who can
decide what -- without that link, 要先問 parks work on nobody.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
MSG
```

---

## Task 3: 清掉計畫一留下的四項 follow-up

四件事，彼此獨立，同一個任務一起做以省下四次派工。**各自一筆提交**，因為它們動的是四個不相干的檔案。

**Files:**
- Modify: `skills/bdd-spec/references/persona-definition.md:117` 附近
- Modify: `PLAN.md` 未決事項 #13
- Modify: `docs/ai-sdlc.md:82` 起那一節
- Modify: `skills/bdd-formulation/scripts/check_spec.py` 的 docstring（**只動註解**）

**Interfaces:**
- Consumes: 無
- Produces: 無

- [ ] **Step 1: `persona-definition.md` —— 兩個過期的說法**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -n '114,125p' skills/bdd-spec/references/persona-definition.md
```

那一段現在寫：

> `bdd-spec` 寫「前置（狀態）」規則時，主詞應該是 `spec.md` 的 `## Personas` 裡的角色名。出現一個沒在 `## Personas` 的角色，就是這裡漏了一個——**那是可機械偵測的**。

兩處過期，一處沒有根據：

1. `bdd-spec` 不再寫「前置（狀態）」規則——四象限與 `rule-taxonomy.md` 都在 `bdd-formulation`。歸屬要改成 `bdd-formulation`。
2. 「前置（狀態）」這個象限名在 `bdd-spec` 底下已經沒有定義了。要嘛指向 `../bdd-formulation/references/rule-taxonomy.md`（跨 skill 相對路徑是這個 repo 的既有寫法，見 `skills/bdd-discovery/SKILL.md` 指 `../bdd-spec/references/persona-definition.md`），要嘛改寫成不需要那個詞。
3. **「那是可機械偵測的」目前沒有任何東西在做。** `check_spec.py` 不解析 `## Personas`。改寫成它實際的樣子：這是**人**在寫的時候要對的，不是腳本會抓的。留一句沒有東西在履行的機械保證，比不留更糟。

改完確認：

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "前置（狀態）\|可機械偵測\|bdd-spec 寫" skills/bdd-spec/references/persona-definition.md
grep -n "Personas" skills/bdd-formulation/scripts/check_spec.py || echo "(check_spec.py 確實不碰 Personas)"
```

提交訊息標題：`docs(bdd-spec): persona-definition cites a skill and a rule that both moved`

- [ ] **Step 2: `PLAN.md` 未決事項 #13 的行號引用**

計畫一在這一項上方插入了兩行，所以它每一個自我引用的行號都**偏移兩行**。
**不是四個，是九個**——計畫一的複審只點名了前四個。

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "^13\. " PLAN.md
sed -n '479,496p' PLAN.md
```

**修法：刪掉行號，改用章節名。** 不要重新對齊——行號在這份文件裡每次編輯都會
再爛一次，而章節名不會。每一個引用都對得到一個真實存在的標題，已經查過：

| #13 現在寫 | 它想指的東西 | 改成 |
| --- | --- | --- |
| 六步流程表（**63 行**，仍列 `brief.md`…） | 行 61 | `## 六步流程` 那張表 |
| `` `## 產物佈局`（**75-133 行**，用現在式…） | 行 77 | `` `## 產物佈局` ``（刪掉行號，名字已經在了） |
| `spec.md` 的節表（**108-129 行**） | 行 110 | `` ### `spec.md` 的節 `` |
| **147-159**（「Slice——一批」） | 行 153 | `` ### 1. Slice ——「一批」 `` |
| `步驟 → skills 對照` | 行 192 | 已有名字，不動 |
| **215-227**（Example Mapping 訊號） | 行 220 | `### Example Mapping 的兩個訊號在哪裡看` |
| **229-247**（追問覆蓋…與 `status.py` 矛盾） | 行 234 | `### 追問覆蓋改成算出來的` |
| `` `命名與結構的決定`（**264 行**，一個 `<slice-slug>`） `` | 節在行 254 | `` `## 命名與結構的決定` ``（刪掉行號） |
| **289-296**（可追溯鏈） | 行 294 | `## 可追溯鏈` |

改完之後，這一項裡**不該剩下任何行號**，而且每一個被引用的章節都要存在：

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
S=$(grep -n "^13\. " PLAN.md | cut -d: -f1)
E=$(awk -v s="$S" 'NR>s && /^$/{print NR; exit}' PLAN.md)
echo "--- 還剩幾個行號（應為 0）---"
sed -n "${S},${E}p" PLAN.md | grep -oE '[0-9]+(-[0-9]+)? 行|（[0-9]+(-[0-9]+)?）' || echo "0 個"
echo "--- 每個被引用的章節存不存在（每行右邊都要是 1）---"
sed -n "${S},${E}p" PLAN.md | grep -oE '#{2,3} [^`）」，、。]+' | sort -u | while read -r n; do
  printf "%-45s %s\n" "$n" "$(grep -c "^$n\$" PLAN.md)"
done
```

Expected: 行號 `0 個`；每個章節名右邊都是 **1**。出現 0 代表引用了不存在的節，
出現 2 以上代表那個標題在檔案裡重複，引用不唯一——兩種都要回頭改。

提交訊息標題：`docs: 未決事項 #13 cited its own file by line number, and the lines moved`

- [ ] **Step 3: `docs/ai-sdlc.md` 補上第四條邊界的理由**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -n '80,100p' docs/ai-sdlc.md
```

那一節的標題已經說「四個 skill」，但內文只論證了兩條邊界（Discovery↔Formulation、PLAN）。**`bdd-spec`↔`bdd-formulation` 這一條——剛切出來的那條——在這份負責記錄理由的文件裡沒有理由。**

補一段，論點取自設計 spec 的 D8：`.feature` 是三個朋友真正拿來對答案的產物，它需要一個擁有者；而封閉步驟文法、tag 與覆蓋盤點三件事全部掛在產 `.feature` 的那一邊，跟「把問答綜合成 `spec.md`」是兩種完全不同的判斷。不要複製 D8 的整段——用這份文件自己的語氣寫兩三句。

**注意**：`docs/bdd.md`、`docs/tdd.md`、`docs/sdd.md` 轉錄外部來源，**不得**寫入專案意見。要改的只有 `docs/ai-sdlc.md`。

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
git diff --stat -- docs/bdd.md docs/tdd.md docs/sdd.md; echo "(空=沒動到那三份)"
```

提交訊息標題：`docs(ai-sdlc): give the bdd-spec/bdd-formulation boundary its reason`

- [ ] **Step 4: `check_spec.py` 的 docstring 自相矛盾**

docstring 開頭第 7 行寫「檢查十三件事」並列了 13 項，結尾第 38 行卻寫「不算在上面**十二**項一致性檢查裡」。

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "十三件事\|十二項\|十三項" skills/bdd-formulation/scripts/check_spec.py
```

把「上面十二項」改成「上面十三項」。

**MUST NOT: 改任何一行程式碼。** 只動 docstring 裡的那兩個字。證明：

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
git diff -U0 -- skills/bdd-formulation/scripts/check_spec.py | grep '^[-+][^-+]'
```

Expected: 恰好兩行，一 `-` 一 `+`，差別只在「十二」與「十三」。

提交訊息標題：`docs(bdd-formulation): check_spec.py's docstring counted itself twice, differently`

- [ ] **Step 5: 全套驗收**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for d in skills/*/; do printf "%-18s" "$(basename $d)"; \
  python3 skills/skill-rules/scripts/audit_skill.py "$d" >/dev/null 2>&1 && echo "exit 0" || echo "exit $?"; done
FX="$(mktemp -d)/fx"; mkdir -p "$FX/specs/2026-09-09-parking" "$FX/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FX/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-spec/scripts/status.py "$FX" | grep "^合計"
python3 skills/bdd-formulation/scripts/check_spec.py "$FX" 2>&1 | tail -1
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
claude plugin validate . 2>&1 | tail -3
```

Expected: 8 個 skill 全部 `exit 0`；`合計` 是 `2     1     1`；`4 個問題`；`validate` 通過，**恰好 1 個警告**（plugin root 的 `CLAUDE.md`）。

---

## Self-Review

**1. 規格覆蓋**

| 規格要求（D5、D6、「計畫二」） | 哪個任務 |
| --- | --- |
| 目標使用者併進 `## Personas` | Task 1 Step 3（格式）＋ Step 7a（fixture） |
| 技術偏好與限制 | Task 1 Step 2（章節順序）＋ Step 4（格式）＋ Step 7b（fixture） |
| Known Boundaries 三張清單 | Task 1 Step 2 ＋ Step 5 ＋ Step 7b；Task 2 Step 3 的兩條 MUST |
| Further Notes | Task 1 Step 2 ＋ Step 6 ＋ Step 7c |
| fixture 是原子單位 | Task 1 的 Files 只有兩個檔，Step 11 一筆提交 |
| D6：新 seam 在最高點提出 | Task 2 Step 2a |
| D6：不得變成需求訪談 | Task 2 Step 2b |
| D6：seam 結果寫進 `## Testing Decisions` | **已存在**，`spec-format.md` 的 `## Testing Decisions` 一節已有 seam 表，不動 |
| D6：跟使用者確認 | **已存在**，`SKILL.md` 的 `Ask user` 段落，不動 |
| AC 層級不動 | Global Constraints；Task 1 沒有任何一步碰 `## User Stories` |

**不在範圍內**（設計 spec 編在計畫三）：`bdd-discovery` 換前沿模型、`technical-probes.md` 加欄、`example-mapping` 改輸入。本計畫不碰。

**2. 佔位字掃描**

沒有 TBD／TODO／「之後再補」。三個新節的完整格式規範與 fixture 內容都逐字寫出來了。Task 3 Step 1 與 Step 3 給的是**要達到的狀態與判準**而非逐字稿——那兩處需要配合上下文的語氣，逐字指定反而會產出接不上的段落；兩者都附了改完要跑的驗證指令。

**3. 名稱一致性**

三個新節名在章節順序表、格式規範標題、fixture、`SKILL.md` 的引用裡拼法一致：`## Tech Preferences & Constraints`、`## Known Boundaries`、`## Further Notes`。Task 1 Step 8 的 `diff` 會機械地抓到任何不一致——它比對的是章節順序表宣告的名稱與 fixture 實際的標題，拼錯一個字就會炸。

**4. 我在這份計畫裡最可能犯的錯**

計畫一裡我犯了兩次**驗證預期與自己 mandate 的施工對不上**（Step 8 預測兩行、實際六行；叫人去改不存在的行號）。這一份的對策：

- 每一個數字預期都是**我剛才實際量過的**（13 節、152 行、633 行、「十三節」2 處），不是推算的。
- Task 1 Step 8 用 `diff` 取代數數字——它同時驗名稱、順序與數量，而且不需要我預先知道答案。
- Task 1 Step 9 用 `diff` 比對腳本輸出，基線在 Step 1 由執行者**自己量**，不是我寫死在計畫裡。
- Task 3 Step 2 明講「刪掉行號，不要重新對齊」，並附一條驗證每個被引用章節真的存在的指令——這正是那一項 follow-up 存在的原因。

**而我在寫完這份計畫後的自我複查裡，又犯了同一個錯一次。** Task 3 Step 2 原本照抄計畫一複審點名的**四個**行號，但實際去讀 `PLAN.md` 的未決事項 #13 才發現它引用了**九個**。複審當時只舉了前四個當例子，我把例子當成了清單。已改成九列對照表，每一列都查過對得到哪個真實標題。同一族的第三次：**別人給的清單是樣本，不是母體。**
