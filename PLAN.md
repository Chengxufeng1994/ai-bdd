# ai-bdd 路線圖

> **這份是 ai-bdd 這個 plugin 自己**要建什麼、怎麼切、什麼順序。
> 它**不是**本流程的產物——流程產出的計畫叫 `docs/bdd/<slice>/plan.md`，
> 那是使用端 repo 裡的東西。看到「plan」先分清楚是哪一個。

BDD 實踐本身的定義（三個實踐、雙迴圈、Gherkin 反模式）見
[docs/bdd-workflow.md](./docs/bdd-workflow.md)。

> **這份文件是假設，不是定案。** 六步框架與底下的 skill 清單目前只有三步有實作，
> 而這三步的職責在 2026-09-02 這次重規劃中被重新切過。它應該被下一批真實產物
> 驗證或推翻，而不是被照著蓋完。

---

## 定位

**完整自建。** 各步驟由 ai-bdd 自己提供，不依賴外部 plugin 或 skill。

代價：環境中既有的 `clarify-loop`、`testing-strategy`、`code-review-and-quality`
等 skill 覆蓋部分相同處境，會與 ai-bdd 競爭觸發權。模型會挑一個，而你不會知道
它挑了哪個。緩解方式是在採用 ai-bdd 的專案停用重疊 plugin，並讓所有
`description` 綁定 BDD 脈絡詞彙。

換得的好處：鏈上每個環節的輸入輸出格式**由本專案自行定義**，因此可以真正閉合。
不是功能更多，是介面一致。

---

## 三條紀律

這是 2026-09-02 重規劃的核心。六步之所以站得住，是因為前三步**各自只做一件事**：

| 步驟 | 只做 | 不做 |
| --- | --- | --- |
| **CLARIFY** | 問問題，把問題問到收斂 | 不寫規格。連 story 的定版句子都不寫 |
| **SPEC** | 綜合已經有答案的東西 | **不訪談**。憑空長出來的內容是缺陷 |
| **PLAN** | 切 tracer bullet | 不設計。API、schema、seam 在 SPEC 就定完了 |

三條紀律互相支撐：SPEC 之所以能「不訪談」，是因為 CLARIFY 已經把決定做完；
PLAN 之所以能「不設計」，是因為 SPEC 已經把技術決定寫下來。任何一條破了，
下一條就守不住。

**CLARIFY 的對象不限。** 產品方的業務問題、開發者的技術問題，都是問題——
差別只在問的時機（見下方「三個 pass」）。把技術澄清排除在 CLARIFY 之外，
就等於逼 SPEC 破戒。

紀律的出處：Matt Pocock 的 `to-spec`——「The spec is a record of decisions
already made, not a place where new ones get made.」

---

## 六步流程

| 步驟 | 這一步在問什麼 | 作法 | 產物 |
| --- | --- | --- | --- |
| **CLARIFY** | 還有哪些不知道？問到不能再問為止 | 三個 pass：廣度（PRD 拆解／角色／切 story／分批）→ 深度（每則 story 的業務邊界）→ 技術（打哪一層、動哪些模組） | `brief.md`、`actor.md`、`glossary.md`、`questions/` ＋ 就緒判定 |
| **SPEC** | 這些答案的可執行規格長怎樣？ | 綜合已答問題 → Gherkin ＋ 決定紀錄 | `.feature`（可執行）＋ `spec.md`（Gherkin 表達不了的決定）＋ `domain-model.md` |
| **PLAN** | 怎麼切成一次做得完的步驟？ | tracer bullet 垂直切片 ＋ blocking edges | `plan.md` |
| **IMPLEMENT** | 讓這張票從 red 到 green | step defs ＋ outside-in 雙迴圈 | 程式碼 |
| **VERIFY** | 綠了嗎？紅的是 bug 還是規格過期？ | 跑情境 ＋ 失敗判讀 | 測試結果 ＋ 判讀結論 |
| **REVIEW** | 這變更真的滿足它宣稱的情境嗎？鏈斷了嗎？ | 對照情境驗收 ＋ 鏈稽核 | 審查結論 ＋ 稽核報告 |

框架相依度：CLARIFY／SPEC／PLAN 與框架無關，可先做完並實際使用；
IMPLEMENT／VERIFY 需要綁定技術棧；REVIEW 部分需要。

---

