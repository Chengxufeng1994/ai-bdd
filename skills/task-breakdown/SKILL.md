---
name: task-breakdown
description: >
  SPEC 說用什麼做、PLAN 說怎麼排，TASK BREAKDOWN 把它拆開——讀 `spec.md` 的技術部分、`plan.md`
  的風險與接合點、`.feature`，切成一疊 tracer bullet 寫進一份 `tasks.md`：一張票一節，每張穿過所有層、自己可
  demo、塞得進一個全新 session，宣告被哪幾張票 block、會碰哪些接合點，由上而下就是順序，
  開頭寫哪幾張可以同時做。不改設計、不改需求。BDD 六步流程 PLAN 步驟的第二個 skill。
  觸發詞：「這批要怎麼切」「拆成幾張票」「切票」「拆成 tasks」「先做哪一個」「排實作順序」
  「這些票誰卡誰」「開工順序」「哪幾張可以同時做」「並行跑票」「plan 寫完了下一步」。
  English: break the plan into tasks, slice the work into tickets, what order should
  I build these, vertical slices, which tickets can run in parallel, what can I start first.
disable-model-invocation: true
---

# TASK BREAKDOWN — 把怎麼做拆成票

| | |
| --- | --- |
| **職責** | 把技術設計拆成一疊 tracer bullet：一張票一節，每張穿過所有層、自己可 demo，標明誰卡誰、會碰哪些接合點 |
| **讀** | `spec.md`（`## User Stories` 與技術部分）、`plan.md`（`## Risks`、`## Parallel Work`）、涵蓋的 `.feature` |
| **寫** | `specs/ongoing/<date>-<feature>/tasks.md`（一張票一個票節）；`plan.md` 只在 `## 缺口` 加列 |
| **上游** | `plan`（風險與「要不要並行」寫好之後） |
| **下游** | `implement <票號>`（一張票一個全新 session） |
| **不做** | 不改 `spec.md`（需求部分與技術部分）、不改 `plan.md` 的既有內容、不派工、不寫程式 |

**不改設計。** seam、模組、契約、schema 是 `spec` 在技術部分定的，風險與接合點是 `plan` 定的；這一步只把它們拆開排序。
拆的時候發現設計有洞，寫進 `plan.md` 的 `## 缺口` 並交回 `spec`，不要順手補。

## 使用時機

- `spec.md` 的技術部分與 `plan.md` 都寫好了，要決定切成幾張票、先做哪一張
- 一批工作大到一個 session 做不完，要拆給多個 session 接力或同時做
- 技術設計沒變，但票切得不對或順序要調，只重切票

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| `spec.md` 沒有技術部分（seam、模組、契約） | `spec` 第 6 步 |
| 還沒有 `plan.md` | `plan` |
| 要改 seam、介面、schema | `spec` |
| 要改 FR、AC、story | `spec` |
| 要把太大的 **story** 切小 | `story-splitting`——它切的是需求，本 skill 切的是實作 |
| 要實際做某一張票 | `implement` |
| 派幾個 agent、誰做哪一張 | 不在本 skill；這裡只講哪幾張**可以**同時做 |

## 參考檔案

- `references/tasks-format.md` — `tasks.md` 的骨架、誰寫哪幾行、接合點的三種與新建的東西、三條切票規則、
  驗收條件要能失敗、wide refactor 的例外。**寫檔之前先讀**
- `../spec/references/tech-format.md` — `spec.md` 技術部分各節裝什麼
- `../plan/references/plan-format.md` — `plan.md` 各節裝什麼，`## 缺口` 的寫法

## 前置確認（Ask user）

Check if: 同目錄的 `spec.md` 有 `## Implementation Decisions`、`## Testing Decisions` 與 `## Commands`，
而且 `plan.md` 讀得到、有記了「要不要並行」的 `## Parallel Work`。前者沒有就停下，說先跑 `spec`
第 6 步；後者沒有就說先跑 `plan`——沒有技術部分時切票，票裡的層與內迴圈
全是這一步自己設計的，而沒有人同意過那份設計。

---

## 流程

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| 1 | 一次讀完整批 | — |
| 2 | 切票：最薄的完整路徑先，其餘依 FR | — |
| 3 | 標 Blocked by 與接合點 | — |
| 4 | 攤出來確認 | 要 |
| 5 | 寫 `tasks.md`：票節帶驗證命令，波與波之間放檢查點 | — |

### 1. 讀

