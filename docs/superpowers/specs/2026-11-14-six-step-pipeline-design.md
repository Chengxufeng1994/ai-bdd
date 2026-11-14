# 六步流程：discovery 只記問答、to spec 才寫規格、技術設計進 plan

**狀態**：待批　**日期**：2026-11-14

## 意圖

把流程改成需求方想要的形狀：

```
discovery → to spec → to plan → implement → review → verify
（brainstorming）                              （未實作）（未實作）
```

### 為什麼

需求方希望整體開發走得跟一次 brainstorming → 寫 spec → 寫多 task 的 plan → 實作 → review →
verify 的 session 一樣（2026-11-13 改 CAPTURE 那次就是這樣走的）。對照現況有三處不合：

1. **discovery 邊問邊寫規格**：現在 `bdd-discovery` 每答一題就把 FR／AC 寫進 `spec.md`。需求方要
   discovery 只做 brainstorming——問、記答案——規格等問完由下一步一次寫出。
2. **SPEC 產兩份文件**：`spec.md`（業務）＋ `design.md`（技術）。需求方要拿掉 `design.md`，技術
   決定改進 `plan.md`，跟 superpowers 的 plan 一樣，task 自己帶檔案、介面與測試。
3. **review 排在 verify 之後**：改成 review 先、verify 後——review 抓到的修完，最後一次驗證才算數。

## 目標

- `bdd-discovery` 只寫 `specs/<date>-<feature>/discovery.md`：原文、粗切 story、問答、狀態。
  **不寫 FR／AC**
- `bdd-spec` 讀 `discovery.md`，寫 `spec.md`（業務規格，FR／AC 行尾 `← Q-<n>`／`← PRD §n`），
  寫 `docs/CONTEXT.md`、`specs/domain-model.md`，然後呼叫 `bdd-formulation` 產 `.feature`
- `example-mapping` 對 `spec.md` 判就緒；`已對焦` 寫在 `spec.md`
- `design.md` 不存在。它的技術節併進 `plan.md`，由 `bdd-plan` 寫；seam 的徵詢移到 `bdd-plan`
- `bdd-implement` 從 `plan.md` 讀 seam、契約、測試層級
- `check_spec.py` 第 14 項照抓：`spec.md` 每條 FR／AC 要有來源，且 `← Q-<n>` 必須對到
  `discovery.md` 裡狀態 `已答` 的列，`← PRD §n` 必須對到 `discovery.md` `## 原文` 的段落
- `claude plugin validate .` 過；`python3 -m unittest discover -s tests` 全過；每個改過的 skill
  過 `audit_skill.py`

## 非目標

- **不寫 `bdd-review`、`bdd-verify`**。流程圖與 `PLAN.md` 列出這兩步，標「未實作」
- 不改前沿演算法、十六個面向、代演模式、`technical-probes.md` 的內容
- 不改 `.feature` 的格式、tag、位置；不改 `step-grammar.md`
- 不改 `check_spec.py` 十六項的檢查**內容**，只改它從哪份檔讀哪一節，外加第 14 項的跨檔對照
- 不處理「改既有行為」（truth／delta）——那是下一件事
- 不重寫 `docs/superpowers/` 底下的歷史設計與計畫

## 反目標（做出來也算失敗）

- `spec.md` 出現一條沒有來源的 FR／AC，或 `← Q-<n>` 指到 `discovery.md` 沒有的題——
  B1 的整個理由就是讓這條防線在「discovery 不寫規格」之後仍然機械可查
- discovery 在問的過程中偷寫規則（「先把 FR 草稿放 discovery.md」）。規則只在 to spec 寫
- `bdd-spec` 寫 `spec.md` 時順手問需求方新問題。發現洞就記進 `spec.md` 的
  `## 回 DISCOVERY 補問`，回 discovery
- session 中途斷掉後，下一次 discovery 重問已答的題

## 限制

- 文件繁中、commit 英文 `WHAT:`／`WHY:`／`HOW:`；搬檔用 `git mv`
- `skills/skill-rules/` 是 house style
- `docs/bdd.md`、`docs/tdd.md`、`docs/sdd.md` 不動（轉述外部來源）

## 考慮過的做法

