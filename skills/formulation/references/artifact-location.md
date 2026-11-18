# `.feature` 放哪裡

讀這份的時機：在一個新專案第一次跑這個 skill。同一個專案只需要決定一次。

## 規則

BDD 流程的**文件型產物住在 `specs/`**，有兩個例外：`.feature`，以及 `docs/` 底下的現況檔
（`CONTEXT.md` 等，feature 做完時由 `archive` 寫）。

### `.feature`

它同時是規格與可執行測試，**位置由測試框架決定**，不由本流程決定。

| 規則 | 強度 |
| --- | --- |
| `.feature` 放測試框架找得到的地方 | `MUST` |
| 不修改專案既有的任何 `.feature` | `MUST` |
| 已經有 `.feature` 就跟隨它們的位置 | `SHOULD` |
| 沒有任何 `.feature`，也還沒有框架 → 專案根的 `features/` | `SHOULD` |

技術上可以用 godog 的 `Paths` 把 `.feature` 也塞進 `specs/`，但那會讓
Go 開發者在預期的地方找不到它。**鏈的完整性靠 tag 與 slug 維繫，不是靠同目錄。**

### `docs/` 的現況檔

系統現在長怎樣，只放做出來了的東西。本 skill **只讀** `docs/CONTEXT.md` 取措辭，不寫 `docs/` 底下任何檔案；
寫入者是 `archive`，規則見 `../../archive/references/docs-format.md`。

