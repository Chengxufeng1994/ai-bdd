---
name: story-splitting
description: >
  把太大的 story 切成數則仍能獨立交付價值的小 story——優先沿規則切，切不動時
  用九種切分模式，並用兩條規則挑出最好的切法。判斷「該不該切」也在這裡。
  觸發詞：「這個 story 太大」「要怎麼切」「怎麼拆成小的」「拆 story」
  「這則做不完」「規則太多了」「切成幾則比較好」「這樣切對嗎」。
  English: split this story, this story is too big, how do I slice this,
  break this into smaller stories, is this a good split.
---

# 切分 Story

輸入是一則太大的 story，輸出是數則**各自都還能交付價值**的小 story，以及切分依據。

## 使用時機

- Example Mapping 跑不完，或藍卡（規則）多到一次講不清
- 一則 story 大到一個迭代做不完
- 已經切了，但想確認切法對不對
- 拿到一坨需求，要先分成幾則才能開始澄清

## Skill Boundaries

- 要從需求找出規則與例子 → 改用 `bdd-clarify`，本 skill 只管切分
- 要把紅卡問到收斂 → 改用 `clarify-loop`
- 要決定切完之後**先做哪一則** → 那是排序，屬於 `bdd-plan`。本 skill 只切，不排
- 本 skill 不寫規格、不估點數

---

## 前置確認（Ask user）

Ask user: 若以下不明確，先問清楚。

1. **要切的 story**：一句話或一份 example map
2. **已知的規則**：有 example map 就直接用它的藍卡。沒有的話，先問對方知道哪些
   規則——**規則是最好用的切割線，沒有規則就只能靠猜**

---

## 1. 先確認它值得切

