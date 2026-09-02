#!/usr/bin/env python3
"""印出所有 slice 的澄清進度。

用法：
    python3 status.py [專案根目錄]        # 預設當前目錄

這是一個**算出來的視圖，不是存起來的檔案**。每個數字都直接數自
`questions/` 底下每個問題檔的狀態列：已答／n/a／待答數自狀態欄，story 歸屬
數自 story 欄，追問覆蓋（業務面向＋Pass 3 的技術面向）數自面向欄。
`example-mapping.md` 已經不存在，規則數、例子數、map 檔頭的就緒判定也一併
拿掉——沒有來源的欄位留著只會印出永遠不變的空白，讀的人會誤以為那代表
什麼。存一份儀表板的話，同一個數字會有兩份，而它們遲早不一樣——這個 repo
已經發生過（map 檔頭寫 21 個例子，實際 23）；追問覆蓋改成從問題檔算，就是
為了不讓它也走上同一條路。

找不到資料時不猜、不印出「已就緒」——目錄不存在、目錄底下沒有問題檔、或
問題檔解析不出任何欄位（狀態／story／面向都缺，代表格式跑掉或檔案壞了），
三種都算「找不到」，說清楚是哪一種、離開碼非 0，不拿看得懂的那一部分硬撐
出一張表。

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


def parse_status(text: str) -> dict[str, str]:
    """問題檔第三行的狀態列 -> 欄位字典。

    格式是 `**狀態**: 已答 · **輪次**: 2 · **story**: log-a-workout · **面向**: 邊界`。
    分隔符是全形間隔號，欄名用 `**...**` 包住——兩者都固定，所以一條 regex 掃完
    比逐欄寫 pattern 好維護。
    """
    return {k: v.strip() for k, v in
            re.findall(r"\*\*(.+?)\*\*:\s*([^·\n]+)", text)}


def scan_questions(qdir: Path) -> tuple[dict[str, dict], list[str]]:
    """走一遍問題檔，依 story 分組；順便回報解析不出任何欄位的檔名。

    回傳 (groups, unparseable)。
    groups: {story: {"answered": int, "na": int, "pending": [檔名], "dims": {面向}}}
    `story` 欄缺席時歸到 `全部`——判斷不了時放共通層，代價不對稱：
    共通的問題每則 story 都看得到，掛錯 story 的別則看不到。

    unparseable：`狀態`／`story`／`面向` 三欄一個都讀不到的檔名。這種檔案
    不能套「story 缺席就歸到全部」的規則——那條規則假設的是「有欄位、只是
    沒填 story」，這裡是整份都讀不出東西，通常代表狀態列格式跑掉、檔案被
    手動改壞，或編碼出了問題。呼叫端要能分辨「這批問題真的還沒問」跟「這
    份資料讀不出來」，不能讓後者悄悄併入前者、印出一張看起來正常的表。

    `面向` 只要出現就算進 dims，不管狀態是已答、待答還是 n/a。n/a 的問題檔
    一樣帶著 `面向`，代表這個面向有人問過、也有人給了「不適用」的理由，跟
    從來沒問過完全不同——如果 n/a 不算進 dims，追問覆蓋就會把「問過但不
    適用」跟「沒人問」混成同一個 `—`，那正是這欄本來要讓人看見的東西。
    """
    out: dict[str, dict] = {}
    bad: list[str] = []
    for c in sorted(qdir.glob("*.md")):
        st = parse_status(c.read_text(encoding="utf-8"))
        if not (st.get("狀態") or st.get("story") or st.get("面向")):
            bad.append(c.stem)
            continue
        story = st.get("story", "全部")
        g = out.setdefault(story, {"answered": 0, "na": 0, "pending": [], "dims": set()})
        if dim := st.get("面向"):
            g["dims"].add(dim)
        status = st.get("狀態", "")
        if status.startswith("已答"):
            g["answered"] += 1
        elif status.startswith("n/a"):
            # n/a 是「問過、判定不適用」的結論，不是還沒處理。待答數與待答
            # 清單存在的目的是讓讀者知道還要去追誰要答案——n/a 沒有答案可
            # 追，跟真正卡住的問題混在一起算，會讓一批問題永遠有消不掉的
            # 紅卡、看起來比實際更沒準備好。獨立開一欄，不計入待答。
            g["na"] += 1
        else:
            g["pending"].append(c.stem)
    return out, bad


def coverage(dims: set[str]) -> str:
    """把「有沒有人問過這個面向」壓成一行符號。

    ✓ 這個面向至少有一個問題檔 · — 一題都沒有

    原本讀的是 map 裡手寫的一段覆蓋表。改成算的之後，「寫了覆蓋表但沒真的問」
    這種狀態不再可能存在——那正是手寫彙總遲早會跟來源不一致的地方。
    """
    return "".join("✓" if d in dims else "—" for d in DIMS)


def main(root: Path) -> int:
    bdd = root / "docs" / "bdd"
    if not bdd.is_dir():
        print(f"找不到 {bdd}")
        return 1

    qdirs = sorted(p for p in bdd.glob("*/questions") if p.is_dir())
    if not qdirs:
        print(f"{bdd} 底下沒有任何 slice 帶 questions/ 目錄")
        return 1

    scans = {qdir: scan_questions(qdir) for qdir in qdirs}

    # 兩種情況都不能往下印表：目錄底下一個 .md 都沒有，或者有 .md 卻解析不
    # 出狀態列（哪怕只有一個檔案解析不出來）。後者尤其不能只印看得懂的那
    # 部分——20/31 個檔案算出來的覆蓋表看起來完全正常，但算漏的那 11 個
    # 不會留下任何痕跡，這種「看起來對」的錯比整個失敗更難發現。
    problems = []
    for qdir, (groups, bad) in scans.items():
        if bad:
            problems.append((qdir, bad))
        elif not groups:
            problems.append((qdir, None))
    if problems:
        for qdir, bad in problems:
            if bad is None:
                print(f"{qdir} 底下沒有可解析的狀態列（沒有 .md 檔）")
            else:
                print(f"{qdir} 有 {len(bad)} 個問題檔解析不出狀態列"
                      f"（狀態／story／面向都缺）：")
                for name in bad:
                    print(f"  {name}")
        return 1

    print(f"{_field('story', 42)}{'已答':>6}{'n/a':>6}{'待答':>6}  追問覆蓋")
    print(f"{'':42}{'':>18}  {BIZ_ABBR}{TECH_ABBR}")
    tot = [0, 0, 0]
    blocked = []

    for qdir in qdirs:
        slug = qdir.parent.name
        groups = scans[qdir][0]
        # 全部（跨 story 的問題）排最前面：這批問題每則 story 都看得到，
        # 排最前面才不會被誤讀成只屬於表上緊接著的那一則 story。
        for story in sorted(groups, key=lambda s: (s != "全部", s)):
            g = groups[story]
            answered, na, pending = g["answered"], g["na"], len(g["pending"])
            tot[0] += answered
            tot[1] += na
            tot[2] += pending
            blocked += [f"{slug}/{n}" for n in g["pending"]]
            label = f"{slug}/{story}"
            print(f"{_field(label, 42)}{answered:>6}{na:>6}{pending:>6}  "
                  f"{coverage(g['dims'])}")

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
