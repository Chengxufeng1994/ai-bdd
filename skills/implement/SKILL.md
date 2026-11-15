---
name: implement
description: >
  照 `tasks/` 的票一張一個 session，outside-in 把票做到綠：先把這張票的場景接上 step
  definition、看到對的紅，再把 green 與 refactor 委派給 `bdd`（規則先進 domain、table-driven
  單元測試先紅），場景全綠後把 `.feature` 從 `@ready` 翻到 `@review`、勾掉票檔的驗收條件、
  跑專案的驗證閘門。不改 `spec.md`、不改 `.feature` 的例子、不在 handler 放規則。
  BDD 六步流程的 IMPLEMENT 步驟。
  觸發詞：「做第 N 張票」「開始實作」「implement 03」「把這張票做綠」
  「寫 step definition」「讓場景變綠」「照 plan 做」「下一張票」。
  English: implement ticket N, make the scenarios green, write the step
  definitions, outside-in from the feature file, work the plan.
---

# IMPLEMENT — 一張票，從紅到綠

| | |
| --- | --- |
| **職責** | 把一張票做到綠：測試層先紅、委派 `bdd` 轉綠與整理、閘門綠、勾驗收條件、合回主線 |
| **讀** | `tasks/<票>.md`、`tasks/index.md`、`spec.md`（技術部分與這張票的 FR／AC）、`plan.md`、這張票涵蓋的 `.feature`、專案的架構文件 |
| **寫** | 產品程式碼與 step definition、這張票的票檔（勾選、`## 缺口`）、`.feature` 的狀態 tag（`@wip`／`@review`） |
| **上游** | `task-breakdown`（票切好之後） |
| **下游** | `bdd`（第 3、4 步委派）；做完之後是 REVIEW（未實作） |
| **不做** | 不改 `spec.md`、不改 `.feature` 的例子、不寫 `plan.md` 與 `tasks/index.md`、不做 `spec.md` 技術部分沒做的設計決定 |

**外迴圈先紅，內迴圈先紅，然後才寫程式。** 外迴圈是 `.feature` 的場景（godog／
Cucumber），內迴圈是 domain 與 use case 的單元測試。兩圈都先看到紅，寫的每一行
才有理由。票裡的轉綠與整理委派給 `bdd`——它一次一個場景、一次一個入口。發現技術部分的洞
就記下來，不順手補。

## 使用時機

- `tasks/` 切好了，要做第一張（或下一張 Blocked by 全完成的）票
- 幾張票同時做，每張各開一個 session 與 worktree
- 上一個 session 做到一半，要從 `@wip` 的檔接手

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| 還沒有 `.feature` | `formulation` |
| `spec.md` 沒有技術部分（seam、埠、契約） | `spec` 第 6 步 |
| 還沒有 `tasks/index.md`（包括票還寫在 `plan.md` `## 票` 裡的舊格式） | `task-breakdown` |
| 票之間的順序、要不要再切、哪幾張能同時做 | `task-breakdown` |
| 只推進**一個場景**的 red／green／refactor 其中一步，而不是做一張票 | 使用者手動叫 `/bdd`；本 skill 做票時第 3、4 步本來就委派它 |
| 場景綠了之後「真的滿足它宣稱的情境嗎、鏈斷了嗎」 | REVIEW（`review`，未實作） |
| 紅的是 bug 還是規格過期 | VERIFY（`verify`，未實作）；本 skill 只處理自己這張票讓紅變綠，不判讀別人的紅 |
| 要改規則或例子 | 回 `discovery` 開一題；本 skill 不碰 `spec.md`、不碰 `.feature` 的例子 |

## 參考檔案

- `../formulation/references/state-tags.md` — `@wip`／`@review` 什麼時候寫
- `../bdd/SKILL.md` — 第 3、4 步委派的 green／refactor 入口；`../bdd/rules/red.md` 與
  `../bdd/references/step-definitions.md` — 第 2 步測試層先紅時照它寫
- `../task-breakdown/references/tasks-format.md` — 票的欄位，驗收條件為什麼要能失敗，接合點與同時跑
- `../spec/references/tech-format.md` — `spec.md` 技術部分各節裝什麼
- 專案自己的架構文件（實測場是 `docs/ARCHITECTURE.md`、`docs/DATAFLOW.md`、各層 `doc.go`）
  ——分層、轉換、錯誤分類都照它，本 skill 不另訂

---

## 前置確認（Ask user）

