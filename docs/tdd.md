# TDD 實踐參考

這份文件記錄 **TDD（Test-Driven Development）這套方法本身**是什麼，轉述自
Kent Beck 的 Canon TDD、Martin Fowler 的 bliki 與 Robert C. Martin 的三律。

**它不談這個專案。** 不提六步流程、不提任何 skill、不提我們同不同意——那些是
[ai-sdlc.md](./ai-sdlc.md) 的事。這份文件改動的理由只有兩個：來源改了定義，
或者我們讀錯了。

它跟 [bdd.md](./bdd.md)、[sdd.md](./sdd.md) 是同一類文件，但位置比較特別：
**BDD 是從 TDD 長出來的**（見下方〈TDD 沒有說的〉），而 TDD 就是 `bdd.md` §③
那個雙迴圈圖裡的**內迴圈**。那張圖不在這裡重畫，去那份看。

---

## Canon TDD

Beck 在 2023 年寫了一篇正典定義，理由是他「一直讀到胡說八道」——這套方法被
批評的版本經常不是他定義的版本。五個步驟，逐字：

> 1. "Write a list of the test scenarios you want to cover"
> 2. "Turn exactly one item on the list into an actual, concrete, runnable test"
> 3. "Change the code to make the test (& all previous tests) pass"
> 4. "Optionally refactor to improve the implementation design"
> 5. "Until the list is empty, go back to #2"

