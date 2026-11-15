# green 的落點：內迴圈與接線

`green` 入口補產品程式碼時，規則放哪一層、怎麼接到 seam。讀這份的時機：進 `green` 之前，
與 `rules/green.md` 一起讀。

這裡寫的是實測場（`benchmark/skeleton/go/`）的分層與轉換；專案有自己的架構文件
（`docs/ARCHITECTURE.md`、`docs/DATAFLOW.md`、各層 `doc.go`）時，分層、轉換、錯誤分類都照它，
本文件不另訂。

## 內迴圈：規則先進 domain

`plan.md` 的 `## Testing Decisions` 說了哪些例子在 unit 層。對每一條有分支的規則：

1. 在 domain 寫 **table-driven** 單元測試，列就是 `spec.md` 的 AC 值加上 AC 沒寫但
   規則隱含的邊界（24 小時又 1 分鐘、到期日的最後一秒）——跑，看它因為「還沒有這個
   套件／函式」而紅
2. 寫最少的 domain 程式讓它綠。規則只住這裡；use case 與 handler 裡出現算式就是
   跑錯層
3. use case 的測試只驗**分類**：domain 回的錯變成哪個 `Kind`、哪個 message key、
   帶哪些 Params；repository 失敗變成什麼。用手刻 stub，一個埠一個 struct

MUST: 測試替身要回**副本**，不回活的指標。repository 的 stub 若把同一個聚合指標交出去，
use case 在 `Save` 失敗前做的修改會留在「儲存」裡，測不出「失敗時狀態不變」——
parking 實跑的第一個假紅就是這個。

## 接線：照架構文件的轉換走一遍

依 `docs/DATAFLOW.md`（或專案的同類文件）：埠（`port/out`）→ use case →
driving port 與 service → contract（`openapi.yaml`）→ 產碼 → mapper → operation →
presenter → bootstrap 組裝。每一段都是機械翻譯，**不做決定**；要做的決定都該在
`plan.md` 的技術節，沒有就停下回報，由 `bdd-implement` 記進票檔的 `## 缺口`。

- contract 改了就重新產碼（實測場 `make gen-api`），產出的介面多了方法，`Server`
  不實作就編不過——那是對的，讓編譯器告訴你還缺哪個 operation
- 失敗路徑：operation 用 errmap 認得的 status 選 contract 宣告的 response 型別，
  其餘落 500；不要在 operation 裡 switch `Kind`
- 錯誤訊息的 message key 要有翻譯；場景斷言的 `錯誤訊息應為` 就是那句翻譯
