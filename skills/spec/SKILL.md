---
name: spec
description: >
  讀已收斂的 `discovery.md`，寫一份自成一體的 `spec.md`——需求部分（story、FR／AC 行尾
  `← Q-<n>`／`← PRD §n`、角色、範圍、假設與限制）交 `example-mapping` 判就緒；需求方點頭後
  寫技術部分（Tech Stack、Commands、Project Structure、Code Style、Boundaries、
  Implementation Decisions、Testing Decisions，seam 攤給使用者確認），再呼叫 `formulation`
  產 `.feature`；並把詞彙抬進 `docs/CONTEXT.md`、聚合增修進 `specs/domain-model.md`。
  不訪談，有洞就退回 `discovery`。BDD 六步流程的 to spec 步驟。
  觸發詞：「寫 spec」「to spec」「把問答寫成規格」「discovery 問完了，接下來」
  「寫成 FR 跟 AC」「這個 feature 的規格」「spec.md」「寫詞彙表」「決定 seam」
  「驗收測試打在哪一層」「API 契約要怎麼定」「這個 feature 動哪些模組」「技術設計」「schema 怎麼定」。
  English: write the spec, turn the discovery into a spec, to spec, write FRs and
  acceptance criteria from the answers, write the glossary, decide the acceptance test seam,
  technical design, define the API contract.
---

# SPEC — 把問答寫成規格

| | |
| --- | --- |
| **職責** | 把已收斂的問答寫成一份讀完就能開工的規格：需求部分帶到需求方說「就緒」，再寫技術部分（用什麼做、seam、模組、契約、每個 AC 在哪一層測），交出去產 `.feature` |
| **讀** | `specs/<date>-<feature>/discovery.md`（`已收斂`，或 `澄清中` 且結束方式 `喊停`）；寫技術部分時加專案既有的結構與文件（Makefile、`CLAUDE.md`、架構文件），或上一份 `spec.md` 的技術部分 |
| **寫** | `specs/<date>-<feature>/spec.md`、`specs/domain-model.md`（增修）、`docs/CONTEXT.md`（建立或追加）；退回時 `discovery.md` 的狀態與修訂歷史一列 |
| **上游** | `discovery` |
| **下游** | `example-mapping` ＋ `spec-by-example`（判就緒）→ `已對焦` 後寫技術部分 → `formulation` → `plan` |
| **不做** | 不訪談、FR／AC 不標 `推論`、不自己寫 `已對焦`、不寫 `.feature`、不排風險與並行（那是 `plan`）、不切票 |

**不訪談。** 這一步只寫已經有答案的東西。問答推不出來的，寫進 `spec.md` 的
`## 回 DISCOVERY 補問`，把 `discovery.md` 退回 `澄清中`，不要順手問掉——這一步沒有把答案
寫進 `discovery.md` 的權限，順口問到的答案只活在這一次對話裡。

## 使用時機

- `discovery` 收斂了，要把問答寫成規格
- `discovery.md` 被退回、問完一輪又收斂了，要修補 `spec.md`
- 要把這個 feature 的詞彙與聚合抬進跨 feature 的檔案

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| `discovery.md` 還有待答、或結束方式是 `阻塞`／`零進展` | `discovery` |
| 要把 FR 與 AC 寫成 `.feature` | `formulation`，本 skill 在 `已對焦` 之後呼叫它 |
| 要攤四色卡、判就緒 | `example-mapping`，本 skill 寫完 `spec.md` 之後呼叫它 |
| 要把規格寫成業務方讀的驗收旅程 | `spec-by-example`，本 skill 在診斷之後、就緒問句之前呼叫它 |
| 寫規則時 story 切得不對 | `story-splitting`，`已對焦` 之前本 skill 可以呼叫它重切 |
| 要排風險、決定要不要並行、列接合點 | `plan`（寫在 `plan.md`） |
| 要切票、排實作順序 | `task-breakdown` |

產出是 `spec.md` ＋ `specs/domain-model.md` ＋ `docs/CONTEXT.md` 的追加，**不是**
`discovery.md`、`plan.md`、程式碼；`.feature` 由 `formulation` 寫。

## 參考檔案

- `references/spec-format.md` — `spec.md` 需求部分的骨架、來源標記、編號規則
- `references/tech-format.md` — `spec.md` 技術部分七節的骨架、四選一的來源標記、範例。**第 6 步寫技術部分之前先讀**
- `references/domain-model-format.md` — `specs/domain-model.md` 的骨架
- `examples/minimal-spec.md` — 一份 `已對焦` 的 `spec.md`；它的問答在
  `../discovery/examples/minimal-discovery.md`
