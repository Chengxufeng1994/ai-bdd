---
name: spec
description: >
  讀 `discovery` 的訪談紀錄 `interview-log.md`（已答的題與原文），寫 `spec.md`——原文逐字照抄、問答整理成
  `## Open Questions`（題號不變），再寫 story、FR／AC 行尾 `← Q-<n>`／`← PRD §n`、角色、範圍與不做的事、假設與限制，交 `example-mapping`
  判就緒，需求方點頭就停，下一步由使用者叫 `plan`（`.feature` 由它 DELEGATE `formulation` 寫）。只寫要做成什麼，不寫技術決定（那是 `plan`）；
  詞彙寫進 `## Domain Concepts`，等 `archive` 才進 `docs/`。
  不訪談：有洞就在 `## Open Questions` 開一題待答、狀態改回澄清中、停下。BDD 六步流程的 to spec 步驟。
  觸發詞：「寫 spec」「to spec」「把問答寫成規格」「discovery 問完了，接下來」
  「寫成 FR 跟 AC」「這個 feature 的規格」「spec.md」「寫詞彙表」。
  English: write the spec, turn the discovery into a spec, to spec, write FRs and
  acceptance criteria from the answers, write the glossary.
disable-model-invocation: true
---

# SPEC — 把問答寫成規格

| | |
| --- | --- |
| **職責** | 把已答的問答寫成規格：要做成什麼、不做什麼，帶到需求方說「就緒」 |
| **讀** | `specs/ongoing/<date>-<feature>/interview-log.md`（結束方式 `收斂` 或 `喊停`）；再跑時既有的 `spec.md`；`docs/CONTEXT.md`（有的話） |
| **寫** | 同一個目錄的 `spec.md`，唯一的寫入者：照抄的 `## 原文`、整理的 `## Open Questions`、其餘各節、FR／AC、狀態 `待對焦`；有洞時補問題與狀態 `澄清中` |
| **上游** | `discovery` |
| **下游** | `feature-mapping` ＋ `example-mapping` ＋ `spec-by-example`（判就緒）→ `已對焦` → 交回，使用者接著叫 `plan` |
| **不做** | 不訪談、不寫 `interview-log.md`、不填紀錄裡沒有的答案、不改寫 `## 原文`、FR／AC 不標 `推論`、不自己寫 `已對焦`、不寫 `.feature` 也不呼叫 `formulation` 或 `discovery`、不寫 `docs/`、不寫技術決定與 seam、不排風險與並行（都是 `plan`）、不切票 |

**不訪談。** 這一步只寫已經有答案的東西。問答推不出來的，在 `## Open Questions` 開一題 `待答`、
把狀態改回 `澄清中`、停下來請使用者跑 `/discovery`——不要順手問掉：答案要有 `**誰答的**` 與原話才能被
`← Q-<n>` 引用，那是 `discovery` 的紀律，順口問到的答案只活在這一次對話裡。

## 使用時機

- `discovery` 以 `收斂`（或 `喊停`）結束了，要把問答寫成規格
- 開過補問題、`discovery` 問完又收斂了，要修補 `spec.md`
- 要定這個 feature 的詞彙（寫在 `spec.md`，做完時由 `archive` 併進 `docs/`）
- 下游 skill 說 `spec.md` 還留著舊的聚合一節（`## Domain Model`／`### 聚合與不變條件`），要清掉它：同目錄 `plan.md`
  已有 `## Domain Model` 時只刪這一節，狀態不動（見〈再跑：修補，不重寫〉第 5、6 步）
