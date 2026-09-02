# CLARIFY／SPEC／PLAN 三步重切 實作計畫

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 `PLAN.md` 定下的三條紀律（CLARIFY 只問、SPEC 只綜合、PLAN 只切）落實到三個已實作的 skill 上，並用實測場證明新分工跑得通。

**Architecture:** 依「格式先定 → 實跑驗證 → 再往上游修」推進。先把 `spec.md` 的骨架定出來（它是這條鏈上唯一全新的介面），拿實測場現有的產物硬跑一次；跑不出來的欄位就是 CLARIFY 該補問的，那份缺口清單直接變成 `bdd-clarify` Pass 3 的追問面向。最後才重寫 `bdd-plan`，因為它讀的是前兩步定案後的產物。

**Tech Stack:** Markdown（`SKILL.md` 與 `references/`）、Python 3 標準庫（稽核腳本，無第三方相依）、Go ＋ godog（實測場）

**Spec:** `/home/benny/Workspace/vivotek/ai-bdd/PLAN.md`

## Global Constraints

這些是全域規則，每個任務的要求都隱含包含這一節。

- **skill 本文 500 行以內**（`skills/skill-rules/rules/structure.md` S10）。接近時把只有特定情況才需要讀的段落下放到 `references/`。
- **子目錄只能是** `rules/` `references/` `examples/` `scripts/`（S3），且**每個檔案都必須被 `SKILL.md` 或其他被引用的檔案用文字指到**（S5）——沒指路的檔案模型讀不到。
- **章節順序**：frontmatter → `# 標題` → `## 使用時機` → `## Skill Boundaries` → 工作流程 → 輸出格式。
- **`description` 壓成觸發條件**，中英文觸發詞都要有。它每個 session 都付費，本文只在觸發時付。
- **skill 只寫入使用端 repo 的 `docs/bdd/`**，`.feature` 除外（位置由 runner 決定）。**不修改專案既有的任何檔案**。
- **產物用繁體中文**；Gherkin 關鍵字一律用英文（zh-TW 方言沒有「規則」，`規則:` 會被靜默吞成散文）。
- **實測場的 slice slug 是 `workout-tracking`**，路徑 `lab/go/skeleton/docs/bdd/workout-tracking/`。
- **不新增 pytest 或任何測試框架。** 這個 repo 現在沒有測試基礎設施；兩支腳本的驗證方式是**對實測場實跑並比對輸出**——那是真實資料，確定性且可重跑。這是已知取捨，不是遺漏。
- **每個任務結束都 commit**，訊息用 `<type>: <description>`（type：feat／fix／refactor／docs／test／chore）。

---

## File Structure

| 檔案 | 責任 | 動作 |
| --- | --- | --- |
| `skills/bdd-spec/references/spec-format.md` | `spec.md` 的骨架與每節判準 | **新增** |
| `skills/bdd-spec/SKILL.md` | SPEC 步驟入口：挑 feature → 寫 `.feature` → 決定 seam → 寫 `spec.md` → 增修 `domain-model.md` | 修改（加三步） |
| `skills/bdd-spec/scripts/check_spec.py` | 雙向覆蓋稽核 | 修改（換來源） |
| `skills/bdd-clarify/references/technical-probes.md` | Pass 3 的技術追問面向清單 | **新增** |
| `skills/bdd-clarify/SKILL.md` | CLARIFY 三個 pass ＋ 就緒判定 | 修改（加分批、加 Pass 3、拿掉 map 產物） |
| `skills/bdd-clarify/references/map-format.md` | 問題檔格式 | 修改（狀態列加欄位、拿掉 map 章節） |
| `lab/go/skeleton/docs/bdd/domain-model.md` | 跨批次的聚合與不變條件 | **新增**（從 `plan.md` §2 抽出） |
| `lab/go/skeleton/docs/bdd/glossary.md` | 只剩詞彙 | 修改（`Shared N` 規則搬家） |
| `skills/bdd-clarify/scripts/status.py` | 澄清進度 | 修改（覆蓋改從問題檔算） |
| `skills/bdd-plan/references/plan-format.md` | `plan.md` 的骨架 | **重寫**（技術設計 → 票） |
| `skills/bdd-plan/SKILL.md` | PLAN 步驟入口 | **重寫** |
| `lab/go/skeleton/docs/bdd/workout-tracking/` | 實測場的新佈局 | **新增**（六個 story 目錄併進來） |

---

# Phase 1 — SPEC：定 `spec.md` 的格式，用實測場驗證

## Task 1: `spec.md` 的骨架

**Files:**
- Create: `skills/bdd-spec/references/spec-format.md`
- Modify: `skills/bdd-spec/SKILL.md`（`## 產物格式` 段加指路）

**Interfaces:**
- Consumes: 無（這是鏈上的第一個新介面）
- Produces: `spec.md` 的七個章節標題，其中 `## Example Mapping` 與其下的 `### <story-slug>` 是 Task 4 的 `check_spec.py` 要解析的格式契約

- [ ] **Step 1: 寫 `skills/bdd-spec/references/spec-format.md`**

````markdown
# `spec.md` 的骨架

照抄即可。讀這份的時機：步驟 9，準備寫 `spec.md` 的時候。

一批一份，路徑 `docs/bdd/<slice-slug>/spec.md`。一個 slice ＝ 一次可交付的價值，
裡面裝好幾則 story。

沒有內容的節**保留標題並寫為什麼空**，不要刪掉。空著的節跟被刪掉的節在檔案上
長得一樣，而它們的意思完全不同。

---

```markdown
# <slice 名稱>

**Slice**: <slice-slug>
**日期**: <YYYY-MM-DD>
**涵蓋**: <story-slug>、<story-slug>⋯（N 則）
**跳過**: <story-slug> —— <理由>（未就緒，或這一批不做）

## Problem Statement

使用者遇到的問題，從使用者的角度寫。不寫解法。

## Solution

解法，一樣從使用者的角度。不寫技術。

## User Stories

編號清單，每則一句：

1. 作為 <角色>，我要 <能力>，以便 <價值>

角色名用 `actor.md` 的，不另創同義詞。

## Example Mapping

規則與例子的**定版編號**。編號在這裡誕生，`.feature` 的 tag 回指它們。

MUST: 每則 story 一個 `### <story-slug>` 區塊，slug 與 `features/<slug>.feature`
一致——那是這條鏈唯一的接縫，也是 `check_spec.py` 比對的鍵。

MUST: 例子一律寫成 `- Example N.M <一句話>`。編號在 story 內從 1 開始，
不跨 story 連號。

### <story-slug>

**Story**: 作為 <角色>，我要 <能力>，以便 <價值>

#### Rule 1. <規則敘述>
- Example 1.1 <具體例子，含實際資料>
- Example 1.2 <具體例子>

#### Rule 2. <規則敘述>
- Example 2.1 <具體例子>

## Implementation Decisions

已經決定好的技術事項。每一條要指得出它來自哪個已答的問題或哪條規則。

- **動哪些模組**：新增什麼、改什麼
- **介面**：那些模組對外的形狀（函式名、參數、回傳型別）
- **API 契約**：方法、路徑、request、response、錯誤碼
- **Schema**：資料表與欄位
- **架構決定**：分層、交易邊界、同步或事件驅動
- **互動**：跨模組的呼叫順序

MUST NOT: 寫具體檔案路徑與程式碼片段——它們過期得比什麼都快。
例外：prototype 產出的、比散文更精確地編碼了某個決定的片段（狀態機、schema、
型別形狀），可以內嵌並註明來自 prototype，只留決定的部分。

MUST NOT: 在這裡做新決定。推不出來的寫進 `## Out of Scope` 的「回 CLARIFY」欄，
不要順手決定掉——那正是規格失效的方式。

## Testing Decisions

**Seam**（驗收測試打在哪一層）

