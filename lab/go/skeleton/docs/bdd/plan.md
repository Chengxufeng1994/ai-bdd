# 實作計畫

**來源**: `features/` 底下標 `@ready` 的全部（六份）——
`log-a-workout`、`edit-a-logged-workout`、`session-training-volume`、
`workout-history`、`week-over-week-progress`、`custom-exercise-library`
**日期**: 2026-08-19
**跳過**: 無。這一輪 CLARIFY 把兩張紅卡收斂之後，六則全部就緒。

一份涵蓋整批。每一節裡按 feature 分小節，但**節號只有一套**。

## 專案既有的東西（計畫說的是「要新增什麼」）

| 既有 | 影響 |
| --- | --- |
| `api/openapi.yaml` 只有 `GET /version` | 15 個領域端點全是新增的 |
| RFC 9457 Problem Details 已定形狀，含 `ValidationProblem` | 錯誤回應不必再設計，只要決定 `type` URI |
| 400／422 的分界已定（讀不懂／讀懂但不允許） | 見 §1 的錯誤表 |
| 分層 `domain` ← `application`（port/in、port/out）← `infrastructure`／`interfaces` | 新型別放 `internal/domain` |
| 讀寫路徑不對稱（讀不載入聚合，走自己的 port 回 view） | 三份 `@query` 不經過聚合 |
| `internal/infrastructure/` 目前是空的，沒有資料庫 | 見 §3 |
| godog `Strict: true`、`Concurrency: 4`、`Randomize: -1` | 場景必須真的會紅，且不能相依順序——這直接要求 §2 的 `Clock` port |

---

## 1. API 操作

### 由 log-a-workout 定義

| 方法 | 路徑 | 來源 | request | response |
| --- | --- | --- | --- | --- |
| POST | `/workouts` | `When 開始訓練` | `{}` | 201 `{workoutId, startedAt, attributedDate, status}` |
| POST | `/workouts/{workoutId}/sets` | `When 記錄動作 N 一組，次數 N` | `{exerciseId, reps?, weightKg?, durationSec?, distanceM?, warmup?}` | 201 `{setId}` |
| POST | `/workouts/{workoutId}/completion` | `When 結束訓練 N` | `{}` | 204 |
| GET | `/workouts/{workoutId}` | `When 查詢訓練 N` | — | 200 `WorkoutDetail` |

`weightKg`／`warmup` 是選用的，對應拆成獨立 `And` 的兩個步驟；
`durationSec`／`distanceM` 來自 `session-training-volume` Rule 8 的計時型與距離型。
四個選用欄位**哪些必填由動作型別決定**（Shared 4），不是四種不同的操作。

### 由 edit-a-logged-workout 定義

| 方法 | 路徑 | 來源 | request | response |
| --- | --- | --- | --- | --- |
| PUT | `/workouts/{workoutId}/sets/{setId}` | `When 修改第 N 組` ＋ `When 將第 N 組標記為熱身` | `{reps?, weightKg?, durationSec?, distanceM?, warmup}` | 204 |
| GET | `/workouts/{workoutId}/deletion-summary` | `When 查詢刪除摘要` | — | 200 `{attributedDate, exerciseCount, setCount, totalVolumeKgReps}` |
| DELETE | `/workouts/{workoutId}` | `When 刪除訓練 N` | — | 204 |

`POST /workouts/{id}/sets`（Example 1.2 的補記）不是新端點——它是 log 的同一個操作
用在已結束的訓練上。edit Rule 1 說「不論多久以前」，所以那個端點**不能檢查訓練是否
還在進行中**。這是 edit 對 log 的唯一約束。

### 由 session-training-volume 定義

**沒有新端點。** 它的 `When` 是 `重訓者 "Alice" 查詢訓練 <id>`，跟 log Example 6.2
是同一個操作。這一則貢獻的是 `WorkoutDetail` 的三個欄位（見下方聯集）。

### 由 workout-history 定義

| 方法 | 路徑 | 來源 | request | response |
| --- | --- | --- | --- | --- |
| GET | `/workouts` | `When 查詢訓練歷史` | `?page=N`（選用） | 200 `WorkoutHistoryPage` |

### 由 week-over-week-progress 定義

| 方法 | 路徑 | 來源 | request | response |
| --- | --- | --- | --- | --- |
| GET | `/progress/week-over-week` | `When 查詢本週與上週同期的比較` | — | 200 `WeekOverWeek` |

### 由 custom-exercise-library 定義

| 方法 | 路徑 | 來源 | request | response |
| --- | --- | --- | --- | --- |
| GET | `/exercises` | `When 查詢動作挑選清單` | — | 200 `[ExerciseListItem]` |
| POST | `/exercises` | `When 新增動作 "X"，型別為 "Y"` ＋ `And 該動作為單邊動作` | `{name, type?, unilateral?}` | 201 `{exerciseId}` |
| PATCH | `/exercises/{exerciseId}` | `When 改名為 "X"` ＋ `When 將型別改為 "Y"` | `{name?, type?}` | 204 |
| DELETE | `/exercises/{exerciseId}` | `When 刪除動作 N` | — | 204 |
| POST | `/exercises/{exerciseId}/archive` | `When 封存動作 N` | — | 204 |
| DELETE | `/exercises/{exerciseId}/archive` | `When 解除封存動作 N` | — | 204 |

