# 每個 skill 一份 `## Verification`

**狀態**：已批　**日期**：2026-11-15

## 意圖

每個 SKILL.md 有一份交回前逐項確認的 `- [ ]` 清單，讓執行 skill 的 agent 在交回前確認
「做完了、做對了」。參考 `sirius-cc-wu/sirius-skills` 的 `walkthrough-me`：清單放在檔尾、
每項是非題，而且同時驗產物與「不該發生的事沒有發生」。

### 為什麼

流程章節寫怎麼做，`## 完成後` 寫要告訴對方什麼，兩者都沒寫怎麼確認做對了。本專案 skill
最危險的錯誤（順手補答案、自己翻 `已對焦`、少跑 `check_spec.py`）都會產出一份看起來
完整的產物。

## 決定

- 需求方選 A：清單給**執行 skill 的 agent** 檢查產物，不是給維護者檢查 skill（B）
- `skill-rules` 新增 S15（`MUST`）：有 `## Verification`，底下至少一個 `- [ ]`；寫法規則在
  `rules/structure.md` 的〈S15：Verification〉——可判是或否、只驗本文已有的規則、有腳本就驗
  腳本、至少一項驗邊界、不重述流程、4–10 項
- 位置：輸出格式之後、`## 完成後` 之前
- `audit_skill.py` 機械檢查 S15（`tests/test_audit_skill.py`）
- 十個 skill 全部補上，項目只從各自既有的 MUST／NEVER／IMPORTANT 與腳本來

## 非目標

- 不做 B（每個 skill 的情境測試清單）、不接 `skill-creator` 的 eval 迴圈
- 不在清單裡新增任何本文沒有的規則

## 考慮過的做法

| 做法 | 結果 |
| --- | --- |
| A. 執行者自檢產物，放 SKILL.md | **採用** |
| B. 維護者用的情境清單，放子目錄 | 捨棄：需求方選 A |
| 兩者都做 | 捨棄：兩份清單大量重疊，會漂移 |
