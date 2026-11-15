import contextlib
import io
import tempfile
import unittest
from pathlib import Path

from _load import load

status = load("skills/discovery/scripts/status.py")

DISCOVERY = """# discovery: parking

## Document Overview

**狀態**：澄清中　**最後更新**：2026-10-02

## Open Questions

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q-1 | 免費額度多長 | 邊界 | 已答 | 30 分 |
| Q-2 | 跨午夜怎麼算 | 時間 | 待答 | — |
| Q-3 | 機車費率 | — | n/a | 本期不做機車 |

## 原文

§1 停車場要收費。
"""

PRD = DISCOVERY  # old-layout fixtures only need a file to exist



def run(root: Path) -> tuple[int, str]:
    buf = io.StringIO()
    with contextlib.redirect_stdout(buf):
        code = status.main(root)
    return code, buf.getvalue()


class StatusReadsDiscovery(unittest.TestCase):
    def test_counts_rows_from_discovery(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "2026-10-02-parking" / "discovery.md").write_text(DISCOVERY, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 0)
            self.assertRegex(out, r"2026-10-02-parking\s+1\s+1\s+1\s")  # 已答 / n/a / 待答
            self.assertIn("待答（1）", out)
            self.assertIn("2026-10-02-parking/Q-2", out)

    def test_missing_discovery_exits_1_and_names_discovery(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "2026-10-02-parking" / "plan.md").write_text("# plan", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("沒有任何 discovery.md", out)
            self.assertIn("discovery", out)
            self.assertNotIn("bdd-capture", out)

    def test_old_layout_prd_beside_spec_exits_1_and_names_git_mv(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "2026-10-02-parking" / "prd.md").write_text(PRD, encoding="utf-8")
            (root / "specs" / "2026-10-02-parking" / "spec.md").write_text("# old technical", encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("舊版佈局", out)
            self.assertIn("先 `git mv spec.md design.md`，再 `git mv prd.md spec.md`", out)

    def test_prd_only_dir_names_single_git_mv(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "2026-10-02-parking" / "prd.md").write_text(PRD, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("舊版佈局", out)
            self.assertIn("`git mv prd.md spec.md`", out)
            self.assertNotIn("design.md", out)

    def test_1113_layout_exits_1_with_split_hint(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "specs" / "2026-10-02-parking").mkdir(parents=True)
            (root / "specs" / "2026-10-02-parking" / "spec.md").write_text(DISCOVERY, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("11-13 佈局", out)
            self.assertIn("discovery.md", out)


if __name__ == "__main__":
    unittest.main()
