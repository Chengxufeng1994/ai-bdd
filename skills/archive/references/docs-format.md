# `docs/` 現況檔的骨架與區塊標記

照抄即可。讀這份的時機：第 2 步，準備改 `docs/` 的時候。

`docs/` 是專案的真實來源：系統**現在**長怎樣。只放做出來了的東西——還沒合回的 feature 不進來。

| 檔 | 裝什麼 | 從 `spec.md` 的哪裡來 |
| --- | --- | --- |
| `docs/DOMAIN-MODEL.md` | 聚合、不變條件、跨聚合的規則 | `## Domain Model`；骨架見 `domain-model-format.md` |
| `docs/DATABASE-SCHEMA.md` | 資料表與欄位、索引與約束 | `## Implementation Decisions` 的 **Schema** |
| `docs/TECH.md` | Tech Stack、Commands、Project Structure、Code Style、Boundaries | 技術部分的同名五節 |

`docs/CONTEXT.md`（詞彙）不歸本 skill：`spec` 寫規格時就新建或追加它，因為 `formulation` 寫 `.feature`
時就要用到那些詞。

---

## 區塊標記

本 skill 寫進 `docs/` 的每一段，都包在一對標記裡：

```markdown
<!-- ai-bdd:begin Commands -->
- 閘門：`make verify` ← 專案 Makefile
<!-- ai-bdd:end Commands -->
```

| 情況 | 怎麼做 |
| --- | --- |
| 檔不存在 | 新建；整份都是 ai-bdd 的區塊，第一行寫「這份由 ai-bdd 的 `archive` 維護，區塊外的內容不會被改」 |
| 檔存在、有這個區塊 | 只改這個區塊裡面 |
| 檔存在、沒有這個區塊 | 在檔尾加一個新區塊，不動既有內容 |
| 區塊的內容跟上一次寫的不一樣（有人手改過） | **停下回報**，不覆蓋——見下 |

「上一次寫的」怎麼判斷：上一次 `archive` 的 commit 裡這個區塊的內容（`git log -1 --format=%H -- <檔>`
找到最後一次改它的 commit，看那是不是 `archive` 的 commit、內容是否跟現在一樣）。不一樣代表有人在
區塊裡手寫了東西；覆蓋掉就是靜靜地丟掉別人的修改。

---

## `docs/DATABASE-SCHEMA.md`

```markdown
# Database Schema

這份由 ai-bdd 的 `archive` 維護，區塊外的內容不會被改。

<!-- ai-bdd:begin tickets -->
## tickets

| 欄位 | 型別 | 約束 | 來源 |
| --- | --- | --- | --- |
| id | uuid | PK | 2026-09-14-parking FR-1 |
| entered_at | timestamptz | NOT NULL | 2026-09-14-parking FR-1 |

**索引**: `(plate, exited_at)` 唯一——同一台車同時只能有一筆在場 ← 2026-09-14-parking FR-3
<!-- ai-bdd:end tickets -->

<!-- ai-bdd:begin 修訂紀錄 -->
## 修訂紀錄

| 日期 | feature | 改了什麼 | 為什麼 |
| --- | --- | --- | --- |
<!-- ai-bdd:end 修訂紀錄 -->
```

一張表一個區塊，區塊名就是表名。`來源` 寫 `<date>-<feature> FR-<n>`——那個 feature 的目錄搬到
`specs/sleeping/` 之後仍然找得到。記憶體儲存、還沒有資料庫的專案，整份寫「（無）——<為什麼>」。

## `docs/TECH.md`

```markdown
# Tech

這份由 ai-bdd 的 `archive` 維護，區塊外的內容不會被改。

<!-- ai-bdd:begin Tech Stack -->
## Tech Stack
<照 spec.md 技術部分的 ## Tech Stack，每條保留來源標記>
<!-- ai-bdd:end Tech Stack -->

<!-- ai-bdd:begin Commands -->
## Commands
...
<!-- ai-bdd:end Commands -->

（Project Structure、Code Style、Boundaries 各一個區塊）

<!-- ai-bdd:begin 修訂紀錄 -->
## 修訂紀錄

| 日期 | feature | 改了什麼 | 為什麼 |
| --- | --- | --- | --- |
<!-- ai-bdd:end 修訂紀錄 -->
```

一節一個區塊，區塊名就是節名。每條的格式照 `../../spec/references/tech-format.md`。

## 修訂紀錄

三份檔尾各一張表，每次 `archive` 改了這份檔就記一列：哪天、哪個 feature、改了什麼、為什麼。
「為什麼」照抄那個 feature `spec.md` 的修訂歷史或技術決定的來源，不自己編。這張表也包在區塊裡
（`修訂紀錄`），本 skill 只在表尾追加，不改既有的列。
