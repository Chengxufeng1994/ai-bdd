# ai-bdd 路線圖

> **這份是 ai-bdd 這個 plugin 自己**要建什麼、怎麼切、什麼順序。
> 它**不是**本流程的產物——流程產出的計畫叫 `specs/<date>-<feature>/plan.md`，
> 那是使用端 repo 裡的東西。看到「plan」先分清楚是哪一個。

四份背景筆記，分工不同：

| 文件 | 回答 | 什麼時候會改 |
| --- | --- | --- |
| [docs/bdd.md](./docs/bdd.md) | BDD 這套方法本身是什麼（三個實踐、雙迴圈、Gherkin 反模式） | Cucumber 改了定義，或我們讀錯了 |
| [docs/tdd.md](./docs/tdd.md) | TDD 這套方法本身是什麼（Canon TDD、三律、兩個學派、BDD 從哪裡分出去） | 來源改了定義，或我們讀錯了 |
| [docs/sdd.md](./docs/sdd.md) | SDD 這套方法本身是什麼（spec-first/anchored/as-source 三層、Kiro/Spec Kit/Tessl、EARS） | 來源改了定義，或我們讀錯了 |
| [docs/ai-sdlc.md](./docs/ai-sdlc.md) | 這條流程為什麼長這樣——六步各自的理由、哪幾段不是 BDD、AI 特有的約束 | 隨實跑演進 |

> **這份文件是假設，不是定案。** 六步框架與底下的 skill 清單目前只有四步有實作，
> 而這四步的職責在 2026-10-02 這次重排中被重新切過。它應該被下一批真實產物
> 驗證或推翻，而不是被照著蓋完。

---

## 定位

**完整自建。** 各步驟由 ai-bdd 自己提供，不依賴外部 plugin 或 skill。

代價：環境中既有的 `clarify-loop`、`testing-strategy`、`code-review-and-quality`
等 skill 覆蓋部分相同處境，會與 ai-bdd 競爭觸發權。模型會挑一個，而你不會知道
它挑了哪個。緩解方式是在採用 ai-bdd 的專案停用重疊 plugin，並讓所有
`description` 綁定 BDD 脈絡詞彙。

換得的好處：鏈上每個環節的輸入輸出格式**由本專案自行定義**，因此可以真正閉合。
不是功能更多，是介面一致。

---

## 三種模式，三條紀律

這是 2026-10-02 重排的核心（2026-11-13 把 CAPTURE 併回 DISCOVERY 開場；2026-11-14 DISCOVERY
改成只記問答、技術設計移到 PLAN）。六步之所以站得住，是因為前三步各自只有一種主要的工作模式：

| 步驟 | 模式 | 只做 | 不做 |
| --- | --- | --- | --- |
| **DISCOVERY**（`discovery`） | agent 問、需求方答（開場建檔那一回合是需求方說、agent 記） | 原文逐字、粗切 story、問紅卡、記答案原話與誰答的 | 不寫 FR／AC、不寫 Gherkin、不設計、不自己填答案 |
| **SPEC**（`spec`／`formulation`） | agent 寫、需求方審 | `已收斂` 的 `discovery.md` → `spec.md`（`已對焦`）→ `.feature` | 不訪談，有洞記進 `## 回 DISCOVERY 補問` 並退回 |
| **PLAN**（`plan`／`task-breakdown`） | agent 設計與切、工程師審 | 技術節 → tracer bullet ＋ blocking edges ＋ 接合點 | `plan` 不改需求、`task-breakdown` 不改設計，有洞記進 `## 缺口` 並退回 |

「agent 補一個合理的答案然後繼續」這個風險，舊版靠「DISCOVERY 不寫檔」擋，
而 2026-09-14 的設計承認那條從來沒真的成立過。現在改靠「每條 FR／AC 必須指得出
來源」擋：`← PRD §n` 或 `← Q-<n>`，`check_spec.py` 第 14 項機械地抓，並跨檔確認
`← Q-<n>` 在 `discovery.md` 找得到、而且是已答。擋的是同一件事，
但從「不准做」改成「做了要留證據」。

紀律的出處：Matt Pocock 的 `to-spec`——「The spec is a record of decisions
already made, not a place where new ones get made.」

---

## 六步流程

