---
name: plan
description: >
  SPEC 說要做成什麼、用什麼做，PLAN 說怎麼排——讀 `spec.md`（已對焦、技術部分已寫）與它的
  `.feature`，寫 `plan.md` 的執行計畫：一段概述、風險與影響程度、要不要並行與合回時會相撞的
  接合點。不改 `spec.md`、不切票——拆成票交給 `task-breakdown`。BDD 六步流程的 PLAN 步驟。
  觸發詞：「寫 plan」「執行計畫」「這批有什麼風險」「哪裡最危險」「要不要並行」
  「這批能不能同時做」「接合點有哪些」「formulation 寫完了下一步」。
  English: write the plan, execution plan, what are the risks, can this run in parallel,
  list the join points.
disable-model-invocation: true
---

# PLAN — 怎麼排

| | |
| --- | --- |
| **職責** | 把一份寫好技術部分的規格變成執行計畫：這批在做什麼、哪裡危險、能不能並行、合回時哪些檔會相撞 |
| **讀** | `spec.md`（`已對焦`、技術部分已寫）、涵蓋的 `.feature`、專案既有的結構 |
| **寫** | `specs/ongoing/<date>-<feature>/plan.md` |
| **上游** | `formulation`（`.feature` 寫好之後） |
| **下游** | `task-breakdown`（照 `spec.md` 技術部分與這份計畫拆成票） |
| **不做** | 不改 `spec.md`（需求部分與技術部分都不改）、不切票、不訪談、不寫程式 |

**不改 `spec.md`。** 技術設計已經在 `spec.md` 的技術部分；這一步只排它。發現技術部分有洞，
寫進 `plan.md` 的 `## 缺口` 交回 `spec`，不要順手補。

## 使用時機

- `spec.md` 的技術部分寫好、`.feature` 也寫好了，要知道這批有什麼風險、能不能並行
- `task-breakdown` 拆票前，`plan.md` 還不存在

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| 還沒有 `spec.md`，或狀態不是 `已對焦` | `spec`（連 `discovery.md` 都沒有就先 `discovery`） |
| `spec.md` 還沒有技術部分（seam、模組、契約、測試層級） | `spec` 第 6 步 |
| 還沒有 `.feature` | `formulation` |
| **要改 FR、AC、story，或 seam、模組、契約、schema** | `spec`——不是這裡 |
| 要切成幾張票、排順序、看哪幾張能同時做 | `task-breakdown`，本 skill 寫完 `plan.md` 交給它 |
| 要實際寫 step definition 與產品程式碼 | `implement` |

## 參考檔案

- `references/plan-format.md` — `plan.md` 的骨架、每一節怎麼寫、範例。**寫檔之前先讀**
- `../spec/references/tech-format.md` — `spec.md` 技術部分各節裝什麼

## 前置確認（Ask user）

Check if: 同目錄的 `spec.md` 狀態是 `已對焦`、有技術部分（至少 `## Implementation Decisions` 與
`## Testing Decisions`），而且它涵蓋的 `.feature` 已經寫好。任何一項不成立就停下，說該回上表的
哪個 skill——沒有技術設計時排的計畫，排的是一份還不存在的設計。

---

## 流程

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| 1 | 一次讀完整批 | — |
| 2 | 寫 `## Risks`，再寫 `## Overview` | — |
| 3 | 問要不要並行、列接合點 | 要 |

### 1. 讀

| 讀什麼 | 為了什麼 |
| --- | --- |
| `spec.md` 的 `## User Stories` | 每條 FR 與 AC——風險要指得出從哪條規則長出來 |
| `spec.md` 的技術部分 | `## Implementation Decisions` 推接合點，`## Testing Decisions` 的 seam 寫進概述 |
| `spec.md` 的 `## Assumptions / Constraints`、`## Non-Functional Requirements` | 外部限制與品質屬性，`## Risks` 的來源 |
| 涵蓋的每個 `.feature` | 場景與 tag——哪些場景會一起被一個風險影響 |
| 專案既有的結構 | 哪些檔是產生的、哪些介面會長大——接合點的種類 |

**一次讀完整批，不分批。** 接合點與風險天生跨 story，只有並排才看得見衝突。

### 2. 寫 Risks 與 Overview

照 `references/plan-format.md`：

- `## Risks`：規格沒說、但實作一定會撞到的。每條指出從哪條規則長出來、標影響程度
- `## Overview`：`## Risks` 寫完才寫，一段摘要，不做新決定

### 3. 要不要並行（Ask user）

問一句：這批票要不要同時做。答案記進 `## Parallel Work`，行尾 `← 使用者 <日期>`。

| 回答 | 寫什麼 |
| --- | --- |
| **不要**（預設） | 只寫那一行 |
| **要** | 從 `spec.md` 的 `## Implementation Decisions` 列出接合點：路徑、種類（`產生`／`追加`／`爭用`／`新建`）、為什麼會共用。寫法 → `references/plan-format.md`〈Parallel Work〉 |

IMPORTANT: 接合點要寫**路徑**，不是模組名。`task-breakdown` 只照抄這一節；這裡寫「HTTP 層」，
下一步就只能自己編路徑，而編出來的路徑在試跑中漏掉的正是合回時出事的那幾個檔。

---

## 不改規格：有洞怎麼辦

MUST NOT: 在這一步新增、刪除或改寫任何 FR、AC、story，或 `spec.md` 技術部分的任何決定。
它們已經在 `spec.md` 裡、需求方說過「就緒」、seam 也確認過；這裡再決定一次，只會產生一份
跟規格不一樣的規格。

判準：指著 `plan.md` 裡的任何一句話，說得出它是為了 `spec.md` 的哪條 FR、哪個技術決定嗎？
說不出來就是這一步僭越了。

`spec.md` 沒寫到、排計畫時才發現需要的，寫進 `## 缺口`，再依種類處理：

| 種類 | 例子 | 怎麼處理 |
| --- | --- | --- |
| **執行上的選擇** | 哪張票先做、某個檢查要不要加 | 留給 `task-breakdown` 或 IMPLEMENT 決定 |
| **技術部分的洞** | 某個介面兩條 FR 的假設不一樣、某個 AC 沒標測試層級 | 告訴使用者跑 `spec` 補技術部分；那是需求的洞時，`spec` 會再退回 `discovery` |

---

## 產物

一個 feature 一份 `specs/ongoing/<date>-<feature>/plan.md`。骨架照抄 `references/plan-format.md`。

MUST: 只寫入 `specs/ongoing/<date>-<feature>/plan.md`，**不修改專案既有的任何檔案**，包括 `spec.md`
與 `discovery.md`。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] `spec.md` 沒有任何新增、刪除或改寫——需求部分與技術部分都是
- [ ] `## Risks` 每條都指得出從哪條規則長出來、標了影響程度
- [ ] `## Overview` 只摘要 `spec.md` 與下面各節已經寫的東西，沒有任何新決定
- [ ] `plan.md` 裡沒有票——沒有 `## 票`、沒有編號的工作項
- [ ] `## Parallel Work` 記了使用者對「要不要並行」的回答；答「要」時每個接合點都有路徑與種類
- [ ] 技術部分的洞進了 `## 缺口`，並告訴使用者回 `spec`；沒有在這裡補

## 完成後

告訴對方：

1. 幾條風險、幾條影響程度「高」
2. 要不要並行；要的話有幾個接合點、幾個 `爭用`
3. `## 缺口` 有幾條，其中幾條要回 `spec`

下一步 `task-breakdown` 把 `spec.md` 的技術部分與這份計畫拆成票。有要回 `spec` 的缺口 → 先 `spec`。