一次讀完整批：`spec.md` 的 `## User Stories` 與技術部分、`plan.md` 全文、涵蓋的每個 `.feature`
（場景與 tag）。票之間的依賴只有並排時才看得見；分批切會各自長出一套順序，然後在實作時
發現它們互相卡住。

順便看一眼專案既有的結構——找前置整理（prefactoring）的機會。

### 2. 切票

先切出**最薄的完整路徑**當第一張票：碰到每一層、邏輯最少。它證明的是接線，不是行為。
專案已經有骨架就跳過。

其餘依 `spec.md` 的 FR 切，一張票一組相關的 `@example-N.M`。碰到 `plan.md` `## Risks` 裡影響程度
「高」的票，在 `Blocked by` 允許的範圍內排前面——越早撞到越便宜。每張票問一次「做完之後我能
demo 什麼」，答不出來就是切錯了。每張票的**內迴圈**從 `spec.md` `## Testing Decisions` 的
測試分層來。

每張票再標兩樣：

- **模組**：照抄 `spec.md` `## Implementation Decisions`「動哪些模組」裡這張票碰的那幾個，不寫檔案路徑
- **大小**：XS／S／M／L／XL（尺寸表 → `references/tasks-format.md`〈大小〉）

MUST NOT: 交出標 XL 的票——再拆。L 也 SHOULD 再拆，拆不開就在攤出來確認時講明為什麼。

三條規則、驗收條件要能失敗、wide refactor 的例外 → `references/tasks-format.md`。

### 3. 標 Blocked by 與接合點

每張票兩個欄位，回答兩個不同的問題：

| 欄位 | 回答 | 怎麼判 |
| --- | --- | --- |
| `Blocked by` | 誰要先做完 | 這張票的驗收條件，在那張票做完之前有沒有可能綠 |
| `接合點` | 合回時要做什麼、誰不能跟它同時做 | 照抄 `plan.md` `## Parallel Work` 的路徑與種類（`產生`／`追加`／`爭用`） |

兩張以上的票都需要一個還不存在的 package 或共用的測試輔助函式時，指定一張票擁有它（通常是
`00` 前置整理），其他票 `Blocked by` 它——各自新建，合回時會撞成重複宣告。

IMPORTANT: `plan.md` 的 `## Parallel Work` 沒寫到、切票時才發現的接合點，不要自己編路徑：加
一列到 `plan.md` 的 `## 缺口`。`追加` 與 `爭用` 拿不準時標 `爭用`。

`Blocked by` 都完成、沒有共同的 `爭用` 接合點、沒有兩張都要新建的東西，才可以同時做；
`## Parallel Work` 選了「不要」時一律不並行。細節與例子 → `references/tasks-format.md`
〈接合點與同時跑〉。

### 4. 攤出來確認（Ask user）

MUST: 寫檔之前先把票的清單攤給使用者看——編號、標題、大小、Blocked by、接合點、交付什麼行為、
哪幾張可以同時開始——並問三件事：

1. 粒度對不對（太粗／太細；大小標 L 的是哪幾張）
2. Blocked by 是不是真的卡住，還是只是「先做比較順」
3. 有沒有該合併或該再拆的

IMPORTANT: **切太細是這一步最常見的失效**，而且它不會自己被發現——一疊原子化的小票看起來
很專業，實際上把「這幾件事其實是同一件」的資訊丟掉了。為了多並行而把一件事切成兩張 `爭用`
同一段邏輯的票，換到的是零：它們照樣不能同時做。

使用者確認之後才寫檔。

### 5. 寫 `tasks.md`

照 `references/tasks-format.md`：開頭三行，然後一波一節、依賴序寫票節。

- 每張票寫 `**驗證**`：這張票的場景命令（tag 運算式由 `涵蓋場景` 組成）與閘門命令，都從 `spec.md` `## Commands` 抄——`implement`
  照著跑，不自己組（〈驗證〉）
- 一波一個 `## 第 N 波`，波與波之間放檢查點；不並行時每 2–3 張票一個（〈檢查點〉）

## 不改設計：有洞怎麼辦

MUST NOT: 在這一步新增、刪除或改寫 `spec.md` 技術部分的任何決定，或任何 FR、AC、story。

判準：指著票裡的任何一句話，說得出它來自 `spec.md` 的哪條 FR、技術部分的哪一個決定、
或哪個 `@example-N.M` 嗎？說不出來就是這一步僭越了。