| 步驟 | 這一步在問什麼 | 作法 | 產物 |
| --- | --- | --- | --- |
| **DISCOVERY** | 需求方想要什麼？還有哪些不知道？問到不能再問為止（brainstorming） | 開場當書記：原文逐字、粗切 story；之後設計樹＋前沿，答案逐題記下 | `discovery.md`（澄清中 → 已收斂） |
| **SPEC** | 這些答案寫成規格長怎樣？需求方點頭了嗎？（to spec） | `discovery.md` → `spec.md`（`spec`）→ 就緒判定（`example-mapping`）→ Gherkin（`formulation`） | `spec.md`（待對焦 → 已對焦）＋ `.feature`（可執行）＋ `domain-model.md` |
| **PLAN** | 怎麼做？怎麼切成一次做得完的步驟？（to plan） | 技術設計（seam、契約、schema、測試層級）＋ tracer bullet 垂直切片 ＋ blocking edges | `plan.md` |
| **IMPLEMENT** | 讓這張票從 red 到 green | step defs ＋ outside-in 雙迴圈 | 程式碼 |
| **REVIEW** | 這變更真的滿足它宣稱的情境嗎？鏈斷了嗎？ | 全新 reviewer 對照情境驗收 ＋ 鏈稽核 | 審查結論 ＋ 稽核報告 |
| **VERIFY** | 交付前最後一次：綠了嗎？紅的是 bug 還是規格過期？ | 跑閘門與全部情境 ＋ 失敗判讀 | 測試結果 ＋ 判讀結論 |

REVIEW 排在 VERIFY 前面：review 抓到的修完，最後一次驗證才算數。

框架相依度：DISCOVERY／SPEC／PLAN 與框架無關，可先做完並實際使用；
IMPLEMENT／VERIFY 需要綁定技術棧；REVIEW 部分需要。

---

## 產物佈局

一步一個產物槽，`ls` 就看得出這一批走到哪：

```
specs/
├── domain-model.md     SPEC      聚合、不變條件；跨 feature，活得比任何一份 plan 久
└── <date>-<feature>/
    ├── discovery.md    DISCOVERY 原文、粗切 story、問答。決策史，活得比實作久
    ├── spec.md         SPEC      story、FR、AC、角色、範圍。業務規則，活得比實作久
    ├── acceptance/     SPEC      <旅程>.feature：業務語意的驗收旅程，不執行；就緒判定對著它答
    ├── plan.md         PLAN      seam、模組、契約、schema、測試層級、風險。實作開始就過期
    └── tasks/          PLAN      index.md（順序、Blocked by、接合點）＋ 一張票一個檔

features/<story-slug>.feature    SPEC  可執行規格，一則 story 一個檔
docs/CONTEXT.md                  SPEC  詞彙，跨 feature
```

**`spec.md` 與 `plan.md` 分檔的理由是壽命。** FR、AC、角色在實作完成後仍然是這個
feature 的事實；seam 與 schema 在第一個 commit 之後就開始過期，跟用到它們的票一起放。
`discovery.md` 與 `spec.md` 分檔的理由是**只有一處寫規則**：問答是來源，規格是定案。

可執行的 `.feature` 是唯一不進 `specs/` 的產物（`acceptance/` 的旅程不執行，所以留在 `specs/`）——它同時是規格與可執行測試，
**位置由測試框架決定**。鏈的完整性靠 tag 與 slug 維繫，不是靠同目錄。

### 三份檔的節

`discovery.md`（格式契約在 `skills/discovery/references/discovery-format.md`）：

| 節 | 內容 | 誰寫 |
| --- | --- | --- |
| Document Overview | 狀態（澄清中／已收斂）、結束方式、修訂歷史、核心關係人 | DISCOVERY |
| 粗切 story | 每則一行，`US-<n> · <slug>` ＋ `身為…`，沒有規則 | DISCOVERY 開場 |
| Open Questions | 五欄表 ＋ 每題小節（答案原話、誰答的、歸屬） | DISCOVERY |
| 原文 | 需求方的話，逐字 | DISCOVERY 開場，之後不得改寫 |

`spec.md`（格式契約在 `skills/spec/references/spec-format.md`）：

| 節 | 內容 | 誰寫 |
| --- | --- | --- |
| Document Overview | 狀態（待對焦／已對焦）、修訂歷史、核心關係人（照抄） | SPEC |
| Background／Goal／Scope／Success Metrics | 問題、預期結果、In／Out、怎麼知道做成了 | SPEC，行尾 `← PRD §n`／`← Q-<n>` |
| Personas | `### P-<n>` 五欄 | SPEC，從角色那幾題 |
| User Stories | `### US-<n> · <slug>` ＋ `#### FR`／`#### NFR`／`#### Q` | SPEC |
| Non-Functional Requirements／Assumptions & Constraints | 跨 story 的品質屬性、外部限制 | SPEC |
| 回 DISCOVERY 補問 | 寫規格時發現、問答推不出來的洞 | SPEC |

