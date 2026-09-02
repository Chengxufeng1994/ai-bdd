# Domain Model

**更新於**: 2026-09-02 · **最後一批**: workout-tracking

## 聚合

### Workout

**是什麼**: 使用者的一次訓練紀錄，含開始／結束時間、狀態，以及訓練中記錄的所有組（Set）
**邊界**: Set 附屬於 Workout，隨 Workout 一起變更（新增組、換組都是對同一個 Workout 聚合的操作）；Exercise、Lifter、BodyweightRecord 只被引用（存 ID），不隨 Workout 變更
**來源**: log-a-workout Rule 1

**不變條件**（任何時刻都必須為真）

| 條件 | 來源 | 被誰保證 |
| --- | --- | --- |
| 一個使用者（lifter）同一時間最多一筆進行中（狀態為「進行中」）的訓練 | log-a-workout Rule 1 | Workout 聚合的建構規則 ＋ `workouts` 資料表的唯一性約束（partial unique index on `status = '進行中'`） |

### Exercise

**是什麼**: 可被訓練組引用的動作定義，分內建與自訂兩種來源
**邊界**: 不含引用它的 Set／Workout（那些屬於別的聚合，只用 `exerciseID` 引用）
**來源**: custom-exercise-library Rule 3

**不變條件**（任何時刻都必須為真）

| 條件 | 來源 | 被誰保證 |
| --- | --- | --- |
| 同一使用者（`created_by`）建立的自訂動作名稱，正規化（去頭尾空白、大小寫、全半形）後彼此不重複 | custom-exercise-library Rule 3 | Exercise 聚合的建構規則 ＋ `exercises` 資料表的唯一性索引 `(created_by, name_normalized)` |

## 跨聚合的規則

目前沒有找到「跨兩個以上聚合、且不執行操作就能由現有狀態直接驗證」的規則。
custom-exercise-library Rule 4（被引用過的動作不能刪除，只能封存）與 Rule 6
（被引用過的動作不能改型別）表面上橫跨 Exercise 與 Workout，但兩者都只在**嘗試
刪除／嘗試改型別**這個操作發生的當下才驗得出來，不是靠讀取目前狀態就能判斷
真假的性質，屬於 `.feature` 的 `Rule:`，留在 `spec.md`，不算不變條件。

## 修訂紀錄

| 日期 | 批次 | 改了什麼 | 為什麼 |
| --- | --- | --- | --- |
| 2026-09-02 | workout-tracking | 初版：從 `plan.md` §2「Domain 型別」抽出 Workout、Exercise 兩個聚合與各自的唯一性不變條件 | 本專案第一份 `domain-model.md`；`plan.md` §2 混了型別欄位（留 `spec.md`）與不變條件（跨批次仍然成立），這裡只留下後者 |
