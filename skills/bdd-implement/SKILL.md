---
name: bdd-implement
description: >
  照 `plan.md` 一張票一個 session，outside-in 把票做到綠：先讓這張票的場景紅、
  domain 規則先寫 table-driven 單元測試、再接 use case／契約／step definition，
  場景全綠後把 `.feature` 從 `@ready` 翻到 `@review`、勾掉 `plan.md` 的驗收條件、
  跑專案的驗證閘門。不改 `prd.md`、不改 `.feature` 的例子、不在 handler 放規則。
  BDD 七步流程的 IMPLEMENT 步驟。
  觸發詞：「做第 N 張票」「開始實作」「bdd-implement 03」「把這張票做綠」
  「寫 step definition」「讓場景變綠」「照 plan 做」「下一張票」。
  English: implement ticket N, make the scenarios green, write the step
  definitions, outside-in from the feature file, work the plan.
---

# IMPLEMENT — 一張票，從紅到綠

一張 `plan.md` 的票進來（或者沒有 `plan.md` 時，一則 `.feature`），出去的是：那張票
涵蓋的場景全綠、專案的驗證閘門綠、`.feature` 標 `@review`、`plan.md` 的驗收條件打勾。

**外迴圈先紅，內迴圈先紅，然後才寫程式。** 外迴圈是 `.feature` 的場景（godog／
Cucumber），內迴圈是 domain 與 use case 的單元測試。兩圈都先看到紅，寫的每一行
才有理由。

不改 `prd.md`。不改 `.feature` 的例子。不做 `spec.md` 沒做的設計決定——發現有洞就
記下來，不順手補。

## 使用時機

- `plan.md` 寫好了，要做第一張（或下一張 Blocked by 全完成的）票
- 整批小到不需要 `bdd-plan`，直接從 `.feature` 與 `spec.md` 做
- 上一個 session 做到一半，要從 `@wip` 的檔接手

## Skill Boundaries

- 還沒有 `.feature` → `bdd-formulation`；還沒有 `spec.md`（seam、埠、契約）→ `bdd-spec`
- 票之間的順序、要不要再切 → `bdd-plan`
- 場景綠了之後「真的滿足它宣稱的情境嗎、鏈斷了嗎」→ REVIEW（`bdd-review`，未實作）
- 紅的是 bug 還是規格過期 → VERIFY（`bdd-verify`，未實作）；本 skill 只處理自己這張票
  讓紅變綠，不判讀別人的紅
- 要改規則或例子 → 回 `bdd-discovery` 開一題；本 skill 不碰 `prd.md`、不碰 `.feature` 的例子

## 參考檔案

- `../bdd-formulation/references/state-tags.md` — `@wip`／`@review` 什麼時候寫
- `../bdd-plan/references/plan-format.md` — 票的欄位，驗收條件為什麼要能失敗
- 專案自己的架構文件（實測場是 `docs/ARCHITECTURE.md`、`docs/DATAFLOW.md`、各層 `doc.go`）
  ——分層、轉換、錯誤分類都照它，本 skill 不另訂

---

## 前置確認（Ask user）

1. **哪一張票**：`bdd-implement <票號>`。沒帶票號就取 `plan.md` 裡第一張
   Blocked by 全部完成、驗收條件還沒全勾的票，講出來讓對方改
2. **驗證閘門是什麼命令**：讀專案的 README／Makefile（實測場是 `make verify`）。
   沒有就問——一張票「做完」的定義就是這個命令綠，不能猜
3. 讀得到 `specs/<date>-<feature>/plan.md`、`spec.md`、`prd.md` 與這張票涵蓋的 `.feature`

NEVER: `spec.md` 不存在時開始。seam、埠的形狀、錯誤分類都在那裡；沒有它，這一步會
自己設計，而設計出來的東西沒有任何場景驗得到。

MUST: 一個 session 一張票。撿起票的人沒看過 `spec.md` 以外的脈絡——這是 `bdd-plan`
切票的上界，也是這一步不要一次做三張的理由：做到第二張時第一張的脈絡已經在擠。

---

## 流程

### 1. 讀票，標 `@wip`

讀 `plan.md` 那一節：交付什麼行為、涵蓋哪些 `@example-N.M`、內迴圈要寫哪些單元測試、
驗收條件。再讀 `spec.md` 的 `## Implementation Decisions` 與 `## Testing Decisions`
——這張票動哪些模組、埠長什麼形狀、seam 在哪一層、每個例子在哪一層測。`prd.md` 只讀
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

`spec.md` 的 `## Testing Decisions` 說了哪些例子在 unit 層。對每一條有分支的規則：

1. 在 domain 寫 **table-driven** 單元測試，列就是 `prd.md` 的 AC 值加上 AC 沒寫但
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
`spec.md`，沒有就記進 `plan.md` 的 `## 缺口`。

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

- `.feature`：這張票做完後，**整個檔的場景都綠了**才從 `@wip` 改 `@review`；還有別張
  票涵蓋同一檔的場景就留 `@wip`
- `plan.md`：勾掉這張票的驗收條件。只勾驗過的——「手動驗」要寫是哪個測試驗的，
  沒有測試就不要勾
- 暫定行為（`spec.md` 的 `## 回 DISCOVERY 補問` 裡的洞，實作時不得不選一個）：
  **用一個測試釘住**，測試的註解寫明它是暫定、等哪一題的答案。沒釘住的暫定行為
  會在答案出來之後被當成規格

### 不越界

MUST NOT: 改 `prd.md` 或 `.feature` 的例子。例子跟規則打架、規則有洞、少一個邊界——
記進 `plan.md` 的 `## 缺口` 或 `spec.md` 的 `## 回 DISCOVERY 補問`，回報，然後
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
└── specs/<date>-<feature>/plan.md   只勾驗收條件、寫缺口
```

MUST: 不修改 `prd.md`、`spec.md`（除了讀）、`.feature` 的任何場景文字。

## 完成後

告訴對方五件事：

1. 票號、涵蓋的場景幾個綠、閘門命令與結果
2. 新增了哪些檔、動了哪一層（一句話，不是檔案清單）
3. `plan.md` 勾了哪幾條、哪一條沒勾與為什麼
4. 暫定的決定有哪些、各被哪個測試釘住
5. 下一張票是哪一張（`plan.md` 裡 Blocked by 剛被解開的那張）

全部票做完 → `.feature` 全部 `@review`，下一步 REVIEW（未實作：先人工對照
`prd.md` 的 AC 與 `.feature`，跑 `check_spec.py` 看鏈有沒有斷）。
