# bdd-spec-by-example：把 spec.md 寫成業務語意的驗收旅程

**狀態**：已批　**日期**：2026-11-14

## 意圖

`spec.md` 寫完之後，產出一份**業務語意的 feature file**——驗收旅程：幾條可以 demo 的長流程、
只用業務語言、未定的地方就地標記。參考 aixbdd 的 `spec-by-example`。它還沒進開發階段，所以放在
規格目錄 `specs/<date>-<feature>/acceptance/`，不放 `features/`。

就緒判定照舊：需求方看過旅程沒問題、答「就緒」，才翻 `已對焦`，才產**程式版本**的
`features/*.feature`。

### 為什麼

今天能拿來判就緒的只有 `spec.md` 的 EARS 規則與一行一個的 AC，或規則層級、封閉文法的可執行
`.feature`——都不是業務方會想讀的形狀，裡面也沒有旅程。驗收旅程把同一批規則與例子寫成業務方
讀得懂、能 demo 的流程。（需求方說明：不需要真的有 PM／RD／QA 角色去看，目的只是產生業務語意
的 feature file。）

## 目標

- 新 skill `bdd-spec-by-example`：讀 `待對焦` 的 `spec.md`，寫
  `specs/<date>-<feature>/acceptance/<旅程>.feature`；之後 `example-mapping` 照舊提就緒問句
- 每條旅程（`Example:`）用 `@covers-<n>.<m>` 標出它涵蓋的 AC；`check_spec.py` 擋下引用
  不存在 AC 的旅程
- `claude plugin validate .`、`python3 -m unittest discover -s tests`、每個改過的 skill 的
  `audit_skill.py` 全過

## 非目標

- **不執行**驗收旅程。它是規格，不是測試；沒有 step definition、不進 `features/`
- 不改可執行 `.feature` 的格式、tag、封閉文法；`bdd-formulation` 照舊在 `已對焦` 之後寫，只多讀旅程的用語（D2a）
- 不要求每條 AC 都被某條旅程涵蓋——旅程刻意少，只選主流程與少數高差異變種
- 不做 UI 雛形（aixbdd 的 `ui-plan`）、不做前後端拆分（aixbdd 的 `dsl-refine`）
- 不寫 `bdd-review`／`bdd-verify`

## 反目標

- 旅程出現 `spec.md` 沒有的結果、規則或例子——那是不執行的那層 Gherkin 最典型的漂移，也是
  agent 最容易「順手補」的地方
- 寫旅程時當場改規格，或發現洞就自己補一個答案。洞只能回 discovery

## 限制

- 文件繁中、commit 英文 `WHAT:`／`WHY:`／`HOW:`；搬檔用 `git mv`
- `skills/skill-rules/` 是 house style
- SPEC 不訪談（`docs/ai-sdlc.md`「三種模式」）

## 考慮過的做法

| 做法 | 取捨 | 結果 |
| --- | --- | --- |
| **A. 新增一層驗收旅程，放 `specs/<feature>/acceptance/`**（本文件） | 業務方最好讀、是規格而非測試；代價是多一層不執行的 Gherkin，用 `@covers-` 加 `check_spec.py` 擋漂移 | **採用**；位置由需求方指定——還沒進開發，不放 `features/` |
| B. 不新增檔案，在對話裡把例子演成旅程 | 沒有新產物；但看過什麼沒有落檔，下一個 session 看不到 | 捨棄 |
| C. 先產可執行 `.feature`，對著它判就緒 | 判的就是會跑的東西；但業務方讀規則層場景，且要把「`已對焦` 之後才產 `.feature`」倒過來 | 捨棄 |

驗收旅程放在就緒判定**之前**：需求方看旅程沒問題，才 `已對焦`、才產程式版本。

## 決定

### D1. 流程位置

```
bdd-spec 寫 spec.md（待對焦）→ 收尾檢查 → check_spec.py
→ example-mapping 攤地圖、跑四條診斷
→ bdd-spec-by-example 寫 acceptance/*.feature → check_spec.py
   ├─ 有 # [need clarification] → 寫進 spec.md 回 DISCOVERY 補問 → 退回 discovery
   └─ 沒有 → example-mapping 提就緒問句（照舊，需求方一人答）
              ├─ 就緒 → 已對焦 → bdd-formulation 產 features/*.feature
              └─ 不就緒 → 退回 discovery（照舊）
```

`bdd-spec` 第 5 步從「交 `example-mapping`」改成「交 `example-mapping` 跑診斷 → 交
`bdd-spec-by-example` 寫旅程 → 交回 `example-mapping` 提就緒問句」。補問清單非空時都不叫——跟
今天一樣。`example-mapping` 的就緒問句、`已對焦` 的寫入者、喊停帶著風險走的路、代演規則
**全部不變**；就緒問句多一句「驗收旅程在 `acceptance/`」，讓需求方對著旅程答。

### D2. 驗收旅程的形狀

照 aixbdd 的五條規則，改寫成本專案的措辭，放 `references/journey-format.md`：

1. 只寫業務語言：不得有 API、HTTP、欄位名、selector、資料表
2. 一個 `Example:` 是一條可 demo 的長流程（Journey），不是一個 AC 一個場景；數量刻意少
3. `Feature:` 一個驗收主題、`Rule:` 原子、`Example:` 標題描述情境而非重複規則
4. `Background:` 只在多條旅程真的共用 setup 時用
5. 未定的地方就地標 `# [need clarification] <問題>`，緊貼引發疑問的那一行；同一個問題同時
   寫進 `spec.md` 的 `## 回 DISCOVERY 補問`（標明哪條旅程），走 `bdd-spec` 的〈退回 discovery〉。
   本 skill 跟 SPEC 一樣不訪談

