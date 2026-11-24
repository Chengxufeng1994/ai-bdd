#!/usr/bin/env python3
"""機械性稽核一個 skill 目錄。

只檢查能客觀判定的項目。需要判斷力的部分（觸發詞涵蓋不涵蓋得夠、步驟具不具體、
理由充不充分）交給人與模型，見 ../rules/quality-checklist.md。

已知限制:
  * S5（孤兒檔案）只認 markdown 文字裡出現的路徑。被 Python `import` 或被
    shell 呼叫的輔助模組不會出現在文字中，會被誤判為孤兒。這是 MUST 級檢查
    帶著已知偽陽性——遇到時人工判斷，不要為了讓腳本閉嘴而加無意義的引用。
  * F8（全形冒號）只掃行首的 `key：`，縮排的續行不查。它擋得住最常見的手誤，
    擋不住所有情況。
  * 不檢查名稱衝突。名稱只需在同一個 plugin 內唯一（S12），而腳本看不到
    其他 plugin。

用法:
    python3 audit_skill.py <skill 目錄>
    python3 audit_skill.py <plugin 的 skills 目錄>   # 逐一稽核底下每個 skill

離開碼: 0 = 沒有 MUST 等級的問題; 1 = 有。
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

ALLOWED_SUBDIRS = {"rules", "references", "examples", "scripts"}
REQUIRED_SECTIONS = ("## 使用時機", "## Skill Boundaries")
NAME_PATTERN = re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*$")
CARD_ROWS = ("職責", "讀", "寫", "上游", "下游", "不做")
BODY_LINE_LIMIT = 500

MUST, SHOULD = "MUST", "SHOULD"


class Finding:
    def __init__(self, rule: str, level: str, message: str):
        self.rule, self.level, self.message = rule, level, message

    def __str__(self) -> str:
        mark = "✘" if self.level == MUST else "⚠"
        return f"  {mark} [{self.rule}] {self.level}: {self.message}"


def split_frontmatter(text: str) -> tuple[str | None, str]:
    """回傳 (frontmatter 原文, 本文)。沒有 frontmatter 時前者為 None."""
    if not text.startswith("---"):
        return None, text
    end = text.find("\n---", 3)
    if end == -1:
        return None, text
    return text[3:end], text[end + 4 :]


def parse_scalar(frontmatter: str, key: str) -> str | None:
    """抓出一個純量欄位的值，支援 YAML 的 `>` 折疊區塊。

    刻意不依賴 pyyaml：稽核腳本應該在任何環境都跑得起來，而這裡只需要讀
    兩三個欄位。
    """
    lines = frontmatter.splitlines()
    for i, line in enumerate(lines):
        if not line.startswith(f"{key}:"):
            continue
        inline = line[len(key) + 1 :].strip()
        if inline and inline not in (">", "|", ">-", "|-"):
            return inline.strip("\"'")
        collected = []
        for cont in lines[i + 1 :]:
            if cont.strip() and not cont.startswith((" ", "\t")):
                break
            collected.append(cont.strip())
        return " ".join(c for c in collected if c)
    return None


def referenced_files(text: str) -> set[str]:
    """SKILL.md 及其引用鏈中提到的子目錄檔名。"""
    return set(re.findall(r"(?:rules|references|examples|scripts)/[\w./-]+", text))


def section_body(body: str, title: str) -> str | None:
    """`## <title>` 底下到下一個 `## ` 為止的內容；沒有這一節回 None。"""
    m = re.search(rf"^## {re.escape(title)}[^\n]*\n(.*?)(?=^## |\Z)", body, re.M | re.S)
    return m.group(1) if m else None


def paths(cell: str) -> set[str]:
    """反引號裡像檔案的東西（有 `/` 或副檔名、沒有空白）；`@wip`、`#### 缺口` 不算。"""
    spans = re.findall(r"`([^`]+)`", cell)
    return {t for t in spans if not re.search(r"\s", t) and ("/" in t or re.search(r"\.[A-Za-z]+$", t))}


def audit(skill_dir: Path) -> list[Finding]:
    findings: list[Finding] = []
    add = lambda r, lv, m: findings.append(Finding(r, lv, m))  # noqa: E731

    skill_md = skill_dir / "SKILL.md"
    if not skill_md.is_file():
        add("S1", MUST, "找不到 SKILL.md")
        return findings

    text = skill_md.read_text(encoding="utf-8")
    frontmatter, body = split_frontmatter(text)

    if frontmatter is None:
        add("F0", MUST, "SKILL.md 沒有 --- 包住的 frontmatter")
        return findings

    if "：" in "\n".join(
        ln.split("：")[0] + "：" for ln in frontmatter.splitlines() if re.match(r"^\w+：", ln)
    ):
        add("F8", MUST, "frontmatter 的 key 使用了全形冒號「：」，YAML 會解析失敗")

    name = parse_scalar(frontmatter, "name")
    if not name:
        add("F1", MUST, "frontmatter 缺少 name")
    else:
        if not NAME_PATTERN.match(name):
            add("F1", MUST, f"name「{name}」不是 kebab-case（全小寫、連字號分隔）")
        if name != skill_dir.name:
            add("F2", MUST, f"name「{name}」與目錄名「{skill_dir.name}」不一致")

    description = parse_scalar(frontmatter, "description")
    if not description:
        add("F4", MUST, "frontmatter 缺少 description —— 這是觸發的唯一依據")
    else:
        if "「" not in description:
            add("F5", MUST, "description 沒有用「」列出中文觸發詞")
        if "English:" not in description:
            add("F6", SHOULD, "description 沒有 English: 起頭的英文觸發詞")
        if len(description) < 40:
            add("F4", MUST, f"description 只有 {len(description)} 字，太短不足以判斷觸發時機")

    for section in REQUIRED_SECTIONS:
        if section not in body:
            rule = "S6" if "使用時機" in section else "S7"
            add(rule, MUST, f"缺少「{section}」章節")

    # S15: 標題底下（到下一個 `## ` 為止）至少要有一個 `- [ ]`，只有標題不算清單
    verification = re.search(r"^## Verification[ \t]*\n(.*?)(?=^## |\Z)", body, re.M | re.S)
    if not verification:
        add("S15", MUST, "缺少「## Verification」章節 —— 交回前逐項確認的清單")
    elif not re.search(r"^\s*- \[ \] ", verification.group(1), re.M):
        add("S15", MUST, "「## Verification」底下沒有任何 `- [ ]` 項目")

    # S16: 職責卡在標題底下、第一個 `## ` 之前，六列都要有
    preamble = re.split(r"^## ", body, maxsplit=1, flags=re.M)[0]
    missing = [row for row in CARD_ROWS if not re.search(rf"^\|\s*\*\*{row}\*\*\s*\|", preamble, re.M)]
    if missing:
        add("S16", SHOULD, f"標題底下的職責卡缺少：{'、'.join(missing)}")

    flow = section_body(body, "流程")
    if flow is not None:
        # S17: 編號步驟（`### N.`）的第一行是「產出：」；沒編號的小節（續跑、不越界）不算步驟
        for title, rest in re.findall(r"^### ((?:Step )?\d+\b.*)\n((?:.*\n?)*?)(?=^### |\Z)", flow, re.M):
            first = next((ln for ln in rest.splitlines() if ln.strip()), "")
            if not first.startswith("產出："):
                add("S17", SHOULD, f"步驟「{title.strip()}」第一行不是「產出：…」")

        # S18: 產物表第一欄的檔，與職責卡「寫」裡的檔，是同一份清單
        card_write = re.search(r"^\|\s*\*\*寫\*\*\s*\|(.*)\|\s*$", preamble, re.M)
        written = paths(card_write.group(1)) if card_write else set()
        products = section_body(body, "產物") or ""
        table = re.search(r"^\|\s*檔\s*\|\s*哪一步寫\s*\|\s*格式在哪\s*\|.*\n((?:\|.*\n?)*)", products, re.M)
        if not table:
            add("S18", SHOULD, "「## 產物」沒有「| 檔 | 哪一步寫 | 格式在哪 |」表")
        else:
            rows = {p for ln in table.group(1).splitlines() for p in paths(ln.split("|")[1])}
            if rows - written:
                add("S18", SHOULD, f"產物表有、職責卡「寫」沒有：{'、'.join(sorted(rows - written))}")
            if written - rows:
                add("S18", SHOULD, f"職責卡「寫」有、產物表沒有：{'、'.join(sorted(written - rows))}")

    body_lines = len(body.splitlines())
    if body_lines > BODY_LINE_LIMIT:
        add("S10", SHOULD, f"本文 {body_lines} 行超過 {BODY_LINE_LIMIT} 行，考慮把細節下放到子目錄")

    unexpected = {
        d.name for d in skill_dir.iterdir() if d.is_dir() and d.name not in ALLOWED_SUBDIRS
    }
    if unexpected:
        add("S3", MUST, f"出現不允許的子目錄：{', '.join(sorted(unexpected))}")

    # S5: 子目錄的每個 .md 都要被引用鏈提到。逐層展開，因為 SKILL.md 可能只指到
    # 某個 rules 檔，而那個檔再指到別的。
    mentioned = referenced_files(text)
    seen: set[Path] = set()
    frontier = [skill_dir / m for m in mentioned]
    while frontier:
        p = frontier.pop()
        if p in seen or not p.is_file():
            continue
        seen.add(p)
        if p.suffix == ".md":
            for m in referenced_files(p.read_text(encoding="utf-8")):
                frontier.append(skill_dir / m)
                mentioned.add(m)

    for sub in sorted(ALLOWED_SUBDIRS):
        d = skill_dir / sub
        if not d.is_dir():
            continue
        if not any(d.iterdir()):
            add("S4", SHOULD, f"{sub}/ 是空目錄 —— 有內容才建")
            continue
        for f in sorted(d.rglob("*")):
            if not f.is_file() or f.name.startswith("."):
                continue
            rel = f.relative_to(skill_dir).as_posix()
            if rel not in mentioned:
                add("S5", MUST, f"{rel} 沒有被任何檔案指路，模型永遠不會讀到它")

    return findings


def main(argv: list[str]) -> int:
    if len(argv) != 2:
        print(__doc__)
        return 2

    root = Path(argv[1]).resolve()
    if not root.is_dir():
        print(f"不是目錄：{root}")
        return 2

    targets = [root] if (root / "SKILL.md").is_file() else sorted(
        d for d in root.iterdir() if d.is_dir() and (d / "SKILL.md").is_file()
    )
    if not targets:
        print(f"{root} 底下找不到任何 SKILL.md")
        return 2

    failed = False
    for skill_dir in targets:
        findings = audit(skill_dir)
        musts = [f for f in findings if f.level == MUST]
        status = f"✘ {len(musts)} 項必須修正" if musts else (
            f"⚠ {len(findings)} 項建議" if findings else "✔ 通過"
        )
        print(f"{skill_dir.name}: {status}")
        for f in findings:
            print(f)
        failed = failed or bool(musts)

    return 1 if failed else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
