# CAPTURE → PRD 流程重排：需求先寫成半成品，Discovery 把它問完整

**日期**：2026-10-02
**狀態**：待審
**前一版**：[`2026-09-14-pipeline-restructure-design.md`](./2026-09-14-pipeline-restructure-design.md)
**外部參照**：Claude Academy《The AI-native SDLC playbook》Stage 1 Capture as `intent.md`、
Stage 2 Requirements and design、Stage 3 Plan mode

這一版改的是**流程的分段方式與中間產物**，不是單一 skill 的格式。六步變七步，
`clarify-log.md` 消失，`prd.md` 回來但意思跟 9 月 11 日那版不同，`spec.md` 瘦成只裝
技術設計。

---

## 為什麼改

### 一、鏈的起點沒有產物

可追溯鏈寫的是 `商業目標 → FR → 例子 → Story → …`，但沒有任何一份檔案擁有
「商業目標」。`bdd-discovery` 的開場把業務目標、成功指標、範圍邊界問過一輪之後，
答案散在 `clarify-log.md` 三張表裡、`面向` 欄填 `—` 的幾列。需求方帶來的話只有
`input.md` 一份逐字副本，沒有一份是「需求方認可過的、用他自己的話寫的」陳述。

playbook 把這一格叫 `intent.md`：問題、預期結果、受影響的人與系統、限制、
開放問題。它是 proto-spec，由需求方自己改到認可，**是整條鏈第一個被人簽過名的東西**。

### 二、`clarify-log.md` 的存在理由只剩一個

9 月 14 日那版讓 CLARIFY 只留問答紀錄、不寫文件，理由是追問與書寫是相反的工作模式，
混在一起模型會跳過追問。那份 log 從此承擔三件事：記答案、記還開著什麼、防止重問。

但 Example Mapping 的紅卡本來就是這三件事：一張紅卡是一個開著的問題，答了之後
它變成藍卡或綠卡，而卡片留在地圖上就是「別再問」。把紅卡記在 PRD 裡、答案留在
原地，log 的三個功能一個都不少，卻少了一份要跟 PRD 對齊的檔案。

「模型會跳過追問直接寫」這個風險要另外擋，見 D4。

### 三、playbook 的 Stage 2 是一次做完 requirements 與 design，而 BDD 的 Discovery 剛好插得進那條縫

playbook 的鏈是 `intent.md → spec.md → plan.md`，Stage 2 把需求與設計壓成一個
session。那是給 PO 一個人跑的。本專案的主張是需求那一半不該用生成的，該用問的
——那正是 BDD 的 Discovery。所以 Stage 2 在這裡拆成兩半：

```
playbook    intent.md  →  spec.md（requirements ＋ design）  →  plan.md
ai-bdd      prd.md（draft）→ prd.md（complete）→ .feature ＋ spec.md（design）→ plan.md
            CAPTURE        DISCOVERY            SPEC                            PLAN
```

`prd.md` 就是 playbook 的 requirements 那一半，`spec.md` 是 design 那一半。
兩份的壽命不同：PRD 活得比實作久，design 實作一開始就過期。壽命不同的東西
不放同一個檔，這是 `PLAN.md` 既有的根層／slice 判準。

---

## 決策

### D1 · 七步

```
CAPTURE    /capture（bdd-capture）  → prd.md（狀態 draft）
DISCOVERY  bdd-discovery            → prd.md（狀態 complete）
SPEC       bdd-formulation          → features/*.feature
           bdd-spec                 → spec.md
PLAN       bdd-plan                 → plan.md
IMPLEMENT  未實作
VERIFY     未實作
REVIEW     未實作
```

CAPTURE 是新的一步。它不是 CLARIFY 的開場抽出去，是一個模式不同的步驟：
需求方說、agent 記。DISCOVERY 是 agent 問、需求方答。SPEC 是 agent 寫、需求方審。
三步三種模式，各自一個 skill 擁有，這是 9 月 14 日那版的分段原則，不變。

### D2 · `prd.md` 是一份從 draft 長到 complete 的活文件

