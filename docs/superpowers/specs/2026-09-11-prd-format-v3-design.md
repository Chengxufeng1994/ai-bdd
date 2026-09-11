# `prd.md` 格式 v3 ＋ `example-mapping` skill 設計

**日期**：2026-09-11
**狀態**：待實作
**前一版**：[`2026-09-10-prd-format-v2-design.md`](./2026-09-10-prd-format-v2-design.md)

v2 上線一天。它做對的事是把規則搬到它服務的 story 底下；沒做的事是**讓文件長得像
那張桌子**。這一版補上那件事，並且把 Example Mapping 從一句引用變成一個 skill。

不在本次範圍：`bdd-clarify` → `bdd-discovery`、`clarify-loop` → `clarify` 這兩個
改名。它們是獨立的一份工作（純 `git mv` ＋ 引用更新），與格式無關，**刻意排在
這一份之後**——改名自成一筆提交才不會讓 `git log --follow` 斷掉。

---

## 為什麼改

### 一、文件看不出它是一場 map 的產物

Example Mapping 的桌上有四種卡：🟨 story、🟦 rule、🟩 example、🟥 question。v2 的
`prd.md` 四種都有，但**擺法看不出它們是同一場對話的四種卡**：規則與例子在
`## User Stories` 裡，問題飛到文件最後的全域表。

一則 story 的紅卡跟它的藍卡隔著半份文件。要判斷「這則 story 就緒了嗎」，得在兩個
地方之間來回。

### 二、例子只有一行，把前置狀態藏起來了

v2 的例子是單行：`EX-1.1  停 29 分 → 收 0 元`。

「停 29 分」的前置狀態是什麼？車已經入場了嗎？計時從進柵欄算還是從抽票算？
**單行寫法讓這些不用回答，而 `bdd-spec` 寫 `.feature` 時必須回答**——於是它自己
編一個。

而「SPEC 憑空長出來的內容是缺陷」是這個 repo 寫在 `CLAUDE.md` 裡的鐵律。單行例子
是那條鐵律的一個結構性漏洞：不是 SPEC 不守規矩，是 CLARIFY 沒把話講完。

### 三、Example Mapping 散在三個 skill，沒有一個是權威

| 內容 | v2 住哪 |
| --- | --- |
| 四色卡、25 分鐘、大拇指投票、四條啟發式 | `docs/bdd.md`（來源轉述，**明訂不得帶專案意見**） |
| 「25 分鐘 map 不完 ＝ 範圍問題」 | `clarify-loop/SKILL.md` 拿它判斷零進展 |
| 「藍卡太多 ＝ story 太大」「沿規則切」 | `story-splitting` |
| 「紅卡太多 ＝ 不確定性高」 | `clarify-loop` 的紅卡機制 |
| 規則與例子本身 | `bdd-clarify` 的 Pass 2 |

`docs/bdd.md` 有定義但**不准講怎麼在這個專案做**（那是它的存在條件）。所以「怎麼
在這裡跑一場 map」現在沒有家。

---

## 決策

### D1 · story 底下三個分組，AC 巢狀在它的 FR 底下

```markdown
### US-1  身為 P-1（訪客），我想在出場前知道要付多少 ← 推論

#### FR
**FR-1**  When 訪客車出場, the system shall 依停留時長計費，前 30 分鐘免費。 ← PRD §3
- **AC-1.1**  停 29 分 → 收 0 元 ← PRD §3
  - Given 一台訪客車已入場並抽了紙票
  - When 它在入場後 29 分鐘於繳費機結帳
  - Then 應收金額是 0 元
- **AC-1.2**  停 31 分 → 收 30 元 ← PRD §3
  - Given 一台訪客車已入場並抽了紙票
  - When 它在入場後 31 分鐘於繳費機結帳
  - Then 應收金額是 30 元

**FR-2**  When 月租戶刷卡出場, the system shall 直接開啟柵欄，不計費。 ← PRD §4
- **AC-2.1**  名單內且未到期 → 開柵欄，金額 0 ← PRD §4
  - Given 一台月租車在名單內且未到期
  - When 它刷卡出場
  - Then 柵欄開啟且應收金額是 0 元

#### NFR
**NFR-1**  出場柵欄自刷卡到開啟不超過 2 秒 ← Q-12

#### Q
**Q-7**  跨午夜的停車怎麼算 ← 擋住 FR-1
```