## 產物佈局

一步一個產物槽，`ls` 就看得出這一批走到哪：

```
docs/bdd/
├── brief.md            CLARIFY  PRD 拆解、批次清單與順序、切分依據、隱含假設
├── actor.md            CLARIFY  角色
├── glossary.md         CLARIFY  ubiquitous language（只有詞彙，不含規則）
├── domain-model.md     SPEC     聚合、不變條件
├── questions/          CLARIFY  跨批次的問題（影響分批的）
└── <slice-slug>/                          一批 ＝ 一次可交付的價值
    ├── questions/      CLARIFY  這一批的問題＋答案
    ├── spec.md         SPEC     Gherkin 表達不了的決定
    └── plan.md         PLAN     tracer bullet ＋ blocking edges

features/<story-slug>.feature    SPEC  可執行規格，一則 story 一個檔
```

**根層只留跨批次、且活得比 spec 久的東西。** 一批做完，那個資料夾可以整包封存；
根層的四份繼續長。這條線的判準是 `to-spec` 自己說的：spec 是拋棄式的快照，
實作一開始就會過期；domain model 與詞彙表要活得比它久，塞進 spec 等於陪葬。

`.feature` 是唯一不進 `docs/bdd/` 的產物——它同時是規格與可執行測試，
**位置由測試框架決定**。鏈的完整性靠 tag 與 slug 維繫，不是靠同目錄。

**`glossary.md` 只放詞彙，不放規則。** 它現在（實測場）同時裝了詞彙與
`Shared N` 跨 story 規則，而那些規則是 CLARIFY 抽規則時長出來的——那一段拿掉
之後，規則沒有理由留在詞彙表裡。新的歸屬：跨 story 的**行為**規則寫成
`.feature` 的 `Rule:`，跨 story 的**不變條件**寫進 `domain-model.md`。
詞彙表回到只做一件事：這個字在這個專案裡是什麼意思。

### `spec.md` 的節

| 節 | 內容 | 來源 |
| --- | --- | --- |
| Problem Statement／Solution | 從使用者角度描述問題與解法 | `brief.md` |
| User Stories | 這一批的 story，定版句子 | CLARIFY Pass 1 切出來的 |
| Example Mapping | Rule／Example 的**定版編號**與敘述 | 已答的問題 |
| Implementation Decisions | 動哪些模組、介面長怎樣、API 契約、schema、架構決定 | 技術問題的答案 |
| Testing Decisions | **seam** ＋ 每個 Example 的測試層級 ＋ 既有測試的 prior art | 技術問題的答案 |
| Risks | 規格沒說但實作一定撞到的（併發、冪等、交易邊界） | 從規則推出來 |
| Out of Scope | 明確拒絕的 ＋「留給 IMPLEMENT」的技術決定 | 被否決的選項 |

**一批一份，不是一則 story 一份。** 這不是新決定：`bdd-plan` 的
`references/plan-format.md` 記錄過「每 feature 一份 ＋ 共用檔」試過並且壞掉——
共用檔把實質吸走，各檔退化成「見共用檔」；而且節號互撞，部分對齊比完全不對齊
更危險。API 契約與 schema 天生跨 story，只有並排才看得見衝突。

**Example Mapping 與 `.feature` 的重複是刻意的——重複就是那個檢查。**
`check_spec.py` 的雙向覆蓋稽核（漏做／**發明**）靠比對兩份獨立的表述工作：
map 有而 `.feature` 沒有 ＝ 漏做，`.feature` 有而 map 沒有 ＝ 憑空發明。
把它們併成一份，這個檢查就變成拿檔案跟自己比，恆真。而「發明」正是最沒有人
會懷疑的那種錯——漏一條會被數字抓到，多一條看起來只是很完整。

這是本專案「同一個數字不要兩個地方」原則的例外，而例外要說得出理由：
那條原則防的是**同一個事實**被抄兩份然後漂移；這裡是**兩種不同的表述**
（結構化的規則清單 vs 可執行的場景），漂移本身就是要被偵測的訊號。

**MUST NOT 寫具體檔案路徑與程式碼片段。** 它們過期得比什麼都快。
例外：prototype 產出的、比散文更精確地編碼了某個決定的片段（狀態機、reducer、
schema、型別形狀），可以內嵌並註明來自 prototype。

### `plan.md`

