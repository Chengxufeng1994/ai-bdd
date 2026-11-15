# ai-bdd

一個 Claude Code plugin：用 BDD（Behavior-Driven Development）把需求一路帶到程式碼——
先把需求問清楚，寫成規格與驗收旅程，再翻成可執行的 Gherkin，切成票，一張一張做到綠。

> **狀態：早期。** DISCOVERY、SPEC、PLAN、IMPLEMENT 已有 skill；REVIEW、VERIFY 尚未實作。
> 路線圖見 [PLAN.md](PLAN.md)（它是假設，不是決定）。

流程是六步，每一步裡用的方法是 BDD：

**DISCOVERY → SPEC → PLAN → IMPLEMENT → REVIEW → VERIFY**

| 步驟 | Skill |
| --- | --- |
| DISCOVERY | `bdd-discovery`、`clarify` |
| SPEC | `bdd-spec`、`example-mapping`、`bdd-spec-by-example`、`bdd-formulation`、`story-splitting` |
| PLAN | `bdd-plan`、`bdd-tasks` |
| IMPLEMENT | `bdd-implement`（一張票：測試層先紅、閘門、合回）、`bdd`（由 `bdd-implement` 委派 green／refactor，模型不自動觸發） |
| REVIEW／VERIFY | 尚未實作 |

`skill-rules` 不在流程裡：它是本專案撰寫與稽核 skill 的房規。哪幾步對應 BDD 的哪個實踐、
哪幾步是本專案自己加的，見 [docs/ai-sdlc.md](docs/ai-sdlc.md)。

## 中間產物：誰產生、誰接著用

每一步只讀上一步實際寫出來的檔案。`<feature>` 指 `specs/<date>-<feature>/`。

| 中間產物 | 位置 | 誰產生 | 誰接著用（用來做什麼） |
| --- | --- | --- | --- |
| `discovery.md` | `<feature>/` | `bdd-discovery`（原文逐字、粗切 story、問答）；`clarify` 補答待答題 | `clarify`：收斂待答題；`bdd-spec`：從已收斂的問答寫規格；`bdd-plan`：讀技術面向的答案 |
| `spec.md`（`待對焦`） | `<feature>/` | `bdd-spec` | `example-mapping`：攤四色卡、跑診斷、問就緒；`bdd-spec-by-example`：寫驗收旅程；`story-splitting`：重切 story |
| `spec.md` `## 回 DISCOVERY 補問` | `<feature>/spec.md` | `bdd-spec`、`bdd-spec-by-example`、`example-mapping`（未就緒時） | `bdd-spec`：退回 discovery，把 `discovery.md` 改回 `澄清中`；`bdd-discovery`：再問 |
| 切分紀錄（`## 切分：…`） | `specs/` 底下 | `story-splitting` | `bdd-spec`：照切法寫 story，FR 編號不重編 |
| 驗收旅程 `*.feature`（不執行） | `<feature>/acceptance/` | `bdd-spec-by-example` | 需求方：對著旅程判就緒；`example-mapping`：提就緒問句；`bdd-formulation`：只借用語 |
| `spec.md`（`已對焦`） | `<feature>/` | `example-mapping`（只寫需求方的決定） | `bdd-formulation`：寫可執行 `.feature`；`bdd-plan`：技術設計與切票 |
| `docs/CONTEXT.md` | 專案 `docs/` | `bdd-spec`（建立或追加） | `bdd-spec-by-example`、`bdd-formulation`：用詞對齊 |
| `domain-model.md` | `specs/` | `bdd-spec`（增修） | 下一個 feature 的 `bdd-spec`：聚合與不變條件 |
| 可執行 `.feature`（`@ready`／`@draft`） | `features/` | `bdd-formulation` | `bdd-plan`：每個 AC 的測試層級；`bdd-tasks`：每張票涵蓋哪些場景；`bdd-implement`：讓場景由紅變綠 |
| `plan.md` | `<feature>/` | `bdd-plan`（技術設計、要不要並行與接合點）；`bdd-tasks` 只在 `## 缺口` 加列 | `bdd-tasks`：拆成票；`bdd-implement`：讀技術節；`## 缺口` 的需求洞退回 `bdd-discovery` |
| `tasks/index.md`（順序、Blocked by、接合點） | `<feature>/tasks/` | `bdd-tasks` | `bdd-implement`：挑下一張票、判斷能不能跟別張同時做 |
| `tasks/<編號>-<slug>.md`（一張票一個檔） | `<feature>/tasks/` | `bdd-tasks` 建立；`bdd-implement` 勾驗收條件、寫 `## 缺口` | `bdd-implement`：一張票一個 session；`bdd-tasks` 再跑時把設計的洞搬進 `plan.md` |
| 產品程式碼、step definition、`@review` tag、票檔勾選 | 專案程式碼、`features/` | `bdd-implement` | REVIEW／VERIFY（尚未實作，目前人工對照） |

## 每個 skill：職責與流程