改名（Rule 5）與改型別（Rule 6）是**同一個 `PATCH`，兩個可選欄位**，不是兩個端點：
兩者都是「改這個動作的某個屬性」，而它們的失敗條件不同（改名撞唯一性、改型別撞
引用），那是 domain 的兩條規則，不是兩個資源。

### 端點衝突與擁有者

`GET /workouts/{workoutId}` 由 log 的 Example 6.2 定義，但**六份 feature 的
`Then` 都靠它把狀態讀回來**。它的 response 是所有斷言欄位的聯集：

```
WorkoutDetail
  workoutId, startedAt, endedAt, status     ← log Rule 1、2、6
  attributedDate                            ← log Rule 6
  totalVolumeKgReps                         ← log Rule 4、5、9；edit Rule 2；volume Rule 1–4、8
  totalVolumeDisplay        string          ← volume Rule 7（"3351.0 磅·次"）
  uncountedExerciseCount    int             ← volume Rule 5、8 ＋ Shared 9
  notice                    string | null   ← volume Rule 2（"輔助配重超過體重"）
  sets[] { setId, exerciseId, weightKg?, weightDisplay, reps?, durationSec?, distanceM?, warmup }
                                            ← log Rule 3、4、9；edit Rule 1；volume Rule 7、8
```

`動作區塊應依序恰好為` 不是額外欄位——`sets[]` 保持記錄順序時，區塊是把**相鄰的
同一動作**折起來得到的（log Example 5.1 的三個區塊正是這樣來的）。多存一份會有
兩個真相。

```
WorkoutHistoryPage                          ← history Rule 1–5
  inProgress  WorkoutSummary | null
  items       []WorkoutSummary
  page        int
WorkoutSummary
  workoutId, attributedDate, exerciseCount, setCount, totalVolumeKgReps

WeekOverWeek                                ← week Rule 1–6
  thisWeek  { from, to, totalVolumeKgReps, sessionCount }
  lastWeek  { from, to, totalVolumeKgReps, sessionCount } | null   ← null 對應 Rule 5
  volumeChangePercent  int | null                                  ← null 對應 Rule 4、5
  note                 string | null                               ← Rule 4、5、6.2

ExerciseListItem                            ← custom Rule 1、4、7
  exerciseId, name, source（內建｜自訂）
```

`lastWeek` 為 `null`（week Example 5.1）與 `lastWeek.totalVolumeKgReps` 為 0
（Example 4.1）是**兩種不同的回應**：前者是「沒有上一週」，後者是「上週那幾天沒練」。

**驗證約束**

| 約束 | 來源 | 失敗回應 |
| --- | --- | --- |
| 次數型動作的 `reps` 必填 | log Rule 3 | 400 · `必要參數未提供` |
| `0 ≤ weightKg ≤ 1000`，精度小數點後一位 | log Rule 8 | 422 |
| `0 ≤ reps ≤ 100`，整數 | log Rule 8 | 422 |
| 單一動作的組數 ≤ 20 | log Rule 8 | 422 · `單一動作最多 20 組` |
| 單次訓練的相異動作數 ≤ 50 | log Rule 8 | 422 · `單次訓練最多 50 個動作` |
| 自體重型動作不需要 `weightKg` | log Rule 3（Example 3.4） | — |
| 自訂動作名稱在該重訓者庫內唯一（正規化後） | custom Rule 3 | 422 · 兩條訊息見下 |

前三條在 `api/openapi.yaml` 就表達得出來（`required`、`minimum`／`maximum`、
`multipleOf: 0.1`），中介層會先擋；其餘要看聚合或動作庫現況，只能在 domain 擋。
**兩層都要有**——契約層擋不到的若只寫在 handler，換個入口就繞過去了。

**錯誤訊息**

| 訊息 | 出現在 | HTTP |
| --- | --- | --- |
| `已有進行中的訓練` | log Rule 1（Example 1.2） | 409 |
| `必要參數未提供` | log Rule 3（Example 3.2） | 400 |
| `單一動作最多 20 組` | log Rule 8（Example 8.5） | 422 |
| `單次訓練最多 50 個動作` | log Rule 8（Example 8.6） | 422 |
| `已經有同名的動作` | custom Rule 3（Example 3.1、3.3、3.4） | 422 |
| `內建動作庫已經有這個動作` | custom Rule 3（Example 3.2）、Rule 5（Example 5.2） | 422 |
| `這個動作已經被訓練紀錄使用，只能封存` | custom Rule 4（Example 4.1） | 409 |
| `這個動作已經被訓練紀錄使用，型別不能改` | custom Rule 6（Example 6.1） | 409 |
| `輔助配重超過體重` | volume Rule 2（Example 2.5） | 200（回應內容，不是錯誤） |
| `上週同期沒有紀錄`／`這是你的第一週`／`上週有 1 個動作因為缺少體重而未計入` | week Rule 4、5、6 | 200（同上） |

log Example 8.1／8.3／8.4 只斷言 `操作失敗`，沒有斷言訊息，所以上表沒有它們的列
——**沒問過的東西不要在這裡補一個訊息出來**。見 §6。

---

## 2. Domain 型別

### 聚合與值物件

