# 訓練追蹤

**Slice**: workout-tracking
**日期**: 2026-09-02
**涵蓋**: log-a-workout、session-training-volume、edit-a-logged-workout、workout-history、week-over-week-progress、custom-exercise-library（6 則）
**跳過**: 無 —— 六則 story 這一輪 CLARIFY 全數收斂為已就緒（提案），本批全收

## Problem Statement

重訓者上健身房做重量訓練，想把每次練了什麼記下來，事後才看得出自己有沒有在進步。
現在靠記憶或紙本：動作、組數、重量、次數當場記不住就漏記，紙本寫錯了要劃掉重寫，
版面很快就亂掉。訓練結束想知道這次算不算有練到，得自己心算每一組的重量乘次數再
加總，容易算錯，也沒辦法馬上知道結果。想知道「這週比上週好還是壞」，得翻出上一頁
紙自己對照兩邊的數字。記過的東西之後想改——把打錯的重量改回來、補上忘記記的一組、
刪掉記錯的整次訓練——紙本上也做得到，但留下的痕跡很難看得清楚哪個數字才是最新的。
另外，每家健身房對同一個動作、同一台器材的叫法都不一樣，重訓者得自己在腦中維護一份
翻譯對照表，才能把紙上寫的和自己健身房的說法對起來。

## Solution

重訓者可以在健身房現場當下，把每個動作的每一組（重量、次數，或依動作型別而定的
秒數、距離）記下來；訓練結束的當下就直接看到這次的總容量，不必自己按計算機。
記錄之後，任何欄位都可以修改，也可以把整次訓練刪除；不論是總容量、當日彙總還是
週彙總，全部以目前的資料重新算出來，不會因為改了一個數字而有兩個地方對不上。
重訓者可以隨時翻閱過去每一次訓練的清單，看到每次的動作數、組數與總容量；也可以
看到本週跟上週同一段時間的訓練量與訓練次數並排比較，用來判斷自己有沒有在進步。
最後，重訓者可以用自己健身房慣用的叫法新增動作、之後改名或不再使用時封存它，
不需要遷就系統內建的名字，也不用擔心刪掉一個已經記錄過的動作會讓歷史紀錄跟著消失。

## User Stories

1. 作為重訓者，我要在健身房當場記下每個動作的每一組，以便不必靠記憶或紙本，也不會漏記。
2. 作為重訓者，我要在訓練結束時直接看到這次的總容量，以便不用自己按計算機，也不會算錯。
3. 作為重訓者，我要能修改或刪除已經記完的訓練，以便補上漏記的組、改掉打錯的數字。
4. 作為重訓者，我要看到過去每一次訓練的清單與各自的總容量，以便回想自己練了什麼、練得怎麼樣。
5. 作為重訓者，我要看到本週跟上週同期的訓練量對照，以便判斷自己有沒有在進步。
6. 作為重訓者，我要用自己健身房的叫法新增動作，以便記錄時找得到、不用去猜系統把它叫成什麼。

## Example Mapping

### log-a-workout

**Story**: 作為<重訓者>，我要在健身房當場記下每個動作的每一組，以便不必靠記憶
或紙本，也不會漏記

#### Rule 1. 同一時間最多只有一筆進行中的訓練
- Example 1.1 沒有進行中的訓練時按「開始訓練」→ 建立一筆進行中的訓練，開始時間 2026-08-19 19:12
- Example 1.2 已有一筆 2026-08-18 19:12 開始、尚未結束的訓練，再按「開始訓練」→ 擋下，提示「你有一筆 8/18 的訓練還沒結束」並提供跳過去的入口

#### Rule 2. 進行中的訓練不會自動結束
- Example 2.1 19:12 開始，記了 3 組後離開 App，隔天 09:00 回來 → 那筆訓練仍是進行中，歸屬日仍是 8/19
- Example 2.2 進行中的訓練放了 5 天 → 仍是進行中，不因時間長短而改變

#### Rule 3. 一組必須有次數才成立；重量欄由動作型別決定要不要填
- Example 3.1 槓鈴臥推（負重型）填 60 kg、8 下 → 接受
- Example 3.2 槓鈴臥推填了 60 kg，次數留空 → 擋下，提示「還沒填次數」
- Example 3.3 槓鈴臥推 60 kg、0 下 → 接受，記為一組，容量 0（推不起來的那一組）
- Example 3.4 伏地挺身（自體重型）只填 20 下 → 接受，畫面上沒有重量欄

#### Rule 4. 標記為熱身的組完整記錄，但不計入總容量（Shared 5）
- Example 4.1 深蹲 4 組：20 kg×10（熱身）、40 kg×10（熱身）、80 kg×8、80 kg×8 → 4 組都看得到，總容量 ＝ 80×8＋80×8 ＝ 1280
- Example 4.2 一個動作的 3 組全部標成熱身 → 該動作出現在紀錄裡，對總容量貢獻 0

#### Rule 5. 同一個動作可以在一次訓練中出現多次，各自獨立
- Example 5.1 臥推 3 組 → 划船 3 組 → 臥推再 2 組 → 畫面上是 3 個區塊、順序不變；總容量把 5 組臥推全部算進去

#### Rule 6. 訓練歸屬於開始時間所在的當地日期（Shared 2）
- Example 6.1 2026-08-19 23:40 開始、2026-08-20 00:50 結束 → 歸屬 8/19，不切成兩筆
- Example 6.2 在台北（UTC+8）記完的 8/19 訓練，使用者之後飛到東京（UTC+9）→ 歷史上仍顯示 8/19

#### Rule 7. 同一天可以有多筆訓練，當日彙總是相加
- Example 7.1 8/19 07:00 練腿（總容量 1600）、19:12 練胸（總容量 480）→ 兩筆獨立的訓練，8/19 當日總容量 2080

#### Rule 8. 輸入的合法範圍
- Example 8.1 重量 1000.5 kg → 擋下，提示上限 1000
- Example 8.2 重量 62.5 kg → 接受（小數點後一位）
- Example 8.3 重量 62.55 kg → 擋下，只接受到小數點後一位
- Example 8.4 次數 101 → 擋下，提示上限 100
- Example 8.5 同一動作要新增第 21 組 → 擋下，提示「單一動作最多 20 組」
- Example 8.6 一次訓練要新增第 51 個動作 → 擋下，提示「單次訓練最多 50 個動作」