`plan.md`（格式契約在 `skills/plan/references/plan-format.md`）：

| 節 | 內容 | 來源 |
| --- | --- | --- |
| 檔頭 ＋ 核心關係人 | 來源、seam、關係人（照抄 `discovery.md`） | — |
| Tech Preferences & Constraints／Known Boundaries | 技術選型偏好與禁令、三張清單 | 技術面向的答案 |
| Implementation Decisions | 動哪些模組、介面長怎樣、API 契約、schema、架構決定 | 技術面向的答案 |
| Testing Decisions | **seam** ＋ 每個 AC 的測試層級 ＋ 既有測試的 prior art | 技術面向的答案 |
| Risks | 規格沒說但實作一定撞到的（併發、冪等、交易邊界） | 從規則推出來 |
| 缺口 | 推不出來的（票與前置整理 2026-11-15 起住 `tasks/`） | — |

**一批一份，不是一則 story 一份。** 這不是新決定：`plan` 的
`references/plan-format.md` 記錄過「每 feature 一份 ＋ 共用檔」試過並且壞掉——
共用檔把實質吸走，各檔退化成「見共用檔」；而且節號互撞，部分對齊比完全不對齊
更危險。API 契約與 schema 天生跨 story，只有並排才看得見衝突。

**規則與例子住 `spec.md`。** `check_spec.py` 比對的是 `spec.md` 的 FR／AC 與 `.feature`
的 tag——**兩份獨立的表述，重複就是那個檢查**。`plan.md` 若也抄一份就成了三份。

**MUST NOT 寫具體檔案路徑與程式碼片段。** 它們過期得比什麼都快。
例外：prototype 產出的、比散文更精確地編碼了某個決定的片段（狀態機、reducer、
schema、型別形狀），可以內嵌並註明來自 prototype。

### `tasks/`

格式契約在 `skills/task-breakdown/references/tasks-format.md`。`index.md` 依賴序排列，blockers
在前；一張票一個檔，每張票要有：

- **交付什麼行為**——從使用者角度，不是分層清單
- **Blocked by**——哪幾張票要先完成，或「無，可立即開始」
- **接合點**——這張票會改、別張票也可能改的檔案，標 `產生`／`追加`／`爭用`；只有 `爭用` 與兩張都要新建的東西擋住並行（2026-11-15 試跑後改，`docs/superpowers/specs/2026-11-15-parallel-tickets-design.md`）
- **驗收條件**——每一條都要能在起點 commit 上**失敗**

prefactoring（「make the change easy, then make the easy change」）排最前面。

2026-11-15 之前票住在 `plan.md` 裡。拆出來是為了並行：做票的 session 只寫自己那張票檔，
沒有任何一份檔案被兩個 session 寫（`docs/superpowers/specs/2026-11-15-bdd-tasks-design.md`）。

---

## 三個新概念

現在完全沒有、這次要加進來的：

### 1. Slice ——「一批」

`plan.md` 一批一份，所以「一批」多大是個真的決定。判準：
**這批做完，能不能對使用者交付一次完整的價值。**

實測場的數字給了上限感：6 則 story、69 個場景 → `plan.md` 36.5K。
那已經在可讀性邊緣。`to-tickets` 的文件記錄過超過之後的具體災難：
spec 大到 tracker 讀不回來，agent 燒 tool call 在重抓片段，永遠讀不到結尾。

### 2. Seam ——「驗收測試打在哪一層」

HTTP handler？application service？domain？**這是 BDD 決定 step definition
長相的那個決定**，而現在的 `plan.md` §4 只問「要不要資料庫」「內迴圈是什麼」，
從沒問過外迴圈的高度。

規則（抄 `to-spec`）：優先用既有 seam、取最高的那一層、
**整個變更的理想數量是一個**。

seam 寫進 `plan.md` 的意義是它會往下傳：IMPLEMENT 照著打，
REVIEW 可以把「沒人同意過的 seam」當成 finding 抓出來。這個綁定是間接的
——它穿過 `plan.md`——所以那份文件寫不寫這一節，決定了後面兩步抓不抓得到。

### 3. 一張票 ＝ 一個 fresh context window

現在的「實作順序」有解鎖關係，但沒有尺規，所以切多細沒有判準。
加上這把尺之後，切分有了上下界：**下界**是票要能自己 demo（垂直切片，
不是「先做完 schema」），**上界**是塞得進一個沒看過 spec 的 session。

垂直切片這條是最常被破的。`to-tickets` 記錄過一個案例：26 張票依層切
（corpus／producer／aggregator／selector），每張票平均跑 20 次 agent，
其中四分之三是返工——他們自己的事後檢討把每一類失敗都追回到水平切分。

