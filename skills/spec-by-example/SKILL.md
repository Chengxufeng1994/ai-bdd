---
name: spec-by-example
description: >
  把 `待對焦` 的 `spec.md` 寫成業務語意的驗收旅程——`specs/ongoing/<date>-<feature>/acceptance/*.feature`，
  幾條可 demo 的長流程、只用業務語言、每條用 `@covers-<n>.<m>` 標出走過的 AC、未定的地方
  就地標 `# [need clarification]`——讓需求方對著旅程判就緒。不執行、不訪談、不改 FR／AC。
  BDD 六步流程 SPEC 步驟裡，診斷之後、就緒問句之前的一步。
  觸發詞：「寫驗收旅程」「spec by example」「業務版的 feature」「給業務方看的 feature」
  「把 spec 寫成可以 demo 的流程」「acceptance 旅程」
  「業務語意的 feature」「產生驗收流程」「業務方看得懂的驗收情境」。
  English: write acceptance journeys, specification by example, business-readable
  feature file, demo-able journeys from the spec.
disable-model-invocation: true
---

# 驗收旅程 — 把 `spec.md` 寫成業務方讀得懂的流程

| | |
| --- | --- |
| **職責** | 把 `待對焦` 的規格寫成幾條可 demo、只用業務語言的長流程，讓需求方對著它判就緒 |
| **讀** | `specs/ongoing/<date>-<feature>/spec.md`（`待對焦`、沒有待答的補問題）、`docs/CONTEXT.md`（有的話） |
| **寫** | `specs/ongoing/<date>-<feature>/acceptance/<旅程>.feature`；有標記時 `spec.md` `## Open Questions` 的新題 |
| **上游** | `spec` 第 5 步（`example-mapping` 跑完診斷之後） |
| **下游** | `example-mapping` 提就緒問句；有 `# [need clarification]` → `spec`〈有洞〉 |
| **不做** | 不執行、不訪談、不改 FR／AC／狀態、不填答案、不寫 `spec.md` 沒有的結果 |

**它是規格，不是測試。** 不執行、沒有 step definition、不進 `features/`。程式版本的
`features/*.feature` 照舊由 `formulation` 在 `已對焦` 之後寫（由 `plan` DELEGATE），來源照舊只有 `spec.md`。
每條旅程標出它走過哪幾條 AC，每個結果都是 `spec.md` 裡已經有的。

## 使用時機

- `spec` 寫完 `spec.md`、`example-mapping` 跑完診斷，要給需求方一份能對著判就緒的東西
- `spec.md` 修補過，既有旅程要跟著改
- 使用者要一份業務方讀得懂的 feature

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| `spec.md` 不存在或還有洞 | `spec`；還沒問完（結束方式不是 `收斂`／`喊停`）→ `discovery` |
| 判就緒、寫 `已對焦` | `example-mapping`，本 skill 寫完旅程交回給它 |
| 寫可執行的 `features/*.feature` | `formulation` |

## 參考檔案

- `rules/journey-format.md` — 六條規則、檔名、骨架。**寫之前先讀**
- `examples/借書看到期日.feature` — 對得上 `../spec/examples/minimal-spec.md` 的範例
- `../formulation/scripts/check_spec.py` — 第 17 項稽核 `@covers-`。在隔壁 skill：只跑，不改

## 前置確認（Ask user）

Check if: `spec.md` 狀態是 `待對焦`，`## Open Questions` 沒有待答的補問題（定義見 `../spec/references/spec-format.md`〈寫規格時發現的洞〉）。不是的話停下來，
說該回哪個 skill——有洞的規格寫成旅程，洞會被流暢的句子蓋住。

## 流程

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| 1 | 讀 `spec.md`，挑旅程 | — |
| 2 | 照 `rules/journey-format.md` 寫旅程 | — |
| 3 | 定不下來的，就地標記並同時在 `spec.md` 的 `## Open Questions` 開一題 | 不問——不訪談 |
| 4 | 跑 `check_spec.py`，第 17 項要過 | — |
| 5 | 交回 `spec` 或 `example-mapping` | — |

