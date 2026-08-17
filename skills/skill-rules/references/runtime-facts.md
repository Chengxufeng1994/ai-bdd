# 執行時行為

## 證據等級

本文件每一項都標註依據強度。**不同等級不可混用**——把慣例當實測引用，是這份
規範最容易犯的錯。

| 標記 | 意義 |
| --- | --- |
| 【實測】 | 在本機驗證過，附重現方式 |
| 【文件】 | 官方工具輸出或說明文件明載 |
| 【推論】 | 從前兩者推導，推導過程寫出來供檢驗 |

---

## 1.【實測】skill 名稱在 plugin 內扁平，跨 plugin 可同名

`skills/<分類>/<名稱>/SKILL.md` 這種巢狀結構在執行時**有效**，skill 會正常載入。
`claude plugin details` 列出的是**名字**，目錄層級不出現。

驗證：`mattpocock-skills` 的 25 個 skill 全部住在 `skills/engineering/`、
`skills/productivity/` 等分類目錄下，inventory 卻是扁平列名。

**但名稱不是全域唯一的。** 跨 plugin 同名可以共存，以 plugin 名前綴區分：

```
$ find ~/.claude/plugins/cache -name SKILL.md | ... | sort | uniq -c
     44 verification-loop
     44 tdd-workflow
     44 security-review
     44 coding-standards
```

同時存在的例子：`code-review`、`mattpocock-skills:code-review`、
`pr-review-toolkit:code-reviewer`。

【推論】分組若不寫進名字，在同一個 plugin 內就完全看不見。同組 skill 用共同
前綴（`bdd-clarify`、`bdd-spec`），字母排序會把它們聚在一起——那是扁平命名空間
裡唯一還能保住分組的機制。跨 plugin 的唯一性由 plugin 名負責，不需要自己處理。

> 修訂記錄：本節初版寫成「名稱在全域唯一」。那是從「inventory 顯示扁平名字」
> 直接跳到結論，中間漏了「跨 plugin 同名會怎樣」這一步——而那一步推翻了結論。

## 2.【實測】`claude plugin validate` 不遞迴

巢狀目錄的 skill **不會被驗證**。

驗證方式（對照實驗。validate 沒問題時不出聲，所以靜默證明不了任何事）：

```
把巢狀的 SKILL.md 的 description 拿掉 → ✔ 通過（完全沒檢查）
把扁平的 SKILL.md 的 description 拿掉 → ✘ 失敗並指出問題
```

【推論】巢狀的 skill 能跑，但等於永遠繞過 CI——問題要等到某個 skill 靜默失效
才會浮現。

**這項證據支持的是「巢狀不會被檢查」，不是「巢狀會壞」。** mattpocock-skills
長期使用巢狀結構且運作良好。所以對應的規則 S11 是 `SHOULD` 而非 `MUST`：這是
一個可驗證性的取捨，不是正確性的紅線。

## 3.【文件】description 是常駐成本，本文不是

`claude plugin details <plugin>` 會列出成本。mattpocock-skills（25 個 skill）：

```
Always-on: ~1,620 tok   每個 session 都付
每個 skill：常駐 ~30–160 tok（name ＋ description）
            觸發時 ~20–3.9k tok（本文）
```

IMPORTANT: 工具自己在輸出末尾標註「Token counts are estimates and may differ
from actual usage」。這些是**估計值，不是量測值**，用來比較數量級可以，拿來做
精算不行。

【推論】`description` 每個 session 都付費，本文只在觸發時付。所以**相對於本文，
description 的每字成本高得多**——厚重內容該放本文或子目錄，`description` 壓成
精準的觸發條件。

這個推論**不**支持「拆得細不貴」。常駐成本隨 skill 數線性成長：25 個約 1.6k，
100 個約 6.5k。它支持的只是「同樣的內容，放本文比放 description 便宜」。

## 4. `disable-model-invocation: true` ＝ 使用者專用

【文件】`claude-code-setup` 的 `claude-automation-recommender` 內註明：
`disable-model-invocation: true` — User-only (for side effects: deploy, commit, send)。

【實測】相關性觀察：`mattpocock-skills` 中設了此欄位的 `ask-matt`、`implement`、
`to-spec`、`wayfinder` 全部不在模型可見的 skill 清單中；未設的 `tdd`、
`diagnosing-bugs`、`domain-modeling` 則都在。

【推論】兩項證據方向一致，但第二項只是相關性——不在清單也可能有別的原因。
結論的信心來自文件那一項，觀察只是佐證。

實務含意：若一個 skill「明明存在卻從來不觸發」，先檢查這個欄位再檢查 description。

## 5.【推論】子目錄檔案不會自動載入

只有 SKILL.md 的本文在 skill 觸發時進 context。`rules/`、`references/`、
`examples/`、`scripts/` 裡的檔案必須由 SKILL.md 明確指示去讀。

依據：第 3 項的成本模型只列出「常駐」與「觸發時」兩層，子目錄不在計費項目中；
且 skill-creator 的漸進式揭露說明明確描述為三層載入，第三層是「as needed」。

尚未直接實測。要驗證的話：放一個只在子目錄檔案中出現的獨特字串，觸發 skill
後看模型是否知道它。

【推論】這是最常見的疏漏——放了一堆參考檔卻從沒在主檔裡指路，等於那些檔案不
存在，維護者卻以為內容已經涵蓋。稽核時要特別查（規則 S5）。

---

## 未經驗證的部分

本文件涵蓋的是**執行時機制**。skill-rules 的其餘內容（章節結構、觸發詞寫法、
強度分級）屬於**慣例**，來源是參考教材與設計判斷，**沒有任何證據證明遵守它們
能產出更好的 skill**。

要讓那些規則可被證偽，需要 `skill-creator` 的 eval 迴圈：同一批測試提示，比較
有無套用規則的產出。在那之前，那些規則的地位是「有理由的約定」，不是「已驗證
的最佳實踐」。

## 重測方式

```bash
claude plugin details <plugin-name>     # 第 1、3 項
claude plugin validate <path> --strict  # 第 2 項（需搭配對照實驗）
find ~/.claude/plugins/cache -name SKILL.md | \
  sed -E 's|.*/([^/]+)/SKILL.md|\1|' | sort | uniq -d   # 第 1 項的同名檢查
```
