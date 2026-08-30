# 資料流與轉換

這份文件只回答一件事：**一筆資料從進來到出去，經過哪幾次轉換、每一次歸誰管。**
正要動手寫 operation、mapper、use case 或 presenter 時讀它。

方向感、目錄結構與分層職責在 [ARCHITECTURE.md](./ARCHITECTURE.md)，這裡不重複。

先講一個容易踩到的名詞：**`handler` 在這份文件裡只指泛型型別名**
`usecase.QueryHandler` / `usecase.CommandHandler`。application 的實作叫 **use case**，
住在 `internal/application/usecase/`；interfaces 的入口是 **operation**，是 `Server`
的方法，一個操作一個檔案。同一個字曾經同時當兩層的名字，這裡不再這樣用。

寫入路徑範例用「記錄一次訓練」，**聚合與欄位名稱是示意的**，實際有哪些聚合由
CLARIFY 決定；讀取路徑用的是已經存在的 `/version` 流程。

---

## 寫入路徑 —— 3 次轉換

```
                                                         層
apigen.RecordWorkoutJSONRequestBody                   interfaces
   │  Server 的操作方法取路由參數與 body
   │  ① mapper       純搬移，不會失敗
   ▼
command.RecordWorkout                                 application
   │  只有基本型別，無 json tag
   │  in.WorkoutService.RecordWorkout → service 委派給 usecase/command
   │  ② toDomain     呼叫領域建構子，會失敗 ← 不變式在此生效
   ▼
workout.Workout                                       domain
   │  檢查不變式、改變狀態、發出事件
   │  ← application 透過 out.WorkoutRepository 存檔
   │  ← application 透過 out.EventPublisher 發布
   │  ── 不轉換 ── 只取 w.ID()
   ▼
command.RecordWorkoutResult{ID}                       application
   │  ③ presenter    純搬移
   ▼
apigen.RecordWorkout201JSONResponse                   interfaces
```

第 ③ 次便宜，是因為 **command 只回識別碼**。一旦回傳完整視圖，就多一個 assembler、
多一組欄位對應，而呼叫端拿到的仍是寫側形狀，不是它要顯示的形狀。

## 讀取路徑 —— 2 次轉換

**這條鏈是實際存在的程式碼**，不是示意——寫入路徑與聚合名稱仍然是示意的，等 CLARIFY
決定有哪些聚合。

```
apigen.GetVersionRequestObject                        interfaces
   │  Server.GetVersion（http/version.go）
   │  ① mapper.ToGetVersion（本例轉換零個欄位）
   ▼
query.GetVersion{}                                    application
   │  in.VersionService.GetVersion
   │  service 委派給 usecase/query 的 use case（型別未匯出）
   │  ── 不轉換 ── out.VersionProvider 直接回傳字串
   ▼
query.GetVersionResult{Value}                         application
   │  ② presenter.ToGetVersionResponse
   ▼
apigen.GetVersion200JSONResponse                      interfaces
```

**完全不經過 domain，沒有 assembler。** 這是 CQRS 唯一真正的報酬：讀側可以打去
反正規化的表、快取或搜尋索引，形狀由「呼叫端要顯示什麼」決定。

如果每個 query handler 都是「載入聚合 → 讀欄位」，這個分離只是換了名字。

**這裡的「2 次轉換」是因為 `out.VersionProvider` 直接回傳字串。** 換成由 read model
支撐的查詢就是 3 次：`port/out` 依規則不能回傳 use case 的 result 型別（見
`internal/application/doc.go` 的 MUST NOT），所以 use case 還要自己做一次
`readmodel → result`——第一個真正接 read model 的查詢不該被這第三次轉換嚇到。

---

## 依賴方向

```
interfaces ──▶ application/port/in ──▶ usecase/query ──▶ application/port/out
    │                                        │                     ▲
    │                                        ▼                     │ 實作
    └──▶ application/errors               domain            infrastructure
```

interfaces 也直接 import `usecase/query`，取 `query.GetVersion` 與
`query.GetVersionResult`；圖上走 `port/in` 那條，是因為驅動埠的簽章本來就用到它們。

