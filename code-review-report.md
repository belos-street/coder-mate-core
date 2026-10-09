# coder-mate-core 代码审查报告

> 审查日期：2026-10-09（v1.1，已合并 [project-evaluation-report.md](project-evaluation-report.md) 交叉核对后的采纳项）
> 审查方法：belos-street 八维审查流程（可读性/一致性/正确性/健壮性/耦合/抽象/性能/安全）+ golang-best-practices 语言规范映射 + 实证探针验证
> 审查环境：Go 1.26.5 / Windows / 全量测试 + `-race` + WASM 交叉编译均已实际执行
> 审查范围：全部 Go 源码（core / language×17 / theme / api / cmd/wasm）、22 个测试文件（272 个 `func Test`）、benchmark、README、agents.md、specs 01/03/06/08、git 历史

---

## 总体结论

**架构素质高，但存在 2 个按 README 用法即可触发的 P1 缺陷和 1 个违反项目自家核心约束的输出不确定性问题。**

- FSM 引擎核心（状态栈不可变转换、CRLF 处理、跨行 token 拆分）实现严谨，测试覆盖到了这些难点
- 分层清晰、适配器模式执行到位、零依赖、基准意识好——这些是同类个人项目里少见的优点
- 但**「库入口未注册内置主题」「字节级 fallback」「HTML 输出随机」**三个问题相互独立，且全部通过探针程序实证复现（见下节），建议优先修复

**项目背景（公允性口径）**：git 历史显示全库 11 个提交、两位作者账号（belos-street 10 + 酱汁 1，疑似同一人双账号）、活跃期仅 2026-05-14 → 05-18 约 4 天的单人迁移冲刺，此后休眠近 5 个月，无 tag/release，巴士因子 1。按「单人 4 天 TS→Go 迁移原型」口径，当前质量（零依赖、22 包测试全绿、文档体系完整）高于预期；本报告按「计划对外提供的 Go/WASM 高亮库」口径从严定级。最后一笔提交 `896d07e` 是三合一混合提交（README + 基准 + 测试 + **删除整个 ts/ 目录**）——后者正是 agents.md/specs 引用失效目录的直接原因。

### 实证结果（探针程序 + 复核命令，均已独立验证）

| 探针 | 输入 | 结果 |
| --- | --- | --- |
| 1 | 按 README 示例：`NewHighlighter()` → `CodeToHtml("const x = 1;", "javascript")` | **panic: nil pointer dereference**（纯 Go 消费者未手动注册主题时必现） |
| 2 | 注册主题后 `CodeToHtml("let 值 = 1", "javascript")` | 输出 `let å ¼`——CJK 字符 `值`（UTF-8: E5 80 BC）被按字节切成 3 个 token，字节值被转成 Latin-1 码点（å + U+0080 + ¼），典型 mojibake |
| 3 | 同一 highlighter、同一输入重复高亮 2000 次 | **产生 4 种不同的 HTML 输出**（style 属性顺序随 map 遍历随机） |
| 4 | Grep/git 复核 | JS 关键字表无 `typeof/new/this/super/yield/delete`；Python 关键字表无 `assert/del/global/nonlocal`、数字规则无 `0x`/`_` 分隔；JS 数字正则 `1e+10` 只匹配出 `1e`、运算符表无 `**`；11 提交/最后提交 2026-05-18/无 tag 全部属实 |

其余验证：`go vet` 干净；`go test ./...` 全绿；`go test -race ./...` 全绿（但测试全为串行，race 检测器实际未被触发）；`GOOS=js GOARCH=wasm go build ./cmd/wasm/` 成功。

---

## 问题清单

> 定级沿用 belos-street 约定：P1 必修（合入前必须修）/ P2 重要（尽快处理）/ P3 可选（择机处理）。标注 ⛽ 的条目采纳自 project-evaluation-report.md 并经独立复核；定级与对方不一致处均注明理由。

### P1 必修