| seam | 是既有的還是新的 | 為什麼是這一層 |
| --- | --- | --- |
| `application.WorkoutService` | 既有 | 規則全在 domain 與 application，HTTP 只是轉發 |

判準：優先用既有 seam、取最高的那一層、**整個變更的理想數量是一個**。
seam 一旦寫在這裡就會往下傳：IMPLEMENT 照著打，REVIEW 把沒人同意過的 seam
當成 finding。

**每個 Example 的測試層級**

一個**場景區塊**一行。`Scenario Outline` 算一個，編號寫範圍。

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 2.1 進度不可倒退 | 單調遞增判斷 | ✗ |
| 6.1 查詢回傳完整內容 | — | ✓ 要讀得到前一次寫入 |

`需要資料庫` 的判準：這個場景要驗的行為**本身跨越一次操作的邊界**嗎？
不是「有沒有資料要存」——多數場景在同一次操作內就驗完了。

`—` 代表沒有內迴圈（純轉發、純投影），那是正常的。

**Prior art**

專案裡已經有的同類測試，寫路徑與一句話。新測試跟著它們的形狀寫，不另創。

## Risks

規格沒說、但實作一定會撞到的。每一條指出從哪條規則長出來。
跟 Out of Scope 的缺口不同：缺口是**還不知道要做什麼**，風險是**知道要做什麼
但做對很難**。

| 風險 | 從哪條規則長出來 | 影響什麼 |
| --- | --- | --- |
| 併發：兩個請求同時通過「最多一筆進行中」 | log Rule 1 | schema（唯一索引）、實作順序 |

SHOULD NOT: 把解法寫死。「用 Redis 還是 DB 行鎖」需要知道實際流量與既有基礎
設施，而那些不在規格裡。寫「兩種都可以，看部署形態」比寫「用 Redis」有用——
後者看起來是決定，實際上是猜的。

## Out of Scope

分兩類，因為**誰能回答**不同。

**明確拒絕的**

| 拒絕了什麼 | 為什麼 |
| --- | --- |

**回 CLARIFY 補問**（規格沒講，推不出來）

| 缺口 | 影響什麼定不下來 |
| --- | --- |

**留給 IMPLEMENT**（技術決定，不需要規格授權）

| 決定 | 為什麼現在不決定 |
| --- | --- |
```

---

## 為什麼 Example Mapping 要跟 `.feature` 重複

`check_spec.py` 的雙向覆蓋稽核靠比對兩份獨立的表述工作：這裡有而 `.feature`
沒有 ＝ 漏做；`.feature` 有而這裡沒有 ＝ **憑空發明**。併成一份，這個檢查就變成
拿檔案跟自己比，恆真。

而「發明」是最沒有人會懷疑的那種錯：漏一條會被數字抓到，多一條看起來只是很完整。

## 為什麼一批一份，不是一則 story 一份

試過「每 feature 一份 ＋ 一份共用檔」，壞成兩種：共用檔把實質吸走（各檔的
schema 一節全退化成「見共用檔」），而且兩種檔案各用一套節號、意思只對得上一半
——**部分對齊比完全不對齊更危險**，讀的人會假設對齊然後在三個節上讀錯。

API 契約與 schema 天生跨 story，只有並排才看得見衝突。

---

## `domain-model.md` 的骨架

路徑 `docs/bdd/domain-model.md`——**根層，跨批次**。它比 `spec.md` 活得久：
spec 是拋棄式的快照，實作一開始就會過期；domain model 要留下來。

每一批只**增修**，不重寫。改了既有條目要寫「原本是什麼、為什麼改」。

```markdown
# Domain Model

**更新於**: <YYYY-MM-DD> · **最後一批**: <slice-slug>

## 聚合

### <聚合名>

**是什麼**: 一句話
**邊界**: 哪些東西跟它一起變更、哪些不是
**來源**: <slice-slug> Rule N

**不變條件**（任何時刻都必須為真）

| 條件 | 來源 | 被誰保證 |
| --- | --- | --- |
| 一個使用者最多一筆進行中的訓練 | log-a-workout Rule 1 | Workout 聚合 ＋ DB 唯一索引 |

## 跨聚合的規則

不屬於任何單一聚合的。寫明為什麼放不進去。

## 修訂紀錄

| 日期 | 批次 | 改了什麼 | 為什麼 |
| --- | --- | --- | --- |
```

**不變條件與 `.feature` 的 `Rule:` 有什麼不同**：`Rule:` 是**某次操作要滿足
的行為**，不變條件是**任何時刻都必須為真的狀態**。前者用場景驗，
後者多半要靠型別、建構子或資料庫約束保證——它們的實作位置不同，所以要分開寫。
````

- [ ] **Step 2: 在 `skills/bdd-spec/SKILL.md` 的 `## 產物格式` 段開頭加指路**

在 `## 產物格式` 這一行底下、`**關鍵字用英文⋯**` 那段之前插入：

```markdown
本步驟產兩份東西，格式各有一份參考檔：

| 產物 | 格式 |
| --- | --- |
| `.feature` | 本節（關鍵字、狀態 tag、型別 tag、方言陷阱） |
| `docs/bdd/<slice-slug>/spec.md` | `references/spec-format.md` |
| `docs/bdd/domain-model.md`（增修） | `references/spec-format.md` 末節 |
```

- [ ] **Step 3: 跑機械稽核，確認沒有孤兒檔案**

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec`
Expected: 離開碼 0，輸出不含 `S5`（孤兒檔案）。若報 S5 指到 `spec-format.md`，代表 Step 2 的指路沒寫進去。

- [ ] **Step 4: Commit**

```bash
git add skills/bdd-spec/references/spec-format.md skills/bdd-spec/SKILL.md
git commit -m "feat(bdd-spec): define the spec.md skeleton and its seven sections"
```

---

## Task 2: `bdd-spec` 加 seam、`spec.md` 與 `domain-model.md` 三步

**Files:**
- Modify: `skills/bdd-spec/SKILL.md`（`## 流程` 段尾端加步驟 8、9、10；步驟 6 加一句；`## Skill Boundaries` 與 `## 完成後` 更新）

**Interfaces:**
- Consumes: Task 1 的 `references/spec-format.md`
- Produces: `bdd-spec` 執行後在 `docs/bdd/<slice-slug>/spec.md` 留下檔案

**為什麼是「加在後面」而不是插進中間：** 現有的 `references/` 與本文互相用步驟編號指路（`coverage-report.md` 從步驟 6 指來、骨架範例寫在步驟 4）。重排編號會讓那些引用靜默指到別的步驟——跟 Rule 編號不可重排是同一個道理。

- [ ] **Step 1: 在步驟 6 的標題底下加一句執行時機**

把 `### 6. 稽核，然後才算完成` 底下第一行改成：

```markdown
### 6. 稽核，然後才算完成

**這一步在步驟 9 之後跑**——它同時稽核 `.feature` 與 `spec.md` 的 Example
Mapping 段，兩份都要存在才比對得起來。

```bash
python3 <skill>/scripts/check_spec.py <專案根>
```
```

- [ ] **Step 2: 在 `### 7. 寫檔` 之後、`---` 之前，新增步驟 8**

