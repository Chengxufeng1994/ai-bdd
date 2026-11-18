---
name: plan
description: >
  SPEC 說要做成什麼，PLAN 說用什麼做、怎麼排——讀 `spec.md`（已對焦）與它的 `.feature`，寫 `plan.md`：
  前半是技術決定（Tech Stack、Commands、Project Structure、Code Style、Boundaries、Implementation
  Decisions、Testing Decisions，seam 攤給使用者確認），後半是執行計畫（概述、風險與影響程度、要不要並行與
  合回時會相撞的接合點）。不改 `spec.md`、不改 `.feature`、不切票——拆成票交給 `task-breakdown`。
  BDD 六步流程的 PLAN 步驟。
  觸發詞：「寫 plan」「執行計畫」「技術設計」「決定 seam」「驗收測試打在哪一層」「API 契約要怎麼定」
  「這個 feature 動哪些模組」「schema 怎麼定」「這批有什麼風險」「哪裡最危險」「要不要並行」
  「這批能不能同時做」「接合點有哪些」「formulation 寫完了下一步」。
  English: write the plan, technical design, decide the acceptance test seam, define the API contract,
  execution plan, what are the risks, can this run in parallel, list the join points.
disable-model-invocation: true
---

# PLAN — 用什麼做、怎麼排

| | |
| --- | --- |
| **職責** | 把一份已對焦的規格變成計畫：先定用什麼做、seam、模組、契約、每個 AC 在哪一層測，再排哪裡危險、能不能並行、合回時哪些檔會相撞 |
| **讀** | `spec.md`（`已對焦`）、涵蓋的 `.feature`、`discovery.md` 技術面向的已答題、專案既有的結構與文件、`docs/TECH.md`（有的話） |
| **寫** | `specs/ongoing/<date>-<feature>/plan.md` |
| **上游** | `formulation`（`.feature` 寫好之後，使用者叫 `/plan`） |
| **下游** | `task-breakdown`（照這份計畫拆成票） |
| **不做** | 不改 `spec.md`、不改 `.feature`、不切票、不訪談需求、不寫程式 |

**不改 `spec.md`。** 要做成什麼已經在 `spec.md` 裡、需求方說過「就緒」；這一步只決定用什麼做、怎麼排。
發現需求有洞，寫進 `plan.md` 的 `## 缺口` 交回 `spec`，不要順手補，也不要順口問掉。

## 使用時機

- `spec.md` 已對焦、`.feature` 也寫好了，要決定技術做法並排出計畫
- `task-breakdown` 拆票前，`plan.md` 還不存在
- `spec.md` 修補過、`.feature` 跟著改了，`plan.md` 要跟上

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| 還沒有 `spec.md`，或狀態不是 `已對焦` | `spec`（連 `discovery.md` 都沒有就先 `discovery`） |
| 還沒有 `.feature` | `formulation` |
| **要改 FR、AC、story、範圍** | `spec`——不是這裡 |
| `.feature` 寫了 seam 驗不到的行為 | `formulation`，本 skill 只指出來 |
| 要切成幾張票、排順序、看哪幾張能同時做 | `task-breakdown`，本 skill 寫完 `plan.md` 交給它 |
| 要實際寫 step definition 與產品程式碼 | `implement` |

## 參考檔案

- `references/plan-format.md` — `plan.md` 的骨架、執行計畫各節怎麼寫、範例
- `references/tech-format.md` — 技術決定七節的骨架、四選一的來源標記、範例。**第 3 步寫之前先讀**

## 前置確認（Ask user）

Check if: 同目錄的 `spec.md` 狀態是 `已對焦`、`## Open Questions` 是「（無）」，而且它涵蓋的 `.feature` 已經
寫好。任何一項不成立就停下，說該回上表的哪個 skill——規格還有洞時定的技術決定，定的是一份還不存在的需求。

---

## 流程

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| 1 | 一次讀完整批 | — |
| 2 | 決定 seam | 要，攤給使用者確認 |
| 3 | 寫技術決定七節 | — |
| 4 | 對照 `.feature`：有沒有 seam 驗不到的場景 | — |
| 5 | 寫 `## Risks`，再寫 `## Overview` | — |
| 6 | 問要不要並行、列接合點 | 要 |

### 1. 讀

