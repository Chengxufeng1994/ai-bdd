# BDD 實踐參考

這份文件記錄 BDD 這套方法**本身**是什麼：三個實踐、雙迴圈、常見反模式。
它是設計依據，不是使用說明。

要建什麼、怎麼切、什麼順序，見 [PLAN.md](../PLAN.md)。

> 前提：BDD 的核心是 **Discovery → Formulation → Automation** 三個實踐，
> 測試自動化只是最後一段。若 plugin 只做最後一段，它就只是測試產生器，不是
> BDD 流程。

---

## 外層：三個實踐

### ① Discovery — 需求探索

把模糊需求逼成具體例子。手法：**Example Mapping**（Matt Wynne）、
**Three Amigos**（業務 / 開發 / 測試三方）。

- **輸入**：一則 user story
- **輸出**：規則（rules）＋ 每條規則的具體例子（examples）＋ 未解問題（questions）
- **判準**：紅卡（問題）太多 → 這個 story 還不能進 sprint

Example Mapping 的四色卡：黃＝story、藍＝rule、綠＝example、紅＝question。
它是一場約 25 分鐘的對話，不是文件工作，結束時以**大拇指投票**決定就緒與否。

出處：[Example Mapping introduction](https://cucumber.io/blog/bdd/example-mapping-introduction/)。
以下幾點直接來自該文，不是本專案的發明：

- 25 分鐘內 map 不完 ＝ story 太大或不確定性太高
- 紅卡太多 ＝ 不確定性高；藍卡太多 ＝ story 太大；一條規則掛太多例子 ＝ 規則該拆
- 未就緒有兩條路：切小，或讓產品方回去把問題問到答案
- 切分沿著規則走："Rules make natural fault lines for slicing your story"
- 事後由產品方讀開發寫出的 Gherkin，問「這是我會寫的樣子嗎」——那是在檢驗對話
  有沒有真的發生，屬於 SPEC 步驟的驗收機制

**AI 能做**：扮演提問者，針對每條規則追問邊界、例外、空值、權限、併發。
把「使用者可以上傳檔案」逼成「10MB 剛好 / 10MB+1 byte / 0 byte / 副檔名偽裝」。

**AI 不能做**：取代三方對話。它沒有業務脈絡的最終裁量權。

---

### ② Formulation — 情境撰寫

把選出來的例子寫成業務讀得懂的 Gherkin。

- **輸入**：Discovery 的綠卡
- **輸出**：`.feature` 檔（Feature / Rule / Scenario / Given-When-Then /
  Background / Scenario Outline）
- **判準**：業務方能讀懂並確認「對，就是這樣」

最大陷阱是**命令式 vs 宣告式**：

```gherkin
# ✗ 命令式 — 這是 UI 腳本，不是行為規格
When I click "Login"
And I type "admin" into "#username"
And I click the submit button

# ✓ 宣告式 — 描述行為，不綁實作
When the administrator signs in
Then the audit log records the sign-in
```

命令式寫法讓每次 UI 改動都要改 feature 檔，是 BDD 專案腐爛的頭號原因。

**AI 能做**：草擬 Gherkin、稽核是否命令式化、統一詞彙（ubiquitous language）、
找出重複情境。

---

### ③ Automation — 自動化

把情境接上程式。

- **輸入**：`.feature`
- **輸出**：step definitions（膠水層）＋ 底層 driver（API client / page object /
  測試替身）
- **狀態**：此時情境應該是 **red**

**AI 能做**：產 step 骨架、**優先比對既有 step 避免重複**（step 重複是膠水層
失控的起點）、指出該拆的 conjunction step（`Given A and B`）。

---

## 內層：雙迴圈開發

Automation 之後進入實作，這是 BDD 與 TDD 接合的地方：

```
外迴圈（BDD）：scenario red
   └─ 內迴圈（TDD）：unit test red → green → refactor
   └─ 內迴圈：red → green → refactor
   └─ ...
外迴圈：scenario green ✓
```

外迴圈由業務行為驅動（做對的事），內迴圈由設計驅動（把事做對）。
**outside-in**：從使用者看得到的行為往內推，而不是先建 model 再往外接。

---

## 持續：常被忽略的兩段

### ④ Verification — 驗證與回歸

跑全套、CI 整合、tag 分層（`@smoke` / `@wip` / `@slow`）、flaky 處理。

**關鍵判斷**：情境失敗時，是**真 bug**、還是**規格已改而情境過期**？
兩者處理方式相反，混淆會導致「改測試讓它變綠」的災難。

### ⑤ Living Documentation — 活文件與治理

`.feature` 檔同時是規格、測試、文件。但它會漂移：

- 情境描述的行為與程式實際行為不一致
- 孤兒 step（沒有任何情境使用）
- 重複情境（同一規則被驗證五次）
- 已刪功能的情境還留著

沒有治理，feature 目錄兩年後會變成沒人敢刪的墳場。

---

## 為什麼 Discovery 與 Formulation 要拆成兩個 skill

兩段雖然相鄰，卻是**方向相反的工作模式**：Discovery 是模型不停追問、使用者
回答；Formulation 是模型書寫、使用者審閱。合成一個 skill 會讓模型在兩種模式
間搖擺，實務上的結果通常是跳過追問直接開始寫 —— 因為書寫比追問「看起來更
有生產力」。而跳過 Discovery 正是 BDD 最常見的失敗模式：付出膠水層的維護
成本，卻沒得到任何發現價值。

---

階段如何對應到 plugin 元件、分期順序、以及尚未拍板的事項，見
[PLAN.md](../PLAN.md)。