`application/service` 不在這條鏈上：它實作 `port/in`，持有 `usecase` 的泛型 handler
介面，一個能力委派一次。

**箭頭一律指向內層。** infrastructure 看似反向，但它反向的是「實作」而非「依賴」
——它 import application 的 port 介面，application 不 import 它。

**`port/out` 不回傳 use case 的 result。** 它回的是領域型別，或由儲存端形狀決定的
read model。這條規則就是「result 與它的 query 放同一個檔案」不會長出循環的原因：
`usecase/query` import `port/out` 是必然的——use case 要呼叫埠；反方向那條邊被這條
規則擋掉，而循環需要兩條邊才成立。完整論證見 `internal/application/doc.go`。

---

## 四段轉換

| # | 轉換 | 角色 | 位置 | 簽章 |
| --- | --- | --- | --- | --- |
| 1 | 協定請求 → command／query | **mapper** | `interfaces/<proto>/mapper/` | `func(X) Y` |
| 2 | command → 領域模型 | 建構 `toDomain` | `application/usecase/command/` | `func(X) (Y, error)` |
| 3 | 領域模型 → result | **assembler** | `application/assembler/`（尚未建立） | `func(X) Y` |
| 4 | result → 協定回應 | **presenter** | `interfaces/<proto>/presenter/` | `func(X) Y` |

**result 與它的 command／query 放同一個檔案**，所以沒有 `dto/` 這個套件；第 3 列的
`assembler/` 也還不存在——它的簽章是 `domain → result`，而 `domain/` 依規則是空的，
沒有函式可寫就沒有套件。兩者的理由都在 `internal/application/doc.go`。

**入向與出向的性質不對稱。** 入向是建構：呼叫領域建構子，不變式在此生效，會失敗。
出向是投影：讀取已經成立的狀態，不會失敗。看到 assembler 回傳 error，先問它是不是
偷偷做了別的事。

**mapper 與 presenter 是方向不同的同一件事**，都是機械翻譯。分兩個名字是為了一眼
看出方向，不是因為性質不同。

### 1. mapper

```go
// interfaces/http/mapper/workout_request_mapper.go
func ToRecordWorkout(req apigen.RecordWorkoutJSONRequestBody) command.RecordWorkout {
	sets := make([]command.SetInput, len(req.Sets))
	for i, s := range req.Sets {
		sets[i] = command.SetInput{Reps: s.Reps, WeightKg: s.WeightKg}
	}
	return command.RecordWorkout{ExerciseID: req.ExerciseId, Sets: sets}
}
```

### 2. toDomain

```go
// application/usecase/command/record_workout.go
func (c RecordWorkout) toDomain() (*workout.Workout, error) {
	sets := make([]workout.Set, 0, len(c.Sets))
	for i, s := range c.Sets {
		weight, err := workout.NewWeight(s.WeightKg) // 驗證在領域建構子裡
		if err != nil {
			return nil, fmt.Errorf("set %d: %w", i, err)
		}
		sets = append(sets, workout.NewSet(s.Reps, weight))
	}
	return workout.New(c.ExerciseID, sets)
}
```

**驗證屬於領域建構子，這一段只負責串接與標註位置。** 這裡出現 `if weightKg < 0`
就是把規則搬出了領域，那條規則從此有兩個家，而它們會漂移。

### 3. assembler

```go
// application/assembler/workout_assembler.go
func ToRecordWorkoutResult(w *workout.Workout) command.RecordWorkoutResult {
	return command.RecordWorkoutResult{
		ID:            w.ID(),
		TotalVolumeKg: w.TotalVolume().Kilograms(), // 領域算的，不在這裡算
	}
}
```

三條規則：**不回傳 error**（聚合已是有效狀態）、**不呼叫 port**（要更多資料由
use case 取好傳進來）、**不做計算**（算式出現在這裡代表業務規則跑錯層）。

