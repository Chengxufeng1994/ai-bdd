#!/usr/bin/env python3
"""稽核 .feature 與 prd.md／spec.md 的一致性。

用法：
    python3 check_spec.py [專案根目錄]        # 預設當前目錄

檢查六件事，全部是機械性的：

  1. 覆蓋——雙向。prd.md 有例子而 feature 沒有的（漏做），
     以及 feature 指向 prd.md 裡不存在的例子（發明出來的驗收條件）。
  2. 缺口——漏掉的例子有沒有在 .feature 裡就地留註解交代。
  3. 狀態 tag——每個有對應 story 的 .feature 恰好一個。
  4. 方言陷阱——中文「規則:」不是關鍵字，會被解析成散文。
  5. 步驟樣板重用率——封閉文法有沒有真的套上。沒有及格線，只報數字：
     散文式實測 76 場景／208 樣板；封閉文法 7 場景／10 樣板。
  6. FR 完整性——prd.md 裡有沒有 FR 完全沒掛任何例子；沒有例子，
     SPEC 就無從寫出場景，只能發明。

退出碼 0 = 全過，1 = 有問題。適合放進 CI。

只讀不寫。找不到 specs/ 或 .feature 時直接說找不到，不猜。
例子的目錄讀自 specs/<date>-<feature>/prd.md 的 `## Functional Requirements` 段
（FR 編號在那裡定版）；story 由哪些 FR 組成讀自 specs/<date>-<feature>/spec.md 的
`## Stories` 段——切 story 是 SPEC 的事，不再是 CLARIFY 的事。
"""
import re
import sys
from pathlib import Path

STATES = ("@draft", "@ready", "@wip", "@review", "@done")


def find_features(root: Path) -> Path:
    """.feature 的位置由測試框架決定，所以用找的，不用猜。"""
    for c in (root / "features", root / "test" / "features", root / "tests" / "features"):
        if c.is_dir():
            return c
    hits = {p.parent for p in root.rglob("*.feature")}
    return sorted(hits, key=lambda p: len(p.parts))[0] if hits else None


