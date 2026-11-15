---
name: bdd-implement
description: >
  照 `tasks/` 的票一張一個 session，outside-in 把票做到綠：先把這張票的場景接上 step
  definition、看到對的紅，再把 green 與 refactor 委派給 `bdd`（規則先進 domain、table-driven
  單元測試先紅），場景全綠後把 `.feature` 從 `@ready` 翻到 `@review`、勾掉票檔的驗收條件、
  跑專案的驗證閘門。不改 `spec.md`、不改 `.feature` 的例子、不在 handler 放規則。
  BDD 六步流程的 IMPLEMENT 步驟。
  觸發詞：「做第 N 張票」「開始實作」「bdd-implement 03」「把這張票做綠」
  「寫 step definition」「讓場景變綠」「照 plan 做」「下一張票」。
  English: implement ticket N, make the scenarios green, write the step
  definitions, outside-in from the feature file, work the plan.
---

# IMPLEMENT — 一張票，從紅到綠

一張 `tasks/` 的票進來，出去的是：那張票涵蓋的場景全綠、專案的驗證閘門綠、票檔的驗收
條件打勾；涵蓋同一個 `.feature` 的票都做完時，`.feature` 標 `@review`。

**外迴圈先紅，內迴圈先紅，然後才寫程式。** 外迴圈是 `.feature` 的場景（godog／
Cucumber），內迴圈是 domain 與 use case 的單元測試。兩圈都先看到紅，寫的每一行
才有理由。

本 skill 管**整張票**：標 tag、讓測試層先紅、跑閘門、勾驗收條件、合回主線。票裡的轉綠與整理
委派給 `bdd`——它一次一個場景、一次一個入口。

不改 `spec.md`。不改 `.feature` 的例子。不做 `plan.md` 技術節沒做的設計決定——發現有洞就
記下來，不順手補。

## 使用時機

- `tasks/` 切好了，要做第一張（或下一張 Blocked by 全完成的）票
- 幾張票同時做，每張各開一個 session 與 worktree
- 上一個 session 做到一半，要從 `@wip` 的檔接手

## Skill Boundaries

- 還沒有 `.feature` → `bdd-formulation`；`plan.md` 沒有技術節（seam、埠、契約）→ `bdd-plan`
- 還沒有 `tasks/index.md`（包括票還寫在 `plan.md` `## 票` 裡的舊格式）→ `bdd-tasks`
- 票之間的順序、要不要再切、哪幾張能同時做 → `bdd-tasks`
- 只推進**一個場景**的 red／green／refactor 其中一步，而不是做一張票 → 使用者手動叫 `/bdd`。
  本 skill 做票時第 3、4 步本來就委派它
- 場景綠了之後「真的滿足它宣稱的情境嗎、鏈斷了嗎」→ REVIEW（`bdd-review`，未實作）
- 紅的是 bug 還是規格過期 → VERIFY（`bdd-verify`，未實作）；本 skill 只處理自己這張票
  讓紅變綠，不判讀別人的紅
- 要改規則或例子 → 回 `bdd-discovery` 開一題；本 skill 不碰 `spec.md`、不碰 `.feature` 的例子

## 參考檔案

- `../bdd-formulation/references/state-tags.md` — `@wip`／`@review` 什麼時候寫
- `../bdd/SKILL.md` — 第 3、4 步委派的 green／refactor 入口；`../bdd/rules/red.md` 與
  `../bdd/references/step-definitions.md` — 第 2 步測試層先紅時照它寫
- `../bdd-tasks/references/tasks-format.md` — 票的欄位，驗收條件為什麼要能失敗，接合點與同時跑
- `../bdd-plan/references/plan-format.md` — 技術節各節裝什麼
- 專案自己的架構文件（實測場是 `docs/ARCHITECTURE.md`、`docs/DATAFLOW.md`、各層 `doc.go`）
  ——分層、轉換、錯誤分類都照它，本 skill 不另訂

---

## 前置確認（Ask user）

1. **哪一張票**：`bdd-implement <票號>`。沒帶票號就取 `tasks/index.md` 裡第一張
   Blocked by 全部完成（那些票檔的驗收條件全勾）、自己還沒全勾的票，講出來讓對方改
2. **驗證閘門是什麼命令**：讀專案的 README／Makefile（實測場是 `make verify`）。
   沒有就問——一張票「做完」的定義就是這個命令綠，不能猜
3. 讀得到 `specs/<date>-<feature>/tasks/` 的這張票檔、`plan.md`、`spec.md` 與這張票涵蓋的 `.feature`

NEVER: `plan.md` 沒有技術節時開始。seam、埠的形狀、錯誤分類都在那裡；沒有它，這一步會
自己設計，而設計出來的東西沒有任何場景驗得到。

MUST: 一個 session 一張票。撿起票的人只看過票檔與 `plan.md`——這是 `bdd-tasks`
切票的上界，也是這一步不要一次做三張的理由：做到第二張時第一張的脈絡已經在擠。

