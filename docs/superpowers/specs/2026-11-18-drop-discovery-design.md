# `discovery.md` 併進 `spec.md`；`clarify` 併進 `discovery`

**狀態**：已批　**日期**：2026-11-18

## 意圖

需求方的話：「discovery 比較像是記錄問題，最後整理成 spec.md」「spec.md 就不再管 discovery.md 就可以利用
open question」；給了 aixbdd 的 `clarify` 與 addyosmani 的 `interview-me` 當參考：「或許根本不需要
discovery skill？」；最後定案：「discovery 應該較類似於 (interview-me) 使用 clarify 流程這樣，所以我覺得是
不需要 clarify skill」；並指明 interview-me 的定義：「One-question-at-a-time interview that extracts what the
user actually wants instead of what they think they should want, until ~95% confidence」。

兩個參考都沒有一份獨立的問答紀錄檔：aixbdd 的 `specify` 寫 `spec.md`、遇到高影響缺口委派 `clarify` 去問；
`interview-me` 的產物是一段需求方明說 yes 的意圖陳述。問答是寫 spec 的過程，不是另一份要維護的文件。

DISCOVERY 與 SPEC 分開的原始理由（`docs/ai-sdlc.md`〈為什麼是這幾個 skill〉）是怕 agent 寫規格時順手補答案；
那條防線早已換成機械的：`check_spec.py` 第 14 項抓沒有來源、或來源不是已答題的 FR／AC。分開的是**誰問**
（`discovery`）與**誰寫規則**（`spec`）；不必分開的是**檔案**。試跑退回 discovery 三次，每次都在兩份檔之間
改狀態。

## 決定

### D1. 一個 feature 一份 `spec.md`

`specs/ongoing/<date>-<feature>/` 底下不再有 `discovery.md`。`spec.md` 收進它的兩樣東西：

| 節 | 從哪來 | 位置 |
| --- | --- | --- |
| `## Open Questions` | `discovery.md` 的問答表＋每題小節（狀態、面向、答案、誰答的、原話、決策史） | 照原順序，取代「要送回去問的洞」清單 |
| `## 原文` | `discovery.md` 的 `## 原文`，逐字、§ 編號 | 最後，附錄 |

題號維持 `Q-<n>`：FR／AC 照舊標 `← Q-<n>`，`#### Q` 照舊只放待答的指標。`check_spec.py`（第 10、11、14 項與
`## 原文` 的檢查）與 `status.py` 改讀 `spec.md`。寫規格時發現的洞直接在 `## Open Questions` 開一題 `待答`，
影響到的 story 就地標 `[NEEDS CLARIFICATION]`。

### D2. 狀態

| 狀態 | 意思 | 誰寫 |
| --- | --- | --- |
| `澄清中` | 還有該問的沒問完 | `discovery` 建檔時；`spec` 開了新題時改回 |
| `待對焦` | 前沿空了（或喊停），規則寫好了 | `spec` |
| `已對焦` | 需求方說就緒 | `example-mapping` |
| `已凍結` | 尚無寫入者 | — |

`已收斂` 不再是一個狀態：`discovery` 問完時寫 `結束方式`（`收斂`／`喊停`／`阻塞`／`零進展`），狀態留 `澄清中`，
由 `spec` 寫完規則時改 `待對焦`。

### D3. `discovery` 問（像 interview-me、照 clarify 的流程），`spec` 寫

```
discovery   建 spec.md（原文、粗切 story、核心關係人、空的 Open Questions）→ 一輪一輪問 → 結束方式
spec        照 Open Questions 的已答題寫 FR／AC → 收尾檢查 → 自檢 → Domain Concepts → 診斷 → 旅程 → 就緒問句
            寫的時候發現洞 → Open Questions 開一題 → 狀態改回 澄清中 → STOP，告訴使用者跑 /discovery
```

`discovery` 寫 `spec.md` 的 `## 原文`、`## Open Questions`、`## Document Overview` 的關係人與結束方式、
`## User Stories` 的粗切（只有 `### US-<n> · <slug>` 與「身為…」一句）；**不寫 FR／AC**。`spec` 不改
`## 原文`，不自己填答案。

一題的寫法採兩個參考：

| 來自 | 做法 |
| --- | --- |
| interview-me | **一次只問一題**，附上自己的猜測與理由，等回答再問下一題（下一題常取決於這一題）；開場先寫一句假設與信心（缺什麼），信心隨答案更新；**停在 95%**：預測得出需求方對接下來三題的反應，且十六個面向都掃過；聽出「應該要的」而不是「真正要的」時追問；「你決定就好」不算答案；問完回寫意圖陳述（Outcome／User／Why now／Success／Constraint／Out of scope），要需求方明說 yes，然後停 |
| aixbdd clarify | 每題先給判斷所需的上下文、一句問題、推薦答案與理由；依影響排序（試跑第 1 輪一次問了 15 題——一次一題取代整批） |
| 本專案原有 | 前沿與十六個面向、能自己查的不准問、答案原話與誰答的、稽核分類是否窮盡、代演、喚醒 |

`clarify` 的獨有規則（每題分可解／阻塞／失效三堆、每輪結束給進度表）併進 `discovery`，`clarify` 整個刪掉。

### D4. 檔案怎麼搬

| 檔 | 去處 |
| --- | --- |
| `discovery/references/discovery-format.md` | 內容併進 `spec/references/spec-format.md` 的 `## Open Questions`、`## 原文`；原檔 `git rm` |
| `discovery/examples/minimal-discovery.md` | 內容併進 `spec/examples/minimal-spec.md`；原檔 `git rm` |
| `clarify/` | 有用的規則併進 `discovery/SKILL.md`；`git rm -r` |
| 其餘 `discovery/` 的檔 | 原位不動 |

### D5. 下游

| skill | 改什麼 |
| --- | --- |
| `plan` | 技術面向的已答題從 `spec.md` 的 `## Open Questions` 讀 |
| `feature-mapping`、`example-mapping`、`spec-by-example` | 缺口在 `## Open Questions` 開成 `Q-<n>` 待答題，不是一列描述 |
| `story-splitting` | 粗切寫進 `spec.md` 的 `## User Stories` |
| `archive`、`formulation`、`implement` | 不再提 `discovery.md` |

### 不在這次

- `system-analysis`
- `tests/` 的 fixture 仍有停車場文字（skill 不讀它們）

## 驗證

- `audit_skill.py` 每個 skill 通過；`claude plugin validate .` 通過；`python3 -m unittest discover -s tests` 通過
  （`check_spec.py` 與 `status.py` 的測試改成單檔 fixture）
- `grep -rn "discovery\.md" skills tests` 只剩說明舊佈局的遷移句
- 重跑 `parking-billing` 試跑：`discovery` → `spec` → `formulation` → `plan`，比對上一次的 grading 表