| 做法 | 取捨 | 結果 |
| --- | --- | --- |
| A. discovery 照舊邊問邊寫 `spec.md`；to spec 只產 `.feature`＋技術設計 | 防線最強；但 discovery 不是需求方要的 brainstorming | 需求方選 B |
| **B1. discovery 只寫問答紀錄 `discovery.md`，to spec 從紀錄寫 `spec.md`**（本文件） | 規則只在一處寫；問答落檔所以 session 斷了接得上、`← Q-<n>` 有地方指 | **採用** |
| B2. discovery 完全不寫檔，問完一次交給 to spec | 最乾淨；session 斷就全沒了，`← Q-<n>` 無處可指，第 14 項得拿掉 | 捨棄 |
| 技術決定進 `spec.md` 的一節 | 一份講完；但技術決定第一個 commit 就過期，業務規則活得比實作久 | 需求方選進 `plan.md` |
| **技術決定進 `plan.md`**（本文件） | task 帶自己的檔案、介面、測試，實作者讀一張票就夠 | **採用** |

## 決定

### D1. 產物

```
specs/
├── domain-model.md               to spec   聚合、不變條件；跨 feature
└── <date>-<feature>/
    ├── discovery.md              discovery 原文、粗切 story、問答、狀態
    ├── spec.md                   to spec   業務規格：story、FR／AC（行尾來源）、角色、範圍
    └── plan.md                   to plan   技術設計（seam、模組、契約、schema、測試層級、風險）＋ 票

features/<story-slug>.feature     to spec（bdd-formulation）
docs/CONTEXT.md                   to spec   詞彙，跨 feature
```

### D2. `discovery.md` 的形狀

就是 `clarify` 單獨模式今天產的那個三節檔，加一節粗切 story：

| 節 | 內容 |
| --- | --- |
| `## Document Overview` | 狀態、最後更新、結束方式、版本修訂歷史（含「需求方確認建檔」一列）、核心關係人 |
| `## 粗切 story` | 每則一行：`US-<n> · <slug>` ＋ `身為…我想…以便…` ＋ 行尾 `← PRD §n`；標「暫定」。**沒有 FR／AC** |
| `## Open Questions` | 今天的五欄表 ＋ 小節，格式不變（`Q`／`問題`／`面向`／`狀態`／`答案`；小節的 `誰答的`、`歸屬`、`擋住`） |
| `## 原文` | 逐字、`§<n>` 編號，之後不得改寫 |

`擋住` 今天寫「FR-2」；之後 FR 還不存在，改寫成 story slug 或原文段落（`擋住：visitor-billing`）。
`歸屬` 的值域從 `design.md` 的章節名改成 `spec.md` 的章節名或 `plan.md` 的技術節名（D5）。

開場的四題（Background／Personas／Goal／Scope 空著時才問）照舊問，答案成為 `面向` 為 `—` 的
`Q-<n>` 列；**不再**由 discovery 寫成開頭幾節——那是 to spec 的事。Personas 五欄的問法照舊，
答案留在問答裡，to spec 寫進 `spec.md` 的 `## Personas`。

**狀態值**：`澄清中` → `已收斂`。`已收斂` ＝ 前沿空了、『待答』零列、結束方式 `收斂`。
`待對焦`／`已對焦` 不再出現在 `discovery.md`。

### D3. `spec.md` 的形狀

今天 `spec-format.md` 的骨架，**拿掉 `## Open Questions` 與 `## 原文`**（它們在 `discovery.md`），
加 `## 回 DISCOVERY 補問`（從 `design.md` 搬來，語意不變）。其餘節不變：Document Overview、
Background、Goal、Scope、Personas、User Stories（FR／NFR／Q 三分組）、NFR、Assumptions /
Constraints、Success Metrics。`#### Q` 分組仍放指標，指向 `discovery.md` 的 `Q-<n>`。

**狀態值**：`待對焦` → `已對焦`。`bdd-spec` 寫完是 `待對焦`；`example-mapping` 對它攤卡、跑
四條診斷、需求方答「就緒」才翻 `已對焦`。補問清單非空時 `example-mapping` 不提就緒問句。

**來源標記**：`← PRD §n` 指 `discovery.md` `## 原文` 的第 n 段；`← Q-<n>` 指 `discovery.md`
`## Open Questions` 的那一列。標記寫法不變，只是指向另一個檔。

### D4. skill 的變化

