# `bdd-formulation` 拆出來 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 Gherkin 的產出從 `bdd-spec` 整段搬到新的 `bdd-formulation` skill，`bdd-spec` 只留 `spec.md`／`docs/CONTEXT.md`／`domain-model.md` 三份產物。

**Architecture:** 十八個檔案先做**純搬移**（一筆提交、零內容變更，讓 `git log --follow` 活下來），再從還完整的 `bdd-spec/SKILL.md` **抄**出 `bdd-formulation/SKILL.md`，最後才**剪**掉 `bdd-spec/SKILL.md` 的 Gherkin 段落。先抄後剪，中間沒有任何一刻那段文字不存在於工作樹裡。第四個任務掃 repo 其餘指路。

**Tech Stack:** Markdown skill 文件、Python 3 稽核腳本（`audit_skill.py`、`check_spec.py`、`status.py`）、`claude plugin validate`。沒有 CI，全部靠手跑。

**Spec:** `docs/superpowers/specs/2026-09-14-grilling-and-formulation-design.md`（D8 是本計畫的依據；「計畫一」那一節是它的範圍聲明）

## Global Constraints

- **格式完全不動。** `spec.md` 還是現在的樣子，`.feature` 還是現在的樣子，只是產它的人換了。本計畫不得改動 `references/spec-format.md`、`examples/minimal-spec.md` 的任何一個位元組，也不得改 `check_spec.py`／`status.py` 的任何一行程式碼。
- **接縫不變：** `FR-<n>` ↔ `@rule-<n>`、`AC-<n>.<m>` ↔ `@example-<n>.<m>`。
- **搬移與改寫必須分屬不同提交。** git 的 rename 是從內容相似度**推斷**出來的；把搬移和大改寫綁進同一筆提交，`git log --follow` 會安靜地斷掉，而且不會有任何錯誤訊息。驗證方式是**比對 blob SHA**，不是 `git diff <range> -- <new-path>`——後者在 git 看不到舊路徑時會吐出整份檔案的 diff，看起來像改過但其實沒有。
- **散文用繁體中文，commit message 與 PR 用英文**，conventional-commit 標題加 `WHAT:` / `WHY:` / `HOW:` 正文。本計畫的 scope 用 `bdd-spec` 與 `bdd-formulation`。
- **搬移一律用 `git mv`。**
- 每個 skill 目錄底下只准有 `rules/`、`references/`、`examples/`、`scripts/` 四種子目錄（`audit_skill.py` 的 S3，MUST）。
- **每一個子目錄裡的檔案都必須被同一個 skill 內部指到**（S5，MUST）。指路鏈從 `SKILL.md` 出發，沿 `(?:rules|references|examples|scripts)/[\w./-]+` 這個正則遞迴展開。搬進新 skill 的檔案如果沒有在新的 `SKILL.md` 裡被指到，就是孤兒。
- `SKILL.md` 本文超過 500 行會拿到 S10（SHOULD，不影響離開碼）。兩份都要壓在 500 行以下。
- frontmatter 必備：`name`（kebab-case，且**等於目錄名**）、`description`（含「」包住的中文觸發詞、`English:` 起頭的英文觸發詞，長度 ≥ 40）。本文必備章節：`## 使用時機`、`## Skill Boundaries`（S6／S7，MUST）。

### 基線（動手前已量過，收工要對得上）

| 量什麼 | 現值 |
| --- | --- |
| `audit_skill.py` 跑 7 個 skill | 全部 exit 0；`bdd-spec` 有一項 S10 SHOULD（本文 622 行） |
| `claude plugin validate .` | 通過，**恰好 1 個警告**（plugin root 的 `CLAUDE.md`，這是刻意的取捨） |
| `skills/bdd-spec/SKILL.md` | 635 行，`grep -c "^#\{2,4\} "` ＝ **32**（其中 2 個是 code fence 裡的 `### US-1 · visitor-billing` / `### US-2 · monthly-pass-exit`，不是真標題） |
| `status.py` 對 fixture | 合計 `2 1 1`，exit 0 |
| `check_spec.py` 對 fixture | `4 個問題`，exit 1 |

**fixture** 指下面這個目錄，每個任務都會用到，建法固定：

```bash
FX="$(mktemp -d)/fx"
mkdir -p "$FX/specs/2026-09-09-parking" "$FX/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FX/specs/2026-09-09-parking/spec.md"
```

`features/` 是**空的**，這是故意的——`check_spec.py` 因此會報「兩則 story 各缺狀態 tag、各漏了例子」共 4 個問題。這份輸出是本計畫的迴歸基準：搬完之後一個字都不該變。

---

## File Structure

### 搬過去的十八個檔案（Task 1，純搬移）

| 從 | 到 |
| --- | --- |
| `skills/bdd-spec/references/step-grammar.md` | `skills/bdd-formulation/references/step-grammar.md` |
| `skills/bdd-spec/references/state-tags.md` | `skills/bdd-formulation/references/state-tags.md` |
| `skills/bdd-spec/references/rule-taxonomy.md` | `skills/bdd-formulation/references/rule-taxonomy.md` |
| `skills/bdd-spec/references/coverage-report.md` | `skills/bdd-formulation/references/coverage-report.md` |
| `skills/bdd-spec/references/artifact-location.md` | `skills/bdd-formulation/references/artifact-location.md` |
| `skills/bdd-spec/examples/1-video-progress.md` | `skills/bdd-formulation/examples/1-video-progress.md` |
| `skills/bdd-spec/examples/2-submit-assignment.md` | `skills/bdd-formulation/examples/2-submit-assignment.md` |
| `skills/bdd-spec/examples/3-deliver-course.md` | `skills/bdd-formulation/examples/3-deliver-course.md` |
| `skills/bdd-spec/examples/4-query-progress.md` | `skills/bdd-formulation/examples/4-query-progress.md` |
| `skills/bdd-spec/examples/5-query-course-list.md` | `skills/bdd-formulation/examples/5-query-course-list.md` |
| `skills/bdd-spec/examples/6-query-product-list.md` | `skills/bdd-formulation/examples/6-query-product-list.md` |
| `skills/bdd-spec/examples/7-create-order.md` | `skills/bdd-formulation/examples/7-create-order.md` |
| `skills/bdd-spec/examples/8-query-order.md` | `skills/bdd-formulation/examples/8-query-order.md` |
| `skills/bdd-spec/examples/9-pay-order.md` | `skills/bdd-formulation/examples/9-pay-order.md` |
| `skills/bdd-spec/examples/10-cancel-order.md` | `skills/bdd-formulation/examples/10-cancel-order.md` |
| `skills/bdd-spec/examples/11-query-user-roles.md` | `skills/bdd-formulation/examples/11-query-user-roles.md` |
| `skills/bdd-spec/examples/anti-patterns.md` | `skills/bdd-formulation/examples/anti-patterns.md` |
| `skills/bdd-spec/scripts/check_spec.py` | `skills/bdd-formulation/scripts/check_spec.py` |

**留在 `bdd-spec`：** `references/spec-format.md`、`references/persona-definition.md`、`examples/minimal-spec.md`、`scripts/status.py`。

> **規格對帳：** D8 的表格左欄寫「`examples/` 的 12 個 Gherkin 範例 ＋ `anti-patterns.md`」。實際目錄裡編號範例只有 11 個（`1-` 到 `11-`），加上 `anti-patterns.md` 共 **12 個檔**，`minimal-spec.md` 留下。上面的清單是權威的：`examples/` 搬 12 個檔，留 1 個。

搬移之間的互相指路**全部落在搬移集合內**，已經查過：`coverage-report.md`、`examples/10-cancel-order.md`、`examples/11-query-user-roles.md` 指向 `rule-taxonomy.md`（三者都搬）；`persona-definition.md` 指向 `spec-format.md`（兩者都留）。沒有跨邊界的指路鏈。

