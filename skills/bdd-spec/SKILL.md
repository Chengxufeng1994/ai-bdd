---
name: bdd-spec
description: >
  讀已對焦的 `spec.md`，寫 `design.md` 的技術設計——seam、模組邊界、介面與 API 契約、
  schema、每個 AC 的測試層級、風險——並把詞彙抬進 `docs/CONTEXT.md`、聚合與不變條件
  增修進 `specs/domain-model.md`。不訪談、不切 story、不寫 FR／AC；
  `.feature` 由 `bdd-formulation` 產出。BDD 六步流程 SPEC 步驟的技術那一半。
  觸發詞：「寫 design.md」「技術設計」「決定 seam」「驗收測試打在哪一層」
  「API 契約要怎麼定」「這個 feature 動哪些模組」「spec.md 已對焦了，接下來」。
  English: write the technical spec, decide the acceptance test seam,
  implementation decisions from the prd, write the glossary.
---

# SPEC — 把已對焦的 PRD 變成技術設計

一份狀態 `已對焦` 的 `spec.md` 進來，出去的是同目錄的 `design.md`：**驗收測試打在哪一層
seam、動哪些模組、介面與 API 契約、schema、每個 AC 在哪一層測、風險。** FR／AC／story
已經在 `spec.md` 裡定案，這一步**讀**它們，不改它們。

不訪談。不切 story。不寫 FR／AC。**不產 `.feature`**——那是 `bdd-formulation` 的事。

## 使用時機

- `spec.md` 已對焦，要決定 seam 與技術設計
- `.feature` 寫好了但還沒有 seam，`bdd-plan` 無法切票
- 要把這個 feature 的詞彙與聚合抬進跨 feature 的檔案

## Skill Boundaries

- **要把 FR 與 AC 寫成 `.feature` → 改用 `bdd-formulation`**
- `spec.md` 還有紅卡、或狀態不是 `已對焦` → 回 `bdd-discovery`／`example-mapping`
- story 切得不對 → `story-splitting`，然後回 `bdd-discovery`；本 skill 不碰 `## User Stories`
- 要決定實作順序 → 改用 `bdd-plan`
- 要稽核既有 `.feature` 寫得好不好 → `bdd-spec-review` 未實作（見 `PLAN.md`），先跑 `bdd-formulation` 的 `scripts/check_spec.py`
- 產出是 `specs/<date>-<feature>/design.md` ＋ `specs/domain-model.md` ＋ `docs/CONTEXT.md` 的追加，**不是** `spec.md`、`.feature`、`openapi.yaml`、migration

## 參考檔案

- `references/design-format.md` — `design.md` 的骨架（只有技術節）
- `../bdd-discovery/references/spec-format.md` — 讀 `spec.md` 時對照它的節
- `examples/minimal-design.md` — 一份完整的 `design.md` 範例
- `../bdd-formulation/scripts/check_spec.py` — 交件前自檢（第 6 步）。**在隔壁 skill**：本 skill 只跑它，不改它
- `../bdd-discovery/scripts/status.py` — 確認 `spec.md` 沒有待答的題

---

## story 在 `spec.md` 裡，不在這裡切

`## User Stories` 由 DISCOVERY 開場粗切、輪次沿規則細切，到這一步已經定案。
一則 story ＝ 一個 `.feature` 檔 ＝ 一次可獨立驗收的交付，slug 就是檔名。

MUST NOT: 改 `spec.md` 的 `## User Stories`。寫技術設計時覺得切法不對，寫進
`design.md` 的 `## 回 DISCOVERY 補問`，由 DISCOVERY 用 `story-splitting` 重切。
悄悄重切比公開改糟，因為 `.feature` 的檔名與 tag 已經掛在舊的切法上。

## 詞彙表 —— `docs/CONTEXT.md`

`bdd-formulation` 寫 `.feature` 時措辭必須一致，所以詞彙在這一步就要定下來。把
`spec.md` `## Open Questions` 裡狀態 `已答`、答案是一個定義的詞抬進 `docs/CONTEXT.md`。

