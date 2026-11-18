import contextlib
import io
import tempfile
import unittest
from pathlib import Path

from _load import load

status = load("skills/discovery/scripts/status.py")

# 2026-11-18 起 discovery 只寫 interview-log.md；問答一題一節，面向與狀態是兩個粗體欄
LOG = """# interview-log：parking

## 原文

§1 停車場要收費。

## Session

### 2026-10-02 · 建檔

- **需求方**：場長

## 問答

### Q-1 · 免費額度多長

**面向**：邊界　**狀態**：已答
- 2026-10-02 場長 →「30 分」
**誰答的**：場長

### Q-2 · 跨午夜怎麼算

**面向**：時間　**狀態**：待答
**該問誰**：場長

### Q-3 · 機車費率

**面向**：—　**狀態**：n/a
**收斂結果**：本期不做機車

### Q-4 · 月租車要不要另算

**面向**：重複　**狀態**：失效
**收斂結果**：被 Q-3 取代

## Confirmed Decisions

（無）

## Remaining Open Items

**結束方式**：—
"""

# 11-18 單檔佈局（discovery 寫進 spec.md）與舊佈局的 fixture 只要檔案在就好
SPEC = """# spec: parking

## Open Questions

| Q | 問題 | 面向 | 狀態 | 答案 |
| --- | --- | --- | --- | --- |
| Q-1 | 免費額度多長 | 邊界 | 已答 | 30 分 |
"""

PRD = SPEC  # old-layout fixtures only need a file to exist



def run(root: Path) -> tuple[int, str]:
    buf = io.StringIO()
    with contextlib.redirect_stdout(buf):
        code = status.main(root)
    return code, buf.getvalue()


FEATURE = ("specs", "ongoing", "2026-10-02-parking")


def write(root: Path, name: str, text: str) -> None:
    d = root.joinpath(*FEATURE)
    d.mkdir(parents=True, exist_ok=True)
    (d / name).write_text(text, encoding="utf-8")


class StatusReadsInterviewLog(unittest.TestCase):
    def test_counts_questions_from_interview_log(self):
        # 失效跟 n/a 一樣是結論，不是還要去追誰要答案
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "interview-log.md", LOG)

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertRegex(out, r"2026-10-02-parking\s+1\s+2\s+1\s")  # 已答 / n/a / 待答
            self.assertIn("待答（1）", out)
            self.assertIn("2026-10-02-parking/Q-2", out)

    def test_coverage_comes_from_dimension_field(self):
        # 邊界、重複、時間各有一題；— 不算任何面向
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "interview-log.md", LOG)

            code, out = run(root)

            self.assertIn("—✓✓———✓—————————", out)

    def test_spec_beside_log_is_the_normal_layout(self):
        # spec 讀 log 寫出 spec.md；兩份並存是常態，不是舊佈局
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "interview-log.md", LOG)
            write(root, "spec.md", SPEC)

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertNotIn("佈局", out)

    def test_counts_ignore_spec_table(self):
        # 進度只數 log：spec.md 的表是 spec 整理過的副本，不是第二個來源
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "interview-log.md", LOG)
            write(root, "spec.md", SPEC.replace("| 已答 |", "| 待答 |"))

            code, out = run(root)

            self.assertIn("待答（1）", out)

    def test_missing_qa_section_exits_1(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "interview-log.md", LOG.replace("## 問答", "## 問題"))

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("找不到 `## 問答`", out)

    def test_empty_qa_section_exits_1(self):
        # 一節存在但零題：DISCOVERY 沒問過，不是需求沒有疑點
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            text = LOG[:LOG.index("### Q-1")] + LOG[LOG.index("## Confirmed Decisions"):]
            write(root, "interview-log.md", text)

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("零題", out)

    def test_question_without_status_is_reported_and_others_still_counted(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "interview-log.md", LOG.replace("**面向**：時間　**狀態**：待答", "**面向**：時間"))

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("Q-2", out)
            self.assertRegex(out, r"2026-10-02-parking\s+1\s+2\s+0\s")

    def test_spec_without_log_is_single_file_layout(self):
        # 11-18 的單檔佈局：discovery 把問答寫進 spec.md；不擋下來，這個 feature 會安靜地不見
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "spec.md", SPEC)

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("interview-log.md", out)
            self.assertIn("spec.md", out)

    def test_root_domain_model_exits_1_and_names_git_mv(self):
        # 2026-11-15 起聚合的現況住 docs/DOMAIN-MODEL.md，specs/ 根層不再放它
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "interview-log.md", LOG)
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
            (root / "specs" / "sleeping" / "2026-10-02-parking" / "interview-log.md").write_text(LOG, encoding="utf-8")

            code, out = run(root)

            self.assertEqual(code, 0, out)
            self.assertIn("沒有進行中的 feature", out)

    def test_missing_log_exits_1_and_names_discovery(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "plan.md", "# plan")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("沒有任何 interview-log.md", out)
            self.assertIn("discovery", out)
            self.assertNotIn("bdd-capture", out)

    def test_old_layout_prd_beside_spec_exits_1_and_names_git_mv(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "prd.md", PRD)
            write(root, "spec.md", "# old technical")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("舊版佈局", out)
            self.assertIn("先 `git mv spec.md design.md`，再 `git mv prd.md spec.md`", out)

    def test_prd_only_dir_names_single_git_mv(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "prd.md", PRD)

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("舊版佈局", out)
            self.assertIn("`git mv prd.md spec.md`", out)
            self.assertNotIn("design.md", out)

    def test_legacy_two_file_layout_exits_1_with_rename_hint(self):
        # 11-18 以前問答住 discovery.md；它現在的位置是 interview-log.md
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            write(root, "discovery.md", SPEC)
            write(root, "spec.md", "# spec: parking\n")

            code, out = run(root)

            self.assertEqual(code, 1)
            self.assertIn("git mv discovery.md interview-log.md", out)
            self.assertIn("spec.md", out)


if __name__ == "__main__":
    unittest.main()
