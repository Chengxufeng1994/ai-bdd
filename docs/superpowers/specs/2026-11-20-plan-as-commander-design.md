# `plan` 當指揮官：先設計 API／schema／domain model，再排計畫

**狀態**：待批　**日期**：2026-11-20

## 意圖

需求方的話：「plan 應該要像 addyosmani、mattpocock 那樣做」「plan 去分析 api、database schema、domain model
→ 拆 ticket」「應該新增一個技能叫做 api-and-interface-design，可以直接抄過來；plan 技能比較像是指揮官：
1. 收集資料 2. api design／database design 3. 就開始是 addyosmani 的流程」。選了：API 契約與 schema
寫在 `plan.md` 裡、`implement` 才搬進專案（(b)）；改寫 `plan`、保留 `task-breakdown`（做法 1）；`plan` 做到
依賴、風險、順序、checkpoint 與並行，寫票交給 `task-breakdown`（(i)）。

現在的 `plan.md` 是七節技術決定加執行計畫。API 契約只是 `## Implementation Decisions` 裡一段散文，
schema 也是；domain model 住在 `spec.md` 的 `## Domain Concepts`。第二次試跑的 `plan.md` 寫「HTTP 契約
以 `api/openapi.yaml` 為唯一來源」，但整個 PLAN 沒有一份能拿來對的契約。

流程不變：`discovery → spec → spec-by-example → formulation → plan → task-breakdown → implement → archive`。

## 目標、非目標、反目標

| | |
| --- | --- |
| **目標** | `plan.md` 結束時有三份可以直接照著寫程式的設計：完整的 OpenAPI 片段、資料表、聚合與不變條件；有依賴、風險、順序、checkpoint 與並行判斷，`task-breakdown` 讀它就能切票 |
| **非目標** | `data-plan`、`ui-plan`、`system-analysis` 這類分開的 planner；把契約拆成多檔；PLAN 階段改 `api/openapi.yaml` 或寫 migration |
| **反目標** | `plan.md` 裡的契約與實作時寫進 `api/openapi.yaml` 的不一樣，而沒人發現；三份設計各寫各的、互相矛盾（聚合邊界與表的切法對不上）；`api-and-interface-design` 抄過來後被改成另一套原則 |

## 做法

考慮過三個：

| 做法 | 取捨 | 結果 |
| --- | --- | --- |
| 1. 改寫 `plan`，保留 `task-breakdown` | 改動集中；`tasks.md`、`implement` 的平行與合回規則不動 | **選這個** |
| 2. 拆成 `api-plan`／`data-plan`／`domain-plan` 由 `plan` 協調（aixbdd） | 每份設計有自己的擁有者；但三份要並排才看得見衝突，拆開反而看不到；每一步都要使用者叫，PLAN 變四步 | 不做 |
| 3. `plan` 一路寫到票（addyosmani 原樣），拿掉 `task-breakdown` | 少一步；`plan` 變很大，剛定好的 `tasks.md` 格式與 `implement` 都要重來 | 不做 |

## 決定

### D1. 新增 `skills/api-and-interface-design/`

照 addyosmani/agent-skills 的 `skills/api-and-interface-design/SKILL.md`（MIT）**內容原樣保留**，英文不翻。只補
本 repo 稽核要的東西，不改任何一條原則：

- frontmatter 加 `disable-model-invocation: true`
- 標題下加職責卡，以及一段說明它在本流程裡的位置：由 `plan` 第 2 步 DELEGATE，輸入是 `spec.md` 的 FR／AC、涵蓋的
  `.feature`、專案既有的 `api/openapi.yaml`；產出是交回 `plan` 的 OpenAPI 片段，由 `plan` 寫進 `plan.md`；
  單獨叫用時只給設計建議，不寫檔
