# CAPTURE → PRD 流程重排 實作計畫

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 新增 CAPTURE 步驟與 `bdd-capture` skill，把 `clarify-log.md`／`input.md` 併成一份從 draft 長到 complete 的 `prd.md`，`spec.md` 瘦成只裝技術設計，兩支腳本改讀 `prd.md` 並新增來源標記檢查。

**Architecture:** `prd.md` 的機械格式**完全沿用**今天 `spec.md` 的業務節（`### US-n · slug`、`#### FR`／`#### NFR`／`#### Q`、`## Open Questions` 五欄表、`← PRD §x`／`← Q-n` 來源標記），只是換檔名、多一節 `## 原文`。所以 `status.py` 與 `check_spec.py` 的解析器不動，只改 glob 與加一項檢查。Skill 的改動是**誰在什麼時候寫哪一節**：CAPTURE 建檔、DISCOVERY 填 FR/AC/Q、`example-mapping` 翻狀態、`bdd-spec` 只寫技術節。

**Tech Stack:** Markdown skill 文件（繁體中文）、Python 3 標準庫腳本（`unittest`，環境沒有 pytest）、`claude plugin validate`、`skills/skill-rules/scripts/audit_skill.py`。

**Spec:** `docs/superpowers/specs/2026-10-02-capture-prd-pipeline-design.md`

## Global Constraints

- 散文文件一律繁體中文；commit 訊息英文，conventional-commit 標題 ＋ `WHAT:`／`WHY:`／`HOW:` 本文，結尾 `Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>`。
- 搬檔與改名一律 `git mv`，`git log --follow` 要能接上。
- `claude plugin validate .` 必須通過（恰好一個已知警告：根層 `CLAUDE.md`）。
- `python3 skills/skill-rules/scripts/audit_skill.py skills` 每個 skill 必須 `✔ 通過`：子目錄只能是 `rules/`、`references/`、`examples/`、`scripts/`；每個子目錄檔案都要被 SKILL.md 的引用鏈指到；必備 `## 使用時機` 與 `## Skill Boundaries`；description 要有「」中文觸發詞與 `English:`。
- 編號契約不動：`FR-<n>` 全域唯一、`AC-<n>.<m>` 掛在 `FR-<n>` 底下、`.feature` 的 `@rule-<n>`／`@example-<n>.<m>` 回指它們。
- `## Open Questions` 表固定五欄 `| Q | 問題 | 面向 | 狀態 | 答案 |`，狀態值只有 `已答`／`待答`／`n/a`，`面向` 值域是 `status.py` 的 `DIMS` 十六個或 `—`。
- 來源標記三選一不變（`PRD §x`／`Q-<n>`／`推論`），但 **FR 與 AC 這兩種行只准前兩種**。
- 腳本 fail-loud：找不到輸入就指名缺什麼並 exit 1，不從零輸入推結論。
- 每個 task 結尾一個 commit；不用 `--no-verify`。

## Review Focus

1. `prd.md` 裡一條 FR 行尾標 `← 推論`：`check_spec.py` 必須報「沒有來源」而不是放行（Task 2 的測試 `test_inference_marker_counts_as_unsourced`）。
2. `specs/<x>/` 底下只有舊的 `spec.md`、沒有 `prd.md`：`status.py` 與 `check_spec.py` 都要 exit 1 並指名「先跑 bdd-capture」，不能印出空儀表板（Task 1 `test_missing_prd_exits_1_and_names_capture`、Task 2 `test_spec_only_dir_exits_1`）。
3. `prd.md` 存在但 `## 原文` 一節缺席：CAPTURE 的產物不完整，`check_spec.py` 要報（Task 2 `test_missing_verbatim_section_reported`）。
4. 使用者帶著現成 PRD 來跑 `bdd-capture`：原文必須逐字進 `## 原文`，不得改寫（Task 4 的 MUST，無法機械驗，寫進 SKILL.md 的 NEVER）。
5. 改名 `clarify-loop → clarify` 之後仍有檔案引用舊名：Task 12 的 grep 閘門要抓到零筆（`grep -rn "clarify-loop" skills/ docs/ai-sdlc.md PLAN.md CLAUDE.md benchmark/README.md`，允許的例外只有 `PLAN.md` 第 26 行講環境中**外部**同名 skill 的那句）。

---

### Task 0: 設計文件校正

**Files:**
- Modify: `docs/superpowers/specs/2026-10-02-capture-prd-pipeline-design.md`

**Interfaces:**
- Produces: 後面每個 task 引用的五條校正（格式沿用、狀態值、來源標記、技術答案落點、`spec.md` 的補問節）。

- [ ] **Step 1: 在設計文件末尾追加校正節**

在 `## 待實跑驗證` 之前插入：

```markdown
## 2026-10-02 校正：對齊既有格式

寫實作計畫時讀了 `spec-format.md`、`status.py`、`check_spec.py`，發現 D2 畫的
`#### 規則／#### 例子／#### 問題` 格式，今天的 `spec.md` 業務節已經有一份機械等價
且**有解析器在讀**的版本。重畫格式會逼兩支腳本重寫解析器，換不到任何東西。以下
五條覆蓋前面對應的決策：

1. **`prd.md` 的格式 ＝ 今天 `spec.md` 的業務節，原樣搬。** `### US-<n> · <slug>`、
   `#### FR`（`**FR-<n>**` 粗體項目、`- **AC-<n>.<m>**` 巢狀其下）、`#### NFR`、
   `#### Q`（只放待答指標）、`## Open Questions` 五欄表＋決策史小節，一個字不改。
   新增的只有 `## 原文` 一節（取代 `input.md`）。D2 的章節草圖作廢，以
   `skills/bdd-capture/references/prd-format.md` 為準。
2. **狀態值沿用 `澄清中／待對焦／已對焦／已凍結`。** D2 的 `draft` ＝ `澄清中`
   （CAPTURE 建檔時寫），`complete` ＝ `已對焦`。寫入者：`澄清中` CAPTURE；
   `待對焦` DISCOVERY 在前沿空了時寫；`已對焦` `example-mapping` 在需求方點頭後寫；
   `已凍結` 仍無寫入者。
3. **來源標記三選一不變，但 FR 與 AC 只准 `PRD §x` 或 `Q-<n>`。** `推論` 仍可用在
   story 句子、成功指標、persona 的「他要什麼」這類推導；一條規則或例子標 `推論`
   就是「agent 自己補的」，`check_spec.py` 第 14 項抓。`PRD §x` 從此指 `prd.md`
   自己的 `## 原文` 一節。D4 的「只有兩種來源」限縮為「FR／AC 只有兩種來源」。
4. **技術面向的答案住 `prd.md` 的 `## Open Questions`，不另開 `## 限制`。** 它們本來
   就是題，`面向` 欄填技術面向的名字；`bdd-spec` 讀已答的那幾列寫 `spec.md`，
   行尾 `← Q-<n>`。D6 的落點改為此。
5. **`spec.md` 保留一節 `## 回 DISCOVERY 補問`。** 原本 `## Scope — In / Out` 的
   「回 CLARIFY 補問」清單跟著 Scope 搬去 `prd.md` 之後，SPEC 與 FORMULATION 寫到
   一半發現的洞需要一個不是 `prd.md` 的落點（它們不得寫 `prd.md`）。DISCOVERY
   下一輪把這一節的每一條開成 `Q-<n>`。
```

- [ ] **Step 2: Commit**

```bash
git add docs/superpowers/specs/2026-10-02-capture-prd-pipeline-design.md
git commit -m "docs: align the capture/prd design with the existing spec.md machine format" -m "WHAT: Append five corrections to the 2026-10-02 design: prd.md reuses spec.md's business-section format verbatim plus a 原文 section, status values stay 澄清中/待對焦/已對焦, FR/AC source markers are limited to PRD §x or Q-n, technical answers live in Open Questions, and spec.md keeps a 回 DISCOVERY 補問 section.

WHY: The parsers in status.py and check_spec.py already read exactly the shape D2 sketched under new names. Redrawing the format would force a parser rewrite for no gain.

HOW: Corrections appended as a dated section so the original decisions stay readable.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 1: `status.py` 搬家並改讀 `prd.md`

**Files:**
- Create: `tests/_load.py`
- Create: `tests/test_status.py`
- Move: `skills/bdd-spec/scripts/status.py` → `skills/bdd-discovery/scripts/status.py`（`git mv`）
- Modify: `skills/bdd-discovery/scripts/status.py`

**Interfaces:**
- Consumes: 無。
- Produces: `status.main(root: Path) -> int`（簽名不變）；glob 改成 `specs/*/prd.md`；找不到時印 `{specs} 底下沒有任何 prd.md —— 先跑 bdd-capture`。`tests/_load.py` 的 `load(rel_path: str) -> module` 供 Task 2 重用。

- [ ] **Step 1: 寫載入輔助與失敗的測試**

`tests/_load.py`：

```python
"""把 skills/<x>/scripts/*.py 當模組載入。腳本不是套件，用路徑載。"""
import importlib.util
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def load(rel: str):
    path = ROOT / rel
    spec = importlib.util.spec_from_file_location(path.stem, path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module
```

`tests/test_status.py`：

```python
import contextlib
import io
import tempfile
import unittest
from pathlib import Path

from _load import load

status = load("skills/bdd-discovery/scripts/status.py")

PRD = """# PRD: parking

## Document Overview

**狀態**：澄清中　**最後更新**：2026-10-02

## User Stories

### US-1 · visitor-billing

身為 P-1（訪客），我想在出場前知道要付多少 ← 推論

#### FR

**FR-1**  When 訪客車出場, the system shall 依停留時長計費。 ← PRD §3
- **AC-1.1**  停 29 分 → 收 0 元 ← Q-1

#### Q

**Q-2**  跨午夜怎麼算 ← 擋住 FR-1

## Open Questions

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q-1 | 免費額度多長 | 邊界 | 已答 | 30 分 |
| Q-2 | 跨午夜怎麼算 | 時間 | 待答 | — |
| Q-3 | 機車費率 | — | n/a | 本期不做機車 |

## 原文

§1 停車場要收費。
"""


def run(root: Path) -> tuple[int, str]:
    buf = io.StringIO()
    with contextlib.redirect_stdout(buf):
        code = status.main(root)
    return code, buf.getvalue()


class StatusReadsPrd(unittest.TestCase):
    def test_counts_rows_from_prd(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "2026-10-02-parking" / "prd.md").write_text(PRD, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 0)
            self.assertIn("2026-10-02-parking", out)
            self.assertIn("待答（1）", out)
            self.assertIn("2026-10-02-parking/Q-2", out)

    def test_missing_prd_exits_1_and_names_capture(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "2026-10-02-parking" / "spec.md").write_text("# old", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("沒有任何 prd.md", out)
            self.assertIn("bdd-capture", out)


if __name__ == "__main__":
    unittest.main()
```

- [ ] **Step 2: 跑測試確認失敗**

Run: `python3 -m unittest discover -s tests -v`
Expected: 兩個測試都 ERROR，訊息含 `FileNotFoundError` 或 `No such file`（`skills/bdd-discovery/scripts/status.py` 還不存在）。

- [ ] **Step 3: 搬檔並改 glob**

```bash
mkdir -p skills/bdd-discovery/scripts
git mv skills/bdd-spec/scripts/status.py skills/bdd-discovery/scripts/status.py
```

在 `skills/bdd-discovery/scripts/status.py` 做四處替換：

1. 模組 docstring 第 6 行起：`` `specs/<date>-<feature>/spec.md` 的 `## Open Questions` 表 `` → `` `specs/<date>-<feature>/prd.md` 的 `## Open Questions` 表 ``；同段 `` 底下沒有任何 `spec.md` `` → `` 底下沒有任何 `prd.md` ``；`` 單一 `spec.md` 也有三種算「壞」的情況 `` → `` 單一 `prd.md` 也有三種算「壞」的情況 ``；`CLARIFY 沒問過任何問題` → `DISCOVERY 沒問過任何問題`（docstring 與 `main()` 裡的訊息各一處）。
2. `parse_open_questions` 的 docstring 第一行：`spec.md 的` → `prd.md 的`。
3. `main()`：

```python
    specs_files = sorted(specs.glob("*/prd.md"))
    if not specs_files:
        print(f"{specs} 底下沒有任何 prd.md —— 先跑 bdd-capture")
        return 1
```

4. `TECH_DIMS` 上方註解 `references/technical-probes.md` 的路徑說明不變（現在同 skill）；`BIZ_DIMS` 註解不變。變數名 `spec_path`／`specs_files` 保留，只是讀的檔變了——不做無關的改名。

- [ ] **Step 4: 跑測試確認通過**

Run: `python3 -m unittest discover -s tests -v`
Expected: `Ran 2 tests ... OK`

- [ ] **Step 5: Commit**

```bash
git add tests/_load.py tests/test_status.py skills/bdd-discovery/scripts/status.py
git commit -m "feat(bdd-discovery): move status.py here and read prd.md" -m "WHAT: git mv status.py from bdd-spec/scripts to bdd-discovery/scripts, glob specs/*/prd.md instead of spec.md, and name bdd-capture as the step to run when none exists. Add a unittest harness under tests/.

WHY: The dashboard measures DISCOVERY progress and prd.md is now the file DISCOVERY writes; the Open Questions table it parses is unchanged.

HOW: Parser untouched; only the glob, the missing-input message and the docstring change. Tests load the script by path since it is not a package.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 2: `check_spec.py` 改讀 `prd.md`，新增來源標記與 `## 原文` 檢查

**Files:**
- Create: `tests/test_check_spec.py`
- Modify: `skills/bdd-formulation/scripts/check_spec.py`

**Interfaces:**
- Consumes: `tests/_load.py` 的 `load`。
- Produces: `check(root: Path) -> int`（不變）；新函式 `unsourced_rules(text: str) -> list[str]` 與 `has_verbatim_section(text: str) -> bool`；glob `specs/*/prd.md`。

- [ ] **Step 1: 寫失敗的測試**

`tests/test_check_spec.py`：

