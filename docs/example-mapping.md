# Example Mapping 筆記

這份文件記錄 **Example Mapping 這個手法本身**是什麼，轉述自 Matt Wynne 的原文
[Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/)
（Cucumber blog，2015-12-08）。引號內是原文逐字。

**它不談這個專案。** 四色卡在 `spec.md` 裡叫什麼、診斷在這裡改了什麼意義、就緒怎麼判——
那些是 `skills/example-mapping/SKILL.md` 與 [ai-sdlc.md](./ai-sdlc.md) 的事。這份文件改動
的理由只有兩個：原文改了，或者我們讀錯了。它跟 [bdd.md](./bdd.md) 是同一類文件；
`bdd.md` 只用一段話交代 Example Mapping 在 Discovery 裡的位置，細節在這裡。

---

## 它要解決什麼

story 拉進開發前，需要一場對話來釐清並確認驗收條件。這場對話有很多名字——backlog
refinement、three amigos、specification workshop、discovery workshop——但原文指出共同
的問題：

> "it's unstructured, it takes too long and gets boring."

結果是團隊不定期做、做得不一致，或乾脆放棄。Example Mapping 是一個「simple, low-tech」
的方法，讓這場對話變短、變有產出。

## 四色卡

對話中除了例子，還會冒出另外三種值得記下的東西。四種各用一色索引卡：

| 卡 | 內容 | 放哪 |
| --- | --- | --- |
| 🟨 黃 | **story**：正在討論的那則；也包括對話中新發現、或切出去延後的 story | 桌面最上方 |
| 🟦 藍 | **rule**：已知的驗收條件或規則——"rules that summarise a bunch of examples, or express other agreed constraints about the scope of the story" | 黃卡下方一字排開 |
| 🟩 綠 | **example**：說明某條規則的具體例子，一條規則可以有一個或多個 | 對應的藍卡下方 |
| 🟥 紅 | **question**：在場沒有人知道正確結果的情境，或為了往下談而先做的**假設** | 記下來，繼續談 |

注意藍卡的定義：原文寫的是 "each of the acceptance criteria, or **rules**"——在 Example
Mapping 裡，**驗收條件是藍卡，不是綠卡**。

進行方式：寫下黃卡 → 擺出已知的藍卡 → 每條規則補綠卡 → 談不出答案的就寫紅卡、跳過。
持續到 "the group is satisfied that the scope of the story is clear, or we run out of time"。

## 桌面就是即時回饋

攤開的卡片本身就是團隊目前理解的視覺化，形狀會說話：

| 看到 | 意思 | 原文建議 |
| --- | --- | --- |
| 滿桌紅卡 | 對這則 story 還有很多要學 | — |
| 滿桌藍卡 | story 又大又複雜 | 切掉：拿一張新黃卡，把切出的 story 放回 backlog |
| 一條規則下掛很多綠卡 | 規則可能太複雜 | 是不是有好幾條規則混在一起，該拆開 |

反方向也有一條，常被忽略：

> "some rules are so obvious you don't need examples at all."

大家都懂的規則就不用硬擠例子——原文的說法是不必 "grind out examples like BDD-brainwashed
automatons"。

## 25 分鐘的框

一小組人應該能在**約 25 分鐘**內攤完一則理解清楚、大小合適的 story。攤不完有三種可能：

1. 還在學這個手法（沒關係）
2. story 太大（"definitely not fine"）
3. 不確定性還太高

對應的兩條路：**切小**，或讓產品方回去做功課，之後再排一次 mapping。

時間到了用 **thumb-vote（大拇指投票）** 決定能不能拉進開發。還剩紅卡不一定就是不行——
如果大家覺得問題夠小、可以邊做邊解，"Let the group decide."

## 例子要粗，不要 Gherkin

原文明確反對在 mapping 當下寫正式的 Gherkin：圍著投影機、有人在 IDE 裡打 Given/When/Then，
"generally I think it's a bad idea"，因為它會讓人從對話的真正目的分心。

真正的目的不是「把 AC 變成例子再變成測試」，而是：

> "to reach a *shared understanding* of what it will take for the story to be done."

所以綠卡用 **friends episode** 命名法寫粗例子：

- The one where the customer forgot his receipt
- The one where the product is damaged
- The one where the product was bought 15 days ago

想更具體時也還不必用 Given/When/Then。而有一條判準直接決定一張卡是綠還是紅：

> "When the outcome (Then) is unclear, you don't have an example, you have a question."

## 紅卡：把未知變成已知的未知

對話繞圈，是因為資訊不夠——缺了某個人、需要使用者研究、或需要 spike。這時不要讓每個人
輪流發表意見，**寫下問題、繼續談**：

> "You've just turned an *unknown unknown* into a *known unknown*."

原文說，很多團隊光靠這一點，就讓 discovery workshop 從冗長的閒聊變成短而有效的對話。

## 誰該來

最少是 three amigos：開發、測試、產品各一。歡迎找營運、UX 或任何能「發現新問題、或在
對話中把問題變成決定」的人。

剛學的時候可以指定一位 facilitator，專門確保講出來的每件事都落在卡片上——例子和問題
在房間裡飛得很快，要紀律才能都擺上桌。

## 那什麼時候寫 Gherkin

原文不否定一起寫 Gherkin 的價值，特別是專案初期——那是建立 ubiquitous language 的地方。
但它指出兩件事需要不同的思考模式：

- 決定哪些例子在範圍內、釐清背後的規則（mapping）
- 把例子寫成 Gherkin（formulation）

對已經上軌道、領域語言成熟的團隊，作者偏好：產品方把精力花在 mapping，Gherkin 交給另外
兩位 amigo 起草，再由產品方回饋：

> "Is that how I would have written it?"

這一步同時檢驗了 mapping 有沒有真的把產品方的知識傳出去。（這就是 [bdd.md](./bdd.md)
§② 說的 Formulation 的 *check for agreement*。）

## 頻率與遠端

- **頻率**：產品方很忙，尊重他們的時間。作者看過有效的節奏是**每隔一天**一次，一次挑一則
  story、給 25 分鐘，然後回去工作。一次塞很多只會耗盡精力。
- **分散式團隊**：共用文件的條列、用儲存格顏色代表卡片的試算表、心智圖都行。關鍵是又快
  又好操作，"so you can focus on the conversation"。

## 最後幾條提醒

- 用之前要先分清楚**規則與例子**的差別（原文指向 Liz Keogh 的
  [Acceptance Criteria vs Scenarios](http://lizkeogh.com/2011/06/20/acceptance-criteria-vs-scenarios/)）
- 這場對話的整個目的是 "discover the stuff you don't already know"，所以沒有笨問題
- **切 story 沿著規則切**："rules make natural fault lines for slicing your story"——
  能延後的盡量延後，先解決問題的核心，複雜度以後再加

---

這個手法在本專案怎麼落地（沒有三方、沒有時鐘、票不是 agent 投的）見
`skills/example-mapping/SKILL.md`；為什麼這樣改見 [ai-sdlc.md](./ai-sdlc.md)。
