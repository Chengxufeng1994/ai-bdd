# `prd.md` 格式 v2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 `prd.md` 從七節擴成十節，加入 Document Overview／Background／Goal／Success Metrics／Personas／User Stories，並把 FR 從平鋪的 `## Functional Requirements` 改成掛在 `### US-<n>` 底下的 `#### FR-<n>`。

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
| `skills/bdd-clarify/references/prd-format.md` | 格式契約：十節的順序與各節內容、三條參照完整性規則、已知弱點 | 1 |
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

Expected：`status.py` 印出 **已答 2／n/a 1／待答 1**，覆蓋為 `邊界` ＋ `降級`；`frs_in_prd` 印出 `{'1': {'1.1','1.2','1.3'}, '2': {'2.1'}}`。**把這兩個結果抄下來——Step 8 要逐字比對。**

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

- [ ] **Step 4: 遷移 `prd-format.md` 其餘各處的舊章節名**

舊名稱織在全文，不只在那張表裡。逐處改：

**開頭「三個 Pass」那段**（`Pass 1 開骨架（Problem／Actors／Scope）` 起到
`## Assumptions / Constraints`。止）整段換成：

```
三個 Pass 都對同一份文件寫：Pass 1 開骨架（Document Overview／Background／
Goal／Success Metrics／Scope／Personas），Pass 2 把 `## User Stories` 填滿
——先寫 `### US-<n>`，再把 `#### FR-<n>` 掛進去，Pass 3 把量得出來的品質屬性
填成 NFR（只服務一則 story 的掛在那則底下，跨 story 的進
`## Non-Functional Requirements`）、把不可協商的外部限制填進
`## Assumptions / Constraints`。
```

**「來源標記：三選一」的適用範圍**（`## Functional Requirements` 的每條 FR／AC／EX
那句起，到 `行尾都要標來源，三選一：` 止）換成：

```
`## User Stories` 底下的每條 FR／AC／EX 與 story 專屬的 NFR、
`## Non-Functional Requirements` 的每條 NFR、`## Scope — In / Out` 的每條邊界、
`## Assumptions / Constraints` 的每條，行尾都要標來源，三選一：
```

同一節的程式碼範例 `### FR-2  When 月租戶刷卡出場, …` 井號改成四個：
`#### FR-2  When 月租戶刷卡出場, the system shall 直接開啟柵欄，不計費。 ← PRD §4`

**該節末尾的豁免段**（`## Actors` 與 `## Open Questions` 不在這條規則的範圍內
那段）整段換成：

```
`## Open Questions` 不在這條規則的範圍內——問題本身就是在記錄還沒有答案的
東西，不需要再疊一層來源標記。

`## Personas` **在**範圍內。v1 的角色只有「是誰／怎麼取得／跟誰容易混」，
來源就是表格自己的欄位；v2 多了「他要什麼」與「現在什麼讓他痛」，而那兩件事
可以是問出來的，也可以是編出來的——正是這個標記存在的理由。
```

**「Non-Functional Requirements 與 Constraints：怎麼分邊」的第一句**：
`一件事該進 `## Functional Requirements` 底下的某條 FR` 改成
`一件事該進 `## User Stories` 底下的某條 FR`。

**同節的 NFR 程式碼範例**改成跨 story 的那個家，並在範例後補一段：

````markdown
## Non-Functional Requirements

- NFR-2  車輛進出紀錄保存至少一年
````

**只服務一則 story 的 NFR 不寫在這裡**，寫成那則 story 底下的 `#### NFR-<n>`。
判準是跨不跨 story，不是重不重要——「柵欄兩秒內開啟」只跟出場那則有關，
「紀錄保存一年」對每一則都成立。

**「角色與詞義：兩條界線」的第一句**：`角色留在 `prd.md` 的 `## Actors`` 改成
`角色留在 `prd.md` 的 `## Personas``。

改完自查：

```bash
grep -n "Functional Requirements\|## Actors\|Problem／Actors" skills/bdd-clarify/references/prd-format.md
```

Expected：只剩 `## Non-Functional Requirements` 的各處（含「與 Constraints：
怎麼分邊」那個標題），`## Functional Requirements` 與 `## Actors` 一個都不剩。

- [ ] **Step 5: 改寫「編號規則」一節——標題與內文都還在講 v1**

