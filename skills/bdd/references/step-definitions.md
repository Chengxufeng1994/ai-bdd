# step definitions：先找再寫

`red` 寫步驟骨架、`green` 補步驟實作、`refactor` 整理步驟時都適用。`implement` 在票的
測試層先紅時也照這份寫。

MUST: 寫新步驟之前 `grep` 既有的 step definition（實測場 `make steps`；它不加 `-v` 只印 `ok`，
這時直接 `grep -rn 'sc\.\(Given\|When\|Then\|Step\)' test/acceptance/`）。同一句話兩種
寫法就是兩支函式，而它們會分岔。

- 一個樣板一支函式；`操作成功`／`操作失敗`／`應付金額應為` 這類跨 feature 的句子
  一定已經有了
- 每個場景的狀態住 world（`context.Context` 裡），不住 package 變數——並行跑的套件
  一個共用變數就是間歇性紅
- 時間、ID、會失敗的 repository 這類「場景要控制的東西」用測試端的替身注入
  bootstrap，走**跟 binary 同一個組裝函式**；不要為測試另開一條接線
- 播種步驟直接寫 store，不經 use case——場景要擺出的狀態（付過款、票已過期）
  不一定是 use case 走得到的