def frs_in_prd(text: str) -> dict[str, set[str]]:
    """prd.md 的 `## Functional Requirements` 段：FR 編號 -> EX 編號集合。

    來源是 prd.md 而不是 spec.md：例子的定版編號誕生在 CLARIFY，
    spec.md 只記哪些 FR 湊成一則 story，不重述例子。
    """
    section = re.search(
        r"^## Functional Requirements\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return {}
    out: dict[str, set[str]] = {}
    current = None
    for line in section.group(1).splitlines():
        fr = re.match(r"^### FR-(\d+)\b", line)
        if fr:
            current = fr.group(1)
            out.setdefault(current, set())
            continue
        if current:
            for ex in re.findall(r"\bEX-(\d+\.\d+)\b", line):
                out[current].add(ex)
    return out


def stories_in_spec(text: str) -> dict[str, set[str]]:
    """spec.md 的 `## Stories` 段：story-slug -> 它涵蓋的 FR 編號集合。"""
    section = re.search(r"^## Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return {}
    out: dict[str, set[str]] = {}
    current = None
    for line in section.group(1).splitlines():
        st = re.match(r"^### (\S+)\s*$", line)
        if st:
            current = st.group(1)
            out.setdefault(current, set())
            continue
        if current:
            for fr in re.findall(r"\bFR-(\d+)\b", line):
                out[current].add(fr)
    return out


def step_templates(texts: list[str]) -> tuple[int, int, int]:
    """把引號內容、佔位符、數字正規化之後，數不重複的步驟樣板。

    正規化是重點：`「臥推」` 與 `「深蹲」` 是同一支 step definition 的兩次呼叫，
    不該算成兩個樣板。沒有正規化的話這個指標會永遠很難看，也就沒人會看。
    """
    seen: dict[str, int] = {}
    steps = 0
    for t in texts:
        for raw in re.findall(r"^\s*(?:Given|When|Then|And|But)\s+(.*)$", t, re.M):
            steps += 1
            n = re.sub(r"「[^」]*」", "「X」", raw.strip())
            n = re.sub(r'"[^"]*"', '"X"', n)
            # Outline 的 <佔位符> 在執行前就被代換成值，所以它和具體值版本
            # 由同一支 step definition 接住——正規化成同一個 token，否則
            # 每個 Outline 都會被算成多一個樣板，指標系統性高估。
            n = re.sub(r"<[^>]*>", "N", n)
            n = re.sub(r"\d+(\.\d+)?", "N", n).rstrip("：:")
            seen[n] = seen.get(n, 0) + 1
    return steps, len(seen), sum(1 for v in seen.values() if v == 1)


def outcome_coverage(text: str) -> list[tuple[str, int, int]]:
    """每條 Rule 走過幾個成功結果、幾個失敗結果。

    只數 `Then 操作成功` / `Then 操作失敗`——邊界是判斷題，機器數不出來，
    所以這裡不假裝數得出來，只把成功與失敗的分布攤開。
    """
    out = []
    for blk in re.split(r"^\s*Rule:", text, flags=re.M)[1:]:
        name = blk.split("\n", 1)[0].strip()
        ok = len(re.findall(r"^\s*Then 操作成功", blk, re.M))
        ng = len(re.findall(r"^\s*Then 操作失敗", blk, re.M))
        out.append((name, ok, ng))
    return out


def check(root: Path) -> int:
    specs_dir = root / "specs"
    feat_dir = find_features(root)
    if not specs_dir.is_dir():
        print(f"找不到 {specs_dir} —— 沒有 CLARIFY 的產物可以比對")
        return 1
    if feat_dir is None:
        print("找不到任何 .feature")
        return 1

    problems = 0
    print(f"map: {specs_dir}    feature: {feat_dir}\n")

    prds = sorted(specs_dir.glob("*/prd.md"))
    if not prds:
        print(f"{specs_dir} 底下沒有 prd.md —— 先跑 bdd-clarify")
        return 1

    covered: set[str] = set()      # 有出現在某份 spec.md 裡的 story slug
    for prd_path in prds:
        ptext = prd_path.read_text(encoding="utf-8")
        frs = frs_in_prd(ptext)
        if not frs:
            print(f"✗ {prd_path.parent.name} 的 prd.md 沒有 `## Functional Requirements` "
                  f"一節、或該節底下沒有任何 FR —— 先跑 bdd-clarify")
            problems += 1
            continue

        # 這是 D3 兩層結構唯一的機械防線：完全沒有例子的 FR，SPEC 無從寫出場景，
        # 只能發明——而發明的驗收條件正是這支腳本要抓的另一個方向。FR 編號是
        # 每個 feature 各自從 1 編起，「FR-3」在多 feature 的 repo 裡定位不到是
        # 哪一份，所以訊息要帶上是哪個 feature 目錄。
        barren = sorted(fr for fr, exs in frs.items() if not exs)
        if barren:
            print(f"✗ {prd_path.parent.name} 這些 FR 沒有任何 EX，SPEC 無從寫出場景："
                  f"{', '.join('FR-' + b for b in barren)}")
            problems += 1

        spec_path = prd_path.parent / "spec.md"
        if not spec_path.exists():
            print(f"✗ {prd_path.parent.name} 有 prd.md 但沒有 spec.md —— 先跑 bdd-spec")
            problems += 1
            continue
        stext = spec_path.read_text(encoding="utf-8")
        if not re.search(r"^## Stories\s*$", stext, re.M):
            print(f"✗ {prd_path.parent.name} 的 spec.md 找不到 `## Stories` 一節 —— 檔案格式損壞")
            problems += 1
            continue
        stories = stories_in_spec(stext)
        if not stories:
            print(f"✗ {prd_path.parent.name} 的 spec.md `## Stories` 一節底下沒有任何 story "
                  f"—— SPEC 還沒切 story")
            problems += 1
            continue

        # 逐 story 比，不可把所有 story 的例子聯集起來跟單一 .feature 比：
        # 一則 story 一個 .feature，聯集會讓 A 的例子出現在 B 檔裡也算通過。
        for slug, fr_ids in sorted(stories.items()):
            covered.add(slug)
            expected = {ex for fr in fr_ids for ex in frs.get(fr, set())}
            fpath = feat_dir / f"{slug}.feature"
            ftext = fpath.read_text(encoding="utf-8") if fpath.exists() else ""

            tags = set(re.findall(r"@example-(\d+\.\d+)", ftext))
            key = lambda s: tuple(map(int, s.split(".")))
            missing = sorted(expected - tags, key=key)
            invented = sorted(tags - expected, key=key)

            states = [s for s in STATES if re.search(rf"^{s}\b", ftext, re.M)]
            zh_rule = re.findall(r"^\s*規則:", ftext, re.M)

            issues = []
            # 發明優先於漏做：憑空的驗收條件比缺一條更難發現，因為它看起來很完整。
            if invented:
                issues.append(f"指向 prd.md 裡不存在的例子 {invented}")
            if len(states) != 1:
                issues.append(f"狀態 tag {states or '缺'} —— 每個檔恰好要一個")
            if zh_rule:
                issues.append(f"用了中文「規則:」{len(zh_rule)} 處 —— 會被解析成散文，不會報錯")

            # 漏做不一定是錯：標「暫定」的例子本來就不該寫成場景。要求就地留註解交代。
            unexplained = [e for e in missing
                           if not re.search(rf"^\s*#.*Example {re.escape(e)}", ftext, re.M)]
            if unexplained:
                issues.append(f"漏了 {unexplained} 且檔案裡沒有註解說明")

            # 一條規則只走過一種結果不是錯，但要看得見——多數時候它代表沒問過
            # 「這條規則被違反時會怎樣」。邊界機器判不了，所以只報成功／失敗。
            lop = [n for n, ok, ng in outcome_coverage(ftext) if bool(ok) != bool(ng)]

            explained = sorted(set(missing) - set(unexplained), key=key)
            status = "✗" if issues else "✓"
            print(f"{status} {slug:30} {len(expected & tags):>3}/{len(expected)} 例子 · {states[0] if len(states)==1 else '?':7}"
                  + (f" · 已交代不寫 {explained}" if explained else "")
                  + (f" · 單一結果的規則 {len(lop)}" if lop else ""))
            for i in issues:
                print(f"    ✗ {i}")
                problems += 1

    written = [p.read_text(encoding="utf-8") for p in feat_dir.glob("*.feature")
               if p.stem in covered]
    if written:
        n_steps, n_tpl, n_once = step_templates(written)
        n_scen = sum(len(re.findall(r"^\s*(?:Scenario|Scenario Outline|Example):", t, re.M))
                     for t in written)
        print(f"\n步驟 {n_steps} 行 · 不重複樣板 {n_tpl} · 只出現一次 {n_once} · 場景 {n_scen}")
        # 只出現一次的樣板 ≈ 一支只會被呼叫一次的 step definition。
        # 沒有及格線（那會是編出來的），但接近場景數就代表文法沒套上。
        if n_once >= n_scen:
            print(f"  ⚠ 只出現一次的樣板({n_once}) 已達場景數({n_scen})——"
                  f"封閉文法可能沒真的套上，step definition 會比場景還多")

    orphans = [p.name for p in feat_dir.glob("*.feature")
               if p.stem not in covered]
    if orphans:
        print(f"\n沒有對應 story 的 .feature（不在本檢查範圍）：{orphans}")

    print(f"\n{'全部通過' if not problems else f'{problems} 個問題'}")
    return 1 if problems else 0


if __name__ == "__main__":
    sys.exit(check(Path(sys.argv[1] if len(sys.argv) > 1 else ".").resolve()))
