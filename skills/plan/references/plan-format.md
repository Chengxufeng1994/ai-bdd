# `plan.md` 的骨架

照抄即可。讀這份的時機：`plan` 第 2 步開始寫 `plan.md` 之前，以及第 4 步寫執行計畫之前。

一個 feature 一份，路徑 `specs/ongoing/<date>-<feature>/plan.md`，跟同目錄的 `spec.md` 配對——目錄名已經帶 feature 的識別資訊，檔頭不重複列。

它裝四件事：`## 系統介面`——這次碰到哪些介面、照什麼順序分析；**技術決定**——用什麼做、seam、模組、
測試層級，七節的寫法見 `tech-format.md`；**三份設計**——`## Domain Model`、`## Database Schema`、
`## API 契約`，寫法見 `design-format.md`；**執行計畫**——依賴、哪裡危險、順序與 checkpoint、能不能並行，
以及 `## 缺口`。要做成什麼在 `spec.md`，這裡不重寫。票不在這裡——`task-breakdown` 從這份計畫拆到同目錄的
`tasks.md`，格式見 `../../task-breakdown/references/tasks-format.md`。

執行計畫是拋棄式快照，實作一開始就會過期；技術決定與三份設計不是——它們改了要在 `## Further Notes` 留紀錄，
做完時 `archive` 從它們更新 `docs/TECH.md`、`docs/DOMAIN-MODEL.md`、`docs/DATABASE-SCHEMA.md`。

```markdown
# 執行計畫

**來源**: `spec.md`（已對焦，<日期>）＋ `features/` 底下標 `@ready` 的 <N> 個檔
**日期**: <YYYY-MM-DD>

## Overview

## 系統介面          → design-format.md

## Tech Stack        → tech-format.md（這七節）
## Commands
## Project Structure
## Code Style
## Boundaries
## Implementation Decisions
## Testing Decisions

## Domain Model      → design-format.md（這三節）
## Database Schema
## API 契約

## 依賴圖            → design-format.md
## Risks
## 順序與 Checkpoint
## Parallel Work

## 缺口

這份計畫自己的洞、三份設計並排時對不上的，以及 `task-breakdown` 拆票時發現技術決定沒定案的。
沒有就寫「沒有——`spec.md` 與這份計畫涵蓋了這批需要的一切」。

- <缺口>：<需求的洞（回 `spec`，它開成補問題交給 `discovery`）、技術決定或設計的洞（回 `plan` 補）、`.feature` 寫了 seam
  驗不到的行為（回 `formulation`），或執行上的選擇（留給 IMPLEMENT 決定）>

## Further Notes

- **推論的技術決定**：技術決定每條標 `← 推論` 的再列一次，一行一條指回它在哪一節
- **留給 IMPLEMENT 的技術選擇**：不影響任何 AC 的選擇（哪個 library、欄位叫什麼）
- **技術決定的修訂**：票切出去之後改了 seam、介面、schema、聚合、API 契約的，一行一條：日期、改了什麼、為什麼
```

骨架裡的 `→` 只是告訴你去哪讀寫法，不要抄進 `plan.md`。

---

## Overview

一段話，三到五句：這批 story 做完之後使用者能做什麼，加上技術上大致怎麼做（seam 在哪一層、
新增或改動的主要模組、三份設計各新增了什麼——聚合、表、operation，照下面各節）。給 `task-breakdown` 與人第一眼看懂這份計畫在做什麼。

只摘要 `spec.md` 與下面各節已經寫的東西，**不在這裡做新決定**，所以不標來源。寫著寫著發現要
講一件下面各節沒有的事，那件事屬於某一節技術決定——寫進那一節，不是這一段。

## Risks

規格沒說、但實作一定會撞到的。每一條指出從哪條規則長出來。
跟 `## 缺口` 不同：缺口是**還不知道要做什麼**，風險是
**知道要做什麼但做對很難**。

| 風險 | 影響程度 | 從哪條規則長出來 | 影響什麼 |
| --- | --- | --- | --- |
| 併發：兩個請求同時通過「最多一筆進行中」 | 高 | log-a-workout FR-1 | schema（唯一索引）、實作順序 |

`影響程度` 看做錯的代價，不看發生機率：

| 影響程度 | 做錯會怎樣 |
| --- | --- |
| 高 | 某條 AC 不成立，或資料錯了救不回來 |
| 中 | 要返工，但不會壞資料 |
| 低 | 只影響效能、維運或日後擴充 |

`task-breakdown` 會把碰到「高」的票排在前面：越早撞到，越便宜。

SHOULD NOT: 把解法寫死。「用 Redis 還是 DB 行鎖」需要知道實際流量與既有基礎
設施，而那些不在規格裡。寫「兩種都可以，看部署形態」比寫「用 Redis」有用——
後者看起來是決定，實際上是猜的。

## 順序與 Checkpoint

照 `## 依賴圖` 排：被依賴的節點先做，碰到「高」風險的節點盡量往前。每兩三個節點之後放一個
checkpoint——停下來驗系統還是好的，再往下。