| skill | 讀 | 寫 | 變化 |
| --- | --- | --- | --- |
| `bdd-discovery` | 需求本身、`discovery.md`、回頭清單（D4c） | `discovery.md` | 拿掉「答案當場長成 FR／AC」；第 0 件事建 `discovery.md`；收尾翻 `已收斂`，下一步 `bdd-spec` |
| `clarify` | `discovery.md` | `discovery.md` | 檔名 `spec.md` → `discovery.md`；單獨模式的三節檔就是 `discovery.md`（少粗切 story 一節，第 0 件事補） |
| `story-splitting` | `discovery.md` 的粗切、`spec.md` 的 FR | `discovery.md` 的 `## 粗切 story`（discovery 時）或 `spec.md` 的 `## User Stories`（to spec 時） | 兩個呼叫點，各改各的檔 |
| `example-mapping` | `spec.md` | `spec.md` 的狀態 | 呼叫者從 `bdd-discovery` 改成 `bdd-spec`；「流程內·逼問中」時機消失 |
| `bdd-spec` | `discovery.md`（`已收斂` 或 `喊停`） | `spec.md`、`docs/CONTEXT.md`、`specs/domain-model.md` | **改寫**：從「寫技術設計」改成「從問答寫業務規格」，寫完做收尾檢查、呼叫 `example-mapping`，`已對焦` 後才呼叫 `bdd-formulation`（D4a） |
| `bdd-formulation` | `spec.md`（`已對焦`） | `.feature` | 只改「`spec.md` 由誰維護」的措辭；仍要求 `已對焦` |
| `bdd-plan` | `spec.md`、`.feature`、`discovery.md` | `plan.md` | **吸收 `design-format.md` 的技術節**與 seam 的四條規則、唯一一次徵詢 |
| `bdd-implement` | `plan.md`、`spec.md`、`.feature` | 程式碼 | `design.md` 的讀取點全改 `plan.md` |

`bdd-spec` 的 `references/design-format.md` 拆開：技術節（Tech Preferences & Constraints、Known
Boundaries、Implementation Decisions、Testing Decisions、Risks）進 `bdd-plan/references/plan-format.md`
的票之前；`domain-model.md` 骨架留在 `bdd-spec`；`## 回 DISCOVERY 補問` 進 `spec-format.md`。
`examples/minimal-design.md` 併進 `bdd-plan` 的範例或刪除，由實作者判斷。

`bdd-discovery/references/spec-format.md` 搬到 `bdd-spec/references/spec-format.md`（`git mv`），
因為擁有者換了；`bdd-discovery` 新增 `references/discovery-format.md` 描述 D2。

### D4a. `bdd-spec` 的呼叫順序

```
讀 discovery.md（已收斂 或 喊停）→ 寫 spec.md（待對焦）→ 收尾檢查（D4b）
→ example-mapping 攤卡、四條診斷、就緒問句 → 需求方答「就緒」→ 已對焦
→ bdd-formulation 產 .feature（@ready）
```

`bdd-formulation` 仍只收 `已對焦`，`state-tags.md` 的 `@ready` 定義不變。就緒問句答「不就緒」或
補問清單非空 → 不呼叫 `bdd-formulation`，回 discovery（D4c）。

### D4b. 原本 discovery 收尾的檢查搬到 `bdd-spec`

discovery 不寫 FR 之後，下面這些只看得到 FR 才能做的檢查改由 `bdd-spec` 在寫完 `spec.md`、
呼叫 `example-mapping` 之前做：

- 每個角色都被某條 FR 的主詞或 story 用到
- 一則 story 的 FR 全是 command 或全是 query，否則用 `story-splitting` 拆（一個 `.feature`
  只能是 `@command` 或 `@query`）
- 藍卡太多 → `story-splitting`
- `check_spec.py` 第 14 項為零

因此 **`bdd-spec` 在 `已對焦` 之前可以重切 story**（透過 `story-splitting`，改 `spec.md` 的
`## User Stories`），因為那時還沒有 `.feature` 掛在舊切法上。`已對焦` 之後要重切只能回 discovery
再走一輪，並記修訂。`discovery.md` 的 `## 粗切 story` 是起點，不跟著同步——`spec.md` 的
`## User Stories` 是 story 的定案位置。

discovery 自己的收尾保留：前沿空、『待答』零列、稽核分類是否窮盡、需求方確認。`example-mapping`
的「流程內·逼問中」時機刪除。

### D4c. 回 discovery 的路

三種回頭：`example-mapping` 得到「不就緒」、`spec.md` 的 `## 回 DISCOVERY 補問` 非空、`plan.md`
的 `## 缺口` 有需求洞。

- 回頭的那個 skill（`bdd-spec` 或 `bdd-plan`）把 `discovery.md` 狀態改回 `澄清中`，修訂歷史記一列
  「回頭：<原因>」，並告訴使用者跑 `bdd-discovery`
- `bdd-discovery` 續跑時**讀同目錄 `spec.md` 的 `## 回 DISCOVERY 補問` 與 `plan.md` 的 `## 缺口`**，
  每一條開成新的 `Q-<n>`（`待答`），小節寫來自哪一份；問完翻回 `已收斂`
- `bdd-spec` 再跑時**修補 `spec.md`，不重寫**：FR／AC 編號不重排（`.feature` 的 `@example-n.m`
  掛在上面），新增的接著編號，改動的記修訂歷史；補問清單裡被答掉的那幾條刪除；狀態回 `待對焦`，
  重跑 D4a 後半
