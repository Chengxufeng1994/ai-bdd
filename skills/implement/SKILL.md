---
name: implement
description: >
  照 `tasks.md` 的票一張一個 session，outside-in 把票做到綠：先把這張票的場景接上 step
  definition、看到對的紅，再把 green 與 refactor 委派給 `behavior-driven-development`（規則先進 domain、table-driven
  單元測試先紅），場景全綠後把 `.feature` 從 `@ready` 翻到 `@review`、勾掉票節的驗收條件、
  跑專案的驗證閘門。不改 `spec.md`、不改 `.feature` 的例子、不在 handler 放規則。
  BDD 六步流程的 IMPLEMENT 步驟。
  觸發詞：「做第 N 張票」「開始實作」「implement 03」「把這張票做綠」
  「寫 step definition」「讓場景變綠」「照 plan 做」「下一張票」。
  English: implement ticket N, make the scenarios green, write the step
  definitions, outside-in from the feature file, work the plan.
disable-model-invocation: true
---

# IMPLEMENT — 一張票，從紅到綠

| | |
| --- | --- |
| **職責** | 把一張票做到綠：測試層先紅、委派 `behavior-driven-development` 轉綠與整理、閘門綠、勾驗收條件、合回主線 |
| **讀** | `tasks.md`（開頭三行與這張票的 `### <票號>` 票節）、`spec.md`（這張票的 FR／AC）、`plan.md`（技術決定、三份設計與執行計畫）、這張票涵蓋的 `.feature`、專案的架構文件 |
| **寫** | 產品程式碼與 step definition、`api/openapi.yaml`（照抄 `plan.md` 的片段）與 migration、`tasks.md` 裡這張票的票節（勾選、`**檢查點**:`、`#### 缺口`）、`.feature` 的狀態 tag（`@wip`／`@review`） |
| **上游** | `plan`（票切好之後） |
| **下游** | `behavior-driven-development`（第 3、4 步委派）；全部票做完之後是 `archive` |
| **不做** | 不改 `spec.md`、不改 `.feature` 的例子、不寫 `plan.md`、不改 `tasks.md` 裡這張票的票節以外的任何一行、不做 `plan.md` 的技術決定沒做的設計決定 |

**外迴圈先紅，內迴圈先紅，然後才寫程式。** 外迴圈是 `.feature` 的場景（godog／
Cucumber），內迴圈是 domain 與 use case 的單元測試。兩圈都先看到紅，寫的每一行
才有理由。票裡的轉綠與整理委派給 `behavior-driven-development`——它一次一個場景、一次一個入口。發現技術決定的洞
就記下來，不順手補。

## 使用時機

- `tasks.md` 切好了，要做第一張（或下一張 Blocked by 全完成的）票
- 幾張票同時做，每張各開一個 session 與 worktree
- 上一個 session 做到一半，要從 `@wip` 的檔接手

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| 還沒有 `.feature` | `formulation` |
| `plan.md` 沒有技術決定（seam、埠、契約） | `plan` |
| 要改 operation、`components` 或表的形狀 | `plan`（改 `## API 契約`／`## Database Schema`）——本 skill 只照抄 |
| 還沒有 `tasks.md`（包括票還寫在 `plan.md` `## 票` 裡、或還是 `tasks/` 目錄的舊格式） | `plan` |
| 票之間的順序、要不要再切、哪幾張能同時做 | `plan` |
| 只推進**一個場景**的 red／green／refactor 其中一步，而不是做一張票 | 使用者手動叫 `/behavior-driven-development`；本 skill 做票時第 3、4 步本來就委派它 |
| 場景綠了之後「真的滿足它宣稱的情境嗎、鏈斷了嗎」 | REVIEW（`review`，未實作） |
| 紅的是 bug 還是規格過期 | VERIFY（`verify`，未實作）；本 skill 只處理自己這張票讓紅變綠，不判讀別人的紅 |
| 要改規則或例子 | 回 `discovery` 開一題；本 skill 不碰 `spec.md`、不碰 `.feature` 的例子 |

## 參考檔案