- `../discovery/references/discovery-format.md` — 讀 `discovery.md` 時對照它的節
- `../formulation/scripts/check_spec.py` — 交件前自檢（第 3 步）。**在隔壁 skill**：本 skill 只跑它，不改它

---

## 前置確認（Ask user）

不必問輸入的位置。`discovery.md` 在 `specs/<date>-<feature>/`。收兩種：

| 狀態／結束方式 | 怎麼做 |
| --- | --- |
| `已收斂`（結束方式 `收斂`） | 照下面的流程 |
| `澄清中`，結束方式 `喊停` | 使用者決定帶著開著的題往下走。照流程寫，還『待答』的題在對應 story 的 `#### Q` 放指標 |

NEVER: 結束方式 `阻塞`、`零進展`，或 `澄清中` 而沒有喊停時開始。那代表還有人該答的題沒答，
寫出來的規格會建立在猜測上。回 `discovery`。

NEVER: `discovery.md` 沒有 `## 原文` 時開始。`← PRD §x` 指的就是它。

## 流程

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| 1 | 寫 `spec.md`（`待對焦`） | — |
| 2 | 收尾檢查，必要時重切 story | — |
| 3 | 自檢 `check_spec.py` | — |
| 4 | 詞彙 → `docs/CONTEXT.md`、聚合 → `specs/domain-model.md` | — |
| 5 | `example-mapping` 診斷 → `spec-by-example` 寫旅程 → 就緒問句 | 要，由 `example-mapping` 問需求方 |
| 6 | 需求方答「就緒」→ `已對焦` → 寫技術部分 → `formulation` 產 `.feature` | 要，seam 攤給使用者確認 |
| 7 | 不就緒、補問清單非空、或旅程有 `# [need clarification]` → 退回 discovery | — |

### 1. 寫 `spec.md`

骨架照抄 `references/spec-format.md`。來源只有 `discovery.md`：

| 讀什麼 | 寫成什麼 |
| --- | --- |
| `## 原文` | Background／Goal／Scope／Success Metrics 的句子，行尾 `← PRD §n`；原文裡的假設與限制 |
| `## Open Questions` 已答、`歸屬` 是 `spec.md` 某一節的題 | 那一節的內容，行尾 `← Q-<n>` |
| 改變了行為的已答題 | 所屬 story 底下的 `**FR-<n>**`；小節裡記下的例子寫成 `- **AC-<n>.<m>**`，行尾 `← Q-<n>` |
| 角色那幾題的小節 | `## Personas` 的五欄 |
| `## 粗切 story` | `### US-<n> · <slug>` 與 `身為…` 那一句（slug 照抄，`.feature` 檔名靠它） |
| `### 核心關係人` | **照抄**，不增列也不改寫 `決策權` |
| `歸屬` 是技術部分某一節的題 | 這一步不寫，留到第 6 步寫技術部分；`歸屬` 是 `Risks` 的留給 `plan` |

MUST: FR 與 AC 的行尾**只准** `← Q-<n>`（那一題必須是 `已答`）或 `← PRD §n`。一條規則
找不到來源，就是問答裡沒有它——寫進 `## 回 DISCOVERY 補問`，不要標 `推論` 把它寫出來。

MUST: AC 的值照小節裡記下的例子寫。小節只有結論、沒有例子的，那條 FR 的 AC 寫不出來
——進補問清單，不要從結論編一個例子。

狀態寫 `待對焦`。

### 2. 收尾檢查

只有寫出 FR 之後才做得了的檢查，原本在 discovery 收尾，現在在這裡：

| 檢查 | 要成立的 | 不成立時 |
| --- | --- | --- |
| **角色** | `## Personas` 的每個角色至少被一條 FR 的主詞或一則 story 的 `身為` 用到；FR 與 story 用到的角色都在 `## Personas` 裡 | 要的資訊問答裡沒有 → 第 7 步 |
| **型別** | 每則 story 底下的 FR 要嘛全是改變世界的（命令），要嘛全是只讀的（查詢） | 用 `story-splitting` 沿型別切開——`.feature` 一檔只能是 `@command` 或 `@query` |
| **大小** | 一則 story 的 FR（藍卡）不多到 `example-mapping` 會判「太大」 | 先用 `story-splitting` 切 |