- 補 `## 使用時機`（中文，指向原文的 When to Use）、`## Skill Boundaries`（跟 `plan`、`implement` 的分界）
- 原文的 `## Verification` 保留，另加一項本流程的：每個 operation 指得出它服務哪幾條 AC
- 檔頭一行註明出處與授權；`RESOURCES.md` 加一列；repo 根若沒有對應的授權聲明，在 skill 目錄放 `LICENSE` 或 `NOTICE`

### D2. `plan` 的流程：收集 → 盤點介面 → 設計 → 排

需求方：「aixbdd 的 system-analysis 就等於 addyosmani 的 plan」。所以 `plan` 照 `system-analysis` 的角色：不自己
做每一份設計，而是盤點這次碰到哪些系統介面、排分析順序、把每個介面交給負責的那一步，再照 addyosmani 排計畫。

```
1. 收集   spec.md（FR／AC、## Domain Concepts 的詞彙、## Open Questions 技術面向的已答題）、.feature、
          docs/TECH.md、docs/DOMAIN-MODEL.md、docs/DATABASE-SCHEMA.md、api/openapi.yaml、既有程式碼
2. 盤點   系統介面（照 system-analysis 的兩條規則）：從需求部位出發，不先預設「前端、後端、資料庫」；
          每個介面寫端點類型、主要介面、需求原文依據（FR／AC 或 Q-<n>）
          分析順序（Wave）：哪個介面的結論是另一個的輸入就先做；同一波只放互不等待的
          對應：API 契約、事件、請求回應、錯誤碼 → api-and-interface-design；實體、欄位、生命週期、
          儲存 → plan 自己（Domain Model、Database Schema）；這次沒有 UI，不設 ui 那一路
3. 設計   seam 攤給使用者確認（照舊）
          Domain Model：聚合、不變條件、跨聚合規則
          Database Schema：表、欄位、索引、約束——每條不變條件指得出被哪個約束或聚合守住
          API 契約：DELEGATE api-and-interface-design → OpenAPI 片段
          照第 2 步的 Wave 順序做；三份並排對一次：聚合邊界 ↔ 表 ↔ operation 對得上
4. 排     依 addyosmani planning-and-task-breakdown 的 Step 2、5 與 Parallelization：依賴圖、
          Risks & Mitigations、順序與 checkpoint、能不能並行（「先定契約，再並行」）
          對照 .feature：seam 驗得到嗎（照舊）
          → STOP，告訴使用者 /task-breakdown
```

垂直切片與每張票的寫法（addyosmani Step 3、4）留給 `task-breakdown`。

### D3. `plan.md` 的骨架

```markdown
# 計畫
## Overview
## 系統介面        盤點結果與 Wave（system-analysis 的格式）
## Tech Stack／## Commands／## Project Structure／## Code Style／## Boundaries   照舊完整寫出
## Implementation Decisions   動哪些模組、程式碼層的介面與埠、架構決定、互動——不再放 API 契約與 schema
## Testing Decisions seam ＋ 每個 AC 在哪一層測（照舊）
## Domain Model
## Database Schema
## API 契約          ```yaml 的 OpenAPI 片段（新增與改動的 paths、components）
## 依賴圖
## Risks
## 順序與 Checkpoint
## Parallel Work
## 缺口
## Further Notes
```

來源標記照 `tech-format.md` 的五選一。`tech-format.md` 改寫成這份骨架的各節寫法。

### D4. domain model 從 `spec` 移到 `plan`

`spec.md` 的 `## Domain Concepts` 只留 `### 詞彙`（名詞、定義，不含型別），`### 聚合與不變條件` 移到 `plan.md` 的
`## Domain Model`。理由：聚合邊界是設計判斷，要跟 schema 與 API 並排才看得見衝突。`check_spec.py` 第 18 項照舊
要 `## Domain Concepts` 一節。

### D5. 契約從 `plan.md` 搬進專案

`implement` 做一張票時，把這張票涵蓋的 operation 與它用到的 schema，從 `plan.md` 的 `## API 契約` 照抄進
`api/openapi.yaml`（追加型接合點，照 `tasks-format.md` 的合回規則）；資料表照 `## Database Schema` 寫 migration。
`plan.md` 是來源：搬的時候不改形狀，要改就先回 `plan` 改片段。`task-breakdown` 在每張票的驗收條件加一項：搬進
`api/openapi.yaml` 的部分與 `plan.md` 片段一致。

