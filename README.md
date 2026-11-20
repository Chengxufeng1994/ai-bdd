# ai-bdd

一個 Claude Code plugin：用 BDD（Behavior-Driven Development）把需求一路帶到程式碼——
先把需求問清楚，寫成規格與驗收旅程，再翻成可執行的 Gherkin，切成票，一張一張做到綠。

> **狀態：早期。** DISCOVERY、SPEC、PLAN、IMPLEMENT 已有 skill；REVIEW、VERIFY 尚未實作。
> 路線圖見 [PLAN.md](PLAN.md)（它是假設，不是決定）。

流程是六步，每一步裡用的方法是 BDD：

**DISCOVERY → SPEC → PLAN → IMPLEMENT → REVIEW → VERIFY**

| 步驟 | Skill |
| --- | --- |
| DISCOVERY | `discovery` |
| SPEC | `spec`、`feature-mapping`、`example-mapping`、`spec-by-example`、`formulation`、`story-splitting` |
| PLAN | `plan`（設計 Domain Model／Database Schema／API 契約、依賴圖、垂直切片成票、排序與 Checkpoint；寫 `plan.md` 與 `tasks.md`）、`api-and-interface-design`（由 `plan` 委派寫 API 契約） |
| IMPLEMENT | `implement`（一張票：測試層先紅、閘門、合回）、`behavior-driven-development`（由 `implement` 委派 green／refactor，模型不自動觸發）、`archive`（票全部合回後：寫回 `docs/` 現況、搬進 `specs/sleeping/`） |
| REVIEW／VERIFY | 尚未實作 |

`skill-rules` 不在流程裡：它是本專案撰寫與稽核 skill 的房規。哪幾步對應 BDD 的哪個實踐、
哪幾步是本專案自己加的，見 [docs/ai-sdlc.md](docs/ai-sdlc.md)。

## 中間產物：誰產生、誰接著用

每一步只讀上一步實際寫出來的檔案。`<feature>` 指 `specs/ongoing/<date>-<feature>/`。

| 中間產物 | 位置 | 誰產生 | 誰接著用（用來做什麼） |
| --- | --- | --- | --- |
| `interview-log.md`（訪談紀錄，只追加：`## 原文`、`## Session`、`## 問答`、`## Confirmed Decisions`、`## Remaining Open Items`） | `<feature>/` | `discovery`（原文逐字、關係人、粗切 story、逐題的問答原話、結束方式） | `spec`：照抄原文與問答、從已答的題寫規則；`status.py`：算進度 |
| `spec.md`（`待對焦`） | `<feature>/` | `spec`（唯一的寫入者） | `feature-mapping`：攤整個 feature 的流程、找缺的替代流程與反例；`example-mapping`：攤四色卡、跑診斷、問就緒；`spec-by-example`：寫驗收旅程；`story-splitting`：重切 story |
| `spec.md` `## Open Questions` 的補問題（寫規格時發現的洞，`Q-<n>` 待答） | `<feature>/spec.md` | `spec`、`feature-mapping`、`spec-by-example`、`example-mapping`（未就緒時） | `spec`：狀態改回 `澄清中`、停；`discovery`：只讀，在紀錄追加同號的題先問 |
| `splitting.md`（切分紀錄） | `<feature>/` | `story-splitting` | `spec`：照切法寫 story，FR 編號不重編 |
| 驗收旅程 `*.feature`（不執行） | `<feature>/acceptance/` | `spec-by-example` | 需求方：對著旅程判就緒；`example-mapping`：提就緒問句；`formulation`：只借用語 |
| `spec.md`（`已對焦`） | `<feature>/` | `example-mapping`（只寫需求方的決定） | `formulation`：寫可執行 `.feature`；`plan`：寫系統分析與技術決定 |
| `spec.md` `## Domain Concepts`（只有 `### 詞彙`；聚合與不變條件在 `plan.md`） | `<feature>/spec.md` | `spec`（第 4 步，先讀 `docs/CONTEXT.md`） | `spec-by-example`、`formulation`：用詞對齊；`plan`：設計 Domain Model 的起點；`archive`：併進 `docs/CONTEXT.md` |
| 可執行 `.feature`（`@ready`／`@draft`） | `features/` | `formulation` | `plan`：每個 AC 的測試層級、每張票涵蓋哪些場景；`implement`：讓場景由紅變綠 |
| `plan.md`（技術決定：Tech Stack、Commands、Project Structure、Code Style、Boundaries、Implementation／Testing Decisions；三份設計：Domain Model、Database Schema、API 契約（OpenAPI 片段）；依賴圖、風險、要不要並行與接合點） | `<feature>/` | `plan`（seam 經使用者確認；有洞加進 `## 缺口`） | `plan`：讀三份設計切票、每張票列碰到的 operation 與表、抄模組、內迴圈、驗證命令，高風險排前面；`implement`：把票涵蓋的 OpenAPI operation 照抄進 `api/openapi.yaml`、照 Database Schema 寫 migration，其餘照著寫程式；`archive`：聚合併進 `docs/DOMAIN-MODEL.md`、表併進 `docs/DATABASE-SCHEMA.md`、技術現況併進 `docs/TECH.md` |
| `tasks.md`（一波一節、一張票一個票節：行為、模組、大小、Blocked by、接合點、驗收條件、驗證命令；波與波之間是檢查點） | `<feature>/` | `plan` 建立；`implement` 只改自己的票節：勾驗收條件、記檢查點、寫 `#### 缺口` | `implement`：挑下一張票、判斷能不能跟別張同時做、開新一波前跑檢查點，一張票一個 session；`plan` 再跑時把設計的洞在設計裡補上（記進 `## Further Notes`）；屬於需求的才進 `## 缺口` 交回 `spec` |
| 產品程式碼、step definition、`@review` tag、票節勾選 | 專案程式碼、`features/` | `implement` | `archive`：確認全部做完；REVIEW／VERIFY（尚未實作，目前人工對照） |
| `docs/DOMAIN-MODEL.md`、`docs/DATABASE-SCHEMA.md`、`docs/TECH.md`（專案現況，只動 ai-bdd 區塊） | 專案 `docs/` | `archive` | 下一個 feature 的 `plan`：聚合、表與技術現況從這裡讀 |
| `specs/sleeping/<date>-<feature>/`（做完的 feature 整個目錄） | `specs/sleeping/` | `archive`（`git mv`） | 要修正時，`spec` 喚醒、搬回 `specs/ongoing/`、狀態改回 `澄清中` |

