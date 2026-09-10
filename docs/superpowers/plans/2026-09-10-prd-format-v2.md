# `prd.md` 格式 v2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 `prd.md` 從八節擴成十一節，加入 Document Overview／Background／Goal／Success Metrics／Personas／User Stories，並把 FR 從平鋪的 `## Functional Requirements` 改成掛在 `### US-<n>` 底下的 `#### FR-<n>`。

**Architecture:** 契約、canonical fixture、解析器三者是同一個原子單位——分開改會留下「腳本解析不到任何 FR 卻不報錯」的中間狀態，那正是這個 repo 反覆在防的靜默失敗。所以 Task 1 三個一起改、一起驗；Task 2、3 只改描述它的 skill 文件。

**Tech Stack:** Markdown（skill 與 reference）、Python 3 標準庫（`check_spec.py`，無外部依賴）

**Spec:** [`docs/superpowers/specs/2026-09-10-prd-format-v2-design.md`](../specs/2026-09-10-prd-format-v2-design.md)

## Global Constraints

- **不引進測試框架。** repo 沒有測試基礎設施，`pytest` 沒安裝。驗證是「建 fixture → 跑腳本 → 驗 exit code 與 stdout」。
- **`check_spec.py` 只用標準庫。**
- **fail loud 不可退化。** 找不到輸入時必須非零離開並指名缺什麼。
- **`status.py` 不得改動。** 它只讀 `## Open Questions`，那一節的位置與格式在 v2 完全不變。改到它就是改壞了。
- **FR 與 NFR 的編號全域唯一**，不是每則 story 各自從 1 起算。
- **沿用 v1 不變**：來源標記三選一（`PRD §x`／`Q<n>`／`推論`）、`## Open Questions` 的五欄表與 `面向` 值域、AC 與 EX 的兩層讀者、EARS 引用 `docs/sdd.md` 不重述。
- 散文繁體中文，commit 訊息英文、`WHAT:`／`WHY:`／`HOW:` 三段。
- 每個 task 結束時工作區乾淨，`audit_skill.py` 對動過的 skill exit 0。

## File Structure

| 檔案 | 責任 | Task |
| --- | --- | --- |
| `skills/bdd-clarify/references/prd-format.md` | 格式契約：十一節的順序與各節內容、三條參照完整性規則、已知弱點 | 1 |
| `skills/bdd-clarify/examples/minimal-prd.md` | **canonical 解析對象**，兩支腳本唯一的驗證基準 | 1 |
| `skills/bdd-spec/scripts/check_spec.py` | `frs_in_prd()` 改讀 `## User Stories` 段的 `#### FR-<n>` | 1 |
| `skills/bdd-clarify/SKILL.md` | 產物一節、流程圖的 Pass 3 路由、各處章節名引用 | 2 |
| `skills/bdd-spec/SKILL.md` | 加 D2 的規則：PRD 的 story 分組是線索不是約束 | 3 |
| `skills/bdd-plan/SKILL.md` | 第 63 行引用了即將消失的 `Functional Requirements` | 3 |

---

### Task 1: 契約、fixture、解析器（原子）

三者一起改。**不可拆**——分開改會留下「`check_spec.py` 找不到 `## Functional Requirements`、回傳空 dict、不報錯」的中間狀態。

**Files:**
- Modify: `skills/bdd-clarify/references/prd-format.md`
- Modify: `skills/bdd-clarify/examples/minimal-prd.md`
- Modify: `skills/bdd-spec/scripts/check_spec.py`

**Interfaces:**
- Produces: `## User Stories` → `### US-<n>` → `#### FR-<n>` → `##### AC` / `##### Examples` 的階層，Task 2、3 的文件都要跟它一致
- Produces: `frs_in_prd(text) -> dict[str, set[str]]`（簽章不變，來源段改變）

- [ ] **Step 1: 記錄改動前的基準，之後每一步都要對得回來**

```bash
rm -rf /tmp/v2 && mkdir -p /tmp/v2/specs/2026-09-09-x
cp skills/bdd-clarify/examples/minimal-prd.md /tmp/v2/specs/2026-09-09-x/prd.md
python3 skills/bdd-clarify/scripts/status.py /tmp/v2
python3 -c "
import sys; sys.path.insert(0,'skills/bdd-spec/scripts')
import check_spec, io
print(check_spec.frs_in_prd(io.open('skills/bdd-clarify/examples/minimal-prd.md',encoding='utf-8').read()))
"
```

