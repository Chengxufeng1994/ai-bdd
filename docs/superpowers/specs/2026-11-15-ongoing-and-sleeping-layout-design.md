# 進行中放 `specs/ongoing/`，現況放 `docs/`

**狀態**：已批　**日期**：2026-11-15

## 意圖

需求方的話：`specs/<date>-<feature>/` 改成 `specs/intent/<date>-<feature>/`；「這可能只是還在計劃
階段的，等到全部完成，才會更新」到一份「整個專案的真實來源」——`docs/` 底下的
`CONTEXT.md`、`DOMAIN-MODEL.md`、`DATABASE-SCHEMA.md`、`TECH.md`。做完的 feature 寫一份紀錄到
`specs/archive/`，`intent/` 留著。目錄名後來定為 `specs/ongoing/`（原提的 `intent`）與 `specs/sleeping/`
（原提的 `archive`）；再後來改成搬家：做完的目錄整個搬進 `sleeping/`，「有一天想修正，可以搬回
`ongoing`」；交付紀錄 `delivered.md` 也拿掉了，見 D2。

現在一個 feature 做完，產物就停在 `specs/<date>-<feature>/`。下一個 feature 要知道系統現在長怎樣，
得自己翻每一個舊目錄，而舊目錄裡的規則可能已經被後來的 feature 改掉。

先例是 OpenSpec（`Fission-AI/OpenSpec`）：`openspec/changes/` 放進行中的變更，`openspec/specs/`
放現況，一個變更做完就 archive，把它的增修併進現況。差別在於這裡的現況放專案自己的 `docs/`。

## 決定

### D1. 三個地方，三種壽命

```
specs/
├── ongoing/<date>-<feature>/       還在做的：discovery.md、spec.md、splitting.md、acceptance/、plan.md、tasks/
└── sleeping/<date>-<feature>/      做完的：整個目錄 git mv 過來，內容不變

docs/                               系統現在長怎樣——專案的真實來源
├── CONTEXT.md                      詞彙（已有，規則不變）
├── DOMAIN-MODEL.md                 聚合與不變條件（取代 specs/domain-model.md）
├── DATABASE-SCHEMA.md              資料表與欄位
└── TECH.md                         Tech Stack、Commands、Project Structure、Code Style、Boundaries

features/<story-slug>.feature       行為的現況（不動——跑得過，所以不會過期）
```

| 問題 | 去哪找 |
| --- | --- |
| 現在還有哪些 feature 在做 | `ls specs/ongoing/` |
| 這個 feature 當時為什麼這樣決定 | 它的目錄（`ongoing/` 或 `sleeping/`）裡的 `discovery.md`、`spec.md` |
| 這個 feature 最後交付了什麼 | `specs/sleeping/<date>-<feature>/spec.md`（最後一版）與它的票檔 |
| 系統現在的詞彙、聚合、資料表、技術現況 | `docs/` |
| 系統現在的行為 | `features/*.feature` |

一個 feature 在任何時刻只住一個地方：`ongoing/` 或 `sleeping/`。所有搬動都用 `git mv`，
`git log --follow` 才追得到搬家前後的歷史。

### D2. 不另寫交付紀錄

做完的目錄裡已經有這個 feature 交付了什麼，不另寫一份：

| 想知道 | 在哪 |
| --- | --- |
| 最終的 story 與 FR／AC、技術決定 | `spec.md` 的最後一版 |
| 跟原本規格不一樣的地方 | `spec.md` 的 `### 版本修訂歷史` |
| 實作時定下的暫定行為、用哪個測試釘住 | 各票檔的 `## 缺口` |
| 行為 | `features/*.feature` |

再寫一份就是同一個事實兩份，遲早不一樣。代價是暫定行為散在各票檔，沒有一張總表。

### D3. `docs/` 四份怎麼更新

