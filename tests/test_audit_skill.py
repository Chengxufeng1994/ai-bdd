import tempfile
import unittest
from pathlib import Path

from _load import load

audit_skill = load("skills/skill-rules/scripts/audit_skill.py")

HEAD = """---
name: demo
description: >
  示範用的 skill，用來測試稽核腳本能不能抓到缺少的章節。
  觸發詞：「示範」「demo」。
  English: demo skill.
---

# demo

## 使用時機

- 測試

## Skill Boundaries

- 無
"""


def rules_hit(body: str) -> set[str]:
    with tempfile.TemporaryDirectory() as tmp:
        d = Path(tmp) / "demo"
        d.mkdir()
        (d / "SKILL.md").write_text(HEAD + body, encoding="utf-8")
        return {f.rule for f in audit_skill.audit(d) if f.level == audit_skill.MUST}


def rules_hit_any(text: str) -> set[str]:
    with tempfile.TemporaryDirectory() as tmp:
        d = Path(tmp) / "demo"
        d.mkdir()
        (d / "SKILL.md").write_text(text, encoding="utf-8")
        return {f.rule for f in audit_skill.audit(d)}


CARD = """
| | |
| --- | --- |
| **職責** | 示範 |
| **讀** | 無 |
| **寫** | 無 |
| **上游** | 無 |
| **下游** | 無 |
| **不做** | 不寫檔 |
"""


class ResponsibilityCard(unittest.TestCase):
    # S16: 讀的人要在同一個位置找到「這個 skill 讀什麼、寫什麼、交給誰」，
    # 不必讀完流程才拼得出來。

    def with_card(self, card: str) -> str:
        return HEAD.replace("# demo\n", "# demo\n" + card)

    def test_no_card_is_flagged(self):
        self.assertIn("S16", rules_hit_any(HEAD))

    def test_complete_card_passes(self):
        self.assertNotIn("S16", rules_hit_any(self.with_card(CARD)))

    def test_card_missing_a_row_is_flagged(self):
        # 少「不做」一列：邊界是卡片最不能省的一格
        card = CARD.replace("| **不做** | 不寫檔 |\n", "")
        self.assertIn("S16", rules_hit_any(self.with_card(card)))

    def test_card_below_the_first_section_does_not_count(self):
        # 卡片要在標題底下、第一個章節之前——放到後面就不是「同一個位置」了
        self.assertIn("S16", rules_hit_any(HEAD + CARD))


class Verification(unittest.TestCase):
    # S15: 沒有交回前的檢查清單，模型會在「看起來寫完了」時交回。

    def test_missing_section_is_a_must(self):
        self.assertIn("S15", rules_hit("\n## 流程\n\n1. 做事\n"))

    def test_heading_without_checkboxes_is_a_must(self):
        # 只有標題、底下是散文，等於沒有清單——項目才是能逐一答是或否的東西
        self.assertIn("S15", rules_hit("\n## Verification\n\n記得檢查。\n"))

    def test_checkbox_in_another_section_does_not_count(self):
        body = "\n## Verification\n\n記得檢查。\n\n## 完成後\n\n- [ ] 告訴對方\n"
        self.assertIn("S15", rules_hit(body))

    def test_section_with_checkboxes_passes(self):
        body = "\n## Verification\n\n- [ ] 產物存在\n\n## 完成後\n\n告訴對方。\n"
        self.assertNotIn("S15", rules_hit(body))


def skill(write: str, flow: str, products: str = "") -> str:
    card = CARD.replace("| **寫** | 無 |", f"| **寫** | {write} |")
    return HEAD.replace("# demo\n", "# demo\n" + card) + flow + products


TABLE = "\n## 產物\n\n| 檔 | 哪一步寫 | 格式在哪 |\n| --- | --- | --- |\n"


class StepOutputs(unittest.TestCase):
    # S17: 每一步交出什麼要寫在步驟開頭，讀的人不必讀完整步才知道它產出什麼。

    def test_numbered_step_without_output_is_flagged(self):
        self.assertIn("S17", rules_hit_any(skill("無", "\n## 流程\n\n### 1. 讀\n\n- x\n")))

    def test_step_with_output_passes(self):
        self.assertNotIn("S17", rules_hit_any(skill("無", "\n## 流程\n\n### 1. 讀\n\n產出：無，只讀\n")))

    def test_step_n_heading_counts_as_numbered(self):
        # plan 的步驟叫「Step 2 依賴圖」，別的 skill 指它時用這個名字，不改成「2.」
        self.assertIn("S17", rules_hit_any(skill("無", "\n## 流程\n\n### Step 2 依賴圖\n\n- x\n")))

    def test_unnumbered_subsection_is_exempt(self):
        # 「續跑」「不越界」這類不是步驟，沒有自己的產出
        self.assertNotIn("S17", rules_hit_any(skill("無", "\n## 流程\n\n### 續跑\n\n- x\n")))

    def test_no_flow_section_skips_both(self):
        # 照搬原文的 skill 沒有 `## 流程`，不受這兩條約束
        rules = rules_hit_any(skill("`a.md`", "\n## Overview\n\n### 1. Do\n\n- x\n"))
        self.assertFalse({"S17", "S18"} & rules)


class ProductTable(unittest.TestCase):
    # S18: 產物表與職責卡的「寫」是同一份清單的兩個位置，對不上就有一邊是錯的。

    FLOW = "\n## 流程\n\n### 1. 寫\n\n產出：`a.md`\n"

    def test_product_row_missing_from_card_is_flagged(self):
        table = TABLE + "| `a.md` | 1 | x |\n| `b.md` | 1 | x |\n"
        self.assertIn("S18", rules_hit_any(skill("`a.md`", self.FLOW, table)))

    def test_card_path_missing_from_table_is_flagged(self):
        table = TABLE + "| `a.md` | 1 | x |\n"
        self.assertIn("S18", rules_hit_any(skill("`a.md`、`b.md`", self.FLOW, table)))

    def test_missing_table_is_flagged(self):
        self.assertIn("S18", rules_hit_any(skill("`a.md`", self.FLOW)))

    def test_span_with_space_does_not_shift_backtick_pairs(self):
        # `#### 缺口` 不是檔；它後面的 `b.md` 仍要被認出來
        table = TABLE + "| `a.md` | 1 | x |\n| `b.md` | 1 | x |\n"
        self.assertNotIn("S18", rules_hit_any(skill("`a.md`（`#### 缺口`）、`b.md`", self.FLOW, table)))

    def test_only_the_product_table_is_read(self):
        # 產物節底下別的表（產物隔離的規則表）提到的路徑不是產物
        table = TABLE + "| `a.md` | 1 | x |\n\n### 產物隔離\n\n| 規則 | 強度 |\n| --- | --- |\n| 只寫入 `specs/` 底下 | MUST |\n"
        self.assertNotIn("S18", rules_hit_any(skill("`a.md`", self.FLOW, table)))

    def test_matching_card_and_table_pass(self):
        table = TABLE + "| `a.md` | 1 | x |\n"
        self.assertNotIn("S18", rules_hit_any(skill("`a.md`", self.FLOW, table)))


if __name__ == "__main__":
    unittest.main()