**只轉錄，不裁決。** 詞義有爭議時不要自己定一個——那是訪談，SPEC 不訪談。
回 `bdd-discovery` 把它變成一題。

### 條目格式

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

---

## 前置確認（Ask user）

不必問輸入的位置。`spec.md` 在 `specs/<date>-<feature>/`。唯一要徵詢的是 seam，在流程步驟 3。

MUST NOT: 在這裡問需求。**不訪談**——`spec.md` 推不出來的，寫進 `design.md` 的
`## 回 DISCOVERY 補問`，不要順手決定掉。

NEVER: `spec.md` 不存在、或狀態不是 `已對焦` 時開始。`澄清中` 代表還有紅卡，
`待對焦` 代表需求方還沒點頭——兩種都可能改 FR，而 seam 與 schema 是照 FR 定的。
回 `bdd-discovery` 或 `example-mapping`。

NEVER: `spec.md` 沒有 `## 原文` 時開始。`← PRD §x` 指的就是它。

## 流程

### 1. 挑 feature

| 呼叫方式 | 範圍 |
| --- | --- |
| `bdd-spec` | **預設：狀態 `已對焦` 的每一份 `spec.md`** |
| `bdd-spec <feature>` | 只做指定的那一個 |

MUST: 狀態不是 `已對焦` 的 `spec.md` 整份跳過，並明講跳過的理由與該回哪個 skill。

### 2. 讀 `spec.md` 的這幾節

| 讀什麼 | 為了什麼 |
| --- | --- |
| `## User Stories` | 每條 FR 與 AC——seam 要看過具體的 AC 才推得出來；`## Testing Decisions` 要替每個 AC 標測試層級 |
| `## Open Questions` 裡 `面向` 是技術面向、狀態 `已答` 的列 | **技術決定的原始素材**——小節的 `歸屬` 說它該進 `design.md` 哪一節，行尾 `← Q-<n>` |
| `## Assumptions / Constraints`、`## Non-Functional Requirements` | 外部限制與品質屬性，`## Tech Preferences & Constraints` 與 `## Risks` 的來源 |
| `## Document Overview` 的核心關係人 | **照抄**進 `design.md` 的同名表 |
| `## 原文` | `← PRD §x` 的依據 |

IMPORTANT: `spec.md` 由 DISCOVERY 維護，本 skill **不得寫入**。推導不出來的
寫進 `design.md` 的 `## 回 DISCOVERY 補問`。

### 3. 決定 seam —— 驗收測試打在哪一層

`.feature` 決定驗什麼，seam 決定**打進系統的哪一層**：HTTP handler？
application service？domain？這個決定會長成 step definition 的形狀，
所以它屬於規格，不屬於實作。

四條規則：

| 規則 | 為什麼 |
| --- | --- |
| 優先用既有 seam | 每多一個 seam 就多一套測試替身與資料建構 |
| 取**最高**的那一層 | 越高越接近使用者真正做的事，也越不綁實作細節 |
| 整個變更的理想數量是**一個** | 兩個 seam 代表這批行為的入口不只一個，多半是切分沒切乾淨 |
| 需要新 seam 就在**你能做到的最高點**提出 | 往下挪一層通常是為了好寫，代價是驗到的東西離使用者更遠 |

**Ask user: 把 seam 攤出來確認再寫。** 這是本 skill 唯一一次徵詢。DISCOVERY 的技術
面向通常已經問過 seam（`## Open Questions` 裡 `面向` 為 `seam` 的題）——有的話這裡
只是確認，行尾 `← Q-<n>`；沒有的話，提出你的判斷、標 `← 推論`，並在
`## 回 DISCOVERY 補問` 記一條「seam 沒問過」。