### `skills/bdd-spec/SKILL.md` 的逐段去向

行號以**動手前的 635 行版本**為準。

| 行 | 段落 | 去向 |
| --- | --- | --- |
| 1–15 | frontmatter | 留（改寫：description 拿掉 Gherkin） |
| 17–23 | 標題與開場 | 留（改寫） |
| 25–30 | `## 使用時機` | 留（改寫第 1、3 條） |
| 31–38 | `## Skill Boundaries` | 留（新增一條指向 `bdd-formulation`） |
| 39–47 | `## 參考檔案` | 留（原樣——這張表本來就只列 `spec-format.md`、`persona-definition.md`、`minimal-spec.md`、`status.py`） |
| 48–88 | `## 切 story` | 留（第 54、81 行的 `check_spec.py` 要改成指向 `bdd-formulation`） |
| 89–116 | `## 詞彙表 —— docs/CONTEXT.md` | 留（第 90 行的理由句要改寫，見 Ruling A） |
| 117–148 | `## 這份 .feature 是寫給誰讀的` | **搬** |
| 149–151 | `## 前置確認` 的「不必問。`.feature` 的位置偵測得到」 | **搬** |
| 152–165 | `## 前置確認` 的 MUST NOT／兩條 NEVER | 留 |
| 166–167 | `## 流程` | 兩邊各有一個 |
| 168–190 | `### 1. 挑 story` | 留（編號不變，仍是 1） |
| 191–218 | `### 2. 讀三樣東西` | 留（標題改成 `### 2. 讀兩樣東西`，表格刪掉「既有的 `.feature`」那一列與底下第 3 段，見 Task 3 Step 6） |
| 219–240 | `### 3. 決定 seam` | 留 |
| 241–262 | `### 4. 寫 spec.md` | 留 |
| 263–314 | `### 5. 規則對應 Rule，例子對應 Example` | **搬** |
| 315–334 | `#### 規則歸到四個象限` | **搬** |
| 335–355 | `#### 每條規則要看三種結果` | **搬** |
| 356–368 | `#### 場景名字要能定位` | **搬** |
| 369–398 | `### 規則與它自己的例子打架時` | **搬** |
| 399–430 | `### 6. 步驟文法 —— 四種形狀` | **搬** |
| 431–459 | `### 每一條規則` | **搬** |
| 460–487 | `### 7. 什麼時候用 Scenario Outline` | **搬** |
| 488–494 | `### 8. 寫檔` | **搬** |
| 495–507 | `### 9. 增修 specs/domain-model.md` | 留（**重新編號成 `### 5.`**） |
| 508–546 | `### 10. 稽核，然後才算完成` | **搬** |
| 547–574 | `## 產物格式` | 拆（見 Task 3 步驟 5） |
| 575–586 | `### 狀態 tag` | **搬** |
| 587–603 | `### 型別 tag` | **搬** |
| 604–616 | `### 方言陷阱` | **搬** |
| 617–623 | `## 產物隔離` | 拆（`artifact-location.md` 跟著走，`bdd-spec` 自己重寫一段） |
| 624–635 | `## 完成後` | 拆（見 Task 3 步驟 7） |

整段搬走的標題有 **14 個**：`## 這份 .feature 是寫給誰讀的`、`### 5.`、`#### 規則歸到四個象限`、`#### 每條規則要看三種結果`、`#### 場景名字要能定位`、`### 規則與它自己的例子打架時`、`### 6.`、`### 每一條規則`、`### 7.`、`### 8.`、`### 10.`、`### 狀態 tag`、`### 型別 tag`、`### 方言陷阱`。所以 Task 3 收工時 `grep -c "^#\{2,4\} " skills/bdd-spec/SKILL.md` 應該是 **32 − 14 ＝ 18**。

---

## Rulings（規格沒有裁定、本計畫替它裁定的事）

**Ruling A — `## 詞彙表`（`docs/CONTEXT.md`）留在 `bdd-spec`。**
D8 的表格沒有提到它。現行 `SKILL.md` 第 90 行給的理由是「寫 `.feature` 時措辭必須一致，所以詞彙在這一步才變成承重的東西」——照字面讀會把它拉去 `bdd-formulation`。但詞彙的**來源**是 `clarify-log.md` 的「已答」表，而 `bdd-formulation` 只讀 `spec.md`，讀不到 `clarify-log.md`。把產出留在讀得到來源的那一邊，讓 `bdd-formulation` 改成**讀** `docs/CONTEXT.md`。理由句要跟著改寫成「下一步寫 `.feature` 時措辭必須一致」。
**錯了的代價：** `docs/CONTEXT.md` 會晚一步才寫，但內容與格式完全一樣；要翻案只需把一節搬過去。

**Ruling B — `bdd-formulation` 要讀 `docs/CONTEXT.md`，這不算新格式。**
若不加這條，`.feature` 的措辭就沒有詞彙來源，封閉文法的重用會從第一份就開始漂移。它只是**新增一個輸入**，沒有改動任何產物的格式，所以不違反「格式完全不動」。

**Ruling C — 修掉「十二件事 vs 十三件事」的落差。**
`check_spec.py` 的 docstring 列了 13 項檢查（第 13 項是「無主的 FR」），`SKILL.md` 第 508 行起的稽核段寫「檢查十二件事」並且只列了 12 項。這段文字要被搬進一個全新的檔案；把一句已知是錯的話搬進沒有歷史的地方，等於讓它**重新變成新的**。所以在 Task 2 補上第 13 項並把數字改成十三，**獨立成一筆提交**以便單獨翻案。
同一份 docstring 第 38 行「不算在上面十二項一致性檢查裡」也是同一個落差，但那在 `check_spec.py` 裡——**不動**，因為 Global Constraints 禁止改腳本內容。改 `SKILL.md` 而不改腳本，落差從「兩處都錯」變成「腳本的 docstring 自相矛盾」，這件事寫進 Task 4 的收工回報，留給後續計畫。
**錯了的代價：** 若使用者其實想要一字不改的搬移，revert 那一筆提交即可。

**Ruling D — 純搬移那一筆提交會讓工作樹暫時壞掉，這是可以的。**
Task 1 結束時 `skills/bdd-formulation/` 沒有 `SKILL.md`（S1 MUST 失敗），而 `bdd-spec/SKILL.md` 指向一批已經不在的檔案（S5 MUST 失敗）。這是刻意的：rename 的完整性只有在那一筆提交不含任何內容變更時才成立。Task 1 的驗收標準是 **blob SHA 相同**，不是稽核通過。Task 2 收尾時稽核就會回綠。

---

## Task 1: 十八個檔案純搬移

**Files:**
- Move（`git mv`，內容零變更）：上面「搬過去的十八個檔案」表列的 18 個
- Create（由 `git mv` 自動建）：`skills/bdd-formulation/references/`、`skills/bdd-formulation/examples/`、`skills/bdd-formulation/scripts/`

**Interfaces:**
- Consumes: 無（第一個任務）
- Produces: 18 個檔案的新路徑，Task 2 的 `SKILL.md` 要指向它們；`skills/bdd-formulation/scripts/check_spec.py` 是 Task 2、Task 4 要跑的腳本

- [ ] **Step 1: 記下搬移前的 blob SHA**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
git rev-parse --verify HEAD > /tmp/base-commit.txt
for f in references/step-grammar.md references/state-tags.md references/rule-taxonomy.md \
         references/coverage-report.md references/artifact-location.md \
         examples/1-video-progress.md examples/2-submit-assignment.md examples/3-deliver-course.md \
         examples/4-query-progress.md examples/5-query-course-list.md examples/6-query-product-list.md \
         examples/7-create-order.md examples/8-query-order.md examples/9-pay-order.md \
         examples/10-cancel-order.md examples/11-query-user-roles.md examples/anti-patterns.md \
         scripts/check_spec.py; do
  printf "%s %s\n" "$(git rev-parse "HEAD:skills/bdd-spec/$f")" "$f"