## 每個 skill：職責與流程

下面摘自每個 skill 的職責卡與流程總表；完整的讀寫檔案、規則與 Verification 在各自的 `skills/<name>/SKILL.md`。

### DISCOVERY

#### `discovery`

**職責**：把一段需求（想法或別人寫的 PRD）收成逐字原文與粗切 story，再一題一題問到猜得中需求方的反應，每個答案逐題記下，最後讓需求方對一段意圖陳述說 yes

- **上游**：無（可單獨使用）；`spec` 開了補問題、把狀態改回 `澄清中` 時再跑
- **下游**：`spec`（從已答的題寫 FR／AC），由使用者叫
- **不做**：不寫 FR／AC／Gherkin／規則草稿、不自己填答案、不改寫 `## 原文`、不寫 `spec.md` 的其他節、不寫 `plan.md`／`.feature`、不翻 `待對焦`

0. 建檔：原文逐字、粗切 story、最多四題、改到需求方點頭
1. 開場其餘：掃原文找隱含假設、補核心關係人、識別角色——問出來的都開成『待答』題；寫一句開場假設與信心
2. 算前沿，過「值不值得問」；能自己查的事實派 sub-agent 查；照影響排序
3. 問影響最大的那一題：上下文、一句問題、猜測的答案與理由
4. 答案回來：記下、確認夠具體、更新信心、看停下來的訊號 → 回到 2
5. 停的條件成立：收尾核對 → 意圖陳述 → 需求方說 yes → 結束方式，停

### SPEC

#### `spec`

**職責**：把已答的問答寫成規格：要做成什麼、不做什麼，帶到需求方說「就緒」

- **上游**：`discovery`
- **下游**：`feature-mapping` ＋ `example-mapping` ＋ `spec-by-example`（判就緒）→ `已對焦` → 交回，使用者接著叫 `formulation`
- **不做**：不訪談、不填答案、不改 `## 原文`、FR／AC 不標 `推論`、不自己寫 `已對焦`、不寫 `.feature` 也不呼叫 `formulation`、不寫 `docs/`、不寫技術決定與 seam、不排風險與並行（都是 `plan`）、不切票