SHOULD: 跟別張票同時做時，一張票一個 worktree。開工前看 `tasks/index.md`：這張票跟正在做
的票有共同的 `爭用` 接合點、或兩張都要新建同一個東西，就不要同時開——它們合不回來。

MUST: 開工時工作區是乾淨的；做完這張票一個 commit，只 stage 這張票動到的檔。合回與回滾都
以這個 commit 為單位。

---

## 流程

### 1. 讀票，標 `@wip`

讀票檔 `tasks/<票號>-<slug>.md`：交付什麼行為、涵蓋哪些 `@example-N.M`、內迴圈要寫哪些單元測試、
驗收條件。再讀 `plan.md` 的 `## Implementation Decisions` 與 `## Testing Decisions`
——這張票動哪些模組、埠長什麼形狀、seam 在哪一層、每個例子在哪一層測。`spec.md` 只讀
對應的 FR 與 AC，看規則的原話。

把這張票涵蓋的 `.feature` 檔從 `@ready` 改成 `@wip`（`state-tags.md`）。下一個 session
看到 `@wip` 就知道有人做到一半。

### 2. 測試層先紅（本 skill 自己做，不委派）

把這張票涵蓋的每個場景都接上 step definition，直到**沒有任何 undefined step**、失敗只能是
斷言或產品行為。寫步驟照 `../bdd/references/step-definitions.md`；判斷紅得對不對照
`../bdd/rules/red.md`〈什麼是「對的紅」〉。只寫測試層，不寫產品程式碼。

```bash
go test -count=1 ./test/acceptance/ -godog.tags=@rule-2          # 一條規則
go test -test.v ./test/acceptance/ -test.run "^TestFeatures/停留"  # 用名字選
```

MUST: 沒看到紅就不要寫程式。一個一寫就綠的場景沒有驗到任何東西——通常是 step
definition 已經被別的場景定義過，而它的斷言沒抓到這張票要做的事。

這張票所有場景都紅得對了，才進第 3 步。測試層沒對齊就改產品程式碼，綠了也說不出是哪個
斷言在保護它。

### 3. 委派 `bdd` 的 green

讀 `../bdd/SKILL.md`，照它的 `green` 入口做，範疇與輸入由本 skill 給：

| 交給 `bdd` 的 | 內容 |
| --- | --- |
| 入口 | `green` |
| 範疇（Test Scope） | 這張票檔的 `涵蓋場景`，組成一個 tag 運算式（例 `@example-8.1 or @example-8.2`） |
| 票 | `tasks/<票號>-<slug>.md` 的路徑 |
| 設計 | `plan.md` 的 `## Implementation Decisions` 與 `## Testing Decisions` |

`bdd` 一個場景一個場景轉綠：規則先進 domain、先看到單元測試紅（`../bdd/references/inner-loop.md`）。
它停下來回報規格的洞或範疇外的耦合時，本 skill 把那一條記進票檔的 `## 缺口`，照
〈不越界〉處理——`bdd` 自己不寫 `specs/` 底下任何檔案。

### 4. 委派 `bdd` 的 refactor

範疇全綠之後，同樣的範疇交給 `bdd` 的 `refactor` 入口：在綠燈保護下整理這張票留下的重複
步驟、轉換與支撐層。沒有要整理的就跳過，在回報裡寫「無」。

### 5. 外迴圈綠，閘門綠

```bash
go test -count=1 ./test/acceptance/ -godog.tags=@rule-2   # 這張票的場景
make verify                                               # 專案的全部閘門
```

MUST: 兩個都綠才算做完。閘門包含 fmt、產碼是否過期、lint、`-race`——其中任何一項紅
而只回報「場景綠了」，是這一步最常見的「完成」謊言。

跑閘門時產碼工具會重跑（mockery 之類），新的埠要讓它產一次，否則「generated code
is stale」。

### 6. 收尾

- `.feature`：在主線上做時，這張票做完後**整個檔的場景都綠了**才從 `@wip` 改 `@review`；
  還有別張票涵蓋同一檔的場景就留 `@wip`。在 worktree 裡做時一律留 `@wip`，到〈合回主線〉
  第 4 步再判
- 票檔：勾掉這張票的驗收條件。只勾驗過的——「手動驗」要寫是哪個測試驗的，
  沒有測試就不要勾
- 暫定行為（票檔或 `plan.md` 的 `## 缺口`、或 `spec.md` 的 `## 回 DISCOVERY 補問` 裡的洞，實作時不得不選一個）：
  **用一個測試釘住**，測試的註解寫明它是暫定、等哪一題的答案。沒釘住的暫定行為
  會在答案出來之後被當成規格

### 7. 合回主線（只在 worktree 裡做時）

這張票由做它的這個 session 自己合回，不交給別人——衝突要靠這張票的脈絡才解得對。

