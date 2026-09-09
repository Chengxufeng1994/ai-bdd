# CLARIFY 的產物改成一份 PRD，切 story 搬到 SPEC

**日期**：2026-09-09
**狀態**：設計已確認，待實作
**動的東西**：`bdd-clarify`、`bdd-spec`、`story-splitting` 的職責，兩支腳本與 `clarify-loop` 的資料來源，`PLAN.md` 未決事項 #1

---

## 要解的問題

CLARIFY 現在在 Pass 1 就切 story 並分批，而**切分的依據那時還沒穩定**。

實測證據（2026-09-09，`parking-billing` case，agent 對 agent）：CLARIFY 問了 15 題
範圍題，需求方的答案在過程中三次改變形狀——

| 改口 | 連帶效果 |
| --- | --- |
| 「入場不擋」→「入場擋」 | 抽票與柵欄的矛盾才浮現 |
| 出口不擋救護車 | 入口變成救護車進不來 |
| 入場改成擋 | 「滿場要不要拒絕進入」從不存在變成必答 |

若在第一批之後就切 story，這些切線全部是錯的。

### 這不是流程問題，是結構問題

現行主文件 `clarify.md` 住在 `specs/<slice-slug>/` 底下。**要有 slug 才有目錄，
所以切分必須先發生。** 目錄結構逼出流程順序，流程順序又逼出職責分配——三者是
同一個決定的三個面向，只改流程不改目錄，改不動。

同一個結構還造成第二個症狀：`prd.md` 的前身把 users 放 `actor.md`、user stories 放
`brief.md`、決策放 `questions/`、規則放 `clarify.md`，**散在四個檔兩個層級**。
要拿去跟 PM 對焦，得同時開四個檔。

---

## 已確認的決策

### D1 — 切分有兩種，依據不同，所以住不同地方

| | 靠什麼判斷 | 住哪 |
| --- | --- | --- |
| **切 story**（哪些行為湊成一則可交付的 story） | 領域知識 | **SPEC** |
| **切票**（tracer bullet、依賴、順序） | 依賴關係與規模，**不需要**領域知識 | PLAN（不變） |

依據是 `ai-sdlc.md` 已經寫下的話：

> PLAN 要獨立的理由不同：它是唯一一個**不需要領域知識**的步驟。切票靠的是依賴
> 關係與規模判斷，不是懂不懂停車場怎麼收費。

把切 story 放進 PLAN 會違反這句；放進 SPEC 則是**把它已經在做的事講明白**——現行
契約本來就寫著 `features/<story-slug>.feature　SPEC　一則 story 一個檔`，SPEC 必須
決定哪些場景進哪個檔，那就是切 story。

### D2 — CLARIFY 一個 feature 一份 `prd.md`，不切分

三個 Pass 都對同一份文件寫：`Pass 1 廣度 → Pass 2 深度（對整個 feature）→ Pass 3 技術`。
拿掉現行 Pass 1 的第 4 步（切 story）與第 5 步（分批）。

**Pass 3 的產出落在哪**：技術追問問出來的東西分兩類——量得出來的（延遲、併發、
資料量、可用性）進 `## Non-Functional Requirements`；不可協商的外部限制（既有硬體、
法規、已發包的規格）進 `## Constraints`。答不出來的一律進 `## Open Questions`，
跟業務面的紅卡同一張表。Pass 3 不新增 `FR-N`。

### D3 — AC 與 Examples 分兩層，編號同源

- **AC-N.M** 給 PM 讀：使用者看得到什麼
- **EX-N.M** 給 SPEC 讀：具體到含數字

兩層共用 `FR-N` 的編號，所以可追溯鏈不斷：
`FR-3 → EX-3.2 → @example-3.2 → Scenario → Ticket → Step`。

漂移風險（改了 AC 沒改 EX）用機械檢查擋，見〈腳本影響〉。

### D4 — 外層全部消失，`prd.md` 自足

刪掉根層的 `brief.md`、`actor.md`、`glossary.md`、`questions/`。

它們存在的理由是「跨批次累積」，而**批次是早切的產物**。切分搬走之後 CLARIFY
眼裡只有一個 feature，「跨批次」這個概念本身就沒有了。

### D5 — 目錄名 `specs/<date>-<feature-name>/`

日期優於序號：序號要協調（誰佔了 001），日期不用，且天然照時間排。

### D6 — 問題與決策史收進 `prd.md`，表當索引、小節裝追問串