```
Workout                                        ← 聚合根，log-a-workout 定義
  id, lifterID, startedAt, endedAt, status    ← 播種表
  attributedDate  Date                        ← log Rule 6：開始當下算出後凍結
  sets            []Set                       ← 「訓練 N 有以下組」播種表

WorkoutStatus = 進行中 | 已結束                 ← 播種表出現過的值（見 §6）

Set
  id           SetID       ← 定址用；規格沒給組的 identity（見 §6）
  exerciseID   ExerciseID  ← 組播種表
  weightKg     *Decimal1   ← 組播種表；可空（log Example 3.4、volume Rule 8）
  reps         *int        ← 組播種表；計時型與距離型為空（volume Rule 8）
  durationSec  *int        ← 組播種表「秒數」（volume Example 8.1）
  distanceM    *Decimal1   ← 組播種表「距離」（volume Example 8.2）
  warmup       bool        ← 組播種表「熱身」欄

Exercise                                       ← custom-exercise-library 定義
  id           ExerciseID  ← 播種表
  name         string      ← 播種表
  type         ExerciseType ← 播種表
  unilateral   bool        ← 播種表「單邊」欄
  source       ExerciseSource ← 播種表「來源」欄（custom Rule 1、7）
  createdBy    *LifterID   ← 播種表「建立者」欄；內建為空（custom Rule 1）
  archived     bool        ← custom Rule 4（Example 4.2、4.4）

ExerciseType   = 負重型 | 自體重型 | 自體重加負重型 | 輔助型 | 計時型 | 距離型
                 ← Shared 4；六種都有例子（volume Rule 2、Rule 8）
ExerciseSource = 內建 | 自訂                    ← custom 播種表「來源」欄

Lifter
  id, name, displayUnit, timezone             ← 播種表
DisplayUnit = 公斤 | 磅                         ← volume Rule 7（兩種都有例子）

BodyweightRecord
  lifterID, recordedOn, weightKg              ← 播種表

WeekWindow                                     ← 值物件，week Rule 1、2 ＋ Shared 7、8
  SamePeriod(now, tz) → (thisFrom, thisTo, lastFrom, lastTo)
```

`attributedDate` 是**存下來的欄位而不是每次算**——log Rule 6 的「凍結」在
Example 6.2 直接可見：換時區之後仍要是 8/19。每次由 `startedAt` 現算就會跟著跑。

**Workout 的行為**（來自 `後置（狀態）` 規則）

| 方法 | 規則 | 說明 |
| --- | --- | --- |
| `Start(lifter, now)` | log Rule 1、6 | 檢查沒有進行中的；算出並凍結 `attributedDate` |
| `AddSet(exercise, measures, warmup)` | log Rule 3、5、8、9 | 附加到尾端（順序即區塊）；上限檢查；欄位由型別決定 |
| `ReplaceSet(setID, ...)` | edit Rule 1 | 就地換掉，**不動 position** |
| `Complete(now)` | log Rule 6 | 設 `endedAt`、`status`；**不動** `attributedDate` |

**Exercise 的行為**

| 方法 | 規則 | 說明 |
| --- | --- | --- |
| `Rename(newName)` | custom Rule 5 | 名稱是 Exercise 自己的狀態；唯一性檢查不在這裡（跨實體，見下） |
| `ChangeType(newType)` | custom Rule 6 | 被引用時擋下——「被引用」跨聚合，見下 |
| `Archive()` / `Unarchive()` | custom Rule 4 | 只改自己的旗標 |

log Rule 2（不自動結束）、Rule 7（同一天可多筆）、edit Rule 2（彙總即時重算）
**沒有對應的方法**——它們是「不存在的行為」：沒有逾時排程、開始時不檢查當天筆數、
不存快照。這三條靠場景防止有人之後加上去，是刻意的，不是漏列。

### Domain Service

| 服務 | 規則 | 為什麼不放進聚合 |
| --- | --- | --- |
| `VolumeCalculator` | Shared 1、4、5、6、9；volume Rule 1–8；log Rule 4、5、9；week Rule 6 | 有效重量要讀**動作型別**（Exercise）與**訓練日當下的體重**（BodyweightRecord），橫跨三個聚合 |
| `ExerciseNameNormalizer` | custom Rule 3、Rule 7 | 去頭尾空白、大小寫、全半形正規化。**同時被寫入端（Rule 3 的唯一性檢查）與讀取端（Rule 7 的遮蔽判定）使用**，見 §4 |
| `ExerciseCatalog` | custom Rule 1、4、7 | 「這個重訓者看得到哪些動作」要合併內建與自訂、排除封存、處理同名遮蔽。它橫跨多個 Exercise 實體，不屬於任何一個 |
| `ExerciseUsage` | custom Rule 4、6 | 「這個動作被幾筆訓練引用」要查 Workout 那一側。不屬於 Exercise 聚合 |

```
totalVolume(workout, exercises, bodyweightAt) =
  Σ over sets where !warmup and 型別屬於次數型:
      effectiveWeight(set) × reps × (unilateral ? 2 : 1)

effectiveWeight(set) =
  負重型          → set.weightKg                        ← volume Example 2.1
  自體重型        → 歸屬日當下最新體重，沒有則 0          ← volume Example 2.2、5.1
  自體重加負重型  → 歸屬日當下最新體重 ＋ set.weightKg    ← volume Example 2.3
  輔助型          → max(0, 歸屬日當下最新體重 − set.weightKg) ← volume Example 2.4、2.5
  計時型 / 距離型 → 不計入                               ← volume Rule 8
```

