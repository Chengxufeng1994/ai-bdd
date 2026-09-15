---
name: bdd-spec
description: >
  把 CLARIFY 的問答綜合成 `spec.md`——十六節、完整自足的規格文件：切 story、
  抽出 FR 與 AC、決定驗收測試打在哪一層 seam、把詞彙抬進 `docs/CONTEXT.md`。
  不訪談，只綜合已經有答案的東西；`.feature` 由下一步的 `bdd-formulation` 產出。
  BDD 六步流程的 SPEC 步驟。
  觸發詞：「產出規格文件」「寫 spec.md」「clarify-log.md 變成規格」
  「把問答整理成規格」「切 story」「決定 seam」「驗收測試打在哪一層」
  「這個 feature 的規格長怎樣」。
  English: synthesize the clarification log into spec.md, turn answered
  questions into a spec, slice stories, decide the acceptance test seam,
  write the glossary.
---

# SPEC — 把答案綜合成一份規格

一份收斂的 `clarify-log.md`（連同它問過的原始輸入）進來，出去的是
`spec.md`：**先切 story——決定哪幾條 FR 湊成一則可獨立驗收的 story——再把
已經有答案的東西綜合成 FR 與 AC**，並決定驗收測試打在哪一層 seam。

不訪談。不發明新的例子。**不產 `.feature`**——那是下一步 `bdd-formulation` 的事。

## 使用時機

- CLARIFY 判定收斂，要把 `clarify-log.md` 變成 `spec.md`
- 手上有規則與具體例子，需要一份完整自足的規格文件
- 準備跟 PM 對答案，但還沒有 `spec.md`

## Skill Boundaries

- **要把 FR 與 AC 寫成 `.feature` → 改用 `bdd-formulation`**
- 要稽核既有 `.feature` 寫得好不好 → 改用 `bdd-spec-review`
- 要決定實作順序 → 改用 `bdd-plan`（情境跑在哪一層測試是 seam，本 skill 步驟 3 決定）
- 例子還不夠、還有紅卡 → 回 `bdd-discovery`
- 產出是 `specs/<date>-<feature>/spec.md` ＋ `specs/domain-model.md` ＋ `docs/CONTEXT.md` 的追加，**不是** `.feature`、`openapi.yaml`、migration 或任何可執行的檔案
- 要把 `spec.md` 拆成可執行的票 → 改用 `bdd-plan`

## 參考檔案

- `references/spec-format.md` — `spec.md` 的骨架
- `references/persona-definition.md` — 角色怎麼定義
- `examples/minimal-spec.md` — 一份完整的 `spec.md` 範例
- `scripts/status.py` — 從 `## Open Questions` 算澄清進度
- `../bdd-formulation/scripts/check_spec.py` — 交件前自檢（第 6 步）。**在隔壁 skill**：本 skill 只跑它，不改它

---

## 切 story

**只有一個 story 分組，你就是分它的人。** 舊版有兩個——CLARIFY 先分敘事分組、
SPEC 再分交付切片——所以需要一條規則仲裁。CLARIFY 不再寫文件之後，那條規則已刪除。

每則 story 的標題是 `### US-<n> · <slug>`，`<slug>` 是它的 `.feature` 檔名。
kebab-case，一份 `spec.md` 裡唯一——`bdd-formulation` 的 `check_spec.py` 用它對應檔案。

一則 story ＝ 一個 `.feature` 檔 ＝ 一次可獨立驗收的交付。

沿規則切，不沿使用者旅程或畫面切：後者容易切出「做一半」的 story，前半段
交付了但沒有任何一條 FR 被完整滿足。

每一則要通過 INVEST（**Small 除外**，那是切分本身要解決的）。切不動、或切完
仍然太大時 → `story-splitting`，那裡有九種切分模式。

切分結果寫進 `specs/<date>-<feature>/spec.md` 的 `## User Stories`：

```markdown
### US-1 · visitor-billing

身為 P-1（訪客），我想在出場前知道要付多少，以便付完就能開走 ← 推論

涵蓋 FR-1、FR-4。沿規則切：訪客計費是一條獨立成立的約束，自己可驗收。

### US-2 · monthly-pass-exit

身為 P-2（月租戶），我想刷卡就直接出去，以便不用每次都停下來付款 ← 推論

涵蓋 FR-2。
```

