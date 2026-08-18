# ai-bdd 藍圖

`ai-bdd` 提供一套**完整的開發流程**，以 BDD 為主軸貫穿。它不是測試工具，是把
「為什麼做」一路傳遞到「線上實際行為」的流程骨架。

BDD 實踐本身的定義（三個實踐、雙迴圈、Gherkin 反模式）見
[bdd-workflow.md](./bdd-workflow.md)。本文件談的是**要建什麼、怎麼切、什麼順序**。

> **這份文件是假設，不是定案。** 六步框架與底下的 skill 清單目前只有一個
> 已實作。它應該被第一批真實產物驗證或推翻，而不是被照著蓋完。

---

## 定位

**完整自建。** 各步驟由 ai-bdd 自己提供，不依賴外部 plugin 或 skill。

代價：環境中既有的 `clarify-loop`、`testing-strategy`、`code-review-and-quality`
等 skill 覆蓋部分相同處境，會與 ai-bdd 競爭觸發權。模型會挑一個，而你不會知道
它挑了哪個。緩解方式是在採用 ai-bdd 的專案停用重疊 plugin，並讓所有
`description` 綁定 BDD 脈絡詞彙。

換得的好處：鏈上每個環節的輸入輸出格式**由本專案自行定義**，因此可以真正閉合。
不是功能更多，是介面一致。

---

## 六步流程

| 步驟 | 這一步在問什麼 | BDD 作法 | 產物 |
| --- | --- | --- | --- |
| **CLARIFY** | 我們到底要解決什麼？哪些還不知道？ | Example Mapping：規則／例子／問題 | `example-mapping.md` ＋ 就緒判定 |
| **SPEC** | 業務讀得懂的驗收條件長怎樣？ | Gherkin，反命令式稽核 | `.feature` ＋ 詞彙表 |
| **PLAN** | 每條情境在哪層驗？先做哪一條？ | 情境分層 ＋ 依行為排序 | 分層表 ＋ 實作順序 |
| **IMPLEMENT** | 讓情境從 red 到 green | step defs ＋ outside-in 雙迴圈 | 程式碼 |
| **VERIFY** | 綠了嗎？紅的是 bug 還是規格過期？ | 跑情境 ＋ 失敗判讀 | 測試結果 ＋ 判讀結論 |
| **REVIEW** | 這變更真的滿足它宣稱的情境嗎？鏈斷了嗎？ | 對照情境驗收 ＋ 鏈稽核 | 審查結論 ＋ 稽核報告 |

PLAN 這一步值得特別說明：BDD 的計畫**按可交付的行為排序，不按模組排**。
「先做完整的登入情境」而非「先做完 User model」——前者每一步都有可驗收的東西，
後者要到最後才知道對不對。

框架相依度：CLARIFY／SPEC／PLAN 與框架無關，可先做完並實際使用；
IMPLEMENT／VERIFY 需要綁定技術棧；REVIEW 部分需要。

---

## 步驟 → skills 對照

**檔案結構是扁平的**（`skills/<name>/SKILL.md`），層級只存在於下面這張表。

| 步驟 | skill | 職責 | 狀態 |
| --- | --- | --- | --- |
| CLARIFY | `bdd-clarify` | 切分 → 規則／例子／問題 → 就緒判定 | **已實作** |
| | `bdd-clarify-loop` | 多輪把紅卡問到收斂；也可單獨使用 | **已實作** |
| | `bdd-clarify-story-splitting` | 沿規則切；九種模式；兩條選法規則 | **已實作** |
| SPEC | `bdd-spec` | 例子 → `.feature` ＋ 詞彙表 | 未實作 |
| | `bdd-spec-review` | 反命令式、conjunction step、情境爆炸稽核 | 未實作 |
| PLAN | `bdd-plan` | 分層與順序的總體決策 | 未實作 |
| | `bdd-plan-layering` | 每條情境該在 unit／integration／e2e | 未實作 |
| | `bdd-plan-ordering` | 按可交付行為排序 | 未實作 |
| IMPLEMENT | `bdd-implement` | 實作流程總入口 | 未實作 |
| | `bdd-implement-step-definitions` | 比對既有 step 再產新的 | 未實作 |
| | `bdd-implement-outside-in` | 雙迴圈 | 未實作 |
| VERIFY | `bdd-verify` | 跑情境並回報 | 未實作 |
| | `bdd-verify-failure-triage` | 真 bug vs 規格過期 | 未實作 |
| REVIEW | `bdd-review` | 審查總入口 | 未實作 |
| | `bdd-review-acceptance` | 這 diff 真的滿足它宣稱的情境嗎 | 未實作 |
| | `bdd-review-traceability` | 鏈稽核：孤兒、漂移、未覆蓋 | 未實作 |

子 skill 是否真的需要，等步驟 skill 跑過真實案例、發現它太長或處境確實可區分時
再拆。**現在全部列在這裡是假設，不是待辦清單。**

---

## 命名與結構的決定

**扁平目錄 ＋ `bdd-` 前綴。** 三個理由：

1. skill 名稱在執行時是**全域扁平**的——`claude plugin details` 列出的是名字，
   目錄層級不存在。分組若不寫進名字就完全看不見。
2. 前綴讓它們在字母排序時自然聚在一起，這是扁平命名空間裡唯一還能保住分組的方式。
3. **`claude plugin validate` 不遞迴。** 巢狀目錄（`skills/<分類>/<名稱>/SKILL.md`）
   在 runtime 有效——mattpocock-skills 就是這樣放的——但 validate 掃不到，
   等於那些 skill 永遠不會被檢查。扁平換來的是每個 skill 都在 CI 的守備範圍內。

