# 三段式流程重構：discovery 只問、spec 只寫、example-mapping 求共識

**日期**：2026-09-14
**狀態**：待實作
**前一版**：[`2026-09-11-prd-format-v3-design.md`](./2026-09-11-prd-format-v3-design.md)

這一版不是格式調整，是**流程的分段方式改了**。`prd.md` 這個檔名消失，CLARIFY 停止
寫規格，Formulation 多出一個「給人審閱」的產物。

---

## 為什麼改

### 一、`CLARIFY 不寫規格` 這條紀律從來沒有真的成立過

repo 自己有三份文件這樣寫：

| 文件 | 寫著什麼 |
| --- | --- |
| `CLAUDE.md` | CLARIFY **does not write specs**, and it does not slice stories |
| `PLAN.md:42` | CLARIFY 只做問問題 ｜ **不寫規格。連 story 的定版句子都不寫** |
| `docs/ai-sdlc.md:81` | Discovery 是模型追問、使用者回答；Formulation 是模型書寫、使用者審閱——**方向相反的工作模式** |

而 v3 的 `bdd-clarify` 產出的 `prd.md` 裡有
`### US-1  身為 P-1（訪客），我想在出場前知道要付多少` ——**那正是「story 的定版
句子」**。2026-09-09 的重構、v2、v3 三次都在讓 CLARIFY 寫規格，每一次都靠措辭掩飾，
沒有一次回頭看那條紀律。

根源是**文件描述三段、實作只有兩段**：

```
文件說：  追問（產出答案）→ 書寫（產出規格）→ 審閱（確認共識）
實作是：  追問並書寫（bdd-clarify → prd.md）→ 轉成 Gherkin（bdd-spec）
```

第一段被壓成一段，所以那條規則沒有立足點；第三段根本不存在。

### 二、「確認共識」是 BDD 定義的一部分，而它沒有產物

`docs/ai-sdlc.md:99` 寫著 Formulation 的「agent 不能做」是
**確認共識——定義裡的 check for agreement 需要一個真的會反對的人**。

但目前沒有任何一份產物是為了那件事存在的。`.feature` 帶著 tag 與封閉步驟文法，
是給 runner 讀的；`prd.md` 是給 SPEC 讀的。**沒有一份是給會反對的那個人讀的。**

### 三、`prd.md` 與 `spec.md` 的分界已經站不住

v3 之後 `prd.md` 有十節、`spec.md` 只有一節（`## Stories`，記哪幾條 FR 湊成一則）。
兩份文件、一個實質內容，而分界的理由（一個是 CLARIFY 的、一個是 SPEC 的）在
CLARIFY 不再寫文件之後就消失了。

---

## 決策

### D1 · 三段，各自一件事

```
bdd-discovery    →  clarify-log.md                        追問的問答紀錄
bdd-spec         →  spec.md                               綜合成一份文件
example-mapping  →  example-map/<story-slug>.feature      寫成例子，給 PM 審閱
                 →  spec.md 的 `狀態`                      PM 點頭之後才寫
（未實作）        →  features/<story-slug>.feature         加 tag、套封閉步驟文法
```

**`bdd-discovery` 不寫規格。** 它問問題、記答案，就這樣。三個 pass（廣度／深度／
技術）從「文件的三個階段」變成「問題的三種來源」。

**`bdd-spec` 不訪談。** 它讀問答紀錄，綜合成 `spec.md`。憑空長出來的內容是缺陷——
這條沒變，但現在它有意義了，因為 SPEC 真的是第一個寫字的人。

**`example-mapping` 不做決定。** 它把 `spec.md` 寫成 Gherkin 草稿，交給 PM，然後
把 PM 的答覆寫進 `狀態`。

### D2 · `bdd-discovery` 的產物是問答紀錄，不是規格

`specs/<date>-<feature>/clarify-log.md`。形狀沿用 `clarify-loop` 單獨使用時的產物：
已解決的問題與答案、還開著的、被推翻的提案。

**它不是規格**，所以 `CLAUDE.md` 那條紀律成立。它也不是暫存——換一個 session，
`bdd-spec` 讀得到它。

MUST NOT: `clarify-log.md` 出現 `US-<n>`、`FR-<n>`、`AC-<n>.<m>` 任何一種編號。
**編號是規格的東西**，出現在這裡就代表 discovery 又在寫規格了。這一條可機械檢查。

### D3 · `spec.md` 的十節，順序固定

```markdown
## Document Overview              狀態／版本修訂歷史／核心關係人
## Background
## Goal
## Scope — In / Out
## Personas                       ### P-<n>
## User Stories                   ### US-<n> ＋ #### Q（紅卡指標）
## Non-Functional Requirements
## Assumptions / Constraints
## Success Metrics
## Open Questions                 五欄表
```