done | sort > /tmp/before-shas.txt
wc -l /tmp/before-shas.txt
```

Expected: `18 /tmp/before-shas.txt`

- [ ] **Step 2: 建目錄並搬移**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
mkdir -p skills/bdd-formulation/references skills/bdd-formulation/examples skills/bdd-formulation/scripts
for f in references/step-grammar.md references/state-tags.md references/rule-taxonomy.md \
         references/coverage-report.md references/artifact-location.md \
         examples/1-video-progress.md examples/2-submit-assignment.md examples/3-deliver-course.md \
         examples/4-query-progress.md examples/5-query-course-list.md examples/6-query-product-list.md \
         examples/7-create-order.md examples/8-query-order.md examples/9-pay-order.md \
         examples/10-cancel-order.md examples/11-query-user-roles.md examples/anti-patterns.md \
         scripts/check_spec.py; do
  git mv "skills/bdd-spec/$f" "skills/bdd-formulation/$f"
done
git status --porcelain
```

Expected: 18 行 `R  skills/bdd-spec/... -> skills/bdd-formulation/...`，沒有任何 `M`。

- [ ] **Step 3: 確認工作樹裡零內容變更**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
git diff --cached --numstat | awk '$1 != 0 || $2 != 0'
```

Expected: **沒有任何輸出。** 有輸出就代表某個檔被改到了，停下來還原再重做。

- [ ] **Step 4: 提交**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
git commit -F - <<'MSG'
refactor(bdd-formulation): move the Gherkin half of bdd-spec, verbatim

WHAT: git mv five references, twelve examples and check_spec.py from
skills/bdd-spec/ into a new skills/bdd-formulation/. Not one byte of any
file changed.

WHY: bdd-spec produces two unrelated things -- a spec.md written from the
clarification log, and .feature files written from that spec.md. Splitting
the second half into its own skill is D8 of the grilling-and-formulation
design. This commit carries only the move so that git can infer the renames
from identical blobs; bundling the move with the rewrite makes git log
--follow break silently, with no error and no warning.

HOW: The tree is deliberately broken after this commit -- bdd-formulation
has no SKILL.md yet and bdd-spec/SKILL.md still points at files that left.
The next two commits close both halves. Verified by comparing blob SHAs
before and after, not by git diff on the new path: git cannot see a rename
when only the new path is given, so that diff reports the whole file as
added and proves nothing.
Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
MSG
git log --oneline -1
```

- [ ] **Step 5: 用 blob SHA 證明搬移是無損的**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
for f in references/step-grammar.md references/state-tags.md references/rule-taxonomy.md \
         references/coverage-report.md references/artifact-location.md \
         examples/1-video-progress.md examples/2-submit-assignment.md examples/3-deliver-course.md \
         examples/4-query-progress.md examples/5-query-course-list.md examples/6-query-product-list.md \
         examples/7-create-order.md examples/8-query-order.md examples/9-pay-order.md \
         examples/10-cancel-order.md examples/11-query-user-roles.md examples/anti-patterns.md \
         scripts/check_spec.py; do
  printf "%s %s\n" "$(git rev-parse "HEAD:skills/bdd-formulation/$f")" "$f"
done | sort > /tmp/after-shas.txt
diff /tmp/before-shas.txt /tmp/after-shas.txt && echo "ALL 18 BLOBS IDENTICAL"
```

Expected: `ALL 18 BLOBS IDENTICAL`，`diff` 無輸出。

- [ ] **Step 6: 證明 `git log --follow` 還活著**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
git log --follow --oneline -- skills/bdd-formulation/references/step-grammar.md | wc -l
git log --follow --oneline -- skills/bdd-formulation/scripts/check_spec.py | wc -l
git show --stat -M HEAD | head -25
```

Expected: 兩個數字都 **> 1**（不是只有剛剛那一筆）；`git show --stat -M` 每一行都以 `skills/bdd-spec/... => skills/bdd-formulation/...` 的形式呈現，代表 git 認出了 rename。

---

## Task 2: 寫 `skills/bdd-formulation/SKILL.md`

這個任務**只讀不剪** `bdd-spec/SKILL.md`——它此刻還是完整的 635 行，是抄寫的來源。剪除留給 Task 3。

**Files:**
- Create: `skills/bdd-formulation/SKILL.md`
- Read-only 來源: `skills/bdd-spec/SKILL.md`（635 行版本）

**Interfaces:**
- Consumes: Task 1 搬好的 18 個檔案路徑
- Produces: `skills/bdd-formulation/SKILL.md`；Task 3 與 Task 4 會從別處指向 `bdd-formulation` 這個名字

- [ ] **Step 1: 先讓稽核紅燈，記下它現在在抱怨什麼**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-formulation; echo "exit=$?"
```

Expected: `[S1] MUST: 找不到 SKILL.md`，`exit=1`。這就是這個任務要轉綠的那盞燈。

- [ ] **Step 2: 寫 frontmatter 與開場**

檔案開頭原樣寫成下面這樣：

```markdown
---
name: bdd-formulation
description: >
  把定案的 `spec.md` 翻成可執行的 Gherkin `.feature`——FR 變成 `Rule:` 區塊、
  AC 變成 `Example:`、編號變成 tag 帶下去，步驟套一組封閉文法讓 step definition
  能重用，最後跑一次雙向覆蓋稽核。不訪談、不發明例子、不改 `spec.md`。
  BDD 三實踐裡的 Formulation。
  觸發詞：「寫成 feature 檔」「轉成 Gherkin」「產出驗收情境」「把例子寫成場景」
  「寫 .feature」「產生驗收測試規格」「這些規則的 scenario 長怎樣」
  「spec 寫好了下一步」「跟 PM 對答案」。
  English: turn spec.md into Gherkin, write the feature file, generate
  acceptance scenarios, convert rules and examples into scenarios,
  formulation step of BDD.
---

# FORMULATION — 把 `spec.md` 翻成可執行的規格

一份定案的 `spec.md` 進來，出去的是 `.feature`：**FR 變成 `Rule:` 區塊、
AC 變成 `Example:`、編號變成 tag**，步驟套一組封閉文法，外加一張覆蓋表。

不訪談。不發明新的例子。不寫 step definition。不決定情境跑在哪一層——seam 在
`bdd-spec` 就決定完了，寫在 `spec.md` 裡。

這一步也是**跟 PM 對答案**的那一步。`.feature` 就是對答案用的那一份，不另做第二
份給人讀的版本——那是 Cucumber 的前提。PM 讀不動就修 `references/step-grammar.md`。
```

- [ ] **Step 3: 寫 `## 使用時機`、`## Skill Boundaries`、`## 參考檔案`、`## 前置確認`**

接在開場之後，原樣寫成：

