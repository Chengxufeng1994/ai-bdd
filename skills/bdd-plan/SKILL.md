---
name: bdd-plan
description: >
  SPEC 說要做成什麼，PLAN 說怎麼做——讀已對焦的 `spec.md`、它的 `.feature` 與
  `discovery.md` 裡技術面向的答案，寫 `plan.md`：技術設計（seam、模組、介面與 API 契約、
  schema、每個 AC 的測試層級、風險）。不改需求、不切票——拆成票交給 `bdd-tasks`。
  BDD 六步流程的 PLAN 步驟。
  觸發詞：「決定 seam」「驗收測試打在哪一層」「API 契約要怎麼定」
  「這個 feature 動哪些模組」「技術設計」「schema 怎麼定」「spec 寫完了下一步」「寫 plan」。
  English: write the plan, decide the acceptance test seam, technical design,
  which modules does this touch, define the API contract.
---

# PLAN — 怎麼做

一份已對焦的 `spec.md`（FR 與 story）與涵蓋的 `.feature` 進來，出去的是 `plan.md`：
**技術設計**——驗收測試打在哪一層 seam、動哪些模組、介面與契約、schema、每個 AC 在哪一層
測、風險。把這份設計拆成票是下一個 skill `bdd-tasks` 的事。

**不改需求。** 技術決定是這一步的工作；FR、AC、story 不是。需求有洞寫進 `## 缺口`，
回 discovery，不要順手補。

## 使用時機

- `spec.md` 已對焦、`.feature` 寫好了，要決定 seam 與技術設計
- `bdd-tasks` 拆票時在 `## 缺口` 記了設計的洞，要補技術節

## Skill Boundaries

- 還沒有 `spec.md`、或不是 `已對焦` → `bdd-spec`（`discovery.md` 也沒有就先 `bdd-discovery`）；
  還沒有 `.feature` → `bdd-formulation`
- **要改 FR、AC、story → 那是 `bdd-spec` 的工作**，不是這裡
- 要切成幾張票、排順序、看哪幾張能同時做 → `bdd-tasks`，本 skill 寫完 `plan.md` 交給它
- 整批塞得進一個 context window 而且技術決定一句話講得完 → 技術節寫短
- 要實際寫 step definition 與產品程式碼 → `bdd-implement`

## 參考檔案

- `references/plan-format.md` — `plan.md` 的骨架、技術節怎麼寫、範例
- `../bdd-discovery/references/technical-probes.md` — 五個技術面向問了什麼

---

## 不改需求

MUST NOT: 在這一步新增、刪除或改寫任何 FR、AC、story。它們已經在 `spec.md` 裡、
需求方說過「就緒」；這裡再決定一次，只會產生一份跟規格不一樣的規格。

判準：指著技術節裡的任何一個決定，說得出它是為了 `spec.md` 的哪條 FR、哪個 AC，
或來自 `discovery.md` 的哪一題嗎？說不出來就是這一步僭越了。

`spec.md` 沒寫到而做技術設計時發現需要的，寫進 `## 缺口`：

- **技術選擇**（哪個 library、欄位叫什麼）→ 留給 IMPLEMENT 決定
- **需求洞**（某條規則在某個情況下該怎麼做沒人說過）→ 把同目錄 `discovery.md` 改回
  `澄清中`、`### 版本修訂歷史` 記一列「回頭：<原因>」，請對方跑 `bdd-discovery`——它續跑時
  會讀這份 `## 缺口`，並在新題的 `**來自**` 逐字引用那一條

**再跑時清掉答過的需求洞。** `discovery.md` 裡某一題的 `**來自**` 逐字引了 `## 缺口` 的
某一條、而且那題已答，就把那一條從 `## 缺口` 刪掉，答案照常進技術節或交給 `bdd-spec`
（改了規則的話）。留著的話，下一次 discovery 續跑雖然不會重開它，讀的人卻分不出它還開不開。

---

## 流程

### 1. 讀

一次讀完整批：