**`### US-<n> · <slug>` 與它底下提到的 `FR-<n>` 是機械可解析的**——
`bdd-formulation` 的 `check_spec.py` 用它算「這則 story 該有哪些 `@example` tag」。

MUST: 寫下**為什麼是這個切法**。半年後有人要加一則 story 時，第一個該讀的
就是這一段。

MUST NOT: 切出 `spec.md` `## Scope — In / Out` 明確排除的東西。範圍要擴張就回
`bdd-discovery` 公開改範圍——悄悄擴張比公開改糟，因為沒有人有機會反對。

## 詞彙表 —— `docs/CONTEXT.md`

**下一步**寫 `.feature` 時措辭必須一致，所以詞彙在這一步就要定下來。把
`clarify-log.md`「已答」表裡**已經有答案**的詞抬進 `docs/CONTEXT.md`。

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

不必問輸入的位置。`clarify-log.md` 與 `input.md` 都在 `specs/<date>-<feature>/`
底下。唯一要徵詢的是 seam，在流程步驟 3。

MUST NOT: 在這裡問問題。**不訪談**——憑空長出來的內容是缺陷。`clarify-log.md`
的「待答」表裡還有東西時，那些對應的規則就寫不出來，寫進 `## Scope — In / Out`
的「回 CLARIFY 補問」，不要順手決定掉。

NEVER: `clarify-log.md` 不存在時，憑需求描述直接寫 `spec.md`。那會產出一份
沒有人同意過的規格，而它看起來跟真的一模一樣。沒有 `clarify-log.md` 就先跑
`bdd-discovery`。

NEVER: `input.md` 不存在時開始。`PRD §x` 這個來源標記回指的就是它——缺了它，
每一個 `← PRD §x` 都指向一份從未被寫下的文件，而且**不會報錯，只會悄悄指錯**，
於是「這句抄自需求方原文」與「這句是我們自己補的」再也分不開。回 `bdd-discovery`
補存一份逐字副本再開始。

## 流程

### 1. 挑 story —— 只做已就緒的

`spec.md` `## User Stories` 底下還沒有任何 `### US-<n>` 時，先做「切 story」
那件事——這一步的「story」指切出來的結果，不是預先存在的東西。已經切過一次
的 feature 再跑這個 skill 時，`<story-slug>` 指 `## User Stories` 底下已經
有的那個 slug。

| 呼叫方式 | 範圍 |
| --- | --- |
| `bdd-spec` | **預設：`結束方式` 已是 `收斂` 的每一份 `clarify-log.md`，整份做完——包含從它切出來的每一則 story** |
| `bdd-spec <story-slug>` | 只做指定的那一則 |

MUST: `結束方式` 還不是 `收斂` 的 `clarify-log.md` **整份**跳過，並明講跳過的
理由與該回哪個 skill——就緒住在整份 log 一格，不是每則 story 一格。

未就緒只有一種成因會擋住這一步：**`clarify-log.md` 的「待答」表裡還有東西**
→ 回 `clarify-loop` 收斂。「規則太多、範圍太大」不是未就緒——story 怎麼分本來
就是這一步自己的工作（見「切 story」），不是抄現成的分組。把還沒定案的東西
寫成場景，等於**把不確定性從一個顯眼的問題檔搬進一份看起來已完成的
規格**——之後沒有人會回頭質疑它。

### 2. 讀兩樣東西

| 讀什麼 | 為了什麼 |
| --- | --- |
| `specs/<date>-<feature>/clarify-log.md` | **規則與例子的綜合來源**——「已答」表的答案是原始素材，`歸屬` 欄說這句話該進 `spec.md` 哪一節，不必從問題文字反推；「待答」「n/a」兩張表劃出這一步碰不得的界線；`## 核心關係人` 表是 `## Document Overview` 那一節**唯一**的來源 |
| `specs/<date>-<feature>/input.md` | **`PRD §x` 這個來源標記唯一的依據**——`clarify-log.md` 只記問答，不記哪一句是原文照抄；`input.md` 是需求方帶來的原始輸入的逐字副本，沒有它就分不出「這句抄自需求方原文」與「這句是我們自己補的」 |

