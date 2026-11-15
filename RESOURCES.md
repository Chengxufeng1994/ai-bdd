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
| [Waterball-Software-Academy/aixbdd](https://github.com/Waterball-Software-Academy/aixbdd) | [`bdd`](https://github.com/Waterball-Software-Academy/aixbdd/tree/main/.agents/skills/bdd) | `skills/bdd`、`skills/implement` | red／green／refactor 三個入口、一次一個 slice、最窄的重跑入口、規格有洞就回交；委派方式照它的 `implement`：模型不自動觸發 `bdd`，只委派 green 與 refactor，測試層先紅由呼叫者自己做。本專案沒有 `dsl.md`，上游改成 `.feature` 與 `spec.md` 的技術部分；它按層切的 Phase 3（先寫完所有步驟再寫產品碼）沒有採用，票仍是垂直切片 |
| 同上 | `spec-by-example` | `skills/spec-by-example` | 業務語意的驗收旅程、`# [need clarification]` 標記。不做它的 UI 雛形與前後端拆分 |
| [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `planning-and-task-breakdown`、`incremental-implementation` | `skills/task-breakdown`、`skills/implement` | 一票一 commit、合併後重跑閘門；票檔的大小（XS–XL，XL 要拆）、驗證命令、檢查點；`plan.md` 的 `## Overview` 與風險的影響程度，高風險的票排前面。風險的 Mitigation 欄不採用（`plan-format.md` 不把解法寫死） |
| 同上 | [`spec-driven-development`](https://github.com/addyosmani/agent-skills/blob/main/skills/spec-driven-development/SKILL.md) | `skills/spec`（`references/tech-format.md`） | 一個 feature 一份自成一體的 spec；技術部分的 Tech Stack、Commands、Project Structure、Code Style、Boundaries 五節。不採用它的 Phase 0 能力地圖、`SPEC.md` 放在專案根的做法 |
| [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec) | `openspec/changes/`、`openspec/specs/` 與 archive | `skills/archive`、`specs/ongoing/` 與 `specs/sleeping/` | 進行中的變更與現況分開放，做完就把變更併進現況、移走變更目錄。改成：現況放專案自己的 `docs/`（只動 ai-bdd 區塊）、做完的目錄整個搬進 `specs/sleeping/` 而且可以搬回來修正，不另寫 delta |
| [sirius-cc-wu/sirius-skills](https://github.com/sirius-cc-wu/sirius-skills) | [`walkthrough-me`](https://github.com/sirius-cc-wu/sirius-skills/blob/main/skills/walkthrough-me/SKILL.md) | `skills/skill-rules`（S15）與每個 skill 的 `## Verification` | 交回前逐項確認的 `- [ ]` 清單，同時驗產物與「不該發生的事沒有發生」 |
| [mattpocock/skills](https://github.com/mattpocock/skills) | `grilling` | `skills/discovery`、`skills/clarify` | 一題一題問到底的問法。改掉每題附建議答案的做法——它會把人推向接受 |
| 同上 | [`to-spec`](https://github.com/mattpocock/skills/blob/main/skills/engineering/to-spec/SKILL.md) | `skills/spec`（第 6 步、`references/tech-format.md`） | Implementation Decisions、Testing Decisions 兩節與「不寫檔案路徑」的規則；seam 的四條規則（既有優先、取最高、理想是一個、新 seam 提在最高點）與攤給使用者確認。不發佈到 issue tracker |
| 同上 | `to-tickets` | `skills/task-breakdown` | tracer bullet 垂直切片、blocking edges、依層切票的失敗案例 |
| [obra/superpowers](https://github.com/obra/superpowers) | `brainstorming`、`writing-plans`、`executing-plans` | 開發本專案的流程；`docs/superpowers/` | 設計文件先批再做、計畫再執行。DISCOVERY「像一場 brainstorming」的定位也來自這裡 |
| [anthropics/skills](https://github.com/anthropics/skills) | `skill-creator` | `skills/skill-rules` | SKILL.md 500 行的門檻、漸進式揭露。`skill-rules` 只管規範與稽核，eval 迴圈仍用 `skill-creator` |
