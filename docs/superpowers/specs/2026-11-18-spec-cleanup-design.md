# `spec.md` 只寫要做成什麼；技術決定回到 `plan.md`

**狀態**：已批　**日期**：2026-11-18

## 意圖

需求方的話，依時間順序：

- 「不應該跑 formulation」「整體流程應該會是 spec 選擇跑 spec-by-example 對焦，都沒問題才做 formulation
  才做 plan」
- 「這個 spec 還不用寫回去 docs」「就保持在這個 spec.md」
- 「所以我們的 spec 這個技能要好好修正」「我希望每一個步驟都清楚表達」，並給了 aixbdd 的 `specify` 與
  addyosmani 的 `spec-driven-development` 當參考
- 「discovery.md 不要再出現 spec.md」
- 「Tech Stack … Testing Decisions 感覺這應該放到 plan」「回 DISCOVERY 改成 Open Question」「新增
  Out of Scope」「我這邊也覺得不應該有技術部分」
- 「可以按照 aixbdd 重寫 format」並給了一張章節順序表；「可以不要把現在的優點拿掉嗎」

來源是一次試跑（`runs/parking-billing`，discovery → spec → plan）：`spec` 照它第 6 步的指示自己呼叫了
`formulation`，繞過了每個 skill 都由使用者叫用的決定（`f166d51`）；`spec` 第 4 步在 feature 還沒做完時
就寫進 `docs/CONTEXT.md`；`formulation` 不知道 seam，把 seam 上驗不到的設備端行為寫成了場景。

## 決定

### D1. 流程

```
discovery → spec（需求 → 診斷 → 旅程 → 就緒問句 → 已對焦，停）→ formulation → plan → task-breakdown → implement
```

每一步由使用者打 `/<skill>`；沒有 skill 呼叫下一個。`spec` 第 5 步裡的 `feature-mapping`、
`example-mapping`、`spec-by-example` 仍由 `spec` 帶著跑——它們是 `spec` 自己的診斷，不是下一步。

### D2. `spec.md` 不寫技術決定，`plan.md` 寫

推翻 2026-11-15 的 `spec-carries-tech`。Tech Stack 到 Testing Decisions 七節與 seam 確認搬到 `plan`；
`tech-format.md` 用 `git mv` 搬到 `skills/plan/references/`。`plan.md` 前半是技術決定、後半是執行計畫。

`plan` 在 `.feature` 之後跑，所以多一步：每個場景的 Then 在 seam 上看不看得到；看不到的進 `## 缺口`，
交回 `formulation` 換說法或交回 `spec` 換驗法。

下游（`task-breakdown`、`implement`、`behavior-driven-development`、`archive`）改從 `plan.md` 讀技術決定。

### D3. `spec.md` 照 aixbdd `specify` 的順序重排，現有規則全部保留

Document Overview → Context / Problem → Personas → Goals / Non-Goals → User Stories → Non-Functional
Requirements → Domain Concepts → Edge Cases → Assumptions / Constraints → Open Questions → Success
Metrics → Further Notes。

| 新的 | 做法 |
| --- | --- |
| story 的優先級 | 標題下一行（`check_spec.py` 精確比對標題）；來源跟 FR 一樣只准 `← Q-<n>` 或 `← PRD §n`；`discovery` 收尾時有兩則以上 story 就問 |
| Goals / Non-Goals | Non-Goals 一行一項，分延後與不做，寫為什麼 |
| Domain Concepts | `### 詞彙` ＋ `### 聚合與不變條件`，不含欄位型別；`check_spec.py` 第 18 項改查這一節 |
| Edge Cases | 只是指向 AC 的索引；找不到 AC 的是洞，進 `## Open Questions`，不在這裡寫規則 |
| Open Questions | 取代 `## 回 DISCOVERY 補問`；影響到的 story 就地標 `[NEEDS CLARIFICATION: …]` |
| NFR | 每條要有數字 |

保留的：來源標記三選一與 FR／AC 不准 `推論`、FR 底下巢狀掛 AC 與 `AC-<n>.<m>` 編號（`.feature` 的
`@example-n.m`、`check_spec.py` 的覆蓋比對、`@covers-` 都靠它）、EARS（`shall` 即 MUST）、Personas 五欄、
`#### Q` 只放指標、空節寫為什麼空、修訂歷史的觸發條件、NFR 與 Constraints 的分界。aixbdd 把 AC 與 FR
分開列；這裡不跟，因為那會拆掉整條追溯鏈。

### D4. 詞彙與聚合留在 `spec.md`，`archive` 才寫 `docs/`

`spec` 當初直接寫 `docs/CONTEXT.md` 的理由是「刪掉 spec 目錄後詞彙還在」；`archive` 改成整個目錄搬進
`specs/sleeping/` 之後不成立。現在只有 `archive` 寫 `docs/`：`CONTEXT.md` 只追加、不包區塊，已有同名
而定義不同的詞就停下。

### D5. `discovery.md` 不知道 `spec.md` 的格式

拿掉 `**歸屬**`：答案進哪一節由 `spec`／`plan` 依 `面向` 判斷。被退回時重開的題寫
`**來自**：補問 —「…」`，只引原文。

### D6. 試跑找到的其他問題

| 問題 | 改成 |
| --- | --- |
| `spec:44`「還有待答 → discovery」vs `:75`「喊停帶著待答往下」 | `:44` 加「而沒有喊停」 |
| `spec` 說補問非空也跑 `example-mapping`，`example-mapping` 說補問非空 `spec` 不叫它 | 照跑地圖與診斷，只是不提就緒問句 |
| `feature-mapping` 只有補問一個出口 | 多一個：答案在來源、只是 `spec` 漏寫 AC 的，寫進判讀交回 `spec`，不開 Open Questions |
| AC 能不能組合兩題的例子 | 可以，兩個來源都標；要新算一個數字的才進 Open Questions |

### D7. `spec` 的流程寫成動詞開頭的步驟

照 aixbdd：每一步一串 `READ`／`WRITE`／`CHECK`／`RUN`／`DELEGATE`／`ASK`／`STOP`，最後一行「做完的樣子」與
下一步。只先改 `spec`，其他 skill 不跟。

### 不在這次

- benchmark 污染：好幾個 skill 的範例就是 `parking-billing` 這個 case（試跑回報 #1–#4），另一個改動
- `system-analysis` 獨立成一步、`spec-by-example` 拆成獨立一步
- 名詞與狀態提早列舉：試跑回報說第 4 步沒找到任何洞，證據不支持

## 驗證

- 每個 skill 的 `audit_skill.py` 通過；`claude plugin validate .` 通過；`python3 -m unittest discover -s tests` 通過
- 還沒有實跑：要用新流程重跑一次 `parking-billing` 才知道 `plan` 的「seam 驗得到嗎」攔不攔得到上次那幾個場景
