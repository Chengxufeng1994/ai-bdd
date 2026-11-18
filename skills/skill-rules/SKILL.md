---
name: skill-rules
description: >
  本專案撰寫與稽核 SKILL 的規範：資料夾結構、frontmatter 寫法、觸發詞設計、
  Skill Boundaries、品質清單，並可對現成 skill 做逐項稽核。
  觸發詞：「寫一個 skill」「建立技能」「新增 skill」「skill 要怎麼寫」
  「SKILL.md 格式」「技能規範」「skill 規範」「檢查這個 skill」
  「稽核 skill」「這個 skill 有問題嗎」「skill 的資料夾結構」。
  English: write a skill, create a new skill, author a SKILL.md,
  skill conventions, review this skill, audit a skill, skill folder structure.
disable-model-invocation: true
---

# SKILL 撰寫與稽核規範

| | |
| --- | --- |
| **職責** | 定本專案 skill 的房規（結構、frontmatter、規則關鍵字），並對現成 skill 逐項稽核 |
| **讀** | 稽核：`skills/<name>/` 整個目錄；寫新 skill：既有 skill 的 `## Skill Boundaries`（找重疊） |
| **寫** | 稽核：不寫檔，報告照 `rules/quality-checklist.md` 的格式；寫新 skill：`skills/<name>/SKILL.md` 與有內容的子目錄 |
| **上游** | 無（不在流程裡；寫或改任何 skill 之前讀） |
| **下游** | 無 |
| **不做** | 不跑 eval 迴圈、不最佳化觸發率（都是 `skill-creator`）；職責重疊時不自己決定合併或劃界 |

本專案所有 skill 的房規。寫新 skill 時照著長，改既有 skill 時照著檢查。

## 使用時機

- 要新增一個 skill，需要知道資料夾怎麼擺、SKILL.md 要有哪些章節
- 寫完一個 skill，想知道它合不合格
- 既有 skill 觸發不了、或在不該觸發時觸發
- 想確認 description 該寫多長、觸發詞該放幾個

## Skill Boundaries

| 情況 | 改用 |
| --- | --- |
| 要**完整的建立流程**（產草稿 → 跑測試案例 → 人工評分 → 迭代） | `skill-creator`，那是一整套 eval 迴圈，本 skill 不做 |
| 要**最佳化 description 的觸發率**（自動跑 20 題觸發測試） | 也在 `skill-creator` 裡 |

本 skill 只負責兩件事：**規範**（怎麼寫）與**稽核**（寫得對不對）。兩者可以接力：先用
`skill-creator` 迭代出內容，再用本 skill 校正結構與房規。

## 參考檔案

- `rules/` — `structure.md`（目錄與章節 S1–S16、章節順序、範本）、`frontmatter.md`（`description` 與觸發詞）、
  `keywords.md`（標強度）、`quality-checklist.md`（人工判斷的完整清單與稽核報告格式）
- `references/runtime-facts.md` — 執行期行為的證據（巢狀目錄、命名空間）
- `examples/anti-patterns.md` — 反面教材，每一項都是實際會導致 skill 失效的寫法；`examples/minimal-skill.md` —
  合格的最小範本，可以直接複製改
- `scripts/audit_skill.py` — 機械檢查

## 流程

### 稽核一個 skill

先跑機械檢查，再做人工判斷。順序不要顛倒——機械檢查會先淘汰掉大部分問題，省下人工判斷的力氣。

#### 1. 機械檢查

```bash
python3 scripts/audit_skill.py <skill 目錄路徑>
```

檢查目錄結構、frontmatter 欄位、name 與目錄名是否一致、必要章節是否存在、是否有未被 SKILL.md
引用的孤兒參考檔。若該 skill 屬於某個 plugin，再跑：

```bash
claude plugin validate <plugin 路徑>
```

本專案的閘門是不加 `--strict` 的版本：`--strict` 會因為 plugin 根目錄的 `CLAUDE.md` 失敗，那是
刻意的取捨（見 repo 的 `CLAUDE.md`）。想看全部警告時再加。

IMPORTANT: `claude plugin validate` **不會遞迴巢狀目錄**。放在 `skills/<分類>/<名稱>/` 的 skill 在執行時
有效，但永遠不會被驗證。這是選擇扁平結構的實質理由，不只是美觀。詳見 `references/runtime-facts.md`。

#### 2. 人工判斷

機械檢查過不了的東西，才需要讀：

- 觸發詞涵蓋動詞變化、名詞同義、口語與正式、中英文了嗎？
- Skill Boundaries 指到的 skill 真的存在嗎？
- 流程步驟具體到能執行，還是停在「分析一下」？
- 有沒有解釋**為什麼**，還是只有一串命令？
- 長流程有沒有使用者檢查點？

完整清單見 `rules/quality-checklist.md`。

### 寫新 skill 的流程

| 步驟 | 做什麼 | 要問使用者嗎 |
| --- | --- | --- |
| 1 | **確認它不該是別的東西**。只是一段常用提示詞 → 不需要 skill。只有你會手動叫 → 可能該是 command | — |
| 2 | **先寫 `description`**。寫不出清楚的觸發條件，代表這個 skill 的職責還沒想清楚 | — |
| 3 | **寫 `## Skill Boundaries`**，特別是與既有 skill 的分界。這一步會逼你面對重疊 | 職責明顯重疊時要 |
| 4 | 寫 `## 流程` 與 `## 產物`，再寫 `## Verification`——寫不出可判是或否的項目，代表流程還沒說清楚「做完」長什麼樣 | — |
| 5 | 判斷需不需要子目錄；需要時建，並在 SKILL.md 裡明確指路 | — |
| 6 | 跑稽核 | — |