**`已對焦` 之前可以重切 story**：還沒有 `.feature` 掛在舊的切法上。`discovery.md` 的
`## 粗切 story` 是起點，不跟著改——story 的定案位置是 `spec.md` 的 `## User Stories`。
`已對焦` 之後要重切只能回 discovery 開一題，走一輪再改，並記修訂。

### 3. 自檢

```bash
python3 <ai-bdd>/skills/formulation/scripts/check_spec.py <專案根>
```

`.feature` 還不存在時跑得到十二項：`spec.md` 的形狀、它對 `discovery.md` 的引用，以及
`acceptance/` 旅程的 `@covers-` tag，
包含第 14 項：每條 FR／AC 有來源、`← Q-<n>` 在 `discovery.md` 找得到而且是 `已答`。
抓到的每一條要嘛修 `spec.md`（寫錯了），要嘛進補問清單（問答裡沒有）。

**再跑（修補）時 `.feature` 已經存在**，`check_spec.py` 會多跑覆蓋比對：這一輪新增或改動
的 AC 在既有 `.feature` 裡還找不到 `@example-n.m`，會報「漏了」。**只有這幾條**是預期的——
`formulation` 在第 6 步才補——不要為它們開補問、也不要跳過。其餘每一條照上面處理。
第 6 步之後再跑一次，這時要全過才算完成。

### 4. 詞彙與聚合

補問清單非空、接下來要走第 7 步退回時，這一步**先不做**：詞彙與聚合要建立在收斂的規則上，退回
之後的答案可能改掉它們。退回、再收斂、修補完之後再做。

#### 詞彙表 —— `docs/CONTEXT.md`

`formulation` 寫 `.feature` 時措辭必須一致，所以詞彙在這一步就要定下來。把
`discovery.md` `## Open Questions` 裡狀態 `已答`、答案是一個定義的詞抬進 `docs/CONTEXT.md`。

**只轉錄，不裁決。** 詞義有爭議時不要自己定一個——那是訪談，SPEC 不訪談。
進補問清單，回 `discovery` 把它變成一題。

欄位固定，這樣不同 feature 各自執行的 SPEC 才接得上同一份文件：

| 詞 | 定義 | 來源 feature | 為什麼要定義 |
| --- | --- | --- | --- |
| 訪客 | 沒有月租證的臨停車輛 | 2026-09-09-parking-fee-collection | 跟「月租戶」對照才看得出計費規則分岔在哪 |

第四欄不是裝飾——一個不需要解釋為什麼要定義的詞，通常本來就沒有歧義，
列進來只是雜訊。

MUST: 只得建立 `docs/CONTEXT.md`，或在其中**追加**條目。
MUST NOT: 改寫該檔既有的任何段落；寫入 `docs/` 底下其他任何檔案。

這是「文件型產物一律寫入 `specs/`」的第二個例外（第一個是 `.feature`，位置由
runner 決定）。理由是**壽命**：刪掉某個 feature 的 spec 目錄之後，「訪客」的
定義應該還在。

#### 聚合 —— `specs/domain-model.md`

根層一份，跨 feature 累積，骨架見 `references/domain-model-format.md`。這個 feature
長出來的聚合與不變條件加進去。

MUST: 只增修，不重寫。動到既有條目要在修訂紀錄寫「原本是什麼、為什麼改」——
那張表記錄的是這個領域被理解的過程，而下一個 feature 的人靠它判斷某個決定還算不算數。

判準（哪些該進來）：這條規則**跨 feature 仍然成立**嗎？只在這個 feature 的情境下為真的，
留在 `spec.md` 的 FR。

### 5. 交給 `example-mapping` 與 `spec-by-example`

補問清單是「（無）」時，依序：

1. `example-mapping` 攤地圖、跑四條診斷
2. `spec-by-example` 把 `spec.md` 寫成 `acceptance/*.feature` 的驗收旅程、跑
   `check_spec.py`。它把定不下來的地方標成 `# [need clarification]` 並寫進補問清單
3. 補問清單仍是「（無）」→ 交回 `example-mapping` 提就緒問句，需求方對著旅程答

`example-mapping` 是唯一把狀態寫成 `已對焦` 的 skill，寫的是需求方的決定。補問清單非空時
不提就緒問句——就緒問句建立在已知有洞的規格上，答「就緒」也沒有意義；直接走第 7 步。

### 6. `已對焦` → 寫技術部分 → 呼叫 `formulation`

需求方答「就緒」、`example-mapping` 寫了 `已對焦`，接著寫技術部分。需求方判就緒時只看需求
部分；技術決定要建立在講定的需求上，所以技術部分在這之後才寫。

