import contextlib
import io
import tempfile
import unittest
from pathlib import Path

from _load import ROOT, load

check_spec = load("skills/formulation/scripts/check_spec.py")

GOOD = """# spec: parking

## Document Overview

**狀態**：待對焦

## User Stories

### US-1 · visitor-billing

身為 P-1（訪客），我想在出場前知道要付多少 ← 推論

#### FR

**FR-1**  When 訪客車出場, the system shall 依停留時長計費。 ← PRD §3
- **AC-1.1**  停 29 分 → 收 0 元 ← Q-1

#### Q

**Q-2**  跨午夜怎麼算 ← 擋住 FR-1
"""

DISCOVERY = """# discovery: parking

## Document Overview

**狀態**：已收斂

## Open Questions

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q-1 | 免費額度多長 | 邊界 | 已答 | 30 分 |
| Q-2 | 跨午夜怎麼算 | 時間 | 待答 | — |

## 原文

§3 前 30 分鐘免費。
"""

COMBINED_1113 = GOOD + "\n## Open Questions" + DISCOVERY.split("## Open Questions", 1)[1]



def run(root: Path) -> tuple[int, str]:
    buf = io.StringIO()
    with contextlib.redirect_stdout(buf):
        code = check_spec.check(root)
    return code, buf.getvalue()