```markdown
### 8. 決定 seam —— 驗收測試打在哪一層

`.feature` 決定驗什麼，seam 決定**打進系統的哪一層**：HTTP handler？
application service？domain？這個決定會長成 step definition 的形狀，
所以它屬於規格，不屬於實作。

三條規則：

| 規則 | 為什麼 |
| --- | --- |
| 優先用既有 seam | 每多一個 seam 就多一套測試替身與資料建構 |
| 取**最高**的那一層 | 越高越接近使用者真正做的事，也越不綁實作細節 |
| 整個變更的理想數量是**一個** | 兩個 seam 代表這批行為的入口不只一個，多半是切分沒切乾淨 |

**Ask user: 把 seam 攤出來確認再寫。** 這是本 skill 唯一一次徵詢——
其餘全部只綜合已有的答案。之所以是例外：seam 要看過 `.feature` 才推得出來，
而它推錯的話，底下所有 step definition 都打在錯的高度，改起來是整批重寫。

IMPORTANT: seam 一旦寫進 `spec.md` 就往下傳——IMPLEMENT 照著打，
REVIEW 把「沒人同意過的 seam」當成 finding。這個綁定是**間接的**，
穿過 `spec.md` 這份文件；所以這一節寫不寫，決定了後面兩步抓不抓得到。

### 9. 寫 `spec.md`

一批一份，路徑 `docs/bdd/<slice-slug>/spec.md`。骨架照抄
`references/spec-format.md`。

MUST NOT: 在這一步做新決定。**本 skill 只綜合已經有答案的東西。**
推不出來的寫進 `## Out of Scope` 的「回 CLARIFY 補問」欄，不要順手決定掉。

判準：指著 `spec.md` 的任何一句話，說得出它來自哪個已答的問題、哪條規則、
或哪個 `.feature` 的哪一段嗎？說不出來的就是憑空長出來的，那是缺陷。

### 10. 增修 `docs/bdd/domain-model.md`

根層一份，跨批次累積。這一批長出來的聚合與不變條件加進去。

MUST: 只增修，不重寫。動到既有條目要在修訂紀錄寫「原本是什麼、為什麼改」——
那張表記錄的是這個領域被理解的過程，而下一批的人靠它判斷某個決定還算不算數。

MUST NOT: 寫進 `spec.md`。`spec.md` 是這一批的快照，實作一開始就會過期；
domain model 要活得比它久。塞進去等於陪葬。

判準（哪些該進來）：這條規則**跨批次仍然成立**嗎？只在這一批的情境下為真的，
留在 `spec.md`。
```

- [ ] **Step 3: 更新 `## Skill Boundaries`**

把 `- 產出是 `.feature`⋯` 那條改成：

```markdown
- 產出是 `.feature` ＋ `docs/bdd/<slice-slug>/spec.md`，**不是** `openapi.yaml`、migration 或任何可執行的檔案
- 要把 `spec.md` 拆成可執行的票 → 改用 `bdd-plan`
```

- [ ] **Step 4: 更新 `## 完成後`**

把三件事改成四件：

```markdown
告訴對方四件事：

1. 產了哪些 `.feature`、各幾個場景、**幾個不重複步驟樣板**
2. **覆蓋表**——哪些例子還沒有場景，為什麼
3. **seam 是哪一層**、`domain-model.md` 這一批新增了哪些聚合與不變條件
4. `spec.md` 的 `## Out of Scope` 裡「回 CLARIFY 補問」有幾條；缺口多 → 回 `clarify-loop`，缺口少 → `bdd-plan` 切票
```

- [ ] **Step 5: 確認本文沒超過 500 行**

Run: `wc -l skills/bdd-spec/SKILL.md`
Expected: ≤ 500。超過的話，把步驟 8 的三條 seam 規則下放到 `references/spec-format.md`，本文只留一句指路。

- [ ] **Step 6: 機械稽核**

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec`
Expected: 離開碼 0。

- [ ] **Step 7: Commit**

```bash
git add skills/bdd-spec/SKILL.md
git commit -m "feat(bdd-spec): add the seam decision, spec.md and domain-model.md steps"
```

---

## Task 3: 實測場硬跑 —— 產出第一份 `spec.md`

這是這份計畫的**驗證核心**。`PLAN.md` 自己寫了「這次重規劃是在紙上改的，所以實測場重跑不能跳」——就是這一步。

**Files:**
- Create: `lab/go/skeleton/docs/bdd/workout-tracking/spec.md`
- Move: `lab/go/skeleton/docs/bdd/<六個 story 目錄>/questions/*.md` → `lab/go/skeleton/docs/bdd/workout-tracking/questions/`
- Create: `/tmp/claude-1000/-home-benny-Workspace-vivotek-ai-bdd/cfe2d8e3-7ecc-4246-9f8a-9b137e4e99bd/scratchpad/gaps.md`（缺口清單，Task 5 的輸入，不進 repo）

**Interfaces:**
- Consumes: Task 1 的骨架；實測場現有的六份 `example-mapping.md`、六個 `.feature`、`plan.md` §1–5
- Produces: `docs/bdd/workout-tracking/spec.md`，其 `## Example Mapping` 段有六個 `### <story-slug>` 區塊，slug 與 `features/*.feature`（不含 `version.feature`）一一對應

- [ ] **Step 1: 建 slice 目錄，把六個 story 的問題檔集中**

```bash
cd lab/go/skeleton/docs/bdd
mkdir -p workout-tracking/questions
for d in custom-exercise-library edit-a-logged-workout log-a-workout \
         session-training-volume week-over-week-progress workout-history; do
  git mv "$d"/questions/*.md workout-tracking/questions/ 2>/dev/null || \
    mv "$d"/questions/*.md workout-tracking/questions/
done
ls workout-tracking/questions | wc -l
```

Expected: 26 個檔（六則 story 的問題檔全部集中）。實際數字以現場為準，記下來——Step 5 要用。

- [ ] **Step 2: 把六份 `example-mapping.md` 的規則與例子併進 `spec.md` 的 Example Mapping 段**

照 `skills/bdd-spec/references/spec-format.md` 的骨架寫
`lab/go/skeleton/docs/bdd/workout-tracking/spec.md`。

`## Example Mapping` 段的每則 story 一個區塊，**編號原封不動照搬**——
`.feature` 的 `@example-N.M` tag 已經指著它們，重編等於讓那些 tag 靜默指錯：

```markdown
## Example Mapping

### log-a-workout

**Story**: <照抄 example-mapping.md 檔頭的 Story 句>

#### Rule 1. <照抄>
- Example 1.1 <照抄>
```

MUST NOT: 改任何 Rule 或 Example 的編號與敘述。這一步是搬家，不是重寫。

- [ ] **Step 3: 把現有 `plan.md` 的 §1–5 翻成 Implementation Decisions 與 Testing Decisions**

來源對照：

| 現有 `plan.md` 的節 | 搬進 `spec.md` 的 |
| --- | --- |
| §1 API 操作（含驗證約束、錯誤訊息） | Implementation Decisions · API 契約 |
| §2 Domain 型別（含 domain service、port/out） | Implementation Decisions · 介面、架構決定 |
| §3 Schema | Implementation Decisions · Schema |
| §4 測試分層 | Testing Decisions · 每個 Example 的測試層級 |
| §5 技術風險 | Risks |
| §7「回 CLARIFY 補問」 | Out of Scope · 回 CLARIFY 補問 |
| §7「留給 IMPLEMENT」 | Out of Scope · 留給 IMPLEMENT |

- [ ] **Step 4: 補 `plan.md` 裡沒有、但骨架要求的欄位——這一步的產出就是缺口清單**

現有 `plan.md` **沒有**下列欄位，逐一嘗試填寫：

| 欄位 | 現有產物裡有沒有答案 |
| --- | --- |
| Testing Decisions · **Seam** | 沒有。看 `lab/go/skeleton/` 的分層，提一個並記下判斷依據 |
| Testing Decisions · **Prior art** | 部分。掃 `lab/go/skeleton` 底下既有測試 |
| Implementation Decisions · **動哪些模組** | 沒有。`plan.md` 有型別但沒有模組邊界 |
| Implementation Decisions · **介面** | 只有 `port/out` 那一小塊 |
| Implementation Decisions · **互動** | 沒有 |

**填不出來的每一格寫進 scratchpad 的 `gaps.md`**，格式：

```markdown
| 欄位 | 為什麼填不出來 | 該問誰 | 這是哪一類技術問題 |
| --- | --- | --- | --- |
| 介面 · WorkoutService | 規格只說行為，沒說這批行為要收在幾個服務裡 | 開發者 | 模組邊界 |
```

