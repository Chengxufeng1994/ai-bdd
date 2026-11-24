---
name: formulation
description: >
  把已對焦的 `spec.md` 翻成可執行的 Gherkin `.feature`——FR 變成 `Rule:` 區塊、
  AC 變成 `Example:`、編號變成 tag 帶下去，步驟套一組封閉文法讓 step definition
  能重用，最後跑一次雙向覆蓋稽核。不訪談、不發明例子、不改 `spec.md`。
  BDD 三實踐裡的 Formulation；由 `plan` 在定了 seam 之後、設計 API 與 schema 之前委派。
  觸發詞：「寫成 feature 檔」「轉成 Gherkin」「產出驗收情境」「把例子寫成場景」
  「寫 .feature」「產生驗收測試規格」「這些規則的 scenario 長怎樣」「跟 PM 對答案」。
  English: turn spec.md into Gherkin, write the feature file, generate
  acceptance scenarios, convert rules and examples into scenarios,
  formulation step of BDD.
disable-model-invocation: true
---

# FORMULATION — 把 `spec.md` 翻成可執行的規格

| | |
| --- | --- |
| **職責** | 把已對焦的 `spec.md` 翻成 Gherkin `.feature`：FR 變 `Rule:`、AC 變 `Example:`、編號變 tag，步驟套一組封閉文法 |
| **讀** | `specs/ongoing/<date>-<feature>/spec.md`（`已對焦`）、`docs/CONTEXT.md`（有的話）、既有的 `.feature`、`specs/ongoing/<date>-<feature>/acceptance/*.feature`（只取用語） |
| **寫** | `features/<story-slug>.feature`（測試框架已有位置就跟隨）；覆蓋表只回報、不落地 |
| **上游** | `plan`（〈設計〉定了 seam 之後 DELEGATE；也可以使用者單獨叫 `/formulation`） |
| **下游** | 回 `plan`，它接著設計 Domain Model、schema、API 契約；覆蓋缺口由 `plan` 記進 `## 缺口` 交回 `spec` |
| **不做** | 不訪談、不發明例子、不寫 `spec.md`、不寫 step definition、不讀 `plan.md` |

**不發明例子。** `Rule:`、`Example:` 與 `@rule-`／`@example-` tag 只從 `spec.md` 來；覺得少了
一個邊界，那是 DISCOVERY 的發現，回 `spec`，不在這裡順手補。覆蓋表要產，但它是回報給人的，
不是落地的檔案——缺口同時要就地註解進 `.feature`。

## 使用時機

- `spec.md` 已對焦，要把 FR 與 AC 變成 `.feature`
- 要跟 PM 或需求方對答案——攤開規則與例子請他點頭。`.feature` 就是對答案用的那一份，不另做第二
  份給人讀的版本——那是 Cucumber 的前提。PM 讀不動就修 `references/step-grammar.md`
- `plan` 定了 seam、要設計 API 與 schema 之前，需要先有行為的場景（`plan` DELEGATE 進來）

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| `spec.md` 不存在 | `discovery`（它寫 `interview-log.md`），再 `spec`（它從紀錄寫 `spec.md`） |
| `spec.md` 狀態不是 `已對焦` | `spec`——它會先交 `example-mapping` 判就緒 |
| 要決定 seam（情境跑在哪一層測試） | `plan` |
| 要定詞彙或切 story | `spec`（詞彙寫在 `spec.md` 的 `## Domain Concepts`；切 story 用 `story-splitting`） |
| 要決定實作順序、切票 | `plan` |
| 要稽核既有 `.feature` 寫得好不好 | 目前沒有 skill 接手（`spec-review` 未實作，見 `PLAN.md`）；在那之前跑本 skill 的 `scripts/check_spec.py`（十八項裡五項稽核 `.feature`），對照 `examples/anti-patterns.md` |

產出是 `.feature`，**不是** `spec.md`、`plan.md`、`openapi.yaml`、migration 或任何可執行的檔案。

## 參考檔案