def write_spec(root: Path, text: str, discovery: str | None = DISCOVERY) -> None:
    d = root / "specs" / "2026-10-02-parking"
    d.mkdir(parents=True)
    (d / "spec.md").write_text(text, encoding="utf-8")
    if discovery is not None:
        (d / "discovery.md").write_text(discovery, encoding="utf-8")


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
        self.assertEqual(check_spec.missing_cited_paragraphs(GOOD + DISCOVERY), [])

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

    def test_technical_part_after_stories_is_not_read_as_rules(self):
        # spec.md 的技術部分接在需求部分之後；它的決定標 `← 推論`／`← 專案 <路徑>`，
        # 不是 FR／AC，不該被第 14 項當成沒有來源的規則，也不該被當成 story
        tech = (
            "\n## Implementation Decisions\n\n"
            "- **FR-9 不是規則**：只是一行長得像 FR 的技術決定 ← 推論\n"
            "- **動哪些模組**：新增計費模組 ← 推論\n\n"
            "## Commands\n\n- 閘門：`make verify` ← 專案 Makefile\n"
        )
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD + tech)

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertNotIn("FR-9", out)

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
            write_spec(root, GOOD, DISCOVERY.split("## 原文")[0])

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("## 原文", out)

    def test_spec_without_rules_points_at_bdd_spec(self):
        fresh = GOOD.replace(
            "**FR-1**  When 訪客車出場, the system shall 依停留時長計費。 ← PRD §3\n"
            "- **AC-1.1**  停 29 分 → 收 0 元 ← Q-1\n", "（空）\n")
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, fresh)

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("spec", out)
            self.assertNotIn("先跑 bdd-capture", out)

    def test_citation_to_absent_paragraph_fails_check(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD.replace("← PRD §3", "← PRD §9"))

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("§9", out)

    def test_shipped_example_spec_passes_spec_only_checks(self):
        example = (ROOT / "skills" / "spec" / "examples" / "minimal-spec.md").read_text(encoding="utf-8")
        disc = (ROOT / "skills" / "discovery" / "examples" / "minimal-discovery.md").read_text(encoding="utf-8")
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, example, disc)

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
            self.assertIn("discovery", out)
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

    def test_q_marker_missing_from_discovery_fails(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD.replace("收 0 元 ← Q-1", "收 0 元 ← Q-9"))

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("Q-9", out)
            self.assertIn("discovery.md", out)

    def test_q_marker_pointing_at_unanswered_fails(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD.replace("收 0 元 ← Q-1", "收 0 元 ← Q-2"))

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("Q-2", out)
            self.assertIn("已答", out)

    def test_discovery_only_dir_is_a_notice(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "x").mkdir(parents=True)
            (root / "specs" / "x" / "discovery.md").write_text(DISCOVERY, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertIn("先跑 spec", out)

    def test_1113_layout_exits_1_with_split_hint(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, COMBINED_1113, discovery=None)

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("11-13 佈局", out)
            self.assertIn("discovery.md", out)


JOURNEY = """Feature: 訪客出場

  Rule: 前 30 分鐘免費

    @covers-1.1
    Example: 訪客停一下就走
      Given 訪客停了 29 分鐘
      When 訪客開到出口
      Then 不用付錢
"""


def write_journey(root: Path, name: str, text: str) -> None:
    d = root / "specs" / "2026-10-02-parking" / "acceptance"
    d.mkdir(parents=True, exist_ok=True)
    (d / name).write_text(text, encoding="utf-8")


class AcceptanceJourneys(unittest.TestCase):
    """Journeys are not executed, so @covers- is the only thing tying them to
    spec.md. A journey citing an AC that does not exist is exactly the drift
    the unexecuted layer invites: an outcome nobody agreed to."""

    def test_journey_citing_known_ac_passes(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)
            write_journey(root, "訪客出場.feature", JOURNEY)

            code, out = run(root)

            self.assertEqual(code, 0, out)

    def test_journey_citing_unknown_ac_fails(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)
            write_journey(root, "訪客出場.feature", JOURNEY.replace("@covers-1.1", "@covers-9.9"))

            code, out = run(root)

            self.assertEqual(code, 1, out)
            self.assertIn("@covers-9.9", out)
            self.assertIn("訪客出場.feature", out)

    def test_example_without_covers_fails(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)
            write_journey(root, "訪客出場.feature", JOURNEY.replace("    @covers-1.1\n", ""))

            code, out = run(root)

            self.assertEqual(code, 1, out)
            self.assertIn("訪客停一下就走", out)

    def test_no_acceptance_dir_is_silent(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertNotIn("acceptance", out)

    def test_feature_level_covers_do_not_anchor_examples(self):
        text = "@covers-1.1\nFeature: x\n  Example: y\n    Then z\n"
        self.assertEqual(check_spec.journey_covers(text), ({"1.1"}, ["y"]))

    def test_multiple_covers_on_one_line(self):
        text = "Feature: x\n  @covers-1.1 @covers-1.2\n  Example: y\n"
        self.assertEqual(check_spec.journey_covers(text), ({"1.1", "1.2"}, []))

    def test_malformed_covers_reported_as_unknown(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)
            write_journey(root, "訪客出場.feature", JOURNEY.replace("@covers-1.1", "@covers-1"))

            code, out = run(root)

            self.assertEqual(code, 1, out)
            self.assertIn("@covers-1", out)

    def test_journey_without_recognisable_example_fails(self):
        # zh-TW keywords parse as nothing to journey_covers: no covers, no bare
        # Examples — check 17 would otherwise pass a journey it never read.
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)
            write_journey(root, "訪客出場.feature", "# language: zh-TW\n功能: x\n  場景: 沒標\n")

            code, out = run(root)

            self.assertEqual(code, 1, out)
            self.assertIn("訪客出場.feature", out)

    def test_scenario_template_counts_as_example(self):
        self.assertEqual(check_spec.journey_covers("Feature: x\n  Scenario Template: y\n"),
                         (set(), ["y"]))

    def test_shipped_journey_example_passes(self):
        skills = ROOT / "skills"
        spec = (skills / "spec" / "examples" / "minimal-spec.md").read_text(encoding="utf-8")
        disc = (skills / "discovery" / "examples" / "minimal-discovery.md").read_text(encoding="utf-8")
        journey = (skills / "spec-by-example" / "examples" / "訪客出場結帳.feature").read_text(encoding="utf-8")
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, spec, disc)
            write_journey(root, "訪客出場結帳.feature", journey)

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertNotIn("acceptance", out)


class FindFeatures(unittest.TestCase):
    """Journeys are written before any features/ dir exists. The rglob
    fallback must not mistake specs/<dir>/acceptance/ for the runner's
    feature directory, or every journey is reported as an orphan."""

    def test_acceptance_only_returns_none(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_journey(root, "a.feature", JOURNEY)

            self.assertIsNone(check_spec.find_features(root))

    def test_rglob_still_finds_features_outside_specs(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_journey(root, "b.feature", JOURNEY)
            (root / "app" / "features").mkdir(parents=True)
            (root / "app" / "features" / "a.feature").write_text("Feature: a\n", encoding="utf-8")

            self.assertEqual(check_spec.find_features(root), root / "app" / "features")

    def test_acceptance_not_listed_as_orphan_when_features_dir_exists(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write_spec(root, GOOD)
            (root / "features").mkdir()
            (root / "features" / "version.feature").write_text("Feature: v\n", encoding="utf-8")
            write_journey(root, "訪客出場.feature", JOURNEY)

            code, out = run(root)

            self.assertEqual(code, 0, out)
            orphan_lines = [l for l in out.splitlines() if "沒有對應 story" in l]
            self.assertTrue(orphan_lines)
            self.assertNotIn("訪客出場", orphan_lines[0])


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