`Success Metrics` 排在 `Assumptions / Constraints` 之後而不是 `Goal` 之後：
前七節是「要做什麼」，最後三節是「怎麼知道做成了、還有什麼不知道」。

`prd.md` 這個檔名消失。`prd-format.md` 併進 `bdd-spec/references/spec-format.md`
（那份現在只描述舊的單節 `spec.md`，要整份重寫）。

### D4 · 規則與例子只住 `.feature` 草稿，`spec.md` 不掛 FR

v3 把 `#### FR` / `#### NFR` / `#### Q` 三個分組掛在 story 底下。v3 之後：

- `## User Stories` 的 `### US-<n>` 只有 **story 句子 ＋ `#### Q` 紅卡指標**
- **FR 與 AC 不在 `spec.md` 裡**，它們第一次出現是在 `.feature` 草稿
- 跨 story 的 NFR 仍在 `## Non-Functional Requirements`

**這順手溶掉 v3 的 D2。** v3 有一條規則：「PRD 的 story 分組是敘事分組，不是交付
切片，SPEC 得重新分組」。它存在是因為 CLARIFY 先分了一次、SPEC 又分一次，兩個分組
並存所以需要仲裁。現在**只有一個分組，而且 SPEC 就是分它的人**——那條規則沒有東西
可管，刪掉。

> **一條規則若只是為了調解兩個本不該同時存在的東西，正確的修法是拿掉其中一個。**

### D5 · 編號寫在 `Rule:` / `Example:` 的標題文字裡

```gherkin
Feature: 訪客計費
  身為 P-1（訪客），我想在出場前知道要付多少，以便付完就能開走

  Rule: FR-1  訪客車出場依停留時長計費，前 30 分鐘免費

    Example: AC-1.1  停 29 分 → 收 0 元
      Given 一台訪客車已入場並抽了紙票
      When 它在入場後 29 分鐘於繳費機結帳
      Then 應收金額是 0 元
```

**為什麼不用 tag。** 草稿是給 PM 讀的，`@rule-1` 對他沒有意義。但編號**必須在草稿
裡誕生**，否則可執行版本的 tag 是最後一步憑空編出來的，而
`prd-format.md` 現在寫著「編號在這裡定版——`.feature` 的 tag 回指它們，**這是這條鏈
唯一的接縫**」。接縫不能斷。

`Rule: FR-1  訪客車出場…` 對 PM 完全可讀，只是多一個前綴。

**可執行版本怎麼接**：同名檔案，`Rule:` 上方加 `@rule-1`、`Example:` 上方加
`@example-1.1`，並把步驟收進封閉文法。**加 tag，不改編號。**

### D6 · 一則 story 一個草稿檔

```
specs/<date>-<feature>/
├── clarify-log.md
├── spec.md
└── example-map/
    ├── visitor-billing.feature
    └── monthly-pass.feature
```

檔名對齊未來的 `features/<story-slug>.feature`。最後那一步是「同名檔案加 tag」，
不是「一個檔拆成三個」。

草稿放 `specs/` 而不是 `features/`：testbed 的 godog 掃 `features/`，草稿放那裡會
變成一堆 undefined scenario。

### D6b · 草稿是**寫出來的**，不是從 `spec.md` 算出來的

設計過程中一度說草稿是「從 `spec.md` 產生」的衍生檔案，可以用「重新產生一次、
跟簽入的比對」來守。**那是錯的，因為 D4。**

`spec.md` 裡沒有 FR 也沒有 AC——規則與例子第一次出現就是在草稿裡。**算不出來的
東西不能靠重新算來守。** `example-mapping` 在這裡做的是真正的 Formulation 書寫，
不是排版轉換。

好處是它不會跟 `spec.md` 漂移——兩份裝的是不同的東西，沒有第二份可以對不上。
代價是沒有「重新產生就對了」這條退路，草稿的正確性只能靠 D9 的結構檢查與 PM 的
審閱。

### D7 · 兩份 Gherkin 是刻意的，各自錨在 `spec.md`

| 檔案 | 誰寫 | 給誰 | 特徵 |
| --- | --- | --- | --- |
| `example-map/<slug>.feature` | `example-mapping` | PM | 無 tag、自由措辭、編號在標題文字 |
| `features/<slug>.feature` | （未實作） | runner | `@rule-n`／`@example-n.m`、封閉步驟文法 |

**已知弱點：沒有東西檢查兩份一致。** 草稿改了、可執行版本沒跟上，不會有人發現。

緩解的方向是「同名檔案、同編號」——可執行版本的每個 `@rule-<n>` 都應該在草稿的
同名檔裡找得到 `Rule: FR-<n>`。**那是可機械檢查的，但要等可執行那一步存在才有東西
可檢查**，所以本次不做，只把它寫下來。

### D8 · `狀態` 住 `spec.md`，由 `example-mapping` 在 PM 點頭後寫

