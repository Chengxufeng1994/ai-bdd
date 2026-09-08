# Benchmark Restructure Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 `lab/go/skeleton/` 拆成三個壽命不同的東西——`benchmark/cases/`（題目與評分方案）、`benchmark/skeleton/go/`（可複製的範本）、`runs/`（gitignored 的工作區）。

**Architecture:** 純目錄搬移，零程式碼變更。已確認 `skills/` 底下三支 Python 都沒有寫死 `lab/`、`skeleton` 或 `testbed`，所以這次只動文件與路徑。全部用 `git mv` 以保住 `git log --follow`。

**Tech Stack:** git、Go 1.27 ＋ godog（僅驗證用）、Python 3（連結檢查）

**Spec:** [`docs/superpowers/specs/2026-09-08-benchmark-restructure-design.md`](../specs/2026-09-08-benchmark-restructure-design.md)

## Global Constraints

- **全部搬移用 `git mv`**，不可 `cp` ＋ `rm`——`git log --follow` 是這次搬移的驗收條件之一。
- **不改 Go module 名。** `module skeleton` 留著；`Makefile` 的 `-X skeleton/pkg/version.value` 依賴它，`verify-stamp` 是守它的 gate。
- **不碰 `skills/` 底下任何檔案。**
- **`docs/bdd/` → `specs/` 的既有陳舊引用不在本次範圍**，唯一例外：搬進新檔的那一份要就地改對（Task 1 Step 6）。
- **散文用繁體中文，commit 訊息用英文**，WHAT/WHY/HOW 三段。
- 每個 task 結束時工作區乾淨、驗證全過。

## File Structure

| 檔案 | 責任 | Task |
| --- | --- | --- |
| `benchmark/README.md` | **新建。** case 的定義、餵食線、快照與污染兩條規則、怎麼跑一次 | 1 |
| `benchmark/cases/*.md` | 題目 ＋ 評分方案，語言無關 | 1 搬移、2 補模糊點 |
| `benchmark/skeleton/go/` | Go harness 範本 | 1 搬移 |
| `benchmark/skeleton/go/README.md` | 純 harness 說明 | 1 刪舊節、2 搬出情境節 |
| `.gitignore` | `runs/` 取代 `eval-runs/` | 1 |
| `README.md`／`PLAN.md`／`CLAUDE.md` | 路徑更新 | 1 |
| `benchmark/skeleton/go/docs/ARCHITECTURE.md` | 路徑更新 | 1 |

---

### Task 1: 搬移目錄樹並修正每一處引用

一次搬完並修完所有引用。**不可拆**——只做 `git mv` 會留下 13 處壞掉的引用，那是不能提交的中間狀態。

**Files:**
- Move: `lab/go/skeleton/` → `benchmark/skeleton/go/`
- Move: `lab/prompts/1-fitness-tracker-clarify.md` → `benchmark/cases/fitness-tracker-clarify.md`
- Move: `lab/prompts/1-fitness-tracker-clarify-prd.md` → `benchmark/cases/fitness-tracker-clarify-prd.md`
- Move: `lab/prompts/2-parking-billing-clarify.md` → `benchmark/cases/parking-billing-clarify.md`
- Create: `benchmark/README.md`
- Modify: `.gitignore`, `README.md`, `PLAN.md`, `CLAUDE.md`
- Modify: `benchmark/skeleton/go/README.md`, `benchmark/skeleton/go/docs/ARCHITECTURE.md`

**Interfaces:**
- Produces: 目錄路徑 `benchmark/cases/`、`benchmark/skeleton/go/`、`runs/`；Task 2 會編輯 `benchmark/cases/fitness-tracker-clarify.md` 與 `benchmark/skeleton/go/README.md`。

- [ ] **Step 1: 建立驗證基準線**

搬移前先確認檢查是綠的，這樣之後變紅才歸因得了。把這段存成 `/tmp/linkcheck.py`（整個 task 會重複用到）：

```python
import re, pathlib, sys
bad = []
for md in pathlib.Path('.').rglob('*.md'):
    if '.git' in md.parts: continue
    for link in re.findall(r'\]\((\.[^)]+)\)', md.read_text(encoding='utf-8')):
        if not (md.parent / link.split('#')[0]).exists():
            bad.append(f"{md}: {link}")
print('\n'.join(bad) or 'all links resolve')
sys.exit(1 if bad else 0)
```