一張票一節，依賴序排列，blockers 在前。每張票要有：

- **交付什麼行為**——從使用者角度，不是分層清單
- **Blocked by**——哪幾張票要先完成，或「無，可立即開始」
- **驗收條件**——每一條都要能在起點 commit 上**失敗**

prefactoring（「make the change easy, then make the easy change」）排最前面。

---

## 三個新概念

現在完全沒有、這次要加進來的：

### 1. Slice ——「一批」

`spec.md` 一批一份，所以「一批」多大是個真的決定。判準：
**這批做完，能不能對使用者交付一次完整的價值。**

實測場的數字給了上限感：6 則 story、69 個場景 → `plan.md` 36.5K。
那已經在可讀性邊緣。`to-tickets` 的文件記錄過超過之後的具體災難：
spec 大到 tracker 讀不回來，agent 燒 tool call 在重抓片段，永遠讀不到結尾。

### 2. Seam ——「驗收測試打在哪一層」

HTTP handler？application service？domain？**這是 BDD 決定 step definition
長相的那個決定**，而現在的 `plan.md` §4 只問「要不要資料庫」「內迴圈是什麼」，
從沒問過外迴圈的高度。

規則（抄 `to-spec`）：優先用既有 seam、取最高的那一層、
**整個變更的理想數量是一個**。

seam 寫進 `spec.md` 的意義是它會往下傳：IMPLEMENT 照著打，
REVIEW 可以把「沒人同意過的 seam」當成 finding 抓出來。這個綁定是間接的
——它穿過 `spec.md`——所以那份文件寫不寫這一節，決定了後面兩步抓不抓得到。

### 3. 一張票 ＝ 一個 fresh context window

現在的「實作順序」有解鎖關係，但沒有尺規，所以切多細沒有判準。
加上這把尺之後，切分有了上下界：**下界**是票要能自己 demo（垂直切片，
不是「先做完 schema」），**上界**是塞得進一個沒看過 spec 的 session。

垂直切片這條是最常被破的。`to-tickets` 記錄過一個案例：26 張票依層切
（corpus／producer／aggregator／selector），每張票平均跑 20 次 agent，
其中四分之三是返工——他們自己的事後檢討把每一類失敗都追回到水平切分。

**例外：wide refactor。** 一個機械性改動（改欄位名、換共用型別）的影響面
擴散到全 codebase 時，沒有任何垂直切片能綠。這種改 expand–contract：
先加新形式（不破壞）→ 分批遷移呼叫端（每批一張票，CI 保持綠）→
最後刪舊形式（被全部遷移票 block）。

---

## 步驟 → skills 對照

**檔案結構是扁平的**（`skills/<name>/SKILL.md`），層級只存在於下面這張表。

| 步驟 | skill | 職責 | 狀態 ／ 這次要做什麼 |
| --- | --- | --- | --- |
| CLARIFY | `clarify-loop` | 多輪把問題問到收斂 | **已實作，不動**。地位從配角變主角——它就是 CLARIFY 的主體 |
| | `story-splitting` | 沿規則切；九種模式 | **已實作**。觸發訊號要改（見下） |
| | `bdd-clarify` | 三個 pass 的總入口 ＋ 就緒判定 | **已實作，要改（中大）**：Pass 1 加「分批」；Pass 2 砍成純問答、拿掉 `example-mapping.md` 這個產物；新增 Pass 3 技術澄清 |
| SPEC | `bdd-spec` | 答案 → `.feature` ＋ `spec.md` | **已實作，要擴張（大）**：新增 `spec.md` 七節、seam 決定、`domain-model.md` 維護；`check_spec.py` 的覆蓋來源從 `example-mapping.md` 改成 `spec.md` |
| | `bdd-spec-review` | 反命令式、conjunction step、情境爆炸稽核 | 未實作 |
| PLAN | `bdd-plan` | `.feature` → tracer bullet ＋ blocking edges | **已實作，幾乎重寫（大）**：現在的 §1–5（API／domain／schema／測試分層／風險）全部搬進 `spec.md` |
| IMPLEMENT | `bdd-implement` | 實作流程總入口 | 未實作 |
| | `bdd-implement-step-definitions` | 比對既有 step 再產新的 | 未實作 |
| | `bdd-implement-outside-in` | 雙迴圈 | 未實作 |
| VERIFY | `bdd-verify` | 跑情境並回報 | 未實作 |
| | `bdd-verify-failure-triage` | 真 bug vs 規格過期 | 未實作 |
| REVIEW | `bdd-review` | 審查總入口 | 未實作 |
| | `bdd-review-acceptance` | 這 diff 真的滿足它宣稱的情境嗎 | 未實作 |
| | `bdd-review-traceability` | 鏈稽核：孤兒、漂移、未覆蓋 | 未實作 |