Expected：`status.py` 印出 **已答 2／n/a 1／待答 1**，覆蓋為 `邊界` ＋ `降級`；`frs_in_prd` 印出 `{'1': {'1.1','1.2','1.3'}, '2': {'2.1'}}`。**把這兩個結果抄下來——Step 6 要逐字比對。**

- [ ] **Step 2: 改寫 `prd-format.md` 的「章節順序」表**

整張表換成：

```markdown
| 章節 | 裝什麼 |
| --- | --- |
| `## Document Overview` | 狀態、最後更新、版本修訂歷史、核心關係人——見下方三小節 |
| `## Background` | 業務問題與現在什麼在痛 |
| `## Goal` | 做成之後的樣子 |
| `## Success Metrics` | 能量化就寫數字；量不出來就寫「怎麼知道做成了」 |
| `## Scope — In / Out` | In 是這次要交付的能力；Out 是**明確邊界，不是免責條款** |
| `## Personas` | `### P-<n>`，這個 feature 用到的角色——不外流到 `docs/CONTEXT.md` |
| `## User Stories` | `### US-<n>`，每則底下掛它的 `#### FR-<n>` 與 story 專屬的 `#### NFR-<n>` |
| `## Non-Functional Requirements` | 跨 story 的品質屬性（個資保存、可用性） |
| `## Open Questions` | 表是機械可解析的索引，小節是人讀的決策史——見下 |
| `## Assumptions / Constraints` | 沒驗證過的假設，以及做不到／不准這樣做的外部限制 |

順序固定，不因為某節內容少而調換或省略標題。
```

原本「Out 要分辨兩類」那一段**保留不動**。

- [ ] **Step 3: 在 `prd-format.md` 新增五節**

在「章節順序」表之後、「來源標記：三選一」之前插入。逐字：

````markdown
## Document Overview 的三個部分

```markdown
## Document Overview

**狀態**：澄清中　**最後更新**：2026-09-10

### 版本修訂歷史

| 版本 | 日期 | 改了什麼 | 為什麼 | 誰 |
| --- | --- | --- | --- | --- |
| 0.3 | 2026-09-10 | FR-1 免費額度從「每趟 30 分」改成「一天總共 30 分」 | Q15：早晚各來一趟的客人在舊規則下反而被多收 | 管委會主委 |

### 核心關係人

| 關係人 | 負責什麼 | 決策權 |
| --- | --- | --- |
| 管委會主委 | 需求方，日常庶務與月租帳 | 預算內、不動住戶權益的事可自行決定 |
| 區權會 | 住戶大會 | 動到住戶權益或預算外支出；兩個月開一次 |
| 設備廠商 | 已發包的柵欄／繳費機 | 只答技術現況，不做決定 |
```

**狀態只有四個值**，判準是「下游能不能開工」：

| 狀態 | 意思 |
| --- | --- |
| `澄清中` | 還有待答問題，SPEC 不得開始 |
| `待對焦` | 問題都答完了，等需求方過目 |
| `已對焦` | 需求方同意，SPEC 可以開始 |
| `已凍結` | SPEC 已經消費過它，之後任何改動都要開新版本 |

MUST NOT: 狀態欄寫題數。「9 題已答／10 題待答」是 `scripts/status.py` **算出來的**，
寫進文件等於同一個事實兩份，而它們遲早不一樣。

**版本修訂歷史的觸發條件**：文件到 `已對焦` 之後，任何動到 FR／AC／EX 的改動都要
留一列——SPEC 可能已經照著寫了 `.feature`。`已對焦` 之前的來回不記，那是澄清過程
本身，決策史在 `## Open Questions` 的小節裡。

**「為什麼」那欄不可省。** 只記改了什麼，半年後有人會把它改回去。

**核心關係人的「決策權」不是裝飾。** 同一個人給的答案，有些是定案、有些只是提案
（要開會表決），下游要分得出來，否則會拿沒表決過的東西去寫規格。