**例外：wide refactor。** 一個機械性改動（改欄位名、換共用型別）的影響面
擴散到全 codebase 時，沒有任何垂直切片能綠。這種改 expand–contract：
先加新形式（不破壞）→ 分批遷移呼叫端（每批一張票，CI 保持綠）→
最後刪舊形式（被全部遷移票 block）。

---

## 步驟 → skills 對照

**檔案結構是扁平的**（`skills/<name>/SKILL.md`），層級只存在於下面這張表。

| 步驟 | skill | 職責 | 狀態 |
| --- | --- | --- | --- |
| DISCOVERY | `discovery` | 開場建檔，再設計樹＋前沿問到完整，只寫 `discovery.md` | **已實作**；2026-11-13 併入原 `bdd-capture` 的建檔；2026-11-14 改成只記問答 |
| | `clarify` | 一疊紅卡逐題收斂 | **已實作**；原 `clarify-loop`，由 `discovery` 呼叫 |
| | `story-splitting` | 沿規則切；九種模式 | **已實作**；由 `discovery`（粗切）與 `spec`（`已對焦` 前沿規則重切）呼叫 |
| SPEC | `spec` | `discovery.md` → `spec.md`，判就緒後呼叫 `formulation` | **已實作**；2026-11-14 從「寫技術設計」改成「從問答寫規格」 |
| | `example-mapping` | 四色卡、四條診斷、就緒判定；唯一寫 `已對焦` 的 skill | **已實作**；由 `spec` 呼叫 |
| | `spec-by-example` | `spec.md` → `acceptance/*.feature` 業務語意驗收旅程，`@covers-` 標 AC | **已實作**；2026-11-14，參考 aixbdd `spec-by-example`；由 `spec` 在診斷與就緒問句之間呼叫 |
| | `formulation` | `spec.md` → `.feature` | **已實作** |
| | `spec-review` | 反命令式、conjunction step、情境爆炸稽核 | 未實作 |
| PLAN | `plan` | `spec.md` ＋ `.feature` ＋ 技術面向的答案 → 技術節 | **已實作**；2026-11-14 接手原 `design.md` 的技術設計；2026-11-15 切票移到 `task-breakdown` |
| | `task-breakdown` | `plan.md` ＋ `.feature` → `tasks/`：tracer bullet ＋ blocking edges ＋ 接合點 | **已實作**（2026-11-15，從 `plan` 拆出；試跑一次，接合點與合回主線依試跑結果重寫） |
| IMPLEMENT | `implement` | 一張票一個 session，outside-in：外迴圈紅 → domain 單元測試 → 接線 → step definition → 閘門綠 → `@review` | **已實作**（2026-10-03，從 parking 實跑的做法抽出來）；原本規劃的 `-step-definitions`／`-outside-in` 兩個子 skill 併進本文，等實跑證明太長再拆 |
| | `bdd` | 一個範疇內一次一個 slice，red／green／refactor 三個入口擇一；`implement` 委派 green 與 refactor，red 由 `implement` 自己做 | **已實作**（2026-11-15，照 aixbdd 的 `bdd` 與 `implement` 的委派方式；還沒實跑過） |
| REVIEW | `review`、`review-acceptance`、`review-traceability` | 審查總入口、對照情境驗收、鏈稽核 | 未實作 |
| VERIFY | `verify`、`verify-failure-triage` | 跑閘門與情境並回報、真 bug vs 規格過期 | 未實作 |

`spec` 加了技術設計之後會不會太肥、要不要拆成兩個 skill——**先不拆**（2026-11-14 技術設計已移到 `plan`）。
子 skill 是否真的需要，等步驟 skill 跑過真實案例、發現它太長時再拆。
現在拆是猜的。**2026-09-14 翻案：拆了。** 見
`docs/superpowers/specs/2026-09-14-grilling-and-formulation-design.md` 的 D8。

### Example Mapping 的兩個訊號在哪裡看

規則與例子住在 `spec.md` 每則 story 的 `#### FR`，所以 Example Mapping
原本靠數卡片的兩個訊號都還在，只是各有兩個來源：

| 訊號 | 卡片 | 現在看哪裡 |
| --- | --- | --- |
| story 太大 → 該切 | 藍卡（規則）太多 | `spec.md` 該 story 底下 `#### FR` 的條數 |
| 不確定性高 → 未就緒 | 紅卡（問題）太多 | `## Open Questions` 的列數，等於 `status.py` 的未答數 |

