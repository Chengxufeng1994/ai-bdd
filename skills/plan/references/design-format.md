# `plan.md` 三份設計的骨架

照抄即可。讀這份的時機：`plan`〈設計〉寫 `## Domain Model`、`## Database Schema`、`## API 契約` 之前，
Step 2 寫 `## 依賴圖` 之前。整份 `plan.md` 的順序見 `plan-format.md`。

四節各一個骨架、一段圖書館借還書的範例（照 `../../spec/examples/minimal-spec.md`：FR-1 到期日、FR-2 續借、
FR-3 逾期擋續借）。來源標記照 `tech-format.md`〈來源標記：五選一〉；**需求依據**另外指 FR／AC 或 `Q-<n>`，那是
「這一項為了哪條規則存在」，跟「這個決定是誰定的」是兩件事。

空著的節保留標題，寫為什麼空（「（無）——這個 feature 不存任何資料」）。

---

## Domain Model

骨架同 `../../archive/references/domain-model-format.md`，標題降一級：`# Domain Model` 變成 `## Domain Model`
（拿掉 `## 聚合` 那一層，聚合用 `### <聚合名>`）、`## 跨聚合的規則` 變成 `###`。只寫這個 feature 新增或改動的
聚合，沒有修訂紀錄。`archive` 做完時照它併進
`docs/DOMAIN-MODEL.md`。

```markdown
## Domain Model

### <聚合名>

**是什麼**: 一句話
**邊界**: 哪些東西跟它一起變更、哪些不是
**來源**: FR-N
**原本是什麼、為什麼改**: <動到 `docs/DOMAIN-MODEL.md` 既有的聚合才寫；新的聚合寫「（新）」>

**不變條件**（任何時刻都必須為真）

| 條件 | 來源 | 被誰保證 |
| --- | --- | --- |

### 跨聚合的規則

不屬於任何單一聚合的。寫明為什麼放不進去。
```

- **判準**：這條規則**跨 feature 仍然成立**才是不變條件；只在某次操作要滿足的，留在 FR 與 `.feature` 的 `Rule:`
- `被誰保證` **只寫哪個聚合**（或哪條跨聚合的規則），不寫資料庫索引、鎖——那些寫在 `## Database Schema` 的
  `守住哪條不變條件`
- 動到 `docs/DOMAIN-MODEL.md` 既有的聚合，`原本是什麼、為什麼改` 一定要寫：`archive` 照它記修訂紀錄，下一個
  feature 的人靠那張表判斷某個決定還算不算數

範例：

```markdown
## Domain Model

### 借閱

**是什麼**: 一本書的一次借出，從館員借出到讀者還書
**邊界**: 到期日、續借過沒、歸還時間跟它一起變；讀者證的有效與否不是——那屬於縣裡的讀者證系統
**來源**: FR-1、FR-2
**原本是什麼、為什麼改**: （新）——第一個 feature，`docs/DOMAIN-MODEL.md` 還不存在

**不變條件**（任何時刻都必須為真）

| 條件 | 來源 | 被誰保證 |
| --- | --- | --- |
| 一筆借閱最多續借一次 | FR-2 | 借閱聚合 |
| 到期日只由借出日或續借算出，不被直接改寫 | FR-1、FR-2 | 借閱聚合 |

### 跨聚合的規則

- 名下有任何一筆逾期未還的借閱，就不能續借任何一筆 ← FR-3。放不進借閱聚合：要看同一位讀者的**其他**借閱，
  一筆借閱看不到別筆
```

---

## Database Schema

表的形狀同 `docs/DATABASE-SCHEMA.md`（`../../archive/references/docs-format.md`），多一欄
`守住哪條不變條件`：指回 `## Domain Model` 的一條不變條件，這個欄位或約束沒有守住任何一條就寫 `—`。
這一欄只寫在 `plan.md`；`archive` 併進 `docs/DATABASE-SCHEMA.md` 時拿掉它。

```markdown
## Database Schema

### <表名>

屬於 <聚合名> 聚合。

| 欄位 | 型別 | 約束 | 來源 | 守住哪條不變條件 |
| --- | --- | --- | --- | --- |

**索引**: <欄位>——<為什麼> ← <來源>
```

- **一張表屬於一個聚合**，寫在表名底下；屬於兩個聚合的表，代表聚合邊界或表的切法有一邊錯了
- `來源` 寫需要這個欄位的 FR／AC；`archive` 搬進 `docs/` 時前面加 `<date>-<feature>`
- 記憶體儲存、還沒有資料庫的專案，整節寫「（無）——<為什麼>」

範例：

```markdown
## Database Schema

### loans

屬於借閱聚合。

| 欄位 | 型別 | 約束 | 來源 | 守住哪條不變條件 |
| --- | --- | --- | --- | --- |
| id | uuid | PK | FR-1 | — |
| reader_card_no | text | NOT NULL | FR-1 | — |
| book_id | text | NOT NULL | FR-1 | — |
| lent_on | date | NOT NULL | FR-1 | — |
| due_on | date | NOT NULL | FR-1、FR-2 | — |
| renewed | boolean | NOT NULL DEFAULT false | FR-2 | 一筆借閱最多續借一次（布林只有一次可翻） |
| returned_at | timestamptz | NULL | FR-3 | — |

**索引**: `(reader_card_no, returned_at)`——續借時查這位讀者名下還沒還的借閱 ← FR-3
```

