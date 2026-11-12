# 砍掉 CAPTURE：PRD 是輸入，不是產物

**狀態**：待批　**日期**：2026-11-13

## 意圖

把流程的第一個產物從「我們自己維護的 `prd.md`」改成「需求方帶來的東西原封不動，
流程產出的第一份文件是 `spec.md`」。七步變六步：

```
DISCOVERY → SPEC → PLAN → IMPLEMENT → VERIFY → REVIEW
```

### 為什麼

需求方手上的 PRD（或一段想法）是**進來**的東西。現在的設計把它抄成一份自己格式的
`prd.md`，再由 DISCOVERY 填，等於替別人的文件代筆——所以才需要 `bdd-capture` 這一步
先「建」它。需求方的看法是：拿著 PRD 或想法，透過 DISCOVERY 收斂；收斂不了再 clarify
一輪，直到收斂為止。收斂的結果才是規格。這跟 10-02 參考的 playbook
（intent → spec → plan）更貼：輸入就是 intent，不必再翻成一份自己的格式。

CAPTURE 當初獨立成 skill 的理由（`docs/ai-sdlc.md`「為什麼是五個 skill」：09-09～09-11
三版都在「記」的時候順手補規則）已經不靠 skill 邊界守了——10-02 起由 `check_spec.py`
第 14 項機械地抓沒來源的 FR／AC。所以把「抄原文、粗切 story」併進 DISCOVERY 開場，
那條防線不受影響。

## 目標

- 流程不再產出或要求 `prd.md`；`bdd-capture` 從 `skills/` 消失
- DISCOVERY 的產物是 `specs/<date>-<feature>/spec.md`（業務規格）
- 技術設計另一份檔：`specs/<date>-<feature>/design.md`（= 今天的 `spec.md`，只改名）
- 「每條 FR／AC 指得出誰說的」原則不變：`← PRD §n` 仍指 `spec.md` 的 `## 原文`，
  `← Q-<n>` 仍指 `## Open Questions`；第 14 項照抓
- `claude plugin validate .` 過；`python3 -m unittest discover -s tests` 全過

## 非目標

- 不改 `spec.md`（原 `prd.md`）的**章節結構**。`prd-format.md` 的骨架整份搬過去，只
  改檔名與「誰第一次寫」那一欄（CAPTURE 的工作變 DISCOVERY 開場）
- 不改 `design.md`（原 `spec.md`）的章節結構
- 不改 `.feature`、`plan.md`、`docs/CONTEXT.md`、`specs/domain-model.md` 的格式與位置
- 不改 DISCOVERY 的前沿演算法、十六個面向、代演模式
- 不改 `check_spec.py` 的十六項檢查內容，只改它找的檔名與訊息字串
- 不重寫 `docs/superpowers/` 底下的歷史設計與計畫——那些是當時的紀錄
- 不動 `docs/bdd.md`、`docs/tdd.md`、`docs/sdd.md`（轉述外部來源，與此無關）

## 反目標（做出來也算失敗）

- DISCOVERY 開場變成一個新的訪談階段。CAPTURE 的「最多四題」限制跟著搬過去：
  開場抄原文、抽開頭幾節、粗切 story，空著的節最多問四題，然後才進輪次
- 第 14 項的來源標記出現第三種合法寫法。合法的仍只有 `PRD §n` 與 `Q-<n>`
- 歷史文件被「順手」改掉名字，之後沒人看得出 10-02 到 11-13 之間流程長什麼樣

## 限制

- 文件繁中、commit 英文，`WHAT:`／`WHY:`／`HOW:` 正文（`CLAUDE.md`）
- 搬檔用 `git mv`，`git log --follow` 要接得上
- `skills/skill-rules/` 是 house style；改過的每個 skill 要過 `audit_skill.py`
- `runs/go-parking-billing/` 是 gitignored 的實跑紀錄，不動

## 考慮過的做法

| 做法 | 取捨 | 結果 |
| --- | --- | --- |
| **A. 只砍 skill，檔名不動**：capture 併進 discovery，`prd.md`／`spec.md` 照舊 | diff 最小；但 `prd.md` 還在，跟「PRD 是輸入不是產物」直接衝突 | 捨棄 |
| **B. 一份 `spec.md`**：業務與技術合在同一檔，DISCOVERY 寫前半、`bdd-spec` 補後半 | PLAN 只讀一份；但業務與技術兩種判斷混在一份文件裡，10-02 前就是這個形狀，當時是刻意拆開的 | 需求方選兩份 |
| **C. 兩份：`spec.md` ＋ `design.md`**（本文件） | 多一個改名；每份文件只有一種擁有者（DISCOVERY／SPEC） | **採用** |
| 輸入另存一檔（像以前的 `input.md`） | 多一個檔；`## 原文` 逐字貼進 `spec.md` 已做到同一件事，第 14 項也是對它查 | 捨棄 |