一份檔案，路徑 `specs/<date>-<feature>/prd.md`，狀態只有兩個值：

| 狀態 | 誰寫的 | 條件 |
| --- | --- | --- |
| `draft` | CAPTURE | 需求方認可了開頭五節 |
| `complete` | `example-mapping`，在 DISCOVERY 結束、需求方點頭後 | **每一則 story 的紅卡都有答案** |

章節順序固定，空的節保留標題並寫為什麼空：

```markdown
# PRD: <feature>
**狀態**：draft　**最後更新**：<date>

## 問題              現在什麼在痛、誰在痛
## 預期結果          做成之後的樣子
## 角色              ### P-<n>，判準沿用 persona-definition.md
## 限制              外部限制、不准這樣做的事；技術限制也在這裡（見 D6）
## 範圍              In / Out
## 成功指標
## User Stories
### US-<n> · <slug>
身為 P-<k>，我想 …，以便 …
#### 規則            FR-<n>  一條規則一列
#### 例子            AC-<n>.<m>  掛在它說明的那條規則底下
#### 問題            Q-<n>  待答 / 已答（答案留在原地）
## 原文              需求方帶來的描述，逐字
```

**`## 原文` 取代 `input.md`。** 它仍然是「這句是需求方說的」唯一依據，只是搬進
同一份檔案。需求方帶著現成 PRD 來時，CAPTURE 把整份逐字放進這一節，開頭五節
從它抽取並請需求方確認，不改寫原文。

**規則與例子的編號跟現在一樣**：`FR-<n>` 全域、`AC-<n>.<m>` 掛在 `FR-<n>` 底下，
`.feature` 的 `@rule-<n>` / `@example-<n>.<m>` 回指它們。編號方案不動，
`bdd-formulation` 的 tag 契約就不用動。

**每一則 story 就是一張 Example Mapping 的地圖。** 黃卡是 `### US-<n>`，藍卡是
`#### 規則`，綠卡是 `#### 例子`，紅卡是 `#### 問題`。`example-mapping` 的四條診斷
直接數這三節：藍卡太多 → `story-splitting`；紅卡還在 → 未就緒；一條規則掛太多例子
→ 規則該拆；舉不出例子的規則 → 還在想像階段。

### D3 · CAPTURE 當書記，最多問四題

`bdd-capture` 的輸入是需求方的一段描述，或一份現成文件。它做四件事，順序固定：

1. 原文逐字存進 `## 原文`。
2. 從原文抽開頭五節。抽不出來的節留空並寫「原文沒提」。
3. 最多問四題，而且只在對應的節抽不出來時問：這解決誰的什麼問題？受影響的人與
   系統有誰？做成之後長怎樣？什麼明確不做？
4. 粗切 story：沿原文裡看得出的使用者目標切，每則只有一句 `身為…我想…以便…`，
   三個子節全空。切不出來就只開一則，slug 用 feature 名。

需求方改到認可，狀態寫 `draft`。

MUST NOT：替 story 填規則或例子。CAPTURE 看到的只有原文，原文沒有的規則是編的。
粗切 story 標明「暫定」，DISCOVERY 可以重切。

### D4 · DISCOVERY 只問紅卡；藍綠卡必須有來源

`bdd-discovery` 的輸入與輸出都是 `prd.md`。前沿模型保留：每個已定案的決定岔出
下一批可問的題，前沿空了就結束。改的是**單位與落點**：

- 前沿以 story 為單位算。一輪可以橫跨多則 story，但每張紅卡掛在它屬於的那則底下。
- 答案回來，直接改 `prd.md`：答案留在 `#### 問題` 原地標 `已答`，由它推出的規則
  寫進 `#### 規則`、例子寫進 `#### 例子`。
- 開場的「隱含假設機械掃描」對 `## 原文` 跑，掃出來的題開成紅卡。原文裡標題叫
  「假設」的每一條都是一張紅卡，不是前提。

