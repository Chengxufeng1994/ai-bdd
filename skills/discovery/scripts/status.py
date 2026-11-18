#!/usr/bin/env python3
"""印出所有 feature 的澄清進度。

用法：
    python3 status.py [專案根目錄]        # 預設當前目錄

這是一個**算出來的視圖，不是存起來的檔案**。每個數字都直接數自
`specs/ongoing/<date>-<feature>/spec.md` 的 `## Open Questions` 表：已答／n/a／待答數自
狀態欄，追問覆蓋（業務面向＋技術面向）數自面向欄。存一份儀表板的話，
同一個數字會有兩份，而它們遲早不一樣——這個 repo 已經發生過（map 檔頭寫 21
個例子，實際 23）；追問覆蓋改成從表格算，就是為了不讓它也走上同一條路。

找不到資料時不猜、不印出「已就緒」。`specs/` 不存在、或底下沒有任何
`spec.md`，兩種都算「找不到」。單一 `spec.md` 也有三種算「壞」的情況，訊息
分開講因為修法不同：`## Open Questions` 一節整個缺席（標題缺漏或打錯）是
檔案格式損壞；一節存在但零列，代表 DISCOVERY 沒問過任何問題，不是需求沒有
疑點；表格裡有列的欄數不是 5，是那一列本身壞了。三種都指名是哪個檔案、
離開碼非 0，但不因為某一份壞掉就把其他份、或同一份裡還解析得出來的列藏
起來不印——安靜的「一切正常」比看得出來的「壞掉」更危險。

只讀不寫。
"""
import re
import sys
from pathlib import Path
from unicodedata import east_asian_width

# 前沿檢查表的十一個業務面向。
BIZ_DIMS = ["空與零", "邊界", "重複", "時序", "權限", "失敗", "時間", "規模",
            "降級", "時限", "可觀測"]
# 前沿檢查表的五個技術面向，名稱與順序取自 references/technical-probes.md 的
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


