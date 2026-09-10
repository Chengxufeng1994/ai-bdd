# Benchmark

固定輸入 ＋ 評分方案，用來量測 ai-bdd 的 skill 有沒有變好。

## 一個 case 是什麼

`cases/` 底下一個**目錄**就是一個 case，裡面的檔案分成兩類，界線是實體的：

```
cases/parking-billing/
├── input-brief.md    ← 餵給 skill 的全部內容
└── grading.md        ← 答案卷：模糊點、及格標準、怎麼比較、實測基準
```

一個情境可以有多份輸入（`personal-memo/` 有 `input-brief.md` 與 `input-prd.md`
兩種形式），共用同一份 `grading.md`。

**輸入檔是純酬載——沒有標題、沒有說明、沒有圍欄。** 任何「這是評估用的固定輸入、
請逐字複製」之類的抬頭，本身就會告訴被測者它正在被評估，而那會改變它的行為。
操作說明放在 `grading.md` 與這份 README，不放在輸入裡。

**為什麼是兩個檔而不是一個檔的兩節。** 評分方案必須進版控（不然沒有東西可以比較），
但它同時**就是答案**——「模糊點」那一節列的正是你期待 skill 自己問出來的東西。
最初的設計是同檔靠章節隔離，實測第一次就失效：把 case 檔整份 @ 給一個 session，
章節界線擋不住任何東西，那次 eval 當場作廢。**靠章節的界線需要有人記得遵守；
靠檔案的界線不需要**——`input-brief.md` 可以安全地交給任何人，因為它裡面沒有答案。

## 為什麼 case 與 skeleton 分開

`cases/` 語言無關，`skeleton/` 每個語言一份。這條軸不是隨便切的：**IMPLEMENT
以上游的一切都是共用的，只有 harness 隨語言不同**。CLARIFY、SPEC、PLAN 的產物
（規則、例子、`.feature`）換個語言仍然成立；`Makefile`、分層規則、測試 runner
不行。

## 怎麼跑一次

```bash
mkdir -p runs
cp -r benchmark/skeleton/go runs/go-personal-memo
```

`runs/` 是 gitignored，所以新 clone 上不存在——`mkdir -p` 不能省。

然後把該 case 的輸入檔**整份**交給 `bdd-clarify`：

```bash
cat benchmark/cases/personal-memo/input-brief.md
```

那個檔案裡沒有答案，所以整份給是安全的——這正是它跟 `grading.md` 分開的理由。

**跑 eval 的那個 session 不可以讀到 `grading.md`。** 最保險的做法是在 repo 外的空
目錄跑；退而求其次，開場就明確禁止它讀 `benchmark/cases/`。禁令是散文，擋不住
工具，所以能用前者就用前者。

一次跑會產出：

```
runs/<lang>-<case>/
├── specs/<date>-<feature>/prd.md    CLARIFY 的產物
└── features/                        SPEC 產出的 .feature
```

## 兩條規則

**快照。** `benchmark/skeleton/<lang>/` 永遠凍結在「只有走路骨架」的狀態。一次跑的
複本可以長出任何東西，但**沒有任何東西流回範本**。要改進範本就直接改
`benchmark/skeleton/`，之後複製的跑會拿到，正在跑的不會——harness 若在跑到一半被
換掉，比較就失效了。

**不進版控。** `runs/` 是 gitignored。一次 eval 的答案若躺在 git 裡，下一次跑就能
**讀**到而不是**推導**出來，而推導出來的答案和抄來的答案，在產物裡長得一模一樣。
