# `prd.md` 格式 v3 ＋ `example-mapping` Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 讓 `prd.md` 長得像一場 Example Mapping 的桌面——story 底下三個分組 `#### FR` / `#### NFR` / `#### Q`，例子巢狀在它的規則底下並帶 Given/When/Then——並把 Example Mapping 從散在三個 skill 的引用收攏成一個 skill。

**Architecture:** 契約、canonical fixture、解析器仍然是同一個原子單位（Task 1）。新 skill 是純新增（Task 2）。其餘三個 task 只是把重複的段落改成指向新 skill，以及更新文件。

**Tech Stack:** Markdown（skill 與 reference）、Python 3 標準庫（`check_spec.py`，無外部依賴）

**Spec:** [`docs/superpowers/specs/2026-09-11-prd-format-v3-design.md`](../specs/2026-09-11-prd-format-v3-design.md)

## Global Constraints

- **不引進測試框架。** repo 沒有測試基礎設施，`pytest` 沒安裝。驗證是「建 fixture → 跑腳本 → 驗 exit code 與 stdout」。
- **`check_spec.py` 只用標準庫。**
- **fail loud 不可退化。** 找不到輸入時必須非零離開並指名缺什麼。
- **`skills/bdd-clarify/scripts/status.py` 不得改動。** `## Open Questions` 的表在 v3 只有編號寫法從 `Q1` 變成 `Q-1`，而它不解讀那一欄的格式（`cells[0] in ("Q", "---")` 只擋表頭）。改到它就是改壞了。
- **`.feature` 的 tag 一個都不改。** `@rule-<n> ↔ FR-<n>`、`@example-<n>.<m> ↔ AC-<n>.<m>`——兩份產物各用母語，接縫寫明對應。`check_spec.py` 第 209 行的 `@example-` 正則不動。
- **FR 與 NFR 的編號全域唯一**，不是每則 story 各自從 1 起算。
- **AC 的 `<n>` 必須等於它所屬 FR 的編號。**
- 沿用不變：來源標記三選一（`PRD §x`／`Q-<n>`／`推論`）、`## Open Questions` 的五欄表與 `面向` 值域、EARS 引用 `docs/sdd.md` 不重述。
- 散文繁體中文，commit 訊息英文、`WHAT:`／`WHY:`／`HOW:` 三段。
- 每個 task 結束時工作區乾淨，`audit_skill.py` 對動過的 skill exit 0（跑過 python 要清掉 `__pycache__`，它會擋住 S5）。

## File Structure

| 檔案 | 責任 | Task |
| --- | --- | --- |
| `skills/bdd-clarify/references/prd-format.md` | 格式契約：十節順序、story 三分組、AC 帶 G/W/T、接縫對應表 | 1 |
| `skills/bdd-clarify/examples/minimal-prd.md` | **canonical 解析對象**，兩支腳本唯一的驗證基準 | 1 |
| `skills/bdd-spec/scripts/check_spec.py` | `frs_in_prd()` 改抓 `**FR-<n>**`／`AC-<n>.<m>`；新增 Q 懸空檢查 | 1 |
| `skills/example-mapping/SKILL.md` | **新增**：四色卡的操作定義、問題順序、四條診斷、就緒判定 | 2 |
| `skills/bdd-clarify/SKILL.md` | Pass 2 改成呼叫 example-mapping；就緒判定搬走 | 3 |
| `skills/clarify-loop/SKILL.md` | 刪 25 分鐘那段，指向 example-mapping | 4 |
| `skills/story-splitting/SKILL.md` | 刪三個訊號表，指向 example-mapping | 4 |
| `skills/bdd-spec/SKILL.md` | 第 232 行的 `EX-2.2`；寫明接縫兩側對應 | 4 |
| `docs/ai-sdlc.md`、`PLAN.md`、`benchmark/skeleton/go/README.md` | Discovery 那列、skill 表、EX 對應表 | 5 |

---

### Task 1: 契約、fixture、解析器（原子）

三者一起改。**不可拆**——分開改會留下「`check_spec.py` 找不到 `#### FR-<n>` 標題、回傳空 dict、而 Task 1 的其餘檢查一片綠」的中間狀態。

**Files:**
- Modify: `skills/bdd-clarify/references/prd-format.md`
- Modify: `skills/bdd-clarify/examples/minimal-prd.md`
- Modify: `skills/bdd-spec/scripts/check_spec.py`

**Interfaces:**
- Produces: `### US-<n>` → `#### FR` → `**FR-<n>**` → `- **AC-<n>.<m>**` → `  - Given/When/Then` 的階層，Task 2–5 的文件都要跟它一致
- Produces: `frs_in_prd(text) -> dict[str, set[str]]`（簽章不變，FR→AC）
- Produces: `qs_in_stories(text) -> set[str]`、`qs_in_table(text) -> set[str]`（新）

- [ ] **Step 1: 記錄改動前的基準**

```bash
rm -rf /tmp/v3 && mkdir -p /tmp/v3/specs/2026-09-09-x
cp skills/bdd-clarify/examples/minimal-prd.md /tmp/v3/specs/2026-09-09-x/prd.md
python3 skills/bdd-clarify/scripts/status.py /tmp/v3
python3 -c "
import sys; sys.path.insert(0,'skills/bdd-spec/scripts')
import check_spec, io
print({k: sorted(v) for k,v in sorted(check_spec.frs_in_prd(io.open('skills/bdd-clarify/examples/minimal-prd.md',encoding='utf-8').read()).items())})
"
```

Expected：`status.py` 印 **已答 2／n/a 1／待答 1**，覆蓋 `邊界`＋`降級`；第二條印 `{'1': ['1.1', '1.2', '1.3'], '2': ['2.1']}`。**抄下來——Step 7 要求兩者逐字相同。**

- [ ] **Step 2: 遷移 `examples/minimal-prd.md` 到 v3**

整份換成（逐字，這是兩支腳本唯一的 canonical 解析對象）：

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

明確排除的能力：
- 月租費線上收款 ← Q-1（範圍決定：帳務接在社區總帳上）
- 機車計費 ← Q-3

接受的風險：
- 不做月租證掛失名單管理，因此擋不住已作廢的月租證繼續感應出入 ← 推論
- 機車格沒有柵欄，日後要收費是另一次發包 ← Q-3

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