- 要改一個已經在 `specs/sleeping/` 的 feature：先喚醒它，見 `references/wake-up.md`

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| 結束方式是 `阻塞`／`零進展`／`—`（還沒問完） | `discovery` |
| 要把 FR 與 AC 寫成 `.feature` | `plan`——它在定了 seam 之後 DELEGATE `formulation`；本 skill 不呼叫它 |
| 要決定 seam、模組、契約、測試層級 | `plan`（寫在 `plan.md`） |
| 要定聚合、不變條件、schema、API 契約 | `plan` |
| 要攤整個 feature 的流程、找缺的替代流程與反例 | `feature-mapping`，本 skill 寫完 `spec.md` 之後第一個呼叫它 |
| 要攤四色卡、判就緒 | `example-mapping`，本 skill 在 `feature-mapping` 之後呼叫它 |
| 要把規格寫成業務方讀的驗收旅程 | `spec-by-example`，本 skill 在診斷之後、就緒問句之前呼叫它 |
| 寫規則時 story 切得不對 | `story-splitting`，`已對焦` 之前本 skill 可以呼叫它重切 |
| 要排風險、決定要不要並行、列接合點 | `plan`（寫在 `plan.md`） |
| 要切票、排實作順序 | `plan` |

產出是 `spec.md`，**不是** `interview-log.md`（那是 `discovery` 的證據，只讀）、`plan.md`、`docs/` 底下任何檔案、程式碼；
`.feature` 由 `formulation` 寫。

## 參考檔案

- `references/spec-format.md` — `spec.md` 的骨架（含 `## Open Questions` 與 `## 原文`）、來源標記、編號規則
- `examples/minimal-spec.md` — 一份 `已對焦` 的 `spec.md`，照 `../discovery/examples/interview-log.md` 寫成（題號、原話、原文對得上）
- `../discovery/references/interview-log-format.md` — 訪談紀錄的五節與一題的欄位。**第 1 步讀紀錄時讀**
- `references/wake-up.md` — 把 `specs/sleeping/` 的 feature 搬回來、狀態改回 `澄清中`。只在喚醒時讀
- `../formulation/scripts/check_spec.py` — 交件前自檢（第 3 步）。**在隔壁 skill**：本 skill 只跑它，不改它

---

## 前置確認（Ask user）

不必問輸入的位置。紀錄在 `specs/ongoing/<date>-<feature>/interview-log.md`，由 `discovery` 寫；看它
`## Remaining Open Items` 的 `**結束方式**`，收兩種：

| 結束方式 | 怎麼做 |
| --- | --- |
| `收斂` | 照下面的流程 |
| `喊停` | 使用者決定帶著開著的題往下走。照流程寫，還『待答』的題在對應 story 的 `#### Q` 放指標；下面說的「有待答的補問題」不數喊停帶過來的題（定義見 `references/spec-format.md`〈寫規格時發現的洞〉） |

NEVER: 結束方式 `阻塞`、`零進展`、`—` 時開始。那代表還有人該答的題沒答，
寫出來的規格會建立在猜測上。告訴使用者跑 `/discovery`。

NEVER: 紀錄沒有 `## 原文` 時開始。`← PRD §x` 指的就是它照抄過來的那一節。

只有 `discovery.md`（2026-11-18 以前）、或只有 `spec.md` 沒有紀錄（11-18 當天 `discovery` 寫進 `spec.md` 的單檔佈局）時，
先照 `../discovery/scripts/status.py` 印的退件訊息把紀錄補出來，再開始。

## 流程

每一步都寫成一串動作，開頭的詞說明這一步在做什麼：

| 詞 | 意思 |
| --- | --- |
| `READ` | 讀檔，不寫 |
| `WRITE` | 寫進 `spec.md`（本 skill 不寫別的檔） |
| `CHECK` | 對照一條規則判是或否；不成立就照那一行寫的去處 |
| `RUN` | 跑腳本，輸出留到回報 |
| `DELEGATE` | 讀另一個 skill 的 `SKILL.md`，照它做完再回來；只用在第 5 步 |
| `ASK` | 問使用者，等他回答 |
| `STOP` | 回報後結束這一輪，不往下做 |

