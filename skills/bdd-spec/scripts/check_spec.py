#!/usr/bin/env python3
"""稽核 .feature 與 example map 的一致性。

用法：
    python3 check_spec.py [專案根目錄]        # 預設當前目錄

檢查五件事，全部是機械性的：

  1. 覆蓋——雙向。map 有而 feature 沒有的（漏做），
     以及 feature 指向 map 裡不存在的例子（發明出來的驗收條件）。
  2. 缺口——漏掉的例子有沒有在 .feature 裡就地留註解交代。
  3. 狀態 tag——每個有 map 的 .feature 恰好一個。
  4. 方言陷阱——中文「規則:」不是關鍵字，會被解析成散文。
  5. 步驟樣板重用率——封閉文法有沒有真的套上。沒有及格線，只報數字：
     散文式實測 76 場景／208 樣板；封閉文法 7 場景／10 樣板。

退出碼 0 = 全過，1 = 有問題。適合放進 CI。

只讀不寫。找不到 docs/bdd/ 或 .feature 時直接說找不到，不猜。
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


def examples_in_map(text: str) -> set[str]:
    """map 的例子。備註區回指例子的句子不算，所以要求編號後面接空白再接內容。"""
    return set(re.findall(r"^- Example (\d+\.\d+) \S", text, re.M))


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


def check(root: Path) -> int:
    bdd = root / "docs" / "bdd"
    feat_dir = find_features(root)
    if not bdd.is_dir():
        print(f"找不到 {bdd} —— 沒有 CLARIFY 的產物可以比對")
        return 1
    if feat_dir is None:
        print("找不到任何 .feature")
        return 1

    problems = 0
    print(f"map: {bdd}    feature: {feat_dir}\n")

    for map_path in sorted(bdd.glob("*/example-mapping.md")):
        slug = map_path.parent.name
        fpath = feat_dir / f"{slug}.feature"
        if not fpath.exists():
            print(f"{slug:30} —— 還沒寫成 .feature")
            continue

        mtext = map_path.read_text(encoding="utf-8")
        ftext = fpath.read_text(encoding="utf-8")

        exs = examples_in_map(mtext)
        tags = set(re.findall(r"@example-(\d+\.\d+)", ftext))
        key = lambda s: tuple(map(int, s.split(".")))
        missing = sorted(exs - tags, key=key)
        invented = sorted(tags - exs, key=key)

        states = [s for s in STATES if re.search(rf"^{s}\b", ftext, re.M)]
        zh_rule = re.findall(r"^\s*規則:", ftext, re.M)

        issues = []
        # 發明優先於漏做：憑空的驗收條件比缺一條更難發現，因為它看起來很完整。
        if invented:
            issues.append(f"指向 map 裡不存在的例子 {invented}")
        if len(states) != 1:
            issues.append(f"狀態 tag {states or '缺'} —— 每個檔恰好要一個")
        if zh_rule:
            issues.append(f"用了中文「規則:」{len(zh_rule)} 處 —— 會被解析成散文，不會報錯")

        # 漏做不一定是錯：標「暫定」的例子本來就不該寫成場景。要求就地留註解交代。
        unexplained = [e for e in missing
                       if not re.search(rf"^\s*#.*Example {re.escape(e)}", ftext, re.M)]
        if unexplained:
            issues.append(f"漏了 {unexplained} 且檔案裡沒有註解說明")

        explained = sorted(set(missing) - set(unexplained), key=key)
        status = "✗" if issues else "✓"
        print(f"{status} {slug:30} {len(exs & tags):>3}/{len(exs)} 例子 · {states[0] if len(states)==1 else '?':7}"
              + (f" · 已交代不寫 {explained}" if explained else ""))
        for i in issues:
            print(f"    ✗ {i}")
            problems += 1

    written = [p.read_text(encoding="utf-8") for p in feat_dir.glob("*.feature")
               if (bdd / p.stem / "example-mapping.md").exists()]
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
               if not (bdd / p.stem / "example-mapping.md").exists()]
    if orphans:
        print(f"\n沒有對應 map 的 .feature（不在本檢查範圍）：{orphans}")

    print(f"\n{'全部通過' if not problems else f'{problems} 個問題'}")
    return 1 if problems else 0


if __name__ == "__main__":
    sys.exit(check(Path(sys.argv[1] if len(sys.argv) > 1 else ".").resolve()))