`## 編號規則：FR 平鋪，AC／EX 掛在底下` 從標題到內文都跟 v2 相反，內文還寫著
「這個 feature 不再切 story」。整節（標題到 `沒有從具體情境長出來。` 止）換成：

`````markdown
## 編號規則：FR 掛在 story 底下，AC／EX 掛在 FR 底下

`FR-<n>` 住在它服務的 `### US-<n>` 底下，但**編號全域唯一**——不是每則 story
各自從 1 起算。`FR-1` 在 US-1 底下、`FR-4` 在 US-2 底下都可以，但不得有兩個
`FR-1`，否則 `EX-1.1` 指向兩個地方。

NFR 的編號同樣全域唯一，而且跨越它的兩個家：story 專屬的 `#### NFR-<n>` 與
跨 story 的 `## Non-Functional Requirements`。

`AC-<n>.<m>` 與 `EX-<n>.<m>` 掛在各自的 `FR-<n>` 底下，`<n>` 跟著父層的 FR：

````markdown
### US-1  身為 P-1（訪客），我想在出場前知道要付多少，以便付完就能開走

#### FR-1  When 訪客車出場, the system shall 依停留時長計費，前 30 分鐘免費。

##### AC
- AC-1.1  訪客在繳費機看得到金額與停留時長

##### Examples
- EX-1.1  停 29 分 → 收 0 元
- EX-1.2  停剛好 30 分 → 收 0 元
- EX-1.3  停 31 分 → 收 30 元
````

**編號在這裡定版**——不是在 SPEC。`.feature` 的 `@example-<n>.<m>` tag 回指
`EX-<n>.<m>`，這是這條鏈唯一的接縫。

MUST NOT: 重排既有的 FR／AC／EX 編號。重排會讓那些引用**靜默**指向別的東西
——不會報錯，只會對錯。刪掉一條規則就留下空號（例如只剩 `FR-1`、`FR-3`，
沒有 `FR-2`），不要把後面的號碼往前遞補：**空號看得出來，重排看不出來。**

**把一條 FR 從一則 story 搬到另一則，編號不變。** 巢狀是給人讀的分組，編號是
給機器追的身分——搬家換分組，不換身分。搬完在版本修訂歷史留一列，因為 SPEC
可能已經照著舊分組切過 story 了。

每條 FR 底下至少要有一個 EX——完全沒有例子的 FR，通常是規則還停在想像階段，
沒有從具體情境長出來。
`````

- [ ] **Step 6: 遷移 `examples/minimal-prd.md` 到 v2**

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

- [ ] **Step 7: 改 `check_spec.py` 的 `frs_in_prd()`**

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

- [ ] **Step 8: 回歸——兩個結果必須跟 Step 1 逐字相同**

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

- [ ] **Step 9: 端到端——`check_spec.py` 對完整 fixture 跑一次**

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

- [ ] **Step 10: 缺 `## User Stories` 要大聲失敗**

```bash
rm -rf /tmp/v2bad && mkdir -p /tmp/v2bad/specs/2026-09-09-y /tmp/v2bad/features
grep -v "^## User Stories" skills/bdd-clarify/examples/minimal-prd.md > /tmp/v2bad/specs/2026-09-09-y/prd.md
cp /tmp/v2/specs/2026-09-09-x/spec.md /tmp/v2bad/specs/2026-09-09-y/
python3 skills/bdd-spec/scripts/check_spec.py /tmp/v2bad; echo "exit=$?"
```

Expected：`exit=1`，訊息指名該 feature 目錄與「沒有 `## User Stories`」。

- [ ] **Step 11: fail loud 沒退化**

```bash
python3 skills/bdd-spec/scripts/check_spec.py /tmp/nonexistent-v2; echo "exit=$?"
python3 skills/bdd-clarify/scripts/status.py /tmp/nonexistent-v2; echo "exit=$?"
```

Expected：兩次都 `exit=1` 並指名缺什麼。

- [ ] **Step 12: 稽核與 diff 範圍確認**

```bash
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify; echo "exit=$?"
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec; echo "exit=$?"
git diff --stat -- skills/bdd-clarify/scripts/status.py
```

Expected：兩個 audit exit 0；`status.py` 的 diff **為空**（Global Constraints 明訂不得改動）。

- [ ] **Step 13: Commit**

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