每一步最後一行是**做完的樣子**與下一步；做完的樣子不成立，就不往下走。

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| 1 | 從紀錄抄原文與問答，寫需求部分（`待對焦`） | — |
| 2 | 收尾檢查，必要時重切 story | — |
| 3 | 自檢 `check_spec.py` | — |
| 4 | 詞彙 → `## Domain Concepts` | — |
| 5 | `feature-mapping` 與 `example-mapping` 診斷 → `spec-by-example` 寫旅程 → 就緒問句 | 要，由 `example-mapping` 問需求方 |
| 6 | `已對焦` → 交回，停 | — |
| 7 | 有洞 → 開補問題、改回 `澄清中`，停 | — |

### 1. 寫需求部分

1. READ `references/spec-format.md`，骨架照抄；`examples/minimal-spec.md` 是寫好的樣子。
2. READ `interview-log.md` 的五節（格式 `../discovery/references/interview-log-format.md`）。
3. WRITE 先把紀錄抄進 `spec.md`——這幾樣是**照抄**，不是改寫：

   | 紀錄 | `spec.md` |
   | --- | --- |
   | `## 原文` | `## 原文`，逐字，`§` 編號不變 |
   | `## 問答` 的每一題 | `## Open Questions` 一列＋一個小節，**題號不變**：狀態照抄（`失效` 寫 `n/a`），答案欄寫它的 `**收斂結果**`；小節放 `**誰答的**`／`**該問誰**`、`**來自**` 與答案原話、例子、被推翻的猜測 |
   | 各場 Session 的核心關係人表 | `### 核心關係人`，合成一張，不增列也不改寫 `決策權` |
   | 建檔那場的 `**需求方確認建檔**` | `### 版本修訂歷史` 的 0.1 那一列 |
   | 粗切 story | `## User Stories` 的 `### US-<n> · <slug>` 與「身為…」一句 |
   | `## Remaining Open Items` 的 `**結束方式**` | `## Document Overview` 的 `結束方式` |

   **同一個答案在兩處，是刻意的。** 紀錄是逐字的原始證據，跟 `## 原文` 一樣沒有人改；`spec.md` 的表是整理過、
   給規則引用的版本——`← Q-<n>` 指它，`check_spec.py` 只讀它，`spec.md` 因此自成一體。兩處對不上時**以紀錄的原話為準**，
   照紀錄改 `spec.md`；不是反過來。
4. WRITE 照下表把每一樣寫到它的位置；沒有的節照骨架的順序插進去。來源只有剛抄進來的 `## 原文` 與 `## Open Questions`：

   | 讀什麼 | 寫成什麼 |
   | --- | --- |
   | `## 原文` | `## Context / Problem`、`## Goals / Non-Goals`、`## Success Metrics` 的句子，行尾 `← PRD §n`；原文裡的假設與限制 |
   | `## Open Questions` 已答、答案是範圍／角色／限制／品質屬性的題 | 對應的那一節（`## Goals / Non-Goals`、`## Personas`、`## Assumptions / Constraints`、`## Non-Functional Requirements`），行尾 `← Q-<n>` |
   | 改變了行為的已答題 | 所屬 story 底下的 `**FR-<n>**`；小節裡記下的例子寫成 `- **AC-<n>.<m>**`，行尾 `← Q-<n>` |
   | 角色那幾題的小節 | `## Personas` 的五欄 |
   | `## User Stories` 的粗切 | 標題不動（slug 是 `.feature` 檔名）；`身為` 那一句的角色換成 `P-<n>`，底下加 FR／AC |
   | 問優先級的那一題 | 每則 story 標題下一行的 `**優先級**：P1／P2／P3 —— <為什麼>`；沒問過寫「（沒問過）」並開補問題 |
   | 意圖陳述那一題（六行在紀錄的 `## Confirmed Decisions`，抄進它的小節） | `## Goals / Non-Goals`、`## Success Metrics` 可以引它 |
   | 技術面向的題（`面向` 是 seam、模組邊界、契約、既有資產、測試慣例） | 不寫，`plan` 寫技術決定時讀 |

