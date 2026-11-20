---
name: plan
description: >
  SPEC 說要做成什麼，PLAN 說用什麼做、怎麼排——讀 `spec.md`（已對焦）與它的 `.feature`，像指揮官一樣寫 `plan.md`：
  先盤點這次碰到的系統介面與分析順序（Wave），再設計三份可以直接照著寫程式的東西——Domain Model（聚合與不變條件）、
  Database Schema（表、欄位、約束）、API 契約（DELEGATE `api-and-interface-design` 寫 OpenAPI 片段，沒有 HTTP
  時寫程式碼層介面契約），並排對過；技術決定七節（seam 攤給使用者確認）；最後排依賴圖、風險、順序與
  checkpoint、要不要並行與接合點。不改 `spec.md`、不改 `.feature`、不動 `api/openapi.yaml`、不切票——拆成票交給
  `task-breakdown`。BDD 六步流程的 PLAN 步驟。
  觸發詞：「寫 plan」「執行計畫」「技術設計」「決定 seam」「驗收測試打在哪一層」「API 契約要怎麼定」
  「這個 feature 動哪些模組」「schema 怎麼定」「domain model 怎麼切」「聚合怎麼定」「盤點系統介面」
  「這批有什麼風險」「哪裡最危險」「依賴圖」「要不要並行」「這批能不能同時做」「接合點有哪些」
  「formulation 寫完了下一步」。
  English: write the plan, technical design, decide the acceptance test seam, define the API contract,
  design the database schema, design the domain model, list the system interfaces, dependency graph,
  execution plan, what are the risks, can this run in parallel, list the join points.
disable-model-invocation: true
---

# PLAN — 用什麼做、怎麼排

| | |
| --- | --- |
| **職責** | 把一份已對焦的規格變成計畫：盤點系統介面，設計聚合、資料表、API 契約並排對過，定技術決定與 seam，再排依賴、風險、順序與 checkpoint、能不能並行 |
| **讀** | `spec.md`（`已對焦`）、涵蓋的 `.feature`、`docs/TECH.md`、`docs/DOMAIN-MODEL.md`、`docs/DATABASE-SCHEMA.md`、`api/openapi.yaml`（有的話）、專案既有的結構與程式碼 |
| **寫** | `specs/ongoing/<date>-<feature>/plan.md` |
| **上游** | `formulation`（`.feature` 寫好之後，使用者叫 `/plan`） |
| **下游** | `task-breakdown`（照這份計畫拆成票）；`api-and-interface-design` 由第 3 步 DELEGATE |
| **不做** | 不改 `spec.md`、不改 `.feature`、不動 `api/openapi.yaml`、不寫 migration、不切票、不訪談需求、不寫程式 |

**不改 `spec.md`。** 要做成什麼已經在 `spec.md` 裡、需求方說過「就緒」；這一步只決定用什麼做、怎麼排。
發現需求有洞，寫進 `plan.md` 的 `## 缺口` 交回 `spec`，不要順手補，也不要順口問掉。

**指揮官，不是一人包辦。** `plan` 盤點這次碰到哪些系統介面、排分析順序、把每個介面交給負責的那一步——API 契約交給
`api-and-interface-design`，聚合與資料表自己做——再把三份並排對一次。三份分開寫、各自看起來都對，並排才看得見
聚合邊界與表的切法對不上。

## 使用時機

- `spec.md` 已對焦、`.feature` 也寫好了，要設計 API、schema、domain model 並排出計畫
- `task-breakdown` 拆票前，`plan.md` 還不存在
- `spec.md` 修補過、`.feature` 跟著改了，`plan.md` 要跟上

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| 還沒有 `spec.md`，或狀態不是 `已對焦` | `spec`（連訪談紀錄 `interview-log.md` 都沒有就先 `discovery`） |
| 還沒有 `.feature` | `formulation` |
| **要改 FR、AC、story、範圍** | `spec`——不是這裡 |
| `.feature` 寫了 seam 驗不到的行為 | `formulation`，本 skill 只指出來 |
| 只要 API 設計的建議，不寫 `plan.md` | `api-and-interface-design`（單獨叫用） |
| 要切成幾張票、排每張票的波次、寫票的驗收條件 | `task-breakdown`，本 skill 寫完 `plan.md` 交給它 |
| 要把契約搬進 `api/openapi.yaml`、寫 migration、寫 step definition 與產品程式碼 | `implement` |

## 參考檔案

