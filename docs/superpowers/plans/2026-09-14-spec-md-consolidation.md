# 計畫一：`prd.md` 併成 `spec.md` Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 讓 CLARIFY 停止寫規格——`bdd-clarify` 改名為 `bdd-discovery` 且只產出問答紀錄，`prd.md` 與 `spec.md` 併成一份十三節的 `spec.md` 由 `bdd-spec` 產出，兩支腳本改讀它。

**Architecture:** 純搬移先自成一筆提交，`git log --follow` 才追得過改名。契約、canonical fixture、兩支解析器仍是同一個原子單位。**這一步 FR 與 AC 仍留在 `spec.md` 裡**——把它們搬進 `.feature` 草稿是計畫二的事。

**Tech Stack:** Markdown（skill 與 reference）、Python 3 標準庫（兩支腳本，無外部依賴）

**Spec:** [`docs/superpowers/specs/2026-09-14-pipeline-restructure-design.md`](../specs/2026-09-14-pipeline-restructure-design.md)

## Global Constraints

- **不引進測試框架。** repo 沒有測試基礎設施。驗證是「建 fixture → 跑腳本 → 驗 exit code 與 stdout」。
- **兩支腳本只用標準庫。**
- **fail loud 不可退化。** 找不到輸入時必須非零離開並指名缺什麼。
- **改名與大幅改寫不得放在同一筆提交。** `git mv` 自成一筆，內容改寫另一筆——否則相似度跌破門檻，`git log --follow` 會靜默斷掉。
- **`.feature` 的 tag 一個都不改。** `@rule-<n> ↔ FR-<n>`、`@example-<n>.<m> ↔ AC-<n>.<m>`。
- **FR／AC 這一步不搬家。** 它們留在 `spec.md` 的 `## User Stories` 底下，維持 v3 的 `**FR-<n>**` ／ `- **AC-<n>.<m>**` 寫法。搬進草稿是計畫二。
- 來源標記三選一不變：`PRD §x`（進來的那份文件這樣寫）／`Q-<n>`／`推論`。
- 散文繁體中文；commit 訊息英文、`WHAT:`／`WHY:`／`HOW:` 三段，trailer 署實際執行的模型。
- 每個 task 結束時工作區乾淨，`audit_skill.py` 對動過的 skill exit 0（跑過 python 要清 `__pycache__`，它會擋住 S5）。

## File Structure

| 檔案 | 責任 | Task |
| --- | --- | --- |
| `skills/bdd-discovery/`（原 `bdd-clarify/`） | 改名後的目錄 | 1 |
| `skills/bdd-spec/references/spec-format.md` | 十三節契約（吸收 `prd-format.md`） | 1 搬、2 併 |
| `skills/bdd-spec/examples/minimal-spec.md` | **canonical 解析對象**，兩支腳本唯一的驗證基準 | 1 搬、2 改 |
| `skills/bdd-spec/references/persona-definition.md` | Personas 現在由 SPEC 寫 | 1 |
| `skills/bdd-spec/scripts/status.py` | 讀 `spec.md` 的 `## Open Questions` | 1 搬、2 改 |
| `skills/bdd-spec/scripts/check_spec.py` | 單檔讀取；`stories_in_spec` 改讀結構 | 2 |
| `skills/bdd-discovery/SKILL.md` | 只問不寫，產出 `clarify-log.md` | 3 |
| `skills/bdd-discovery/references/clarify-log-format.md` | **新增**：問答紀錄的格式 | 3 |
| `skills/bdd-spec/SKILL.md` | 吸收整份文件的產出職責 | 4 |
| `skills/example-mapping/`、`bdd-plan/`、`clarify-loop/`、`CLAUDE.md` | 檔名與路徑指向 | 5 |

---

### Task 1: 純搬移（零內容改寫）

**只有 `git mv` 與路徑字串更新。** 被搬動的檔案內容一個字都不改——這是 `--follow` 能穿過改名的條件。

**Files:**
- Rename: `skills/bdd-clarify/` → `skills/bdd-discovery/`
- Rename: `skills/bdd-discovery/references/persona-definition.md` → `skills/bdd-spec/references/`
- Rename: `skills/bdd-discovery/references/prd-format.md` → `skills/bdd-spec/references/prd-format.md`
- Rename: `skills/bdd-discovery/examples/minimal-prd.md` → `skills/bdd-spec/examples/minimal-prd.md`
- Rename: `skills/bdd-discovery/scripts/status.py` → `skills/bdd-spec/scripts/status.py`
- Modify（僅路徑字串）: `skills/bdd-discovery/SKILL.md`、`skills/clarify-loop/SKILL.md`、`CLAUDE.md`

**Interfaces:**
- Produces: 目錄名 `skills/bdd-discovery/`；`skills/bdd-spec/scripts/status.py`；`skills/bdd-spec/references/{prd-format,persona-definition}.md`；`skills/bdd-spec/examples/minimal-prd.md`

- [ ] **Step 1: 記錄基準**

```bash
rm -rf /tmp/p1 && mkdir -p /tmp/p1/specs/2026-09-09-x
cp skills/bdd-clarify/examples/minimal-prd.md /tmp/p1/specs/2026-09-09-x/prd.md
python3 skills/bdd-clarify/scripts/status.py /tmp/p1 | grep -E "^合計"
python3 -c "
import sys; sys.path.insert(0,'skills/bdd-spec/scripts')
import check_spec, io
print({k: sorted(v) for k,v in sorted(check_spec.frs_in_prd(io.open('skills/bdd-clarify/examples/minimal-prd.md',encoding='utf-8').read()).items())})
"
```

Expected：`合計` 那行是 **2　1　1**；第二條印 `{'1': ['1.1', '1.2', '1.3'], '2': ['2.1']}`。**抄下來——Task 2 Step 8 要求兩者仍然相同。**

- [ ] **Step 2: 目錄改名**

```bash
git mv skills/bdd-clarify skills/bdd-discovery
```

- [ ] **Step 3: 四個檔搬進 `bdd-spec`**

