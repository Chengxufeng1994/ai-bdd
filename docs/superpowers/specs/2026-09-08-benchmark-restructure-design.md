# 把實測場拆成 benchmark／skeleton／runs

**日期**：2026-09-08
**狀態**：設計已確認，待實作
**動的東西**：`lab/` 底下全部，加上 5 份指向它的文件

---

## 要解的問題

`lab/go/skeleton/` 一個目錄同時扮演三個角色，而它們的壽命完全不同：

| 角色 | 內容 | 隨什麼變 |
| --- | --- | --- |
| harness | `Makefile`、`go.mod`、分層規則、godog runner、`api/` | 語言 |
| 情境章程 | 為什麼選訓練容量、10 個模糊點、及格標準 | 應用 |
| 活的工作區 | `internal/`、`features/`、`specs/` | 一次執行 |

症狀已經寫在 README 裡：「Once IMPLEMENT fills `internal/`, the code belongs to
one scenario — so **archive a completed run before starting another** rather than
trying to keep two live.」——「跑第二個情境前要先封存」不是設計，是缺少結構時的
變通辦法。

另一個症狀是跨目錄引用：`prompts/1-fitness-tracker-clarify.md` 的及格標準寫著
「最低門檻（**見 README**）」，也就是**情境知識躺在語言 harness 裡**。

### 三個驅動

1. **同時留多個情境** — fitness 與 parking 並存，不必先封存
2. **評估可比較** — 改進 skill 之後能知道有沒有變好
3. **名實相符** — 目錄名講實話

---

## 已確認的決策

### D1 — 複本模型：純複本，快照語意

`benchmark/skeleton/<lang>/` 是範本，每次跑 `cp -r` 一份到 `runs/`。**沒有同步
機制，沒有東西流回範本。**

凍結不是缺陷而是 benchmark 的定義：harness 若在跑到一半被換掉，比較就失效。

複本**不改 module 名**（`module skeleton` 留著）。各 benchmark 是互不相干的獨立
module，Go 不在意重名；`Makefile` 寫死的 `-X skeleton/pkg/version.value` 因此
繼續有效，而 `verify-stamp` 正是守它的 gate。複製動作是純粹的 `cp -r`，沒有事後
修補步驟可以忘記。

漂移檢查（`make verify-skeleton-drift`）**現在不做**。活的 benchmark 是零個，
現在建「哪些檔案算 harness」的清單只會是猜的；等第一次「範本修好了要回填到跑到
一半的 benchmark」真的發生再說。這條沿用 repo 既有的 rule of three。

### D2 — 跑出來的產物絕對不進版控

依 `d699aa9` 已決議且有理由的約束：

> An eval whose answers sit beside it in git is one a future run can read instead of
> derive, and a derived answer and a copied one look identical in the artifact.

所以 `runs/` 必須 gitignored。這一條否決了「`benchmark/go-<app>/` 裡裝著 `specs/`
與 `features/` 並進版控」的形狀。

### D3 — 評分方案與輸入同檔，靠章節隔離

同一則 `d699aa9` 的另一半：

> The ambiguity list and pass criteria stay in `prompts/` on purpose — those are the
> **marking scheme, not the input**.

評分方案必須進版控（不然沒得比較），但它**同時就是答案**。兩個需求直接衝突，
解法是在同一個檔案裡劃一條餵食線：**只有 `## 輸入` 底下的內容進 skill 的
context**，其餘給人看。

因此「及格標準該不該跟 prompt 分居」是假問題——它們必須同檔。分居反而會把
「怎麼比較兩次執行」這種語言無關的東西複製到每個語言底下。

### D4 — case 檔就是計分板，不另設 `results/`

`953e1b6`（record both eval runs side by side）已經證明工作流：兩次跑的結果寫進
case 檔的 `## 實測基準`。第三個軸（第幾次跑）**由散文承載，不由目錄承載**——
要比較的不是兩棵產物樹，是「這次有沒有問到自體重」這種判定。

### D5 — 檔名保留 `-clarify`，去掉數字前綴

`-clarify` 帶真資訊（這是測 CLARIFY 的 case，將來會有測 SPEC 的）；`1-`／`2-`
只是建立順序，git 已經記著。等第一個 SPEC case 出現，再分成 `cases/clarify/`、
`cases/spec/`。

---

## 目標結構