`clarify-log.md`「已答」表的答案是規則與例子的原始素材，但 FR／AC 的定版
編號都是**這一步**第一次寫下來的——story 由哪些 FR 組成同樣是這一步自己切
出來的，不是讀到的（見「切 story」）。

IMPORTANT: `clarify-log.md` 由 CLARIFY 維護，本 skill **不得寫入**。問答是
澄清對話的產物，在寫規格時偷偷改一句答案，等於繞過了當初定義它的那場對話。
推導不出來的寫進 `spec.md` 的 `## Scope — In / Out`「回 CLARIFY 補問」。

MUST NOT: 把 `clarify-log.md`「待答」表裡的列當成已經有答案。那一列還沒有
答案，把它寫成場景或抬進詞彙表，等於把不確定性從一個顯眼的清單搬進一份
看起來已完成的規格，之後沒有人會回頭質疑它。**只有「已答」表可以用**——
這也是「詞彙表」一節「只轉錄，不裁決」的同一條界線。

角色要雙向對得上：FR 出現 `spec.md` `## Personas` 沒有的角色 →
CLARIFY 漏了一個；`## Personas` 有但沒有任何 FR 提到 → 那個角色是憑空的。
兩種都寫進 `spec.md` 的 `## Scope — In / Out`「回 CLARIFY 補問」。

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

**Ask user: 把 seam 攤出來確認再寫。** 這是本 skill 唯一一次徵詢——
其餘全部只綜合已有的答案。之所以是例外：seam 要看過具體的 AC 才推得出來，
而它推錯的話，底下所有 step definition 都打在錯的高度，改起來是整批重寫。

MUST NOT: 把 seam 的討論變成需求訪談。這一次徵詢問的是**技術判斷**——
「打這一層對嗎」——不是「這個規則該怎麼定」。談的時候發現需求有洞，
寫進 `## Scope — In / Out` 的「回 CLARIFY 補問」，**不要順手問掉**。

理由是這一步沒有 `clarify-log.md` 的記錄機制：訪談的問答會被寫進 log、
被後面每一步讀到；在 seam 討論裡順口問到的答案只活在這一次對話裡，
下一個人打開 `spec.md` 只會看到一條沒有出處的規則。

IMPORTANT: seam 一旦寫進 `spec.md` 就往下傳——IMPLEMENT 照著打，
REVIEW 把「沒人同意過的 seam」當成 finding。這個綁定是**間接的**，
穿過 `spec.md` 這份文件；所以這一節寫不寫，決定了後面兩步抓不抓得到。

### 4. 寫 `spec.md`

一個 feature 一份，路徑 `specs/<date>-<feature>/spec.md`。骨架照抄
`references/spec-format.md`。

MUST: `## Document Overview` 的 `核心關係人` 表**照抄** `clarify-log.md` 的
`## 核心關係人`，三欄原樣搬過來，不增列也不改寫 `決策權`。誰能決定什麼是問出來的，
不是看得出來的——log 裡沒有那個人，這裡就不該有那一列。log 缺這張表時**不要自己
補三列看起來合理的**，寫一行說明它缺了、回 `bdd-discovery` 登記；`## Open Questions`
每一題的「該問誰」都要指向這張表裡的人，兩邊一起編的話那條規則會輕鬆通過而毫無意義。

MUST: `狀態` 的初值由這一步寫。建檔時還有『待答』的列就寫 `澄清中`，一列都不剩
就寫 `待對焦`。**這個檔是本 skill 建的，所以初值只能由本 skill 寫**——
`bdd-discovery` 明文不得寫 `spec.md`。之後 `已對焦` 由 `example-mapping` 依使用者
的答覆改寫，本 skill 不碰。

MUST: `## Known Boundaries` 的三張清單要填，空的也要寫為什麼空——判準與
`要先問` 的格式都在 `references/spec-format.md`，**不要在這裡抄一份**。

MUST NOT: 在這一步做新決定。**本 skill 只綜合已經有答案的東西。**
推不出來的寫進 `## Scope — In / Out` 的「回 CLARIFY 補問」，不要順手決定掉。