下面摘自每個 skill 的職責卡與流程總表；完整的讀寫檔案、規則與 Verification 在各自的 `skills/<name>/SKILL.md`。

### DISCOVERY

#### `bdd-discovery`

**職責**：把一段需求（想法或別人寫的 PRD）收成逐字原文與粗切 story，再一輪輪問到前沿空了，每個答案逐題記下

- **上游**：無（可單獨使用）；`bdd-spec`／`bdd-plan` 把它退回 `澄清中` 時再跑
- **下游**：`bdd-spec`（從這份紀錄寫 `spec.md`）
- **不做**：不寫 FR／AC／Gherkin／規則草稿、不自己填答案、不改寫 `## 原文`、不寫 `spec.md`／`plan.md`／`.feature`

0. 建檔：原文逐字、粗切 story、最多四題、改到需求方點頭
1. 開場其餘三件：掃原文找隱含假設、補核心關係人、識別角色——問出來的都開成『待答』題
2. 算前沿，過「值不值得問」；能自己查的事實派 sub-agent 查
3. 一次問完整個前沿，編號，每題附建議答案
4. 答案回來：四件事一起記、確認夠具體、看停下來的訊號 → 回到 2
5. 前沿空了：收尾核對，狀態翻 `已收斂`，交 `bdd-spec`

#### `clarify`

**職責**：把一疊未解的問題逐題問到收斂，交出「哪些解了、哪些還開著、還缺誰」

- **上游**：`bdd-discovery`（它的紅卡）；也可單獨使用
- **下游**：收斂後進 SPEC（`bdd-spec`）
- **不做**：不做決定、不寫 `spec.md`／`plan.md`／`.feature`、不改寫 `## 原文`、不改程式碼、不自己編題

1. 把手上的問題分成可解／阻塞／失效三堆
2. 讀已填的答案與跳過的原因，答定一題就立刻寫回
3. 可解的題一次一題問，四個選項＋逃生口
4. 給這一輪的進度表，對照「什麼時候停」決定要不要再跑一輪

### SPEC

#### `bdd-spec`

**職責**：把已收斂的問答寫成業務規格，帶到需求方說「就緒」，再交出去產 `.feature`

- **上游**：`bdd-discovery`
- **下游**：`example-mapping` ＋ `bdd-spec-by-example`（判就緒）→ `已對焦` 後 `bdd-formulation` → `bdd-plan`
- **不做**：不訪談、不標 `推論`、不自己寫 `已對焦`、不寫 `.feature`、不寫 seam／API／schema／測試層級

1. 寫 `spec.md`（`待對焦`）
2. 收尾檢查，必要時重切 story
3. 自檢 `check_spec.py`
4. 詞彙 → `docs/CONTEXT.md`、聚合 → `specs/domain-model.md`
5. `example-mapping` 診斷 → `bdd-spec-by-example` 寫旅程 → 就緒問句
6. 需求方答「就緒」→ `已對焦` → `bdd-formulation` 產 `.feature`
7. 不就緒、補問清單非空、或旅程有 `# [need clarification]` → 退回 discovery

#### `example-mapping`

**職責**：把一則 story 攤成四色卡地圖、跑四條診斷，把就緒問句交給使用者

- **上游**：`bdd-spec`（`spec.md` 寫完後呼叫）；也可單獨使用
- **下游**：就緒 → `bdd-spec` 接著呼叫 `bdd-formulation`；未就緒 → `bdd-spec`「退回 discovery」
- **不做**：不替使用者判就緒、不把紅卡問到收斂、不切 story、不寫 Gherkin、不產出 `example-map.md`

1. 認模式：流程內或單獨
2. 攤卡，在終端印出地圖
3. 跑四條診斷
4. 就緒判定：給判讀、提就緒問句、依回答寫狀態

#### `bdd-spec-by-example`

**職責**：把 `待對焦` 的規格寫成幾條可 demo、只用業務語言的長流程，讓需求方對著它判就緒

- **上游**：`bdd-spec` 第 5 步（`example-mapping` 跑完診斷之後）
- **下游**：`example-mapping` 提就緒問句；有 `# [need clarification]` → `bdd-spec`「退回 discovery」
- **不做**：不執行、不訪談、不改 FR／AC／狀態、不寫 `discovery.md`、不寫 `spec.md` 沒有的結果

1. 讀 `spec.md`，挑旅程
2. 照 `rules/journey-format.md` 寫旅程
3. 定不下來的，就地標記並同時寫進補問清單
4. 跑 `check_spec.py`，第 17 項要過
5. 交回 `bdd-spec` 或 `example-mapping`

#### `bdd-formulation`

**職責**：把已對焦的 `spec.md` 翻成 Gherkin `.feature`：FR 變 `Rule:`、AC 變 `Example:`、編號變 tag，步驟套一組封閉文法

