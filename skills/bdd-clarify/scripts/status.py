#!/usr/bin/env python3
"""印出所有 feature 的澄清進度。

用法：
    python3 status.py [專案根目錄]        # 預設當前目錄

這是一個**算出來的視圖，不是存起來的檔案**。每個數字都直接數自來源：
規則與例子數自 example-mapping.md，紅卡數自 questions/ 的狀態列（含根層的
Pass 1 問題），就緒判定讀自 map 的檔頭。存一份儀表板的話，同一個數字會有兩份，
而它們遲早不一樣——這個 repo 已經發生過（map 檔頭寫 21 個例子，實際 23）。

只讀不寫。
"""
import re
import sys
from pathlib import Path

DIMS = ["空與零", "邊界", "重複", "時序", "權限", "失敗", "時間", "規模",
        "降級", "時限", "可觀測"]


def readiness(text: str) -> str:
    """map 檔頭的就緒判定。找不到就說找不到，不猜。"""
    m = re.search(r"^\*\*就緒判定\*\*:\s*(.+)$", text, re.M)
    if not m:
        return "？未寫"
    line = m.group(1)
    if "未就緒" in line:
        return "未就緒"
    return "已就緒" if "就緒" in line else line[:12]


def coverage(text: str) -> str:
    """追問覆蓋壓成一行符號。

    那一段會跨行、每個面向後面常接一段括號理由，所以先把整段接成一行再解析。
    只讀緊接在面向名稱後面的那個記號，不然括號裡的 ✓（例如「Rule 3 ✓」）會誤判。
    """
    m = re.search(r"^## 追問覆蓋\s*\n(.*?)(?=\n## |\Z)", text, re.M | re.S)
    if not m:
        return "?" * len(DIMS)
    flat = " ".join(m.group(1).split())
    out = []
    for d in DIMS:
        hit = re.search(re.escape(d) + r"\s*(n/a|✓|—|-)", flat)
        out.append({"n/a": "n", "✓": "✓"}.get(hit.group(1), "—") if hit else "?")
    return "".join(out)


def count_questions(qdir: Path) -> tuple[int, list[str]]:
    """回傳（已答數, 待答的檔名）。

    狀態只認檔內的狀態列，不認所在目錄——問題檔不因狀態改變而搬家。
    """
    answered, pending = 0, []
    for c in sorted(qdir.glob("*.md")):
        if re.search(r"^\*\*狀態\*\*:\s*已答", c.read_text(encoding="utf-8"), re.M):
            answered += 1
        else:
            pending.append(c.stem)
    return answered, pending


def main(root: Path) -> int:
    bdd = root / "docs" / "bdd"
    if not bdd.is_dir():
        print(f"找不到 {bdd}")
        return 1

    maps = sorted(bdd.glob("*/example-mapping.md"))
    if not maps:
        print(f"{bdd} 底下沒有 example-mapping.md")
        return 1

    print(f"{'feature':30}{'規則':>5}{'例子':>5}{'已答':>5}{'待答':>5}  {'狀態':8} 追問覆蓋")
    print(f"{'':30}{'':>20}  {'':8} 空邊複序權敗期模降限觀")
    tot = [0, 0, 0, 0]
    blocked = []

    # Pass 1 的問題住在根層，不屬於任何一則 story。漏掉這一列等於漏掉
    # 跨 feature 的未知，而那正是單看某一個 feature 永遠看不見的部分。
    root_q = bdd / "questions"
    if root_q.is_dir():
        answered, pending = count_questions(root_q)
        tot[2] += answered
        tot[3] += len(pending)
        blocked += [f"questions/{n}" for n in pending]
        print(f"{'questions/ (Pass 1)':30}{'-':>5}{'-':>5}{answered:>5}{len(pending):>5}"
              f"  {'跨 feature':8}")

    for m in maps:
        slug = m.parent.name
        t = m.read_text(encoding="utf-8")
        rules = len(re.findall(r"^### Rule ", t, re.M))
        exs = len(re.findall(r"^- Example \d+\.\d+ \S", t, re.M))

        answered, pending = count_questions(m.parent / "questions")
        blocked += [f"{slug}/{n}" for n in pending]

        for i, v in enumerate((rules, exs, answered, len(pending))):
            tot[i] += v
        print(f"{slug:30}{rules:>5}{exs:>5}{answered:>5}{len(pending):>5}  "
              f"{readiness(t):8} {coverage(t)}")

    print(f"{'合計':28}{tot[0]:>5}{tot[1]:>5}{tot[2]:>5}{tot[3]:>5}")

    # 待答的問題就是紅卡。列出來，因為「還剩什麼」比「已經做了多少」有用。
    if blocked:
        print(f"\n待答（{len(blocked)}）：")
        for b in blocked:
            print(f"  {b}")
    else:
        print("\n待答：0")

    print("\n追問覆蓋 空與零/邊界/重複/時序/權限/失敗/時間/規模/降級/時限/可觀測："
          "✓ 問過 · — 還沒問 · n 不適用 · ? 沒寫這一段")
    return 0


if __name__ == "__main__":
    sys.exit(main(Path(sys.argv[1] if len(sys.argv) > 1 else ".").resolve()))