5. CHECK 每條 FR 與 AC 的行尾**只准** `← Q-<n>`（那一題是 `已答`）或 `← PRD §n`。找不到來源 → WRITE
   補問題（下面〈補問題〉），不寫那條規則。
6. CHECK 每條 AC 的值都照小節裡記下的例子。一條 AC 用到兩題的例子（例如一題給規則、另一題給
   「第 14 天休館 → 順延一天」）可以組合，兩個來源都標；要**新算一個數字**才寫得出來的 → 補問題。
7. WRITE `## Edge Cases`：照下表一類一類找，每一條指向涵蓋它的 AC；找不到 AC 的寫「（沒有 AC）」並開補問題：

   | 類別 | 例子 |
   | --- | --- |
   | 失敗 | 付款失敗、外部系統沒回應 |
   | 併發 | 兩個人同時做同一件事 |
   | 空值 | 一筆都沒有、欄位沒填 |
   | 重複 | 同一個動作做兩次 |
   | 權限邊界 | 沒權限的人、資格剛過期的人 |

8. WRITE `## Document Overview` 的 `狀態` 寫 `待對焦`。

MUST: FR 與 AC 不准標 `推論`。一條規則找不到來源，就是問答裡沒有它。

MUST NOT: 小節只有結論、沒有例子時，從結論編一個例子。那條 FR 的 AC 寫不出來，就開補問題。

#### 補問題

寫的時候發現的洞，在 `## Open Questions` 表接著編號開一題，格式照 `references/spec-format.md`〈寫規格時發現的洞〉：
狀態 `待答`、`面向` `—`、小節寫 `**該問誰**`、`**來自**：補問 —「<發現它的那句話，逐字>」` 與 `**影響**`；影響到的
story 句子下面就地標 `[NEEDS CLARIFICATION: … —— 見 Q-<n>]`、`#### Q` 放它的指標。本文說「開補問題」都是這個動作。

MUST NOT: 填補問題的答案欄，或寫任何一題紀錄裡沒有的答案——答案只從紀錄抄。也不寫 `interview-log.md`：補問題只開在 `spec.md`，
`discovery` 續跑時讀它、在紀錄追加同號的題。

**做完的樣子**：每條 FR／AC 都有來源、每個 AC 的值都對得到例子、每則 story 有優先級、五類邊界都找過；寫不出來的都開成了補問題。→ 第 2 步

### 2. 收尾檢查

只有寫出 FR 之後才做得了的檢查：

1. CHECK **角色**：`## Personas` 的每個角色至少被一條 FR 的主詞或一則 story 的 `身為` 用到，FR 與 story
   用到的角色都在 `## Personas` 裡；五欄裡寫「沒問過」的也算缺。不成立 → 補問題。
2. CHECK **型別**：每則 story 底下的 FR 要嘛全是改變世界的（命令），要嘛全是只讀的（查詢）。不成立 →
   READ `../story-splitting/SKILL.md`，沿型別切開——`.feature` 一檔只能是 `@command` 或 `@query`。
3. CHECK **大小**：一則 story 的 FR（藍卡）不多到 `example-mapping` 會判「太大」。不成立 → 同上，先切。
4. WRITE 重切的結果寫進 `## User Stories`，取代粗切——粗切只是起點。

**`已對焦` 之前才可以重切**：還沒有 `.feature` 掛在舊的切法上。`已對焦` 之後要重切，只能開一題補問題，
`discovery` 問過再改，並記修訂。

**做完的樣子**：三項都成立，或不成立的都切了、開成了補問題。→ 第 3 步

### 3. 自檢

