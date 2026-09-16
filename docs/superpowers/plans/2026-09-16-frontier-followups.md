# 前沿模型的五項 follow-up Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 清掉四份計畫（`bdd-formulation` 拆出、`spec.md` 加四節、follow-up 清理、前沿模型）留下的五項，其中一項是**一條宣稱了沒有東西能履行的核對**。

**Architecture:** 四個任務。Task 1 修最重要的那一項——`clarify-log.md` 的 `n/a` 表加一欄，並讓兩個 skill 對「不該問的問題」處置一致。Task 2 給 `example-mapping` 唯一沒有祈使式呼叫的那個模式一個呼叫者。Task 3 改 `CLAUDE.md` 的版面。Task 4 處理 benchmark，**而它有一半不是我能決定的**。

**Tech Stack:** Markdown skill 文件、Python 3 稽核腳本、`claude plugin validate`。沒有 CI，全部手跑。

**Spec:** 無單一 spec。依據是 `docs/superpowers/specs/2026-09-14-grilling-and-formulation-design.md` 落地之後四輪 final review 累積的 park 項，逐項在下方複述並標明出處與現況（已對 `main` 的 `e7d2eda` 驗證過）。

## Global Constraints

- **`status.py` 與 `check_spec.py` 一行都不准改。** 兩支腳本都已完成並凍結。三支腳本對 fixture 的輸出從頭到尾必須不變。
- **`spec.md` 的十六節不增不減、標題文字不動、順序不動。** AC 的層級不動。
- **16 個追問面向不增不減、不改名。** 它們是 `status.py` 的 `DIMS`、`references/technical-probes.md`、`bdd-discovery/SKILL.md` 與 `spec-format.md` 四方的共用契約。
- 接縫不變：`FR-<n>` ↔ `@rule-<n>`、`AC-<n>.<m>` ↔ `@example-<n>.<m>`。
- 散文用繁體中文；commit message 用英文，conventional-commit 標題加 `WHAT:` / `WHY:` / `HOW:` 正文，結尾兩行：
  `Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>`
  `Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN`
- skill 目錄只准有 `rules/`、`references/`、`examples/`、`scripts/`（S3，MUST）；子目錄檔案要被同 skill 指到（S5，MUST）；`SKILL.md` 本文超過 500 行拿 S10（SHOULD）。
- 提交前清 `__pycache__`：`find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null`

### 基線（動手前已量過）

| 量什麼 | 現值 |
| --- | --- |
| `audit_skill.py` 跑 8 個 skill | 全部 exit 0，沒有任何 S10 |
| `claude plugin validate .` | 通過，**恰好 1 個警告**（plugin root 的 `CLAUDE.md`） |
| `status.py` 對 fixture | `合計  2  1  1`，exit 0 |
| `check_spec.py`（空 `features/`） | `4 個問題`，exit 1 |
| `check_spec.py`（無 `features/`） | `全部通過`，exit 0 |
| `skills/bdd-discovery/SKILL.md` | 458 行（本文 443） |
| `skills/bdd-discovery/references/clarify-log-format.md` | 86 行 |
| `skills/example-mapping/SKILL.md` | 162 行 |

**fixture** 建法固定：

```bash
FA="$(mktemp -d)/fx"; mkdir -p "$FA/specs/2026-09-09-parking" "$FA/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FA/specs/2026-09-09-parking/spec.md"
```

**這份計畫不碰任何被腳本解析的檔案**——`clarify-log.md` 沒有任何腳本讀它（已驗證：`grep -rln "clarify-log" skills/*/scripts/` 零命中），`spec.md` 這一份完全不動。三支腳本的輸出必須紋風不動。

---

## 五項與它們的出處

| # | 是什麼 | 出處 | Task |
| --- | --- | --- | --- |
| 3 | `bdd-discovery/SKILL.md:365-367` 的覆蓋率核對要求十六個面向每一個都出現在某一列的 `面向` 欄「**包含標 `n/a` 的**」，但 `clarify-log.md` 的 `n/a` 表是 `Q ／ 問題 ／ 為什麼不用問`，**沒有 `面向` 欄**——標成 n/a 的面向永遠滿足不了那條核對 | 計畫四 final 的 scoped re-review | 1 |
| 6 | `bdd-discovery/SKILL.md:301` 說「這題不該問 → 搬進『n/a』表」，`clarify-loop/SKILL.md:87` 說「這題不該問 → **刪掉它**，不要留著佔位」——同一個決定，兩個相反的處置 | 計畫四 final fix wave | 1 |
| 5 | `example-mapping` 宣告的四個模式裡，`流程內·對照 spec.md` 是唯一沒有祈使式呼叫的——它只能從 `bdd-spec/SKILL.md:228` 一句**描述句**到達 | 計畫四 final review | 2 |
| 7 | `CLAUDE.md:62-63` 兩行共用 `python3 skills/bdd-` 前綴，**六個獨立的 agent 誤讀過**，全都以為 `check_spec.py` 還在 `bdd-spec` 底下。文字是對的 | 計畫一到四，累計六次 | 3 |
| 4 | `benchmark/cases/*/grading.md` 用 Pass 1／Pass 2 當評分層級，那個結構已經不存在 | 計畫四 Ruling B | 4 |