`VolumeCalculator` 同時被寫入路徑（log 的斷言）與讀取路徑（history、week）用到，
**必須是同一份實作**——history Rule 6 要的正是「查詢當下重算」，兩份實作會在
edit 之後分岔。

### 外部依賴（`port/out`）

| port | 從哪看出來 | 為什麼不能直接呼叫 |
| --- | --- | --- |
| `Clock` | **六份 feature 的 `Background` 都有 `Given 現在時間為 "..."`**；log Rule 6 的歸屬日、Rule 2 的「隔天」「放置五天」、week Rule 1／2 的「今天是週幾」、volume Rule 6 的「歸屬日當下最新體重」 | 直接用 `time.Now()` 的場景不可能穩定重跑。godog 設了 `Randomize: -1` 與 `Concurrency: 4` |
| `IDGenerator` | 播種表的「訓練 ID」「動作 ID」欄；`POST /workouts`、`POST /sets`、`POST /exercises` 都要回一個新 ID | 同上。ID 要能在場景裡被固定住才斷言得了 |

**沒有其他外部依賴**：六份 feature 沒有任何場景出現「應寄出」「應通知」、金流或
回呼。`第 N 頁`（history Rule 3）是查詢參數不是外部依賴。內建動作庫是**打包在
App 裡隨版本更新**（`built-in-library-source` 的答案），不是外部服務——所以它是
一份種子資料，不是一個 port。

### 讀取路徑的 view（不是聚合）

```
port/out.WorkoutHistoryView
  ListCompleted(lifterID, page) ([]WorkoutSummary, error)     ← history Rule 1–4
  FindInProgress(lifterID) (*WorkoutSummary, error)           ← history Rule 5

port/out.WeeklyVolumeView
  VolumeInRange(lifterID, from, to Date) (VolumeKgReps, sessionCount, error)
                                                              ← week Rule 2、3
```

`VolumeInRange` 是**一個方法被呼叫兩次**（本週、上週），不是兩個方法——兩邊套的是
完全同一套計入規則，那正是 week Rule 6 要求的。寫成兩個方法就是給規則分岔留了門。

---

## 3. Schema

判準：跨兩次使用者操作之間，什麼必須還在。使用者先開始訓練、隔一段時間才記下
一組（log Rule 2 明說可以隔天），所以**要持久化**。

### 由 log-a-workout 建立

| 資料表 | 欄位 | 來源 |
| --- | --- | --- |
| `workouts` | `id, lifter_id, started_at, ended_at, status, attributed_date` | log Rule 1、2、6 |
| `workout_sets` | `id, workout_id, position, exercise_id, weight_kg, reps, duration_sec, distance_m, warmup` | log Rule 3、4、5、9；volume Rule 8 |

`workout_sets.position` 是必要的：log Rule 5 的動作區塊靠**相鄰**判定，沒有順序欄
就分不出「臥推 3 組、划船、臥推 2 組」與「臥推 5 組、划船」。

`workouts` 需要「每個 lifter 最多一筆 `進行中`」的唯一性約束（log Rule 1）。
domain 也要擋，資料庫那層是最後一道——見 §4。

### 由 custom-exercise-library 建立

| 資料表 | 欄位 | 來源 |
| --- | --- | --- |
| `exercises` | `id, name, name_normalized, type, unilateral, source, created_by, archived` | custom Rule 1–7 |

`name_normalized` 是**存下來的欄位**而不是查詢時算——custom Rule 3 的唯一性約束要
建在正規化後的名字上（`(created_by, name_normalized)` 唯一），而資料庫的唯一索引
只能建在存下來的欄位上。這也讓 Rule 7 的遮蔽查詢用得到索引。見 §4。

### 由播種表建立

| 資料表 | 欄位 | 來源 |
| --- | --- | --- |
| `lifters` | `id, name, display_unit, timezone` | 播種表；volume Rule 7 用到 `display_unit` |
| `bodyweight_records` | `lifter_id, recorded_on, weight_kg` | 播種表；volume Rule 6 |

### 由三份 @query 讀取，不新增

`session-training-volume`、`workout-history`、`week-over-week-progress` 都只讀
上面的表。

history Rule 3 的分頁需要 `(lifter_id, status, started_at DESC)` 索引，
week Rule 1／2 的區間查詢需要 `(lifter_id, attributed_date)` 索引——效能決定不是
規格要求，等有量再說。

week 的區間查詢建在 `attributed_date` 上，**不是** `started_at`：Example 1.2
（8/16 週日那筆屬於上一週）與跨午夜那筆是同一個道理。

### 不存的東西，各有來源

| 不存 | 為什麼 |
| --- | --- |
| `workouts.total_volume` | edit Rule 2、history Rule 6 明說不保留快照，一律重算 |
| 週彙總表 | 同上；week Rule 6 要兩週套同一套規則，預先彙總會凍住當時的規則 |
| `workout_sets.effective_weight` | 自體重型的有效重量隨體重紀錄而定，存了會跟 volume Rule 6 的不回溯打架 |
| `deleted_at`（軟刪除） | edit Rule 3 明說不可復原。軟刪除欄位一旦存在，每個查詢都要記得過濾 |
| `exercises.usage_count` | custom Rule 4／6 的「被引用」要看當下的 `workout_sets`。存一份計數就要在每次記錄／刪除時維護，而漏一次就會擋掉可以刪的動作、或放行不該刪的 |
| 挑選清單的物化檢視 | custom Rule 7 的遮蔽在改名後要立刻「還原」（Example 7.3），存起來就要處理失效 |

