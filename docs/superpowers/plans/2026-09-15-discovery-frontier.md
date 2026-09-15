# `bdd-discovery` 換成前沿模型 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 `bdd-discovery` 的三個 pass 換成**設計樹 ＋ 前沿**：每輪算出所有前提已定案的問題，一次問完，答案回來重算前沿，前沿空了就結束。

**Architecture:** 四個任務按「誰讀誰」排。Task 1 重寫 `bdd-discovery/SKILL.md`（最大的一份）；Task 2 改它的兩份 reference；Task 3 改 `example-mapping`（D4：它從下游步驟變成 discovery 的手法）；Task 4 掃掉散在其他九個檔裡的 pass 說法。前三個任務彼此獨立，Task 4 必須最後跑——它要掃的東西有一部分是前三個任務剛寫下的。

**Tech Stack:** Markdown skill 文件、Python 3 稽核腳本（`audit_skill.py`、`status.py`、`check_spec.py`）、`claude plugin validate`。沒有 CI，全部手跑。

**Spec:** `docs/superpowers/specs/2026-09-14-grilling-and-formulation-design.md` — D1、D2、D3、D4、D9 與「計畫三」那一節。參考來源是 [`mattpocock-skills:grilling`](~/.claude/plugins/cache/mattpocock/mattpocock-skills/1.2.3/skills/productivity/grilling/SKILL.md)，22 行。

## Global Constraints

- **這一份不動任何格式。** `clarify-log.md`、`input.md`、`spec.md` 的欄位與結構完全不變；改的只有**怎麼問**。
- **`status.py` 的解析邏輯一行都不准改。** 它讀 `## Open Questions` 的 `面向` 欄，那一欄照樣填。**唯一允許動的是兩行使用者看得到的標籤與註解裡的 pass 字樣**（見下方 Ruling A）。
- **16 個追問面向不增不減、不改名。** `BIZ_DIMS` 11 個 ＋ `TECH_DIMS` 5 個，名稱與順序都是 `status.py` 與 `technical-probes.md` 的共用契約。
- **`benchmark/cases/*/grading.md` 完全不碰。** 那是評估用的答案卷，用 Pass 1／Pass 2 當評分層級；改它等於改被評的東西。見 Ruling B。
- 接縫不變：`FR-<n>` ↔ `@rule-<n>`、`AC-<n>.<m>` ↔ `@example-<n>.<m>`。
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
| `check_spec.py` 對 fixture（無 `features/`） | `全部通過`，exit 0 |
| `skills/bdd-discovery/SKILL.md` | **435 行** |
| `skills/bdd-discovery/references/technical-probes.md` | 105 行 |
| `skills/bdd-discovery/references/clarify-log-format.md` | 85 行 |
| `skills/example-mapping/SKILL.md` | 約 150 行 |

**fixture** 建法固定：

```bash
FA="$(mktemp -d)/fx"; mkdir -p "$FA/specs/2026-09-09-parking" "$FA/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FA/specs/2026-09-09-parking/spec.md"
```

這份計畫**不碰任何被腳本解析的檔案**，所以三個腳本的輸出從頭到尾都必須紋風不動。那是本計畫最外圈的護欄。

---

## 兩條動手前就定下的裁決

**Ruling A — `status.py` 的兩行輸出標籤要改，但只改字串。**
D3 寫「`status.py` 一行都不用改」，那句話講的是**解析行為**，而它寫下來的時候沒有人去看這支腳本印什麼。它實際印的是：

```
追問覆蓋 · 業務面向（Pass 2）空與零/邊界/…
追問覆蓋 · 技術面向（Pass 3）seam/模組邊界/…
```

這一份計畫之後，Pass 2 與 Pass 3 在任何 skill 裡都不存在了。使用者跑 `status.py` 會看到一個指向已刪結構的標籤——正是這一系列計畫一直在清的那種缺陷。**改兩個字串與幾行註解，零行邏輯**，並在提交訊息裡明講這超出 D3 字面但符合它的用意。

**Ruling B — `benchmark/cases/*/grading.md` 不碰，即使它們用 Pass 當評分層級。**
`parking-billing/grading.md` 有 3 處、`personal-memo/grading.md` 有 2 處。那兩份是**評估用的答案卷**：改掉評分層級的名字，等於改動被評的標準，而且改完之後拿舊結果跟新結果比就失去意義。這件事要在收工時回報給使用者，由他決定什麼時候連同評估設計一起更新。

---

## 新的 `bdd-discovery` 長什麼樣

