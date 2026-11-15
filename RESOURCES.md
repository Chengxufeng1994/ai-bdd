# 參考來源

這個專案的 skill 與文件借鏡了哪些外部來源。分兩類：**方法論的原文**（`docs/` 轉述的對象）與
**其他 skill 專案**（個別 skill 參考的寫法）。

新增或改寫 skill、文件時用到新的來源，在這裡加一列——寫清楚是哪個 skill 或文件用了它、
取了什麼。只列真的讀過並用上的，不列「相關閱讀」。

## 方法論

| 來源 | 作者 | 用在哪 | 取了什麼 |
| --- | --- | --- | --- |
| [Behaviour-Driven Development](https://cucumber.io/docs/bdd/) | Cucumber | `docs/bdd.md` | Discovery／Formulation／Automation 三個實踐的定義 |
| [Introducing BDD](https://dannorth.net/introducing-bdd/) | Dan North | `docs/tdd.md` | BDD 怎麼從 TDD 長出來 |
| [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/) | Matt Wynne | `docs/example-mapping.md`、`skills/example-mapping` | 四色卡、桌面即回饋、25 分鐘的框、例子要粗 |
| [Acceptance Criteria vs Scenarios](http://lizkeogh.com/2011/06/20/acceptance-criteria-vs-scenarios/) | Liz Keogh | `docs/example-mapping.md` | 規則與例子的差別 |
| [Canon TDD](https://newsletter.kentbeck.com/p/canon-tdd) | Kent Beck | `docs/tdd.md` | TDD 的標準流程 |
| [The Cycles of TDD](https://blog.cleancoder.com/uncle-bob/2014/12/17/TheCyclesOfTDD.html) | Robert C. Martin | `docs/tdd.md` | 三條定律與各層循環 |
| [Test Driven Development](https://martinfowler.com/bliki/TestDrivenDevelopment.html) | Martin Fowler | `docs/tdd.md` | TDD 的定義 |
| [Mocks Aren't Stubs](https://martinfowler.com/articles/mocksArentStubs.html) | Martin Fowler | `docs/tdd.md` | 兩個學派（classicist／mockist） |
| [Is TDD Dead?](https://martinfowler.com/articles/is-tdd-dead/) | Fowler、Beck、DHH | `docs/tdd.md` | 對 TDD 的批評與回應 |
| [Spec Kit](https://github.com/github/spec-kit)、[spec-driven.md](https://github.com/github/spec-kit/blob/main/spec-driven.md) | GitHub | `docs/sdd.md` | Spec-Driven Development 的流程與產物（spec／plan／tasks） |
| [Kiro Feature Specs](https://kiro.dev/docs/specs/feature-specs/) | AWS Kiro | `docs/sdd.md` | SDD 的另一種實作 |
| [Spec-driven development with Tessl](https://docs.tessl.io/use/spec-driven-development-with-tessl) | Tessl | `docs/sdd.md` | SDD 的另一種實作 |
| [Understanding SDD: Kiro, spec-kit, and Tessl](https://martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html) | Birgitta Böckeler | `docs/sdd.md` | 三種工具的比較 |
| [EARS](https://alistairmavin.com/ears/) | Alistair Mavin | `docs/sdd.md` | 需求句型 |
| [The Humanizing Work Guide to Splitting User Stories](https://www.humanizingwork.com/the-humanizing-work-guide-to-splitting-user-stories/) | Humanizing Work | `skills/story-splitting` | 九種切分模式與挑選切法的兩條規則 |

## 其他 skill 專案

| 專案 | 參考的 skill | 用在哪 | 取了什麼、改了什麼 |
| --- | --- | --- | --- |
| [Waterball-Software-Academy/aixbdd](https://github.com/Waterball-Software-Academy/aixbdd) | [`bdd`](https://github.com/Waterball-Software-Academy/aixbdd/tree/main/.agents/skills/bdd) | `skills/bdd` | red／green／refactor 三個入口、一次一個 slice、最窄的重跑入口、規格有洞就回交。本專案沒有 `dsl.md`，上游改成 `.feature` 與 `plan.md` 的技術節；回交對象改成 `bdd-implement` 的票檔 |
| 同上 | `spec-by-example` | `skills/bdd-spec-by-example` | 業務語意的驗收旅程、`# [need clarification]` 標記。不做它的 UI 雛形與前後端拆分 |
| [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `planning-and-task-breakdown`、`incremental-implementation` | `skills/bdd-tasks`、`skills/bdd-implement` | 一票一 commit、合併後重跑閘門。對照後列出的其餘差距（票的大小上限、高風險先做）尚未採用 |
| [sirius-cc-wu/sirius-skills](https://github.com/sirius-cc-wu/sirius-skills) | [`walkthrough-me`](https://github.com/sirius-cc-wu/sirius-skills/blob/main/skills/walkthrough-me/SKILL.md) | `skills/skill-rules`（S15）與每個 skill 的 `## Verification` | 交回前逐項確認的 `- [ ]` 清單，同時驗產物與「不該發生的事沒有發生」 |
| [mattpocock/skills](https://github.com/mattpocock/skills) | `grilling` | `skills/bdd-discovery`、`skills/clarify` | 一題一題問到底的問法。改掉每題附建議答案的做法——它會把人推向接受 |
| 同上 | `to-tickets` | `skills/bdd-tasks` | tracer bullet 垂直切片、blocking edges、依層切票的失敗案例 |
| [obra/superpowers](https://github.com/obra/superpowers) | `brainstorming`、`writing-plans`、`executing-plans` | 開發本專案的流程；`docs/superpowers/` | 設計文件先批再做、計畫再執行。DISCOVERY「像一場 brainstorming」的定位也來自這裡 |
| [anthropics/skills](https://github.com/anthropics/skills) | `skill-creator` | `skills/skill-rules` | SKILL.md 500 行的門檻、漸進式揭露。`skill-rules` 只管規範與稽核，eval 迴圈仍用 `skill-creator` |