1. RUN `python3 <ai-bdd>/skills/formulation/scripts/check_spec.py <專案根>`。
2. CHECK 抓到的每一條：寫錯了 → 修 `spec.md`；問答裡沒有 → 補問題。
3. CHECK 再跑（修補）時 `.feature` 已經存在，覆蓋比對會報這一輪新增或改動的 AC「漏了」。**只有這幾條**是
   預期的：本 skill 交回之後由 `plan` DELEGATE `formulation` 補，不開補問題、也不跳過。

`.feature` 還不存在時跑得到十三項：`spec.md` 的形狀、它的 `← PRD §n` 對 `## 原文` 的引用、`acceptance/` 旅程的
`@covers-` tag，包含第 14 項（每條 FR／AC 有來源、`← Q-<n>` 在 `## Open Questions` 找得到而且是 `已答`）。

**做完的樣子**：除了上面預期的「漏了」，全過；輸出留著貼進回報。→ 有待答的補問題：第 7 步；否則第 4 步

### 4. 詞彙

有待答的補問題時**不做這一步**，直接第 7 步：詞彙要建立在收斂的規則上，補問的答案可能改掉它們。

1. READ `docs/CONTEXT.md`（專案的現況詞彙；沒有就是第一個 feature）。
2. WRITE `## Domain Concepts` 的 `### 詞彙`：`## Open Questions` 裡狀態 `已答`、答案是一個定義的詞，一詞一列。欄位固定，
   `archive` 才能原樣追加進 `docs/CONTEXT.md`：

   | 詞 | 定義 | 來源 | 為什麼要定義 |
   | --- | --- | --- | --- |
   | 到期日 | 借出當天加 14 天，遇休館日順延到下一個開館日 | Q-6 | 「兩週內還」的兩週不是 14 個開館日 |

   第四欄不是裝飾——一個不需要解釋為什麼要定義的詞，通常本來就沒有歧義，列進來只是雜訊。
3. CHECK 詞義有爭議、或跟 `docs/CONTEXT.md` 裡同一個詞的定義不同 → 補問題，問清楚是改定義還是
   換一個詞。**只轉錄，不裁決**：自己挑一個定義就是訪談，SPEC 不訪談。

MUST NOT: 寫 `docs/CONTEXT.md` 或 `docs/` 底下任何檔案。`docs/` 只放做出來了的
東西：詞彙等這個 feature 做完，由 `archive` 併進去。

**做完的樣子**：`### 詞彙` 寫了（沒有內容的寫為什麼空）；有爭議的詞開成了補問題。→ 有待答的補問題：第 7 步；
否則第 5 步

### 5. 診斷、寫旅程、問就緒

沒有待答的補問題時才進來。

1. DELEGATE `../feature-mapping/SKILL.md`：攤整個 feature 的流程、跑三條診斷；缺口開成補問題，`spec`
   漏寫例子的列在它的判讀裡。
2. DELEGATE `../example-mapping/SKILL.md`：攤四色卡、跑四條診斷，**只到診斷，先不提就緒問句**。上一步開了
   補問題也照樣跑——兩張圖的洞一次交給 `discovery`，需求方一次問完。
3. WRITE 兩張圖指出「`spec` 漏寫了例子」的（答案在來源裡，AC 沒寫）→ 回第 1 步補 AC，再從第 2 步跑；
   診斷要切 story 的 → 回第 2 步。
4. CHECK 有待答的補問題 → 第 7 步，不寫旅程。
5. DELEGATE `../spec-by-example/SKILL.md`：把 `spec.md` 寫成 `acceptance/*.feature` 的驗收旅程、跑
   `check_spec.py`。定不下來的地方它標 `# [need clarification]`、開成補問題。
6. CHECK 有待答的補問題 → 第 7 步。
7. DELEGATE `../example-mapping/SKILL.md` 的〈就緒判定〉：給判讀、提就緒問句，需求方對著旅程答。
8. ASK （由 `example-mapping` 問）需求方答「就緒」→ `example-mapping` 寫 `已對焦`；答未就緒 → 第 7 步。