Run: `python3 /tmp/linkcheck.py && git status --porcelain`
Expected: `all links resolve`，且工作區乾淨

- [ ] **Step 2: 搬移目錄與 case 檔**

```bash
mkdir -p benchmark/skeleton benchmark/cases
git mv lab/go/skeleton benchmark/skeleton/go
git mv lab/prompts/1-fitness-tracker-clarify.md     benchmark/cases/fitness-tracker-clarify.md
git mv lab/prompts/1-fitness-tracker-clarify-prd.md benchmark/cases/fitness-tracker-clarify-prd.md
git mv lab/prompts/2-parking-billing-clarify.md     benchmark/cases/parking-billing-clarify.md
rmdir lab/prompts lab/go lab
```

- [ ] **Step 3: 確認 git 認得出這是搬移而非刪除重建**

Run: `git status --short | grep '^R' | wc -l`
Expected: `4`（skeleton 整棵樹算一筆或多筆皆可，重點是出現 `R` 而非 `D`＋`??`；若顯示 `D`／`??` 請改用 `git add -A` 後再看 `git diff --cached --stat`，rename 偵測發生在 diff 階段）

- [ ] **Step 4: 更新 `.gitignore`**

`eval-runs/` 與新的 `runs/` 是同一個概念，合併成一個名字。

把 `.gitignore` 改成：

```
.superpowers/
runs/
```

- [ ] **Step 5: 建立 `benchmark/README.md`**

Create `benchmark/README.md`：

````markdown
# Benchmark

固定輸入 ＋ 評分方案，用來量測 ai-bdd 的 skill 有沒有變好。

## 一個 case 是什麼

`cases/` 底下一個檔案就是一個 case。它同時是**題目**與**改考卷的標準**：

| 章節 | 給誰 |
| --- | --- |
| `## 輸入` | **只有這一段餵給 skill** |
| `## 為什麼是這樣措辭` | 人 |
| `## 模糊點` | 人 |
| `## 及格標準` | 人 |
| `## 怎麼比較兩次執行` | 人 |
| `## 實測基準` | 人；每跑一次累積一次 |

**餵食線就是 `## 輸入`。** 評分方案必須進版控——不然沒有東西可以比較——但它同時
就是答案：「模糊點」那一節列的正是你期待 skill 自己問出來的東西。所以隔離靠章節，
不靠目錄。**把 case 檔整份貼進 context 等於直接給答案。**

## 怎麼跑一次

```bash
mkdir -p runs
cp -r benchmark/skeleton/go runs/go-fitness-tracker
```

`runs/` 是 gitignored，所以新 clone 上不存在——`mkdir -p` 不能省。

然後把該 case 的 `## 輸入` 那一段——只有那一段——餵給 `bdd-clarify`。

一次跑會產出：

```
runs/<lang>-<case>/
├── specs/        那次的 example map 與問題檔
└── features/     它產出的 .feature
```

## 兩條規則

**快照。** `benchmark/skeleton/<lang>/` 永遠凍結在「只有走路骨架」的狀態。一次跑的
複本可以長出任何東西，但**沒有任何東西流回範本**。要改進範本就直接改
`benchmark/skeleton/`，之後複製的跑會拿到，正在跑的不會——harness 若在跑到一半被
換掉，比較就失效了。

**不進版控。** `runs/` 是 gitignored。一次 eval 的答案若躺在 git 裡，下一次跑就能
**讀**到而不是**推導**出來，而推導出來的答案和抄來的答案，在產物裡長得一模一樣。
````

- [ ] **Step 6: 從 skeleton README 刪掉 `## One scenario at a time`**

**這一步偏離 spec 的提交切法，是刻意的：** spec 的表格把這一節的刪除放進 commit 2，但它描述的是「一次只能裝一個情境所以要先封存」——那是舊結構的限制，刪它屬於結構變更本身。放進 Task 1 也免掉「先更新這一節的路徑、下一個 task 再刪掉」的白工。commit 2 的標題（move the fitness ambiguity list into its case file）因此更精準。