Ask user: 若新 skill 與既有 skill 的職責明顯重疊，先問使用者要合併還是劃界，不要自己決定——重疊的兩個
skill 會互搶觸發，而使用者無從得知模型挑了哪一個。

## 房規

### 資料夾結構

```
<skill-name>/
├── SKILL.md        必備
├── rules/          該 skill 執行時要套用到產出上的規則
├── references/     模型需要時才讀的說明文件
├── examples/       正面與反面範例
└── scripts/        機械性工作的程式碼
```

**目錄名就是 skill 的識別名**，必須與 frontmatter 的 `name` 一致。

四個子目錄**有內容才建**。空目錄 git 不追蹤，佔位檔只是雜訊。原則上只用這四種（S3 是 `SHOULD NOT`）——需要第五種時，
先問這東西是不是其實屬於前四種之一，真的不是就明講理由（`rules/structure.md`〈S3 為何是 `SHOULD NOT` 而非 `MUST NOT`〉）。

#### rules/ 與 references/ 的分界

兩者都是 Markdown，差別在**服務對象**：

| | 內容 | 判準 |
| --- | --- | --- |
| `rules/` | 產出必須符合的規則 | 「照這條做出來的東西對不對」 |
| `references/` | 模型理解任務所需的背景 | 「模型讀完更懂，但不直接約束產出」 |

例：Gherkin 的命名規範屬於 `rules/`；Gherkin 語法有哪些關鍵字屬於 `references/`。分不清時問自己：**這段文字能不能拿來判定一份產出合格與否？** 能，就是 rule。

#### 這些檔案不會自動載入

**只有 SKILL.md 的本文會在 skill 觸發時進 context。** 子目錄的檔案必須由 SKILL.md 明確指示去讀，否則
等於不存在。這是漸進式揭露省 token 的原理，也是寫 skill 時最常漏掉的一步——放了一堆參考檔卻從沒指路。

### SKILL.md 的必要章節

| 章節 | 作用 | |
| --- | --- | --- |
| frontmatter | `name` ＋ `description`（觸發的唯一依據） | 必要 |
| 職責卡 | 標題底下六列：職責／讀／寫／上游／下游／不做 | 建議（S16） |
| `## 使用時機` | 什麼情況該用它，以及**方向**（避免選錯） | 必要 |
| `## Skill Boundaries` | 什麼情況該改用別的，本 skill 不做什麼 | 必要 |
| `## 前置確認（Ask user）` | 缺哪些輸入就不能開始，以及為什麼需要 | 有輸入需求時 |
| `## 流程` | 可執行的步驟，不是「分析一下」這種模糊指令 | 必要 |
| `## 產物` | 產物長什麼樣 | 建議 |
| `## Verification` | 交回前逐項確認的 `- [ ]` 清單：做完了沒、做對了沒 | 必要 |

`## 前置確認` 防的是 skill 最常見的失效方式：**用假設的輸入產出真實的結果**。模型不會停下來說「我不知道」，
它會挑一個合理的猜測繼續做，錯誤就藏在一份看起來正常的成品裡。

細節、範本與章節順序見 `rules/structure.md`。

### frontmatter

`name` 用 kebab-case，全小寫，望名知義，且與目錄名一致。

`description` 是**觸發的唯一依據**，也是唯一每個 session 都付費的部分。寫法、觸發詞公式與數量取捨見
`rules/frontmatter.md`——那裡也說明了為什麼「觸發詞越多越好」在 skill 數量多的專案裡不成立。

### 規則關鍵字

規則要標示**強度**，用有共識的英文關鍵字，說明文字用繁體中文。由強至弱：

| 關鍵字 | 意義 |
| --- | --- |
| `MUST NOT` / `NEVER` | 絕對不可。完全禁止，沒有例外。 |
| `MUST` | 必須。硬性規定，沒有選擇餘地。 |
| `IMPORTANT` | 高優先。特別容易被漏掉、漏掉代價很高，但仍留有判斷空間。 |
| `SHOULD` | 應該。正確的路，但不是法律；偏離時要說明理由。 |
| `SHOULD NOT` | 不應該。不建議，有正當理由時仍可為之。 |
| `COULD` | 可以。軟性選項，由執行者決定。 |

輔助標記（不表示強度）：`Check if:` 標條件、`Ask user:` 標需要確認。

**強度必須稀有才有意義。** 一份 skill 通常只有一兩條真正的 `MUST` 與 `NEVER`，其餘都是 `SHOULD`；整份
都是 `MUST` 時模型分不出哪條真的不能違反，等同全部沒標。分級判準、寫法與常見錯誤見 `rules/keywords.md`。

## 產物

稽核是一份報告，照 `rules/quality-checklist.md` 的格式（必須修正／建議／通過）；寫新 skill 是 `skills/<name>/`，結構照〈房規〉。

## Verification

交回前逐項確認；任何一項答「否」就不交——修好，或在回報裡明說哪一項沒過：

- [ ] 跑了 `scripts/audit_skill.py`，輸出貼在回報裡；人工判斷是在它之後做的
- [ ] skill 屬於 plugin 時，也跑了 `claude plugin validate`，結果貼出
- [ ] 稽核報告用 `rules/quality-checklist.md` 的格式：必須修正／建議／通過，通過項沒有逐條列出
- [ ] 每一條「必須修正」都指到一條規則編號或清單項，並寫了怎麼改
- [ ] 新 skill 與既有 skill 職責重疊時，問過使用者要合併還是劃界，沒有自己決定