```bash
git mv skills/bdd-discovery/references/persona-definition.md skills/bdd-spec/references/persona-definition.md
git mv skills/bdd-discovery/references/prd-format.md        skills/bdd-spec/references/prd-format.md
git mv skills/bdd-discovery/examples/minimal-prd.md         skills/bdd-spec/examples/minimal-prd.md
git mv skills/bdd-discovery/scripts/status.py               skills/bdd-spec/scripts/status.py
rmdir skills/bdd-discovery/examples 2>/dev/null || true
```

`technical-probes.md` **留在 `bdd-discovery`**——它是問題清單，屬於追問那一側。

`prd-format.md` 這一步只是搬過去，Task 2 才併進 `spec-format.md` 並刪除。

- [ ] **Step 4: 更新指向被搬動檔案的路徑**

只改路徑字串，不改任何句子的意思。

```bash
grep -rn "bdd-clarify\|references/prd-format\|references/persona-definition\|scripts/status\.py\|examples/minimal-prd" skills/ CLAUDE.md | grep -v "^docs/"
```

逐處判斷後改：

| 現在 | 改成 |
| --- | --- |
| `skills/bdd-clarify/scripts/status.py` | `skills/bdd-spec/scripts/status.py` |
| `bdd-discovery/SKILL.md` 裡的 `scripts/status.py` | `../bdd-spec/scripts/status.py` |
| `bdd-discovery/SKILL.md` 裡的 `references/prd-format.md` | `../bdd-spec/references/prd-format.md` |
| `bdd-discovery/SKILL.md` 裡的 `references/persona-definition.md` | `../bdd-spec/references/persona-definition.md` |
| `persona-definition.md` 裡的 `../examples/minimal-prd.md` | 該檔已搬到 `bdd-spec/references/`，相對路徑不變，**確認它仍解得開** |
| `CLAUDE.md:61` | `python3 skills/bdd-spec/scripts/status.py [root]` |

**`bdd-clarify` 這個名字在 SKILL.md 的散文裡也要改成 `bdd-discovery`**——但只改名稱，職責敘述留給 Task 3。

- [ ] **Step 5: 驗證搬移沒弄壞任何東西**

```bash
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for s in bdd-discovery bdd-spec bdd-plan story-splitting clarify-loop example-mapping skill-rules; do
  python3 skills/skill-rules/scripts/audit_skill.py skills/$s >/dev/null 2>&1; echo "  $s exit=$?"
done
claude plugin validate .; echo "validate exit=$?"
rm -rf /tmp/p1 && mkdir -p /tmp/p1/specs/2026-09-09-x
cp skills/bdd-spec/examples/minimal-prd.md /tmp/p1/specs/2026-09-09-x/prd.md
python3 skills/bdd-spec/scripts/status.py /tmp/p1 | grep -E "^合計"
git log --follow --oneline -- skills/bdd-spec/scripts/status.py | wc -l
python3 - <<'EOF'
import re, pathlib, sys
bad = []
for md in pathlib.Path('skills').rglob('*.md'):
    out, fence = [], None
    for line in md.read_text(encoding='utf-8').splitlines():
        m = re.match(r'^(`{3,})', line)
        if m and fence is None: fence = m.group(1); continue
        if m and line.startswith(fence): fence = None; continue
        if fence is None: out.append(line)
    for link in re.findall(r'\]\((\.[^)]+)\)', '\n'.join(out)):
        if not (md.parent / link.split('#')[0]).exists(): bad.append(f"{md}: {link}")
print('\n'.join(bad) or 'all links resolve'); sys.exit(1 if bad else 0)
EOF
git status --porcelain
```

Expected：七個 audit exit 0;`validate` exit 0（只有 root `CLAUDE.md` 警告）;`合計` 仍是 **2　1　1**;`--follow` 的行數 **大於 1**（追得到改名前）;連結全解得到;工作區只有這次的改動。

- [ ] **Step 6: Commit**

```bash
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
git add -A skills CLAUDE.md
git commit -F - <<'EOF'
refactor: move the document-producing files to the step that will write them

WHAT: Rename the clarification skill to discovery and move the format
      contract, the persona reference, the fixture and the progress
      script into the specification skill
WHY: Three documents in this repo say clarification does not write
     specifications, and it has been writing one for three revisions.
     Moving the files is the half of that correction that must not be
     mixed with anything else: git infers a rename from content
     similarity rather than recording it, so a move bundled with a
     rewrite drops below the threshold and log --follow stops at the
     rename without saying so.
HOW: Contents are untouched — only paths, and the references that point
     at them. The technical probe list stays with discovery because it is
     question material, not document material. The format contract lands
     in its new home unchanged and gets merged in the next commit.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 2: 契約、fixture、兩支腳本（原子）

三者一起改。**不可拆**——分開改會留下「腳本找不到 `prd.md`、或找到了卻解析不出 story slug」的中間狀態。

**Files:**
- Modify: `skills/bdd-spec/references/spec-format.md`（併入 `prd-format.md`，重排成十三節）
- Delete: `skills/bdd-spec/references/prd-format.md`
- Rename ＋ Modify: `skills/bdd-spec/examples/minimal-prd.md` → `minimal-spec.md`
- Modify: `skills/bdd-spec/scripts/check_spec.py`
- Modify: `skills/bdd-spec/scripts/status.py`

**Interfaces:**
- Produces: `spec.md` 的十三節與 `### US-<n> · <slug>` 標題形式
- Produces: `frs_in_prd(text) -> dict[str, set[str]]`（簽章不變，來源檔改變）
- Produces: `stories_in_spec(text) -> dict[str, set[str]]`（簽章不變，改讀 `## User Stories` 的結構）

- [ ] **Step 1: 先把 fixture 改名，自成一筆**

```bash
git mv skills/bdd-spec/examples/minimal-prd.md skills/bdd-spec/examples/minimal-spec.md
git commit -q -m "refactor(bdd-spec): rename the fixture ahead of rewriting it

WHAT: Rename minimal-prd.md to minimal-spec.md with no content change
WHY: The next commit rewrites most of this file. Bundling the rename
     with the rewrite would put the similarity below git's rename
     detection threshold and silently break log --follow.
HOW: Pure git mv, verified by --follow before the rewrite lands.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN"
git log --follow --oneline -- skills/bdd-spec/examples/minimal-spec.md | wc -l
```