最後一欄是 Task 5 的原料——**技術追問面向清單就是從這一欄長出來的**。

- [ ] **Step 5: 從 `plan.md` §2 抽出跨批次的部分，寫 `docs/bdd/domain-model.md`**

現有 `plan.md` §2「Domain 型別」裡混了兩種東西，這一步把它們分開：

| 屬於 | 判準 | 去哪 |
| --- | --- | --- |
| 聚合與不變條件 | 跨批次仍然成立 | `docs/bdd/domain-model.md`（根層） |
| 型別欄位、read model 投影、port/out | 只在這一批的情境下成立 | `spec.md` 的 Implementation Decisions |

實測場最明顯的例子：「一個使用者最多一筆進行中的訓練」是不變條件（任何時刻都
必須為真，靠唯一索引保證）；「`VolumeCalculator` 要讀動作型別與體重快照」是這一批
的實作決定。

Expected: `lab/go/skeleton/docs/bdd/domain-model.md` 存在，且**不在**
`workout-tracking/` 底下——它是根層產物。

- [ ] **Step 6: 跑 `status.py` 確認問題檔搬家沒搞丟東西**

Run: `python3 skills/bdd-clarify/scripts/status.py lab/go/skeleton`
Expected: 這一步**會壞**——`status.py` 現在找 `*/example-mapping.md`，而問題檔已經搬走。預期輸出是各 story 的「已答／待答」歸零或整列消失。

把實際輸出貼進 `gaps.md`。**這是預期中的失敗，Task 8 修它。** 不要在這裡順手改腳本。

- [ ] **Step 7: Commit**

```bash
git add lab/go/skeleton/docs/bdd/
git commit -m "refactor(lab): fold six story dirs into one workout-tracking slice with a spec.md"
```

---

## Task 4: `check_spec.py` 的覆蓋來源改讀 `spec.md`

**Files:**
- Modify: `skills/bdd-spec/scripts/check_spec.py:37-40`（`examples_in_map` → `stories_in_spec`）、`:92-155`（`check` 的主迴圈）

**Interfaces:**
- Consumes: Task 3 產出的 `lab/go/skeleton/docs/bdd/workout-tracking/spec.md`
- Produces: `check_spec.py <專案根>` 對新佈局回報雙向覆蓋，離開碼 0 = 沒問題

- [ ] **Step 1: 先跑一次，確認它現在是壞的**

Run: `python3 skills/bdd-spec/scripts/check_spec.py lab/go/skeleton`
Expected: 因為 `bdd.glob("*/example-mapping.md")` 找不到東西，六則 story 全部不被檢查；輸出只剩 `沒有對應 map 的 .feature` 那一行列出全部七個檔。**記下這個輸出**，它是修好之後的對照組。

- [ ] **Step 2: 用 `stories_in_spec` 取代 `examples_in_map`**

把 `skills/bdd-spec/scripts/check_spec.py` 的：

```python
def examples_in_map(text: str) -> set[str]:
    """map 裡的例子編號。"""
    return set(re.findall(r"^- Example (\d+\.\d+) \S", text, re.M))
```

換成：