#### Rule 9. Dropset 記成連續的多組，沒有專屬概念
- Example 9.1 臥推 60 kg×8 力竭後立刻降到 45 kg×5 → 記成兩組（60×8、45×5），總容量 480＋225 ＝ 705

### session-training-volume

**Story**: 作為<重訓者>，我要在訓練結束時直接看到這次的總容量，以便不用自己
按計算機，也不會算錯

#### Rule 1. 總容量 ＝ 所有計入的組的（有效重量 × 次數）之和（Shared 1）
- Example 1.1 臥推 60 kg×10、60 kg×8、55 kg×8 → 600＋480＋440 ＝ 1520 公斤·次
- Example 1.2 一次訓練只有熱身組 → 總容量 0，畫面顯示 0 而非空白

#### Rule 2. 有效重量由動作型別決定（Shared 4）
- Example 2.1 負重型：槓鈴臥推 60 kg×10，使用者填 60 → 有效重量 60，容量 600
- Example 2.2 自體重型：伏地挺身×20，使用者體重 70 kg → 有效重量 70，容量 1400
- Example 2.3 自體重加負重型：負重引體向上 加掛 10 kg×5，體重 70 → 有效重量 80，容量 400
- Example 2.4 輔助型：輔助引體向上 輔助配重 25 kg×8，體重 70 → 有效重量 45，容量 360
- Example 2.5 輔助型且輔助配重 75 kg ≥ 體重 70 → 有效重量以 0 計，該組容量 0，畫面提示「輔助配重超過體重」

#### Rule 3. 單邊動作的容量乘以 2（Shared 6）
- Example 3.1 單手啞鈴划船 20 kg，次數欄填 8（左右各 8 下）→ 容量 20×8×2 ＝ 320；畫面顯示「8 下 ×2 邊」
- Example 3.2 單邊的自體重型（單腳臀橋）×12，體重 70 → 70×12×2 ＝ 1680

#### Rule 4. 標記為熱身的組不計入（Shared 5，規則本體在 log-a-workout Rule 4）
- Example 4.1 深蹲 20×10（熱身）、40×10（熱身）、80×8、80×8 → 總容量 1280

#### Rule 5. 沒有體重資料時，需要體重的型別一律以 0 計並在畫面標示
- Example 5.1 沒填過體重的使用者做伏地挺身×20 → 該動作容量 0，畫面顯示「補上體重才能算這個動作的容量」與設定入口
- Example 5.2 同一次訓練裡臥推 60×10、伏地挺身×20（無體重）→ 總容量 600，並標示「1 個動作未計入」

#### Rule 6. 用訓練日當下最新的一筆體重，之後更新體重不回溯
- Example 6.1 8/01 記體重 70 kg，8/10 做伏地挺身×20（容量 1400），8/15 把體重改成 68 → 8/10 的容量仍是 1400
- Example 6.2 8/16 之後的自體重動作改用 68 kg 計算

#### Rule 7. 顯示單位只影響顯示，不影響儲存（Shared 3）
- Example 7.1 總容量 1520 公斤·次，顯示單位設為磅 → 顯示 3351.0 磅·次
- Example 7.2 以磅輸入 45 lb → 存為 20.4 kg；切回磅顯示 45.0 lb；來回切換十次數值不漂移

#### Rule 8. 計時型與距離型動作完整記錄，但對總容量的貢獻為 0，並套用 Shared 9 的標示
- Example 8.1 棒式 60 秒、60 秒、45 秒共 3 組 → 三組的秒數都看得到，總容量貢獻 0，畫面標示「1 個動作未計入」
- Example 8.2 農夫走路 40 kg 走 30 公尺 ×2 組 → 重量與距離都記下來，總容量貢獻 0
- Example 8.3 一次訓練只有臥推 60 kg×10（容量 600）與棒式 60 秒×3 → 總容量 600，畫面標示「1 個動作未計入」
- Example 8.4 一次訓練只有棒式 → 總容量 0，畫面標示「1 個動作未計入」而不是空白

### edit-a-logged-workout

**Story**: 作為<重訓者>，我要能修改或刪除已經記完的訓練，以便補上漏記的組、
改掉打錯的數字

#### Rule 1. 已結束的訓練可以修改任何欄位，沒有時間限制
- Example 1.1 8/12（一週前）的臥推某組從 60 kg×8 改成 70 kg×8 → 接受
- Example 1.2 在 8/12 的訓練裡補一個 8/12 當天忘記記的動作 → 接受
- Example 1.3 把 8/12 某組改標為熱身 → 接受

#### Rule 2. 所有彙總以當前資料重算，不保留快照
- Example 2.1 8/12 原本總容量 1520，把某組從 60×8 改成 70×8 → 8/12 總容量變成 1600
- Example 2.2 承上，8/12 所在那一週的總容量同步變成 1600
- Example 2.3 把 60×10 那組改標為熱身 → 總容量從 1520 降為 920，明細裡仍看得到那一組

#### Rule 3. 刪除訓練不可復原，且刪除前要看得到這筆的摘要
- Example 3.1 對 8/12 的訓練要求刪除 → 先看到摘要「2026-08-12，1 個動作、3 組、總容量 1520」，此時訓練還在
- Example 3.2 刪除後 → 該筆從歷史消失，8/12 當日與該週的總容量都變成 0，沒有垃圾桶可復原

### workout-history

**Story**: 作為<重訓者>，我要看到過去每一次訓練的清單與各自的總容量，以便回想
自己練了什麼、練得怎麼樣

#### Rule 1. 以「一次訓練」為單位，依開始時間由新到舊排序
- Example 1.1 8/19 07:00 練腿、8/19 19:00 練胸 → 19:00 那筆排在 07:00 那筆前面
- Example 1.2 8/19 23:40 開始、8/20 00:50 結束的訓練 → 歸屬日是 8/19（Shared 2），但排序看開始時間，所以排在 8/19 19:00 那筆之前

#### Rule 2. 每一列顯示歸屬日、動作數、總組數、總容量；組數含熱身而總容量不含
- Example 2.1 8/19 19:00 那筆：臥推 20×10（熱身）、60×10、60×8，划船 50×10、50×10 → 顯示「2 個動作 · 5 組 · 2080 公斤·次」

