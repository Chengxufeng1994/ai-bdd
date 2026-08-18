# 提交挑戰題作業

**命令 · 型別檢查：同一動作施加在錯型別目標上被拒**

> 出處是一個線上課程平台。`@rule-N` / `@example-N.M` 是照本 skill 的規定補上的
> （原專案沒有用 tag，編號為示意）。其餘逐字保留。

> 其他範例：
> `1-video-progress` `3-deliver-course` `4-query-progress`
> `5-query-course-list` `6-query-product-list` `7-create-order`
> `8-query-order` `9-pay-order` `10-cancel-order`
> `11-query-user-roles` `anti-patterns`

---

```gherkin
@ready
Feature: 提交挑戰題作業

  Background:
    Given 系統中有以下旅程：
      | 旅程 ID                              | 名稱             |
      | a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 | 物件導向設計之旅 |
    And 系統中有以下課程：
      | 課程 ID                              | 名稱         | 類型       | 旅程 ID                              | 章節 ID                              | 獎勵.經驗值 |
      | 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 | 物件導向基礎 | 影片課程   | a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 | d4c3f5a7-ae61-4fb4-8c5d-30819eafbd24 | 300         |
      | 2b3c4d5e-6f70-4819-8b2c-3d4e5f607182 | 設計模式問卷 | 問卷課程   | a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 | d4c3f5a7-ae61-4fb4-8c5d-30819eafbd24 | 50          |
      | 3c4d5e6f-7081-492a-9c3d-4e5f60718293 | 實作練習     | 挑戰題課程 | a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 | d4c3f5a7-ae61-4fb4-8c5d-30819eafbd24 | 200         |
    And 系統中有以下用戶：
      | 名稱  | 等級 | 經驗值 |
      | Alice | 1    | 0      |
    And 用戶 "Alice" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
      | 角色類型     |
      | 旅程購買     |
      | 旅程訂閱狀態 |

  Rule: 前置（狀態）- 用戶必須同時擁有「旅程購買」和「旅程訂閱狀態」角色才能提交作業

    @rule-1 @example-1.1
    Example: 未擁有旅程購買角色無法提交作業
      Given 系統中有以下用戶：
        | 名稱 | 等級 | 經驗值 |
        | Bob  | 1    | 0      |
      And 用戶 "Bob" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
        | 角色類型     |
        | 旅程訂閱狀態 |
      When 用戶 "Bob" 提交課程 3c4d5e6f-7081-492a-9c3d-4e5f60718293 的挑戰題作業
      Then 操作失敗
      And 錯誤訊息應為 "權限不足"

    @rule-1 @example-1.2
    Example: 未擁有旅程訂閱狀態角色無法提交作業
      Given 系統中有以下用戶：
        | 名稱 | 等級 | 經驗值 |
        | Bob  | 1    | 0      |
      And 用戶 "Bob" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
        | 角色類型 |
        | 旅程購買 |
      When 用戶 "Bob" 提交課程 3c4d5e6f-7081-492a-9c3d-4e5f60718293 的挑戰題作業
      Then 操作失敗
      And 錯誤訊息應為 "權限不足"

  Rule: 提交挑戰題作業直接完成

    @rule-2 @example-2.1
    Example: 提交挑戰題作業成功
      Given 用戶 "Alice" 在課程 3c4d5e6f-7081-492a-9c3d-4e5f60718293 的狀態為 "進行中"
      When 用戶 "Alice" 提交課程 3c4d5e6f-7081-492a-9c3d-4e5f60718293 的挑戰題作業
      Then 操作成功
      And 用戶 "Alice" 在課程 3c4d5e6f-7081-492a-9c3d-4e5f60718293 的狀態應為 "已完成"

    @rule-2 @example-2.2
    Example: 無法提交非挑戰題課程的作業
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態為 "進行中"
      When 用戶 "Alice" 提交課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的挑戰題作業
      Then 操作失敗
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的狀態應為 "進行中"
```

### 這一份示範的重點

**`Background` 播了三種課程，但這則 story 只動到課程「物件導向基礎」 和 3。**
課程「設計模式問卷」 沒被任何例子用到——留著是對的：它讓「課程有不同型別」這件事在世界裡
成立，而 `無法提交非挑戰題課程的作業` 需要一個非挑戰題的課程存在才驗得了。

**「用錯型別」的失敗，狀態斷言指的是目標沒被改動。**
`用戶 "Alice" 在課程「物件導向基礎」 的狀態應為 "進行中"`——證明系統不只回了錯誤，
而且**沒有順手把課程「物件導向基礎」 標成已完成**。這正是只斷言錯誤訊息會漏掉的 bug。

**權限規則的兩個例子與範例一逐字相同，只有動作那一行不同。**
兩份 `.feature` 因此共用權限相關的 step definition。這不是巧合——
封閉文法讓跨 feature 的重用自然發生，而散文式的寫法會讓同一件事在
兩個檔案裡長成兩個樣子。