**但這批情境的驗收測試還不需要資料庫。** 82 個場景裡只有 2 個的斷言要讀「前一次
操作」寫進去的東西（見 §3 的表），其餘換成 in-memory adapter 一樣紅、一樣綠。
建議 `internal/infrastructure/` 先只放 in-memory adapter——提早建表是最容易被
當成進度的浪費。

---

## 4. 測試分層

一個場景區塊一行。`需要資料庫` 的判準：這個場景要驗的行為本身跨越一次操作的邊界嗎？

### log-a-workout（19）

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1 沒有進行中時可開始 | 歸屬日換算（時區 → Date） | ✗ |
| 1.2 已有進行中時失敗 | 唯一性檢查 | ✗ |
| 3.2 缺少次數（Outline） | — 契約層擋掉 | ✗ |
| 3.1／3.3 提供次數時成立（Outline） | 容量公式（含 reps=0） | ✗ |
| 3.4 自體重型不帶重量 | 有效重量（自體重型分支） | ✗ |
| 8.2 小數點後一位被接受 | 精度驗證 | ✗ |
| 8.1／8.3 超範圍或精度（Outline） | 精度與值域驗證 | ✗ |
| 8.4 次數超上限 | 值域驗證 | ✗ |
| 8.5 第 21 組 | 組數上限 | ✗ |
| 8.6 第 51 個動作 | 相異動作數計算 | ✗ |
| 2.1 隔天仍可繼續記錄 | 歸屬日不隨當下時間改變 | ✗ |
| 2.2 放置五天仍進行中 | — | ✗ |
| 4.1 熱身組不計入總容量 | 容量公式（排除熱身） | ✗ |
| 4.2 全部熱身時容量為零 | 同上，邊界 | ✗ |
| 5.1 同動作多次記為獨立區塊 | 相鄰折疊 | ✗ |
| 6.1 跨午夜結束仍歸屬開始日 | 結束不動 attributedDate | ✗ |
| 6.2 換時區不改變歸屬日 | — 純讀取 | ✗ |
| 7.1 同一天兩筆相加 | 當日加總 | ✗ |
| 9.1 dropset 記為兩組 | 容量公式 | ✗ |

### edit-a-logged-workout（8）

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1 修改一週前的某一組 | 就地取代（position 不變） | ✗ |
| 1.2 補記漏掉的動作 | 附加到尾端 | ✗ |
| 1.3 把已記錄的組改標熱身 | 就地取代 | ✗ |
| 2.1 改重量後該次總容量變大 | 容量重算 | ✗ |
| 2.2 改重量後該週總容量變大 | 週彙總重算 | ✗ |
| 2.3 改標熱身後扣除但留在明細 | 容量重算（排除熱身） | ✗ |
| 3.1 查詢刪除摘要且訓練仍在 | 摘要投影 | ✗ |
| 3.2 刪除後歷史與週彙總都消失 | — | **✓** |

### session-training-volume（16）

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1 三組相加 | 容量公式 | ✗ |
| 1.2 只有熱身時為零 | 容量公式邊界 | ✗ |
| 2.1–2.4 各型別的有效重量（Outline） | **有效重量對照表**（四個分支） | ✗ |
| 2.5 輔助配重不小於體重 | 有效重量下限夾擠 | ✗ |
| 3.1／3.2 單邊 ×2（Outline） | 單邊乘數 | ✗ |
| 4.1 熱身不計入 | 容量公式（排除熱身） | ✗ |
| 5.1 無體重時為零並標示 | 未計入計數 | ✗ |
| 5.2 混合時只算得出來的 | 未計入計數 | ✗ |
| 6.1 舊訓練用當時的體重 | **體重的日期解析** | ✗ |
| 6.2 新訓練用新體重 | 同上 | ✗ |
| 7.1 磅顯示總容量 | 單位換算 | ✗ |
| 7.2 兩種單位來回不漂移（Outline） | 單位換算的可逆性 | ✗ |
| 8.1 只有棒式 | 計時型不計入 | ✗ |
| 8.2 農夫走路 | 距離型不計入 | ✗ |
| 8.3 混合時只算次數型 | 未計入計數 | ✗ |
| 8.4 全是計時型 | 未計入計數邊界 | ✗ |

`2.1–2.4` 的有效重量對照表掛了四個分支，加上 `2.5` 的夾擠、`8.x` 的兩種不計入型別
——**這是整批裡邏輯最集中的一塊**，六種型別各要一個 unit 測試。

### workout-history（8）

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1／1.2 依開始時間排序 | 排序鍵（開始時間，非歸屬日） | ✗ |
| 2.1 組數含熱身而總容量不含 | 投影計算 | ✗ |
| 3.1 第二頁只回剩下的 | 分頁切片 | ✗ |
| 3.2 不滿一頁時就是全部 | 分頁邊界 | ✗ |
| 4.1／4.2 沒有訓練時回空清單 | — | ✗ |
| 5.1 進行中不在清單而另外回傳 | 分流 | ✗ |
| 5.2 結束後回到清單的對應位置 | 同上 | ✗ |
| 6.1 修改過的顯示新的總容量 | — 純投影 | **✓** |