四個值不變（`澄清中`／`待對焦`／`已對焦`／`已凍結`），但寫入者重新分配：

| 狀態 | 誰寫 | 何時 |
| --- | --- | --- |
| `澄清中` | `bdd-spec` | 開 `spec.md` 時的預設 |
| `待對焦` | `example-mapping` | 草稿產出、等 PM 看 |
| `已對焦` | `example-mapping` | **PM 說可以之後** |
| `已凍結` | （未實作） | 可執行版本產出之後 |

MUST NOT: agent 自己把狀態推到 `已對焦`。`docs/ai-sdlc.md:99` 的「確認共識需要一個
真的會反對的人」就是這一條的依據。

### D9 · `check_spec.py` 一分為二，feature 側休眠

目前它做兩類事：

| 類別 | 新流程下 |
| --- | --- |
| **`spec.md` 側**：紅卡指標懸空（`#### Q` 引用的 `Q-<n>` 要在 `## Open Questions` 表裡）、`身為 P-<n>` 指向存在的 persona | 改讀 `spec.md`，**照樣有效** |
| **草稿側**：規則有沒有例子、編號對不對得上 | **搬家**——這些原本讀 `prd.md`，現在讀草稿 `.feature`（見下） |
| **可執行 feature 側**：`@example-n.m` 涵蓋每一條 AC、有沒有孤兒 scenario | **沒有輸入**——可執行那一步還不存在 |

`v2 形式殘留` 那個檢查（v3 加的，抓 `## User Stories` 底下的 `#### FR-<n>` 標題）
要改成抓**這一版的舊形式**：`spec.md` 的 story 底下不該再有 `#### FR` / `#### NFR`
分組，出現就是沒遷移完。

feature 側的檢查**進入休眠**，而且要在腳本裡明說是休眠不是壞掉。

新的檢查對象變成草稿：
- 每個 `Rule:` 的標題都要帶 `FR-<n>`，每個 `Example:` 都要帶 `AC-<n>.<m>`
- `AC-<n>.<m>` 的 `<n>` 要等於它所在 `Rule:` 的 FR 編號
- 每個 `Rule:` 底下至少一個 `Example:`
- 每個 `Example:` 的 Given／When／Then 三行齊全
- 草稿檔名要對得上 `spec.md` 的某一則 story

`status.py` 讀 `spec.md` 的 `## Open Questions`——**那一節的格式完全不變**，所以
它只需要改讀檔名。

---

## 連帶要改的

| 檔案 | 改什麼 |
| --- | --- |
| `skills/bdd-clarify/` → `skills/bdd-discovery/` | `git mv`（**自成一筆提交**，否則 `--follow` 斷掉）；SKILL.md 重寫成「只問不寫」；產出 `clarify-log.md` |
| `skills/bdd-clarify/references/prd-format.md` | 併進 `bdd-spec/references/spec-format.md`，原檔刪除 |
| `skills/bdd-clarify/examples/minimal-prd.md` | 變成 `bdd-spec/examples/minimal-spec.md`；另加一份草稿 fixture |
| `skills/bdd-clarify/references/persona-definition.md` | 移到 `bdd-spec/`（Personas 現在是 SPEC 寫的） |
| `skills/bdd-clarify/scripts/status.py` | 移到 `bdd-spec/`；只改讀 `spec.md` |
| `skills/bdd-spec/SKILL.md` | 大改：吸收整份文件的產出職責；不再寫 `.feature` |
| `skills/bdd-spec/references/spec-format.md` | 整份重寫（十節） |
| `skills/bdd-spec/scripts/check_spec.py` | D9 的一分為二 |
| `skills/example-mapping/SKILL.md` | 從「讀 prd 攤地圖」改成「產生 Gherkin 草稿 ＋ 求共識 ＋ 寫狀態」 |
| `skills/clarify-loop/SKILL.md` | 與 `bdd-discovery` 的邊界要重畫——兩個都只問問題 |
| `skills/story-splitting/SKILL.md` | 觸發點從 CLARIFY 移到 SPEC |
| `skills/bdd-plan/SKILL.md` | 讀 `spec.md` ＋ 草稿，不再讀 `prd.md` |
| `docs/ai-sdlc.md` | 六步表的 CLARIFY／SPEC 兩列重寫；Formulation 那列補「確認共識有產物了」 |
| `CLAUDE.md`、`PLAN.md` | 三個紀律的敘述、產物表、skill 表 |
| `benchmark/skeleton/go/README.md` | 鏈的對應表 |

---

## 驗證