1. 寫需求部分（`待對焦`）：story 與優先級、FR／AC、Edge Cases
2. 收尾檢查，必要時重切 story
3. 自檢 `check_spec.py`
4. 詞彙 → `spec.md` 的 `## Domain Concepts`（聚合與不變條件留給 `plan`）
5. `feature-mapping` 與 `example-mapping` 診斷 → `spec-by-example` 寫旅程 → 就緒問句
6. 需求方答「就緒」→ `已對焦` → 交回，停
7. 有洞（不就緒、寫不出來源、旅程有 `# [need clarification]`）→ 開補問題、狀態改回 `澄清中`，停

#### `feature-mapping`

**職責**：把整個 feature 攤成「規則 × 例子 × 步驟 × 結果」的列，逐步找出缺的流程與反例，寫成補問

- **上游**：`spec`（第 5 步第一個呼叫）；也可單獨使用
- **下游**：`example-mapping`（照樣跑完，兩張圖的洞一次交給 `discovery`）
- **不做**：不訪談、不改 US／FR／AC、不寫狀態、不判就緒、不寫 Gherkin、不產出 `feature-map.md`

1. 認模式：流程內或單獨
2. 排主流程：從 `## 原文` 抽出整個 feature 的步驟順序
3. 攤列：每條 AC 一列，依 FR 分組，對齊到主流程
4. 跑三條診斷
5. 流程內：缺口寫進補問；印地圖與結果，交回 `spec`

#### `example-mapping`

**職責**：把一則 story 攤成四色卡地圖、跑四條診斷，把就緒問句交給使用者

- **上游**：`spec`（`spec.md` 寫完、`feature-mapping` 跑完後呼叫）；也可單獨使用
- **下游**：就緒 → `spec` 交回，使用者接著叫 `formulation`；未就緒 → `spec`〈有洞〉
- **不做**：不替使用者判就緒、不把紅卡問到收斂、不切 story、不寫 Gherkin、不產出 `example-map.md`

1. 認模式：流程內或單獨
2. 攤卡，在終端印出地圖
3. 跑四條診斷
4. 就緒判定：給判讀、提就緒問句、依回答寫狀態

#### `spec-by-example`

**職責**：把 `待對焦` 的規格寫成幾條可 demo、只用業務語言的長流程，讓需求方對著它判就緒

- **上游**：`spec` 第 5 步（`example-mapping` 跑完診斷之後）
- **下游**：`example-mapping` 提就緒問句；有 `# [need clarification]` → `spec`〈有洞〉
- **不做**：不執行、不訪談、不改 FR／AC／狀態、不填答案、不寫 `spec.md` 沒有的結果

1. 讀 `spec.md`，挑旅程
2. 照 `rules/journey-format.md` 寫旅程
3. 定不下來的，就地標記並同時在 `spec.md` 的 `## Open Questions` 開一題
4. 跑 `check_spec.py`，第 17 項要過
5. 交回 `spec` 或 `example-mapping`

#### `formulation`

**職責**：把已對焦的 `spec.md` 翻成 Gherkin `.feature`：FR 變 `Rule:`、AC 變 `Example:`、編號變 tag，步驟套一組封閉文法

- **上游**：`spec`（`已對焦` 之後呼叫本 skill）
- **下游**：`plan`（覆蓋表沒有大洞時）；洞多回 `spec`
- **不做**：不訪談、不發明例子、不寫 `spec.md`、不寫 step definition、不讀 `plan.md`

1. 讀四樣東西，只做 `已對焦` 的 feature
2. FR → `Rule:`、AC → `Example:`、編號 → tag；歸象限、看三種結果、取場景名
3. 每個步驟套四種形狀之一
4. 句子只差資料的例子收成 `Scenario Outline`
5. 一則 story 一個 `.feature`
6. 跑 `check_spec.py`、產覆蓋表

#### `story-splitting`

**職責**：把一則太大的 story 切成數則**各自都還能交付價值**的小 story，並寫下切分依據

- **上游**：`discovery`（粗切）、`spec`（`已對焦` 前重切）；也可單獨使用
- **下游**：回到呼叫它的 skill；`spec.md` 改完由 `formulation` 寫成 `.feature`
- **不做**：不排順序、不估點數、不寫規格、`spec.md` 已經 `已對焦` 時不切