```markdown
## 使用時機

- `spec.md` 已經寫好，要把 FR 與 AC 變成 `.feature`
- 要跟 PM 或需求方對答案——攤開規則與例子請他點頭
- 準備進 IMPLEMENT，需要先有紅燈的情境

## Skill Boundaries

- `spec.md` 還不存在、或還有沒答完的問題 → 回 `bdd-spec`（更上游是 `bdd-discovery`）
- 要切 story、決定 seam、寫 `spec.md` 或 `docs/CONTEXT.md` → 那是 `bdd-spec` 的工作
- 要稽核既有 `.feature` 寫得好不好 → 改用 `bdd-spec-review`
- 要決定情境跑在哪一層測試、實作順序、切票 → 改用 `bdd-plan`
- 產出是 `.feature`，**不是** `spec.md`、`openapi.yaml`、migration 或任何可執行的檔案

## 參考檔案

- `references/step-grammar.md` — 四種步驟形狀的完整理由與邊界
- `references/rule-taxonomy.md` — 四象限分類與覆蓋盤點
- `references/state-tags.md` — 狀態 tag 的完整詞彙表
- `references/coverage-report.md` — 覆蓋表的格式
- `references/artifact-location.md` — `.feature` 該放哪裡
- `examples/anti-patterns.md` — 寫完自檢
- `scripts/check_spec.py` — `.feature` ↔ `spec.md` 的一致性稽核

十一份逐案範例的索引在下面「這份 `.feature` 是寫給誰讀的」。

## 前置確認（Ask user）

不必問。`.feature` 的位置偵測得到（見「產物隔離」），要跑哪些 story 有預設值。

MUST NOT: 在這裡問問題。**不訪談**——憑空長出來的內容是缺陷。`spec.md` 的
`## Scope — In / Out` 裡「回 CLARIFY 補問」還有東西時，那些對應的規則就寫不
出來，不要順手決定掉。

NEVER: `spec.md` 不存在時，憑需求描述直接寫 `.feature`。那會產出一份沒有人同意
過的規格，而它看起來跟真的一模一樣。沒有 `spec.md` 就先跑 `bdd-spec`。
```

- [ ] **Step 4: 抄「這份 `.feature` 是寫給誰讀的」**

把 `skills/bdd-spec/SKILL.md` 的 **117–148 行**原樣抄過來，**一字不改**。那張表裡的 `references/step-grammar.md`、`references/rule-taxonomy.md`、`examples/anti-patterns.md`、`examples/1-` 到 `examples/11-` 在新 skill 底下是完全相同的相對路徑，所以不需要動——這張表同時也是 S5 指路鏈的主幹，抄漏一列就會多一個孤兒檔。

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -n '117,148p' skills/bdd-spec/SKILL.md
```

- [ ] **Step 5: 寫 `## 流程`，六個步驟**

`### 1. 讀三樣東西` 是新寫的，原樣寫成：

```markdown
## 流程

### 1. 讀三樣東西

| 讀什麼 | 為了什麼 |
| --- | --- |
| `specs/<date>-<feature>/spec.md` | **唯一的規範來源**——`## User Stories` 底下的 FR 與 AC 是這一步要翻譯的全部內容，`## Personas` 是 `Feature:` 敘述裡角色的來源，`## Scope — In / Out` 劃出碰不得的界線 |
| `docs/CONTEXT.md` | **措辭的來源**——同一個概念在不同 `.feature` 裡換一個說法，封閉文法的重用就從第一份開始漂移 |
| 既有的 `.feature` | 已經定下的步驟樣板，能套就不要另造 |

第三項最容易被跳過，代價最貴，而且**在封閉文法下更貴**：文法的價值全在重用，
另造一個形狀等於白做。寫新樣板之前先 grep 一次既有的 `.feature`。

IMPORTANT: `spec.md` 由 `bdd-spec` 維護，本 skill **不得寫入**。要補一條規則、
改一個例子、重排一個編號，都回 `bdd-spec`。

MUST: 只做 `spec.md` 的 `## Document Overview` 裡 `狀態` 已經是 `待對焦` 或
`已對焦` 的 feature。還停在 `澄清中` 的整份跳過，並明講該回哪個 skill——
把還沒定案的東西寫成場景，等於**把不確定性從一個顯眼的問題檔搬進一份看起來
已完成的規格**，之後沒有人會回頭質疑它。

使用者要求「未就緒的也先寫」時就寫，標成 `@draft`——見下方「狀態 tag」。
```

接著把 `bdd-spec/SKILL.md` 的下列行段原樣抄進來，只改標題編號：

| 抄哪幾行 | 新標題 |
| --- | --- |
| 263–314 | `### 2. 規則對應 Rule，例子對應 Example` |
| 315–334 | `#### 規則歸到四個象限`（不變） |
| 335–355 | `#### 每條規則要看三種結果`（不變） |
| 356–368 | `#### 場景名字要能定位`（不變） |
| 369–398 | `### 規則與它自己的例子打架時`（不變） |
| 399–430 | `### 3. 步驟文法 —— 四種形狀` |
| 431–459 | `### 每一條規則`（不變） |
| 460–487 | `### 4. 什麼時候用 Scenario Outline` |
| 488–494 | `### 5. 寫檔` |
| 508–546 | `### 6. 稽核，然後才算完成` |

抄完之後只有**一段**內文要改，其餘一字不動：

原 369–398「規則與它自己的例子打架時」那一節寫的是本 skill 不得回頭改上游。搬過來之後上游多了一層，所以：

- 原文 `MUST: 同時寫進 \`spec.md\` 的 \`## Scope — In / Out\` → 回 CLARIFY 補問，標明來源要更正。`
  → 改成 `MUST: 同時回報給 \`bdd-spec\`，由它寫進 \`spec.md\` 的 \`## Scope — In / Out\` → 回 CLARIFY 補問，標明來源要更正。`
- 原文 `MUST NOT: 直接改原始輸入，也不得寫入 \`clarify-log.md\`——那是 CLARIFY 的產物，本 skill 不得寫入。`
  → 改成 `MUST NOT: 直接改原始輸入，也不得寫入 \`spec.md\` 或 \`clarify-log.md\`——那是 \`bdd-spec\` 與 CLARIFY 的產物，本 skill 不得寫入。`

「步驟 N」這種內文回指在 263–546 這個範圍裡**一次都沒有出現**（整份 `SKILL.md` 只有第 570 行有一處「完整骨架見步驟 6 的範例」，那在 `## 產物格式` 裡，已經在下一個 Step 改成「步驟 3」）。抄完跑一次確認：

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "步驟 [0-9]" skills/bdd-formulation/SKILL.md
```

Expected: 只有 `## 產物格式` 那一行的「完整骨架見步驟 3 的範例」。

- [ ] **Step 6: 抄產物格式、狀態 tag、型別 tag、方言陷阱、產物隔離、完成後**

`## 產物格式` 原樣寫成（這是從 `bdd-spec` 的 547–574 行**拆**出來的 `.feature` 那一半）：

````markdown
---

## 產物格式

```
features/<story-slug>.feature     FORMULATION  ★ 本步驟唯一的產物

（spec.md 是 bdd-spec 的產物，本 skill 只讀不寫）
```

`.feature` 的格式全在本節（關鍵字、狀態 tag、型別 tag、方言陷阱）與
`references/step-grammar.md`。

規則與例子的**定版編號在 `bdd-spec` 誕生**：`FR-<n>`、`AC-<n>.<m>`。`.feature` 的
`@rule-<n>`／`@example-<n>.<m>` 回指它們——本 skill 只把它們原封不動地搬進 tag。

**關鍵字用英文，名稱與步驟用中文。** 不加 `# language:` 那一行（英文是預設方言）。

完整骨架見步驟 3 的範例。`Feature:` 底下那段敘述寫的是這則 story 的
「作為⋯我要⋯以便⋯」——角色抄 `spec.md` 的 `## Personas`，能力與價值從這則
story 涵蓋的那幾條 FR 濃縮而來（`## User Stories` 裡的切法說明就是濃縮的
起點）。它是這份規格存在的理由，而讀 `.feature` 的人不會同時開著 `spec.md`。
````

然後把 `bdd-spec/SKILL.md` 的 **575–616 行**（狀態 tag、型別 tag、方言陷阱三節）原樣抄過來，**一字不改**。

`## 產物隔離` 原樣寫成：

```markdown
## 產物隔離

MUST: 只新增 `.feature`，**不修改專案既有的任何檔案**（含既有的 `.feature`、
`spec.md`、`docs/CONTEXT.md`）。`.feature` 的位置由測試框架決定——有 `.feature`
就跟隨，沒有就用專案根的 `features/`。理由與例外 → `references/artifact-location.md`。
```