| 讀哪一份 | 什麼時候 |
| --- | --- |
| `references/step-grammar.md` | 四種步驟形狀的完整理由與邊界；某條步驟規則想不通為什麼，或準備破例之前 |
| `references/rule-taxonomy.md` | 四象限分類與覆蓋盤點；排規則之前，以及寫完做覆蓋盤點 |
| `references/rule-example-conflict.md` | 規則與它自己的例子打架時（步驟 2） |
| `references/state-tags.md` | 狀態 tag 的完整詞彙表 |
| `references/coverage-report.md` | 覆蓋表的格式 |
| `references/artifact-location.md` | `.feature` 該放哪裡 |
| `examples/anti-patterns.md` | 寫完自檢；或樣板數看起來不對 |
| `scripts/check_spec.py` | 稽核 `spec.md` 自己的形狀；有 `.feature` 之後再加上兩者的一致性 |

十一份逐案範例：

| 讀哪一份 | 什麼時候 |
| --- | --- |
| `examples/1-video-progress.md` | 第一次用。命令 · 狀態機，最完整的一份 |
| `examples/2-submit-assignment.md` | 型別檢查；Example 層播種 |
| `examples/3-deliver-course.md` | 副作用與序列；兩個 `When` 何時對 |
| `examples/4-query-progress.md` | 寫查詢時。第四種形狀不一樣 |
| `examples/5-query-course-list.md` | 查詢回傳多筆時。順序與「有沒有多餘的」 |
| `examples/6-query-product-list.md` | 巢狀結構、參數篩選、空象限怎麼註明 |
| `examples/7-create-order.md` | 最複雜的命令。算出來的預期值、選用參數 |
| `examples/8-query-order.md` | 一個情境被斷言切開的反例；錯誤訊息的粒度 |
| `examples/9-pay-order.md` | 副作用最重的一份；七個 Example 如何收斂成兩個 |
| `examples/10-cancel-order.md` | 狀態守衛。`Given` 用狀態還是用動作鋪陳 |
| `examples/11-query-user-roles.md` | 空清單、軟刪除；弱動詞如何逼出多餘的例子 |

## 前置確認（Ask user）

不必問。`.feature` 的位置偵測得到（見「產物隔離」），要跑哪些 story 有預設值。

MUST NOT: 在這裡問問題。**不訪談**——憑空長出來的內容是缺陷。`spec.md` 的
`## Open Questions` 還有待答時，那些對應的規則就寫不出來，不要順手決定掉。

NEVER: `spec.md` 不存在、或狀態不是 `已對焦` 時，憑需求描述直接寫 `.feature`。那會
產出一份沒有人同意過的規格，而它看起來跟真的一模一樣。唯一的例外是下面的 `@draft`：
`spec.md` 存在、狀態 `待對焦`，而且使用者明說要先寫。沒有 `spec.md` 就先跑
`discovery`、再 `spec`。

---

## 流程

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| 1 | 讀四樣東西，只做 `已對焦` 的 feature | 不問 |
| 2 | FR → `Rule:`、AC → `Example:`、編號 → tag；歸象限、看三種結果、取場景名 | 不問；規則與例子打架時回報 |
| 3 | 每個步驟套四種形狀之一 | 不問 |
| 4 | 句子只差資料的例子收成 `Scenario Outline` | 不問 |
| 5 | 一則 story 一個 `.feature` | 不問 |
| 6 | 跑 `check_spec.py`、產覆蓋表 | 不問；缺口回報 |

### 1. 讀四樣東西

產出：無，只讀——這次要做哪幾個 `已對焦` 的 feature、可用的詞彙、既有的步驟樣板

| 讀什麼 | 為了什麼 |
| --- | --- |
| `specs/ongoing/<date>-<feature>/spec.md` | **唯一的規範來源**——`## User Stories` 底下的 FR 與 AC 是這一步要翻譯的全部內容，`## Personas` 是 `Feature:` 敘述裡角色的來源，`## Goals / Non-Goals` 劃出碰不得的界線 |
| `spec.md` 的 `## Domain Concepts` 與 `docs/CONTEXT.md` | **措辭的來源**——前者是這個 feature 定下的詞，後者是專案既有的詞。同一個概念在不同 `.feature` 裡換一個說法，封閉文法的重用就從第一份開始漂移 |
| 既有的 `.feature` | 已經定下的步驟樣板，能套就不要另造 |
| `specs/ongoing/<date>-<feature>/acceptance/*.feature` | **用語的參考**——挑步驟句型、寫 `Feature:` 敘述、取參數名時，優先用業務方在驗收旅程裡讀過的詞，讓兩層講同一套話 |