兩個訊號都不是硬門檻——實測中 `visitor-fee-and-exit` 觸發了「story 太大」
（10 條規則、三則裡最高）而**被有理由地否決**：那些規則是同一個行為的不同分支，
拆開會產生半條規則的 story。訊號亮了要看，不是要照做。

### 追問覆蓋改成算出來的

十六個追問面向的覆蓋**從 `spec.md` 的 `## Open Questions` 算**：每一列的 `面向` 欄
宣告它蓋哪個面向，`skills/discovery/scripts/status.py` 聚合出「哪幾個面向一題
都沒有」。這符合本專案既有的「進度是算出來的，不另存一份」原則——存一份彙總
等於同一個數字兩個地方，而它們遲早不一樣。

`n/a`（這個面向不適用）的理由寫在那一列的答案欄。否則它跟「懶得問」分不出來。

---

## 命名與結構的決定

**扁平目錄，名字不加前綴。** 2026-11-15 拿掉了原本的 `bdd-` 前綴（使用者的決定）：

1. skill 名稱在執行時是**全域扁平**的——`claude plugin details` 列出的是名字，
   目錄層級不存在。分組因此只寫在上面的表與 README，不寫進名字。安裝後的 plugin skill
   會帶 plugin 名（`ai-bdd:plan`），跟別的 plugin 撞名時用它區分。
2. **`claude plugin validate` 不遞迴。** 巢狀目錄（`skills/<分類>/<名稱>/SKILL.md`）
   在 runtime 有效，但 validate 掃不到，等於那些 skill 永遠不會被檢查。
   扁平換來的是每個 skill 都在 CI 的守備範圍內。

拿掉前綴的代價是字母排序不再把它們聚在一起，而 `plan`、`spec` 這類名字更通用。
`tasks` 因此改名 `task-breakdown`：它寫的 `tasks/` 目錄與它同名，讀起來分不清是 skill 還是目錄。

通用字眼（`story-splitting`、`failure-triage`）要在 `description` 裡寫明**排除
條件**，否則會在非 BDD 情境被誤觸發。

**`slice` 與 `feature` 是兩個粒度，不共用同一個字。**
`specs/<slice-slug>/` 裝的是一批（好幾則 story）；`features/*.feature` 是
一則 story 一個檔。用同一個字會製造「部分對齊」——讀的人會假設兩邊指同一件事，
然後在其中幾處讀錯。`feature` 這個字純粹留給 `.feature`。

**三個 plan。** 根層 `PLAN.md`（這個 repo 要建什麼）、步驟 PLAN、產物
`specs/<date>-<feature>/plan.md`。後兩者保留，因為「PLAN ＝ 把 spec 拆成可逐條
執行的任務、寫進 `plan.md`」是既有慣例；歧義靠這份文件開頭那段話消解。

---

## Token 成本

`claude plugin details` 會估算成本。實測參考值（mattpocock-skills，25 個 skill）：

```
Always-on: ~1,620 tok  每個 session 都付
每個 skill：常駐 ~30–160 tok（名稱＋description）
            觸發時 ~20–3.9k tok（本文）
```

推論：**拆得細不貴，description 寫得囉嗦才貴。** description 每次 session 都付費，
本文用到才付。所以 description 壓成一行觸發條件，厚重內容放本文或參考檔。

---

## 可追溯鏈

```
需求方的話（discovery.md ## 原文）→ Story → FR → AC → Scenario → Ticket → Step → 程式碼
                                      ↑        ↑                   ↑       ↑
                               開場粗切      SPEC 定版            一票一   seam 決定
                                             行尾 ← Q-n / PRD §n  session  在 PLAN
                                                  ↓
            回饋 ← 線上行為 ← 發布 ← 審查 ← 測試執行 ←┘
```

每個箭頭都是一個意圖遺失點。保住這條鏈是 ai-bdd 存在的理由，也是最難被取代的
部分——因為鏈的完整性檢查**幾乎全是機械性的，人只是懶得做**。

這條鏈的定版點被改過三次，三次的理由都值得留著：

- **2026-09-02**：編號從 CLARIFY 移到 SPEC 定版。理由是「凍在 CLARIFY 等於在最
  不確定的時候做最不可逆的事」。
- **2026-09-09**：**移回 CLARIFY**。CLARIFY 的產物變成一份自足的 `prd.md`，
  `FR-<n>` 與 `AC-<n>.<m>` 在那裡誕生，`.feature` 的 `@example-<n>.<m>` 回指它們。
  上面那個顧慮改用別的方式解掉：**不確定的東西根本不會拿到編號**——它留在
  `## Open Questions` 當一列待答，而不是先變成 FR 再回頭改。