#### Rule 3. 每頁最多 30 筆
- Example 3.1 使用者有 45 筆訓練 → 第 1 頁 30 筆，第 2 頁 15 筆
- Example 3.2 使用者有 3 筆 → 第 1 頁就是全部，沒有第 2 頁

#### Rule 4. 一筆紀錄都沒有時查詢成功並回空清單
- Example 4.1 全新使用者開啟歷史 → 查詢成功、清單為空，畫面顯示「還沒有訓練紀錄」與開始訓練的入口
- Example 4.2 使用者刪掉了唯一一筆紀錄 → 回到同一個空清單狀態

#### Rule 5. 進行中的訓練不在清單裡，另外回傳
- Example 5.1 有一筆 8/18 開始、尚未結束的訓練 → 清單從 8/19 那三筆開始，另外回傳那筆進行中的
- Example 5.2 該筆結束後 → 不再另外回傳，它出現在清單裡 8/18 的位置

#### Rule 6. 清單顯示的總容量以查詢當下的資料重算，不使用記錄當時的快照
- Example 6.1 把 8/19 07:00 那筆的第 1 組從 100×8 改成 110×8 之後再開歷史 → 那一列的總容量從 1600 變成 1680

### week-over-week-progress

**Story**: 作為<重訓者>，我要看到本週跟上週同期的訓練量對照，以便判斷自己
有沒有在進步

#### Rule 1. 一週從當地時間週一 00:00 起算（Shared 7）
- Example 1.1 8/19 是週三 → 本週指 8/17（週一）到 8/23（週日）
- Example 1.2 8/16（週日）的訓練 → 屬於 8/10–8/16 那一週，不是本週

#### Rule 2. 比較兩週的同期：第 1 到第 N 天，N ＝ 本週已過的天數含今天（Shared 8）
- Example 2.1 今天是 8/19（週三，N＝3）→ 比較 8/17–8/19 與 8/10–8/12，畫面標明比較的是週一到週三
- Example 2.2 今天是 8/23（週日，N＝7）→ 比較 8/17–8/23 與 8/10–8/16，等於整週對整週
- Example 2.3 今天是 8/17（週一，N＝1）→ 比較 8/17 與 8/10 兩天

#### Rule 3. 回傳週總容量與訓練次數；容量變化 ＝（本週 − 上週）÷ 上週，四捨五入到整數百分比
- Example 3.1 本週同期 3600／2 次，上週同期 3000／2 次 → 顯示「容量 +20%」與「訓練 2 次，與上週相同」
- Example 3.2 本週同期 4200／3 次，上週同期 3000／2 次 → 顯示「容量 +40%」與「訓練 3 次，比上週多 1 次」——讓使用者自己看出成長有一部分來自多練一天
- Example 3.3 本週同期 3800、上週同期 3000 → (3800−3000)÷3000 ＝ 26.67% → 顯示 +27%

#### Rule 4. 上週同期總容量為 0 時不顯示百分比
- Example 4.1 上週同期沒有任何訓練、本週同期 2200 → 顯示「上週同期沒有紀錄」，不顯示 +100% 也不顯示 +∞
- Example 4.2 上週同期只練了熱身組（容量 0，但有 1 次訓練）→ 容量不顯示百分比，訓練次數仍照常比較

#### Rule 5. 完全沒有上一週的資料時，顯示「這是你的第一週」
- Example 5.1 使用者 8/19 第一次使用，8/17 之前沒有任何紀錄 → 顯示「這是你的第一週，下週就能比較了」，上週那一欄整個沒有值

#### Rule 6. 兩週用同一套計入規則
- Example 6.1 上週熱身 20×10、本週熱身 40×10，兩週正式組都是 1000 → 兩週總容量都是 1000，變化 0%（熱身若計入，會變成 +17%）
- Example 6.2 上週與本週各做伏地挺身 20 下，但體重是 8/18 才補記的 → 上週那筆以 0 計、本週用 70 kg 計，比較顯示 +140%，並標示「上週有 1 個動作因為缺少體重而未計入」

### custom-exercise-library

**Story**: 作為<重訓者>，我要用自己健身房的叫法新增動作，以便記錄時找得到、
不用去猜系統把它叫成什麼

#### Rule 1. 自訂動作只有建立者看得到
- Example 1.1 A 建立「史密斯臥推（三樓）」→ B 的動作庫裡沒有這一項
- Example 1.2 A 建立的動作出現在 A 記錄時的挑選清單，排在內建動作之後

#### Rule 2. 建立時必須選型別，預設負重型（Shared 4）
- Example 2.1 建立「輔助引體向上（藍色帶）」選輔助型 → 之後記錄這個動作時，畫面上的欄位是「輔助配重」而不是「重量」
- Example 2.2 建立時不動型別 → 存成負重型，記錄時出現「重量」欄
- Example 2.3 建立時勾選「單邊動作」→ 記錄時次數欄標示「單邊次數」，容量 ×2（Shared 6）
- Example 2.4 建立「棒式（抬腿）」選計時型 → 記錄時的欄位是「秒數」而不是「次數」，且該動作不計入總容量（Shared 4）

#### Rule 3. 名稱在該使用者的動作庫內唯一，比對前先正規化（去頭尾空白、大小寫、全半形）
- Example 3.1 已有「啞鈴臥推」，再建「啞鈴臥推 」（尾端一個空白）→ 擋下，提示已經有同名的動作
- Example 3.2 已有內建的「槓鈴臥推」，使用者建同名 → 擋下，提示改用內建的那一個
- Example 3.3 已有「Bench Press」，再建「bench press」→ 擋下
- Example 3.4 已有「臥推（三樓）」（全形括號），再建「臥推(三樓)」（半形）→ 擋下

#### Rule 4. 被訓練紀錄引用過的動作不能刪除，只能封存
- Example 4.1 「史密斯臥推（三樓）」已被 8/12 的訓練引用，按刪除 → 提示不能刪，改提供「封存」
- Example 4.2 封存後 → 記錄時的挑選清單看不到它；8/12 的歷史紀錄仍完整顯示它與它的容量
- Example 4.3 從未被任何訓練引用過的動作 → 可以直接刪除
- Example 4.4 封存過的動作可以解除封存，重新回到挑選清單

#### Rule 5. 改名會套用到所有歷史紀錄
- Example 5.1 「史密斯臥推（三樓）」改名為「史密斯臥推」→ 8/12 的歷史紀錄顯示新名字
- Example 5.2 改名後的新名字仍要通過 Rule 3 的唯一性檢查