MUST NOT: 自己寫 `已對焦`。`example-mapping` 是唯一寫它的 skill，寫的是需求方的決定。有待答的補問題時不提
就緒問句——已知有洞的規格答「就緒」也沒有意義。

**做完的樣子**：`已對焦`，由 `example-mapping` 照需求方的回答寫。→ 第 6 步

### 6. 交回，然後停

需求方答「就緒」、`example-mapping` 寫了 `已對焦`，這份規格就完成了。

1. RUN `check_spec.py`：除了第 3 步說的覆蓋比對「漏了」，其餘全過。
2. STOP 照〈完成後〉回報，告訴使用者下一步是 `/plan`，然後結束這一輪。

MUST NOT: 呼叫 `formulation`、自己寫 `.feature`，或接著做 `plan`。每個 skill 都由使用者明確叫用，一個 skill
在內部接著跑下一個，等於繞過了那個決定。

MUST NOT: 寫任何技術決定——Tech Stack、seam、模組、契約、測試層級。那些由 `plan` 寫在 `plan.md`；需求方
判就緒讀的是整份 `spec.md`，裡面不該有他判斷不了的東西。

**做完的樣子**：`已對焦`、`check_spec.py` 只剩預期的「漏了」；回報後停下。

### 7. 有洞：開了補問題，停

四種情況會到這裡：開了補問題、`spec-by-example` 在旅程裡標了 `# [need clarification]`、
`example-mapping` 判不就緒、收尾檢查要的資訊問答裡沒有。後三種也都要先開成補問題。

1. CHECK 每一題待答的補問題都有 `**影響**`（什麼定不下來）與 `**該問誰**`（格式見 `references/spec-format.md`）。
2. WRITE `## Document Overview` 的狀態改回 `澄清中`；文件到過 `已對焦` 的，照修訂歷史的觸發條件記一列。
3. STOP 告訴使用者跑 `/discovery`——它續跑時讀 `spec.md` 的待答題、問完記進紀錄。本 skill 不叫它。

**做完的樣子**：每題補問題都寫了影響與該問誰、狀態已改回 `澄清中`；回報後停下。

### 再跑：修補，不重寫

`discovery` 問完補問題、又以 `收斂` 結束之後再跑本 skill，**修補**既有的 `spec.md`，從第 2 步重跑：

0. READ 紀錄的 `## 問答`；WRITE 照第 1 步的抄法更新 `## Open Questions`——紀錄有、表上沒有的題補上，狀態或答案不同的照紀錄改。
1. WRITE 新的 FR／AC 接著編號，**不重排**——`.feature` 的 `@example-n.m` 掛在上面，重排會讓 tag 安靜地
   指到另一條規則。
2. WRITE 文件到過 `已對焦` 的，改動的每一條在 `### 版本修訂歷史` 記一列（改了什麼、為什麼、`← Q-<n>`）；
   從沒到過 `已對焦` 的不記——觸發條件見 `references/spec-format.md`。
3. WRITE 補問題答掉的，把答案寫成 FR／AC，刪掉它在 story 底下的 `[NEEDS CLARIFICATION]` 與 `#### Q` 指標——題目本身
   留在表裡；引用了被「取代」的舊題的規則，改引新題。
4. CHECK 舊版 `spec.md` 的章節照 `references/spec-format.md` 的順序改名併節（`## Background` →
   `## Context / Problem`、`## Goal`／`## Scope`／Out → `## Goals / Non-Goals`、`## 回 DISCOVERY 補問` → `## Open Questions`），不改內容；
   缺的節補上。舊的 `## Open Questions` 若是一串沒有編號的洞，每一條開成一題補問題。
   `check_spec.py` 第 18 項在已對焦時擋缺 `## Domain Concepts`。