`## 完成後` 原樣寫成：

```markdown
## 完成後

告訴對方四件事：

1. 產了哪些 `.feature`、各幾個場景、**幾個不重複步驟樣板**
2. **覆蓋表**——哪些例子還沒有場景，為什麼
3. 哪些象限是空的，空得合不合理
4. `spec.md` 有沒有需要更正的地方（規則與例子打架、缺了角色）——本 skill 不改它，
   要回 `bdd-spec`

這些場景現在應該**全部是紅的**（step definition 還不存在）。這是對的：
綠燈要等 IMPLEMENT。一跑就綠代表這些場景沒有驗到任何東西。

覆蓋表沒有大洞 → `bdd-plan` 切票。洞多 → 回 `bdd-spec`，更上游是 `bdd-discovery`。
```

- [ ] **Step 7: 跑稽核，要轉綠**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-formulation; echo "exit=$?"
```

Expected: `exit=0`，而且**沒有任何 S5**（S5 會逐一點名哪個檔沒被指到）。

本文行數不必自己算——**`audit_skill.py` 只在超過 500 行時才印 S10**，所以上面那條
指令沒有印出 S10 就代表過了。印出來的話，把「這份 `.feature` 是寫給誰讀的」的表格
下放到 `references/`——但先回報，不要自己決定。

- [ ] **Step 8: 逐一確認 18 個檔案都被指到**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
for f in $(cd skills/bdd-formulation && find references examples scripts -type f | sort); do
  grep -q "$f" skills/bdd-formulation/SKILL.md && echo "ok   $f" || echo "MISS $f"
done
```

Expected: 18 行全部 `ok`。`MISS` 代表 S5 會炸——即使 Step 7 因為遞迴指路鏈而過關，直接指路仍然比較穩。

- [ ] **Step 9: 用 fixture 證明 `check_spec.py` 從新家跑出來一模一樣**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FX="$(mktemp -d)/fx"
mkdir -p "$FX/specs/2026-09-09-parking" "$FX/features"
cp skills/bdd-spec/examples/minimal-spec.md "$FX/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-formulation/scripts/check_spec.py "$FX" | sed "s|$FX|FX|g"
echo "exit=${PIPESTATUS[0]}"
```

Expected（`4 個問題`、`exit=1`，兩則 story 各兩條）：

```
specs: FX/specs    feature: FX/features

✗ monthly-pass                     0/1 例子 · ?      
    ✗ 狀態 tag 缺 —— 每個檔恰好要一個
    ✗ 漏了 ['2.1'] 且檔案裡沒有註解說明
✗ visitor-billing                  0/3 例子 · ?      
    ✗ 狀態 tag 缺 —— 每個檔恰好要一個
    ✗ 漏了 ['1.1', '1.2', '1.3'] 且檔案裡沒有註解說明

4 個問題
```

- [ ] **Step 10: 提交**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
git add skills/bdd-formulation/SKILL.md
git commit -F - <<'MSG'
feat(bdd-formulation): give the moved half its own SKILL.md

WHAT: Add skills/bdd-formulation/SKILL.md -- frontmatter, 使用時機,
Skill Boundaries, 參考檔案, 前置確認, a six-step 流程, and the .feature
artifact format. The Gherkin sections are copied verbatim from
bdd-spec/SKILL.md, which is still intact at this commit; only step numbers,
the skill names in cross-references, and the input list changed.

WHY: After the previous commit the eighteen moved files had no SKILL.md
pointing at them, which is an S5 MUST failure -- a file nothing points at is
a file the model never reads. This commit closes that half. bdd-spec keeps
its copy of the same text until the next commit, so the text never exists in
zero places.

HOW: New input list is spec.md + docs/CONTEXT.md + existing .feature files;
bdd-spec keeps ownership of spec.md and docs/CONTEXT.md because their source
is the clarification log, which this skill does not read. Verified with
audit_skill.py (exit 0, no S5), a per-file pointer check over all eighteen
files, and check_spec.py run from its new path against a fixture, producing
byte-identical output to the pre-move baseline.
Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
MSG
```

- [ ] **Step 11: 補上稽核段的第 13 項（Ruling C），獨立提交**

`### 6. 稽核，然後才算完成` 那段裡，把「檢查十二件事」改成「檢查十三件事」，並在清單末尾 `AC-<n>.<m>` 那一項之後補進：

```
、`## User Stories` 底下有沒有 FR 不屬於任何一則 story（沒有 story 認領它，
就沒有任何 `.feature` 被要求帶它的 `@example-` tag，那些 AC 於是安靜地不見了）
```

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "十三件事\|無主\|不屬於任何一則 story" skills/bdd-formulation/SKILL.md
git add skills/bdd-formulation/SKILL.md
git commit -F - <<'MSG'
docs(bdd-formulation): the audit section listed 12 checks, the script has 13

WHAT: Change 「檢查十二件事」 to 十三件事 and add the missing item -- an FR
that belongs to no story, which check_spec.py has checked as item 13 all
along.

WHY: The sentence was already wrong in bdd-spec/SKILL.md. It was about to be
carried into a brand-new file, where a claim with no history reads as freshly
verified. An orphan FR is exactly the defect this repo keeps producing:
every line is well-formed, nothing errors, and the ACs simply stop being
covered.

HOW: Kept as its own commit so it can be reverted without touching the split.
check_spec.py's own docstring carries the same drift at its closing paragraph
(「不算在上面十二項一致性檢查裡」) but is left alone -- this plan forbids
editing the scripts. Reported as a follow-up.
Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
MSG
```

---

## Task 3: 把 Gherkin 從 `skills/bdd-spec/SKILL.md` 剪掉

**Files:**
- Modify: `skills/bdd-spec/SKILL.md`（635 行 → 約 300 行）

**Interfaces:**
- Consumes: Task 2 建好的 `skills/bdd-formulation/SKILL.md`（`bdd-spec` 要指向這個名字）
- Produces: 一份只講 `spec.md` 的 `bdd-spec`；Task 4 的 repo 掃描以此為準

- [ ] **Step 1: 先讓稽核紅燈，記下它現在在抱怨什麼**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
```

Expected: `exit=0`（S5 只檢查**存在的檔案有沒有被指到**，不檢查**指到的檔案存不存在**，所以指向已搬走的檔案不會被抓出來——這正是這一步必須靠人做 grep 的理由）。同時應該還看得到那項 S10（本文 622 行）。

- [ ] **Step 2: 剪掉整段搬走的十四個標題所轄範圍**

下面那條 `sed` 的六個 `-e` **全部以原始 635 行的行號為準**——sed 一次掃過輸入，範圍指的是輸入行號，不是刪到一半的行號，所以順序無所謂。改用編輯器手動刪的話才需要從檔尾往前。已實跑驗證：刪完 290 行、18 個標題。

| 刪哪幾行 |
| --- |
| 604–616（方言陷阱） |
| 587–603（型別 tag） |
| 575–586（狀態 tag） |
| 508–546（`### 10. 稽核`） |
| 263–494（`### 5.` 到 `### 8.`，含 `#### 規則歸到四個象限`、`#### 每條規則要看三種結果`、`#### 場景名字要能定位`、`### 規則與它自己的例子打架時`、`### 6.`、`### 每一條規則`、`### 7.`）——**注意 495–507 的 `### 9. 增修 specs/domain-model.md` 要留下** |
| 117–148（`## 這份 .feature 是寫給誰讀的`） |

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -i -e '604,616d' -e '587,603d' -e '575,586d' -e '508,546d' -e '263,494d' -e '117,148d' skills/bdd-spec/SKILL.md
grep -c "^#\{2,4\} " skills/bdd-spec/SKILL.md
```

Expected: **18**（32 − 14；其中 2 個仍是 code fence 裡的 `### US-1` / `### US-2`）。

