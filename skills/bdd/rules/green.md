# green：最小補碼轉綠

判定「這次補的程式是不是只為了眼前這個失敗」。讀這份的時機：進 `green` 入口之前。

| # | 強度 | 規則 | 違反的後果 |
| --- | --- | --- | --- |
| G1 | `MUST` | 從**看過的那個失敗**出發。失敗跟這個 slice 無關、或被別的錯遮住，先回 red 修訊號 | 補的是想像中的完整功能，不是這個失敗要的東西 |
| G2 | `MUST` | 只補目前這個 slice 的失敗需要的最小行為。預見下一個 slice 會用到的分支、旗標、角色，留給下一個紅燈來要 | 沒有測試驅動的分支沒有人驗過，而它看起來跟驗過的一樣 |
| G3 | `MUST` | step definition 保持薄：匹配步驟、解析參數、呼叫 helper 或斷言。共用的行為下沉到 helper、fixture、domain | step 變成完整劇本，改一個流程要改十支 step |
| G4 | `MUST` | 有分支的規則住 domain，而且先有一個看過紅的單元測試；step definition、handler、mapper 裡不放規則 | 規則散在接線層，場景綠了卻沒有任何單元測試說得出它為什麼對 |
| G5 | `SHOULD` | 連續做幾個 slice 時，每個 slice 各有一次獨立的綠燈，才切下一個 | 「幾個相關的測試好像都過了」——哪一個真的綠了說不清 |

G4 的層級照 `plan.md` 的 `## Testing Decisions`（每個例子在哪一層測）；規則放哪一層、怎麼
接到 seam、測試替身怎麼寫 → `../references/inner-loop.md`。

## G2 的例子

```
slice：@example-8.1 付清後可以出場
失敗：付款後查詢狀態應為「已付清」，得到「未付款」

只補：  付款金額等於應付金額時，把停留標成已付清
不補：  金額不足的拒絕（8.2 的事）、逾時補繳（9.1 的事）、退款、部分付款
```

## G3 的例子

```go
// 薄：只匹配、解析、委派
func (w *world) 訪客付款(amount int) error {
    return w.pay(w.currentPlate, amount)
}

// 厚：step 自己變成一整段測試劇本
func (w *world) 訪客付款(amount int) error {
    w.enter(plate); w.advance(31*time.Minute); fee := w.queryFee(plate)
    if fee != 30 { return errors.New("...") }
    if err := w.pay(plate, amount); err != nil { return err }
    return w.assertPaid(plate)
}
```

步驟按領域名詞分檔（付款、月租卡），不按 `.feature` 的檔名分——同一句話會被好幾個
feature 用到。

## 綠了之後

綠了但留下重複的步驟、寫死的資料或暫時碼時，先做最小的整理再切下一個 slice
（`rules/slice.md` L5）。需要動到這個 slice 以外的結構，那是 `refactor` 入口的事——停下，
在回報裡建議。

## 回報長什麼樣

```
入口：green　範疇：visitor-exit-payment.feature 的 Rule 8
slice：@example-8.1
起點的失敗：狀態應為「已付清」，得到「未付款」
單元測試：domain/parking/stay_test.go TestPay（先紅：undefined: Stay.Pay）→ 綠
重跑：go test -count=1 ./test/acceptance/ -godog.tags=@example-8.1 → 1 scenario passed
範疇內其餘：@rule-8 全部重跑，原本綠的仍綠
動了：domain/parking/stay.go（Pay）、usecase/command/pay.go、steps_payment_test.go
建議下一步：red，@example-8.2
```