- **上游**：`bdd-spec`（`已對焦` 之後呼叫本 skill）
- **下游**：`bdd-plan`（覆蓋表沒有大洞時）；洞多回 `bdd-spec`
- **不做**：不訪談、不發明例子、不寫 `spec.md`／`discovery.md`、不寫 step definition、不決定 seam、不讀 `plan.md`

1. 讀四樣東西，只做 `已對焦` 的 feature
2. FR → `Rule:`、AC → `Example:`、編號 → tag；歸象限、看三種結果、取場景名
3. 每個步驟套四種形狀之一
4. 句子只差資料的例子收成 `Scenario Outline`
5. 一則 story 一個 `.feature`
6. 跑 `check_spec.py`、產覆蓋表

#### `story-splitting`

**職責**：把一則太大的 story 切成數則**各自都還能交付價值**的小 story，並寫下切分依據

- **上游**：`bdd-discovery`（粗切）、`bdd-spec`（`已對焦` 前重切）；也可單獨使用
- **下游**：回到呼叫它的 skill；`spec.md` 改完由 `bdd-formulation` 寫成 `.feature`
- **不做**：不排順序、不估點數、不寫規格、`spec.md` 已經 `已對焦` 時不切

1. 過 INVEST（Small 除外），Valuable 不成立就先修、不切
2. 判斷是真的大（規則太多），還是規則講得太粗
3. 優先沿規則切
4. 規則切不動時，換九種模式
5. 有好幾種切法時，過規則零再套規則一、二
6. 逐則驗收是不是垂直切片

### PLAN

#### `bdd-plan`

**職責**：把已對焦的規格變成技術設計：驗收測試打在哪一層、動哪些模組、契約與 schema 長怎樣、每個 AC 在哪一層測

- **上游**：`bdd-formulation`（`.feature` 寫好之後）
- **下游**：`bdd-tasks`（把這份設計拆成票）
- **不做**：不改 FR／AC／story、不切票、不訪談、不寫程式

1. 一次讀完整批
2. 決定 seam
3. 寫技術節
4. 問要不要並行、列接合點

#### `bdd-tasks`

**職責**：把技術設計拆成一疊 tracer bullet：一張票一個檔，每張穿過所有層、自己可 demo，標明誰卡誰、會碰哪些接合點

- **上游**：`bdd-plan`（技術節與「要不要並行」寫好之後）
- **下游**：`bdd-implement <票號>`（一張票一個全新 session）
- **不做**：不改 `plan.md` 的技術節、不改 FR／AC／story、不派工、不寫程式

1. 一次讀完整批
2. 切票：最薄的完整路徑先，其餘依 FR
3. 標 Blocked by 與接合點
4. 攤出來確認
5. 寫 `tasks/`

### IMPLEMENT

#### `bdd-implement`

**職責**：把一張票做到綠：測試層先紅、委派 `bdd` 轉綠與整理、閘門綠、勾驗收條件、合回主線

- **上游**：`bdd-tasks`（票切好之後）
- **下游**：`bdd`（第 3、4 步委派）；做完之後是 REVIEW（未實作）
- **不做**：不改 `spec.md`、不改 `.feature` 的例子、不寫 `plan.md` 與 `tasks/index.md`、不做技術節沒做的設計決定

1. 讀票，標 `@wip`
2. 測試層先紅：場景都接上 step definition，沒有 undefined step
3. green：一個場景一個場景轉綠
4. refactor：在綠燈下整理
5. 外迴圈綠，閘門綠
6. 收尾：勾驗收條件、寫缺口、判斷 `@review`
7. 合回主線（只在 worktree 裡做時）

#### `bdd`

**職責**：在指定的 `.feature` 範疇裡，用 red／green／refactor 其中一個入口，一次一個 slice 把場景推進到那個入口的完成條件

- **上游**：`bdd-implement`（第 3 步委派 green、第 4 步委派 refactor，範疇是票的涵蓋場景），或使用者手動叫 `/bdd`
- **下游**：回報給呼叫它的人，並建議下一次從哪個入口繼續
- **不做**：不產或改 `.feature`、不改 `spec.md`／`plan.md`、不翻狀態 tag、不勾票、不跨出指定範疇、一次不進兩個入口

1. 對齊範疇與入口
2. 盤點測試入口與既有抽象
3. 選下一個 slice
4. 照入口把這個 slice 做到完成條件
5. 還有下一個 slice 就回到 3，否則停
6. 回報

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
│   ├── tdd.md            TDD 本身是什麼——外部來源轉述
│   ├── sdd.md            Spec-Driven Development 是什麼——外部來源轉述
│   ├── ai-sdlc.md        本專案的立場：為什麼是這六步
│   └── superpowers/      設計文件（specs/）與實作計畫（plans/）
├── skills/
│   ├── bdd-discovery/
│   ├── clarify/
│   ├── bdd-spec/
│   ├── example-mapping/
│   ├── bdd-spec-by-example/
│   ├── bdd-formulation/
│   ├── story-splitting/
│   ├── bdd-plan/
│   ├── bdd-tasks/
│   ├── bdd-implement/
│   ├── bdd/
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
