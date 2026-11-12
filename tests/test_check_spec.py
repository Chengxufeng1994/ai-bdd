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


def write_spec(root: Path, text: str) -> None:
    d = root / "specs" / "2026-10-02-parking"
    d.mkdir(parents=True)
    (d / "spec.md").write_text(text, encoding="utf-8")


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


class CheckReadsSpec(unittest.TestCase):
    def test_good_spec_passes_in_spec_only_mode(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertIn("全部通過", out)

    def test_unsourced_fr_fails(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD.replace("← PRD §3", "← 推論"))

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("沒有來源", out)
            self.assertIn("FR-1", out)

    def test_missing_verbatim_section_reported(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD.split("## 原文")[0])

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("## 原文", out)

    def test_fresh_spec_points_at_discovery(self):
        fresh = GOOD.replace(
            "**FR-1**  When 訪客車出場, the system shall 依停留時長計費。 ← PRD §3\n"
            "- **AC-1.1**  停 29 分 → 收 0 元 ← Q-1\n", "（空）\n")
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, fresh)

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("bdd-discovery", out)
            self.assertNotIn("先跑 bdd-capture", out)

    def test_citation_to_absent_paragraph_fails_check(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD.replace("← PRD §3", "← PRD §9"))

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("§9", out)

    def test_shipped_example_spec_passes_spec_only_checks(self):
        example = (ROOT / "skills" / "bdd-discovery" / "examples" / "minimal-spec.md").read_text(encoding="utf-8")
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, example)

            code, out = run(root)

            self.assertEqual(code, 0, out)

    def test_spec_only_dir_exits_1(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "x").mkdir(parents=True)
            (root / "specs" / "x" / "design.md").write_text("# old", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("沒有 spec.md", out)
            self.assertIn("bdd-discovery", out)
            self.assertNotIn("bdd-capture", out)

    def test_old_layout_prd_beside_spec_exits_1_and_names_git_mv(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)
            (root / "specs" / "2026-10-02-parking" / "prd.md").write_text(GOOD, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("舊版佈局", out)
            self.assertIn("先 `git mv spec.md design.md`，再 `git mv prd.md spec.md`", out)

    def test_prd_only_dir_names_single_git_mv(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "2026-10-02-parking" / "prd.md").write_text(GOOD, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("舊版佈局", out)
            self.assertIn("`git mv prd.md spec.md`", out)
            self.assertNotIn("design.md", out)


if __name__ == "__main__":
    unittest.main()


class StoriesNotYetFormulated(unittest.TestCase):
    """A features/ dir that holds only unrelated files (the walking skeleton's
    version.feature) means FORMULATION has not run for these stories yet —
    that is a notice, not a coverage failure. Missing examples are only a
    problem once the story's own .feature exists."""

    def test_missing_feature_file_is_a_notice_not_a_problem(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)
            (root / "features").mkdir()
            (root / "features" / "version.feature").write_text("Feature: v\n", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertIn("還沒有 .feature", out)
            self.assertIn("visitor-billing", out)

    def test_existing_feature_file_missing_examples_is_still_a_problem(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)
            (root / "features").mkdir()
            (root / "features" / "visitor-billing.feature").write_text(
                "@ready @query\nFeature: x\n", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("漏了", out)