```
frontmatter                       改寫：不再提三個 pass、四選項
# CLARIFY — 把需求逼問到收斂        改寫開場
## 使用時機                        大致保留
## Skill Boundaries                大致保留
## 前置確認（Ask user）             保留
## 流程                            重寫：開場一次 ＋ 每輪一次
### 開場：只做一次的四件事           ← 原 Pass 1 的步驟 1、2 搬過來
### 每一輪                          ← 新：算前沿 → 一次問完 → 等 → 重算
### 續跑                            改寫：不再是「跳過 Pass 1」
## 前沿怎麼算                       ← 新的核心章節
### 設計樹                          ← D1
### 十六個面向是前沿的檢查表          ← D3
### 能自己查的事實不准問使用者        ← D9
## 一題長什麼樣                     ← grilling 的 ❓／➡️ 格式 ＋ 逃生口（D2）
## 結束：前沿空了                   ← 原「收尾」搬過來 ＋ example-mapping 就緒判定（D4）
### 稽核分類是否窮盡                 原樣保留
## 產物                            原樣保留
## 產物隔離                        原樣保留
## 完成後                          原樣保留（四種結束方式）
```

**必須活下來的東西**（逐項在 Task 1 的驗收裡點名）：

| 什麼 | 現在在哪 | 為什麼不能掉 |
| --- | --- | --- |
| `input.md` 逐字複製的 MUST | Pass 1 步驟 1 開頭 | `bdd-spec` 的 `PRD §x` 回指它；晚一步等於中間那段時間指向不存在的文件 |
| `## 核心關係人` 表 ＋ `決策權` 的理由 | Pass 1 步驟 1 | **只有這裡寫得出 `決策權`**；下游只抄不編 |
| 隱含假設的機械掃描表（名詞／數字／「自動」／指標／排除項） | Pass 1 步驟 1 末 | 假設對寫的人是常識，常識不會被寫下來——靠讀找不到 |
| PRD 的「假設」段每條都當沒問過 | Pass 1 步驟 1 | 它是未驗證的斷言，不是前提 |
| 識別角色的五個問題 ＋ 指向 `persona-definition.md` | Pass 1 步驟 2 | 角色名不統一會讓同一題問兩次 |
| 「重問一題已經答過的，比漏問一題更傷信任」 | 續跑 | |
| 收尾的三項核對（面向覆蓋、persona 雙向、分類窮盡） | Pass 2 步驟 5 | 面向覆蓋是「要嘛問過要嘛標 n/a」那條規則的驗收 |
| 稽根分類是否窮盡 ＋ 健身 App 棒式那個實例 | Pass 2 步驟 5 底下 | 十一個面向問的是已知事物的極端值，問不到「漏掉一整類」 |
| `clarify-log.md` 保留的理由 | Pass 2 步驟 5 底下 | |
| 四種結束方式 | 完成後 | |

**該消失的東西**：三個 pass 的框架本身、流程圖、「一次一題」的節奏、「4 ＋ 1」四選項格式、Pass 1 的「只問會改變範圍」篩選（前沿在結構上就做到了）、Pass 3 作為獨立階段、`續跑` 裡「跳過 Pass 1」的邏輯。

**行數目標**：435 → **300 上下**。D1 說「435 行應該大幅縮短」。低於 260 要回報——可能刪掉了上表某一項。

---

## Task 1: 重寫 `bdd-discovery/SKILL.md`

**Files:**
- Modify: `skills/bdd-discovery/SKILL.md`

**Interfaces:**
- Consumes: 無
- Produces: 前沿模型的用語（設計樹、前沿、一輪），Task 2、3、4 都要跟它一致

