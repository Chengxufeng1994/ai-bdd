---
name: plan
description: >
  SPEC 說要做成什麼，PLAN 說用什麼做、怎麼排、切成哪幾張票——照 addyosmani 的 planning-and-task-breakdown 五步走：
  唯讀收集 `spec.md`（已對焦）與專案現況；Step 2 之前先定 seam、DELEGATE `formulation` 寫 `.feature`，再設計三份可以直接照著寫程式的東西——Domain Model
  （聚合與不變條件）、Database Schema（表、欄位、約束）、API 契約（DELEGATE `api-and-interface-design` 寫 OpenAPI
  片段，沒有 HTTP 時寫程式碼層介面契約），並排對過，再定技術決定七節（seam 攤給使用者確認）；然後排依賴圖、
  切成 tracer bullet、寫每張票、排序與 checkpoint、問要不要並行，把切法攤給使用者看、點頭才寫
  `tasks.md`。不改 `spec.md`、`.feature` 只經 `formulation` 寫、不動 `api/openapi.yaml`、不寫程式。BDD 六步流程的 PLAN 步驟。
  觸發詞：「寫 plan」「執行計畫」「技術設計」「決定 seam」「驗收測試打在哪一層」「API 契約要怎麼定」
  「這個 feature 動哪些模組」「schema 怎麼定」「domain model 怎麼切」「聚合怎麼定」「依賴圖」
  「這批有什麼風險」「哪裡最危險」「這批要怎麼切」「拆成幾張票」「切票」「拆成 tasks」「先做哪一個」
  「排實作順序」「這些票誰卡誰」「要不要並行」「哪幾張可以同時做」「接合點有哪些」
  「spec 已對焦了下一步」「寫 .feature 之後要做什麼」。
  English: write the plan, technical design, decide the acceptance test seam, define the API contract,
  design the database schema, design the domain model, dependency graph, what are the risks,
  break the plan into tasks, slice the work into tickets, vertical slices, what order should I build these,
  which tickets can run in parallel, list the join points.
disable-model-invocation: true
---

# PLAN — 用什麼做、怎麼排、切成哪幾張票

| | |
| --- | --- |
| **職責** | 把一份已對焦的規格變成計畫與票：定 seam、交 `formulation` 寫 `.feature`，設計聚合、資料表、API 契約並排對過，定技術決定，排依賴圖，切成一疊 tracer bullet、排好順序與 checkpoint |
| **讀** | `spec.md`（`已對焦`）、`formulation` 寫出的 `.feature`、`docs/TECH.md`、`docs/DOMAIN-MODEL.md`、`docs/DATABASE-SCHEMA.md`、`api/openapi.yaml`（有的話）、專案既有的結構與程式碼；再跑時同目錄既有的 `plan.md` 與 `tasks.md`（或舊的 `tasks/`） |
| **寫** | `specs/ongoing/<date>-<feature>/plan.md`、`specs/ongoing/<date>-<feature>/tasks.md`（一張票一個票節）；`features/*.feature` 由 DELEGATE 的 `formulation` 寫 |
| **上游** | `spec`（`已對焦` 之後，使用者叫 `/plan`） |
| **下游** | `implement <票號>`（一張票一個全新 session）；`formulation` 與 `api-and-interface-design` 由〈設計〉DELEGATE |
| **不做** | 不改 `spec.md`、不自己寫或改 `.feature`（只經 `formulation`）、不動 `api/openapi.yaml`、不寫 migration、不訪談需求、不寫程式、不動開工了的票（只准補 `接合點`／`Blocked by`／`大小`）、不派工 |

**Plan Mode：只寫 `plan.md` 與 `tasks.md`，不寫程式、不改 `spec.md`、切法點頭才寫。**
- 程式碼、migration、`api/openapi.yaml` 都是 `implement` 做票時才碰
- 需求有洞寫進 `## 缺口` 交回 `spec`，不順手補、不順口問掉——要做成什麼，需求方已經在 `spec.md` 說過「就緒」
- 票寫進 `tasks.md` 就有人照著開工，切錯的代價在實作時才付——所以先攤給使用者

## 使用時機

- `spec.md` 已對焦，要寫 `.feature`、設計 API、schema、domain model、排計畫、切成票
- `spec.md` 修補過，`.feature`、`plan.md` 與 `tasks.md` 要跟上
- 技術設計沒變，但票切得不對或順序要調，只重切票（〈再跑〉「只重切」那一條）
- 一批工作大到一個 session 做不完，要拆給多個 session 接力或同時做

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| 還沒有 `spec.md`，或狀態不是 `已對焦` | `spec`（連訪談紀錄 `interview-log.md` 都沒有就先 `discovery`） |
| **要改 FR、AC、story、範圍** | `spec`——不是這裡 |
| `.feature` 寫了 seam 驗不到的行為 | 〈設計〉第 10 步：換說法看得到就再 DELEGATE `formulation` 改寫，怎麼換都看不到就交回 `spec` |
| 只想把 `spec.md` 翻成 `.feature`，不排計畫 | `formulation`（單獨叫用） |
| 要把太大的 **story** 切小 | `story-splitting`——它切的是需求，本 skill 切的是實作 |
| 只要 API 設計的建議，不寫 `plan.md` | `api-and-interface-design`（單獨叫用） |
| 要把契約搬進 `api/openapi.yaml`、寫 migration、寫 step definition 與產品程式碼、實際做某一張票 | `implement` |
| 派幾個 agent、誰做哪一張 | 不在本 skill；這裡只講哪幾張**可以**同時做 |

