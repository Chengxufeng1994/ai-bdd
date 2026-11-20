# `plan` 照 addyosmani 重寫，併掉 `task-breakdown`

**狀態**：已批　**日期**：2026-11-20　**取代**：`2026-11-20-plan-as-commander-design.md` 的 D2（流程）與「`task-breakdown` 維持獨立」

## 意圖

需求方的話：「plan 感覺要重新寫了」「addyosmani 的 planning-and-task-breakdown，我覺得 follow 他的設計很不錯，
只是在 step 2 之前可能要設計 api（有需要的話）、database schema、domain model」；選了 (a)：`plan` 一路寫到票，
拿掉 `task-breakdown`。

## 決定

### D1. `plan` 的流程＝addyosmani 的五步，Step 2 前插一步設計

```
Step 1  唯讀收集（Plan Mode）：spec.md、.feature、docs/ 現況、api/openapi.yaml、既有程式碼與慣例；不寫程式
設計    seam 攤給使用者確認 → Domain Model → Database Schema → API 契約（有 HTTP 對外介面時；交給
        api-and-interface-design）→ 三份並排對一次（設計側的錯當場修，需求側的進 ## 缺口）
        → 五節技術現況與 Implementation／Testing Decisions → .feature 在 seam 上驗得到嗎
Step 2  依賴圖（節點是模組與契約）
Step 3  垂直切片（tracer bullet；wide refactor 用 expand–contract）
Step 4  寫任務：每張票的行為、驗收條件、驗證命令、Blocked by、會碰的模組、契約、資料表、接合點、大小
Step 5  排序與 checkpoint：高風險在前、每 2–3 張票一個 checkpoint；要不要並行（ASK）
ASK     把切法給使用者看（mattpocock to-tickets 的 Quiz）：粒度、Blocked by、要不要合併或再切；點頭才寫
WRITE   plan.md 與 tasks.md → STOP，告訴使用者 /implement <票號>
```

### D2. 產物

- `plan.md`：Overview、Tech Stack 到 Testing Decisions（照舊，自成一體）、Domain Model、Database Schema、API 契約、
  依賴圖、Risks & Mitigations、Parallel Work、缺口（＝addyosmani 的 Open Questions）、Further Notes。**拿掉
  `## 系統介面`（盤點＋Wave）**——addyosmani 沒有這一步，要盤點的由 Step 1 與依賴圖涵蓋
- `tasks.md`：格式不變（`task-breakdown/references/tasks-format.md` 用 `git mv` 搬到 `plan/references/`）；
  `implement` 照舊一張票一個 session
- addyosmani 的「不覆寫還有沒勾的計畫」：再跑時已勾的票不動、沒勾的照新設計重切（沿用 `task-breakdown` 原規則）

### D3. 搬家與刪除

| 檔 | 去處 |
| --- | --- |
| `task-breakdown/SKILL.md` 的規則（切法三條、Blocked by、接合點、wide refactor、再跑） | 併進 `plan/SKILL.md` Step 3–5，細節放 `plan/references/` |
| `task-breakdown/references/tasks-format.md` | `git mv` → `plan/references/tasks-format.md` |
| `task-breakdown/` | `git rm -r` |
| `plan/references/design-format.md` 的 `## 系統介面` 一節 | 刪 |

`plan/SKILL.md` 維持 500 行以內，超過就把細節下放到 references。

### D4. 下游與文件

`implement`、`archive`、`spec`、`formulation` 等提到 `task-breakdown` 的地方改指 `plan`；README、CLAUDE.md、PLAN.md、
docs/ai-sdlc.md、RESOURCES.md（addyosmani `planning-and-task-breakdown` 與 mattpocock `to-tickets` 兩列）跟著改。

## 驗證

- `audit_skill.py` 全過（`api-and-interface-design` 的 F6 照舊接受）；`claude plugin validate .`；單元測試
- `grep -rn "task-breakdown" skills README.md CLAUDE.md PLAN.md docs/ai-sdlc.md` 只剩歷史說明
- 重跑一次試跑到 `implement` 第一張票