---

## Task 1: `n/a` 表加 `面向` 欄，並讓兩個 skill 對「不該問」處置一致

這兩項是同一件事的兩半：一條核對需要 `n/a` 的列帶著面向，而另一個 skill 正在叫人把那些列刪掉。

**Files:**
- Modify: `skills/bdd-discovery/references/clarify-log-format.md`
- Modify: `skills/clarify-loop/SKILL.md`
- Modify: `skills/bdd-discovery/SKILL.md`（只在需要時，見 Step 4）

**Interfaces:**
- Consumes: 無
- Produces: 無

- [ ] **Step 1: 先證明那條核對現在做不到**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
echo "--- 核對要求什麼 ---"; sed -n '363,368p' skills/bdd-discovery/SKILL.md
echo "--- 待答表有 面向 欄 ---"; sed -n '20,24p' skills/bdd-discovery/references/clarify-log-format.md
echo "--- n/a 表沒有 ---"; sed -n '38,42p' skills/bdd-discovery/references/clarify-log-format.md
echo "--- 有腳本讀 clarify-log.md 嗎 ---"
grep -rln "clarify-log" skills/*/scripts/ 2>/dev/null || echo "零命中——加欄不影響任何解析"
```

Expected: 核對說「包含標 `n/a` 的」；`待答` 表有 `面向` 欄；`n/a` 表只有三欄且沒有它；沒有腳本讀這個檔。**這就是這一步要修的東西：一條要求了格式提供不了的東西的核對。**

- [ ] **Step 2: `n/a` 表加 `面向` 欄**

`skills/bdd-discovery/references/clarify-log-format.md` 的 `## n/a` 一節，表頭從
`| Q | 問題 | 為什麼不用問 |` 改成 `| Q | 問題 | 面向 | 為什麼不用問 |`，範例列跟著補。

欄位順序**照同一份檔案裡另外兩張表的既有慣例**——`## 已答`（`:28`）是
`Q ／ 問題 ／ 面向 ／ 歸屬 ／ 答案 ／ 誰答的`，`## 待答`（`:34`）是
`Q ／ 問題 ／ 面向 ／ 該問誰 ／ 為什麼還沒答`。**兩張表的 `面向` 都在第三欄**，
所以 `n/a` 也放第三欄。同名欄位在每張表的相對位置一致，讀的人才不必每張表重新找。

`面向` 的值域**完全不變**：十一個業務面向、五個技術面向，或 `—`。這一欄在 `n/a` 表裡
裝的是「**這一列原本要問的是哪個面向**」——正是覆蓋率核對要數的東西。

同一節補一句說明它為什麼在這裡：**標 `n/a` 不等於那個面向沒被想過，而覆蓋率核對要分開
這兩件事**。沒有這一欄，「問過了、判定不適用」跟「從來沒想到」在檔案上長得一樣。

- [ ] **Step 3: 讓 `clarify-loop` 與 `bdd-discovery` 對「這題不該問」講同一件事**

兩處現況：

| 檔案 | 現在說 |
| --- | --- |
| `skills/bdd-discovery/SKILL.md:301` | 這題不該問 → 回頭確認它真的不影響任何行為。是 → **搬進『n/a』表**，寫清楚為什麼不用問 |
| `skills/clarify-loop/SKILL.md:87` | 這題不該問 → 先確認它真的不影響任何規則。是 → **刪掉它**，不要留著佔位 |

**先讀完 `clarify-loop` 整份再決定**，不要假設哪一邊該讓。兩個 skill 的處境不同：
`bdd-discovery` 在建立那份 log，覆蓋率核對靠 `n/a` 的列存在；`clarify-loop` 處理的是
**一疊已經在檯面上的題**，它說的「刪掉」可能指的是**從它自己那一疊裡拿掉**，不是從
log 裡刪除那一列。

判準：**刪掉之後，`bdd-discovery` 的覆蓋率核對還數得到那個面向嗎？** 數不到就是錯的，
不管那句話原本想講什麼。

改完把你的判斷寫進報告：兩邊原本是不是真的矛盾、還是措辭造成的誤讀、你選了哪一種
統一寫法。

- [ ] **Step 4: 檢查 `bdd-discovery` 那一側要不要跟著改**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "n/a" skills/bdd-discovery/SKILL.md
```

逐行看過。Step 2 之後 `n/a` 的列多了一欄，所以任何一處**描述那張表有哪些欄位**的句子
都要跟著改；只是**提到 `n/a` 這個表**的句子不必動。這兩種句子長得很像。

**不要在 `SKILL.md` 裡重述欄位清單**——格式住在 `clarify-log-format.md`，這裡指路就好。
這四份計畫花了大半力氣在刪同一件事的第二份副本。

- [ ] **Step 5: 驗收**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
echo "--- n/a 表現在有 面向 欄 ---"
sed -n '/^## n\/a/,/^## /p' skills/bdd-discovery/references/clarify-log-format.md | grep '^|'
echo "--- 兩邊對「這題不該問」一致了嗎 ---"
grep -n "這題不該問" -A 1 skills/bdd-discovery/SKILL.md skills/clarify-loop/SKILL.md
echo "--- 覆蓋率核對現在做得到了嗎（人讀）---"
sed -n '363,368p' skills/bdd-discovery/SKILL.md
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for d in skills/*/; do printf "%-18s" "$(basename $d)"; \
  python3 skills/skill-rules/scripts/audit_skill.py "$d" >/dev/null 2>&1 && echo "exit 0" || echo "exit $?"; done
