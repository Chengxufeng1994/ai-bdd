---
name: bdd-spec
description: >
  讀已收斂的 `discovery.md`，寫 `spec.md`——story、FR／AC（行尾 `← Q-<n>`／`← PRD §n`）、
  角色、範圍、假設與限制——跑收尾檢查、交 `example-mapping` 判就緒，需求方點頭後呼叫
  `bdd-formulation` 產 `.feature`；並把詞彙抬進 `docs/CONTEXT.md`、聚合與不變條件增修進
  `specs/domain-model.md`。不訪談，有洞就退回 `bdd-discovery`。BDD 六步流程的 to spec 步驟。
  觸發詞：「寫 spec」「to spec」「把問答寫成規格」「discovery 問完了，接下來」
  「寫成 FR 跟 AC」「這個 feature 的規格」「spec.md」「寫詞彙表」。
  English: write the spec, turn the discovery into a spec, to spec, write FRs and
  acceptance criteria from the answers, write the glossary.
---

# SPEC — 把問答寫成規格

一份 `已收斂` 的 `discovery.md` 進來，出去的是同目錄的 `spec.md`（`已對焦`）與它的
`.feature`：**每則 story 底下有規則有例子、每一條行尾指得出 `discovery.md` 裡的哪一題或
哪一段原文、需求方看過說「就緒」。**

**不訪談。** 這一步只寫已經有答案的東西。問答推不出來的，寫進 `spec.md` 的
`## 回 DISCOVERY 補問`，把 `discovery.md` 退回 `澄清中`，不要順手問掉——這一步沒有把答案
寫進 `discovery.md` 的權限，順口問到的答案只活在這一次對話裡。

不寫 seam、API、schema、測試層級——那是 `bdd-plan` 的事，寫在 `plan.md`。

## 使用時機

- `bdd-discovery` 收斂了，要把問答寫成規格
- `discovery.md` 被退回、問完一輪又收斂了，要修補 `spec.md`
- 要把這個 feature 的詞彙與聚合抬進跨 feature 的檔案

## Skill Boundaries

- `discovery.md` 還有待答、或結束方式是 `阻塞`／`零進展` → 回 `bdd-discovery`
- 要把 FR 與 AC 寫成 `.feature` → `bdd-formulation`，本 skill 在 `已對焦` 之後呼叫它
- 要攤四色卡、判就緒 → `example-mapping`，本 skill 寫完 `spec.md` 之後呼叫它
- 要把規格寫成業務方讀的驗收旅程 → `bdd-spec-by-example`，本 skill 在診斷之後、就緒問句之前呼叫它
- 寫規則時 story 切得不對 → `story-splitting`，`已對焦` 之前本 skill 可以呼叫它重切
- 要決定 seam、模組、契約 → `bdd-plan`；實作順序與切票 → `bdd-tasks`
- 產出是 `spec.md` ＋ `specs/domain-model.md` ＋ `docs/CONTEXT.md` 的追加，**不是**
  `discovery.md`、`plan.md`、程式碼；`.feature` 由 `bdd-formulation` 寫

## 參考檔案

- `references/spec-format.md` — `spec.md` 的骨架、來源標記、編號規則
- `references/domain-model-format.md` — `specs/domain-model.md` 的骨架
- `examples/minimal-spec.md` — 一份 `已對焦` 的 `spec.md`；它的問答在
  `../bdd-discovery/examples/minimal-discovery.md`
- `../bdd-discovery/references/discovery-format.md` — 讀 `discovery.md` 時對照它的節
- `../bdd-formulation/scripts/check_spec.py` — 交件前自檢（第 3 步）。**在隔壁 skill**：本 skill 只跑它，不改它

---

## 前置確認（Ask user）

不必問輸入的位置。`discovery.md` 在 `specs/<date>-<feature>/`。

收兩種 `discovery.md`：

| 狀態／結束方式 | 怎麼做 |
| --- | --- |
| `已收斂`（結束方式 `收斂`） | 照下面的流程 |
| `澄清中`，結束方式 `喊停` | 使用者決定帶著開著的題往下走。照流程寫，還『待答』的題在對應 story 的 `#### Q` 放指標 |

NEVER: 結束方式 `阻塞`、`零進展`，或 `澄清中` 而沒有喊停時開始。那代表還有人該答的題沒答，
寫出來的規格會建立在猜測上。回 `bdd-discovery`。

NEVER: `discovery.md` 沒有 `## 原文` 時開始。`← PRD §x` 指的就是它。

