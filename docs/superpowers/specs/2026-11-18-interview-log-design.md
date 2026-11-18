# `discovery` 不寫規格，只留一份訪談紀錄

**狀態**：已批　**日期**：2026-11-18　**修正**：`2026-11-18-drop-discovery-design.md` 的 D1、D3

## 意圖

需求方的話：「discovery 好像還是很奇怪，他應該是單純的問問題，然後再把東西丟給 spec，產出 spec.md，
discovery 不應該有任何產物？」「如果怕忘記可以參考 aixbdd clarify 的作法」「沒錯，我覺得這樣會更好」。

上一版讓 `discovery` 直接建 `spec.md`、寫它的 `## 原文` 與 `## Open Questions`——一份檔兩個寫入者，
`discovery` 還是在寫規格。aixbdd 的 `clarify` 不寫規格；呼叫者要求詳記時，它寫一份**問答紀錄**
（`templates/clarify-log.md`：Session、每輪每題的總結之提問／使用者選擇／補充／收斂結果、Confirmed
Decisions、Remaining Open Items），規格由呼叫者 `specify` 寫。

## 決定

### D1. `discovery` 只寫 `interview-log.md`

路徑 `specs/ongoing/<date>-<feature>/interview-log.md`。它是**流水帳**，不是規格：

| 節 | 寫什麼 | 怎麼改 |
| --- | --- | --- |
| `## 原文` | 需求方帶來的話，逐字，§ 編號 | 建檔時寫一次，之後不改 |
| `## Session` | 主題、需求方、核心關係人與決策權、代演與否、粗切 story（`US-<n> · <slug>` ＋ 身為一句） | 只追加 |
| `## 問答` | 一題一節 `### Q-<n> · <問題>`：面向、上下文、猜測、答案原話、誰答的、收斂結果、日期 | 只追加；唯一可改的是那一題的 `狀態` 從 `待答` 到 `已答`／`n/a`／`失效`，答案一律追加不覆寫 |
| `## Confirmed Decisions` | 意圖陳述（Outcome／User／Why now／Success／Constraint／Out of scope）與已拍板的決定 | 每次結束時重寫 |
| `## Remaining Open Items` | 還待答的題號與它掛在誰身上；`結束方式` | 每次結束時重寫 |

`discovery` 不寫 `spec.md`、不寫 FR／AC。`status.py` 讀 `interview-log.md` 算待答與面向覆蓋。

### D2. `spec` 是 `spec.md` 唯一的寫入者

`spec` 讀 `interview-log.md`，建 `spec.md`：`## 原文` 逐字照抄、`## Open Questions` 由 log 的每一題整理成表與
小節（題號 `Q-<n>` 不變）、FR／AC 照已答的題寫。`spec.md` 自成一體：`← Q-<n>` 指向它自己的
`## Open Questions`，`check_spec.py` 只讀 `spec.md`（不變）。

### D3. 同一個答案兩處

log 是**逐字的原始紀錄**，跟 `## 原文` 一樣不改；`spec.md` 的表是**整理過、規則引用的**。兩處對不上時以
log 的原話為準，`spec` 修補時照 log 改。這是刻意的例外，不是同一個事實兩份：一份是證據，一份是用證據
寫成的規格。

### D4. 洞的路線

```
spec 寫規格時發現洞 → spec.md 的 ## Open Questions 開一題 Q-<n> 待答、狀態改 澄清中 → STOP，告訴使用者 /discovery
discovery 讀 spec.md 待答的題（只讀）→ 在 log 追加成同號的題 → 問 → 記在 log
spec 再跑 → 讀 log 的新答案 → 修補 spec.md
```

`discovery` 讀 `spec.md` 只為了知道要補問哪幾題，不寫它。

### D5. 不變

一次一題、附猜測、開場假設與信心、95% 停、意圖陳述要明說 yes、十六個面向、前沿、能自己查的不准問、
稽核分類是否窮盡、代演、喚醒、三堆（可解／阻塞／失效）。

## 驗證

- `audit_skill.py` 全過；`claude plugin validate .` 通過；`python3 -m unittest discover -s tests` 通過
  （`status.py` 的測試改讀 `interview-log.md`）
- `discovery` 的 SKILL.md 與 references 裡沒有任何「寫 `spec.md`」的指示
- 範例：`discovery/examples/` 有一份 `interview-log.md`，`spec/examples/minimal-spec.md` 照它寫成，兩份的題號與原話對得上