- `../formulation/references/state-tags.md` — `@wip`／`@review` 什麼時候寫
- `../behavior-driven-development/SKILL.md` — 第 3、4 步委派的 green／refactor 入口；`../behavior-driven-development/rules/red.md` 與
  `../behavior-driven-development/references/step-definitions.md` — 第 2 步測試層先紅時照它寫
- `../plan/references/tasks-format.md` — 票的欄位，驗收條件為什麼要能失敗，接合點與同時跑
- `../plan/references/tech-format.md` — `plan.md` 的技術決定各節裝什麼
- `../plan/references/design-format.md` — `## Domain Model`、`## Database Schema`、`## API 契約` 裝什麼
- 專案自己的架構文件（實測場是 `docs/ARCHITECTURE.md`、`docs/DATAFLOW.md`、各層 `doc.go`）
  ——分層、轉換、錯誤分類都照它，本 skill 不另訂

---

## 前置確認（Ask user）

1. **哪一張票**：`implement <票號>`。沒帶票號就取 `tasks.md` 裡第一張
   Blocked by 全部完成（那些票節的驗收條件全勾）、自己還沒全勾的票，講出來讓對方改
2. **驗證閘門是什麼命令**：票節 `**驗證**` 寫了就用它；寫「（問使用者）」或沒寫，讀 `plan.md`
   `## Commands` 的閘門，再讀專案的 README／Makefile（實測場是 `make verify`），再沒有就問——一張票「做完」的定義就是這個命令綠，不能猜
3. **前面有沒有檢查點**：`tasks.md` 裡這張票的上方有 `**檢查點**` 時，看票節有沒有記 `**檢查點**:`
   那一行。沒有、也沒有人在開工說明裡帶來，就在主線上照檢查點跑一次（已勾完的票的場景、閘門），
   結果給使用者看，使用者點頭才開工，並把那一行寫進這張票節。一波只跑一次
   （`../plan/references/tasks-format.md`〈檢查點〉）
4. 讀得到 `specs/ongoing/<date>-<feature>/tasks.md` 的這張票節、`plan.md`、`spec.md` 與這張票涵蓋的 `.feature`

NEVER: `plan.md` 沒有技術決定時開始。seam、埠的形狀、錯誤分類都在那裡；沒有它，這一步會
自己設計，而設計出來的東西沒有任何場景驗得到。

Check if: 舊格式——`plan.md` 沒有 `## Domain Model`，或 `spec.md` 還有 `## Domain Model` 或 `### 聚合與不變條件`。是就停下，請使用者
重跑 `/plan`（`spec.md` 的舊節由再跑 `/spec` 刪掉；`plan` 照〈再跑〉重切 `tasks.md`）——舊格式沒有可以照抄的契約與資料表。

MUST: 一個 session 一張票。撿起票的人只看過票節與 `plan.md` 的技術決定——這是 `plan` 的 Step 3
切票的上界，也是這一步不要一次做三張的理由：做到第二張時第一張的脈絡已經在擠。

SHOULD: 跟別張票同時做時，一張票一個 worktree。開工前看 `tasks.md`：這張票跟正在做
的票有共同的 `爭用` 接合點、或兩張都要新建同一個東西，就不要同時開——它們合不回來。

MUST: 開工時工作區是乾淨的；做完這張票一個 commit，只 stage 這張票動到的檔。合回與回滾都
以這個 commit 為單位。

---

## 流程

| 步驟 | 做什麼 | 誰做 | 要問使用者嗎 |
| --- | --- | --- | --- |
| 1 | 讀票，標 `@wip` | 本 skill | — |
| 2 | 測試層先紅：場景都接上 step definition，沒有 undefined step | 本 skill | — |
| 3 | green：一個場景一個場景轉綠 | 委派 `behavior-driven-development` | — |
| 4 | refactor：在綠燈下整理 | 委派 `behavior-driven-development` | — |
| 5 | 外迴圈綠，閘門綠 | 本 skill | — |
| 6 | 收尾：勾驗收條件、寫缺口、判斷 `@review` | 本 skill | — |
| 7 | 合回主線（只在 worktree 裡做時） | 本 skill | 碰到 `爭用` 或沒列到的衝突時停下回報 |

### 1. 讀票，標 `@wip`

