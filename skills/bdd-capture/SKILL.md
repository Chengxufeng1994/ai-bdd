---
name: bdd-capture
description: >
  把一個想法或一份現成需求文件寫成 `prd.md` 的半成品——原文逐字存檔、抽出問題／
  預期結果／角色／範圍／成功指標、粗切 story、需求方自己列的開放問題開成紅卡。
  當書記不當分析師：最多問四題，不替 story 填規則與例子。
  BDD 七步流程的 CAPTURE 步驟，產物交給 `bdd-discovery` 問到完整。
  觸發詞：「我有一個想法」「記一下這個需求」「幫我把需求寫下來」「開一份 PRD」
  「建立 prd.md」「capture」「這是我想做的」「需求進來了」「先把它寫成文件」
  「把這份 PRD 收進來」。
  English: capture this idea, write down this requirement, start a PRD,
  create prd.md, intake this requirement document.
---

# CAPTURE — 把需求寫成半成品的 PRD

需求方的一段話、或一份現成文件進來，出去的是一份狀態 `澄清中` 的
`specs/<date>-<feature>/prd.md`：**原文逐字在裡面、開頭幾節從原文抽出來、story
粗切好但底下是空的、需求方自己知道的開放問題已經是紅卡。**

這一步是**書記**。需求方說、你記；原文沒講的留空並說它沒講。追問是下一步
`bdd-discovery` 的事，規則與例子是那一步問出來才寫的。

## 使用時機

- 有人有一個想法，想開始走流程，但還沒有任何文件
- 手上有一份別人寫的 PRD、ticket、或會議紀錄，要收進流程
- `bdd-discovery` 說找不到 `prd.md`

## Skill Boundaries

- 要把需求**問**清楚、找隱含假設、填規則與例子 → 改用 `bdd-discovery`，它讀本 skill 的產物
- 要把一則太大的 story 切開 → 粗切是本 skill 的事，但沿規則細切要等規則問出來，那時用 `story-splitting`
- 要攤四色卡、判定就緒 → 改用 `example-mapping`
- 本 skill 不問十六個面向、不寫 FR／AC、不寫 Gherkin、不碰 `spec.md`

## 參考檔案

- `references/prd-format.md` — `prd.md` 的骨架，每一節誰第一次寫
- `references/persona-definition.md` — 角色五欄怎麼填、哪些可以留「（無資料）」
- `examples/minimal-prd.md` — 一份已對焦的 `prd.md` 長什麼樣（本 skill 交出去的是它的
  澄清中版本：story 底下三個分組空著）

---

## 前置確認（Ask user）

**唯一必要的輸入是需求本身**：一段描述，或一份文件的內容／路徑。

NEVER: 需求還沒拿到就自己想一個來跑。那會產出一份看似完整、實則憑空的
`prd.md`，而錯誤全藏在專業的格式底下。

Ask user: 這個 feature 的 slug（kebab-case，用來命名 `specs/<date>-<slug>/`）。
需求方沒有偏好就從原文的主詞取一個，講出來讓他改。

`specs/<date>-<slug>/prd.md` 已經存在時 **不重建**，告訴對方它在、狀態是什麼，
該跑 `bdd-discovery` 繼續。唯一例外是 `clarify` 單獨模式留下的三節檔（只有
`## Document Overview`、`## Open Questions`、`## 原文`）：那是一份還沒有開頭幾節的
`prd.md`，本 skill 從它的 `## 原文` 補第 2 到第 4 件事，**不碰 `## Open Questions`
與 `## 原文`**——問答已經在那裡了，搬家只會弄丟 `Q-<n>`。

---

## 流程：四件事，順序固定

### 1. 原文逐字進 `## 原文`

MUST: 這是第一個動作，不是收尾。需求方帶著文件來，整份貼進去；口述的，照講的記。
每一段加 `§<n>` 編號（原文自己有編號就沿用），這是之後每個 `← PRD §n` 指的東西。

MUST NOT: 潤飾、補主詞、把「大概半小時」寫成「30 分鐘」。那個「大概」正是下一步
要問的題；改掉它，題目就消失了。

### 2. 從原文抽開頭幾節