> "If you start with something that isn't an increment of value, there's no way
> to slice it smaller and get an increment of value."
> —— [Humanizing Work Guide to Splitting User Stories](https://www.humanizingwork.com/the-humanizing-work-guide-to-splitting-user-stories/)

切之前先檢查這則 story 符合 INVEST（**Small 除外，那正是要解決的**）：

| | 檢查 |
| --- | --- |
| **V**aluable | 對使用者有可觀察的價值嗎？ |
| **I**ndependent | 能不能不依賴其他技術前置就排進迭代？ |
| **N**egotiable | 細節還有討論空間，不是規格書？ |
| **E**stimable | 團隊估得出來嗎？ |
| **T**estable | 有辦法驗證做完了嗎？ |

IMPORTANT: Valuable 不成立時**先修它，不要切**。把一個沒有價值的東西切小，只會
得到數個更小的沒有價值的東西——而且它們看起來像進度。

## 2. 判斷該不該切

來自 Example Mapping 的三個訊號
（[出處](https://cucumber.io/blog/bdd/example-mapping-introduction/)）：

| 訊號 | 意義 |
| --- | --- |
| 藍卡（規則）太多 | story 太大 |
| 25 分鐘內 map 不完 | 太大**或**太不確定——先分清是哪一種 |
| 一條規則掛太多例子 | 那條規則本身該拆成數條，不一定要切 story |

第三項容易誤判：**例子多不代表 story 大**，可能只是規則講得太粗。先試著把那條
規則拆成兩三條，再看還需不需要切。

---

## 3. 優先沿規則切

> "Rules make natural fault lines for slicing your story."
> —— [Example Mapping introduction](https://cucumber.io/blog/bdd/example-mapping-introduction/)

規則是**獨立成立的約束**，所以沿它切，切出來的每一則都還是可驗收的。

在 BDD 流程裡這是預設做法，因為 Example Mapping 已經把規則攤在桌上了——別的方法
要先猜切割線在哪，這裡不用猜。

```
Story: 記錄訓練並看到總容量
  Rule 1 每組要有重量與次數        ┐
  Rule 2 重量不可為負              ├→ Story A：記錄一次訓練
  Rule 3 一次訓練至少一組          ┘
  Rule 4 容量 = Σ(重量 × 次數)     ┐
  Rule 5 熱身組不計入              ├→ Story B：看到總容量
  Rule 6 自體重動作另計            ┘
```

判準：**切完之後，每一則的規則集合是不是自己就完整？** A 的規則不需要知道容量
怎麼算；B 的規則不需要知道資料怎麼進來。互相不需要，就切對了。

## 4. 規則切不動時的九種模式

規則數量少、或規則之間纏在一起時，換這幾種
（[出處](https://www.humanizingwork.com/the-humanizing-work-guide-to-splitting-user-stories/)）：

| 模式 | 什麼時候用 |
| --- | --- |
| **Workflow Steps** | 流程有多個步驟——先做最簡單的頭尾貫通，再補中間與特例 |
| **Operations (CRUD)** | story 裡出現「管理」這種字眼，其實是好幾個操作 |
| **Business Rule Variations** | 幾種同樣複雜的情境走不同規則——每個變體一則 |
| **Variations in Data** | 大小來自資料的複雜度——先做「夠用」的資料模型，複雜的後補 |
| **Data Entry Methods** | 複雜度集中在 UI——先做最陽春的輸入方式 |
| **Major Effort** | 工作量集中在第一個實作——先不決定哪個變體先做 |
| **Simple/Complex** | 把變體抽出去，核心留最簡單的 |
| **Defer Performance** | 複雜度來自非功能需求——先做慢的，再做快的 |
| **Break Out a Spike** | **最後手段**：實作方式根本不清楚，先 time-box 探索 |

Spike 排最後是有理由的：它**不交付價值**，只交付知識。用它代表你承認自己還不
知道怎麼做，那是誠實的，但不該是第一反應。

## 5. 有好幾種切法時，選哪一種

兩條規則，依序套用（同上出處）：

**規則一：選那個讓你能丟掉一則的切法。**

切分的目的不只是變小，是**讓低價值的部分現形**。如果某種切法切出一則你看了之後
說「這個其實可以不做」，那就是最好的切法——你剛剛省下的是整則的成本。

**規則二：選那個切出來大小比較平均的。**

四則各 2 點，優於一則 5 點加一則 3 點。大小懸殊的切法通常代表切在錯的地方，
大的那則裡還藏著沒被分開的東西。

---

## 6. 驗收：垂直切片

每一則都必須是**垂直切片**——穿過 UI、業務邏輯、資料庫，交付可觀察的價值。

```
✗ 水平切（依架構層）        ✓ 垂直切（依行為）
┌─────────────┐            ┌───┬───┬───┐
│  前端 story │            │ A │ B │ C │  每則都從上到下
├─────────────┤            │   │   │   │  各自能交付
│  後端 story │            │   │   │   │
├─────────────┤            │   │   │   │
│  DB story   │            └───┴───┴───┘
└─────────────┘
```

「先做後端 API，前端下個 sprint」是最常見的錯誤切法：兩則都無法單獨驗收，而且
第一則做完時**沒有任何使用者行為改變**——你無法知道它對不對，只能等第二則。

檢查清單：

- [ ] 每則都能單獨排進迭代並交付
- [ ] 每則的規則集合自己完整，不需要另一則的規則才說得通
- [ ] 沒有任何一則的名字裡有「API」「資料表」「後端」這種實作字眼
- [ ] 至少有一則是你敢考慮不做的

最後一項是規則一的反向檢驗。**每一則都非做不可，通常表示切分只是把工作分段，
沒有把價值分開。**

---

## 產物

切分結果寫進 `docs/bdd/` 底下，不修改專案既有檔案：

```markdown
## 切分：<原始 story>

**依據**: 沿規則切 / <九種模式之一>
**為什麼選這個切法**: <規則一或規則二的理由>

| | Story | 規則 | 敢不敢不做 |
|---|---|---|---|
| A | <story> | Rule 1、2、3 | 不敢——核心 |
| B | <story> | Rule 4、5 | 敢，可延後 |

### 被否決的切法
<試過但不採用的，連同理由——通常是「切出來不是垂直切片」或「大小太懸殊」>
```

「被否決的切法」值得留：半年後有人想重切時，那段會告訴他哪些路已經走過。

已有 example map 時，切分後**每則各自一份 map**，規則編號沿用原本的（`Rule 4` 到了
Story B 仍然叫 `Rule 4`），**不要重編**——下游的 Gherkin 情境會回指這些編號。

---

## 完成後

- 告訴對方切成幾則、依據是什麼、哪一則你認為可以不做
- **不要順便排順序**——那是 `bdd-plan` 的事。切分只回答「分成哪幾則」，
  不回答「先做哪一則」
- 每則若還沒有 example map，下一步是對每則跑 `bdd-clarify`