## 參考檔案

- `references/plan-format.md` — `plan.md` 的骨架（整份的節序）、`## Overview`、`## Risks`、`## Parallel Work`、
  `## 缺口` 怎麼寫、範例。**〈設計〉開始寫 `plan.md` 之前先讀**
- `references/design-format.md` — `## Domain Model`、`## Database Schema`、`## API 契約`、`## 依賴圖` 的骨架、三份並排對照、
  範例。**〈設計〉寫三份之前先讀**
- `references/tech-format.md` — 技術決定七節的骨架、五選一的來源標記、範例。**〈設計〉寫技術決定之前先讀**
- `references/tasks-format.md` — `tasks.md` 的骨架、誰寫哪幾行、契約與資料表、接合點與同時跑、大小、驗證、檢查點、
  三條切票規則、驗收條件要能失敗、wide refactor 的例外。**Step 3 切票之前先讀**
- `../formulation/SKILL.md` — 〈設計〉DELEGATE，把 `spec.md` 寫成 `.feature`
- `../api-and-interface-design/SKILL.md` — 〈設計〉DELEGATE，設計 `## API 契約`
- addyosmani/agent-skills 的 `planning-and-task-breakdown`（MIT）——整個流程的骨架：Step 1 唯讀收集（Plan Mode）、
  Step 2 依賴圖、Step 3 垂直切片、Step 4 每張票的寫法、Step 5 順序與 checkpoint；Task Sizing（XS–XL，L 以上要拆）、
  「先定契約，再並行」、「不覆寫還有沒勾的計畫」、`## Open Questions`（本專案的 `## 缺口`）、Common Rationalizations
  與 Red Flags。本專案加上：Step 2 之前插一步〈設計〉
- mattpocock/skills 的 `to-tickets`——tracer bullet 的三條規則、blocking edges、前置整理、wide refactor 的
  expand–contract，以及〈ASK 切法〉的 Quiz：粒度、Blocked by、合併或再拆三問，點頭才寫

## 前置確認（Ask user）

Check if: 同目錄的 `spec.md` 狀態是 `已對焦`、`## Open Questions` 沒有待答的補問題（定義見 `../spec/references/spec-format.md`〈寫規格時發現的洞〉）。
任何一項不成立就停下，說該回上表的哪個 skill——規格還有洞時定的技術決定，定的是一份還不存在的需求。

---

## 流程

每一步都寫成一串動作，開頭的詞說明這一步在做什麼：

| 詞 | 意思 |
| --- | --- |
| `READ` | 讀檔，不寫 |
| `WRITE` | 寫進 `plan.md` 的一節，或最後寫 `tasks.md`（本 skill 不寫別的檔） |
| `CHECK` | 對照一條規則判是或否；不成立就照那一行寫的去處 |
| `DELEGATE` | 讀另一個 skill 的 `SKILL.md`，照它做完再回來；只用在〈設計〉 |
| `ASK` | 問使用者，等他回答 |
| `STOP` | 回報後結束這一輪，不往下做 |

每一步第一行是**產出**，最後一行是**做完的樣子**與下一步；做完的樣子不成立，就不往下走。

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| Step 1 收集 | 唯讀：一次讀完整批與專案現況 | — |
| 設計 | seam → `formulation` 寫 `.feature` → Domain Model、Database Schema、API 契約 → 並排對照 → 技術決定七節 → 場景在 seam 上驗得到嗎 | 要，seam 攤給使用者確認 |
| Step 2 依賴圖 | 節點是模組與契約；Risks | — |
| Step 3 垂直切片 | 最薄的完整路徑先，其餘依 FR；wide refactor 用 expand–contract | — |
| Step 4 寫任務 | 每張票的行為、驗收條件、驗證命令、Blocked by、模組、契約、資料表、接合點、大小 | — |
| Step 5 排序與 Checkpoint | 高風險在前、每 2–3 張票一個 checkpoint；要不要並行 | 要，問要不要並行 |
| ASK 切法 | 把票的清單攤給使用者：粒度、Blocked by、合併或再拆 | 要，點頭才寫 |
| WRITE | 寫 `tasks.md`、補 `## Overview` → 停 | — |

### Step 1 收集

產出：無，只讀——知道每一份輸入有沒有、舊格式（`spec.md` 的聚合一節、`tasks/` 目錄）有沒有