判準：指著 `spec.md` 的任何一句話，說得出它來自哪個已答的問題、哪條規則嗎？
說不出來的就是憑空長出來的，那是缺陷。

### 5. 增修 `specs/domain-model.md`

根層一份，跨 feature 累積。這個 feature 長出來的聚合與不變條件加進去。

MUST: 只增修，不重寫。動到既有條目要在修訂紀錄寫「原本是什麼、為什麼改」——
那張表記錄的是這個領域被理解的過程，而下一個 feature 的人靠它判斷某個決定還算不算數。

MUST NOT: 寫進 `spec.md`。`spec.md` 是這個 feature 的快照，實作一開始就會過期；
domain model 要活得比它久。塞進去等於陪葬。

判準（哪些該進來）：這條規則**跨 feature 仍然成立**嗎？只在這個 feature 的情境下為真的，
留在 `spec.md`。

### 6. 自檢，然後才算完成

```bash
python3 <ai-bdd>/skills/bdd-formulation/scripts/check_spec.py <專案根>
```

`.feature` 還不存在，所以它只跑得到八項——全部是 `spec.md` 自己的形狀：FR 有沒有掛
AC、`## User Stories` 底下有沒有 v2 的舊標題、story 標題形式對不對、`#### Q` 的形式、
紅卡指向不存在的題號、紅卡列了已答的題、`AC-<n>.<m>` 的 `<n>` 對不對得上它的 FR、
有沒有 FR 不屬於任何一則 story。它會明講另外五項跳過了。

**這八項全是「安靜地錯」的形狀**：每一行都合規、檔案照樣 parse、讀起來通順，
但下游的解析器看不到那條 FR，於是它的 AC 從覆蓋率裡消失，而**沒有任何東西會報錯**。
交件前跑一次是這條鏈上唯一一個能在寫的人還記得為什麼這樣寫的時候抓到它們的位置。

MUST: `exit` 非 0 就修完再交件。修不掉的（例如缺的答案要回頭問）寫進
`## Scope — In / Out` 的「回 CLARIFY 補問」，不要留給下一步去撞。

## 產物格式

```
specs/<date>-<feature>/
└── spec.md           SPEC  ★ 十六節，完整自足

（clarify-log.md 是 CLARIFY 的產物，本 skill 只讀不寫）
（.feature 是 bdd-formulation 的產物，本 skill 不產）
```

本步驟產三份東西，格式各有一份參考檔：

| 產物 | 格式 |
| --- | --- |
| `specs/<date>-<feature>/spec.md` | `references/spec-format.md` |
| `specs/domain-model.md`（增修） | `references/spec-format.md` 末節 |
| `docs/CONTEXT.md`（建立或追加） | 「詞彙表」一節 |

規則與例子的**定版編號在這裡誕生**：`FR-<n>`、`AC-<n>.<m>`。`bdd-formulation` 寫的
`@rule-<n>`／`@example-<n>.<m>` 回指它們——所以**這批編號一旦寫定就不准重排**。
下游回指的是 tag，不是那句規則的文字本身；重排等於讓那些引用**靜默**指向別的
東西，不會報錯，只會對錯。

## 產物隔離

MUST: 只新增本節列出的三種產物，**不修改專案既有的任何檔案**（含既有的
`.feature`、`docs/CONTEXT.md` 既有段落）。`.feature` 完全不碰——本 skill 不產它，
也不改它。

## 完成後

告訴對方五件事：

1. 切了哪些 story、為什麼這樣切（`## User Stories` 的切法說明）
2. **seam 是哪一層**、`domain-model.md` 這一批新增了哪些聚合與不變條件、
   `docs/CONTEXT.md` 追加了哪些詞彙
3. `## Open Questions` 還剩幾題待答、`## Scope — In / Out` 的「回 CLARIFY 補問」
   有幾條
4. `## Document Overview` 的 `狀態` 這一步寫成什麼（`澄清中` 還是 `待對焦`）
5. `check_spec.py` 那八項跑出什麼——過了就說「八項全過」，抓到就說抓到哪幾條、
   當場修掉了還是寫進「回 CLARIFY 補問」

缺口多 → 回 `clarify-loop` 收斂。缺口少 → 跑 `bdd-formulation` 產 `.feature`，
拿它跟 PM 對答案。
