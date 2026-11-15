---
name: bdd-tasks
description: >
  PLAN 說怎麼做，TASKS 把怎麼做拆開——讀 `plan.md` 的技術設計、已對焦的 `spec.md` 與
  `.feature`，切成一疊 tracer bullet 寫進 `tasks/`：一張票一個檔，每張穿過所有層、自己可
  demo、塞得進一個全新 session，宣告被哪幾張票 block、會動哪些共用檔案，`index.md` 記順序
  與哪幾張可以同時做。不改設計、不改需求。BDD 六步流程 PLAN 步驟的第二個 skill。
  觸發詞：「這批要怎麼切」「拆成幾張票」「切票」「拆成 tasks」「先做哪一個」「排實作順序」
  「這些票誰卡誰」「開工順序」「哪幾張可以同時做」「並行跑票」「plan 寫完了下一步」。
  English: break the plan into tasks, slice the work into tickets, what order should
  I build these, vertical slices, which tickets can run in parallel, what can I start first.
---

# TASKS — 把怎麼做拆成票

一份寫好技術節的 `plan.md`、已對焦的 `spec.md` 與涵蓋的 `.feature` 進來，出去的是同目錄的
`tasks/`：**一張票一個檔**——每張是穿過所有層的 tracer bullet——加一份 `index.md`，記誰卡誰、
誰跟誰會動到同一個檔、哪幾張可以同時開始。

**不改設計。** seam、模組、契約、schema 是 `bdd-plan` 定的；這一步只把它們拆開排序。
拆的時候發現設計有洞，寫進 `plan.md` 的 `## 缺口`，不要順手補。

## 使用時機

- `plan.md` 的技術節寫好了，要決定切成幾張票、先做哪一張
- 一批工作大到一個 session 做不完，要拆給多個 session 接力或同時做
- 技術設計沒變，但票切得不對或順序要調，只重切票

## Skill Boundaries

- 還沒有 `plan.md`、或它沒有技術節（seam、模組、契約）→ `bdd-plan`
- 要改 seam、介面、schema → `bdd-plan`；要改 FR、AC、story → `bdd-spec`
- 要把太大的 **story** 切小 → `story-splitting`，它切的是需求，本 skill 切的是實作
- 要實際做某一張票 → `bdd-implement`
- 派幾個 agent、誰做哪一張 → 不在本 skill；這裡只講哪幾張**可以**同時做

## 參考檔案

- `references/tasks-format.md` — `index.md` 與票檔的骨架、共用檔案怎麼列、三條切票規則、
  驗收條件要能失敗、wide refactor 的例外。**寫檔之前先讀**
- `../bdd-plan/references/plan-format.md` — `plan.md` 各節裝什麼，`## 缺口` 的寫法

## 前置確認（Ask user）

Check if: `specs/<date>-<feature>/plan.md` 讀得到，而且有 `## Implementation Decisions` 與
`## Testing Decisions`。沒有就停下，說先跑 `bdd-plan`——沒有技術節時切票，票裡的層與內迴圈
全是這一步自己設計的，而沒有人同意過那份設計。

---

## 不改設計

MUST NOT: 在這一步新增、刪除或改寫 `plan.md` 的技術節，或任何 FR、AC、story。

判準：指著票裡的任何一句話，說得出它來自 `spec.md` 的哪條 FR、技術節的哪一個決定、
或哪個 `@example-N.M` 嗎？說不出來就是這一步僭越了。

切票時發現需要、但技術節沒定案的（兩張票對同一個介面的假設不一樣、某個場景不知道在
哪一層測），在 `plan.md` 的 `## 缺口` 加一列，告訴使用者，請他跑 `bdd-plan` 補——這是本
skill 唯一寫 `plan.md` 的地方。

## 流程

### 1. 讀

一次讀完整批：`plan.md` 全文、`spec.md` 的 `## User Stories`、涵蓋的每個 `.feature`
（場景與 tag）。票之間的依賴只有並排時才看得見；分批切會各自長出一套順序，然後在實作時
發現它們互相卡住。

順便看一眼專案既有的結構——找前置整理（prefactoring）的機會。

### 2. 切票

先切出**最薄的完整路徑**當第一張票：碰到每一層、邏輯最少。它證明的是接線，不是行為。
專案已經有骨架就跳過。