一個目錄一個檔。格式見〈產物契約〉。

### D7 — 詞彙表由 SPEC 寫進 `docs/CONTEXT.md`

**誰寫**：CLARIFY 發現詞義有歧義 → 變成一題 → 答案進 `prd.md`；SPEC 把**已經有
答案**的詞抬進 `CONTEXT.md`，不新增定義。這符合「SPEC 不訪談」——SPEC 做的是
轉錄，不是裁決。

理由是 SPEC 本來就負責 `domain-model.md` 與 step grammar，兩者都是詞彙工作；詞彙
在 SPEC 才變成承重的東西，step definition 靠一致措辭重用。

**放哪**：`docs/CONTEXT.md`，成為第二個「不進 `specs/`」的例外，見〈約束的例外〉。

---

## 目標結構

```
specs/
└── 2026-09-09-parking-billing/
    ├── prd.md      CLARIFY  完整自足，拿去跟 PM 對焦
    ├── spec.md     SPEC     切出哪幾則 story、依據、Gherkin 表達不了的決定
    └── plan.md     PLAN     tracer bullet ＋ blocking edges

features/
├── visitor-billing.feature       SPEC  一則 story 一個檔
└── monthly-pass-entry.feature

docs/CONTEXT.md                   SPEC  跨 feature 的 ubiquitous language
```

---

## 產物契約：`prd.md`

```markdown
## Problem / Goal / Success Metrics
## Actors                            這個 feature 用到的角色，不外流到 CONTEXT.md
## Scope — In / Out                  Out of Scope 是明確邊界，不是免責條款
## Functional Requirements
    ### FR-N  <EARS 句式>
        #### AC     AC-N.M  給 PM：使用者看得到什麼
        #### Examples  EX-N.M  給 SPEC：具體到含數字
## Non-Functional Requirements
## Open Questions
## Assumptions / Constraints
```

相對現行 `clarify.md` 的變動：

| 變動 | 理由 |
| --- | --- |
| `Business Rules → ### <story-slug>` 那層消失，規則平鋪成 `FR-N` | story 分組是早切的痕跡 |
| Rules 改寫成 EARS 句式 | 比自由格式更有結構，`When <trigger>, the system shall <response>` 天然對應 `Rule:` 區塊。六個 pattern 見 [`docs/sdd.md`](../../sdd.md) 的 `## EARS`，本文與 skill 都不重述 |
| 新增 `## Non-Functional Requirements` | 現行完全沒有；CLARIFY 是唯一還在跟需求方對話的步驟，這裡不問後面沒人問 |
| 新增 AC 層 | PM 對焦用 |
| `## Domain Glossary` 移除 | 交給 SPEC 寫 `docs/CONTEXT.md`（D7） |

**兩個容易混淆的界線，在這裡講死：**

- **角色留在 `prd.md`，不進 `CONTEXT.md`。** 角色是這個 feature 的流程裡誰在做事，
  換一個 feature 就換一批；詞彙不是。角色若日後也需要跨 feature 累積，套同一條
  rule of three，不預先建。
- **詞義的決策住 `## Open Questions`，不另設一節。** CLARIFY 發現「訪客」有歧義時，
  它就是一題；答案進表格與小節。SPEC 之後從那裡抬進 `CONTEXT.md`。所以 `prd.md`
  沒有 glossary 一節，但**詞義的來源與時間點都查得到**。

### `## Open Questions` 的格式

表是機械可解析的索引，小節是人讀的決策史：

```markdown
## Open Questions

| Q | 問題 | 狀態 | 答案 |
| --- | --- | --- | --- |
| Q1 | 月租費帳務在不在範圍 | 已答 | 名單與期限在，錢不在 |
| Q2 | 消防法規對閘門的要求 | **待答** | — |

### Q2 消防法規對閘門的要求
**狀態**：待答　**該問誰**：消防設備師、柵欄廠商
2026-09-09 問管委會 → 「我知道有規定但背不出來，不敢猜」
**擋住**：FR-11（緊急車輛進場）
```

一個檔同時服務兩種讀者，靠章節而不是靠檔案——與 `benchmark/cases/` 的 `## 輸入`
餵食線同一個手法。

---

## 職責重畫