- `references/plan-format.md` — `plan.md` 的骨架（整份的節序）、執行計畫各節（Risks、順序與 Checkpoint、Parallel Work、缺口）怎麼寫、範例
- `references/design-format.md` — `## 系統介面`、`## Domain Model`、`## Database Schema`、`## API 契約`、`## 依賴圖` 的骨架、三份並排對照、範例。**第 2 步寫之前先讀**
- `references/tech-format.md` — 技術決定七節的骨架、五選一的來源標記、範例。**第 3 步寫技術決定之前先讀**
- `../api-and-interface-design/SKILL.md` — 第 3 步 DELEGATE，設計 `## API 契約`
- addyosmani/agent-skills 的 `planning-and-task-breakdown`（MIT）——第 4 步的依賴圖、順序與 checkpoint、「先定契約，再並行」照它的 Step 2、Step 5 與 Parallelization Opportunities；垂直切片與每張票的寫法（Step 3、4）留給 `task-breakdown`
- Waterball-Software-Academy/aixbdd 的 `system-analysis`（Apache-2.0）——第 2 步盤點介面與 Wave 的兩條規則

## 前置確認（Ask user）

Check if: 同目錄的 `spec.md` 狀態是 `已對焦`、`## Open Questions` 沒有待答的補問題（定義見 `../spec/references/spec-format.md`〈寫規格時發現的洞〉），而且它涵蓋的 `.feature` 已經
寫好。任何一項不成立就停下，說該回上表的哪個 skill——規格還有洞時定的技術決定，定的是一份還不存在的需求。

---

## 流程

每一步都寫成一串動作，開頭的詞說明這一步在做什麼：

| 詞 | 意思 |
| --- | --- |
| `READ` | 讀檔，不寫 |
| `WRITE` | 寫進 `plan.md`（本 skill 不寫別的檔） |
| `CHECK` | 對照一條規則判是或否；不成立就照那一行寫的去處 |
| `DELEGATE` | 讀另一個 skill 的 `SKILL.md`，照它做完再回來；只用在第 3 步 |
| `ASK` | 問使用者，等他回答 |
| `STOP` | 回報後結束這一輪，不往下做 |

每一步最後一行是**做完的樣子**與下一步；做完的樣子不成立，就不往下走。

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| 1 | 收集：一次讀完整批 | — |
| 2 | 盤點：系統介面、Wave、對應 | — |
| 3 | 設計：seam → Domain Model、Database Schema、API 契約 → 並排對照 → 技術決定七節 | 要，seam 攤給使用者確認 |
| 4 | 排：依賴圖、Risks、順序與 Checkpoint、Parallel Work、seam 驗得到嗎 → Overview → 停 | 要，問要不要並行 |

### 1. 收集

1. READ 照下表一次讀完：

   | 讀什麼 | 為了什麼 |
   | --- | --- |
   | `spec.md` 的 `## User Stories` | 每條 FR 與 AC——介面要指得出需求依據，seam 要看過具體的 AC 才推得出來，測試層級要逐個 AC 標；風險要指得出從哪條規則長出來 |
   | `spec.md` 的 `## Domain Concepts` 的 `### 詞彙` | 聚合、表、operation 的命名照這些詞，不另起一套 |
   | `spec.md` 的 `## Assumptions / Constraints`、`## Non-Functional Requirements` | 外部限制與品質屬性，`## Tech Stack`、`## Boundaries`、`## Risks` 的來源 |
   | `spec.md` `## Open Questions` 裡技術面向、狀態 `已答` 的題 | **技術決定的原始素材**（seam、模組邊界、契約、既有資產、測試慣例），行尾標 `← Q-<n>` |
   | 涵蓋的每個 `.feature` | 場景與 tag——第 4 步對照 seam；哪些場景會一起被一個風險影響 |
   | `docs/TECH.md`（有的話） | 技術現況：Tech Stack 到 Boundaries 從這裡抄，逐條確認現在還成立，標 `← 專案 docs/TECH.md` |
   | `docs/DOMAIN-MODEL.md`、`docs/DATABASE-SCHEMA.md`（有的話） | 既有的聚合與表——這次是新增還是改動；改動要寫「原本是什麼、為什麼改」 |
   | `api/openapi.yaml`（有的話） | 既有的契約——`## API 契約` 只放新增與改動的，錯誤格式與命名照它 |
   | 專案的 Makefile、README、`CLAUDE.md`、架構文件、既有程式碼 | `## Commands`、`## Project Structure`、`## Code Style` 的來源，標 `← 專案 <路徑>`；找既有的 seam；哪些檔是產生的、哪些介面會長大 |

