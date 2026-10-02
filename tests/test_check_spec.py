import contextlib
import io
import tempfile
import unittest
from pathlib import Path

from _load import ROOT, load

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

    def test_parenthetical_annotation_after_marker_is_still_sourced(self):
        text = (GOOD.replace("← PRD §3", "← PRD §3（限制：問過）")
                    .replace("← Q-1", "← Q-1（30 分整算免費）"))
        self.assertEqual(check_spec.unsourced_rules(text), [])

    def test_prose_after_marker_is_unsourced(self):
        text = GOOD.replace("← PRD §3", "← PRD §3 — 另見 Q-2")
        self.assertEqual(len(check_spec.unsourced_rules(text)), 1)


class CitedParagraphs(unittest.TestCase):
    def test_all_cited_paragraphs_exist(self):
        self.assertEqual(check_spec.missing_cited_paragraphs(GOOD), [])

    def test_citation_to_absent_paragraph_is_reported(self):
        text = GOOD.replace("← PRD §3", "← PRD §9")
        self.assertEqual(check_spec.missing_cited_paragraphs(text), ["9"])


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

    def test_capture_fresh_prd_points_at_discovery(self):
        fresh = GOOD.replace(
            "**FR-1**  When 訪客車出場, the system shall 依停留時長計費。 ← PRD §3\n"
            "- **AC-1.1**  停 29 分 → 收 0 元 ← Q-1\n", "（空）\n")
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_prd(root, fresh)

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("bdd-discovery", out)
            self.assertNotIn("先跑 bdd-capture", out)

    def test_citation_to_absent_paragraph_fails_check(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_prd(root, GOOD.replace("← PRD §3", "← PRD §9"))

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("§9", out)

    def test_shipped_example_prd_passes_prd_only_checks(self):
        example = (ROOT / "skills" / "bdd-capture" / "examples" / "minimal-prd.md").read_text(encoding="utf-8")
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_prd(root, example)

            code, out = run(root)

            self.assertEqual(code, 0, out)

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