```python
import contextlib
import io
import tempfile
import unittest
from pathlib import Path

from _load import load

check_spec = load("skills/bdd-formulation/scripts/check_spec.py")

GOOD = """# PRD: parking

## User Stories

### US-1 · visitor-billing

身為 P-1（訪客），我想在出場前知道要付多少 ← 推論

#### FR

**FR-1**  When 訪客車出場, the system shall 依停留時長計費。 ← PRD §3
- **AC-1.1**  停 29 分 → 收 0 元 ← Q-1

#### Q

**Q-2**  跨午夜怎麼算 ← 擋住 FR-1

## Open Questions

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q-1 | 免費額度多長 | 邊界 | 已答 | 30 分 |
| Q-2 | 跨午夜怎麼算 | 時間 | 待答 | — |

## 原文

§3 前 30 分鐘免費。
"""


def run(root: Path) -> tuple[int, str]:
    buf = io.StringIO()
    with contextlib.redirect_stdout(buf):
        code = check_spec.check(root)
    return code, buf.getvalue()


def write_prd(root: Path, text: str) -> None:
    d = root / "specs" / "2026-10-02-parking"
    d.mkdir(parents=True)
    (d / "prd.md").write_text(text, encoding="utf-8")


class UnsourcedRules(unittest.TestCase):
    def test_sourced_lines_are_clean(self):
        self.assertEqual(check_spec.unsourced_rules(GOOD), [])

    def test_inference_marker_counts_as_unsourced(self):
        text = GOOD.replace("← PRD §3", "← 推論")
        self.assertEqual(check_spec.unsourced_rules(text),
                         ["**FR-1**  When 訪客車出場, the system shall 依停留時長計費。 ← 推論"])

    def test_missing_marker_on_ac_counts_as_unsourced(self):
        text = GOOD.replace("- **AC-1.1**  停 29 分 → 收 0 元 ← Q-1",
                            "- **AC-1.1**  停 29 分 → 收 0 元")
        self.assertEqual(check_spec.unsourced_rules(text),
                         ["- **AC-1.1**  停 29 分 → 收 0 元"])


class CheckReadsPrd(unittest.TestCase):
    def test_good_prd_passes_in_spec_only_mode(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_prd(root, GOOD)

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertIn("全部通過", out)

    def test_unsourced_fr_fails(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_prd(root, GOOD.replace("← PRD §3", "← 推論"))

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("沒有來源", out)
            self.assertIn("FR-1", out)

    def test_missing_verbatim_section_reported(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_prd(root, GOOD.split("## 原文")[0])

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("## 原文", out)

    def test_spec_only_dir_exits_1(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "x").mkdir(parents=True)
            (root / "specs" / "x" / "spec.md").write_text("# old", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("沒有 prd.md", out)
            self.assertIn("bdd-capture", out)


if __name__ == "__main__":
    unittest.main()
```

- [ ] **Step 2: 跑測試確認失敗**

Run: `python3 -m unittest tests.test_check_spec -v` （在 repo 根目錄，先 `cd tests` 或用 `-s tests`：`python3 -m unittest discover -s tests -p 'test_check_spec.py' -v`）
Expected: `UnsourcedRules` 三個 ERROR（`AttributeError: module has no attribute 'unsourced_rules'`）；`CheckReadsPrd` 四個 FAIL（glob 還讀 `spec.md`，印 `沒有 spec.md`）。

- [ ] **Step 3: 實作**

在 `check_spec.py` 的 `orphan_frs()` 之後、`step_templates()` 之前加入：