2. CHECK `spec.md` 還有 `### 聚合與不變條件`（舊格式，聚合寫在 SPEC）→ 當成第 3 步 `## Domain Model` 的素材讀，
   照本 skill 的規則重新判，不照抄；不改 `spec.md`。

**一次讀完整批，不分批。** API 契約、schema、聚合、接合點與風險天生跨 story，只有並排才看得見衝突。

**做完的樣子**：上表每一列讀過；沒有的檔（第一個 feature 沒有 `docs/`、專案沒有 `api/openapi.yaml`）知道它沒有。→ 第 2 步

### 2. 盤點

1. READ `references/design-format.md`〈系統介面〉，骨架照抄。
2. WRITE `## 系統介面`：從需求部位出發——FR／AC 裡的使用者行為、外部依賴、資料責任——切出介面，每個寫
   `端點類型`、`主要介面`、`需求原文依據`（FR／AC 或 `Q-<n>`）。不先預設「前端、後端、資料庫」。
3. CHECK 每個介面的 `需求原文依據` 指得出一條 FR／AC 或一題 `已答` 的 `Q-<n>`。理由是「系統通常都有」的 →
   刪掉那個介面；需求裡看得到責任、卻沒有一條規則說它 → `## 缺口` 交回 `spec`。
4. WRITE `### Wave <n>`：哪個介面的結論是另一個的輸入就先做；每波寫 `平行分析介面` 與 `安排理由`。
5. CHECK 每一波只放互不等待的介面：某介面要先知道同一波另一個的流程、外部限制、輸入輸出或資料邊界 → 拆到下一波。
6. WRITE `### 對應`：API 契約、事件、請求回應、錯誤碼 → `api-and-interface-design`；實體、欄位、生命週期、
   儲存 → `plan`（`## Domain Model`、`## Database Schema`）。這次沒有畫面就不設 UI 那一路。

**做完的樣子**：每個介面都有需求依據，每一波只放互不等待的，每個面向都對到一個設計者與一節。→ 第 3 步

### 3. 設計

#### 決定 seam —— 驗收測試打在哪一層

`.feature` 決定驗什麼，seam 決定**打進系統的哪一層**：HTTP handler？application service？domain？這個決定會長成
step definition 的形狀。先問它，是因為 `## Testing Decisions` 照它寫、第 4 步對照 `.feature` 也是對著它驗；
`## API 契約` 寫 OpenAPI 片段還是程式碼層介面契約，看的是有沒有 HTTP 對外介面，不是 seam。

| 規則 | 為什麼 |
| --- | --- |
| 優先用既有 seam | 每多一個 seam 就多一套測試替身與資料建構 |
| 取**最高**的那一層 | 越高越接近使用者真正做的事，也越不綁實作細節 |
| 整個變更的理想數量是**一個** | 兩個 seam 代表這批行為的入口不只一個，多半是切分沒切乾淨 |
| 需要新 seam 就在**你能做到的最高點**提出 | 往下挪一層通常是為了好寫，代價是驗到的東西離使用者更遠 |

1. ASK 把 seam 攤出來確認再寫。三種情況：

   | `## Open Questions` 裡 `面向` 為 `seam` 的題 | 怎麼做 |
   | --- | --- |
   | `已答` | 只是確認，行尾標 `← Q-<n>` |
   | `n/a`（答案欄寫了為什麼不適用） | 照那個理由選最外層能驗到行為的一層，行尾標 `← Q-<n>` |
   | 沒有這題，或還 `待答` | 提出你的判斷給使用者確認：他確認或改了，標 `← 使用者 <日期>`；沒得到回答才標 `← 推論`，`## Further Notes` 的「推論的技術決定」記一條「seam 沒問過」 |

MUST NOT: 把 seam 的討論變成需求訪談。談的時候發現需求有洞，寫進 `## 缺口` 交回 `spec`，**不要順手問掉**——
這一步沒有把答案寫進 `## Open Questions` 的權限，順口問到的答案只活在這一次對話裡。

IMPORTANT: seam 一旦寫進 `plan.md` 就往下傳——IMPLEMENT 照著打，REVIEW 把「沒人同意過的 seam」當成 finding。

#### 三份設計

三份的順序固定：Domain Model → Database Schema → API 契約。Wave 排的是**介面的分析**順序：寫每一份時，照第 2 步的
Wave 順序逐個介面處理。格式與範例都在 `references/design-format.md`。

