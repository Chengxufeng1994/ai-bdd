# 三段式：discovery 換成逼問、spec 只寫、formulation 出 Gherkin

**日期**：2026-09-14
**狀態**：待實作
**取代**：[`2026-09-14-pipeline-restructure-design.md`](./2026-09-14-pipeline-restructure-design.md) 的**計畫二**
（那份假設 `bdd-spec` 仍然產 Gherkin，這一版把它移出去了；計畫一已經實作並合併）

---

## 最終形狀

```
bdd-discovery    設計樹 ＋ 前沿，用四色卡逼問   →  clarify-log.md ＋ input.md
bdd-spec         綜合，不訪談；先勾 seam 再寫   →  spec.md
bdd-formulation  spec.md → *.feature（DSL）     →  PM 也用它對答案
```

`example-mapping` 保留，成為 `bdd-discovery` 追問時用的手法。

---

## 為什麼改

### 一、三個 pass 是猜好的順序，前沿是算出來的

現在的 `bdd-discovery` 有 435 行 SKILL.md，三個 pass 固定順序：廣度 → 深度 → 技術。

參考 [`mattpocock-skills:grilling`](~/.claude/plugins/cache/mattpocock/mattpocock-skills/1.2.3/skills/productivity/grilling/SKILL.md)（**22 行**）的模型：

> 把它畫成一棵**設計樹**：每個決策底下分岔出依賴它的決策。**前沿**是所有前提都已定案的
> 決策——你現在就問得動、不必猜任何還沒聽到的答案的那些問題。把整個前沿一次問完，
> 然後等。每輪答案把前沿往外推。**一題的答案若取決於本輪還開著的另一題，它屬於下一輪。**

三個 pass 把兩件事混在一起了：技術問題裡有一半不依賴任何業務答案（「這個模組現在有沒有
交易邊界」）——那些**第一輪就該問**；而業務問題裡有些依賴範圍決定——那些**不該在第一輪問**。
`Pass 3` 把前者押到最後，`Pass 1` 把後者提到最前。

前沿模型讓順序從「猜的」變成「算的」。

### 二、「找事實是你的工作，永遠不是使用者的」

`grilling` 有這條，我們沒有：

> 前沿問題需要環境裡的事實（檔案系統、工具）時，**派 sub-agent 去找**——不要問使用者
> 任何你自己查得到的東西。決策才是使用者的。

`references/technical-probes.md`（105 行）有一半是可以自己查的：既有測試慣例、模組邊界、
現有的介面形狀。我們在問使用者他自己的 codebase。

### 三、`.feature` 本來就是三方拿來對答案的東西

設計過程中一度規劃第二份給 PM 讀的產物（草稿 `.feature`、`example-mapping.md`、
`example.md`）。那個規劃預設了「加上 tag 與封閉步驟文法之後 PM 就讀不懂」。

**但 Cucumber 存在的全部理由就是「可執行的規格，業務讀得懂」。** 如果 PM 讀不動
`.feature`，真正的問題是 `step-grammar.md` 太技術，不是需要第二份檔案。

拿掉那個規劃，順便把一個風險變成診斷：**PM 讀不動 = 步驟文法該修**。

---

## 決策

### D1 · `bdd-discovery` 換成設計樹 ＋ 前沿

三個 pass 消失。一輪的流程：

1. 算出前沿——所有前提已定案的問題
2. 一次問完，編號，每題附建議答案
3. 等使用者回答
4. 重算前沿，下一輪

問題格式沿用 `grilling` 的：

```
❓ **Q1** - **<問題標題>**：<問題本體，可以多段，可以含選項>

➡️ <你的建議答案>
```

**前沿空了就結束**——每一條分支都走過，沒有東西被默默假設。

### D2 · 逃生口要保住，而 `grilling` 的建議答案會推人

`grilling` 每題附一個建議答案。那對「有領域慣例可循」的題目是好的，但它**把人推向接受**。

我們現行的四選項格式有第五個逃生口（「不知道／先跳過」），存在的理由寫在現行 SKILL.md 裡：
沒有逃生口時，人在不確定時會挑一個看起來最合理的，而**被記成決定的猜測比開著的問題更糟**。