讀 `tasks.md` 的 `### <票號>` 票節：交付什麼行為、涵蓋哪些 `@example-N.M`、內迴圈要寫哪些單元測試、
驗收條件。再讀 `plan.md` 的 `## Implementation Decisions`、`## Testing Decisions`、
`## Project Structure`、`## Code Style`——這張票動哪些模組、埠長什麼形狀、seam 在哪一層、每個 AC
在哪一層測、程式碼放哪、長什麼樣——以及三份設計裡這張票碰的部分：`## Domain Model` 的聚合、不變條件、
票節 `**資料表**` 在 `## Database Schema` 的那幾張表、`**契約**` 在 `## API 契約` 的那幾個 operation。需求部分只讀對應的 FR 與 AC，看規則的原話。

再讀 `Blocked by` 那幾張票的票節 `#### 缺口`：點名這張票的列（「留給 03」「到期日碰巧不是休館日」）是前一個
session 的猜測，當成**待驗證的線索**——用這張票的第一個紅驗證它，不直接信。喚醒之後切出來的票
（改的是已經做出來的規則），改讀**跟這張票的接合點或模組重疊的所有舊票**的 `#### 缺口`：被推翻的暫定
行為常記在一張不在 `Blocked by` 裡的票上。

**沒有 `涵蓋場景` 的票**（前置整理）：沒有 `.feature` 要標 `@wip`，第 2 步沒有場景可紅。紅用驗收條件
寫的單元測試；第 3、4 步不委派 `behavior-driven-development`，由本 skill 照同樣的紅→綠→整理自己做；第 5 步只跑閘門。

把這張票涵蓋的 `.feature` 檔從 `@ready` 改成 `@wip`（`state-tags.md`）；已經是 `@wip`（別張涵蓋
同一檔的票先開工了）就不動。下一個 session 看到 `@wip` 就知道有人做到一半。

### 2. 測試層先紅（本 skill 自己做，不委派）

把這張票涵蓋的每個場景都接上 step definition，直到**沒有任何 undefined step**、失敗只能是
斷言或產品行為。寫步驟照 `../behavior-driven-development/references/step-definitions.md`；判斷紅得對不對照
`../behavior-driven-development/rules/red.md`〈什麼是「對的紅」〉。只寫測試層，不寫產品程式碼——播種步驟需要的空型別、
回傳零值的空殼、以及讓請求打得到 seam 的接線，照 `red.md`〈R3：最小通路的界線〉寫，不算產品程式碼。
seam 是 HTTP 時，契約的 operation、重新產碼、**回錯誤**（不是回成功）的 handler 空殼都算最小通路——
停在 404 只證明路由不存在，不是這張票的紅。

**契約從 `plan.md` 搬進來。** 票節有 `**契約**`、`plan.md` 的 `## API 契約` 是 OpenAPI 片段時：把那幾個 operation
與它們用到的 `components`，從片段照抄進 `api/openapi.yaml`，再重新產碼。某個 component 前面的票已經搬進
`api/openapi.yaml`：跟片段一樣 → 留著，不再搬一次；不一樣 → 記進票節 `#### 缺口`、停下回報，交回 `plan`。`api/openapi.yaml` 是 `追加` 接合點，
合回照〈合回主線〉。`## API 契約` 是程式碼層契約（沒有 OpenAPI 片段）時沒有要搬的，跳過這段，埠照它的方法寫。

MUST NOT: 搬的時候改形狀——欄位、型別、必填、狀態碼、錯誤碼都跟片段一樣。要改，記進票節 `#### 缺口`、停下回報：
`plan` 再跑時在設計裡補上（改好片段、記進 `plan.md` 的 `## Further Notes`）再搬；屬於需求的才進 `## 缺口` 交回 `spec`。在 `api/openapi.yaml` 改了，`plan.md` 與
專案的契約從此分家，而沒有人發現。

票節有 `**資料表**` 時，migration 照 `plan.md` `## Database Schema` 寫：欄位、型別、約束、索引都一樣，
`守住哪條不變條件` 那一欄不進 migration。它是產品程式碼，在第 3 步 repository 需要它時寫；要改表的形狀同上，回 `plan`。

