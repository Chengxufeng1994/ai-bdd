# SPEC 加一張 Feature Map

**狀態**：已批　**日期**：2026-11-18

## 意圖

需求方的話：「是不是應該再寫一個 feature-mapping 的 skills，在 spec 的時候可以使用
feature-mapping + example-mapping？」「在 spec 階段就可以先產出 feature-mapping +
example-mapping → spec-by-example 去對焦，感覺上更完整」。

現在 SPEC 的診斷只有 `example-mapping`：一則 story 一張四色卡，看規則有沒有例子撐住。它是
深度的，看不到**整個 feature 的流程**——少了哪一步、某一步的替代流程沒寫、某條規則只有正例
沒有反例。這些洞要到 `spec-by-example` 寫旅程、甚至 IMPLEMENT 才撞到。

出處是 John Ferguson Smart 的 Feature Mapping，轉述見 `docs/feature-mapping.md`：例子 → 步驟
→ 結果排成一列，逐步問「這裡還可能發生什麼」長出新的列，再依規則分組、補反例。原文說它是
廣度優先。

## 決定

### D1. 新增 `feature-mapping`，在 `example-mapping` 之前

`spec` 第 5 步改成：

```
1. feature-mapping  攤整個 feature 的流程，缺口寫進 ## 回 DISCOVERY 補問
2. example-mapping  每則 story 攤四色卡、跑四條診斷
   ── 補問清單非空 → 第 7 步退回 discovery，不寫旅程
3. spec-by-example  寫 acceptance/*.feature 旅程
4. example-mapping  提就緒問句（寫 已對焦 的仍然只有它）
```

先廣後深：Feature Map 找到的洞（缺流程、缺規則）會改動 story 與 FR，先找出來，四色卡才不會
攤一張馬上要重畫的。

### D2. 不新增產物，`spec.md` 不加章節

起初提議把兩張圖都寫進 `spec.md`。需求方指出：「但我已經有 User Story 了」。對照之後兩張圖的
內容都已經住在別處：

| 圖上的東西 | 已經住在哪 |
| --- | --- |
| Example Map 的 story／rule／example／question | `## User Stories` 的 `US`／`FR`／`AC`／`Q` |
| Feature Map 的規則 | `FR-<n>` |
| Feature Map 的例子、步驟、結果 | `AC-<n>.<m>` 的 Given／When／Then |
| 一列列串成整個 feature 的流程 | `acceptance/*.feature` 的旅程 |

再寫一節就是同一件事兩份，遲早不一樣。`feature-mapping` 新增的是**一種找洞的方法**：地圖只印
在終端，找到的洞寫進 `## 回 DISCOVERY 補問`，跟 `example-mapping` 一樣。`check_spec.py` 不改。

### D3. 「這一步還可能發生什麼」不當場問

原文在工作坊現場問。SPEC 不訪談：當場問到的答案寫不進 `discovery.md`，FR 也標不出
`← Q-<n>`。所以每個問題都寫成補問一列，由 `discovery` 開成 `Q-<n>` 去問。

### D4. 單獨模式

跟 `example-mapping` 一樣：對一段描述跑時只印地圖與缺口，不寫任何檔案。

## 改動

- 新增 `skills/feature-mapping/SKILL.md`
- `skills/spec/SKILL.md`：下游、Skill Boundaries、第 5 步
- `skills/example-mapping/SKILL.md`：上游
- `skills/spec/references/spec-format.md`：`## 回 DISCOVERY 補問` 的寫入者
- `README.md`、`PLAN.md`、`CLAUDE.md`、`docs/ai-sdlc.md`

## 驗證

- `python3 skills/skill-rules/scripts/audit_skill.py skills/feature-mapping` 通過
- `claude plugin validate .` 通過
- `python3 -m unittest discover -s tests` 通過（沒動腳本，應不變）
- 還沒有實跑：拿 `runs/go-parking-billing/` 跑一次 Feature Map，看它有沒有抓到 `example-mapping`
  抓不到的洞