1. 過 INVEST（Small 除外），Valuable 不成立就先修、不切
2. 判斷是真的大（規則太多），還是規則講得太粗
3. 優先沿規則切
4. 規則切不動時，換九種模式
5. 有好幾種切法時，過規則零再套規則一、二
6. 逐則驗收是不是垂直切片

### PLAN

#### `plan`

**職責**：把一份已對焦的規格變成計畫與票：設計 Domain Model、Database Schema、API 契約三份並排對過，定技術決定與 seam，排依賴圖，垂直切片成一疊 tracer bullet，排好順序與 Checkpoint、能不能並行

- **上游**：`formulation`（`.feature` 寫好之後）
- **下游**：`implement <票號>`（一張票一個全新 session）；〈設計〉把 API 契約委派給 `api-and-interface-design`
- **不做**：不改 `spec.md`、不改 `.feature`、不動 `api/openapi.yaml`、不寫 migration、不訪談需求、不寫程式

1. Step 1 收集：讀 `spec.md`、`.feature` 與專案現況
2. 設計：seam 攤給使用者確認，再寫 Domain Model → Database Schema → API 契約，三份並排對一次
3. Step 2 依賴圖
4. Step 3 垂直切片、Step 4 寫任務、Step 5 排序與 Checkpoint
5. ASK 切法：攤給使用者，點頭才寫 `tasks.md`（`plan.md` 在前面各步已寫好），停，交 `/implement <票號>`

#### `api-and-interface-design`

**職責**：設計不容易被誤用的 API 與介面（REST、GraphQL、OpenAPI、程式碼層契約）。內容照 addyosmani/agent-skills 原樣保留，本專案只補職責卡、使用時機與 Skill Boundaries

- **上游**：`plan`（〈設計〉委派，輸入是 FR／AC、涵蓋的 `.feature`、既有 `api/openapi.yaml`）
- **下游**：交回 `plan`，由 `plan` 把 OpenAPI 片段寫進 `plan.md` 的 `## API 契約`
- **不做**：不寫檔、不動 `api/openapi.yaml`（單獨叫用時只給設計建議）；`disable-model-invocation`，只由 `plan` 委派或手動 `/api-and-interface-design`

### IMPLEMENT

#### `implement`

**職責**：把一張票做到綠：測試層先紅、委派 `behavior-driven-development` 轉綠與整理、閘門綠、勾驗收條件、合回主線

- **上游**：`plan`（票切好之後）
- **下游**：`behavior-driven-development`（第 3、4 步委派）；全部票做完之後是 `archive`
- **不做**：不改 `spec.md`、不改 `.feature` 的例子、不寫 `plan.md`、不改 `tasks.md` 裡這張票的票節以外的任何一行、不做 `plan.md` 的設計沒做的設計決定

1. 讀票，標 `@wip`
2. 測試層先紅：場景都接上 step definition，沒有 undefined step
3. green：一個場景一個場景轉綠
4. refactor：在綠燈下整理
5. 外迴圈綠，閘門綠
6. 收尾：勾驗收條件、寫缺口、判斷 `@review`
7. 合回主線（只在 worktree 裡做時）

#### `behavior-driven-development`

**職責**：在指定的 `.feature` 範疇裡，用 red／green／refactor 其中一個入口，一次一個 slice 把場景推進到那個入口的完成條件

- **上游**：`implement`（第 3 步委派 green、第 4 步委派 refactor，範疇是票的涵蓋場景），或使用者手動叫 `/behavior-driven-development`
- **下游**：回報給呼叫它的人，並建議下一次從哪個入口繼續
- **不做**：不產或改 `.feature`、不改 `spec.md`／`plan.md`、不翻狀態 tag、不勾票、不跨出指定範疇、一次不進兩個入口

1. 對齊範疇與入口
2. 盤點測試入口與既有抽象
3. 選下一個 slice
4. 照入口把這個 slice 做到完成條件
5. 還有下一個 slice 就回到 3，否則停
6. 回報

#### `archive`

**職責**：把一個做完的 feature 併進專案的真實來源：`docs/` 四份現況檔照它增修，再把它的目錄搬進 `specs/sleeping/`