1. **哪一張票**：`implement <票號>`。沒帶票號就取 `tasks/index.md` 裡第一張
   Blocked by 全部完成（那些票檔的驗收條件全勾）、自己還沒全勾的票，講出來讓對方改
2. **驗證閘門是什麼命令**：票檔 `**驗證**` 寫了就用它；寫「（問使用者）」或沒寫，讀 `spec.md`
   `## Commands` 的閘門，再讀專案的 README／Makefile（實測場是 `make verify`），再沒有就問——一張票「做完」的定義就是這個命令綠，不能猜
3. **前面有沒有檢查點**：`index.md` 裡這張票的上方有 `**檢查點**` 時，看票檔有沒有記 `**檢查點**:`
   那一行。沒有、也沒有人在開工說明裡帶來，就在主線上照檢查點跑一次（已勾完的票的場景、閘門），
   結果給使用者看，使用者點頭才開工，並把那一行寫進這張票檔。一波只跑一次
   （`../task-breakdown/references/tasks-format.md`〈檢查點〉）
4. 讀得到 `specs/ongoing/<date>-<feature>/tasks/` 的這張票檔、`plan.md`、`spec.md` 與這張票涵蓋的 `.feature`

NEVER: `spec.md` 沒有技術部分時開始。seam、埠的形狀、錯誤分類都在那裡；沒有它，這一步會
自己設計，而設計出來的東西沒有任何場景驗得到。

MUST: 一個 session 一張票。撿起票的人只看過票檔與 `spec.md` 的技術部分——這是 `task-breakdown`
切票的上界，也是這一步不要一次做三張的理由：做到第二張時第一張的脈絡已經在擠。

SHOULD: 跟別張票同時做時，一張票一個 worktree。開工前看 `tasks/index.md`：這張票跟正在做
的票有共同的 `爭用` 接合點、或兩張都要新建同一個東西，就不要同時開——它們合不回來。

MUST: 開工時工作區是乾淨的；做完這張票一個 commit，只 stage 這張票動到的檔。合回與回滾都
以這個 commit 為單位。

---

## 流程

| 步驟 | 做什麼 | 誰做 | 要問使用者嗎 |
| --- | --- | --- | --- |
| 1 | 讀票，標 `@wip` | 本 skill | — |
| 2 | 測試層先紅：場景都接上 step definition，沒有 undefined step | 本 skill | — |
| 3 | green：一個場景一個場景轉綠 | 委派 `bdd` | — |
| 4 | refactor：在綠燈下整理 | 委派 `bdd` | — |
| 5 | 外迴圈綠，閘門綠 | 本 skill | — |
| 6 | 收尾：勾驗收條件、寫缺口、判斷 `@review` | 本 skill | — |
| 7 | 合回主線（只在 worktree 裡做時） | 本 skill | 碰到 `爭用` 或沒列到的衝突時停下回報 |

### 1. 讀票，標 `@wip`

讀票檔 `tasks/<票號>-<slug>.md`：交付什麼行為、涵蓋哪些 `@example-N.M`、內迴圈要寫哪些單元測試、
驗收條件。再讀 `spec.md` 技術部分的 `## Implementation Decisions`、`## Testing Decisions`、
`## Project Structure`、`## Code Style`——這張票動哪些模組、埠長什麼形狀、seam 在哪一層、每個 AC
在哪一層測、程式碼放哪、長什麼樣。需求部分只讀對應的 FR 與 AC，看規則的原話。

把這張票涵蓋的 `.feature` 檔從 `@ready` 改成 `@wip`（`state-tags.md`）；已經是 `@wip`（別張涵蓋
同一檔的票先開工了）就不動。下一個 session 看到 `@wip` 就知道有人做到一半。

### 2. 測試層先紅（本 skill 自己做，不委派）

把這張票涵蓋的每個場景都接上 step definition，直到**沒有任何 undefined step**、失敗只能是
斷言或產品行為。寫步驟照 `../bdd/references/step-definitions.md`；判斷紅得對不對照
`../bdd/rules/red.md`〈什麼是「對的紅」〉。只寫測試層，不寫產品程式碼——播種步驟需要的空型別、
回傳零值的空殼照 `red.md`〈R3：最小通路的界線〉寫，不算產品程式碼。

```bash
go test -count=1 ./test/acceptance/ -v -godog.format=progress -godog.tags=@rule-2   # 一條規則
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
| 範疇（Test Scope） | 這張票檔的 `涵蓋場景`，組成一個 tag 運算式（實測場 godog 寫成 `@example-8.1,@example-8.2`） |
| 票 | `tasks/<票號>-<slug>.md` 的路徑 |
| 設計 | `spec.md` 的 `## Implementation Decisions` 與 `## Testing Decisions` |