`bdd-spec` 加了 `spec.md` 之後會不會太肥、要不要拆成兩個 skill——**先不拆**。
子 skill 是否真的需要，等步驟 skill 跑過真實案例、發現它太長時再拆。
現在拆是猜的。

### 拿掉抽規則之後，兩個訊號怎麼保住

CLARIFY 不再產出 `example-mapping.md`，Example Mapping 原本靠數卡片的兩個訊號
要換來源：

| 訊號 | 原本 | 改成 |
| --- | --- | --- |
| story 太大 → 該切 | 藍卡（規則）太多 | **已答問題數**太多（答案多 ≈ 規則多） |
| 不確定性高 → 未就緒 | 紅卡（問題）太多 | **未答問題數**太多（不變） |

兩個都是 `status.py` 現在就在數的東西。而 `references/map-format.md` 自己寫過
「澄清的進展看已答／待答，不是看規則數」——所以這個換法跟它自己的說法一致。

### 追問覆蓋改成算出來的

十一個追問面向原本記在 `example-mapping.md` 的覆蓋段。那份檔案沒了之後，
覆蓋改成**從問題檔算**：每個問題檔的狀態列宣告它屬於哪則 story、蓋哪個面向。

```
**狀態**: 待答 · **輪次**: 2 · **story**: session-training-volume · **面向**: 邊界
```

`status.py` 聚合出「這則 story 有哪幾個面向一題都沒有」。這符合本專案既有的
「進度是算出來的，不另存一份」原則——存一份彙總等於同一個數字兩個地方，
而它們遲早不一樣。

`n/a`（這個面向不適用）需要一個寫理由的地方：開一個狀態為 `n/a` 的問題檔，
理由寫在本文。否則它跟「懶得問」分不出來。

**待實跑驗證。** 這個改法會動到 `status.py`，而它現在能跑。

---

## 命名與結構的決定

**扁平目錄 ＋ `bdd-` 前綴。** 三個理由：

1. skill 名稱在執行時是**全域扁平**的——`claude plugin details` 列出的是名字，
   目錄層級不存在。分組若不寫進名字就完全看不見。
2. 前綴讓它們在字母排序時自然聚在一起，這是扁平命名空間裡唯一還能保住分組的方式。
3. **`claude plugin validate` 不遞迴。** 巢狀目錄（`skills/<分類>/<名稱>/SKILL.md`）
   在 runtime 有效，但 validate 掃不到，等於那些 skill 永遠不會被檢查。
   扁平換來的是每個 skill 都在 CI 的守備範圍內。

通用字眼（`story-splitting`、`failure-triage`）要在 `description` 裡寫明**排除
條件**，否則會在非 BDD 情境被誤觸發。

**`slice` 與 `feature` 是兩個粒度，不共用同一個字。**
`docs/bdd/<slice-slug>/` 裝的是一批（好幾則 story）；`features/*.feature` 是
一則 story 一個檔。用同一個字會製造「部分對齊」——讀的人會假設兩邊指同一件事，
然後在其中幾處讀錯。`feature` 這個字純粹留給 `.feature`。

**三個 plan。** 根層 `PLAN.md`（這個 repo 要建什麼）、步驟 PLAN、產物
`docs/bdd/<slice>/plan.md`。後兩者保留，因為「PLAN ＝ 把 spec 拆成可逐條執行的
任務、寫進 `plan.md`」是既有慣例；歧義靠這份文件開頭那段話消解。

---

## Token 成本

`claude plugin details` 會估算成本。實測參考值（mattpocock-skills，25 個 skill）：

```
Always-on: ~1,620 tok  每個 session 都付
每個 skill：常駐 ~30–160 tok（名稱＋description）
            觸發時 ~20–3.9k tok（本文）
```