切票時發現需要、但技術部分沒定案的（兩張票對同一個介面的假設不一樣、某個 AC 不知道在
哪一層測），在 `plan.md` 的 `## 缺口` 加一列，告訴使用者跑 `spec` 補技術部分；接合點的洞請使用者
跑 `plan` 補——這是本 skill 唯一寫 `plan.md` 的地方。

## 再跑：重切，不重編

`tasks.md` 已經存在時（還是 2026-11-16 以前的 `tasks/` 目錄時，先併成一份：照 `index.md` 的順序把每張
票檔的內容放進對應的票節、標題 `# 01 — …` 改 `### 01 — …`、`## 缺口` 改 `#### 缺口`，其餘一字不改，
`git rm -r tasks/`，單獨一個 commit）：

- 驗收條件已經有勾的票**不動**——有人做過了，改它等於改掉一份已經驗過的約定。判斷看**主線上**的
  票節；正在 worktree 裡做的票（`.feature` 是 `@wip`、主線上的票節還沒勾）也算開工了
- 票節的欄位各有主人：`接合點`、`Blocked by`、`大小` 歸本 skill，開工了的票也只准改這三欄（多半是
  補漏列的接合點）；驗收條件的勾、`**驗證**` 的勾、`**檢查點**:`、`#### 缺口` 歸 `implement`，本 skill 不動（唯二例外是下兩條在列尾加的標記）——
  兩邊改的是不同行，合回時才不會互相蓋掉
- 還沒開始的票可以改、合併、再拆；作廢的票留檔、標題後加「（作廢：<原因>）」，編號不重用，
  新票接在最大編號之後
- 票節 `#### 缺口` 裡 `implement` 記下的洞：屬於設計的，搬一列到 `plan.md` 的 `## 缺口`
  並在票節那一條後面標「→ 已記入 plan.md」
- 喚醒或修補之後：舊票 `#### 缺口` 裡已經被 `discovery.md` 新題回答的列（判準：那一題的答案正是那一列
  在等的），在列尾標「→ 已回答 Q-<n>」，交給新票做；不標，下一個讀缺口的人會再猜一次

## 產物

```
specs/ongoing/<date>-<feature>/tasks.md   開頭：來源、哪幾張可以同時開始；一波一節，一張票一個票節（00 前置整理有才寫）
```

MUST: 只寫入 `specs/ongoing/<date>-<feature>/tasks.md`，以及 `plan.md` `## 缺口` 的新列。不修改專案
既有的任何其他檔案。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 票裡每一句都指得出來自哪條 FR、技術部分哪個決定、或哪個 `@example-N.M`
- [ ] 每張票答得出「做完之後我能 demo 什麼」
- [ ] 碰到 `## Risks` 影響程度「高」的票，在 `Blocked by` 允許的範圍內排在前面
- [ ] 每張票都有 `Blocked by` 與 `接合點` 兩欄（沒有就寫「無」），每個接合點標了種類且照抄 `plan.md` `## Parallel Work`
- [ ] 每張票標了模組（照抄 `spec.md`，沒有檔案路徑）與大小，沒有 XL
- [ ] 每張票的 `**驗證**` 有閘門命令；有 `涵蓋場景` 的票另有場景命令，tag 運算式跟 `涵蓋場景` 一一對應
- [ ] `tasks.md` 每一波（不並行時每 2–3 張票）之後有檢查點
- [ ] `tasks.md` 開頭寫了哪幾張可以同時開始（或「不適用」），那幾張沒有共同的 `爭用` 接合點、也沒有兩張都要新建的東西
- [ ] 寫檔前把票的清單攤給使用者、問過粒度／Blocked by／合併或再拆三件事，使用者確認了
- [ ] `spec.md` 的需求部分與技術部分沒有任何新增、刪除或改寫；設計的洞進了 `plan.md` 的 `## 缺口`
- [ ] 再跑時，驗收條件已有勾的票沒有被改，編號沒有重排或重用

## 完成後

告訴對方：

1. 幾張票、第一張是哪一張（以及**為什麼是它**）、有沒有前置整理
2. 哪幾張可以同時開始；哪幾張邏輯上獨立卻因為 `爭用` 或新建的東西排成先後
3. `plan.md` 的 `## 缺口` 新增了幾列（技術部分的洞先跑 `spec`，接合點的洞先跑 `plan`）

下一步 `implement <票號>`，一張票一個全新 session；同時做的票各開一個 worktree。