- **2026-10-02**：**留在 DISCOVERY，但換了載體。** 編號在 `prd.md` 誕生，
  `clarify-log.md` 與 `input.md` 消失，`spec.md` 只剩技術節。來源標記取代
  「不寫檔」成為守門規則：沒有 `← Q-<n>` 或 `← PRD §n` 的 FR／AC 是缺陷。

還有一處改動至今未變：

- **IMPLEMENT 的輸入從「一個場景」變成「一張票」。** 票是垂直切片，
  可能蓋到多個場景；票自己可 demo，單一場景不一定。

### 狀態 tag

`.feature` 在 `Feature:` 上掛恰好一個狀態 tag（`@draft` `@ready` `@wip`
`@review` `@done`），記錄它走到六步的哪一站。

完整詞彙表與三個設計決定在 `skills/formulation/references/state-tags.md`——
**它是 skill 之間的契約，所以跟著 skill 走，不放在這份路線圖裡**。
別人安裝這個 plugin 時拿不到 repo 的文件，契約寫在這裡就等於斷了。

依賴方向：repo 的文件可以指向 skill，skill 不可以指向 repo 的文件。

---

## 進度

- [x] `bdd-clarify`
- [x] `clarify-loop`
- [x] `story-splitting`
- [x] `benchmark/skeleton/go` — Go ＋ godog 實測場（骨架綠，零業務）
- [x] 拿健身追蹤需求跑一次 `bdd-clarify`，依實跑結果修正格式
- [x] `spec` — 依 11 份真實 `.feature` 反覆修正，封閉文法把
      76 個場景的步驟樣板從 208 降到 41
- [x] `plan` — 拿實測場的 5 個 `.feature`（69 個場景）跑過，依執行者回報的
      三處摩擦修正
- [x] 2026-09-02 重規劃：三條紀律、slice、seam、tracer bullet 尺規
- [x] `bdd-clarify` 改版（Pass 3、分批、拿掉 example-mapping）
- [x] `spec` 擴張（`spec.md`、seam、`domain-model.md`）
- [x] `plan` 重寫（tracer bullet ＋ blocking edges）
- [x] 實測場整批重跑一次，驗證新的三步分工
- [x] 2026-10-02 重排：CAPTURE 步驟、`prd.md`、來源標記守門、`clarify-loop` 改名 `clarify`
- [x] 實測場重跑一次，驗證 CAPTURE → DISCOVERY → SPEC 的新分工——2026-10-03 用 parking-billing 跑到 IMPLEMENT，八個發現在 `runs/go-parking-billing/REPORT.md`
- [x] `implement`——2026-10-03，從那次實跑抽出
- [x] 2026-11-13 砍掉 CAPTURE 與 `prd.md`；2026-11-14 六步流程：DISCOVERY 只寫 `discovery.md`、SPEC 從問答寫 `spec.md`、`design.md` 併進 `plan.md`、REVIEW 排到 VERIFY 前
- [x] 2026-11-14 `spec-by-example`：就緒判定前先寫 `acceptance/` 驗收旅程，`check_spec.py` 第 17 項稽核 `@covers-`
- [ ] `verify`、`review`
- [ ] 用 `implement` 重跑一次 parking 的某張票，驗證 skill 寫的跟做的一樣

順序刻意是「做一個 → 實跑 → 再做下一個」。**產物格式就是介面**：SPEC 讀已答的
問題、PLAN 讀 `spec.md` 與 `.feature`、REVIEW 稽核整條鏈，下游全都依賴上游實際
吐出什麼。在紙上設計那個格式是猜的。

這次重規劃違反了上面那句話——它是在紙上改的。所以「實測場整批重跑」那一項不能跳。

---

## 實測場

`benchmark/skeleton/go/` 是六步流程的 dogfooding 對象：Go ＋ godog ＋ 分層
架構（cmd／domain／application／infrastructure／interfaces）＋ 三層測試
（unit／integration／acceptance）＋ `api/openapi.yaml`。

它的規則是**在 DISCOVERY 跑完前不寫任何業務程式碼**。題目選「訓練總容量計算」
而非「記錄訓練」，因為前者才有規則與模糊點——自體重動作、單邊動作乘不乘二、
熱身組算不算、dropset 是幾組。詳見該目錄的 README，那裡也寫明了
`discovery` 的及格標準：**沒問到自體重與單邊動作，就算失敗**。