**為什麼 FR／NFR／Q 平輩，AC 不平輩。** 那張桌子是一棵樹，不是四個欄位——綠卡擺在
**它所屬的藍卡正下方**：

```
        [🟨 story]
   [🟦 rule-1]   [🟦 rule-2]
   [🟩][🟩]        [🟩]
        [🟥][🟥]
```

Markdown 攤不平二維，只能二選一。FR、NFR、Q 都是 **story 層級**的東西，並排成三個
分組是對的；AC 屬於某一條 FR，平輩化會讓「`#### FR` 列五條規則、`#### AC` 列十二條
AC，讀的人靠編號在腦裡接回去」。巢狀讓 `AC-1.1 屬於 FR-1` 是**結構上不可能弄錯**的，
不是編號約定。

**FR／NFR／AC／Q 全部改用粗體項目而非標題**，總深度停在 H4。理由是 AC 的 G/W/T 還要
再往下一層，用標題會走到 H6。

### D2 · AC 吃掉 v2 的 Examples，並帶 Given/When/Then

`EX-<n>.<m>` 消失。`AC-<n>.<m>` 同時是驗收條件與具體例子，底下掛三行 G/W/T。

v2 的 AC 是「給 PM 讀、不含數字」的那層（`AC-1.1 訪客在繳費機看得到金額與停留時長`），
**這一層刪掉**，由 story 那一句（`身為 P-1（訪客），我想在出場前知道要付多少`）承擔。

**取捨要記下來，因為它有真實成本**：PM 對焦時不想讀三十條 Given/When/Then。v2 的
兩層設計就是為了這個。選擇合併的理由是**一層寫不清楚的東西，兩層只會讓它在兩個
地方各寫一半**——v2 的 AC 層在實測中沒有攔下任何東西，而單行例子讓 SPEC 編造前置
狀態的洞是真的。

MUST: 每條 AC 底下的 G/W/T 三行齊全。Given 寫不出來就是這個例子的前置狀態還沒被
問過——**那是一張紅卡，不是一句可以省略的話**。

### D3 · 用 PRD 的詞彙，不用 Example Mapping 的詞彙

`prd.md` 的讀者帶著 PRD 的詞彙來：User Story、Functional Requirement、Acceptance
Criteria。**Example Mapping 是手法，PRD 是產物，手法的詞彙不該覆蓋產物的詞彙。**

對照表（`example-mapping` skill 的第一張表就是它）：

| 卡 | Example Mapping | `prd.md` |
| --- | --- | --- |
| 🟨 | story | `US-<n>` |
| 🟦 | rule | `FR-<n>` |
| 🟩 | example | `AC-<n>.<m>` |
| 🟥 | question | `Q-<n>` |

**一個已知的詞義衝突，必須寫明。** Matt Wynne 把藍卡描述成 "a business rule or
acceptance criterion"——在 Example Mapping 的血統裡，**AC 是藍卡**。在主流 PRD／Jira
用法裡，AC 是那串 Given/When/Then。兩個都對，但同一份文件裡只能有一個，這裡選 PRD 的。

`example-mapping` skill 要把這件事講出來，否則讀過 Cucumber 文章的人會以為對照表寫錯了。

### D4 · tag 不改名：兩份產物各用母語，接縫寫明對應

```
prd.md（PRD 詞彙）        .feature（Gherkin 詞彙）
FR-<n>              ↔     @rule-<n>
AC-<n>.<m>          ↔     @example-<n>.<m>
```

`@rule-<n> ↔ FR-<n>` **今天就是這樣運作的**，只是從來沒寫明。`@example-<n>.<m>` 照
同一條規則走：`.feature` 是 Gherkin 的產物，`@example` 是 Gherkin 的字。

**「兩個名字」在這裡不是缺陷。** 同一份產物裡有兩個名字才是缺陷；兩份產物各自用
母語、接縫上寫明對應，那是設計。

實際效果：`skills/bdd-spec/examples/` 十二個教學檔約 90 處 `@example-N.M` **一個
都不用動**，`check_spec.py` 第 209 行的 tag 正則不動。

MUST: 這條對應要寫進 `prd-format.md` 與 `bdd-spec/SKILL.md` 兩邊。只寫一邊，另一邊
的讀者就得自己猜。

### D5 · 紅卡：story 底下只放指標，`## Open Questions` 不動

story 底下的 `#### Q` 放 **`Q-<n>` ＋ 問題一句話 ＋ 擋住誰**。狀態、面向、答案、
決策史全部留在 `## Open Questions` 的五欄表。

