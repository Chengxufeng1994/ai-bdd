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
| IMPLEMENT | `bdd-implement` |
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
| `plan.md` | `<feature>/` | `bdd-plan`（技術設計）；`bdd-tasks` 只在 `## 缺口` 加列 | `bdd-tasks`：拆成票；`bdd-implement`：讀技術節；`## 缺口` 的需求洞退回 `bdd-discovery` |
| `tasks/index.md`（順序、Blocked by、共用檔案） | `<feature>/tasks/` | `bdd-tasks` | `bdd-implement`：挑下一張票、判斷能不能跟別張同時做 |
| `tasks/<編號>-<slug>.md`（一張票一個檔） | `<feature>/tasks/` | `bdd-tasks` 建立；`bdd-implement` 勾驗收條件、寫 `## 缺口` | `bdd-implement`：一張票一個 session；`bdd-tasks` 再跑時把設計的洞搬進 `plan.md` |
| 產品程式碼、step definition、`@review` tag、票檔勾選 | 專案程式碼、`features/` | `bdd-implement` | REVIEW／VERIFY（尚未實作，目前人工對照） |

## 目錄

```
ai-bdd/
├── .claude-plugin/
│   └── plugin.json       manifest：name、description、version、author
├── PLAN.md               要做什麼、先後順序
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