## 流程

```
discovery.md（已收斂 或 喊停）
   ↓ 1. 寫 spec.md（待對焦）
   ↓ 2. 收尾檢查（可重切 story）
   ↓ 3. 自檢 check_spec.py
   ↓ 4. 詞彙 → docs/CONTEXT.md、聚合 → specs/domain-model.md
   ↓ 5. example-mapping 地圖、四條診斷 → bdd-spec-by-example 寫 acceptance/ 旅程
        → example-mapping 就緒問句
   ├─ 需求方答「就緒」→ 已對焦 → 6. bdd-formulation 產 .feature
   └─ 不就緒／補問清單非空／旅程有 # [need clarification] → 7. 退回 discovery
```

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

MUST: FR 與 AC 的行尾**只准** `← Q-<n>`（那一題必須是 `已答`）或 `← PRD §n`。一條規則
找不到來源，就是問答裡沒有它——寫進 `## 回 DISCOVERY 補問`，不要標 `推論` 把它寫出來。

MUST: AC 的值照小節裡記下的例子寫。小節只有結論、沒有例子的，那條 FR 的 AC 寫不出來
——進補問清單，不要從結論編一個例子。

`歸屬` 是 `plan.md` 技術節的題不寫進這裡，留給 `bdd-plan`。

狀態寫 `待對焦`。

### 2. 收尾檢查

只有寫出 FR 之後才做得了的檢查，原本在 discovery 收尾，現在在這裡：

- **角色**：`## Personas` 的每個角色至少被一條 FR 的主詞或一則 story 的 `身為` 用到；
  FR 與 story 用到的角色都在 `## Personas` 裡
- **型別**：每則 story 底下的 FR 要嘛全是改變世界的（命令），要嘛全是只讀的（查詢）。
  混著的那則用 `story-splitting` 沿型別切開——`.feature` 一檔只能是 `@command` 或 `@query`
- **大小**：一則 story 的 FR（藍卡）多到 `example-mapping` 會判「太大」，先用
  `story-splitting` 切

**`已對焦` 之前可以重切 story**：還沒有 `.feature` 掛在舊的切法上。`discovery.md` 的
`## 粗切 story` 是起點，不跟著改——story 的定案位置是 `spec.md` 的 `## User Stories`。
`已對焦` 之後要重切只能回 discovery 開一題，走一輪再改，並記修訂。

### 3. 自檢

```bash
python3 <ai-bdd>/skills/bdd-formulation/scripts/check_spec.py <專案根>
```

`.feature` 還不存在時跑得到十二項，全部是 `spec.md` 的形狀與它對 `discovery.md` 的引用，
包含第 14 項：每條 FR／AC 有來源、`← Q-<n>` 在 `discovery.md` 找得到而且是 `已答`。
抓到的每一條要嘛修 `spec.md`（寫錯了），要嘛進補問清單（問答裡沒有）。

**再跑（修補）時 `.feature` 已經存在**，`check_spec.py` 會多跑覆蓋比對：這一輪新增或改動
的 AC 在既有 `.feature` 裡還找不到 `@example-n.m`，會報「漏了」。**只有這幾條**是預期的——
`bdd-formulation` 在第 6 步才補——不要為它們開補問、也不要跳過。其餘每一條照上面處理。
第 6 步之後再跑一次，這時要全過才算完成。

### 4. 詞彙與聚合

#### 詞彙表 —— `docs/CONTEXT.md`

`bdd-formulation` 寫 `.feature` 時措辭必須一致，所以詞彙在這一步就要定下來。把
`discovery.md` `## Open Questions` 裡狀態 `已答`、答案是一個定義的詞抬進 `docs/CONTEXT.md`。

**只轉錄，不裁決。** 詞義有爭議時不要自己定一個——那是訪談，SPEC 不訪談。
進補問清單，回 `bdd-discovery` 把它變成一題。

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

### 5. 交給 `example-mapping` 與 `bdd-spec-by-example`

補問清單是「（無）」時，依序：

1. `example-mapping` 攤地圖、跑四條診斷
2. `bdd-spec-by-example` 把 `spec.md` 寫成 `acceptance/*.feature` 的驗收旅程、跑
   `check_spec.py`。它把定不下來的地方標成 `# [need clarification]` 並寫進補問清單
3. 補問清單仍是「（無）」→ 交回 `example-mapping` 提就緒問句，需求方對著旅程答