保留那個性質，但不保留四選項的形式（一輪問五題、每題四選項會長到沒人讀得完）：

MUST: 沒有站得住的預設時，建議欄要寫「**這題我不建議猜——要去問 <誰>**」，而不是挑一個。

MUST NOT: 每一題都給建議。純業務價值取捨（「要比較週與週、還是同一動作的進步」）沒有
「常見答案」，給了只是把偏好偽裝成建議。

### D3 · 16 個追問面向留著，但從「三個 pass」變成「產生前沿問題的清單」

`skills/bdd-spec/scripts/status.py` 的 `DIMS` 是 11 個業務面向 ＋ 5 個技術面向，而
`status.py` 用它算追問覆蓋率——**那是這個 repo 唯一看得出「有沒有真的問過邊界」的東西**。

面向不動，改的是它們怎麼被用：

| | 現在 | 改成 |
| --- | --- | --- |
| 業務 11 個 | Pass 2 逐項掃過 | 每輪算前沿時，對每個已定案的決策問「這 11 個面向有哪個還沒問」 |
| 技術 5 個 | Pass 3 逐項掃過 | 同上，而且**能自己查的先派 sub-agent 查**（D9） |

`status.py` 一行都不用改——它讀的是 `## Open Questions` 表的 `面向` 欄，那一欄照樣填。

### D4 · `example-mapping` 是 discovery 的手法，不是下游步驟

四色卡（🟨 story／🟦 rule／🟩 example／🟥 question）是**追問時攤在桌上的東西**，不是事後的
檢查。放在 spec 之後等於把已經知道的事重新推導一次。

`example-mapping` skill 保留，由 `bdd-discovery` 在逼問時呼叫。它現在的內容（四色卡的
操作定義、四條診斷、就緒判定）大致適用，要改的是**輸入從 `spec.md` 變成正在進行的訪談**。

就緒判定留在它身上——前沿空了之後，由它攤開地圖、跑診斷、把就緒問句交給使用者。

### D5 · `spec.md` 吸收 `to-spec` 的章節 ＋ 四件新東西

參考 [`mattpocock-skills:to-spec`](~/.claude/plugins/cache/mattpocock/mattpocock-skills/1.2.3/skills/engineering/to-spec/SKILL.md)。
它的模板（Problem Statement／Solution／User Stories／Implementation Decisions／
Testing Decisions／Out of Scope／Further Notes）正是這個 repo **併掉之前的舊
`spec-format.md`**——血緣在這裡。

現行 `spec-format.md` 有 24 個 `## ` 小節。這一版新增四件：

| 新增 | 裝什麼 |
| --- | --- |
| **目標使用者** | 併進現有的 `## Personas`，但要寫明「這個 feature 是為誰做的」而不只是「誰在流程裡出現」 |
| **技術偏好與限制** | 語言、框架、不准引入的依賴——現在散在 `## Implementation Decisions` 與 `## Assumptions / Constraints` 兩邊 |
| **Known Boundaries** | **全新**：`總是要做` ／ `要先問` ／ `絕對不做` 三張清單 |
| **Further Notes** | 放不進其他節、但下一個人會想知道的 |

**Known Boundaries 是這一版最值得加的東西。** 它等於每個 feature 自帶一份 agent 約束：

```markdown
## Known Boundaries

**總是要做**
- 金額一律用整數分計算，不用浮點 ← Q-7

**要先問**
- 動到既有的出場閘門控制流程 ← 推論（那是已發包的設備）

**絕對不做**
- 寫入車牌辨識相關的任何欄位 ← PRD §5（區權會否決）
```

每一條照樣標來源。這三個分類的判準是**誰承擔後果**：自己做、問了再做、做了要負責的事不做。

### D6 · `bdd-spec` 寫之前先勾 seam 並跟使用者確認

`to-spec` 的第二步，我們沒有：

> 勾出你要在哪些 seam 測這個 feature。**優先用既有的 seam**。用**最高的**那一層。需要新的
> 就在你能做到的最高點提出。**跨 codebase 的 seam 越少越好，理想是一個。**
>
> 跟使用者確認這些 seam 符合他的預期。