- [ ] **Step 3: 把 `### 9.` 重新編號成 `### 5.`**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -i 's|^### 9\. 增修 `specs/domain-model.md`$|### 5. 增修 `specs/domain-model.md`|' skills/bdd-spec/SKILL.md
grep -n "^### [0-9]\." skills/bdd-spec/SKILL.md
```

Expected 恰好五行，依序是 `### 1. 挑 story`、`### 2. 讀三樣東西`、`### 3. 決定 seam`、`### 4. 寫 \`spec.md\``、`### 5. 增修 \`specs/domain-model.md\``。

- [ ] **Step 4: 改寫 frontmatter 與開場**

`description` 與標題底下的開場段換成：

```markdown
---
name: bdd-spec
description: >
  把 CLARIFY 的問答綜合成 `spec.md`——十三節、完整自足的規格文件：切 story、
  抽出 FR 與 AC、決定驗收測試打在哪一層 seam、把詞彙抬進 `docs/CONTEXT.md`。
  不訪談，只綜合已經有答案的東西；`.feature` 由下一步的 `bdd-formulation` 產出。
  BDD 六步流程的 SPEC 步驟。
  觸發詞：「產出規格文件」「寫 spec.md」「clarify-log.md 變成規格」
  「把問答整理成規格」「切 story」「決定 seam」「驗收測試打在哪一層」
  「這個 feature 的規格長怎樣」。
  English: synthesize the clarification log into spec.md, turn answered
  questions into a spec, slice stories, decide the acceptance test seam,
  write the glossary.
---

# SPEC — 把答案綜合成一份規格

一份收斂的 `clarify-log.md`（連同它問過的原始輸入）進來，出去的是
`spec.md`：**先切 story——決定哪幾條 FR 湊成一則可獨立驗收的 story——再把
已經有答案的東西綜合成 FR 與 AC**，並決定驗收測試打在哪一層 seam。

不訪談。不發明新的例子。**不產 `.feature`**——那是下一步 `bdd-formulation` 的事。
```

- [ ] **Step 5: 改寫 `## 使用時機`、`## Skill Boundaries`、`## 前置確認`、`## 產物格式`、`## 產物隔離`、`## 完成後`**

`## 使用時機` 三條改成：

```markdown
- CLARIFY 判定收斂，要把 `clarify-log.md` 變成 `spec.md`
- 手上有規則與具體例子，需要一份完整自足的規格文件
- 準備跟 PM 對答案，但還沒有 `spec.md`
```

`## Skill Boundaries` 改成：

```markdown
- **要把 FR 與 AC 寫成 `.feature` → 改用 `bdd-formulation`**
- 要稽核既有 `.feature` 寫得好不好 → 改用 `bdd-spec-review`
- 要決定情境跑在哪一層測試、實作順序 → 改用 `bdd-plan`
- 例子還不夠、還有紅卡 → 回 `bdd-discovery`
- 產出是 `specs/<date>-<feature>/spec.md` ＋ `specs/domain-model.md` ＋ `docs/CONTEXT.md` 的追加，**不是** `.feature`、`openapi.yaml`、migration 或任何可執行的檔案
- 要把 `spec.md` 拆成可執行的票 → 改用 `bdd-plan`
```

`## 前置確認（Ask user）` 的第一段（原 151 行「不必問。`.feature` 的位置偵測得到⋯」）換成：

```markdown
不必問輸入的位置。`clarify-log.md` 與 `input.md` 都在 `specs/<date>-<feature>/`
底下。唯一要徵詢的是 seam，在流程步驟 3。
```

下面的 MUST NOT 與兩條 NEVER **原樣留著**，只把最後一條 NEVER 裡「每一個 `← PRD §x` 都指向一份從未被寫下的文件」那段保留不動。

`## 產物格式` 的區塊圖與表格改成：

````markdown
## 產物格式

```
specs/<date>-<feature>/
└── spec.md           SPEC  ★ 十三節，完整自足

（clarify-log.md 是 CLARIFY 的產物，本 skill 只讀不寫）
（.feature 是 bdd-formulation 的產物，本 skill 不產）
```

本步驟產三份東西，格式各有一份參考檔：

| 產物 | 格式 |
| --- | --- |
| `specs/<date>-<feature>/spec.md` | `references/spec-format.md` |
| `specs/domain-model.md`（增修） | `references/spec-format.md` 末節 |
| `docs/CONTEXT.md`（建立或追加） | 「詞彙表」一節 |

規則與例子的**定版編號在這裡誕生**：`FR-<n>`、`AC-<n>.<m>`。`bdd-formulation` 寫的
`@rule-<n>`／`@example-<n>.<m>` 回指它們——所以**這批編號一旦寫定就不准重排**。
下游回指的是 tag，不是那句規則的文字本身；重排等於讓那些引用**靜默**指向別的
東西，不會報錯，只會對錯。
````

`## 產物隔離` 改成：

```markdown
## 產物隔離

MUST: 只新增本節列出的三種產物，**不修改專案既有的任何檔案**（含既有的
`.feature`、`docs/CONTEXT.md` 既有段落）。`.feature` 完全不碰——本 skill 不產它，
也不改它。
```

`## 完成後` 改成：

```markdown
## 完成後

告訴對方四件事：

1. 切了哪些 story、為什麼這樣切（`## User Stories` 的切法說明）
2. **seam 是哪一層**、`domain-model.md` 這一批新增了哪些聚合與不變條件、
   `docs/CONTEXT.md` 追加了哪些詞彙
3. `## Open Questions` 還剩幾題待答、`## Scope — In / Out` 的「回 CLARIFY 補問」
   有幾條
4. `## Document Overview` 的 `狀態` 這一步寫成什麼（`澄清中` 還是 `待對焦`）

缺口多 → 回 `clarify-loop` 收斂。缺口少 → 跑 `bdd-formulation` 產 `.feature`，
拿它跟 PM 對答案。
```

- [ ] **Step 6: 改掉三處指向搬走的東西的內文**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "check_spec\|既有的 \`\.feature\`\|寫 \`\.feature\` 時措辭" skills/bdd-spec/SKILL.md
```

逐一改成：

1. 「切 story」裡 `kebab-case，一份 `spec.md` 裡唯一——`check_spec.py` 用它對應檔案。」→ 把 `check_spec.py` 改成 ``bdd-formulation` 的 `check_spec.py``。
2. 「切 story」裡 ``scripts/check_spec.py` 用它算「這則 story 該有哪些 `@example` tag」。」→ 改成 ``bdd-formulation` 的 `check_spec.py` 用它算⋯」。**這一處務必把 `scripts/` 前綴拿掉**——留著會讓 S5 的正則把 `scripts/check_spec.py` 收進 `bdd-spec` 的 `mentioned` 集合，雖然不會報錯，但那是一條指向不存在檔案的路。
3. 「詞彙表」第一句「寫 `.feature` 時措辭必須一致，所以詞彙在這一步才變成承重的東西。」→ 改成「**下一步**寫 `.feature` 時措辭必須一致，所以詞彙在這一步就要定下來。」（Ruling A）
4. `### 2. 讀三樣東西` 的表格刪掉「既有的 `.feature`」那一列，標題改成 `### 2. 讀兩樣東西`，並刪掉表格底下那一段「第三項最容易被跳過⋯先 grep 一次既有的 `.feature`。」（那整段搬去 `bdd-formulation` 了）。

- [ ] **Step 7: 負向檢查——`bdd-spec` 不該再提到任何搬走的東西**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "references/step-grammar\|references/state-tags\|references/rule-taxonomy\|references/coverage-report\|references/artifact-location\|scripts/check_spec\|examples/[0-9]\|examples/anti-patterns" skills/bdd-spec/SKILL.md
```

Expected: **沒有任何輸出。**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -n "Scenario Outline\|Rule:\|@rule-\|@example-\|@ready\|@draft\|@command\|@query\|gherkin\|Gherkin" skills/bdd-spec/SKILL.md
```