| # | 驗什麼 | 怎麼驗 |
| --- | --- | --- |
| 1 | `audit_skill.py` 對每個動過的 skill exit 0 | 含改名後的 `bdd-discovery` |
| 2 | `git log --follow` 穿得過改名 | 改名自成一筆提交，`--follow` 要追得到改名前 |
| 3 | `status.py` 對新 fixture 給出正確題數與面向覆蓋 | 格式沒變，只換檔名 |
| 4 | **草稿的五條新檢查各自先證明會失敗** | 缺編號的 `Rule:`、編號對不上的 `Example:`、沒有 `Example:` 的 `Rule:`、缺 Given 的 `Example:`、對不上 story 的檔名 |
| 5 | **`clarify-log.md` 出現規格編號要報錯** | D2 的那條 MUST NOT，可機械檢查 |
| 6 | feature 側休眠是**明說的**，不是靜默跳過 | 腳本要印出「這幾項待可執行步驟實作後啟用」 |
| 7 | fail loud 沒退化 | 兩支腳本對不存在路徑 exit 1 |
| 8 | `claude plugin validate .` exit 0 | 只有 root `CLAUDE.md` 那則已知警告 |
| 9 | 六種形態的殘留掃描 | 名稱／標題形式／散文術語／裸字／圍籬內／收窄後未通知——**不帶 `^` 錨點**，逐一判斷 |

---

## 已知弱點

- **兩份 Gherkin 沒有一致性檢查**（D7）。要等可執行那一步存在。
- **`example-map/` 的草稿沒有人消費**。它是給人讀的，過時了不會弄壞任何下游——但
  一張沒人維護的舊地圖比沒有地圖更糟，因為它看起來是現況。
- **`bdd-discovery` 與 `clarify-loop` 的邊界很薄**。兩個都只問問題；前者產生問題集、
  後者收斂問題集。這條界線在實作時要寫清楚，否則兩個 skill 會互相搶。
- **`clarify-loop` 的名稱衝突仍未解**——使用者全域有一個同名 skill，plugin 的這個
  很可能叫不到。不在本次範圍。

---

## 這份 spec 要切成兩個實作計畫

15 個檔案、三個 skill 換角色，塞成一份計畫會得到一份沒人能審的東西。可以切，
而且切點很自然——**兩段各自都留下一條走得通的流程**。

### 計畫一：`prd.md` 併成 `spec.md`

- `bdd-clarify` → `bdd-discovery`（`git mv` 自成一筆），改寫成「只問不寫」，
  產出 `clarify-log.md`
- `bdd-spec` 吸收整份文件的產出職責，寫 `spec.md`——**這一步 FR 與 AC 仍然留在
  `spec.md` 裡**，就是 v3 的結構換個檔名、照 D3 重排章節
- `prd-format.md` 併進 `spec-format.md`；`persona-definition.md`、`status.py` 搬家
- `check_spec.py` 與 `status.py` 改讀 `spec.md`
- `example-mapping` 只改讀檔名，行為不變（仍是攤終端地圖 ＋ 就緒判定）

做完之後流程完整可跑：discovery 問、spec 寫、example-mapping 審。**只是還沒有
Gherkin 草稿。**

### 計畫二：規則與例子搬進 `.feature` 草稿

- D4：FR／AC 從 `spec.md` 移出，`## User Stories` 只留 story 句子與紅卡指標
- D5、D6：`example-mapping` 改成產生 `example-map/<story-slug>.feature`
- D8：`狀態` 的寫入者重新分配
- D9：草稿的五條新檢查；feature 側明說休眠
- `docs/ai-sdlc.md`、`CLAUDE.md`、`PLAN.md`、benchmark README

**為什麼這樣切。** 計畫一動的是「文件住哪、誰寫」；計畫二動的是「規則與例子住哪」。
把它們混在一起，任何一個檢查失敗時都分不清是哪一半造成的——而這條 repo 前兩次遷移
的教訓都是「半遷移的狀態最危險」。

## 這份工作的規模，以及一個誠實的提醒

比 v3 大：三個 skill 改角色、一個改名、兩支腳本換輸入、一個檔名消失、兩份 reference
合併、四份文件重寫。

**而 v3 三天前才合併，一次都沒被實測過。** `benchmark/cases/personal-memo/` 藏了三個
陷阱（兩句直接矛盾的話、「消失」的定義、單機 vs 同步），專門測「模型會不會用先驗
填空而不是問」——**那正是這次重構要強化的能力**。

先跑一次會拿到真實證據：v3 的 CLARIFY 到底是在問還是在寫。若它其實在猜，這次重構的
理由會更硬；若它問得很好，會知道哪些部分不該動。

這份 spec 不預設順序，但那個 case 需要一個未被污染的 session（本 session 讀過
`grading.md`，只能當關係人或評分者）。

---

## 不在範圍

- 產生 `features/<story-slug>.feature` 的那一步（可執行版本）
- 兩份 Gherkin 的一致性檢查（等上一項存在）
- `clarify-loop` → `clarify` 改名與名稱衝突
- `runs/` 底下既有的 v2／v3 產物遷移