其餘依 `spec.md` 的 FR 切，一張票一組相關的 `@example-N.M`。每張票問一次「做完之後我能
demo 什麼」，答不出來就是切錯了。每張票的**內迴圈**從 `plan.md` `## Testing Decisions` 的
測試分層來。

三條規則、驗收條件要能失敗、wide refactor 的例外 → `references/tasks-format.md`。

### 3. 標 Blocked by 與共用檔案

每張票兩個欄位，回答兩個不同的問題：

| 欄位 | 回答 | 怎麼判 |
| --- | --- | --- |
| `Blocked by` | 誰要先做完 | 這張票的驗收條件，在那張票做完之前有沒有可能綠 |
| `共用檔案` | 誰不能跟它同時做 | 這張票會改、別張票也可能改的檔案——照 `plan.md` `## Implementation Decisions` 的模組與契約推 |

IMPORTANT: `共用檔案` 拿不準就列。它是預測，票還沒做；多列只是少一點並行，漏列是兩個
worktree 做完合不回來——而那時兩張票都已經「完成」了。

兩張票 `Blocked by` 都完成、`共用檔案` 沒有交集，才可以同時做。細節與例子 →
`references/tasks-format.md`〈共用檔案與同時跑〉。

### 4. 攤出來確認（Ask user）

MUST: 寫檔之前先把票的清單攤給使用者看——編號、標題、Blocked by、共用檔案、交付什麼行為、
哪幾張可以同時開始——並問三件事：

1. 粒度對不對（太粗／太細）
2. Blocked by 是不是真的卡住，還是只是「先做比較順」
3. 有沒有該合併或該再拆的

IMPORTANT: **切太細是這一步最常見的失效**，而且它不會自己被發現——一疊原子化的小票看起來
很專業，實際上把「這幾件事其實是同一件」的資訊丟掉了。為了多並行而把一件事切成兩張共用
同一個檔的票，換到的是零：它們照樣不能同時做。

使用者確認之後才寫檔。

### 5. 寫 `tasks/`

照 `references/tasks-format.md`：先票檔，後 `index.md`；`index.md` 的 `Blocked by`、
`共用檔案` 兩欄跟票檔一字不差。

## 再跑：重切，不重編

`tasks/` 已經存在時：

- 驗收條件已經有勾的票**不動**——有人做過了，改它等於改掉一份已經驗過的約定
- 還沒開始的票可以改、合併、再拆；作廢的票留檔、標題後加「（作廢：<原因>）」，編號不重用，
  新票接在最大編號之後
- 票檔 `## 缺口` 裡 `bdd-implement` 記下的洞：屬於設計的，搬一列到 `plan.md` 的 `## 缺口`
  並在票檔那一條後面標「→ 已記入 plan.md」

## 產物

```
specs/<date>-<feature>/tasks/
├── index.md         順序表 ＋ 哪幾張可以同時開始
├── 00-<slug>.md     前置整理（有才建）
└── 01-<slug>.md     一張票一個檔
```

MUST: 只寫入 `specs/<date>-<feature>/tasks/`，以及 `plan.md` `## 缺口` 的新列。不修改專案
既有的任何其他檔案。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 票裡每一句都指得出來自哪條 FR、技術節哪個決定、或哪個 `@example-N.M`
- [ ] 每張票答得出「做完之後我能 demo 什麼」
- [ ] 每張票都有 `Blocked by` 與 `共用檔案` 兩欄（沒有就寫「無」），`index.md` 的同名欄位跟票檔一字不差
- [ ] `index.md` 寫了哪幾張可以同時開始，而且那幾張的 `共用檔案` 沒有交集
- [ ] 寫檔前把票的清單攤給使用者、問過粒度／Blocked by／合併或再拆三件事，使用者確認了
- [ ] `plan.md` 的技術節與 `spec.md` 的 FR、AC、story 沒有任何新增、刪除或改寫；設計的洞進了 `plan.md` 的 `## 缺口`
- [ ] 再跑時，驗收條件已有勾的票沒有被改，編號沒有重排或重用

## 完成後

告訴對方：

1. 幾張票、第一張是哪一張（以及**為什麼是它**）、有沒有前置整理
2. 哪幾張可以同時開始；哪幾張邏輯上獨立卻因為共用檔案排成先後
3. `plan.md` 的 `## 缺口` 新增了幾列（有的話先跑 `bdd-plan`）

下一步 `bdd-implement <票號>`，一張票一個全新 session；同時做的票各開一個 worktree。
