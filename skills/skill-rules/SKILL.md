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
---

# SKILL 撰寫與稽核規範

本專案所有 skill 的房規。寫新 skill 時照著長，改既有 skill 時照著檢查。

## 使用時機

- 要新增一個 skill，需要知道資料夾怎麼擺、SKILL.md 要有哪些章節
- 寫完一個 skill，想知道它合不合格
- 既有 skill 觸發不了、或在不該觸發時觸發
- 想確認 description 該寫多長、觸發詞該放幾個

## Skill Boundaries

- 要**完整的建立流程**（產草稿 → 跑測試案例 → 人工評分 → 迭代）
  → 改用 `skill-creator`，那是一整套 eval 迴圈，本 skill 不做
- 要**最佳化 description 的觸發率**（自動跑 20 題觸發測試）
  → 也在 `skill-creator` 裡
- 本 skill 只負責兩件事：**規範**（怎麼寫）與**稽核**（寫得對不對）

兩者可以接力：先用 `skill-creator` 迭代出內容，再用本 skill 校正結構與房規。

---

## 資料夾結構

```
<skill-name>/
├── SKILL.md        必備
├── rules/          該 skill 執行時要套用到產出上的規則
├── references/     模型需要時才讀的說明文件
├── examples/       正面與反面範例
└── scripts/        機械性工作的程式碼
```

**目錄名就是 skill 的識別名**，必須與 frontmatter 的 `name` 一致。

四個子目錄**有內容才建**。空目錄 git 不追蹤，佔位檔只是雜訊。它們是唯一允許的
子目錄——需要第五種時，先問這東西是不是其實屬於前四種之一。

### rules/ 與 references/ 的分界

兩者都是 Markdown，差別在**服務對象**：

| | 內容 | 判準 |
| --- | --- | --- |
| `rules/` | 產出必須符合的規則 | 「照這條做出來的東西對不對」 |
| `references/` | 模型理解任務所需的背景 | 「模型讀完更懂，但不直接約束產出」 |

例：Gherkin 的命名規範屬於 `rules/`；Gherkin 語法有哪些關鍵字屬於 `references/`。
分不清時問自己：**這段文字能不能拿來判定一份產出合格與否？** 能，就是 rule。

### 這些檔案不會自動載入

**只有 SKILL.md 的本文會在 skill 觸發時進 context。** 子目錄的檔案必須由
SKILL.md 明確指示去讀，否則等於不存在。這是漸進式揭露省 token 的原理，也是寫
skill 時最常漏掉的一步——放了一堆參考檔卻從沒指路。

---

## SKILL.md 的必要章節

| 章節 | 作用 | |
| --- | --- | --- |
| frontmatter | `name` ＋ `description`（觸發的唯一依據） | 必要 |
| `## 使用時機` | 什麼情況該用它，以及**方向**（避免選錯） | 必要 |
| `## Skill Boundaries` | 什麼情況該改用別的，本 skill 不做什麼 | 必要 |
| `## 前置確認（Ask user）` | 缺哪些輸入就不能開始，以及為什麼需要 | 有輸入需求時 |
| 工作流程 | 可執行的步驟，不是「分析一下」這種模糊指令 | 必要 |
| 輸出格式 | 產物長什麼樣 | 建議 |

`## 前置確認` 防的是 skill 最常見的失效方式：**用假設的輸入產出真實的結果**。
模型不會停下來說「我不知道」，它會挑一個合理的猜測繼續做，錯誤就藏在一份看起來
正常的成品裡。

細節與範本見 `rules/structure.md`。

## frontmatter

`name` 用 kebab-case，全小寫，望名知義，且與目錄名一致。

`description` 是**觸發的唯一依據**，也是唯一每個 session 都付費的部分。寫法、
觸發詞公式與數量取捨見 `rules/frontmatter.md`——那裡也說明了為什麼「觸發詞越多
越好」在 skill 數量多的專案裡不成立。

## 規則關鍵字

規則要標示**強度**，用有共識的英文關鍵字，說明文字用繁體中文。

由強至弱：

| 關鍵字 | 意義 |
| --- | --- |
| `MUST NOT` / `NEVER` | 絕對不可。完全禁止，沒有例外。 |
| `MUST` | 必須。硬性規定，沒有選擇餘地。 |
| `IMPORTANT` | 高優先。特別容易被漏掉、漏掉代價很高，但仍留有判斷空間。 |
| `SHOULD` | 應該。正確的路，但不是法律；偏離時要說明理由。 |
| `SHOULD NOT` | 不應該。不建議，有正當理由時仍可為之。 |
| `COULD` | 可以。軟性選項，由執行者決定。 |

輔助標記（不表示強度）：`Check if:` 標條件、`Ask user:` 標需要確認。

**強度必須稀有才有意義。** 一份 skill 通常只有一兩條真正的 `MUST` 與 `NEVER`，
其餘都是 `SHOULD`；整份都是 `MUST` 時模型分不出哪條真的不能違反，等同全部沒標。
分級判準、寫法與常見錯誤見 `rules/keywords.md`。

---

## 稽核一個 skill

先跑機械檢查，再做人工判斷。順序不要顛倒——機械檢查會先淘汰掉大部分問題，
省下人工判斷的力氣。

### 1. 機械檢查

```bash
python3 scripts/audit_skill.py <skill 目錄路徑>
```

檢查目錄結構、frontmatter 欄位、name 與目錄名是否一致、必要章節是否存在、
是否有未被 SKILL.md 引用的孤兒參考檔。

若該 skill 屬於某個 plugin，再跑：

```bash
claude plugin validate <plugin 路徑> --strict
```

IMPORTANT: `claude plugin validate` **不會遞迴巢狀目錄**。放在
`skills/<分類>/<名稱>/` 的 skill 在執行時有效，但永遠不會被驗證。這是選擇扁平
結構的實質理由，不只是美觀。詳見 `references/runtime-facts.md`。

### 2. 人工判斷

機械檢查過不了的東西，才需要讀：

- 觸發詞涵蓋動詞變化、名詞同義、口語與正式、中英文了嗎？
- Skill Boundaries 指到的 skill 真的存在嗎？
- 流程步驟具體到能執行，還是停在「分析一下」？
- 有沒有解釋**為什麼**，還是只有一串命令？
- 長流程有沒有使用者檢查點？

完整清單見 `rules/quality-checklist.md`。反面教材見
`examples/anti-patterns.md`——那裡的每一項都是實際會導致 skill 失效的寫法。
合格的最小範本見 `examples/minimal-skill.md`，可以直接複製改。

---

## 寫新 skill 的流程

1. **確認它不該是別的東西**。只是一段常用提示詞 → 不需要 skill。只有你會手動叫
   → 可能該是 command。
2. **先寫 `description`**。寫不出清楚的觸發條件，代表這個 skill 的職責還沒想清楚。
3. **寫 `## Skill Boundaries`**，特別是與既有 skill 的分界。這一步會逼你面對重疊。
4. 寫工作流程與輸出格式。
5. 判斷需不需要子目錄；需要時建，並在 SKILL.md 裡明確指路。
6. 跑稽核。

Ask user: 若新 skill 與既有 skill 的職責明顯重疊，先問使用者要合併還是劃界，
不要自己決定——重疊的兩個 skill 會互搶觸發，而使用者無從得知模型挑了哪一個。