### US-1  身為 P-1（訪客），我想在出場前知道要付多少，以便付完就能開走 ← 推論

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

### US-2  身為 P-2（月租戶），我想刷卡就直接出去，以便不用每次都停下來付款 ← 推論

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

## Assumptions / Constraints

- 沒有車牌辨識，去年區權會否決 ← PRD §5（限制：問過，答案是不准）
````

**四個不可破壞的東西**：`## Open Questions` 的五欄表與表頭那格的 `Q`（`status.py` 靠它們）；`**FR-<n>**` 的粗體寫法；`- **AC-<n>.<m>**` 的層級；`#### FR` / `#### NFR` / `#### Q` 三個分組標題。

**v2 刪掉了什麼**：原本給 PM 讀、不含數字的那層 AC（`AC-1.1 訪客在繳費機看得到金額與停留時長`、`AC-1.2 付款完成前柵欄不開`、`AC-2.1 月租戶出場不需經過繳費機`）。它由 story 那一句承擔。

**修掉的既有缺陷**：Q-2 的 `擋住` 從 `FR-11（緊急車輛進場）` 改成 `FR-2`——這份 fixture 裡從來沒有 FR-11。不修的話沒有任何一則 story 有紅卡指標可放，v3 的新結構在 canonical fixture 裡一次都不會被示範。

- [ ] **Step 3: 改 `check_spec.py` 的 `frs_in_prd()`**

整個函式換成：