2. WRITE `## Domain Model`：這次新增或改動的聚合、不變條件、跨聚合的規則。`被誰保證` 只寫聚合；動到
   `docs/DOMAIN-MODEL.md` 既有的聚合，寫「原本是什麼、為什麼改」。
3. WRITE `## Database Schema`：表、欄位、索引、約束；每張表寫它屬於哪個聚合，每個欄位寫 `守住哪條不變條件`。
4. DELEGATE `../api-and-interface-design/SKILL.md`：交給它這次的 FR／AC、涵蓋的 `.feature`、`api/openapi.yaml`
   與剛寫的 Domain Model；它交回 OpenAPI 片段（有 HTTP 對外介面時）或程式碼層介面契約（沒有時）。WRITE 進
   `## API 契約`，每個 operation 寫出它服務哪幾條 AC。
5. CHECK 三份並排對一次（判準表見 `references/design-format.md`〈三份並排對照〉）：每條不變條件有一個聚合、資料庫
   守得住的有一個約束；每張表屬於一個聚合；每個 operation 最多寫一個聚合，除非一條跨聚合的規則說要。
   對不上 → `## 缺口`，不順手改掉其中一份。
6. CHECK 每個 operation 指得出 AC。指不出 → 刪掉那個 operation；AC 要求的行為沒有 operation → 補上，或是需求的洞就
   進 `## 缺口`。

MUST NOT: 在這一步改 `api/openapi.yaml` 或寫 migration。`plan.md` 是契約與 schema 的來源；`implement` 做票時
才照抄進專案。

#### 技術決定

7. READ `references/tech-format.md`。
8. WRITE 七節：`## Tech Stack`、`## Commands`、`## Project Structure`、`## Code Style`、`## Boundaries`、
   `## Implementation Decisions`、`## Testing Decisions`。前五節照舊完整寫出——`plan.md` 自成一體；
   `## Implementation Decisions` 寫動哪些模組、程式碼層的介面與埠、架構決定、互動，**不放** API 契約與 schema；
   `## Testing Decisions` 寫 seam 與每個 AC 在哪一層測。每條行尾標來源（五選一）；可以用 `推論`，但每一條同時
   列進 `## Further Notes` 的「推論的技術決定」。

**做完的樣子**：seam 使用者確認過（或標了 `推論` 並記進 Further Notes）；三份設計都寫了、並排對過，對不上的進了
`## 缺口`；技術決定七節都在。→ 第 4 步

### 4. 排

照 addyosmani `planning-and-task-breakdown` 的 Step 2（依賴圖）、Step 5（順序與 checkpoint）與 Parallelization
Opportunities；垂直切片與每張票怎麼寫留給 `task-breakdown`。

1. WRITE `## 依賴圖`（格式 `references/design-format.md`〈依賴圖〉）：節點是模組與契約，不是票——哪個先存在，別的
   才寫得出來。
2. WRITE `## Risks`（格式 `references/plan-format.md`〈Risks〉）：規格沒說、但實作一定會撞到的。每條指出從哪條規則
   長出來、標影響程度。
3. WRITE `## 順序與 Checkpoint`（格式 `references/plan-format.md`〈順序與 Checkpoint〉）：被依賴的節點先做，
   「高」風險的節點往前；每兩三個節點一個 checkpoint，命令照 `## Commands`。
4. ASK 這批票要不要同時做。答案記進 `## Parallel Work`，行尾 `← 使用者 <日期>`：

   | 回答 | 寫什麼 |
   | --- | --- |
   | **不要**（預設） | 只寫那一行 |
   | **要** | 從 `## Implementation Decisions` 與 `## API 契約` 列出接合點：路徑、種類（`產生`／`追加`／`爭用`／`新建`）、為什麼會共用。**先定契約，再並行**：共用同一份契約的票等 `## API 契約` 定了才同時做。寫法 → `references/plan-format.md`〈Parallel Work〉 |

   IMPORTANT: 接合點要寫**路徑**，不是模組名。`task-breakdown` 只照抄這一節；這裡寫「HTTP 層」，下一步就只能
   自己編路徑，而編出來的路徑在試跑中漏掉的正是合回時出事的那幾個檔。