1. READ 照下表一次讀完；只讀，不寫任何檔、不寫程式：

   | 讀什麼 | 為了什麼 |
   | --- | --- |
   | `spec.md` 的 `## User Stories` | 每條 FR 與 AC——seam 要看過具體的 AC 才推得出來，測試層級要逐個 AC 標，票要指得出涵蓋哪幾條；風險要指得出從哪條規則長出來 |
   | `spec.md` 的 `## Domain Concepts` 的 `### 詞彙` | 聚合、表、operation、票的標題照這些詞，不另起一套 |
   | `spec.md` 的 `## Assumptions / Constraints`、`## Non-Functional Requirements` | 外部限制與品質屬性，`## Tech Stack`、`## Boundaries`、`## Risks` 的來源 |
   | `spec.md` `## Open Questions` 裡技術面向、狀態 `已答` 的題 | **技術決定的原始素材**（seam、模組邊界、契約、既有資產、測試慣例），行尾標 `← Q-<n>` |
   | 既有的 `.feature`（有的話） | 已經在跑的場景——這次的 `.feature` 由〈設計〉DELEGATE `formulation` 寫，不在這裡讀 |
   | `docs/TECH.md`（有的話） | 技術現況：Tech Stack 到 Boundaries 從這裡抄，逐條確認現在還成立，標 `← 專案 docs/TECH.md` |
   | `docs/DOMAIN-MODEL.md`、`docs/DATABASE-SCHEMA.md`（有的話） | 既有的聚合與表——這次是新增還是改動；改動要寫「原本是什麼、為什麼改」 |
   | `api/openapi.yaml`（有的話） | 既有的契約——`## API 契約` 只放新增與改動的，錯誤格式與命名照它 |
   | 專案的 Makefile、README、`CLAUDE.md`、架構文件、既有程式碼 | `## Commands`、`## Project Structure`、`## Code Style` 的來源，標 `← 專案 <路徑>`；找既有的 seam；哪些檔是產生的、哪些介面會長大；前置整理（prefactoring）的機會 |
   | 同目錄既有的 `plan.md`、`tasks.md`（再跑時） | 哪些票開工了——見〈再跑：重切，不重編〉 |

2. CHECK `spec.md` 還有 `## Domain Model` 或 `### 聚合與不變條件`（舊格式，聚合寫在 SPEC）→ 當成〈設計〉
   `## Domain Model` 的素材讀，照本 skill 的規則重新判，不照抄；不改 `spec.md`。切票之前停下（Step 3 第 2 步）。
3. CHECK 同目錄是 2026-11-16 以前的 `tasks/` 目錄 → 記下，切票之前照〈再跑〉先併成一份 `tasks.md`（Step 3）。

**一次讀完整批，不分批**——契約、schema、聚合、接合點、風險與依賴天生跨 story，只有並排才看得見衝突；分批切會各自長出一套順序，實作時才發現互相卡住。

**做完的樣子**：上表每一列讀過；沒有的檔（第一個 feature 沒有 `docs/`、專案沒有 `api/openapi.yaml`）知道它沒有。→ 〈設計〉

### 設計

產出：使用者確認的 seam；`features/*.feature`（經 `formulation`）；`plan.md` 的 `## Domain Model`、`## Database Schema`、`## API 契約`、
技術決定七節；`## 缺口` 與 `## Further Notes` 的列

addyosmani 的 Step 2 之前插這一步——依賴圖的節點是模組與契約，契約與 schema 還沒定，圖就只能猜。

#### 決定 seam —— 驗收測試打在哪一層

seam 是驗收測試**打進系統的哪一層**（HTTP handler、application service、domain），會長成 step definition 的形狀——先定它，因為
`## Testing Decisions` 照它寫、場景對著它驗。`## API 契約` 寫 OpenAPI 片段還是程式碼層契約，看有沒有 HTTP 對外介面，不看 seam。

| 規則 | 為什麼 |
| --- | --- |
| 優先用既有 seam | 每多一個 seam 就多一套測試替身與資料建構 |
| 取**最高**的那一層 | 越高越接近使用者真正做的事，也越不綁實作細節 |
| 整個變更的理想數量是**一個** | 兩個 seam 代表這批行為的入口不只一個，多半是切分沒切乾淨 |
| 需要新 seam 就在**你能做到的最高點**提出 | 往下挪一層通常是為了好寫，代價是驗到的東西離使用者更遠 |

1. ASK 把 seam 攤出來確認再寫（只重切時跳過，見〈再跑〉）。三種情況：

   | `## Open Questions` 裡 `面向` 為 `seam` 的題 | 怎麼做 |
   | --- | --- |
   | `已答` | 只是確認，行尾標 `← Q-<n>` |
   | `n/a`（答案欄寫了為什麼不適用） | 照那個理由選最外層能驗到行為的一層，行尾標 `← Q-<n>` |
   | 沒有這題，或還 `待答` | 提出你的判斷給使用者確認：他確認或改了，標 `← 使用者 <日期>`；沒得到回答才標 `← 推論`，`## Further Notes` 的「推論的技術決定」記一條「seam 沒問過」 |

MUST NOT: 把 seam 的討論變成需求訪談。發現需求有洞寫進 `## 缺口` 交回 `spec`，**不要順手問掉**——這一步沒有寫
`## Open Questions` 的權限，順口問到的答案只活在這一次對話裡。

IMPORTANT: seam 一旦寫進 `plan.md` 就往下傳——IMPLEMENT 照著打，REVIEW 把「沒人同意過的 seam」當成 finding。

#### `.feature` —— 先寫行為，再設計

2. DELEGATE `../formulation/SKILL.md`（只重切時跳過，見〈再跑〉）：照它把 `spec.md` 寫成 `features/*.feature`、跑
   `check_spec.py`、產覆蓋表。交給它的只有 `spec.md` 與它自己要讀的檔——**不給**剛定的 seam 以外的任何技術決定。
   它回報的覆蓋缺口（`spec.md` 漏了例子、規則與例子打架）WRITE 進 `## 缺口`，交回 `spec`。