MUST: `## Open Questions` 每一題的「該問誰」必須指向這張表裡的某個關係人。指到
表外的人＝那個人沒被登記為關係人，是文件的洞不是筆誤。

MUST NOT: 在這一節列「誰卡著哪幾題」。那是算出來的——每一題自己宣告掛在誰身上。

## Personas：在三欄上加兩欄

```markdown
### P-1  訪客

**是誰**：不在月租名單上的任何一台車 ← Q8
**怎麼取得這個身分**：進場時未刷月租卡 ← Q5（待答，取決於辨識方式）
**跟誰容易混**：跟「住戶的客人」不是同一件事——規約說訪客車位給住戶親友，
但實務上沒有東西在分辨 ← Q7
**他要什麼**：出場前知道要付多少、付得掉 ← 推論
**現在什麼讓他痛**：（無資料）
```

前三欄沿用舊的 `## Actors`。新增的是**「他要什麼」**與**「現在什麼讓他痛」**——
FR 掛在 story、story 掛在 persona，**沒有目標的 persona 會產出沒有理由的 story**。

**「無資料」是合法的值，而且比編一個好。** 需求方常常是在轉述別人的痛（主委轉述
住戶、住戶轉述他的客人）。標「無資料」讓那件事在文件上看得見；填「訪客希望流程
順暢」則讓它消失。

## User Stories：強制指向 Persona，FR 掛在底下

```markdown
### US-1  身為 P-1（訪客），我想在出場前知道要付多少，以便準備好再開到閘門 ← 推論

#### FR-1  When 訪客車出場, the system shall 依停留時長計費，前 30 分鐘免費。 ← PRD §3

##### AC
- AC-1.1  訪客在繳費機看得到金額與停留時長 ← 推論

##### Examples
- EX-1.1  停 29 分 → 收 0 元 ← PRD §3

#### NFR-1  出場柵欄自刷卡到開啟不超過 2 秒 ← Q12
```

MUST: `身為` 後面必須是一個 `P-<n>`。指到不存在的 persona＝角色沒被定義，是洞。

MUST: **FR 與 NFR 的編號全域唯一。** `FR-1` 在 US-1 底下、`FR-4` 在 US-2 底下都
可以，但不得有兩個 `FR-1`——否則 `EX-1.1` 指向兩個地方，`.feature` 的
`@example-1.1` 失去意義。

**NFR 有兩個家，判準是跨不跨 story**：只服務某一則 story 的掛在那則底下；跨全部的
留在 `## Non-Functional Requirements`。搬家時編號不變，這樣「NFR-3 從 story 專屬
改成跨 story」在版本修訂歷史裡才追得到。

### 一條 FR 服務多則 story

巢狀強迫一對多，共用的 FR 只能住一個地方。**住主要的那則，其他宣告依賴**：

```markdown
### US-2  身為 P-2（主委），我想知道訪客格有沒有被外人占走 ← Q2

**另外依賴**：FR-3（定義在 US-1 底下）
```

**已知弱點：這行宣告沒有東西在檢查它。** US-2 忘了宣告依賴 FR-3，文件看起來完全
正常，只是 SPEC 為 US-2 切 story 時會漏掉一條規則。這是巢狀換可讀性的代價，不是
可以設計掉的東西——除非把 FR 平鋪回去。

## PRD 的 story 分組不是交付切片

`## User Stories` 是**需求方視角的敘事分組**。SPEC **得重新分組，不必解釋為什麼跟
PRD 不一樣**；但如果它照抄 PRD 的分組，要說得出為什麼那個分組剛好也是好的交付
邊界。

這條規則存在的理由：2026-09-09 的重構把切 story 從 CLARIFY 搬到 SPEC，因為切分
需要範圍先穩定——實測中需求方在 15 題範圍題裡改口三次。**那次的根源是目錄不是
章節**（主文件住在 `specs/<slice-slug>/`，要有 slug 才有目錄），現在 story 只是
一份檔案裡的章節，改邊界＝搬一段。逼迫早切的力量已經不存在，剩下的風險只有
「SPEC 被現成的分組綁住」，由這條規則管。
````

- [ ] **Step 4: 遷移 `examples/minimal-prd.md` 到 v2**

整份換成（逐字，這是兩支腳本的 canonical 解析對象）：

