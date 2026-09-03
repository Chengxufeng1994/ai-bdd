# BDD 實踐參考

這份文件記錄 **BDD 這套方法本身**是什麼，轉述自 Cucumber 官方文件與
Example Mapping 的原始出處。

**它不談這個專案。** 不提六步流程、不提任何 skill、不提產物放哪裡——那些是
[ai-sdlc.md](./ai-sdlc.md) 的事。這份文件改動的理由只有兩個：Cucumber 改了
定義，或者我們讀錯了。

> 這條分界不是潔癖。上一版把本專案的主張（「若 plugin 只做最後一段，它就只是
> 測試產生器」）寫進轉述裡，結果是**沒有人回頭去對原文**——而原文對其中一個
> 實踐的定義，跟上一版寫的不一樣。

---

## 三個實踐

BDD 的核心是三個實踐，Cucumber 用三個時態區分它們：

| 實踐 | 問的是 | 定義（原文） |
| --- | --- | --- |
| **Discovery** | What it **could** do | "Take a small upcoming change to the system — a User Story — and talk about concrete examples of the new functionality to explore, discover and agree on the details of what's expected to be done." |
| **Formulation** | What it **should** do | "Document those examples in a way that can be automated, and check for agreement." |
| **Automation** | What it **actually does** | "Implement the behaviour described by each documented example, starting with an automated test to guide the development of the code." |

出處：[cucumber.io/docs/bdd](https://cucumber.io/docs/bdd/)。

---

### ① Discovery — 探索

- **輸入**：一則 user story
- **輸出**：達成共識的具體例子
- **起訖**：從探索問題開始，到「找出有價值的例子」為止

**手法：Example Mapping**（Matt Wynne）與 **Three Amigos**（業務／開發／測試三方）。

Example Mapping 的四色卡：黃＝story、藍＝rule、綠＝example、紅＝question。
它是一場約 25 分鐘的對話，不是文件工作，結束時以**大拇指投票**決定就緒與否。

出處：[Example Mapping introduction](https://cucumber.io/blog/bdd/example-mapping-introduction/)。
以下幾點直接來自該文：

- 25 分鐘內 map 不完 ＝ story 太大或不確定性太高
- 紅卡太多 ＝ 不確定性高；藍卡太多 ＝ story 太大；一條規則掛太多例子 ＝ 規則該拆
- 未就緒有兩條路：切小，或讓產品方回去把問題問到答案
- 切分沿著規則走："Rules make natural fault lines for slicing your story"

---

### ② Formulation — formulate

- **輸入**：Discovery 選出來的例子（綠卡）
- **輸出**：人與機器都讀得懂的文件——可被自動化的規格
- **起訖**：從「已經有例子」開始，到「規格寫好且各方同意」為止

注意定義裡的第二個動詞：**check for agreement**。Formulation 不只是把例子寫成
Gherkin，還包含拿寫好的東西回去確認「這是我會寫的樣子嗎」。少了那一步，
它就只是換一種格式重寫一次，沒有驗證共識有沒有真的存在。

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

命令式寫法讓每次 UI 改動都要改 feature 檔，是 BDD 專案腐爛的常見原因。

---

### ③ Automation — 自動化

- **輸入**：已寫成規格的例子
- **輸出**：**可運作的系統程式碼，且自動化測試全綠**
- **起訖**：從寫下自動化測試開始，到**實作完成並通過測試**為止

**這一段最容易被讀窄。** Automation 不是「把情境接上膠水層，然後留在紅色」——
定義寫得很清楚：*implement the behaviour*，輸出是 *working system code*，
結束於 *implementation is complete and tested*。

寫 step definition 只是它的起手式。**產品程式碼是它的產物。**

自動化測試在這裡的角色是 **guide the development of the code**——先寫測試，
用它驅動實作，而不是實作完再補測試。這就是雙迴圈的位置：

```
外迴圈：scenario red
   └─ 內迴圈：unit test red → green → refactor
   └─ 內迴圈：red → green → refactor
   └─ ⋯
外迴圈：scenario green ✓
```

外迴圈由業務行為驅動（做對的事），內迴圈由設計驅動（把事做對）。
**outside-in**：從使用者看得到的行為往內推，而不是先建 model 再往外接。

雙迴圈**發生在 Automation 內部**，不是它之後的另一個階段。

---

## BDD 沒有說的

這一節跟上面同等重要，因為它劃出方法的邊界。以下三件事在 BDD 的三個實踐裡
**找不到任何指引**，不是因為它們不重要，而是因為它們不屬於這套方法：

| 沒說的 | 說明 |
| --- | --- |
| **順序與排程** | 哪個情境先做、工作怎麼切成任務、誰卡誰——Cucumber 的文件只提到「working in smaller increments」，沒有任何排序或切分的指引 |
| **治理與維護** | `.feature` 會漂移（孤兒 step、重複情境、已刪功能的殘留），但怎麼稽核與清理不在三個實踐裡 |
| **產物存放** | 規格檔放哪、澄清紀錄要不要留、留成什麼格式——完全沒有規定 |

**Living documentation** 常跟 BDD 一起被提到，但它是產物的一個**性質**
（`.feature` 同時是規格、測試與文件），不是第四個實踐。維持那個性質需要治理，
而治理如上所述不在 BDD 的範圍內。

同理，**Verification**（跑全套、CI、flaky 處理、失敗判讀）是 Automation 產出
之後的持續活動，不是獨立的第四個實踐。

---

需要這些沒被涵蓋的東西時，得自己補。本專案怎麼補的見
[ai-sdlc.md](./ai-sdlc.md)。