這一步在寫 `spec.md` 之前跑，結果寫進 `## Testing Decisions`。它是 `bdd-spec` 唯一一個
**與使用者互動**的動作——不是訪談（不問需求），是確認一個技術判斷。

MUST NOT: 把 seam 的討論變成需求訪談。**發現需求有洞就寫進 `## Out of Scope` 的
「回 CLARIFY 補問」**，不要順手問掉。

### D7 · `spec.md` 是單一來源，AC 維持現行的一層

`bdd-formulation` 只讀 `spec.md`，所以 `spec.md` 必須含**具體的**例子——含數字的那種。
現行格式已經是這樣，**不動**：

```markdown
**FR-1**  When 訪客車出場, the system shall 依停留時長計費，前 30 分鐘免費。 ← PRD §3
- **AC-1.1**  停 29 分 → 收 0 元 ← PRD §3
  - Given 一台訪客車已入場並抽了紙票
  - When 它在入場後 29 分鐘於繳費機結帳
  - Then 應收金額是 0 元
```

**設計過程中一度規劃兩層**（抽象 AC ＋ 具體 EX），目的是把「哪幾條併成一個
`Scenario Outline`」這個決定寫進格式。那個規劃屬於 `example.md` 還存在的版本——它拿掉
之後，承載那個決定的需求也跟著消失。

MUST NOT: 在這一版動 AC 的層級。它是 2026-09-14 那條分支剛出貨的格式，而現在沒有已知的
需求逼我們加一層。**多一層容易，少一層難。**

`bdd-formulation` 看到骨架幾乎相同的數條 AC 時，自己決定要不要併成 `Scenario Outline`。
併不併都是合法的 Gherkin，而三個獨立的 `Scenario:` 對 PM 甚至更好讀——所以這是它的判斷，
不需要格式替它決定。

### D8 · Gherkin 移出 `bdd-spec`，`bdd-formulation` 是新 skill

`bdd-spec` 不再產 `.feature`。新 skill 讀 `spec.md`，產出套好封閉步驟文法與 tag 的
`features/<slug>.feature`。

| 移到 `bdd-formulation` | 留在 `bdd-spec` |
| --- | --- |
| `references/step-grammar.md`（15.7K） | `references/spec-format.md`（28.7K） |
| `references/state-tags.md` | `references/persona-definition.md` |
| `references/rule-taxonomy.md` | `scripts/status.py` |
| `references/coverage-report.md` | `examples/minimal-spec.md` |
| `references/artifact-location.md` | |
| `examples/` 的 12 個 Gherkin 範例 ＋ `anti-patterns.md` | |
| `scripts/check_spec.py`（23.5K） | |

**不做第二份給 PM 的版本。** `.feature` 就是對答案用的那一份——那是 Cucumber 的前提。
PM 讀不動就修 `step-grammar.md`。

接縫不變：`FR-<n>` ↔ `@rule-<n>`、`AC-<n>.<m>` ↔ `@example-<n>.<m>`。

### D9 · 能自己查的事實不准問使用者

MUST: 前沿問題需要的事實若在環境裡查得到（既有測試怎麼寫、某個模組的介面長怎樣、
有沒有現成的 seam），**派 sub-agent 去查**，不要問使用者。

MUST NOT: 因為 sub-agent 還在跑就停下來。**一個進行中的探查是一個未定案的前提**——只有
依賴它的問題要等，前沿其餘的照問。

`references/technical-probes.md` 要重寫成兩欄：**這個面向要知道什麼** ／ **自己查還是問人**。

---

## 連帶要改的

| 檔案 | 改什麼 |
| --- | --- |
| `skills/bdd-discovery/SKILL.md` | 三個 pass → 前沿模型；435 行應該大幅縮短 |
| `skills/bdd-discovery/references/technical-probes.md` | 加「自己查還是問人」欄 |
| `skills/example-mapping/SKILL.md` | 輸入從 `spec.md` 改成進行中的訪談 |
| `skills/bdd-spec/SKILL.md` | 加 seam 步驟；拿掉 Gherkin 的產出職責 |
| `skills/bdd-spec/references/spec-format.md` | 加四節；AC／EX 兩層 |
| `skills/bdd-spec/examples/minimal-spec.md` | **canonical 解析對象**，跟著改 |
| `skills/bdd-formulation/` | **全新**，接收五份 reference ＋ 13 個 examples ＋ `check_spec.py` |
| `skills/bdd-plan/SKILL.md` | 讀 `spec.md` ＋ `.feature`，指向更新 |
| `docs/ai-sdlc.md`、`CLAUDE.md`、`PLAN.md` | 三段式的敘述；skill 表 |
| `benchmark/skeleton/go/README.md` | 鏈的對應表 |