| Skill | 拿掉 | 新增 |
| --- | --- | --- |
| **bdd-clarify** | Pass 1 第 4 步（切 story）、第 5 步（分批）；根層四個產物 | EARS FR、AC／Examples 兩層、NFR、問題與決策史寫進 `prd.md` |
| **bdd-spec** | — | 切 story（掛 `story-splitting`）；把已有答案的詞抬進 `docs/CONTEXT.md` |
| **bdd-plan** | — | 無變動 |
| **story-splitting** | 從 CLARIFY 掛到 **SPEC** | — |

---

## 腳本影響

兩支腳本加一個 skill 的查詢方式。三者都是**換資料來源，不換邏輯**。
（`clarify-loop` 不是腳本，它是 SKILL.md 裡的一行 grep。）

| 誰 | 現行 | 改成 |
| --- | --- | --- |
| `status.py` | 數 `questions/*.md` 的狀態列 | 解析 `prd.md` 的 Open Questions 表 |
| `clarify-loop` | `grep -L '狀態.*已答' questions/*.md` | 解析同一張表的「待答」列 |
| `check_spec.py` | map ↔ `.feature` 雙向覆蓋 | `prd.md` 的 `EX-N.M` ↔ `.feature` 的 `@example-N.M` |

`status.py` 的核心主張不變——它的 docstring 寫著「這是一個**算出來的視圖，不是存
起來的檔案**」，改的只是從哪裡數。

**`check_spec.py` 新增一條檢查**：每條 `FR-N` 至少要有一個 `EX-N.M`。這是 D3 兩層
結構的漂移防線——AC 改了而 Examples 沒跟上時，至少「完全沒有例子的 FR」會被抓到。

---

## 約束的例外

`PLAN.md` 未決事項 #1 現行寫著：

> 已定：文件型產物一律寫入使用端 repo 的 `specs/`，一個目錄裝完，刪掉即乾淨。
> skill 不得寫入其他位置、不得修改專案既有檔案。

改成明列**兩個例外**，各有具體理由：

| 例外 | 理由 |
| --- | --- |
| `features/*.feature` | 位置由測試 runner 的慣例決定（原未決事項 #2） |
| `docs/CONTEXT.md` | **壽命長於任何 feature。** 刪掉某個 feature 的 spec 目錄之後，「訪客」的定義應該還在——跟著 feature 被刪是錯的 |

**這個例外必須寫得很窄。** 約束原文還有「不得修改專案既有檔案」，而 `CONTEXT.md`
在使用端 repo 可能本來就存在。所以例外的措辭是：

> SPEC 得建立 `docs/CONTEXT.md`，或在其中**追加**詞彙條目。不得改寫該檔既有的
> 任何段落，不得寫入 `docs/` 底下其他檔案。

否則「不得修改專案既有檔案」會被這一個例外整條打穿。

---

## 驗證

沒有既有產物要遷移（`110819e` 清掉了訓練追蹤那一批），所以不需要相容層。

| # | 檢查 | 通過條件 |
| --- | --- | --- |
| 1 | `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify`（`bdd-spec` 同） | exit 0 |
| 2 | `claude plugin validate .` | 通過（`--strict` 會有 CLAUDE.md 那則已知警告） |
| 3 | `status.py` 與 `check_spec.py` 對一份**手寫的最小 `prd.md`** 跑 | `status.py` 數得出待答題數；`check_spec.py` 抓得到缺 EX 的 FR |
| 4 | 同兩支對**不存在的路徑**跑 | exit 1 並指名缺什麼——fail loud 的行為不可因為改資料來源而消失 |
| 5 | `grep -rn "clarify\.md\|brief\.md\|actor\.md\|<slice-slug>" skills/ docs/ PLAN.md CLAUDE.md` | 只剩刻意保留的歷史敘述 |

第 4 項不是形式。`ai-sdlc.md` 說這條鏈的特徵性失敗是「靜默地錯」而不是「壞掉」，
而這次改的正是這兩支腳本讀資料的地方——最容易改出「找不到就當成零」的那類 bug。

---

## 明確不做的事

- **跨 feature 的詞彙衝突偵測** —— `CONTEXT.md` 只是累積，不檢查同一個詞有沒有被
  兩個 feature 定義得不一樣。等真的撞到再說
- **既有產物遷移** —— 沒有既有產物
- **`bdd-plan` 的任何變動** —— 它讀 `spec.md` 與 `.feature`，兩者的格式都不變
- **`step-grammar.md` 的重寫** —— 詞彙表換家不影響封閉文法本身
- **重跑 benchmark** —— 改完之後 `parking-billing` 的實測基準會失效（量的是不同的
  流程），但重跑是另一件事，不在本次範圍