MUST NOT: 從驗收旅程新增 `spec.md` 沒有的例子、場景或 `Then`。旅程刻意少、句子不受
封閉文法約束——當成來源的話，漏掉的 AC 與自由發揮的句子會流進可執行的那一層。
`Rule:`、`Example:` 與 `@rule-`／`@example-` tag 只從 `spec.md` 來。旅程的用語跟
`spec.md` 打架（同一件事兩個名字）時以 `spec.md` 的 `## Domain Concepts` 為準、再來是 `docs/CONTEXT.md`，
回報給使用者，不改旅程。

寫新樣板之前先 grep 一次既有的 `.feature`——第三項最容易被跳過、代價最貴：封閉文法的價值全在重用，另造一個形狀等於白做。

IMPORTANT: `spec.md` 由 `spec` 維護，本 skill **不得寫入**。要補一條
規則、改一個例子，都回 `spec`；問答裡沒有的，它會開一題補問題、請使用者跑 `discovery`。本 skill 通常是
`plan` 在定了 seam 之後、設計 API 與 schema 之前 DELEGATE 的（為什麼排在設計之前，見 `plan`〈設計〉）。

MUST: 只做 `spec.md` `## Document Overview` 的 `狀態` 已經是 `已對焦` 的 feature；`澄清中`、`待對焦` 的整份跳過（例外只有下一段的 `@draft`），
並明講該回哪個 skill——還沒定案的東西寫成場景，等於**把不確定性從顯眼的問題檔搬進一份看起來已完成的規格**，之後沒有人會回頭質疑它。

使用者明說「未就緒的也先寫」、而 `spec.md` 狀態是 `待對焦` 時就寫，標成 `@draft`——見下方
「狀態 tag」。`澄清中` 的不寫：那時連 FR 都還沒有。

### 2. 規則對應 Rule，例子對應 Example

產出：每條 FR 一個 `Rule:`、每個 AC 至少一個 `Example:`（草稿），掛 `@rule-<n>`／`@example-<n>.<m>`；每條規則的象限與三種結果

```
spec.md 的 `## User Stories` 段      .feature
─────────────────────────────────────────────────────────────
這則 story 的切法說明                →    Feature: 與其下的敘述
**FR-2**  影片進度不可回退            →    Rule: 影片進度必須單調遞增
- **AC-2.2** 進度 70% 想改 60%        →    @rule-2 @example-2.2
                                          Example: 進度回退時操作失敗
```

| `spec.md`（PRD 詞彙） | `.feature`（Gherkin 詞彙） |
| --- | --- |
| `FR-<n>` | `@rule-<n>` |
| `AC-<n>.<m>` | `@example-<n>.<m>` |

- tag 不跟著 `spec.md` 改名——`@example` 是 Gherkin 的字、FR／AC 是 PRD 的字，接縫寫明對應就夠了
- 用 `Example:` 不用 `Scenario:`——同一個東西在鏈上換名字，讀的人就得自己對應

MUST: 每個 Example 至少對應一個場景，且**不新增 `spec.md` 裡沒有的例子**。

覺得少了一個邊界時，那是 DISCOVERY 的發現：回去補一題，讓它經過同樣的討論——在這裡順手加一個場景，等於一個人決定了驗收條件。

**唯一的例外是中性的填充值。** 最常見的是識別碼：`spec.md` 寫「物件導向基礎」，`.feature` 需要
一個 ID 才能穩定指涉。編一個不算發明驗收條件——**ID 是定址，不是行為**。同一類還有：可執行場景
需要、但 AC 沒寫的值（借出那一刻的時分、「R-001 有效」的讀者證期限、刷書的時刻），只要它不影響斷言。
判準只有一條：把這個值換成別的，有沒有任何一條斷言的真假會變？會變就不是填充值，是行為——
那個例子跳過、就地註解、回報要回 `spec` 在 `spec.md` 的 `## Open Questions` 開一題補問題。補了填充值的場景，在 `.feature` 就地留一行
註解說哪個值是補的。

