# SDD 實踐參考

這份文件記錄 **SDD（Spec-Driven Development）這套方法本身**是什麼，轉述自
Spec Kit、Kiro、Tessl 的官方文件，以及 Martin Fowler 對三者的比較。

**它不談這個專案。** 不提六步流程、不提任何 skill、不提我們同不同意——那些是
[ai-sdlc.md](./ai-sdlc.md) 的事。這份文件改動的理由只有兩個：來源改了定義，
或者我們讀錯了。

它跟 [bdd.md](./bdd.md) 是同一類文件：外部方法的轉述。要並排讀這兩份時請注意
一件事——**兩套方法都在講「規格」，但指的不是同一種東西**。BDD 的規格是可執行
的例子，SDD 的規格是給 agent 讀的散文。

---

## 一句話定義

> "writing a 'spec' before writing code with AI ('documentation first').
> The spec becomes the source of truth for the human and the AI."

Fowler 對「spec」本身的界定：

> "a structured, behavior-oriented artifact … written in natural language that
> expresses software functionality and serves as guidance to AI coding agents."

出處：[Understanding Spec-Driven-Development: Kiro, spec-kit, and
Tessl](https://martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html)。

Spec Kit 的官方說法更強硬，直接點出這套方法的主張：

> "Spec-Driven Development (SDD) inverts this power structure. Specifications
> don't serve code—code serves specifications."
>
> "Specifications must be precise, complete, and unambiguous enough to generate
> working systems."

出處：[github/spec-kit — spec-driven.md](https://github.com/github/spec-kit/blob/main/spec-driven.md)。

注意第二句是**條件句，不是現況描述**。整套方法的效力取決於散文規格能不能做到
precise、complete、unambiguous，而這正是下面幾節反覆出現的爭點。

---

## 一個詞，三種主張

「SDD」底下其實有三種強度不同的立場，差別在**規格在實作完成之後的命運**。
這個分類出自 Fowler：

| 層級 | 規格的命運 | 人改什麼 | 代表 |
| --- | --- | --- | --- |
| **spec-first** | 實作完就丟 | 改程式碼 | Kiro |
| **spec-anchored** | 保留，隨功能演進 | 改兩邊 | Spec Kit（志向） |
| **spec-as-source** | 唯一維護的產物 | **只改規格** | Tessl |

三者的成本與風險完全不同，**混用同一個詞是這個領域最常見的誤讀來源**。
spec-first 只是把需求想清楚再動工；spec-as-source 則要求放棄直接編輯程式碼。

Fowler 觀察到 Spec Kit 雖有 spec-anchored 的志向（每個 spec 開一個 branch），
實際行為仍偏向 spec-first。

---

## 三家工具

### Kiro（AWS）

每個功能是一組 spec，三個檔案：

| 檔案 | 內容 |
| --- | --- |
| `requirements.md` | user story 與驗收條件，用 EARS 語法寫 |
| `design.md` | 架構、sequence diagram、實作考量 |
| `tasks.md` | 拆好的實作任務，可勾選追蹤 |

兩條路徑：Requirements-First（需求 → 設計 → 任務）或 Design-First（設計 →
需求 → 任務）。**階段之間有人工核可 gate**；`Quick Spec` 會一口氣跑完三個階段
而不停下來確認。另有 steering 檔（`product.md`、`tech.md`、`structure.md`）
充當跨 session 的記憶。

出處：[kiro.dev/docs/specs/feature-specs](https://kiro.dev/docs/specs/feature-specs/)。

### Spec Kit（GitHub）

不是 IDE，是一組灌進既有 coding agent 的 slash command：

| 命令 | 做什麼 |
| --- | --- |
| `/speckit.constitution` | 建立專案的治理原則與開發準則 |
| `/speckit.specify` | 定義要建什麼（需求與 user story） |
| `/speckit.clarify` | 追問規格不足之處（選配，建議在 plan 之前） |
| `/speckit.plan` | 產出技術實作計畫與技術選型 |
| `/speckit.tasks` | 產出可執行的任務清單 |
| `/speckit.analyze` | 跨產物的一致性與覆蓋檢查 |
| `/speckit.implement` | 執行所有任務 |
| `/speckit.checklist` | 產出驗證需求用的品質檢查表 |
| `/speckit.converge` | 拿現況比對 spec/plan/tasks，補上未完成的部分 |

**constitution 是它獨有的東西**：一份凌駕於個別功能之上的原則檔，官方描述是
"Create your project's governing principles and development guidelines that
will guide all subsequent development."

出處：[github/spec-kit](https://github.com/github/spec-kit)。

### Tessl

三家裡唯一明確瞄準 spec-as-source 的：一份 spec 對應一個程式碼檔案，產生的
程式碼帶生成標記。它支援從既有程式碼反推 spec，並且**會刪掉所有宣告為 target
的檔案再從 spec 重建**——用重建結果來驗證規格是否真的完整。

Fowler 的評語是它同時最激進也最未經驗證，截至 2026 年中「still mostly a
thesis」。

出處：[docs.tessl.io](https://docs.tessl.io/use/spec-driven-development-with-tessl)。

---

## EARS

Kiro 的 `requirements.md` 用 EARS（Easy Approach to Requirements Syntax）寫。
它由 Alistair Mavin 與 Rolls-Royce 的同事在分析噴射引擎控制系統的適航法規時
發展出來，2009 年發表——**比 AI coding 早了十幾年**，目的是壓掉自然語言的歧義。

六個 pattern：

| 類型 | 樣板 | 關鍵字 |
| --- | --- | --- |
| Ubiquitous | `The <system> shall <response>` | 無 |
| State-driven | `While <precondition>, the <system> shall <response>` | While |
| Event-driven | `When <trigger>, the <system> shall <response>` | When |
| Optional feature | `Where <feature is included>, the <system> shall <response>` | Where |
| Unwanted behaviour | `If <trigger>, then the <system> shall <response>` | If / Then |
| Complex | `While <precondition>, When <trigger>, the <system> shall <response>` | While + When |

出處：[alistairmavin.com/ears](https://alistairmavin.com/ears/)。

EARS 是**有結構的自然語言，不是可執行的格式**。它約束的是句型，不是語意——
一句合法的 EARS 仍然可以指涉不存在的行為，而且沒有任何東西會告訴你。

---

## SDD 沒有說的

這一節劃出方法的邊界。以下幾件事在三家的流程裡**找不到指引**：

| 沒說的 | 說明 |
| --- | --- |
| **共識從哪裡來** | 三家的流程都是「agent 產 markdown、人審 markdown」。沒有任何一家規定多方對話，也沒有對應 BDD Three Amigos 的角色 |
| **規格怎麼被驗證** | 規格是散文，不可執行。確認它成立的手段只有人 review 與 agent 自述。唯一的例外是 Tessl 的「刪掉重建」，而那是最未經驗證的一支 |
| **規格漂移怎麼辦** | spec-first 直接放棄（丟掉規格）；spec-anchored 保留規格，但沒有機制保證它跟程式碼一致 |
| **規模怎麼縮放** | Fowler 指出 Kiro 與 spec-kit 都沒有處理問題大小的差異：小 bug 走完整套流程過度繁瑣，中型功能則造成 review overload |

Fowler 另外記下三項實跑觀察：

- **markdown 負擔**——"I'd rather review code than all these markdown files."
- **控制的錯覺**——即使規格與指令都寫得很細，agent 仍常忽略指引，或過度套用
  規則而產生重複與非預期的實作
- **歷史類比**——他把 spec-as-source 對照 Model-Driven Development 的失敗，
  警告它可能同時繼承「inflexibility and non-determinism」

最後一項值得單獨記住：**MDD 的模型是可編譯的，SDD 的規格是散文。** 把不確定性
放進生成環節，不會因為輸入變成自然語言就消失。

---

## 跟 BDD 的關係

這裡比的是兩套方法各自的文獻主張，不是本專案的取捨。

| | SDD | BDD |
| --- | --- | --- |
| **一份規格的範圍** | 整個 feature：需求＋技術計畫＋任務拆解 | 單一行為的一個例子 |
| **規格的形式** | 散文 markdown（Kiro 用 EARS 約束句型） | Given-When-Then，可執行 |
| **主要讀者** | AI coding agent | 人與機器都要讀得懂 |
| **共識機制** | 人 review agent 的產出 | Three Amigos 的對話（見 [bdd.md](./bdd.md) §①） |

兩者不互斥。常見的組合是 feature 層用 SDD、unit 層用 [TDD](./tdd.md)，並借 BDD 的
Given-When-Then 來寫人與 agent 都讀得懂的驗收條件。

也有一項針對 SDD 的批評直接指向這裡：**SDD 沒有解決 BDD 當年沒解決的那個協作
問題。** 由 agent 生成的情境反映的是 LLM 的訓練分佈，而不是這個領域真正的邊界
案例——換句話說，把追問的工作交出去，得到的是看起來完整的規格，不是問對問題
的規格。

---

需要這些沒被涵蓋的東西時，得自己補。本專案怎麼補的見 [ai-sdlc.md](./ai-sdlc.md)。