Expected：行數大於 1。

- [ ] **Step 2: `spec-format.md` 重排成十三節**

新的「章節順序」表放在檔案最前面（`# `spec.md` 的骨架` 標題之後）：

```markdown
| 章節 | 裝什麼 |
| --- | --- |
| `## Document Overview` | 狀態、最後更新、版本修訂歷史、核心關係人 |
| `## Background` | 業務問題與現在什麼在痛 |
| `## Goal` | 做成之後的樣子 |
| `## Scope — In / Out` | In 是這次要交付的能力；Out 分三類，見下 |
| `## Personas` | `### P-<n>` |
| `## User Stories` | `### US-<n> · <slug>`，底下三個分組 `#### FR`／`#### NFR`／`#### Q` |
| `## Non-Functional Requirements` | 跨 story 的品質屬性 |
| `## Assumptions / Constraints` | 沒驗證過的假設，以及做不到／不准這樣做的外部限制 |
| `## Success Metrics` | 能量化就寫數字；量不出來就寫「怎麼知道做成了」 |
| `## Implementation Decisions` | 動哪些模組、介面、API 契約、Schema、架構決定 |
| `## Testing Decisions` | 測試分層與策略 |
| `## Risks` | 風險 |
| `## Open Questions` | 表是機械可解析的索引，小節是人讀的決策史 |

順序固定，不因為某節內容少而調換或省略標題。
`domain-model.md` 仍是獨立檔案，格式見本文最後一節。
```

**併法**（都是整段搬移，不重寫內容）：

| 來源 | 去處 |
| --- | --- |
| `prd-format.md` 的 `## Document Overview 的三個部分`、`## Personas：在三欄上加兩欄`、`## User Stories：三個分組…`、`## 編號規則…`、`## 用 EARS 寫規則`、`## AC 就是綠卡`、`## Non-Functional Requirements 與 Constraints：怎麼分邊`、`## Open Questions：表是索引…`、`## 角色與詞義：兩條界線`、`## 來源標記：三選一` | 整段搬進 `spec-format.md`，順序照上表 |
| `spec-format.md` 原有的 `## Implementation Decisions`、`## Testing Decisions`、`## Risks`、`## 為什麼規則不寫在這裡`、`## 為什麼一個 feature 一份…`、`## domain-model.md 的骨架`及其子節 | **原位保留，內容不動** |
| `spec-format.md` 原有的 `## Problem Statement`、`## Solution` | 刪除——由 `## Background`／`## Goal` 取代 |
| `spec-format.md` 原有的 `## Stories`、`## Acceptance Criteria` | 刪除——併進 `## User Stories` |
| `prd-format.md` 的 `## PRD 的 story 分組不是交付切片` | **刪除**——見下方 D4 說明 |

**刪掉 `## PRD 的 story 分組不是交付切片` 要寫一句話交代**，加在 `## User Stories` 那一節末尾：

```markdown
**只有一個 story 分組，而且寫它的人就是 SPEC。** 舊版有兩個——CLARIFY 先分一次
敘事分組、SPEC 再分一次交付切片——所以需要一條規則仲裁哪個算數。CLARIFY 不再寫
文件之後，那條規則沒有東西可管，已刪除。
```

**`## Scope — In / Out` 的 Out 併成三類**（原本 prd 兩類、spec 兩類，其中一類重複）：

```markdown
**Out** 分三類，因為**誰能回答**不同：

明確排除的能力：
- <這次不做什麼>

接受的風險：
- <因此接受了什麼後果>

回 CLARIFY 補問：
- <推不出來、需要回頭問的缺口> —— 影響什麼定不下來
```

- [ ] **Step 3: `## User Stories` 的標題加 slug**

在該節的說明裡寫明：

```markdown
### US-<n> · <slug>

`<slug>` 是這則 story 的交付識別名，**也是它的 `.feature` 檔名**。用 kebab-case，
不含空白。story 的句子寫在標題下一行，不寫進標題。

```markdown
### US-1 · visitor-billing

身為 P-1（訪客），我想在出場前知道要付多少，以便付完就能開走 ← 推論
```

MUST: slug 在一份 `spec.md` 裡唯一。`check_spec.py` 用它對應 `.feature` 檔名。
```

- [ ] **Step 4: 遷移 fixture 到十三節**

`skills/bdd-spec/examples/minimal-spec.md` 整份換成：

````markdown
# SPEC：社區地下停車場計費系統

## Document Overview

**狀態**：澄清中　**最後更新**：2026-09-14

### 版本修訂歷史

| 版本 | 日期 | 改了什麼 | 為什麼 | 誰 |
| --- | --- | --- | --- | --- |
| 0.1 | 2026-09-14 | 初版 | — | 管委會主委 |

### 核心關係人

| 關係人 | 負責什麼 | 決策權 |
| --- | --- | --- |
| 管委會主委 | 需求方 | 預算內、不動住戶權益的事可自行決定 |
| 消防設備師 | 消防法規 | 只答法規，不做決定 |
| 柵欄廠商 | 已發包的設備 | 只答技術現況，不做決定 |

## Background

訪客臨停目前無法計費，訪客車位常被外部車輛長時間占用。 ← PRD §1

## Goal

讓訪客離場前完成付款，把長時間占位的外部車輛擋在外面。 ← PRD §1

## Scope — In / Out

**In**
- 訪客計費與出場放行 ← PRD §2

**Out**

明確排除的能力：
- 月租費線上收款 ← Q-1（範圍決定：帳務接在社區總帳上）
- 機車計費 ← Q-3

接受的風險：
- 不做月租證掛失名單管理，因此擋不住已作廢的月租證繼續感應出入 ← 推論
- 機車格沒有柵欄，日後要收費是另一次發包 ← Q-3

回 CLARIFY 補問：
- 尖峰時段要不要限制單次停留上限 —— 影響費率上限定不下來 ← 推論

## Personas

### P-1  訪客