## 決定

### D1. 檔案

| 今天 | 之後 | 擁有者 |
| --- | --- | --- |
| `specs/<date>-<feature>/prd.md` | `specs/<date>-<feature>/spec.md` | DISCOVERY（`bdd-discovery`、`example-mapping`、`clarify`、`story-splitting`） |
| `specs/<date>-<feature>/spec.md` | `specs/<date>-<feature>/design.md` | SPEC（`bdd-spec`） |
| 需求方帶來的 PRD 檔 | 不搬、不改。內容逐字進 `spec.md` 的 `## 原文` | 需求方 |

`design.md` 這個名字是本文件作者定的：跟 playbook 的 intent／spec／plan 不撞名，
也說得出它裝的是設計而非規格。

### D2. `bdd-capture` 併入 `bdd-discovery`

`skills/bdd-capture/` 整個目錄消失。它的內容去向：

| 原位置 | 去向 |
| --- | --- |
| `SKILL.md` 流程 1–4（抄原文、抽開頭幾節、最多四題、粗切 story）與收尾 | `bdd-discovery/SKILL.md` 開場新增第 0 件事「建檔」，**只在 `spec.md` 不存在時做**；三件既有開場（掃原文、補關係人、識別角色）接在後面 |
| `references/prd-format.md` | `git mv` → `bdd-discovery/references/spec-format.md`；內文 `prd.md` → `spec.md`，「誰第一次寫」欄的 CAPTURE → DISCOVERY 開場 |
| `references/persona-definition.md` | `git mv` → `bdd-discovery/references/` |
| `examples/minimal-prd.md` | `git mv` → `bdd-discovery/examples/minimal-spec.md` |
| description 的觸發詞（「我有一個想法」「記一下這個需求」「把這份 PRD 收進來」…） | 併進 `bdd-discovery` 的 description |

`bdd-discovery` 的前置確認改成：**唯一必要的輸入是需求本身**——一段描述、一份文件
的內容或路徑、或已經存在的 `specs/<date>-<feature>/spec.md`。沒有 `spec.md` 就建
（第 0 件事），有就續跑。「NEVER: `prd.md` 不存在時憑一段描述開始問」這條改為
「NEVER: 還沒把原文逐字寫進 `## 原文` 就開始問」——擋的事一樣，只是建檔的人換了。

CAPTURE 的收尾「把整份給需求方看一遍，改到他說『這是我的意思』」**保留**：第 0 件事
是自己的一個回合，需求方點頭之後才跑開場掃描與第一輪前沿。不保留的話四題會混進
第一輪，而開頭幾節沒經他確認就被當成前提。

`clarify` 單獨模式留下的三節檔（`## Document Overview`、`## Open Questions`、`## 原文`）
的處理方式照搬：第 0 件事從 `## 原文` 補開頭幾節，不碰另外兩節。

### D3. `bdd-spec` 改名產物

`bdd-spec` 讀 `spec.md`、寫 `design.md`。`references/spec-format.md` →
`references/design-format.md`，`examples/minimal-spec.md` → `examples/minimal-design.md`。
description、觸發詞「寫 spec.md」→「寫 design.md」。

`## 回 DISCOVERY 補問` 留在 `design.md`，語意不變。

### D4. 下游 skill 只改路徑

`bdd-formulation`、`bdd-plan`、`bdd-implement`、`example-mapping`、`clarify`、
`story-splitting` 的 SKILL.md、references、examples：`prd.md` → `spec.md`，
（原）`spec.md` → `design.md`，`bdd-capture` 的指引（「沒有 `prd.md` → 先跑
`bdd-capture`」）→「沒有 `spec.md` → 跑 `bdd-discovery`，它會建」。

改名時**一定先改 `spec.md` → `design.md`，再改 `prd.md` → `spec.md`**——反過來會把
兩個字串撞在一起。

### D5. 兩支 script 與測試

- `status.py`：glob `*/prd.md` → `*/spec.md`；「先跑 bdd-capture」→「先跑 bdd-discovery」
- `check_spec.py`：同上；docstring 與訊息裡的 `prd.md` → `spec.md`。第 494 行附近那段
  註解提到的「`spec.md` 點名 `prd.md` 裡沒有的 FR」是歷史說明，改成 `design.md`／`spec.md`