#### Rule 6. 被引用過的動作不能改型別
- Example 6.1 「引體向上」原本存成負重型、已被 3 次訓練引用，想改成自體重型 → 擋下，提示建立一個新動作並把舊的封存
- Example 6.2 從未被引用過的動作 → 型別可以改

#### Rule 7. 內建動作庫隨 App 版本更新；新增的內建動作若與使用者既有的自訂動作同名，使用者的自訂動作優先
- Example 7.1 使用者已自訂「單槓引體向上」，App 更新後內建庫也收錄同名動作 → 他的清單只出現自己那一個，8/12 的歷史紀錄不受影響
- Example 7.2 沒有自訂同名動作的使用者 → 更新後清單出現新收錄的內建動作
- Example 7.3 使用者把自訂的「單槓引體向上」改名為「引體向上（單槓）」之後 → 原本被讓位的內建動作出現在他的清單裡

## Implementation Decisions

- **動哪些模組**：全部新增，專案目前只有 `GET /version` 這一條走骨架路徑。
  domain 層新增 `Workout`（聚合根，含 `Set`）與 `Exercise` 兩個聚合，以及
  `WeekWindow` 等值物件；四個橫跨聚合或無單一擁有者的規則放進 domain
  service：`VolumeCalculator`、`ExerciseNameNormalizer`、`ExerciseCatalog`、
  `ExerciseUsage`。application 層新增對應每個 API 操作的 command／query
  use case、實作驅動埠的 service，以及兩個新的驅動外埠（Clock、
  IDGenerator）加上訓練與動作的讀寫埠。interfaces/http 依 API 契約新增每個
  操作的 handler，沿用既有的 mapper／presenter／errmap 三段式。
  infrastructure 依 §3 的建議先只新增 in-memory adapter（訓練、動作、體重
  的讀寫，Clock、IDGenerator），因為 82 個場景中有 80 個不需要真的資料庫。
  以上模組切分是把 API 契約、domain 型別與既有分層（沿用各層既有的套件文件慣例）
  三者對齊推出來的，不是新決定。

- **介面**：domain 層的行為已經有明確的方法形狀（來自既有規則推導，非新決定）：

  | 型別 | 方法 | 對應規則 |
  | --- | --- | --- |
  | `Workout` | `Start(lifter, now)` | log Rule 1、6 |
  | `Workout` | `AddSet(exercise, measures, warmup)` | log Rule 3、5、8、9 |
  | `Workout` | `ReplaceSet(setID, ...)` | edit Rule 1 |
  | `Workout` | `Complete(now)` | log Rule 6 |
  | `Exercise` | `Rename(newName)` | custom Rule 5 |
  | `Exercise` | `ChangeType(newType)` | custom Rule 6 |
  | `Exercise` | `Archive()` / `Unarchive()` | custom Rule 4 |

  `VolumeCalculator` 的形狀（比散文精確，逐字保留自 plan.md，因為它就是
  容量規則本身的編碼，不是實作細節）：

  ```
  totalVolume(workout, exercises, bodyweightAt) =
    Σ over sets where !warmup and 型別屬於次數型:
        effectiveWeight(set) × reps × (unilateral ? 2 : 1)

  effectiveWeight(set) =
    負重型          → set.weightKg
    自體重型        → 歸屬日當下最新體重，沒有則 0
    自體重加負重型  → 歸屬日當下最新體重 ＋ set.weightKg
    輔助型          → max(0, 歸屬日當下最新體重 − set.weightKg)
    計時型 / 距離型 → 不計入
  ```

  讀取路徑的兩個埠（既有結論，非新決定）：

  ```
  port/out.WorkoutHistoryView
    ListCompleted(lifterID, page) ([]WorkoutSummary, error)
    FindInProgress(lifterID) (*WorkoutSummary, error)

  port/out.WeeklyVolumeView
    VolumeInRange(lifterID, from, to Date) (VolumeKgReps, sessionCount, error)
  ```

  應用層對外公開的驅動埠（六個 story 合起來要收在幾個服務、每個服務叫什麼
  名字）沒有答案——見 Out of Scope 的「回 CLARIFY 補問」。`Clock`、
  `IDGenerator` 兩個新驅動外埠已知一定要存在（六份 feature 的 `Background`
  都靠 `Given 現在時間為 "..."` 播種、四個建立操作都要回新 ID），但確切的
  方法簽章沒有來源，同樣列在 Out of Scope。