#### 讀

| 讀什麼 | 為了什麼 |
| --- | --- |
| `spec.md` 的 `## User Stories` | 每條 FR 與 AC——seam 要看過具體的 AC 才推得出來，測試層級要逐個 AC 標 |
| `discovery.md` `## Open Questions` 裡技術面向、狀態 `已答` 的題 | **技術決定的原始素材**——小節的 `歸屬` 說它該進技術部分哪一節，行尾標 `← Q-<n>` |
| `spec.md` 的 `## Assumptions / Constraints`、`## Non-Functional Requirements` | 外部限制與品質屬性，`## Tech Stack` 與 `## Boundaries` 的來源 |
| 專案的 Makefile、README、`CLAUDE.md`、架構文件、既有程式碼 | `## Commands`、`## Project Structure`、`## Code Style` 的來源，標 `← 專案 <路徑>`；找既有的 seam |
| 上一份 `spec.md` 的技術部分（有的話） | 抄過來，逐條確認現在還成立 |

**一次讀完整批，不分批。** API 契約與 schema 天生跨 story，只有並排才看得見衝突。

#### 決定 seam —— 驗收測試打在哪一層

`.feature` 決定驗什麼，seam 決定**打進系統的哪一層**：HTTP handler？application service？
domain？這個決定會長成 step definition 的形狀。

| 規則 | 為什麼 |
| --- | --- |
| 優先用既有 seam | 每多一個 seam 就多一套測試替身與資料建構 |
| 取**最高**的那一層 | 越高越接近使用者真正做的事，也越不綁實作細節 |
| 整個變更的理想數量是**一個** | 兩個 seam 代表這批行為的入口不只一個，多半是切分沒切乾淨 |
| 需要新 seam 就在**你能做到的最高點**提出 | 往下挪一層通常是為了好寫，代價是驗到的東西離使用者更遠 |

**Ask user: 把 seam 攤出來確認再寫。** 三種情況：

| `discovery.md` 裡 `面向` 為 `seam` 的題 | 怎麼做 |
| --- | --- |
| `已答` | 這裡只是確認，行尾標 `← Q-<n>` |
| `n/a`（答案欄寫了為什麼不適用，例如沒有對外介面） | 照那個理由選最外層能驗到行為的一層，行尾標 `← Q-<n>` |
| 沒有這題，或還 `待答` | 提出你的判斷、標 `← 推論`，並在 `## Further Notes` 的「推論的技術決定」記一條「seam 沒問過」 |

MUST NOT: 把 seam 的討論變成需求訪談。談的時候發現需求有洞，寫進 `## 回 DISCOVERY 補問`，
**不要順手問掉**——這一步沒有把答案寫進 `discovery.md` 的權限，順口問到的答案只活在這一次
對話裡。

IMPORTANT: seam 一旦寫進 `spec.md` 就往下傳——IMPLEMENT 照著打，REVIEW 把「沒人同意過的
seam」當成 finding。這個綁定穿過 `spec.md` 這份文件；所以這一節寫不寫，決定了後面兩步
抓不抓得到。

#### 寫技術部分

照 `references/tech-format.md` 填七節：`## Tech Stack`、`## Commands`、`## Project Structure`、
`## Code Style`、`## Boundaries`、`## Implementation Decisions`、`## Testing Decisions`。每條行尾標
來源（四選一）；技術部分可以用 `推論`，但每一條同時列進 `## Further Notes` 的「推論的技術決定」——
FR／AC 不准 `推論` 的規則不變。

寫技術部分時發現的洞：

| 種類 | 例子 | 怎麼處理 |
| --- | --- | --- |
| **技術選擇** | 哪個 library、欄位叫什麼 | `## Further Notes`，留給 IMPLEMENT 決定 |
| **需求洞** | 某條規則在某個情況下該怎麼做，沒人說過 | `## 回 DISCOVERY 補問` 加一列、狀態改回 `待對焦`，走第 7 步 |

#### 呼叫 `formulation`

技術部分寫完，呼叫 `formulation` 把 FR 與 AC
寫成 `.feature`（`@ready`），`acceptance/` 的旅程只供它參考用語。`.feature` 已經存在
（修補的情況）時，它照自己的〈再跑：`spec.md` 改過之後〉，依修訂歷史只改對到的場景。寫完再跑一次
`check_spec.py`，沒有任何一項失敗才算 to spec 完成。