Expected: 只剩下**接縫說明**那幾行——`@rule-<n>`／`@example-<n>.<m>` 出現在 `## 產物格式` 的編號段與「切 story」的 `@example` tag 那句。**不該**再有 `Scenario Outline`、`@ready`、`@draft`、`@command`、`@query`、`Rule:`。逐行看過，說得出每一行為什麼還在。

- [ ] **Step 8: 跑稽核與行數**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-formulation; echo "exit=$?"
wc -l skills/bdd-spec/SKILL.md skills/bdd-formulation/SKILL.md
```

Expected: 兩個都 `exit=0`，而且**兩份都不再印出 S10** —— `bdd-spec` 原本有一項
S10（本文 622 行），剪完應該消失。`bdd-spec/SKILL.md` 整份約 300 行（純刪除之後
實測 290 行，步驟 4–6 的改寫會再加回十幾行）。

- [ ] **Step 9: 用 fixture 證明 `status.py` 沒被動到**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
FX="$(mktemp -d)/fx"
mkdir -p "$FX/specs/2026-09-09-parking"
cp skills/bdd-spec/examples/minimal-spec.md "$FX/specs/2026-09-09-parking/spec.md"
python3 skills/bdd-spec/scripts/status.py "$FX" | grep "^合計"
echo "exit=${PIPESTATUS[0]}"
git diff --stat HEAD -- skills/bdd-spec/scripts/ skills/bdd-spec/references/ skills/bdd-spec/examples/
```

Expected: `合計` 那行是 `2     1     1`，`exit=0`；最後一個 `git diff --stat` **沒有輸出**（Global Constraints：格式檔與腳本一個位元組都不准動）。

- [ ] **Step 10: 提交**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
git add skills/bdd-spec/SKILL.md
git commit -F - <<'MSG'
refactor(bdd-spec): stop producing Gherkin

WHAT: Remove the fourteen Gherkin sections from bdd-spec/SKILL.md -- the
examples index, steps 5 through 8, step 10, the state/type tag sections and
the dialect trap -- and rewrite the frontmatter, 使用時機, Skill Boundaries,
前置確認, 產物格式, 產物隔離 and 完成後 around the three artifacts that are
left. 635 lines down to roughly 300.

WHY: A skill that produces both a spec.md and the .feature files derived from
it has two jobs, and the second one is where every step-grammar decision
lives. The previous commit gave that half its own home; leaving a duplicate
here would drift -- the same table in two files is the failure mode this repo
already named.

HOW: Sections were copied into bdd-formulation first and cut here second, so
the text was never absent from the tree. 詞彙表 (docs/CONTEXT.md) stays here
because its source is the clarification log, which bdd-formulation does not
read; the justification sentence now says the wording matters for the NEXT
step. Step 9 renumbered to 5. Verified with a heading count reconciliation
(32 minus 14 equals 18), a negative grep for every moved path and every
Gherkin keyword, audit_skill.py exit 0 on both skills with S10 now gone, and
git diff proving no reference, example or script file was touched.
Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
MSG
```

---

## Task 4: 掃掉 repo 其餘還指向舊位置的地方

**Files:**
- Modify: `CLAUDE.md:62`
- Modify: `README.md:30` 附近的目錄樹
- Modify: `PLAN.md:200`（SPEC 那一列）
- Modify: `skills/bdd-discovery/SKILL.md:40`
- Modify: `skills/clarify-loop/SKILL.md:30`
- Modify: `skills/example-mapping/SKILL.md:48`
- Modify: `skills/story-splitting/SKILL.md:197`
- Modify: `skills/bdd-plan/SKILL.md:28`
- Modify: `benchmark/skeleton/go/README.md:149`

**Interfaces:**
- Consumes: Task 2、Task 3 定下的 skill 名稱 `bdd-formulation` 與檔案新路徑
- Produces: 無下游任務

- [ ] **Step 1: 列出所有還需要改的地方**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -rn "bdd-spec" CLAUDE.md README.md PLAN.md docs/ benchmark/ skills/ 2>/dev/null \
  | grep -v "^docs/superpowers/" | grep -v "^skills/bdd-spec/"
```

這份輸出就是這一步的工作清單。逐行判斷：講的是**綜合 `spec.md`** 的就保持 `bdd-spec`，講的是**寫 Gherkin／`.feature`／`check_spec.py`** 的改成 `bdd-formulation`。

- [ ] **Step 2: 改 `CLAUDE.md`**

第 61–62 行那個區塊改成：

````markdown
```bash
python3 skills/bdd-spec/scripts/status.py [root]                # clarification progress per feature
python3 skills/bdd-formulation/scripts/check_spec.py [root]     # .feature ↔ spec.md consistency
```
````

同一份檔案的「The three disciplines」那一節底下，在 **PLAN** 那一條之前補一條：

```markdown
- **FORMULATION** turns the settled `spec.md` into `.feature` files — it does not
  interview, does not invent examples, and does not write back to `spec.md`.
```

並把 **SPEC** 那一條的「**SPEC** slices stories and synthesises what already has answers」保持不動（它本來就沒提 Gherkin）。

Commands 那一節的開頭句「Both take a project root and default to `.`」保持不動——兩支腳本的行為沒變，只是住在不同 skill 底下。

- [ ] **Step 3: 改 `README.md` 的目錄樹**

在 `│   ├── bdd-spec/` 之後補一行 `│   ├── bdd-formulation/`（縮排與相鄰行對齊；若該樹有一行說明就一併補上，措辭照相鄰行的風格）。

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -n '25,40p' README.md
```

- [ ] **Step 4: 改 `PLAN.md`**

第 200 行 SPEC 那一列，把產物欄的 `答案 → \`.feature\` ＋ \`spec.md\`` 改成 `答案 → \`spec.md\``，並在它下面補一列：

```markdown
| FORMULATION | `bdd-formulation` | `spec.md` → `.feature` | **已實作**：從 `bdd-spec` 拆出來，封閉步驟文法與 `check_spec.py` 跟著走 |
```

第 212 行「`bdd-spec` 加了 `spec.md` 之後會不會太肥、要不要拆成兩個 skill——**先不拆**。」改寫成記錄這個決定已經翻案：保留原句，後面接一句「**2026-09-14 翻案：拆了。** 見 `docs/superpowers/specs/2026-09-14-grilling-and-formulation-design.md` 的 D8。」

第 323 行 ``skills/bdd-spec/references/state-tags.md`` 改成 ``skills/bdd-formulation/references/state-tags.md``。

第 338、344 行的兩個 `- [x]` 條目講的是 `bdd-spec` 過去做完的事，是歷史紀錄，**不動**。

- [ ] **Step 5: 改四個 skill 的「要把例子寫成 Gherkin」那一行**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -i 's|- 要把例子寫成 Gherkin → 改用 `bdd-spec`|- 要把例子寫成 Gherkin → 改用 `bdd-formulation`|' \
  skills/bdd-discovery/SKILL.md skills/clarify-loop/SKILL.md skills/example-mapping/SKILL.md
grep -rn "要把例子寫成 Gherkin" skills/
```

Expected: 三行全部指向 `bdd-formulation`。

- [ ] **Step 6: 改 `story-splitting` 與 `bdd-plan`**

`skills/story-splitting/SKILL.md:197`「每則切完後回 `bdd-spec`，把各自的 FR 與例子寫成 `.feature`」→ 改成「每則切完後回 `bdd-spec` 寫進 `spec.md`，再由 `bdd-formulation` 把各自的 FR 與例子寫成 `.feature`」。

`skills/bdd-plan/SKILL.md:28`「還沒有 `spec.md` → 先跑 `bdd-spec`；規則還沒定案 → 回 `bdd-discovery`」→ 改成「還沒有 `spec.md` → 先跑 `bdd-spec`；有 `spec.md` 但還沒有 `.feature` → 先跑 `bdd-formulation`；規則還沒定案 → 回 `bdd-discovery`」。

`skills/bdd-plan/SKILL.md:19` 與 `:29` 講的是「API、domain 型別、schema、seam 是 `bdd-spec` 的工作」——seam 確實留在 `bdd-spec`，**不動**。

- [ ] **Step 7: 改 `benchmark/skeleton/go/README.md:149`**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
sed -n '143,155p' benchmark/skeleton/go/README.md
```