| 檔 | 什麼時候寫 | 從哪來 |
| --- | --- | --- |
| `docs/CONTEXT.md` | **照舊**：`spec` 寫規格時新建或追加 | `discovery.md` 已答的定義——`formulation` 寫 `.feature` 時就要用到，不能等做完 |
| `docs/DOMAIN-MODEL.md` | feature 做完時 | `spec.md` 的聚合與不變條件 |
| `docs/DATABASE-SCHEMA.md` | feature 做完時 | `spec.md` Implementation Decisions 的 Schema |
| `docs/TECH.md` | feature 做完時 | `spec.md` 技術部分的 Tech Stack 到 Boundaries，只改有變的條目 |

`docs/` 只放**做出來了的**東西：還沒實作的聚合與資料表不該先出現在現況裡，所以後三份等 feature
做完才寫。`CONTEXT.md` 是唯一例外，理由在表裡。

**專案自己已經有同名的檔時**，只增修標了 ai-bdd 維護的區塊，其他內容一律不動：

```markdown
<!-- ai-bdd:begin Commands -->
- 閘門：`make verify` ← 專案 Makefile
<!-- ai-bdd:end Commands -->
```

檔不存在就新建，整份都是 ai-bdd 的區塊。專案自己的架構文件（實測場的 `docs/ARCHITECTURE.md`）
照舊只讀不寫。

### D4. 新 skill `archive` 負責最後一步

一個 feature 的票全部合回、它涵蓋的 `.feature` 都是 `@review` 之後跑，一次一個 feature：

1. 更新 `docs/DOMAIN-MODEL.md`、`docs/DATABASE-SCHEMA.md`、`docs/TECH.md`（D3）；每份檔尾的修訂紀錄
   記一列：哪個 feature、改了什麼
2. `git mv specs/ongoing/<date>-<feature> specs/sleeping/<date>-<feature>`——搬家跟第 1 步同一個 commit

### D5. 喚醒：從 `sleeping/` 搬回 `ongoing/`

已經做完的 feature 要修正（規則改了、發現交付有錯、要加一條）時，不開新目錄，把它搬回來：

1. `git mv specs/sleeping/<date>-<feature> specs/ongoing/<date>-<feature>`，目錄名不改——日期是它出生的
   日子，不是這次修正的日子
2. `discovery` 把 `discovery.md` 改回 `澄清中`、`### 版本修訂歷史` 記一列「喚醒：<為什麼>」，照現在的
   續跑規則問新的題
3. 之後照現有的再跑規則走：`spec` 修補（FR／AC 編號不重排、改動記修訂歷史）、`formulation`〈再跑：
   `spec.md` 改過之後〉、`plan`、`task-breakdown` 重切（有勾的舊票不動，新票接著編號）
4. 「上一次交付了什麼」看搬回來那個 commit 的 `spec.md`（`git show <commit>:<路徑>`）；修正完照 D4
   更新 `docs/`、再搬回 `sleeping/`

由誰搬：`discovery` 被叫來處理一個 `sleeping/` 底下的 feature 時，先問使用者要不要喚醒——搬回來是
使用者的決定，不是 skill 自己判斷「這個需要改」。

**不喚醒的情況**：要做的是一個**新的能力**、只是碰到同一批程式碼——開新的 `ongoing/<date>-<feature>/`，
做完時 `archive` 照樣更新 `docs/`。判準：這次的需求是改舊 feature 的 FR／AC，還是長出新的 story？

不訪談、不改需求、不改程式碼：只搬已經講定、已經做完的東西。規則跟 `docs/` 現況衝突時（例如
`TECH.md` 說 Go 1.23、這個 feature 的 spec 說 1.24）不自己選，回報給使用者。

### D6. 上游跟著改

- **路徑**：所有 skill 的 `specs/<date>-<feature>/` 改成 `specs/ongoing/<date>-<feature>/`
- **`spec`**：
  - 技術部分的 Tech Stack 到 Boundaries 改從 `docs/TECH.md` 抄來再確認，標
    `← 專案 docs/TECH.md`；沒有 `docs/TECH.md`（第一個 feature）時照舊從專案檔讀
  - 不再寫 `specs/domain-model.md`：聚合與不變條件寫在 `spec.md` 裡（新的一節），做完時由
    `archive` 併進 `docs/DOMAIN-MODEL.md`；寫之前讀 `docs/DOMAIN-MODEL.md` 知道現況