該節若列出 `prd.md` 的章節，換成 Task 1 的十節；若只寫「格式 → `references/prd-format.md`」則不動。**先讀再改，不要憑印象。**

- [ ] **Step 3: 全檔掃一次章節名引用**

```bash
grep -n "## Functional Requirements\|## Actors\|Problem / Goal" skills/bdd-clarify/SKILL.md
```

每一處逐一判斷：`## Actors` → `## Personas`；`## Functional Requirements` → `## User Stories` 底下的 `#### FR-<n>`；`Problem / Goal / Success Metrics` → 拆成的三節。**逐處判讀語意再改，不要無差別取代**——有些句子講的是「功能需求」這個概念而不是那個章節名。

- [ ] **Step 4: 把 `actor-definition.md` 改名**

```bash
git mv skills/bdd-clarify/references/actor-definition.md \
       skills/bdd-clarify/references/persona-definition.md
```

用 `git mv` 而不是刪了重建，`git log --follow` 才追得到——這是 repo 慣例。

- [ ] **Step 5: 改寫 `persona-definition.md`**

這份是**被改名那一節的專屬 reference**，從標題到範例整份都在講 v1。逐處改：

**標題**：`# 怎麼定義 Actor` → `# 怎麼定義 Persona`

**第一段**：`` `prd.md` 的 `## Actors` 節的寫法 `` → `` `prd.md` 的 `## Personas` 節的寫法 ``

**第三段開頭**：`Actor 是**角色**，不是帳號。` → `Persona 是**角色**，不是帳號。`
（該段其餘不動）

**`## 每個角色三件事` 整節**（標題到 `那是正常的。` 止）換成：

````markdown
## 每個角色五件事

| 欄位 | 內容 | 為什麼必填 |
| --- | --- | --- |
| **是誰** | 一句話，用業務語言 | 沒有這句，角色名稱會被各自解讀 |
| **怎麼取得這個身分** | 具體條件 ＋ 來源 | **最常漏的一格**，見下 |
| **跟誰容易混** | 區分它的那條規則 ＋ 來源 | 沒有區分的兩個角色是同一個角色 |
| **他要什麼** | 這個角色想達成什麼 | story 掛在 persona 底下，**沒有目標的 persona 會長出沒有理由的 story** |
| **現在什麼讓他痛** | 現況哪裡不好 | 沒有痛點的角色，通常是從組織圖抄來的，不是從需求長出來的 |

v1 問的是「跟誰不同」，v2 問「跟誰容易混」——**答案的內容契約沒變**，仍然要
講出區分它們的那條規則與來源。換問法是因為「跟誰不同」會得到「他們不一樣」
這種同義反覆，「跟誰容易混」會逼出實際會混淆的那一組。

**最後兩欄允許寫「（無資料）」，而且比編一個好。** 需求方常常在轉述別人的痛
（主委轉述住戶、住戶轉述他的客人）。標「無資料」讓那件事在文件上看得見；
填「訪客希望流程順暢」則讓它消失。

v1 有一個選填的 **能做什麼** 欄（每條回指 `FR-<n>`），v2 拿掉了：FR 現在住在
`### US-<n>` 底下，而每則 story 都寫明服務哪個 `P-<n>`——**這個角色能做什麼，
讀 story 就有了**。手抄一份等於兩個真相來源，遲早不一樣。

**「沒問過」與「（無資料）」本身就說明了來源狀態**，不再疊 `←` 標記。其餘每
一欄都要標，規則見 `prd-format.md` 的「來源標記：三選一」。
````

**中段的程式碼範例**（```markdown 圍籬包住、`## 旅程購買` 到 `## 旅程老師` 那塊）
換成 v2 的形狀：

````markdown
### P-1  旅程購買者

**是誰**：已買下某趟旅程、可永久存取其內容的學員 ← `pay-order` Rule 5
**怎麼取得這個身分**：付款成功後依商品的方案項目授予 ← `pay-order` Rule 5
**跟誰容易混**：跟「旅程訂閱狀態」容易混——訂閱會過期，這個不會 ← `query-user-roles` Rule 2
**他要什麼**：買過的內容隨時看得到，不必擔心過期 ← 推論
**現在什麼讓他痛**：（無資料）

### P-2  旅程老師