| 讀什麼 | 為了什麼 |
| --- | --- |
| `spec.md` 的 `## User Stories` | 每條 FR 與 AC——seam 要看過具體的 AC 才推得出來，測試層級要逐個 AC 標；風險要指得出從哪條規則長出來 |
| `spec.md` 的 `## Assumptions / Constraints`、`## Non-Functional Requirements` | 外部限制與品質屬性，`## Tech Stack`、`## Boundaries`、`## Risks` 的來源 |
| `discovery.md` `## Open Questions` 裡技術面向、狀態 `已答` 的題 | **技術決定的原始素材**（seam、模組邊界、契約、既有資產、測試慣例），行尾標 `← Q-<n>` |
| 涵蓋的每個 `.feature` | 場景與 tag——第 4 步對照 seam；哪些場景會一起被一個風險影響 |
| 專案的 Makefile、README、`CLAUDE.md`、架構文件、既有程式碼 | `## Commands`、`## Project Structure`、`## Code Style` 的來源，標 `← 專案 <路徑>`；找既有的 seam；哪些檔是產生的、哪些介面會長大 |
| `docs/TECH.md`（有的話） | 技術現況：Tech Stack 到 Boundaries 從這裡抄，逐條確認現在還成立，標 `← 專案 docs/TECH.md` |

**一次讀完整批，不分批。** API 契約、schema、接合點與風險天生跨 story，只有並排才看得見衝突。

### 2. 決定 seam —— 驗收測試打在哪一層

`.feature` 決定驗什麼，seam 決定**打進系統的哪一層**：HTTP handler？application service？domain？這個決定會長成
step definition 的形狀。

| 規則 | 為什麼 |
| --- | --- |
| 優先用既有 seam | 每多一個 seam 就多一套測試替身與資料建構 |
| 取**最高**的那一層 | 越高越接近使用者真正做的事，也越不綁實作細節 |
| 整個變更的理想數量是**一個** | 兩個 seam 代表這批行為的入口不只一個，多半是切分沒切乾淨 |
| 需要新 seam 就在**你能做到的最高點**提出 | 往下挪一層通常是為了好寫，代價是驗到的東西離使用者更遠 |

**Ask user: 把 seam 攤出來確認再寫。** 三種情況：

| `discovery.md` 裡 `面向` 為 `seam` 的題 | 怎麼做 |
| --- | --- |
| `已答` | 只是確認，行尾標 `← Q-<n>` |
| `n/a`（答案欄寫了為什麼不適用） | 照那個理由選最外層能驗到行為的一層，行尾標 `← Q-<n>` |
| 沒有這題，或還 `待答` | 提出你的判斷、標 `← 推論`，`## Further Notes` 的「推論的技術決定」記一條「seam 沒問過」 |

MUST NOT: 把 seam 的討論變成需求訪談。談的時候發現需求有洞，寫進 `## 缺口` 交回 `spec`，**不要順手問掉**——
這一步沒有把答案寫進 `discovery.md` 的權限，順口問到的答案只活在這一次對話裡。

IMPORTANT: seam 一旦寫進 `plan.md` 就往下傳——IMPLEMENT 照著打，REVIEW 把「沒人同意過的 seam」當成 finding。

### 3. 寫技術決定

照 `references/tech-format.md` 填七節：`## Tech Stack`、`## Commands`、`## Project Structure`、`## Code Style`、
`## Boundaries`、`## Implementation Decisions`、`## Testing Decisions`。每條行尾標來源（四選一）；可以用
`推論`，但每一條同時列進 `## Further Notes` 的「推論的技術決定」。

### 4. 對照 `.feature`：seam 驗得到嗎

逐個場景問：它的 Then 在第 2 步定的 seam 上**看得到**嗎？例如 seam 是 HTTP API，「繳費機退回 500 元」「畫面
顯示請改投小鈔」是設備端的行為，HTTP 回應裡看不到。

| 結果 | 怎麼做 |
| --- | --- |
| 看得到 | 照常，`## Testing Decisions` 的測試層級表標它在哪一層 |
| 看不到，但換一個說法看得到（「API 回應這次不算付款、要求改付」） | `## 缺口` 一列，交回 `formulation` 改寫那個場景的 Then |
| 怎麼換都看不到 | `## 缺口` 一列，交回 `spec`：這條 AC 驗不了，要需求方決定是換一個驗法還是拿掉 |

MUST NOT: 自己改 `.feature` 的場景，或把驗不到的場景默默標成手動測試。

### 5. 寫 Risks 與 Overview