本專案多一條：

6. **每個 `Example:` 上方標 `@covers-<n>.<m>`**（一個或多個），列出它走過的 AC。旅程裡
   每一個 `Then` 都要對得到某個被 cover 的 AC；對不到的就是 `spec.md` 沒有的結果——不寫，
   改寫進補問清單

檔名用業務旅程命名（`訪客出場結帳.feature`），不按 story slug——一條旅程常跨兩則 story。
語言：`Feature`／`Rule`／`Example` 關鍵字英文、內容中文，跟 `features/` 的慣例一致，
但**不受封閉步驟文法約束**（它不執行，句子以業務方讀得順為準）。

### D2a. `bdd-formulation` 讀旅程的用語，不讀它的內容

（需求方選 A）`bdd-formulation` 的輸入變成 `spec.md` ＋ `acceptance/*.feature`，但兩者角色不同：

- **`spec.md` 是唯一來源**：`Rule:` 與 `Example:` 只從 FR／AC 來，`@rule-`／`@example-` tag 照舊，
  `check_spec.py` 照舊比對 `spec.md` 的 AC 與 `features/` 的 tag
- **`acceptance/` 只參考用語**：挑封閉文法的步驟句型、`Feature:` 敘述、參數名時，優先用旅程裡
  業務方看過的詞，讓兩層講同一套話

MUST NOT: 從旅程新增 `spec.md` 沒有的例子、場景或 `Then`。旅程刻意少、句子不受文法約束——當
來源的話，漏掉的 AC 與自由發揮的句子會流進可執行的那一層。

旅程的用語與 `spec.md` 或 `docs/CONTEXT.md` 打架（同一件事兩個名字）→ 以 `docs/CONTEXT.md` 為準，
回報給使用者，不改旅程。

### D3. 再跑

`bdd-spec` 修補了 `spec.md` 之後再跑本 skill：既有旅程**修補，不重寫**——動到的旅程改、
`@covers-` 跟著改；答掉的 `# [need clarification]` 刪掉；新的 AC 視需要加一條旅程或併進
既有的。

### D4. `check_spec.py`

加一項檢查（第 17 項，**驗收旅程的引用**）：

- 讀 `specs/*/acceptance/*.feature`；目錄不存在就跳過並不報錯（舊 feature、或還沒走到這一步）
- 每個 `@covers-<n>.<m>` 都要在同目錄 `spec.md` 的 AC 裡存在；不存在 → ✗，訊息帶檔名與 tag
- 一個 `Example:` 上方沒有任何 `@covers-` → ✗（沒有錨點的旅程無法被稽核）

不檢查「每條 AC 都被涵蓋」（非目標），不檢查 `Then` 的文字對不對得上 AC（機器判不了，交給
D2 第 6 條與就緒問句）。

**`find_features` 要改**：它在找不到 `features/`、`test/features/`、`tests/features/` 時會退回
`rglob("*.feature")` 取最淺的目錄（`check_spec.py:63-69`），而寫旅程、判就緒的當下正好還沒有
`features/`——`specs/<dir>/acceptance/` 會被當成可執行 `.feature`，旅程被報成孤兒。退回搜尋
要排除 `specs/` 底下，並加一個測試：只有 `specs/x/acceptance/*.feature`、沒有 `features/` 時，
`find_features` 回 `None`。

### D5. 文件

- `docs/ai-sdlc.md`：SPEC 那一列加驗收旅程
- `PLAN.md`：skill 表加一列、產物佈局加 `acceptance/`、進度加一條
- `CLAUDE.md`：SPEC 的描述與產物清單

### D6. 範例

新增 `bdd-spec-by-example/examples/` 的旅程範例，`@covers-` 對得上
`bdd-spec/examples/minimal-spec.md` 的 AC。測試新增一個 fixture：同一個 `specs/<dir>/` 放
`minimal-spec.md`、`minimal-discovery.md` 與旅程範例（`acceptance/`），`check_spec.py` 全過。
既有兩份範例不用改。

## 交給實作者的

- `bdd-spec-by-example/SKILL.md` 的長度與 references 怎麼拆（至少 `references/journey-format.md`、
  一份範例 `examples/<旅程>.feature`，用停車場的例子，對得上 `bdd-spec/examples/minimal-spec.md`）
- 第 17 項的訊息措辭
- `@covers-` 的 regex 細節（`@covers-1.1` 一個 tag 一個 AC）

## 驗證

1. 新測試：引用不存在 AC 的旅程 → ✗；`Example:` 沒有 `@covers-` → ✗；沒有 `acceptance/` →
   不報錯；`acceptance/` 下的檔案不被當成可執行 `.feature`
2. 範例旅程對 `minimal-spec.md` 跑 `check_spec.py` 全過
3. `example-mapping` 的就緒問句提到 `acceptance/`；`bdd-spec` 第 5 步的呼叫順序是
   診斷 → `bdd-spec-by-example` → 就緒問句
4. `claude plugin validate .`、`python3 -m unittest discover -s tests`、改過的每個 skill 的
   `audit_skill.py` 全過