**為什麼不全搬進 story。** 兩個理由：

1. **不是每張紅卡都屬於某則 story。** 範圍題（「月租費帳務在不在範圍」）與技術題
   （「消防法規對閘門的要求」）在 Pass 1 與 Pass 3 產生，那時候還沒有任何 story。
   全搬會逼出一個「未掛 story 的問題」小節——兩個地方裝問題。
2. **`status.py` 算的十六個面向覆蓋率會整支重寫。** 那一列是這個 repo 唯一看得出
   「Pass 2 有沒有真的問過邊界」的東西，也是少數真的有機械檢查的規則之一。

代價是**問題的文字重複一次**（story 底下一句、表裡一句）。這是明知的重複，換到的是
文件像地圖 ＋ `status.py` 完全不動。

**編號改成 `Q-<n>`**（v2 是 `Q<n>`），與 `US-<n>`／`FR-<n>`／`NFR-<n>`／`AC-<n>.<m>`
一致。`status.py:83` 用 `cells[0] in ("Q", "---")` 擋表頭，**表頭那一格仍然是 `Q`**，
資料列的編號格式它不解讀——已確認安全。

### D6 · 新 skill：`example-mapping`

**它擁有**：四色卡在這個 repo 的操作定義（D3 的對照表）、攤卡片的問題順序、四條
診斷啟發式、就緒判定。

| 模式 | 輸入 | 輸出 |
| --- | --- | --- |
| 流程內 | `prd.md` | 攤開的地圖 ＋ 四條診斷 ＋ 就緒問句；**使用者投票後**才寫 `狀態` |
| 單獨 | 一則 story 的描述 | 同樣的地圖與診斷，**不寫檔** |

雙模式與 `clarify-loop` 同型。單獨模式不產出檔案——2026-09-09 拿掉
`example-mapping.md` 這個產物的理由仍然成立：**第二個裝規則與例子的地方會跟
`prd.md` 漂移**。

**攤開的地圖長這樣**（終端輸出，不是檔案）：

```
US-1  身為 P-1（訪客），我想在出場前知道要付多少
│
├─🟦 FR-1  前 30 分鐘免費
│   ├─🟩 AC-1.1  停 29 分 → 收 0 元
│   ├─🟩 AC-1.2  停剛好 30 分 → 收 0 元
│   └─🟩 AC-1.3  停 31 分 → 收 30 元
├─🟦 FR-2  月租戶不計費
│   └─🟩 AC-2.1  名單內且未到期 → 開柵欄，金額 0
└─🟥 Q-7  跨午夜的停車怎麼算              ← 擋住 FR-1

診斷：🟦 2  🟩 4  🟥 1
  ⚠ 還有 1 張紅卡 —— 未就緒
  ✓ 每條 FR 都有例子
  ✓ 沒有規則堆了過多例子
```

**攤卡片的問題順序**：從 story 開始，之後**規則與例子哪個先都可以**——這是
Example Mapping 的原始建議，而且它給了一條判準：**談不攏規則時就寫例子，寫到規則
自己浮出來**。反過來也成立：規則講得很順但舉不出例子，那條規則多半還停在想像階段。

MUST NOT: 先把規則寫滿再回頭補例子。那個順序會得到一批「聽起來對、但沒有人驗證過
它在具體情境下成立」的規則——`prd-format.md` 的「每條 FR 底下至少要有一個 EX」
本來就是在擋這件事，但它只擋得住完全沒有例子的，擋不住事後補的。

**其他三個 skill 要刪掉重複的段落改成指向它**：

| skill | 刪什麼 |
| --- | --- |
| `bdd-clarify` | Pass 2 的手法部分 → 呼叫 `example-mapping`；Pass 1、Pass 3 不動 |
| `clarify-loop` | 「25 分鐘 map 不完」那一段（現在在 SKILL.md 第 187-190 行） |
| `story-splitting` | 「藍卡太多」的判準 |

### D7 · 兩條啟發式在這裡改了意義，照抄會給出錯的建議

| 卡 | 原始意思 | 這裡 | 傳給 |
| --- | --- | --- | --- |
| 🟥 還有紅卡 | 不確定性高，未就緒 | 同左 | `clarify-loop` |
| 🟦 一則 story 掛太多 FR | **story 太大，該切** | **不阻塞** ← 見下 | 只回報 |
| 🟩 AC 堆在同一條 FR | 這條規則該拆 | 同左 | `bdd-clarify` |
| ⏱ 25 分鐘 map 不完 | 太大或太不確定 | **沒有時鐘** ← 見下 | `clarify-loop` |

