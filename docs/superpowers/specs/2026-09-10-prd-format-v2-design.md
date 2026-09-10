# `prd.md` 格式 v2：加上 PRD 該有的那幾節，FR 改掛在 User Story 底下

**日期**：2026-09-10
**狀態**：設計已確認，待實作
**動的東西**：`references/prd-format.md`、`examples/minimal-prd.md`、`bdd-clarify/SKILL.md`、`bdd-spec/SKILL.md`、`scripts/check_spec.py`

---

## 要解的問題

v1 的 `prd.md`（2026-09-09 上線）只有八節：Problem/Goal/Success Metrics、Actors、
Scope、Functional Requirements、NFR、Open Questions、Assumptions/Constraints。

它**能被解析、能驅動 SPEC**，但它不像一份可以拿去跟 PM 對焦的 PRD。2026-09-10 的
實測（`parking-billing`，19 題問了 9 題後暫停）暴露了三個缺口：

| 缺的 | 實測撞到的樣子 |
| --- | --- |
| **核心關係人** | 四個待答問題掛在四個不同的人身上（設備廠商／會計師／律師／區權會），但文件裡只有「該問誰」散在各題底下，**沒有一個地方看得到「這個案子卡在哪四個人身上」** |
| **文件狀態** | 跑到一半暫停，檔案本身看不出它是暫停中還是完成了 |
| **背景與痛點** | 成功指標是需求方講抱怨時被問出來的，但 `## Problem / Goal / Success Metrics` 三件事擠在一節，痛點沒有自己的位置 |

---

## 已確認的決策

### D1 — FR 掛在 User Story 底下，`## Functional Requirements` 這一節消失

使用者的決定。我提出過兩個反對意見，兩個都被看過之後仍然選這個：

1. `check_spec.py` 的 `^### FR-(\d+)` 抓不到 H4，會解析出零條、不報錯、覆蓋檢查
   靜默失效 → **本設計把它列為必改**
2. FR 按 story 分組，形式上等於 CLARIFY 又在分組 → 見 D2

### D2 — PRD 的 story 分組是敘事分組，不是交付切片

**這是 D1 能成立的前提，必須寫進兩個 skill。**

2026-09-09 的重構把切 story 從 CLARIFY 搬到 SPEC，理由是切分需要範圍先穩定——
而實測證實了那個理由：15 題範圍題裡需求方改口三次（入場擋不擋、月租錢進不進
系統、免費額度怎麼算），第一批之後畫的切線都會是錯的。

**但那次的根源是目錄，不是章節。** 舊結構的主文件住在 `specs/<slice-slug>/`——
要有 slug 才有目錄，所以切分被迫先發生，切錯要重建目錄樹。現在 story 是一份檔案
裡的章節，改一則 story 的邊界＝把一條 FR 從一個 `###` 搬到另一個，代價接近零。
**逼迫切分的那個力量已經不存在。**

剩下的風險只有一個：SPEC 被現成的分組綁住。用一條規則管：

> PRD 的 `## User Stories` 是需求方視角的敘事分組，不是交付切片。
> **SPEC 得重新分組，不必解釋為什麼跟 PRD 不一樣**；但如果它照抄 PRD 的分組，
> 要說得出為什麼那個分組剛好也是好的交付邊界。

### D3 — FR 編號全域唯一

`FR-1` 在 US-1 底下、`FR-4` 在 US-2 底下都可以，但**不得有兩個 `FR-1`**——否則
`EX-1.1` 指向兩個地方，`.feature` 的 `@example-1.1` 失去意義。

### D4 — NFR 有兩個家，判準是「跨不跨 story」

- 某則 story 專屬的（「出場柵欄自刷卡到開啟不超過 2 秒」）→ 掛在那則 story 底下
- 跨全部的（個資保存期限、可用性）→ 留在 `## Non-Functional Requirements`

放錯的成本很低（搬一段），但**兩個家一定要有判準**，否則會漂移。

**NFR 編號跟 FR 一樣全域唯一。** 兩個家不代表兩套編號——`NFR-1` 掛在某則 story
底下、`NFR-2` 在跨 story 那一節，是合法的；兩個 `NFR-1` 不是。搬家時編號不變，
這樣「NFR-3 從 story 專屬改成跨 story」在版本修訂歷史裡才追得到。