### week-over-week-progress（12）

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1／1.2 週從當週週一起算 | `WeekWindow`：找週一 | ✗ |
| 2.1 週三時取週一到週三 | `WeekWindow`：N＝3 | ✗ |
| 2.2 週日時整週對整週 | `WeekWindow`：N＝7 邊界 | ✗ |
| 2.3 週一時只取一天 | `WeekWindow`：N＝1 邊界 | ✗ |
| 3.1 次數相同時的成長 | 百分比計算 | ✗ |
| 3.2 多練一天時兩者都上升 | 百分比計算＋次數 | ✗ |
| 3.3 除不盡時四捨五入 | 四捨五入 | ✗ |
| 4.1 上週同期沒訓練 | 除以零的分支 | ✗ |
| 4.2 上週只有熱身組 | 除以零＋熱身排除 | ✗ |
| 5.1 沒有上一週 | 「沒有上一週」與「上週為 0」的分辨 | ✗ |
| 6.1 兩週熱身都不計入 | 容量公式（共用） | ✗ |
| 6.2 本週才補體重 | 有效重量的日期解析 | ✗ |

### custom-exercise-library（19）

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1 別人的自訂動作看不到 | `ExerciseCatalog`：可見性過濾 | ✗ |
| 1.2 自訂排在內建之後 | `ExerciseCatalog`：排序 | ✗ |
| 3.1／3.3／3.4 正規化後同名失敗（Outline） | **`ExerciseNameNormalizer`**（空白、大小寫、全半形） | ✗ |
| 3.2 與內建同名失敗 | 同上＋來源判定 | ✗ |
| 2.1 選輔助型 | 型別 → 輸入欄位對照 | ✗ |
| 2.2 不指定型別存成負重型 | 預設值 | ✗ |
| 2.3 單邊動作 | — | ✗ |
| 2.4 選計時型 | 型別 → 輸入欄位對照 | ✗ |
| 4.1 刪除被引用過的失敗 | `ExerciseUsage` | ✗ |
| 4.2 封存後不在清單但歷史看得到 | `ExerciseCatalog`：排除封存 | ✗ |
| 4.3 沒被引用過的可直接刪 | `ExerciseUsage` | ✗ |
| 4.4 解除封存後回到清單 | `ExerciseCatalog` | ✗ |
| 5.1 改名後歷史顯示新名字 | — 名稱是共用的實體屬性 | ✗ |
| 5.2 改成與內建同名失敗 | `ExerciseNameNormalizer` | ✗ |
| 6.1 改被引用過的型別失敗 | `ExerciseUsage` | ✗ |
| 6.2 沒被引用過的可改型別 | `ExerciseUsage` | ✗ |
| 7.1 同名時只出現自訂的 | **`ExerciseCatalog`：遮蔽** | ✗ |
| 7.2 沒有同名自訂的看得到內建 | 同上 | ✗ |
| 7.3 改名讓位後內建出現 | 同上（遮蔽是算出來的） | ✗ |

`7.3` 值得單獨看：它證明遮蔽**不能是存下來的旗標**——改名之後被蓋住的內建動作要
自己回來。這一格如果用「建立時標記 hidden」實作，會通過 7.1 與 7.2 而在 7.3 紅。

**只有 2 個場景需要資料庫**（edit 3.2、history 6.1），兩個都是 `Given` 是一次真正
的前一次操作。其餘 80 個在 in-memory 下一樣紅、一樣綠。

---

## 5. 技術風險

| 風險 | 從哪條規則長出來 | 影響什麼 |
| --- | --- | --- |
| **併發：訓練唯一性**。兩個 `POST /workouts` 同時通過「沒有進行中的訓練」檢查 | log Rule 1 | schema（要不要 partial unique index on `status='進行中'`）。單行程用應用層鎖就夠，多副本必須落到資料庫 |
| **併發：動作名稱唯一性**。兩個 `POST /exercises` 同時通過正規化後的同名檢查 | custom Rule 3 | schema（`(created_by, name_normalized)` 唯一索引）。跟上一條同一類，但它多一個轉折：唯一性建在**正規化後**的欄位上，所以那個欄位必須存下來（見 §3） |
| **重送／冪等：無法分辨**。`POST /sets` 重送會多一組，而**領域上無法分辨重送與真的又做了一組** | log Rule 5（同一動作可以出現多次，各自獨立）直接造成 | API（要不要冪等鍵）。這是**規則選擇的後果不是實作難題**：Rule 5 把「去重」這條路關掉了。健身房收訊差、使用者連按兩次，都會讓總容量默默偏高 |
| **正規化必須只有一份實作**。custom Rule 3 在**寫入時**用它擋同名，Rule 7 在**讀取時**用它決定遮蔽 | custom Rule 3 ＋ Rule 7 | 兩份實作會產生一種很難查的矛盾：新增被擋下（寫入端認為同名），但清單沒有遮蔽（讀取端認為不同名），於是使用者看到一個他建不出來、也看不到的名字 |
| **TOCTOU：引用檢查與寫入之間**。「這個動作沒被引用 → 刪掉」與「記錄一組用這個動作」同時發生 | custom Rule 4（Example 4.3）、Rule 6（Example 6.2） | 交易邊界。刪除動作與記錄組是兩個不同的聚合，跨聚合的一致性要嘛用交易，要嘛接受並補償。規格沒說哪一個 |
| **讀取路徑的重算成本**。每次查歷史／週比較都要從原始組重算容量；自體重型每一組要解析一次體重紀錄；每一組還要查一次動作型別 | edit Rule 2、history Rule 6、week Rule 6、volume Rule 2／6 的「不存快照」 | 分層與索引。N+1 是最直接的形式。不存快照是規格要求，解法只能在讀取那一側（一次撈進來、或可失效的快取），不能回頭去存 |
| **時區與歸屬日**。`attributedDate` 由 `startedAt` ＋ lifter 時區算出後凍結；日光節約的地區會有不存在或重複的當地時刻 | log Rule 6、Shared 2 | 時間換算的實作。台北沒有日光節約，但 `timezone` 是每個 lifter 一欄，規格上允許任何時區 |