排在三份設計之前——`.feature` 只寫業務行為，設計怎麼定都不用改它；而 context 裡還沒有 operation 與表名，場景才不會寫成「Then 回應 201」「寫入 sessions 表」（綁死實作、需求方讀不懂）；
綁設計的是 step definition，`implement` 才寫。設計反過來讀 `.feature`：每個 operation 要指得出它服務哪幾個場景。

#### 三份設計

三份的順序固定：Domain Model → Database Schema → API 契約，都排在 `.feature` 之後。格式與範例都在 `references/design-format.md`。

3. WRITE `## Domain Model`：這次新增或改動的聚合、不變條件、跨聚合的規則。`被誰保證` 只寫聚合；動到
   `docs/DOMAIN-MODEL.md` 既有的聚合，寫「原本是什麼、為什麼改」。
4. WRITE `## Database Schema`：表、欄位、索引、約束；每張表寫它屬於哪個聚合，每個欄位寫 `守住哪條不變條件`。
5. DELEGATE `../api-and-interface-design/SKILL.md`（只重切時跳過，見〈再跑〉）：交給它這次的 FR／AC、第 2 步寫的 `.feature`、`api/openapi.yaml`
   與剛寫的 Domain Model；它交回 OpenAPI 片段（有 HTTP 對外介面時）或程式碼層介面契約（沒有時）。WRITE 進
   `## API 契約`，每個 operation 寫出它服務哪幾條 AC。
6. CHECK 三份並排對一次（判準表見 `references/design-format.md`〈三份並排對照〉）：每條不變條件有一個聚合、資料庫
   守得住的有一個約束；每張表屬於一個聚合；每個 operation 最多寫一個聚合，除非一條跨聚合的規則說要。
   對不上 → 先判斷錯在哪：錯在某一份設計 → 這一步就找出是哪一份、改掉它，對不上的地方與怎麼改的記進
   `## Further Notes` 的「技術決定的修訂」；錯在需求（規則沒說、或兩條 FR 打架）→ `## 缺口` 交回 `spec`。
   不要不判斷就挑一份去遷就另一份。
7. CHECK 每個 operation 指得出 AC。指不出 → 刪掉那個 operation；AC 要求的行為沒有 operation → 補上，或是需求的洞就
   進 `## 缺口`。

MUST NOT: 在這一步改 `api/openapi.yaml` 或寫 migration。`plan.md` 是契約與 schema 的來源；`implement` 做票時
才照抄進專案。

#### 技術決定

8. READ `references/tech-format.md`。
9. WRITE 七節：`## Tech Stack`、`## Commands`、`## Project Structure`、`## Code Style`、`## Boundaries`、
   `## Implementation Decisions`、`## Testing Decisions`。前五節（技術現況）照舊完整寫出——`plan.md` 自成一體；
   `## Implementation Decisions` 寫動哪些模組、程式碼層的介面與埠、架構決定、互動，**不放** API 契約與 schema；
   `## Testing Decisions` 寫 seam 與每個 AC 在哪一層測。每條行尾標來源（五選一）；可以用 `推論`，但每一條同時
   列進 `## Further Notes` 的「推論的技術決定」。

#### `.feature` 在 seam 上驗得到嗎

10. CHECK 逐個場景問：它的 Then 在剛定的 seam 上**看得到**嗎？例如 seam 是 HTTP API，「借書單印出到期日」
   「櫃台刷書機嗶一聲」是設備端的行為，HTTP 回應裡看不到。

   | 結果 | 怎麼做 |
   | --- | --- |
   | 看得到 | 照常，`## Testing Decisions` 的測試層級表標它在哪一層 |
   | 看不到，但換一個說法看得到（「API 回應帶這次借閱的到期日」） | 再 DELEGATE `formulation` 改寫那個場景的 Then：交給它那個場景與「seam 上看不到什麼」，措辭由它定；改完回這一步再對一次 |
   | 怎麼換都看不到 | `## 缺口` 一列，交回 `spec`：這條 AC 驗不了，要需求方決定是換一個驗法還是拿掉 |

   MUST NOT: 不經 `formulation` 自己改 `.feature` 的場景，或把驗不到的場景默默標成手動測試。

**做完的樣子**：seam 使用者確認過（或標了 `推論` 並記進 Further Notes）；`.feature` 由 `formulation` 寫好、覆蓋缺口進了 `## 缺口`；三份設計都寫了、並排對過，設計這邊對不上的
當場改了並記進 Further Notes、需求那邊的進了 `## 缺口`；技術決定七節都在；每個場景對照過 seam。→ Step 2

### Step 2 依賴圖

產出：`plan.md` 的 `## 依賴圖`、`## Risks`

1. WRITE `## 依賴圖`（格式 `references/design-format.md`〈依賴圖〉）：節點是模組與契約，不是票——哪個先存在，別的
   才寫得出來。建構順序照圖由下而上：地基先做。
2. WRITE `## Risks`（格式 `references/plan-format.md`〈Risks〉）：規格沒說、但實作一定會撞到的。每條指出從哪條規則
   長出來、標影響程度、寫它影響圖上哪些節點。

