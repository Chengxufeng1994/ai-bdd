import contextlib
import io
import tempfile
import unittest
from pathlib import Path

from _load import load

status = load("skills/discovery/scripts/status.py")

# 2026-11-18 起問答表住 spec.md；discovery 建檔時它還只有粗切 story
SPEC = """# spec: parking

## Document Overview

**狀態**：澄清中　**最後更新**：2026-10-02　**結束方式**：—

## User Stories

### US-1 · visitor-billing

身為 P-1（訪客），我想在出場前知道要付多少 ← PRD §1

## Open Questions

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q-1 | 免費額度多長 | 邊界 | 已答 | 30 分 |
| Q-2 | 跨午夜怎麼算 | 時間 | 待答 | — |
| Q-3 | 機車費率 | — | n/a | 本期不做機車 |

## 原文

§1 停車場要收費。
"""

PRD = SPEC  # old-layout fixtures only need a file to exist



def run(root: Path) -> tuple[int, str]:
    buf = io.StringIO()
    with contextlib.redirect_stdout(buf):
        code = status.main(root)
    return code, buf.getvalue()


class StatusReadsSpec(unittest.TestCase):
    def test_counts_rows_from_spec(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "ongoing" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "ongoing" / "2026-10-02-parking" / "spec.md").write_text(SPEC, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 0)
            self.assertRegex(out, r"2026-10-02-parking\s+1\s+1\s+1\s")  # 已答 / n/a / 待答
            self.assertIn("待答（1）", out)
            self.assertIn("2026-10-02-parking/Q-2", out)

    def test_root_domain_model_exits_1_and_names_git_mv(self):
        # 2026-11-15 起聚合的現況住 docs/DOMAIN-MODEL.md，specs/ 根層不再放它
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "ongoing" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "ongoing" / "2026-10-02-parking" / "spec.md").write_text(SPEC, encoding="utf-8")
            (root / "specs" / "domain-model.md").write_text("# Domain Model", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("git mv specs/domain-model.md docs/DOMAIN-MODEL.md", out)

    def test_flat_layout_exits_1_and_names_git_mv(self):
        # 2026-11-15 以前 feature 直接放在 specs/ 底下；新佈局在 specs/ongoing/
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "2026-10-02-parking" / "spec.md").write_text(SPEC, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("git mv", out)
            self.assertIn("specs/ongoing/2026-10-02-parking", out)

    def test_all_features_archived_exits_0(self):
        # 全部 feature 都歸檔、specs/ongoing/ 被 git 清掉是正常狀態
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "sleeping" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "sleeping" / "2026-10-02-parking" / "spec.md").write_text(SPEC, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertIn("沒有進行中的 feature", out)

    def test_missing_spec_exits_1_and_names_discovery(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "ongoing" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "ongoing" / "2026-10-02-parking" / "plan.md").write_text("# plan", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("沒有任何 spec.md", out)
            self.assertIn("discovery", out)
            self.assertNotIn("bdd-capture", out)

    def test_old_layout_prd_beside_spec_exits_1_and_names_git_mv(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "ongoing" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "ongoing" / "2026-10-02-parking" / "prd.md").write_text(PRD, encoding="utf-8")
            (root / "specs" / "ongoing" / "2026-10-02-parking" / "spec.md").write_text("# old technical", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("舊版佈局", out)
            self.assertIn("先 `git mv spec.md design.md`，再 `git mv prd.md spec.md`", out)

    def test_prd_only_dir_names_single_git_mv(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "ongoing" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "ongoing" / "2026-10-02-parking" / "prd.md").write_text(PRD, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("舊版佈局", out)
            self.assertIn("`git mv prd.md spec.md`", out)
            self.assertNotIn("design.md", out)

    def test_legacy_two_file_layout_exits_1_with_merge_hint(self):
        # 11-18 以前問答住 discovery.md；不擋下來，spec.md 會被報成「找不到 `## Open Questions`」
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            feature = root / "specs" / "ongoing" / "2026-10-02-parking"
            feature.mkdir(parents=True)
            (feature / "discovery.md").write_text(SPEC, encoding="utf-8")
            (feature / "spec.md").write_text("# spec: parking\n", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("discovery.md", out)
            self.assertIn("spec.md", out)
            self.assertIn("git rm", out)

    def test_single_file_spec_is_not_legacy(self):
        # 11-13 的退件訊息把「spec.md 有 `## Open Questions`、旁邊沒有 discovery.md」當舊佈局；11-18 起它就是新佈局
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "ongoing" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "ongoing" / "2026-10-02-parking" / "spec.md").write_text(SPEC, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertNotIn("佈局", out)


if __name__ == "__main__":
    unittest.main()
