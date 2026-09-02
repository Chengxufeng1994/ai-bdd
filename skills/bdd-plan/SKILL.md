---
name: bdd-plan
description: >
  SPEC 說要做成什麼，PLAN 說先做哪一塊——把一份 `spec.md` 與它的 `.feature`
  切成一疊 tracer bullet：每張票穿過所有層、自己可 demo、塞得進一個全新 session，
  並宣告它被哪幾張票 block。BDD 六步流程的 PLAN 步驟。
  觸發詞：「這批要怎麼切」「拆成幾張票」「先做哪一個」「排實作順序」
  「切成小步驟」「spec 寫完了下一步」「這些票誰卡誰」「開工順序」。
  English: break this spec into tickets, slice the work, what order should I
  build these, vertical slices, what can I start first.
---

# PLAN — SPEC 說要做成什麼，PLAN 說先做哪一塊

一份就緒的 `spec.md` 與它涵蓋的 `.feature` 進來，出去的是 `plan.md`：切成幾張
tracer-bullet 票、每張穿過所有層、彼此怎麼 block、先做哪一張。

不設計。API、domain 型別、schema、seam 都已經在 `spec.md` 裡定案——那是
`bdd-spec` 的工作。

## 使用時機

- `spec.md` 與 `.feature` 都寫好了，要決定切成幾塊、先做哪一塊
- 一批工作大到一個 session 做不完，需要拆給多個 session 接力

## Skill Boundaries

- 還沒有 `spec.md` → 先跑 `bdd-spec`；規則還沒定案 → 回 `bdd-clarify`
- **要決定 API、domain 型別、schema、seam → 那是 `bdd-spec` 的工作**，不是這裡
- 整批塞得進一個 context window → 不需要這一步，直接 `bdd-implement`
- 要實際寫 step definition 與產品程式碼 → `bdd-implement`

---

## 只切，不設計

MUST NOT: 在這一步決定任何 API、型別、欄位、資料表或 seam。
**它們已經在 `spec.md` 裡了**；這裡再決定一次，只會產生一份跟規格不一樣的規格。

判準：指著票裡的任何一句話，說得出它來自 `spec.md` 的哪一節、
或哪個 `@example-N.M` 嗎？說不出來就是這一步僭越了。

`spec.md` 沒寫到而切票時發現需要的，寫進 `## 缺口` 回報，不要順手補。

---

## 流程

### 1. 讀

`docs/bdd/<slice-slug>/spec.md` 全部，加上它涵蓋的每個 `.feature`。

**一次讀完整批。** 票之間的依賴只有並排時才看得見；分批切會各自長出一套順序，
然後在實作時發現它們互相卡住。

順便看一眼專案既有的結構——找 prefactoring 的機會。

### 2. 切

先切出**最薄的完整路徑**當第一張票：碰到每一層、邏輯最少。它證明的是接線，
不是行為。專案已經有骨架就跳過。

其餘依 `spec.md` 的 Example Mapping 切，一張票一組相關的 `@example-N.M`。
每張票問一次「做完之後我能 demo 什麼」，答不出來就是切錯了。

三條規則與 wide refactor 的例外 → `references/plan-format.md`。

### 3. 攤出來確認（Ask user）

MUST: 寫檔之前先把票的清單攤給使用者看——編號、標題、Blocked by、
交付什麼行為——並問三件事：

1. 粒度對不對（太粗／太細）
2. Blocked by 是不是真的卡住，還是只是「先做比較順」
3. 有沒有該合併或該再拆的

IMPORTANT: **切太細是這一步最常見的失效**，而且它不會自己被發現——
一疊原子化的小票看起來很專業，實際上把「這幾件事其實是同一件」的資訊丟掉了。
這個確認點就是為它存在的。

使用者確認之後才寫檔。

---

## 產物

一批一份 `docs/bdd/<slice-slug>/plan.md`。骨架照抄 `references/plan-format.md`。

MUST: 只寫入 `docs/bdd/`，**不修改專案既有的任何檔案**。

## 完成後

告訴對方：幾張票、第一張是哪一張（以及**為什麼是它**）、有沒有 prefactoring、
`spec.md` 裡推不出來的缺口有幾條。下一步 `bdd-implement <票號>`，一張票一個
全新 session。