````markdown
# PRD：社區地下停車場計費系統

## Document Overview

**狀態**：澄清中　**最後更新**：2026-09-09

### 版本修訂歷史

| 版本 | 日期 | 改了什麼 | 為什麼 | 誰 |
| --- | --- | --- | --- | --- |
| 0.1 | 2026-09-09 | 初版 | — | 管委會主委 |

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

## Success Metrics

訪客車位在尖峰時段仍有空位可供住戶的客人使用。 ← 推論

## Scope — In / Out

**In**
- 訪客計費與出場放行 ← PRD §2

**Out**
- 月租費線上收款 ← Q1（範圍決定：帳務接在社區總帳上）
- 機車計費 ← Q3（範圍決定；風險：機車格沒有柵欄，日後要收費是另一次發包）

## Personas

### P-1  訪客

**是誰**：不在月租名單上的任何一台車 ← Q1
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

### US-1  身為 P-1（訪客），我想在出場前知道要付多少，以便付完就能開走 ← 推論

#### FR-1  When 訪客車出場, the system shall 依停留時長計費，前 30 分鐘免費。 ← PRD §3

##### AC
- AC-1.1  訪客在繳費機看得到金額與停留時長 ← 推論
- AC-1.2  付款完成前柵欄不開 ← PRD §3

##### Examples
- EX-1.1  停 29 分 → 收 0 元 ← PRD §3
- EX-1.2  停剛好 30 分 → 收 0 元 ← Q4
- EX-1.3  停 31 分 → 收 30 元 ← PRD §3

### US-2  身為 P-2（月租戶），我想刷卡就直接出去，以便不用每次都停下來付款 ← 推論

#### FR-2  When 月租戶刷卡出場, the system shall 直接開啟柵欄，不計費。 ← PRD §4

##### AC
- AC-2.1  月租戶出場不需經過繳費機 ← PRD §4

##### Examples
- EX-2.1  名單內且未到期 → 開柵欄，金額 0 ← PRD §4

#### NFR-1  出場柵欄自刷卡到開啟不超過 2 秒 ← 推論

## Non-Functional Requirements

- NFR-2  車輛進出紀錄保存至少一年 ← 推論

## Open Questions

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q1 | 月租費帳務在不在範圍 | — | 已答 | 名單與期限在，錢不在 |
| Q2 | 消防法規對閘門的要求 | 降級 | 待答 | — |
| Q3 | 機車費率 | — | n/a | 本期不做機車 |
| Q4 | 剛好停滿 30 分鐘算不算免費 | 邊界 | 已答 | 算免費，30 分整仍在免費區間內 |

### Q2 消防法規對閘門的要求

**狀態**：待答　**輪次**：3　**該問誰**：消防設備師、柵欄廠商

2026-09-09 問管委會 → 「我知道有規定但背不出來，不敢猜」

**擋住**：FR-11（緊急車輛進場）

## Assumptions / Constraints

- 沒有車牌辨識，去年區權會否決 ← Q2（限制：問過，答案是不准）
````

**三個不可破壞的東西**：`## Open Questions` 的五欄表（`status.py` 用 `len(cells) != 5` 篩）、`#### FR-<n>` 的層級、`EX-<n>.<m>` 的寫法。

- [ ] **Step 5: 改 `check_spec.py` 的 `frs_in_prd()`**

整個函式換成：

