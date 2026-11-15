# `spec.md` 技術部分的骨架

照抄即可。讀這份的時機：第 6 步，需求部分已經 `已對焦`、準備寫技術部分的時候。

技術部分接在需求部分的 `## Success Metrics` 之後、`## 回 DISCOVERY 補問` 之前，七節，順序固定。
它回答「用什麼做、怎麼驗」；要做成什麼在需求部分，怎麼排、哪裡危險、能不能並行在 `plan.md`
（`../../plan/references/plan-format.md`）。

```markdown
## Tech Stack
## Commands
## Project Structure
## Code Style
## Boundaries
## Implementation Decisions
## Testing Decisions
```

空著的節保留標題，寫為什麼空（「（無）——還沒有程式碼」）；空清單與沒想過在檔案上長得一樣。

---

## 技術部分怎麼寫

上面骨架七節的寫法。

**每份 spec 自成一體。** `## Tech Stack`、`## Commands`、`## Project Structure`、`## Code Style`
跟這個 feature 無關，每份照樣寫——讀者不必去翻別的檔。專案自己有這些文件（`CLAUDE.md`、
`docs/ARCHITECTURE.md`、Makefile）時照它們寫，標 `← 專案 <路徑>`；有 `docs/TECH.md`（前面的 feature
做完時 `archive` 寫的現況）就從它抄過來，逐條確認現在還成立再寫，標 `← 專案 docs/TECH.md`。這次
要改的條目照常標它自己的來源——`archive` 做完時照它更新 `docs/TECH.md`。它們是寫的當下的快照：專案之後改了，舊 spec 不會跟著改。

**技術部分改了要留紀錄**：票已經切出去之後，任何動到 seam、介面、schema 的改動都在
`### 版本修訂歷史` 記一列——實作者可能已經照舊的做了一半。

## 來源標記：四選一

`## Tech Stack`、`## Commands`、`## Project Structure`、`## Code Style` 的每一條、
`## Boundaries` 三張清單的每一條、
`## Implementation Decisions` 的每條決定、`## Testing Decisions` 的 Seam 決定與框架，
行尾都要標來源，四選一：

| 標記 | 意思 |
| --- | --- |
| `PRD §x` | 同目錄 `discovery.md` 的 `## 原文` 這樣寫（引用它的章節或段落）——需求方帶來的原始輸入的逐字副本 |
| `Q-<n>` | 來自 `discovery.md` `## Open Questions` 某一題的答案（通常是技術面向、`歸屬` 指向這一節的題），`<n>` 對應那一列的 `Q` 欄 |
| `專案 <路徑>` | 從專案既有的檔讀來的事實（`← 專案 Makefile`、`← 專案 docs/ARCHITECTURE.md`）——Commands、Project Structure、Code Style 多半是這種 |
| `推論` | 以上都不是，是從別的地方推導出來的 |

寫法：行尾加 `← <標記>`，例如：

```markdown
- **動哪些模組**：只新增計費模組；出場閘門沿用既有控制介面 ← Q-9
```

MUST: 這個標記不是裝飾。少了它，「這個決定是從某條 FR 或某題答案推出來的」跟
「這個決定是寫的人順手定的」在 `spec.md` 上長得一模一樣——完整範例見本文最後的
〈範例：技術部分〉。整份文件的來源欄全是 `推論` 時，代表 DISCOVERY 的技術那一輪沒有真的
發生。`推論` 可以用，但每一條同時列進 `## Further Notes` 的「推論的技術決定」。

**已知缺口**：技術部分的 `← Q-<n>` 目前沒有 script 檢查——`check_spec.py` 只查 `## User Stories`
底下的 FR／AC。等 `review` 做鏈稽核時接手。

`## 回 DISCOVERY 補問` 與 `## Further Notes` 不在這條規則的範圍內——它們記的正是還沒有來源的東西。

`## Testing Decisions` 只有 **Seam** 的決定在範圍內，**「每個 Example 的
測試層級」那張表不在**——它的單位是場景，不是決定，而且它的標記永遠只可能是
`推論`：沒有任何 PRD 或澄清答案答得出「這個場景需不需要資料庫」。一格永遠
同一個值就沒有鑑別力，鑑別力正是下面這條 MUST 說這個標記存在的理由。