- **上游**：`implement`（讓最後一張票完成的那次合回之後）
- **下游**：無；要修正時由 `spec` 喚醒，搬回 `specs/ongoing/`
- **不做**：不訪談、不改 `spec.md`／票節／`.feature`、不改程式碼、不改寫 `docs/CONTEXT.md` 既有的段落、不動區塊外的內容、不判斷實作對不對

1. 讀這個 feature 的 `spec.md`、`plan.md` 與 `docs/` 現況，列出要改的條目
2. 增修 `docs/` 四份檔，記修訂紀錄
3. `git mv` 到 `specs/sleeping/`，跟第 2 步同一個 commit
### 不在流程裡

#### `skill-rules`

**職責**：定本專案 skill 的房規（結構、frontmatter、規則關鍵字），並對現成 skill 逐項稽核

- **上游**：無（不在流程裡；寫或改任何 skill 之前讀）
- **下游**：無
- **不做**：不跑 eval 迴圈、不最佳化觸發率（都是 `skill-creator`）；職責重疊時不自己決定合併或劃界

稽核一個 skill：

1. 機械檢查：`audit_skill.py`
2. 人工判斷：照 `rules/quality-checklist.md`

寫一個新 skill：

1. 確認它不該是一段提示詞或一個 command
2. 先寫 `description`
3. 寫 `## Skill Boundaries`，職責跟既有 skill 重疊時問使用者
4. 寫流程、產物與 `## Verification`
5. 需要子目錄才建，並在 SKILL.md 指路
6. 跑稽核

## 目錄

```
ai-bdd/
├── .claude-plugin/
│   └── plugin.json       manifest：name、description、version、author
├── PLAN.md               要做什麼、先後順序
├── RESOURCES.md          參考來源：方法論的原文、借鏡的 skill 專案
├── docs/
│   ├── bdd.md            BDD 本身是什麼——外部來源轉述
│   ├── example-mapping.md  Example Mapping 原文轉述
│   ├── feature-mapping.md  Feature Mapping 原文轉述
│   ├── tdd.md            TDD 本身是什麼——外部來源轉述
│   ├── sdd.md            Spec-Driven Development 是什麼——外部來源轉述
│   ├── ai-sdlc.md        本專案的立場：為什麼是這六步
│   └── superpowers/      設計文件（specs/）與實作計畫（plans/）
├── skills/
│   ├── discovery/
│   ├── spec/
│   ├── feature-mapping/
│   ├── example-mapping/
│   ├── spec-by-example/
│   ├── formulation/
│   ├── story-splitting/
│   ├── plan/
│   ├── api-and-interface-design/
│   ├── implement/
│   ├── behavior-driven-development/
│   ├── archive/
│   └── skill-rules/
├── tests/                skill 腳本的 stdlib unittest
├── benchmark/
│   ├── cases/            每個情境一個目錄：input-*.md ＋ grading.md
│   └── skeleton/go/      Go ＋ godog 的試驗專案，每次跑時複製
└── README.md
```

skill 目錄是扁平的，步驟分組只放在上面的表與 `PLAN.md`：skill 名稱在執行時是全域、
扁平的，而 `claude plugin validate` 不會遞迴進巢狀目錄。

目錄名就是 skill 的識別名，frontmatter 的 `description` 決定它什麼時候被觸發——兩者都要寫對。
寫或改 skill 之前先讀 `skills/skill-rules/`。

`commands/`、`agents/`、`hooks/` 是 plugin 慣例目錄，放進去就會被自動載入；目前都還沒有，
有第一個真正的能力時再建。

## 驗證

```bash
claude plugin validate .                                          # 閘門
python3 skills/skill-rules/scripts/audit_skill.py skills/<name>   # 單一 skill 的機械稽核
python3 -m unittest discover -s tests                             # skill 腳本的單元測試
```

`--strict` 在這裡過不了，這是刻意的取捨：它只會因為一個警告失敗——plugin 根目錄的
`CLAUDE.md` 在安裝後不會載入。這個 repo 同時是 plugin 本身與開發它的地方，那份檔案在後者
有它的用處。閘門用一般的 `validate`。

## 安裝

這個 repo 是單一 plugin，不是 marketplace，要從某個 marketplace 的條目引用：

```json
{
  "name": "ai-bdd",
  "source": { "source": "github", "repo": "<owner>/ai-bdd" }
}
```

然後：

```bash
claude plugin install ai-bdd@<marketplace-name>
```