簡單到只有 `command.X{ID: w.ID()}` 時直接寫在 use case 裡，不必為它開檔——這也是
目前還沒有 `assembler/` 的原因之一。

### 4. presenter

```go
// interfaces/http/presenter/workout_presenter.go
func ToRecordWorkoutResponse(r command.RecordWorkoutResult) apigen.RecordWorkout201JSONResponse {
	return apigen.RecordWorkout201JSONResponse{Id: r.ID}
}
```

presenter 可以**丟掉** result 的欄位——各協定要少給就少給。若某協定需要 result 沒有
的欄位，不是那個欄位本來就該加進 result，就是那個協定需要的其實是另一個查詢。兩者
都不是讓 adapter 伸手進 domain 的理由。

格式化長出規則時（依地區換算單位、四捨五入、破千縮寫），presenter 才從純搬移變成
有內容的東西。此時先確認那條規則屬於呈現而非業務——「使用者偏好的單位」若系統要
記住並用來比較，那是 domain。

### 不要用反射式 mapper

copier、mapstructure 這類靠欄位名自動對應的工具，在這四個位置都是壞主意：欄位改名
時手寫版**編譯失敗**，反射版**安靜留一個零值**，而那個零值會一路流進領域。

---

## 三層驗證，各管各的

| 層 | 檢查什麼 | 失敗 | 誰寫 |
| --- | --- | --- | --- |
| interfaces | **形狀**：JSON 合不合法、必填在不在、型別對不對 | 400 | `OapiRequestValidator` 依 `openapi.yaml` 自動擋掉，不必手寫 |
| application | **前提**：識別碼格式、權限、資源存不存在 | 401／403／404 | 手寫 |
| domain | **意義**：業務規則、聚合不變式 | 422（由 kind 對應） | 手寫 |

**領域驗證不可以在邊界重複一份。** HTTP 層若也檢查「重量不可為負」，這條規則就同時
活在兩個地方，而它們一定會漂移——通常是改了領域卻忘了邊界，於是同一個請求在某個
協定通過、在另一個被擋。

口訣：**邊界驗證「形狀」，領域驗證「意義」。**

表裡「失敗」那欄寫的是呼叫端最後看到的狀態；application 與 domain 兩列本身回的是
kind，不是狀態碼，對應由 adapter 做——見下一節。第一列的 400 只有 interfaces 層產得
出來，這也是 kind 清單裡沒有 400 的原因。

---

## 錯誤的回流

```
errors.Error{Kind, Err}（application/errors）
   │
   │  use case  以 %w 包裝加脈絡，只分類、不翻譯
   ▼
   ├─ interfaces/http/errmap   kind → HTTP status ＋ RFC 9457 Problem
   ├─ interfaces/grpc/errmap   kind → gRPC status code
   └─ interfaces/cli/errmap    kind → 離開碼 ＋ 訊息
```

| kind | HTTP | gRPC | GraphQL | CLI |
| --- | --- | --- | --- | --- |
| NotFound | 404 | NOT_FOUND | NOT_FOUND | 4 |
| Invalid | 422 | INVALID_ARGUMENT | BAD_USER_INPUT | 2 |
| Conflict | 409 | ABORTED | CONFLICT | 3 |
| Unauthorized | 401 | UNAUTHENTICATED | UNAUTHENTICATED | 5 |
| Forbidden | 403 | PERMISSION_DENIED | FORBIDDEN | 6 |
| Unavailable | 503 | UNAVAILABLE | SERVICE_UNAVAILABLE | 7 |

application **永遠不回傳傳輸層代碼**。對應表各協定自己維護，**分類只有一份**——這是
同一個失敗不會在 HTTP 是 404、在 gRPC 卻是 500 的原因。kind 定義在
`internal/application/errors`，HTTP 那一欄實作在 `internal/interfaces/http/errmap`；
只有 HTTP 那一欄是程式碼，其餘三欄先寫在這裡，好讓第一個接上的協定沿用同一份分類。

### 先認得，再兜底

adapter 的順序是**規則，不是風格**：