`bdd` 一個場景一個場景轉綠：規則先進 domain、先看到單元測試紅（`../bdd/references/inner-loop.md`）。
它停下來回報規格的洞或範疇外的耦合時，本 skill 把那一條記進票檔的 `## 缺口`，照
〈不越界〉處理——`bdd` 自己不寫 `specs/` 底下任何檔案。

### 4. 委派 `bdd` 的 refactor

範疇全綠之後，同樣的範疇交給 `bdd` 的 `refactor` 入口：在綠燈保護下整理這張票留下的重複
步驟、轉換與支撐層。沒有要整理的就跳過，在回報裡寫「無」。

### 5. 外迴圈綠，閘門綠

跑票檔 `**驗證**` 寫好的命令，再加一行前面的票：

```bash
<票檔驗證：這張票的場景>      # 例 go test -count=1 ./test/acceptance/ -v -godog.format=progress -godog.tags='@example-1.1,@example-1.2'
<已勾完的票的場景>            # 前面的票沒有被改壞
<票檔驗證：閘門>              # 例 make verify
```

第二行把 `index.md` 列出的、驗收條件已全勾的票，各自票檔 `**驗證**` 裡的場景運算式接成一個
「或」（實測場 godog 用逗號）。沒有場景運算式的票（前置整理）跳過；還沒有這種票就整行省略。

MUST: 每次跑場景都看輸出的**場景數**，並寫進回報。運算式寫錯時 runner 常常選到 0 個場景、
照樣成功結束——實測場的 godog 對 `@a || @b` 與 `@a or @b` 就是這樣。0 個場景的綠是假綠。票檔沒有 `**驗證**`（舊格式）時，從 `涵蓋場景` 自己組，並在
`## 缺口` 記一列。

MUST: 兩個都綠才算做完。閘門包含 fmt、產碼是否過期、lint、`-race`——其中任何一項紅
而只回報「場景綠了」，是這一步最常見的「完成」謊言。

跑閘門時產碼工具會重跑（mockery 之類），新的埠要讓它產一次，否則「generated code
is stale」。

### 6. 收尾

- `.feature`：在主線上做時，這張票做完後**整個檔的場景都綠了**才從 `@wip` 改 `@review`；
  還有別張票涵蓋同一檔的場景就留 `@wip`。在 worktree 裡做時一律留 `@wip`，到〈合回主線〉
  第 4 步再判
- 票檔：勾掉這張票的驗收條件與 `**驗證**`。只勾驗過的——「手動驗」要寫是哪個測試驗的，
  沒有測試就不要勾；`**驗證**` 的手動項要使用者看過才勾
- 暫定行為（票檔或 `plan.md` 的 `## 缺口`、`spec.md` 的 `## 回 DISCOVERY 補問` 裡的洞、`.feature` 裡
  「沒問過」的註解，以及實作時才發現、規格沒說的邊界——實作時不得不選一個）：先記進票檔的
  `## 缺口`，選一個，**用一個測試釘住**，測試的註解寫明它是暫定、等哪一題的答案。沒釘住的暫定行為
  會在答案出來之後被當成規格
- 爭用接合點：這張票的票檔列了 `爭用` 的檔時，核對那段邏輯實際落在哪個檔。落在別的檔，
  在票檔 `## 缺口` 記一列（列的是哪個檔、實際在哪個檔）——`task-breakdown` 再跑時據此改正接合點

### 7. 合回主線（只在 worktree 裡做時）

這張票由做它的這個 session 自己合回，不交給別人——衝突要靠這張票的脈絡才解得對。

1. **先把主線拉進 worktree**：在 worktree 裡 `git rebase main`，衝突在這裡解，不在主線上解。
   用 rebase 不用 merge——merge 會多一個反方向的 merge commit，打破「一張票一個 commit」