- **API 契約**：

  | 方法 | 路徑 | 來源 story |
  | --- | --- | --- |
  | POST | `/workouts` | log-a-workout |
  | POST | `/workouts/{workoutId}/sets` | log-a-workout |
  | POST | `/workouts/{workoutId}/completion` | log-a-workout |
  | GET | `/workouts/{workoutId}` | log-a-workout（六份 story 的 `Then` 共用這個查詢） |
  | PUT | `/workouts/{workoutId}/sets/{setId}` | edit-a-logged-workout |
  | GET | `/workouts/{workoutId}/deletion-summary` | edit-a-logged-workout |
  | DELETE | `/workouts/{workoutId}` | edit-a-logged-workout |
  | GET | `/workouts` | workout-history |
  | GET | `/progress/week-over-week` | week-over-week-progress |
  | GET | `/exercises` | custom-exercise-library |
  | POST | `/exercises` | custom-exercise-library |
  | PATCH | `/exercises/{exerciseId}` | custom-exercise-library（改名與改型別是同一個操作的兩個可選欄位，因為兩者的失敗條件不同屬於 domain 的兩條規則，不是兩個資源） |
  | DELETE | `/exercises/{exerciseId}` | custom-exercise-library |
  | POST | `/exercises/{exerciseId}/archive` | custom-exercise-library |
  | DELETE | `/exercises/{exerciseId}/archive` | custom-exercise-library |

  `plan.md` 對「改名與改型別是一個 `PATCH` 還是兩個端點」自相矛盾：§1 把上面這行
  當已決定，§7「留給 IMPLEMENT」卻把同一題列成未決（「規格只說『改名』『改型別』，
  沒說資源怎麼切」）。這裡保留 §1 的框架——失敗條件不同是 domain 的兩條規則，不是
  兩個資源，兩個可選欄位仍分得清——但把分歧本身寫下來：§7 的顧慮沒有被反駁，
  只是沒有被採用；如果之後要拆成兩個端點，改的地方在這裡。

  `session-training-volume` 沒有自己的端點，它貢獻的是 `GET /workouts/{workoutId}`
  回應裡的欄位（`totalVolumeKgReps`、`totalVolumeDisplay`、
  `uncountedExerciseCount`、`notice`）。

  驗證約束（兩層都要有：契約層擋得到的先擋，domain 才擋得到的留給 domain）：

  | 約束 | 來源 |
  | --- | --- |
  | 次數型動作的 `reps` 必填 | log Rule 3 |
  | `0 ≤ weightKg ≤ 1000`，精度小數點後一位 | log Rule 8 |
  | `0 ≤ reps ≤ 100`，整數 | log Rule 8 |
  | 單一動作的組數 ≤ 20 | log Rule 8 |
  | 單次訓練的相異動作數 ≤ 50 | log Rule 8 |
  | 自體重型動作不需要 `weightKg` | log Rule 3（Example 3.4） |
  | 自訂動作名稱在該重訓者庫內唯一（正規化後） | custom Rule 3 |

  錯誤訊息（僅列有例子明確斷言訊息文字的；沒斷言訊息的不在這裡發明一個，見
  Out of Scope）：

  | 訊息 | 出現在 |
  | --- | --- |
  | 已有進行中的訓練 | log Rule 1（Example 1.2） |
  | 必要參數未提供 | log Rule 3（Example 3.2） |
  | 單一動作最多 20 組 | log Rule 8（Example 8.5） |
  | 單次訓練最多 50 個動作 | log Rule 8（Example 8.6） |
  | 已經有同名的動作 | custom Rule 3（Example 3.1、3.3、3.4） |
  | 內建動作庫已經有這個動作 | custom Rule 3（Example 3.2）、Rule 5（Example 5.2） |
  | 這個動作已經被訓練紀錄使用，只能封存 | custom Rule 4（Example 4.1） |
  | 這個動作已經被訓練紀錄使用，型別不能改 | custom Rule 6（Example 6.1） |
  | 輔助配重超過體重 | volume Rule 2（Example 2.5），200 回應內容，不是錯誤 |
  | 上週同期沒有紀錄／這是你的第一週／上週有 1 個動作因為缺少體重而未計入 | week Rule 4、5、6，200 回應內容，不是錯誤 |

- **Schema**：判準是「跨兩次使用者操作之間，什麼必須還在」——使用者可能隔一段
  時間才記下一組（log Rule 2），所以要持久化。

  | 資料表 | 欄位 | 來源 |
  | --- | --- | --- |
  | `workouts` | `id, lifter_id, started_at, ended_at, status, attributed_date` | log Rule 1、2、6 |
  | `workout_sets` | `id, workout_id, position, exercise_id, weight_kg, reps, duration_sec, distance_m, warmup` | log Rule 3、4、5、9；volume Rule 8 |
  | `exercises` | `id, name, name_normalized, type, unilateral, source, created_by, archived` | custom Rule 1–7 |
  | `lifters` | `id, name, display_unit, timezone` | 播種表；volume Rule 7 |
  | `bodyweight_records` | `lifter_id, recorded_on, weight_kg` | 播種表；volume Rule 6 |

  兩個欄位是刻意存下來而非查詢時現算：`workout_sets.position`（log Rule 5
  的動作區塊靠相鄰判定，沒有順序欄就分不出「臥推 3 組、划船、臥推 2 組」
  與「臥推 5 組、划船」）、`exercises.name_normalized`（custom Rule 3 的
  唯一性索引要建在正規化後的名字上，資料庫的唯一索引只能建在存下來的欄位；
  這也讓 Rule 7 的遮蔽查詢用得到索引）。`workouts` 需要「每個 lifter 最多
  一筆進行中」的唯一性約束（log Rule 1），domain 也要擋，資料庫是最後一道。

  刻意不存的欄位，各有規則依據：`workouts.total_volume`（edit Rule 2、
  history Rule 6 明說不保留快照）、週彙總表（同上，且 week Rule 6 要兩週
  套同一套規則，預先彙總會凍住當時的規則）、`workout_sets.effective_weight`
  （自體重型的有效重量隨體重紀錄而定，存了會跟 volume Rule 6 的不回溯打架）、
  `deleted_at` 軟刪除欄位（edit Rule 3 明說不可復原，軟刪除欄位一旦存在，
  每個查詢都要記得過濾）、`exercises.usage_count`（custom Rule 4／6 的
  「被引用」要看當下的 `workout_sets`，存一份計數要在每次記錄／刪除時維護，
  漏一次就會誤放行或誤擋）、挑選清單的物化檢視（custom Rule 7 的遮蔽在改名後
  要立刻還原，見 Example 7.3，存起來就要處理失效）。

  三份 `@query` story（session-training-volume、workout-history、
  week-over-week-progress）都只讀上面的表，不新增資料表。

- **架構決定**：分層沿用專案既有的 `domain ← application（port/in、
  port/out）← infrastructure／interfaces`（沿用各層既有的套件文件慣例），這一批不改動
  分層本身。是否需要交易邊界（例如「刪除未被引用的動作」與「用同一個動作
  記一組」同時發生時如何一致）規格沒有講清楚該用交易還是接受並補償，屬於
  未解的風險，見 Risks；因為六份 feature 沒有任何「應寄出」「應通知」等
  外部依賴（見下方外部依賴），這一批全部走同步請求／回應，不需要事件驅動。

- **互動**：以操作類型分類的呼叫順序（由既有 mapper／presenter 分工與
  domain 方法表組合推出，非新決定）：

  - 寫入一次訓練的狀態（開始／記一組／結束／修改一組）：HTTP handler 用
    mapper 把請求轉成 command → use case 呼叫 `WorkoutRepository` 讀出（或
    新建）聚合 → 呼叫 `Workout` 對應的方法 → 存回 `WorkoutRepository`。
  - 查詢一次訓練詳情（含總容量）：use case 呼叫 `WorkoutRepository` 讀出
    聚合 → 呼叫 `VolumeCalculator`（需要讀 `Exercise` 的型別與
    `BodyweightRecord` 的歸屬日當下最新體重）→ presenter 組裝回應。
  - 刪除或改動作型別：use case 先呼叫 `ExerciseUsage` 確認有沒有被訓練引用
    → 通過才呼叫 `Exercise` 對應的方法。
  - 歷史清單／週對照查詢：不經過聚合，直接呼叫 `WorkoutHistoryView` 或
    `WeeklyVolumeView` 讀投影（沿用既有「讀寫路徑不對稱」的決定）。

  `ExerciseUsage` 檢查與後續動作之間如果有並發的寫入，呼叫順序要不要放進
  同一個交易裡沒有答案，見 Risks 的 TOCTOU 風險。

