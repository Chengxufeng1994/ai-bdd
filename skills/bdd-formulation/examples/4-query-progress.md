# 查詢課程進度

**查詢：第四種形狀接回傳的投影，不是實體狀態**

> 出處是一個線上課程平台。`@rule-N` / `@example-N.M` 是照本 skill 的規定補上的
> （原專案沒有用 tag，編號為示意）。其餘逐字保留。

> 其他範例：
> `1-video-progress` `2-submit-assignment` `3-deliver-course`
> `5-query-course-list` `6-query-product-list` `7-create-order`
> `8-query-order` `9-pay-order` `10-cancel-order`
> `11-query-user-roles` `anti-patterns`

---

前三份都是**命令**（改變世界）。這份是**查詢**——第四種形狀因此不同：
`Then` 接的是回傳的投影，不是實體狀態。

```gherkin
@ignore @query
Feature: 查詢課程進度

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
      | Bob   | 1    | 0      |
    And 用戶 "Alice" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
      | 角色類型     |
      | 旅程購買     |
      | 旅程訂閱狀態 |

  Rule: 前置（狀態）- 學生必須同時擁有「旅程購買」和「旅程訂閱狀態」角色才能查詢課程進度

    @rule-1 @example-1.1
    Example: 完全沒有旅程角色時查詢失敗
      Given 用戶 "Bob" 沒有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的任何角色
      When 用戶 "Bob" 查詢課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度
      Then 操作失敗

    @rule-1 @example-1.2
    Example: 只有「旅程訂閱狀態」角色時查詢失敗
      Given 用戶 "Bob" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
        | 角色類型     |
        | 旅程訂閱狀態 |
      When 用戶 "Bob" 查詢課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度
      Then 操作失敗

    @rule-1 @example-1.3
    Example: 只有「旅程購買」角色時查詢失敗
      Given 用戶 "Bob" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
        | 角色類型 |
        | 旅程購買 |
      When 用戶 "Bob" 查詢課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度
      Then 操作失敗

  Rule: 後置（回應）- 成功查詢應回傳完整的課程進度資訊

    @rule-2 @example-2.1
    Example: 查詢進行中的影片課程時回傳當前進度
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度為 70%，狀態為 "進行中"
      When 用戶 "Alice" 查詢課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度
      Then 操作成功
      And 課程進度的查詢結果應包含：
        | 用戶名稱 | 課程 ID                              | 課程名稱     | 課程類型 | 進度 | 狀態   |
        | Alice    | 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 | 物件導向基礎 | 影片課程 | 70   | 進行中 |

    @rule-2 @example-2.2
    Example: 查詢已交付的課程時狀態為已交付
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度為 100%，狀態為 "已交付"
      When 用戶 "Alice" 查詢課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度
      Then 操作成功
      And 課程進度的查詢結果應包含：
        | 用戶名稱 | 課程 ID                              | 課程名稱     | 課程類型 | 進度 | 狀態   |
        | Alice    | 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 | 物件導向基礎 | 影片課程 | 100  | 已交付 |

    @rule-2 @example-2.3
    Example: 查詢挑戰題課程時不回傳進度百分比
      Given 用戶 "Alice" 在課程 3c4d5e6f-7081-492a-9c3d-4e5f60718293 的狀態為 "已完成"
      When 用戶 "Alice" 查詢課程 3c4d5e6f-7081-492a-9c3d-4e5f60718293 的進度
      Then 操作成功
      And 課程進度的查詢結果應包含：
        | 用戶名稱 | 課程 ID                              | 課程名稱 | 課程類型   | 狀態   |
        | Alice    | 3c4d5e6f-7081-492a-9c3d-4e5f60718293 | 實作練習 | 挑戰題課程 | 已完成 |

    @rule-2 @example-2.4
    Example: 查詢從未開始的課程時回傳未開始
      When 用戶 "Alice" 查詢課程 2b3c4d5e-6f70-4819-8b2c-3d4e5f607182 的進度
      Then 操作成功
      And 課程進度的查詢結果應包含：
        | 用戶名稱 | 課程 ID                              | 課程名稱     | 課程類型 | 進度 | 狀態   |
        | Alice    | 2b3c4d5e-6f70-4819-8b2c-3d4e5f607182 | 設計模式問卷 | 問卷課程 | 0    | 未開始 |
```

### 查詢的四個不同之處

**1. 失敗只到 `操作失敗` 為止。**
命令的失敗要斷言「狀態沒變」，查詢不必——它本來就不改變任何東西。硬寫一句
「狀態應為原值」是在驗一件不可能發生的事，那種斷言永遠是綠的。

**2. 結果用資料表，不用一串 `And`。**
回傳的是一個投影（六個欄位湊成一列）。拆成六行斷言就看不出它們屬於同一筆資料，
而且六行各要一支 step definition。一張表一支。

**3. 斷言表的欄位隨案例變動。**
挑戰題那一列沒有「進度」欄——挑戰題沒有百分比這個概念。用固定欄位硬填一個
`N/A`，等於發明一個業務上不存在的值。

**4. 「未開始」那個例子完全沒有 `Given`。**
未開始就是 `Background` 播完之後的樣子。多寫一行
`Given 用戶 "Alice" 在課程「設計模式問卷」 的狀態為 "未開始"`，驗的就變成自己的佈置——
系統若其實把新用戶初始化成別的狀態，這個例子會因為被你覆蓋而測不出來。

### 三個權限例子為什麼不合併成 Outline

它們差在 Bob 擁有哪些角色，而角色是**資料表**——`Scenario Outline` 的
`Examples:` 每格只能放單一值，塞不進一張表。而且第一個例子（完全沒有角色）
連表都沒有，形狀本來就不同。

這是 Outline 的真實邊界：**參數是純量才用得上**。