從 `benchmark/skeleton/go/README.md` 刪除整節（`## One scenario at a time` 到 `## Ground rule` 之前）。

該節內含的產出清單**已在 Step 5 併入 `benchmark/README.md`**，並就地把 `docs/bdd/` 改成 `specs/`（`5dacedc` 早已改名，照抄會把已知錯誤寫進新檔）。該節其餘內容——「archive a completed run before starting another」與 `lab/python/skeleton` 的預告——直接刪除，前者的限制已不存在，後者由 `benchmark/skeleton/` 的結構自己表達。

- [ ] **Step 7: 改 skeleton README 的標題**

`benchmark/skeleton/go/README.md` 第 1 行：

```
- # lab/go/skeleton
+ # benchmark/skeleton/go
```

- [ ] **Step 8: 更新 `ARCHITECTURE.md` 的四處**

檔案：`benchmark/skeleton/go/docs/ARCHITECTURE.md`。用內容定位，不要靠行號。

`## Related documents` 表格裡那一列（相對深度不變，只換目錄名）：

```
- | [../../../prompts/1-fitness-tracker-clarify.md](../../../prompts/1-fitness-tracker-clarify.md) | Evaluating the skills | It is a test fixture, not documentation |
+ | [../../../cases/fitness-tracker-clarify.md](../../../cases/fitness-tracker-clarify.md) | Evaluating the skills | It is a test fixture, not documentation |
```

`## 2. High-Level System Diagram` 的資料流圖起點（`docs/bdd/` 保持原樣，不在本次範圍）：

```
- lab/prompts/ ──▶ CLARIFY ──▶ docs/bdd/*/example-mapping.md
+ benchmark/cases/ ──▶ CLARIFY ──▶ docs/bdd/*/example-mapping.md
```

`## 9. Future Considerations / Roadmap` 的 runbook 第 1 步：

```
- 1. Run `bdd-clarify` against `../../prompts/1-fitness-tracker-clarify.md`
+ 1. Run `bdd-clarify` against `../../cases/fitness-tracker-clarify.md`
```

`## 10. Project Identification`：

```
- | Repository | none — lives inside the `ai-bdd` repo at `lab/go/skeleton` |
+ | Repository | none — lives inside the `ai-bdd` repo at `benchmark/skeleton/go` |
```

- [ ] **Step 9: 更新 `PLAN.md` 的四處**

```
- - [x] `lab/go/skeleton` — Go ＋ godog 實測場（骨架綠，零業務）
+ - [x] `benchmark/skeleton/go` — Go ＋ godog 實測場（骨架綠，零業務）
```

```
- `lab/go/skeleton/` 是六步流程的 dogfooding 對象：Go ＋ godog ＋ 分層
+ `benchmark/skeleton/go/` 是六步流程的 dogfooding 對象：Go ＋ godog ＋ 分層
```

```
- `lab/prompts/` 現在有兩份固定輸入，測互補的失敗模式：
+ `benchmark/cases/` 現在有兩份固定輸入，測互補的失敗模式：
```

```
- > `lab/` 不是 plugin 的慣例目錄，不會被載入，對 plugin 行為零影響。
+ > `benchmark/` 與 `runs/` 都不是 plugin 的慣例目錄（只有 `skills/`、`commands/`、
+ > `agents/`、`hooks/` 是），不會被載入，對 plugin 行為零影響。
```

- [ ] **Step 10: 更新 `CLAUDE.md` 的四處**

```
- to make them fail visibly when they are weak (`lab/`).
+ to make them fail visibly when they are weak (`benchmark/`).
```

```
- ### Testbed (`lab/go/skeleton/`)
+ ### Testbed (`benchmark/skeleton/go/`)
```

```
- `lab/go/skeleton/` is a Go + godog project used to dogfood the pipeline: layered
+ `benchmark/skeleton/go/` is a Go + godog project used to dogfood the pipeline: layered
```

```
- and `api/openapi.yaml`. `lab/` is not a plugin convention directory, so nothing in it is
+ and `api/openapi.yaml`. `benchmark/` is not a plugin convention directory, so nothing in it is
```

- [ ] **Step 11: 更新 root `README.md` 的 Layout 樹**