---

## 切成三份實作計畫

八個 skill、兩份 reference 大改、一個全新 skill、一支腳本搬家。塞成一份計畫會得到一份
沒人能審的東西。切點讓**每一段做完都留下一條走得通的流程**：

### 計畫一：`bdd-formulation` 拆出來

純搬移 ＋ 新 SKILL.md。`bdd-spec` 停止產 Gherkin，新 skill 接手。**格式完全不動**——
`spec.md` 還是現在的樣子，`.feature` 還是現在的樣子，只是產它的人換了。

做完之後流程完整可跑。

### 計畫二：`spec.md` 加四節 ＋ AC／EX 兩層

`spec-format.md` 加四節、fixture 跟著改（原子單位）。`bdd-spec` 加 seam 步驟。

**AC 的層級不動**——這一份只加章節，不改既有結構。

### 計畫三：`bdd-discovery` 換成前沿模型

最大的一份，也是最獨立的——它不動任何格式，只動「怎麼問」。`example-mapping` 的輸入
跟著改。

**為什麼這樣切。** 計畫一動「誰產出」、計畫二動「格式」、計畫三動「怎麼問」。混在一起
的話，任何檢查失敗都分不清是哪一層造成的——而這條 repo 前三次遷移的教訓都是
**半遷移的狀態最危險**。

---

## 驗證

| # | 驗什麼 | 怎麼驗 |
| --- | --- | --- |
| 1 | `audit_skill.py` 對八個 skill 全 exit 0 | 含新的 `bdd-formulation` |
| 2 | `git log --follow` 穿得過搬移 | 五份 reference ＋ 13 個 examples ＋ `check_spec.py`，改名自成一筆提交 |
| 3 | `status.py` 對 fixture 的數字不變 | `面向` 欄的用法沒改，數字變了就是改壞了 |
| 4 | **Known Boundaries 的每一條都要標來源** | 新檢查；先證明它會失敗 |
| 5 | `check_spec.py` 搬家後對既有 fixture 的結果不變 | FR→AC 的對應與 story slug 都不受這次影響 |
| 6 | 健康路徑仍 exit 0，新檢查在正確文件上沉默 | 偽陽性探針 |
| 7 | 每個契約規定的欄位說得出供給者 | 第七種缺陷形態的機械測試 |
| 8 | 全 repo 掃描（排除 `docs/superpowers/`、`runs/`） | 逐一判斷，歷史紀錄要說明 |
| 9 | 引文逐行比對來源 | 每個 `>` 區塊的文字要真的在被引用的檔案裡 |

---

## 已知弱點

- **前沿模型沒有機械檢查。** 「這一輪該問哪些」是判斷，不是計算。`status.py` 的面向覆蓋率
  是唯一的間接訊號——它看得出哪個面向一題都沒問，看不出順序對不對。
- **一輪問五題會不會太多**，沒有實測過。`grilling` 沒說上限；我們現行的是「一輪 3–5 題，
  但真正的判準是回答品質」。前沿可能一次算出十題。
- **`bdd-formulation` 拆出來之後，`bdd-spec` 剩 spec-format.md（28.7K）＋ persona-definition
  ＋ status.py**。那份 28.7K 的契約可能太大，但拆它不在這一版範圍。
- **八個 skill** 之後，`skill-rules` 的房規（S10 的 500 行建議）會有更多檔案踩到。

---

## 不在範圍

- `clarify-loop` 與使用者全域同名 skill 的衝突
- `spec-format.md` 的拆分
- `runs/` 底下既有產物的遷移
- 實測（`personal-memo` case 需要未被污染的 session）