推論：**拆得細不貴，description 寫得囉嗦才貴。** description 每次 session 都付費，
本文用到才付。所以 description 壓成一行觸發條件，厚重內容放本文或參考檔。

---

## 可追溯鏈

```
商業目標 → 批次 → Story → 規則 → 例子 → Scenario → Ticket → Step → 程式碼
                                            ↑         ↑        ↑
                                        定版編號   一票一   seam 決定
                                        在 SPEC    session   在 SPEC
                    ↓
   回饋 ← 線上行為 ← 發布 ← 審查 ← 測試執行 ←┘
```

每個箭頭都是一個意圖遺失點。保住這條鏈是 ai-bdd 存在的理由，也是最難被取代的
部分——因為鏈的完整性檢查**幾乎全是機械性的，人只是懶得做**。

這次重規劃改了鏈上兩處：

- **Rule／Example 的編號在 SPEC 定版，不在 CLARIFY。** 原本編號凍在 CLARIFY，
  等於在最不確定的時候做最不可逆的事。改成在寫 Gherkin、真正要用 tag 引用它們
  的那一刻才定版。「MUST NOT 重排編號」這條契約因此變簡單：編號誕生在會用到
  它的地方。
- **IMPLEMENT 的輸入從「一個場景」變成「一張票」。** 票是垂直切片，
  可能蓋到多個場景；票自己可 demo，單一場景不一定。

### 狀態 tag

`.feature` 在 `Feature:` 上掛恰好一個狀態 tag（`@draft` `@ready` `@wip`
`@review` `@done`），記錄它走到六步的哪一站。

完整詞彙表與三個設計決定在 `skills/bdd-spec/references/state-tags.md`——
**它是 skill 之間的契約，所以跟著 skill 走，不放在這份路線圖裡**。
別人安裝這個 plugin 時拿不到 repo 的文件，契約寫在這裡就等於斷了。

依賴方向：repo 的文件可以指向 skill，skill 不可以指向 repo 的文件。

---

## 進度

- [x] `bdd-clarify`
- [x] `clarify-loop`
- [x] `story-splitting`
- [x] `lab/go/skeleton` — Go ＋ godog 實測場（骨架綠，零業務）
- [x] 拿健身追蹤需求跑一次 `bdd-clarify`，依實跑結果修正格式
- [x] `bdd-spec` — 依 11 份真實 `.feature` 反覆修正，封閉文法把
      76 個場景的步驟樣板從 208 降到 41
- [x] `bdd-plan` — 拿實測場的 5 個 `.feature`（69 個場景）跑過，依執行者回報的
      三處摩擦修正
- [x] 2026-09-02 重規劃：三條紀律、slice、seam、tracer bullet 尺規
- [x] `bdd-clarify` 改版（Pass 3、分批、拿掉 example-mapping）
- [x] `bdd-spec` 擴張（`spec.md`、seam、`domain-model.md`）
- [x] `bdd-plan` 重寫（tracer bullet ＋ blocking edges）
- [x] 實測場整批重跑一次，驗證新的三步分工
- [ ] `bdd-implement`

順序刻意是「做一個 → 實跑 → 再做下一個」。**產物格式就是介面**：SPEC 讀已答的
問題、PLAN 讀 `.feature` 與 `spec.md`、REVIEW 稽核整條鏈，下游全都依賴上游實際
吐出什麼。在紙上設計那個格式是猜的。

這次重規劃違反了上面那句話——它是在紙上改的。所以「實測場整批重跑」那一項不能跳。

---

## 實測場

`lab/go/skeleton/` 是六步流程的 dogfooding 對象：Go ＋ godog ＋ 分層
架構（cmd／domain／application／infrastructure／interfaces）＋ 三層測試
（unit／integration／acceptance）＋ `api/openapi.yaml`。

它的規則是**在 CLARIFY 跑完前不寫任何業務程式碼**。題目選「訓練總容量計算」
而非「記錄訓練」，因為前者才有規則與模糊點——自體重動作、單邊動作乘不乘二、
熱身組算不算、dropset 是幾組。詳見該目錄的 README，那裡也寫明了
`bdd-clarify` 的及格標準：**沒問到自體重與單邊動作，就算失敗**。