- **`status.py`、`check_spec.py`**：找 `specs/ongoing/*/`（`sleeping/` 不算進度）；還是舊佈局（`specs/<date>-<feature>/`、
  `specs/domain-model.md`）時 exit 1 並印出要跑的 `git mv`——跟以前處理 `prd.md` 佈局一樣，不兩種
  都支援
- **`CLAUDE.md` 的規矩**：「skill 只准新建或追加 `docs/CONTEXT.md`」改成「`spec` 新建或追加
  `docs/CONTEXT.md`；`archive` 新建或增修 `docs/DOMAIN-MODEL.md`、`DATABASE-SCHEMA.md`、`TECH.md`
  裡 ai-bdd 維護的區塊；`docs/` 其他檔不碰」

## 分兩期做

| 期 | 做什麼 | 動到 |
| --- | --- | --- |
| **一** | 搬家：`specs/<date>-<feature>/` → `specs/ongoing/<date>-<feature>/`；兩支腳本認新路徑、舊佈局報 `git mv` | 每個 skill 的路徑、兩支腳本與測試、README／CLAUDE.md／PLAN.md |
| **二** | `archive` skill、`specs/sleeping/`、`docs/` 三份新檔與區塊標記、喚醒；`spec` 改讀 `docs/TECH.md`、`domain-model.md` 退場 | 新 skill、`spec`、`discovery`（喚醒）、`domain-model-format.md`、CLAUDE.md 的 `docs/` 規矩 |

第一期是純搬家，行為不變，可以先做先驗；第二期才是新能力。

## 非目標

- 不刪任何 feature 目錄——做完搬到 `sleeping/`，不刪
- 不改 `.feature` 的位置與狀態 tag
- 不做 REVIEW／VERIFY；`archive` 不判斷實作對不對，只搬已經合回的東西
- 不寫專案自己的架構文件

## 考慮過的做法

| 做法 | 取捨 | 結果 |
| --- | --- | --- |
| A. 只改路徑成 `specs/ongoing/` | 最小；但「做完之後更新真實來源」沒有人做 | 只當第一期 |
| B. 現況放 `specs/current/`（上一版提案） | skill 不碰 `docs/`；但專案的人找文件不會去翻 `specs/` | 捨棄——需求方要現況住 `docs/` |
| C. 現況放 `docs/` 四份，做完的紀錄放 `specs/sleeping/`，`ongoing/` 的目錄留著 | 不用搬；但 `ongoing/` 越積越多，看不出哪些還在做 | 捨棄——需求方要做完的能搬回來修正 |
| **D. 現況放 `docs/` 四份，做完的目錄整個搬到 `specs/sleeping/`，要修正再搬回 `ongoing/`**（本文件） | `ongoing/` 一眼看出還在做什麼、修正有固定的路；`docs/` 的規矩要放寬，同名檔要用區塊標記隔開，每次搬動都得用 `git mv` | **建議** |

## 已定

1. **`archive` 在票全部合回就跑**（需求方 2026-11-15）——REVIEW／VERIFY 還沒實作，不等它們
2. **先只做第一期**（需求方 2026-11-15）：純搬家，第二期另行開工

## 最不確定的地方

- **區塊標記會被人手改壞。** 有人在 `<!-- ai-bdd:begin -->` 區塊裡手寫了東西，下一次 `archive` 會
  覆蓋掉；提議 `archive` 先比對區塊內容跟上一次寫的是否一樣，不一樣就停下回報
- **喚醒期間 `docs/` 已經是上一次交付的樣子。** 修正中的 feature 改了某條規則，`docs/` 要等它再做完才
  更新；這段期間別的 feature 讀到的是舊現況。這跟一般 feature 進行中的情況一樣，但修正通常是在改已經
  有人依賴的東西，要寫明
- **`domain-model.md` 從「寫規格時當場增修」改成「做完才寫」**：兩個並行的 feature 寫規格時都看
  不到對方新增的聚合。現在的 `specs/domain-model.md` 也只在 `spec` 寫完時才更新，差別不大，但要寫明