SHOULD: 用 UUID，不用 1、2、3——序號**可以猜**，「查別人的訂單回權限不足、查不存在的回訂單不存在」在序號下就是一個枚舉工具；ID 的形式是規格的一部分，不是排版偏好。

**編號用 tag 帶下去。** `Rule:` 區塊帶規則的文字，編號在 tag 裡：

```gherkin
@rule-2 @example-2.2
Example: 進度回退時操作失敗
```

兩個 tag 都要——`@example-2.1` 讓你機械稽核哪個例子還沒有場景，`@rule-2` 讓你一次跑完一條規則的所有場景（含 Outline 展開的每一列）。

MUST NOT: 事後重排 `spec.md` 的 FR 與 AC 編號。**編號在 `spec` 誕生**，本 skill 原封不動搬進 tag——下游回指的是 `@example-2.2`
這個 tag，重排等於讓那些引用**靜默**指向別的東西：不會報錯，只會對錯。

#### 規則歸到四個象限

| | 操作**前** | 操作**後** |
| --- | --- | --- |
| **狀態** | 前置（狀態）：什麼情況下不該被允許 | 後置（狀態）：世界變成什麼樣，含副作用 |
| **輸入輸出** | 前置（參數）：邊界、格式、範圍 | 後置（回應）：回傳的形狀 |

命名帶上象限：`Rule: 前置（狀態）- <主詞>必須<單一條件>`。**單一條件**——「必須已購買**且**已訂閱」是兩條規則，拆開之後
失敗的例子才指得出是哪個條件沒滿足；判準：兩半的**錯誤訊息不同**，就是兩條規則。

象限是**分類，不是規則的名字**：`前置（參數）` 底下常常有好幾條。

空掉的象限是規格的缺口。MUST NOT: 自己補一條規則填它——那是 DISCOVERY 的工作。
MUST: 在 `.feature` 裡留一行註解說明它是**合理地空**還是**沒問過**，並寫進覆蓋表——報告會隨對話消失，檔案會留著。

每個象限逼你想什麼、覆蓋盤點怎麼寫 → `references/rule-taxonomy.md`。

#### 每條規則要看三種結果

象限分類問的是「有哪幾類規則」，這一條問的是**每條規則被走過哪幾種結果**：

| | 問什麼 |
| --- | --- |
| **成功** | 這條規則被滿足時會怎樣 |
| **邊界** | 剛好在界線上、空的、第一次、最大值 |
| **失敗** | 這條規則被違反時會怎樣，錯誤訊息是什麼 |

MUST: 每條規則的三種結果，在覆蓋表裡各標**有／無／不適用**，
`不適用` 要寫理由。

不是每條規則三種都要有（查詢常沒有失敗，例如 `workout-history` 的分頁規則；常數規則沒有邊界），但**「沒有」與「沒想到」在 `.feature` 裡長得一樣**——
這一條把它們分開。實測：一批 55 條規則只有 5 條同時有成功與失敗的例子，一個 `@command` 的 feature 四條規則全部沒有失敗場景——不是那些 story 不會失敗，是沒有人問過它怎麼失敗。

#### 場景名字要能定位

**句型**：`<情境條件>時<預期結果>`（失敗）或 `<操作>後<預期結果>`（成功）。
規則就寫在上一行，場景再抄一次只是佔位置；測試紅了時讀到的是這一行，它要能單獨說明壞了什麼。

| | |
| --- | --- |
| ✓ | `Example: 已購買旅程的用戶再次購買時操作失敗` |
| ✗ | `Example: 使用折扣券` —— 讀不出結果 |
| ✗ | `Example: 影片進度必須單調遞增` —— 這是規則，不是案例 |

#### 規則與它自己的例子打架時

把規則套回它的例子，有時候會發現**算不出來**——例子裡的數字在它自己那條規則下不可能出現。