```python
def frs_in_prd(text: str) -> dict[str, set[str]]:
    """prd.md 的 `## User Stories` 段：FR 編號 -> EX 編號集合。

    FR 掛在 `### US-<n>` 底下，所以標題層級是 `#### FR-<n>` 而不是 `###`。

    來源是 prd.md 而不是 spec.md：例子的定版編號誕生在 CLARIFY，
    spec.md 只記哪些 FR 湊成一則 story，不重述例子。

    FR 編號全域唯一，不是每則 story 各自從 1 起算——否則 EX-1.1 會指向
    兩個地方，`.feature` 的 @example-1.1 就失去意義。

    碰到新的 `### US-` 就把 current 清掉：story 標題與第一條 FR 之間的散文
    （例如「另外依賴：FR-3」）不該被算成上一則 story 最後一條 FR 的例子。
    """
    section = re.search(
        r"^## User Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return {}
    out: dict[str, set[str]] = {}
    current = None
    for line in section.group(1).splitlines():
        if re.match(r"^### US-", line):
            current = None
            continue
        fr = re.match(r"^#### FR-(\d+)\b", line)
        if fr:
            current = fr.group(1)
            out.setdefault(current, set())
            continue
        if current:
            for ex in re.findall(r"\bEX-(\d+\.\d+)\b", line):
                out[current].add(ex)
    return out
```

同時改兩處文字：
- 檔頭 docstring 第 22 行的「`## Functional Requirements` 段」→「`## User Stories` 段」
- 第 145 行的錯誤訊息「沒有 `## Functional Requirements`」→「沒有 `## User Stories`」

- [ ] **Step 6: 回歸——兩個結果必須跟 Step 1 逐字相同**

```bash
rm -rf /tmp/v2 && mkdir -p /tmp/v2/specs/2026-09-09-x
cp skills/bdd-clarify/examples/minimal-prd.md /tmp/v2/specs/2026-09-09-x/prd.md
python3 skills/bdd-clarify/scripts/status.py /tmp/v2
python3 -c "
import sys; sys.path.insert(0,'skills/bdd-spec/scripts')
import check_spec, io
print(check_spec.frs_in_prd(io.open('skills/bdd-clarify/examples/minimal-prd.md',encoding='utf-8').read()))
"
```

Expected：**已答 2／n/a 1／待答 1、覆蓋 `邊界`＋`降級`**（與 Step 1 相同——`## Open Questions` 沒動，變了就是改壞了）；`{'1': {'1.1','1.2','1.3'}, '2': {'2.1'}}`（與 Step 1 相同）。

- [ ] **Step 7: 端到端——`check_spec.py` 對完整 fixture 跑一次**

```bash
mkdir -p /tmp/v2/features
cat > /tmp/v2/specs/2026-09-09-x/spec.md <<'EOF'
# SPEC

## Stories

### visitor-billing
涵蓋 FR-1。

### monthly-pass-exit
涵蓋 FR-2。
EOF
cat > /tmp/v2/features/visitor-billing.feature <<'EOF'
Feature: Visitor billing

  Rule: 前 30 分鐘免費

    @example-1.1
    Scenario: 停 29 分
      Given 一台訪客車停了 29 分鐘
      When 它在繳費機結帳
      Then 應收金額是 0 元
EOF
python3 skills/bdd-spec/scripts/check_spec.py /tmp/v2; echo "exit=$?"
```

Expected：非零離開，且輸出指出 `visitor-billing` 漏了 `EX-1.2`、`EX-1.3`，`monthly-pass-exit` 漏了 `EX-2.1`。**這證明 FR→EX 的對應在新層級下真的被抓到了**，不是回傳空 dict 之後「什麼都沒漏」。

- [ ] **Step 8: 缺 `## User Stories` 要大聲失敗**

```bash
rm -rf /tmp/v2bad && mkdir -p /tmp/v2bad/specs/2026-09-09-y /tmp/v2bad/features
grep -v "^## User Stories" skills/bdd-clarify/examples/minimal-prd.md > /tmp/v2bad/specs/2026-09-09-y/prd.md
cp /tmp/v2/specs/2026-09-09-x/spec.md /tmp/v2bad/specs/2026-09-09-y/
python3 skills/bdd-spec/scripts/check_spec.py /tmp/v2bad; echo "exit=$?"
```

Expected：`exit=1`，訊息指名該 feature 目錄與「沒有 `## User Stories`」。

- [ ] **Step 9: fail loud 沒退化**

```bash
python3 skills/bdd-spec/scripts/check_spec.py /tmp/nonexistent-v2; echo "exit=$?"
python3 skills/bdd-clarify/scripts/status.py /tmp/nonexistent-v2; echo "exit=$?"
```

Expected：兩次都 `exit=1` 並指名缺什麼。

- [ ] **Step 10: 稽核與 diff 範圍確認**

```bash
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify; echo "exit=$?"
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
git diff --stat -- skills/bdd-clarify/scripts/status.py
```

Expected：兩個 audit exit 0；`status.py` 的 diff **為空**（Global Constraints 明訂不得改動）。

- [ ] **Step 11: Commit**

```bash
git add skills/bdd-clarify/references/prd-format.md \
        skills/bdd-clarify/examples/minimal-prd.md \
        skills/bdd-spec/scripts/check_spec.py
git commit -F - <<'EOF'
feat(bdd-clarify): prd.md v2 — requirements move under user stories

WHAT: Rewrite the prd.md contract to eleven sections, migrate the
      canonical fixture, and repoint frs_in_prd at the User Stories
      section where FRs now sit at H4
WHY: v1 parses and drives SPEC but does not read like a PRD anyone would
     take to a product manager: no document status, no revision history,
     no list of who the open questions are blocked on, and pain sharing a
     heading with goal and metrics. Nesting requirements under stories is
     safe now for a reason that was not true in v1 — the force that made
     slicing come first was the directory, not the heading, and one
     prd.md per feature makes moving a requirement between stories an
     edit rather than a restructure.
HOW: Contract, fixture and parser change together because splitting them
     leaves a state where the parser finds no requirements and reports
     success — the silent failure this pipeline exists to catch. The
     parser resets its current requirement at each story heading, so
     prose between a story title and its first requirement cannot have
     its example numbers attributed to the previous story's last one.
     Two sections deliberately store less than they could: status carries
     no question counts and the stakeholder table carries no blocking
     list, both replaced by referential rules a script can check. status.py
     is untouched and its numbers verified identical before and after,
     because Open Questions did not move.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 2: `bdd-clarify/SKILL.md` 跟上新格式

**Files:**
- Modify: `skills/bdd-clarify/SKILL.md`

**Interfaces:**
- Consumes: Task 1 的 `prd-format.md` 與 `minimal-prd.md`（不一致時以它們為準）

- [ ] **Step 1: 改流程圖的 Pass 3 路由**

`## 流程` 的 code block 裡，現在寫著：

```
Pass 3 · 技術 —— 逐項掃過 technical-probes.md
  量得出來的 → ## Non-Functional Requirements
  不可協商的外部限制 → ## Assumptions / Constraints
```

改成：

```
Pass 3 · 技術 —— 逐項掃過 technical-probes.md
  量得出來的 → NFR：只服務一則 story 的掛在那則底下，
              跨 story 的進 ## Non-Functional Requirements
  不可協商的外部限制 → ## Assumptions / Constraints
```

- [ ] **Step 2: 改 `## 產物` 一節裡的章節清單**

該節若列出 `prd.md` 的章節，換成 Task 1 的十一節；若只寫「格式 → `references/prd-format.md`」則不動。**先讀再改，不要憑印象。**

- [ ] **Step 3: 全檔掃一次章節名引用**

```bash
grep -n "## Functional Requirements\|## Actors\|Problem / Goal" skills/bdd-clarify/SKILL.md
```

每一處逐一判斷：`## Actors` → `## Personas`；`## Functional Requirements` → `## User Stories` 底下的 `#### FR-<n>`；`Problem / Goal / Success Metrics` → 拆成的三節。**逐處判讀語意再改，不要無差別取代**——有些句子講的是「功能需求」這個概念而不是那個章節名。

- [ ] **Step 4: 驗證**

```bash
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify; echo "exit=$?"
grep -n "## Functional Requirements\|## Actors" skills/bdd-clarify/
python3 skills/bdd-clarify/scripts/status.py /tmp/v2
```

Expected：audit exit 0；grep 只剩刻意保留的歷史敘述；`status.py` 仍是 **已答 2／n/a 1／待答 1**。

- [ ] **Step 5: Commit**

```bash
git add skills/bdd-clarify/SKILL.md
git commit -F - <<'EOF'
docs(bdd-clarify): bring the skill onto the v2 section names

WHAT: Update the flow diagram's Pass 3 routing for NFRs now having two
      homes, and repoint every section-name reference at v2's headings
WHY: The skill is the operating instructions for the step that writes
     this document; a heading it names that no longer exists sends the
     reader looking for a section the format does not have. Pass 3
     needed more than a rename — a measurable requirement that serves one
     story now belongs under that story, and only cross-cutting ones stay
     in the top-level section, so the routing had to say which.
HOW: Each reference judged in place rather than swept: some sentences
     name the concept "functional requirement" and some name the heading,
     and only the second kind changes.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 3: `bdd-spec` 與 `bdd-plan` 跟上

**Files:**
- Modify: `skills/bdd-spec/SKILL.md`
- Modify: `skills/bdd-plan/SKILL.md`

- [ ] **Step 1: 在 `bdd-spec/SKILL.md` 的 `## 切 story` 一節開頭加入 D2 的規則**

該節現在開頭是「`prd.md` 的 FR 是平鋪的。第一件事是決定哪幾條湊成一則 story」。整段換成：

```markdown
`prd.md` 已經把 FR 分在 `### US-<n>` 底下了，**但那是需求方視角的敘事分組，不是
交付切片**。

**你得自己重新分組，不必解釋為什麼跟 PRD 不一樣。** 但如果你照抄 PRD 的分組，
要說得出為什麼那個分組剛好也是好的交付邊界——「PRD 就是這樣分的」不是理由。

理由是這兩種分組的判準不同：PRD 的 US 是「誰在什麼情境下想做什麼」，你的 story 是
「哪幾條 FR 湊成一次可獨立驗收的交付」。前者可以是一個很大的願望，後者必須做得完。

一則 story ＝ 一個 `.feature` 檔 ＝ 一次可獨立驗收的交付。
```

其餘（沿規則切、INVEST、`## Stories` 的格式、MUST NOT 切出 Scope Out 的東西）**不動**。

- [ ] **Step 2: 全檔掃 `bdd-spec` 的章節名引用**

```bash
grep -n "## Functional Requirements\|## Actors" skills/bdd-spec/
```

逐處改成 v2 的名稱（`## Personas`、`## User Stories`）。

- [ ] **Step 3: 改 `bdd-plan/SKILL.md` 第 63 行**

```
- 其餘依 `prd.md` 的 Functional Requirements 切，一張票一組相關的 `@example-N.M`。
+ 其餘依 `prd.md` 的 FR 切，一張票一組相關的 `@example-N.M`。
```

PLAN 的職責與它讀的檔（`spec.md` ＋ `.feature`）都不變，這是純引用更新。

- [ ] **Step 4: 驗證**

```bash
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-plan; echo "exit=$?"
claude plugin validate .; echo "exit=$?"
grep -rn "## Functional Requirements\|## Actors" skills/ --include="*.md" --include="*.py"
python3 - <<'EOF'
import re, pathlib, sys
bad = []
for md in pathlib.Path('.').rglob('*.md'):
    if '.git' in md.parts or 'runs' in md.parts: continue
    text = re.sub(r'^(```|````).*?^\1', '', md.read_text(encoding='utf-8'), flags=re.S|re.M)
    for link in re.findall(r'\]\((\.[^)]+)\)', text):
        if not (md.parent / link.split('#')[0]).exists(): bad.append(f"{md}: {link}")
print('\n'.join(bad) or 'all links resolve'); sys.exit(1 if bad else 0)
EOF
```

Expected：兩個 audit exit 0；`plugin validate` exit 0（`--strict` 仍有 CLAUDE.md 那則已知警告）；grep 只剩刻意保留的歷史敘述；連結全解得到。

- [ ] **Step 5: Commit**

```bash
git add skills/bdd-spec/SKILL.md skills/bdd-plan/SKILL.md
git commit -F - <<'EOF'
docs(bdd-spec): say that the PRD's story grouping is a clue, not a constraint

WHAT: Rewrite the opening of the story-slicing section now that prd.md
      arrives pre-grouped, and repoint bdd-plan's one stale section name
WHY: Requirements now sit under user stories in the PRD, and a grouping
     that already exists is one SPEC will follow without deciding to.
     The two groupings answer different questions — the PRD's asks who
     wants what in which situation, SPEC's asks which requirements make
     one independently verifiable delivery — and a wish can be larger
     than anything finishable. Without the rule the distinction survives
     only as long as whoever reads the document happens to notice it.
HOW: Copying the PRD's grouping stays allowed; what changes is that it
     now needs a reason of its own, and "the PRD grouped it that way" is
     explicitly not one. Everything else in the section — slicing along
     rules, INVEST, the Stories format check_spec.py parses, the ban on
     slicing out what Scope excluded — is untouched.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```