照 `references/prd-format.md` 的順序填 `## Background`、`## Goal`、`## Scope — In / Out`、
`## Personas`、`## Assumptions / Constraints`、`## Success Metrics`。每一句行尾標
`← PRD §n`——這幾節在這一步只能有兩種來源：原文（`PRD §n`），或第 3 件事那四題的
答案（`Q-<n>`）。出現 `推論` 就是你補了原文沒有的話。

原文抽不出來的節，保留標題、寫「原文沒提」。不要替它想一個。

`## Personas` 五欄照 `references/persona-definition.md`：原文有講的填，沒講的寫
「（無資料）」，**「他要什麼」與「現在什麼讓他痛」尤其不要編**。

`## Assumptions / Constraints`：原文裡標題叫「假設」或語氣像既定事實的數字、門檻、
SLA，都抄進來並標 `← PRD §n`——它們是**沒驗證過的**，下一步會逐條變成題。

### 3. 最多問四題，而且只在對應的節空著時問

| 節空著 | 問 |
| --- | --- |
| `## Background` | 這解決誰的什麼問題？現在怎麼做、哪裡痛？ |
| `## Personas` | 受影響的人與系統有誰？ |
| `## Goal` | 做成之後長怎樣？怎麼知道做成了？ |
| `## Scope — In / Out` | 什麼明確不做？ |

答案照講的記，行尾標 `← Q-<n>`，並在 `## Open Questions` 開對應的一列（狀態
`已答`、`面向` 填 `—`）。這四題是分析師開場會問的，不是十六個面向——後者是
`bdd-discovery` 的事。

MUST NOT: 問第五題。問到第五題代表你在做 discovery；停下來，把它記成紅卡。

### 4. 粗切 story，底下留空

沿原文看得出的**使用者目標**切，一個目標一則：

```markdown
## User Stories

### US-1 · visitor-billing

身為 P-1（訪客），我想在出場前知道要付多少，以便付完就能開走 ← PRD §2

**切法**：暫定。沿原文 §2–§3 的訪客計費切；規則問出來之後可能重切。

#### FR

（空——規則在 DISCOVERY 問出來才寫）

#### NFR

（空）

#### Q

（空）
```

`### US-<n> · <slug>` 的形式與 `身為 P-<n>` 的要求見 `references/prd-format.md`——
slug 就是日後的 `.feature` 檔名，第一次就要是 kebab-case。

切不出來就只開一則，slug 用 feature 名。**標「暫定」**：這一刀是沿目標切的，沿規則
切要等 `bdd-discovery` 把規則問出來，那時 `story-splitting` 會重切。

MUST NOT: 替 story 填 FR 或 AC。你看到的只有原文，原文沒有規則；寫了就是編的，
`check_spec.py` 第 14 項會把沒有 `← Q-<n>` 或 `← PRD §n` 的規則抓出來。

### 收尾

- `## Open Questions`：需求方自己列的「還不確定」每一條開一列，狀態 `待答`、`面向` `—`，
  小節寫 `**該問誰**`。原文裡的「假設」段**不**在這裡開題——那是 `bdd-discovery`
  開場掃描的工作，這裡只抄進 `## Assumptions / Constraints`
- `## Document Overview`：`**狀態**：澄清中　**最後更新**：<今天>　**結束方式**：—`；
  核心關係人表登記需求方自己與他點名的人，`決策權` 照他講的填，不知道就寫「（未問）」
- 把整份 `prd.md` 給需求方看一遍，改到他說「這是我的意思」為止。改的是開頭幾節的
  措辭，**不是 `## 原文`**

## 產物

```
specs/<date>-<feature>/
└── prd.md     CAPTURE  ★ 狀態 澄清中；story 底下三個分組空著；## 原文 逐字
```

MUST: 只寫這一個檔。不建 `spec.md`、不建 `.feature`、不寫 `docs/` 底下任何東西。

## 完成後

告訴對方四件事：

1. `prd.md` 在哪、粗切了幾則 story、切法為什麼是暫定
2. 哪幾節原文沒提、留空了
3. `## Open Questions` 開了幾張紅卡（全部來自需求方自己列的）
4. 下一步：`bdd-discovery`，它會從 `## 原文` 掃隱含假設、把每則 story 問到有規則有例子