1. 用標準庫的 `errors.As` 找出 `errors.Error`——不是型別斷言。use case 一路用 `%w`
   包裝，到 adapter 手上的已經不是 `errors.Error` 本身；型別斷言只看得到最外層那個
   殼，找不到就把每個已分類的失敗都變成 500。
2. 認得的 kind 選對應的 status。
3. **其餘一律 500，配一份不含原始訊息的通用 body。**

反過來寫——「認得的翻譯掉，其餘原樣往外送」——讀起來等價，實際上是第一個沒分類的錯
誤出現時就把 `err.Error()` 交給了呼叫端。`errmap_test.go` 用整份 JSON 做斷言守住這
條兜底路徑，而不是只檢查某個欄位。

`KindUnclassified` 是零值也是同一件事的一部分：忘了設 kind 會掉進 500，而不是悄悄
變成 NotFound。

### 為什麼 kind 只有這六個

沒有對應 400 的 kind。上面那張三層驗證表已經把「形狀」判給 interfaces 層，請求在
`OapiRequestValidator` 就被擋掉，根本到不了 use case——application 產不出「我看不懂
這個請求」。成功狀態（200／201／204）也不在這裡分類，它們來自 operation 的合約與
presenter。照著 `components/responses/` 把清單「補齊」，只會多出永遠沒有程式碼能設
定的值。

**目前每個 kind 在 `/version` 都還是 500。** 理由與現狀寫在
`internal/interfaces/http/version.go` 的 `GetVersion` 註解裡，這裡不重複——下一個
宣告多種失敗回應的 operation 出現時，只有一份說法要改。

---

## 交易邊界

**一個 command ＝ 一個交易 ＝ 一個聚合。**

交易由 application 開啟與提交，透過 port/out 取得（domain 不知道有交易這回事）。
需要同時改兩個聚合時：

- 兩者其實是同一個一致性邊界 → 應該合併成一個聚合
- 兩者可以最終一致 → 第一個發出領域事件，第二個由事件處理器跟進

跨聚合的單一交易看似方便，實際上是聚合邊界劃錯了的徵兆。

---

## 完整 DTO 的正確歸屬

application 的 result 是**不完整**的——每個只帶那個 use case 需要的欄位。真正需要
**完整雙向轉換**的是持久化模型，它是 repository adapter 的私有實作：

| | application result | persistence model |
| --- | --- | --- |
| 完整性 | 不完整，只帶所需 | 完整，全部欄位 |
| 方向 | 單向 | 雙向 |
| 數量 | 每個 use case 一個 | 每個聚合一個 |
| 可見性 | 公開 | 私有於 repository |

在 application 放一個完整鏡射聚合的 DTO，會重新製造它本該防止的耦合，並且很快變成
大家實際傳遞的模型——**貧血領域模型通常就是這樣長出來的**，不是有人決定要貧血。

再往旁邊一格是 **read model**：為了顯示而反正規化的投影，有自己的表或快取，由
`port/out` 的專屬 reader 讀出來。跟 result 的分界可以檢查——**read model 不依附任何
單一查詢，好幾個查詢可以各取它的一片；result 是某一個 use case 的輸出形狀，別無其
他。** `query.GetVersionResult` 屬於後者。

---

## 還有第 5 次轉換，但它被關起來了

infrastructure 內部的「DB row → domain」。它不算進摩擦，因為對 application 不可見
——**前提是沒讓 ORM 的 struct 兼任領域物件**。一旦兼任，這次轉換會反向把資料庫的
形狀壓進領域。

## 剩下的省不掉

| 想省 | 代價 |
| --- | --- |
| ①④ 直接拿 `apigen` 型別當 command／result | application 依賴 HTTP，其餘協定無法共用 use case |
| ② command 直接放領域值物件 | interfaces 得建構領域物件，等於在邊界做驗證 |
| ③ command 回傳完整視圖 | 多一個 assembler，且形狀是寫側的 |

第一項值得記住：**這是為多協定付的保費**。若確定永遠只有 HTTP，直接拿產碼型別當
command 是完全正當的簡化——省掉的正是這一項。