- `tests/test_check_spec.py`、`tests/test_status.py` 的 fixture 檔名跟著改
- 來源標記 `← PRD §n` **不改**：它指的是「原文第 n 段」，原文仍叫原文
- **舊佈局防呆**：兩支 script 在 `specs/<x>/` 底下看到 `prd.md` 與 `spec.md` 並存時，
  退出 1 並說「舊版佈局：先 `git mv spec.md design.md`，再 `git mv prd.md spec.md`」（只有 `prd.md` 時只說後一句），
  附一個測試。這不是第十七項檢查，是跟「沒有 specs/」同級的輸入防線——否則新 script
  會把舊的技術 `spec.md` 當業務規格讀，報「找不到 `## Open Questions`」，訊息指錯修法，
  正是 `ai-sdlc.md` 說的「靜默地錯」。磁碟上 `runs/go-parking-billing/` 就是這種佈局

### D6. 文件

- `docs/ai-sdlc.md`：七步圖改六步；CAPTURE 那一列刪除，它的「存在的理由」併進
  DISCOVERY 列（「原文逐字、粗切 story」變成 DISCOVERY 開場做的事）；「三種模式，四條
  紀律」改「三種模式，三條紀律」，刪 CAPTURE 列；「為什麼是五個 skill」改「四個」，
  CAPTURE 那一段改寫成「為什麼 CAPTURE 曾經獨立、現在為什麼併回去」（防線從 skill
  邊界換成第 14 項）；可追溯鏈的第一節 `需求方的話（prd.md ## 原文）` → `spec.md ## 原文`
- `PLAN.md`：步驟表、skill 表、進度列；未決事項裡「`spec.md` 649 行」那類**仍然開著**
  的條目把檔名改成 `design.md`（它不是歷史，改名後會被讀成業務規格）；「重跑 parking
  某張票」的條目加一句：那個 run 要先 `git mv` 成新佈局。可追溯鏈「定版點被改過三次」
  與 2026-10-02 的進度 bullet 是歷史，**不動**
- `benchmark/cases/*/grading.md` 裡描述 10-02 改動的段落是歷史，不動；
  `benchmark/cases/personal-memo/input-prd.md` 檔名不改——它正是「需求方帶來的 PRD」
- `CLAUDE.md`：pipeline 一行、三模式一節、script 說明、產物路徑
- `benchmark/README.md`、`benchmark/skeleton/go/README.md`、
  `benchmark/skeleton/go/docs/ARCHITECTURE.md`：同樣的改名。注意 skeleton README
  228–258 行的 `spec.md` 是 10-02 前的用法（FR／AC 在裡面），改名後反而變對，
  不要機械套用「先 spec→design」的順序
- `docs/superpowers/specs/*` 與 `plans/*`：**不動**

### D7. 狀態流轉不變

`澄清中` → `待對焦` → `已對焦` 三個狀態、誰能翻哪一個（`example-mapping` 獨占
`已對焦`）全部照舊，只是寫在 `spec.md` 的 `## Document Overview`。

## 交給實作者的

- 各 SKILL.md 改完後的措辭潤飾，以 `audit_skill.py` 過為準
- `bdd-discovery/SKILL.md` 已 472 行，併入開場後會更長；要不要把「第 0 件事」拆成
  `references/opening.md` 由實作者判斷，標準是 `skill-rules` 的長度指引
- 改名後 `check_spec.py` 這個檔名要不要跟著變——不改也說得通（它查的仍是 spec 與
  feature 的一致性），實作者決定，傾向不改

## 驗證

1. `grep -rn "prd.md\|bdd-capture" --include=*.md --include=*.py . | grep -v docs/superpowers | grep -v "^./runs"` 剩下的每一筆都是**歷史敘述**（`ai-sdlc.md` 的「CAPTURE 曾經獨立」、`PLAN.md` 的定版點與 10-02 進度、`grading.md` 的 10-02 段落、`input-prd.md` 檔名、舊佈局防呆的訊息）——實作者在 PR 裡逐筆列出，沒有一筆是指引讀者去找 `prd.md` 或跑 `bdd-capture` 的
2. `claude plugin validate .` 過
3. `python3 -m unittest discover -s tests` 全過
4. 每個改過的 skill：`python3 skills/skill-rules/scripts/audit_skill.py skills/<name>` 過
5. `git log --follow skills/bdd-discovery/references/spec-format.md` 接得回 `prd-format.md`