def parse_open_questions(
    text: str,
) -> tuple[bool, list[dict[str, str]], list[str]]:
    """spec.md 的 `## Open Questions` 表 -> (找到一節沒、題目列表、解析不出的列)。

    只讀那一節的表格列，不掃整份檔案：`## Document Overview` 底下也有表格，
    而**版本修訂歷史正好也是五欄**，跟下面用來擋壞列的欄數判斷相同——對整份
    掃，它的每一列都會安靜地被當成一題（核心關係人是三欄，會被擋下來，所以
    危險的不是它）。題數會隨著版本紀錄變長而虛增，且沒有任何一列進得了
    unparseable，也就沒有任何地方看得出來。

    第一個回傳值分開「這一節根本不存在」跟「存在但零列」——前者是標題
    缺漏或打錯，檔案格式跑掉；後者是 DISCOVERY 沒問過任何問題。呼叫端要能
    分開講，不能都印成同一種「沒有」。

    表頭固定五欄 `| Q | 問題 | 面向 | 狀態 | 答案 |`；表頭列與分隔列
    （純 `-`／空白組成）之外，欄數不是 5 的列收進第三個回傳值，不猜、不
    硬套進五欄格式——呼叫端要能報數，不能讓這種列悄悄從統計裡消失，看起
    來只是「題目比較少」。
    """
    section = re.search(
        r"^## Open Questions\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return False, [], []
    rows = []
    unparseable = []
    for line in section.group(1).splitlines():
        line = line.strip()
        if not line.startswith("|"):
            continue
        cells = [c.strip() for c in line.strip("|").split("|")]
        if cells[0] in ("Q", "---") or set(cells[0]) <= {"-", " "}:
            continue
        if len(cells) != 5:
            unparseable.append(line)
            continue
        rows.append(dict(zip(("q", "問題", "面向", "狀態", "答案"), cells)))
    return True, rows, unparseable


def coverage(dims: set[str]) -> str:
    """把「有沒有人問過這個面向」壓成一行符號。

    ✓ 這個面向至少有一個問題檔 · — 一題都沒有

    原本讀的是 map 裡手寫的一段覆蓋表。改成算的之後，「寫了覆蓋表但沒真的問」
    這種狀態不再可能存在——那正是手寫彙總遲早會跟來源不一致的地方。
    """
    return "".join("✓" if d in dims else "—" for d in DIMS)


def flat_layout_hint(specs: Path) -> str | None:
    """2026-11-15 以前的平鋪佈局（feature 直接放在 `specs/` 底下）的退件訊息；不是就回 None。

    新佈局的 feature 住 `specs/ongoing/`，做完搬到 `specs/sleeping/`。平鋪的目錄不擋下來
    就會被安靜地略過——腳本報「沒有任何 feature」，看起來像是還沒開始。
    """
    if (specs / "domain-model.md").exists():
        # 聚合的現況改住 docs/DOMAIN-MODEL.md，由 archive 在 feature 做完時增修
        return "specs/domain-model.md 是舊佈局：`git mv specs/domain-model.md docs/DOMAIN-MODEL.md`"
    for d in sorted(p for p in specs.iterdir() if p.is_dir() and p.name not in ("ongoing", "sleeping")):
        if any((d / f).exists() for f in ("discovery.md", "spec.md", "prd.md")):
            return (f"{d} 是舊佈局（feature 直接放在 specs/ 底下）："
                    f"`git mv specs/{d.name} specs/ongoing/{d.name}`")
    return None


def legacy_hint(specs_dir: Path) -> str | None:
    """舊佈局的退件訊息；沒有舊佈局回 None。

    兩種舊佈局都會被新邏輯讀錯而不報錯，所以要在讀之前擋下：10-02 的
    prd.md，與 11-18 以前把問答和原文放在旁邊 discovery.md 的兩檔佈局。
    """
    for prd in sorted(specs_dir.glob("*/prd.md")):
        # 技術 spec.md 還在旁邊時要先讓位，否則第一個 mv 就撞檔。
        hint = ("先 `git mv spec.md design.md`，再 `git mv prd.md spec.md`"
                if (prd.parent / "spec.md").exists() else "`git mv prd.md spec.md`")
        return f"{prd.parent} 是舊版佈局（prd.md）：{hint}"
    for disc in sorted(specs_dir.glob("*/discovery.md")):
        return (f"{disc.parent} 是 11-18 以前的兩檔佈局：把 `discovery.md` 的 `## Open Questions`"
                f"（表與小節）與 `## 原文` 搬進 `spec.md`、`## 粗切 story` 併進 `## User Stories`，"
                f"`spec.md` 舊的 `## Open Questions` 清單每條開成一題 `Q-<n>` 待答，再 `git rm discovery.md`")
    return None


def main(root: Path) -> int:
    specs = root / "specs"
    if not specs.is_dir():
        print(f"找不到 {specs}")
        return 1

    ongoing = specs / "ongoing"
    hint = flat_layout_hint(specs) or (legacy_hint(ongoing) if ongoing.is_dir() else None)
    if hint:
        print(hint)
        return 1

    specs_files = sorted(ongoing.glob("*/spec.md"))
    if not specs_files:
        # 全部歸檔後 git 不留空的 ongoing/：那是做完了，不是還沒開始
        archived = [p for p in (specs / "sleeping").glob("*") if p.is_dir()]
        if archived:
            print(f"沒有進行中的 feature（specs/sleeping/ 有 {len(archived)} 個已歸檔）")
            return 0
        print(f"{ongoing} 底下沒有任何 spec.md —— 先跑 discovery")
        return 1

    # 逐份先解析、逐份先報壞——三種壞法訊息分開講，因為修法不同：
    # 一節缺席是檔案格式跑掉，一節存在但零列是 DISCOVERY 沒問過問題，欄數
    # 不對的列是那一列本身壞了。任何一份壞掉都不能把其他份藏起來，或把
    # 同一份裡還解析得出來的列也一起蓋掉——那正是「看起來正常」的來源。
    parsed = {}
    broken = False
    for spec_path in specs_files:
        found, rows, unparseable = parse_open_questions(
            spec_path.read_text(encoding="utf-8"))
        parsed[spec_path] = rows
        if not found:
            print(f"{spec_path} 找不到 `## Open Questions` 一節 —— 檔案格式損壞")
            broken = True
        elif not rows and not unparseable:
            print(f"{spec_path} 的 `## Open Questions` 表零列 —— "
                  f"DISCOVERY 還沒問過任何問題，不是需求沒有疑點")
            broken = True
        if unparseable:
            print(f"{spec_path} 有 {len(unparseable)} 列解析不出欄位（欄數不是 5）：")
            for line in unparseable:
                print(f"  {line}")
            broken = True

    print(f"{_field('feature', 42)}{'已答':>6}{'n/a':>6}{'待答':>6}  追問覆蓋")
    print(f"{'':42}{'':>18}  {BIZ_ABBR}{TECH_ABBR}")
    tot = [0, 0, 0]
    blocked = []

    for spec_path in specs_files:
        slug = spec_path.parent.name
        rows = parsed[spec_path]
        answered = na = pending = 0
        dims: set[str] = set()
        pending_qs = []
        for row in rows:
            # `—` 是開場的範圍題，不屬於任何探針面向，不計入追問覆蓋
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
                # 有答案可追，跟真正卡住的問題混在一起算，會讓一份
                # spec.md 永遠有消不掉的紅卡、看起來比實際更沒準備好。
                # 獨立開一欄，不計入待答。
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

    print("\n追問覆蓋 · 業務面向（前沿檢查表）"
          "空與零/邊界/重複/時序/權限/失敗/時間/規模/降級/時限/可觀測："
          "✓ 問過（含 n/a）· — 還沒問")
    print("追問覆蓋 · 技術面向（前沿檢查表）"
          "seam/模組邊界/介面與型別契約/排序契約／決定性/既有資產／測試慣例："
          "✓ 問過（含 n/a）· — 還沒問")
    return 1 if broken else 0


if __name__ == "__main__":
    sys.exit(main(Path(sys.argv[1] if len(sys.argv) > 1 else ".").resolve()))