- `.feature` 已存在時，`bdd-formulation` 照它今天處理「規則改了」的方式改場景

### D4d. 阻塞與喊停

`bdd-spec` 收兩種 `discovery.md`：`已收斂`，或結束方式 `喊停`（使用者決定帶著開著的題往下走）。
喊停時還『待答』的題在 `spec.md` 對應 story 的 `#### Q` 放指標——今天 `#### Q` 分組與第 10、11 項
檢查就是為這種情況存在的，所以照留。`阻塞`／`零進展` 的 `discovery.md` 拒收。

### D5. `plan.md` 的形狀

```
# <feature> — plan
**來源**、**seam**（照抄下面 Testing Decisions）
## Tech Preferences & Constraints
## Known Boundaries
## Implementation Decisions
## Testing Decisions       ← seam、每個 AC 的測試層級
## Risks
## 前置整理（prefactoring）
## 票
## 缺口
```

技術節的寫法、來源標記規則（三選一，`推論` 允許，但要列進 `## 缺口`）照搬 `design-format.md`。
票的切法規則不變（tracer bullet、一票一 session、blocking edges）。

**`bdd-plan` 也讀 `discovery.md`**：技術面向已答的題是技術節的原始素材（小節的 `歸屬` 指向
`plan.md` 的哪一節），核心關係人表（`決策權`）照抄進 `plan.md` 開頭。`plan.md` 裡的 `← Q-<n>`
目前沒有 script 檢查——這次不加，列為已知缺口，等 `bdd-review` 做鏈稽核時接手。

`bdd-plan` 的紀律從「不設計」改成「**不改需求**」：可以做技術決定，不能改 FR／AC／story；需求有洞
寫進 `## 缺口` 回 discovery。

### D6. 兩支 script

- `status.py`：讀 `specs/*/discovery.md` 的 `## Open Questions`。沒有 `discovery.md` → 「先跑
  bdd-discovery」
- `check_spec.py`：
  - 讀 `spec.md` 的形狀檢查（FR／AC、story、編號、來源標記）
  - `## 原文` 相關的兩項（存在、`← PRD §n` 的段落存在）改讀同目錄 `discovery.md`
  - 第 14 項之外加一條跨檔對照：每個 `← Q-<n>` 在 `discovery.md` 存在且狀態 `已答`。這算第 14 項
    的延伸，不是第十七項
  - 只有 `discovery.md`、還沒有 `spec.md` → 不是錯誤，印「還沒寫 spec，先跑 bdd-spec」並退出 0
    （跟今天「還沒有 `.feature`」的處理同一個道理：跑得到的先跑）——實作者可改成退出 1，
    但要在訊息裡講清楚
- **舊佈局防呆**擴充為兩種：
  1. `prd.md` 存在（10-02 佈局）：照今天的提示
  2. `spec.md` 存在、內含 `## Open Questions` 或 `## 原文`、同目錄沒有 `discovery.md`
     （11-13 佈局）：退出 1，提示「把 `## Open Questions` 與 `## 原文` 搬到 `discovery.md`，
     `design.md` 的技術節併進 `plan.md`」
  各附一個測試

### D7. 文件

- `docs/ai-sdlc.md`：六步圖（含 implement；review、verify 順序對調）；步驟表重寫 DISCOVERY／SPEC／
  PLAN 三列；「三種模式」改成 discovery 問、to spec 寫、plan 切＋設計；PLAN 的紀律改「不改需求」；
  說明為什麼 discovery 只記問答（B1 vs B2 的理由）
- `PLAN.md`、`CLAUDE.md`、benchmark 文件：同步
- `runs/go-parking-billing/`：手動拆成新佈局（`spec.md` 的 Open Questions 與原文搬到
  `discovery.md`，`design.md` 併進 `plan.md`），兩支 script 重跑過

## 交給實作者的

- `bdd-spec` 改寫後的長度與 references 拆法
- `minimal-design.md` 的去留
- `check_spec.py` 在「只有 `discovery.md`」時退出 0 還是 1
- 每份 SKILL.md 的措辭，以 `audit_skill.py` 過為準

## 驗證

1. `grep -rn "design\.md" skills/ tests/` 無結果（舊佈局防呆訊息除外）
2. `grep -rn "已對焦\|待對焦" skills/bdd-discovery/` 無結果
3. `bdd-discovery/SKILL.md` 沒有任何「寫 FR」「寫 AC」的指示
4. 兩支 script 對 `runs/go-parking-billing`（拆完後）退出 0；對 11-13 佈局的 fixture 退出 1 並給拆法
5. `claude plugin validate .`、`python3 -m unittest discover -s tests`、每個 skill 的 `audit_skill.py` 全過