**只指出風險，不寫死解法。** 前兩條要 Redis 還是 DB 唯一索引，需要知道實際流量與
既有基礎設施，而那些不在 `.feature` 裡。

第三條最該現在知道：它不是實作難題，是**規則選擇的後果**。要不要處理（加冪等鍵）
是產品決定——如果決定不處理，那應該是一個寫下來的決定，不是一個沒人發現的洞。

---

## 6. 實作順序

專案已經有走骨架（`GET /version` 打通 openapi → apigen → handler → godog），
不必再做一次，直接進垂直切片。

**跨 feature 的順序**：

```
custom-exercise-library（只做 Rule 2 的建立）
  → log-a-workout → session-training-volume → workout-history
  → edit-a-logged-workout → week-over-week-progress
  → custom-exercise-library（其餘）
```

`custom-exercise-library` 被切成兩段，這是這一輪跟上一輪最大的差別：
**log 的 `Background` 需要動作存在**，而動作的建立是 custom 的 Rule 2。
但 custom 的其餘規則（刪除、封存、遮蔽）反過來需要「動作被訓練引用」才驗得了
（Rule 4、6 的 `已被 N 筆訓練引用`）。硬要整則一次做完，會卡在自己身上。

- volume 排在 log 之後、history 之前：history Rule 2 的「總容量」欄要先算得對
- edit 排在 history 之後：沒有讀取路徑，「重算」跟「存了快照剛好一樣」分不出來
- week 最後：它建在歸屬日與容量上，那兩個錯了它會錯得很難看出來

**場景層級的順序**：

```
1. custom @example-2.2 不指定型別時存成負重型
   內迴圈：預設值
   解鎖：log 的全部（Background 要有動作）
   資料庫：不需要
   為什麼是它：最薄的完整路徑——一個 POST 進來，建一個實體，回一個 ID。
   它證明的是接線（openapi → apigen → handler → application → domain）
   與 IDGenerator port，不是行為。

2. log @example-1.1 沒有進行中的訓練時可以開始訓練
   內迴圈：歸屬日換算；把 Clock port 接起來
   解鎖：其餘 80 個場景的 Given
   資料庫：不需要

3. log @example-3.1 提供次數時該組成立
   內迴圈：容量公式最簡分支（負重型、非單邊、非熱身）
   解鎖：3.3、8.x、4.x、5.1、9.1，以及 volume 全部
   資料庫：不需要

4. volume @example-2.1–2.4 各型別的有效重量（Outline）
   內迴圈：有效重量對照表（六種型別的四個次數型分支）
   解鎖：volume 幾乎全部、week 6.2、custom 2.1／2.4 的欄位對照
   資料庫：不需要
   （排這麼前面是刻意的：型別表是整批的地基，六份 feature 都間接依賴它）

5. log @example-4.1 熱身組不計入 → volume 5.1（未計入計數）→ log 6.1（歸屬日凍結）

6. log 其餘 → volume 其餘（6.x 體重日期解析 → 7.x 單位 → 8.x 計時距離型）

7. history @example-2.1 每一筆帶出四個欄位
   內迴圈：投影計算
   解鎖：history 1.1、1.2、3.x、5.x（它們斷言同一張表的欄位子集）
   資料庫：不需要
   （排在 1.1 之前：排序測試需要有東西可排，2.1 先把每一筆長什麼樣定下來）

8. history 其餘 → edit 1.1 → edit 2.1（彙總重算，這是 edit 真正的目的）
   → edit 其餘 → history 6.1（要等 edit 的修改操作存在）

9. week @example-2.1 週三時兩邊都取週一到週三
   內迴圈：WeekWindow
   解鎖：week 1.x、2.2、2.3、3.1
   資料庫：不需要
   （排在 1.1 之前：1.1 只斷言兩個區間字串，2.1 斷言區間**加上**兩邊的彙總，
     是最薄的完整路徑）

10. week 其餘：3.1 → 4.1 → 5.1 → 2.2、2.3 → 1.1／1.2 → 3.2 → 3.3 → 4.2 → 6.1 → 6.2
    （4.1 與 5.1 相鄰做是刻意的：先做 4.1 會讓人想用同一個分支處理 5.1，
      而 5.1 的回應少了整個 lastWeek 物件）

11. custom 其餘：3.x（正規化）→ 1.x（可見性）→ 7.1、7.2 → 7.3（遮蔽要能還原）
    → 4.x、6.x（要 ExerciseUsage，所以要先有訓練引用動作）
    → 2.1、2.3、2.4、5.x
```