### 1. 讀 `spec.md`，挑旅程

讀 `## User Stories` 的 FR／AC、`## Personas`、`## Goals / Non-Goals`、`## Domain Concepts`，以及 `docs/CONTEXT.md`
的詞彙。挑的是**業務方會想看人走一遍的流程**：主流程一條，加少數結果明顯不同的變種。
一條旅程可以跨好幾則 story、好幾條 AC。

### 2. 寫旅程

照 `rules/journey-format.md`。每個 `Example:` 上方 `@covers-` 列出它走過的 AC。

MUST NOT: 寫出 `spec.md` 沒有的結果、規則或例子。旅程不執行，沒有任何東西會在之後
擋下一句「順手補上」的 `Then`——每個 `Then` 寫之前，先指出它是哪一條 AC。指不出來
就不寫，改走第 3 步。

### 3. 定不下來的，就地標記並寫進補問

IMPORTANT: 標記與 `spec.md` `## Open Questions` 的那一題要**同時**寫。只標在旅程裡，`spec` 與
`example-mapping` 看的 `spec.md` 仍然沒有待答的補問題，就緒問句照樣會提出來——洞安靜地留在
一份沒有人會再讀的檔案裡。

寫到一半發現 `spec.md` 推不出來的東西（某個結果沒有 AC、兩條 AC 在同一條流程裡打架、
流程缺一個前置狀態）：

- 旅程裡緊貼那一行標 `# [need clarification] <問題>`
- `spec.md` 的 `## Open Questions` 開一題補問題（`Q-<n>` 待答），照 `../spec/references/spec-format.md`
  〈寫規格時發現的洞〉的寫法，`**影響**` 那一行標明哪條旅程（`—— 見 acceptance/<檔名>`）

不要自己答、不要問需求方——本 skill 跟 SPEC 一樣不訪談。

### 4. 跑 `check_spec.py`

```bash
python3 <ai-bdd>/skills/formulation/scripts/check_spec.py <專案根>
```

第 17 項要過：每個 `@covers-` 都指得到 AC，每個 `Example:` 都有 `@covers-`。

### 5. 交回

| 旅程裡 | 交給 |
| --- | --- |
| 有任何 `# [need clarification]` | 每一處都已開成補問題；交回 `spec` 的〈有洞〉，由它把狀態改回 `澄清中` |
| 沒有 | `example-mapping` 提就緒問句，告訴需求方旅程在 `acceptance/` |

### 再跑：修補，不重寫

`spec` 修補了 `spec.md` 之後再跑：動到的旅程改、`@covers-` 跟著改；答掉的
`# [need clarification]` 刪掉；新的 AC 視需要加一條旅程或併進既有的。沒動到的旅程不碰——
需求方已經讀過它們，重寫等於要他再讀一次。

## 產物

`specs/ongoing/<date>-<feature>/acceptance/<旅程>.feature`，一個驗收主題一個檔，檔名用業務旅程
命名。骨架與六條規則見 `rules/journey-format.md`，完整範例見
`examples/借書看到期日.feature`。有標記時另有 `spec.md` `## Open Questions` 的新題。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 每個 `Example:` 上方都有 `@covers-`，且每一個都指得到 `spec.md` 的 AC
- [ ] 每個 `Then` 都指得出它是哪一條被 cover 的 AC；指不出來的沒有寫
- [ ] 旅程裡沒有 API、HTTP、欄位名、selector、資料表（`rules/journey-format.md` 第 1 條）
- [ ] 每個 `# [need clarification]` 在 `spec.md` 的 `## Open Questions` 都有對應的一題 `待答`，標明哪條旅程
- [ ] 跑了 `check_spec.py`，第 17 項過，輸出貼在回報裡
- [ ] 沒有改 `spec.md` 的 FR／AC，也沒有自己答或去問需求方
- [ ] 再跑時沒動到的旅程沒有被重寫

## 完成後

告訴使用者：寫了幾個檔、幾條旅程、各自 cover 哪些 AC、有沒有 `# [need clarification]`
（有的話列出來）、交回給誰。