```bash
go test -count=1 ./test/acceptance/ -v -godog.format=progress -godog.tags=@rule-2   # 一條規則
go test -test.v ./test/acceptance/ -test.run "^TestFeatures/借書"  # 用名字選
```

MUST: 沒看到紅就不要寫程式。一個一寫就綠的場景沒有驗到任何東西——通常是 step
definition 已經被別的場景定義過，而它的斷言沒抓到這張票要做的事。

這張票所有場景都紅得對了，才進第 3 步。測試層沒對齊就改產品程式碼，綠了也說不出是哪個
斷言在保護它。唯一的例外是**零值例子**（定義見 `red.md` R5：在起點就綠、或只紅在共同的接線失敗）：
它可以帶著綠進第 3 步，但轉綠之後要照 R5 做突變驗證，回報裡寫出改壞了什麼、它怎麼變紅。

### 3. 委派 `behavior-driven-development` 的 green

讀 `../behavior-driven-development/SKILL.md`，照它的 `green` 入口做，範疇與輸入由本 skill 給：

| 交給 `behavior-driven-development` 的 | 內容 |
| --- | --- |
| 入口 | `green` |
| 範疇（Test Scope） | 這張票節的 `涵蓋場景`，組成一個 tag 運算式（實測場 godog 寫成 `@example-8.1,@example-8.2`） |
| 票 | `tasks.md` 的路徑與票號（`### <票號>`） |
| 設計 | `plan.md` 的 `## Implementation Decisions`、`## Testing Decisions`、`## Domain Model`，與票節 `**資料表**` 在 `## Database Schema` 的表 |

`behavior-driven-development` 一個場景一個場景轉綠：規則先進 domain、先看到單元測試紅（`../behavior-driven-development/references/inner-loop.md`）。
它停下來回報規格的洞或範疇外的耦合時，本 skill 把那一條記進票節的 `#### 缺口`，照
〈不越界〉處理——`behavior-driven-development` 自己不寫 `specs/` 底下任何檔案。

### 4. 委派 `behavior-driven-development` 的 refactor

範疇全綠之後，同樣的範疇交給 `behavior-driven-development` 的 `refactor` 入口：在綠燈保護下整理這張票留下的重複
步驟、轉換與支撐層。沒有要整理的就跳過，在回報裡寫「無」。

### 5. 外迴圈綠，閘門綠

跑票節 `**驗證**` 寫好的命令，再加一行前面的票：

```bash
<票節驗證：這張票的場景>      # 例 go test -count=1 ./test/acceptance/ -v -godog.format=progress -godog.tags='@example-1.1,@example-1.2'
<已勾完的票的場景>            # 前面的票沒有被改壞
<票節驗證：閘門>              # 例 make verify
```

第二行把 `tasks.md` 列出的、驗收條件已全勾的票，各自票節 `**驗證**` 裡的場景運算式接成一個
「或」（實測場 godog 用逗號）。沒有場景運算式的票（前置整理）跳過；還沒有這種票就整行省略。

MUST: 每次跑場景都看輸出的**場景數**，並寫進回報。運算式寫錯時 runner 常常選到 0 個場景、
照樣成功結束——實測場的 godog 對 `@a || @b` 與 `@a or @b` 就是這樣。0 個場景的綠是假綠。票節沒有 `**驗證**`（舊格式）時，從 `涵蓋場景` 自己組，並在
`#### 缺口` 記一列。

MUST: 兩個都綠才算做完。閘門包含 fmt、產碼是否過期、lint、`-race`——其中任何一項紅
而只回報「場景綠了」，是這一步最常見的「完成」謊言。閘門裡跑場景的那一段同樣要看得到場景數；
看不到就另外帶旗標重跑一次。

跑閘門時產碼工具會重跑（mockery 之類），新的埠要讓它產一次，否則「generated code
is stale」。

### 6. 收尾

- `.feature`：在主線上做時，這張票做完後**整個檔的場景都綠了**（`plan.md` `## Commands` 的
  `跑單一 .feature`，看場景數）才從 `@wip` 改 `@review`；
  還有別張票涵蓋同一檔的場景就留 `@wip`。在 worktree 裡做時一律留 `@wip`，到〈合回主線〉
  第 4 步再判
