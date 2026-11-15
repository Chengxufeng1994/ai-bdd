---
name: archive
description: >
  一個 feature 的票全部合回之後收尾：把 `spec.md` 講定、已經做出來的東西併進專案的真實來源——
  `docs/DOMAIN-MODEL.md`（聚合與不變條件）、`docs/DATABASE-SCHEMA.md`（資料表）、`docs/TECH.md`
  （Tech Stack、Commands、Project Structure、Code Style、Boundaries）——再把整個
  `specs/ongoing/<date>-<feature>/` 用 `git mv` 搬進 `specs/sleeping/`。只動 ai-bdd 維護的區塊，
  不訪談、不改需求、不改程式碼。BDD 六步流程 IMPLEMENT 之後的收尾。
  觸發詞：「歸檔」「archive」「這個 feature 做完了」「票都合回了，下一步」「更新 docs」
  「把現況寫回文件」「收尾」「搬到 sleeping」「更新 DOMAIN-MODEL」「更新 TECH.md」。
  English: archive this feature, the feature is done, update the docs, write back the current
  state, move it to sleeping.
---

# ARCHIVE — 做完了，寫回現況

| | |
| --- | --- |
| **職責** | 把一個做完的 feature 併進專案的真實來源：`docs/` 三份現況檔照它增修，再把它的目錄搬進 `specs/sleeping/` |
| **讀** | `specs/ongoing/<date>-<feature>/` 的 `spec.md`（`## Domain Model` 與技術部分）、`tasks/`（全部勾完）、涵蓋的 `.feature`（`@review`）、`docs/` 既有的三份檔 |
| **寫** | `docs/DOMAIN-MODEL.md`、`docs/DATABASE-SCHEMA.md`、`docs/TECH.md` 裡 ai-bdd 維護的區塊；`git mv` 這個 feature 的目錄到 `specs/sleeping/` |
| **上游** | `implement`（讓最後一張票完成的那次合回之後） |
| **下游** | 無；要修正時由 `discovery` 喚醒，搬回 `specs/ongoing/` |
| **不做** | 不訪談、不改 `spec.md`／票檔／`.feature`、不改程式碼、不寫 `docs/CONTEXT.md`、不動區塊外的內容、不判斷實作對不對 |

**只搬已經講定、已經做出來的東西。** `docs/` 是系統現在長怎樣；這一步不做任何新決定。
`spec.md` 跟 `docs/` 的現況衝突時（`TECH.md` 說 Go 1.23、這次的 spec 說 1.24），不自己選——停下回報。

## 使用時機

- 一個 feature 的票全部合回主線、涵蓋的 `.feature` 都是 `@review`，要把它寫回現況
- 一個被喚醒、修正完的 feature 又做完了，要再寫一次現況、搬回 `specs/sleeping/`

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| 還有票沒勾完、或某個 `.feature` 還是 `@wip`／`@ready` | `implement` 把票做完 |
| 要改詞彙 | `spec`（寫規格時就追加 `docs/CONTEXT.md`） |
| 要修正一個已經在 `specs/sleeping/` 的 feature | `discovery`（先問要不要喚醒） |
| 要判斷實作是不是真的滿足規格 | REVIEW／VERIFY（未實作） |

## 參考檔案

- `references/docs-format.md` — `docs/` 三份檔的骨架、區塊標記、修訂紀錄。**改 `docs/` 之前先讀**
- `references/domain-model-format.md` — `docs/DOMAIN-MODEL.md` 與 `spec.md` `## Domain Model` 共用的骨架
- `../spec/references/tech-format.md` — 技術部分各節的格式，`docs/TECH.md` 照它

## 前置確認（Ask user）

Check if: 這個 feature 在 `specs/ongoing/` 底下；`tasks/index.md` 列的每張票驗收條件都勾完；它涵蓋的
每個 `.feature` 都是 `@review`——涵蓋哪些檔，看各票 `涵蓋場景` 寫的 `.feature` 路徑的聯集（不屬於任何票的，
例如 walking skeleton 的 `version.feature`，不算）。任何一項不成立就停下，說還差哪張票——還沒做完的東西寫進現況，
`docs/` 就不再是系統真正的樣子。

---

## 流程

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| 1 | 先查區塊有沒有被手改過，再讀 `spec.md` 與 `docs/` 現況，列出要改的條目 | 手改過、沒有區塊的舊檔、或衝突時要 |
| 2 | 增修 `docs/` 三份檔的區塊，記修訂紀錄 | — |
| 3 | `git mv` 到 `specs/sleeping/`，跟第 2 步同一個 commit | — |

### 1. 讀，列出要改的

**先查手改，再比 `spec.md`。** 手改過的區塊跟 `spec.md` 不一樣，照字面比會被列成「改回 spec」——
靜靜地丟掉別人的修改。查法見 `references/docs-format.md`〈手改過沒有〉；兩件事同時成立時以手改為準：

| 碰到 | 怎麼辦 |
| --- | --- |
| 區塊被手改過 | **停**，列出差異問使用者要保留哪個。答「保留手改」：區塊照現在的樣子不動，`修訂紀錄` 記一列「保留手改：<區塊>（<手改的 commit>）」——這次的 archive commit 因此碰到這份檔，成為下一次的比對基準，下一次不會在同一行再停。`spec.md` 不改：下一個 feature 的技術部分從 `docs/TECH.md` 抄，自然跟上 |
| 檔存在、一個 ai-bdd 區塊都沒有，內容卻是舊版 ai-bdd 寫的（從 `specs/domain-model.md` 搬來的 `DOMAIN-MODEL.md`） | **停**，問要不要先把每個聚合包進一個區塊（區塊名＝聚合名）、檔頭加「由 `archive` 維護」那一行，單獨一個 commit。不包就照「沒有這個區塊」在檔尾加——同一個聚合會寫兩次 |