### D6. 下游

| 檔 | 改什麼 |
| --- | --- |
| `archive` | 聚合改從 `plan.md` 的 `## Domain Model` 讀；schema 從 `## Database Schema`；`TECH.md` 從 `## 沿用的現況` 的差異與 `## Commands` |
| `task-breakdown` | 讀三份設計切票；每張票列它碰的 operation 與表 |
| `implement` | D5 |
| `spec`、`spec-format.md` | D4 |
| README、CLAUDE.md、PLAN.md、docs/ai-sdlc.md、RESOURCES.md | 跟著 |

## reviewer 問到、由我決定的

一位 reviewer 讀過這份文件後問了七題；都不改變要做的東西，以下是我的決定（不是需求方的）：

| 問題 | 決定 |
| --- | --- |
| `## Implementation Decisions` 去哪 | 留著，名字不變；只把 API 契約與 schema 移出去。`task-breakdown` 照舊從它抄「動哪些模組」、`implement` 照舊從它讀埠的形狀、`## Parallel Work` 照舊從它推接合點 |
| seam 不是 HTTP、或專案沒有 `api/openapi.yaml` | `## API 契約` 有 HTTP 對外介面時寫 OpenAPI 片段；沒有時寫程式碼層的介面契約（埠的方法、輸入輸出、錯誤），D5 的搬進 `api/openapi.yaml` 與那一項驗收條件只在有 OpenAPI 片段時適用 |
| 依賴圖與順序的節點是什麼、`task-breakdown` 跟不跟 | `plan` 的依賴圖以模組與契約為節點（哪個先存在別的才寫得出來），不是票；`task-breakdown` 切票時把它當**約束**：一張票不得排在它依賴的節點之前，波與 checkpoint 仍由 `task-breakdown` 排 |
| 「沿用的現況」縮成一句 | 不縮——那是我先前的提議，不是需求方要的。五節照舊完整寫出，`plan.md` 仍自成一體，`implement` 仍只讀票節與 `plan.md` |
| 不變條件與約束的對應進不進 `docs/DOMAIN-MODEL.md` | 不進。對應只寫在 `plan.md` 的 `## Database Schema`；`domain-model-format.md` 的「被誰保證只寫聚合」不變；「原本是什麼、為什麼改」跟著聚合從 `spec` 移到 `plan` |
| 舊格式要不要繼續吃 | 不吃，跟前幾次一樣標 BREAKING：`task-breakdown`／`implement`／`archive` 遇到舊格式（`spec.md` 有 `### 聚合與不變條件`、`plan.md` 沒有 `## Domain Model`）就停下，請使用者重跑 `plan` |
| 抄過來的 description 過不了稽核 F5（沒有中文觸發詞） | description 可以加中文觸發詞；**本文**不改。驗證改成「本文與原文的 diff 只有 D1 列的增補」 |

## 留給實作的人

- `plan/SKILL.md` 的步驟用 `READ`／`WRITE`／`CHECK`／`DELEGATE`／`ASK`／`STOP` 動詞格式寫（照 `spec`）
- 依賴圖用文字樹還是表
- `plan` 與 `api-and-interface-design` 的範例用圖書館借還書（不用 benchmark 的兩個 case）

## 驗證

- `audit_skill.py` 每個 skill 通過（含新的 `api-and-interface-design`）；`claude plugin validate .`；`python3 -m unittest discover -s tests`
- `api-and-interface-design/SKILL.md` 本文與原文 diff：只有 D1 列的增補（frontmatter 的 description 另外加了中文觸發詞）
- 重跑一次試跑到 `task-breakdown`，看 `plan.md` 的三份設計能不能直接切票、`implement` 第一張票搬契約時有沒有對不上
