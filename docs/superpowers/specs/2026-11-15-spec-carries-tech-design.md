# 技術設計搬進 `spec.md`

**狀態**：已批　**日期**：2026-11-15

## 意圖

需求方的話：「每一個新功能應該是獨立的一份 spec.md」，以及 addyosmani 的
`spec-driven-development` 和 mattpocock 的 `to-spec` 寫 spec 的方式「比較好」。

兩者的 spec 都是**一份讀完就能開工的文件**：除了要做成什麼，還寫用什麼做、怎麼驗。

- addyosmani：Objective、Tech Stack、Commands、Project Structure、Code Style、Testing Strategy、
  Boundaries、Success Criteria、Out of Scope、Open Questions
- mattpocock：Problem Statement、Solution、User Stories、Implementation Decisions、Testing Decisions、
  Out of Scope、Further Notes

我們現在把技術設計放在另一份 `plan.md`，`spec.md` 只有需求。`spec.md` 本來就是 `discovery.md`
收斂之後寫的，而 DISCOVERY 也會問技術面向，所以這份文件把那些答案留在 `spec.md` 裡。

## 決定

### D1. `spec.md` 分成需求部分與技術部分

需求部分的章節不動。技術部分是新的，接在 `## Success Metrics` 之後、`## 回 DISCOVERY 補問` 之前：

| 章節 | 裝什麼 | 對應來源 |
| --- | --- | --- |
| `## Tech Stack` | 語言、框架、主要依賴與版本；不准引入的依賴 | addyosmani；原 `plan.md` `## Tech Preferences & Constraints` |
| `## Commands` | build、test、lint、跑場景、閘門的完整命令 | addyosmani |
| `## Project Structure` | 程式碼、測試、文件各住哪 | addyosmani |
| `## Code Style` | 命名、格式，加一段專案裡真實的程式碼當例子 | addyosmani |
| `## Boundaries` | 總是要做／先問／絕對不做 三張清單 | addyosmani；原 `plan.md` `## Known Boundaries` |
| `## Implementation Decisions` | 動哪些模組、介面、API 契約、schema、架構決定、互動 | mattpocock；原 `plan.md` 同名節，格式不變 |
| `## Testing Decisions` | 測試框架與測試放哪、seam、每個 AC 在哪一層測、prior art | 兩者；原 `plan.md` 同名節，加上框架與位置 |

**每份 spec 自成一體。** Tech Stack、Commands、Project Structure、Code Style 跟 feature 無關，
每份照樣寫——讀者不必去翻別的檔。第二個 feature 起，從上一份 `spec.md` 抄過來，逐條確認現在
還成立再寫；專案自己有這些文件（`CLAUDE.md`、`ARCHITECTURE.md`、Makefile）時照它們寫，並標來源。

### D2. 技術部分在 `已對焦` 之後才寫

順序：

1. `spec` 寫需求部分（`待對焦`）
2. `example-mapping`、`spec-by-example`：需求方判就緒 → `已對焦`
3. **`spec` 寫技術部分**；seam 攤給使用者確認（原本 `plan` 第 2 步的 Ask user，照搬）
4. `formulation` 寫 `.feature`
5. `plan` 寫執行計畫
6. `task-breakdown` 拆票

需求方判就緒只看需求部分；技術決定一定建立在講定的需求上。狀態值不增加：技術部分寫完沒有
另一個狀態，靠 `spec` 的 Verification 檢查。

### D3. 技術部分的來源標記

沿用 `plan.md` 現在的三選一（`PRD §x`／`Q-<n>`／`推論`），加一種：

| 標記 | 意思 |
| --- | --- |
| `專案 <路徑>` | 從專案既有的檔讀來的事實（`← 專案 Makefile`、`← 專案 docs/ARCHITECTURE.md`） |

`推論` 照舊要同時列進缺口——技術部分的缺口寫進 `## 回 DISCOVERY 補問`（需求的洞）或
`## Further Notes`（留給 IMPLEMENT 的技術選擇）。

**技術部分不訪談**，跟需求部分同一條紀律：問答裡沒有的，退回 `discovery`。唯一的 Ask user 是
seam 確認，與原本 `plan` 相同。