```markdown
## 順序與 Checkpoint

1. 借閱聚合（到期日、續借一次）與讀者證系統的查詢埠——兩者都沒有依賴；兩則 story 都用到到期日的算法，最先做
2. `loans` 表與 repository
3. `## API 契約` 的 `lendBook`、`renewLoan` 與 HTTP handler

**Checkpoint：1–2 之後**
- [ ] `## Commands` 的閘門綠
- [ ] AC-1.1–1.3 的內迴圈測試綠

**Checkpoint：3 之後**
- [ ] `## Commands` 的閘門綠
- [ ] 涵蓋的 `.feature` 全部場景綠
- [ ] 給使用者看過再往下
```

這裡排的是**節點**，不是票：`task-breakdown` 切票時照這個先後與 checkpoint 排波次，一張票不得排在它依賴的節點
之前。checkpoint 的命令照 `## Commands` 的 `閘門` 與 `跑場景` 寫，不另編。

## Parallel Work

這批票要不要同時做，以及同時做時哪些檔會在合回時相撞。`task-breakdown` 照這一節標每張票的
`接合點`，不自己編路徑。

**先定契約，再並行**：共用同一份 API 契約的票，要等 `## API 契約` 定好才能同時做——契約還在動時並行，每張票
都照自己猜的形狀寫。

```markdown
## Parallel Work

**要不要並行**: 要 ← 使用者 2026-11-15

| 接合點 | 種類 | 為什麼會共用 |
| --- | --- | --- |
| `api/openapi.yaml` | 追加 | 每個 operation 都加在這份契約 |
| `internal/interfaces/http/apigen/api.gen.go` | 產生 | 由契約產生，內嵌整份 spec |
| `internal/application/port/in/lending_service.go` ＋ 實作它的測試 stub | 追加 | 每個 use case 在同一個驅動埠加方法 |
| `go.mod`／`go.sum` | 產生 | 第一個 path parameter 會帶進 runtime 套件 |
| `internal/application/usecase/command/`（新建） | 新建 | 至少兩張票的 use case 住這裡，交給一張前置票建 |
```

- **要不要並行**是 `Ask user`，不是推論。預設不要——並行換到的是時間，代價是每張票都要
  在合回時解衝突。選「不要」時只寫這一行，表格省略
- 種類三選一加一個 `新建`，定義與合回時怎麼解 →
  `../../task-breakdown/references/tasks-format.md`〈接合點與同時跑〉。
- 從上面的 `## Implementation Decisions` 推，最常漏的有六類（清單與例子 →
  `tasks-format.md`〈接合點與同時跑〉）：會長大的介面連同它的測試 stub、第一次用到的 library 帶進來的
  依賴檔、兩個以上 use case 都要住的新 package、新的必填組裝欄位連同每個組它的測試、按 feature 分的
  mapper／presenter 與共用的步驟檔、跨 feature 共用的步驟句子
- 兩張票會改同一段既有邏輯的地方標 `爭用`，寫明是哪一段；`task-breakdown` 會把它們排成先後

## 為什麼一個 feature 一份，不是一則 story 一份

試過「每則 story 一份 ＋ 一份共用檔」，壞成兩種：共用檔把實質吸走（各檔的
schema 一節全退化成「見共用檔」），而且兩種檔案各用一套節號、意思只對得上一半
——**部分對齊比完全不對齊更危險**，讀的人會假設對齊然後在三個節上讀錯。

接合點與風險天生跨 story，只有並排才看得見衝突。`spec.md` 的 story 再多，`plan.md` 仍然一個 feature 一份——story
是切分單位，不是檔案邊界。

## 範例

````markdown
# 執行計畫

**來源**: `spec.md`（已對焦，2026-09-14）＋ `features/` 底下標 `@ready` 的 2 個檔
**日期**: 2026-09-14

## Overview

館員借出一本書時算出到期日，遇休館日順延；讀者可在到期前上網續借一次，名下有逾期的書時拒絕。驗收測試打在
`application.借閱` 這一層，新增借閱的 domain 模組與借書、續借 use case，HTTP 只做轉接。三份設計：一個新聚合
「借閱」加一條跨聚合規則（逾期擋續借）、一張 `loans` 表、`lendBook` 與 `renewLoan` 兩個 operation。

（`## 系統介面` 與三份設計略，寫法與範例見 `design-format.md`；技術決定七節略，見 `tech-format.md`〈範例：技術決定〉）

## Risks

| 風險 | 影響程度 | 從哪條規則長出來 | 影響什麼 |
| --- | --- | --- | --- |
| 一本逾期就擋掉全部續借，讀者有正當理由時櫃台沒有例外通道 | 中 | renewal FR-3 | 要不要留館員手動續借的通道，留給 IMPLEMENT |

## 順序與 Checkpoint

（略，見上面〈順序與 Checkpoint〉的範例）

## Parallel Work

**要不要並行**: 不要 ← 使用者 2026-09-14

## 缺口

- FR-2 說每本書限續借一次，但沒有一條 AC 寫第二次續借會怎樣——`## API 契約` 定不出那個錯誤回應（三份並排時找到）：
  需求的洞，回 `spec`
````
