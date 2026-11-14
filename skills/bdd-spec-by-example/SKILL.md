---
name: bdd-spec-by-example
description: >
  把 `待對焦` 的 `spec.md` 寫成業務語意的驗收旅程——`specs/<date>-<feature>/acceptance/*.feature`，
  幾條可 demo 的長流程、只用業務語言、每條用 `@covers-<n>.<m>` 標出走過的 AC、未定的地方
  就地標 `# [need clarification]`——讓需求方對著旅程判就緒。不執行、不訪談、不改 FR／AC。
  BDD 六步流程 SPEC 步驟裡，診斷之後、就緒問句之前的一步。
  觸發詞：「寫驗收旅程」「spec by example」「業務版的 feature」「給業務方看的 feature」
  「把 spec 寫成可以 demo 的流程」「acceptance 旅程」。
  English: write acceptance journeys, specification by example, business-readable
  feature file, demo-able journeys from the spec.
---

# 驗收旅程 — 把 `spec.md` 寫成業務方讀得懂的流程

一份 `待對焦`、補問清單是「（無）」的 `spec.md` 進來，出去的是同目錄
`acceptance/<旅程>.feature`：**幾條可 demo 的長流程，只用業務語言，每條旅程標出它走過哪幾條
AC，每個結果都是 `spec.md` 裡已經有的。**

它是規格，不是測試。不執行、沒有 step definition、不進 `features/`。程式版本的
`features/*.feature` 照舊由 `bdd-formulation` 在 `已對焦` 之後寫，來源照舊只有 `spec.md`。

## 使用時機

- `bdd-spec` 寫完 `spec.md`、`example-mapping` 跑完診斷，要給需求方一份能對著判就緒的東西
- `spec.md` 修補過，既有旅程要跟著改
- 使用者要一份業務方讀得懂的 feature

## Skill Boundaries

- `spec.md` 不存在或還有洞 → `bdd-spec`；`discovery.md` 沒收斂 → `bdd-discovery`
- 判就緒、寫 `已對焦` → `example-mapping`，本 skill 寫完旅程交回給它
- 寫可執行的 `features/*.feature` → `bdd-formulation`
- 產出只有 `acceptance/*.feature`，以及（有標記時）`spec.md` `## 回 DISCOVERY 補問` 的新列。
  **不改** FR、AC、狀態，不寫 `discovery.md`

## 參考檔案

- `references/journey-format.md` — 六條規則、檔名、骨架。**寫之前先讀**
- `examples/訪客出場結帳.feature` — 對得上 `../bdd-spec/examples/minimal-spec.md` 的範例
- `../bdd-formulation/scripts/check_spec.py` — 第 17 項稽核 `@covers-`。在隔壁 skill：只跑，不改

## 前置確認（Ask user）

Check if: `spec.md` 狀態是 `待對焦`，`## 回 DISCOVERY 補問` 是「（無）」。不是的話停下來，
說該回哪個 skill——有洞的規格寫成旅程，洞會被流暢的句子蓋住。

## 流程

### 1. 讀 `spec.md`，挑旅程

讀 `## User Stories` 的 FR／AC、`## Personas`、`## Scope — In / Out`，以及 `docs/CONTEXT.md`
的詞彙。挑的是**業務方會想看人走一遍的流程**：主流程一條，加少數結果明顯不同的變種。
一條旅程可以跨好幾則 story、好幾條 AC。

### 2. 寫旅程

照 `references/journey-format.md`。每個 `Example:` 上方 `@covers-` 列出它走過的 AC。

MUST NOT: 寫出 `spec.md` 沒有的結果、規則或例子。旅程不執行，沒有任何東西會在之後
擋下一句「順手補上」的 `Then`——每個 `Then` 寫之前，先指出它是哪一條 AC。指不出來
就不寫，改走第 3 步。

### 3. 定不下來的，就地標記並寫進補問

寫到一半發現 `spec.md` 推不出來的東西（某個結果沒有 AC、兩條 AC 在同一條流程裡打架、
流程缺一個前置狀態）：

- 旅程裡緊貼那一行標 `# [need clarification] <問題>`
- `spec.md` 的 `## 回 DISCOVERY 補問` 加一列，照 `../bdd-spec/references/spec-format.md`
  的寫法，句尾標明哪條旅程（`—— 見 acceptance/<檔名>`）

不要自己答、不要問需求方——本 skill 跟 SPEC 一樣不訪談。

### 4. 跑 `check_spec.py`

```bash
python3 <ai-bdd>/skills/bdd-formulation/scripts/check_spec.py <專案根>
```

第 17 項要過：每個 `@covers-` 都指得到 AC，每個 `Example:` 都有 `@covers-`。

### 5. 交回

- 有任何 `# [need clarification]` → 交回 `bdd-spec` 的〈退回 discovery〉，由它把
  `discovery.md` 改回 `澄清中`
- 沒有 → 交回 `example-mapping` 提就緒問句，告訴需求方旅程在 `acceptance/`

## 再跑：修補，不重寫

`bdd-spec` 修補了 `spec.md` 之後再跑：動到的旅程改、`@covers-` 跟著改；答掉的
`# [need clarification]` 刪掉；新的 AC 視需要加一條旅程或併進既有的。沒動到的旅程不碰——
需求方已經讀過它們，重寫等於要他再讀一次。