**是誰**：不在月租名單上的任何一台車 ← Q-1
**怎麼取得這個身分**：入口抽紙票 ← PRD §2
**跟誰容易混**：跟月租戶靠有沒有感應卡分辨 ← PRD §2
**他要什麼**：出場前知道要付多少、付得掉 ← 推論
**現在什麼讓他痛**：（無資料）

### P-2  月租戶

**是誰**：名單上且未到期的車 ← PRD §2
**怎麼取得這個身分**：感應卡 ← PRD §2
**跟誰容易混**：與「住戶」不完全相同——住戶不一定有月租車位 ← 推論
**他要什麼**：進出不要被擋 ← 推論
**現在什麼讓他痛**：（無資料）

## User Stories

### US-1 · visitor-billing

身為 P-1（訪客），我想在出場前知道要付多少，以便付完就能開走 ← 推論

#### FR

**FR-1**  When 訪客車出場, the system shall 依停留時長計費，前 30 分鐘免費。 ← PRD §3
- **AC-1.1**  停 29 分 → 收 0 元 ← PRD §3
  - Given 一台訪客車已入場並抽了紙票
  - When 它在入場後 29 分鐘於繳費機結帳
  - Then 應收金額是 0 元
- **AC-1.2**  停剛好 30 分 → 收 0 元 ← Q-4
  - Given 一台訪客車已入場並抽了紙票
  - When 它在入場後 30 分鐘整於繳費機結帳
  - Then 應收金額是 0 元
- **AC-1.3**  停 31 分 → 收 30 元 ← PRD §3
  - Given 一台訪客車已入場並抽了紙票
  - When 它在入場後 31 分鐘於繳費機結帳
  - Then 應收金額是 30 元

### US-2 · monthly-pass

身為 P-2（月租戶），我想刷卡就直接出去，以便不用每次都停下來付款 ← 推論

**另外依賴**：FR-1（定義在 US-1 底下）——月租證過期的車不是 P-2，走訪客計費那條路

#### FR

**FR-2**  When 月租戶刷卡出場, the system shall 直接開啟柵欄，不計費。 ← PRD §4
- **AC-2.1**  名單內且未到期 → 開柵欄，金額 0 ← PRD §4
  - Given 一台月租車在名單內且未到期
  - When 它刷卡出場
  - Then 柵欄開啟且應收金額是 0 元

#### NFR

**NFR-1**  出場柵欄自刷卡到開啟不超過 2 秒 ← 推論

#### Q

**Q-2**  消防法規對閘門的要求 ← 擋住 FR-2

## Non-Functional Requirements

- **NFR-2**  車輛進出紀錄保存至少一年 ← 推論

## Assumptions / Constraints

- 沒有車牌辨識，去年區權會否決 ← PRD §5（限制：問過，答案是不准）

## Success Metrics

訪客車位在尖峰時段仍有空位可供住戶的客人使用。 ← 推論

## Implementation Decisions

- **動哪些模組**：新增計費模組；改出場閘門控制 ← 推論
- **介面**：`計算應收金額(入場時間, 出場時間) -> 金額` ← 推論

## Testing Decisions

- 費率邊界（29／30／31 分）走驗收層，用 `.feature` 涵蓋 ← 推論
- 柵欄硬體以假物件替代，不進驗收層 ← 推論

## Risks

| 風險 | 影響 | 目前怎麼辦 |
| --- | --- | --- |
| 消防法規可能要求警報時柵欄常開 | FR-2 的行為要改 | Q-2 待答，先不實作 |

## Open Questions

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q-1 | 月租費帳務在不在範圍 | — | 已答 | 名單與期限在，錢不在 |
| Q-2 | 消防法規對閘門的要求 | 降級 | 待答 | — |
| Q-3 | 機車費率 | — | n/a | 本期不做機車 |
| Q-4 | 剛好停滿 30 分鐘算不算免費 | 邊界 | 已答 | 算免費，30 分整仍在免費區間內 |

### Q-2 消防法規對閘門的要求

**狀態**：待答　**輪次**：3　**該問誰**：消防設備師、柵欄廠商

2026-09-09 問管委會 → 「我知道有規定但背不出來，不敢猜」

**擋住**：FR-2（月租戶刷卡出場開柵欄）
````

**四個不可破壞的東西**：`## Open Questions` 的五欄表與表頭那格的 `Q`（`status.py` 靠它們）；`**FR-<n>**` 的粗體寫法；`- **AC-<n>.<m>**` 的層級；`### US-<n> · <slug>` 的分隔符是 ` · `（U+00B7，前後各一個半形空格）。

- [ ] **Step 5: `stories_in_spec()` 改讀結構**

整個函式換成：

