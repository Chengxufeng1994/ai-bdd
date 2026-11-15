---
name: bdd
description: >
  在一個已經寫好的 `.feature`（或它的一段）裡，用 red／green／refactor 三個入口之一，一次
  一個 slice 推進實作：red 只建立對應這個 slice 的失敗訊號、green 只補目前 failure 需要的
  最小程式、refactor 只在綠燈保護下整理。不產 `.feature`、不改場景文字、不翻狀態 tag、
  不跨出指定範疇；規格有洞就停下回交。可由 `bdd-implement` 的一張票委派，也可單獨使用。
  觸發詞：「先讓它紅」「讓這個場景紅」「讓這個場景變綠」「紅綠重構」「bdd red」「bdd green」
  「bdd refactor」「這個 scenario 做到綠」「在綠燈下整理 step definition」「一個 slice 一個 slice 做」。
  English: red green refactor, make this scenario fail first, make this scenario pass,
  refactor under green, drive one scenario at a time, bdd red, bdd green, bdd refactor.
---

# BDD — 一個 slice，一個入口

| | |
| --- | --- |
| **職責** | 在指定的 `.feature` 範疇裡，用 red／green／refactor 其中一個入口，一次一個 slice 把場景推進到那個入口的完成條件 |
| **讀** | 指定的 `.feature`（`@ready` 或 `@wip`）、`plan.md` 的技術節、既有的 step definition／測試／產品程式碼；由票委派時加上那張票檔 |
| **寫** | step definition、測試、產品程式碼。**不寫 `specs/` 底下任何檔案** |
| **上游** | `bdd-implement`（一張票委派一個範疇與一個入口），或使用者直接指定 |
| **下游** | 回報給呼叫它的人，並建議下一次從哪個入口繼續 |
| **不做** | 不產或改 `.feature`、不改 `spec.md`／`plan.md`、不翻狀態 tag、不勾票、不跨出指定範疇、一次不進兩個入口 |

**每次只進一個入口，同一時刻只處理一個 slice。** 紅燈是用來驅動的，不是事後補的：沒看過
對應這個 slice 的失敗，就沒有理由寫下一行產品程式碼。

## 使用時機

- `.feature` 寫好了，要讓某個場景先紅、再綠、再整理
- `bdd-implement` 做一張票時，把其中一個場景的紅／綠／整理交給本 skill
- 場景已經綠了，想在不改行為的前提下整理 step definition 與支撐層

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| 還沒有 `.feature` | `bdd-formulation` |
| 要做完**一整張票**：標 `@wip`、跑專案閘門、勾驗收條件、翻 `@review`、合回主線 | `bdd-implement`——它管票，本 skill 管票裡的一個 slice |
| 要改場景的文字、例子或規則 | 回 `bdd-discovery` 開一題；本 skill 不碰 `.feature` |
| 要決定 seam、模組、契約 | `bdd-plan` |
| 要決定先做哪張票 | `bdd-tasks` |

跟 `bdd-implement` 的分界：**它回答「這張票做完了沒」，本 skill 回答「這個 slice 在這個入口
做到位了沒」。** 由票委派時，只做那張票授權的範疇與入口，不順手做下一個。

## 參考檔案

- `rules/slice.md` — 什麼是一個 slice、先做哪一個、什麼時候算做完、Scenario Outline 怎麼算
- `rules/red.md` — `red` 入口。**進 red 之前讀**
- `rules/green.md` — `green` 入口。**進 green 之前讀**
- `rules/refactor.md` — `refactor` 入口。**進 refactor 之前讀**
- `../bdd-formulation/references/step-grammar.md` — 步驟的封閉文法；新步驟要照它寫

## 前置確認（Ask user）

Ask user: 以下兩項不明確就先問，不要自己挑。

1. **範疇**：哪一個 `.feature`，或它的哪一段（一個 `Rule:`、幾個 `@example-N.M`）。
   範疇決定「順手多做」的界線在哪——沒有它，每個相鄰的場景看起來都該一起做
2. **入口**：`red`、`green` 還是 `refactor`。三者的完成條件不同，猜錯入口等於做錯事

Check if: 指定的 `.feature` 存在，而且狀態 tag 是 `@ready` 或 `@wip`。是 `@draft` 就停下——
那份規格需求方還沒點頭，照它做出來的行為沒有人同意過。

---

## 流程

| 步驟 | 做什麼 | 讀哪個 rule |
| --- | --- | --- |
| 1 | 對齊範疇與入口 | — |
| 2 | 盤點測試入口與既有抽象 | — |
| 3 | 選下一個 slice | `rules/slice.md` |
| 4 | 照入口把這個 slice 做到完成條件 | `rules/red.md`／`green.md`／`refactor.md` |
| 5 | 還有下一個 slice 就回到 3，否則停 | `rules/slice.md` |
| 6 | 回報 | — |

### 1. 對齊範疇與入口