```python
def frs_in_prd(text: str) -> dict[str, set[str]]:
    """prd.md 的 `## User Stories` 段：FR 編號 -> AC 編號集合。

    v3 把 FR 從標題改成 `#### FR` 分組底下的粗體項目，AC 則是掛在它底下的
    清單項。理由是 AC 還要再往下掛 Given/When/Then 三行，用標題會走到 H6。

    來源是 prd.md 而不是 spec.md：例子的定版編號誕生在 CLARIFY，
    spec.md 只記哪些 FR 湊成一則 story，不重述例子。

    FR 編號全域唯一，不是每則 story 各自從 1 起算——否則 AC-1.1 會指向
    兩個地方，`.feature` 的 @example-1.1 就失去意義。

    碰到任何標題就把 current 清掉。`#### FR`／`#### NFR`／`#### Q` 三個分組
    與 `### US-` 都是 FR 的祖先或兄弟，它們底下的散文不屬於上一條 FR——
    `#### Q` 尤其要緊，紅卡的敘述提到某個 AC 編號時，不清掉 current 就會
    把那個例子靜默地記到上一則 story 最後一條 FR 頭上。
    """
    section = re.search(
        r"^## User Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return {}
    out: dict[str, set[str]] = {}
    current = None
    for line in section.group(1).splitlines():
        if line.startswith("#"):
            current = None
            continue
        fr = re.match(r"^\*\*FR-(\d+)\*\*", line)
        if fr:
            current = fr.group(1)
            out.setdefault(current, set())
            continue
        if current:
            for ac in re.findall(r"\bAC-(\d+\.\d+)\b", line):
                out[current].add(ac)
    return out
```

- [ ] **Step 4: 新增 Q 懸空檢查的兩個解析函式**

加在 `frs_in_prd()` 之後、`stories_in_spec()` 之前：

```python
def qs_in_stories(text: str) -> set[str]:
    """`## User Stories` 段的 `#### Q` 分組底下引用的 Q 編號。

    story 底下的紅卡只是指標——狀態、面向、答案、決策史都住在
    `## Open Questions`。這裡只收編號，好跟那張表對帳。
    """
    section = re.search(
        r"^## User Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return set()
    out: set[str] = set()
    in_q = False
    for line in section.group(1).splitlines():
        if line.startswith("#"):
            in_q = re.match(r"^#### Q\s*$", line) is not None
            continue
        if in_q:
            out.update(re.findall(r"\*\*Q-(\d+)\*\*", line))
    return out


def qs_in_table(text: str) -> set[str]:
    """`## Open Questions` 表第一欄的 Q 編號。

    只讀那一節：`## Document Overview` 底下的版本修訂歷史也是五欄表，
    對整份掃會把它的每一列都當成一題。
    """
    section = re.search(
        r"^## Open Questions\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return set()
    out: set[str] = set()
    for line in section.group(1).splitlines():
        if not line.strip().startswith("|"):
            continue
        first = line.strip().strip("|").split("|")[0].strip()
        m = re.fullmatch(r"Q-(\d+)", first)
        if m:
            out.add(m.group(1))
    return out
```

- [ ] **Step 5: 把 Q 懸空檢查接進主流程**

在 barren 檢查那一段（現在是 `barren = sorted(...)` 那幾行）之後、`spec_path = prd_path.parent / "spec.md"` 之前插入：

```python
        # story 底下的紅卡是指標，指到 `## Open Questions` 沒有的題號就是
        # 斷掉的引用——通常是題目被刪了卻沒回頭清 story，或編號打錯。
        # v2 的結構做不到這個檢查（紅卡不掛在 story 上），v3 才有。
        dangling = sorted(qs_in_stories(ptext) - qs_in_table(ptext), key=int)
        if dangling:
            print(f"✗ {prd_path.parent.name} 的 story 引用了 `## Open Questions` "
                  f"裡沒有的問題：{', '.join('Q-' + d for d in dangling)}")
            problems += 1

        # AC 的 <n> 必須等於它掛在底下的那條 FR。巢狀讓這件事在寫的時候
        # 不容易弄錯，但打字錯不會被結構擋下來——`AC-9.1` 打在 `**FR-1**`
        # 底下會被記成 FR-1 的例子，然後在下面的覆蓋比對裡變成一則
        # 「@example-9.1 不見了」的訊息，指著 .feature 說它漏寫，而錯在 prd.md。
        mismatched = sorted(
            (fr, ac) for fr, acs in frs.items() for ac in acs
            if ac.split(".")[0] != fr)
        if mismatched:
            print(f"✗ {prd_path.parent.name} 這些 AC 的編號對不上它所屬的 FR："
                  f"{', '.join(f'AC-{ac} 掛在 FR-{fr} 底下' for fr, ac in mismatched)}")
            problems += 1
```

同時把 barren 檢查的訊息從 `這些 FR 沒有任何 EX` 改成 `這些 FR 沒有任何 AC`，並把該段上方註解裡的「完全沒有例子的 FR」保留（那句仍然成立），只把 `EX` 字樣換掉。

檔頭 docstring 第 16 行 `prd.md 裡有沒有 FR 完全沒掛任何例子` 保留（沒有提 EX），第 50–51 行的 `EX-1.1` 已在 Step 3 換掉。

- [ ] **Step 6: 改寫 `prd-format.md` 的四節**

**(a) 章節順序表的 `## User Stories` 那一列**換成：

```
| `## User Stories` | `### US-<n>`，每則底下三個分組：`#### FR`（AC 掛在各自的 FR 底下）、`#### NFR`、`#### Q` |
```

**(b) `## User Stories：強制指向 Persona，FR 掛在底下` 整節**（第 123 行到 `## PRD 的 story 分組不是交付切片` 之前）換成：

`````markdown
## User Stories：三個分組，AC 巢狀在它的 FR 底下

````markdown
### US-1  身為 P-1（訪客），我想在出場前知道要付多少，以便準備好再開到閘門 ← 推論

#### FR

**FR-1**  When 訪客車出場, the system shall 依停留時長計費，前 30 分鐘免費。 ← PRD §3
- **AC-1.1**  停 29 分 → 收 0 元 ← PRD §3
  - Given 一台訪客車已入場並抽了紙票
  - When 它在入場後 29 分鐘於繳費機結帳
  - Then 應收金額是 0 元

#### NFR

**NFR-1**  出場柵欄自刷卡到開啟不超過 2 秒 ← Q-12

#### Q

**Q-7**  跨午夜的停車怎麼算 ← 擋住 FR-1
````

**這個排版對應 Example Mapping 的桌面**：🟨 story、🟦 FR、🟩 AC、🟥 Q。怎麼跑一場
map、四條診斷怎麼讀 → `example-mapping` skill。

**為什麼 FR／NFR／Q 平輩，AC 不平輩。** 那張桌子是一棵樹：綠卡擺在它所屬的藍卡
正下方。FR、NFR、Q 都是 story 層級的東西，並排成三個分組是對的；AC 屬於某一條 FR，
平輩化會讓讀的人靠編號在腦裡把它們接回去。巢狀讓「`AC-1.1` 屬於 `FR-1`」是**結構上
不可能弄錯**的。

**為什麼用粗體項目而不是標題。** AC 底下還要掛 Given/When/Then，用標題會走到 H6。

MUST: `身為` 後面必須是一個 `P-<n>`。指到不存在的 persona＝角色沒被定義，是洞。

**已知弱點：這一條沒有東西在檢查它。**

MUST: **FR 與 NFR 的編號全域唯一。** `FR-1` 在 US-1 底下、`FR-4` 在 US-2 底下都
可以，但不得有兩個 `FR-1`——否則 `AC-1.1` 指向兩個地方，`.feature` 的
`@example-1.1` 失去意義。

MUST: 每條 AC 底下的 Given／When／Then 三行齊全。**Given 寫不出來就是這個例子的
前置狀態還沒被問過——那是一張紅卡，不是一句可以省略的話。**

v2 的 AC 是「給 PM 讀、不含數字」的那一層，v3 拿掉了：一層寫不清楚的東西，兩層
只會讓它在兩個地方各寫一半。PM 對焦讀 story 那一句。

**NFR 有兩個家，判準是跨不跨 story**：只服務某一則 story 的進那則的 `#### NFR`；
跨全部的留在 `## Non-Functional Requirements`。搬家時編號不變。

### `#### Q` 只放指標

紅卡在 story 底下只寫 **編號 ＋ 問題一句話 ＋ 擋住誰**；狀態、面向、答案、決策史
全部住在 `## Open Questions`。

MUST: `#### Q` 只列**狀態為「待答」**的問題。已答的不是紅卡了——它的答案已經長成
某條 FR 或 AC，留在這裡會讓地圖上的紅卡數量永遠不歸零。

MUST: 這裡引用的每個 `Q-<n>` 都要存在於 `## Open Questions` 表。**這一條
`check_spec.py` 會檢查。**

**不是每張紅卡都屬於某則 story。** 範圍題與技術題在 Pass 1、Pass 3 產生，那時候
還沒有任何 story——它們只住在 `## Open Questions`，不進任何 `#### Q`。

### 一條 FR 服務多則 story

巢狀強迫一對多，共用的 FR 只能住一個地方。**住主要的那則，其他宣告依賴**：

```markdown
### US-2  身為 P-2（主委），我想知道訪客格有沒有被外人占走 ← Q-2

**另外依賴**：FR-3（定義在 US-1 底下）
```

**已知弱點：這行宣告沒有東西在檢查它。**
`````

**(c) `## 編號規則：FR 掛在 story 底下，AC／EX 掛在 FR 底下` 的標題與內文**——標題
改成 `## 編號規則：AC 的編號跟著它的 FR`，內文裡所有 `EX-<n>.<m>` 改成
`AC-<n>.<m>`、`### US-1` 底下的範例改成 v3 的粗體寫法，並把「編號在這裡定版」那段
之後補上接縫對應表：

```markdown
### 接縫：兩份產物各用母語

| `prd.md`（PRD 詞彙） | `.feature`（Gherkin 詞彙） |
| --- | --- |
| `FR-<n>` | `@rule-<n>` |
| `AC-<n>.<m>` | `@example-<n>.<m>` |

**兩個名字在這裡不是缺陷。** 同一份產物裡有兩個名字才是缺陷；兩份產物各自用母語、
接縫上寫明對應，那是設計。`.feature` 是 Gherkin 的產物，`@example` 是 Gherkin 的字。
```

**(d) `## AC 與 Examples：兩層讀者不同` 整節**（第 270 行到 `## Non-Functional Requirements 與 Constraints` 之前）換成：

`````markdown
## AC 就是綠卡

v2 有兩層：`AC` 給 PM（不含數字）、`EX` 給 SPEC（含數字）。v3 只有一層。

`AC-<n>.<m>` 同時是驗收條件與具體例子：**一行摘要 ＋ 三行 Given/When/Then**。

````markdown
**FR-1**  When 訪客車出場, the system shall 依停留時長計費，前 30 分鐘免費。 ← PRD §3
- **AC-1.1**  停 29 分 → 收 0 元 ← PRD §3
  - Given 一台訪客車已入場並抽了紙票
  - When 它在入場後 29 分鐘於繳費機結帳
  - Then 應收金額是 0 元
````

**為什麼合併。** 單行例子（`停 29 分 → 收 0 元`）把前置狀態藏起來了：車已經入場
了嗎？計時從進柵欄算還是從抽票算？單行寫法讓這些不用回答，而 `bdd-spec` 寫
`.feature` 時**必須**回答——於是它自己編一個。那不是 SPEC 不守規矩，是 CLARIFY 沒把
話講完。

**為什麼叫 AC 不叫 Example。** 產物是 PRD，讀者帶著 PRD 的詞彙來。Example Mapping
是手法，PRD 是產物，手法的詞彙不覆蓋產物的詞彙。

**一個詞義衝突要知道**：Matt Wynne 把藍卡描述成 "a business rule or acceptance
criterion"——在 Example Mapping 的血統裡 **AC 是藍卡**；在主流 PRD／Jira 用法裡 AC
是那串 Given/When/Then。這裡用 PRD 的。四色卡與 PRD 術語的完整對照 →
`example-mapping` skill。

**已知弱點**：G/W/T 與 `.feature` 的 `Scenario` 有字面重複。`prd.md` 的是**談定的**，
`.feature` 的是**可執行的**（還要套封閉步驟文法與 `Rule:` 結構）。沒有東西檢查兩者
一致——`check_spec.py` 只檢查 tag 對得上編號。
`````

**(e) 來源標記的適用範圍**那段裡的 `## User Stories 底下的每條 FR／AC／EX`——把
`／EX` 刪掉。

- [ ] **Step 7: 回歸——兩個結果必須跟 Step 1 逐字相同**

```bash
rm -rf /tmp/v3 && mkdir -p /tmp/v3/specs/2026-09-09-x
cp skills/bdd-clarify/examples/minimal-prd.md /tmp/v3/specs/2026-09-09-x/prd.md
python3 skills/bdd-clarify/scripts/status.py /tmp/v3
python3 -c "
import sys; sys.path.insert(0,'skills/bdd-spec/scripts')
import check_spec, io
t = io.open('skills/bdd-clarify/examples/minimal-prd.md',encoding='utf-8').read()
print('frs :', {k: sorted(v) for k,v in sorted(check_spec.frs_in_prd(t).items())})
print('q_st:', sorted(check_spec.qs_in_stories(t)))
print('q_tb:', sorted(check_spec.qs_in_table(t)))
"
```

Expected：`status.py` 仍是 **已答 2／n/a 1／待答 1**、覆蓋 `邊界`＋`降級`（`## Open Questions` 只換了編號寫法，數字變了就是改壞了）；`frs : {'1': ['1.1', '1.2', '1.3'], '2': ['2.1']}`（與 Step 1 逐字相同）；`q_st: ['2']`；`q_tb: ['1', '2', '3', '4']`。

- [ ] **Step 8: 端到端——健康路徑必須綠**

```bash
mkdir -p /tmp/v3/features
printf '# SPEC\n\n## Stories\n\n### visitor-billing\n涵蓋 FR-1。\n\n### monthly-pass\n涵蓋 FR-2。\n' > /tmp/v3/specs/2026-09-09-x/spec.md
cat > /tmp/v3/features/visitor-billing.feature <<'EOF'
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
cat > /tmp/v3/features/monthly-pass.feature <<'EOF'
@ready
Feature: 月租戶出場

  Rule: 月租戶不計費

    @example-2.1
    Scenario: 名單內且未到期
      Given 一台月租車在名單內且未到期
      When 它刷卡出場
      Then 柵欄開啟且應收金額是 0 元
EOF
python3 skills/bdd-spec/scripts/check_spec.py /tmp/v3; echo "exit=$?"
```

Expected：`exit=0`，兩則 story 都印 `✓`、`3/3 例子`、`1/1 例子`、`全部通過`。**tag 一個都沒改，證明接縫沒斷。**

- [ ] **Step 9: 三個新檢查都必須真的會失敗**

```bash
# (a) 懸空的紅卡指標
rm -rf /tmp/v3a && cp -r /tmp/v3 /tmp/v3a
sed -i 's/^\*\*Q-2\*\*  消防法規對閘門的要求 ← 擋住 FR-2/**Q-99**  憑空的問題 ← 擋住 FR-2/' /tmp/v3a/specs/2026-09-09-x/prd.md
python3 skills/bdd-spec/scripts/check_spec.py /tmp/v3a 2>&1 | grep -E "Q-99|通過"; echo "  exit=${PIPESTATUS[0]}"

# (b) FR 沒有任何 AC
rm -rf /tmp/v3b && cp -r /tmp/v3 /tmp/v3b
python3 - <<'PY'
import re, io
p = '/tmp/v3b/specs/2026-09-09-x/prd.md'
t = io.open(p, encoding='utf-8').read()
t = re.sub(r"\n- \*\*AC-2\.1\*\*.*?(?=\n\n#### NFR)", "", t, flags=re.S)
io.open(p, 'w', encoding='utf-8').write(t)
PY
python3 skills/bdd-spec/scripts/check_spec.py /tmp/v3b 2>&1 | grep -E "沒有任何 AC|通過"; echo "  exit=${PIPESTATUS[0]}"

# (c) 真實觸發：FR 留在 v2 的標題寫法（對 v3 解析器而言是「不存在」而非「錯」）
rm -rf /tmp/v3c && cp -r /tmp/v3 /tmp/v3c
sed -i 's/^\*\*FR-2\*\*  When/#### FR-2  When/' /tmp/v3c/specs/2026-09-09-x/prd.md
python3 skills/bdd-spec/scripts/check_spec.py /tmp/v3c 2>&1 | grep -E "FR-2|通過"; echo "  exit=${PIPESTATUS[0]}"
```

Expected：三個都 `exit=1`。(a) 指名 `Q-99`；(b) 指名 `FR-2 沒有任何 AC`；(c) 指名 `spec.md 點名了 prd.md 裡沒有的 FR：FR-2`——**不是印「全部通過」**。第三個是本次改動的真實風險：留在舊寫法的 FR 對解析器是缺席而非錯誤。

- [ ] **Step 10: fail loud 沒退化 ＋ 稽核 ＋ status.py 未被動過**

```bash
python3 skills/bdd-spec/scripts/check_spec.py /tmp/nonexistent-v3; echo "exit=$?"
python3 skills/bdd-clarify/scripts/status.py /tmp/nonexistent-v3; echo "exit=$?"
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify; echo "exit=$?"
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
git diff --stat -- skills/bdd-clarify/scripts/status.py
```

Expected：兩次 `exit=1` 並指名缺什麼；兩個 audit exit 0；`status.py` 的 diff **為空**。

- [ ] **Step 11: Commit**

```bash
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
git add skills/bdd-clarify/references/prd-format.md \
        skills/bdd-clarify/examples/minimal-prd.md \
        skills/bdd-spec/scripts/check_spec.py
git commit -F - <<'EOF'
feat(bdd-clarify): prd.md v3 — the document takes the shape of the map

WHAT: Give each story three groups — requirements, quality attributes and
      open questions — nest each example under the rule it illustrates,
      and have examples carry Given/When/Then
WHY: v2 moved requirements under the story they serve, but a story's open
     questions still sat half a document away from the rules they block,
     so judging one story's readiness meant working across two places.
     And a one-line example hid its own preconditions: "29 minutes
     parked" never said whether the car had entered or when the clock
     started, which the feature file has to answer and therefore invents.
     That is not the specification step misbehaving; it is clarification
     stopping one sentence early.
HOW: Contract, fixture and parser change together because splitting them
     leaves a parser that finds no requirements and reports success. The
     parser now clears its current requirement at every heading, since
     all three group headings are ancestors or siblings of a requirement
     and the question group in particular carries text that names example
     numbers.

     Two names survive on purpose. The document speaks PRD and the
     feature file speaks Gherkin, and the correspondence between them is
     written on both sides rather than resolved by renaming one — which
     keeps roughly ninety existing tags untouched. Two names are a defect
     inside one artifact; across two, stated, they are the design.

     v3 makes one previously unenforceable rule checkable: a story's red
     cards are pointers, so a pointer to a question the table does not
     contain is now an error. The fixture also loses a defect it carried
     from v2 — its only unanswered question blocked a requirement that
     never existed, which would have left the new structure undemonstrated
     in the one document both scripts are verified against.

Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 2: 新增 `example-mapping` skill

**Files:**
- Create: `skills/example-mapping/SKILL.md`

**Interfaces:**
- Consumes: Task 1 的 `prd.md` v3 階層
- Produces: skill 名稱 `example-mapping`，Task 3、4、5 的文件都指向它

**單檔，不開子目錄。** `audit_skill.py` 的 S4（空目錄）與 S5（沒被指路的檔案）都只
在有子目錄時才會踩到；`story-splitting`（198 行）與 `clarify-loop`（241 行）都是單檔。

- [ ] **Step 1: 建立 `skills/example-mapping/SKILL.md`**

逐字：

`````markdown
---
name: example-mapping
description: >
  用四色卡把一則 story 攤成規則與例子——黃卡 story、藍卡規則、綠卡例子、紅卡問題，
  攤完跑四條診斷並提出就緒判定。可讀 `prd.md`，也可單獨對一段描述跑。
  觸發詞：「攤成四色卡」「example mapping」「跑一次 mapping」「這則 story 就緒了嗎」
  「規則跟例子攤開來看」「地圖長什麼樣」「紅卡還剩幾張」「這則可以開工了嗎」。
  English: run example mapping, map this story, lay out the rules and examples,
  is this story ready, four-colour cards, thumb vote.
---

# EXAMPLE MAPPING — 把一則 story 攤成四色卡

輸入是一則 story（或一份 `prd.md`），輸出是**攤開的地圖 ＋ 四條診斷 ＋ 一個就緒
問句**。

出處與原始定義 → [`docs/bdd.md`](../../docs/bdd.md)。那份只轉述 Cucumber 怎麼說，
**不講怎麼在這裡做**——怎麼做是這一份的事。

## 四色卡在這裡叫什麼

| 卡 | Example Mapping | `prd.md` |
| --- | --- | --- |
| 🟨 | story | `US-<n>` |
| 🟦 | rule | `FR-<n>` |
| 🟩 | example | `AC-<n>.<m>` |
| 🟥 | question | `Q-<n>` |

**一個詞義衝突要先講清楚。** Matt Wynne 把藍卡描述成 "a business rule or acceptance
criterion"——在 Example Mapping 的血統裡 **AC 是藍卡**。在主流 PRD／Jira 用法裡，AC
是那串 Given/When/Then。

`prd.md` 用的是後者，因為**產物是 PRD，讀者帶著 PRD 的詞彙來**。手法的詞彙不覆蓋
產物的詞彙。讀過 Cucumber 那篇文章的人看到這張表會覺得綠卡對錯了——沒有，是同一個
詞在兩個社群有兩個意思。

## 使用時機

- `bdd-clarify` 的 Pass 2 要把需求攤成規則與例子
- 想知道一則 story 就緒了沒
- 手上有一份 `prd.md`，想看它的形狀而不是讀它的字
- 單獨拿一則 story 來攤，不走 BDD 流程

## Skill Boundaries

- 要把紅卡逐一問到收斂 → 改用 `clarify-loop`，本 skill 只指出還剩幾張
- 要把 story 切小 → 改用 `story-splitting`
- 要把例子寫成 Gherkin → 改用 `bdd-spec`
- 本 skill **不做決定**：它攤牌、診斷、提問，就緒與否由使用者判

## 怎麼攤

從 story 開始，之後**規則與例子哪個先都可以**。Example Mapping 的原始建議給了一條
判準：**談不攏規則時就寫例子，寫到規則自己浮出來**。

反過來也成立：規則講得很順但舉不出例子，那條規則多半還停在想像階段。

MUST NOT: 先把規則寫滿再回頭補例子。那個順序會得到一批「聽起來對、但沒有人驗證過
它在具體情境下成立」的規則。

紅卡**隨時**可以放——攤的過程中冒出來的問題就是紅卡，不必等到最後才整理。

## 攤開的地圖長這樣

終端輸出，不是檔案：

```
US-1  身為 P-1（訪客），我想在出場前知道要付多少
│
├─🟦 FR-1  前 30 分鐘免費
│   ├─🟩 AC-1.1  停 29 分 → 收 0 元
│   ├─🟩 AC-1.2  停剛好 30 分 → 收 0 元
│   └─🟩 AC-1.3  停 31 分 → 收 30 元
├─🟦 FR-2  月租戶不計費
│   └─🟩 AC-2.1  名單內且未到期 → 開柵欄，金額 0
└─🟥 Q-2  消防法規對閘門的要求            ← 擋住 FR-2

診斷：🟦 2  🟩 4  🟥 1
  ⚠ 還有 1 張紅卡 —— 未就緒
  ✓ 每條 FR 都有例子
  ✓ 沒有規則堆了過多例子
```

## 四條診斷

**其中兩條在這裡改了意義，照原文抄會給出錯的建議。**

| 卡 | 原始意思 | 這裡 | 傳給 |
| --- | --- | --- | --- |
| 🟥 還有紅卡 | 不確定性高，未就緒 | 同左 | `clarify-loop` |
| 🟦 一則 story 掛太多 FR | **story 太大，該切** | **不阻塞** ← 見下 | 只回報 |
| 🟩 AC 堆在同一條 FR | 這條規則該拆 | 同左 | `bdd-clarify` |
| ⏱ 25 分鐘 map 不完 | 太大或太不確定 | **沒有時鐘** ← 見下 | `clarify-loop` |

### 🟦 為什麼不再是阻塞

原始啟發式假設**黃卡上的 story 就是交付切片**。這個流程裡它不是——`bdd-spec` 的
規則寫著：

> `prd.md` 已經把 FR 分在 `### US-<n>` 底下了，**但那是需求方視角的敘事分組，不是
> 交付切片**。你得自己重新分組。

所以「US-2 掛了 9 條 FR」不代表未就緒——SPEC 本來就要重新分組。它代表的是這則敘事
story 可能包了不只一個使用者目標。

MUST NOT: 因為藍卡多就去呼叫 `story-splitting` 改 `prd.md` 的分組。那會讓本 skill
去動 PRD 的 story 邊界，正好違反上面那條規則。**回報，不動手。**

### ⏱ 沒有時鐘，那個框換成什麼

25 分鐘的框存在是因為人會累，而且 map 不完代表太大。Agent 不會累，所以**可轉移的
是尺寸訊號，不是那個鐘**。

操作等價物：**攤的過程中紅卡出現得比綠卡快**——一邊放例子、一邊冒出新問題，而且冒
的比放的多。

`clarify-loop` 有同型的判準（「問了兩三輪、紅卡總數卻不下降」），撞到這一條就轉給它。

## 就緒判定 —— 票不是你投的

原始實踐用**大拇指投票**：三方各自表態。這裡沒有三方，而 `docs/ai-sdlc.md` 已經
定調：agent 取代三方對話當提問者，但「**沒有業務脈絡的最終裁量權——它問得出問題，
答不了**」。

所以流程是：攤開地圖 → 跑完四條診斷 → **給出判讀** → **問使用者**。

```
還有 1 張紅卡（Q-2 消防法規），擋住 FR-2。
我的判讀：未就緒——這條會改變出場放行的行為。

你要 (a) 先去問消防設備師，還是 (b) 認了這個風險先往下走？
```

MUST NOT: agent 自己把 `## Document Overview` 的 `狀態` 推到 `已對焦`。

**本 skill 是唯一寫那個欄位的 skill**，但寫的是使用者的決定：

| 狀態 | 什麼時候寫 |
| --- | --- |
| `待對焦` | 紅卡清空了，等需求方過目 |
| `已對焦` | 使用者在上面那個問句回答「就緒」 |

`澄清中` 與 `已凍結` 不由本 skill 寫——前者是 `bdd-clarify` 開檔時的預設，後者由
`bdd-spec` 消費過之後標記。

## 兩種模式

| 模式 | 輸入 | 輸出 |
| --- | --- | --- |
| 流程內 | `specs/<date>-<feature>/prd.md` | 地圖 ＋ 診斷 ＋ 就緒問句；使用者答完才寫 `狀態` |
| 單獨 | 一段 story 的描述 | 同樣的地圖與診斷，**不寫任何檔案** |

MUST NOT: 單獨模式產出一份 `example-map.md`。**第二個裝規則與例子的地方會跟
`prd.md` 漂移**——2026-09-09 拿掉那個產物就是為了這件事。要留下來就寫進 `prd.md`。
`````

- [ ] **Step 2: 驗證**

```bash
python3 skills/skill-rules/scripts/audit_skill.py skills/example-mapping; echo "exit=$?"
claude plugin validate .; echo "exit=$?"
python3 - <<'EOF'
import re, pathlib, sys
bad = []
for md in pathlib.Path('skills/example-mapping').rglob('*.md'):
    text = md.read_text(encoding='utf-8')
    out, fence = [], None
    for line in text.splitlines():
        m = re.match(r'^(`{3,})', line)
        if m and fence is None: fence = m.group(1); continue
        if m and line.startswith(fence): fence = None; continue
        if fence is None: out.append(line)
    for link in re.findall(r'\]\((\.[^)]+)\)', '\n'.join(out)):
        if not (md.parent / link.split('#')[0]).exists(): bad.append(f"{md}: {link}")
print('\n'.join(bad) or 'all links resolve'); sys.exit(1 if bad else 0)
EOF
```

Expected：audit exit 0；`plugin validate` exit 0（只有 root `CLAUDE.md` 那則已知警告）；連結全解得到（`../../docs/bdd.md`、`../../docs/ai-sdlc.md` 從 `skills/example-mapping/` 出發要解得開）。

- [ ] **Step 3: Commit**

```bash
git add skills/example-mapping/SKILL.md
git commit -F - <<'EOF'
feat(example-mapping): give the technique a skill of its own

WHAT: Add a skill that lays a story out as four kinds of card, runs the
      four diagnostics, and puts the readiness question to the user
WHY: Example Mapping had definitions but no home. The transcription doc
     states what Cucumber says and is barred from saying how to do
     anything here, so the operational half was scattered across three
     skills — the twenty-five minute heuristic in one, the too-many-rules
     signal in another, the cards themselves nowhere — with no authority
     among them.
HOW: Two of the four diagnostics change meaning in this pipeline and the
     skill says so rather than quoting the original. Too many rules on
     one story no longer blocks, because a story here is narrative
     grouping and the specification step reslices anyway; acting on that
     signal would have this skill edit boundaries another skill owns. The
     twenty-five minute box has no equivalent for an agent that does not
     tire, so the transferable signal is questions arriving faster than
     examples.

     The vocabulary clash is stated up front: the original calls the rule
     card an acceptance criterion, mainstream PRD practice calls the
     Given/When/Then one that, and this document uses the second because
     its reader arrives with a PRD in hand. Someone who has read the
     Cucumber article will otherwise think the table is wrong.

     Readiness stays the user's. The skill is the only one that writes
     the document's status field, and it writes what the user decided.

Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 3: `bdd-clarify/SKILL.md` — Pass 2 委派，就緒判定搬走

**Files:**
- Modify: `skills/bdd-clarify/SKILL.md`

**Interfaces:**
- Consumes: Task 1 的 v3 階層、Task 2 的 `example-mapping` skill

- [ ] **Step 1: 先讀，再改**

```bash
grep -n "^## \|^### " skills/bdd-clarify/SKILL.md
sed -n '329,400p' skills/bdd-clarify/SKILL.md
```

`## Pass 2 · 深度` 在第 198 行，底下有 `### 4. 澄清循環`、`### 5. 就緒判定 —— 提出判定，票不是你投的`、`### 6. 收尾`。**第 5 節已經寫著「票不是你投的」——那正是要搬進 `example-mapping` 的東西，它不是新規則，是既有規則換家。**

- [ ] **Step 2: `### 4. 澄清循環` 開頭加一句委派**

在該節標題之後、第一段之前插入：

```markdown
**手法是 Example Mapping** —— 四色卡怎麼對應 `prd.md`、怎麼攤、四條診斷怎麼讀
→ `example-mapping`。本節只講這個 pass 的節奏（一次一題、每輪抽規則），不重述手法。
```

- [ ] **Step 3: `### 5. 就緒判定` 整節換成指向**

整節（標題到 `### 6.` 之前）換成：

```markdown
### 5. 就緒判定 —— 呼叫 `example-mapping`

紅卡清空（或使用者決定帶著紅卡往下走）之後，跑一次 `example-mapping`：它攤開地圖、
跑四條診斷、把就緒問句交給使用者，並依使用者的答覆寫 `## Document Overview` 的
`狀態` 欄。

MUST NOT: 在這裡自己宣告就緒。**本 skill 不寫 `狀態` 欄**——那是 `example-mapping`
唯一負責的一格，而值由使用者決定。
```

- [ ] **Step 4: 全檔掃過時的引用**

```bash
grep -n "EX-\|##### AC\|##### Examples\|#### FR-" skills/bdd-clarify/SKILL.md
```

逐處改成 v3 的寫法（`**FR-<n>**`、`- **AC-<n>.<m>**`）。**逐處判讀語意再改**——有些句子講的是「例子」這個概念而不是那個編號。

- [ ] **Step 5: 驗證**

```bash
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify; echo "exit=$?"
grep -n "EX-" skills/bdd-clarify/
python3 skills/bdd-clarify/scripts/status.py /tmp/v3
```

Expected：audit exit 0；`EX-` 零筆；`status.py` 仍是 **已答 2／n/a 1／待答 1**。

- [ ] **Step 6: Commit**

```bash
git add skills/bdd-clarify/SKILL.md
git commit -F - <<'EOF'
docs(bdd-clarify): hand the mapping technique and the readiness call over

WHAT: Point Pass 2 at the example-mapping skill for the technique, and
      replace the readiness section with a call to it
WHY: The readiness section already said the vote is not yours to cast —
     that rule is not new, it is moving to the skill that now owns the
     whole readiness call, including the document status field it sets.
     Leaving a copy here would give two skills a claim on the same
     decision, and the one with the weaker claim is the one that cannot
     see the map.
HOW: The pass keeps its own rhythm — one question at a time, rules
     extracted each round — and stops describing the technique, which now
     has a single place to live.

Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 4: 三個指向更新（同形狀，一批做完）

**Files:**
- Modify: `skills/clarify-loop/SKILL.md`
- Modify: `skills/story-splitting/SKILL.md`
- Modify: `skills/bdd-spec/SKILL.md`

- [ ] **Step 1: `clarify-loop/SKILL.md` 的 25 分鐘段**

第 187–190 行現在是：

```
第四種最容易被漏掉，因為它偽裝成「還有很多問題要問」。判準來自 Example Mapping
的原始建議：一則大小合適的 story，三方應該在 **25 分鐘左右** map 完；map 不完
表示它太大或太不確定。對應到這裡：**問了兩三輪、紅卡總數卻不下降，就不是問法的
問題，是範圍的問題**。
```

換成：

```
第四種最容易被漏掉，因為它偽裝成「還有很多問題要問」。判準是**問了兩三輪、紅卡
總數卻不下降，就不是問法的問題，是範圍的問題**——這條的出處與它在 agent 情境下
為什麼長這樣 → `example-mapping` 的「四條診斷」。
```

- [ ] **Step 2: `story-splitting/SKILL.md` 的三個訊號表**

`## 2. 判斷該不該切` 一節裡，從 `來自 Example Mapping 的三個訊號` 那句起，經過三列
的表格，到 `再看還需不需要切。` 為止（表格**與它後面那一段一起**）換成：

```markdown
Example Mapping 的診斷訊號裡，跟切分有關的只有**藍卡（規則）太多**——而在這條
流程裡它**不阻塞**：`prd.md` 的 story 是需求方視角的敘事分組，SPEC 本來就要重新
切。完整的四條診斷、它們在這裡各自代表什麼、各自傳給誰 → `example-mapping`。

**一條規則掛太多例子不是切分訊號。** 例子多不代表 story 大，可能只是規則講得太粗
——先試著把那條規則拆成兩三條，再看還需不需要切。

到這裡的訊號是：**切 story 時發現一則怎麼切都太大。**
```

**表格後面那一段必須一起換掉。** 它開頭是「第三項容易誤判……」——「第三項」指的
是表格的第三列，表格一刪它就變成懸空引用。上面的替換文字已經把那段的內容吸收
進去，不要兩邊都留。

第 18 行的 `- Example Mapping 跑不完，或藍卡（規則）多到一次講不清` 保留——那是
使用時機，不是判準。

- [ ] **Step 3: `bdd-spec/SKILL.md` 的 `EX-` 與接縫對應**

第 232 行 `- EX-2.2 進度 70% 想改 60%           →    @rule-2 @example-2.2` 的
`EX-2.2` 改成 `AC-2.2`。

同一節補上接縫對應（放在該行所在的對照表之後）：

```markdown
| `prd.md`（PRD 詞彙） | `.feature`（Gherkin 詞彙） |
| --- | --- |
| `FR-<n>` | `@rule-<n>` |
| `AC-<n>.<m>` | `@example-<n>.<m>` |

**tag 不跟著 `prd.md` 改名。** `.feature` 是 Gherkin 的產物，`@example` 是 Gherkin
的字；`prd.md` 是 PRD，`AC` 是 PRD 的字。接縫寫明對應就夠了。
```

`grep -n "EX-" skills/bdd-spec/SKILL.md` 在這個檔案裡**只會命中第 232 行那一處**。
改完再跑一次應該零筆。

- [ ] **Step 4: 驗證**

```bash
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for s in clarify-loop story-splitting bdd-spec; do
  python3 skills/skill-rules/scripts/audit_skill.py skills/$s; echo "  $s exit=$?"
done
grep -rn "EX-" skills/ || echo "EX- 零筆"
grep -rn "25 分鐘" skills/ | grep -v example-mapping || echo "25 分鐘只剩 example-mapping"
grep -rc "@example-" skills/bdd-spec/examples/ | head -3
```

Expected：三個 audit exit 0；`EX-` 零筆；`25 分鐘` 只出現在 `example-mapping`；`skills/bdd-spec/examples/` 的 `@example-` 筆數**與改動前相同**（tag 一個都沒動）。

- [ ] **Step 5: Commit**

```bash
git add skills/clarify-loop/SKILL.md skills/story-splitting/SKILL.md skills/bdd-spec/SKILL.md
git commit -F - <<'EOF'
docs: point the three skills that quoted Example Mapping at the skill

WHAT: Replace the copies of the mapping heuristics in the clarify loop
      and the splitting skill with pointers, and repoint the one stale
      example id in the specification skill
WHY: Three skills each carried a fragment of the same technique and none
      was authoritative, so a heuristic could be corrected in one and
      stay wrong in the other two. The splitting skill in particular
      quoted a signal that no longer means what it says here: too many
      rules on a story does not block, because the grouping it measures
      is narrative and gets resliced anyway.
HOW: Each skill keeps the sentence that tells you when you have arrived
     at it and drops the reasoning behind the signal, which now lives in
     one place. The Gherkin tags are deliberately untouched — the
     correspondence between the two vocabularies is now written on both
     sides of the seam instead.

Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```

---

### Task 5: 文件與 roadmap

**Files:**
- Modify: `docs/ai-sdlc.md`
- Modify: `PLAN.md`
- Modify: `benchmark/skeleton/go/README.md`

- [ ] **Step 1: `docs/ai-sdlc.md` 的 Discovery 那列**

該表第 98 行的 `**Discovery**` 列，在「agent 能做」那格結尾補上：

```
。手法是 Example Mapping，四色卡怎麼對應 `prd.md`、四條診斷在這裡各自代表什麼 → `example-mapping` skill
```

**不要動「agent 不能做」那格**——「取代三方對話／沒有最終裁量權」正是 `example-mapping`
的就緒判定歸使用者的依據，它必須留在這裡。

- [ ] **Step 2: `PLAN.md` 的 skill 表**

在 `| | `story-splitting` | ...` 那列之後插入一列：

```
| | `example-mapping` | 四色卡、四條診斷、就緒判定 | **已實作**。Pass 2 的手法與就緒判定從 `bdd-clarify` 搬進來 |
```

並把 `bdd-clarify` 那列的「三個 pass 的總入口 ＋ 就緒判定」改成「三個 pass 的總入口」。

- [ ] **Step 3: `benchmark/skeleton/go/README.md` 的兩張表**

第 233 行 `| `EX-1.1` | `@example-1.1` on a `Scenario` | ...` 的 `EX-1.1` 改成 `AC-1.1`。
第 256 行 `| `@example-1.1` | traces to an `EX-1.1` in `prd.md` |` 的 `EX-1.1` 改成 `AC-1.1`。

**tag 欄位不動**——`@example-1.1` 是對的。

- [ ] **Step 4: 全 repo 收尾驗證**

```bash
find . -name __pycache__ -type d -not -path './.git/*' -exec rm -rf {} + 2>/dev/null
for s in bdd-clarify bdd-spec bdd-plan story-splitting clarify-loop example-mapping skill-rules; do
  python3 skills/skill-rules/scripts/audit_skill.py skills/$s >/dev/null 2>&1; echo "  $s exit=$?"
done
claude plugin validate .; echo "validate exit=$?"
grep -rn "EX-" skills/ docs/ai-sdlc.md PLAN.md benchmark/ || echo "EX- 零筆"
grep -rn "##### AC\|##### Examples\|^#### FR-" skills/ || echo "v2 的標題寫法零筆"
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
git status --porcelain
```

Expected：七個 audit exit 0；`validate` exit 0（只有 root `CLAUDE.md` 警告）；`EX-` 與 v2 的標題寫法都零筆（`docs/superpowers/` 底下的歷史文件不在掃描範圍，那裡保留是對的）；連結全解得到；工作區乾淨。

- [ ] **Step 5: Commit**

```bash
git add docs/ai-sdlc.md PLAN.md benchmark/skeleton/go/README.md
git commit -F - <<'EOF'
docs: name the technique in the pipeline doc and the roadmap

WHAT: Say in the SDLC doc that Discovery's technique is Example Mapping
      and where its operational definition lives, add the new skill to
      the roadmap table, and repoint the testbed's two id references
WHY: The pipeline doc described what an agent can and cannot do during
     Discovery without ever naming the technique, so a reader had no way
     to get from the practice to the instructions. What it already said
     about the agent lacking final authority stays exactly as it is —
     that sentence is the reason the readiness vote belongs to the user,
     and the new skill cites it.
HOW: The testbed's tag column is deliberately unchanged; only the
     document-side id it traces to moves, which is the whole point of
     keeping each artifact's own vocabulary.

Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_014CNoiKUBCWt8jxvA9cyCSN
EOF
```