**🟦 為什麼不再是阻塞。** 原始啟發式假設黃卡上的 story **就是交付切片**。v2 的 D2
已經否定了這件事並寫進 `bdd-spec/SKILL.md`：

> `prd.md` 已經把 FR 分在 `### US-<n>` 底下了，**但那是需求方視角的敘事分組，不是
> 交付切片**。你得自己重新分組。

所以「US-2 掛了 9 條 FR」在這裡不代表未就緒——SPEC 本來就要重新分組。它代表的是
這則敘事 story 可能包了不只一個使用者目標。**回報，不阻塞，也不呼叫
`story-splitting`**——那會讓 `example-mapping` 去動 PRD 的分組，正好違反 v2 剛立的
規則。

**⏱ 沒有時鐘，那個框換成什麼。** 25 分鐘的框存在是因為人會累，而且 map 不完代表
太大。Agent 不會累，所以**可轉移的是尺寸訊號，不是那個鐘**。操作等價物：**攤地圖
的過程中紅卡出現得比綠卡快**——你一邊放例子、一邊冒出新問題，而且冒的比放的多。
`clarify-loop` 已經有同型的判準（「問了兩三輪、紅卡總數卻不下降」），這一條指向
它，不重寫。

### D8 · 就緒判定歸使用者，不歸 agent

`docs/ai-sdlc.md:98` 已經定調：agent 取代三方對話當提問者，但「**沒有業務脈絡的最終
裁量權——它問得出問題，答不了**」。

所以 `example-mapping` 攤開地圖、跑完四條啟發式、給出它的判讀，然後**問使用者**。

它是唯一**寫** `## Document Overview` 那個 `狀態` 欄位的 skill，但**決定**是使用者的。

這順便補上 v2 最終審查指出的一個洞：D7 的 `狀態` 是一個新的機器可讀閘門
（`澄清中` → SPEC 不得開始），而當時**沒有任何東西在強制它**，只加了一句「已知弱點」。
現在有了。

MUST NOT: agent 自己把狀態推到 `已對焦`。

---

## 連帶要改的

| 檔案 | 改什麼 |
| --- | --- |
| `skills/bdd-clarify/references/prd-format.md` | 章節順序表、`## User Stories` 一節整段重寫、編號規則、AC 與 Examples 一節（改成只講 AC）、來源標記範圍、D4 的對應表 |
| `skills/bdd-clarify/examples/minimal-prd.md` | **canonical 解析對象**，遷移到 v3：`EX-` 全改 `AC-` 並補 G/W/T；`## Open Questions` 的四列與兩個小節標題從 `Q1` 改成 `Q-1`（**表頭那一格仍是 `Q`**）；US-2 底下補 `#### Q` 指標；修掉下面那個既有缺陷 |

### fixture 現在有一個缺陷，遷移時要一併修掉

`minimal-prd.md` 的 Q2 小節寫著 `**擋住**：FR-11（緊急車輛進場）`，而**這份 fixture
裡沒有 FR-11**——只有 FR-1 與 FR-2。v2 沒有任何東西檢查這件事，所以它一直在。

這很要緊，因為 Q-2 是 fixture 裡**唯一一題待答的**。不修的話，沒有任何一則 story
會有 `#### Q` 指標可放，**v3 新增的紅卡結構在 canonical fixture 裡一次都不會被示範**
——正是 v2 最終審查抓到的那類問題（契約宣稱的東西，唯一的驗證對象沒有展示）。

修法：Q-2（消防法規對閘門的要求）的 `擋住` 從 `FR-11` 改成 **`FR-2`**（月租戶刷卡
出場開柵欄——同樣是閘門行為），並在 US-2 底下加：

```markdown
#### Q
**Q-2**  消防法規對閘門的要求 ← 擋住 FR-2
```

`## Open Questions` 的表**一列都不動**，所以驗證項 2 的數字不受影響。