**做完的樣子**：圖上每個節點是模組或契約；每條風險指得出規則、標了影響程度。→ Step 3

### Step 3 垂直切片

產出：票的清單草稿（還不寫 `tasks.md`）；兩道關沒過時是 `## Overview` 與回報，不切票；圖錯了時是 `## Further Notes` 的修訂

先過兩道關，任何一道不過就 WRITE `## Overview`（〈WRITE〉第 2 步）、照〈完成後〉回報、STOP，**不切票**：

1. CHECK `## 缺口` 有要回 `spec` 的列 → 告訴使用者先回 `/spec`。
   設計還會變時切的票，`**契約**`、`**資料表**`、`**驗證**` 抄的都是會變的東西。
2. CHECK Step 1 記下 `spec.md` 還有舊的聚合一節 → 請使用者跑 `/spec`（`plan.md` 已有 `## Domain Model`，它只刪那一節），
   再叫 `/plan` 接著切票。舊節留著，兩份聚合各說各的，票的 `**契約**`／`**資料表**` 就只能自己猜。

Step 1 記下舊的 `tasks/` 目錄 → 先照〈再跑〉併成一份 `tasks.md`。然後 READ `references/tasks-format.md`，切：

3. 切出**最薄的完整路徑**當第一張票：碰到每一層、邏輯最少。它證明的是接線，不是行為。專案已經有骨架就跳過。
4. 其餘依 `spec.md` 的 FR 切，一張票一組相關的 `@example-N.M`。每張票守三條規則（〈三條切票規則〉）：
   **垂直**——穿過所有層；**自己可 demo**——做完當下就能展示；**塞得進一個全新 session**。每張票問一次「做完之後
   我能 demo 什麼」，答不出來就是水平切片，重切。
5. 現有結構容不下這批行為時，切一張 `00` 前置整理排在最前：「make the change easy, then make the easy change.」
   兩張以上的票都需要一個還不存在的 package、共用的測試輔助函式或同一張還不存在的表時，指定一張票擁有它
   （通常是 `00`），其他票 `Blocked by` 它——各自新建，合回時會撞成重複宣告。
6. CHECK 這批有沒有一個機械性改動（改欄位名、換共用型別）影響面擴散到全 codebase、沒有任何垂直切片能綠 → 改用
   expand–contract：expand 一張、migrate 依 blast radius 每批一張、contract 一張（〈wide refactor 的例外〉）。
7. CHECK 照 `## 依賴圖` 切不出可 demo 的票（圖要求的先後把一條行為拆成兩半）→ 圖錯了：回 Step 2 改圖，記進
   `## Further Notes` 的「技術決定的修訂」。

**做完的樣子**：每張票答得出能 demo 什麼；需要前置整理的有 `00`；wide refactor 照 expand–contract 排了。→ Step 4

### Step 4 寫任務

產出：每張票每一欄的草稿（還不寫 `tasks.md`）；需求的洞進 `## 缺口`

每張票照 `references/tasks-format.md` 的票節擬出每一欄（這一步還不寫檔，〈ASK 切法〉點頭後才寫）：

| 欄位 | 從哪來 |
| --- | --- |
| `交付什麼行為` | 使用者角度做完之後什麼會動，不是分層清單 |
| `涵蓋場景` | 這張票的 `@example-N.M` |
| `Blocked by` | 這張票的驗收條件，在那張票做完之前有沒有可能綠 |
| `接合點` | 這張票會改、別張票也可能改的檔：**路徑**＋種類（`產生`／`追加`／`爭用`／`新建`），從 `## Implementation Decisions`、`## Database Schema`、`## API 契約` 推（〈接合點與同時跑〉） |
| `模組` | 照抄 `## Implementation Decisions`「動哪些模組」裡這張票碰的那幾個，不寫檔案路徑 |
| `契約` | 這張票涵蓋的 `## API 契約` operation（`operationId`；程式碼層契約寫埠的方法），沒有就「無」 |
| `資料表` | 這張票要建或改的 `## Database Schema` 表，沒有就「無」 |
| `大小` | XS／S／M／L／XL（〈大小〉） |
| `內迴圈` | `## Testing Decisions` 的測試分層 |
| 驗收條件 | 每條在起點 commit 上是紅的（〈驗收條件要能失敗〉）；講業務規則的那條行尾標 `← AC-N.M`，只有 FR 文字、沒有 AC 例子的標 `← FR-N` |
| `**驗證**` | 場景命令（tag 運算式由 `涵蓋場景` 組成）與閘門命令，都從 `## Commands` 抄（〈驗證〉） |

1. CHECK `## API 契約` 是 OpenAPI 片段、而這張票有 `**契約**` → 驗收條件照抄
   `- [ ] 搬進 api/openapi.yaml 的 operation 與 plan.md 的片段一致`；程式碼層契約（沒有 HTTP 對外介面）不加這一條。
2. CHECK 大小。MUST NOT: 交出 L 或 XL 的票——回 Step 3 再拆。agent 在 S 與 M 做得最好；L 以上的訊號（一個全新
   session 做不完、驗收條件三條寫不完、碰到兩個互不相干的子系統、標題裡有「和」）見〈大小〉。