FA="$(mktemp -d)/fx"; mkdir -p "$FA/specs/2026-09-09-parking" "$FA/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FA/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-spec/scripts/status.py "$FA" | grep "^合計"
python3 skills/bdd-formulation/scripts/check_spec.py "$FA" 2>&1 | tail -1
git diff --stat -- skills/bdd-spec/scripts skills/bdd-formulation/scripts; echo "(空=腳本沒動)"
```

Expected: `n/a` 表四欄且 `面向` 在第三；兩處對「這題不該問」講同一件事；8 個 skill 全 exit 0；`合計  2  1  1`；`4 個問題`；腳本零變更。

---

## Task 2: 給 `流程內·對照 spec.md` 一個祈使式呼叫

**Files:**
- Modify: `skills/bdd-spec/SKILL.md`

**Interfaces:**
- Consumes: 無
- Produces: 無

- [ ] **Step 1: 確認現況**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "example-mapping" skills/bdd-spec/SKILL.md
sed -n '/^## 四種模式/,/^$/p' skills/example-mapping/SKILL.md | grep '^|'
```

Expected: `bdd-spec` 只在一處提到 `example-mapping`，而且是**描述句**（「之後 `已對焦` 由
`example-mapping` 依使用者的答覆改寫」），不是叫任何人去跑它。四個模式裡另外三個都有
祈使式的呼叫者：`逼問中` 與 `前沿空了` 在 `bdd-discovery`，`單獨` 由使用者直接觸發。

- [ ] **Step 2: 決定它該不該有呼叫者，再決定怎麼寫**

**這一步是判斷題，不是填空題。** 兩種結論都可能對：

| 結論 | 長什麼樣 |
| --- | --- |
| **該有** | `bdd-spec` 寫完 `spec.md` 之後，流程加一步：跑 `example-mapping` 的 `流程內·對照 spec.md`，攤地圖、跑診斷、把就緒問句交給使用者，`狀態` 由它依答覆改寫 |
| **不該有** | 這個模式本來就是**使用者主動觸發**的（「手上有一份 `spec.md`，想看它的形狀」），跟 `單獨` 模式同一類。那樣的話要在 `example-mapping` 的模式表裡**說清楚它由誰觸發**，而不是留給讀者猜 |

判準：**`spec.md` 寫完之後，就緒判定是流程的一部分，還是使用者想看才看？**
證據在 `bdd-spec/SKILL.md` 的 `## 完成後` 與 `example-mapping` 的 `## 就緒判定` 裡，
兩邊都讀過再決定。`狀態` 欄的 `待對焦 → 已對焦` 這個轉換由誰觸發，是同一個問題的另一面。

**寫進報告：你選了哪一種、證據是什麼。** 如果你的結論是「不該有」，那要改的檔案就
變成 `skills/example-mapping/SKILL.md`——**回報，不要自己擴張 Files 清單**。