## Testing Decisions

**Seam**（驗收測試打在哪一層）

| seam | 是既有的還是新的 | 為什麼是這一層 |
| --- | --- | --- |
| `interfaces/http` 的 router（由 `bootstrap.NewHandler` 組出，`httptest` in-process 打 HTTP，不經過真實 socket） | 既有 | `test/acceptance/steps_test.go` 的 `theServiceIsRunningAtVersion` 已經是這個模式：用 `cmd/server` 實際會用的同一個組裝函式建 router，注入 in-memory／stub 依賴，再用 `httptest.NewRecorder()` 送請求——「這是唯一的入口」的判準下六則 story 全部共用同一個 seam，符合「整個變更的理想數量是一個」。§3 也確認 82 個場景裡 80 個不需要真的資料庫，代表 seam 選在這一層不會被資料庫綁住——用 in-memory adapter 換掉 `WorkoutRepository`／`WorkoutHistoryView`／`WeeklyVolumeView` 等埠即可，router 本身不用換 |

判準依據：優先用既有 seam（沿用 version 故事已經建立的模式）、取最高的那一層
（HTTP router，而不是直接呼叫 domain 或 application）、六則 story 共用同一個
seam。

**每個 Example 的測試層級**

一個場景區塊一行，`Scenario Outline` 算一個、編號寫範圍。`需要資料庫` 只在
斷言要讀到「前一次操作」寫入的東西時才打勾；其餘場景在同一次操作內就驗完，
用 in-memory adapter 一樣紅、一樣綠。

### log-a-workout（19 個場景區塊）

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
| 6.1 跨午夜結束仍歸屬開始日 | 結束不動 `attributedDate` | ✗ |
| 6.2 換時區不改變歸屬日 | — 純讀取 | ✗ |
| 7.1 同一天兩筆相加 | 當日加總 | ✗ |
| 9.1 dropset 記為兩組 | 容量公式 | ✗ |

### edit-a-logged-workout（8 個場景區塊）

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1 修改一週前的某一組 | 就地取代（position 不變） | ✗ |
| 1.2 補記漏掉的動作 | 附加到尾端 | ✗ |
| 1.3 把已記錄的組改標熱身 | 就地取代 | ✗ |
| 2.1 改重量後該次總容量變大 | 容量重算 | ✗ |
| 2.2 改重量後該週總容量變大 | 週彙總重算 | ✗ |
| 2.3 改標熱身後扣除但留在明細 | 容量重算（排除熱身） | ✗ |
| 3.1 查詢刪除摘要且訓練仍在 | 摘要投影 | ✗ |
| 3.2 刪除後歷史與週彙總都消失 | — | ✓ 要讀得到刪除前一次操作寫入的資料已經消失 |

### session-training-volume（16 個場景區塊）

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1 三組相加 | 容量公式 | ✗ |
| 1.2 只有熱身時為零 | 容量公式邊界 | ✗ |
| 2.1–2.4 各型別的有效重量（Outline） | 有效重量對照表（四個分支） | ✗ |
| 2.5 輔助配重不小於體重 | 有效重量下限夾擠 | ✗ |
| 3.1／3.2 單邊 ×2（Outline） | 單邊乘數 | ✗ |
| 4.1 熱身不計入 | 容量公式（排除熱身） | ✗ |
| 5.1 無體重時為零並標示 | 未計入計數 | ✗ |
| 5.2 混合時只算得出來的 | 未計入計數 | ✗ |
| 6.1 舊訓練用當時的體重 | 體重的日期解析 | ✗ |
| 6.2 新訓練用新體重 | 同上 | ✗ |
| 7.1 磅顯示總容量 | 單位換算 | ✗ |
| 7.2 兩種單位來回不漂移（Outline） | 單位換算的可逆性 | ✗ |
| 8.1 只有棒式 | 計時型不計入 | ✗ |
| 8.2 農夫走路 | 距離型不計入 | ✗ |
| 8.3 混合時只算次數型 | 未計入計數 | ✗ |
| 8.4 全是計時型 | 未計入計數邊界 | ✗ |

`2.1–2.4` 的有效重量對照表掛了四個分支，加上 `2.5` 的夾擠、`8.x` 的兩種
不計入型別——這是整批裡邏輯最集中的一塊，六種型別各要一個 unit 測試。

### workout-history（8 個場景區塊）

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1／1.2 依開始時間排序 | 排序鍵（開始時間，非歸屬日） | ✗ |
| 2.1 組數含熱身而總容量不含 | 投影計算 | ✗ |
| 3.1 第二頁只回剩下的 | 分頁切片 | ✗ |
| 3.2 不滿一頁時就是全部 | 分頁邊界 | ✗ |
| 4.1／4.2 沒有訓練時回空清單 | — | ✗ |
| 5.1 進行中不在清單而另外回傳 | 分流 | ✗ |
| 5.2 結束後回到清單的對應位置 | 同上 | ✗ |
| 6.1 修改過的顯示新的總容量 | — 純投影 | ✓ 要讀得到前一次修改操作寫入的值 |

### week-over-week-progress（12 個場景區塊）

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

### custom-exercise-library（19 個場景區塊）

| 場景 | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1 別人的自訂動作看不到 | `ExerciseCatalog`：可見性過濾 | ✗ |
| 1.2 自訂排在內建之後 | `ExerciseCatalog`：排序 | ✗ |
| 3.1／3.3／3.4 正規化後同名失敗（Outline） | `ExerciseNameNormalizer`（空白、大小寫、全半形） | ✗ |
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
| 7.1 同名時只出現自訂的 | `ExerciseCatalog`：遮蔽 | ✗ |
| 7.2 沒有同名自訂的看得到內建 | 同上 | ✗ |
| 7.3 改名讓位後內建出現 | 同上（遮蔽是算出來的，不能是存下來的旗標——見備註） | ✗ |

`7.3` 值得單獨看：它證明遮蔽不能是存下來的旗標——改名之後被蓋住的內建動作要
自己回來。這一格如果用「建立時標記 hidden」實作，會通過 7.1 與 7.2 而在 7.3 紅。