| 情況 | 怎麼做 |
| --- | --- |
| 例子的錯**可以從規則推導出來**（算術錯誤） | 規則為準，照規則算出來的數字寫場景；就地註解算式與差異，回報建議回 `spec` |
| 規則本身含糊，兩種讀法都說得通（規格衝突） | 不要挑：跳過那個例子、就地註解、回報要回 `spec` 在 `spec.md` 的 `## Open Questions` 開一題補問題 |

撞到時先讀 `references/rule-example-conflict.md`——規則、理由與判準的界線都在那裡。

### 3. 步驟文法 —— 四種形狀

產出：每個步驟都套進下表四種形狀之一（套不進的留一行為什麼）

`.feature` 寫給**下一個要實作 step definition 的人或模型**讀，目標是**規律**，不是看得懂——實測散文式（封閉文法之前 `spec` 的產出）76 個場景長出 208 個樣板、
172 個只用一次（step def 比場景還多，遲早分岔）；封閉文法（`examples/1-video-progress.md`）7 個場景 10 個樣板、3 個只用一次，再加二十個場景樣板表也幾乎不變長（`check_spec.py` 量的）。

每個步驟套下列四種形狀之一：

| # | 形狀 | 位置 | 樣板 |
| --- | --- | --- | --- |
| 1 | **播種** | `Background`；只有一個例子要用時放它自己的 `Given` | `系統中有以下<實體>：` ＋ 資料表 |
| 2 | **前置狀態** | `Given` | `<實體> 的 <屬性> 為 <值>` |
| 3 | **動作** | `When`，預設一個 | `<行為者> <動作> <目標> <參數>` |
| 4 | **結果** | `Then` | `操作成功` ／ `操作失敗`，再接斷言——命令接**狀態**，查詢接**結果** |

```gherkin
Feature: 增加影片進度

  Background:
    Given 系統中有以下課程：
      | 課程 ID | 名稱         | 類型     | 獎勵.經驗值 |
      | 1       | 物件導向基礎 | 影片課程 | 300         |
    And 系統中有以下用戶：
      | 名稱  | 等級 | 經驗值 |
      | Alice | 1    | 0      |

  Rule: 影片進度必須單調遞增

    @rule-2 @example-2.1
    Example: 進度不可倒退
      Given 用戶 "Alice" 在課程 1 的進度為 70%，狀態為 "進行中"
      When 用戶 "Alice" 更新課程 1 的影片進度為 60%
      Then 操作失敗
      And 用戶 "Alice" 在課程 1 的進度應為 70%，狀態應為 "進行中"
```

#### 每一條規則

| | 規則 | 為什麼 |
| --- | --- | --- |
| `MUST` | 清單查詢用 `應恰好為`／`應依序恰好為`，不用 `應包含` | 弱動詞證明不了「沒有多餘的」，還會逼出負向斷言 |
| `MUST` | 清單查詢要有**空集合**的例子 | 沒資料時回 404 還是空清單，是兩種實作 |
| `MUST` | 副作用改到集合時，斷言集合的**完整內容** | 只驗新增的那一項，覆蓋掉舊資料的實作會通過 |
| `SHOULD` | 兩個參數的**合法組合**是一條 `前置（參數）` 規則 | 各示範一組成立的組合，推不出其他組合會怎樣 |
| `MUST` | 算出來的預期值，計算規則寫進 `Rule:` 敘述 | 讀規則要能自己算出同一個數字 |
| `MUST` | 斷言表的每一格是字面值 | `18位數字` 是規格描述不是資料 |
| `SHOULD` | `Given` 用狀態表達；不變式複雜時才用動作鋪陳 | 直接設狀態可能造出系統到不了的狀態 |
| `MUST` | 失敗案例要讓失敗原因唯一 | 否則刪掉待驗的規則，測試照樣綠 |
| `MUST NOT` | 同一個 `Given`＋`When` 不因斷言不同而拆成多個 Example | 「這些欄位同時存在」會沒被驗到 |
| `MUST` | Example 之間不共享狀態、不依賴順序 | godog 並行跑，順序相依＝間歇性失敗 |
| `SHOULD` | 命令的失敗案例要斷言狀態沒變；查詢不必 | 擋住「回錯誤但已寫進資料庫」 |
| `SHOULD` | 一個 `When`；序列產生新行為時才用多個 | 狀態守衛用 `Given`，累積才用序列 |
| `SHOULD` | 選用參數拆成獨立 `And`，不每種組合各造一句 | 否則 step definition 組合爆炸 |
| `SHOULD` | `Given` 只寫與預設不同的部分 | 鋪陳預設狀態＝驗自己的佈置 |
| `SHOULD` | 播種用資料表，不一句一個實體 | 欄位增加時不必改簽章 |
| `SHOULD` | 一個 Example 三到五步 | 步驟愈多，失敗時要排除的可能愈多 |
| `SHOULD NOT` | 不把成敗與狀態擠進同一句 | 複合斷言的失敗訊息說不出壞在哪 |
| `SHOULD NOT` | 不寫按鈕、欄位、API 路徑 | 綁死操作方式，換介面就要重寫 |
| `SHOULD` | 行為者用名稱、其他實體用 ID；字串加引號、數字不加 | 指涉不一致＝同一件事兩支 step definition |
| `SHOULD NOT` | 資料表裡不放 JSON 字串 | 失去欄位名，diff 看不出改了哪個值 |