`example-mapping` 是唯一把狀態寫成 `已對焦` 的 skill，寫的是需求方的決定。補問清單非空時
不提就緒問句——就緒問句建立在已知有洞的規格上，答「就緒」也沒有意義；直接走第 7 步。

### 6. `已對焦` → 呼叫 `bdd-formulation`

需求方答「就緒」、`example-mapping` 寫了 `已對焦`，呼叫 `bdd-formulation` 把 FR 與 AC
寫成 `.feature`（`@ready`），`acceptance/` 的旅程只供它參考用語。`.feature` 已經存在（修補的情況）時，它照自己處理
「規則改了」的方式改場景。寫完再跑一次 `check_spec.py`，十七項全過才算 to spec 完成。

### 7. 退回 discovery

四種情況：補問清單非空、`bdd-spec-by-example` 在旅程裡標了 `# [need clarification]`、
`example-mapping` 判不就緒、收尾檢查要的資訊問答裡沒有。

- `spec.md`：`## 回 DISCOVERY 補問` 一條一列，寫清楚影響什麼定不下來；狀態留 `待對焦`
- `discovery.md`：狀態改回 `澄清中`，`### 版本修訂歷史` 記一列「回頭：<原因>」
- 告訴使用者跑 `bdd-discovery`——它續跑時會讀這份補問清單，把每一條開成新題

### 再跑：修補，不重寫

`discovery.md` 又收斂之後再跑本 skill，**修補**既有的 `spec.md`：

- FR／AC 編號**不重排**——`.feature` 的 `@example-n.m` 掛在上面，重排會讓 tag 安靜地
  指到另一條規則。新的接著編號
- 改動的每一條在 `### 版本修訂歷史` 記一列（改了什麼、為什麼、`← Q-<n>`）
- 補問清單裡被答掉的列刪除（判準：`discovery.md` 有一題的 `**來自**` 逐字引了它、而且已答）；引用了被「取代」的舊題的規則，改引新題
- 狀態回 `待對焦`，從第 2 步重跑

## 產物格式

```
specs/<date>-<feature>/
├── discovery.md      DISCOVERY 的產物；本 skill 只讀，退回時改狀態與修訂歷史
└── spec.md           SPEC  ★ 業務規格

specs/domain-model.md           增修
docs/CONTEXT.md                 建立或追加
features/<story-slug>.feature   由 bdd-formulation 寫
```

| 產物 | 格式 |
| --- | --- |
| `specs/<date>-<feature>/spec.md` | `references/spec-format.md` |
| `specs/domain-model.md`（增修） | `references/domain-model-format.md` |
| `docs/CONTEXT.md`（建立或追加） | 第 4 步「詞彙表」 |

## 產物隔離

MUST: 只新增或修改本節列出的產物，**不修改專案既有的任何其他檔案**（含
`docs/CONTEXT.md` 既有段落）。`discovery.md` 只動兩處：退回時的狀態與修訂歷史一列。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 每條 FR／AC 行尾都是 `← Q-<n>`（該題 `已答`）或 `← PRD §n`，沒有 `推論`
- [ ] 每條 AC 的值都來自問答裡記下的例子；只有結論沒有例子的 FR 進了補問清單，沒有編一個
- [ ] 跑了 `check_spec.py`，輸出貼在回報裡；`.feature` 已存在時，沒過的只有這一輪新改的 AC
- [ ] 補問清單非空、或旅程有 `# [need clarification]` → 已走〈退回 discovery〉，沒有提就緒問句
- [ ] `spec.md` 狀態是 `待對焦`——本 skill 沒有自己寫 `已對焦`
- [ ] 再跑時 FR／AC 編號沒有重排，改動的每一條都記了修訂
- [ ] `docs/CONTEXT.md` 只有新建或追加，既有段落沒動；`docs/` 底下其他檔案沒碰

## 完成後

告訴對方：

1. `spec.md` 寫了幾則 story、幾條 FR、幾條 AC；收尾檢查重切了哪幾則、為什麼
2. `check_spec.py` 跑出什麼——過了就說「全過」，抓到就說抓到哪幾條、怎麼處理
3. `domain-model.md` 新增了哪些聚合與不變條件、`docs/CONTEXT.md` 追加了哪些詞彙
4. 就緒問句的結果：`已對焦` 並已產 `.feature`，或退回 discovery 的原因與補問條數

`已對焦` → 下一步 `bdd-plan`。退回 → `bdd-discovery`。