- [ ] **Step 3: 驗收**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
echo "--- 四個模式各自的呼叫者 ---"
grep -n "流程內·\|example-mapping" skills/bdd-spec/SKILL.md skills/bdd-discovery/SKILL.md skills/example-mapping/SKILL.md
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
python3 skills/skill-rules/scripts/audit_skill.py skills/example-mapping; echo "exit=$?"
```

Expected: 四個模式每一個都說得出誰觸發它；兩個稽核都 exit 0。

---

## Task 3: `CLAUDE.md` 的兩行

六個獨立的 agent 誤讀過同一行。**文字是對的**——不要改成別的意思。

**Files:**
- Modify: `CLAUDE.md`

**Interfaces:**
- Consumes: 無
- Produces: 無

- [ ] **Step 1: 看清楚誤讀是怎麼發生的**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -n '56,65p' CLAUDE.md
```

兩行並排，共用 `python3 skills/bdd-` 這個前綴，只有中間一個詞不同：

```
python3 skills/bdd-spec/scripts/status.py [root]                # ...
python3 skills/bdd-formulation/scripts/check_spec.py [root]     # ...
```

一個在找「`check_spec.py` 的路徑」的讀者，很容易錨定在上一行的 `bdd-spec` 然後帶下來。
**六次誤讀指向版面，不是內容。**

- [ ] **Step 2: 改版面，不改事實**

三個方向，選一個並說明理由：

| 方向 | 代價 |
| --- | --- |
| 兩行分成兩個區塊，各自一句話說它屬於哪個 skill | 最清楚，但那一節會變長 |
| 調換順序，`check_spec.py` 排前面 | 最小改動，但只是把錨定的對象換一個 |
| 在路徑裡把不同的那一段標出來（例如加粗或反引號） | 在 `bash` 區塊裡辦不到——那是程式碼，不能有 Markdown 標記 |

第三個列出來是因為它是最直覺的那個，而它**在這個位置行不通**——先想清楚再動手。

**上面那句「Both take a project root and default to `.`; both exit 1 and name what is
missing...」不要動。** 前一份計畫判定過它仍然成立（真正的零輸入路徑照樣 exit 1），
現在沒有新事實推翻那個判定。

- [ ] **Step 3: 驗收**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -n '56,68p' CLAUDE.md
echo "--- 路徑還是對的嗎 ---"
grep -o "skills/[a-z-]*/scripts/[a-z_]*\.py" CLAUDE.md | while read -r p; do
  printf "%-52s %s\n" "$p" "$([ -f "$p" ] && echo 存在 || echo '*** 不存在 ***')"
done
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
claude plugin validate . 2>&1 | tail -3
```

Expected: 每一條路徑都存在；`validate` 恰好 1 個警告。

---

## Task 4: benchmark 的評分層級

**這個任務有一半不該由我決定，那一半要停下來回報。**

**Files:**
- Modify: `benchmark/cases/parking-billing/grading.md`（只有標籤，見下）
- Modify: `benchmark/cases/personal-memo/grading.md`（同上）

**Interfaces:**
- Consumes: 無
- Produces: 無

- [ ] **Step 1: 分清楚哪些是標籤、哪些是判準**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "Pass" benchmark/cases/parking-billing/grading.md benchmark/cases/personal-memo/grading.md
```

五處，兩種性質：

| 處 | 內容 | 性質 |
| --- | --- | --- |
| `parking-billing:39`、`personal-memo:49` | `### 第一層 — Pass 1 就該問（答案改變切分方式）` | **標籤**——括號裡的「答案改變切分方式」才是判準，`Pass 1` 只是那個判準的名字 |
| `parking-billing:47`、`personal-memo:56` | `### 第二層 — Pass 2 深度追問` | **標籤**——同上 |
| `parking-billing:91` | `\| 分層 \| Pass 1 只問改變切分的題，沒有一開始就追費率細節 \|` | **判準本身**——它評的是「有沒有在該問範圍題的時候去追細節」 |

**前四處可以改，第五處不行。** 前四處把 `Pass 1`／`Pass 2` 換成不指涉已刪結構的名字，
括號裡的判準一字不動——**得分條件完全不變**，所以舊的評估結果與新的仍然可比。

- [ ] **Step 2: 改前四處**

層級名稱要講它**評的是什麼**，不是它對應流程的哪一段。第一層評的是「該問範圍題的時候
有沒有問範圍題」，第二層評的是「深度」。用判準本身當名字，這樣下次流程再改，這份
答案卷也不必跟著動——**評分標準綁在行為上，不綁在實作結構上**，本來就比較耐放。