出處：[Kent Beck, Canon TDD](https://newsletter.kentbeck.com/p/canon-tdd)。

**第一步經常被整個略過。** 那份 test list 是 Beck 所謂的行為分析——先列出這次
變更要處理哪些情況，而**不是**列出要寫哪些函式。他明確警告不要在這一步混入
實作決定。

TDD 要達成的四件事，逐字：

> "Everything that used to work still works. The new behavior works as expected.
> The system is ready for the next change. The programmer & their colleagues feel
> confident in the above points."

注意第四項：**信心是列出來的目標之一，不是副作用。**

Beck 對批評者設了一道門檻：

> "If you plan on critiquing TDD & you're not critiquing the following workflow,
> then you're critiquing a strawman."

他列出的常見誤解包括：把 test list 上的項目**全部**先寫成測試才開始讓它們通過、
寫沒有斷言的測試來衝覆蓋率、刪掉斷言讓測試假性通過、以及把重構混進「讓它通過」
那一步。

---

## 三律

Robert C. Martin 把同一件事壓成三條逐秒尺度的規則，逐字：

> 1. "You must write a failing test before you write any production code."
> 2. "You must not write more of a test than is sufficient to fail, or fail to compile."
> 3. "You must not write more production code than is sufficient to make the
>    currently failing test pass."

他描述四層巢狀迴圈，時間尺度各差一個數量級：

| 迴圈 | 尺度 | 內容 |
| --- | --- | --- |
| Nano | 逐秒 | 三律本身，逐行套用 |
| Micro | 分鐘 | Red-Green-Refactor |
| Milli | 十分鐘 | Specific/Generic——測試越來越具體，產品程式碼越來越通用 |
| Primary | 小時 | Boundaries——退一步檢查架構 |

關於重構的位置：

> "Refactoring is a continuous in-process activity, not something that is done
> late (and therefore optionally)."

出處：[The Cycles of TDD](https://blog.cleancoder.com/uncle-bob/2014/12/17/TheCyclesOfTDD.html)。

第三層那條 **Specific/Generic** 值得單獨記住：它說的是隨著測試累積，讓它們
通過的正確作法是把實作**一般化**，而不是繼續補 if。這是 TDD 宣稱能驅動設計的
實際機制。

---

## Red-Green-Refactor

Fowler 的三步版本，逐字：

> 1. "Write a test for the next bit of functionality you want to add."
> 2. "Write the functional code until the test passes."
> 3. "Refactor both new and old code to make it well structured."

以及他點名的頭號失敗模式：

> "The most common way that I hear to screw up TDD is neglecting the third step."

略過第三步的結果他寫成 "a messy aggregation of code fragments"。

TDD 之所以被稱為設計技術而不只是測試技術，理由在這句：

> "thinking about the test first forces us to think about the interface to the
> code first"

出處：[martinfowler.com/bliki/TestDrivenDevelopment](https://martinfowler.com/bliki/TestDrivenDevelopment.html)。
方法本身出自 Kent Beck，1990 年代末，Extreme Programming 的一部分。

---

## 兩個學派

同樣照著 red-green-refactor 走，對「一個測試該隔離到什麼程度」有兩種相反答案：

| | Classicist（Detroit／Chicago） | Mockist（London） |
| --- | --- | --- |
| **用不用真物件** | "use real objects if possible and a double if it's awkward to use the real thing" | "always use a mock for any object with interesting behavior" |
| **驗什麼** | State verification——跑完之後檢查狀態 | Behavior verification——檢查有沒有對協作者發出正確的呼叫 |
| **隔離時機** | 只在必要時（外部服務之類） | 系統性地隔離所有有行為的協作者 |

代價寫在 Fowler 這句：

> "Mockist tests are thus more coupled to the implementation of a method.
> Changing the nature of calls to collaborators usually cause a mockist test to break."

他自己的立場也寫得很清楚：

> "I don't see any compelling benefits for mockist TDD, and am concerned about
> the consequences of coupling tests to implementation."

出處：[Mocks Aren't Stubs](https://martinfowler.com/articles/mocksArentStubs.html)。

**這條分歧不是風格偏好，它決定了測試在重構時會不會碎掉。** 綁在互動上的測試，
會把「怎麼做」寫死進斷言裡。

---

## TDD 沒有說的

這一節劃出方法的邊界。以下幾件事在 TDD 的迴圈裡**找不到指引**：

| 沒說的 | 說明 |
| --- | --- |
| **測什麼、從哪裡開始** | Canon TDD 的第一步是「寫下 test list」，但**怎麼想出那份清單完全沒有規定**。迴圈從清單存在之後才開始轉 |
| **測試該叫什麼名字** | 命名不在三律、不在 red-green-refactor 裡 |
| **一次該測多少** | 「sufficient to fail」界定的是測試的大小，不是行為的切分粒度 |
| **誰來確認這是對的行為** | TDD 全程是一個程式設計師跟他的程式碼。沒有第二個角色，也沒有共識機制 |

**這份缺口清單有名字。** Dan North 在 2006 年提出 BDD，正是因為他發現學 TDD 的人
反覆卡在同樣幾個問題上：從哪開始、要測什麼、什麼不用測、一次測多少、測試該
怎麼命名。他的診斷是「test」這個詞本身在誤導——它把注意力推向「驗證已經寫好的
東西」，而不是「先描述期望的行為」；他的處方是改用 behaviour 的語彙，以及用
`should` 開頭的命名慣例。

出處：[Dan North, Introducing BDD](https://dannorth.net/introducing-bdd/)（本段
為轉述，非逐字引用）。BDD 那套方法本身見 [bdd.md](./bdd.md)。

---

## 什麼時候它不適用

2014 年 Kent Beck、Martin Fowler 與 DHH 的系列對談留下的共識，比當時的標題有用：

- 自動化回歸測試很重要——這點三人沒有分歧
- 對**演算法明確、邊界清楚**的問題，TDD 很有效
- 對**還沒被理解清楚**的問題，效果就差得多
- 判斷某個處境適不適合用 TDD，本身是一項需要技巧的判斷

出處：[Is TDD Dead?](https://martinfowler.com/articles/is-tdd-dead/)（本段為轉述）。

把它當唯一正解、或當成必須反對的東西，都會錯過這個結論：**TDD 是達成信心的
其中一條路，不是唯一一條。**

---

需要這些沒被涵蓋的東西時，得自己補。本專案怎麼補的見 [ai-sdlc.md](./ai-sdlc.md)。