### D5 — `## Problem / Goal / Success Metrics` 拆成三個 H2

`## Background`（背景與痛點）、`## Goal`（產品願景與目標）、`## Success Metrics`。

### D6 — `## Actors` 升級為 `## Personas`，加兩欄

保留現有三欄（是誰／怎麼取得／跟誰容易混），新增**「他要什麼」**與
**「現在什麼讓他痛」**。理由：FR 掛在 story、story 掛在 persona，**沒有目標的
persona 會產出沒有理由的 story**。

**「無資料」是合法的值，而且比編一個好。** 實測全程沒有任何訪客的聲音——需求方
是主委，他轉述住戶的抱怨，住戶又轉述客人的抱怨，三手資訊。標「無資料」讓這件事
在文件上看得見；填「訪客希望流程順暢」則讓它消失。

### D7 — 文件狀態只寫質性，不寫數字

四個狀態，判準是「下游能不能開工」：`澄清中`／`待對焦`／`已對焦`／`已凍結`。

**不得寫「9 題已答／10 題待答」。** `bdd-clarify/SKILL.md` 的產物一節已有硬規定：
「MUST NOT: 另存一份進度儀表板。進度是**算出來的**」。`status.py` 的 docstring
記著這個 repo 已經被咬過一次：map 檔頭寫 21 個例子、儀表板寫 23、實際 23。

### D8 — 版本修訂歷史的觸發條件

欄位：`版本｜日期｜改了什麼｜為什麼｜誰`。

**觸發**：文件到 `已對焦` 之後，任何動到 FR／AC／EX 的改動都要留一列——因為 SPEC
可能已經照著寫了 `.feature`。`已對焦` 之前的來回不記，那是澄清過程本身，決策史在
`## Open Questions` 的小節裡。

**「為什麼」不可省。** 實測有現成例子：FR-1 的免費額度從「每趟 30 分」改成
「一天總共 30 分」，因為算出來舊規則對想保護的人反而更兇。只記改了什麼，半年後
有人會把它改回去。

### D9 — 核心關係人列權威，不列阻塞

欄位：`關係人｜負責什麼｜決策權`。

**「決策權」不是裝飾。** 實測撞到過：主委可以自己決定公務車免費，但「住戶幫訪客
折抵」得開區權會。**同一個人給的答案，一個是定案、一個是提案**，下游要分得出來，
否則會拿沒表決過的東西去寫規格。

**這一節不列「誰卡著哪幾題」**，改成一條可機械檢查的規則：

> **MUST**：`## Open Questions` 每一題的「該問誰」必須指向這張表裡的某個關係人。

阻塞的全貌由腳本算，永遠不會過期。指到表外的人＝那個人沒被登記為關係人，是文件
的洞不是筆誤。

### D10 — User Story 必須指向 Persona

```markdown
### US-1  身為 P-1（訪客），我想在出場前知道要付多少，以便準備好再開到閘門 ← 推論
```

`身為` 後面必須是一個 `P-<n>`。跟 D9 同一個機制：可機械檢查的參照完整性。

### D11 — 一條 FR 服務多則 story：住主要的那則，其他宣告依賴

巢狀強迫一對多，共用的 FR 只能住一個地方。實測有現成例子：`FR-3`（訪客區滿了不開
入口閘）同時服務訪客（不白跑一趟）與主委（保護訪客格）。

```markdown
### US-2  身為 P-2（主委），我想知道訪客格有沒有被外人占走…

**另外依賴**：FR-3（定義在 US-1 底下）
```

**這是巢狀的已知弱點，見〈已知弱點〉一節。**

---

## 目標結構

```
## Document Overview            狀態、最後更新、版本修訂歷史、核心關係人
## Background                   背景與痛點
## Goal                         產品願景與目標
## Success Metrics
## Scope — In / Out
## Personas
    ### P-1  <角色名>
## User Stories
    ### US-1  身為 P-<n>（<角色名>），我想<做什麼>，以便<得到什麼>
        #### FR-1  When … the system shall …          ← EARS
            ##### AC        AC-1.1 …    給 PM
            ##### Examples  EX-1.1 …    給 SPEC
        #### NFR-1  …                                  ← 這則 story 專屬
## Non-Functional Requirements  跨 story 的
## Open Questions               表是索引，小節是決策史
## Assumptions / Constraints
```