**取代「CLARIFY 不寫檔」的守門規則：每一條 FR 與 AC 必須有來源。** 來源只有兩種：
`## 原文` 的某一句，或某張 `已答` 的紅卡。寫法是在那一列尾端加 `← Q-<n>` 或
`← 原文 §<n>`。沒有來源的規則是缺陷，`check_spec.py` 機械地抓。這條規則擋的是
同一件事——agent 補一個合理的答案然後繼續——但擋法從「不准寫」改成「寫了要指得出
誰說的」，因為後者可以驗，前者只能靠自律。

IMPORTANT：重問一張已經標 `已答` 的紅卡，比漏問更傷信任。續跑時先讀 `prd.md`
重建設計樹：`已答` 是定案的決定，`待答` 是開著的分支。

DISCOVERY 呼叫兩個 skill，不吸收它們：

| 時機 | 呼叫 | 做什麼 |
| --- | --- | --- |
| 一輪答案記完、想看形狀；前沿空了、要提就緒問句 | `example-mapping` | 攤某則 story 的三節、跑四條診斷。它讀 `prd.md`，不再讀 `spec.md` 或 `clarify-log.md` |
| 某則 story 藍卡太多，或 CAPTURE 的粗切明顯切錯 | `story-splitting` | 沿規則重切，改 `## User Stories` 的分組 |
| 需求方一次跳過了一疊紅卡，要再收斂 | `clarify`（原 `clarify-loop`） | 一題一答、附建議答案，答案回 `prd.md` 的紅卡 |

就緒值仍由 `example-mapping` 在需求方點頭後寫——它把 `prd.md` 的狀態翻成
`complete`。`bdd-discovery` 不碰那一行。

### D5 · SPEC 讀 complete 的 PRD，產兩份，不訪談

`bdd-formulation` 改讀 `prd.md`：`#### 規則` 變 `Rule:`，`#### 例子` 變
`Example:`，tag 契約不變。`check_spec.py` 的雙向稽核改成 `prd.md` ↔ `.feature`，
並加一項：FR／AC 的來源標記是否齊全（D4）。

`bdd-spec` 瘦身成只寫技術設計。`spec.md` 保留現在十六節裡的技術節：

```
## Document Overview          狀態、修訂歷史；核心關係人從 prd.md 抄
## Tech Preferences & Constraints
## Known Boundaries
## Implementation Decisions   模組、介面、API 契約、schema、seam
## Testing Decisions          seam ＋ 每個 AC 的測試層級 ＋ prior art
## Risks
## Out of Scope               留給 IMPLEMENT 的技術決定
## Further Notes
```

搬去 `prd.md` 的節：Background、Goal、Scope、Personas、User Stories、
Non-Functional Requirements、Assumptions / Constraints、Success Metrics、
Open Questions。`domain-model.md` 不動，仍是根層獨立檔案。

NEVER：`prd.md` 狀態不是 `complete` 時開始。有洞就退回 DISCOVERY，不在這裡問。

順序：先 `.feature` 再 `spec.md`。`Testing Decisions` 要替每個 AC 指定測試層級，
AC 與 Example 一對一，所以讀 `prd.md` 就夠，不依賴 `.feature`；但 seam 的決定要
看場景長什麼樣才下得了手，先有 `.feature` 比較不會憑空。

### D6 · 技術問題在 DISCOVERY 的最後一輪問，答案進 `## 限制`

`bdd-discovery` 現有的五個技術面向（seam、模組邊界、契約、資料、測試）保留，
但排成**業務紅卡清空之後的最後一輪**。答案寫進 `prd.md` 的 `## 限制`，當成
需求方給的外部限制，不是設計。`bdd-spec` 照著寫，不訪談。

否決的替代方案：playbook 的 flagged concerns——design 產出時把拿不準的地方標出來
交給 PO 去解。那是多人組織的作法，單人跑時用問的更直接。多人時要改回去，改的
只有這一節。

### D7 · `clarify-loop` 改名 `clarify`

用 `git mv skills/clarify-loop skills/clarify`。SKILL.md 的 `name:` 跟著改，
本文不動。`PLAN.md` 的步驟對照表、`bdd-discovery` 與 `example-mapping` 的
Skill Boundaries 裡的引用一起改。