整批 82 個場景區塊裡只有 2 個需要資料庫（edit 3.2、history 6.1），且兩個都是
`Given` 已經是一次真正的前一次操作；其餘 80 個在 in-memory adapter 下一樣紅、
一樣綠。

**Prior art**

驗收測試層級（本批直接沿用的形狀）：

- `test/acceptance/steps_test.go`、`test/acceptance/world_test.go`——每個場景
  用 `bootstrap.NewHandler` 組出真正的 router，`world` 把 router 與
  `httptest.ResponseRecorder` 放進 scenario 的 `context.Context`（不用套件
  變數，因為 `godog` 設了 `Concurrency: 4`），`Then` 步驟解碼成 `apigen`
  產生的型別而不是 map，改個欄位名字會讓測試編不過而不是靜默通過。

Unit 測試層級（本批要新增 domain／application 測試時可以照抄的形狀，兩個都是
「拿一個小埠自己寫 stub，不用 mocking 套件」）：

- `internal/application/usecase/query/get_version_test.go`——手寫一個
  `stubProvider` 頂替 `out.VersionProvider`，因為埠只有一個方法；用
  `errors.Is`／`errors.As` 斷言錯誤有沒有被正確分類與包裝。
- `internal/interfaces/http/version_test.go`——手寫一個 `stubVersionService`
  頂替驅動埠，直接呼叫 `Server` 的方法斷言回應型別與內容，不透過真的 HTTP
  傳輸。

這一批需要的 domain 聚合（`Workout`、`Exercise`）與 domain service
（`VolumeCalculator` 等）目前在專案裡沒有任何前例——現有測試都停在
「沒有業務規則」的 `GetVersion` 這一層，六則 story 是第一批真正帶規則的
domain 程式碼。`testify` 已經是專案依賴（`.golangci.yml` 特地開了
`testifylint` 檢查 assert／require 的誤用，`mocks/` 也是用 testify 樣板產生
的），但現有手寫測試（`get_version_test.go`、`version_test.go`）在埠只有
一兩個方法時選擇手刻 stub、只用標準庫斷言，沒有走 testify 那條路。這一批
第一個要驗多個型別分支（volume Rule 2 的六種型別）的測試，該延續「埠小就
手刻 stub」還是改用 table-driven ＋ testify，沒有前例可以照抄，見 Out of
Scope。

## Risks

| 風險 | 從哪條規則長出來 | 影響什麼 |
| --- | --- | --- |
| 併發：訓練唯一性。兩個 `POST /workouts` 同時通過「沒有進行中的訓練」檢查 | log Rule 1 | schema（要不要 partial unique index on `status='進行中'`）。單行程用應用層鎖就夠，多副本必須落到資料庫；兩種都可以，看部署形態 |
| 併發：動作名稱唯一性。兩個 `POST /exercises` 同時通過正規化後的同名檢查 | custom Rule 3 | schema（`(created_by, name_normalized)` 唯一索引）。跟上一條同一類，但多一個轉折：唯一性建在正規化後的欄位上，所以那個欄位必須存下來（見 Implementation Decisions · Schema） |
| 重送／冪等：無法分辨。`POST /sets` 重送會多一組，而領域上無法分辨重送與真的又做了一組 | log Rule 5（同一動作可以出現多次，各自獨立）直接造成 | API（要不要冪等鍵）。這是規則選擇的後果不是實作難題：Rule 5 把「去重」這條路關掉了。健身房收訊差、使用者連按兩次，都會讓總容量默默偏高 |
| 正規化必須只有一份實作。custom Rule 3 在寫入時用它擋同名，Rule 7 在讀取時用它決定遮蔽 | custom Rule 3 ＋ Rule 7 | 兩份實作會產生一種很難查的矛盾：新增被擋下（寫入端認為同名），但清單沒有遮蔽（讀取端認為不同名），於是使用者看到一個他建不出來、也看不到的名字 |
| TOCTOU：引用檢查與寫入之間。「這個動作沒被引用 → 刪掉」與「記錄一組用這個動作」同時發生 | custom Rule 4（Example 4.3）、Rule 6（Example 6.2） | 交易邊界。刪除動作與記錄組是兩個不同的聚合，跨聚合的一致性要嘛用交易，要嘛接受並補償。規格沒說哪一個 |
| 讀取路徑的重算成本。每次查歷史／週比較都要從原始組重算容量；自體重型每一組要解析一次體重紀錄；每一組還要查一次動作型別 | edit Rule 2、history Rule 6、week Rule 6、volume Rule 2／6 的「不存快照」 | 分層與索引。N+1 是最直接的形式。不存快照是規格要求，解法只能在讀取那一側（一次撈進來、或可失效的快取），不能回頭去存 |
| 時區與歸屬日。`attributedDate` 由 `startedAt` ＋ lifter 時區算出後凍結；日光節約的地區會有不存在或重複的當地時刻 | log Rule 6、Shared 2 | 時間換算的實作。台北沒有日光節約，但 `timezone` 是每個 lifter 一欄，規格上允許任何時區 |

只指出風險，不寫死解法：前兩條要 Redis 還是 DB 唯一索引，需要知道實際流量與
既有基礎設施，而那些不在 `.feature` 裡。第三條最該現在知道——它不是實作難題，
是規則選擇的後果。要不要處理（加冪等鍵）是產品決定，如果決定不處理，那應該是
一個寫下來的決定，不是一個沒人發現的洞。

## Out of Scope

**明確拒絕的**

| 拒絕了什麼 | 為什麼 |
| --- | --- |
| 進行中的訓練閒置 30 分鐘後自動結束 | 會憑空生出一個「系統」actor 與一整組沒人問過的規則（自動結束的門檻、通知方式……）。見 log-a-workout 備註 |
| 刪除訓練用軟刪除加 30 天垃圾桶 | 每一個查詢都要記得過濾已刪除的資料，漏掉一處就是總容量默默算錯，而使用者無從發現、無法回報。見 edit-a-logged-workout 備註 |
| 自體重型動作的重量欄直接讓使用者自己填 | 違反「不用自己按計算機」的產品定位：輔助引體向上會逼使用者自己算「體重 − 輔助配重」。見 session-training-volume 備註 |
| 週對照用「本週至今 vs 上週整週」比較 | 會讓週一到週五永遠顯示退步——一個固定給錯答案的進步指標，那個負數是日期的訊號不是進步的訊號。見 week-over-week-progress 備註 |
| 自訂動作可以分享成公開庫 | 需求給的動機（各家健身房叫法不同）正好指向反面：共用庫會很快塞滿同義詞，那正是自訂功能要解決的問題本身。見 custom-exercise-library 備註 |

