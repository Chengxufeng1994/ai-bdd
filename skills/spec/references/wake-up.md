# 喚醒——要改的 feature 已經在 `specs/sleeping/`

做完的 feature 整個目錄住在 `specs/sleeping/<date>-<feature>/`。這次的需求是要**改它的 FR／AC**
（規則改了、交付有錯、加一條例外）時，不開新目錄，把它搬回來；這次的需求是**新的 story**、只是碰到
同一批程式碼時，照常開新的 `specs/ongoing/<date>-<feature>/`。

喚醒歸 `spec`，因為它是 `spec.md` 唯一的寫入者：搬回來的那一刻狀態就要改回 `澄清中`，否則問答進行中
`spec.md` 還寫著 `已對焦`，`formulation` 或 `plan` 會把它當成講定的規格。

Ask user: 要不要喚醒 `<date>-<feature>`。搬回來是使用者的決定——不是本 skill 判斷「這個該改」。

| 回答 | 怎麼做 |
| --- | --- |
| 要 | 1. `mkdir -p specs/ongoing && git mv specs/sleeping/<date>-<feature> specs/ongoing/<date>-<feature>`（目錄名不改：日期是它開始的日子；全部歸檔後 git 不留空的 `specs/ongoing/`，所以先建）<br>2. `spec.md` 的狀態改 `澄清中`，`### 版本修訂歷史` 記一列「喚醒：<需求方的原話，逐字>」<br>3. STOP，告訴使用者跑 `/discovery`：它在 `interview-log.md` 追加一場 `### <日期> · 喚醒` Session，把這次要改什麼開成新題（`**來自**`：喚醒 —「<同一句原話>」）再問 |
| 不要 | 告訴使用者跑 `/discovery` 開新的 feature |

喚醒之後全部走現成的再跑規則：`spec` 讀紀錄的新題修補（編號不重排、記修訂）、`formulation`〈再跑：`spec.md` 改過之後〉、
`plan` 重寫被規格推翻的 Risks 與接合點、`task-breakdown` 重切（有勾的舊票不動）、`implement` 做新票。做完由 `archive` 再寫一次現況、搬回 `specs/sleeping/`。
「上一次交付了什麼」看搬回來那個 commit 的 `spec.md`（`git show <commit>:<路徑>`）。
