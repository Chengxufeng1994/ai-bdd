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
