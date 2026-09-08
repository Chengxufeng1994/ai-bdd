# Benchmark

固定輸入 ＋ 評分方案，用來量測 ai-bdd 的 skill 有沒有變好。

## 一個 case 是什麼

`cases/` 底下一個檔案就是一個 case。它同時是**題目**與**改考卷的標準**：

| 章節 | 給誰 |
| --- | --- |
| `## 輸入` | **只有這一段餵給 skill** |
| `## 為什麼是這樣措辭` | 人 |
| `## 模糊點` | 人 |
| `## 及格標準` | 人 |
| `## 怎麼比較兩次執行` | 人 |
| `## 實測基準` | 人；每跑一次累積一次 |

**餵食線就是 `## 輸入`。** 評分方案必須進版控——不然沒有東西可以比較——但它同時
就是答案：「模糊點」那一節列的正是你期待 skill 自己問出來的東西。所以隔離靠章節，
不靠目錄。**把 case 檔整份貼進 context 等於直接給答案。**

## 為什麼 case 與 skeleton 分開

`cases/` 語言無關，`skeleton/` 每個語言一份。這條軸不是隨便切的：**IMPLEMENT
以上游的一切都是共用的，只有 harness 隨語言不同**。CLARIFY、SPEC、PLAN 的產物
（規則、例子、`.feature`）換個語言仍然成立；`Makefile`、分層規則、測試 runner
不行。

## 怎麼跑一次

```bash
mkdir -p runs
cp -r benchmark/skeleton/go runs/go-fitness-tracker
```

`runs/` 是 gitignored，所以新 clone 上不存在——`mkdir -p` 不能省。

然後把該 case 的 `## 輸入` 那一段——只有那一段——餵給 `bdd-clarify`。

一次跑會產出：

```
runs/<lang>-<case>/
├── specs/        那次的 example map 與問題檔
└── features/     它產出的 .feature
```

## 兩條規則

**快照。** `benchmark/skeleton/<lang>/` 永遠凍結在「只有走路骨架」的狀態。一次跑的
複本可以長出任何東西，但**沒有任何東西流回範本**。要改進範本就直接改
`benchmark/skeleton/`，之後複製的跑會拿到，正在跑的不會——harness 若在跑到一半被
換掉，比較就失效了。

**不進版控。** `runs/` 是 gitignored。一次 eval 的答案若躺在 git 裡，下一次跑就能
**讀**到而不是**推導**出來，而推導出來的答案和抄來的答案，在產物裡長得一模一樣。
