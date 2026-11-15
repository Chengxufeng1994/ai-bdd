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


if __name__ == "__main__":
    unittest.main()
