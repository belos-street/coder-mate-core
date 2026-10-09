# coder-mate-core 项目评测与代码审查报告

- **评测日期**：2026-10-09
- **评测方法**：skill-kit `belos-street` — [project-evaluator](file:///D:/code/personal/skill-kit/skills/belos-street/reference/project-evaluator.md)（仓库级十维评测）+ [code-review-dimensions](file:///D:/code/personal/skill-kit/skills/belos-street/reference/code-review-dimensions.md)（八维走查与定级）
- **评测范围**：全仓库（Go 1.23.5，约 7.0k 行 Go 源码 / 21 个测试文件 / 17 个语言适配包），排除 `.agents/skills`（工具链引用文档）
- **取证方式**：两路并行子代理取证 → 全部可计数断言 Grep/Shell 逐条复核（复核中剔除 1 条误报：`stringifyInlineStyle` 尾分号实为合法 CSS，已从问题清单移除）→ `go test ./...` 22 包实测全绿

---

## 0. 红线熔断结论（阻断级）

**命中「静默吞错」红线**（命中即阻断，按约定不降级讨论）：

| 证据 | 内容 | 后果 |
| --- | --- | --- |
| [api.go#L55](file:///D:/code/personal/coder-mate-core/api.go#L55) | `currentTheme, _ := theme.ResolveTheme(...)` 丢弃 error，无日志无上抛 | 未知主题名时 theme 为 nil，沿 [highlighter.go#L128-L133](file:///D:/code/personal/coder-mate-core/highlighter.go#L128-L133) → [theme/registry.go#L95](file:///D:/code/personal/coder-mate-core/theme/registry.go#L95) `theme.Styles[scope]` **nil 解引用确定性崩溃**，且全库无测试覆盖此路径 |
| [language/builtins.go#L46](file:///D:/code/personal/coder-mate-core/language/builtins.go#L46) | `core.RegisterLanguage(lang)` 返回的 error 被整体忽略 | 注册冲突不可见，静默丢失 |

---

## 1. 总评分览表（10 分制）

| # | 维度 | 得分 | 一句话依据 |
| --- | --- | --- | --- |
| 1 | 架构与分层 | **7** | core / language / theme / wasm 四层清晰、零依赖、公共 API 克制（3 方法 + options）；`Theme interface{}` 弱类型、`NewHighlighter(opts ...HighlighterOptions)` 可变参反模式、入口层吞错拉低 |
| 2 | 代码可维护性 | **5** | 17 个语言包同骨架逐字复制（adapter/spec），共享规则正则 ≥10 处重复，命名三制并存（`GRAMMAR_RULES` / `GrammarRules` / `getGrammarRules()+sync.Once`）；无 god class、单文件小是亮点 |
| 3 | 语言/类型纪律 | **6** | Go 惯用法总体正确（error 返回、拷贝防篡改）；但 `string(code[0])` 字节级索引破坏 UTF-8（[tokenizer.go#L158](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L158)）、两处 error 盲丢、`insertPos := 3` 魔法数 |
| 4 | 安全 | **6** | 零第三方依赖无 CVE 面、无注入面；但用户可控的 `preStyle` / `lineClassPrefix` **未转义直注入 HTML 属性**（[highlighter.go#L210](file:///D:/code/personal/coder-mate-core/highlighter.go#L210)、[L226](file:///D:/code/personal/coder-mate-core/highlighter.go#L226)），缓存无上限 |
| 5 | 测试 | **6** | 双轴：覆盖数量好（272 个 `func Test` / 21 文件，每包有测试，10 Benchmark，22 包全绿）× 断言真实性两极（JS 包 33 个真断言，其余 16 包"存在即过"），零 fuzz，缺多字节/panic/并发 case |
| 6 | 依赖健康 | **9** | 纯标准库零第三方依赖，Go 1.23.5 非 EOL，无任何可腐化项 |
| 7 | 工程化与 CI | **3** | 全仓库无 CI / Makefile / lint 配置（Glob 零命中）；spec10 规划的 gofmt + lint、spec08 的 Go/TS token 一致性验证**零落地**——只立法（部分）不执法 |
| 8 | 文档与协作 | **7** | 10 篇 spec + 迁移计划 + 需求索引是真沉淀，README 含基准数据且自曝短板；但漂移明显：agents.md 仍称"待进入实现阶段"、引用不存在的 `ts/` 目录、Go 版本三方打架（agents.md `1.21+` / go.mod `1.23.5` / README `1.23+`） |
| 9 | 提交与版本纪律 | **6** | 11 个提交 Conventional Commits 形式全覆盖、中文描述具体；但存在三合一混合提交（`896d07e` 一条消息含 readme/基准/单测三件事），无 tag / release |
| 10 | 人力风险（巴士因子） | **2** | `git shortlog`：belos-street 10 + 酱汁 1（疑同一人双账号记账），**巴士因子 1**；typescript 对 JS 规则的顺序手术式代码无人敢动，且无 CI 兜底放大风险 |

> 总分口径：按 skill 约定只输出维度分，不加权合成总分。

---

## 2. 问题清单（按严重度定级）

| No. | 严重度 | 维度 | 问题 | 建议 | 位置 |
| --- | --- | --- | --- | --- | --- |
| 1 | **P1 必修** | 正确性 | **UTF-8 字节级处理**：回退分支按字节取字符 `string(code[0])`，列号按 `len()` 字节计算，CJK 等多字节字符被拆成非法碎片，中文代码渲染乱码 | 改 `utf8.DecodeRuneInString`，列号按 rune 计 | [tokenizer.go#L158](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L158)、[L133](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L133) |
| 2 | **P1 必修** | 健壮性 | **nil theme 确定性 panic 链**（红线吞错的下游后果，见第 0 节） | `newHighlighter` 收到 error 即返回 error；`ResolveScopeStyle` 加 nil 防护 | [api.go#L55](file:///D:/code/personal/coder-mate-core/api.go#L55)、[theme/registry.go#L95](file:///D:/code/personal/coder-mate-core/theme/registry.go#L95) |
| 3 | **P1 必修** | 安全 | **HTML 属性注入**：`lineClassPrefix` 直写 `class` 属性、`preStyle` 直接字符串拼接进 `style` 属性，均未转义，引号可破出属性 | 属性值统一走转义函数（至少转 `"` `<` `&`） | [highlighter.go#L210](file:///D:/code/personal/coder-mate-core/highlighter.go#L210)、[L226](file:///D:/code/personal/coder-mate-core/highlighter.go#L226) |
| 4 | **P1 必修** | 正确性 | **并发竞态**：`languageRegistry` / theme registry / `builtinsRegistered`（裸 bool，非 `sync.Once`）/ highlighter cache 全部无锁，Go map 并发读写是进程级 fatal；作为库被服务端并发调用是常态 | 注册表 `sync.RWMutex`，builtins 改 `sync.Once`，cache 加锁或改 `sync.Map` | [core/registry.go#L15](file:///D:/code/personal/coder-mate-core/core/registry.go#L15)、[language/manager.go#L5-L12](file:///D:/code/personal/coder-mate-core/language/manager.go#L5-L12)、[highlighter.go#L21-L22](file:///D:/code/personal/coder-mate-core/highlighter.go#L21-L22) |
| 5 | P2 重要 | 正确性 | JS 规则缺口：数字正则 `\d+\.?\d*e?\d*n?` 不吃 `e` 后符号（`1e+10` 只匹配出 `1e`）；运算符表缺 `**`；`typeof/new/this/super/yield/delete` 关键字全库零命中（Grep 复核），一律落 identifier | 修正则 + 补关键字表，配 golden token 流回归 | [javascript/rules.go#L105](file:///D:/code/personal/coder-mate-core/language/javascript/rules.go#L105)、[L99](file:///D:/code/personal/coder-mate-core/language/javascript/rules.go#L99) |
| 6 | P2 重要 | 正确性 | Python 规则缺口：`assert/del/global/nonlocal` 关键字缺失（Grep 全库零命中），数字不支持 `0x` / `1_000` 字面量 | 同上，随 golden 回归一并修 | [python/rules.go](file:///D:/code/personal/coder-mate-core/language/python/rules.go) |
| 7 | P2 重要 | 抽象 | **强制抽象违规**：17 个包的 adapter.go 同骨架逐字复制（仅包名/ID/别名三处字符串不同，归一化 MD5 一致）、14/17 的 spec.go 归一化一致；块注释收尾正则 `^[\s\S]*?\*/` 在 10 个包重复 | 抽共享规则工厂 / 注册代码生成，语言包只留规则数据 | [language/*/adapter.go](file:///D:/code/personal/coder-mate-core/language/javascript/adapter.go)（17 份） |
| 8 | P2 重要 | 一致性 | typescript 对 JS 规则做"顺序手术"：硬编码 `"global"`、魔法数 `insertPos := 3` 按位置插入规则，JS 规则一改即静默破坏；且 JS 有深拷贝防护、其余 15 包裸暴露全局 map | 显式规则组合 API 替代位置插入 | [typescript/rules.go#L29](file:///D:/code/personal/coder-mate-core/language/typescript/rules.go#L29)、[L60](file:///D:/code/personal/coder-mate-core/language/typescript/rules.go#L60) |
| 9 | P2 重要 | 测试 | 断言两极：16 个语言包仅 `checkScopeExists`"存在即过"，不验证误分类；存在形式主义测试（建 spec 后 `_ = spec`、断言零输出）；零 fuzz，缺多字节字符 / 非法 theme / 并发用例 | 建 golden token 流基准（含中文用例），淘汰空断言测试 | [core/core_test.go#L559-L581](file:///D:/code/personal/coder-mate-core/core/core_test.go#L559-L581)、[language/*/adapter_test.go](file:///D:/code/personal/coder-mate-core/language/c/adapter_test.go)（16 份） |
| 10 | P2 重要 | 工程化 | CI / lint / Makefile 全缺；spec08 的"Go/TS token 一致性验证"、spec10 的 gofmt+lint 规划均未落地 | 最低成本起步：GitHub Actions 跑 `go vet + gofmt -l + go test ./...` 门禁 | 仓库根（无配置） |
| 11 | P2 重要 | 文档 | 文档漂移：agents.md 称"待进入实现阶段"、引用不存在的 `ts/` 目录、Go 版本三方打架 | 一次性校准 agents.md / README / spec10 的事实陈述 | [agents.md#L14-L17](file:///D:/code/personal/coder-mate-core/agents.md#L14-L17)、[L36](file:///D:/code/personal/coder-mate-core/agents.md#L36) |
| 12 | P3 可选 | 可读性 | `Parse` 75 有效行超 50 行告警线；魔法数 `Grow(len+16)`；`escapeHtml` 对 `` ` `` / `$` 过度转义 | 拆行处理分支；删冗余转义 | [tokenizer.go#L182-L256](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L182-L256)、[highlighter.go#L52-L55](file:///D:/code/personal/coder-mate-core/highlighter.go#L52-L55) |
| 13 | P3 可选 | 性能 | 缓存以整段源码为 key 且无淘汰，大文件重复调用无限膨胀 | LRU 上限或按长度短路 | [highlighter.go#L33-L34](file:///D:/code/personal/coder-mate-core/highlighter.go#L33-L34) |
| 14 | P3 可选 | 一致性 | 主题层硬编码色值遍布：darkplus 目录 300 处 hex 字面量（`#569CD6` 等），调色板改动需全局替换 | 色板常量化、按 token 引用 | [theme/darkplus/](file:///D:/code/personal/coder-mate-core/theme/darkplus/base.go) |

---

## 3. 关键证据（双向）

**好证据**：

1. **状态栈拷贝语义防共享篡改**：Push/Pop 均复制栈而非原地改（[tokenizer.go#L137-L149](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L137-L149)），token 样式深拷贝防缓存污染（[highlighter.go#L110-L126](file:///D:/code/personal/coder-mate-core/highlighter.go#L110-L126)）——并发 bug 之外的"数据防御"意识在线。
2. **零依赖 + 22 包测试全绿**：`go test ./...` 实测 21 个测试包全部 ok，仅 `theme/darkplus` 无测试文件；10 个 Benchmark。
3. **公共 API 克制**：3 个接口方法 + 选项模式（[api.go#L13-L42](file:///D:/code/personal/coder-mate-core/api.go#L13-L42)），没有过度设计。
4. **README 诚实**：性能表自曝 CacheHit 与 LargeJS 数字接近的短板，语言/别名表与代码一致（抽查 gorule/bash/c 吻合）。
5. **文档是真沉淀**：10 篇 spec 覆盖架构/引擎/适配器/主题/WASM/测试策略，含迁移计划与分阶段实施清单。

**坏证据**：

1. **核心算法级缺陷**：UTF-8 字节处理（[tokenizer.go#L158](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L158)）——中文输入即产出坏 token。
2. **确定性崩溃路径**：`NewHighlighter(HighlighterOptions{Theme:"不存在"})` → 吞错 → nil 解引用 panic，零测试覆盖。
3. **模板批量生成痕迹**：17 份同骨架 adapter + 16 份"存在即过"测试 helper 复制（第 2 节 No.7/9）。
4. **死代码 + 隐患**：`token == nil` 分支不可达（非空输入必返回回退 token），若可达则同样按字节跳过（[tokenizer.go#L211-L215](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L211-L215)）；零宽正则可匹配空串时 Parse 零长度消费**无限循环**（[tokenizer.go#L248](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L248) 无空匹配防护）。
5. **治理为零**：无 CI、无 lint、无发布流程，最后一笔提交停在 2026-05-18。

---

## 4. 团队画像

**8 分的架构判断力写了 6 分的代码，被"第二遍没做"拖累**：分层、状态机拷贝语义、文档先行的思路证明作者有能力写好库；但 17 份复制骨架、零 CI、吞错式收尾说明这是"一遍过"的迁移工程——第一遍移植能拿 8 分，抽象提取 / golden 回归 / 机器兜底这三件"第二遍"的事一件没做。另需注意：巴士因子 1，typescript 顺序手术式代码（No.8）已成为无人敢碰的模块。

---

## 5. 公允性说明

- **规模/年龄修正**：11 个提交、约 7.0k 行、活跃期仅 2026-05-14 → 05-18 约 4 天（单人冲刺），此后休眠近 5 个月。**按"单人 4 天 TS→Go 迁移原型"口径，本项目质量高于预期**（零依赖、测试全绿、文档齐）；本报告评分按"计划对外提供 WASM 高亮服务的库"口径给出，故工程化与并发安全从严。
- **新写代码口径**：全库皆为新写代码（TS→Go 重写），无"代码老"的豁免空间；UTF-8 / nil panic / 并发竞态均属新引入问题，不适用年龄修正。
- **代码评价 ≠ 个人考绩**：以上问题全部是工程流程缺口的产物，不是个人能力问题——换任何人在"4 天 + 无 CI"的条件下都会留下同类欠账。

---

## 6. 行动优先级 Top5（按风险 × 成本排序）

| 优先级 | 行动 | 风险 × 成本 | 对应失分维度 |
| --- | --- | --- | --- |
| 1 | **UTF-8 修复**：tokenizer 回退分支与列号计算改 rune 化，补中文测试用例 | 高风险（中文全挂）× 低成本（局部改动） | #3 语言纪律、#2 可维护性 |
| 2 | **消除吞错 + nil panic 链**：`newHighlighter` 上抛 error，`ResolveScopeStyle` nil 防护 | 高风险（确定性崩溃）× 低成本 | #3、#4（红线项） |
| 3 | **HTML 属性转义**：`preStyle` / `lineClassPrefix` 入属性前统一转义 | 高风险（注入面）× 低成本 | #4 安全 |
| 4 | **并发保护**：注册表 `RWMutex`、builtins `sync.Once`、cache 加锁 | 中风险（进程 fatal）× 低成本 | #3、#10 |
| 5 | **CI 门禁 + golden token 回归**：Actions 跑 `go vet/gofmt/test`；建含多字节与 JS/Python 规则缺口的 golden 流，随后修 No.5/6 规则 | 中风险（防回归失守）× 中成本 | #7 工程化（3 分）、#5 测试、#2 |

> 1–4 项合计不到一天工作量，可将红线清零并消除全部 P1；第 5 项是从"中"到"良"的分水岭。