---

## 7. 規格裡推不出來的

### 回 CLARIFY 補問（規格沒講）

| 缺口 | 影響 |
| --- | --- |
| **失敗場景整批缺**（bdd-spec 的「三種結果」檢查抓到的） | 43 條規則裡 37 條只走了單一結果。具體：修改／刪除／查詢一筆**不存在**的訓練回什麼？在**進行中**的訓練上做 edit 的操作可以嗎？刪除進行中的訓練之後 log Rule 1 的唯一性怎麼算？三題都沒問過 |
| **以磅輸入的換算沒有任何場景覆蓋** | volume Rule 7 只驗到顯示端（已存的 20.4 kg 在兩種單位下怎麼顯示）。「使用者以磅輸入 45」沒有任何 `@ready` 的命令場景——log 的重訓者播種表一律是公斤。輸入端的換算與捨入完全沒被驗到 |
| **封存的動作還能不能被記錄引用？** | custom Rule 4.2 只說它不出現在挑選清單，沒說 `POST /sets` 帶那個 ID 會怎樣。清單是 UI，API 擋不擋是另一回事 |
| `WorkoutStatus` 只看得到 `進行中`、`已結束` | 沒有證據說還有沒有別的（放棄？取消？）。狀態轉移表可能缺一格 |
| history 的頁碼超出範圍（第 99 頁）回什麼？ | Rule 3 只有「有下一頁」與「不滿一頁」兩個例子，兩個都在範圍內 |
| history 的排序在**開始時間完全相同**時怎麼辦？ | 沒有次鍵。兩筆同秒開始順序就不決定，而 `應依序恰好為` 會間歇性紅（godog `Concurrency: 4`） |
| **體重、時區、顯示單位怎麼被設定？** | 三個都在播種表裡，也都有規則用到它們（volume Rule 6、7；log Rule 6），但**沒有任何 `When` 建立或修改它們**。這跟 `actor.md` 裡「怎麼取得這個身分」空著是同一個洞 |
| log Example 8.1／8.3／8.4 只斷言「操作失敗」，沒有訊息；8.5／8.6 有 | 不一致。是刻意的還是漏了？不定下來，IMPLEMENT 會自己編三條訊息 |
| 負重型動作的重量填 `0` 該成立還是該擋？ | log Rule 8 說「0 到 1000」，但沒有例子用 0。「不帶重量」（Example 3.4）跟「重量 0」不是同一件事 |
| 開始訓練時 `startedAt` 由誰決定？ | 所有場景都用 `現在時間為` 播種，看不出是伺服器時鐘還是客戶端送的。離線補記會直接撞到這裡 |
| 計時型的秒數、距離型的距離有沒有範圍與精度？ | log Rule 8 只約束了重量與次數。新的兩個欄位沒有任何邊界例子——這是 `timed-and-distance-exercises` 答完之後長出來的新缺口 |
| 內建動作庫初版收錄哪些、誰標六種型別、多久更新一次 | `built-in-library-source` 明記這三件沒有來源。不影響行為規則所以不擋就緒，但**是實際要有人做的工作**，排期時不能當成沒有 |

**兩張紅卡的答案是代答的**（`timed-and-distance-exercises`、`built-in-library-source`），
兩份問題檔都標了低信心。它們決定了核心指標的邊界與動作庫的所有權模型，
**應該由產品確認過再開工**——這一節的前兩列與最後一列都是它們的直接後果。

### 留給 IMPLEMENT（技術決定）

| 決定 | 為什麼現在不決定 |
| --- | --- |
| 「組」的定址用位置還是 SetID | `.feature` 寫「第 2 組」是位置，但位置會因補記而位移。ID 是定址不是行為。建議 SetID，由 `POST /sets` 回傳 |
| 改名與改型別用一個 `PATCH` 還是兩個端點 | 規格只說「改名」「改型別」，沒說資源怎麼切 |
| 分頁用 offset 還是 cursor | history Rule 3 只說每頁 30 筆，兩種都滿足 |
| `page` 從 0 還是 1 起算 | 純慣例 |
| `POST /workouts` 回 201＋Location 還是 200＋body | 規格只要求拿得到新訓練 |
| 結束訓練用 `POST .../completion` 還是 `PATCH status` | 規格只說「結束」 |
| Problem 的 `type` URI 命名 | RFC 9457 的形狀已定，命名是慣例 |
| 全半形正規化用 NFKC 還是自訂對照表 | custom Rule 3 只給了「臥推（三樓）」與「臥推(三樓)」一組例子，兩種都滿足 |
| 「當日總容量」「那一週的總容量」用哪條讀取路徑觀察 | 這兩個斷言沒有對應的 `When`。建議用 `GET /workouts` 逐筆加總，**不新增端點** |
| `warmup` 省略時預設 `false` | 所有播種表都明寫「否」 |
| 組數／動作數上限是常數還是設定 | `session-size-limit.md` 信心低，值會調；先做成具名常數 |