```python
def stories_in_spec(text: str) -> dict[str, set[str]]:
    """spec.md 的 `## User Stories` 段：story-slug -> 它涵蓋的 FR 編號集合。

    併檔之前 slug 住在另一個檔的 `## Stories` 一節，story 與 FR 的對應是
    人手維護的清單。併檔之後那個對應是**結構性的**——FR 就掛在它的 story
    底下，所以走一遍同時拿到 slug 與 FR，不需要第二份清單。

    slug 從 `### US-<n> · <slug>` 取，它也是那則 story 的 `.feature` 檔名。
    """
    section = re.search(
        r"^## User Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return {}
    out: dict[str, set[str]] = {}
    current = None
    for line in section.group(1).splitlines():
        st = re.match(r"^### US-\d+ · (\S+)\s*$", line)
        if st:
            current = st.group(1)
            out.setdefault(current, set())
            continue
        if current:
            for fr in re.findall(r"^\*\*FR-(\d+)\*\*", line):
                out[current].add(fr)
    return out
```

- [ ] **Step 6: 主流程改成單檔讀取，並拿掉恆真的 `unknown` 檢查**

`check()` 裡：

```python
    prds = sorted(specs_dir.glob("*/prd.md"))
    if not prds:
        print(f"{specs_dir} 底下沒有 prd.md —— 先跑 bdd-clarify")
        return 1
```

換成：

```python
    specs = sorted(specs_dir.glob("*/spec.md"))
    if not specs:
        print(f"{specs_dir} 底下沒有 spec.md —— 先跑 bdd-spec")
        return 1
```

迴圈變數 `prd_path` 一律改名 `spec_path`，`ptext` 改 `stext`。

原本「讀第二個檔」那一整段——`spec_path = prd_path.parent / "spec.md"`、
`if not spec_path.exists()`、`if not re.search(r"^## Stories\s*$", stext, re.M)`
三個檢查——**刪除**，改成直接 `stories = stories_in_spec(stext)`，並保留
「一則 story 都沒有」那個檢查（訊息改成指 `## User Stories`）。

`unknown` 那一段**整段刪除**，並在它原本的位置留一段註解：

```python
        # 這裡原本有一個「spec.md 點名了 prd.md 裡沒有的 FR」的檢查。併檔之後
        # 它恆真：slug 與 FR 走的是同一節、同一次遍歷，一份文件不可能跟自己
        # 不一致。它原本要防的「FR 靜默缺席」由 v2_forms_in_prd() 接手——那是
        # 正向檢查，主動找不該存在的標題形式，而不是等對帳對不上。
```

檔頭 docstring 第 27 行的第 12 項（`spec.md 點名的 FR 必須真的存在於 prd.md`）刪除，
並把總數從十二件改成十一件；第 2、9–17、29–35 行提到 `prd.md` 的地方一律改 `spec.md`。

- [ ] **Step 7: `status.py` 改讀 `spec.md`**

```python
    prds = sorted(specs.glob("*/prd.md"))
    if not prds:
        print(f"{specs} 底下沒有任何 prd.md —— 先跑 bdd-clarify")
        return 1
```

換成：

```python
    specs_files = sorted(specs.glob("*/spec.md"))
    if not specs_files:
        print(f"{specs} 底下沒有任何 spec.md —— 先跑 bdd-spec")
        return 1
```

其餘 `prds` / `prd` 變數一律更名，檔頭 docstring 第 8、14 行的 `prd.md` 改 `spec.md`。
**`parse_open_questions()` 一個字都不改**——那一節的格式完全沒變。

- [ ] **Step 8: 回歸——兩個結果必須跟 Task 1 Step 1 逐字相同**

```bash
rm -rf /tmp/p1 && mkdir -p /tmp/p1/specs/2026-09-09-x
cp skills/bdd-spec/examples/minimal-spec.md /tmp/p1/specs/2026-09-09-x/spec.md
python3 skills/bdd-spec/scripts/status.py /tmp/p1 | grep -E "^合計"
python3 -c "
import sys; sys.path.insert(0,'skills/bdd-spec/scripts')
import check_spec, io
t = io.open('skills/bdd-spec/examples/minimal-spec.md',encoding='utf-8').read()
print('frs    :', {k: sorted(v) for k,v in sorted(check_spec.frs_in_prd(t).items())})
print('stories:', {k: sorted(v) for k,v in sorted(check_spec.stories_in_spec(t).items())})
"
```

Expected：`合計` 仍是 **2　1　1**（`## Open Questions` 沒動，變了就是改壞了）；
`frs    : {'1': ['1.1', '1.2', '1.3'], '2': ['2.1']}`（與 Task 1 Step 1 相同）；
`stories: {'monthly-pass': ['2'], 'visitor-billing': ['1']}`。

- [ ] **Step 9: 端到端——健康路徑必須綠**

```bash
mkdir -p /tmp/p1/features
cat > /tmp/p1/features/visitor-billing.feature <<'EOF'
@ready
Feature: 訪客計費
  Rule: 前 30 分鐘免費
    @example-1.1
    Scenario: 停 29 分
      Given 一台訪客車已入場並抽了紙票
      When 它在入場後 29 分鐘於繳費機結帳
      Then 應收金額是 0 元
    @example-1.2
    Scenario: 停剛好 30 分
      Given 一台訪客車已入場並抽了紙票
      When 它在入場後 30 分鐘整於繳費機結帳
      Then 應收金額是 0 元
    @example-1.3
    Scenario: 停 31 分
      Given 一台訪客車已入場並抽了紙票
      When 它在入場後 31 分鐘於繳費機結帳
      Then 應收金額是 30 元
EOF
cat > /tmp/p1/features/monthly-pass.feature <<'EOF'
@ready
Feature: 月租戶出場
  Rule: 月租戶不計費
    @example-2.1
    Scenario: 名單內且未到期
      Given 一台月租車在名單內且未到期
      When 它刷卡出場
      Then 柵欄開啟且應收金額是 0 元
EOF
python3 skills/bdd-spec/scripts/check_spec.py /tmp/p1 >/dev/null 2>&1; echo "健康路徑 exit=$?"
```

Expected：`exit=0`。**tag 一個都沒改,證明接縫沒斷。**

- [ ] **Step 10: 三個檢查必須真的會失敗**

```bash
# (a) slug 少了 —— story 收不到，等於一則 story 都沒有
rm -rf /tmp/p1a && cp -r /tmp/p1 /tmp/p1a
sed -i 's/^### US-1 · visitor-billing$/### US-1  visitor-billing/' /tmp/p1a/specs/2026-09-09-x/spec.md
python3 skills/bdd-spec/scripts/check_spec.py /tmp/p1a 2>&1 | grep -E "story|通過" | head -2
python3 skills/bdd-spec/scripts/check_spec.py /tmp/p1a >/dev/null 2>&1; echo "  (a) exit=$?"

# (b) v2 形式的 FR 標題 —— 上一輪加的正向檢查仍然有效
rm -rf /tmp/p1b && cp -r /tmp/p1 /tmp/p1b
sed -i 's/^\*\*FR-2\*\*  When/#### FR-2  When/' /tmp/p1b/specs/2026-09-09-x/spec.md
python3 skills/bdd-spec/scripts/check_spec.py /tmp/p1b 2>&1 | grep -E "v2 形式|通過" | head -1
python3 skills/bdd-spec/scripts/check_spec.py /tmp/p1b >/dev/null 2>&1; echo "  (b) exit=$?"

# (c) 沒有 spec.md
rm -rf /tmp/p1c && mkdir -p /tmp/p1c/specs/2026-09-09-x /tmp/p1c/features
cp /tmp/p1/features/*.feature /tmp/p1c/features/
python3 skills/bdd-spec/scripts/check_spec.py /tmp/p1c 2>&1 | head -2
python3 skills/bdd-spec/scripts/check_spec.py /tmp/p1c >/dev/null 2>&1; echo "  (c) exit=$?"
```

Expected：三個都 `exit=1`。(a) 指出一則 story 都沒有;(b) 指名 `#### FR-2` 那一行;(c) 說「底下沒有 spec.md —— 先跑 bdd-spec」。

- [ ] **Step 11: fail loud、稽核、`prd-format.md` 已刪**

```bash
python3 skills/bdd-spec/scripts/check_spec.py /tmp/nonexistent; echo "exit=$?"
python3 skills/bdd-spec/scripts/status.py /tmp/nonexistent; echo "exit=$?"
ls skills/bdd-spec/references/prd-format.md 2>&1 | tail -1
grep -rn "prd-format\|minimal-prd\|prd\.md" skills/ | grep -v "^skills/.*docs/"
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for s in bdd-discovery bdd-spec; do python3 skills/skill-rules/scripts/audit_skill.py skills/$s; echo "  $s exit=$?"; done
```

Expected：兩次 `exit=1` 並指名缺什麼;`prd-format.md` 不存在;`prd-format`／`minimal-prd`／`prd.md` 在 `skills/` 底下**零筆**（Task 1 Step 4 改過的路徑此時全部要再改成新檔名——若有殘留就是漏了）;兩個 audit exit 0。

- [ ] **Step 12: Commit**

```bash
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
git add -A skills/bdd-spec
git commit -F - <<'EOF'
feat(bdd-spec): one document, thirteen sections, written by the step that slices

WHAT: Merge the requirements contract into the specification contract,
      migrate the canonical fixture, and point both scripts at a single
      spec.md
WHY: Two documents held one subject. The boundary between them was that
      one belonged to clarification and the other to specification, and
      that boundary disappeared the moment clarification stopped writing.
      What remained was a requirements file with ten sections and a
      specification file whose only unique content was a hand-maintained
      list saying which requirements each story covered.

      That list is now structural: requirements sit under the story they
      serve, so one walk yields both the slug and its requirements. The
      check that compared the two documents went with it — not because
      the risk it covered disappeared, but because a document cannot
      disagree with itself, and the case it was partially catching is
      already caught outright by the check that looks for stale heading
      forms.
HOW: Contract, fixture and both parsers change together because splitting
     them leaves a parser that finds no document, or finds one and cannot
     read a story out of it. The rename of the fixture is its own commit
     so that log --follow survives the rewrite that follows it.

     Stories gain a slug in their heading because the delivery name had
     been living in the file that just disappeared, and it is what names
     the feature file. The rule that arbitrated between two story
     groupings is deleted rather than reworded: there is one grouping
     now, and the step that writes it is the step that slices.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 3: `bdd-discovery` 只問不寫

**Files:**
- Modify: `skills/bdd-discovery/SKILL.md`
- Create: `skills/bdd-discovery/references/clarify-log-format.md`

**Interfaces:**
- Produces: `specs/<date>-<feature>/clarify-log.md` 的格式，Task 4 的 `bdd-spec` 讀它

- [ ] **Step 1: 先讀，再改**

```bash
grep -n "^## \|^### " skills/bdd-discovery/SKILL.md
sed -n '/^## 產物/,/^## /p' skills/bdd-discovery/SKILL.md
```

三個 pass 的**提問內容全部保留**——改的是它們寫到哪裡去。

- [ ] **Step 2: 新增 `references/clarify-log-format.md`**

```markdown
# `clarify-log.md` 的格式

`specs/<date>-<feature>/clarify-log.md`。CLARIFY 的**唯一產物**，也是 `bdd-spec`
唯一的輸入。

**它不是規格。** 它記的是「問了什麼、答了什麼、還開著什麼」，不記「系統該做什麼」。
那條界線是這個 repo 的三份文件共同要求的：CLARIFY 不寫規格，連 story 的定版句子
都不寫。

## 章節

```markdown
# Clarify Log：<主題>

**日期**：<YYYY-MM-DD>　**輪次**：<n>　**結束方式**：收斂／阻塞／零進展／喊停

## 已答

| Q | 問題 | 面向 | 答案 | 誰答的 |
| --- | --- | --- | --- | --- |
| Q-1 | 月租費帳務在不在範圍 | — | 名單與期限在，錢不在 | 管委會主委 |

## 待答

| Q | 問題 | 面向 | 該問誰 | 為什麼還沒答 |
| --- | --- | --- | --- | --- |
| Q-2 | 消防法規對閘門的要求 | 降級 | 消防設備師 | 需求方背不出來，不敢猜 |

## n/a

| Q | 問題 | 為什麼不用問 |
| --- | --- | --- |
| Q-3 | 機車費率 | 本期不做機車 |

## 被推翻的提案

<我提過但需求方說「都不是」的，連同正確答案>
```

「被推翻的提案」值得單獨留：它記錄的是**慣例在這個專案不適用的地方**，而那正是
外人最容易踩錯、文件最少提及的部分。

## 編號

`Q-<n>`，全域唯一，**不重排**。`spec.md` 的來源標記 `← Q-<n>` 回指這裡。

## MUST NOT

MUST NOT: 出現 `US-<n>`、`FR-<n>`、`AC-<n>.<m>` 任何一種編號。**那些是規格的東西**，
出現在這裡就代表 discovery 又在寫規格了。這一條可機械檢查。

MUST NOT: 寫「系統應該…」這種句子。答案照需求方講的記，不要順手翻譯成需求。
翻譯是 `bdd-spec` 的工作，而它翻譯的時候看得到原話才有辦法判斷有沒有過度詮釋。
```

- [ ] **Step 3: 改寫 `bdd-discovery/SKILL.md` 的 `## 產物`**

整節換成：

```markdown
## 產物

```
specs/<date>-<feature>/
└── clarify-log.md    CLARIFY  ★ 唯一產物

（spec.md 是 SPEC 的產物，本 skill 不得寫入）
```

格式 → `references/clarify-log-format.md`。

MUST NOT: 寫 `spec.md`。**連 story 的定版句子都不寫**——那是 `bdd-spec` 的第一句話。
這一條不是分工偏好：Discovery 是模型追問、使用者回答；Formulation 是模型書寫、
使用者審閱。**兩者是方向相反的工作模式**，混在一起的結果是模型一邊問一邊自己補答案。

MUST NOT: 另存一份進度儀表板。進度是**算出來的**——`../bdd-spec/scripts/status.py`
從 `spec.md` 的 `## Open Questions` 直接數；`clarify-log.md` 階段還沒有 `spec.md`，
進度就看這份 log 的三張表。
```

- [ ] **Step 4: 三個 pass 的產出去向改掉**

```bash
grep -n "prd\.md\|spec\.md\|寫進\|填進" skills/bdd-discovery/SKILL.md
```

逐處改：所有「寫進 `prd.md` 的某某節」改成「記進 `clarify-log.md` 的對應表」。
**提問的內容、四選項的形式、面向清單全部不動**——改的只有答案落在哪裡。

`## 完成後` 一節改成：結束方式（收斂／阻塞／零進展／喊停）＋「下一步跑 `bdd-spec`」。
**拿掉任何「就緒判定」的語句**——那是 `example-mapping` 的，而且要等 `spec.md` 存在。

- [ ] **Step 5: 驗證**

```bash
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-discovery; echo "exit=$?"
grep -rn "prd\.md\|寫進 spec\.md\|就緒" skills/bdd-discovery/
grep -rn "FR-\|AC-\|US-" skills/bdd-discovery/references/clarify-log-format.md
```

Expected：audit exit 0;`prd.md` 零筆;`就緒` 只出現在指向 `example-mapping` 的句子（若有）;
`clarify-log-format.md` 裡的 `FR-`／`AC-`／`US-` **只出現在那條 MUST NOT 裡**。

- [ ] **Step 6: Commit**

```bash
git add skills/bdd-discovery
git commit -F - <<'EOF'
docs(bdd-discovery): stop writing the specification, start writing the answers

WHAT: Replace the artifact section with a clarification log, add the
      format for it, and repoint all three passes at it
WHY: This skill has been producing the specification that three other
     documents say it must not produce, including the roadmap's own table
     forbidding it from writing even a story's final sentence. The rule
     had no foothold because asking and writing were fused into one
     stage; separating them gives it one.

     The log is not a lesser specification. It records what was asked,
     what was answered and by whom, and keeps the proposals the requester
     rejected — the place where convention does not apply to this
     project, which is what an outsider gets wrong and what documents
     least often say.
HOW: Every question, option format and probe dimension is unchanged. What
     moved is where the answer lands. The log forbids specification
     identifiers outright, which is mechanically checkable and is the
     line that keeps this skill from drifting back.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 4: `bdd-spec` 吸收文件產出職責

**Files:**
- Modify: `skills/bdd-spec/SKILL.md`

- [ ] **Step 1: 先讀骨架**

```bash
grep -n "^## \|^### " skills/bdd-spec/SKILL.md
```

- [ ] **Step 2: `## 前置確認` 的輸入改成 `clarify-log.md`**

輸入從「`prd.md`」改成「`clarify-log.md` 的三張表」。加一條：

```markdown
MUST NOT: 在這裡問問題。**不訪談**——憑空長出來的內容是缺陷。`clarify-log.md`
的「待答」表裡還有東西時，那些對應的規則就寫不出來，寫進 `## Scope — In / Out`
的「回 CLARIFY 補問」，不要順手決定掉。
```

- [ ] **Step 3: `## 產物` 改成十三節的 `spec.md`**

```markdown
## 產物

```
specs/<date>-<feature>/
├── clarify-log.md    CLARIFY 的產物，本 skill 只讀不寫
└── spec.md           SPEC  ★ 十三節，完整自足

docs/CONTEXT.md       詞彙表，本 skill 可建立或追加，不得重寫既有小節
```

格式 → `references/spec-format.md`。

規則與例子的**定版編號在這裡誕生**：`FR-<n>`、`AC-<n>.<m>`。`.feature` 的
`@rule-<n>`／`@example-<n>.<m>` 回指它們。
```

- [ ] **Step 4: 加一節「切 story 與 slug」**

在 `## 切 story` 開頭加：

```markdown
**只有一個 story 分組，你就是分它的人。** 舊版有兩個——CLARIFY 先分敘事分組、
SPEC 再分交付切片——所以需要一條規則仲裁。CLARIFY 不再寫文件之後，那條規則已刪除。

每則 story 的標題是 `### US-<n> · <slug>`，`<slug>` 是它的 `.feature` 檔名。
kebab-case，一份 `spec.md` 裡唯一——`check_spec.py` 用它對應檔案。
```

- [ ] **Step 5: 全檔掃 `prd.md` 與舊章節名**

```bash
grep -n "prd\.md\|## Stories\|## Problem Statement\|## Solution\|## Acceptance Criteria" skills/bdd-spec/SKILL.md
```

逐處判讀語意再改。`## Stories`／`## Problem Statement`／`## Solution`／
`## Acceptance Criteria` 都已在 Task 2 併掉，指向它們的句子要改。

- [ ] **Step 6: 驗證**

```bash
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
grep -rn "prd\.md" skills/bdd-spec/ || echo "prd.md 零筆"
python3 skills/bdd-spec/scripts/status.py /tmp/p1 | grep "^合計"
```

Expected：audit exit 0;`prd.md` 零筆;`合計` 仍是 **2　1　1**。

- [ ] **Step 7: Commit**

```bash
git add skills/bdd-spec/SKILL.md
git commit -F - <<'EOF'
docs(bdd-spec): take over the document, and say what it may not do to get it

WHAT: Point the inputs at the clarification log, replace the artifact
      section with the thirteen-section document, and state the slug rule
      for stories
WHY: This step now writes the first sentence anyone reads, which makes
     its existing prohibition load-bearing rather than decorative: it
     must not interview. A gap it cannot close from the log is not a gap
     to decide away — it goes back as a question, and the scope section
     has a category for exactly that.
HOW: The story grouping rule is stated as one grouping owned by one step,
     rather than as arbitration between two. The slug moves into the
     story heading because the file that used to carry delivery names no
     longer exists, and it is what names the feature file.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 5: 其餘指向（同形狀，一批做完）

**Files:**
- Modify: `skills/example-mapping/SKILL.md`、`skills/bdd-plan/SKILL.md`、`skills/clarify-loop/SKILL.md`、`skills/story-splitting/SKILL.md`、`CLAUDE.md`

- [ ] **Step 1: 四個 skill 的檔名指向**

```bash
grep -rn "prd\.md\|bdd-clarify" skills/example-mapping skills/bdd-plan skills/clarify-loop skills/story-splitting
```

逐處改：`prd.md` → `spec.md`；`bdd-clarify` → `bdd-discovery`。

**`example-mapping` 這一步只改檔名，行為不動**——它仍然讀文件、攤終端地圖、提出
就緒判定。改成產生 `.feature` 草稿是計畫二。

- [ ] **Step 2: benchmark 的兩份 README**

```bash
grep -n "prd\.md\|bdd-clarify" benchmark/skeleton/go/README.md benchmark/README.md
```

`benchmark/skeleton/go/README.md` 有五處（第 12、226、229、249、256 行）、
`benchmark/README.md` 有兩處（第 44 行的 `bdd-clarify`、第 60 行的產物路徑）。

- `prd.md` → `spec.md`
- 第 229 行的表頭 `` | `prd.md` / `spec.md` | Gherkin | `` → `` | `spec.md` | Gherkin | ``——**兩個檔併成一個了，表頭不該再列兩個**
- `bdd-clarify` → `bdd-discovery`；`benchmark/README.md:44` 的「把輸入交給 `bdd-clarify`」
  要改成交給 `bdd-discovery`，且第 60 行的產物從 `specs/<date>-<feature>/prd.md`
  改成 `clarify-log.md` ＋ `spec.md` 兩行

**tag 欄一個都不改**——`@example-1.1` 是對的。

- [ ] **Step 3: `PLAN.md` 只改「還在用的」，歷史紀錄不動**

```bash
grep -n "bdd-clarify\|prd\.md" PLAN.md
```

判準很簡單：

| 這一行在做什麼 | 怎麼處理 |
| --- | --- |
| **指路**（skill 表的一列、待辦項、檔案路徑） | 改成新名稱 |
| **記錄某天發生了什麼**（`- **2026-09-09**：…`、`- [x]` 的已完成項） | **不動**——那些記的是當時的事實 |

具體：第 198 行 skill 表的 `bdd-clarify` 那一列改名並改職責敘述為「三個 pass，
只問不寫」；第 428、463、468、477 行的待辦項改名；第 308 行與所有 `- [x]` 項**不動**。

- [ ] **Step 4: `CLAUDE.md` 的指令與敘述**

```bash
grep -n "bdd-clarify\|prd\.md\|status\.py\|check_spec\.py" CLAUDE.md
```

改：`skills/bdd-clarify/scripts/status.py` → `skills/bdd-spec/scripts/status.py`；
`prd.md` → `spec.md`；skills 清單裡的 `bdd-clarify` → `bdd-discovery`。

**「The three disciplines」那一段不動**——它描述的正是這次要落實的東西，本來就對。

- [ ] **Step 5: 全 repo 收尾驗證**

```bash
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for s in bdd-discovery bdd-spec bdd-plan story-splitting clarify-loop example-mapping skill-rules; do
  python3 skills/skill-rules/scripts/audit_skill.py skills/$s >/dev/null 2>&1; echo "  $s exit=$?"
done
claude plugin validate .; echo "validate exit=$?"
grep -rn "bdd-clarify\|prd\.md\|prd-format\|minimal-prd" skills/ CLAUDE.md | grep -v "^docs/superpowers/" || echo "舊名稱零筆"
grep -rn "## Stories\|## Problem Statement\|## Acceptance Criteria" skills/ || echo "舊章節名零筆"
python3 - <<'EOF'
import re, pathlib, sys
bad = []
for md in pathlib.Path('.').rglob('*.md'):
    if '.git' in md.parts or 'runs' in md.parts: continue
    out, fence = [], None
    for line in md.read_text(encoding='utf-8').splitlines():
        m = re.match(r'^(`{3,})', line)
        if m and fence is None: fence = m.group(1); continue
        if m and line.startswith(fence): fence = None; continue
        if fence is None: out.append(line)
    for link in re.findall(r'\]\((\.[^)]+)\)', '\n'.join(out)):
        if not (md.parent / link.split('#')[0]).exists(): bad.append(f"{md}: {link}")
print('\n'.join(bad) or 'all links resolve'); sys.exit(1 if bad else 0)
EOF
git log --follow --oneline -- skills/bdd-spec/examples/minimal-spec.md | wc -l
git log --follow --oneline -- skills/bdd-spec/scripts/status.py | wc -l
git status --porcelain
```

Expected：七個 audit exit 0;`validate` exit 0（只有 root `CLAUDE.md` 警告）;
舊名稱與舊章節名皆零筆（`docs/superpowers/` 的歷史文件不在掃描範圍）;連結全解得到;
兩個 `--follow` 行數皆 **大於 1**;工作區乾淨。

- [ ] **Step 6: Commit**

```bash
git add -A skills CLAUDE.md
git commit -F - <<'EOF'
docs: repoint the four skills and the guide at the merged document

WHAT: Change every reference to the requirements file and the old
      clarification skill name across the remaining skills and CLAUDE.md
WHY: The file and the skill both changed name in this branch, and a
     reference to either now sends a reader looking for something that is
     not there. The mapping skill keeps its current behaviour — reading
     the document and laying out a map — and only learns the new
     filename; turning it into a producer of Gherkin drafts belongs to
     the next plan and would be invisible in a rename commit.
HOW: The three disciplines section of CLAUDE.md is deliberately
     untouched. It already describes the separation this branch
     implements, which is why the branch exists.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```
