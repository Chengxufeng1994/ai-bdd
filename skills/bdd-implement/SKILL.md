---
name: bdd-implement
description: >
  照 `tasks/` 的票一張一個 session，outside-in 把票做到綠：先讓這張票的場景紅、
  domain 規則先寫 table-driven 單元測試、再接 use case／契約／step definition，
  場景全綠後把 `.feature` 從 `@ready` 翻到 `@review`、勾掉票檔的驗收條件、
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
- 場景綠了之後「真的滿足它宣稱的情境嗎、鏈斷了嗎」→ REVIEW（`bdd-review`，未實作）
- 紅的是 bug 還是規格過期 → VERIFY（`bdd-verify`，未實作）；本 skill 只處理自己這張票
  讓紅變綠，不判讀別人的紅
- 要改規則或例子 → 回 `bdd-discovery` 開一題；本 skill 不碰 `spec.md`、不碰 `.feature` 的例子

## 參考檔案

- `../bdd-formulation/references/state-tags.md` — `@wip`／`@review` 什麼時候寫
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

### 2. 外迴圈先紅

只跑這張票的場景，確認它們是**紅的、而且是對的紅**——`step is undefined` 或斷言失敗，
不是 parse error：

```bash
go test -count=1 ./test/acceptance/ -godog.tags=@rule-2          # 一條規則
go test -test.v ./test/acceptance/ -test.run "^TestFeatures/停留"  # 用名字選
```

MUST: 沒看到紅就不要寫程式。一個一寫就綠的場景沒有驗到任何東西——通常是 step
definition 已經被別的場景定義過，而它的斷言沒抓到這張票要做的事。

### 3. 內迴圈：規則先進 domain

`plan.md` 的 `## Testing Decisions` 說了哪些例子在 unit 層。對每一條有分支的規則：

1. 在 domain 寫 **table-driven** 單元測試，列就是 `spec.md` 的 AC 值加上 AC 沒寫但
   規則隱含的邊界（24 小時又 1 分鐘、到期日的最後一秒）——跑，看它因為「還沒有這個
   套件／函式」而紅
2. 寫最少的 domain 程式讓它綠。規則只住這裡；use case 與 handler 裡出現算式就是
   跑錯層
3. use case 的測試只驗**分類**：domain 回的錯變成哪個 `Kind`、哪個 message key、
   帶哪些 Params；repository 失敗變成什麼。用手刻 stub，一個埠一個 struct

MUST: 測試替身要回**副本**，不回活的指標。repository 的 stub 若把同一個聚合指標交出去，
use case 在 `Save` 失敗前做的修改會留在「儲存」裡，測不出「失敗時狀態不變」——
parking 實跑的第一個假紅就是這個。

### 4. 接線：照架構文件的轉換走一遍

依 `docs/DATAFLOW.md`（或專案的同類文件）：埠（`port/out`）→ use case →
driving port 與 service → contract（`openapi.yaml`）→ 產碼 → mapper → operation →
presenter → bootstrap 組裝。每一段都是機械翻譯，**不做決定**；要做的決定都該在
`plan.md` 的技術節，沒有就記進票檔的 `## 缺口`。

- contract 改了就重新產碼（實測場 `make gen-api`），產出的介面多了方法，`Server`
  不實作就編不過——那是對的，讓編譯器告訴你還缺哪個 operation
- 失敗路徑：operation 用 errmap 認得的 status 選 contract 宣告的 response 型別，
  其餘落 500；不要在 operation 裡 switch `Kind`
- 錯誤訊息的 message key 要有翻譯；場景斷言的 `錯誤訊息應為` 就是那句翻譯

### 5. step definitions：先找再寫

MUST: 寫新步驟之前 `grep` 既有的 step definition（實測場 `make steps`）。同一句話兩種
寫法就是兩支函式，而它們會分岔。

- 一個樣板一支函式；`操作成功`／`操作失敗`／`應付金額應為` 這類跨 feature 的句子
  一定已經有了
- 每個場景的狀態住 world（`context.Context` 裡），不住 package 變數——並行跑的套件
  一個共用變數就是間歇性紅
- 時間、ID、會失敗的 repository 這類「場景要控制的東西」用測試端的替身注入
  bootstrap，走**跟 binary 同一個組裝函式**；不要為測試另開一條接線
- 播種步驟直接寫 store，不經 use case——場景要擺出的狀態（付過款、票已過期）
  不一定是 use case 走得到的

### 6. 外迴圈綠，閘門綠

```bash
go test -count=1 ./test/acceptance/ -godog.tags=@rule-2   # 這張票的場景
make verify                                               # 專案的全部閘門
```

MUST: 兩個都綠才算做完。閘門包含 fmt、產碼是否過期、lint、`-race`——其中任何一項紅
而只回報「場景綠了」，是這一步最常見的「完成」謊言。

跑閘門時產碼工具會重跑（mockery 之類），新的埠要讓它產一次，否則「generated code
is stale」。

### 7. 收尾

- `.feature`：在主線上做時，這張票做完後**整個檔的場景都綠了**才從 `@wip` 改 `@review`；
  還有別張票涵蓋同一檔的場景就留 `@wip`。在 worktree 裡做時一律留 `@wip`，到〈合回主線〉
  第 4 步再判
- 票檔：勾掉這張票的驗收條件。只勾驗過的——「手動驗」要寫是哪個測試驗的，
  沒有測試就不要勾
- 暫定行為（票檔或 `plan.md` 的 `## 缺口`、或 `spec.md` 的 `## 回 DISCOVERY 補問` 裡的洞，實作時不得不選一個）：
  **用一個測試釘住**，測試的註解寫明它是暫定、等哪一題的答案。沒釘住的暫定行為
  會在答案出來之後被當成規格

### 8. 合回主線（只在 worktree 裡做時）

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