```
ai-bdd/
├── skills/                            產品
├── docs/                              bdd.md、tdd.md、sdd.md、ai-sdlc.md
├── CLAUDE.md  PLAN.md  README.md
│
├── benchmark/
│   ├── README.md                      一個 case 是什麼、餵食線、快照規則、污染規則
│   ├── cases/                         語言無關 · 進版控 · 每檔自足
│   │   ├── fitness-tracker-clarify.md
│   │   ├── fitness-tracker-clarify-prd.md
│   │   └── parking-billing-clarify.md
│   └── skeleton/
│       └── go/                        範本，每次跑複製一份
│
└── runs/                              gitignored：skeleton 複本 ＋ 那次的產物
    ├── go-fitness-tracker/
    └── go-parking-billing/
```

曾考慮過把三者攤平在 root（`prompts/`、`skeleton/`、`benchmark/`），那會讓 root
變成 6 個目錄，產品 `skills/` 得跟三個測試目錄搶注意力。收進 `benchmark/` 傘下
之後是 4 個。（已確認 plugin 只認 `skills/`、`commands/`、`agents/`、`hooks/`
四個慣例目錄，所以兩種擺法都不會被誤載入——這不是安全考量，是可讀性考量。）

`runs/` 取代 `eval-runs/`：兩者是同一個概念（進 repo tree 方便瀏覽、不進版控），
合併成一個名字。`eval-runs/` 目前是空的，刪除不會掉東西。

### case 檔的章節（既有結構，不變）

```
## 輸入              ← 唯一餵給 skill 的部分
## 為什麼是這樣措辭
## 模糊點            ← fitness 這份要新增，內容來自 skeleton README
## 及格標準
## 怎麼比較兩次執行
## 實測基準          ← 每跑一次累積一次
```

---

## harness 與情境的界線

今天這條線是乾淨的，而且骨架自己說了它站哪邊：`features/version.feature` 寫著
「This feature carries no domain meaning. It exists to keep one thin path through the
whole chain alive」；`test/acceptance/steps_test.go` 裡恰好三個 step，全屬這條
version 情境。

| 檔案 | 歸屬 |
| --- | --- |
| `Makefile`、`go.mod/sum`、`.golangci.yml`、`.mockery.yml`、`generate.go`、`api/`、`cmd/`、`pkg/` | skeleton |
| `internal/**/doc.go`（分層規則） | skeleton |
| `internal/` 的 `/version` 切片、`features/version.feature`、`test/acceptance/*` | skeleton |
| `docs/ARCHITECTURE.md`、`docs/DATAFLOW.md` | skeleton |
| `README.md` 的 `## Why this domain`（9–31 行） | **搬到 case 檔** |
| `README.md` 的 `## One scenario at a time`（32–51 行） | **刪除**（見下） |
| `README.md` 其餘章節 | skeleton |

### 為什麼 `## One scenario at a time` 是刪不是搬

它描述的是「一次只能裝一個情境所以要先封存」——那是**缺少結構時的變通辦法**，
而這次重構把結構補上了。把變通辦法搬到新家會讓它以規則的姿態活下去，而讀的人
不知道它描述的限制已經不存在。

該節內含一個「一次跑會產出什麼」的清單，那部分**移到 `benchmark/README.md`**；
被刪掉的是封存那條建議本身。

### 界線的規則（寫進 `benchmark/README.md`）

> `benchmark/skeleton/<lang>/` 永遠凍結在「只有走路骨架」的狀態。一次跑的複本
> 可以長出任何東西，**但沒有任何東西流回範本**。要改進範本就直接改
> `benchmark/skeleton/`，之後複製的跑會拿到；正在跑的不會。

這也回答了「IMPLEMENT 之後 `steps_test.go` 混了 fitness step 怎麼辦」——不怎麼
辦，複本本來就該分岔。

---

## 搬遷步驟

前置條件：工作區乾淨。`ae94243` 之後已滿足，**不需要 `git reset`**。

### 步驟 1 — 搬移（全部 `git mv`）

```bash
mkdir -p benchmark/skeleton benchmark/cases
git mv lab/go/skeleton benchmark/skeleton/go
git mv lab/prompts/1-fitness-tracker-clarify.md     benchmark/cases/fitness-tracker-clarify.md
git mv lab/prompts/1-fitness-tracker-clarify-prd.md benchmark/cases/fitness-tracker-clarify-prd.md
git mv lab/prompts/2-parking-billing-clarify.md     benchmark/cases/parking-billing-clarify.md
rmdir lab/prompts lab/go lab
```

case 檔的來源是 `lab/prompts/`，不是 skeleton 底下——`ae94243` 已經把它們搬出去
了。`git log --follow` 因此會走三站：`lab/go/skeleton/prompts/` →
`lab/prompts/` → `benchmark/cases/`。

### 步驟 2 — 內容編輯