通用字眼（`story-splitting`、`failure-triage`）要在 `description` 裡寫明**排除
條件**，否則會在非 BDD 情境被誤觸發。

---

## Token 成本

`claude plugin details` 會估算成本。實測參考值（mattpocock-skills，25 個 skill）：

```
Always-on: ~1,620 tok  每個 session 都付
每個 skill：常駐 ~30–160 tok（名稱＋description）
            觸發時 ~20–3.9k tok（本文）
```

推論：**拆得細不貴，description 寫得囉嗦才貴。** description 每次 session 都付費，
本文用到才付。所以 description 壓成一行觸發條件，厚重內容放本文或參考檔。

---

## 可追溯鏈

```
商業目標 → Story → 規則 → 例子 → Scenario → Step → 程式碼
                                      ↓
   回饋 ← 線上行為 ← 發布 ← 審查 ← 測試執行 ←┘
```

每個箭頭都是一個意圖遺失點。保住這條鏈是 ai-bdd 存在的理由，也是最難被取代的
部分——因為鏈的完整性檢查**幾乎全是機械性的，人只是懶得做**。

`bdd-clarify` 的產物已為此設計：規則與例子有編號（Rule 1、Example 1.1），下游 Gherkin
情境回指這些編號，鏈才接得起來。**編號一旦寫出去就不要重排。**

### 狀態 tag —— 這則 story 走到第幾步

`.feature` 檔在 `Feature:` 上掛**恰好一個**狀態 tag，記錄它在六步流程的哪一站。
這是**唯一的**跨 skill 共用詞彙表；各 skill 不要自己抄一份，會漂移。

| tag | 由誰寫入 | 意義 |
| --- | --- | --- |
| `@draft` | `bdd-spec` | 來源 map 未就緒。可用來對齊理解，**不是**講定的驗收條件 |
| `@ready` | `bdd-spec` | 來源 map 已就緒，三方講定，等實作 |
| `@wip` | `bdd-implement` 開工時 | 正在實作這則 |
| `@review` | `bdd-implement` 全綠後 | 等 REVIEW |
| `@done` | `bdd-review` 通過後 | 六步走完 |

三個設計決定，各有理由：

**掛在 `Feature:` 不掛在場景。** 狀態追蹤的是「這則 story 走到哪一步」，
不是「哪幾條場景綠了」——後者跑一次測試就知道，而**任何重述可觀察狀態的標記都會
漂移**，不一致時錯的是標記，看起來有權威的也是標記。

**必填，而且恰好一個。** 少標一個是可機械偵測的（`grep -L`），
負向旗標（「沒標＝正常」）則分不出「不需要標」與「忘了標」——
第一份真實產出正是在這裡失守：兩則未就緒的 `.feature` 沒有任何標記，
跟已就緒的長得一模一樣。

**由執行那一步的 skill 寫入，不由人維護。** 狀態是流程的副產品，
不是一個要記得更新的欄位。

CI 用 `--tags "~@draft"` 排除未定案的場景——它們不該擋住任何人的 build。

IMPORTANT: 六個 skill 目前只有 `bdd-spec` 存在，所以只有 `@draft` / `@ready`
兩個轉換是實作過的，其餘是**保留的名字**。後面三個 skill 做出來時若發現需要不同
的切分，以那時的實測為準——這張表可以改，但要一次改完，不要各自加自己的。

---

## 進度

- [x] `bdd-clarify`
- [x] `bdd-clarify-loop`
- [x] `bdd-clarify-story-splitting`
- [x] `lab/go/skeleton` — Go ＋ godog 實測場（骨架綠，零業務）
- [ ] 拿健身追蹤需求跑一次 `bdd-clarify`，驗證產物格式
- [ ] 依實跑結果修正格式，再做 `bdd-spec`

順序刻意是「做一個 → 實跑 → 再做下一個」。**產物格式就是介面**：SPEC 讀
example map、PLAN 讀 `.feature`、REVIEW 稽核整條鏈，下游全都依賴上游實際吐出
什麼。在紙上設計那個格式是猜的。

### 實測場

`lab/go/skeleton/` 是六步流程的 dogfooding 對象：Go ＋ godog ＋ 分層
架構（cmd／domain／application／infrastructure／interfaces）＋ 三層測試
（unit／integration／acceptance）＋ `api/openapi.yaml`。

它的規則是**在 CLARIFY 跑完前不寫任何業務程式碼**。題目選「訓練總容量計算」
而非「記錄訓練」，因為前者才有規則與模糊點——自體重動作、單邊動作乘不乘二、
熱身組算不算、dropset 是幾組。詳見該目錄的 README，那裡也寫明了
`bdd-clarify` 的及格標準：**沒問到自體重與單邊動作，就算失敗**。

> `lab/` 不是 plugin 的慣例目錄，不會被載入，對 plugin 行為零影響。

---

## 未決事項

1. ~~**產物存放路徑。**~~ 已定：文件型產物一律寫入使用端 repo 的 `docs/bdd/`，
   一個目錄裝完，刪掉即乾淨。skill 不得寫入其他位置、不得修改專案既有檔案。

2. ~~**`.feature` 的位置。**~~ 已定：留在測試框架的慣例位置（Cucumber 系是專案根
   的 `features/`），是唯一不進 `docs/bdd/` 的產物——它的位置由 runner 決定。
   兩邊靠 `trace.md` 連結，格式待 `bdd-spec` 定義。

3. **IMPLEMENT／VERIFY 的技術棧。** 未定（Python behave/pytest-bdd、
   JS Cucumber.js、Go godog…）。前三步不受影響。

4. **License。** repo 在 `vivotek/` 底下，`plugin.json` 目前不含 `license` 欄位。