### D4. `plan.md` 只留執行計畫

| 章節 | 狀態 |
| --- | --- |
| `## Overview` | 留 |
| `## Risks`（含影響程度） | 留——風險決定票的順序 |
| `## Parallel Work` | 留——要不要並行、接合點 |
| `## 缺口` | 留——執行計畫自己的洞；`task-breakdown` 照舊在這裡加列 |
| `## Tech Preferences & Constraints`、`## Known Boundaries` | 搬到 `spec.md` 的 `## Tech Stack`、`## Boundaries` |
| `## Implementation Decisions`、`## Testing Decisions` | 搬到 `spec.md` |

`plan` 的職責變成「怎麼排、哪裡危險、能不能並行」，讀 `spec.md` 的技術部分，不改它。
`plan` 不改需求的紀律延伸成「不改 `spec.md`」：技術部分的洞回 `spec`。

### D5. 下游跟著改讀的位置

- `discovery`：`歸屬` 的值域從 `plan.md` 技術節名改成 `spec.md` 技術部分的章節名；
  `references/technical-probes.md` 的「進 plan.md 哪一節」跟著改
- `task-breakdown`：模組、內迴圈、驗證命令改從 `spec.md` 的 `## Implementation Decisions`、
  `## Testing Decisions`、`## Commands` 抄
- `implement`、`bdd`：讀技術設計的地方從 `plan.md` 改成 `spec.md`；`implement` 的閘門命令改讀
  `spec.md` `## Commands`
- `check_spec.py`：不用改。原本以為它會把 `## User Stories` 之後全部當成 story，讀了程式碼才知道
  每個解析器都已經停在下一個 `## `；加一支測試釘住這件事。技術部分的來源標記先不檢查，跟原本的
  `plan.md` 一樣——寫進「已知缺口」

## 非目標

- 不改需求部分的章節與 FR／AC 的規則
- 不新增狀態值
- 不檢查技術部分的來源標記（留給 `review`）
- 不把 `plan` 併進 `task-breakdown`

## 考慮過的做法

| 做法 | 取捨 | 結果 |
| --- | --- | --- |
| A. 專案層級的 `specs/project.md` 放 Tech Stack 等，`plan.md` 不動 | 不重複；但 spec 不再自成一體，讀者要翻三份檔 | 捨棄——需求方要每個 feature 一份獨立的 spec |
| **B. 技術部分搬進 `spec.md`，`plan.md` 只留執行計畫**（本文件） | 一份讀完就能開工；`spec` 變大、多一段技術紀律 | **建議** |
| C. 連 `plan` 也併進 `task-breakdown` | 少一個 skill；`task-breakdown` 要兼管風險、並行、拆票 | 捨棄 |

## 會動到的檔案

- `skills/spec/SKILL.md`、`references/spec-format.md`（技術部分的章節、範例、來源標記）
- `skills/plan/SKILL.md`、`references/plan-format.md`（只留執行計畫）
- `skills/discovery/references/technical-probes.md`、`references/discovery-format.md`（`歸屬`）
- `skills/task-breakdown/`、`skills/implement/SKILL.md`、`skills/bdd/`（讀的位置）
- `tests/test_check_spec.py`（技術部分不被當成 story 的回歸測試）
- README、CLAUDE.md、PLAN.md、docs/ai-sdlc.md、RESOURCES.md

## 怎麼驗

v3 跑完、這批與 v2／v3 的規則問題一起改完之後，拍新快照派 worker 從 baseline 重跑：`spec.md`
技術部分每條有來源、seam 有使用者確認的紀錄、`plan.md` 只剩四節、票的模組與驗證命令能從
`spec.md` 抄到。

## 最不確定的地方

- **`spec.md` 會變長。** 需求方判就緒時要跳過技術部分；靠 D2 的順序（判就緒時技術部分還沒寫）
  擋住，但修補重跑時技術部分已經在了
- **每份 spec 抄一次 Commands 等。** 專案改了 Makefile，舊 spec 不會跟著改——它們是寫的當下的快照，
  跟 `plan.md` 現在的定位一樣