- 票節：勾掉這張票的驗收條件與 `**驗證**`。只勾驗過的——「手動驗」要寫是哪個測試驗的，
  沒有測試就不要勾；`**驗證**` 的手動項要使用者看過才勾。在 worktree 裡做時，`**驗證**` 的閘門那一格
  等〈合回主線〉第 3 步在合併後的程式碼上綠了才勾，跟合回前最後一次 `amend` 一起進這張票的 commit
- 暫定行為（票節或 `plan.md` 的 `## 缺口`、`spec.md` 的 `## Open Questions` 裡的洞、`.feature` 裡
  「沒問過」的註解，以及實作時才發現、規格沒說的邊界——實作時不得不選一個）：先記進票節的
  `#### 缺口`，選一個，**用一個測試釘住**，測試的註解寫明它是暫定、等哪一題的答案；沒有 Q 號時引
  `.feature` 的「沒問過」註解或票節 `#### 缺口` 的那一列。沒釘住的暫定行為會在答案出來之後被當成規格
- 技術細節：`plan.md` 的技術決定沒寫的（時鐘放哪個 package、時區住哪一層、到期日怎麼表示）——照專案的
  架構文件推得出來的，照推出來的做，不算設計決定；推不出來的，照暫定行為處理（記缺口、選一個，
  有行為的用測試釘住），`plan` 再跑時在設計裡補上技術決定（記進 `## Further Notes`）；屬於需求的才進 `## 缺口` 交回 `spec`
- 爭用接合點：這張票的票節列了 `爭用` 的檔時，核對那段邏輯實際落在哪個檔。落在別的檔，
  在票節 `#### 缺口` 記一列（列的是哪個檔、實際在哪個檔）——`plan` 再跑時據此改正接合點
- 漏列與多列的接合點：這張票動到、別張票也會動、但票節與 `tasks.md` 都沒列的檔，記一列「沒列到的
  接合點：<檔>（<種類>）」；列了卻根本沒動的，記一列「列了沒動：<檔>」——讓 `plan` 補上或收回，
  不必等合回撞到

### 7. 合回主線（只在 worktree 裡做時）

這張票由做它的這個 session 自己合回，不交給別人——衝突要靠這張票的脈絡才解得對。

1. **先把主線拉進 worktree**：在 worktree 裡 `git rebase main`，衝突在這裡解，不在主線上解。
   用 rebase 不用 merge——merge 會多一個反方向的 merge commit，打破「一張票一個 commit」
2. **照接合點的種類解**（`../plan/references/tasks-format.md`〈接合點與同時跑〉）：
   - `tasks.md` → 兩邊改的一定是不同的票節（各自只改自己那一節）：取聯集——這張票的票節用 worktree 這邊的，
     其餘用主線的。衝突落在同一個票節裡，代表有人改了不歸他的行，停下回報
   - `產生` → 不手動合：先解完原始檔，再照 `plan.md` `## Commands` 的產碼命令重新產生；`go.sum`
     用 `go mod tidy`
   - `追加` → 取聯集，**以宣告為單位**：衝突的切點常落在函式中間，兩邊共用的 `return`、`}`、契約裡的
     `"500"` 只接一次會漏、接兩次會重複；兩張票都從 `plan.md` 照抄的同一個 `components`，一模一樣，留一份。接完先編譯（實測場 `go vet ./...`），再跟主線 `git diff main`
     確認沒有刪掉主線任何一行，最後檢查有沒有重複宣告
   - 沒列到、但這張票對那個檔的**整份改動**都只是新增頂層宣告（新函式、新型別，沒有改到任何既有的行）
     → 當成 `追加` 解，並在票節 `#### 缺口` 記一列「沒列到的接合點：<檔>」，讓 `plan` 補上
   - `爭用`、或票節與 `tasks.md` 都沒列到的其他衝突 → **停下來回報**，不自己判——那是切票時
     沒預料到的耦合，要回 `plan` 重排。停下來時：
     - 撤掉這次合併（`git rebase --abort`），worktree 回到這張票自己的 commit，不留半解的衝突
     - 在票節 `#### 缺口` 記一列：衝突的檔、跟哪張票撞、為什麼不是 `產生`／`追加`；`amend` 進這張票的
       commit，保持一張票一個 commit
     - 「列到」看 `tasks.md` 任何一張票的 `接合點` 欄或 `plan.md` 的 `## Parallel Work`——那個檔的種類
       有人定過；只記在 `#### 缺口` 的預感不算列到
     - 告訴使用者跑 `/plan`：它就地補上
       接合點的種類後再重排；重排完，這張票從本步驟（第 7 步）的第 1 點重來，不必重做第 1–6 步
     - SHOULD: 這一列點名的檔，同一波還有別張票會動時，在那張票合回之前先跑 `/plan`——
       不然它會為同一個原因再停一次