**現在骨架是空的。** 訓練追蹤那一批的產物（`specs/workout-tracking/` 與六個
`.feature`）在 2026-09-03 的 `110819e` 被清掉了，只留 `version.feature` 這條走路
骨架。它們仍在 git 歷史裡，`git show 9e154b9^:` 取得回來。

因此 `status.py` 與 `check_spec.py` 對這個目錄跑會 exit 1 並指名說找不到產物——
那是正確行為，不是壞掉。

`benchmark/cases/` 現在有兩個情境目錄，測互補的失敗模式：

| 情境 | 測什麼 |
| --- | --- |
| `parking-billing/` | 需求裡**寫得像已經定案的數字**，敢不敢當成未定案來問 |
| `personal-memo/` | 需求裡**兩句直接打架的話**看不看得出來，以及**熟悉的領域**會不會憑先驗填空 |

每個情境有 `input-*.md`（純酬載，可以安全地交給被測者）與 `grading.md`（答案卷）。
`personal-memo/` 另有 PRD 版輸入，測「面對看起來已經完成的文件敢不敢質疑它」。

停車場那份的模糊點清單是**跑完才寫的**，不是事先猜的；它分三層，最後一層
（推導才浮現的矛盾、以及 DISCOVERY 漏掉而 SPEC 補抓的那一格）才是有診斷力的部分。
個人待辦那份是**事先設計的**，尚未跑過——它的清單要等第一次實測才知道準不準。

> `benchmark/` 與 `runs/` 都不是 plugin 的慣例目錄（只有 `skills/`、`commands/`、
> `agents/`、`hooks/` 是），不會被載入，對 plugin 行為零影響。

---

## 未決事項

1. ~~**產物存放路徑。**~~ 已定：文件型產物一律寫入使用端 repo 的 `specs/`，
   一個目錄裝完，刪掉即乾淨。**兩個例外**，各有具體理由：

   | 例外 | 理由 |
   | --- | --- |
   | `features/*.feature` | 位置由測試 runner 的慣例決定 |
   | `docs/CONTEXT.md` | 壽命長於任何 feature——刪掉某個 feature 的 spec 目錄之後，詞彙定義應該還在 |

   例外的範圍很窄：SPEC 得建立 `docs/CONTEXT.md` 或在其中追加條目，**不得改寫
   該檔既有段落，不得寫入 `docs/` 底下其他檔案**。

2. ~~**`.feature` 的位置。**~~ 已併入 #1 的例外表——原本獨立的一條，其實
   跟「產物存放路徑」是同一個決定底下的兩個例外，不是兩件各自獨立的事。

3. ~~**SPEC 與 PLAN 的分界。**~~ 已定（2026-09-02）：技術設計歸 SPEC，
   拆解歸 PLAN。見「四條紀律」。

4. **一批到底多大。** 判準已定（一次可交付的價值），已有兩個數字：

   | 批次 | story | 規則 | 例子 | `plan.md` 技術節 | 票 |
   | --- | --- | --- | --- | --- | --- |
   | `workout-tracking` | 6 | 39 | 93 | 649 行 | 9 |
   | `parking-fee-collection` | 3 | 14 | 21 | 187 行 | 未跑 PLAN |

   兩者相差三到四倍，而兩批都通過了「一次可交付的價值」這個判準——所以那個
   判準**不決定大小**，只決定邊界切在哪。真正的上限訊號還是「`plan.md` 大到
   沒人讀得完」，而 649 行那一批當時已經被評為在邊緣。

   （`workout-tracking` 的產物已於 `110819e` 刪除，數字取自當時的稽核輸出，
   檔案在 git 歷史裡。）

5. ~~**技術問題的追問面向。**~~ 已定：技術面向清單見
   `skills/discovery/references/technical-probes.md`，目前有五個面向
   （seam、模組邊界、介面與型別契約、排序契約／決定性、既有資產／測試慣例）。
   這五個面向全部從 `workout-tracking` 這一次實測場的硬跑反推出來，只跑過
   一個技術棧（Go）、一個專案——換一個技術棧不同的專案重跑，如果冒出新的、
   真的卡住的技術缺口要加進去，這份清單不是天花板。

6. **IMPLEMENT／VERIFY 的技術棧。** 未定（Python behave/pytest-bdd、
   JS Cucumber.js、Go godog…）。前三步不受影響。

7. **License。** repo 在 `vivotek/` 底下，`plugin.json` 目前不含 `license` 欄位。