每一條的理由、範例與邊界 → `references/step-grammar.md`。
套不進四種形狀時先當成設計訊號（通常代表這個 Example 混了兩件事）；
真的要破例就在檔案裡留一行說明為什麼——沒有理由的例外會擴散。

### 4. 什麼時候用 Scenario Outline

產出：句子只差資料的例子收成 `Scenario Outline`，其餘維持獨立的 `Example`

**判準只有一個：這些例子除了資料以外，句子完全一樣嗎？** 一樣 → `Scenario Outline`；不一樣 → 各自獨立的 `Example`，不要硬湊。一條規則掛很多例子時多半是在窮舉一張對照表（型別、狀態、單位），那正是 Outline 的用途；
語意不同的例子塞進同一張 Examples 表，會逼出 `| 課程 | 型別 | 進度 | 預期 | 備註 |` 這種越加越寬、沒人讀得出每列在測什麼的表。

```gherkin
@rule-3 @example-3.1 @example-3.2 @example-3.3
Scenario Outline: 有效範圍內的進度值可以更新
  Given 用戶 "Alice" 在課程 1 的進度為 50%，狀態為 "進行中"
  When 用戶 "Alice" 更新課程 1 的影片進度為 <新進度>%
  Then 操作成功
  And 用戶 "Alice" 在課程 1 的進度應為 <新進度>%，狀態應為 "<預期狀態>"

  Examples:
    | 新進度 | 預期狀態 |
    | 60     | 進行中   |
    | 100    | 已完成   |
```

IMPORTANT: 一個 Outline 涵蓋 Example 3.1 到 3.3 時，**三個 tag 都要掛**，
否則覆蓋稽核會誤報漏做。

### 5. 寫檔

產出：`features/<story-slug>.feature`，一則 story 一個

檔名用 feature slug（`video-progress.feature`），跟 `specs/` 的 slug 一致——那是這條鏈唯一的接縫，換名字就得靠人腦記住對應。

### 6. 稽核，然後才算完成

產出：`check_spec.py` 的輸出、覆蓋表（回報，不落地）、`.feature` 裡的缺口註解

**這一步跑在最後**——它比對 `spec.md` `## User Stories` 段（FR／AC 的定版與 story 涵蓋哪些 FR）與寫好的 `.feature`，兩者都要存在才對得起來——合起來才拼得出「這個 `.feature` 該有哪些 `@example` tag」。

```bash
python3 <skill>/scripts/check_spec.py <專案根>
```

**不能因為 `spec` 跑過同一支腳本而省略**——沒有 `.feature` 時它只跑十三項，需要 `.feature` 的五項（雙向覆蓋、狀態 tag、方言陷阱、
缺口註解、樣板重用率）明講跳過；十八項全開只在這裡。

十八項的完整清單與各自的理由在 `scripts/check_spec.py` 開頭的說明：前五項需要 `.feature`
（上面那五項），其餘十三項只讀 `spec.md` 與 `acceptance/`——FR 有沒有掛 AC、
標題形式、`#### Q` 紅卡、AC 編號對帳、無主的 FR、`← PRD §x`／`← Q-<n>` 來源（`推論` 不算）、
`## 原文` 與它的段落、驗收旅程的 `@covers-`。只讀不寫，退出碼可直接進 CI。