```
- ├── lab/
- │   ├── prompts/        # fixed inputs — language-agnostic, shared
- │   └── go/
- │       └── skeleton/   # Go + godog dogfooding ground
+ ├── benchmark/
+ │   ├── cases/          # fixed inputs + marking scheme, language-agnostic
+ │   └── skeleton/go/    # the Go + godog template, copied per run
```

- [ ] **Step 12: 確認沒有殘留舊路徑**

Run: `grep -rn "lab/\|skeleton/prompts\|eval-runs" --exclude-dir=.git --exclude-dir=superpowers .`
Expected: 空輸出

排除 `docs/superpowers/` 是因為 spec 與本計畫記述的是搬移這件事本身，本來就會提到舊路徑；把它們算進來會讓這項檢查永遠紅。

- [ ] **Step 13: 確認每一條相對連結都解得到檔案**

這是本次唯一真正會靜默壞掉的東西——壞掉的 markdown 連結不會讓任何指令非零離開。

Run: `python3 /tmp/linkcheck.py`
Expected: `all links resolve`

- [ ] **Step 14: 確認 harness 在新路徑下仍然完整**

Run: `cd benchmark/skeleton/go && make verify`
Expected: 全過（含 gofmt、generated、stamp、golangci-lint、`go vet`、`go test -race`）。約數分鐘。

若 `verify-stamp` 失敗，代表 module 名被動到了——回頭確認 `go.mod` 仍是 `module skeleton`。

- [ ] **Step 15: 確認 plugin 未受影響**

Run: `claude plugin validate . --strict`
Expected: 通過

- [ ] **Step 16: 確認腳本仍然「正確地失敗」**

搬完之後腳本若突然 exit 0，那才是搞砸了——這條鏈的特徵性失敗是「靜默地錯」而不是「壞掉」。

Run: `python3 skills/bdd-clarify/scripts/status.py benchmark/skeleton/go; echo "exit=$?"`
Expected: `exit=1`，且輸出指名缺少什麼產物

- [ ] **Step 17: Commit**

```bash
git add -A .gitignore README.md PLAN.md CLAUDE.md benchmark
git status --short   # 確認沒有夾帶 .claude/ 或其他無關檔案
git commit -F - <<'EOF'
refactor: restructure the testbed into benchmark/ and runs/

WHAT: Split lab/go/skeleton into benchmark/cases/ (fixed inputs and their
      marking scheme), benchmark/skeleton/go/ (the copyable template) and
      a gitignored runs/ for working copies
WHY: One directory played three roles with different lifetimes — the
     language harness, one scenario's charter, and one run's workspace —
     and its README documented the symptom rather than a design:
     "archive a completed run before starting another rather than trying
     to keep two live". That workaround is what the split removes, so the
     section describing it goes too. Keeping run artifacts out of version
     control is the constraint d699aa9 established: a committed answer is
     one a later run can read instead of derive, and the two look
     identical in the artifact.
HOW: Every move is a git mv so git log --follow survives; the case files
     travel three hops now (lab/go/skeleton/prompts, lab/prompts,
     benchmark/cases). The Go module name stays `skeleton` so the
     Makefile's -X skeleton/pkg/version.value keeps resolving, which
     verify-stamp gates. Copies are frozen snapshots with no path back,
     so no harness-versus-scenario file manifest has to be guessed yet.
     Verified with a markdown link checker rather than by eye — nothing
     here changes code, so a broken relative link is the only failure
     mode and it exits zero everywhere. make verify, plugin validate, and
     status.py still exiting 1 on the empty skeleton all confirmed.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 2: 把 fitness 的模糊點清單搬進它的 case 檔

情境知識目前躺在語言 harness 裡：case 檔的及格標準寫著「最低門檻（**見 README**）」，而那 10 個模糊點在 skeleton 的 README。搬完之後 skeleton 是純 harness、case 檔自足。

**Files:**
- Modify: `benchmark/cases/fitness-tracker-clarify.md`
- Modify: `benchmark/skeleton/go/README.md`

**Interfaces:**
- Consumes: Task 1 產生的 `benchmark/cases/`、`benchmark/skeleton/go/` 路徑。
- Produces: 無下游 task。

- [ ] **Step 1: 在 case 檔加入 `## 模糊點`**