3. **在合併後的程式碼上重跑閘門**，場景範圍是主線上所有驗收條件已勾完的票加上這張票
   （跟第 5 步第二行同一種組法：各票節 `**驗證**` 的場景運算式接成一個「或」，前置整理的票跳過）。
   重跑要繞過測試快取（實測場 `-count=1`）——rebase 沒改到的套件會拿到上一次的快取結果。
   綠了才合回主線：到**主線的工作目錄**（`main` 往往 checkout 在另一個目錄，worktree 裡切不過去），
   確認工作區乾淨、在 `main` 上，再 `git merge --ff-only <票的 branch>`
   - 合完 git 沒報衝突、閘門卻紅，常見三種，都在票節 `#### 缺口` 記一列，`plan` 補成接合點：
     | 情況 | 怎麼辦 |
     | --- | --- |
     | 兩張票各自新建了同名的輔助函式、替身或變數 | 兩份**一模一樣**就刪掉這張票的、用另一張的；不一樣才改這張票自己那份的名字 |
     | 另一張票新加了必填的組裝欄位，這張票新寫的程式或測試沒帶 | 照〈不越界〉補上 |
     | 另一張票讓一個介面長大，這張票新建的實作或 stub 沒有那個方法 | 補上那個方法（`追加`） |
4. **判 `@review`**：照 `tasks.md` 列的票，逐張讀票節的 `涵蓋場景`（它寫了是哪個 `.feature`），找出
   涵蓋這個 `.feature` 的票，再看它們的驗收條件——全勾了，先在主線跑一次**整個檔**、看場景數跟檔裡的
   場景一樣多，全綠了，就是這次合回讓最後一張完成，**由這次合回翻 `@review`**：
   fast-forward 帶不進新的改動，所以在主線上另開一個 commit 只翻 tag（這是「一張票一個 commit」唯一的
   例外）；還有沒勾的就留 `@wip`
5. 先準備好的先合；同時準備好的，票號小的先合。後合的那張回到本步驟第 1 點

IMPORTANT: git 合得乾淨不等於合對了。兩張票各自新建同一個 package，git 不會報衝突，編譯
卻壞在重複宣告——所以第 3 步的閘門跑在**合併後**，不是合併前。

### 不越界

MUST NOT: 改 `spec.md` 或 `.feature` 的例子。例子跟規則打架、規則有洞、少一個邊界——
記進票節的 `#### 缺口`，回報，然後
照現在的規格做。一個人在 IMPLEMENT 改掉驗收條件，等於沒有人同意過那條規格。

前面的票寫的程式碼可以改——這張票的場景需要時，就改。條件是那些票的場景照樣綠：第 5 步的
閘門要跑主線上所有驗收條件已勾完的票的場景，不只這張。改不動、或改了會讓前面的票的
`plan.md` 的技術決定不成立 → 記進票節的 `#### 缺口` 並回報，不要改掉那張票的驗收條件或票節。

前面的票的測試寫的是被 `spec.md` 修訂推翻的規則（喚醒之後常見）：

| 推翻了什麼 | 怎麼做 |
| --- | --- |
| 只有測試名或註解（期望值照新規則仍然成立） | 改名、改註解，期望值不動；在這張票的 `#### 缺口` 記一列 |
| 期望值也是（新規則下它該紅） | 這張票的驗收條件本來就要改掉它：照新規則改期望值，先看它紅、再讓它綠；在 `#### 缺口` 記一列哪張票的哪個測試 |

MUST NOT: 在 handler、mapper、presenter 裡放規則。判準：把這段刪掉，哪個 domain
單元測試會紅？沒有，它就不是規則，是接線；有，它放錯層了。

