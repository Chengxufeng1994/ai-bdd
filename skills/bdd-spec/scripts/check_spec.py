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

只讀不寫。找不到 specs/ 或 .feature 時直接說找不到，不猜。
規則與例子讀自 specs/<slice>/clarify.md 的 `## Business Rules` 段。
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


def readiness(text: str) -> str:
    """clarify.md 檔頭的就緒判定。找不到就回空字串，不猜。

    稽核需要它來分辨兩種「沒有規則」：**還沒澄清**與**澄清完但漏寫**。
    前者是流程的正常中間狀態，後者是缺陷。把它們算成同一件事，
    等於讓每一個還沒輪到的批次都變成一個假警報——而假警報多了，
    真警報就沒人看。
    """
    m = re.search(r"^\*\*就緒判定\*\*[：:]\s*(.+)$", text, re.M)
    return m.group(1).strip() if m else ""


def stories_in_clarify(text: str) -> dict[str, set[str]]:
    """clarify.md 的 Business Rules 段：story-slug -> 例子編號集合。

    一份 clarify.md 涵蓋一批 story，每則一個 `### <story-slug>` 區塊。

    只掃 `## Business Rules` 這一段。其他章節也會出現 `Example N.M`
    （Open Questions 與 Assumptions 都可能引用某個例子），掃全檔會把那些
    引用當成規格來源——而它們是指回這一段的，比對它們等於自己跟自己比。

    來源是 clarify.md 而不是 spec.md：規則與例子的定版編號誕生在 CLARIFY，
    spec.md 不再重述它們。稽核要成立需要兩份**獨立**的表述（規則清單與
    可執行場景）；第三份只會製造兩個可能來源，然後各自漂移。
    """
    m = re.search(r"^## Business Rules\s*\n(.*?)(?=\n## |\Z)", text, re.M | re.S)
    if not m:
        return {}
    out: dict[str, set[str]] = {}
    for blk in re.split(r"^### ", m.group(1), flags=re.M)[1:]:
        slug = blk.splitlines()[0].strip()
        out[slug] = set(re.findall(r"^- Example (\d+\.\d+) \S", blk, re.M))
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

    clarifies = sorted(specs_dir.glob("*/clarify.md"))
    if not clarifies:
        print(f"{specs_dir} 底下沒有 clarify.md —— 先跑 bdd-clarify")
        return 1

    deferred: list[str] = []       # 還沒澄清完的批次，不算問題
    covered: set[str] = set()      # 有出現在某份 clarify.md 裡的 story slug
    for clarify_path in clarifies:
        ctext = clarify_path.read_text(encoding="utf-8")
        stories = stories_in_clarify(ctext)
        if not stories:
            verdict = readiness(ctext)
            if "已就緒" in verdict:
                print(f"✗ {clarify_path.parent.name:30} —— 判定已就緒，卻沒有 `## Business Rules`")
                problems += 1
            else:
                deferred.append(f"{clarify_path.parent.name}（{verdict or '就緒判定未寫'}）")
            continue

        for slug, exs in sorted(stories.items()):
            covered.add(slug)
            fpath = feat_dir / f"{slug}.feature"
            if not fpath.exists():
                print(f"{slug:30} —— 還沒寫成 .feature")
                continue

            ftext = fpath.read_text(encoding="utf-8")
            # ↓ 以下沿用原本的 tags/missing/invented/states/zh_rule 那一整段，
            #   只把 `exs` 的來源換掉，其餘一字不動。
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

            # 一條規則只走過一種結果不是錯，但要看得見——多數時候它代表沒問過
            # 「這條規則被違反時會怎樣」。邊界機器判不了，所以只報成功／失敗。
            lop = [n for n, ok, ng in outcome_coverage(ftext) if bool(ok) != bool(ng)]

            explained = sorted(set(missing) - set(unexplained), key=key)
            status = "✗" if issues else "✓"
            print(f"{status} {slug:30} {len(exs & tags):>3}/{len(exs)} 例子 · {states[0] if len(states)==1 else '?':7}"
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
        print(f"\n沒有對應 map 的 .feature（不在本檢查範圍）：{orphans}")

    if deferred:
        print(f"\n尚未澄清完、本次不稽核的批次：{deferred}")

    print(f"\n{'全部通過' if not problems else f'{problems} 個問題'}")
    return 1 if problems else 0


if __name__ == "__main__":
    sys.exit(check(Path(sys.argv[1] if len(sys.argv) > 1 else ".").resolve()))
