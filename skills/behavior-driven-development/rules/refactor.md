# refactor：綠燈保護下整理

判定「這次整理有沒有改到行為」。讀這份的時機：進 `refactor` 入口之前。

| # | 強度 | 規則 | 違反的後果 |
| --- | --- | --- | --- |
| F1 | `MUST` | 只在綠燈下開始：先跑一次範疇內的場景與相關單元測試，全綠才動。整理完再跑同一個命令，結果要相同 | 紅燈下整理，分不出是整理弄壞的還是本來就壞的 |
| F2 | `SHOULD` | 優先整理三樣：重複或只差幾個字的 step、型別轉換的邊界、支撐層（helper／fixture／world）的命名與結構 | 把力氣花在產品程式碼的大改，而最會腐壞的測試支撐層沒人管 |
| F3 | `MUST NOT` | 改變任何場景驗的行為，或回頭改 `.feature`、`spec.md` 來配合整理後的結構 | 那不是重構，是沒人同意過的規格變更 |
| F4 | `MUST` | 整理需要動到指定範疇以外的東西時，停下回報，不順手擴張 | 別張票、別的 feature 被一次「順便」的整理改到，而沒有人為它跑過測試 |

## F2：先整理什麼

| 訊號 | 整理成 |
| --- | --- |
| 兩支 step 只差一個字（`應付金額為` / `應付金額應為`） | 留一支。照 `../../formulation/references/step-grammar.md` 的句型，另一支的場景**不改文字**——所以只能在兩句本來就該是同一句時合併；不是的話回報 |
| 好幾支 step 各自把字串轉成金額、時間 | 抽成一個轉換函式，或 runner 的參數型別 |
| step 裡直接操作 store 擺狀態，到處重複 | 下沉成 fixture／建態 helper |
| helper 的名字是它怎麼做的（`insertRowAndSetFlag`） | 改成它擺出什麼狀態（`aPaidStay`） |

## F3：行為沒變的判準

整理前後，**同一個命令、同樣的通過數、同樣的場景**。少了一個場景、多了一個 pending、
哪個斷言被放寬，都是行為變了。

```
可以：  把三支 step 共用的建態抽成 aStayOf(duration)
不可以：順手把「付清超過 15 分鐘」改成可設定的寬限時間
        ← 15 分鐘是 spec.md 的一條規則，改它要回 discovery
```

## 回報長什麼樣

```
入口：refactor　範疇：visitor-exit-payment.feature
整理前：go test -count=1 ./test/acceptance/ -godog.paths=features/visitor-exit-payment.feature → 6 passed
整理了：金額解析抽成 parseAmount；三支 step 共用的建態下沉到 aStayOf
整理後：同一個命令 → 6 passed；go test ./internal/... → ok
沒動：範疇外的 steps_pass_test.go 也有同樣的重複——超出範疇，回報不處理
建議下一步：red，下一個 slice
```