**回 CLARIFY 補問**（規格沒講，推不出來）

| 缺口 | 影響什麼定不下來 |
| --- | --- |
| 失敗場景整批缺（bdd-spec 的「三種結果」檢查抓到的）：43 條 `Rule:` 區塊（`.feature` 裡的，含各 Shared 規則自成一塊；example map 的 `#### Rule` 是 39 條）裡 37 條只走了單一結果 | 修改／刪除／查詢一筆不存在的訓練回什麼？在進行中的訓練上做 edit 的操作可以嗎？刪除進行中的訓練之後 log Rule 1 的唯一性怎麼算？三題都沒問過 |
| 以磅輸入的換算沒有任何場景覆蓋 | volume Rule 7 只驗到顯示端（已存的 20.4 kg 在兩種單位下怎麼顯示）。「使用者以磅輸入 45」沒有任何已就緒的命令場景，log 的重訓者播種表一律是公斤。輸入端的換算與捨入完全沒被驗到 |
| 封存的動作還能不能被記錄引用？ | custom Rule 4.2 只說它不出現在挑選清單，沒說 `POST /sets` 帶那個 ID 會怎樣。清單是 UI，API 擋不擋是另一回事 |
| `WorkoutStatus` 只看得到「進行中」「已結束」 | 沒有證據說還有沒有別的（放棄？取消？）。狀態轉移表可能缺一格 |
| history 的頁碼超出範圍（第 99 頁）回什麼？ | Rule 3 只有「有下一頁」與「不滿一頁」兩個例子，兩個都在範圍內 |
| history 的排序在開始時間完全相同時怎麼辦？ | 沒有次鍵。兩筆同秒開始順序就不決定，而「動作區塊應依序恰好為」這類斷言會間歇性紅（godog `Concurrency: 4`） |
| 體重、時區、顯示單位怎麼被設定？ | 三個都在播種表裡，也都有規則用到它們（volume Rule 6、7；log Rule 6），但沒有任何 `When` 建立或修改它們。這跟 `actor.md` 裡「怎麼取得這個身分」空著是同一個洞。**與 `actor.md` 矛盾**：`actor.md`「能做什麼」表卻把「設定體重與顯示單位」列成已確立的能力（引 session-training-volume Rule 6、Shared 3），但這兩個來源都只消費已播種的值，沒有任何 `When` 驅動的設定動作——`actor.md` 這一格讀起來像已經答過，其實沒有；此處不改 `actor.md`，只記下矛盾 |
| log Example 8.1／8.3／8.4 只斷言「操作失敗」，沒有訊息；8.5／8.6 有 | 不一致，是刻意的還是漏了不定下來，IMPLEMENT 會自己編三條訊息 |
| 負重型動作的重量填 0 該成立還是該擋？ | log Rule 8 說「0 到 1000」，但沒有例子用 0。「不帶重量」（Example 3.4）跟「重量 0」不是同一件事 |
| 開始訓練時 `startedAt` 由誰決定？ | 所有場景都用「現在時間為」播種，看不出是伺服器時鐘還是客戶端送的。離線補記會直接撞到這裡 |
| 計時型的秒數、距離型的距離有沒有範圍與精度？ | log Rule 8 只約束了重量與次數。新的兩個欄位沒有任何邊界例子——這是 `timed-and-distance-exercises` 答完之後長出來的新缺口 |
| 內建動作庫初版收錄哪些、誰標六種型別、多久更新一次 | `built-in-library-source` 明記這三件沒有來源。不影響行為規則所以不擋就緒，但是實際要有人做的工作，排期時不能當成沒有 |
| `timed-and-distance-exercises`、`built-in-library-source` 兩題的答案是代答的（信心低） | 它們決定了核心指標（總容量）的邊界與動作庫的所有權模型。應該由產品確認過再開工——若答案翻案，Shared 4 的型別表與 custom Rule 7 的優先權規則都要重寫 |
| 應用層對外的驅動埠要收在幾個服務、每個服務叫什麼名字（例如是不是一個 `WorkoutService`） | 規格只說行為，沒說這批行為要收在幾個服務裡 |
| `Clock`、`IDGenerator` 兩個驅動外埠的確切方法簽章 | 六份 feature 的 `Background` 與四個建立操作證明這兩個埠一定存在，但沒有任何一份產物給過方法名稱與型別 |

**留給 IMPLEMENT**（技術決定，不需要規格授權）

| 決定 | 為什麼現在不決定 |
| --- | --- |
| 「組」的定址用位置還是 SetID | `.feature` 寫「第 2 組」是位置，但位置會因補記而位移。ID 是定址不是行為。建議 SetID，由 `POST /sets` 回傳 |
| 分頁用 offset 還是 cursor | history Rule 3 只說每頁 30 筆，兩種都滿足 |
| `page` 從 0 還是 1 起算 | 純慣例 |
| `POST /workouts` 回 201＋Location 還是 200＋body | 規格只要求拿得到新訓練 |
| 結束訓練用 `POST .../completion` 還是 `PATCH status` | 規格只說「結束」 |
| Problem 的 `type` URI 命名 | RFC 9457 的形狀已定，命名是慣例 |
| 全半形正規化用 NFKC 還是自訂對照表 | custom Rule 3 只給了「臥推（三樓）」與「臥推(三樓)」一組例子，兩種都滿足 |
| 「當日總容量」「那一週的總容量」用哪條讀取路徑觀察 | 這兩個斷言沒有對應的 `When`。建議用 `GET /workouts` 逐筆加總，不新增端點 |
| `warmup` 省略時預設 `false` | 所有播種表都明寫「否」 |
| 組數／動作數上限是常數還是設定 | 上限的實際數值信心低、值會調；先做成具名常數 |
| domain／application 層的 unit test 該延續「埠小就手刻 stub」還是改用 table-driven ＋ testify | 兩種都能驗到規則，專案沒有帶多分支規則的測試前例可以照抄，見 Testing Decisions · Prior art |