- **發明比漏做更該優先看**——漏一條會被覆蓋數字抓到；憑空多一條看起來很完整，沒有人有理由懷疑它
- 樣板重用率沒有及格線，參考點是步驟 3 的兩組實測（`examples/anti-patterns.md`）；只用一次的樣板數逼近場景數，代表文法沒有真的套上
- MUST: 產出覆蓋表，缺的列出來、寫出為什麼，不要只說「完成」
- MUST: 缺口也**在 `.feature` 裡就地留一行註解**——回報會隨對話消失；半年後看到沒解釋的洞，只能當漏做補上（把還在討論的東西變成驗收條件）或花半小時考古

格式照抄 `references/coverage-report.md`。

---

## 產物

| 檔 | 哪一步寫 | 格式在哪 |
| --- | --- | --- |
| `features/<story-slug>.feature`（測試框架已有位置就跟隨） | 5；再跑時只改對到修訂歷史的場景 | 本節（關鍵字、狀態 tag、型別 tag、方言陷阱）與 `references/step-grammar.md` |

覆蓋表只回報、不落地（格式 `references/coverage-report.md`）。`spec.md` 只讀不寫；`plan.md` 不讀也不寫。

規則與例子的**定版編號在 `spec` 誕生**：`FR-<n>`、`AC-<n>.<m>`。`.feature` 的
`@rule-<n>`／`@example-<n>.<m>` 回指它們——本 skill 只把它們原封不動地搬進 tag。

**關鍵字用英文，名稱與步驟用中文。** 不加 `# language:` 那一行（英文是預設方言）。完整骨架見步驟 3 的範例。

`Feature:` 底下的敘述寫這則 story 的「作為⋯我要⋯以便⋯」：角色抄 `spec.md` 的 `## Personas`，能力與價值從涵蓋的 FR 濃縮（起點是
`## User Stories` 的切法說明）——它是這份規格存在的理由，而讀 `.feature` 的人不會同時開著 `spec.md`。

### 狀態 tag —— 每個檔恰好一個

MUST: 每個 `.feature` 在 `Feature:` 上掛**恰好一個**狀態 tag。本 skill 只寫兩個：

| story 的就緒判定 | 寫入 |
| --- | --- |
| 已就緒 | `@ready` |
| 未就緒（使用者仍要求先寫） | `@draft` ＋ 一段說明它建立在哪些未定的答案上 |

完整詞彙表（含 `@wip` `@review` `@done`，由後面的 skill 寫入）在
`references/state-tags.md`。**不要在這裡抄一份**——同一張表存兩處就會漂移。

### 型別 tag —— `@command` 或 `@query`

MUST: 狀態 tag 旁邊再掛一個型別 tag：`@ready @command` 或 `@ready @query`。

它不只是標籤，**它決定哪些象限該有東西**：命令的重心在前置（狀態／參數）與後置
（狀態），查詢在前置（狀態）與後置（回應）。`@query` 卻長出「後置（狀態）」規則，
代表這個查詢有副作用——那多半是設計問題，不是規格問題。

同時是命令又是查詢的 feature 是兩則 story → 回 `spec` 重切（手法見 `story-splitting`）。

必填而不是「有問題才標」——**少標一個要看得出來**，負向旗標分不出「不需要標」與「忘了標」；第一次真實產出就是兩則未就緒的 `.feature` 沒有任何標記，跟已就緒的長得一模一樣。

`@draft` 的說明段要具體到讀者能自己判斷「這份能不能信」：未就緒的成因是什麼、
恢復要跑哪個 skill。只寫「未就緒」等於把問題換個地方藏。

### 方言陷阱 —— 不要用中文的「規則:」

MUST NOT: 寫 `規則:` 當關鍵字。

zh-TW 方言**沒有**「規則」這個字（`zh-CN`、`ja` 有）。實測 `gherkin/go/v42`：`規則:` 產生 **0 個 Rule 節點**，整行被吞進
`Feature.Description`，**不報錯**、照樣 exit 0，規則與場景的從屬關係安靜消失——而那正是這份規格最主要的結構。這也是關鍵字一律用英文的理由：不必逐個記哪些字有翻譯。