「到期日不被直接改寫」沒有出現在 `守住` 欄：資料庫守不住它，靠借閱聚合只開借出與續借兩個入口。

---

## API 契約

有 HTTP 對外介面時，寫一段 ```` ```yaml ```` 的 OpenAPI 片段，**只放這次新增與改動的** `paths` 與 `components`；
沒有 HTTP 時，寫程式碼層的介面契約：埠的方法、輸入、輸出、錯誤。兩種都由 `api-and-interface-design` 設計，
`plan` 寫進來。

**每個 operation 指得出它服務哪幾條 AC**——OpenAPI 寫在 `description`，程式碼層契約寫在方法後面。指不出 AC
的 operation 是猜的，不進契約。

`implement` 做一張票時，把它涵蓋的 operation 與用到的 `components` 照抄進專案的 `api/openapi.yaml`；`plan.md`
是來源，搬的時候不改形狀。

範例（有 HTTP）：

```yaml
paths:
  /loans:
    post:
      operationId: lendBook
      description: |
        館員借出一本書，回到期日。服務 AC-1.1、AC-1.2、AC-1.3
        不冪等、沒有 idempotency key：不可重試。
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LendRequest'
      responses:
        '201':
          description: 借出成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Loan'
  /loans/{loanId}/renewal:
    post:
      operationId: renewLoan
      description: |
        讀者在到期日前續借一次。服務 AC-2.1、AC-3.1
        不冪等、沒有 idempotency key：不可重試。
      parameters:
        - name: loanId
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: 續借成功，回新的到期日（AC-2.1）
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Loan'
        '409':
          description: 名下有逾期未還的書，拒絕續借，到期日不變（AC-3.1）
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
components:
  schemas:
    LendRequest:
      type: object
      required: [readerCardNo, bookId]
      properties:
        readerCardNo: { type: string }
        bookId: { type: string }
    Loan:
      type: object
      required: [id, dueOn, renewed]
      properties:
        id: { type: string, format: uuid }
        dueOn: { type: string, format: date }
        renewed: { type: boolean }
    Error:
      type: object
      required: [error]
      properties:
        error:
          type: object
          required: [code, message]
          properties:
            code: { type: string, example: OVERDUE_LOANS }
            message: { type: string }
```

範例（沒有 HTTP，同一個借書行為寫成驅動埠）：

```markdown
## API 契約

`借閱服務`（驅動埠）

- `借出(讀者證號, 書號) -> 借閱` —— 回到期日。服務 AC-1.1–1.3
- `續借(借閱編號) -> 借閱` —— 回新的到期日；名下有逾期時回 `有逾期未還` 錯誤，到期日不變。服務 AC-2.1、AC-3.1
```

---

## 三份並排對照

〈設計〉寫完 `## Domain Model`、`## Database Schema`、`## API 契約` 之後，並排對一次：

| 對什麼 | 對不上的樣子 |
| --- | --- |
| 每條不變條件有一個聚合守住，資料庫守得住的有一個約束 | 不變條件的 `被誰保證` 空著；能用約束守卻沒有對應欄位 |
| 每張表屬於一個聚合 | 一張表寫了兩個聚合，或不屬於任何聚合 |
| 每個 operation 最多寫一個聚合，除非一條跨聚合的規則說要 | 一個 operation 同時改兩個聚合，`### 跨聚合的規則` 卻沒有那一條 |

對不上通常表示某一份的判斷錯了，先找出是哪一份，不要不判斷就挑一份去遷就另一份。錯在設計 → 這一步就改掉那一份，
對不上的地方與怎麼改的記進 `## Further Notes` 的「技術決定的修訂」；錯在需求（規則沒說、或兩條 FR 打架）→ `## 缺口`
一列，交回 `spec`。

範例：不變條件「一筆借閱最多續借一次」有聚合、也有 `renewed` 守住，但 `## API 契約` 沒有「第二次續借」的錯誤回應——
FR-2 說限一次，卻沒有一條 AC 寫第二次續借會怎樣。這是需求的洞：`## 缺口` 一列，交回 `spec`。

---

## 依賴圖

節點是**模組與契約**，不是票：哪個先存在，別的才寫得出來。用文字樹，子節點依賴它的父節點：

```markdown
## 依賴圖

- 讀者證系統的查詢埠（`## Implementation Decisions`）
- 借閱聚合（算到期日、續借一次）
  - `loans` 表與它的 migration
    - 借閱 repository
  - `## API 契約` 的 `lendBook`、`renewLoan`
    - HTTP handler（用到讀者證查詢埠與借閱 repository）
```

`plan` 切票時把它當**約束**：一張票不得排在它依賴的節點之前（Step 5）；照圖切不出可 demo 的票，是圖錯了，
回 Step 2 改圖（Step 3）。波次與 checkpoint 排在 `tasks.md`。