**是誰**：能維護某趟旅程內容的人 ← 推論
**怎麼取得這個身分**：**沒問過** → 開一題進 `## Open Questions`，狀態待答
**跟誰容易混**：目前沒有任何一條規則區分它與學員——見下方「兩個角色沒有規則區分」
**他要什麼**：（無資料）
**現在什麼讓他痛**：（無資料）
````

**`## 寫進 `prd.md` 的 `## Actors`` 一節**（標題到 `有需要時再開一題」。` 止）換成：

````markdown
## 寫進 `prd.md` 的 `## Personas`

一個 persona 一個 `### P-<n>` 小節，不是表格的一列。**因為 `## User Stories`
的每一則都要寫「身為 P-<n>」，而表格的一列給不出可引用的識別碼。** 欄位與完整
範例見 [`../examples/minimal-prd.md`](../examples/minimal-prd.md)，本文不重複
格式，只講前面幾節那些判準怎麼落進欄位：**是誰**與**怎麼取得這個身分**照抄；
**跟誰容易混**同時是欄位也是收尾時的核對項（見「兩個角色沒有規則區分，就是
同一個角色」）；**他要什麼**是 story 的來源，寫不出來就代表這個角色還沒有
存在的理由。

MUST: `P-<n>` 的編號全域唯一，而且**不重排**——`## User Stories` 用 `P-<n>`
回指，重排會讓那些引用靜默指向別的角色。刪掉一個角色就留空號。

考慮過但判定不是角色的，附一句理由寫在同一節——否則下一輪會有人再提一次，
例如「管理員——需求裡沒有任何一條規則提到它，有需要時再開一題」。
````

**`## 下游怎麼用它` 一節**：把三處 `` `## Actors` `` 全改成 `` `## Personas` ``，
其餘文字不動。

改完自查：

```bash
grep -n "Actor\|## Actors" skills/bdd-clarify/references/persona-definition.md
```

Expected：零筆。標題、內文、範例都不該再出現 `Actor`。

- [ ] **Step 6: 更新唯一的入站連結**

`skills/bdd-clarify/SKILL.md:177` 現在寫 `` → `references/actor-definition.md`。``
改成 `` → `references/persona-definition.md`。``

```bash
grep -rn "actor-definition" skills/
```

Expected：零筆（`docs/superpowers/plans/2026-09-09-*.md` 的兩處是歷史文件，**不要動**）。

- [ ] **Step 7: 驗證**

```bash
python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-clarify; echo "exit=$?"
grep -n "## Functional Requirements\|## Actors" skills/bdd-clarify/
python3 skills/bdd-clarify/scripts/status.py /tmp/v2
```

Expected：audit exit 0；grep 只剩刻意保留的歷史敘述；`status.py` 仍是 **已答 2／n/a 1／待答 1**。

標題深度也要驗，不能只驗章節名——v2 同時改了名稱與深度，而**名稱式的 grep 對
深度變化結構性失明**（Task 1 就是這樣讓兩節維持 v1 深度通過全部驗證的）：

```bash
grep -rn "^#\+ \(FR-\|NFR-\|AC$\|Examples$\|US-\|P-\)" skills/ --include="*.md"
```

Expected：`US-` 與 `P-` 恰好三個井號，`FR-` 與 `NFR-` 四個，`AC` 與 `Examples`
五個，無例外。


- [ ] **Step 8: Commit**

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

- [ ] **Step 4: 改 `story-splitting/SKILL.md` 第 36 行**

```
- 1. **要切的 story**：一句話，或 `prd.md` 的 `## Functional Requirements`
+ 1. **要切的 story**：一句話，或 `prd.md` 的某幾條 `FR-<n>`
```

同一節第 37 行已經寫著「有 `prd.md` 就直接用它的 `FR-<n>`」，所以這行本來就在
講 FR，只是用了一個已經不存在的章節名指它。

- [ ] **Step 5: 驗證**

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

標題深度也要驗，不能只驗章節名——v2 同時改了名稱與深度，而**名稱式的 grep 對
深度變化結構性失明**（Task 1 就是這樣讓兩節維持 v1 深度通過全部驗證的）：

```bash
grep -rn "^#\+ \(FR-\|NFR-\|AC$\|Examples$\|US-\|P-\)" skills/ --include="*.md"
```

Expected：`US-` 與 `P-` 恰好三個井號，`FR-` 與 `NFR-` 四個，`AC` 與 `Examples`
五個，無例外。


- [ ] **Step 6: Commit**

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