讀指定的 `.feature` 範疇、`plan.md` 的 `## Implementation Decisions` 與 `## Testing Decisions`
（seam 在哪一層、每個例子在哪一層測）。由票委派時再讀票檔的 `涵蓋場景`——它就是範疇。

### 2. 盤點測試入口與既有抽象

動手之前先看專案現在有什麼：

| 盤點什麼 | 為什麼 |
| --- | --- |
| **最窄的重跑命令**——能只跑一個場景嗎？一條 `Rule`？一個檔？ | 回饋越窄，失敗越好歸因。實測場：`go test -count=1 ./test/acceptance/ -godog.tags=@example-2.1` |
| 既有的 step definition（實測場 `make steps`） | 同一句話兩種寫法就是兩支函式，而它們會分岔 |
| 既有的 fixture、測試替身、建態的 helper | `Given` 要擺的狀態，多半已經有人擺過 |
| 測試與產品程式碼的落點 | 沿用既有結構，不另起一套平行目錄 |

SHOULD NOT: 把「先跑整包比較保險」當預設。跑不到單一場景時，選最小可接受的範圍，並在回報
裡寫明用了哪個範圍、為什麼不能更窄。

### 3. 選下一個 slice

照 `rules/slice.md`：一個 slice 是一個 `Example:`（一個 `@example-N.M`）；範疇內先做能形成
最小、最獨立回饋的那一個。

### 4. 照入口推進

| 入口 | 起點 | 可以動 | 完成條件 |
| --- | --- | --- | --- |
| `red` | 這個 slice 還沒有對應的失敗 | step definition 骨架、fixture、斷言、讓測試跑得起來的最小通路 | 窄範圍重跑，看到**對應這個 slice** 的失敗：`step is undefined` 或斷言失敗，不是 parse error 或別處的錯 |
| `green` | 這個 slice 現在是紅的，而且看過那個失敗 | 讓目前這個失敗消失所需的最小程式 | 這個 slice 單獨綠，範疇內原本綠的仍然綠 |
| `refactor` | 範疇內全綠 | 重複的 step、型別轉換的邊界、支撐層的命名與結構 | 整理前後各跑一次，結果相同且全綠 |

MUST: 入口的起點不成立就不要開始，改回報該先進哪個入口。`green` 沒有紅燈可依、
`refactor` 不在綠燈下，做出來的東西都沒有任何測試在保護。

### 5. 下一個 slice，或停

當前 slice 達到完成條件、範疇內還有下一個、而且沒有碰到下表任何一項，就回到第 3 步，
**用同一個入口**做下一個。

| 碰到 | 怎麼辦 |
| --- | --- |
| 場景跟規則打架、少一個邊界、`Then` 推不出來 | **停**。回報是哪個場景、缺什麼；由票委派時請 `bdd-implement` 記進票檔的 `## 缺口` |
| 要讓它綠，得改範疇外的場景或別張票的行為 | **停**，回報耦合在哪 |
| 失敗訊號混在一起、分不出是哪個 slice 的 | **停下修紅燈的品質**，不要再鋪下一個 |
| 該換入口了（紅完要綠、綠完要整理） | **停**，在回報裡建議下一個入口——換入口是下一次呼叫的事 |

MUST NOT: 為了讓場景綠而改 `.feature` 的任何文字。場景是需求方同意過的驗收條件；
改它來換進度，等於沒有人同意過那條規格。

### 6. 回報

見〈完成後〉。

---

## 產物

只有程式：step definition、測試、產品程式碼，落在專案既有的結構裡。

MUST: 不寫 `specs/` 底下任何檔案、不改 `.feature`。狀態 tag、票檔的勾選與 `## 缺口` 都是
`bdd-implement` 的——本 skill 只回報，由它記。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 這次只進了一個入口；做過的每個 slice 都在指定範疇內
- [ ] `red`：每個 slice 都用窄範圍重跑看過對應它的失敗，命令與失敗訊息在回報裡；沒有補產品行為
- [ ] `green`：每個 slice 的補碼都從看過的那個失敗出發；每個 slice 各自有一次獨立的綠燈，命令與結果在回報裡
- [ ] `refactor`：整理前後各跑一次且結果相同；沒有改變任何場景驗的行為
- [ ] 寫新步驟前查過既有的 step definition，沒有同一句話兩種寫法
- [ ] `.feature` 的文字、狀態 tag、`spec.md`、`plan.md`、票檔都沒有被改
- [ ] 碰到規格的洞或範疇外的耦合時停下來回報了，沒有自己決定

## 完成後

告訴對方：

1. 進的是哪個入口、範疇是什麼
2. 做了哪幾個 slice，各自的重跑命令與結果（紅的貼失敗訊息，綠的貼通過）
3. 動了哪些檔（step definition／測試／產品程式碼各一句）
4. 停下來的原因：做完了、碰到上表哪一項、或該換入口
5. 建議下一次從哪個入口、哪個 slice 繼續
