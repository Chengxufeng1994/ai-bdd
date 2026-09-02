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
| 被至少一筆 `workout_sets` 引用的 Exercise 不會被刪除（只能封存）——換句話說，`workout_sets.exercise_id` 指向的 Exercise 列永遠存在 | custom-exercise-library Rule 4 | Exercise 刪除操作在動手前透過 `ExerciseUsage` 服務查引用數，引用中即拒絕；資料庫層以外鍵（`workout_sets.exercise_id → exercises.id`）作最後一道防線 |

## 跨聚合的規則

沒有找到「跨兩個以上聚合、且不執行操作就能由現有狀態直接驗證」的規則。上面
Exercise 的第二條不變條件雖然由 custom-exercise-library Rule 4 而來、且要看
Workout 那側的 `workout_sets`，但它本質是 Exercise 這個聚合「被引用中不可刪除」
的邊界（誰在引用是外部事實，不變條件本身仍然只約束 Exercise 自己的生命週期），
所以歸在 Exercise 底下，不是獨立的跨聚合規則。

custom-exercise-library Rule 6（被引用過的動作不能改型別）**不算**不變條件，
留在 `spec.md` 當 `.feature` 的 `Rule:`：它與 Rule 4 表面相似（都提到「被引用」），
但可驗證的方式不同。Rule 4 是一句可以直接對照快照回答的話——「現在有沒有一筆
`workout_sets.exercise_id` 指不到任何 `exercises` 列？」單一時間點的一次查詢
就能回答有或沒有，跟已 promote 的兩個不變條件（唯一性）同構：都是「檢查一個當下
就能為真或為假的陳述」。Rule 6 沒有這種陳述可寫：「這個 Exercise 的 type 從被
引用起就沒變過」不是狀態的函式，是**歷史**的函式——光看現在的 `type` 欄位，
分不出它是「從來沒被試著改過」還是「試著改過但被擋下了」，兩者當下的資料庫內容
完全相同。驗得出真假只能靠真的送一次 `PATCH` 改型別、看它成不成功，這正是
「只能靠嘗試操作驗證」的定義，所以 Rule 6 留在 `spec.md`。

## 修訂紀錄

| 日期 | 批次 | 改了什麼 | 為什麼 |
| --- | --- | --- | --- |
| 2026-09-02 | workout-tracking | 初版：從 `plan.md` §2「Domain 型別」抽出 Workout、Exercise 兩個聚合與各自的唯一性不變條件 | 本專案第一份 `domain-model.md`；`plan.md` §2 混了型別欄位（留 `spec.md`）與不變條件（跨批次仍然成立），這裡只留下後者 |
| 2026-09-02 | workout-tracking | 修正：custom-exercise-library Rule 4 從「不算不變條件」改列為 Exercise 的不變條件；Rule 6 維持不算 | 上一版對 Rule 4／Rule 6 套用同一句理由（「只能靠嘗試操作驗證」），但 review 指出 Rule 4 其實是可由單一時間點快照直接驗真假的參照完整性事實，跟已 promote 的兩個唯一性不變條件同構，不該套用 Rule 6 的理由；重新用「能不能只看當下狀態回答」個別檢驗兩條，得出不同結論 |