3. CHECK 漏列的接合點：會長大的介面連同它的測試 stub、依賴檔、migration 目錄與序號這幾類最常漏（清單見〈接合點與
   同時跑〉）。`追加` 與 `爭用` 拿不準時標 `爭用`。
4. CHECK 講業務規則的驗收條件：標得出 `← AC-N.M`，或用 `← FR-N` 那條 FR 自己的字詞、不多做一個 FR 沒做的選擇？
   標不出、或要多選一個（「跨 00:00 的補繳歸到付款後那天」而 FR 只說「當天」）→ 那是需求的洞：不寫進票，
   寫進 `## 缺口` 交回 `spec`。

IMPORTANT: 接合點要寫**路徑**，不是模組名——寫「HTTP 層」，合回的人只能自己編路徑，試跑中漏掉的正是合回時出事的那幾個檔。

判準：指著票裡的任何一句話，說得出它來自 `spec.md` 的哪條 FR、`plan.md` 的哪一個技術決定、或哪個 `@example-N.M`
嗎？說不出來就是這一步在發明東西。

**做完的樣子**：每張票每一欄都有（沒有就「無」），沒有 L 與 XL，每條驗收條件在起點是紅的。→ Step 5

### Step 5 排序與 Checkpoint

產出：票的順序與檢查點（草稿）；`plan.md` 的 `## Parallel Work`

1. 排順序：blockers 在前；一張票不得排在它依賴的 `## 依賴圖` 節點之前（票碰的模組、operation 或表，在圖上的父節點
   要由同一張或 `Blocked by` 指得到的票做出來）。`## Risks` 影響程度「高」的票，在 `Blocked by` 允許的範圍內往前排——
   越早撞到越便宜。每張票做完，系統都還是好的。
2. ASK 這批票要不要同時做（只重切時跳過，見〈再跑〉）。答案記進 `plan.md` 的 `## Parallel Work`，行尾 `← 使用者 <日期>`：

   | 回答 | 寫什麼 |
   | --- | --- |
   | **不要**（預設） | 只寫那一行；`tasks.md` 的 `**可以同時開始**` 寫「不適用」，`接合點` 照樣列——合回時怎麼解，順序做的票一樣用得到 |
   | **要** | 接合點表：每張票 `接合點` 的聯集，路徑、種類、為什麼會共用（`references/plan-format.md`〈Parallel Work〉）。**先定契約，再並行**：共用同一份契約的票等 `## API 契約` 定了才同時做 |

3. 分波：答「要」時，`Blocked by` 都完成、沒有共同的 `爭用` 接合點、沒有兩張都要新建的東西的票同一波；答「不要」時
   每 2–3 張票一波。每一波之後放一個檢查點，不並行時每 2–3 張票一個（〈檢查點〉）——停下來驗系統還是好的，再往下；命令照
   `## Commands`。

**做完的樣子**：沒有一張票排在它依賴的節點之前；高風險的票在允許的範圍內排在前面；`## Parallel Work` 記了使用者的回答；
每一波（或每 2–3 張票）之後有檢查點。→ ASK 切法

### ASK 切法

產出：使用者對切法點頭（或要改的地方，回 Step 3／4／5）

MUST: 寫 `tasks.md` 之前先把票的清單攤給使用者看——編號、標題、大小、Blocked by、接合點、交付什麼行為、哪幾張可以同時
開始、檢查點落在哪（再跑時只攤重切的票與新票，開工了的票不再問）——並問三件事：

1. 粒度對不對（太粗／太細）
2. Blocked by 是不是真的卡住，還是只是「先做比較順」
3. 有沒有該合併或該再拆的

使用者說要改 → 回 Step 3／4／5 重切，再攤一次、再問一次。**點頭之前不寫 `tasks.md`。**

IMPORTANT: **切太細是這一步最常見的失效**，而且不會自己被發現——一疊原子化的小票看起來很專業，卻丟掉了「這幾件事其實是同一件」；
為了並行把一件事切成兩張 `爭用` 同一段邏輯的票，換到的是零：它們照樣不能同時做。

### WRITE

產出：`tasks.md`、`plan.md` 的 `## Overview`、給使用者的回報

1. WRITE `tasks.md`，照 `references/tasks-format.md`：開頭三行，然後一波一節、依賴序寫票節，每張票的 `#### 缺口`
   寫「（無）」；波與波之間放檢查點那一行。`tasks.md` 已存在時照〈再跑〉就地改，不整份重寫：開工了的票只動
   `接合點`／`Blocked by`／`大小`，新票接在最大編號之後；`#### 缺口` 寫「（無）」只用於新票。
2. WRITE `plan.md` 的 `## Overview`：上面都寫完才寫，一段摘要，點名三份設計各新增了什麼、切成幾張票，不做新決定。
3. STOP 照〈完成後〉回報，告訴使用者下一步是 `/implement <第一張票的編號>`。

**做完的樣子**：`tasks.md` 照骨架寫好、使用者點過頭；`plan.md` 每一節都在；回報後停下。

---

## 再跑：重切，不重編

addyosmani 的「不覆寫還有沒勾的計畫」：`tasks.md` 有人照著開工了，整份重寫等於丟掉一份只存在那裡的進度。

`plan.md` 已經存在時照新的 `spec.md` 改，不另開一份——〈設計〉第 2 步 DELEGATE `formulation` 照它的〈再跑〉改 `.feature`；票切出去之後改了 seam、介面、schema、聚合、API 契約，
記進 `## Further Notes` 的「技術決定的修訂」。