在範圍內的那條 **Seam** 決定，標記寫在那張表的 `來源` 欄，只寫標記本身、
不加 `←`——那一欄就是「行尾」在表格裡的位置。

唯一的豁免：**值本身就是在宣告「這裡沒有東西」（不論後面有沒有接一句為什麼
空），那就已經說明了來源狀態**——「沒問過」、「（無資料）」、「（無）」都是
這種值，不再疊 `←` 標記。其餘每一欄都要標。


---

## Tech Stack

四個具名區塊，每條照樣標來源：

```markdown
## Tech Stack

**語言與框架**
- Go 1.23，驗收層用 godog ← PRD §1

**主要依賴**（名字與版本）
- oapi-codegen v2.4，從 `api/openapi.yaml` 產 HTTP 介面 ← 專案 go.mod

**不准引入的依賴**
- 不引入需要商業授權的套件 ← Q-6

**既有資產要沿用的**
- HTTP 層沿用既有的 chi router，不另起一套 ← 推論
```

**跟 `## Assumptions / Constraints` 怎麼分邊：限制的來源在不在我們的控制範圍內。**

| 這一條 | 進哪一節 |
| --- | --- |
| 區權會否決車牌辨識 | `## Assumptions / Constraints`——外部決定的，我們改不了 |
| 已發包的閘門硬體只吃某種訊號 | `## Assumptions / Constraints`——外部既成事實 |
| 不准引入需要授權費的套件 | `## Tech Stack`——我們自己的規矩 |
| 沿用既有的 router 不另起 | `## Tech Stack`——我們自己的選擇 |

判準是**誰有權改它**，不是它有多硬。外部限制改不了，只能記下來；自己的偏好
改得了，所以要寫清楚「現在的決定是什麼」，讓下一個人知道破例要先問誰。

MUST NOT: 在這一節寫「已經決定好的技術事項」——那是 `## Implementation Decisions`。
判準：這一條是在說**這個 feature 要怎麼做**（決定），還是在說**任何做法都不准
越過的線**（偏好與禁令）？

## Commands

完整、可以直接貼進終端機的命令，帶旗標，不是工具名：

```markdown
## Commands

- 建置：`go build ./...` ← 專案 Makefile
- 單元測試：`go test -race ./internal/... ./pkg/...` ← 專案 Makefile
- 跑場景：`go test -count=1 ./test/acceptance/ -v -godog.format=progress -godog.tags='<tag 運算式，或用逗號>'` ← 專案 Makefile
- 閘門：`make verify` ← 專案 Makefile
- 產碼：`make gen` ← 專案 Makefile
```

`跑場景` 要寫出 runner 的 tag 語法（「或」怎麼寫），並帶會印出場景數的旗標——寫錯的運算式常常
選到 0 個場景、照樣成功結束。`跑場景` 與 `閘門` 兩行一定要有——`task-breakdown` 把它們抄進每張票的 `**驗證**`，`implement`
照著跑。專案裡找不到閘門命令就寫「（沒有——要問使用者）」，並在 `## 回 DISCOVERY 補問` 記一條；
不要猜一個看起來合理的命令。

## Project Structure

程式碼、測試、文件各住哪，一行一個目錄加一句話。只寫這個 feature 會碰到、或新程式碼該放的
那幾層：

```markdown
## Project Structure

- `internal/domain/` → 業務規則，不依賴任何框架 ← 專案 docs/ARCHITECTURE.md
- `internal/application/` → use case 與埠 ← 專案 docs/ARCHITECTURE.md
- `internal/interfaces/http/` → handler、mapper、presenter，只做轉接 ← 專案 docs/ARCHITECTURE.md
- `test/acceptance/` → godog 的 step definition ← 專案 docs/ARCHITECTURE.md
- `features/` → 可執行的 `.feature` ← 專案 docs/ARCHITECTURE.md
```

## Code Style

命名與格式的慣例，**加一段專案裡真實的程式碼當例子**——一段真的程式碼勝過三段描述。
例子從專案既有的檔挑，標它的路徑；不要為了示範自己寫一段。

