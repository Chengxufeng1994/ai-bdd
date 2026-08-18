# 交付課程

**命令 · 副作用與序列：發經驗值、升級、做兩次、做多個**

> 出處是一個線上課程平台。`@rule-N` / `@example-N.M` 是照本 skill 的規定補上的
> （原專案沒有用 tag，編號為示意）。其餘逐字保留。

> 其他範例：
> `1-video-progress` `2-submit-assignment` `4-query-progress`
> `5-query-course-list` `6-query-product-list` `7-create-order`
> `8-query-order` `9-pay-order` `10-cancel-order`
> `11-query-user-roles` `anti-patterns`

---

原始檔案掛的是 `@ignore`（那個專案用來排除執行的 tag，作用類似本流程的 `@draft`）。
`Then` 有兩處沒有先講成敗——那兩行進了
[anti-patterns.md](anti-patterns.md)，其餘照抄。

```gherkin
@ignore
Feature: 交付課程

  Background:
    Given 系統中有以下旅程：
      | 旅程 ID                              | 名稱             |
      | a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 | 物件導向設計之旅 |
    And 系統中有以下課程：
      | 課程 ID                              | 名稱         | 類型       | 旅程 ID                              | 章節 ID                              | 獎勵.經驗值 |
      | 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 | 物件導向基礎 | 影片課程   | a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 | d4c3f5a7-ae61-4fb4-8c5d-30819eafbd24 | 300         |
      | 2b3c4d5e-6f70-4819-8b2c-3d4e5f607182 | 設計模式問卷 | 問卷課程   | a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 | d4c3f5a7-ae61-4fb4-8c5d-30819eafbd24 | 300         |
      | 3c4d5e6f-7081-492a-9c3d-4e5f60718293 | 實作練習     | 挑戰題課程 | a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 | d4c3f5a7-ae61-4fb4-8c5d-30819eafbd24 | 200         |
    And 系統中有以下用戶：
      | 名稱  | 等級 | 經驗值 |
      | Alice | 1    | 0      |
    And 用戶 "Alice" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
      | 角色類型     |
      | 旅程購買     |
      | 旅程訂閱狀態 |

  Rule: 前置（狀態）- 用戶必須同時擁有「旅程購買」和「旅程訂閱狀態」角色才能交付課程

    @rule-1 @example-1.1
    Example: 未擁有旅程購買角色無法交付課程
      Given 系統中有以下用戶：
        | 名稱 | 等級 | 經驗值 |
        | Bob  | 1    | 0      |
      And 用戶 "Bob" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
        | 角色類型     |
        | 旅程訂閱狀態 |
      And 用戶 "Bob" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態為 "已完成"
      When 用戶 "Bob" 交付課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071
      Then 操作失敗
      And 錯誤訊息應為 "權限不足"

    @rule-1 @example-1.2
    Example: 未擁有旅程訂閱狀態角色無法交付課程
      Given 系統中有以下用戶：
        | 名稱 | 等級 | 經驗值 |
        | Bob  | 1    | 0      |
      And 用戶 "Bob" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
        | 角色類型 |
        | 旅程購買 |
      And 用戶 "Bob" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態為 "已完成"
      When 用戶 "Bob" 交付課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071
      Then 操作失敗
      And 錯誤訊息應為 "權限不足"

  Rule: 課程完成後才能交付，交付課程後獲得經驗值和升級

    @rule-2 @example-2.1
    Example: 成功交付已完成的影片課程
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態為 "已完成"
      When 用戶 "Alice" 交付課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071
      Then 操作成功
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態應為 "已交付"
      And 用戶 "Alice" 的經驗值應為 300，等級應為 3

    @rule-2 @example-2.2
    Example: 無法交付未完成的課程
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態為 "進行中"
      When 用戶 "Alice" 交付課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071
      Then 操作失敗
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態應為 "進行中"
      And 用戶 "Alice" 的經驗值應為 0

    @rule-2 @example-2.3
    Example: 課程已交付時再次交付失敗
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態為 "已交付"
      When 用戶 "Alice" 交付課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071
      Then 操作失敗
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態應為 "已交付"
      And 用戶 "Alice" 的經驗值應為 300

    @rule-2 @example-2.4
    Example: 交付多個課程累積經驗值並多次升級
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態為 "已完成"
      And 用戶 "Alice" 在課程 2b3c4d5e-6f70-4819-8b2c-3d4e5f607182 的狀態為 "已完成"
      When 用戶 "Alice" 交付課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071
      And 用戶 "Alice" 交付課程 2b3c4d5e-6f70-4819-8b2c-3d4e5f607182
      Then 操作成功
      And 用戶 "Alice" 的經驗值應為 600，等級應為 5
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態應為 "已交付"
      And 用戶 "Alice" 在課程 2b3c4d5e-6f70-4819-8b2c-3d4e5f607182 的狀態應為 "已交付"
```

### 這一份最重要的一行

```gherkin
And 用戶 "Bob" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態為 "已完成"      ← 權限測試裡的這一行
```

Bob 的測試在驗權限，為什麼要把課程設成「已完成」？

因為**交付有兩個前提：有權限、且課程已完成**。不設這一行，Bob 的交付一樣會失敗
——但是因為課程沒完成，不是因為權限不足。**把權限檢查整段程式刪掉，這個測試
照樣是綠的。**

範例一的權限測試沒有這一行，因為那裡的權限檢查發生在進度規則之前，沒有競爭
原因。這不是不一致，是**各自排除了自己場景裡真正存在的競爭原因**。

### 另外兩個地方

**失敗案例斷言「副作用沒發生」。** `無法交付未完成的課程` 除了斷言狀態沒變，
還斷言 `經驗值應為 0`——擋住「回傳錯誤但已經把經驗值加下去」這種實作。
只驗狀態抓不到它。

**只有 `交付多個課程累積經驗值` 該用兩個 `When`。** 600 點經驗值與等級 5
**只有跑完兩次才存在**，沒有任何單一狀態表達得出來。

`無法重複交付` 原本也寫成兩個 `When`，**已改成狀態守衛**：`Given 狀態為 "已交付"`
＋ 一個 `When`。那條規則講的是「什麼狀態下不能交付」，不是「連續做會怎樣」——
用動作把系統推到那個狀態，會讓測試在交付功能本身壞掉時因為錯的理由變紅。
完整討論見 [`10-cancel-order.md`](10-cancel-order.md)。