5. CHECK 逐個場景問：它的 Then 在第 3 步定的 seam 上**看得到**嗎？例如 seam 是 HTTP API，「借書單印出到期日」
   「櫃台刷書機嗶一聲」是設備端的行為，HTTP 回應裡看不到。

   | 結果 | 怎麼做 |
   | --- | --- |
   | 看得到 | 照常，`## Testing Decisions` 的測試層級表標它在哪一層 |
   | 看不到，但換一個說法看得到（「API 回應帶這次借閱的到期日」） | `## 缺口` 一列，交回 `formulation` 改寫那個場景的 Then |
   | 怎麼換都看不到 | `## 缺口` 一列，交回 `spec`：這條 AC 驗不了，要需求方決定是換一個驗法還是拿掉 |

   MUST NOT: 自己改 `.feature` 的場景，或把驗不到的場景默默標成手動測試。
6. WRITE `## Overview`：上面都寫完才寫，一段摘要，點名三份設計各新增了什麼，不做新決定。
7. STOP 照〈完成後〉回報，告訴使用者下一步是 `/task-breakdown`（`## 缺口` 有要回 `spec`／`formulation` 的就先回那一步）。

**做完的樣子**：依賴圖、Risks、順序與 Checkpoint、Parallel Work、Overview 都寫了；每個場景對照過 seam；回報後停下。

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
| **三份設計對不上** | 第 3 步並排時找到的：表屬於兩個聚合、不變條件沒人守、operation 改兩個聚合 | 判斷錯在設計 → 回 `plan` 改那一份；錯在需求（規則沒說、或兩條 FR 打架）→ 交回 `spec` |
| **需求的洞** | 某條規則在某個情況下該怎麼做，沒人說過；兩條 FR 對同一個介面假設不一樣 | 告訴使用者跑 `spec`；它開成補問題、交給 `discovery` 問 |

---

## 產物

一個 feature 一份 `specs/ongoing/<date>-<feature>/plan.md`。骨架照抄 `references/plan-format.md`，技術決定
七節照 `references/tech-format.md`，`## 系統介面`、三份設計與 `## 依賴圖` 照 `references/design-format.md`。

MUST: 只寫入 `specs/ongoing/<date>-<feature>/plan.md`，**不修改專案既有的任何檔案**，包括 `spec.md`、`.feature` 與 `api/openapi.yaml`。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] `spec.md`、`.feature`、`api/openapi.yaml` 沒有任何新增、刪除或改寫
- [ ] 每個系統介面有需求依據（FR／AC 或 `Q-<n>`）；每一波只放互不等待的介面
- [ ] seam 攤給使用者確認過，回覆記在回報裡；談 seam 時發現的需求洞進了 `## 缺口`，沒有順口問掉
- [ ] 三份設計並排對過，對不上的進了 `## 缺口`
- [ ] API 契約每個 operation 指得出 AC；有 HTTP 時是 OpenAPI 片段、沒有時是程式碼層介面契約
- [ ] 技術決定七節都在（空的寫了為什麼空），`## Implementation Decisions` 裡沒有 API 契約與 schema；每條行尾有來源，每條 `← 推論` 同時列在 `## Further Notes`
- [ ] `## Commands` 有跑場景與閘門兩行，或寫了「（沒有——要問使用者）」並進了 `## 缺口`
- [ ] 每個場景都對照過 seam；驗不到的進了 `## 缺口`，指名交回 `formulation` 或 `spec`
- [ ] `## 依賴圖` 的節點是模組與契約；`## Risks` 每條指得出從哪條規則長出來、標了影響程度；`## Overview` 沒有新決定
- [ ] `plan.md` 裡沒有票——沒有 `## 票`、沒有編號的工作項；`## Parallel Work` 記了使用者的回答，答「要」時每個接合點都有路徑與種類

## 完成後

告訴對方：

1. **系統介面**盤出幾個、分幾波；**seam 是哪一層**、是確認 `## Open Questions` 裡問過的還是推的、使用者確認了沒
2. **三份設計**：幾個聚合（新增／改動）、幾張表、幾個 operation，`## API 契約` 是 OpenAPI 片段還是程式碼層契約；並排時對不上幾處
3. 動哪些模組、幾個決定標了 `推論`；幾個場景 seam 驗不到、交回給誰
4. 幾條風險、幾條影響程度「高」；幾個 checkpoint
5. 要不要並行；要的話有幾個接合點、幾個 `爭用`
6. `## 缺口` 有幾條，其中幾條要回 `spec`、幾條要回 `formulation`

`## 缺口` 沒有要回 `spec`／`formulation` 的 → 下一步由使用者叫 `/task-breakdown`。有 → 先回那一步。