| 讀什麼 | 為了什麼 |
| --- | --- |
| `spec.md` 的 `## User Stories` | 每條 FR 與 AC——seam 要看過具體的 AC 才推得出來，測試層級要逐個 AC 標 |
| 涵蓋的每個 `.feature` | 場景與 tag——每個 Example 的測試層級要對著它標 |
| `discovery.md` `## Open Questions` 技術面向、狀態 `已答` 的題 | **技術決定的原始素材**——小節的 `歸屬` 說它該進技術節哪一節，行尾 `← Q-<n>` |
| `spec.md` 的 `## Assumptions / Constraints`、`## Non-Functional Requirements` | 外部限制與品質屬性，`## Tech Preferences & Constraints` 與 `## Risks` 的來源 |
| `discovery.md` 的 `### 核心關係人` | **照抄**進 `plan.md` |

**一次讀完整批。** API 契約與 schema 天生跨 story，只有並排才看得見衝突。

順便看一眼專案既有的結構——找既有的 seam。

### 2. 決定 seam —— 驗收測試打在哪一層

`.feature` 決定驗什麼，seam 決定**打進系統的哪一層**：HTTP handler？
application service？domain？這個決定會長成 step definition 的形狀。

四條規則：

| 規則 | 為什麼 |
| --- | --- |
| 優先用既有 seam | 每多一個 seam 就多一套測試替身與資料建構 |
| 取**最高**的那一層 | 越高越接近使用者真正做的事，也越不綁實作細節 |
| 整個變更的理想數量是**一個** | 兩個 seam 代表這批行為的入口不只一個，多半是切分沒切乾淨 |
| 需要新 seam 就在**你能做到的最高點**提出 | 往下挪一層通常是為了好寫，代價是驗到的東西離使用者更遠 |

**Ask user: 把 seam 攤出來確認再寫。** DISCOVERY 的技術面向通常已經問過 seam
（`discovery.md` 裡 `面向` 為 `seam` 的題）——有的話這裡只是確認，行尾 `← Q-<n>`；沒有的話，
提出你的判斷、標 `← 推論`，並在 `## 缺口` 記一條「seam 沒問過」。

MUST NOT: 把 seam 的討論變成需求訪談。談的時候發現需求有洞，寫進 `## 缺口`，
**不要順手問掉**——這一步沒有把答案寫進 `discovery.md` 的權限，順口問到的答案只活在這一次
對話裡。

IMPORTANT: seam 一旦寫進 `plan.md` 就往下傳——IMPLEMENT 照著打，REVIEW 把「沒人同意過的
seam」當成 finding。這個綁定穿過 `plan.md` 這份文件；所以這一節寫不寫，決定了後面兩步
抓不抓得到。

### 3. 寫技術節

照 `references/plan-format.md`〈技術節怎麼寫〉填 `## Tech Preferences & Constraints`、
`## Known Boundaries`、`## Implementation Decisions`、`## Testing Decisions`、`## Risks`。
每個決定行尾標來源；`推論` 可以用，但每一條同時列進 `## 缺口`。

---

## 產物

一個 feature 一份 `specs/<date>-<feature>/plan.md`。骨架照抄 `references/plan-format.md`。

MUST: 只寫入 `specs/`，**不修改專案既有的任何檔案**。`discovery.md` 只在退回需求洞時
動兩處：狀態與修訂歷史一列。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] `spec.md` 的 FR、AC、story 沒有任何新增、刪除或改寫
- [ ] 每個技術決定行尾都有來源；每一條 `← 推論` 同時列在 `## 缺口`
- [ ] seam 寫進了 `plan.md`：是 `← Q-<n>` 確認過的，或標 `← 推論` 並在 `## 缺口` 記了「seam 沒問過」
- [ ] `plan.md` 裡沒有票——沒有 `## 票`、沒有編號的工作項
- [ ] 談 seam 時發現的需求洞進了 `## 缺口`，沒有在對話裡順口問掉
- [ ] 只寫入 `specs/`；`discovery.md` 只在退回需求洞時動了狀態（改回 `澄清中`）與修訂歷史一列
- [ ] 再跑時，`discovery.md` 已答且 `**來自**` 逐字引用的需求洞已從 `## 缺口` 刪掉

## 完成後

告訴對方：

1. **seam 是哪一層**、是確認 `discovery.md` 裡問過的還是這一步推的
2. 動哪些模組、新的契約有哪些、幾個決定標了 `推論`
3. `## 缺口` 有幾條，其中幾條是需求洞（已退回 discovery）

下一步 `bdd-tasks` 把這份設計拆成票。退回了需求洞 → 先 `bdd-discovery`。