```python
def stories_in_spec(text: str) -> dict[str, set[str]]:
    """spec.md 的 Example Mapping 段：story-slug -> 例子編號集合。

    一份 spec.md 涵蓋一批 story，每則一個 `### <story-slug>` 區塊。

    只掃 `## Example Mapping` 這一段。其他章節也會出現 `Example N.M`
    （Testing Decisions 的測試分層表就整欄都是），掃全檔會把測試分層表
    誤當成規格來源——而那張表本來就是從規格抄過去的，比對它等於自己跟自己比。
    """
    m = re.search(r"^## Example Mapping\s*\n(.*?)(?=\n## |\Z)", text, re.M | re.S)
    if not m:
        return {}
    out: dict[str, set[str]] = {}
    for blk in re.split(r"^### ", m.group(1), flags=re.M)[1:]:
        slug = blk.splitlines()[0].strip()
        out[slug] = set(re.findall(r"^- Example (\d+\.\d+) \S", blk, re.M))
    return out
```

- [ ] **Step 3: 換 `check()` 的主迴圈**

把 `for map_path in sorted(bdd.glob("*/example-mapping.md")):` 這一整層換成兩層——外層走 slice，內層走 story：

```python
    specs = sorted(bdd.glob("*/spec.md"))
    if not specs:
        print(f"{bdd} 底下沒有 spec.md —— 先跑 bdd-spec")
        return 1

    covered: set[str] = set()      # 有出現在某份 spec.md 裡的 story slug
    for spec_path in specs:
        stories = stories_in_spec(spec_path.read_text(encoding="utf-8"))
        if not stories:
            print(f"{spec_path.parent.name:30} —— spec.md 缺 `## Example Mapping` 段")
            problems += 1
            continue

        for slug, exs in sorted(stories.items()):
            covered.add(slug)
            fpath = feat_dir / f"{slug}.feature"
            if not fpath.exists():
                print(f"{slug:30} —— 還沒寫成 .feature")
                continue

            ftext = fpath.read_text(encoding="utf-8")
            # ↓ 以下沿用原本的 tags/missing/invented/states/zh_rule 那一整段，
            #   只把 `exs` 的來源換掉，其餘一字不動。
```

同時把最後兩處靠 `example-mapping.md` 判斷的地方換成 `covered`：

```python
    written = [p.read_text(encoding="utf-8") for p in feat_dir.glob("*.feature")
               if p.stem in covered]
```

```python
    orphans = [p.name for p in feat_dir.glob("*.feature") if p.stem not in covered]
```

- [ ] **Step 4: 對實測場實跑，比對輸出**

Run: `python3 skills/bdd-spec/scripts/check_spec.py lab/go/skeleton`
Expected:
- 六則 story 各一行，狀態 `✓`，例子數與 Task 3 搬進去的一致
- `沒有對應 map 的 .feature` 只剩 `['version.feature']`（骨架自帶，本來就不在 BDD 鏈上）
- 最後一行 `全部通過`，離開碼 0

Run: `echo $?`
Expected: `0`

- [ ] **Step 5: 反向驗證——確認「發明」抓得到**

臨時在某個 `.feature` 加一個不存在的 tag，確認稽核會叫：

```bash
sed -i 's/@example-1\.1/@example-1.1 @example-9.9/' lab/go/skeleton/features/log-a-workout.feature
python3 skills/bdd-spec/scripts/check_spec.py lab/go/skeleton; echo "exit=$?"
git checkout lab/go/skeleton/features/log-a-workout.feature
```

Expected: 輸出含 `指向 map 裡不存在的例子 ['9.9']`，`exit=1`。還原後再跑一次應回到 `exit=0`。

**這一步不能跳。** 這支腳本的價值全在「發明」那一條，而換來源最容易悄悄弄壞的就是它——迴圈改成兩層之後，`exs` 可能變成空集合，那時所有 tag 都會被算成發明，或反過來全都不檢查。

- [ ] **Step 6: 更新 `check_spec.py` 的 docstring 與 `bdd-spec/SKILL.md` 步驟 6 的說明**

docstring 裡提到 `example-mapping.md` 的地方改成 `spec.md 的 ## Example Mapping 段`。

- [ ] **Step 7: Commit**

```bash
git add skills/bdd-spec/scripts/check_spec.py skills/bdd-spec/SKILL.md
git commit -m "refactor(bdd-spec): read coverage from spec.md instead of example-mapping.md"
```

---

# Phase 2 — CLARIFY：依 Phase 1 的缺口設計技術追問

## Task 5: 技術追問面向清單

**Files:**
- Create: `skills/bdd-clarify/references/technical-probes.md`
- Modify: `skills/bdd-clarify/SKILL.md`（加指路）

**Interfaces:**
- Consumes: Task 3 Step 4 產出的 `gaps.md` 最後一欄
- Produces: 一組具名的技術追問面向，Task 6 的 Pass 3 逐項掃過

**這個任務的原料必須是 `gaps.md`，不能憑空列。** 業務的十一個面向（空與零、邊界、重複、時序、權限、失敗、時間、規模、降級、時限、可觀測）是從真實產出反覆修出來的；技術面向如果用想的，會列出一堆漂亮但問了沒用的題目。

- [ ] **Step 1: 從 `gaps.md` 的「這是哪一類技術問題」欄歸類**

把出現過的類別數一數。只保留**在實測場真的出現過**的類別，每一類寫成一個面向。

預期會長出來的（依 `spec.md` 骨架反推，實際以 `gaps.md` 為準）：

| 面向 | 問什麼 | 少了它會怎樣 |
| --- | --- | --- |
| **seam** | 驗收測試打在哪一層 | step definition 打錯高度，整批重寫 |
| **模組邊界** | 這批行為要收在幾個模組裡、各自負責什麼 | 型別有了但不知道住哪，實作時各憑喜好 |
| **介面** | 那些模組對外的函式名、參數、回傳 | 兩張票同時實作同一個服務，簽名對不起來 |
| **互動** | 跨模組的呼叫順序、誰呼叫誰 | 交易邊界與錯誤傳遞沒有依據 |
| **既有資產** | 專案已經有什麼可以用（既有測試、既有型別、既有端點） | 重造一份，然後兩份行為不一樣 |

- [ ] **Step 2: 寫 `skills/bdd-clarify/references/technical-probes.md`**

每個面向要寫齊四樣：**問什麼**、**問誰**、**什麼時候不用問（n/a 的判準）**、**一個實測場的真實例子**。

只寫面向名稱的話，得到的回答會太模糊而要再問一輪——這是 `structure.md` S14 記錄過的失敗方式。

- [ ] **Step 3: 在 `bdd-clarify/SKILL.md` 的 `## 產物` 段的格式對照表加一列指路**

```markdown
Pass 3 的技術追問面向 → `references/technical-probes.md`
```

- [ ] **Step 4: 機械稽核**

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify`
Expected: 離開碼 0，不含 S5。

- [ ] **Step 5: Commit**

```bash
git add skills/bdd-clarify/references/technical-probes.md skills/bdd-clarify/SKILL.md
git commit -m "feat(bdd-clarify): add the technical probe dimensions derived from the lab run"
```

---

## Task 6: `bdd-clarify` 加「分批」與 Pass 3

**Files:**
- Modify: `skills/bdd-clarify/SKILL.md`（`## 流程` 概觀、Pass 1 的步驟 4 之後加分批、Pass 2 之後加 Pass 3）

**Interfaces:**
- Consumes: Task 5 的 `references/technical-probes.md`
- Produces: `docs/bdd/brief.md` 多一節「批次清單與順序」；`docs/bdd/<slice-slug>/questions/` 底下多出技術類問題檔

- [ ] **Step 1: Pass 1 加第 5 步「分批」**

在 `### 4. 切 story —— 依據來自上面三步` 之後加：

```markdown
### 5. 分批 —— 一批 ＝ 一次可交付的價值

story 切完之後還要再分一次。判準只有一條：
**這一批做完，能不能對使用者交付一次完整的價值？**

一批一份 `spec.md`，所以批次大小直接決定那份檔案會不會大到沒人讀得完。
實測參考點：六則 story、69 個場景 → 一份 36.5K 的文件，已經在邊緣。

MUST: 批次清單與順序寫進 `brief.md`，並寫**為什麼是這個順序**。
沒有它，第二批開工時沒有人知道第一批當初為什麼不做那幾則。

MUST NOT: 對這一批以外的 story 做深度追問。Pass 2 只問這一批的——
問了三十題才發現其中兩則這一季根本不做，是最常見的浪費。
```

- [ ] **Step 2: 在 Pass 2 之後新增 Pass 3**

```markdown
## Pass 3 · 技術 —— 對開發者

Pass 2 收斂之後才跑。順序不能反：**技術問題有一半要看過規則才問得出來**
（「這批行為要收在幾個服務裡」得先知道有哪些行為）。

問法跟 Pass 2 一樣——一次一題、四個選項、答不出來記成待答。
差別只有兩個：**對象是開發者不是產品方**，**面向清單換一份**。

逐項掃過 `references/technical-probes.md` 的每個面向。
標 `n/a` 的要寫理由，否則它跟「懶得問」分不出來。

IMPORTANT: **SPEC 那一步不會再問任何問題。** 它只綜合已經有答案的東西。
所以這一趟漏掉的面向，不會在下游被補起來——它會變成 `spec.md` 裡一個
沒有來源的句子，或一格空白。
```

- [ ] **Step 3: 更新 `## 流程` 概觀的兩 pass 敘述成三 pass**

含「續跑：什麼時候跳過 Pass 1」那一段——續跑時 Pass 3 的處置要寫明（規則沒變就不用重問）。

- [ ] **Step 4: 行數檢查**

Run: `wc -l skills/bdd-clarify/SKILL.md`
Expected: ≤ 500。現在是約 600 行的量級，**這一步很可能超標**——超了就把 Pass 1 的「拆成 brief」與「識別角色」細節下放到既有的 `references/prd-breakdown.md` 與 `references/actor-definition.md`，本文留骨架與指路。

- [ ] **Step 5: 機械稽核**

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify`
Expected: 離開碼 0。

- [ ] **Step 6: Commit**

```bash
git add skills/bdd-clarify/SKILL.md
git commit -m "feat(bdd-clarify): add slice batching to pass 1 and a technical pass 3"
```

---

## Task 7: 拿掉 `example-mapping.md` 這個產物

**Files:**
- Modify: `skills/bdd-clarify/SKILL.md`（步驟 5「澄清循環 —— 一次一題，每輪抽規則」、步驟 6 就緒判定、步驟 7 收尾、`## 產物` 段）
- Modify: `skills/bdd-clarify/references/map-format.md`（拿掉 map 的章節，狀態列加兩個欄位）
- Modify: `skills/bdd-spec/SKILL.md`（步驟 2「讀三樣東西」的來源換掉）
- Modify: `lab/go/skeleton/docs/bdd/workout-tracking/questions/*.md`（狀態列補欄位）

**Interfaces:**
- Consumes: Task 6 的三 pass 結構
- Produces: 問題檔的狀態列多兩個欄位 `story` 與 `面向`，Task 8 的 `status.py` 靠它們算覆蓋

- [ ] **Step 1: `map-format.md` 的狀態列格式加兩欄**

```markdown
## 問題檔的狀態列

每個問題檔的第三行是狀態列。**紅卡靠它辨識，不是靠所在目錄**：

```markdown
**狀態**: 待答（之後再問） · **輪次**: 2 · **story**: log-a-workout · **面向**: 邊界
**狀態**: 已答 · **輪次**: 2 · **信心**: 中 · **story**: log-a-workout · **面向**: 時序
```

| 欄位 | 值 | 為什麼要有 |
| --- | --- | --- |
| `story` | story slug，或 `全部`（這一批共通） | 一個 slice 的問題全放同一個目錄，靠這欄歸屬 |
| `面向` | 十一個業務面向或技術面向之一 | **追問覆蓋改成從這裡算**，不再另存一段 |

`n/a`（這個面向不適用）也開一個問題檔，狀態寫 `n/a`，理由寫在本文。
否則它跟「懶得問」在檔案上長得一樣。
```

- [ ] **Step 2: 拿掉 `map-format.md` 裡關於 `example-mapping.md` 的三個章節**

刪 `## 循環中的 example-mapping.md`、`## 抽規則之後`、`## `<feature-slug>` 怎麼取` 三節，改寫一段說明去向：

```markdown
## 規則與例子去哪了

規則與例子**不再是 CLARIFY 的產物**。它們的定版編號誕生在 `spec.md` 的
`## Example Mapping` 段——也就是真正要用 tag 引用它們的那一刻。

原本凍在 CLARIFY 的編號，是在最不確定的時候做最不可逆的事：那時 story 邊界
還可能合併或再切，而 `MUST NOT 重排編號` 一旦生效就改不動了。

CLARIFY 仍然會在追問過程中把規則講出來——那是問下一題的前提。
只是它們留在問答的紀錄裡，不另存一份檔案。
```

- [ ] **Step 3: `bdd-clarify/SKILL.md` 的步驟 5、6、7 與 `## 產物` 段對齊**

- 步驟 5：拿掉「每輪抽規則寫進 map」的產出動作，保留「抽規則」作為判斷還有什麼沒問的**工作手法**
- 步驟 6 就緒判定：判準從「藍卡／紅卡數」換成 `status.py` 的兩個數字

```markdown
| 訊號 | 看什麼 | 代表 |
| --- | --- | --- |
| story 太大 | **已答問題數**偏高 | 答案多 ≈ 規則多 → `story-splitting` |
| 不確定性高 | **未答問題數**偏高 | 紅卡沒清 → 再跑一輪 `clarify-loop` |
```

- `## 產物` 段的目錄樹換成 `PLAN.md` 的新佈局

- [ ] **Step 4: `bdd-spec/SKILL.md` 步驟 2「讀三樣東西」換來源**

把 `example-mapping.md | 規則、例子、編號——這是規格的唯一來源` 這一列換成：

```markdown
| `<slice>/questions/*.md` | 已答的問題與答案——**規則與例子從這裡抽**，這是規格的唯一來源 |
```

並在表格底下加一句：

```markdown
IMPORTANT: 只讀狀態列標 `已答` 的。待答的問題代表那一塊還沒有結論，
把它寫成場景等於**把不確定性從一個顯眼的問題檔搬進一份看起來已完成的規格**。
```

- [ ] **Step 5: `glossary.md` 瘦成只有詞彙**

`glossary.md` 現在同時裝詞彙與 `Shared N` 跨 story 規則，而那些規則正是
CLARIFY 抽規則時長出來的——那一段拿掉之後，規則沒有理由留在詞彙表裡。

改兩個地方：

1. `skills/bdd-clarify/SKILL.md` 的 `## 產物` 段，`glossary.md` 那一行改成
   `跨批次的詞彙（ubiquitous language）。不放規則。`
2. `lab/go/skeleton/docs/bdd/glossary.md`：`Shared N` 的條目依下表搬家

| 那條 Shared 規則是 | 搬去 |
| --- | --- |
| 某次操作要滿足的**行為** | 對應 `.feature` 的 `Rule:`（若已經有等價的 Rule 就刪掉，不要留兩份） |
| 任何時刻都必須為真的**狀態** | `docs/bdd/domain-model.md` 的不變條件表 |

MUST: 每一條都要找得到新家。找不到的代表它其實是一個還沒被寫成場景的規則——
寫進 `spec.md` 的 `## Out of Scope` 的「回 CLARIFY 補問」，**不要默默刪掉**。

- [ ] **Step 6: 實測場的問題檔補欄位**

給 `lab/go/skeleton/docs/bdd/workout-tracking/questions/` 底下每個檔的狀態列補上 `story` 與 `面向`。`story` 的值從 Task 3 Step 1 搬家前的原目錄名還原（`git log --diff-filter=R --name-status` 查得到）。

- [ ] **Step 7: Commit**

```bash
git add skills/bdd-clarify skills/bdd-spec/SKILL.md lab/go/skeleton/docs/bdd/
git commit -m "refactor(bdd-clarify): drop example-mapping.md, file questions by story and dimension"
```

---

## Task 8: `status.py` 的覆蓋改從問題檔算

**Files:**
- Modify: `skills/bdd-clarify/scripts/status.py:34-56`（`coverage` 與 `count_questions`）、`:79-125`（`main`）

**Interfaces:**
- Consumes: Task 7 加好欄位的問題檔
- Produces: `status.py <專案根>` 對新佈局印出每則 story 的已答／待答／追問覆蓋

- [ ] **Step 1: 先跑一次，確認它現在是壞的**

Run: `python3 skills/bdd-clarify/scripts/status.py lab/go/skeleton`
Expected: `docs/bdd 底下沒有 example-mapping.md`，離開碼 1。這是 Task 3 Step 5 記下的同一個失敗。

- [ ] **Step 2: 加一個解析狀態列的函式**

```python
def parse_status(text: str) -> dict[str, str]:
    """問題檔第三行的狀態列 -> 欄位字典。

    格式是 `**狀態**: 已答 · **輪次**: 2 · **story**: log-a-workout · **面向**: 邊界`。
    分隔符是全形間隔號，欄名用 `**...**` 包住——兩者都固定，所以一條 regex 掃完
    比逐欄寫 pattern 好維護。
    """
    return {k: v.strip() for k, v in
            re.findall(r"\*\*(.+?)\*\*:\s*([^·\n]+)", text)}
```

- [ ] **Step 3: `count_questions` 改成回傳分組結果**

```python
def scan_questions(qdir: Path) -> dict[str, dict]:
    """走一遍問題檔，依 story 分組。

    回傳 {story: {"answered": int, "pending": [檔名], "dims": {面向}}}。
    `story` 欄缺席時歸到 `全部`——判斷不了時放共通層，代價不對稱：
    共通的問題每則 story 都看得到，掛錯 story 的別則看不到。
    """
    out: dict[str, dict] = {}
    for c in sorted(qdir.glob("*.md")):
        st = parse_status(c.read_text(encoding="utf-8"))
        story = st.get("story", "全部")
        g = out.setdefault(story, {"answered": 0, "pending": [], "dims": set()})
        if dim := st.get("面向"):
            g["dims"].add(dim)
        if st.get("狀態", "").startswith("已答"):
            g["answered"] += 1
        else:
            g["pending"].append(c.stem)
    return out
```

- [ ] **Step 4: `coverage` 改成吃面向集合**

```python
def coverage(dims: set[str]) -> str:
    """把「有沒有人問過這個面向」壓成一行符號。

    ✓ 這個面向至少有一個問題檔 · — 一題都沒有

    原本讀的是 map 裡手寫的一段覆蓋表。改成算的之後，「寫了覆蓋表但沒真的問」
    這種狀態不再可能存在——那正是手寫彙總遲早會跟來源不一致的地方。
    """
    return "".join("✓" if d in dims else "—" for d in DIMS)
```

`n/a` 的表示：狀態列寫 `**狀態**: n/a` 的問題檔仍然帶著 `面向` 欄，所以它會算成 `✓`（問過了，結論是不適用）。這是刻意的——`n/a` 是一個有理由的結論，不是一個沒問到的洞。把這句話寫進函式的 docstring。

- [ ] **Step 5: `DIMS` 加上技術面向**

`DIMS` 現在只有十一個業務面向。Pass 3 的技術面向（Task 5 的
`references/technical-probes.md`）也要出現在覆蓋列，否則「seam 沒問到」
這件事在儀表板上完全看不見——而那正是這次改版要防的洞。

```python
BIZ_DIMS = ["空與零", "邊界", "重複", "時序", "權限", "失敗", "時間", "規模",
            "降級", "時限", "可觀測"]
# Pass 3 的技術面向。清單的來源是 references/technical-probes.md——
# 兩邊不一致時以那份為準，這裡只是為了印出固定寬度的欄位。
TECH_DIMS = ["seam", "模組邊界", "介面", "互動", "既有資產"]
DIMS = BIZ_DIMS + TECH_DIMS
```

表頭那兩行的欄位縮寫也要跟著加，並在最底下的圖例裡把兩組分開標示。

MUST: `TECH_DIMS` 的實際內容以 Task 5 產出的檔案為準，不要照抄上面這五個——
上面那組是預期值，Task 5 跑完才知道真的長什麼樣。

- [ ] **Step 6: `main` 改走新佈局**

外層走 `bdd.glob("*/questions")` 找 slice，內層走 `scan_questions` 的分組印每則 story 一行。`規則`／`例子` 兩欄拿掉（來源沒了），`就緒判定` 那欄改印判斷依據——已答／待答兩個數字已經在表上，這欄改成空白或拿掉，不要留一個永遠印 `？未寫` 的欄位。

同步更新檔頭 docstring：現在的第二段整段在講 `example-mapping.md`，要改。

- [ ] **Step 7: 對實測場實跑**

Run: `python3 skills/bdd-clarify/scripts/status.py lab/go/skeleton`
Expected:
- 六則 story 各一行 ＋ 一列 `全部`（跨 story 的問題）
- 已答＋待答的總數 ＝ Task 3 Step 1 數到的問題檔數
- 追問覆蓋欄有 `✓` 也有 `—`（全部 `—` 代表 Task 7 Step 6 的欄位沒補到，全部 `✓` 代表 regex 太寬鬆把別的東西也算進來了）
- 技術面向那幾格**預期大多是 `—`**：實測場的問題檔是舊制跑出來的，還沒有 Pass 3。這正是這次改版要看見的洞

Run: `echo $?`
Expected: `0`

- [ ] **Step 8: Commit**

```bash
git add skills/bdd-clarify/scripts/status.py
git commit -m "refactor(bdd-clarify): compute probe coverage from question files"
```

---

# Phase 3 — PLAN：重寫成 tracer bullet

## Task 9: `plan.md` 的骨架重寫

**Files:**
- Rewrite: `skills/bdd-plan/references/plan-format.md`

**Interfaces:**
- Consumes: `spec.md` 的 Implementation Decisions、Testing Decisions、Risks
- Produces: `plan.md` 的票格式，Task 11 照著產

- [ ] **Step 1: 整份重寫 `skills/bdd-plan/references/plan-format.md`**

原本的七節（API／Domain／Schema／測試分層／風險／順序／缺口）全部作廢——前五節搬去 `spec.md` 了。新骨架：

````markdown
# `plan.md` 的骨架

照抄即可。讀這份的時機：步驟 3 之後，準備寫檔的時候。

一批一份，路徑 `docs/bdd/<slice-slug>/plan.md`。

```markdown
# 實作計畫

**Slice**: <slice-slug>
**來源**: `spec.md`（<日期>）＋ `features/` 底下標 `@ready` 的 <N> 個檔
**日期**: <YYYY-MM-DD>
**seam**: <照抄 spec.md 的 Testing Decisions>

## 前置整理（prefactoring）

「make the change easy, then make the easy change.」先做的改動，
排在所有功能票前面。沒有就寫「沒有——現有結構直接容得下這批行為」。

### 00 — <標題>
**改什麼**: <一句話>
**Blocked by**: 無，可立即開始
- [ ] <驗收條件>

## 票

依賴序排列，blockers 在前。編號一旦寫出去就不重排——
`bdd-implement 03` 靠它定位。

### 01 — <票的標題>

**交付什麼行為**: 從使用者角度寫這張票做完之後**什麼事會動**。
不是分層清單。「新增 WorkoutRepository」不是行為，
「使用者可以記錄一次只有一個動作的訓練並讀回來」是。

**涵蓋場景**: `@example-1.1`、`@example-1.2`
**Blocked by**: 無，可立即開始
**內迴圈**: <這張票底下要寫哪些 unit test，來自 spec.md 的測試分層>

- [ ] <驗收條件 1>
- [ ] <驗收條件 2>

### 02 — <票的標題>

**Blocked by**: 01
⋯
```

---

## 三條切票規則

| 規則 | 破了會怎樣 |
| --- | --- |
| **垂直切片**：一張票穿過所有層（schema／API／UI／測試） | 依層切的實測案例：26 張票、每張平均跑 20 次 agent，四分之三是返工 |
| **自己可 demo**：做完當下就能展示 | 驗收條件會伸手進別張票擁有的工作，變成不可能獨立驗 |
| **塞得進一個全新 session**：撿起這張票的東西沒看過 `spec.md` | 票太大時實作者會在中途失去脈絡，然後開始猜 |

檢驗每張票只要問一句：**做完之後我能 demo 什麼？** 答不出來的就是水平切片。

## 驗收條件要能失敗

每一條驗收條件，說得出「什麼觀察會顯示它是假的」，並確認**在實作者的起點
commit 上它是紅的**。

三種常見的失效寫法：

| 寫法 | 為什麼是壞的 |
| --- | --- |
| 起點就已經成立的條件 | 它不驗任何東西 |
| 只有別張票做完才可能成立 | 這張票不擁有它，等於把驗收外包 |
| 把需求原文抄一遍 | 沒有從產物推導，只是換句話說 |

垂直切片本身就擋掉大部分：**交付了原本不存在的行為的切片，在起點必然是紅的。**

## wide refactor 的例外

一個機械性改動（改欄位名、換共用型別）影響面擴散到全 codebase 時，
沒有任何垂直切片能綠。這種改 expand–contract：

1. **expand** — 加新形式在舊的旁邊，什麼都不破
2. **migrate** — 呼叫端分批遷移，每批一張票、依 blast radius 分（每個 package、
   每個目錄），每張 blocked by expand。CI 保持綠，因為舊形式還在
3. **contract** — 沒有呼叫者之後刪掉舊形式，blocked by 全部 migrate 票

連分批都無法各自保持綠時，讓它們共用一條整合分支，全部 block 一張最後的
整合驗證票——綠只在那裡承諾。

## 為什麼不是「每張票一個檔」

上游（`to-tickets`）把票寫成一檔一張，因為它要餵一個會去 tracker 撿票的
agent 車隊。本專案的票留在同一份 `plan.md` 裡，理由是**順序是排序關係，
拆進 N 個檔案之後就不存在了**——而順序正是這一步最主要的產出。

真的要平行跑票時再拆。拆一份現成的清單很便宜，把散落的檔案重新排序很貴。
````

- [ ] **Step 2: Commit**

```bash
git add skills/bdd-plan/references/plan-format.md
git commit -m "refactor(bdd-plan): rewrite plan-format.md around tracer-bullet tickets"
```

---

## Task 10: `bdd-plan/SKILL.md` 重寫

**Files:**
- Rewrite: `skills/bdd-plan/SKILL.md`

**Interfaces:**
- Consumes: Task 9 的 `references/plan-format.md`
- Produces: `bdd-plan` 執行後在 `docs/bdd/<slice-slug>/plan.md` 留下票清單

- [ ] **Step 1: 改 frontmatter 的 `description`**

原本的觸發詞全是技術設計（「要建哪些 API」「domain model 長怎樣」「資料表要怎麼設計」）——那些現在屬於 `bdd-spec`。留著會讓兩個 skill 搶同一批處境，而模型會挑一個、你不會知道它挑了哪個。

```yaml
description: >
  SPEC 說要做成什麼，PLAN 說先做哪一塊——把一份 `spec.md` 與它的 `.feature`
  切成一疊 tracer bullet：每張票穿過所有層、自己可 demo、塞得進一個全新 session，
  並宣告它被哪幾張票 block。BDD 六步流程的 PLAN 步驟。
  觸發詞：「這批要怎麼切」「拆成幾張票」「先做哪一個」「排實作順序」
  「切成小步驟」「spec 寫完了下一步」「這些票誰卡誰」「開工順序」。
  English: break this spec into tickets, slice the work, what order should I
  build these, vertical slices, what can I start first.
```

- [ ] **Step 2: 改 `## 使用時機` 與 `## Skill Boundaries`**

```markdown
## 使用時機

- `spec.md` 與 `.feature` 都寫好了，要決定切成幾塊、先做哪一塊
- 一批工作大到一個 session 做不完，需要拆給多個 session 接力

## Skill Boundaries

- 還沒有 `spec.md` → 先跑 `bdd-spec`；規則還沒定案 → 回 `bdd-clarify`
- **要決定 API、domain 型別、schema、seam → 那是 `bdd-spec` 的工作**，不是這裡
- 整批塞得進一個 context window → 不需要這一步，直接 `bdd-implement`
- 要實際寫 step definition 與產品程式碼 → `bdd-implement`
```

倒數第二條是新的，值得留：**票有下界**。整批做得完的時候切票只是多付一次綜合成本，還多一個模型會漂移的環節。

- [ ] **Step 3: 用「只切，不設計」取代「只翻譯，不設計」整節**

```markdown
## 只切，不設計

MUST NOT: 在這一步決定任何 API、型別、欄位、資料表或 seam。
**它們已經在 `spec.md` 裡了**；這裡再決定一次，只會產生一份跟規格不一樣的規格。

判準：指著票裡的任何一句話，說得出它來自 `spec.md` 的哪一節、
或哪個 `@example-N.M` 嗎？說不出來就是這一步僭越了。

`spec.md` 沒寫到而切票時發現需要的，寫進 `## 缺口` 回報，不要順手補。
```

- [ ] **Step 4: 流程改成三步**

```markdown
## 流程

### 1. 讀

`docs/bdd/<slice-slug>/spec.md` 全部，加上它涵蓋的每個 `.feature`。

**一次讀完整批。** 票之間的依賴只有並排時才看得見；分批切會各自長出一套順序，
然後在實作時發現它們互相卡住。

順便看一眼專案既有的結構——找 prefactoring 的機會。

### 2. 切

先切出**最薄的完整路徑**當第一張票：碰到每一層、邏輯最少。它證明的是接線，
不是行為。專案已經有骨架就跳過。

其餘依 `spec.md` 的 Example Mapping 切，一張票一組相關的 `@example-N.M`。
每張票問一次「做完之後我能 demo 什麼」，答不出來就是切錯了。

三條規則與 wide refactor 的例外 → `references/plan-format.md`。

### 3. 攤出來確認（Ask user）

MUST: 寫檔之前先把票的清單攤給使用者看——編號、標題、Blocked by、
交付什麼行為——並問三件事：

1. 粒度對不對（太粗／太細）
2. Blocked by 是不是真的卡住，還是只是「先做比較順」
3. 有沒有該合併或該再拆的

IMPORTANT: **切太細是這一步最常見的失效**，而且它不會自己被發現——
一疊原子化的小票看起來很專業，實際上把「這幾件事其實是同一件」的資訊丟掉了。
這個確認點就是為它存在的。

使用者確認之後才寫檔。
```

- [ ] **Step 5: 改 `## 產物` 與 `## 完成後`**

```markdown
## 產物

一批一份 `docs/bdd/<slice-slug>/plan.md`。骨架照抄 `references/plan-format.md`。

MUST: 只寫入 `docs/bdd/`，**不修改專案既有的任何檔案**。

## 完成後

告訴對方：幾張票、第一張是哪一張（以及**為什麼是它**）、有沒有 prefactoring、
`spec.md` 裡推不出來的缺口有幾條。下一步 `bdd-implement <票號>`，一張票一個
全新 session。
```

- [ ] **Step 6: 行數與機械稽核**

Run: `wc -l skills/bdd-plan/SKILL.md && python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-plan`
Expected: 行數 ≤ 500（重寫後應該比原本的 10.9K 短很多）；離開碼 0。

- [ ] **Step 7: Commit**

```bash
git add skills/bdd-plan/SKILL.md
git commit -m "refactor(bdd-plan): rewrite as a tracer-bullet slicer, drop the design sections"
```

---

## Task 11: 實測場整鏈重跑

**Files:**
- Create: `lab/go/skeleton/docs/bdd/workout-tracking/plan.md`
- Modify: `PLAN.md`（`## 進度` 打勾、`## 未決事項` #4 與 #5 填上實測數字）

**Interfaces:**
- Consumes: Task 3 的 `spec.md`、Task 9/10 的新 `bdd-plan`
- Produces: 三步產物齊全的一個 slice，證明新分工跑得通

- [ ] **Step 1: 依新的 `bdd-plan` 產出 `plan.md`**

輸入是 `lab/go/skeleton/docs/bdd/workout-tracking/spec.md` ＋ 六個 `.feature`。
攤出票清單，自己扮演確認者走一次步驟 3，記下每張票「做完能 demo 什麼」的答案。

- [ ] **Step 2: 逐票檢查驗收條件能不能失敗**

對每張票的每條驗收條件，在**起點 commit**（`git stash` 掉 plan.md 之外的所有改動，或直接看 `HEAD`）上確認它是紅的。做不到的條件就地改掉。

Expected: 每張票至少有一條在起點必然為假的條件。

- [ ] **Step 3: 三支腳本全綠**

```bash
python3 skills/bdd-clarify/scripts/status.py lab/go/skeleton; echo "status=$?"
python3 skills/bdd-spec/scripts/check_spec.py lab/go/skeleton; echo "check=$?"
for d in skills/*/; do python3 skills/skill-rules/scripts/audit_skill.py "$d" || echo "FAIL $d"; done
```

Expected: `status=0`、`check=0`、audit 全部無輸出（沒有 FAIL 行）。

- [ ] **Step 4: `claude plugin validate`**

Run: `claude plugin validate . --strict`
Expected: 通過。這是唯一會抓到 frontmatter 與 manifest 層級錯誤的檢查。

- [ ] **Step 5: 把實測數字回填 `PLAN.md`**

- `## 進度`：`bdd-clarify 改版`、`bdd-spec 擴張`、`bdd-plan 重寫`、`實測場整批重跑` 四項打勾
- `## 未決事項` #4「一批到底多大」：填上這一批的實際數字（幾則 story、幾個場景、`spec.md` 幾 K、切出幾張票）
- `## 未決事項` #5「技術問題的追問面向」：改成已定，指向 `references/technical-probes.md`

- [ ] **Step 6: 記錄執行中發現的摩擦**

`bdd-clarify`、`bdd-spec`、`bdd-plan` 三次改版都是**在紙上做的**。這一趟必然會撞到紙上想不到的東西。把它們寫進 `PLAN.md` 的 `## 未決事項`，不要當場改 skill——當場改會讓「改了什麼」與「為什麼改」一起消失在對話裡。

- [ ] **Step 7: Commit**

```bash
git add lab/go/skeleton/docs/bdd/ PLAN.md
git commit -m "test(lab): run the re-cut three-step chain end to end on workout-tracking"
```

---

## 已知風險

寫在這裡而不是散在任務裡，因為它們**跨任務**：

| 風險 | 從哪裡長出來 | 撞到的時候 |
| --- | --- | --- |
| `bdd-clarify/SKILL.md` 已經接近 500 行，Task 6 還要再加 Pass 3 | S10 是訊號不是硬上限，但這份確實已經該拆 | Task 6 Step 4 檢查行數，超了就把 Pass 1 的細節下放 |
| 技術追問面向清單只有一個實測案例（Go／分層／REST） | Task 5 的原料只有一份 `gaps.md` | 先寫下來，第二個實測場（不同技術棧）再修。不要為了看起來完整而補想像出來的面向 |
| `check_spec.py` 與 `status.py` 沒有單元測試 | 這個 repo 沒有測試基礎設施 | Task 4 Step 5 的反向驗證是唯一擋住「改壞了但看起來正常」的機制，不能跳 |
| 六則 story 併成一個 slice 之後 `spec.md` 可能過大 | Task 3 是把 36.5K 的 `plan.md` 內容搬進去 | 真的太大就在 Task 3 之後把 slice 拆成兩批，並把數字寫進 `PLAN.md` 未決事項 #4——那正是那一項存在的理由 |