- [ ] **Step 3: 第五處——停下來回報，不要自己決定**

`parking-billing:91` 的 `分層` 這一列評的是：「**Pass 1 只問改變切分的題，沒有一開始
就追費率細節**」。

前沿模型裡沒有一個「只問範圍題」的階段。範圍題之所以先問，是因為**其他幾乎所有題的
前提都掛在它身上**——那是算出來的結果，不是規定。所以這一列的判準要改成什麼，
是一個**評估設計的決定**：前沿模型底下，「一開始就追費率細節」還算不算失分？
如果算，靠什麼判斷？

**不要自己改它。** 把這一列原文、你的分析、以及你認為的兩三種可能寫法寫進報告，
交給使用者裁決。理由：寫這份計畫的人**讀過這兩份答案卷**，由他來重寫一條評分判準，
會有把標準往新流程擅長的方向挪的風險——而那種偏差在評估裡是最貴的，因為它不會
表現成錯誤，只會表現成分數變好。

- [ ] **Step 4: 驗收**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
echo "--- 剩下的 Pass 應該只有 parking-billing:91 一處 ---"
grep -n "Pass" benchmark/cases/*/grading.md
echo "--- 判準（括號裡那句）有沒有被動到 ---"
git diff -- benchmark/cases/ | grep "^[-+][^-+]"
```

Expected: 只剩一處；diff 裡每一行都只改了層級名稱，括號裡的判準一字未動。

---

## Self-Review

**1. 覆蓋**

五項全部有歸屬：Task 1 收第 3、6 項，Task 2 收第 5 項，Task 3 收第 7 項，Task 4 收第 4 項的可做部分並把不可做的那半交回去。

**2. 佔位字掃描**

沒有 TBD／TODO。Task 1 Step 2 給了完整的欄位形狀與位置理由。Task 1 Step 3、Task 2 Step 2、Task 3 Step 2 給的是**判準與選項**而非逐字稿——那三處都需要讀過上下文才判得出來，逐字指定會產出接不上的段落，而每一處都附了「寫進報告」與「不確定就回報」的出口。

**3. 任務相依**

四個任務彼此獨立，順序隨意。Task 1 的兩項之間有相依（欄位要先存在，處置才統一得起來），已排在同一個任務的相鄰步驟。

**4. 這份計畫刻意不做的兩件事**

- **不跑 benchmark。** 那需要一個沒讀過 `grading.md` 的 session；寫這份計畫的人對兩個案例都已污染。
- **不碰 `clarify-loop` 的同名衝突**（`~/.claude/skills/clarify-loop` 與 `skills/clarify-loop`）。那是安裝層的問題，不是文件問題。

**5. 我在這份計畫裡最可能犯的錯**

前四份計畫我的缺陷全是同一族：**驗證範圍與施工範圍對不齊**，而最後一次最徹底——「必須活下來」清單的完整性，上限是寫清單的人讀過多少。這一份的對策：

- 每一項的現況都是剛才**實際跑指令量過的**，不是憑記憶（`n/a` 表的三欄、五處 Pass 的行號與性質、`bdd-spec` 只有一處提到 `example-mapping`）。
- Task 1 Step 1、Task 3 Step 1 都是**先證明缺陷存在**再修。
- Task 4 Step 1 把五處**逐處分類**成標籤或判準，而不是說「把 Pass 換掉」——這是這份計畫裡唯一一個「看起來可以整批取代、實際上不行」的地方。
- Task 1 Step 4 明講「描述欄位的句子要改、只提到表的句子不必」，因為那兩種句子長得很像，而這正是前一份計畫在 `example-mapping` 上踩過的那個坑。

**而我在寫完之後的自我複查裡，又踩了一次同族的坑——這次是在檢查裡，不是在計畫裡。**
我原本寫「欄位順序照 `待答` 表的慣例」，然後用 `sed -n '20,22p'` 去驗它——那三行是
`## 核心關係人` 表，不是 `待答`。主張本身是對的（`待答` 在 `:34`，`面向` 確實是第三欄），
而且事實比我寫的更強：`已答` 與 `待答` **兩張**都把 `面向` 放第三。**引用一個行範圍卻
不讀它周圍，是我這個 session 重複最多次的失誤**——計畫二把 `歸屬` 欄寫成 `面向` 欄是
同一個成因。已改成同時引用兩張表與它們的真實行號。