`tasks.md` 已經存在時（還是 2026-11-16 以前的 `tasks/` 目錄時，先併成一份：照 `index.md` 的順序把每張票檔的內容放進
對應的票節、標題 `# 01 — …` 改 `### 01 — …`、`## 缺口` 改 `#### 缺口`，其餘一字不改，`git rm -r tasks/`，單獨一個 commit）：

- 驗收條件已經有勾的票**不動**——有人做過了，改它等於改掉一份已經驗過的約定。判斷看**主線上**的票節；正在 worktree
  裡做的票（`.feature` 是 `@wip`、主線上的票節還沒勾）也算開工了
- 票節的欄位各有主人：`接合點`、`Blocked by`、`大小` 歸本 skill，開工了的票也只准改這三欄（多半是補漏列的接合點）；
  驗收條件的勾、`**驗證**` 的勾、`**檢查點**:`、`#### 缺口` 歸 `implement`，本 skill 不動（唯二例外是下兩條在列尾加的
  標記）——兩邊改的是不同行，合回時才不會互相蓋掉
- 還沒開始的票照新設計重切：可以改、合併、再拆；作廢的票留著、標題後加「（作廢：<原因>）」，編號不重用，新票接在最大
  編號之後。重切的票一樣過〈ASK 切法〉
- 票節 `#### 缺口` 裡 `implement` 記下的洞：屬於設計的，這一輪在設計裡補上、記進 `## Further Notes` 的「技術決定的修訂」；
  屬於需求的，進 `plan.md` 的 `## 缺口` 交回 `spec`。兩種都在票節那一條後面標「→ 已記入 plan.md」
- 只重切：`spec.md` 與 `.feature` 在 `plan.md` 的日期之後都沒改，新的輸入只有 `tasks.md` 的 `#### 缺口` 列 → 跳過〈設計〉的
  seam ASK、`formulation` 與 API 兩個 DELEGATE、Step 5 的並行 ASK（記下的回答照舊成立），屬於設計的列當場補上，從 Step 3 開始
- 喚醒或修補之後：舊票 `#### 缺口` 裡已經被 `spec.md` `## Open Questions` 新題回答的列（判準：那一題的答案正是那一列在等的），
  在列尾標「→ 已回答 Q-<n>」，交給新票做；不標，下一個讀缺口的人會再猜一次

## 不改規格：有洞怎麼辦

MUST NOT: 在這一步新增、刪除或改寫任何 FR、AC、story、範圍；`.feature` 只由 `formulation` 照 `spec.md` 寫，本 skill
不自己動任何場景。需求已經在 `spec.md` 裡、需求方說過「就緒」；這裡再決定一次，只會產生一份跟規格不一樣的規格。

判準：指著 `plan.md` 或 `tasks.md` 裡的任何一句話，說得出它是為了 `spec.md` 的哪條 FR 或 NFR 嗎？說不出來就是這一步
僭越了。

寫計畫時才發現的，依種類處理：

| 種類 | 例子 | 怎麼處理 |
| --- | --- | --- |
| **技術選擇** | 哪個 library、欄位叫什麼 | `## Further Notes` 的「留給 IMPLEMENT 的技術選擇」 |
| **執行上的選擇** | 某個檢查要不要加 | 留給 IMPLEMENT 決定 |
| **場景驗不到** | 〈設計〉第 10 步找到的 | 換說法看得到 → 再 DELEGATE `formulation`；怎麼換都看不到 → `## 缺口` 交回 `spec`（換驗法） |
| **`.feature` 的覆蓋缺口** | 〈設計〉第 2 步 `formulation` 回報的：`spec.md` 漏了例子、規則與例子打架 | `## 缺口`，交回 `spec` |
| **三份設計對不上、依賴圖切不出票** | 表屬於兩個聚合、不變條件沒人守、operation 改兩個聚合；照圖排會把一條行為拆成兩半 | 錯在設計 → 不進 `## 缺口`：這一步當場改掉錯的那一份，記進 `## Further Notes` 的「技術決定的修訂」；錯在需求 → `## 缺口` 交回 `spec` |
| **需求的洞** | 某條規則在某個情況下該怎麼做，沒人說過；兩條 FR 對同一個介面假設不一樣 | `## 缺口`，告訴使用者跑 `spec`；它開成補問題、交給 `discovery` 問 |

## 常見的藉口與紅旗

| 藉口 | 實際上 |
| --- | --- |
| 「邊做邊想就好」 | 那就是返工與一團亂的來源；切票時十分鐘，省下實作時好幾個小時 |
| 「票很明顯，不用寫」 | 照樣寫下來。寫出來的票才看得見藏著的依賴與被忘掉的邊界情況 |
| 「我記得住」 | context 有限。寫下來的計畫活得過 session 邊界與 compaction |
| 「舊的 `tasks.md` 過期了，整份換掉」 | 沒勾的票可能正在另一個 worktree 裡做；照〈再跑〉重切，不重編 |

看到這些就停下來修：票說「實作這個功能」卻沒有驗收條件；票沒有 `**驗證**`；每張都是 L 或 XL；票與票之間沒有檢查點；
沒照依賴圖排順序；沒問過就改了有勾的票；規劃到一半開始寫程式。