再填這張表：

| `docs/` 的檔 | 從 `spec.md` 的哪裡來 | 怎麼比 |
| --- | --- | --- |
| `DOMAIN-MODEL.md` | `## Domain Model` | 新的聚合加進去；改了既有聚合的，`spec.md` 那一節寫了「原本是什麼、為什麼改」 |
| `DATABASE-SCHEMA.md` | `## Implementation Decisions` 的 Schema | 新表加區塊；改了既有表的，照新的寫 |
| `TECH.md` | `## Tech Stack` 到 `## Boundaries` | 只改跟現況不一樣的條目——多數 feature 這五節跟 `TECH.md` 一字不差，那就不動。`## Commands` 的 `替代閘門（本環境）` 不寫進去：它是那次環境的變通，不是專案的現況 |

把要改的條目先攤成一張表（檔、區塊、原本、改成、來源），再動手。

### 2. 增修 `docs/`

照 `references/docs-format.md`〈區塊標記〉：檔不存在就新建，存在就只改 ai-bdd 的區塊、沒有就在檔尾
加。每份改過的檔在 `修訂紀錄` 區塊追加一列。

MUST: 只改 `<!-- ai-bdd:begin … -->` 到 `<!-- ai-bdd:end … -->` 之間的內容。區塊外是專案自己的文件，
有人在維護；改到它就是替別人做決定。唯一例外是第 1 步使用者答應的「把舊檔包成區塊」：只加標記與檔頭那一行，
內容一字不改。

MUST: 區塊的內容跟上一次 `archive` 寫的不一樣（有人手改過）時停下回報，不覆蓋——第 1 步查過了，
這裡只照使用者的答案做。

| 碰到 | 怎麼辦 |
| --- | --- |
| `spec.md` 跟 `docs/` 現況衝突（同一條，兩邊說法不同，而 `spec.md` 沒有記這是一次改動） | **停**，列出兩邊原文問使用者——可能是 spec 抄錯，也可能是現況過期 |
| `spec.md` 的 `## Domain Model` 寫「（無）」、或某節寫「（無）」 | 那份檔不動，回報裡說明 |

### 3. 搬進 `specs/sleeping/`

```bash
git mv specs/ongoing/<date>-<feature> specs/sleeping/<date>-<feature>
```

MUST: 用 `git mv`，不用複製再刪——`git log --follow` 靠它追到搬家前的歷史；喚醒時搬回來也一樣。

目錄裡的東西一個字都不改：`discovery.md`、`spec.md`、`plan.md`、`tasks/` 就是這個 feature 交付了什麼、
為什麼這樣決定的紀錄，不另寫一份交付紀錄。

第 2 步與第 3 步放在同一個 commit：`docs/` 的現況與 feature 搬進 `sleeping/` 是同一件事，分開提交
中間就有一個「現況已更新、feature 卻還在進行中」的狀態。

MUST: commit 標題固定寫 `docs: archive <date>-<feature>`——下一次查手改時，靠它認出哪個 commit 是
`archive` 寫的。

`docs/` 不必改時（修正只動了 FR 與 Implementation Decisions 裡不進 `docs/` 的部分），這個 commit 只有
搬家；這次交付了什麼的紀錄就是 `spec.md` 的 `### 版本修訂歷史`，不為了留痕去改 `docs/`。

## 產物

| 產物 | 格式 |
| --- | --- |
| `docs/DOMAIN-MODEL.md` | `references/domain-model-format.md` |
| `docs/DATABASE-SCHEMA.md`、`docs/TECH.md` | `references/docs-format.md` |
| `specs/sleeping/<date>-<feature>/` | 從 `specs/ongoing/` `git mv` 過來，內容不變 |

MUST NOT: 寫 `docs/` 底下這三份以外的任何檔，或改 `specs/` 底下任何檔的內容。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 開工前每張票都勾完、涵蓋的 `.feature` 都是 `@review`
- [ ] `docs/` 只改了三份檔、只改了 ai-bdd 區塊裡的內容；區塊外的 diff 是空的
- [ ] 每條寫進 `docs/` 的內容都指得出來自 `spec.md` 的哪一節、保留了來源
- [ ] 改過的檔都在 `修訂紀錄` 追加了一列
- [ ] 先查了手改才填表；衝突、手改過的區塊、沒有區塊的舊檔都停下來問了，沒有自己選
- [ ] commit 標題是 `docs: archive <date>-<feature>`
- [ ] feature 目錄是 `git mv` 搬的，目錄裡的檔一個字都沒改；`docs/` 的改動與搬家在同一個 commit
- [ ] 沒有寫 `docs/CONTEXT.md`、沒有改程式碼、`spec.md`、票檔或 `.feature`

## 完成後

告訴對方：

1. `docs/` 改了哪幾份、哪些區塊；哪幾份沒動、為什麼
2. 有沒有停下來問的衝突或手改過的區塊
3. 這個 feature 已搬到 `specs/sleeping/<date>-<feature>/`；之後要修正它，跑 `discovery` 喚醒
4. `docs/` 只放結構（聚合、資料表、技術現況）；行為的現況——規則怎麼算——住在 `features/*.feature`
   與 `specs/sleeping/` 裡那份 `spec.md`