8. **十一個業務追問面向需要重新檢視。** 這次實測場一輪跑下來冒出
   16 個缺口，其中 11 個是業務缺口而非技術缺口，裡面至少兩個對不上現有十一
   個面向裡的任何一個——包括「這個能力到底在不在這批範圍內」這種連要不要問
   都還沒有位置放的問題。這裡只能誠實記下這個發現：詳細的證據在執行過程中
   工作目錄被清空，沒能留下來，所以這一項記的是結論，不是分析過程本身。第二
   個實測場跑的時候，要重新對照這十一個面向，看看夠不夠、還是要加。

9. **`audit_skill.py` 只查得到孤兒檔案，查不到反過來的情況。** S5 規則抓得到
   子目錄裡一個沒人指到的檔案，但抓不到 `SKILL.md` 指到、實際上不存在的
   章節。這一批 `plan` 整份重寫，就活生生撞過一次同樣形狀的缺陷——只是
   換了個樣子：舊版 `SKILL.md` 把缺口那一節標成「第 6 節」，但缺口那一節
   實際編號是第 7 節；這次重寫把那個引用整段拿掉了，缺陷是在人工複查時才被
   抓到的，不是任何一支腳本抓到的。要不要幫 `audit_skill.py` 加一條反向規則
   （`SKILL.md` 引用的章節號／檔名是否真的存在），是懸而未決的。

10. **這條鏈最容易出的錯是靜默地錯，不是壞掉。** 這一輪實跑撞到四個例子：一個
    dashboard 回傳 exit 0、顯示「已就緒」，但其實什麼都沒找到；一張覆蓋率表，
    只要每個問題都歸錯 story，看起來一樣完整可信；重新編號過的 Rule／Example
    標籤不會噴任何錯誤，只會讓答案悄悄變成錯的；一支腳本靠比對內文裡剛好加粗
    的文字去解析一行已經壞掉的狀態列，照樣跑得動。四個都是 exit 0、看起來健康。
    這條鏈要不要為這一類錯誤加一層系統性的檢查，還沒有答案。

11. **`feature` 這個字還有九處用在退休的舊義上。** 這一輪把 `feature` 收窄成
    「Gherkin 的 `.feature` 檔，一則 story 一個」，把交付批次改叫 `slice`。但
    `clarify-loop/SKILL.md` 有四處、`discovery/SKILL.md` 有五處仍寫「跨 feature
    問題」「換 feature 時」，那裡的 `feature` 指的是 story。最終審查判定不要把
    這個掃描併進 merge 前的修復波——散在兩個當時範圍外的檔案，只修一半比不修
    更容易誤導。它需要自己一次變更，連同 `story` 與 `slice` 的用法一起校準。

12. **`discovery/SKILL.md` 需要一次結構性的下放，不是再刮幾行。** 它連續三輪
    卡在 500 行的天花板上：每加一樣東西（技術澄清、分批、`clarify.md`）就得回頭
    刮掉幾行別的，現在是 504 行。逐行刮是在對付症狀——S10 這個訊號要講的正是
    「該把只有特定情況才需要讀的段落下放了」。候選是十一個業務面向的追問手法與
    五個技術面向的細節，兩者都只有真的在問那一輪時才需要讀。這件事跟第 11 項可以
    一起做，因為兩者都要通篇改一遍。

13. **2026-10-02 已處理**——下面點名的章節在那次重排全部重寫，留著是為了記錄當時
    發現了什麼。原文：**這份文件自己有大段仍在描述 2026-09-09 PRD 改版前的舊世界，位置比已修好
    的未決事項 #1 更靠前。** 第 11 項只點名 `feature` 這個字在
    `clarify-loop`／`discovery` 兩份 SKILL.md 裡的舊義，沒有涵蓋這裡：
    `## 六步流程` 那張表（仍列 `brief.md`、`actor.md`、`glossary.md`、
    `questions/`）、`## 產物佈局`（用現在式描述整棵已經不存在的
    舊目錄樹）、`` ### `spec.md` 的節 `` 那張表（規則定版位置仍寫
    `clarify.md`）、`### 1. Slice ——「一批」`（「一批」這個字本身已經廢除）、
    `步驟 → skills 對照`、`### Example Mapping 的兩個訊號在哪裡看`（引用
    `clarify.md` 的 `## Business Rules` 與 `#### Rule` 數）、
    `### 追問覆蓋改成算出來的` 那一節（「追問覆蓋改成從問題檔算」，
    與現行 `status.py` 直接矛盾）、`## 命名與結構的決定`（一個
    `<slice-slug>`）、`## 可追溯鏈`。讀者由上往下讀，會先撞到這些錯的版本
    才讀到已修好的未決事項。量級跟第 11、12 項一樣大，是通篇改一遍才處理得完的事，
    這裡先記下範圍，不在這次一併修——未被追蹤的債等於不存在的債。