### D8 · 文件與腳本

| 檔案 | 改什麼 |
| --- | --- |
| `PLAN.md` | 六步改七步；產物佈局整段重寫（順便移除早已不存在的 `brief.md`／`actor.md`）；步驟→skill 對照表加 `bdd-capture`、改 `clarify`；可追溯鏈前段改成 `prd.md → FR → AC → Scenario` |
| `docs/ai-sdlc.md` | 六步圖改七步，CAPTURE 標為「不是 BDD」；四條紀律改成 D1 的三種模式 ＋ PLAN |
| `CLAUDE.md` | pipeline 那段與四條紀律 |
| `skills/bdd-spec/scripts/status.py` | `git mv` 到 `bdd-discovery/scripts/`（它量的是 DISCOVERY 的進度）；改讀 `prd.md`，數每則 story 的紅卡；零 story 時 exit 1 |
| `skills/bdd-spec/references/persona-definition.md` | `git mv` 到 `bdd-capture/references/`，角色現在在 CAPTURE 第一次寫下 |
| `skills/bdd-discovery/references/clarify-log-format.md` | 刪除，由 `prd.md` 的格式參考檔取代（放 `bdd-capture/references/prd-format.md`，三個 skill 共讀） |
| `skills/bdd-formulation/scripts/check_spec.py` | 改讀 `prd.md`；加來源標記檢查 |
| `benchmark/README.md` | 及格標準裡提到 `clarify-log.md` 或 `input.md` 的地方改成 `prd.md` |
| `skills/bdd-formulation/references/state-tags.md`、`rule-taxonomy.md`、`skills/example-mapping/SKILL.md`、`skills/clarify/SKILL.md` | 引用 `clarify-log.md`／`input.md`／`spec.md` 作為輸入的地方改成 `prd.md`。這是 `grep -rl 'clarify-log\|input.md'` 列出的全部檔案，實作時再跑一次確認沒漏 |

---

## 產物佈局

```
specs/<date>-<feature>/
  prd.md           CAPTURE 寫 draft，DISCOVERY 寫到 complete。活得比實作久
  spec.md          SPEC    技術設計。實作開始就過期
  plan.md          PLAN
specs/
  domain-model.md  SPEC    不動
features/<story-slug>.feature   SPEC   一則 story 一個檔，位置由 runner 決定
docs/CONTEXT.md                 不動
```

消失的檔案：`input.md`（併入 `prd.md` 的 `## 原文`）、`clarify-log.md`（紅卡取代）。

---

## 三種模式，四條紀律

| 步驟 | 模式 | 只做 | 不做 |
| --- | --- | --- | --- |
| CAPTURE | 需求方說、agent 記 | 抄原文、抽五節、最多問四題、粗切 story | 不填規則與例子 |
| DISCOVERY | agent 問、需求方答 | 問紅卡，答過的才能變藍綠卡 | 不寫 Gherkin、不設計、不自己填答案 |
| SPEC | agent 寫、需求方審 | complete 的 PRD → `.feature` ＋ `spec.md` | 不訪談，有洞退回 DISCOVERY |
| PLAN | agent 切、工程師審 | tracer bullet ＋ blocking edges | 不設計 |

---

## 不在這一版

- IMPLEMENT／VERIFY／REVIEW 的 skill 仍未實作，這版不碰。
- `bdd-capture` 只處理「一個人有想法」與「帶著現成文件來」兩種入口。playbook 的
  ticket 與 incident 入口等有實跑需求再加。
- `prd.md` 狀態的 accept／reject 門不做。單人跑時門是「我讀過且它說的是我的意思」。

## 待實跑驗證

- 紅卡掛在 story 底下之後，**跨 story 的共用規則**會不會晚才浮現。對策是前沿一輪
  橫跨多則 story，但要跑過才知道夠不夠。
- CAPTURE 的粗切錯得多嚴重、DISCOVERY 重切的成本有多高。
- 來源標記（D4）在實際寫 PRD 時會不會太囉嗦到被省略。省略就等於沒有這條規則。
