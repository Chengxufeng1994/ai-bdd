# 驗收旅程的格式

`specs/ongoing/<date>-<feature>/acceptance/<旅程>.feature`。它是規格，不是測試：不執行、沒有
step definition、不進 `features/`。寫給業務方讀，讓他對著能 demo 的流程判就緒。

## 六條規則

1. **只寫業務語言。** 不得出現 API、HTTP、欄位名、selector、資料表。業務方唸出來
   會卡住的詞，就是不該在這裡的詞
2. **一個 `Example:` 是一條可 demo 的長流程（Journey）**，不是一個 AC 一個場景。
   數量刻意少：主流程一條，再加少數結果明顯不同的變種。一條旅程可以有好幾組
   `When`／`Then`
3. **`Feature:` 一個驗收主題、`Rule:` 原子、`Example:` 標題描述情境**，不重複規則。
   「訪客多停了一分鐘」是情境；「超過 30 分鐘要收費」是規則，已經寫在 `Rule:` 上
4. **`Background:` 只在多條旅程真的共用 setup 時用。** 只有一條旅程用到的前置寫在
   那條旅程裡
5. **未定的地方就地標 `# [need clarification] <問題>`**，緊貼引發疑問的那一行。同一個
   問題同時寫進 `spec.md` 的 `## Open Questions`，標明是哪條旅程。不訪談——
   標記是交給 `discovery` 問的題，不是你要回答的題
6. **每個 `Example:` 上方標 `@covers-<n>.<m>`**，一個 tag 一個 AC，列出它走過的 AC。
   旅程裡每一個 `Then` 都要對得到某個被 cover 的 AC；對不到的就是 `spec.md` 沒有的
   結果——**不寫**，改寫進 `spec.md` 的 `## Open Questions`

`check_spec.py` 第 17 項擋兩件事：`@covers-` 指到 `spec.md` 沒有的 AC、`Example:` 上方
沒有任何 `@covers-`。`Then` 的文字對不對得上 AC，機器判不了——那是第 6 條與就緒問句的事。

不要求每條 AC 都被某條旅程涵蓋。旅程選的是業務方想看的流程，不是覆蓋率。

## 檔名與語言

- 檔名用業務旅程命名（`訪客出場結帳.feature`），不按 story slug——一條旅程常跨兩則 story
- `Feature`／`Rule`／`Example`／`Given`／`When`／`Then` 關鍵字英文、內容中文，跟
  `features/` 的慣例一致
- **不受封閉步驟文法約束。** 它不執行，句子以業務方讀得順為準

## 骨架

```gherkin
Feature: <驗收主題>
  <一兩句：誰、要完成什麼>

  Rule: <一條 FR 的業務說法>

    @covers-<n>.<m> @covers-<n>.<m>
    Example: <情境>
      Given <業務上的前置狀態>
      When <角色做了什麼>
      Then <業務上看得到的結果>
      # [need clarification] <問題>   ← 只在真的定不下來時
```

範例：`../examples/訪客出場結帳.feature`，對得上 `../../spec/examples/minimal-spec.md`。