| 檔案 | 動作 |
| --- | --- |
| `benchmark/skeleton/go/README.md` | 標題改為 `# benchmark/skeleton/go`；`## Why this domain` 整節搬出；`## One scenario at a time` 刪除（產出清單移入 `benchmark/README.md`） |
| `benchmark/cases/fitness-tracker-clarify.md` | 收下 10 個模糊點成 `## 模糊點`；`## 及格標準` 的「（見 README）」改成自我引用 |
| `benchmark/skeleton/go/docs/ARCHITECTURE.md` | 四處，用內容找而非行號：`## Related documents` 表格裡指向 case 檔的連結、`## 2.` 的資料流圖起點、`## 9.` 的 runbook 第 1 步、`## 10. Project Identification` 的 repo 位置 |
| `PLAN.md` | 實測場整節路徑；`lab/ 不是 plugin 慣例目錄` 那條註記改指 `benchmark/`、`runs/` |
| `README.md` | Layout 樹 |
| `CLAUDE.md` | 全部 `lab/go/skeleton` 路徑；Testbed 一節標題 |
| `.gitignore` | `eval-runs/` → `runs/` |
| `benchmark/README.md` | **新檔**：case 的定義、`## 輸入` 餵食線、快照規則、污染規則、一次跑會產出什麼 |

---

## 驗證

這次重構全部是連結，所以驗證的重點也是連結。

```bash
# 1. 沒有殘留舊路徑
grep -rn "lab/\|skeleton/prompts\|eval-runs" --exclude-dir=.git .     # 應為空

# 2. 每一條相對 markdown 連結都解得到檔案
python3 - <<'EOF'
import re, pathlib, sys
bad = []
for md in pathlib.Path('.').rglob('*.md'):
    if '.git' in md.parts: continue
    for link in re.findall(r'\]\((\.[^)]+)\)', md.read_text(encoding='utf-8')):
        if not (md.parent / link.split('#')[0]).exists():
            bad.append(f"{md}: {link}")
print('\n'.join(bad) or 'all links resolve'); sys.exit(1 if bad else 0)
EOF

# 3. harness 在新路徑下仍然完整（含 lint 與 race test，數分鐘）
cd benchmark/skeleton/go && make verify

# 4. plugin 未受影響
claude plugin validate . --strict

# 5. 腳本對空骨架仍然「正確地失敗」
python3 skills/bdd-clarify/scripts/status.py benchmark/skeleton/go     # 預期 exit 1 並指名缺什麼
```

第 2 項是唯一真正需要自動化的：這次沒有任何程式碼被改動，build 與 lint 幾乎
不可能失敗；會壞的是相對連結，而壞掉的 markdown 連結**不會讓任何指令非零離開**
——GitHub 顯示 404、編輯器沒反應、CI 全綠。

第 5 項不是形式。`ai-sdlc.md` 說這條鏈的特徵性失敗是「靜默地錯」而不是「壞掉」，
所以驗證裡必須有一項是確認**該紅的地方還在紅**。搬完後腳本若突然 exit 0，那才是
搞砸了。

---

## 提交切法

| # | 標題 | 內容 |
| --- | --- | --- |
| 1 | `refactor: restructure the testbed into benchmark/ and runs/` | 全部 `git mv`、路徑更新、`.gitignore`、新的 `benchmark/README.md` |
| 2 | `docs(benchmark): move the fitness ambiguity list into its case file` | skeleton README 搬出一節刪一節、case 檔收下 |

1 是機械性的（搬移＋改路徑），2 是編輯性的（散文在檔案間移動），各自能獨立讀懂
也能獨立 revert。

---

## 明確不做的事

- **漂移檢查**（`make verify-skeleton-drift`）— 見 D1，等它真的痛
- **第二個語言的 skeleton** — 這次只讓結構容得下 `benchmark/skeleton/python/`，不建立它
- **全面修 `docs/bdd/` → `specs/` 的既有陳舊引用** — `5dacedc` 已把該目錄改名為
  `specs/`，但 skeleton 的 `README.md` 與 `docs/ARCHITECTURE.md` 都還寫著舊名。
  這是既有缺陷、與本次搬移無關，**另開一筆修**。

  **例外：搬進新檔的那一份要就地改對。** `## One scenario at a time` 裡的產出清單
  會被移進 `benchmark/README.md`，若照抄就等於把已知錯誤的路徑寫進一個全新檔案，
  比留在舊檔更糟。搬的時候直接寫 `specs/`。剩下未搬動的那些留給另一筆
- **改任何 skill 或腳本** — 已確認 `skills/` 底下三支 Python 都沒有寫死 `lab/`、
  `skeleton` 或 `testbed`，本次重構不需要碰它們