---

## 產物

| 檔 | 哪一步寫 | 格式在哪 |
| --- | --- | --- |
| `specs/ongoing/<date>-<feature>/plan.md` | 設計、Step 2–5、WRITE（`## Overview` 最後寫；Step 3 兩道關沒過時提早寫） | `references/plan-format.md`（整份節序）；技術決定七節 `references/tech-format.md`；三份設計與 `## 依賴圖` `references/design-format.md` |
| `specs/ongoing/<date>-<feature>/tasks.md` | WRITE（〈ASK 切法〉點頭後） | `references/tasks-format.md`：開頭三行（來源、哪幾張可以同時開始）、一波一節、一張票一個票節（`00` 前置整理有才寫） |
| `features/*.feature` | 設計第 2、10 步，由 DELEGATE 的 `formulation` 寫 | `../formulation/SKILL.md` |

MUST: 只寫入 `specs/ongoing/<date>-<feature>/plan.md` 與 `tasks.md`，**不修改專案既有的任何檔案**，包括 `spec.md` 與 `api/openapi.yaml`。例外兩個：
`features/*.feature` 由〈設計〉DELEGATE 的 `formulation` 照它自己的規則寫；併舊的 `tasks/` 目錄時 `git rm -r tasks/`，單獨一個 commit（〈再跑〉）。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] `spec.md`、`api/openapi.yaml` 與任何程式碼沒有新增、刪除或改寫；`.feature` 只經 `formulation` 寫或改
- [ ] `.feature` 在三份設計之前寫好，場景裡沒有 endpoint、表名、欄位；`formulation` 回報的覆蓋缺口進了 `## 缺口`
- [ ] seam 攤給使用者確認過，回覆記在回報裡；談 seam 時發現的需求洞進了 `## 缺口`，沒有順口問掉
- [ ] 三份設計並排對過；錯在設計的當場改了、記進 `## Further Notes`；錯在需求的進了 `## 缺口` 交回 `spec`；API 契約每個 operation 指得出 AC
- [ ] 技術決定七節都在（空的寫了為什麼空），`## Implementation Decisions` 裡沒有 API 契約與 schema；每條行尾有來源，每條 `← 推論` 同時列在 `## Further Notes`
- [ ] 每個場景都對照過 seam；有要回 `spec` 的 `## 缺口` 時沒有切票
- [ ] `## 依賴圖` 的節點是模組與契約；`## Risks` 每條指得出從哪條規則長出來、標了影響程度
- [ ] 每張票答得出能 demo 什麼；每張票每一欄都在（沒有就「無」），模組照抄、沒有檔案路徑，接合點是路徑加種類；沒有 L 與 XL
- [ ] `## API 契約` 是 OpenAPI 片段時，有契約的票驗收條件裡有「搬進 api/openapi.yaml 的 operation 與 plan.md 的片段一致」；程式碼層契約的票沒有這一條
- [ ] 每張票的 `**驗證**` 有閘門命令；有 `涵蓋場景` 的票另有場景命令，tag 運算式跟 `涵蓋場景` 一一對應
- [ ] 指著票裡的任何一句話，說得出它來自哪條 FR、`plan.md` 的哪一個技術決定、或哪個 `@example-N.M`
- [ ] `## Risks` 影響程度「高」的票，在 `Blocked by` 允許的範圍內排在前面
- [ ] `**可以同時開始**` 的票之間沒有共同的 `爭用` 接合點，也沒有兩張都要新建的東西
- [ ] 沒有一張票排在它依賴的 `## 依賴圖` 節點之前；每一波（不並行時每 2–3 張票）之後有檢查點；`## Parallel Work` 記了使用者的回答
- [ ] 寫 `tasks.md` 之前把切法攤給使用者、問過粒度／Blocked by／合併或再拆三件事，使用者點了頭
- [ ] 再跑時，驗收條件已有勾的票只動了 `接合點`／`Blocked by`／`大小`，編號沒有重排或重用

## 完成後

告訴對方：

1. **seam 是哪一層**、是確認 `## Open Questions` 裡問過的還是推的、使用者確認了沒
2. **三份設計**：幾個聚合（新增／改動）、幾張表、幾個 operation，`## API 契約` 是 OpenAPI 片段還是程式碼層契約；並排時對不上幾處、其中幾處當場改了設計、幾處交回 `spec`
3. `formulation` 寫了幾個 `.feature`、幾個場景、覆蓋表有沒有缺口；動哪些模組、幾個決定標了 `推論`；幾個場景 seam 驗不到、幾個請 `formulation` 換了說法、幾個交回 `spec`
4. 幾條風險、幾條影響程度「高」
5. 幾張票、第一張是哪一張（以及**為什麼是它**）、有沒有前置整理；幾個檢查點
6. 要不要並行；要的話哪幾張可以同時開始、哪幾張邏輯上獨立卻因為 `爭用` 或新建的東西排成先後
7. `## 缺口` 有幾條，其中幾條要回 `spec`

`## 缺口` 有要回 `spec` 的 → 先回 `/spec`（這一輪沒有切票）。沒有 → 下一步 `/implement <票號>`，一張票一個
全新 session；同時做的票各開一個 worktree。