在 `benchmark/cases/fitness-tracker-clarify.md` 的 `## 為什麼是這樣措辭` 之後、`## 及格標準` 之前插入。內容逐字取自 skeleton README 的 `## Why this domain`，改成與 `parking-billing-clarify.md` 一致的章節名：

```markdown
## 模糊點

訓練紀錄看起來像 CRUD，所以刻意把重點放在**訓練總容量**而不是記錄本身。容量才是
規則所在，而這些規則是真的模糊：

- 引體向上沒有外加重量——容量算零，還是體重 × 次數？
- 輔助引體向上 −20 kg——負容量，還是（體重 − 20）× 次數？
- 單手划船 20 kg × 8，左右各做——算 160 還是 320？
- 20 kg 空槓熱身組算不算進容量？
- 目標 10 次但只做到 7 次——記成 7，還是「10，失敗」？
- Dropset 60 kg × 8 直接接 40 kg × 6——算一組還是兩組？
- 同一次訓練混用 kg 與 lb
- 史密斯機 60 kg 與槓鈴 60 kg——能不能相加？
- 60 秒棒式既沒有次數也沒有負重——它的容量是多少？
- 「總容量」是對什麼而言——一次訓練、一個動作、一週、還是一個肌群？

大多數人未經提示只會想到其中三、四個。
```

- [ ] **Step 2: 把及格標準的跨檔引用改成自我引用**

同一個檔案的 `## 及格標準`：

```
- 最低門檻（見 README）：一次執行若沒問到**自體重動作**與**單邊動作要不要乘二**，
+ 最低門檻：一次執行若沒問到**自體重動作**與**單邊動作要不要乘二**，
```

- [ ] **Step 3: 從 skeleton README 刪掉 `## Why this domain`**

刪除 `benchmark/skeleton/go/README.md` 的整節——從 `## Why this domain` 到 `## Ground rule` 之前。刪完之後開頭的兩段（「A Go + godog project used to dogfood...」與「It is not a product...」）直接接 `## Ground rule`。

- [ ] **Step 4: 確認 skeleton README 已無情境痕跡**

Run: `grep -in "fitness\|training\|volume\|pull-up\|dropset\|bodyweight" benchmark/skeleton/go/README.md`
Expected: 無輸出

- [ ] **Step 5: 確認連結與既有引用仍然完好**

這個 task 可能在全新 session 執行，`/tmp/linkcheck.py` 不一定還在。重建它：

```python
import re, pathlib, sys
bad = []
for md in pathlib.Path('.').rglob('*.md'):
    if '.git' in md.parts: continue
    for link in re.findall(r'\]\((\.[^)]+)\)', md.read_text(encoding='utf-8')):
        if not (md.parent / link.split('#')[0]).exists():
            bad.append(f"{md}: {link}")
print('\n'.join(bad) or 'all links resolve')
sys.exit(1 if bad else 0)
```

Run: `python3 /tmp/linkcheck.py`
Expected: `all links resolve`

Run: `grep -c "見 README" benchmark/cases/fitness-tracker-clarify.md`
Expected: `0`

- [ ] **Step 6: Commit**

```bash
git add benchmark/cases/fitness-tracker-clarify.md benchmark/skeleton/go/README.md
git commit -F - <<'EOF'
docs(benchmark): move the fitness ambiguity list into its case file

WHAT: Move the ten training-volume ambiguities out of the Go skeleton's
      README into benchmark/cases/fitness-tracker-clarify.md, and drop
      the case file's cross-file "see README" pointer
WHY: The list is scenario knowledge that was sitting in the language
     harness, so the case file could not state its own pass criteria
     without pointing somewhere else — and a second language's skeleton
     would either have duplicated the list or inherited a dangling
     reference. After this the harness README says nothing about
     fitness, which is what lets it be copied for any scenario.
HOW: Rendered into Traditional Chinese under the heading `## 模糊點`, the
     same heading parking-billing-clarify.md already uses, so both case
     files now have the same shape and can be skimmed side by side.
     Verified by grepping the harness README for every scenario term and
     for the removed pointer.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```
