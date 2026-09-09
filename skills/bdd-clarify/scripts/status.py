#!/usr/bin/env python3
"""印出所有 feature 的澄清進度。

用法：
    python3 status.py [專案根目錄]        # 預設當前目錄

這是一個**算出來的視圖，不是存起來的檔案**。每個數字都直接數自
`specs/<date>-<feature>/prd.md` 的 `## Open Questions` 表：已答／n/a／待答數自
狀態欄，追問覆蓋（業務面向＋Pass 3 的技術面向）數自面向欄。存一份儀表板的話，
同一個數字會有兩份，而它們遲早不一樣——這個 repo 已經發生過（map 檔頭寫 21
個例子，實際 23）；追問覆蓋改成從表格算，就是為了不讓它也走上同一條路。

找不到資料時不猜、不印出「已就緒」——`specs/` 不存在、或底下沒有任何
`prd.md`，兩種都算「找不到」，說清楚是哪一種、離開碼非 0，不拿看得懂的那一
部分硬撐出一張表。

只讀不寫。
"""
import re
import sys
from pathlib import Path
from unicodedata import east_asian_width

# Pass 2 的十一個業務面向。
BIZ_DIMS = ["空與零", "邊界", "重複", "時序", "權限", "失敗", "時間", "規模",
            "降級", "時限", "可觀測"]
# Pass 3 的五個技術面向，名稱與順序取自 references/technical-probes.md 的
# 章節標題，不是憑空編的縮寫——兩邊不一致以那份為準。沒有這一半，「seam
# 沒問到」在儀表板上完全看不見，而那正是這次改版要補的洞。
TECH_DIMS = ["seam", "模組邊界", "介面與型別契約", "排序契約／決定性",
             "既有資產／測試慣例"]
DIMS = BIZ_DIMS + TECH_DIMS

# 表頭縮寫，跟 DIMS 一一對應。技術面向用單字或字母另挑一個沒被業務面向占走
# 的符號，讀表時才分得出「這格是業務還是技術」——圖例另外把兩組分開講一次。
BIZ_ABBR = "空邊複序權敗期模降限觀"
TECH_ABBR = "S組型定慣"


def _field(label: str, width: int) -> str:
    """等寬字型裡全形字元占兩格，Python 算字串長度只算一格——
    欄寬要扣掉全形字元數，不然含中文的列（例如「全部」「合計」）
    會把後面的數字欄位擠歪。"""
    wide = sum(1 for ch in label if east_asian_width(ch) in ("W", "F"))
    return f"{label:<{max(width - wide, 0)}}"


def parse_open_questions(text: str) -> list[dict[str, str]]:
    """prd.md 的 `## Open Questions` 表 -> 每題一個欄位字典。

    只讀那一節的表格列，不掃整份檔案：prd.md 的其他章節也有表格
    （Actors、NFR），對整份掃會把它們算成問題。

    表頭固定五欄 `| Q | 問題 | 面向 | 狀態 | 答案 |`；欄數不符的列
    直接丟進 unparseable，不猜。
    """
    section = re.search(
        r"^## Open Questions\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return []
    rows = []
    for line in section.group(1).splitlines():
        line = line.strip()
        if not line.startswith("|"):
            continue
        cells = [c.strip() for c in line.strip("|").split("|")]
        if len(cells) != 5 or cells[0] in ("Q", "---"):
            continue
        if set(cells[0]) <= {"-", " "}:
            continue
        rows.append(dict(zip(("q", "問題", "面向", "狀態", "答案"), cells)))
    return rows


def coverage(dims: set[str]) -> str:
    """把「有沒有人問過這個面向」壓成一行符號。

    ✓ 這個面向至少有一個問題檔 · — 一題都沒有

    原本讀的是 map 裡手寫的一段覆蓋表。改成算的之後，「寫了覆蓋表但沒真的問」
    這種狀態不再可能存在——那正是手寫彙總遲早會跟來源不一致的地方。
    """
    return "".join("✓" if d in dims else "—" for d in DIMS)


def main(root: Path) -> int:
    specs = root / "specs"
    if not specs.is_dir():
        print(f"找不到 {specs}")
        return 1

    prds = sorted(specs.glob("*/prd.md"))
    if not prds:
        print(f"{specs} 底下沒有任何 prd.md —— 先跑 bdd-clarify")
        return 1

    print(f"{_field('feature', 42)}{'已答':>6}{'n/a':>6}{'待答':>6}  追問覆蓋")
    print(f"{'':42}{'':>18}  {BIZ_ABBR}{TECH_ABBR}")
    tot = [0, 0, 0]
    blocked = []

    for prd in prds:
        slug = prd.parent.name
        rows = parse_open_questions(prd.read_text(encoding="utf-8"))
        answered = na = pending = 0
        dims: set[str] = set()
        pending_qs = []
        for row in rows:
            # `—` 是 Pass 1 的範圍題，不屬於任何探針面向，不計入追問覆蓋
            # ——跟 n/a 不同，n/a 代表「問過某個面向、判定不適用」。
            dim = row.get("面向", "")
            if dim and dim != "—":
                dims.add(dim)
            status = row.get("狀態", "")
            if status.startswith("已答"):
                answered += 1
            elif status.startswith("n/a"):
                # n/a 是「問過、判定不適用」的結論，不是還沒處理。待答數與
                # 待答清單存在的目的是讓讀者知道還要去追誰要答案——n/a 沒
                # 有答案可追，跟真正卡住的問題混在一起算，會讓一份 PRD 永
                # 遠有消不掉的紅卡、看起來比實際更沒準備好。獨立開一欄，
                # 不計入待答。
                na += 1
            else:
                pending += 1
                pending_qs.append(row.get("q", "?"))
        tot[0] += answered
        tot[1] += na
        tot[2] += pending
        blocked += [f"{slug}/{q}" for q in pending_qs]
        print(f"{_field(slug, 42)}{answered:>6}{na:>6}{pending:>6}  "
              f"{coverage(dims)}")

    print(f"{_field('合計', 42)}{tot[0]:>6}{tot[1]:>6}{tot[2]:>6}")

    # 待答的問題就是紅卡。列出來，因為「還剩什麼」比「已經做了多少」有用。
    if blocked:
        print(f"\n待答（{len(blocked)}）：")
        for b in blocked:
            print(f"  {b}")
    else:
        print("\n待答：0")

    print("\n追問覆蓋 · 業務面向（Pass 2）"
          "空與零/邊界/重複/時序/權限/失敗/時間/規模/降級/時限/可觀測："
          "✓ 問過（含 n/a）· — 還沒問")
    print("追問覆蓋 · 技術面向（Pass 3）"
          "seam/模組邊界/介面與型別契約/排序契約／決定性/既有資產／測試慣例："
          "✓ 問過（含 n/a）· — 還沒問")
    return 0


if __name__ == "__main__":
    sys.exit(main(Path(sys.argv[1] if len(sys.argv) > 1 else ".").resolve()))