那一段講的是 `@ready` 這個狀態 tag 的規則出處。狀態 tag 現在歸 `bdd-formulation`，所以把該句的 `bdd-spec` 改成 `bdd-formulation`。**先讀上下文再改**——如果那句其實在講 story 就緒判定（那仍然是 `bdd-spec`），就不要改，並在回報裡說明為什麼不改。

- [ ] **Step 8: 全域負向掃描**

掃「同一行裡既提到 `bdd-spec`、又提到 Gherkin 相關的東西」。`skills/bdd-spec/`
與 `skills/bdd-formulation/` 兩個目錄排除在外——它們是 Task 3 與 Task 2 的施工面，
各自有自己的負向檢查；把它們掃進來會讓這一步的驗證範圍比它的 Files 清單寬，
而那正是這份計畫最想避免的錯。

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
grep -rn "bdd-spec" CLAUDE.md README.md PLAN.md docs/ benchmark/ skills/ 2>/dev/null \
  | grep -v "^docs/superpowers/" \
  | grep -v "^skills/bdd-spec/" | grep -v "^skills/bdd-formulation/" \
  | grep -i "gherkin\|\.feature\|check_spec\|state-tags\|step-grammar\|rule-taxonomy\|coverage-report\|artifact-location\|@ready\|@draft\|scenario"
```

Expected: **恰好兩行**，而且兩行都是刻意寫成「兩個 skill 一起出現」的交棒句：

| 檔案 | 為什麼還在 |
| --- | --- |
| `skills/bdd-plan/SKILL.md` | Step 6 改出來的「還沒有 `spec.md` → 先跑 `bdd-spec`；有 `spec.md` 但還沒有 `.feature` → 先跑 `bdd-formulation`」 |
| `skills/story-splitting/SKILL.md` | Step 6 改出來的「每則切完後回 `bdd-spec` 寫進 `spec.md`，再由 `bdd-formulation` 把各自的 FR 與例子寫成 `.feature`」 |

出現第三行就是還有一處把 Gherkin 掛在 `bdd-spec` 名下。**少於兩行**同樣是問題——
代表 Step 6 的兩處改寫有一處沒做到，或做成了只提一個 skill 的形式。

- [ ] **Step 9: 全套驗收**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for d in skills/*/; do printf "%-22s " "$(basename $d)"; \
  python3 skills/skill-rules/scripts/audit_skill.py "$d" >/dev/null 2>&1 && echo "exit 0" || echo "exit $?"; done
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
claude plugin validate . 2>&1 | tail -5
```

Expected: **8 個 skill 全部 `exit 0`**（原本 7 個，加上 `bdd-formulation`）；`claude plugin validate` 通過，**恰好 1 個警告**，就是 plugin root 那個 `CLAUDE.md`。

`__pycache__` 的清除不能省：稽核腳本會把 `scripts/__pycache__/` 底下的檔案當成沒有被指路的孤兒，報出假的 S5。

- [ ] **Step 10: 提交**

```bash
cd /home/benny/Workspace/vivotek/ai-bdd
git add -A
git commit -F - <<'MSG'
docs: point the rest of the repo at bdd-formulation

WHAT: Update CLAUDE.md's check_spec.py path and discipline list, README.md's
skill tree, PLAN.md's pipeline table and state-tags path, the "要把例子寫成
Gherkin" boundary line in bdd-discovery / clarify-loop / example-mapping,
the handoff sentences in story-splitting and bdd-plan, and the @ready
citation in the Go testbed README.

WHY: A pointer to a skill that no longer does the job is worse than no
pointer: it sends the reader somewhere real, which is why nobody checks it.
PLAN.md also recorded the opposite decision ("先不拆") and needed the
reversal written down beside it rather than silently edited away.

HOW: Split by what each line is actually about -- lines about synthesising
spec.md, slicing stories, or deciding the seam keep pointing at bdd-spec;
only Gherkin, .feature, tags and check_spec.py moved. Verified with a
repo-wide negative grep that pairs "bdd-spec" against every Gherkin term,
audit_skill.py exit 0 on all eight skills, and claude plugin validate passing
with exactly the one known warning.
Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
MSG
git log --oneline -6
```

---

## Self-Review

**1. 規格覆蓋**

| 規格要求（D8 與「計畫一」） | 哪個任務 |
| --- | --- |
| 五份 reference 搬過去 | Task 1 |
| 12 個 examples 搬過去 | Task 1 |
| `check_spec.py` 搬過去 | Task 1 |
| `spec-format.md`／`persona-definition.md`／`status.py`／`minimal-spec.md` 留下 | Task 1（不搬）＋ Task 3 Step 9（`git diff` 證明沒動到） |
| `bdd-spec` 不再產 `.feature` | Task 3 |
| 新 skill 讀 `spec.md`、產 `features/<slug>.feature` | Task 2 Step 5 |
| 接縫不變（`FR-<n>` ↔ `@rule-<n>`、`AC-<n>.<m>` ↔ `@example-<n>.<m>`） | Task 2 Step 5（原樣抄）＋ Task 3 Step 7（負向檢查只准留接縫說明） |
| 不做第二份給 PM 的版本 | Task 2 Step 2 開場段明講 |
| 格式完全不動 | Global Constraints ＋ Task 3 Step 9 的 `git diff --stat` |
| 做完流程完整可跑 | Task 4 Step 9 |

**未被任何任務覆蓋而我知道的事**（刻意留在範圍外，設計 spec 把它們編在計畫二、三）：
`spec.md` 加四節、`bdd-spec` 加 seam 草擬步驟、`bdd-discovery` 換前沿模型、`example-mapping` 改輸入。本計畫不碰。

**2. 佔位字掃描**

全文沒有 TBD／TODO／「之後再補」／「參照 Task N」。每一處要抄的文字都給了**來源行號或完整原文**；每一處要改的文字都給了**改前與改後**。唯一一處帶判斷的是 Task 4 Step 7（testbed README 那句要先讀上下文），而它明寫了「不改的話要在回報裡說明為什麼」，不是把判斷丟給執行者自己消化。

**3. 名稱一致性**

skill 名 `bdd-formulation` 在 frontmatter `name`、目錄名、與所有交叉引用裡拼法一致（`audit_skill.py` 的 F2 會抓 `name` 與目錄名不符）。步驟編號：`bdd-formulation` 是 1–6，`bdd-spec` 是 1–5，兩份各自連續，內文回指（「見步驟 3 的範例」）都跟著改了。fixture 的建法在 Task 2 Step 9 與 Task 3 Step 9 一字不差。

**4. 我自己在這份計畫裡最可能犯的錯**

上一個分支裡我 26 條裁決有 12 條是自己的計畫缺陷，全是同一種：**驗證範圍比施工範圍寬**——Files 清單只列了三個檔，卻要求執行者跑一個會掃到第四個檔的 grep。這份計畫的對策是每個任務的負向檢查都限定在**該任務 Files 清單列出的檔案**裡（Task 3 Step 7 只 grep `skills/bdd-spec/SKILL.md`；Task 4 Step 8 的全域 grep 排除了 `skills/bdd-spec/` 與 `docs/superpowers/`，而 Task 4 的 Files 清單正好列滿了剩下的每一個檔）。
