# 增加影片進度

**命令 · 狀態機：合法範圍、單向移動、閾值觸發轉換**

> 出處是一個線上課程平台。`@rule-N` / `@example-N.M` 是照本 skill 的規定補上的
> （原專案沒有用 tag，編號為示意）。其餘逐字保留。

> 其他範例：
> `2-submit-assignment` `3-deliver-course` `4-query-progress`
> `5-query-course-list` `6-query-product-list` `7-create-order`
> `8-query-order` `9-pay-order` `10-cancel-order`
> `11-query-user-roles` `anti-patterns`

---

```gherkin
@ready
Feature: 增加影片進度

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

  Rule: 前置（狀態）- 用戶必須同時擁有「旅程購買」和「旅程訂閱狀態」角色才能更新課程進度

    @rule-1 @example-1.1
    Example: 未擁有旅程購買角色無法更新進度
      Given 系統中有以下用戶：
        | 名稱 | 等級 | 經驗值 |
        | Bob  | 1    | 0      |
      And 用戶 "Bob" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
        | 角色類型     |
        | 旅程訂閱狀態 |
      When 用戶 "Bob" 更新課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的影片進度為 80%
      Then 操作失敗
      And 錯誤訊息應為 "權限不足"

    @rule-1 @example-1.2
    Example: 未擁有旅程訂閱狀態角色無法更新進度
      Given 系統中有以下用戶：
        | 名稱 | 等級 | 經驗值 |
        | Bob  | 1    | 0      |
      And 用戶 "Bob" 擁有旅程 a1f0c2d4-7b3e-4c81-9f2a-0d5e6b7c8a91 的以下角色：
        | 角色類型 |
        | 旅程購買 |
      When 用戶 "Bob" 更新課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的影片進度為 80%
      Then 操作失敗
      And 錯誤訊息應為 "權限不足"

  Rule: 影片進度必須單調遞增

    @rule-2 @example-2.1
    Example: 成功增加影片進度
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度為 70%，狀態為 "進行中"
      When 用戶 "Alice" 更新課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的影片進度為 80%
      Then 操作成功
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度應為 80%，狀態應為 "進行中"

    @rule-2 @example-2.2
    Example: 進度不可倒退
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度為 70%，狀態為 "進行中"
      When 用戶 "Alice" 更新課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的影片進度為 60%
      Then 操作失敗
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度應為 70%，狀態應為 "進行中"

    @rule-2 @example-2.3
    Example: 相同進度值的更新應被接受但不改變狀態
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度為 70%，狀態為 "進行中"
      When 用戶 "Alice" 更新課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的影片進度為 70%
      Then 操作成功
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度應為 70%，狀態應為 "進行中"

  Rule: 進度值必須在 0-100% 之間

    @rule-3 @example-3.1 @example-3.2 @example-3.3
    Scenario Outline: 有效範圍內的進度值可以更新
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度為 50%，狀態為 "進行中"
      When 用戶 "Alice" 更新課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的影片進度為 <新進度>%
      Then 操作成功
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度應為 <新進度>%，狀態應為 "<預期狀態>"

      Examples:
        | 新進度 | 預期狀態 |
        | 60     | 進行中   |
        | 100    | 已完成   |
        | 51     | 進行中   |

    @rule-3 @example-3.4 @example-3.5
    Scenario Outline: 超出範圍的進度值無法更新
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度為 50%，狀態為 "進行中"
      When 用戶 "Alice" 更新課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的影片進度為 <新進度>%
      Then 操作失敗
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度應為 50%，狀態應為 "進行中"

      Examples:
        | 新進度 |
        | 101    |
        | -10    |

  Rule: 觀看至 100% 則課程完成

    @rule-4 @example-4.1
    Example: 影片進度達到 100% 時課程自動完成
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度為 90%，狀態為 "進行中"
      When 用戶 "Alice" 更新課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的影片進度為 100%
      Then 操作成功
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度應為 100%，狀態應為 "已完成"

  Rule: 已交付的課程仍可更新進度，但狀態不變

    @rule-5 @example-5.1
    Example: 已交付的課程更新進度不改變狀態
      Given 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度為 100%，狀態為 "已交付"
      When 用戶 "Alice" 更新課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的影片進度為 100%
      Then 操作成功
      And 用戶 "Alice" 在課程 1a2b3c4d-5e6f-4708-9a1b-2c3d4e5f6071 的進度應為 100%，狀態應為 "已交付"
```

### 這份檔案的樣板清單（實測）

7 個場景、33 行步驟，**10 個不重複樣板**，其中只出現一次的有 3 個：

| 次數 | 樣板                                            |
| ---  | ---                                             |
| 6    | `用戶 "X" 在課程 N 的進度為 N%，狀態為 "X"`     |
| 7    | `用戶 "X" 更新課程 N 的影片進度為 N%`           |
| 6    | `用戶 "X" 在課程 N 的進度應為 N%，狀態應為 "X"` |
| 4    | `操作成功`                                      |
| 3    | `操作失敗`                                      |
| 2    | `系統中有以下用戶：`                            |
| 2    | `用戶 "X" 擁有旅程 N 的以下角色：`              |
| 1    | `系統中有以下旅程：`                            |
| 1    | `系統中有以下課程：`                            |
| 1    | `錯誤訊息應為 "X"`                              |

**再加二十個場景，這張表幾乎不會變長。** 這就是封閉文法要買的東西——
`用 check_spec.py` 量得到，不必憑感覺。

### 值得注意的四個地方

**1. 權限規則排在第一條，而且用「前置（狀態）」開頭。**
它不是業務規則，是所有其他規則的前提。放第一條，讀者先知道「什麼情況下這些
規則根本輪不到」。

**2. 兩個權限例子幾乎一模一樣，但沒有合併成 Outline。**
它們差在缺少哪一個角色——那是兩個不同的原因，只是碰巧結果相同。
合併之後，某天其中一個的錯誤訊息改了，就得把 Outline 再拆開。

**3. Bob 在 Example 層播種，不在 `Background`。**
只有權限那兩個例子需要他。放進 `Background`，其餘五個例子都要背著一個
跟自己無關的用戶。

**4. 「相同進度值」那個例子的成敗判斷。**
`70% → 70%` 判定為**操作成功**（而不是失敗），這是設計決定，不是自然結果。
封閉文法逼你在 `Then` 的第一句就表態——散文式的寫法很容易把這一題含糊過去。