MUST NOT: 把 seam 的討論變成需求訪談。談的時候發現需求有洞，寫進
`## 回 DISCOVERY 補問`，**不要順手問掉**——這一步沒有把答案寫進 `spec.md` 的權限，
順口問到的答案只活在這一次對話裡。

IMPORTANT: seam 一旦寫進 `design.md` 就往下傳——IMPLEMENT 照著打，
REVIEW 把「沒人同意過的 seam」當成 finding。這個綁定是**間接的**，
穿過 `design.md` 這份文件；所以這一節寫不寫，決定了後面兩步抓不抓得到。

### 4. 寫 `design.md`

路徑 `specs/<date>-<feature>/design.md`，骨架照抄 `references/design-format.md`——只有技術節。

MUST: `## Document Overview` 的核心關係人表**照抄** `spec.md`，不增列也不改寫 `決策權`。

MUST: `## Known Boundaries` 的三張清單要填，空的也要寫為什麼空。

MUST: `## 回 DISCOVERY 補問` 要有——沒有洞就寫「（無）」。有洞就一條一列，
寫清楚影響什麼定不下來。

MUST NOT: 在這一步做需求決定。**本 skill 只綜合已經有答案的東西。** 判準：指著
`design.md` 的任何一句話，說得出它來自 `spec.md` 的哪條 FR、哪題已答的 Q 嗎？

### 5. 增修 `specs/domain-model.md`

根層一份，跨 feature 累積。這個 feature 長出來的聚合與不變條件加進去。

MUST: 只增修，不重寫。動到既有條目要在修訂紀錄寫「原本是什麼、為什麼改」——
那張表記錄的是這個領域被理解的過程，而下一個 feature 的人靠它判斷某個決定還算不算數。

MUST NOT: 寫進 `design.md`。`design.md` 是這個 feature 的快照，實作一開始就會過期；
domain model 要活得比它久。塞進去等於陪葬。

判準（哪些該進來）：這條規則**跨 feature 仍然成立**嗎？只在這個 feature 的情境下為真的，
留在 `design.md`。

### 6. 自檢，然後才算完成

```bash
python3 <ai-bdd>/skills/bdd-formulation/scripts/check_spec.py <專案根>
```

它讀的是 `spec.md`。`.feature` 還不存在時只跑得到十項——全部是 `spec.md` 自己的形狀，
包含第 14 項「每條 FR／AC 有來源」。它在這裡跑是因為**你剛讀完 `spec.md`**：抓到的
每一條都是 DISCOVERY 漏的，寫進 `## 回 DISCOVERY 補問`，不要自己修 `spec.md`。

## 產物格式

```
specs/<date>-<feature>/
└── design.md           SPEC  ★ 技術設計

（spec.md 是 DISCOVERY 的產物，本 skill 只讀不寫）
（.feature 是 bdd-formulation 的產物，本 skill 不產）
```

| 產物 | 格式 |
| --- | --- |
| `specs/<date>-<feature>/design.md` | `references/design-format.md` |
| `specs/domain-model.md`（增修） | `references/design-format.md` 末節 |
| `docs/CONTEXT.md`（建立或追加） | 「詞彙表」一節 |

## 產物隔離

MUST: 只新增本節列出的三種產物，**不修改專案既有的任何檔案**（含既有的
`.feature`、`docs/CONTEXT.md` 既有段落）。`.feature` 完全不碰——本 skill 不產它，
也不改它。

## 完成後

告訴對方四件事：

1. **seam 是哪一層**、是確認 `spec.md` 裡問過的還是這一步推的
2. `domain-model.md` 這一批新增了哪些聚合與不變條件、`docs/CONTEXT.md` 追加了哪些詞彙
3. `## 回 DISCOVERY 補問` 有幾條
4. `check_spec.py` 那十項跑出什麼——過了就說「十項全過」，抓到就說抓到哪幾條、
   寫進了補問清單

補問清單非空 → 回 `bdd-discovery`。空 → 跑 `bdd-formulation`（若還沒跑）或 `bdd-plan`。