| No. | 维度 | 问题 | 建议 | 位置 |
| --- | --- | --- | --- | --- |
| 1 | 健壮性 | **README 首选用法即 panic**（该吞错链按 belos-street 红线口径命中「静默吞错」）。`newHighlighter` 用 `currentTheme, _ :=` 吞掉 `ResolveTheme` 的错误，而 Go 库路径上没有任何代码调用 `theme.RegisterBuiltInThemes()`（仅测试 init 和 `cmd/wasm/main.go:23` 调用）。主题注册表为空 → `currentTheme = nil` → 首次渲染非空代码时 `ResolveScopeStyle` 解引用 nil theme panic | 在 `newHighlighter` 中通过 `sync.Once` 注册内置主题（与 `language.EnsureBuiltinsRegistered()` 对称），并且不吞 `ResolveTheme` 错误——让 `NewHighlighter` 返回 `(*Highlighter, error)` 或对无效 Theme fail-fast | [api.go#L55](file:///D:/code/personal/coder-mate-core/api.go#L55)、[theme/registry.go#L95](file:///D:/code/personal/coder-mate-core/theme/registry.go#L95) |
| 2 | 正确性 | **非 ASCII 字符被字节级切碎成 mojibake**。`MatchToken` 的 fallback 是 `char := string(code[0])`——取单字节并把**字节值当 Unicode 码点**转字符。12 个语言的 global 状态末条规则是 ASCII 标识符而非 `^.` 兜底，使该路径可达（探针 2 实证） | fallback 改用 `utf8.DecodeRuneInString(code)` 按 rune 解码并按 rune 宽度推进；同时明确 `Col` 语义（当前为字节偏移，与 rune/UTF-16 均不一致，文档应注明）。可另为 12 个语言的 global 态补 `^.` catch-all 作防御（修复本条后非必须） | [core/tokenizer.go#L158](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L158)、[core/tokenizer.go#L133](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L133)；缺 catch-all 的语言见「八维小结·正确性」 |
| 3 | 正确性（一致性） | **同一输入产生不同 HTML**。`stringifyInlineStyle` 用 `range style` 遍历 map，多属性 style（如块注释的 `color` + `font-style`）属性顺序随机。探针 3 实证 2000 次产生 4 种输出，直接违反 agents.md 核心约束「相同代码 → 相同 HTML/Token 流」，也使快照测试、HTTP 缓存/ETag 全部不可行 | 按 key 排序后序列化（约 3 行改动），随后即可补确定性快照测试 | [highlighter.go#L98](file:///D:/code/personal/coder-mate-core/highlighter.go#L98) |

### P2 重要

| No. | 维度 | 问题 | 建议 | 位置 |
| --- | --- | --- | --- | --- |
| 4 | 健壮性 | **全链路无并发保护**：highlighter 的 `cache`/`parsedStyleCache` 读写、`UpdateTheme` 边遍历边写、`core.languageRegistry`/`theme.themeRegistry` 全局 map、`builtinsRegistered` 裸 bool check-then-act 全部无锁。WASM 单线程无碍，但 README 将 Go API 作为一等卖点，HTTP 服务多 goroutine 复用同一 Highlighter 即 data race（Go map 并发读写是进程级 fatal）；现有测试全串行，`-race` 从未真正覆盖此面 ⛽ | 最小改法：`highlighter` 加 `sync.Mutex`，注册表注册完成后只读（但需 `sync.Once` 保证首次注册的可见性）；或至少在 README/doc 注明「非并发安全」 | [highlighter.go#L162](file:///D:/code/personal/coder-mate-core/highlighter.go#L162)、[highlighter.go#L237-L242](file:///D:/code/personal/coder-mate-core/highlighter.go#L237-L242)、[core/registry.go#L15](file:///D:/code/personal/coder-mate-core/core/registry.go#L15)、[language/manager.go#L5-L13](file:///D:/code/personal/coder-mate-core/language/manager.go#L5-L13) |
| 5 | 一致性 | **注册错误被静默丢弃（3 处）**：`RegisterBuiltins` 忽略 `core.RegisterLanguage` 返回值；`RegisterBuiltInThemes` 忽略 `RegisterTheme` 返回值；加上 P1-1 的 `currentTheme, _ :=`。语言/主题别名冲突时将无声失败，排查困难 | 循环内 `if err := ...; err != nil { panic/log }` 或聚合返回 error；配合 P1-1 一并处理 | [language/builtins.go#L45-L47](file:///D:/code/personal/coder-mate-core/language/builtins.go#L45-L47)、[theme/presets.go#L289-L292](file:///D:/code/personal/coder-mate-core/theme/presets.go#L289-L292) |
| 6 | 一致性（文档） | **文档与代码严重漂移**：agents.md 仍描述 `ts/` 参考目录（实际不存在，且 `896d07e` 提交正是删除它的那笔）并声称「待进入实现阶段」（实现已完成）；specs 01/03/06 大量引用 `ts/lib/` 路径；**Go 版本三方打架**（agents.md `1.21+` / go.mod `1.23.5` / README `1.23+`）⛽；spec 08 §2 承诺的 TS 一致性 golden 验证（每语言 3-5 片段 JSON diff）**完全未落地**——仓库无 testdata、无 golden 文件，而「与 TS 版输出一致」是 agents.md 列的第一条核心约束，目前只能靠肉眼；spec 08 的基准函数名与实际命名也不符 | 首选：删除 agents.md/specs 中对 `ts/` 的引用或注明「参考实现已随 `896d07e` 移除，以 golden 文件为准」；统一 Go 版本陈述；落地最小 golden 测试（确定性修复 P1-3 之后即可做） | [agents.md#L14-L17](file:///D:/code/personal/coder-mate-core/agents.md#L14-L17)、[.agents/docs/specs/08-testing-strategy.md](file:///D:/code/personal/coder-mate-core/.agents/docs/specs/08-testing-strategy.md) |
| 7 | 一致性（工程化） | **全仓库 100% gofmt 不合规**：仓库无 `.gitattributes`，`git ls-files --eol` 证实 index 为 LF、工作区为 CRLF，`gofmt -l` 列出全部 122 个 .go 文件。任何 gofmt CI 门禁会直接全红，多人/多机协作会持续产生伪 diff | 添加根级 `.gitattributes`（`*.go text eol=lf`）→ `git add --renormalize .` → 一次性 `gofmt -w .` 提交（与 kb-vault 仓库同类问题的既定修法一致） | 仓库根目录 |
| 8 | 耦合 | **TS 适配器对 JS 规则顺序的位置耦合**：TS 规则插入到过滤后 JS global 规则的第 3 位（`insertPos := 3` 魔法数，无注释），正确性完全依赖「JS 前 3 条是注释/动态 import 规则」这一隐式契约。JS 规则表重排会静默破坏 TS 语法且无测试报警 | 加注释说明「插在注释规则之后」的意图，并改为按 Scope/锚定规则定位插入点；或干脆追加到注释规则之后用显式断言锁住 JS 规则顺序 | [language/typescript/rules.go#L60-L69](file:///D:/code/personal/coder-mate-core/language/typescript/rules.go#L60-L69) |
| 9 | 性能（文档） | **缓存卖点名不副实**：README 宣称「Token caching — repeated highlights return instantly」，但实测基准 CacheHit 165,755 ns/op ≈ LargeJS 全流程 165,096 ns/op——缓存只省了 parse，clone + 重渲染占绝对主导（README L168 自己也承认了这一点，与 L17 的卖点文案互相矛盾） | 卖点改为「避免重复解析（re-parse）」；若要真收益，可缓存渲染后的 HTML 字符串（同 key 命中直接返回）而非 token 流 | [README.md#L17](file:///D:/code/personal/coder-mate-core/README.md#L17)、[benchmark/bench_test.go#L138-L154](file:///D:/code/personal/coder-mate-core/benchmark/bench_test.go#L138-L154) |
| 10 | 安全 ⛽ | **HTML 属性注入面**：`lineClassPrefix` 直写 `class` 属性、`preStyle` 直接拼进 `style` 属性，均未转义，含 `"` 即可破出属性边界 | 属性值统一转义（至少 `"` `<` `&`）。**定级说明**：对方报告定 P1；本报告降为 P2——这两个选项默认是开发者可控配置而非外部攻击者输入（token 文本这一真正的高危面已被 `escapeHtml` 覆盖），但上层应用若将用户偏好透传到这两个选项即成真实 XSS，作为库应默认转义 | [highlighter.go#L210](file:///D:/code/personal/coder-mate-core/highlighter.go#L210)、[highlighter.go#L226](file:///D:/code/personal/coder-mate-core/highlighter.go#L226) |
| 11 | 正确性（语法保真）⛽ | **JS/Python 关键字与字面量缺口**（经 Grep 复核）：javascript 关键字表**无 `typeof/new/this/super/yield/delete`**（TS global 态继承同一缺口，唯 `typeof` 经 TS 类型运算符规则覆盖）；JS 数字正则不吃 `e` 后符号（`1e+10` 只匹配出 `1e`，指数部分被拆成 `+`/`10`）；JS 运算符表无 `**`；Python 关键字表**无 `assert/del/global/nonlocal`**；Python 数字规则无 `0x`/`0o`/`0b` 前缀与 `1_000` 分隔符（`0x1F` 会被切成 `0` + `x1F` 标识符）。另合并原清单的规则顺序问题：java Type 先于 Function（`Foo(` 恒判为类型名）、yaml numeric 先于 `-`（`-5` 整体吃成数字） | 补关键字表与数字正则（Python 的 `[eE][+-]?\d+` 与 `\*\*`/`//`/`:=` 处理可作参照），修正两条规则顺序；随 golden token 流回归一并落地 | [language/javascript/rules.go#L99](file:///D:/code/personal/coder-mate-core/language/javascript/rules.go#L99)、[language/javascript/rules.go#L105](file:///D:/code/personal/coder-mate-core/language/javascript/rules.go#L105)、[language/python/rules.go#L10](file:///D:/code/personal/coder-mate-core/language/python/rules.go#L10)、[language/python/rules.go#L49](file:///D:/code/personal/coder-mate-core/language/python/rules.go#L49)、[language/java/rules.go#L31-L32](file:///D:/code/personal/coder-mate-core/language/java/rules.go#L31-L32)、[language/yaml/rules.go#L26-L27](file:///D:/code/personal/coder-mate-core/language/yaml/rules.go#L26-L27) |
| 12 | 工程化 ⛽ | **零工程化设施**：无 CI、无 lint 配置、无 Makefile。spec 10 规划的 gofmt + lint、spec 08 的一致性验证均「只立法不执法」；结合 #7（gofmt 全红）与休眠 5 个月的现状，任何回归都不会被机器拦截 | 最低成本起步：GitHub Actions 三件套 `go vet + gofmt -l + go test ./...` 门禁；后续补 `-race` 与 golden 一致性 job | 仓库根（无配置文件） |

### P3 可选

| No. | 维度 | 问题 | 建议 | 位置 |
| --- | --- | --- | --- | --- |
| 13 | 可读性 | 死代码与重复逻辑：`PushState`/`PopState` 仅被 core_test 使用，`MatchToken` 内联重写了 push/pop 且行为不一致（pop 到底：内联版静默忽略，独立版返回 error）；`Parse` 中 `token == nil` 分支不可达（`MatchToken` 仅在空输入返回 nil，而外层循环已保证非空——但若可达则同样按字节跳过 ⛽）；`MatchToken` 内 `len(code) > 0` 守卫恒真 | 删除 `Parse` 死分支与恒真守卫；`MatchToken` 改调 `PushState`/`PopState` 消除双实现（pop 到底语义需统一：建议静默，与 TS 版对齐） | [core/tokenizer.go#L17-L35](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L17-L35)、[core/tokenizer.go#L138-L150](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L138-L150)、[core/tokenizer.go#L211-L215](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L211-L215) |
| 14 | 健壮性 ⛽ | **零宽正则可致 Parse 无限循环**：`Parse` 按消费 `len(token.Text)` 推进，若规则正则可匹配空串则零长度消费死循环。经核查当前 17 个语言包的内置规则均无零宽匹配（`^.*` 类兜底被 Parse 的换行预检屏蔽），属潜在隐患而非现存 bug——但库开放自定义 `TokenizerSpec` 时缺少防护 | `Parse` 内加 `if len(token.Text) == 0 { 消费 1 字节 }` 防护（一行），杜绝自定义规则触发挂死 | [core/tokenizer.go#L248](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L248) |
| 15 | 一致性 | Go 命名惯例：`GRAMMAR_RULES`、`THEME_ALIAS_MAP` 全大写下划线风格（应为 camelCase）；`copy_` 下划线后缀；`ResolveScopeStyle` 参数名 `theme` 遮蔽包名；spec 导出模式两套并存——16 个语言用包级 `var Spec`（且为**可变导出全局 map**，外部可篡改语法表），唯 typescript 用 `sync.Once` + `GetSpec()` 封装，javascript 的 `GetRules()` 又做了深拷贝防护 ⛽ | 改为 `grammarRules` / `themeAliasMap` / `copyRules` / 参数名 `t`；统一 spec 封装模式（要么都深拷贝、要么都不可变） | [language/javascript/rules.go#L11](file:///D:/code/personal/coder-mate-core/language/javascript/rules.go#L11)、[language/javascript/rules.go#L20](file:///D:/code/personal/coder-mate-core/language/javascript/rules.go#L20)、[theme/registry.go#L10](file:///D:/code/personal/coder-mate-core/theme/registry.go#L10)、[theme/registry.go#L94](file:///D:/code/personal/coder-mate-core/theme/registry.go#L94)、[language/typescript/spec.go#L14-L24](file:///D:/code/personal/coder-mate-core/language/typescript/spec.go#L14-L24) |
| 16 | 抽象 ⛽ | **17 个语言包的 adapter.go 同骨架逐字复制**（仅 ID/别名/Parse 三处字符串不同），注释块收尾正则 `^[\s\S]*?\*/` + `^[\s\S]+` 对与字符串状态四件套在 7 个包近乎逐字重复。按 belos-street 抽象规则（>4 行逻辑出现 ≥3 次强制抽象）够 P2；考虑到语言包内容稳定、改动频率低，降为 P3 择机 | `core` 提供 `NewSimpleAdapter(id, aliases, spec)` 构造函数与公共规则工厂（注释块/字符串四件套），语言包只留规则数据 | [language/javascript/adapter.go](file:///D:/code/personal/coder-mate-core/language/javascript/adapter.go)（17 份同构） |
| 17 | API 细节 | `UpdateTheme(t interface{})` / `HighlighterOptions.Theme` 用 `interface{}` 而非 `any`（Go 1.23）；`NewHighlighter(opts ...HighlighterOptions)` 可变参收单个配置属反模式（直接收 `HighlighterOptions` 即可）⛽；`newHighlighter` 里 `ResolveAlias` + `ResolveTheme` 双重解析（`GetTheme` 内部已处理别名），冗余；README `UpdateTheme` 示例忽略 error 返回值 | `interface{}` → `any`；去掉冗余的 `ResolveAlias` 调用；README 示例处理 error | [api.go#L16](file:///D:/code/personal/coder-mate-core/api.go#L16)、[api.go#L44](file:///D:/code/personal/coder-mate-core/api.go#L44)、[api.go#L55](file:///D:/code/personal/coder-mate-core/api.go#L55) |
| 18 | 正确性 | 非确定性 API：`core.ListLanguages` / `theme.ListThemes` 直接遍历 map，返回顺序随机 → WASM `getLanguages()` / `getThemes()` 每次调用顺序可能不同 | 按 ID 排序后返回 | [core/registry.go#L56-L70](file:///D:/code/personal/coder-mate-core/core/registry.go#L56-L70)、[theme/registry.go#L58-L64](file:///D:/code/personal/coder-mate-core/theme/registry.go#L58-L64) |
| 19 | 性能 ⛽ | **缓存以整段源码为 key 且无淘汰上限**：服务端对任意用户输入调用时 cache 无限膨胀（叠加 #4 的并发问题） | LRU 上限或按长度短路（超长代码不缓存） | [highlighter.go#L33-L34](file:///D:/code/personal/coder-mate-core/highlighter.go#L33-L34) |
| 20 | 性能 | 微优化：`fmt.Sprintf("%d", rowIndex+1)` → `strconv.Itoa`；空 style 的 token 仍输出 `<span style="">`（可跳过属性）；`escapeHtml` 转义反引号/`$` 非必需（无害）；`b.Grow(len+16)` 魔法数 ⛽；`Parse` 有效逻辑行数约 60 行、超 50 行警线（拆分换行处理分支可解）⛽；`ensureValidID` 在库代码里 panic（宜改返回 error） | 择机处理 | [highlighter.go#L211](file:///D:/code/personal/coder-mate-core/highlighter.go#L211)、[highlighter.go#L214-L221](file:///D:/code/personal/coder-mate-core/highlighter.go#L214-L221)、[highlighter.go#L39](file:///D:/code/personal/coder-mate-core/highlighter.go#L39)、[core/tokenizer.go#L182-L256](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L182-L256)、[core/registry.go#L21-L27](file:///D:/code/personal/coder-mate-core/core/registry.go#L21-L27) |

---

## 八维小结

| 维度 | 评价 | 说明 |
| --- | --- | --- |
| 可读性 | 良好 | 规则表按「注释/模块化/关键字/运算符/值类型/标识符」分节注释，命名达意；扣分项：`insertPos := 3` 魔法数（#8）、`Parse` 超 50 行警线（#20） |
| 一致性 | 中等 | 17 个适配器五件套结构与 `var Language` 模式高度统一（亮点）；但 spec 缓存两套模式并存且防护不对齐（#15）、错误处理三种方式混用（panic/error/吞，#5）、C 风格全大写命名混入 |
| 正确性 | 需修复 | 三个 P1 均在此维（panic / mojibake / 输出随机）；核心 FSM 逻辑本身（状态栈 copy-on-write、CRLF、跨行拆分）实现严谨且有边界测试；语法保真缺口单列（#11） |
| 健壮性 | 需修复 | 错误静默传导成延迟 panic（P1-1）、并发裸奔（#4）、nil theme 无防护（P1-1）、零宽正则无防护（#14） |
| 耦合 | 良好 | core/language/theme 单向分层，无反向依赖；TS→JS 规则派生 DRY 但引入位置耦合（#8） |
| 抽象 | 良好 | 适配器接口统一、主题用 DarkPlus 基座 + 12 色板替换派生出 9 个变体（DRY 意识强）；17 份 adapter 骨架复制属「提取即可复用」欠账（#16） |
| 性能 | 良好 | 有独立 benchmark 包且 README 数据与基准名一一对应；逐规则线性正则匹配是与 TS 版一致的设计取舍；缓存收益被 clone 抵消（#9）且无淘汰上限（#19） |
| 安全 | 良好偏上 | `escapeHtml` 覆盖 token 文本（真正的高危面）且单测覆盖 8 种字符；零第三方依赖无 CVE 面；遗留属性注入面（#10）；无密钥/SQL 注入面 |

## 测试缺口 Top 5

1. **TS 一致性 golden 验证零落地**——spec 08 §2 的核心承诺，「与 TS 输出一致」这一首要约束目前无任何自动化验证
2. **无确定性快照测试**——所有 HTML 断言均为 `strings.Contains`，恰好掩盖了 P1-3 的输出随机性
3. **并发安全零测试**——无 `-race` 下多 goroutine 调用 `CodeToHtml`/`UpdateTheme` 的用例
4. **无多字节用例**——17 个语言包的 adapter_test 均未覆盖中文/emoji/变音字母（P1-2 的风险面完全裸奔）
5. **断言两极分化**——javascript 包有 30+ 个精确断言（`tok.Text != ...` 逐 token 校验），其余 16 包普遍为 `checkScopeExists`「存在即过」式弱断言，不验证误分类；另有零输出形式主义测试（建 spec 后 `_ = spec`）

补充：测试总量不低（22 个文件 / 272 个 `func Test`），但测试卫生存在问题——`ClearRegistry`/`ResetBuiltins`/`RegisterBuiltInThemes` 在 init 和各测试中反复操作全局注册表，测试无法并行。

## 建议修复路线（最小改动优先）

1. **P1-1**（~5 行）：`newHighlighter` 内 `sync.Once` 注册内置主题 + 不吞 `ResolveTheme` 错误 → 消灭 README 用法 panic
2. **P1-2**（~10 行）：fallback 改 `utf8.DecodeRuneInString` 按 rune 切分 → 消灭 mojibake
3. **P1-3**（~3 行）：`stringifyInlineStyle` 按 key 排序 → 输出确定
4. **#7/#12**：`.gitattributes` + renormalize + 一次性 `gofmt -w .`；GitHub Actions 三件套门禁
5. **#4/#5/#10**：highlighter 加 mutex（或文档声明非并发安全）+ 注册错误不再丢弃 + 属性值转义
6. 补测试：快照（依赖 P1-3）→ 多字节用例 → `-race` 并发用例 → golden token 流（含 #11 的关键字缺口与中文用例），随后修 #11
7. **#6**：清理 agents.md/specs 对 `ts/` 的过期引用，统一 Go 版本陈述，spec 08 基准名对齐

以上 1–5 项合计约一天内工作量，可清零全部 P1。

## 亮点（值得保持）

- **零依赖纯标准库**，`go vet` 干净，WASM 交叉编译开箱即用
- **copy-on-write 状态栈**：所有状态转换返回新 `ParserContext` 而非原地修改，token 样式深拷贝防缓存污染——「数据防御」意识贯穿核心
- **两层基座派生**：TS 规则派生自 JS 基座（过滤+重排+追加）、9 个主题派生自 DarkPlus 基座（palette 色板替换），把「同类事物」的变更收敛到单点
- `SplitTokenByLineBreak` 的 CRLF 边界处理完整（连续换行、前导换行、CRLF 混合均有专项测试）
- WASM 层有 mutex 防护实例注册表、`dispose` 生命周期管理，JS API 错误以 `{error}` 结构返回而非抛异常
- agents.md → specs → plans 的 AI 协作文档体系完整（虽然内容需要更新，但骨架是对的）
- README 诚实：性能表自曝 CacheHit 与 LargeJS 数字接近的短板，语言/别名表与代码逐项一致

---

## 附录：与 project-evaluation-report.md 的交叉核对

对方报告的**事实主张全部复核通过**（git 历史、关键字缺失、测试数量 272/22 等逐条 Grep/Shell 验证），以下为采纳与分歧记录：

| 对方条目 | 处理 | 说明 |
| --- | --- | --- |
| 红线「静默吞错」（api.go:55、builtins.go:46） | ✅ 采纳 | 本报告 P1-1（含红线口径标注）+ #5 |
| UTF-8 字节级处理 | ✅ 采纳 | 本报告 P1-2，探针已实证（对方未做运行时探针） |
| HTML 属性注入 | ⚠️ 采纳但降级 P1→P2 | 缺陷真实；但 `preStyle`/`lineClassPrefix` 默认为开发者配置而非攻击者可控输入，真高危面（token 文本）已被转义。红线定义要求「参数外部可控」，不满足 |
| 并发竞态 | ✅ 采纳 | 本报告 #4，与对方结论一致 |
| JS/Python 规则缺口 | ✅ 采纳 | 本报告 #11。修正一处表述：`typeof` 并非「全库零命中」——TS global 态经类型运算符规则覆盖了它，但 javascript 包确实全缺；`new/this/super/yield/delete` 在 JS/TS 两侧均缺失属实 |
| CI/lint/Makefile 全缺 | ✅ 采纳 | 本报告 #12 |
| Go 版本三方打架 | ✅ 采纳 | 并入本报告 #6 |
| 17 份 adapter 骨架复制 | ✅ 采纳但降级 P2→P3 | 按 belos-street 抽象规则够 P2；考虑语言包内容稳定、改动频率低，择机处理即可 |
| 缓存无上限 | ✅ 采纳 | 本报告 #19 |
| 零宽正则无限循环 | ✅ 采纳但修正范围 | 缺陷机理成立；经核查当前内置规则无零宽匹配路径（`^.*` 类兜底被 Parse 换行预检屏蔽），定性为「自定义 spec 防护缺失」而非现存 bug，见 #14 |
| 可变参 NewHighlighter / interface{} / Parse 行数 / Grow(16) | ✅ 采纳 | 并入 #17、#20 |
| 团队画像 / 巴士因子 / 公允性口径 | ✅ 采纳 | 并入「总体结论·项目背景」 |
| 10 维评分卡（含巴士因子 2 分等） | ❌ 不并入 | 评分属对方「仓库级评测」方法论，无法逐分审计；本报告保持八维证据式清单，只吸收其可验证的事实发现 |
| darkplus 300 处 hex 硬编码 | ❌ 不采纳 | 色值是主题的数据本体，必须存放于某处；presets.go 的 palette 派生机制已解决「调色板全局替换」诉求（12 个映射色改一处即生效），该条与现有设计矛盾 |
| 「stringifyInlineStyle 尾分号」误报剔除 | ✅ 认可 | 对方复核剔除的这条与本报告结论一致（尾分号为合法 CSS，未列入问题） |