- [ ] **Step 1: 先把要保住的東西抄出來**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
mkdir -p /tmp/keep
sed -n '132,167p' skills/bdd-discovery/SKILL.md > /tmp/keep/01-骨架與關係人.md
sed -n '168,177p' skills/bdd-discovery/SKILL.md > /tmp/keep/02-識別角色.md
sed -n '324,364p' skills/bdd-discovery/SKILL.md > /tmp/keep/03-收尾與窮盡.md
sed -n '383,435p' skills/bdd-discovery/SKILL.md > /tmp/keep/04-產物與完成後.md
wc -l /tmp/keep/*
```

這四段是「必須活下來」表的載體。改寫時從這裡貼，不要憑記憶重寫——**憑記憶重寫一段
論證，掉的永遠是那句解釋為什麼的話**，而這份檔案的價值幾乎都在那些句子裡。

- [ ] **Step 2: 讀 22 行的 `grilling`**

```bash
cat ~/.claude/plugins/cache/mattpocock/mattpocock-skills/1.2.3/skills/productivity/grilling/SKILL.md
```

整個模型就這 22 行。要抄的是**輪次 ＋ 前沿**這個結構與 `❓`／`➡️` 的題目格式，
不是它的英文散文。

- [ ] **Step 3: 改寫 frontmatter 與開場**

`description` 現在寫「再逐題四選項問到不確定收斂」——那個作法要消失。改成：

```markdown
---
name: bdd-discovery
description: >
  把 PRD 或一段模糊的需求逼問到收斂——把它攤成一棵設計樹，每輪算出所有前提
  已定案的問題一次問完，答案回來重算前沿，前沿空了就結束。問答記進
  `clarify-log.md`，原始輸入逐字存進 `input.md`。規則由 `bdd-spec` 從答案抽出，
  就緒判定交給 `example-mapping`。
  BDD 六步流程的 CLARIFY 步驟。
  觸發詞：「這個需求要怎麼拆」「需求講不清楚」「幫我釐清需求」
  「這份 PRD 有什麼沒寫清楚」「PRD 的隱含假設」「這份需求哪裡有洞」「PRD 怎麼拆解」
  「有哪些角色」「actor 要怎麼定義」「這份 PRD 要怎麼變成 story」
  「寫規格前先釐清」「需求有什麼漏掉的」。
  English: clarify this requirement, break down a PRD, find the unstated
  assumptions, grill me about this spec.
---

# CLARIFY — 把需求逼問到收斂

一份 PRD 或一段模糊的敘述進來，出去的是一份 `clarify-log.md`：**問了什麼**、
**答了什麼**，以及**還開著什麼**。

**把需求攤成一棵設計樹。** 每個決定都會岔出依賴它的下一批決定。**前沿**是所有
**前提已經定案**的決定——現在就問得出口、不必先猜任何還沒聽到的答案的那些題。
一輪問完整個前沿，等答案，重算，再問下一輪。**前沿空了就結束**：每一條分支都
走過，沒有東西被默默假設。

不寫 Gherkin。不談實作。不設計資料表。
```

- [ ] **Step 4: 重寫 `## 流程`**

整段（原第 62–110 行，含那張流程圖與「續跑：什麼時候跳過 Pass 1」）換成：

```markdown
## 流程

```
新的需求敘述                        或    繼續澄清（bdd-discovery [feature]）
        ↓                                          ↓
開場（只做一次）                              （開場已經做過，直接進輪次）
  原文存檔 · 拆骨架 · 登記關係人 · 識別角色        ↓
        ↓                                          ↓
每一輪 ←──────────────────────────────────────────┘
  1. 算前沿：所有前提已定案的決定
  2. 一次問完，編號，每題附建議答案
  3. 等使用者回答
  4. 答案重塑這棵樹 → 回到 1
        ↓
前沿空了 → 收尾核對 → example-mapping 攤地圖、跑診斷、提就緒問句
        ↓
一份 specs/<date>-<feature>/clarify-log.md，記下結束方式（收斂／阻塞／零進展／喊停）
```

### 續跑

| 呼叫方式 | 範圍 |
| --- | --- |
| `bdd-discovery` | **預設：所有 `clarify-log.md` 裡『待答』表非空的 feature** |
| `bdd-discovery <feature>` | 只跑指定的那一個 |

**不要問對方要跑哪一個。** 預設全部，想限定的人自己會帶參數——多問一次等於把
決定推回去，而多數時候答案就是「都繼續」。

`clarify-log.md` 不存在才跑開場。存在就只讀不重跑，並且**用它重建設計樹**：
『已答』表的每一列是一個已定案的決定，『待答』是還開著的分支，『n/a』是走過
但判定不適用的。前沿從這三張表算得出來。

IMPORTANT: 重問一題已經答過的，比漏問一題更傷信任——它證明你沒有讀既有產物。

換 feature 時明講「現在開始 X」，並在每一輪結束回報還剩幾個。跑多個 feature 時，
**跨 feature 的共用規則會浮現**——那是一次跑完全部的主要收穫，單獨跑一個看不到。
```

**注意**：上面那段含一個三引號區塊（流程圖），寫進檔案時內外都要保留。

- [ ] **Step 5: 寫 `### 開場：只做一次的四件事`**

把 `/tmp/keep/01-骨架與關係人.md` 與 `/tmp/keep/02-識別角色.md` 貼進來，
標題從 `### 1. 拆成骨架` / `### 2. 識別角色` 改成這一節底下的四個小節：
**原文存檔**、**拆骨架與隱含假設**、**登記核心關係人**、**識別角色**。

內文除了兩處，其餘**一字不改**：

1. 原文裡「答案記進 `clarify-log.md` 的 `已答`／`待答`／`n/a` 三張表——這幾類問題
   沒有自己的章節⋯」那句提到「十一個業務面向或五個技術面向」，保留。
2. 刪掉原 `### 3. 只問「會改變範圍」的題目` 整節。**前沿在結構上就做到了它想做的事**
   ——一個前提還沒定案的問題本來就不在這一輪的前沿裡，不需要另一條篩選規則。
   在新的「前沿怎麼算」一節裡用一句話交代這件事，讓讀過舊版的人知道它去哪了。

- [ ] **Step 6: 寫 `## 前沿怎麼算`——這一節是整份改寫的核心**

新寫，三個小節：

```markdown
## 前沿怎麼算

### 設計樹

每個已定案的決定會岔出依賴它的下一批決定。**前沿 ＝ 所有前提已經定案的決定。**

判準只有一句：**這題的答案，會不會因為另一題還沒答而改變？** 會，它屬於**後面**
的輪次，不是這一輪。

這取代了舊版 Pass 1 的「只問會改變範圍的題目」那條篩選——那條規則在猜哪些題
該先問，前沿是**算**出來的。範圍題之所以先問，不是因為有人規定，是因為其他
幾乎所有題的前提都掛在它身上。

MUST NOT: 因為某題「感覺比較重要」就把它拉進這一輪。重要不等於前提已定案；
拉進來的結果是你拿到一個建立在還沒問到的答案上的回答。

### 十六個面向是前沿的檢查表

11 個業務面向（空與零／邊界／重複／時序／權限／失敗／時間／規模／降級／時限／可觀測）
與 5 個技術面向（seam／模組邊界／介面與型別契約／排序契約／決定性／既有資產／測試慣例）
不再是三趟掃描，而是**每輪算前沿時的檢查表**：

> 對每一個**已經定案**的決定問——這 16 個面向裡，有哪個還沒對它問過？

問得出來的就是這一輪前沿的題。技術面向多一條：**能自己查的先查**（見下一節）。

面向名稱與順序是 `../bdd-spec/scripts/status.py` 的 `DIMS` 與
`references/technical-probes.md` 的共用契約，**不得增減或改名**。每一題的
`面向` 欄照樣填——`status.py` 的追問覆蓋率完全靠那一欄，而它是這個 repo 裡
唯一看得出「有沒有真的問過邊界」的東西。

### 能自己查的事實不准問使用者

MUST: 前沿問題需要的**事實**若在環境裡查得到——既有測試怎麼寫、某個模組的介面
長怎樣、有沒有現成的 seam——**派 sub-agent 去查**，不要問使用者。
哪些面向該自己查、哪些該問人 → `references/technical-probes.md` 的第二欄。

MUST NOT: 因為 sub-agent 還在跑就停下來。**一個進行中的探查是一個未定案的前提**
——只有依賴它的問題要等，**前沿其餘的照問**。

**決定是使用者的，事實是你的。** 把事實問題丟回去，等於要對方做你做得到的功課，
而且他的答案未必比環境裡的真相準。
```

**注意**：這一節含一個引用區塊（`>`），不是程式碼區塊。

- [ ] **Step 7: 寫 `## 一題長什麼樣`**

取代舊版的「4 ＋ 1」四選項格式。**逃生口（D2）必須保住**：

```markdown
## 一題長什麼樣

格式固定，一輪裡每題都這樣：

```
❓ **Q1** - **<問題標題>**：<問題本體，可以多段，可以含選項>

➡️ <你的建議答案>
```

編號讓對方可以只回「Q2 選第二個、Q3 你的建議可以」，不必重述題目。

**建議欄是逃生口，不是裝飾。** 一個沒有建議的問題把全部認知負擔丟回去；一個
有建議的問題只要求對方判斷「對不對」。

MUST: **沒有站得住的預設時，建議欄要寫「這題我不建議猜——要去問 <誰>」**，
`<誰>` 指向 `## 核心關係人` 表裡的人。編一個聽起來合理的建議，比承認不知道更糟
——對方多半會直接接受，而那個答案從此看起來像是問過的。

一輪問完整個前沿。**不要一次只問一題**——前沿上的題彼此沒有依賴，分開問只是
多幾趟來回，而每一趟來回都是對方重新載入上下文的成本。
```

**注意**：這一節含一個三引號區塊（題目格式），寫進檔案時內外都要保留。

- [ ] **Step 8: 寫 `## 結束：前沿空了`**

把 `/tmp/keep/03-收尾與窮盡.md` 貼進來，標題從 `### 5. 收尾：標記完成⋯` 改成
這一節，開頭那句「業務面向的『待答』列都清空之後（技術面留給 Pass 3 掃）」改成
「**前沿空了之後**」。三項核對與「稽核分類是否窮盡」整段**一字不改**。

核對清單末尾加第四項（D4）：

```markdown
4. **交給 `example-mapping`**：前沿空了之後，由它攤開四色卡、跑四條診斷、把就緒
   問句交給使用者。就緒判定不在本 skill——`clarify-log.md` 沒有就緒欄，本 skill
   不碰它。
```

- [ ] **Step 9: 貼回 `## 產物` / `## 產物隔離` / `## 完成後`**

`/tmp/keep/04-產物與完成後.md` **一字不改**貼回。這三節裡沒有任何 pass 語言，
已經確認過。

- [ ] **Step 10: 逐項驗「必須活下來」**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
for s in "逐字複製" "核心關係人" "決策權" "隱含假設" "每個名詞" "每個數字" \
         "怎麼取得這個身分" "persona-definition" "比漏問一題更傷信任" \
         "稽核分類是否窮盡" "棒式" "clarify-log.md 保留" "收斂" "阻塞" "零進展" "喊停"; do
  printf "%-28s %s\n" "$s" "$(grep -c "$s" skills/bdd-discovery/SKILL.md)"
done
```

Expected: **每一項都 ≥ 1**。任何一個是 0，代表改寫掉了一樣「必須活下來」表裡的東西。

- [ ] **Step 11: 負向檢查——pass 語言零殘留**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "Pass 1\|Pass 2\|Pass 3\|四選項\|4 ＋ 1\|一次一題" skills/bdd-discovery/SKILL.md || echo "零殘留 ✓"
grep -n "前沿\|設計樹" skills/bdd-discovery/SKILL.md | wc -l
wc -l skills/bdd-discovery/SKILL.md
```

Expected: 第一條零輸出；「前沿」與「設計樹」合計出現 **≥ 10** 次；行數落在 **260–340**。
低於 260 就回頭對 Step 10 的清單。

- [ ] **Step 12: 稽核與提交**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-discovery; echo "exit=$?"
```

Expected: `✔ 通過`，`exit=0`，**沒有 S10**。

---

## Task 2: `bdd-discovery` 的兩份 reference

**Files:**
- Modify: `skills/bdd-discovery/references/technical-probes.md`
- Modify: `skills/bdd-discovery/references/clarify-log-format.md`

**Interfaces:**
- Consumes: Task 1 定下的用語（前沿、一輪、設計樹）
- Produces: `technical-probes.md` 的第二欄，Task 1 的「能自己查的事實」一節指向它

- [ ] **Step 1: `technical-probes.md` 加第二欄（D9）**

檔案現在是 105 行、五個面向各一節（`## seam`、`## 模組邊界`、`## 介面與型別契約`、
`## 排序契約／決定性`、`## 既有資產／測試慣例`）。

**五個 `## ` 標題與它們的名稱不准改**——它們是 `status.py` 的 `TECH_DIMS` 的來源，
兩邊不一致以這份為準，改名會讓追問覆蓋率安靜地少一格。

改兩件事：

1. **檔頭**（第 1–19 行）現在寫「# Pass 3 的技術追問面向」與「Pass 3 逐題掃過這份
   清單⋯」。改成講它在前沿模型裡的角色：**每輪算前沿時，對已定案的決定逐項對照這
   五個面向**；並說明第二欄是什麼。
2. **每一節加一段 `**自己查還是問人**`**，寫清楚這個面向的事實哪些查得到、哪些
   只有人知道。例如 seam：既有的 seam 有哪些、測試怎麼掛——**查得到**；這批行為
   該打在哪一層——**問人**（那是決定，不是事實）。

判準寫在檔頭，一句話：**事實查得到，決定問得到。查得到的去查，查不到的才問。**

- [ ] **Step 2: `clarify-log-format.md` 的兩處 pass 語言**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "Pass" skills/bdd-discovery/references/clarify-log-format.md
```

第 70–71 行講 `面向` 欄的值域，寫著「最後三個是 Pass 3 答案的家⋯Pass 3 每一題都
只能填 `—`」。改寫成前沿模型的說法，**不要改那一欄的值域本身**——值域是
`status.py` 的契約。

- [ ] **Step 3: 驗收**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
echo "--- 五個技術面向的標題還在、名稱沒變 ---"
grep -n "^## " skills/bdd-discovery/references/technical-probes.md
python3 - <<'PY'
import re, pathlib
probes = pathlib.Path("skills/bdd-discovery/references/technical-probes.md").read_text(encoding="utf-8")
heads = [h.split("——")[0].strip() for h in re.findall(r"^## (.+)$", probes, re.M)]
status = pathlib.Path("skills/bdd-spec/scripts/status.py").read_text(encoding="utf-8")
dims = re.search(r"TECH_DIMS = \[(.*?)\]", status, re.S).group(1)
dims = re.findall(r'"([^"]+)"', dims)
print("probes:", heads)
print("DIMS  :", dims)
print("一致" if [h for h in heads if h in dims] == dims else "*** 不一致 ***")
PY
echo "--- 第二欄有沒有寫 ---"; grep -c "自己查還是問人" skills/bdd-discovery/references/technical-probes.md
echo "--- pass 零殘留 ---"; grep -n "Pass [0-9]" skills/bdd-discovery/references/*.md || echo "零殘留 ✓"
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-discovery; echo "exit=$?"
```

Expected: `TECH_DIMS` 的五個名稱全部出現在 probes 的標題裡且順序一致；
`自己查還是問人` 出現 **5** 次；pass 零殘留；稽核 exit 0。

**那段 Python 是這個任務的重點檢查**：`status.py` 的註解明寫「名稱與順序取自
`references/technical-probes.md` 的章節標題，不是憑空編的縮寫——兩邊不一致以那份
為準」。那是一個**跨檔的契約，沒有任何腳本在驗它**。改這份檔案時最容易踩的就是
順手改一個標題的措辭。

---

## Task 3: `example-mapping` 的輸入改成進行中的訪談（D4）

**Files:**
- Modify: `skills/example-mapping/SKILL.md`

**Interfaces:**
- Consumes: Task 1 的用語；Task 1 的 `## 結束：前沿空了` 第 4 項會指向這裡
- Produces: 無

- [ ] **Step 1: 改 `## 使用時機`**

第一條現在是「`bdd-discovery` 的 Pass 2 要把需求攤成規則與例子」。改成前沿模型的
說法：**逼問進行中，要把已經問出來的東西攤在桌上看形狀**；以及**前沿空了，要跑
就緒判定**。

第三條「手上有一份 `spec.md`，想看它的形狀而不是讀它的字」——**保留**。D4 說
「四色卡是追問時攤在桌上的東西，不是事後的檢查」，講的是它在流程裡的**主要**位置，
不是禁止對一份既有的 `spec.md` 跑。`## 兩種模式` 那一節本來就講這件事。

- [ ] **Step 2: 檢查整份檔案還有沒有把自己當下游步驟的說法**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "Pass\|spec.md\|之後\|下一步" skills/example-mapping/SKILL.md
```

逐行判斷：講「對一份既有 `spec.md` 跑」的**留著**（那是它的第二種模式）；
講「本 skill 在 `spec.md` 之後才跑」的要改。

- [ ] **Step 3: 驗收**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "Pass [0-9]" skills/example-mapping/SKILL.md || echo "pass 零殘留 ✓"
grep -n "就緒" skills/example-mapping/SKILL.md | head -3
python3 skills/skill-rules/scripts/audit_skill.py skills/example-mapping; echo "exit=$?"
```

Expected: pass 零殘留；就緒判定還在；稽核 exit 0。

---

## Task 4: 掃掉其他九個檔的 pass 說法

**必須最後跑**——它要掃的東西有一部分是前三個任務剛寫下的。

**Files:**
- Modify: `skills/bdd-spec/scripts/status.py`（**只改兩個字串與註解，零行邏輯**）
- Modify: `skills/bdd-spec/references/persona-definition.md`
- Modify: `skills/bdd-spec/references/spec-format.md`
- Modify: `skills/clarify-loop/SKILL.md`
- Modify: `PLAN.md`

**Interfaces:**
- Consumes: 前三個任務定下的全部用語
- Produces: 無

- [ ] **Step 1: 列出工作清單**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -rn "Pass 1\|Pass 2\|Pass 3" CLAUDE.md PLAN.md README.md docs/ skills/ 2>/dev/null \
  | grep -v "^docs/superpowers/"
```

這份輸出就是工作清單。**`benchmark/` 刻意不在掃描範圍內**——見下面 Step 5。

- [ ] **Step 2: `status.py`——兩個字串與幾行註解，零行邏輯**

使用者看得到的兩行（約第 184、187 行）：

```python
    print("\n追問覆蓋 · 業務面向（Pass 2）"
    print("追問覆蓋 · 技術面向（Pass 3）"
```

`（Pass 2）`／`（Pass 3）` 改成不指向已刪結構的說法——業務／技術這個分組本身還在，
只是它不再對應某一趟掃描。同樣處理第 9、28、31、149 行註解裡的 pass 字樣。

**證明只動了字串與註解：**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
git diff -U0 -- skills/bdd-spec/scripts/status.py | grep '^[-+][^-+]'
```

逐行看過：每一行都必須是註解（`#` 開頭）或 `print(` 裡的字面字串。有任何一行不是，
停下來回報。

**然後證明行為沒變：**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FA="$(mktemp -d)/fx"; mkdir -p "$FA/specs/2026-09-09-parking"
cp skills/bdd-spec/examples/minimal-spec.md "$FA/specs/2026-09-09-parking/spec.md"
git show HEAD:skills/bdd-spec/scripts/status.py > /tmp/old_status.py
diff <(python3 /tmp/old_status.py "$FA" | grep -v "追問覆蓋") \
     <(python3 skills/bdd-spec/scripts/status.py "$FA" | grep -v "追問覆蓋") \
  && echo "除了那兩行標籤，輸出完全相同 ✓"
python3 skills/bdd-spec/scripts/status.py "$FA" >/dev/null 2>&1; echo "exit=$? (基線 0)"
```

- [ ] **Step 3: 四份 Markdown**

| 檔案 | 哪幾行 | 改什麼 |
| --- | --- | --- |
| `skills/bdd-spec/references/persona-definition.md` | 3、74、75 | 「`bdd-discovery` 在 Pass 1 步驟 2」→ 開場那一步；74–75 講「Pass 1 驗不完、等 Pass 2 抽完規則再回來判」→ 改成前沿模型的說法（那條核對留到前沿空了的收尾） |
| `skills/bdd-spec/references/spec-format.md` | 243、452、626、642 | 243「範圍題與技術題在 Pass 1、Pass 3 產生」、452「Pass 3 逐項掃過技術面向時」、626「Pass 2 的〈追問清單〉」、642「`—`：Pass 1 的範圍題」 |
| `skills/clarify-loop/SKILL.md` | 103 | `### 3. 一次一題，四個選項` —— 這一節整個是舊模型。**改之前先讀完整個 skill**：它是「把一疊跳過的問題再收斂」的獨立循環，不是 `bdd-discovery` 的一部分，所以「一次一題」在它身上**可能仍然是對的**。判準：這個 skill 自己有沒有前沿可算？沒有（它處理的是一疊已經在檯面上的題），那就只拿掉「四個選項」那半、留下節奏。把你的判斷寫進報告 |
| `PLAN.md` | 200、201、347、431、443、473、475 | 200／201 是 skill 對照表的現況欄——改成前沿模型；347 是已完成的歷史條目，**不動**；431、443、473、475 在未決事項裡，逐條判斷是仍然成立、還是被這份計畫解決了 |

`PLAN.md` 的每一處都要先讀上下文再決定。**歷史條目保持歷史**，前瞻性的敘述才改。

- [ ] **Step 4: 全域負向掃描**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -rn "Pass 1\|Pass 2\|Pass 3" CLAUDE.md PLAN.md README.md docs/ skills/ 2>/dev/null \
  | grep -v "^docs/superpowers/"
```

Expected: 只剩 `PLAN.md` 那條已完成的歷史條目（第 347 行附近）。其餘零命中。
**多出任何一行都要說明它為什麼還在**，不要為了讓輸出變空而改歷史。

- [ ] **Step 5: `benchmark/` 的兩份 grading.md——不碰，但要回報**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -c "Pass" benchmark/cases/parking-billing/grading.md benchmark/cases/personal-memo/grading.md
```

那兩份用 Pass 1／Pass 2 當評分層級。**一個字都不要改**：它們是評估用的答案卷，
改掉層級名稱等於改動被評的標準，而且改完之後舊結果與新結果不可比。
把這兩個檔名與命中數寫進報告，收工時要回報給使用者。

- [ ] **Step 6: 全套驗收**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for d in skills/*/; do printf "%-18s" "$(basename $d)"; \
  python3 skills/skill-rules/scripts/audit_skill.py "$d" >/dev/null 2>&1 && echo "exit 0" || echo "exit $?"; done