5. CHECK 舊版 `spec.md` 還有 `## Domain Model` 或 `### 聚合與不變條件` → 同目錄 `plan.md` 已有 `## Domain Model`
   就從 `spec.md` 刪掉這一節；沒有就留著並**停下**，請使用者先重跑 `/plan`（`plan` 讀這一節當輸入），再重跑 `/spec`。
   `spec` 自己不搬。
6. WRITE 狀態回 `待對焦`，從第 2 步重跑。**例外**：這次唯一的改動是第 5 步刪掉舊的聚合一節（同目錄 `plan.md`
   已有 `## Domain Model`）→ 狀態不動、不重跑——規則一條都沒變，不用再對焦一次。

## 產物

| 產物 | 格式 |
| --- | --- |
| `specs/ongoing/<date>-<feature>/spec.md` | `references/spec-format.md` |
| `spec.md` 的 `## Domain Concepts` | `references/spec-format.md` |
| `interview-log.md` | DISCOVERY 寫的；本 skill 只讀，`spec.md` 的 `## 原文` 與 `## Open Questions` 照它抄 |
| `features/<story-slug>.feature` | 由 `plan` DELEGATE `formulation` 寫 |

MUST: 只新增或修改本節列出的產物，**不修改專案既有的任何其他檔案**（含
`docs/` 底下的每一個檔案）。`interview-log.md` 一個字都不動。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 每條 FR／AC 行尾都是 `← Q-<n>`（該題 `已答`）或 `← PRD §n`，沒有 `推論`
- [ ] 每條 AC 的值都來自問答裡記下的例子；只有結論沒有例子的 FR 開成了補問題，沒有編一個
- [ ] `## 原文` 與紀錄的逐字相同；`## Open Questions` 每一題的題號、狀態、答案原話對得上紀錄；新開的補問題都是 `待答`，答案欄 `—`
- [ ] `interview-log.md` 沒有被本 skill 改過
- [ ] 跑了 `check_spec.py`，輸出貼在回報裡；沒過的只有覆蓋比對裡這一輪新改、`.feature` 還沒跟上的 AC
- [ ] 有待答的補問題、或旅程有 `# [need clarification]` → 狀態改回 `澄清中`、停下，沒有提就緒問句
- [ ] `已對焦` 是 `example-mapping` 照需求方的回答寫的——本 skill 沒有自己寫
- [ ] `spec.md` 裡沒有任何技術決定（Tech Stack、seam、模組、契約、測試層級）
- [ ] 再跑時 FR／AC 編號沒有重排，改動的每一條都記了修訂
- [ ] 再跑時舊版的 `## Domain Model`／`### 聚合與不變條件`：`plan.md` 已有 `## Domain Model` 才刪，否則留著並停下請使用者先重跑 `/plan`；只刪了這一節時狀態沒有改回 `待對焦`
- [ ] 詞彙寫在 `## Domain Concepts`，跟 `docs/CONTEXT.md` 定義不同的詞開成了補問題；`docs/` 底下沒有任何檔案被碰
- [ ] 每則 story 有優先級（或「（沒問過）」並開了補問題）；`## Edge Cases` 每一條指向 AC 或一題補問題
- [ ] 沒有呼叫 `formulation` 或 `discovery`，也沒有寫任何 `.feature`

## 完成後

告訴對方：

1. `spec.md` 寫了幾則 story、幾條 FR、幾條 AC；收尾檢查重切了哪幾則、為什麼
2. `check_spec.py` 跑出什麼——過了就說「全過」，抓到就說抓到哪幾條、怎麼處理
3. `## Domain Concepts` 新增或改了哪些詞彙
4. 就緒問句的結果：`已對焦`，或開了哪幾題補問題（`Q-<n>`）、為什麼
5. `### Non-Goals` 列了幾項延後、幾項不做；`## Edge Cases` 有幾條找不到 AC

`已對焦` → 下一步由使用者叫 `/plan`（它先 DELEGATE `formulation` 寫 `.feature`，再做設計與切票）。開了補問題 → 使用者叫 `/discovery`。