```python
SOURCED = re.compile(r"←\s*(PRD §\S+|Q-\d+)\s*$")


def unsourced_rules(text: str) -> list[str]:
    """`## User Stories` 段裡沒有 `← PRD §x` 或 `← Q-<n>` 來源的 FR／AC 行。

    來源標記三選一，但規則與例子只准前兩種：`推論` 在這兩種行上的意思正是
    「agent 自己補了一條合理的」——那是 DISCOVERY 不寫檔那條舊紀律要擋的
    東西，現在改用這個檢查擋。story 句子、NFR、persona 的「他要什麼」不在
    此列，它們本來就允許推導。

    只認契約宣告的兩種行首：`**FR-<n>**` 與 `- **AC-<n>.<m>**`。散文裡提到
    FR-2 不算一條規則。
    """
    section = re.search(
        r"^## User Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return []
    out: list[str] = []
    for line in section.group(1).splitlines():
        is_rule = re.match(r"^\*\*FR-\d+\*\*", line) is not None
        is_example = re.match(r"^\s*-\s*\*\*AC-\d+\.\d+\*\*", line) is not None
        if (is_rule or is_example) and not SOURCED.search(line.rstrip()):
            out.append(line.strip())
    return out


def has_verbatim_section(text: str) -> bool:
    """`## 原文` 一節在不在。

    它是 `PRD §x` 標記唯一的依據。沒有這一節，每個 `← PRD §x` 都指向一段從未
    寫下的話，而且不會報錯——只會悄悄指錯。CAPTURE 建檔時寫，之後不得改寫。
    """
    return re.search(r"^## 原文\s*$", text, re.M) is not None
```

在 `check()` 裡做三處改動：

1. glob 與訊息：

```python
    specs = sorted(specs_dir.glob("*/prd.md"))
    if not specs:
        print(f"{specs_dir} 底下沒有 prd.md —— 先跑 bdd-capture")
        return 1
```

2. 進入 `for spec_path in specs:` 迴圈、讀完 `stext` 之後、`frs = frs_in_spec(stext)` 之前加：

```python
        if not has_verbatim_section(stext):
            print(f"✗ {spec_path.parent.name} 的 prd.md 沒有 `## 原文` 一節，"
                  f"每個 `← PRD §x` 都指向不存在的話 —— CAPTURE 的產物不完整")
            problems += 1

        unsourced = unsourced_rules(stext)
        if unsourced:
            print(f"✗ {spec_path.parent.name} 這些 FR／AC 沒有來源（只准 `← PRD §x` "
                  f"或 `← Q-<n>`，`推論` 不算）：")
            for line in unsourced:
                print(f"    {line}")
            problems += 1
```

3. 其餘訊息裡的字串替換（只動訊息與 docstring，不動邏輯）：`spec.md 沒有 \`## User Stories\`` → `prd.md 沒有 \`## User Stories\``；`—— 先跑 bdd-spec` 兩處 → `—— 先跑 bdd-capture`；`SPEC 還沒切 story` → `CAPTURE 還沒切 story`；`SPEC 無從寫出場景` 兩處 → `FORMULATION 無從寫出場景`；`只跑 spec.md 自己的八項檢查` → `只跑 prd.md 自己的十項檢查`；模組 docstring 開頭改成「稽核 .feature 與 prd.md 的一致性」，「檢查十三件事」改「檢查十五件事」並在清單末尾加：

```
 14. 來源——`## User Stories` 底下每條 FR 與 AC 行尾都要有 `← PRD §x` 或
     `← Q-<n>`。`推論` 不算：規則或例子標推論，就是 agent 自己補的。
 15. `## 原文` 一節在不在——它是 `PRD §x` 唯一的依據。
```

docstring 其餘 `spec.md` 一律改 `prd.md`、`SPEC 交件到 FORMULATION 開跑之間` 改 `DISCOVERY 交件到 FORMULATION 開跑之間`、`第 6-13 項只讀 spec.md` 改 `第 6-15 項只讀 prd.md`。函式 docstring 裡的 `spec.md` 也全改 `prd.md`（`frs_in_spec`、`v2_forms_in_spec`、`qs_in_table`、`stories_in_spec` 等）。函式名保留 `_in_spec`——改名會牽動測試與 bdd-spec SKILL.md 的引用，不值得。

- [ ] **Step 4: 跑測試確認通過**

Run: `python3 -m unittest discover -s tests -v`
Expected: `Ran 9 tests ... OK`

- [ ] **Step 5: Commit**

```bash
git add tests/test_check_spec.py skills/bdd-formulation/scripts/check_spec.py
git commit -m "feat(bdd-formulation): check_spec reads prd.md and rejects unsourced rules" -m "WHAT: Glob specs/*/prd.md, add check 14 (every FR/AC line ends with ← PRD §x or ← Q-n; 推論 counts as unsourced) and check 15 (## 原文 section present), rename the step names in messages.

WHY: The design replaces the 'DISCOVERY must not write' guard with 'every rule names its source', which only works if a script enforces it. 原文 is what PRD §x points at.

HOW: Two new pure functions with unit tests; parsers for stories, FR, AC and Open Questions untouched.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 3: `prd-format.md` 建立、`spec-format.md` 瘦身、persona 與範例搬家

**Files:**
- Create: `skills/bdd-capture/references/prd-format.md`
- Move: `skills/bdd-spec/references/persona-definition.md` → `skills/bdd-capture/references/persona-definition.md`
- Create: `skills/bdd-capture/examples/minimal-prd.md`
- Modify: `skills/bdd-spec/references/spec-format.md`
- Modify: `skills/bdd-spec/examples/minimal-spec.md`
- Delete: `skills/bdd-discovery/references/clarify-log-format.md`

**Interfaces:**
- Produces: `prd-format.md` 是 CAPTURE／DISCOVERY／`example-mapping`／FORMULATION 共讀的格式契約；`spec-format.md` 只剩技術節。這個 task 結束時 `audit_skill.py` 對 `bdd-capture` 會報 S1（還沒有 SKILL.md），那是預期，Task 4 補。

- [ ] **Step 1: 搬 persona 與建立 `prd-format.md`**

```bash
mkdir -p skills/bdd-capture/references skills/bdd-capture/examples
git mv skills/bdd-spec/references/persona-definition.md skills/bdd-capture/references/persona-definition.md
```

`persona-definition.md` 開頭第一段改成：

```markdown
`prd.md` 的 `## Personas` 節的寫法。三個讀者：`bdd-capture` 從原文第一次抽角色時讀它；
`bdd-discovery` 在開場的「識別角色」那一步、或後面某一輪冒出沒見過的角色時讀它；
`bdd-spec` 在 `spec.md` 引用角色時讀它。
```

其餘出現 `spec.md` 的地方全改 `prd.md`（用 `grep -n "spec.md" skills/bdd-capture/references/persona-definition.md` 列出逐一改）。

建立 `skills/bdd-capture/references/prd-format.md`，內容來自 `spec-format.md` 的這幾段**原樣搬過來**（用編輯器剪貼，不重寫）：`## 章節順序`（刪掉技術節那幾列）、`## Document Overview 的三個部分`、`## Scope — In / Out`、`## Personas：在三欄上加兩欄`、`## User Stories：三個分組，AC 巢狀在它的 FR 底下`（含 `### #### Q 只放指標`、`### 一條 FR 服務多則 story`）、`## 來源標記：三選一`、`## 編號規則：AC 的編號跟著它的 FR`、`## 用 EARS 寫規則`、`## AC 就是綠卡`、`## Non-Functional Requirements 與 Constraints：怎麼分邊`、`## Open Questions：表是索引，小節是決策史`（含 `面向` 值域、狀態值、輪次與信心）、`## Further Notes`、`## 角色與詞義：兩條界線`。搬完之後對搬過來的文字做這些修改：

開頭換成：

```markdown
# `prd.md` 的骨架

照抄即可。三個 skill 在不同時刻讀它：`bdd-capture` 建檔時、`bdd-discovery`
每一輪寫答案時、`bdd-formulation` 翻 `.feature` 時。`example-mapping` 讀它攤地圖。

一個 feature 一份，路徑 `specs/<date>-<feature>/prd.md`。一份 `prd.md` 裝這個
feature 的全部 story——story 由 CAPTURE 粗切、DISCOVERY 細切，見各自的 SKILL.md。

它是**一份從 draft 長到 complete 的活文件**：CAPTURE 建檔時只有開頭幾節與空的
story；DISCOVERY 每答一題就往 story 底下長一條 FR 或 AC；紅卡清空、需求方點頭，
狀態翻成 `已對焦`，`bdd-formulation` 與 `bdd-spec` 才讀它。

沒有內容的節**保留標題並寫為什麼空**，不要刪掉。空著的節跟被刪掉的節在檔案上
長得一樣，而它們的意思完全不同。
```

`## 章節順序` 表改成：

```markdown
| 章節 | 裝什麼 | 誰第一次寫 |
| --- | --- | --- |
| `## Document Overview` | 狀態、最後更新、結束方式、版本修訂歷史、核心關係人 | CAPTURE |
| `## Background` | 業務問題與現在什麼在痛 | CAPTURE |
| `## Goal` | 做成之後的樣子 | CAPTURE |
| `## Scope — In / Out` | In 是這次要交付的能力；Out 分兩類，見下 | CAPTURE |
| `## Personas` | `### P-<n>` | CAPTURE（五欄可留「（無資料）」），DISCOVERY 補 |
| `## User Stories` | `### US-<n> · <slug>`，底下三個分組 `#### FR`／`#### NFR`／`#### Q` | CAPTURE 寫標題與句子，DISCOVERY 填三個分組 |
| `## Non-Functional Requirements` | 跨 story 的品質屬性 | DISCOVERY |
| `## Assumptions / Constraints` | 沒驗證過的假設，以及做不到／不准這樣做的外部限制 | CAPTURE 抄原文的假設，DISCOVERY 改寫成問過的 |
| `## Success Metrics` | 能量化就寫數字；量不出來就寫「怎麼知道做成了」 | CAPTURE |
| `## Open Questions` | 表是機械可解析的索引，小節是人讀的決策史 | CAPTURE 填需求方自己列的開放問題，DISCOVERY 填其餘 |
| `## Further Notes` | 放不進其他節、但下一個接手的人會想知道的 | 任何一步 |
| `## 原文` | 需求方帶來的描述，逐字 | CAPTURE，之後**不得改寫** |

順序固定，不因為某節內容少而調換或省略標題。技術設計（seam、API、schema、
測試層級、風險）不在這份文件裡，那是 `spec.md`，格式見
`../../bdd-spec/references/spec-format.md`。
```

`## Document Overview` 範例的狀態列改成 `**狀態**：澄清中　**最後更新**：2026-09-10　**結束方式**：—`，並在狀態值表之後加一段：

```markdown
**結束方式**是 DISCOVERY 最近一輪怎麼停的，四選一：`收斂`／`阻塞`／`零進展`／
`喊停`，CAPTURE 建檔時寫 `—`。它跟 `狀態` 不同：`狀態` 說這份文件走到哪，
`結束方式` 說上一輪為什麼停在那裡。`阻塞` 要在 `## Open Questions` 對應小節寫清楚
卡在誰身上。
```

狀態值表的「誰寫」欄改成：`澄清中｜CAPTURE 建檔時`、`待對焦｜DISCOVERY 在前沿空了時`、`已對焦｜example-mapping，依使用者的答覆`、`已凍結｜尚無寫入者`。那段「這四個值講的是這份文件自己走到哪……」的三段舊版辯解全刪，換成一句：`這四個值講的是這份文件自己走到哪。下游的閘門：bdd-formulation 與 bdd-spec 只讀 已對焦 的 prd.md。`

核心關係人那兩段「整張表照抄 `clarify-log.md`……」改成：

```markdown
**核心關係人由 CAPTURE 登記第一批（需求方自己與他點名的人），DISCOVERY 每冒出
一個新的受訪者就補一列。** `決策權` 這一格只有問的人寫得出來——`bdd-spec` 把整張表
原樣抄進 `spec.md`，不增列也不改寫。
```

`## Scope — In / Out` 的 Out 改成兩類（刪掉「回 CLARIFY 補問」那一類），並加一句：`推不出來、需要回頭問的缺口不在這裡列——它是一張紅卡，開進 \`## Open Questions\`。`

`## 來源標記：三選一` 表的 `PRD §x` 那列改成：`` 同一份 `prd.md` 的 `## 原文` 這樣寫（引用它的章節或段落）——需求方帶來的原始輸入的逐字副本 ``。範圍清單裡刪掉技術節（Tech Preferences、Known Boundaries、Implementation Decisions、Testing Decisions 那幾行）。在 `MUST: 這個標記不是裝飾` 之後加：

```markdown
MUST: `## User Stories` 底下的 FR 與 AC **只准 `PRD §x` 或 `Q-<n>`**。一條規則或
例子標 `推論`，意思就是「沒有人說過這句話，是我補的」——那正是 DISCOVERY 不得
做的事。`check_spec.py` 第 14 項會抓。story 的句子、NFR、persona 的「他要什麼」
仍可標 `推論`。
```

`## Open Questions` 小節範例的狀態列後面加一行說明：小節可以多兩個選填欄位 `**誰答的**：<關係人>` 與 `**歸屬**：<這個答案要進哪一節>`，`歸屬` 的值域是本文件的章節名或 `spec.md` 的章節名（`Implementation Decisions`／`Testing Decisions`／`Risks`／`Tech Preferences & Constraints`），技術面向的答案靠這一欄讓 `bdd-spec` 知道該抄進哪一節。

`## 角色與詞義：兩條界線` 裡的 `clarify-log.md` 那句改成：`答案進 \`## Open Questions\` 的表與小節`。

最後加一節：

```markdown
## 原文

```markdown
## 原文

需求方 2026-10-02 口述，由 bdd-capture 逐字記錄；§ 編號是記錄時加的，原文沒有。

§1 停車場要開始收費，訪客跟月租戶分開算。
§2 前半小時免費，之後每小時 30 元。
§3 月租戶刷卡就直接出去。
```

MUST: 逐字。需求方帶著現成文件來，整份貼進來；口述的，照講的記，不潤飾、
不補主詞、不把「大概半小時」改成「30 分鐘」——那個「大概」正是下一步要問的題。

`§<n>` 是記錄時加的段落編號，給 `← PRD §n` 用。原文自己有章節編號就沿用，
不另編。

MUST NOT: 建檔之後改寫這一節。要補充就在 `## Open Questions` 開一題，答案
長成 FR；原文被改寫，每個 `← PRD §x` 從此指向一句沒人說過的話，而且不會報錯。
```

- [ ] **Step 2: 瘦身 `spec-format.md`**

保留並修改：

開頭改成：

```markdown
# `spec.md` 的骨架

照抄即可。讀這份的時機：`bdd-spec` 準備寫 `spec.md` 的時候。

一個 feature 一份，路徑 `specs/<date>-<feature>/spec.md`。它只裝**技術設計**：
seam、模組、介面、API 契約、schema、測試層級、風險。業務內容（story、FR、AC、
角色、範圍、開放問題）住同目錄的 `prd.md`，格式見
`../../bdd-capture/references/prd-format.md`——本文件不重複那些節。

`spec.md` 是拋棄式快照，實作一開始就會過期；`prd.md` 活得比它久。壽命不同的
東西不放同一個檔。

沒有內容的節**保留標題並寫為什麼空**，不要刪掉。
```

`## 章節順序` 表改成只剩：`## Document Overview`（最後更新、版本修訂歷史、核心關係人——照抄 `prd.md`）、`## Tech Preferences & Constraints`、`## Known Boundaries`、`## Implementation Decisions`、`## Testing Decisions`、`## Risks`、`## 回 DISCOVERY 補問`、`## Further Notes`。

`## Document Overview 的三個部分` 改成兩個部分：刪掉 `狀態` 那一列與整段狀態值的討論（就緒住在 `prd.md`）；範例狀態列改成 `**最後更新**：2026-09-10　**對應 PRD**：prd.md @ 已對焦`；核心關係人段改成一句：`整張表照抄 prd.md 的 ## Document Overview，不增列、不改寫 決策權；prd.md 缺那張表就寫一行說缺了、回 bdd-discovery 登記，不要自己補。`

刪掉整段：`## Scope — In / Out`、`## Personas`、`## User Stories`（含兩個子節）、`## 編號規則`、`## 用 EARS 寫規則`、`## AC 就是綠卡`、`## Non-Functional Requirements 與 Constraints`、`## Open Questions`（含三個子節）、`## 角色與詞義`。

`## 來源標記：三選一` 保留，範圍清單改成只列技術節：`## Tech Preferences & Constraints` 的每條、`## Known Boundaries` 三張清單的每一條、`## Implementation Decisions` 的每條決定、`## Testing Decisions` 的 Seam 決定、`## Further Notes`。表的 `PRD §x` 列改成指 `prd.md` 的 `## 原文`；`Q-<n>` 列改成指 `prd.md` 的 `## Open Questions`。刪掉 `## Personas 在範圍內` 與 `## Assumptions / Constraints 是這條豁免的例外` 兩段。

在 `## Risks` 之後、`## Further Notes` 之前新增：

```markdown
## 回 DISCOVERY 補問

寫技術設計時發現的洞：規格推不出來、兩條 FR 打架、seam 要的前置狀態沒人問過。
一條一列，寫清楚**影響什麼定不下來**：

```markdown
## 回 DISCOVERY 補問

- FR-2 沒說月租到期當天算不算有效，seam 的 Given 寫不出來 —— 影響 AC-2.1 的前置狀態
```

MUST NOT: 在這裡順手決定掉。`bdd-spec` 與 `bdd-formulation` 都不得寫 `prd.md`，
所以洞只能記在這裡；DISCOVERY 下一輪把每一條開成 `Q-<n>`，答案長成 FR／AC 之後，
這一節對應的列才刪。沒有洞就寫「（無）」。
```

`## domain-model.md 的骨架` 保留不動。

- [ ] **Step 3: 拆範例**

建立 `skills/bdd-capture/examples/minimal-prd.md`：把 `minimal-spec.md` 的 `## Document Overview`（狀態列改成 `**狀態**：已對焦　**最後更新**：2026-09-10　**結束方式**：收斂`）、`## Background`、`## Goal`、`## Scope — In / Out`（刪掉「回 CLARIFY 補問」小清單）、`## Personas`、`## User Stories`、`## Non-Functional Requirements`、`## Assumptions / Constraints`、`## Success Metrics`、`## Open Questions`、`## Further Notes` 原樣搬過來，標題第一行改 `# PRD: parking-fee-collection`，末尾加：

```markdown
## 原文

管委會主委 2026-09-09 帶來的一頁需求，逐字；§ 編號是記錄時加的。

§1 社區停車場要開始對訪客收費，月租戶不受影響。
§2 訪客前 30 分鐘免費，之後每 30 分鐘 30 元，不足 30 分以 30 分計。
§3 訪客出場在繳費機付款，付完才開柵欄。
§4 月租戶刷卡出場，柵欄直接開，不計費。
§5 消防車進出要能不受柵欄阻擋（主委說有規定，但不確定內容）。
```

`minimal-spec.md` 只留 `## Document Overview`（改成兩部分：最後更新與對應 PRD、版本修訂歷史、核心關係人）、`## Tech Preferences & Constraints`、`## Known Boundaries`、`## Implementation Decisions`、`## Testing Decisions`、`## Risks`、新增 `## 回 DISCOVERY 補問`（寫 `（無）——FR 與 AC 在 prd.md 已對焦時全部有前置狀態`）、`## Further Notes`。標題第一行改 `# Spec: parking-fee-collection`。

- [ ] **Step 4: 刪 `clarify-log-format.md`**

```bash
git rm skills/bdd-discovery/references/clarify-log-format.md
```

- [ ] **Step 5: 驗證搬運沒有漏字**

Run: `grep -c "" skills/bdd-capture/references/prd-format.md skills/bdd-spec/references/spec-format.md` 兩份行數加總應接近原本 `spec-format.md` 的行數（原本約 760 行；差距在 ±80 行內，多的是新增段落，少的是刪掉的舊版辯解）。
Run: `grep -n "clarify-log\|input.md\|回 CLARIFY" skills/bdd-capture/references/*.md skills/bdd-spec/references/spec-format.md skills/bdd-capture/examples/*.md skills/bdd-spec/examples/*.md`
Expected: 零筆。

- [ ] **Step 6: Commit**

```bash
git add -A skills/bdd-capture skills/bdd-spec/references skills/bdd-spec/examples skills/bdd-discovery/references
git commit -m "refactor(bdd-spec): split spec-format into prd-format (business) and spec-format (technical)" -m "WHAT: Move the business sections of spec-format.md verbatim into skills/bdd-capture/references/prd-format.md with a new ## 原文 section and per-section writers; leave spec-format.md with the technical sections plus a ## 回 DISCOVERY 補問 section; git mv persona-definition.md to bdd-capture; split minimal-spec.md into minimal-prd.md and a technical minimal-spec.md; delete clarify-log-format.md.

WHY: prd.md and spec.md have different lifespans and different writers. The machine format of the business sections is unchanged so the scripts keep parsing.

HOW: Cut-and-paste rather than rewrite, so git log --follow on persona-definition.md still works and the rules text stays identical where it moved.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 4: `bdd-capture` SKILL.md

**Files:**
- Create: `skills/bdd-capture/SKILL.md`

**Interfaces:**
- Consumes: `references/prd-format.md`、`references/persona-definition.md`、`examples/minimal-prd.md`（Task 3）。
- Produces: `specs/<date>-<feature>/prd.md`，狀態 `澄清中`，結束方式 `—`。

- [ ] **Step 1: 寫 SKILL.md**

```markdown
---
name: bdd-capture
description: >
  把一個想法或一份現成需求文件寫成 `prd.md` 的半成品——原文逐字存檔、抽出問題／
  預期結果／角色／範圍／成功指標、粗切 story、需求方自己列的開放問題開成紅卡。
  當書記不當分析師：最多問四題，不替 story 填規則與例子。
  BDD 七步流程的 CAPTURE 步驟，產物交給 `bdd-discovery` 問到完整。
  觸發詞：「我有一個想法」「記一下這個需求」「幫我把需求寫下來」「開一份 PRD」
  「建立 prd.md」「capture」「這是我想做的」「需求進來了」「先把它寫成文件」
  「把這份 PRD 收進來」。
  English: capture this idea, write down this requirement, start a PRD,
  create prd.md, intake this requirement document.
---

# CAPTURE — 把需求寫成半成品的 PRD

需求方的一段話、或一份現成文件進來，出去的是一份狀態 `澄清中` 的
`specs/<date>-<feature>/prd.md`：**原文逐字在裡面、開頭幾節從原文抽出來、story
粗切好但底下是空的、需求方自己知道的開放問題已經是紅卡。**

這一步是**書記**。需求方說、你記；原文沒講的留空並說它沒講。追問是下一步
`bdd-discovery` 的事，規則與例子是那一步問出來才寫的。

## 使用時機

- 有人有一個想法，想開始走流程，但還沒有任何文件
- 手上有一份別人寫的 PRD、ticket、或會議紀錄，要收進流程
- `bdd-discovery` 說找不到 `prd.md`

## Skill Boundaries

- 要把需求**問**清楚、找隱含假設、填規則與例子 → 改用 `bdd-discovery`，它讀本 skill 的產物
- 要把一則太大的 story 切開 → 粗切是本 skill 的事，但沿規則細切要等規則問出來，那時用 `story-splitting`
- 要攤四色卡、判定就緒 → 改用 `example-mapping`
- 本 skill 不問十六個面向、不寫 FR／AC、不寫 Gherkin、不碰 `spec.md`

## 參考檔案

- `references/prd-format.md` — `prd.md` 的骨架，每一節誰第一次寫
- `references/persona-definition.md` — 角色五欄怎麼填、哪些可以留「（無資料）」
- `examples/minimal-prd.md` — 一份已對焦的 `prd.md` 長什麼樣（本 skill 交出去的是它的
  澄清中版本：story 底下三個分組空著）

---

## 前置確認（Ask user）

**唯一必要的輸入是需求本身**：一段描述，或一份文件的內容／路徑。

NEVER: 需求還沒拿到就自己想一個來跑。那會產出一份看似完整、實則憑空的
`prd.md`，而錯誤全藏在專業的格式底下。

Ask user: 這個 feature 的 slug（kebab-case，用來命名 `specs/<date>-<slug>/`）。
需求方沒有偏好就從原文的主詞取一個，講出來讓他改。

`specs/<date>-<slug>/prd.md` 已經存在時 **不重建**，告訴對方它在、狀態是什麼，
該跑 `bdd-discovery` 繼續。

---

## 流程：四件事，順序固定

### 1. 原文逐字進 `## 原文`

MUST: 這是第一個動作，不是收尾。需求方帶著文件來，整份貼進去；口述的，照講的記。
每一段加 `§<n>` 編號（原文自己有編號就沿用），這是之後每個 `← PRD §n` 指的東西。

MUST NOT: 潤飾、補主詞、把「大概半小時」寫成「30 分鐘」。那個「大概」正是下一步
要問的題；改掉它，題目就消失了。

### 2. 從原文抽開頭幾節

照 `references/prd-format.md` 的順序填 `## Background`、`## Goal`、`## Scope — In / Out`、
`## Personas`、`## Assumptions / Constraints`、`## Success Metrics`。每一句行尾標
`← PRD §n`——這幾節在這一步**只能**有這一種來源，標記全是 `PRD §` 才對；出現 `推論`
就是你補了原文沒有的話。

原文抽不出來的節，保留標題、寫「原文沒提」。不要替它想一個。

`## Personas` 五欄照 `references/persona-definition.md`：原文有講的填，沒講的寫
「（無資料）」，**「他要什麼」與「現在什麼讓他痛」尤其不要編**。

`## Assumptions / Constraints`：原文裡標題叫「假設」或語氣像既定事實的數字、門檻、
SLA，都抄進來並標 `← PRD §n`——它們是**沒驗證過的**，下一步會逐條變成題。

### 3. 最多問四題，而且只在對應的節空著時問

| 節空著 | 問 |
| --- | --- |
| `## Background` | 這解決誰的什麼問題？現在怎麼做、哪裡痛？ |
| `## Personas` | 受影響的人與系統有誰？ |
| `## Goal` | 做成之後長怎樣？怎麼知道做成了？ |
| `## Scope — In / Out` | 什麼明確不做？ |

答案照講的記，行尾標 `← Q-<n>`，並在 `## Open Questions` 開對應的一列（狀態
`已答`、`面向` 填 `—`）。這四題是分析師開場會問的，不是十六個面向——後者是
`bdd-discovery` 的事。

MUST NOT: 問第五題。問到第五題代表你在做 discovery；停下來，把它記成紅卡。

### 4. 粗切 story，底下留空

沿原文看得出的**使用者目標**切，一個目標一則：

```markdown
## User Stories

### US-1 · visitor-billing

身為 P-1（訪客），我想在出場前知道要付多少，以便付完就能開走 ← PRD §2

**切法**：暫定。沿原文 §2–§3 的訪客計費切；規則問出來之後可能重切。

#### FR

（空——規則在 DISCOVERY 問出來才寫）

#### NFR

（空）

#### Q

（空）
```

`### US-<n> · <slug>` 的形式與 `身為 P-<n>` 的要求見 `references/prd-format.md`——
slug 就是日後的 `.feature` 檔名，第一次就要是 kebab-case。

切不出來就只開一則，slug 用 feature 名。**標「暫定」**：這一刀是沿目標切的，沿規則
切要等 `bdd-discovery` 把規則問出來，那時 `story-splitting` 會重切。

MUST NOT: 替 story 填 FR 或 AC。你看到的只有原文，原文沒有規則；寫了就是編的，
`check_spec.py` 第 14 項會把沒有 `← Q-<n>` 或 `← PRD §n` 的規則抓出來。

### 收尾

- `## Open Questions`：需求方自己列的「還不確定」每一條開一列，狀態 `待答`、`面向` `—`，
  小節寫 `**該問誰**`。原文裡的「假設」段**不**在這裡開題——那是 `bdd-discovery`
  開場掃描的工作，這裡只抄進 `## Assumptions / Constraints`
- `## Document Overview`：`**狀態**：澄清中　**最後更新**：<今天>　**結束方式**：—`；
  核心關係人表登記需求方自己與他點名的人，`決策權` 照他講的填，不知道就寫「（未問）」
- 把整份 `prd.md` 給需求方看一遍，改到他說「這是我的意思」為止。改的是開頭幾節的
  措辭，**不是 `## 原文`**

## 產物

```
specs/<date>-<feature>/
└── prd.md     CAPTURE  ★ 狀態 澄清中；story 底下三個分組空著；## 原文 逐字
```

MUST: 只寫這一個檔。不建 `spec.md`、不建 `.feature`、不寫 `docs/` 底下任何東西。

## 完成後

告訴對方四件事：

1. `prd.md` 在哪、粗切了幾則 story、切法為什麼是暫定
2. 哪幾節原文沒提、留空了
3. `## Open Questions` 開了幾張紅卡（全部來自需求方自己列的）
4. 下一步：`bdd-discovery`，它會從 `## 原文` 掃隱含假設、把每則 story 問到有規則有例子
```

- [ ] **Step 2: 跑機械稽核**

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-capture`
Expected: `bdd-capture: ✔ 通過`

- [ ] **Step 3: Commit**

```bash
git add skills/bdd-capture/SKILL.md
git commit -m "feat(bdd-capture): add the CAPTURE step that writes prd.md as a draft" -m "WHAT: New skill. Stores the originator's words verbatim in ## 原文, extracts the opening sections with ← PRD §n markers, asks at most four analyst questions, coarse-slices stories with empty FR/NFR/Q groups, and files the originator's own open questions as red cards.

WHY: The traceability chain started at a business goal no artifact owned. CAPTURE is the scribe mode (originator speaks, agent writes); DISCOVERY is the interrogation mode. Separate skills keep the modes from blending.

HOW: Format contract lives in references/prd-format.md shared with bdd-discovery and bdd-formulation.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 5: `bdd-discovery` 改寫：讀寫 `prd.md`，答案當場長成 FR／AC

**Files:**
- Modify: `skills/bdd-discovery/SKILL.md`
- Modify: `skills/bdd-discovery/references/technical-probes.md`（只改檔名引用）

**Interfaces:**
- Consumes: `../bdd-capture/references/prd-format.md`、`../bdd-capture/references/persona-definition.md`、`scripts/status.py`（Task 1）。
- Produces: `prd.md` 的 `## Open Questions`、`#### FR`／`#### NFR`／`#### Q`、`## Personas`、核心關係人表、`狀態` 由 `澄清中` 翻 `待對焦`。

- [ ] **Step 1: frontmatter 與開頭**

description 改成：

```yaml
description: >
  把 `bdd-capture` 寫好的半成品 `prd.md` 問到完整——從 `## 原文` 掃隱含假設，把每則
  story 攤成設計樹，每輪算出所有前提已定案的問題一次問完，答案回來當場長成
  `prd.md` 裡的 FR／AC（行尾標 `← Q-<n>`），紅卡清空就把狀態翻成待對焦。
  就緒判定交給 `example-mapping`，紅卡逐題收斂交給 `clarify`。
  BDD 七步流程的 DISCOVERY 步驟。
  觸發詞：「這個需求要怎麼拆」「需求講不清楚」「幫我釐清需求」
  「這份 PRD 有什麼沒寫清楚」「PRD 的隱含假設」「這份需求哪裡有洞」「PRD 怎麼拆解」
  「有哪些角色」「actor 要怎麼定義」「這份 PRD 要怎麼變成 story」
  「寫規格前先釐清」「需求有什麼漏掉的」「把 prd.md 問完整」。
  English: clarify this requirement, break down a PRD, find the unstated
  assumptions, grill me about this spec, fill in the prd.
```

標題下第一段改成：

```markdown
# DISCOVERY — 把半成品的 PRD 問到完整

一份狀態 `澄清中` 的 `prd.md` 進來（`bdd-capture` 的產物：原文逐字、開頭幾節抽好、
story 粗切、底下空著），出去的是同一份檔案，狀態 `待對焦`：**每則 story 底下有
規則有例子、每條規則與例子行尾指得出是誰說的、紅卡清空。**

**把需求攤成一棵設計樹。** 每個決定都會岔出依賴它的下一批決定。**前沿**是所有
**前提已經定案**的決定——現在就問得出口、不必先猜任何還沒聽到的答案的那些題。
一輪問完整個前沿，等答案，重算，再問下一輪。**前沿空了就結束**：每一條分支都
走過，沒有東西被默默假設。

不寫 Gherkin。不談實作。不設計資料表。**不自己填答案**：`prd.md` 裡每一條 FR 與 AC
都要指得出來源——`## 原文` 的某一段，或某一題已答的紅卡。
```

- [ ] **Step 2: 使用時機、Skill Boundaries、前置確認**

`## 使用時機` 第一條改成「`bdd-capture` 建好了 `prd.md`，要把它問到完整」，其餘保留。

`## Skill Boundaries` 改成：

```markdown
- 還沒有 `prd.md` → 先跑 `bdd-capture`，本 skill 不建檔
- 只是要把一則太大的 story 切開 → 改用 `story-splitting`，本 skill 在藍卡太多時呼叫它
- 已經有一疊跳過的紅卡要逐題收斂 → 改用 `clarify`，本 skill 在對方一次跳過多題時呼叫它
- 要把例子寫成 Gherkin → 改用 `bdd-formulation`
- 要寫 seam、API、schema、測試層級 → 改用 `bdd-spec`
- 要攤四色卡、或判定這則 story 就不就緒 → 改用 `example-mapping`——它攤地圖、跑四條
  診斷，也是唯一把狀態寫成 `已對焦` 的 skill；本 skill 只寫到 `待對焦`
- 產出是填完整的 `prd.md`，**不是** `spec.md`、程式碼或資料模型
```

`## 前置確認（Ask user）` 改成：

```markdown
## 前置確認（Ask user）

**唯一必要的輸入是 `specs/<date>-<feature>/prd.md`。** 它要有 `## 原文` 一節、
至少一則 `### US-<n> · <slug>`、狀態 `澄清中` 或 `待對焦`。

NEVER: `prd.md` 不存在時憑一段描述開始問。問出來的答案沒有地方寫、`← Q-<n>` 沒有
表可以指，整場訪談會留在對話裡然後消失。回 `bdd-capture`。

NEVER: `prd.md` 存在但沒有 `## 原文` 時開始。隱含假設的掃描對著它跑，`← PRD §x`
指著它；缺了它這一步的第一件事做不了。

產物路徑不必問，就是那份 `prd.md`。
```

- [ ] **Step 3: 流程圖與續跑**

`## 流程` 的圖改成：

```
prd.md（澄清中，bdd-capture 的產物）      或    繼續澄清（bdd-discovery [feature]）
        ↓                                          ↓
開場（只做一次）                              （開場已經做過，直接進輪次）
  掃原文找隱含假設 · 補核心關係人 · 補角色          ↓
        ↓                                          ↓
每一輪 ←──────────────────────────────────────────┘
  1. 算前沿：所有前提已定案的決定
  2. 一次問完，編號，每題附建議答案
  3. 等使用者回答
  4. 答案進 ## Open Questions，長出的 FR／AC 當場寫進它的 story → 回到 1
        ↓
前沿空了 → 收尾核對 → 狀態翻 待對焦 → example-mapping 攤地圖、跑診斷、提就緒問句
        ↓
同一份 prd.md，結束方式四選一（收斂／阻塞／零進展／喊停）寫在 Document Overview
```

`### 續跑` 的表改成：`bdd-discovery` → **預設：所有 `prd.md` 裡 `## Open Questions` 有『待答』列、狀態不是 `已對焦` 的 feature**；`bdd-discovery <feature>` 只跑那一個。

「`clarify-log.md` 不存在才跑開場……」那段改成：

```markdown
`prd.md` 的 `## Open Questions` 一列都沒有、`## Personas` 只有 CAPTURE 抽的那幾個，
就跑開場。否則只讀不重跑，並且**用它重建設計樹**：『已答』的每一列是一個已定案的
決定，『待答』是還開著的分支，『n/a』是走過但判定不適用的；每則 story 底下已有的
FR／AC 是那些決定長出來的樣子。前沿從這三種列算得出來。
```

- [ ] **Step 4: 開場改成三件事**

`### 開場：只做一次的四件事` 改標題為 `### 開場：只做一次的三件事`，刪掉 `#### 原文存檔` 整段。

`#### 拆骨架與隱含假設` 改成 `#### 掃原文找隱含假設`，第一段改成：

```markdown
`bdd-capture` 已經把業務目標、成功指標、範圍邊界從 `## 原文` 抽出來了，但它是書記，
抽的時候不質疑。這一步對 `## 原文` 與那幾節逐句質疑，問出來的題開進
`## Open Questions`，`面向` 欄填 `—`（它們不屬於十六個面向）。
```

`MUST: PRD 裡標題叫「假設」的那一段……` 改成：

```markdown
MUST: `## Assumptions / Constraints` 裡 CAPTURE 抄進來的每一條，**都當成沒問過**——
它是未驗證的斷言，不是前提。每一條開一題進『待答』；答案出來才在那一條行尾把
`← PRD §n` 換成 `← Q-<n>`，並依答案改寫成「問過的」版本。
```

掃描表保留，開頭那句「對 PRD 裡的每一個元素問同一句話」改成「對 `## 原文` 裡的每一個元素問同一句話」。

`#### 登記核心關係人` 改成：

```markdown
#### 補核心關係人

CAPTURE 登記了需求方與他點名的人。這一步每冒出一個新的受訪者就在 `prd.md`
`## Document Overview` 的 `### 核心關係人` 表補一列，三欄 `關係人／負責什麼／決策權`。
『待答』題的 `該問誰` 只能指向這張表裡的人——指到表外就補一列，不要當筆誤略過。

**`決策權` 這一格只有這裡寫得出來。** `bdd-spec` 不訪談，它把整張表原樣抄進
`spec.md`；這裡沒登記，它就分不出哪些答案是定案、哪些只是提案。
```

`#### 識別角色` 保留，第二段的路徑改成 `../bdd-capture/references/persona-definition.md`，並加一句：「CAPTURE 抽出來的角色五欄可能有『（無資料）』——這一步把它們問成有資料的，行尾 `← Q-<n>`；問不出來就留著，不要編。」

- [ ] **Step 5: 面向契約的路徑**

`### 十六個面向是前沿的檢查表` 末段「面向名稱與順序是 `../bdd-spec/scripts/status.py` 的 `DIMS`……」改成 `scripts/status.py`（同 skill）。

- [ ] **Step 6: 記錄改成「答案當場長成規則」**

`### 記錄：三件事一起做` 的表改成：

```markdown
| # | 動作 |
| --- | --- |
| 1 | 每一題答定的，`## Open Questions` 那一列狀態 `待答` → `已答`、答案欄填上結論；小節補 `**誰答的**`，技術面向的題再補 `**歸屬**`（進 `spec.md` 哪一節） |
| 2 | 答案改變了行為的，**當場**把它寫成那則 story 底下的 `**FR-<n>**`，例子寫成 `- **AC-<n>.<m>**`，每一行行尾 `← Q-<n>`；寫不進任何一則 story 的，先用 `story-splitting` 開或重切 story，再寫 |
| 3 | 那題原本在某則 story 的 `#### Q` 底下有指標的，把指標拿掉——紅卡答了就不是紅卡 |
| 4 | 若某個答案消掉了其他還沒問的問題，說出來並跳過它們 |
| 5 | 使用者對你的建議答案回「都不是」時，把那個提案連同正確答案記進該題的小節（決策史） |
```

表後面那三段（`歸屬` 欄、第 1 項不可留到最後、`為什麼還沒答`）改成：

```markdown
欄位長什麼樣、`歸屬` 的值域、FR／AC 的寫法與編號規則 →
`../bdd-capture/references/prd-format.md`。

MUST: FR 與 AC 的行尾只准 `← Q-<n>` 或 `← PRD §n`。這是這一步**唯一**的寫檔紀律，
取代了舊版「不寫規格」那條：你可以寫規則，但每一條要指得出誰說的。一條規則標
`推論` 就是你補的；`../bdd-formulation/scripts/check_spec.py` 第 14 項會抓出來，
但抓到時那條規則可能已經長出三個場景了——**寫的當下就標，不要事後補**。

第 2 項不可以留到最後一次做。剛答完那一刻你還記得這個答案長成哪條規則；隔了二十題
再回來寫，會寫成「大致上是這個意思」的版本，而那正是推論。而且澄清很少一次做完，
**停在任何一輪，`prd.md` 都該是完整可用的**。

編號：`FR-<n>` 全域唯一、`AC-<n>.<m>` 掛在 `FR-<n>` 底下、`Q-<n>` 全域唯一且**不重排**。
```

`### 跳過：要記為什麼跳` 第一段改成「留在 `## Open Questions`，狀態 `待答`、答案欄 `—`，小節記下**跳過的原因**」；表裡 `clarify-loop` 改 `clarify`；最後一段 `clarify-log.md 就會變成` 改 `prd.md 的紅卡就會變成`。

`### 攤開看形狀` 第一段「把目前已經定案的東西交給 `example-mapping`」改成「把這份 `prd.md` 交給 `example-mapping`」。

`### 停下來的訊號` 表後面「**沒有「把答案收回去」這個動作。** 三張表是單向的。『已答』的一列……`clarify-log.md` 記的是問答本身」整段改成：

```markdown
**沒有「把答案收回去」這個動作。** `## Open Questions` 的列是單向的：『已答』的一列
記的是**當時問了什麼、答了什麼**，那件事確實發生過。定錯的邊界靠**新開一題**取代，
舊那一列留著不動，它長出的那條 FR 改寫並把行尾換成新題的 `← Q-<n>`，小節寫
「取代 Q-<舊>」。`Q-<n>` 全域唯一且不重排，整份 `prd.md` 因此讀得出是哪一輪改了主意。
```

- [ ] **Step 7: 結束、為什麼保留、產物、完成後**

`## 結束：前沿空了` 的清單改成：

```markdown
1. **核對覆蓋**：跑 `python3 scripts/status.py <專案根>`。十一個業務面向每一個都至少
   出現在某一列的 `面向` 欄——包含標 `n/a` 的。技術那五個一樣各印一欄，但這一關
   不擋它們（理由見 `references/technical-probes.md` 開頭）
2. **核對 persona**：`## Personas` 的每個角色至少被一條 FR 的主詞或一則 story 的
   `身為` 用到；FR 與 story 用到的角色都在 `## Personas` 裡。只對一邊會漏掉憑空的角色
3. **核對來源**：跑 `python3 ../bdd-formulation/scripts/check_spec.py <專案根>`，
   第 14 項零筆才算前沿真的空了——有一條 FR 沒有來源，就是有一題沒問
4. **稽核分類是否窮盡**（見下）
5. **狀態翻成 `待對焦`**：`## Document Overview` 的 `狀態` 由 `澄清中` 改 `待對焦`，
   `結束方式` 寫 `收斂`。這是本 skill 唯一寫狀態的時刻；`已對焦` 不是你寫的
6. **交給 `example-mapping`**：由它攤開四色卡、跑四條診斷、把就緒問句交給使用者
```

`### 為什麼保留` 整節改成：

```markdown
### 決策史住在 `## Open Questions` 的小節

規則寫在 story 底下之後，**為什麼**留在它的 `← Q-<n>` 指向的那題小節裡：問了誰、
答案怎麼來的、被推翻過的提案。下次改動第一個會問的問題往往是「當初為什麼不選
第二個選項」，答案只在那裡找得到。

IMPORTANT: 小節不是規則的定案來源——規則以 story 底下那一行為準；小節留著的是
**為什麼**，不是規則本身。
```

`## 產物` 整節改成：

```markdown
## 產物

```
specs/<date>-<feature>/
└── prd.md     DISCOVERY  ★ 同一份檔案，從 澄清中 寫到 待對焦
```

格式 → `../bdd-capture/references/prd-format.md`。

MUST NOT: 改寫 `## 原文`。要補充就開一題；原文被改寫，每個 `← PRD §x` 從此指向
一句沒人說過的話，而且不會報錯。

MUST NOT: 寫 `spec.md` 或 `.feature`。seam、API、schema 是 `bdd-spec` 的事；技術面向
問出來的答案留在 `## Open Questions`（`歸屬` 欄標它該進 `spec.md` 哪一節），由
`bdd-spec` 去讀。

MUST NOT: 另存一份進度儀表板或 `clarify-log.md`。進度是**算出來的**——
`scripts/status.py` 從 `## Open Questions` 直接數；問答史就在那張表與它的小節。
```

`## 完成後` 改成：

```markdown
## 完成後

在 `prd.md` `## Document Overview` 的 `結束方式` 填上這一輪怎麼停的，四選一：

- **收斂** → 『待答』零列，狀態已翻 `待對焦`，下一步 `example-mapping` 提就緒問句
- **阻塞** → 卡在等外部答案，對應小節寫清楚卡在誰身上、等什麼；狀態留 `澄清中`
- **零進展** → 連續幾輪答不出新東西，回頭檢查是不是問錯層次；狀態留 `澄清中`
- **喊停** → 使用者決定帶著還開著的問題往下走；狀態留 `澄清中`，`bdd-formulation`
  與 `bdd-spec` 會拒收，要他自己決定要不要硬翻 `已對焦`

不要不講就往下跑——收斂以外的三種都代表 `prd.md` 還有洞。
```

- [ ] **Step 8: 全文殘留替換**

Run: `grep -n "clarify-log\|input.md\|clarify-loop\|bdd-spec/scripts\|bdd-spec/references/persona" skills/bdd-discovery/SKILL.md skills/bdd-discovery/references/technical-probes.md`
逐筆改：`clarify-log.md` → `prd.md`（指檔案）或 `## Open Questions`（指表）；`clarify-loop` → `clarify`；路徑改 Task 1／3 的新位置。改完再跑同一條 grep，Expected: 零筆。

- [ ] **Step 9: 稽核與 commit**

Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-discovery`
Expected: `✔ 通過`（`scripts/status.py` 已在本文被引用）。

```bash
git add skills/bdd-discovery
git commit -m "feat(bdd-discovery): read and write prd.md; answers become sourced FR/AC on the spot" -m "WHAT: Input and output are the same prd.md. Opening scans ## 原文 for unstated assumptions, stakeholders and personas are filled into prd.md, each answered question is written as FR/AC under its story with ← Q-n, and the frontier-empty close flips 狀態 to 待對焦. clarify-log.md and input.md are gone.

WHY: The red cards in Open Questions already record what was asked, what was answered and what is still open. The 'do not write' rule is replaced by 'every rule names its source', which check_spec.py enforces.

HOW: Frontier model, sixteen dimensions, round discipline and stop signals unchanged; only the recording target and the artifact section change.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 6: `clarify-loop` 改名 `clarify`

**Files:**
- Move: `skills/clarify-loop/` → `skills/clarify/`（`git mv`）
- Modify: `skills/clarify/SKILL.md`

**Interfaces:**
- Produces: skill 名 `clarify`；單獨模式寫 `specs/<topic-slug>/prd.md`（只有 `## Document Overview`、`## Open Questions`、`## 原文` 三節）。

- [ ] **Step 1: 改名**

```bash
git mv skills/clarify-loop skills/clarify
```

frontmatter `name: clarify-loop` → `name: clarify`。description 開頭保留，末尾加觸發詞「clarify」。

- [ ] **Step 2: 本文改動**

標題 `# CLARIFY LOOP — 把問題問到收斂` → `# CLARIFY — 把紅卡問到收斂`。

`## Skill Boundaries` 第三條改成：

```markdown
- 通用的需求訪談、不涉及 `prd.md` 與紅卡 → 環境中另有同名的外部 skill `clarify-loop`，
  以「最多五題」為終止條件；本 skill 以「紅卡清空」為終止條件
```

`## 前置確認` 第 1 點改成：

```markdown
1. **要收斂哪一疊問題**：`prd.md` `## Open Questions` 裡狀態為「待答」的列。
   `python3 skills/bdd-discovery/scripts/status.py` 列得出每個 feature 還剩哪幾題。
   單獨使用時，一段模糊的描述也可以——本 skill 會先把它拆成問題，逐題記錄答案
```

`### 2. 先讀已填的答案` 第一段改成：

```markdown
`prd.md` 的 `## Open Questions` 表裝的是這個 feature 問過的**每一題**，狀態三選一
`已答`／`待答`／`n/a`。本 skill 的對象是『待答』那些列；每題的小節記了跳過的原因，
先讀那個原因再決定怎麼問。
```

同節的 `MUST: 一個問題確定就立刻……` 改成：

```markdown
MUST: 一個問題確定就**立刻**在 `prd.md` `## Open Questions` 把那一列狀態改 `已答`、
答案欄填結論，小節補 `**誰答的**`。答案改變了行為的，**當場**寫成所屬 story 底下的
`**FR-<n>**`／`- **AC-<n>.<m>**`，行尾 `← Q-<n>`——寫法與編號規則見
`../bdd-capture/references/prd-format.md`。這跟 `bdd-discovery` 的記錄紀律是同一條。
```

`MUST NOT: 把答案寫進 spec.md……` 整段改成：

```markdown
MUST NOT: 寫 `spec.md`、`.feature`，或改寫 `prd.md` 的 `## 原文`。答案只有一個家：
`## Open Questions` 的列與它長出的 FR／AC。
```

`### 4. 每輪結束給進度` 範例裡「→ 成為 FR-4」「→ 成為 FR-1（其他 feature 也適用，各自 `spec.md` 加一條）」改成「→ 已寫成 FR-4 ← Q-3」「→ 已寫成 FR-1 ← Q-5（跨 feature 的規則另記進 `specs/domain-model.md`，由 `bdd-spec` 處理）」。

`## 什麼時候停` 裡「多大交給 `bdd-spec` 切 story 時決定」改成「多大由 `story-splitting` 沿規則切時決定」。

`## 產物 › 單獨使用時` 改成：

```markdown
### 單獨使用時

寫入 `specs/<topic-slug>/prd.md`，只有三節：

```markdown
# PRD: <主題>

## Document Overview

**狀態**：澄清中　**最後更新**：<YYYY-MM-DD>　**結束方式**：收斂／阻塞／零進展／喊停

## Open Questions

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q-1 | <問題> | — | 已答 | <答案> |
| Q-7 | <問題> | — | 待答 | — |

### Q-7 <問題>

**狀態**：待答　**該問誰**：教練

## 原文

<對方給的那段描述，逐字>
```

被推翻的提案記在對應題的小節。用同一個檔名與格式，是為了之後跑 `bdd-capture`
補開頭幾節時，問答不必搬家。
```

- [ ] **Step 3: 殘留檢查、稽核、commit**

Run: `grep -n "clarify-log\|spec.md\|clarify-loop" skills/clarify/SKILL.md`
Expected: 只剩 Skill Boundaries 那一句講外部 `clarify-loop` 的，以及 `MUST NOT` 裡列為禁寫對象的 `spec.md`。
Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/clarify`
Expected: `✔ 通過`

```bash
git add -A skills/clarify skills/clarify-loop
git commit -m "refactor(clarify): rename clarify-loop to clarify and point it at prd.md" -m "WHAT: git mv skills/clarify-loop to skills/clarify, rename in frontmatter, read red cards from prd.md's Open Questions table, write answers back there and as sourced FR/AC, standalone mode writes a three-section prd.md.

WHY: The loop is the red-card convergence mode bdd-discovery hands off to; clarify-log.md no longer exists. Shorter name, same job.

HOW: git mv keeps history; body edits limited to artifact names and the recording rule.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 7: `example-mapping` 改讀 `prd.md`、三種模式

**Files:**
- Modify: `skills/example-mapping/SKILL.md`

**Interfaces:**
- Produces: 唯一把 `prd.md` 狀態寫成 `已對焦` 的 skill。

- [ ] **Step 1: 改動**

description：`可讀 spec.md、bdd-discovery 進行中的訪談` → `可讀 prd.md、bdd-discovery 進行中的訪談`。

第一段 `輸入是一則 story、一份 spec.md、或` → `輸入是一則 story、一份 prd.md、或`；「四個時機」→「三種模式」。

`## 四色卡在這裡叫什麼` 表頭 `spec.md` → `prd.md`；下一段「`spec.md` 用的是後者」→「`prd.md` 用的是後者」。

`## 使用時機` 的「手上有一份 `spec.md`」→「手上有一份 `prd.md`」。

`## Skill Boundaries` 第一條 `clarify-loop` → `clarify`。

`## 攤開的地圖長這樣` 保留。

`## 四條診斷` 表的 `clarify-loop` 兩處 → `clarify`。表後那段「🟦 這條之所以照原意當阻塞：分組現在只剩一層，`bdd-spec` 寫 `spec.md` 時就是照交付切片寫」改成「分組只有一層，`prd.md` 的一個 `### US-<n> · <slug>` 直接對應一個 `.feature` 檔」。

`## 就緒判定` 的狀態表與前後兩段改成：

```markdown
MUST NOT: agent 自己把 `prd.md` `## Document Overview` 的 `狀態` 推到 `已對焦`。

**本 skill 是唯一寫 `已對焦` 的 skill**，但寫的是使用者的決定：

| 狀態 | 誰寫 | 什麼時候 |
| --- | --- | --- |
| `澄清中` | `bdd-capture` | 建檔時 |
| `待對焦` | `bdd-discovery` | 前沿空了、紅卡清空 |
| `已對焦` | 本 skill | 使用者在上面那個問句回答「就緒」 |
| `已凍結` | 尚無寫入者 | PLAN 消費過之後 |

使用者回「未就緒」或選了去問別人，狀態留 `待對焦`，紅卡開回 `## Open Questions`，
回 `bdd-discovery` 或 `clarify`。
```

`## 四種模式` 改成 `## 三種模式`：

```markdown
| 模式 | 輸入 | 輸出 |
| --- | --- | --- |
| 流程內·逼問中 | `prd.md` 狀態 `澄清中`，`bdd-discovery` 一輪記完後呼叫 | 攤開地圖看形狀；不給判讀、不提就緒問句、不寫狀態 |
| 流程內·前沿空了 | `prd.md` 狀態 `待對焦`，由 `bdd-discovery` 收尾時呼叫或使用者指定 | 地圖 ＋ 四條診斷 ＋ 就緒問句；使用者答「就緒」才寫 `已對焦` |
| 單獨 | 一段 story 的描述 | 同樣的地圖與診斷，**不寫任何檔案** |

MUST NOT: 單獨模式產出一份 `example-map.md`。**第二個裝規則與例子的地方會跟
`prd.md` 漂移**。要留下來就先跑 `bdd-capture` 開一份 `prd.md`，再寫進去。
```

`## 就緒判定` 開頭那句「這一節只在會走到就緒判定的三個時機成立……」改成「這一節只在 `流程內·前沿空了` 與 `單獨` 兩種模式成立；`流程內·逼問中` 停在前兩步」。

- [ ] **Step 2: 殘留檢查、稽核、commit**

Run: `grep -n "spec.md\|clarify-log\|clarify-loop\|四種模式\|對照 spec" skills/example-mapping/SKILL.md`
Expected: 零筆。
Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/example-mapping`
Expected: `✔ 通過`

```bash
git add skills/example-mapping/SKILL.md
git commit -m "docs(example-mapping): read prd.md and own the 已對焦 status write" -m "WHAT: Inputs are prd.md or a live bdd-discovery session; the four modes collapse to three because prd.md exists from CAPTURE onward; the status table names each writer; clarify-loop references become clarify.

WHY: The map is read off prd.md's story groups, which carry the same four-card shape spec.md's business section did.

HOW: Text-only edits; diagnostics unchanged.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 8: `bdd-spec` 瘦身成技術設計

**Files:**
- Modify: `skills/bdd-spec/SKILL.md`

**Interfaces:**
- Consumes: `prd.md`（狀態 `已對焦`）、`references/spec-format.md`（Task 3）、`examples/minimal-spec.md`、`../bdd-formulation/scripts/check_spec.py`。
- Produces: `spec.md`（技術節）、`specs/domain-model.md`、`docs/CONTEXT.md` 追加。

- [ ] **Step 1: frontmatter 與開頭**

description 改成：

```yaml
description: >
  讀已對焦的 `prd.md`，寫 `spec.md` 的技術設計——seam、模組邊界、介面與 API 契約、
  schema、每個 AC 的測試層級、風險——並把詞彙抬進 `docs/CONTEXT.md`、聚合與不變條件
  增修進 `specs/domain-model.md`。不訪談、不切 story、不寫 FR／AC；
  `.feature` 由 `bdd-formulation` 產出。BDD 七步流程 SPEC 步驟的技術那一半。
  觸發詞：「寫 spec.md」「技術設計」「決定 seam」「驗收測試打在哪一層」
  「API 契約要怎麼定」「這個 feature 動哪些模組」「prd.md 已對焦了，接下來」。
  English: write the technical spec, decide the acceptance test seam,
  implementation decisions from the prd, write the glossary.
```

標題下第一段改成：

```markdown
# SPEC — 把已對焦的 PRD 變成技術設計

一份狀態 `已對焦` 的 `prd.md` 進來，出去的是同目錄的 `spec.md`：**驗收測試打在哪一層
seam、動哪些模組、介面與 API 契約、schema、每個 AC 在哪一層測、風險。** FR／AC／story
已經在 `prd.md` 裡定案，這一步**讀**它們，不改它們。

不訪談。不切 story。不寫 FR／AC。**不產 `.feature`**——那是 `bdd-formulation` 的事。
```

- [ ] **Step 2: 使用時機、Boundaries、參考檔案**

`## 使用時機` 改成：`prd.md` 已對焦，要決定 seam 與技術設計；`.feature` 寫好了但還沒有 seam，`bdd-plan` 無法切票；要把這個 feature 的詞彙與聚合抬進跨 feature 的檔案。

`## Skill Boundaries` 改成：

```markdown
- **要把 FR 與 AC 寫成 `.feature` → 改用 `bdd-formulation`**
- `prd.md` 還有紅卡、或狀態不是 `已對焦` → 回 `bdd-discovery`／`example-mapping`
- story 切得不對 → `story-splitting`，然後回 `bdd-discovery`；本 skill 不碰 `## User Stories`
- 要決定實作順序 → 改用 `bdd-plan`
- 要稽核既有 `.feature` 寫得好不好 → `bdd-spec-review` 未實作（見 `PLAN.md`），先跑 `bdd-formulation` 的 `scripts/check_spec.py`
- 產出是 `specs/<date>-<feature>/spec.md` ＋ `specs/domain-model.md` ＋ `docs/CONTEXT.md` 的追加，**不是** `prd.md`、`.feature`、`openapi.yaml`、migration
```

`## 參考檔案` 改成：

```markdown
- `references/spec-format.md` — `spec.md` 的骨架（只有技術節）
- `../bdd-capture/references/prd-format.md` — 讀 `prd.md` 時對照它的節
- `examples/minimal-spec.md` — 一份完整的 `spec.md` 範例
- `../bdd-formulation/scripts/check_spec.py` — 交件前自檢（第 5 步）。**在隔壁 skill**：本 skill 只跑它，不改它
- `../bdd-discovery/scripts/status.py` — 確認 `prd.md` 沒有待答的題
```

- [ ] **Step 3: 刪切 story，改詞彙表**

刪掉整節 `## 切 story`，換成：

```markdown
## story 在 `prd.md` 裡，不在這裡切

`## User Stories` 由 CAPTURE 粗切、DISCOVERY 沿規則細切，到這一步已經定案。
一則 story ＝ 一個 `.feature` 檔 ＝ 一次可獨立驗收的交付，slug 就是檔名。

MUST NOT: 改 `prd.md` 的 `## User Stories`。寫技術設計時覺得切法不對，寫進
`spec.md` 的 `## 回 DISCOVERY 補問`，由 DISCOVERY 用 `story-splitting` 重切。
悄悄重切比公開改糟，因為 `.feature` 的檔名與 tag 已經掛在舊的切法上。
```

`## 詞彙表` 第一段「把 `clarify-log.md`「已答」表裡**已經有答案**的詞抬進」改成「把 `prd.md` `## Open Questions` 裡狀態 `已答`、答案是一個定義的詞抬進」。

- [ ] **Step 4: 前置確認與流程**

`## 前置確認` 改成：

```markdown
## 前置確認（Ask user）

不必問輸入的位置。`prd.md` 在 `specs/<date>-<feature>/`。唯一要徵詢的是 seam，在流程步驟 3。

MUST NOT: 在這裡問需求。**不訪談**——`prd.md` 推不出來的，寫進 `spec.md` 的
`## 回 DISCOVERY 補問`，不要順手決定掉。

NEVER: `prd.md` 不存在、或狀態不是 `已對焦` 時開始。`澄清中` 代表還有紅卡，
`待對焦` 代表需求方還沒點頭——兩種都可能改 FR，而 seam 與 schema 是照 FR 定的。
回 `bdd-discovery` 或 `example-mapping`。

NEVER: `prd.md` 沒有 `## 原文` 時開始。`← PRD §x` 指的就是它。
```

`### 1. 挑 story` 改成 `### 1. 挑 feature`：

```markdown
| 呼叫方式 | 範圍 |
| --- | --- |
| `bdd-spec` | **預設：狀態 `已對焦` 的每一份 `prd.md`** |
| `bdd-spec <feature>` | 只做指定的那一個 |

MUST: 狀態不是 `已對焦` 的 `prd.md` 整份跳過，並明講跳過的理由與該回哪個 skill。
```

`### 2. 讀兩樣東西` 改成 `### 2. 讀 prd.md 的這幾節`：

```markdown
| 讀什麼 | 為了什麼 |
| --- | --- |
| `## User Stories` | 每條 FR 與 AC——seam 要看過具體的 AC 才推得出來；`## Testing Decisions` 要替每個 AC 標測試層級 |
| `## Open Questions` 裡 `面向` 是技術面向、狀態 `已答` 的列 | **技術決定的原始素材**——小節的 `歸屬` 說它該進 `spec.md` 哪一節，行尾 `← Q-<n>` |
| `## Assumptions / Constraints`、`## Non-Functional Requirements` | 外部限制與品質屬性，`## Tech Preferences & Constraints` 與 `## Risks` 的來源 |
| `## Document Overview` 的核心關係人 | **照抄**進 `spec.md` 的同名表 |
| `## 原文` | `← PRD §x` 的依據 |

IMPORTANT: `prd.md` 由 CAPTURE 與 DISCOVERY 維護，本 skill **不得寫入**。推導不出來的
寫進 `spec.md` 的 `## 回 DISCOVERY 補問`。
```

`### 3. 決定 seam` 保留四條規則表，`Ask user` 段改成：

```markdown
**Ask user: 把 seam 攤出來確認再寫。** 這是本 skill 唯一一次徵詢。DISCOVERY 的技術
面向通常已經問過 seam（`## Open Questions` 裡 `面向` 為 `seam` 的題）——有的話這裡
只是確認，行尾 `← Q-<n>`；沒有的話，提出你的判斷、標 `← 推論`，並在
`## 回 DISCOVERY 補問` 記一條「seam 沒問過」。

MUST NOT: 把 seam 的討論變成需求訪談。談的時候發現需求有洞，寫進
`## 回 DISCOVERY 補問`，**不要順手問掉**——這一步沒有把答案寫進 `prd.md` 的權限，
順口問到的答案只活在這一次對話裡。
```

`### 4. 寫 spec.md` 改成：

```markdown
### 4. 寫 `spec.md`

路徑 `specs/<date>-<feature>/spec.md`，骨架照抄 `references/spec-format.md`——只有技術節。

MUST: `## Document Overview` 的核心關係人表**照抄** `prd.md`，不增列也不改寫 `決策權`。

MUST: `## Known Boundaries` 的三張清單要填，空的也要寫為什麼空。

MUST: `## 回 DISCOVERY 補問` 要有——沒有洞就寫「（無）」。有洞就一條一列，
寫清楚影響什麼定不下來。

MUST NOT: 在這一步做需求決定。**本 skill 只綜合已經有答案的東西。** 判準：指著
`spec.md` 的任何一句話，說得出它來自 `prd.md` 的哪條 FR、哪題已答的 Q 嗎？
```

`### 5. 增修 domain-model.md` 保留。`### 6. 自檢` 改成 `### 5. 自檢`：

```markdown
```bash
python3 <ai-bdd>/skills/bdd-formulation/scripts/check_spec.py <專案根>
```

它讀的是 `prd.md`。`.feature` 還不存在時只跑得到十項——全部是 `prd.md` 自己的形狀，
包含第 14 項「每條 FR／AC 有來源」。它在這裡跑是因為**你剛讀完 `prd.md`**：抓到的
每一條都是 DISCOVERY 漏的，寫進 `## 回 DISCOVERY 補問`，不要自己修 `prd.md`。
```

`## 產物格式` 改成：

```markdown
```
specs/<date>-<feature>/
└── spec.md           SPEC  ★ 技術設計

（prd.md 是 CAPTURE／DISCOVERY 的產物，本 skill 只讀不寫）
（.feature 是 bdd-formulation 的產物，本 skill 不產）
```

| 產物 | 格式 |
| --- | --- |
| `specs/<date>-<feature>/spec.md` | `references/spec-format.md` |
| `specs/domain-model.md`（增修） | `references/spec-format.md` 末節 |
| `docs/CONTEXT.md`（建立或追加） | 「詞彙表」一節 |
```

刪掉「規則與例子的定版編號在這裡誕生」那一段（編號現在在 DISCOVERY 誕生）。

`## 完成後` 改成四件事：seam 是哪一層；`domain-model.md` 與 `docs/CONTEXT.md` 增加了什麼；`## 回 DISCOVERY 補問` 有幾條；`check_spec.py` 十項跑出什麼。結尾：「補問清單非空 → 回 `bdd-discovery`。空 → 跑 `bdd-formulation`（若還沒跑）或 `bdd-plan`」。

- [ ] **Step 5: 殘留檢查、稽核、commit**

Run: `grep -n "clarify-log\|input.md\|切 story\|clarify-loop\|scripts/status.py\|persona-definition" skills/bdd-spec/SKILL.md`
Expected: 只剩 Boundaries 裡 `story-splitting` 那句、`../bdd-discovery/scripts/status.py` 的新路徑。
Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-spec`
Expected: `✔ 通過`

```bash
git add skills/bdd-spec/SKILL.md
git commit -m "refactor(bdd-spec): write only the technical spec from an 已對焦 prd.md" -m "WHAT: Input is prd.md at 已對焦; story slicing and FR/AC numbering leave this skill; output spec.md keeps seam, implementation and testing decisions, risks and a 回 DISCOVERY 補問 section; stakeholders copied from prd.md.

WHY: Business content now has its own longer-lived file written by CAPTURE and DISCOVERY. SPEC stays the no-interview writing mode.

HOW: Section-level rewrite; seam rules and domain-model handling unchanged.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 9: `bdd-formulation` 與其參考檔改讀 `prd.md`

**Files:**
- Modify: `skills/bdd-formulation/SKILL.md`
- Modify: `skills/bdd-formulation/references/state-tags.md`
- Modify: `skills/bdd-formulation/references/rule-taxonomy.md`
- Modify: `skills/bdd-formulation/references/coverage-report.md`
- Modify: `skills/bdd-formulation/references/step-grammar.md`
- Modify: `skills/bdd-formulation/examples/anti-patterns.md`

**Interfaces:**
- Consumes: `prd.md`（`## User Stories`、`## Personas`、`## Scope`、狀態）。

- [ ] **Step 1: SKILL.md 替換規則**

先列：`grep -n "spec.md\|clarify-log\|bdd-spec" skills/bdd-formulation/SKILL.md`。逐行依**意思**改：

- 指 FR／AC／story／Personas／Scope／狀態／`## User Stories`／`## Open Questions` 的 `spec.md` → `prd.md`（這是絕大多數：description、標題、第 1 步讀的表、第 2 步對應表、第 6 步稽核說明、產物格式、完成後）。
- 「`spec.md` 由 `bdd-spec` 維護，本 skill 不得寫入」→「`prd.md` 由 CAPTURE 與 DISCOVERY 維護，本 skill 不得寫入」。
- Skill Boundaries：「`spec.md` 還不存在、或還有沒答完的問題 → 回 `bdd-spec`（更上游是 `bdd-discovery`）」→「`prd.md` 不存在 → `bdd-capture`；還有紅卡或狀態不是 `已對焦` → `bdd-discovery`／`example-mapping`」；「要切 story、決定 seam、寫 `spec.md` 或 `docs/CONTEXT.md` → 那是 `bdd-spec` 的工作」→「要決定 seam、寫 `spec.md` 或 `docs/CONTEXT.md` → `bdd-spec`；要切 story → `story-splitting`」。
- 前置確認的 NEVER：「`spec.md` 不存在時憑需求描述直接寫 `.feature`……沒有 `spec.md` 就先跑 `bdd-spec`」→「`prd.md` 不存在或狀態不是 `已對焦` 時……先跑 `bdd-capture` → `bdd-discovery`」；「`## Scope — In / Out` 裡「回 CLARIFY 補問」還有東西時」→「`## Open Questions` 還有待答時」。
- 第 1 步 MUST「只做 `spec.md` 的 `## Document Overview` 裡 `狀態` 已經是 `待對焦` 或……」→「只做 `prd.md` 狀態 `已對焦` 的」。
- 「規則與它自己的例子打架時」一節：「MUST: 同時回報給 `bdd-spec`，由它寫進 `spec.md` 的 `## Scope — In / Out` → 回 CLARIFY」→「MUST: 回報給使用者並建議回 `bdd-discovery` 開一題；本 skill 與 `bdd-spec` 都不得改 `prd.md`」；「MUST NOT: 直接改原始輸入，也不得寫入 `spec.md` 或 `clarify-log.md`」→「MUST NOT: 寫入 `prd.md`，包括它的 `## 原文`」。
- 第 6 步稽核說明裡「十三項」→「十五項」，在列舉末尾加「每條 FR／AC 行尾有沒有來源、`## 原文` 在不在」；「`spec.md` 裡有沒有 FR 完全沒掛任何 AC（沒有例子，SPEC 就無從寫出場景）」→「（沒有例子，FORMULATION 就無從寫出場景）」。
- 產物格式註解「（spec.md 是 bdd-spec 的產物，本 skill 只讀不寫）」→「（prd.md 是 CAPTURE／DISCOVERY 的產物，本 skill 只讀不寫；spec.md 是 bdd-spec 的，本 skill 不讀也不寫）」。
- 完成後第 4 點「`spec.md` 有沒有需要更正的地方」→「`prd.md` 有沒有需要回 DISCOVERY 補的地方」。
- 「這份 `.feature` 是寫給誰讀的」一節與 Feature 敘述那段：「角色抄 `spec.md` 的 `## Personas`」→「角色抄 `prd.md` 的 `## Personas`」。

- [ ] **Step 2: 參考檔**

`state-tags.md` 表：`@draft｜bdd-spec｜來源 clarify-log.md 未收斂` → `@draft｜bdd-formulation｜來源 prd.md 狀態 待對焦（需求方還沒點頭）`；`@ready｜bdd-spec｜來源 clarify-log.md 已收斂` → `@ready｜bdd-formulation｜來源 prd.md 狀態 已對焦`。「六步流程」→「七步流程」。

`rule-taxonomy.md` 第 3、6、102 行 `spec.md` → `prd.md`；第 121 行「規則來自 `clarify-log.md` 已答的問題」→「規則來自 `prd.md` 裡標了 `← Q-<n>` 或 `← PRD §n` 的那幾條」。

`coverage-report.md` 第 28 行、`step-grammar.md` 第 222 行、`anti-patterns.md` 第 91 行：`spec.md` → `prd.md`。

- [ ] **Step 3: 殘留檢查、稽核、commit**

Run: `grep -n "clarify-log\|input.md\|回 CLARIFY\|六步" skills/bdd-formulation/`
Expected: 零筆。
Run: `grep -n "spec.md" skills/bdd-formulation/SKILL.md`
Expected: 只剩 Boundaries 與產物格式註解裡指 `bdd-spec` 擁有的 `spec.md` 那幾處。
Run: `python3 skills/skill-rules/scripts/audit_skill.py skills/bdd-formulation`
Expected: `✔ 通過`

```bash
git add skills/bdd-formulation
git commit -m "docs(bdd-formulation): translate prd.md, not spec.md, into .feature" -m "WHAT: Every reference to the FR/AC/story/persona source becomes prd.md; the readiness gate is prd.md at 已對焦; conflicts between a rule and its examples are reported back to DISCOVERY instead of written into spec.md; state-tags map @draft/@ready to prd.md statuses.

WHY: The business half of spec.md moved to prd.md with its format intact, so the translation rules are unchanged and only the file name and gate move.

HOW: Meaning-based substitution; spec.md stays where it means bdd-spec's technical file.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 10: `bdd-plan` 與 `story-splitting` 的檔名引用

**Files:**
- Modify: `skills/bdd-plan/SKILL.md`
- Modify: `skills/story-splitting/SKILL.md`

- [ ] **Step 1: `bdd-plan`**

`grep -n "spec.md\|clarify" skills/bdd-plan/SKILL.md` 列出的行，依意思改：

- description 與第 15 行「一份就緒的 `spec.md` 與它涵蓋的 `.feature`」→「一份 `prd.md`（FR 與 story）、它的 `spec.md`（seam 與技術決定）與 `.feature`」。
- 第 18 行「API、domain 型別、schema、seam 都已經在 `spec.md` 裡定案」保留。
- 第 23 行使用時機「`spec.md` 與 `.feature` 都寫好了」→「`prd.md`、`spec.md` 與 `.feature` 都寫好了」。
- 第 28 行「還沒有 `spec.md` → 先跑 `bdd-spec`；有 `spec.md` 但還沒有 `.feature` → 先跑 `bdd-formulation`」→「還沒有 `prd.md` → `bdd-capture`；`prd.md` 不是 `已對焦` → `bdd-discovery`；還沒有 `.feature` → `bdd-formulation`；還沒有 `spec.md` → `bdd-spec`」。
- 第 41 行判準「來自 `spec.md` 的哪一節」→「來自 `prd.md` 的哪條 FR、或 `spec.md` 的哪一節」。
- 第 44 行「`spec.md` 沒寫到而切票時發現需要的，寫進 `## 缺口` 回報」保留，補一句「業務面的缺口請對方回 `bdd-discovery`，技術面的回 `bdd-spec`」。
- 第 52 行讀什麼「`spec.md` 全部，加上它涵蓋的每個 `.feature`」→「`prd.md` 的 `## User Stories`、`spec.md` 全部、每個 `.feature`」。
- 第 64 行「依 `spec.md` 的 FR 切」→「依 `prd.md` 的 FR 切」。
- 第 95 行「`spec.md` 裡推不出來的缺口」→「`prd.md`／`spec.md` 裡推不出來的缺口」。

- [ ] **Step 2: `story-splitting`**

- 第 26 行 `clarify-loop` → `clarify`。
- 第 37–38 行「或 `spec.md` 的某幾條 `FR-<n>`」「有 `spec.md` 就直接用它的 `FR-<n>`」→ `prd.md`。
- 第 65 行「`bdd-spec` 寫 `spec.md` 時，一個……」→「`prd.md` 裡一個 `### US-<n> · <slug>` 直接對應一個 `.feature` 檔，一個……」。
- 第 187 行「規則編號沿用 `spec.md` 原本的」→「沿用 `prd.md` 原本的」。
- 第 197–198 行「每則切完後回 `bdd-spec` 寫進 `spec.md`，再由 `bdd-formulation`……——規則與例子在 `spec.md` 已經就緒，不是回 `bdd-discovery` 補問」→「每則切完後直接改 `prd.md` 的 `## User Stories`（搬 FR 到它的新 story 底下，編號不變），再由 `bdd-formulation` 把各自的 FR 與例子寫成 `.feature`——規則與例子已經就緒，這一刀不需要回 `bdd-discovery` 補問」。

- [ ] **Step 3: 稽核、commit**

Run: `grep -n "clarify-loop\|clarify-log" skills/bdd-plan/SKILL.md skills/story-splitting/SKILL.md`
Expected: 零筆。
Run: `python3 skills/skill-rules/scripts/audit_skill.py skills`
Expected: 每個 skill `✔ 通過`（此時 `bdd-capture` 在內，`clarify-loop` 已不存在）。

```bash
git add skills/bdd-plan/SKILL.md skills/story-splitting/SKILL.md
git commit -m "docs(bdd-plan,story-splitting): read FR and stories from prd.md" -m "WHAT: bdd-plan reads prd.md for FR/stories and spec.md for seam; story-splitting edits prd.md's User Stories directly and points at clarify instead of clarify-loop.

WHY: Follow-through of the prd.md/spec.md split.

HOW: Reference-only edits.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 11: `PLAN.md`、`docs/ai-sdlc.md`、`CLAUDE.md`、`benchmark/README.md`

**Files:**
- Modify: `PLAN.md`
- Modify: `docs/ai-sdlc.md`
- Modify: `CLAUDE.md`
- Modify: `benchmark/README.md`

- [ ] **Step 1: `PLAN.md`**

`## 四條紀律` 整節改成 `## 三種模式，四條紀律`：

```markdown
## 三種模式，四條紀律

這是 2026-10-02 重排的核心。七步之所以站得住，是因為前四步底下的每個 skill
**各自只有一種工作模式**——SPEC 由兩個 skill 分擔，所以是四條：

| 步驟 | 模式 | 只做 | 不做 |
| --- | --- | --- | --- |
| **CAPTURE**（`bdd-capture`） | 需求方說、agent 記 | 原文逐字、抽開頭幾節、最多問四題、粗切 story | 不填規則與例子 |
| **DISCOVERY**（`bdd-discovery`） | agent 問、需求方答 | 問紅卡；答過的才能變 FR／AC，行尾標來源 | 不寫 Gherkin、不設計、不自己填答案 |
| **SPEC**（`bdd-formulation`／`bdd-spec`） | agent 寫、需求方審 | `已對焦` 的 `prd.md` → `.feature` ＋ `spec.md` | 不訪談，有洞記進 `## 回 DISCOVERY 補問` |
| **PLAN**（`bdd-plan`） | agent 切、工程師審 | tracer bullet ＋ blocking edges | 不設計 |

「agent 補一個合理的答案然後繼續」這個風險，舊版靠「DISCOVERY 不寫檔」擋，
而 2026-09-14 的設計承認那條從來沒真的成立過。現在改靠「每條 FR／AC 必須指得出
來源」擋：`← PRD §n` 或 `← Q-<n>`，`check_spec.py` 第 14 項機械地抓。擋的是同一件事，
但從「不准做」改成「做了要留證據」。

紀律的出處：Matt Pocock 的 `to-spec`——「The spec is a record of decisions
already made, not a place where new ones get made.」
```

`## 六步流程` 改成 `## 七步流程`，表改成：

```markdown
| 步驟 | 這一步在問什麼 | 作法 | 產物 |
| --- | --- | --- | --- |
| **CAPTURE** | 需求方想要什麼、為什麼、有什麼限制？ | 當書記：原文逐字、抽開頭幾節、粗切 story | `prd.md`（澄清中） |
| **DISCOVERY** | 還有哪些不知道？問到不能再問為止 | 設計樹＋前沿，每則 story 攤成 FR／AC／Q | 同一份 `prd.md`（待對焦 → 已對焦） |
| **SPEC** | 這些答案的可執行規格與技術設計長怎樣？ | `prd.md` → Gherkin（`bdd-formulation`）；`prd.md` → seam、契約、schema（`bdd-spec`） | `.feature`（可執行）＋ `spec.md` ＋ `domain-model.md` |
| **PLAN** | 怎麼切成一次做得完的步驟？ | tracer bullet 垂直切片 ＋ blocking edges | `plan.md` |
| **IMPLEMENT** | 讓這張票從 red 到 green | step defs ＋ outside-in 雙迴圈 | 程式碼 |
| **VERIFY** | 綠了嗎？紅的是 bug 還是規格過期？ | 跑情境 ＋ 失敗判讀 | 測試結果 ＋ 判讀結論 |
| **REVIEW** | 這變更真的滿足它宣稱的情境嗎？鏈斷了嗎？ | 對照情境驗收 ＋ 鏈稽核 | 審查結論 ＋ 稽核報告 |
```

`## 產物佈局` 整段改成：

```markdown
## 產物佈局

一步一個產物槽，`ls` 就看得出這一批走到哪：

```
specs/
├── domain-model.md     SPEC     聚合、不變條件；跨 feature，活得比任何一份 spec 久
└── <date>-<feature>/
    ├── prd.md          CAPTURE 建、DISCOVERY 填到 已對焦。story、FR、AC、角色、範圍、
    │                   開放問題、原文。活得比實作久
    ├── spec.md         SPEC     seam、模組、契約、schema、測試層級、風險。實作開始就過期
    └── plan.md         PLAN     tracer bullet ＋ blocking edges

features/<story-slug>.feature    SPEC  可執行規格，一則 story 一個檔
docs/CONTEXT.md                  SPEC  詞彙，跨 feature
```

**`prd.md` 與 `spec.md` 分檔的理由是壽命。** FR、AC、角色在實作完成後仍然是這個
feature 的事實；seam 與 schema 在第一個 commit 之後就開始過期。`to-spec` 說 spec 是
拋棄式快照——所以快照裡只放會過期的東西。

`.feature` 是唯一不進 `specs/` 的產物——它同時是規格與可執行測試，
**位置由測試框架決定**。鏈的完整性靠 tag 與 slug 維繫，不是靠同目錄。
```

`### spec.md 的節` 改成 `### prd.md 與 spec.md 的節`，兩張表：`prd.md`：Document Overview、Background、Goal、Scope、Personas、User Stories（`#### FR`／`#### NFR`／`#### Q`）、NFR、Assumptions/Constraints、Success Metrics、Open Questions、Further Notes、原文；`spec.md`：Document Overview、Tech Preferences & Constraints、Known Boundaries、Implementation Decisions、Testing Decisions、Risks、回 DISCOVERY 補問、Further Notes。保留「一批一份」與「MUST NOT 寫具體檔案路徑」兩段；刪掉「規則與例子不在 `spec.md` 裡……`clarify.md`」那段（它描述的是更舊的佈局），換成一句「規則與例子住 `prd.md`，`check_spec.py` 比對的是 `prd.md` 的 FR／AC 與 `.feature` 的 tag」。

`## 步驟 → skills 對照` 表改成：

```markdown
| 步驟 | skill | 職責 | 狀態 |
| --- | --- | --- | --- |
| CAPTURE | `bdd-capture` | 當書記，`prd.md` 建成半成品 | **已實作**（2026-10-02） |
| DISCOVERY | `bdd-discovery` | 設計樹＋前沿，把 `prd.md` 問到完整 | **已實作**；2026-10-02 改成讀寫 `prd.md` |
| | `clarify` | 一疊紅卡逐題收斂 | **已實作**；原 `clarify-loop`，由 `bdd-discovery` 呼叫 |
| | `story-splitting` | 沿規則切；九種模式 | **已實作**；由 `bdd-capture`（粗切）與 `bdd-discovery`（細切）呼叫 |
| | `example-mapping` | 四色卡、四條診斷、就緒判定；唯一寫 `已對焦` 的 skill | **已實作**；由 `bdd-discovery` 呼叫 |
| SPEC | `bdd-formulation` | `prd.md` → `.feature` | **已實作** |
| | `bdd-spec` | `prd.md` → `spec.md`（技術設計） | **已實作**；2026-10-02 瘦身 |
| | `bdd-spec-review` | 反命令式、conjunction step、情境爆炸稽核 | 未實作 |
| PLAN | `bdd-plan` | `.feature` ＋ `spec.md` → tracer bullet ＋ blocking edges | **已實作** |
| IMPLEMENT | `bdd-implement`、`bdd-implement-step-definitions`、`bdd-implement-outside-in` | | 未實作 |
| VERIFY | `bdd-verify`、`bdd-verify-failure-triage` | | 未實作 |
| REVIEW | `bdd-review`、`bdd-review-acceptance`、`bdd-review-traceability` | | 未實作 |
```

表後「`bdd-spec` 加了 `spec.md` 之後會不會太肥……」段保留。`### Example Mapping 的兩個訊號在哪裡看` 表裡 `clarify.md` → `prd.md`，`status.py` 路徑改 `bdd-discovery`。`### 追問覆蓋改成算出來的` 裡「每個問題檔的狀態列」→「`## Open Questions` 每一列的 `面向` 欄」，範例那行狀態列刪掉，`n/a` 段改成「`n/a` 的理由寫在答案欄」。

`## 可追溯鏈` 的圖與說明改成：

```
需求方的話（prd.md ## 原文）→ Story → FR → AC → Scenario → Ticket → Step → 程式碼
                               ↑         ↑                  ↑
                          CAPTURE 粗切  DISCOVERY 定版       一票一 seam，在 SPEC
                                       行尾 ← Q-n / PRD §n
```

兩次定版點的歷史後面加第三條：「**2026-10-02**：留在 DISCOVERY，但換了載體——編號在 `prd.md` 誕生，`clarify-log.md` 與 `input.md` 消失，`spec.md` 只剩技術節。來源標記取代『不寫檔』成為守門規則。」

`## 進度` 加三項：`- [x] 2026-10-02 重排：CAPTURE 步驟、prd.md、來源標記守門`、`- [ ] 實測場重跑一次，驗證 CAPTURE → DISCOVERY → SPEC 的新分工`、`- [ ] bdd-implement`（原本那條保留）。

第 26 行「環境中既有的 `clarify-loop`」保留——那句講的是外部 skill。

- [ ] **Step 2: `docs/ai-sdlc.md`**

`## 六步，以及哪幾段是 BDD` 改成 `## 七步，以及哪幾段是 BDD`，圖改成：

```
CAPTURE  →  DISCOVERY  →  SPEC  →  PLAN  →  IMPLEMENT  →  VERIFY  →  REVIEW
   ✗        └──── BDD ──────┘       ✗       └─ BDD ─┘        ✗          ✗
```

表加一列在最前：

```markdown
| **CAPTURE** | 把需求方的話變成一份他簽過名的半成品 `prd.md`：原文逐字、問題／預期結果／角色／範圍、粗切 story。這是 playbook 的 `intent.md`——鏈上第一個有人擁有的東西 | **無**——BDD 假設 story 已經存在 |
```

CLARIFY 那列改名 DISCOVERY，內容改成「兩件事：對 `## 原文` 掃隱含假設（不是 BDD）＋把每則 story 逼成規則與例子（Discovery）。**答案當場寫成 `prd.md` 的 FR／AC，每一條行尾指得出是誰說的。**」；SPEC 那列改成「兩件事：把已對焦的 `prd.md` 翻成可執行規格（Formulation，`bdd-formulation`）＋把技術決定寫成 `spec.md`（不是 BDD，`bdd-spec`）。切 story 不再在這裡」。

「六塊不屬於 BDD 的東西」那段改成「七塊……（CAPTURE 整步、PLAN、REVIEW、DISCOVERY 的隱含假設掃描、SPEC 的技術設計、story 的粗切、綜合成文）」。

`## 四條紀律` 整節換成 `PLAN.md` Step 1 那張 `三種模式，四條紀律` 表與其後兩段（兩份文件各自有一份，這是既有做法：`PLAN.md` 講建什麼，這裡講為什麼）。

`### 為什麼是四個 skill，不是一個` 改成 `### 為什麼是五個 skill，不是一個`，開頭加一段：

```markdown
CAPTURE 與 DISCOVERY 相鄰，也是**方向相反的工作模式**：CAPTURE 是需求方說、模型記；
DISCOVERY 是模型問、需求方答。合成一個 skill，模型會在記的時候順手補規則——那正是
2026-09-09 到 09-11 三個版本都犯的事，每一次都靠措辭掩飾。拆開之後 CAPTURE 連
FR 都不准寫，補不了。
```

其餘保留。`## agent 能替代什麼` 表的 Discovery 列「四色卡怎麼對應 `spec.md`」→ `prd.md`。

`## AI 特有的四條` 第 3 條的四個例子保留；在第 2 條 Fail loud 後加一句「2026-10-02 起多一條機械防線：`check_spec.py` 第 14 項拒絕沒有來源的 FR／AC」。

- [ ] **Step 3: `CLAUDE.md`**

「What this repo is」段：`six-step pipeline` → `seven-step pipeline`，`**CLARIFY → SPEC → PLAN → IMPLEMENT → VERIFY → REVIEW**` → `**CAPTURE → DISCOVERY → SPEC → PLAN → IMPLEMENT → VERIFY → REVIEW**`，「Only the first three steps exist as skills」→「Only the first four steps exist as skills」。

「Skill scripts」段：`status.py belongs to bdd-spec` → `status.py belongs to bdd-discovery`，路徑 `skills/bdd-discovery/scripts/status.py`，註解改 `# red cards and probe coverage per feature, read from prd.md`；`check_spec.py` 註解改 `# prd.md shape incl. source markers; adds .feature ↔ prd.md coverage once features exist`。加一行 `python3 -m unittest discover -s tests   # script tests`。

「The four disciplines」段改成「The three modes, four disciplines」：

```markdown
### The three modes, four disciplines

Each of the first four steps holds because every skill under it runs in exactly one
working mode. SPEC has two skills, so there are four disciplines:

- **CAPTURE** is scribe mode: the originator speaks, the agent writes `prd.md` as a
  draft. It asks at most four analyst questions and **never fills in rules or examples**.
- **DISCOVERY** is interrogation mode: the agent asks, the originator answers, and each
  answer is written on the spot as an FR/AC under its story with a `← Q-<n>` source.
  **An FR or AC without `← PRD §n` or `← Q-<n>` is a defect** — `check_spec.py` check 14
  rejects it. This replaces the older "CLARIFY does not write" rule, which could not be
  verified.
- **SPEC** is writing mode: `bdd-formulation` turns an `已對焦` `prd.md` into `.feature`
  files; `bdd-spec` writes the technical `spec.md` — **neither interviews**, and holes go
  into `spec.md`'s `## 回 DISCOVERY 補問`.
- **PLAN** slices tracer bullets — it does not design.
```

「The document contract」表保留。「Artifact format is the interface」的鏈改成 `需求方的話 → Story → FR → AC → Scenario → Ticket → Step → 程式碼`，前一句「SPEC reads answered questions」→「DISCOVERY reads CAPTURE's draft, SPEC reads the `已對焦` `prd.md`」。

「Conventions」第二條：「Documentation artifacts the skills produce go to the consuming repo's `specs/<date>-<feature>/` (e.g. `spec.md`)」→「(`prd.md`, `spec.md`, `plan.md`)」；後面「SPEC may create it or append entries」保留。

「The testbed」段：「`internal/domain/` currently holds only `doc.go`」保留；「no business code until CLARIFY has run」→「until DISCOVERY has run」。

- [ ] **Step 4: `benchmark/README.md`**

第 60–62 行改成：

```
├── specs/<date>-<feature>/prd.md           CAPTURE 建、DISCOVERY 填到 已對焦
├── specs/<date>-<feature>/spec.md          SPEC（bdd-spec）的技術設計
└── features/                               SPEC（bdd-formulation）產出的 .feature
```

「然後把該 case 的輸入檔**整份**交給 `bdd-discovery`」→「先交給 `bdd-capture` 建 `prd.md`（它會逐字收進 `## 原文`），再跑 `bdd-discovery`」。其餘 `grep -n "CLARIFY\|clarify-log\|input.md" benchmark/README.md` 列出的行依意思改成 DISCOVERY／`prd.md`。

- [ ] **Step 5: 殘留檢查、commit**

Run: `grep -rn "clarify-log\|input\.md\|六步\|six-step\|brief\.md\|actor\.md" PLAN.md docs/ai-sdlc.md CLAUDE.md benchmark/README.md`
Expected: 只剩 `PLAN.md` 可追溯鏈歷史那三條（它們記的是過去）與 `## 進度` 已勾選的舊項目。

```bash
git add PLAN.md docs/ai-sdlc.md CLAUDE.md benchmark/README.md
git commit -m "docs: describe the seven-step pipeline with CAPTURE and prd.md" -m "WHAT: PLAN.md, docs/ai-sdlc.md, CLAUDE.md and benchmark/README.md now describe CAPTURE → DISCOVERY → SPEC → PLAN → IMPLEMENT → VERIFY → REVIEW, the prd.md/spec.md artifact split, the three modes and four disciplines, the source-marker guard, and the renamed clarify skill. The stale brief.md/actor.md layout is gone.

WHY: The documents are the roadmap and the position; they must match what the skills do.

HOW: Section rewrites; history entries kept as history.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>"
```

---

### Task 12: 全部閘門

**Files:** 無新檔。

- [ ] **Step 1: 跑四個閘門**

```bash
python3 -m unittest discover -s tests -v
python3 skills/skill-rules/scripts/audit_skill.py skills
claude plugin validate .
grep -rn "clarify-loop" skills/ docs/ai-sdlc.md PLAN.md CLAUDE.md benchmark/README.md
grep -rn "clarify-log\|input\.md" skills/ docs/ai-sdlc.md CLAUDE.md benchmark/README.md
```

Expected，依序：`Ran 9 tests ... OK`；每個 skill `✔ 通過`，含 `bdd-capture` 與 `clarify`，不含 `clarify-loop`；`Validation passed with warnings` 且恰好一個警告（根層 `CLAUDE.md`）；`clarify-loop` 只剩 `PLAN.md` 第 26 行與 `skills/clarify/SKILL.md` Boundaries 講外部 skill 的那句；`clarify-log`／`input.md` 零筆。

- [ ] **Step 2: 腳本對實測場跑**

```bash
python3 skills/bdd-discovery/scripts/status.py benchmark/skeleton/go
python3 skills/bdd-formulation/scripts/check_spec.py benchmark/skeleton/go
```

Expected: 兩者 exit 1，訊息指名找不到 `specs/`（實測場是空的，fail-loud 是正確行為）。

- [ ] **Step 3: 把 `tasks/todo.md` 不存在這件事確認一次**

本計畫沒有建 `tasks/todo.md`（計畫本身就是進度表）；`ls tasks/ 2>/dev/null` 應為空。沒有要 commit 的東西就結束。