MUST: `#### Q` 只列**狀態為「待答」**的問題。已答的不是紅卡了——它的答案已經長成
某條 FR 或 AC，留在那裡會讓地圖上的紅卡數量永遠不歸零。
| `skills/bdd-spec/scripts/check_spec.py` | `frs_in_prd()` 改抓 `**FR-<n>**` 與 `AC-<n>.<m>`；barren 檢查訊息；兩段 docstring。**第 209 行 tag 正則不動** |
| `skills/bdd-clarify/SKILL.md` | Pass 2 改成呼叫 `example-mapping`；產物一節的編號說明 |
| `skills/bdd-spec/SKILL.md` | `EX-` → `AC-`；寫明 D4 的兩側對應 |
| `skills/clarify-loop/SKILL.md` | 刪 25 分鐘那段，指向 `example-mapping` |
| `skills/story-splitting/SKILL.md` | 刪「藍卡太多」判準，指向 `example-mapping` |
| `skills/example-mapping/SKILL.md` | **新增** |
| `docs/ai-sdlc.md` | Discovery 那一列補一句：手法是 Example Mapping，操作定義在該 skill |
| `PLAN.md` | skill 表加一列 |
| `benchmark/skeleton/go/README.md` | `EX-1.1` 的兩張對應表 |
| `skills/bdd-clarify/scripts/status.py` | **不動**（D5 已確認 Q 欄格式它不解讀） |

---

## 這一版新增的機械檢查

v2 的最終審查指出「三條參照完整性 MUST 有兩條沒人檢查」。v3 的結構讓其中一條變成
可檢查的，應該一併做：

| 規則 | 怎麼檢查 |
| --- | --- |
| story 底下 `#### Q` 引用的每個 `Q-<n>` 都要存在於 `## Open Questions` 表 | `check_spec.py` 新增，**這是 v3 才有辦法做的** |
| 每條 FR 至少一條 AC | v2 的 barren 檢查改個名字就是 |
| `AC-<n>.<m>` 的 `<n>` 對得上同一則 story 的 FR | D1 的巢狀讓它結構性成立，仍加檢查擋打字錯 |

---

## 驗證

| # | 驗什麼 | 怎麼驗 |
| --- | --- | --- |
| 1 | `audit_skill.py` 對每個動過的 skill exit 0 | 含新的 `example-mapping` |
| 2 | `status.py` 對遷移後的 fixture **數字不變** | 已答 2／n/a 1／待答 1、覆蓋 `邊界`＋`降級`。`## Open Questions` 沒動，變了就是改壞了 |
| 3 | `check_spec.py` 解析得到 FR→AC 對應 | **`{'1': {'1.1','1.2','1.3'}, '2': {'2.1'}}`** —— 與 v2 的 EX 對應**逐字相同**，只是鍵名從 EX 換成 AC。不同就是遷移漏了一條 |
| 4 | **新檢查確實會失敗** | `spec.md` 點名 prd.md 沒有的 FR → 非零；`#### Q` 引用表裡沒有的 `Q-<n>` → 非零；FR 沒有 AC → 非零 |
| 5 | 深度／粗體格式掃描全 `skills/` 零例外 | v2 學到的教訓：**名稱式 grep 對結構變化失明**，這次連粗體項目格式一起掃 |
| 6 | fail loud 沒退化 | 兩支腳本對不存在路徑 exit 1 |
| 7 | `claude plugin validate .` exit 0 | 只有 root `CLAUDE.md` 那則已知警告 |
| 8 | 連結全解得到 | 用逐行圍籬追蹤，不用正則（v2 學到的） |

---

## 已知弱點

- **問題的文字在兩個地方各一份**（D5）。這是明知的重複，換到文件像地圖 ＋
  `status.py` 不動。沒有東西檢查兩份是否一致——只檢查 `Q-<n>` 存不存在。
- **`身為 P-<n>` 必須指向存在的 persona** 仍然沒有檢查。v2 已揭露，v3 不處理。
- **`**另外依賴**：FR-<n>` 的宣告仍然沒有檢查**。v2 已揭露，v3 不處理。
- **G/W/T 與 `.feature` 的 Scenario 有字面重複**。這是 D2 刻意付的代價：prd.md 的
  是**談定的**，`.feature` 的是**可執行的**，後者還要套封閉步驟文法與 `Rule:` 結構。
  沒有東西檢查兩者是否一致——`check_spec.py` 只檢查 tag 對得上編號。
- **PM 那層讀者消失**（D2）。對焦會議上要讀 G/W/T。如果實測發現這真的礙事，回頭加
  一層比現在就預留一層便宜。

---

## 不在範圍

- `bdd-clarify` → `bdd-discovery`、`clarify-loop` → `clarify`（獨立的一份工作）
- 把 `身為 P-<n>` 與 `另外依賴` 變成有檢查的
- `runs/` 底下既有的 v2 產物遷移