1. **先把主線拉進 worktree**（rebase 或 merge 主線），衝突在這裡解，不在主線上解
2. **照接合點的種類解**（`../bdd-tasks/references/tasks-format.md`〈接合點與同時跑〉）：
   - `產生` → 不手動合：先解完原始檔，再跑產碼工具重新產生（實測場 `make gen`）；`go.sum`
     用 `go mod tidy`
   - `追加` → 取聯集：兩邊加的都留，然後檢查有沒有重複宣告
   - `爭用`、或票檔與 `index.md` 都沒列到的衝突 → **停下來回報**，不自己判——那是切票時
     沒預料到的耦合，要回 `bdd-tasks` 重排
3. **在合併後的程式碼上重跑閘門**，場景範圍是主線上所有驗收條件已勾完的票加上這張票
   （照 `index.md` 把它們的 `@example-N.M` 組成一個 tag 運算式）。綠了才合回主線
4. **判 `@review`**：讀 `index.md` 找出涵蓋這個 `.feature` 的票，逐一看票檔的驗收條件——
   全勾了，就是這次合回讓最後一張完成，**由這次合回翻 `@review`**，跟合回的 commit 一起進主線；
   還有沒勾的就留 `@wip`
5. 同時有兩張票要合回時，票號小的先合；後合的那張回到第 1 步

IMPORTANT: git 合得乾淨不等於合對了。兩張票各自新建同一個 package，git 不會報衝突，編譯
卻壞在重複宣告——所以第 3 步的閘門跑在**合併後**，不是合併前。

### 不越界

MUST NOT: 改 `spec.md` 或 `.feature` 的例子。例子跟規則打架、規則有洞、少一個邊界——
記進票檔的 `## 缺口`，回報，然後
照現在的規格做。一個人在 IMPLEMENT 改掉驗收條件，等於沒有人同意過那條規格。

MUST NOT: 在 handler、mapper、presenter 裡放規則。判準：把這段刪掉，哪個 domain
單元測試會紅？沒有，它就不是規則，是接線；有，它放錯層了。

## 產物

```
<專案>/
├── internal/…                 這張票穿過的每一層
├── api/openapi.yaml           新的 operation（若 seam 是 HTTP）
├── test/acceptance/steps_*.go 新的 step definition，按領域名詞分檔
├── features/<story>.feature   只改第一行的狀態 tag
└── specs/<date>-<feature>/tasks/<票>.md   只勾驗收條件、寫缺口
```

MUST: 不修改 `spec.md`、`discovery.md`、`plan.md`、`tasks/index.md`、別張票的票檔、`.feature` 的任何場景文字。
同時做的幾個 session 各寫各的票檔才不會互相覆寫；設計的洞由 `bdd-tasks` 再跑時從票檔搬進
`plan.md`。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 寫程式前看過這張票的場景紅（`step is undefined` 或斷言失敗，不是 parse error），每條有分支的規則的 table-driven 單元測試也先看過紅
- [ ] 這張票的場景綠、專案的驗證閘門（實測場 `make verify`）綠，兩個命令與結果都在回報裡
- [ ] 寫新步驟前 grep 過既有的 step definition（實測場 `make steps`），沒有同一句話兩種寫法
- [ ] 測試替身回的是副本，不是活的指標
- [ ] handler、mapper、presenter 裡沒有規則
- [ ] `spec.md`、`discovery.md`、`plan.md`、`tasks/index.md`、別張票的票檔、`.feature` 的場景文字都沒有被改
- [ ] `.feature` 只有涵蓋它的票全勾完、整個檔的場景都綠了才改 `@review`，否則留 `@wip`
- [ ] 這張票是一個 commit，只含這張票動到的檔
- [ ] 在 worktree 裡做的：衝突照種類解（`產生` 重新產生、`追加` 取聯集），碰到 `爭用` 或沒列到的衝突停下來回報；閘門在合併後的程式碼上重跑過、綠了才合回
- [ ] 票檔只勾了有測試驗過的驗收條件
- [ ] 每個暫定行為都有一個測試釘住，註解寫明等哪一題的答案
- [ ] 這個 session 只做了一張票

## 完成後

告訴對方這幾件事：

1. 票號、涵蓋的場景幾個綠、閘門命令與結果
2. 新增了哪些檔、動了哪一層（一句話，不是檔案清單）
3. 票檔勾了哪幾條、哪一條沒勾與為什麼；票檔 `## 缺口` 記了什麼
4. 暫定的決定有哪些、各被哪個測試釘住
5. 下一張票是哪一張（`tasks/index.md` 裡 Blocked by 剛被解開的那張）；這個 `.feature` 還在 `@wip` 的話，等哪幾張票
6. 在 worktree 裡做的：合回時解了哪些衝突、各是哪一種，合併後的閘門命令與結果

全部票做完 → `.feature` 全部 `@review`，下一步 REVIEW（未實作：先人工對照
`spec.md` 的 AC 與 `.feature`，跑 `check_spec.py` 看鏈有沒有斷）。