```markdown
## Code Style

- 錯誤包成 `apperrors.Error{Kind, MessageKey, Where}.Wrap(err)`，不直接回傳底層錯誤 ← 專案 internal/application/errors/
- 介面命名用角色名（`VersionService`），不加 `I` 前綴 ← 專案 internal/application/port/in/

例子（`internal/application/usecase/query/get_version.go`）：

    func (h getVersionHandler) Handle(ctx context.Context, _ GetVersion) (GetVersionResult, error) {
        v, err := h.versions.Version(ctx)
        if err != nil {
            return GetVersionResult{}, apperrors.Error{
                Kind:       apperrors.KindUnavailable,
                MessageKey: "version.unavailable",
                Where:      "usecase.GetVersion",
            }.Wrap(fmt.Errorf("read build version: %w", err))
        }
        return GetVersionResult{Value: v}, nil
    }
```

專案是全新的、還沒有任何程式碼時，例子寫「（無——還沒有程式碼）」，慣例照 `## Tech Stack`
的語言社群慣例寫並標 `← 推論`。

## Boundaries

三張清單。**這一節等於每個 feature 自帶一份 agent 約束**。

```markdown
## Boundaries

**總是要做**
- 金額一律用整數分計算，不用浮點 ← Q-7

**要先問**
- 動到既有的出場閘門控制流程 —— 問管委會主委 ← 推論（那是已發包的設備）

**絕對不做**
- 寫入車牌辨識相關的任何欄位 ← PRD §5（區權會否決）
```

**三個分類的判準是誰承擔後果：** 自己做、問了再做、做了要負責的事不做。

| 分類 | 意思 | 判準 |
| --- | --- | --- |
| `總是要做` | 不用問就照做 | 做錯的後果落在這個 feature 內，而且改得回來 |
| `要先問` | 可以做，但動手前要有人點頭 | 做錯會影響這個 feature 以外的東西 |
| `絕對不做` | 不論理由 | 已經有人否決過，或後果不可逆 |

MUST: 每一條照樣標來源。`要先問` 的條目要寫**該問誰**——寫「要先問」而不寫問誰，
等於把一件事掛起來而沒有人接。該問誰要指向誰才算數，判準見 `spec-format.md`
〈Document Overview 的三個部分〉那條 MUST，**不要在這裡抄一份**。

MUST NOT: 把 `## Scope — In / Out` 的「明確排除的能力」搬進 `絕對不做`。
兩者不同：Out of Scope 是**這次不做**，`絕對不做` 是**做了要負責**。
下一期可能會做的東西進 Out of Scope；下一期也不准做的才進這裡。

**空著的分類保留標題並寫為什麼空。** 三張清單都空的 feature 是存在的（純查詢、
沒有副作用），但那要寫出來——空清單與沒想過在檔案上長得一樣。

```markdown
**絕對不做**
（無）——這個 feature 是純查詢，沒有做了要負責的副作用
```

## Implementation Decisions

已經決定好的技術事項。每一條都要指得出理由，不准憑空——`推論` 合法，但推論
本身要指得出是從什麼推來的。標記寫法見 `## 來源標記：四選一`。

- **動哪些模組**：新增什麼、改什麼 ← <標記>
- **介面**：那些模組對外的形狀（函式名、參數、回傳型別） ← <標記>
- **API 契約**：方法、路徑、request、response、錯誤碼 ← <標記>
- **Schema**：資料表與欄位 ← <標記>
- **架構決定**：分層、交易邊界、同步或事件驅動 ← <標記>
- **互動**：跨模組的呼叫順序 ← <標記>

MUST NOT: 在這一節寫具體檔案路徑與程式碼片段——它們過期得比什麼都快。`plan.md` 的
`## Parallel Work` 接合點是例外：合回時要對著檔案解衝突，那裡必須寫路徑。
例外：prototype 產出的、比散文更精確地編碼了某個決定的片段（狀態機、schema、
型別形狀），可以內嵌並註明來自 prototype，只留決定的部分。

MUST NOT: 在這裡做新決定。推不出來的寫進 `## 回 DISCOVERY 補問`（需求的洞）或
`## Further Notes`（留給 IMPLEMENT 的技術選擇），不要順手決定掉——那正是規格失效的方式。

## Testing Decisions

**框架與位置**