**現存產物與新佈局不一致。** `lab/go/skeleton/docs/bdd/` 底下的六個目錄
（`log-a-workout`、`session-training-volume`⋯⋯）是**六則 story，不是六個 slice**。
照新佈局它們會併成一到兩個 slice，六則 story 變成 `spec.md` 裡的小節、
`features/` 裡的六個 `.feature`。這個遷移就是驗證新分工的那次實跑。

> `lab/` 不是 plugin 的慣例目錄，不會被載入，對 plugin 行為零影響。

---

## 未決事項

1. ~~**產物存放路徑。**~~ 已定：文件型產物一律寫入使用端 repo 的 `docs/bdd/`，
   一個目錄裝完，刪掉即乾淨。skill 不得寫入其他位置、不得修改專案既有檔案。

2. ~~**`.feature` 的位置。**~~ 已定：留在測試框架的慣例位置（Cucumber 系是專案根
   的 `features/`），是唯一不進 `docs/bdd/` 的產物——它的位置由 runner 決定。

3. ~~**SPEC 與 PLAN 的分界。**~~ 已定（2026-09-02）：技術設計歸 SPEC，
   拆解歸 PLAN。見「三條紀律」。

4. **一批到底多大。** 判準已定（一次可交付的價值），已有第一個數字：實測場
   `workout-tracking` 這一批是 6 則 story、39 條規則、93 個場景（`@example`
   標籤數）、`spec.md` 649 行，切成 9 張票。這只是一個數字，不是上限——第二個
   數量級不同的實測場出現之前，先當經驗值用，不當硬性天花板。

5. ~~**技術問題的追問面向。**~~ 已定：Pass 3 的技術面向清單見
   `skills/bdd-clarify/references/technical-probes.md`，目前有五個面向
   （seam、模組邊界、介面與型別契約、排序契約／決定性、既有資產／測試慣例）。
   這五個面向全部從 `workout-tracking` 這一次實測場的硬跑反推出來，只跑過
   一個技術棧（Go）、一個專案——換一個技術棧不同的專案重跑，如果冒出新的、
   真的卡住的技術缺口要加進去，這份清單不是天花板。

6. **IMPLEMENT／VERIFY 的技術棧。** 未定（Python behave/pytest-bdd、
   JS Cucumber.js、Go godog…）。前三步不受影響。

7. **License。** repo 在 `vivotek/` 底下，`plugin.json` 目前不含 `license` 欄位。

8. **Pass 2 的十一個業務追問面向需要重新檢視。** 這次實測場一輪跑下來冒出
   16 個缺口，其中 11 個是業務缺口而非技術缺口，裡面至少兩個對不上現有十一
   個面向裡的任何一個——包括「這個能力到底在不在這批範圍內」這種連要不要問
   都還沒有位置放的問題。這裡只能誠實記下這個發現：詳細的證據在執行過程中
   工作目錄被清空，沒能留下來，所以這一項記的是結論，不是分析過程本身。第二
   個實測場跑的時候，要重新對照這十一個面向，看看夠不夠、還是要加。

9. **`audit_skill.py` 只查得到孤兒檔案，查不到反過來的情況。** S5 規則抓得到
   子目錄裡一個沒人指到的檔案，但抓不到 `SKILL.md` 指到、實際上不存在的
   章節。這一批 `bdd-plan` 整份重寫，就活生生撞過一次同樣形狀的缺陷——只是
   換了個樣子：舊版 `SKILL.md` 把缺口那一節標成「第 6 節」，但缺口那一節
   實際編號是第 7 節；這次重寫把那個引用整段拿掉了，缺陷是在人工複查時才被
   抓到的，不是任何一支腳本抓到的。要不要幫 `audit_skill.py` 加一條反向規則
   （`SKILL.md` 引用的章節號／檔名是否真的存在），是懸而未決的。

10. **這條鏈最容易出的錯是靜默地錯，不是壞掉。** 這一輪實跑撞到四個例子：一個
    dashboard 回傳 exit 0、顯示「已就緒」，但其實什麼都沒找到；一張覆蓋率表，
    只要每個問題都歸錯 story，看起來一樣完整可信；重新編號過的 Rule／Example
    標籤不會噴任何錯誤，只會讓答案悄悄變成錯的；一支腳本靠比對內文裡剛好加粗
    的文字去解析一行已經壞掉的狀態列，照樣跑得動。四個都是 exit 0、看起來健康。
    這條鏈要不要為這一類錯誤加一層系統性的檢查，還沒有答案。
