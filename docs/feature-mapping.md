# Feature Mapping 筆記

這份文件記錄 **Feature Mapping 這個手法本身**是什麼，轉述自 John Ferguson Smart 的原文
[Feature Mapping – a lightweight requirements discovery practice for agile teams](https://johnfergusonsmart.com/feature-mapping-a-lightweight-requirements-discovery-practice-for-agile-teams/)
（作者網站，2019-11-06；行文是書的一章，提到「前一章的 BDD 生命週期」）。引號內是原文逐字。

**它不談這個專案。** 這個手法在這裡要不要用、用在哪一步，是 [ai-sdlc.md](./ai-sdlc.md)
的事。這份文件改動的理由只有兩個：原文改了，或者我們讀錯了。它跟
[example-mapping.md](./example-mapping.md) 是同一類文件。

---

## 它是什麼

> "Feature Mapping is a simple but effective way for teams to discover, explore and deeply understand"

一則 feature 的需求。原文把它放在需求探索工作坊與 three amigos 的場合，位置是書中 BDD
生命週期的 *Speculate* 與 *Illustrate* 兩段；也可以拉到策略層級，在規劃版本時畫使用者
旅程、切 story、排優先序。它被拿來替代「比較自由形式的需求探索工作坊」。

> "Feature mapping is a breadth-first approach - we want to get a high level understanding across the whole feature."

**原文沒有提到 Example Mapping 或 Story Mapping**，也沒有拿它們比較。二手來源（DZone 轉載版
的介紹、Medium 文章）說它借用了 Jeff Patton 的 Story Mapping 與 Matt Wynne 的 Example
Mapping，但那不是這篇原文說的。

## 起點：story 與錨定敘事

開工前，產品負責人與 BA 先寫好 user story 與粗的驗收條件，再寫一段**錨定敘事**（anchoring
narrative）：一段具體、寫實的情境文字，印給每位與會者。它的用途是讓人開始想、開始問——
有哪些變化、有哪些假設——**不是**要列出所有路徑或規則。

## 四色卡與版面

| 卡 | 內容 | 放哪 |
| --- | --- | --- |
| 🟩 綠 | **example**：一句話講一個具體例子（原文例：Joe 把一個跟高風險國家有往來的客戶升級處理） | 一列的開頭 |
| 🟨 黃 | **step**：這個例子的一個步驟，用業務語言寫；另外也當**標題卡**，放在第一列各步驟的上方 | 綠卡右邊，一張接一張 |
| 🟪 紫 | **consequence**：情境成功走完後預期發生的結果 | 步驟卡右邊 |
| 🟦 藍 | **business rule**：這一組例子在說明的規則 | 該組綠卡的左邊 |

整張圖是一列一列疊起來的：每列從左到右是 規則 → 例子 → 步驟 → 結果，一列就是一個情境。
原文沒有替「列首的情境前提」另定顏色，也沒有替反例那一列另定畫法。

## 進行方式

原文用一場實際的工作坊示範，並說明「活動不一定照這個順序做」：

1. **畫一個例子**（綠卡）。具體例子帶動討論，引出反例與說不清的地方，之後也是自動化驗收
   測試的基礎。
2. **拆成步驟**（黃卡），擺在例子右邊。逐步討論每一步發生什麼，最容易挖出反例與不確定的
   地方。
3. **寫出結果**（紫卡）。

   > "Consequence are outcomes that we expect to happen when the scenario finishes successfully."

   （原文就寫 "Consequence are"。）結果卡讓討論停在業務價值與結果上，不提早設計畫面：

   > "In Feature Mapping, we use the idea of consequences to help us avoid this trap."

   原文的團隊刻意不在這時談畫面、勾選框或用哪個前端元件。
4. **找替代流程**，每個變化多一列。對每一步問「這裡還可能發生什麼」「使用者還能做什麼」。
   原文例：換成中風險國家，多出一列，結果是一般優先而不是高優先。
5. **依規則分組**（藍卡）。藍卡放在它那組例子的左邊，把具體例子連到它們說明的驗收條件。
6. **補反例**。除了正向的「Happy-Day」例子，也放反面例子，畫出規則的適用範圍——找出規則
   **不**適用的地方。

原文強調這只是「一場典型工作坊的一部分」，各列就是產出的情境。

## 從地圖到 Gherkin

> "Feature Maps are also easy to convert into high quality, highly readable acceptance criteria."

幾乎每張卡在 Given-When-Then 裡都有對應的位置，規則用與解法無關的語言寫，實作改了情境
也還站得住。原文的對照組是團隊先前寫的 Cucumber 情境：滿是點擊、選取的步驟，跟畫面綁死，
難讀也難維護。

原文的範例（大意，不是逐字）：

- **Feature**：高風險或中風險國家的審查要升級處理。敘事是：為了……，作為法遵人員，我希望
  複雜的案子升級給資深分析師。
- **角色**：Joe 是初級分析師，Sally 是資深分析師。
- **Background**：Bill 經營進口生意，從伊朗進口地毯。
- **Scenario**：升級一個跟高風險國家有往來的客戶。規則：跟高或中風險國家有往來的客戶要
  升級處理。
  - Given：Bill 請 Joe 開一個企業帳戶
  - When：Joe 被提示要升級這次審查，他升級了
  - Then：審查進到 Sally 的審查清單，標成高優先

Gherkin 是工作坊**之後**由一位開發者與一位測試人員寫的，不是在現場寫。地圖上其餘各列
（中風險國家、可選擇升級並在六個月後自動標記、另一個高風險案例、低風險的反例）原文只放了
圖，沒有給 Gherkin。

> "Feature Mapping offers a clean, natural path to high quality automated acceptance tests."

## 頻率、時間與適用情境

- **時間**：原文沒有給時間框，只說一場可能比傳統的 backlog grooming 稍長。
- **節奏**：「越早做、越常做越好」。Scrum 團隊通常在 sprint 開始前不久做；範例團隊是兩週
  一個 sprint，在開始前的週二做。
- **適用**：新的或複雜的領域、還在學 BDD 的團隊，以及想避免還沒弄懂問題就先設計解法的時候。
- **在場的人**：範例裡有主持人（法遵團隊的人）、產品負責人（領域專家）、兩位開發者、一位
  資深測試人員。白板加彩色便利貼。原文沒有給正式的主持規則。

---

這個手法在本專案要不要用、怎麼用，見 [ai-sdlc.md](./ai-sdlc.md)。