## 產物

```
<專案>/
├── internal/…                 這張票穿過的每一層
├── api/openapi.yaml           這張票的 operation 與 components，從 plan.md 的 OpenAPI 片段照抄（有片段時）
├── <migration>                這張票的表，照 plan.md 的 ## Database Schema（票節有 **資料表** 時）
├── test/acceptance/steps_*.go 新的 step definition，按領域名詞分檔
├── features/<story>.feature   只改第一行的狀態 tag
└── specs/ongoing/<date>-<feature>/tasks.md   只改這張票的票節：勾驗收條件、寫缺口
```

MUST: 不修改 `spec.md`、`plan.md`、`.feature` 的任何場景文字；`tasks.md` 只改這張票的票節裡歸 `implement` 的行（勾、`**檢查點**:`、`#### 缺口`）。
同時做的幾個 session 各寫各的票節才不會互相覆寫；設計的洞由 `plan` 再跑時在設計裡補上
（記進 `plan.md` 的 `## Further Notes`）；屬於需求的才進 `## 缺口` 交回 `spec`。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 寫程式前看過這張票的場景紅（`step is undefined` 或斷言失敗，不是 parse error），每條有分支的規則的 table-driven 單元測試也先看過紅；零值例子沒有紅的，回報裡有它的突變驗證（改壞了什麼、怎麼變紅）
- [ ] 這張票的場景綠、專案的驗證閘門（實測場 `make verify`）綠，兩個命令與結果都在回報裡
- [ ] 寫新步驟前 grep 過既有的 step definition（實測場 `make steps`），沒有同一句話兩種寫法
- [ ] 測試替身回的是副本，不是活的指標
- [ ] handler、mapper、presenter 裡沒有規則
- [ ] `spec.md`、`plan.md`、`.feature` 的場景文字都沒有被改；`tasks.md` 的 diff 只落在這張票的票節
- [ ] `.feature` 只有涵蓋它的票全勾完、整個檔的場景都綠了才改 `@review`，否則留 `@wip`
- [ ] 這張票是一個 commit，只含這張票動到的檔
- [ ] 在 worktree 裡做的：衝突照種類解（`產生` 重新產生、`追加` 取聯集），碰到 `爭用` 或沒列到的衝突停下來回報；閘門在合併後的程式碼上重跑過、綠了才合回
- [ ] 跑的是票節 `**驗證**` 的命令，另跑了已勾完的票的場景；這張票前面有檢查點時，使用者看過檢查點的結果才開工
- [ ] 票節只勾了有測試驗過的驗收條件
- [ ] 每個暫定行為都有一個測試釘住，註解寫明等哪一題的答案
- [ ] 這個 session 只做了一張票
- [ ] 有 OpenAPI 片段時，`api/openapi.yaml` 裡這張票的 operation 與 `components` 跟 `plan.md` 的片段一致；migration 跟
  `## Database Schema` 一致；要改形狀的記進了 `#### 缺口`，沒有在專案裡改
- [ ] `plan.md` 有 `## Domain Model`、`spec.md` 沒有 `## Domain Model` 與 `### 聚合與不變條件`；不是就停下請使用者重跑 `/plan`

## 完成後

告訴對方這幾件事：

1. 票號、涵蓋的場景幾個綠、閘門命令與結果
2. 新增了哪些檔、動了哪一層（一句話，不是檔案清單）
3. 票節勾了哪幾條、哪一條沒勾與為什麼；票節 `#### 缺口` 記了什麼
4. 暫定的決定有哪些、各被哪個測試釘住
5. 下一張票是哪一張（`tasks.md` 裡 Blocked by 剛被解開的那張）；這個 `.feature` 還在 `@wip` 的話，等哪幾張票
6. 在 worktree 裡做的：合回時解了哪些衝突、各是哪一種，合併後的閘門命令與結果

全部票做完 → `.feature` 全部 `@review`，下一步 `archive`：把這個 feature 寫回 `docs/` 的現況、搬進
`specs/sleeping/`。REVIEW 還沒實作，在那之前先人工對照 `spec.md` 的 AC 與 `.feature`，跑
`check_spec.py` 看鏈有沒有斷。