---

## 產物隔離

MUST: 只新增 `.feature`，**不修改專案既有的任何檔案**（含既有的 `.feature`、
`spec.md`、`plan.md`、`docs/CONTEXT.md`）。唯一的例外是下面〈再跑：`spec.md` 改過之後〉——
那時只改那些 `.feature` 裡對到改動的場景。`.feature` 的位置由測試框架決定——有 `.feature`
就跟隨，沒有就用專案根的 `features/`。理由與例外 → `references/artifact-location.md`。

### 再跑：`spec.md` 改過之後

`spec.md` 到過 `已對焦`、`.feature` 已經寫了，之後 `spec` 照新的答案修補了 FR／AC，
需求方又答「就緒」——這時使用者再叫 `/plan`，它在〈設計〉再 DELEGATE 本 skill。照 `spec.md` 的 `### 版本修訂歷史`
新增的列，只動對到的場景：

| 修訂歷史說 | 對 `.feature` 做什麼 |
| --- | --- |
| 新增 FR／AC | 照步驟 2–4 加 `Rule:`／`Example:`，tag 接著編 |
| 改了 FR／AC | 改掛同一個 `@example-<n>.<m>` 的場景，tag 不變 |
| 刪了 FR／AC | 刪掉掛那個 tag 的場景，在原處留一行註解：刪了什麼、修訂歷史哪一列 |
| 改了 FR／AC，但上次那個 AC 被跳過（只有一行缺口註解，沒有場景） | 當成新增：刪掉那行缺口註解，照步驟 2–4 寫場景，tag 用原本的 `@example-<n>.<m>` |

改過場景的 `.feature`，狀態 tag 改回 `@ready`——改過的場景還沒有人實作，舊的 `@wip`／
`@review` 已經不是真的。回報裡列出改了哪些 tag，`plan` 才知道要為它們開票。
沒有對到修訂歷史的場景一個字都不動。`plan` 交來一個場景與「seam 上看不到什麼」時，只改那個 `Then`：換成 seam 上看得到、意思不變的業務說法，tag 不變，**不寫** endpoint、狀態碼、欄位名；意思會變或怎麼換都看不到 → 不改，回報給 `plan`，由它交回 `spec`。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] `spec.md` 狀態是 `已對焦` 才開始寫；`待對焦` 而寫了的，是使用者明說要先寫，且標 `@draft`
- [ ] 每個 AC 都有掛 `@example-<n>.<m>` 的場景，或在覆蓋表與 `.feature` 註解裡寫了為什麼沒有
- [ ] 沒有 `spec.md` 以外的例子、場景或 `Then`——包括從 `acceptance/` 旅程抄來的
- [ ] 每個 `.feature` 恰好一個狀態 tag、一個型別 tag；沒有 `規則:`
- [ ] 跑了 `check_spec.py`，沒有任何一項失敗，輸出貼在回報裡（只報數字的項目與跳過的項目照實列出）
- [ ] 沒有寫入 `spec.md`、`plan.md`、`docs/CONTEXT.md`；既有的 `.feature` 只在再跑時、只改對到修訂歷史的場景或 `plan` 指名的那個 `Then`

## 完成後

告訴對方四件事：

1. 產了哪些 `.feature`、各幾個場景、**幾個不重複步驟樣板**
2. **覆蓋表**——哪些例子還沒有場景，為什麼
3. 哪些象限是空的，空得合不合理
4. `spec.md` 有沒有需要回 `spec` 補的地方（規則與例子打架、`check_spec.py` 抓到的
   FR 沒掛任何 AC 或沒有來源）——本 skill 不改它

這些場景現在應該**全部是紅的**（step definition 還不存在）。這是對的：
綠燈要等 IMPLEMENT。一跑就綠代表這些場景沒有驗到任何東西。

由 `plan` DELEGATE 進來的 → 回 `plan`，覆蓋缺口交它記進 `## 缺口`；單獨叫用的 → 下一步 `/plan`，洞多回 `spec`。