- 驗收：godog，step definition 在 `test/acceptance/` ← 專案 Makefile
- 單元：`go test`，table-driven，測試檔跟程式碼同目錄 ← 專案 docs/ARCHITECTURE.md
- 覆蓋率要求：<有就寫數字與來源；沒有就寫「（無）」> 

**Seam**（驗收測試打在哪一層）

| seam | 是既有的還是新的 | 為什麼是這一層 | 來源 |
| --- | --- | --- | --- |
| `application.WorkoutService` | 既有 | 規則全在 domain 與 application，HTTP 只是轉發 | 推論 |

判準見 `SKILL.md` 第 6 步〈決定 seam〉的那張表，**不要在這裡抄一份**。
seam 一旦寫在這裡就會往下傳：IMPLEMENT 照著打，REVIEW 把沒人同意過的 seam
當成 finding。

**每個 AC 的測試層級**

一個 AC 一行；`formulation` 之後它就是 `@example-<n>.<m>` 那個場景。幾個 AC 會收成同一個
`Scenario Outline` 時，編號寫範圍。

| AC | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 2.1 進度不可倒退 | 單調遞增判斷 | ✗ |
| 6.1 查詢回傳完整內容 | — | ✓ 要讀得到前一次寫入 |

`需要資料庫` 的判準：這個 AC 要驗的行為**本身跨越一次操作的邊界**嗎？
不是「有沒有資料要存」——多數場景在同一次操作內就驗完了。

`—` 代表沒有內迴圈（純轉發、純投影），那是正常的。

**Prior art**

專案裡已經有的同類測試，寫路徑與一句話。新測試跟著它們的形狀寫，不另創。

## 範例：技術部分

````markdown
## Tech Stack

**語言與框架**
- （無資料）——需求方沒有指定，技術選型由實作端決定

**主要依賴**（名字與版本）
- （無）——虛構教學專案，還沒有任何依賴

**不准引入的依賴**
- 不引入需要連外網的服務，停車場網路是封閉的 ← PRD §5

**既有資產要沿用的**
- 出場閘門的既有控制介面沿用，不另起 ← PRD §5

## Commands

- （沒有——要問使用者）：虛構教學專案沒有建置與測試命令

## Project Structure

- （無）——還沒有程式碼

## Code Style

- （無——還沒有程式碼）

## Boundaries

**總是要做**
- 金額一律用整數元計算，不用浮點 ← 推論

**要先問**
- 動到既有的出場閘門控制流程 —— 問管委會主委 ← PRD §5（已發包的設備）

**絕對不做**
- 寫入車牌辨識相關的任何欄位 ← PRD §5（去年區權會否決）

## Implementation Decisions

- **動哪些模組**：只新增計費模組；出場閘門沿用既有控制介面、不改它——改它是「要先問」且還沒問 ← 推論
- **介面**：`計算應收金額(入場時間, 出場時間) -> 金額`——從 FR-1 與 AC-1.1–1.3 推來：它們只拿入場到出場的停留時長決定一個金額，沒有第三個輸入 ← 推論

## Testing Decisions

**Seam**（驗收測試打在哪一層）

| seam | 是既有的還是新的 | 為什麼是這一層 | 來源 |
| --- | --- | --- | --- |
| `application.出場結帳` | 新 | 出場閘門原本只有硬體控制介面，沒有既有 seam 可掛；FR-1（訪客計費）與 FR-2（月租戶跳過計費）在這裡合流，一個 seam 覆蓋兩則 story 的全部 AC——比只打 `計算應收金額` 高一層（不然漏掉 FR-2 那條分支），比出場閘門硬體低一層（不必為了驗收多開一個假物件） | 推論 |

**每個 AC 的測試層級**

| AC | 內迴圈（unit） | 需要資料庫 |
| --- | --- | --- |
| 1.1–1.3 費率邊界（29／30／31 分） | 分鐘數換算金額 | ✗ |
| 2.1 月租戶跳過計費 | — | ✓ 要讀到月租名單前一次寫入的資料 |
| 3.1 警報中柵欄常開 | — | ✗ 警報訊號用假物件餵 |

**Prior art**

（無）——虛構教學專案，沒有既有同類驗收測試可比照，新測試不套用任何既有形狀。

（`## 回 DISCOVERY 補問`、`## Further Notes` 略）
````
