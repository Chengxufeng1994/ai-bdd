# 把切票從 `bdd-plan` 拆成 `bdd-tasks`

**狀態**：已批　**日期**：2026-11-15

## 意圖

需求方的話：「Plan 是怎麼做，Tasks 是將怎麼做拆解。」以及：「我想要併行處理。」

`bdd-plan` 只留技術設計；新 skill `bdd-tasks` 把設計拆成票，一張票一個檔，讓幾張票能在
各自的 session／worktree 同時做而不互相覆寫。

### 為什麼

現在的格式下，兩張票同時跑會撞三個地方：

| 共寫的檔案 | 現在誰寫 | 並行時 |
| --- | --- | --- |
| `plan.md` | 每張票的 `bdd-implement` 勾驗收條件、寫 `## 缺口` | 每個 session 改同一份檔 |
| `features/<story>.feature` 的狀態 tag | 做完票的 session 翻 `@review` | 一檔一個 tag，同一則 story 的幾張票誰翻沒有規定 |
| `api/openapi.yaml`、`steps_*.go` 等 | 每張票各自加 | `Blocked by` 看不出這種衝突 |

`plan-format.md` 的〈為什麼不是「每張票一個檔」〉已經預留這一步：「真的要平行跑票時再拆。」

## 決定

### D1. 兩個 skill、兩種紀律

| | `bdd-plan` | `bdd-tasks`（新） |
| --- | --- | --- |
| 回答 | 怎麼做 | 把怎麼做拆成幾張、誰先誰後、誰能同時 |
| 讀 | `spec.md`、`.feature`、`discovery.md` 技術面向的答案 | `plan.md`、`spec.md`、`.feature` |
| 寫 | `plan.md`（技術節＋`## 缺口`） | `tasks/` |
| 紀律 | 不改需求 | **不改設計**——拆的時候發現設計有洞，加一列到 `plan.md` 的 `## 缺口`，不自己補 |

`bdd-plan` 現在的第 4 步（切票）、第 5 步（攤票確認）、三條切票規則、驗收條件要能失敗、
wide refactor 的例外，整段搬到 `bdd-tasks`，內容不改。

### D2. 產物

```
specs/<date>-<feature>/
├── plan.md            技術設計 ＋ ## 缺口（不再有 ## 票、## 前置整理）
└── tasks/
    ├── index.md       順序表：編號｜標題｜Blocked by｜共用檔案
    ├── 00-<slug>.md   前置整理（有的話）
    ├── 01-<slug>.md
    └── 02-<slug>.md
```

- **順序只住 `index.md`。** 當初不拆檔的理由是「順序拆進 N 個檔就不存在了」；索引表保住它。
  編號一旦寫出去不重排
- **票檔**沿用現在票的欄位（交付什麼行為、涵蓋場景、Blocked by、內迴圈、驗收條件），加兩樣：
  - `**共用檔案**`：這張票會改、而別張票也可能改的檔案（契約檔、step definition 檔、
    migration）。沒有就寫「無」
  - 票檔自己的 `## 缺口`：實作時發現的洞寫這裡，不寫 `plan.md`

### D3. 什麼叫可以同時跑

兩張票可以同時跑，要兩個條件都成立：

1. 各自的 `Blocked by` 都已完成
2. `共用檔案` 沒有交集

`bdd-tasks` 攤票給使用者確認時，把「哪幾張可以同時開始」講出來。

### D4. `bdd-implement` 跟著改

- 讀 `tasks/<票>.md` ＋ `plan.md` 的技術節；**只寫自己那張票檔**（勾驗收條件、寫缺口），
  不再寫 `plan.md`、不寫 `index.md`
- 沒帶票號時，從 `index.md` 取第一張 `Blocked by` 全完成、驗收條件還沒全勾的票
- 並行時一張票一個 worktree（`SHOULD`）
- **狀態 tag**：一張票做完，只有在「涵蓋同一個 `.feature` 的票全都勾完」時才翻 `@review`；
  否則留 `@wip`。並行時這個判斷在票合回主線之後做——worktree 裡看不到別張票的進度

### D5. 流程仍是六步

`bdd-tasks` 是 PLAN 步驟的第二個 skill，就像 SPEC 有好幾個 skill。不新增第七步——
「六步流程」寫在每個流程 skill 的 description 與 `docs/` 裡，改它的代價大於多一個名字的好處。
**這一條是我的判斷，不是需求方說的；要 TASKS 成為獨立的一步請講。**

### D6. 舊格式

`plan.md` 還有 `## 票` 的 feature（含 `runs/` 底下的）是舊格式。`bdd-implement` 遇到沒有
`tasks/index.md` 的 feature 時停下，說先跑 `bdd-tasks`，不兩種都支援。

## 非目標

- 不做自動派工（誰開幾個 agent、怎麼分配）——這份設計只讓產物**可以**並行，派工由使用者
  或外部工具決定
- 不寫計算波次的腳本——D3 兩個條件先由 skill 講出來；實跑發現算錯再加
- 不改三條切票規則與票的既有欄位
- 不動 `spec.md`、`discovery.md` 的格式與 `check_spec.py`、`status.py`

## 考慮過的做法

| 做法 | 取捨 | 結果 |
| --- | --- | --- |
| A. 不拆 skill，`bdd-plan` 把票寫成多檔 | 改動最小；但「只重排票、不動設計」得重跑整個 `bdd-plan`，而並行讓重排變常態 | 捨棄 |
| **B. 拆出 `bdd-tasks`，票一檔一張＋索引**（本文件） | 多一個交接點與一條紀律；兩半可各自重跑，票的 session 不共寫檔案 | **建議** |
| C. 票留在 `plan.md`，並行靠合併時解衝突 | 不用改格式；每次合併都衝突在同一份檔，缺口互相蓋掉 | 捨棄 |

## 會動到的檔案

- 新增 `skills/bdd-tasks/`（`SKILL.md`、`references/tasks-format.md`）
- `skills/bdd-plan/SKILL.md`、`references/plan-format.md`：拿掉票的部分
- `skills/bdd-implement/SKILL.md`：D4
- `skills/bdd-formulation/references/state-tags.md`：`@review` 的寫入時機
- `CLAUDE.md`、`README.md`（產物表）、`PLAN.md` 的 skill 表、`docs/ai-sdlc.md` 提到切票處

## 最不確定的地方

- **D4 的狀態 tag**：「合回主線後才判斷」需要有人在合併後做這件事，目前沒有 skill 擁有
  「合併」這個動作
- **`共用檔案` 靠切票時預測**：票還沒做，哪些檔會被改是推論；預測漏了，兩張票照樣會撞