2. **照接合點的種類解**（`../task-breakdown/references/tasks-format.md`〈接合點與同時跑〉）：
   - `產生` → 不手動合：先解完原始檔，再跑產碼工具重新產生（實測場 `make gen`）；`go.sum`
     用 `go mod tidy`
   - `追加` → 取聯集：兩邊加的都留，然後檢查有沒有重複宣告
   - 票檔與 `index.md` 都沒列到、但兩邊**只各自新增頂層宣告**（新函式、新型別，沒有改到任何既有的行）
     → 當成 `追加` 解，並在票檔 `## 缺口` 記一列「沒列到的接合點：<檔>」，讓 `task-breakdown` 補上
   - `爭用`、或票檔與 `index.md` 都沒列到的其他衝突 → **停下來回報**，不自己判——那是切票時
     沒預料到的耦合，要回 `task-breakdown` 重排。停下來時：
     - 撤掉這次合併（`git rebase --abort`），worktree 回到這張票自己的 commit，不留半解的衝突
     - 在票檔 `## 缺口` 記一列：衝突的檔、跟哪張票撞、為什麼不是 `產生`／`追加`；`amend` 進這張票的
       commit，保持一張票一個 commit
     - 「列到」只看票檔與 `index.md` 的 `接合點` 欄；先前記在 `## 缺口` 的預感不算列到
     - 告訴使用者跑 `task-breakdown`：它把這一列搬進 `plan.md` 的 `## 缺口`，`plan` 補上
       接合點的種類後再重排；重排完，這張票從本步驟（第 7 步）的第 1 點重來，不必重做第 1–6 步
3. **在合併後的程式碼上重跑閘門**，場景範圍是主線上所有驗收條件已勾完的票加上這張票
   （跟第 5 步第二行同一種組法：各票檔 `**驗證**` 的場景運算式接成一個「或」，前置整理的票跳過）。
   綠了才合回主線：回到主線 `git merge --ff-only <票的 branch>`
   - 合完 git 沒報衝突、閘門卻紅，常見兩種：兩張票在同一個測試 package 各自新建的檔宣告了同名的
     輔助函式或變數 → 改**這張票自己新建的檔**裡的名字；另一張票新加了必填的組裝欄位，這張票新建的
     測試沒帶 → 照〈不越界〉補上。兩種都在票檔 `## 缺口` 記一列，`task-breakdown` 補成接合點
4. **判 `@review`**：照 `index.md` 列的票，逐張讀票檔的 `涵蓋場景`，找出涵蓋這個 `.feature` 的票，
   再看它們的驗收條件——全勾了，就是這次合回讓最後一張完成，**由這次合回翻 `@review`**：
   fast-forward 帶不進新的改動，所以在主線上另開一個 commit 只翻 tag（這是「一張票一個 commit」唯一的
   例外）；還有沒勾的就留 `@wip`
5. 先準備好的先合；同時準備好的，票號小的先合。後合的那張回到本步驟第 1 點

IMPORTANT: git 合得乾淨不等於合對了。兩張票各自新建同一個 package，git 不會報衝突，編譯
卻壞在重複宣告——所以第 3 步的閘門跑在**合併後**，不是合併前。

### 不越界

MUST NOT: 改 `spec.md` 或 `.feature` 的例子。例子跟規則打架、規則有洞、少一個邊界——
記進票檔的 `## 缺口`，回報，然後
照現在的規格做。一個人在 IMPLEMENT 改掉驗收條件，等於沒有人同意過那條規格。

前面的票寫的程式碼可以改——這張票的場景需要時，就改。條件是那些票的場景照樣綠：第 5 步的
閘門要跑主線上所有驗收條件已勾完的票的場景，不只這張。改不動、或改了會讓前面的票的
`spec.md` 技術決定不成立 → 記進票檔的 `## 缺口` 並回報，不要改掉那張票的驗收條件或票檔。

MUST NOT: 在 handler、mapper、presenter 裡放規則。判準：把這段刪掉，哪個 domain
單元測試會紅？沒有，它就不是規則，是接線；有，它放錯層了。

## 產物

```
<專案>/
├── internal/…                 這張票穿過的每一層
├── api/openapi.yaml           新的 operation（若 seam 是 HTTP）
├── test/acceptance/steps_*.go 新的 step definition，按領域名詞分檔
├── features/<story>.feature   只改第一行的狀態 tag
└── specs/ongoing/<date>-<feature>/tasks/<票>.md   只勾驗收條件、寫缺口
```

MUST: 不修改 `spec.md`、`discovery.md`、`plan.md`、`tasks/index.md`、別張票的票檔、`.feature` 的任何場景文字。
同時做的幾個 session 各寫各的票檔才不會互相覆寫；設計的洞由 `task-breakdown` 再跑時從票檔搬進
`plan.md` 的 `## 缺口`，再交回 `spec`。

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
- [ ] 跑的是票檔 `**驗證**` 的命令，另跑了已勾完的票的場景；這張票前面有檢查點時，使用者看過檢查點的結果才開工
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
