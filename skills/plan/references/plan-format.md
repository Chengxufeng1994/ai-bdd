# `plan.md` 的骨架

照抄即可。讀這份的時機：步驟 2，準備寫執行計畫的時候。

一個 feature 一份，路徑 `specs/ongoing/<date>-<feature>/plan.md`，跟同目錄的 `spec.md`、
`discovery.md` 配對——目錄名已經帶 feature 的識別資訊，檔頭不重複列。

它只裝**執行計畫**：這批怎麼排、哪裡危險、能不能並行，以及 `## 缺口`。技術設計（Tech Stack、
seam、模組、契約、schema、測試層級）在 `spec.md` 的技術部分，格式見
`../../spec/references/tech-format.md`。票不在這裡——`task-breakdown` 從 `spec.md` 的技術部分與這份
計畫拆到同目錄的 `tasks/`，格式見 `../../task-breakdown/references/tasks-format.md`。
`plan.md` 是拋棄式快照，實作一開始就會過期。

```markdown
# 執行計畫

**來源**: `spec.md`（已對焦、技術部分已寫，<日期>）＋ `features/` 底下標 `@ready` 的 <N> 個檔
**日期**: <YYYY-MM-DD>

## Overview
## Risks
## Parallel Work

## 缺口

執行計畫自己的洞，以及 `task-breakdown` 拆票時發現 `spec.md` 技術部分沒定案的。
沒有就寫「沒有——`spec.md` 與這份計畫涵蓋了這批需要的一切」。

- <缺口>：<技術部分的洞（回 `spec` 補），或執行上的選擇（留給 IMPLEMENT 決定）>
```

---

## Overview

一段話，三到五句：這批 story 做完之後使用者能做什麼，加上技術上大致怎麼做（seam 在哪一層、
新增或改動的主要模組，照 `spec.md` 技術部分）。給 `task-breakdown` 與人第一眼看懂這份計畫在做什麼。

只摘要 `spec.md` 與下面各節已經寫的東西，**不在這裡做新決定**，所以不標來源。寫著寫著發現要
講一件下面各節沒有的事，那件事屬於 `spec.md` 的技術部分——回 `spec` 補，不是這一段。

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

## Parallel Work

這批票要不要同時做，以及同時做時哪些檔會在合回時相撞。`task-breakdown` 照這一節標每張票的
`接合點`，不自己編路徑。

```markdown
## Parallel Work

**要不要並行**: 要 ← 使用者 2026-11-15

| 接合點 | 種類 | 為什麼會共用 |
| --- | --- | --- |
| `api/openapi.yaml` | 追加 | 每個 operation 都加在這份契約 |
| `internal/interfaces/http/apigen/api.gen.go` | 產生 | 由契約產生，內嵌整份 spec |
| `internal/application/port/in/parking_service.go` ＋ 實作它的測試 stub | 追加 | 每個 use case 在同一個驅動埠加方法 |
| `go.mod`／`go.sum` | 產生 | 第一個 path parameter 會帶進 runtime 套件 |
| `internal/application/usecase/command/`（新建） | 新建 | 至少兩張票的 use case 住這裡，交給一張前置票建 |
```

- **要不要並行**是 `Ask user`，不是推論。預設不要——並行換到的是時間，代價是每張票都要
  在合回時解衝突。選「不要」時只寫這一行，表格省略
- 種類三選一加一個 `新建`，定義與合回時怎麼解 →
  `../../task-breakdown/references/tasks-format.md`〈接合點與同時跑〉。
- 從 `spec.md` 的 `## Implementation Decisions` 推，最常漏的有六類（清單與例子 →
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

**來源**: `spec.md`（已對焦、技術部分已寫，2026-09-14）＋ `features/` 底下標 `@ready` 的 2 個檔
**日期**: 2026-09-14

## Overview

訪客出場時依停留時間算出應付金額、繳費後放行；月租車在有效期內不計費。驗收測試打在
`application.出場結帳` 這一層，新增計費的 domain 模組與出場結帳 use case，HTTP 只做轉接。

## Risks

| 風險 | 影響程度 | 從哪條規則長出來 | 影響什麼 |
| --- | --- | --- | --- |
| 警報期間出場一律不計費，警報頻繁或被誤觸發時收入流失 | 中 | monthly-pass FR-3 | 要不要記錄警報期間的每一次出場，留給 IMPLEMENT |

## Parallel Work

**要不要並行**: 不要 ← 使用者 2026-09-14

## 缺口

沒有——`spec.md` 與這份計畫涵蓋了這批需要的一切
````
