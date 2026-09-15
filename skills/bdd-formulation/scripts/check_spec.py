#!/usr/bin/env python3
"""稽核 .feature 與 spec.md 的一致性。

用法：
    python3 check_spec.py [專案根目錄]        # 預設當前目錄

檢查十三件事，全部是機械性的：

  1. 覆蓋——雙向。spec.md 有例子而 feature 沒有的（漏做），
     以及 feature 指向 spec.md 裡不存在的例子（發明出來的驗收條件）。
  2. 缺口——漏掉的例子有沒有在 .feature 裡就地留註解交代。
  3. 狀態 tag——每個有對應 story 的 .feature 恰好一個。
  4. 方言陷阱——中文「規則:」不是關鍵字，會被解析成散文。
  5. 步驟樣板重用率——封閉文法有沒有真的套上。沒有及格線，只報數字：
     散文式實測 76 場景／208 樣板；封閉文法 7 場景／10 樣板。
     十三項裡只有這一項不會把退出碼變成 1。
  6. FR 完整性——spec.md 裡有沒有 FR 完全沒掛任何例子；沒有例子，
     SPEC 就無從寫出場景，只能發明。
  7. v2 殘留——`## User Stories` 底下還有 `#### FR-2` 這種標題形式的規則。
     對 v3 的解析器而言那是「不存在」不是「錯」，會安靜地掉出覆蓋率。
  8. story 標題形式——`### US-<n> · <slug>` 差一點的標題，對
     stories_in_spec() 而言也是「不存在」；文件有兩則以上 story 時，那則
     story 連同它的 `.feature` 會安靜地從覆蓋比對裡消失，跟第 7 項是同一
     種缺陷。層級寫錯（`#### US-1 · x`）與行內寫錯（少了 `·`）都算。
  9. 紅卡群組的形式——`#### Q` 與 `**Q-<n>**` 都是精確形式；偏掉的行會讓
     第 10、11 項一張紅卡都收不到，於是那兩項無法失敗。
 10. 懸空的紅卡指標——`#### Q` 指向 `## Open Questions` 表裡沒有的題號。
 11. 紅卡只列待答——`#### Q` 指向表裡狀態不是「待答」的題號。已答的答案
     已經長成某條 FR 或 AC，留著會讓地圖上的紅卡數量永遠不歸零。
 12. AC 編號對帳——`AC-<n>.<m>` 的 `<n>` 必須等於它掛在底下的那條 FR。
 13. 無主的 FR——有 FR 不屬於任何一則 story。沒有 story 認領它，就沒有任何
     `.feature` 被要求帶它的 `@example-` tag，那些 AC 於是安靜地不見了。
     第 8 項抓標題本身寫壞，這一項抓對應關係斷掉：FR 擺在第一個 story 標題
     之前，每一行都合規，照樣沒有人認領。

退出碼 0 = 全過，1 = 有問題。適合放進 CI。找不到 specs/、沒有 spec.md、
`## User Stories` 底下沒有任何 story 也各自退出 1——那些是「先跑上一步」，
不算在上面十三項一致性檢查裡。

找不到 .feature 不在此列。第 6-13 項只讀 spec.md，而 SPEC 交件到 FORMULATION
開跑之間正是它們最該跑的時候，所以那時只跳過需要 .feature 的第 1-5 項並在開頭
明講跳過了哪幾項；spec.md 乾淨就退出 0。

只讀不寫。找不到什麼就說找不到，不猜——但 specs/ 與 .feature 不對稱：
前者是停，後者是照說不誤，然後把跑得到的八項跑完。
FR／AC 的定版編號與 story 由哪些 FR 組成，都讀自
specs/<date>-<feature>/spec.md 同一個 `## User Stories` 段——FR 掛在它的
story 底下，一次遍歷就同時拿到兩者，不再是兩份文件、兩份清單。
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


def frs_in_spec(text: str) -> dict[str, set[str]]:
    """spec.md 的 `## User Stories` 段：FR 編號 -> AC 編號集合。

    v3 把 FR 從標題改成 `#### FR` 分組底下的粗體項目，AC 則是掛在它底下的
    清單項。理由是 AC 還要再往下掛 Given/When/Then 三行，用標題會走到 H6。

    FR 編號全域唯一，不是每則 story 各自從 1 起算——否則 AC-1.1 會指向
    兩個地方，`.feature` 的 @example-1.1 就失去意義。

    碰到任何標題就把 current 清掉。`#### FR`／`#### NFR`／`#### Q` 三個分組
    與 `### US-` 都是 FR 的祖先或兄弟，它們底下的項目不屬於上一條 FR。

    AC 只認契約宣告的 `- **AC-<n>.<m>**` 項目形式。用自由文字掃會把散文裡
    合法的交叉引用（「此規則與 AC-2.1 的計時起點相同」）記成這條 FR 的例子，
    然後對著 `.feature` 報一則假的「漏了 AC-2.1」——叫作者去修一個沒有錯的
    句子，而契約從來沒有禁止在散文裡引用 AC。
    """
    section = re.search(
        r"^## User Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return {}
    out: dict[str, set[str]] = {}
    current = None
    for line in section.group(1).splitlines():
        if line.startswith("#"):
            current = None
            continue
        fr = re.match(r"^\*\*FR-(\d+)\*\*", line)
        if fr:
            current = fr.group(1)
            out.setdefault(current, set())
            continue
        if current:
            ac = re.match(r"^\s*-\s*\*\*AC-(\d+\.\d+)\*\*", line)
            if ac:
                out[current].add(ac.group(1))
    return out


def v2_forms_in_spec(text: str) -> list[str]:
    """`## User Stories` 段裡殘留的 v2 標題形式。

    v3 把 FR／NFR 從標題改成粗體項目，所以 `#### FR-2` 這種寫法對
    `frs_in_spec()` 而言是「不存在」而不是「錯」——它會安靜地連同它的例子
    一起掉出覆蓋率計算，而整份文件印出「全部通過」。

    這個檢查是唯一能把那個缺席變回錯誤的東西。半遷移的文件才是暴露面：
    整份 v2 已經會在「沒有 `## User Stories` 一節」那裡大聲失敗。
    """
    section = re.search(
        r"^## User Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return []
    return [line.strip() for line in section.group(1).splitlines()
            if re.match(r"^#+\s*(FR|NFR|AC|EX)-\d", line)]


def malformed_story_headings(text: str) -> list[str]:
    """`## User Stories` 段裡不合規的 story 標題。

    slug 住在標題裡（`### US-<n> · <slug>`），而 stories_in_spec() 比對的是
    那個確切形式。差一點的標題不會報錯，只會讓那則 story 連同它的 `.feature`
    一起從覆蓋比對裡消失——文件有兩則以上 story 時，「一則都沒有」那道防線
    也不會觸發，於是整份印出「全部通過」。

    **層級也算差一點。** 只看 `### ` 開頭的話，`#### US-1 · <slug>` 兩個條件
    都不成立：它不是 `### ` 開頭，所以這裡看不到；它也不是 stories_in_spec()
    要的形式，所以那裡也收不到——那則 story 與它的 `.feature` 一起靜默消失。
    比照 v2_forms_in_spec() 用 `^#+` 匹配任何層級，再要求標題文字以
    `US-<數字>` 起始，就把它收回來了。

    `### ` 那一條仍然保留：這一節底下的 `### ` 只能是 story，寫成別的字也是
    一則收不到的 story。

    這個檢查把那個缺席變回錯誤。跟 v2_forms_in_spec() 是同一種東西：解析器
    找它要的，這個找它不要的，只有後者抓得到缺席。
    """
    section = re.search(
        r"^## User Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return []
    return [line.strip() for line in section.group(1).splitlines()
            if (line.startswith("### ") or re.match(r"^#+\s*US-\d", line))
            and not re.match(r"^### US-\d+ · \S+\s*$", line)]


def qs_in_stories(text: str) -> tuple[set[str], list[str]]:
    """`## User Stories` 段的 `#### Q` 分組：引用的 Q 編號，以及形式不對的行。

    story 底下的紅卡只是指標——狀態、面向、答案、決策史都住在
    `## Open Questions`。這裡只收編號，好跟那張表對帳。

    形式稍偏就回傳空集合是這個檢查最危險的失敗方式：`#### Q（待答）` 不是
    `#### Q`，`- Q-77  …` 不是 `- **Q-77**  …`，兩者都讓底下的紅卡一張都收
    不到，於是懸空指標檢查**無法失敗**。契約寫的是精確形式，所以偏掉的行要
    報錯，不是被容納。

    形式與編號一起回傳，是因為「哪幾行算在這個分組裡」只有這一份狀態機。
    拆成兩支函式就會有兩份，而它們遲早對分組的邊界有不同意見——那時形式
    檢查會安靜地掃錯範圍。
    """
    section = re.search(
        r"^## User Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return set(), []
    out: set[str] = set()
    malformed: list[str] = []
    in_q = False
    for line in section.group(1).splitlines():
        if line.startswith("#"):
            in_q = re.match(r"^#### Q\s*$", line) is not None
            if not in_q and line.startswith("#### Q"):
                malformed.append(line.strip())
            continue
        if in_q:
            found = re.findall(r"\*\*Q-(\d+)\*\*", line)
            if found:
                out.update(found)
            elif re.search(r"Q-\d", line):
                malformed.append(line.strip())
    return out, malformed


def qs_in_table(text: str) -> dict[str, str]:
    """`## Open Questions` 表：Q 編號 -> 「狀態」欄。

    只讀那一節：`## Document Overview` 底下的版本修訂歷史也是五欄表，
    對整份掃會把它的每一列都當成一題。

    狀態一起讀出來，是因為契約要求 `#### Q` 只列「待答」的題目——已答的
    答案已經長成某條 FR 或 AC，留在紅卡群組裡會讓地圖上的紅卡數量永遠
    不歸零。欄序固定 `| Q | 問題 | 面向 | 狀態 | 答案 |`，跟 `status.py`
    讀的是同一張表的同一欄。
    """
    section = re.search(
        r"^## Open Questions\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return {}
    out: dict[str, str] = {}
    for line in section.group(1).splitlines():
        line = line.strip()
        if not line.startswith("|"):
            continue
        cells = [c.strip() for c in line.strip("|").split("|")]
        m = re.fullmatch(r"Q-(\d+)", cells[0])
        if m:
            out[m.group(1)] = cells[3] if len(cells) > 3 else ""
    return out


def stories_in_spec(text: str) -> dict[str, set[str]]:
    """spec.md 的 `## User Stories` 段：story-slug -> 它涵蓋的 FR 編號集合。

    併檔之前 slug 住在另一個檔的 `## Stories` 一節，story 與 FR 的對應是
    人手維護的清單。併檔之後那個對應是**結構性的**——FR 就掛在它的 story
    底下，所以走一遍同時拿到 slug 與 FR，不需要第二份清單。

    slug 從 `### US-<n> · <slug>` 取，它也是那則 story 的 `.feature` 檔名。
    """
    section = re.search(
        r"^## User Stories\s*$(.*?)(?=^## |\Z)", text, re.M | re.S)
    if not section:
        return {}
    out: dict[str, set[str]] = {}
    current = None
    for line in section.group(1).splitlines():
        st = re.match(r"^### US-\d+ · (\S+)\s*$", line)
        if st:
            current = st.group(1)
            out.setdefault(current, set())
            continue
        if current:
            for fr in re.findall(r"^\*\*FR-(\d+)\*\*", line):
                out[current].add(fr)
    return out


def orphan_frs(text: str) -> list[str]:
    """`## User Stories` 段裡不屬於任何一則 story 的 FR 編號。

    覆蓋比對是逐 story 走的：一則 story 認領哪些 FR，就決定它的 `.feature`
    該帶哪些 `@example-` tag。沒有 story 認領的 FR 因此不會被任何檔案要求
    覆蓋——它的 AC 一條都不必寫，而腳本照樣印出「全部通過」。

    最常見的來源是位置：FR 擺在第一個 `### US-` 標題之前，stories_in_spec()
    的 current 還是 None，於是它兩邊都在（frs_in_spec 看得到它）卻不屬於誰。
    那不是「錯的 FR」而是「沒有人認領的 FR」，只有用減法找得到——正向解析
    永遠看不見自己沒收進來的東西。

    另一個來源是 story 標題寫壞，那由 malformed_story_headings() 指名；
    這裡會跟著報一次，訊息指的是不同的東西：那邊說標題錯了，這邊說哪幾條
    FR 因此沒有人覆蓋。
    """
    owned = {fr for frs in stories_in_spec(text).values() for fr in frs}
    return sorted(set(frs_in_spec(text)) - owned, key=int)


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
        print(f"找不到 {specs_dir} —— 沒有 SPEC 的產物可以比對")
        return 1
    # 十三項裡有八項只讀 spec.md（第 6-13 項）。SPEC 交件到 FORMULATION
    # 開跑之間還沒有任何 .feature，而那正是 spec.md 最需要被檢查的時刻——
    # 早退會讓那八項在整段空窗期一項都跑不到，壞掉的 spec.md 於是要等到
    # 下一個 skill 才被發現。所以缺 .feature 不再是「不能跑」，是「跑得到
    # 的先跑，跑不到的明講跳過」。
    spec_only = feat_dir is None

    problems = 0
    if spec_only:
        print(f"specs: {specs_dir}    feature: 找不到任何 .feature\n"
              f"只跑 spec.md 自己的八項檢查；覆蓋比對、狀態 tag、方言陷阱、"
              f"缺口註解、樣板重用率五項需要 .feature，跳過。")
    else:
        print(f"specs: {specs_dir}    feature: {feat_dir}\n")

    specs = sorted(specs_dir.glob("*/spec.md"))
    if not specs:
        print(f"{specs_dir} 底下沒有 spec.md —— 先跑 bdd-spec")
        return 1

    covered: set[str] = set()      # 有出現在某份 spec.md 裡的 story slug
    for spec_path in specs:
        stext = spec_path.read_text(encoding="utf-8")
        frs = frs_in_spec(stext)

        # 要先於 barren 檢查：v2 形式的標題會讓那條 FR 連同它的例子一起
        # 從 frs 消失，barren 於是對著一份殘缺的 dict 發出誤導訊息。
        stale = v2_forms_in_spec(stext)
        if stale:
            print(f"✗ {spec_path.parent.name} 的 `## User Stories` 底下有 v2 形式的"
                  f"標題，v3 要用粗體項目：{', '.join(stale)}")
            problems += 1

        malformed = malformed_story_headings(stext)
        if malformed:
            print(f"✗ {spec_path.parent.name} 的 `## User Stories` 底下有不合規的 "
                  f"story 標題，形式要是 `### US-<n> · <slug>`：{', '.join(malformed)}")
            problems += 1

        if not frs:
            print(f"✗ {spec_path.parent.name} 的 spec.md 沒有 `## User Stories` "
                  f"一節、或該節底下沒有任何 FR —— 先跑 bdd-spec")
            problems += 1
            continue

        # 完全沒有 AC 的 FR，SPEC 無從寫出場景，只能發明——而發明的驗收條件
        # 正是這支腳本要抓的另一個方向。FR 編號是每個 feature 各自從 1 編起，
        # 「FR-3」在多 feature 的 repo 裡定位不到是哪一份，所以訊息要帶上是
        # 哪個 feature 目錄。
        barren = sorted(fr for fr, exs in frs.items() if not exs)
        if barren:
            print(f"✗ {spec_path.parent.name} 這些 FR 沒有任何 AC，SPEC 無從寫出場景："
                  f"{', '.join('FR-' + b for b in barren)}")
            problems += 1

        # 形式不對的紅卡群組收不到任何編號，下面兩個檢查因此無法失敗——
        # 所以形式本身要先報錯，不能被容納成「這則 story 沒有紅卡」。
        q_ids, q_malformed = qs_in_stories(stext)
        q_table = qs_in_table(stext)
        if q_malformed:
            print(f"✗ {spec_path.parent.name} 的 `#### Q` 形式不對，紅卡收不到："
                  f"{', '.join(q_malformed)}")
            problems += 1

        # story 底下的紅卡是指標，指到 `## Open Questions` 沒有的題號就是
        # 斷掉的引用——通常是題目被刪了卻沒回頭清 story，或編號打錯。
        # v2 的結構做不到這個檢查（紅卡不掛在 story 上），v3 才有。
        dangling = sorted(q_ids - set(q_table), key=int)
        if dangling:
            print(f"✗ {spec_path.parent.name} 的 story 引用了 `## Open Questions` "
                  f"裡沒有的問題：{', '.join('Q-' + d for d in dangling)}")
            problems += 1

        # 已答的不是紅卡了——它的答案已經長成某條 FR 或 AC。留在 `#### Q`
        # 會讓地圖上的紅卡數量永遠不歸零，而那正是就緒判定要數的東西。
        # 懸空的題號在上面報過了，這裡不重複報。
        answered = sorted((q for q in q_ids if q_table.get(q, "待答") != "待答"),
                          key=int)
        if answered:
            named = ", ".join(f"Q-{q}（{q_table[q] or '表缺狀態欄'}）" for q in answered)
            print(f"✗ {spec_path.parent.name} 的 `#### Q` 列了狀態不是「待答」的"
                  f"問題：{named}")
            problems += 1

        # AC 的 <n> 必須等於它掛在底下的那條 FR。巢狀讓這件事在寫的時候
        # 不容易弄錯，但打字錯不會被結構擋下來——`AC-9.1` 打在 `**FR-1**`
        # 底下會被記成 FR-1 的例子，然後在下面的覆蓋比對裡變成一則
        # 「@example-9.1 不見了」的訊息，指著 .feature 說它漏寫，而錯在 spec.md。
        mismatched = sorted(
            (fr, ac) for fr, acs in frs.items() for ac in acs
            if ac.split(".")[0] != fr)
        if mismatched:
            print(f"✗ {spec_path.parent.name} 這些 AC 的編號對不上它所屬的 FR："
                  f"{', '.join(f'AC-{ac} 掛在 FR-{fr} 底下' for fr, ac in mismatched)}")
            problems += 1

        stories = stories_in_spec(stext)
        if not stories:
            print(f"✗ {spec_path.parent.name} 的 spec.md `## User Stories` 一節底下沒有任何 story "
                  f"—— SPEC 還沒切 story")
            problems += 1
            continue

        # 這裡原本有一個「spec.md 點名了 PRD 裡沒有的 FR」的檢查。它比的是
        # **兩份文件互相點名**，併檔之後那個方向不存在了，所以刪掉它沒有損失
        # 保護。但「story 與 FR 的對應可能不一致」是另一回事，併檔並沒有消掉
        # 它：同一次遍歷讀到的兩份資料仍然各有各的錨點（`### US-` 與
        # `**FR-**`），FR 擺在第一個 story 標題之前就誰也不屬於。那個方向由
        # 下面這個檢查接手。
        unowned = orphan_frs(stext)
        if unowned:
            print(f"✗ {spec_path.parent.name} 這些 FR 不屬於任何一則 story，"
                  f"沒有 .feature 會被要求涵蓋它們的 AC："
                  f"{', '.join('FR-' + f for f in unowned)}")
            problems += 1

        # slug 記的是「spec.md 提到哪些 story」，跟 .feature 在不在無關，
        # 所以留在守衛外面——孤兒 .feature 那一段要用它。
        for slug in sorted(stories):
            covered.add(slug)

        # 逐 story 比，不可把所有 story 的例子聯集起來跟單一 .feature 比：
        # 一則 story 一個 .feature，聯集會讓 A 的例子出現在 B 檔裡也算通過。
        if not spec_only:
            for slug, fr_ids in sorted(stories.items()):
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
                    issues.append(f"指向 spec.md 裡不存在的例子 {invented}")
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

    if not spec_only:
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

    if not spec_only:
        orphans = [p.name for p in feat_dir.glob("*.feature")
                   if p.stem not in covered]
        if orphans:
            print(f"\n沒有對應 story 的 .feature（不在本檢查範圍）：{orphans}")

    print(f"\n{'全部通過' if not problems else f'{problems} 個問題'}")
    return 1 if problems else 0


if __name__ == "__main__":
    sys.exit(check(Path(sys.argv[1] if len(sys.argv) > 1 else ".").resolve()))