照 `references/plan-format.md`：

- `## Risks`：規格沒說、但實作一定會撞到的。每條指出從哪條規則長出來、標影響程度
- `## Overview`：`## Risks` 寫完才寫，一段摘要，不做新決定

### 6. 要不要並行（Ask user）

問一句：這批票要不要同時做。答案記進 `## Parallel Work`，行尾 `← 使用者 <日期>`。

| 回答 | 寫什麼 |
| --- | --- |
| **不要**（預設） | 只寫那一行 |
| **要** | 從 `## Implementation Decisions` 列出接合點：路徑、種類（`產生`／`追加`／`爭用`／`新建`）、為什麼會共用。寫法 → `references/plan-format.md`〈Parallel Work〉 |

IMPORTANT: 接合點要寫**路徑**，不是模組名。`task-breakdown` 只照抄這一節；這裡寫「HTTP 層」，
下一步就只能自己編路徑，而編出來的路徑在試跑中漏掉的正是合回時出事的那幾個檔。

---

## 不改規格：有洞怎麼辦

MUST NOT: 在這一步新增、刪除或改寫任何 FR、AC、story、範圍，或任何 `.feature` 場景。它們已經在 `spec.md`
與 `.feature` 裡、需求方說過「就緒」；這裡再決定一次，只會產生一份跟規格不一樣的規格。

判準：指著 `plan.md` 裡的任何一句話，說得出它是為了 `spec.md` 的哪條 FR 或 NFR 嗎？說不出來就是這一步
僭越了。

寫計畫時才發現的，寫進 `## 缺口`，再依種類處理：

| 種類 | 例子 | 怎麼處理 |
| --- | --- | --- |
| **技術選擇** | 哪個 library、欄位叫什麼 | `## Further Notes` 的「留給 IMPLEMENT 的技術選擇」 |
| **執行上的選擇** | 哪張票先做、某個檢查要不要加 | 留給 `task-breakdown` 或 IMPLEMENT 決定 |
| **場景驗不到** | 第 4 步找到的 | 交回 `formulation`（換說法）或 `spec`（換驗法） |
| **需求的洞** | 某條規則在某個情況下該怎麼做，沒人說過；兩條 FR 對同一個介面假設不一樣 | 告訴使用者跑 `spec`；它會再退回 `discovery` |

---

## 產物

一個 feature 一份 `specs/ongoing/<date>-<feature>/plan.md`。骨架照抄 `references/plan-format.md`，技術決定
七節照 `references/tech-format.md`。

MUST: 只寫入 `specs/ongoing/<date>-<feature>/plan.md`，**不修改專案既有的任何檔案**，包括 `spec.md`、
`discovery.md` 與 `.feature`。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] `spec.md` 與 `.feature` 沒有任何新增、刪除或改寫
- [ ] seam 攤給使用者確認過，回覆記在回報裡；談 seam 時發現的需求洞進了 `## 缺口`，沒有順口問掉
- [ ] 技術決定七節都在（空的寫了為什麼空）；每條行尾有來源，每條 `← 推論` 同時列在 `## Further Notes`
- [ ] `## Commands` 有跑場景與閘門兩行，或寫了「（沒有——要問使用者）」並進了 `## 缺口`
- [ ] 每個場景都對照過 seam；驗不到的進了 `## 缺口`，指名交回 `formulation` 或 `spec`
- [ ] `## Risks` 每條都指得出從哪條規則長出來、標了影響程度
- [ ] `## Overview` 只摘要已經寫的東西，沒有任何新決定
- [ ] `plan.md` 裡沒有票——沒有 `## 票`、沒有編號的工作項
- [ ] `## Parallel Work` 記了使用者對「要不要並行」的回答；答「要」時每個接合點都有路徑與種類

## 完成後

告訴對方：

1. **seam 是哪一層**、是確認 `discovery.md` 裡問過的還是推的、使用者確認了沒；動哪些模組、新的契約有哪些、
   幾個決定標了 `推論`
2. 幾個場景 seam 驗不到、交回給誰
3. 幾條風險、幾條影響程度「高」
4. 要不要並行；要的話有幾個接合點、幾個 `爭用`
5. `## 缺口` 有幾條，其中幾條要回 `spec`、幾條要回 `formulation`

`## 缺口` 沒有要回 `spec`／`formulation` 的 → 下一步由使用者叫 `/task-breakdown`。有 → 先回那一步。