### 7. 退回 discovery

四種情況：補問清單非空、`spec-by-example` 在旅程裡標了 `# [need clarification]`、
`example-mapping` 判不就緒、收尾檢查要的資訊問答裡沒有。

| 檔案 | 寫什麼 |
| --- | --- |
| `spec.md` | `## 回 DISCOVERY 補問` 一條一列，寫清楚影響什麼定不下來；狀態留 `待對焦` |
| `discovery.md` | 狀態改回 `澄清中`，`### 版本修訂歷史` 記一列「回頭：<原因>」 |

告訴使用者跑 `discovery`——它續跑時會讀這份補問清單，把每一條開成新題。

### 再跑：修補，不重寫

`discovery.md` 又收斂之後再跑本 skill，**修補**既有的 `spec.md`：

- FR／AC 編號**不重排**——`.feature` 的 `@example-n.m` 掛在上面，重排會讓 tag 安靜地
  指到另一條規則。新的接著編號
- 文件到過 `已對焦` 的，改動的每一條在 `### 版本修訂歷史` 記一列（改了什麼、為什麼、`← Q-<n>`）；
  從沒到過 `已對焦` 的不記——觸發條件見 `references/spec-format.md`
- 補問清單裡被答掉的列刪除（判準：`discovery.md` 有一題的 `**來自**` 逐字引了它、而且已答）；
  引用了被「取代」的舊題的規則，改引新題
- 狀態回 `待對焦`，從第 2 步重跑

## 產物

| 產物 | 格式 |
| --- | --- |
| `specs/<date>-<feature>/spec.md`（需求部分＋技術部分） | `references/spec-format.md`、`references/tech-format.md` |
| `specs/domain-model.md`（增修） | `references/domain-model-format.md` |
| `docs/CONTEXT.md`（建立或追加） | 第 4 步「詞彙表」 |
| `specs/<date>-<feature>/discovery.md` | DISCOVERY 的產物；本 skill 只讀，退回時改狀態與修訂歷史 |
| `features/<story-slug>.feature` | 由 `formulation` 寫 |

MUST: 只新增或修改本節列出的產物，**不修改專案既有的任何其他檔案**（含
`docs/CONTEXT.md` 既有段落）。`discovery.md` 只動兩處：退回時的狀態與修訂歷史一列。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 每條 FR／AC 行尾都是 `← Q-<n>`（該題 `已答`）或 `← PRD §n`，沒有 `推論`
- [ ] 每條 AC 的值都來自問答裡記下的例子；只有結論沒有例子的 FR 進了補問清單，沒有編一個
- [ ] 跑了 `check_spec.py`，輸出貼在回報裡；`.feature` 已存在時，沒過的只有這一輪新改的 AC
- [ ] 補問清單非空、或旅程有 `# [need clarification]` → 已走〈退回 discovery〉，沒有提就緒問句
- [ ] `已對焦` 是 `example-mapping` 照需求方的回答寫的——本 skill 沒有自己寫
- [ ] 技術部分在 `已對焦` 之後才寫，七節都在（空的寫了為什麼空）；每條行尾有來源，每條 `← 推論` 同時列在 `## Further Notes`
- [ ] seam 攤給使用者確認過，回覆記在回報裡；談 seam 時發現的需求洞進了補問清單，沒有順口問掉
- [ ] `## Commands` 有跑場景與閘門兩行，或寫了「（沒有——要問使用者）」並進了補問清單
- [ ] 再跑時 FR／AC 編號沒有重排，改動的每一條都記了修訂
- [ ] `docs/CONTEXT.md` 只有新建或追加，既有段落沒動；`docs/` 底下其他檔案沒碰

## 完成後

告訴對方：

1. `spec.md` 寫了幾則 story、幾條 FR、幾條 AC；收尾檢查重切了哪幾則、為什麼
2. `check_spec.py` 跑出什麼——過了就說「全過」，抓到就說抓到哪幾條、怎麼處理
3. `domain-model.md` 新增了哪些聚合與不變條件、`docs/CONTEXT.md` 追加了哪些詞彙
4. 就緒問句的結果：`已對焦` 並已產 `.feature`，或退回 discovery 的原因與補問條數
5. 技術部分：**seam 是哪一層**、是確認 `discovery.md` 裡問過的還是推的、使用者確認了沒；
   動哪些模組、新的契約有哪些、幾個決定標了 `推論`

`已對焦`、技術部分寫完、`.feature` 產了 → 下一步 `plan`。退回 → `discovery`。