FA="$(mktemp -d)/fx"; mkdir -p "$FA/specs/2026-09-09-parking" "$FA/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FA/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-spec/scripts/status.py "$FA" | grep "^合計"
python3 skills/bdd-formulation/scripts/check_spec.py "$FA" 2>&1 | tail -1
FB="$(mktemp -d)/fx"; mkdir -p "$FB/specs/2026-09-09-parking"
cp skills/bdd-spec/examples/minimal-spec.md "$FB/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-formulation/scripts/check_spec.py "$FB" >/dev/null 2>&1; echo "無 features/ 閘門: exit=$? (必須 0)"
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
claude plugin validate . 2>&1 | tail -3
```

Expected: 8 個 skill 全 `exit 0`；`合計  2  1  1`；`4 個問題`；閘門 `exit=0`；
`validate` 恰好 1 個警告。

---

## Self-Review

**1. 規格覆蓋**

| 規格要求 | 哪個任務 |
| --- | --- |
| D1 設計樹 ＋ 前沿，三個 pass 消失 | Task 1 Step 4、6 |
| D1 `❓`／`➡️` 題目格式 | Task 1 Step 7 |
| D2 逃生口保住 | Task 1 Step 7 的 MUST |
| D3 16 面向留著，改成前沿的檢查表 | Task 1 Step 6；Task 2 Step 3 的跨檔契約驗證 |
| D3 `status.py` 解析不改 | Global Constraints；Task 4 Step 2 的兩條證明 |
| D4 `example-mapping` 是 discovery 的手法 | Task 1 Step 8 第 4 項；Task 3 |
| D9 能自己查的派 sub-agent | Task 1 Step 6 第三小節 |
| D9 `technical-probes.md` 改成兩欄 | Task 2 Step 1 |
| 計畫三「不動任何格式」 | Global Constraints；Task 4 Step 6 的三腳本驗收 |

**2. 佔位字掃描**

沒有 TBD／TODO。Task 1 的改寫給了完整的替換文字與「從 `/tmp/keep/` 貼，不要憑記憶重寫」
的指示。Task 2、3、4 有幾處給的是**要達到的狀態與判準**而非逐字稿（`technical-probes.md`
的第二欄、`PLAN.md` 的逐條判斷）——那幾處需要讀上下文，逐字指定會產出接不上的段落，
而且每一處都附了驗證指令與「不確定就回報」的出口。

**3. 任務相依**

Task 4 必須最後。Task 1、2、3 彼此獨立，但 Task 2 的 `technical-probes.md` 第二欄
會被 Task 1 的「能自己查的事實」一節指向，所以 Task 1 寫那一節時要用 Task 2 會建立的
小節名 `**自己查還是問人**`——已在兩邊寫明。

**4. 一個我在自我複查裡抓到的錯**

我原本在 Task 4 的表裡寫著 `skills/clarify-loop/` 「是使用者自己的全域 skill，不是本
plugin 的」，並要執行者先確認再動。**前提是錯的**——`git log -- skills/clarify-loop`
顯示它就在本 repo 的版控裡（`61c75fc`、`4a015d6`）。真正的情況是**同名的兩個 skill**：
本 plugin 這一份，以及使用者 `~/.claude/skills/clarify-loop` 那一份。本 repo 這份是
我們的，可以改。已改成真正該問的問題——那一節的「一次一題」在它自己的脈絡裡是不是
還成立。

**5. 我在這份計畫裡最可能犯的錯**

前三份計畫我犯過四次同一種：**驗收預期與自己 mandate 的施工對不上**。這一份的對策：

- Task 1 Step 10 的「必須活下來」檢查用 `grep -c` 對**16 個字串**逐一計數，清單來自
  我實際讀過那 435 行之後列的表，不是憑印象。
- Task 1 Step 11 的行數區間給的是**範圍**（260–340）而不是一個數字，因為改寫的行數
  本來就無法精確預測；並寫明「低於 260 要回頭對清單」。
- Task 2 Step 3 的跨檔契約驗證是一段**實際會跑的 Python**，不是「請確認一致」。
- Task 4 Step 4 明講「多出任何一行都要說明它為什麼還在，不要為了讓輸出變空而改歷史」
  ——前一份計畫我就寫過一個逼執行者重排整棵樹來湊綠燈的檢查。

**6. 我沒有解決的問題**

`benchmark/cases/*/grading.md` 用 Pass 1／Pass 2 當評分層級。這份計畫刻意不碰，
但那兩份答案卷從此描述一個不存在的流程結構。**下一次要跑那兩個 benchmark case 之前，
評分層級要連同評估設計一起重想**——那不是一次文字替換，而是「這個評分層級現在對應
到前沿模型的什麼」這個問題。收工時要明確回報給使用者。