**沿用 v1 不變的**：來源標記三選一（`PRD §x`／`Q<n>`／`推論`）適用每一條 FR／AC／
EX／NFR／Scope／Assumption；`## Open Questions` 的五欄表（`Q｜問題｜面向｜狀態｜
答案`）與 `面向` 的值域；AC 與 EX 的兩層讀者；EARS 句式引用 `docs/sdd.md`。

---

## 連帶要改的

| 檔案 | 改什麼 |
| --- | --- |
| `check_spec.py` 的 `frs_in_prd()` | 從「`## Functional Requirements` 段找 `^### FR-(\d+)`」改成「`## User Stories` 段找 `^#### FR-(\d+)`」。EX 的抓法不變 |
| `check_spec.py` 的錯誤訊息 | 「沒有 `## Functional Requirements`」改成「沒有 `## User Stories`」 |
| `examples/minimal-prd.md` | **必改，而且是最關鍵的一項**——它是兩支腳本唯一的 canonical 解析對象。不遷移它，`check_spec.py` 改完之後沒有東西可以驗 |
| `status.py` | **不動**。它只讀 `## Open Questions`，那一節的位置與格式都不變 |
| `bdd-clarify/SKILL.md` | 產物一節、流程圖的 Pass 3 路由（NFR 現在有兩個家） |
| `bdd-spec/SKILL.md` | 加入 D2 的對應規則：讀 PRD 的 story 分組當線索，但得自己重新分組 |

---

## 已知弱點

**D11 的依賴宣告沒有東西在檢查它。** 如果 US-2 忘了宣告依賴 FR-3，文件看起來完全
正常，只是 SPEC 為 US-2 切 story 時會漏掉一條規則。

這是巢狀換可讀性的代價，不是可以設計掉的東西——除非把 FR 平鋪回去，而那是 D1
已經否決的方向。**寫進格式文件，不藏起來。**

（可以做但本次不做的緩解：`check_spec.py` 檢查每條 FR 至少被一則 story 涵蓋。那
只抓得到「完全沒人用的 FR」，抓不到「該宣告依賴卻沒宣告」。）

---

## 驗證

**兩份要遷移的檔案，性質不同：**

- `skills/bdd-clarify/examples/minimal-prd.md` —— **必須遷移**。它是兩支腳本唯一的
  canonical 解析對象，不動它就沒有東西可以驗。
- `runs/go-parking-billing/specs/2026-09-10-parking-billing/prd.md` —— **選用**。它有
  6 條 FR、19 題、完整來源標記，是**現成的真實產物**；把它改寫成 v2 可以驗證新格式
  裝得下真實規模的東西，而不只是裝得下最小範例。它是 gitignored 的，改壞了不影響 repo。

| # | 檢查 | 通過條件 |
| --- | --- | --- |
| 1 | `audit_skill.py` 對 `bdd-clarify`、`bdd-spec` | exit 0 |
| 2 | `status.py` 對改寫後的 fixture | 數字與改寫前**相同**（`## Open Questions` 沒動，變了就是改壞了） |
| 3 | `check_spec.py` 對改寫後的 fixture | FR/EX 解析結果與 v1 相同 |
| 4 | 兩支腳本對**不存在的路徑** | exit 1 並指名缺什麼 |
| 5 | `check_spec.py` 對一份**缺 `## User Stories`** 的 prd.md | 指名該檔案、非零離開 |
| 6 | 全 repo 相對 markdown 連結 | 全解得到 |

第 2 項是回歸檢查的重點：這次改的是 FR 那一段，`## Open Questions` 不該受影響，
**但 fixture 是同一個檔案**，改 A 弄壞 B 是最容易發生的事。

---

## 明確不做的事

- **`spec.md` 的格式** —— SPEC 仍然自己切 story、寫 `## Stories`，格式不變
- **D11 依賴宣告的機械檢查** —— 見〈已知弱點〉
- **Persona 的更多欄位**（使用頻率、技術熟練度、引述）—— 兩欄夠了；等真的有 persona
  因為缺欄位而做錯決定，再加
- **重跑 benchmark** —— 改完之後 `parking-billing` 那次的產物是 v1 格式，遷移它是
  驗證手段，不是重新評分
