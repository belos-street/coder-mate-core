# coder-mate-core 优化 TODO

> 综合两份审查产物整理（八维审查 + 运行时探针实证、仓库级十维评测）。两份报告的事实发现已互相复核，**裁定结论已内联于下方裁定表，本文件自包含、不依赖任何外部文件**。原文备查：[code-review-report.md](code-review-report.md)、[project-evaluation-report.md](project-evaluation-report.md)（链接仅供深挖细节用，缺失不影响本 TODO 执行）。
>
> 验收总口径见文末；每个阶段内条目可独立实施、独立验证。

---

## 定级分歧裁定（先看这里）

| 分歧点 | 评测报告 | 代码审查报告 | 最终裁定 |
| --- | --- | --- | --- |
| HTML 属性注入（preStyle/lineClassPrefix） | P1 | P2 | **P2**。采纳降级理由：两项属开发者配置而非攻击者输入，真高危面（token 文本）已被 `escapeHtml` 覆盖；但库应默认转义——若上层应用把用户偏好透传到这两个选项，立即升 P1 |
| 17 份 adapter 骨架复制 | P2 | P3 | **P3**。语言包内容稳定、改动频率低，择机抽象 |
| 并发竞态 | P1 | P2 | **P2 但排期靠前**。Go map 并发读写是进程级 fatal，库被服务端多 goroutine 调用是常态；排在三个确定性缺陷（阶段 1）之后立即做 |
| darkplus 色值硬编码 | P3 | 不采纳 | **移除**。色值是主题数据本体，palette 派生机制已解决调色板全局替换诉求 |
| 「`typeof` 全库零命中」表述 | 过强 | 修正 | **采纳修正**：TS 经类型运算符规则覆盖了 `typeof`，真实缺口在 javascript 包（`new/this/super/yield/delete` 两侧均缺） |

另采纳代码审查报告对缺陷触发面的**升级**：nil theme panic 不止"传错主题名才崩"——库路径从不注册内置主题，**README 首选用法即 panic**（探针 1）。此为阶段 1 第一优先。

---

## 阶段 1：P1 清零（确定性缺陷，逐项小改、逐项验收）

### 1.1 主题自动注册 + 消除吞错（红线项，探针 1）

- [x] 新增 `theme.EnsureBuiltinsRegistered()`（`sync.Once` 包装 `RegisterBuiltInThemes`），在 `newHighlighter` 中与 `language.EnsureBuiltinsRegistered()` 对称调用 — [presets.go#L288](file:///D:/code/personal/coder-mate-core/theme/presets.go#L288)、[api.go#L53-L55](file:///D:/code/personal/coder-mate-core/api.go#L53-L55)
- [x] **方案 A（推荐，零 API 破坏）**：仅凭上一条 + 主题解析兜底即可消除 panic——默认主题必已注册后，string/nil 主题路径永不返回 nil；未知主题名与无效类型主题**回退默认主题**（此回退是文档化行为而非吞错，须在 godoc/README 注明"未知主题回退默认主题"，与红线区分）。`NewHighlighter` 签名保持不变
- [ ] （备选方案 B，挪至阶段 4 与 API 细节一并做）`NewHighlighter` 改返回 `(*Highlighter, error)`。**影响面（变更时必须同步）**：[cmd/wasm/main.go#L25](file:///D:/code/personal/coder-mate-core/cmd/wasm/main.go#L25)、README 示例、[api_test.go](file:///D:/code/personal/coder-mate-core/api_test.go)、[benchmark/bench_test.go#L16](file:///D:/code/personal/coder-mate-core/benchmark/bench_test.go#L16) init
- [x] `RegisterBuiltins` / `RegisterBuiltInThemes` 循环内不再丢弃注册 error（聚合返回或 panic）— [language/builtins.go#L45-L47](file:///D:/code/personal/coder-mate-core/language/builtins.go#L45-L47)、[theme/presets.go#L289-L292](file:///D:/code/personal/coder-mate-core/theme/presets.go#L289-L292)
- [x] 验收：README 用法 `NewHighlighter()` → `CodeToHtml("const x = 1;", "javascript")` 不 panic（含方案 A 的回退行为单测）；新增单测覆盖"未手动注册主题"路径

### 1.2 UTF-8 rune 化（探针 2）

- [x] `MatchToken` fallback 改 `utf8.DecodeRuneInString(code)` 按 rune 解码并按 rune 宽度推进 — [tokenizer.go#L158](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L158)
- [x] 明确 `Col` 语义（当前为字节偏移，与 rune/UTF-16 均不一致）：文档注明或改按 rune 计 — [tokenizer.go#L133](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L133)
- [x] 补多字节测试用例：至少 17 个语言包各 1 条中文/emoji 断言（当前 adapter_test 多字节覆盖为零）
- [x] 验收：`CodeToHtml("let 值 = 1", ...)` 输出完整 `值` 字符（非 mojibake）；`go test ./...` 全绿

### 1.3 输出确定性（探针 3）

- [x] `stringifyInlineStyle` 按 key 排序后序列化（约 3 行改动）— [highlighter.go#L96-L107](file:///D:/code/personal/coder-mate-core/highlighter.go#L96-L107)
- [x] `ListLanguages` / `ListThemes` 按 ID 排序后返回（消除 map 遍历随机序）— [core/registry.go#L56-L70](file:///D:/code/personal/coder-mate-core/core/registry.go#L56-L70)、[theme/registry.go#L58-L64](file:///D:/code/personal/coder-mate-core/theme/registry.go#L58-L64)
- [x] 补确定性快照测试：同一输入连续 2000 次高亮输出 byte-equal（当前所有 HTML 断言均为 `strings.Contains`，恰好掩盖了此问题）
- [ ] 验收：探针 3 转绿（2000 次唯一输出）；快照测试进 CI

### 1.4 零宽正则防护

- [x] `Parse` 内加一行防护：`if len(token.Text) == 0 { 按 1 rune 消费 }`，杜绝自定义 `TokenizerSpec` 触发零长度消费死循环 — [tokenizer.go#L248](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L248)
- [x] 验收：构造零宽匹配规则的 spec 不再挂死（新增测试）

---

## 阶段 2：P2 重要（阶段 1 完成后立即做）

### 2.1 并发安全（定级分歧裁定为 P2，但排期最前）

- [x] `highlighter` 加 `sync.Mutex` 保护 `cache` / `parsedStyleCache`，`UpdateTheme` 边遍历边写处一并加锁 — [highlighter.go#L162](file:///D:/code/personal/coder-mate-core/highlighter.go#L162)、[L237-L242](file:///D:/code/personal/coder-mate-core/highlighter.go#L237-L242)
- [x] 全局注册表（`core.languageRegistry` / `theme.themeRegistry`）`sync.RWMutex`，或"注册完成后只读"+ `sync.Once` 保证首次注册可见性 — [core/registry.go#L15](file:///D:/code/personal/coder-mate-core/core/registry.go#L15)
- [x] `builtinsRegistered` 裸 bool check-then-act → `sync.Once`；`theme.EnsureBuiltinsRegistered` 的"查 map 再注册"同为 check-then-act 竞态面，加锁时一并覆盖（可并入注册表 RWMutex：存在性检查放读锁下）— [language/manager.go#L5-L13](file:///D:/code/personal/coder-mate-core/language/manager.go#L5-L13)、[theme/presets.go#L292](file:///D:/code/personal/coder-mate-core/theme/presets.go#L292)
- [x] 补 `-race` 并发测试：多 goroutine 同时调 `CodeToHtml` / `UpdateTheme`（现有测试全串行，`-race` 从未真正触发）
- [x] 验收：`go test -race ./...` 在并发用例真实触发检测的前提下全绿

### 2.2 HTML 属性转义

- [x] `lineClassPrefix` 入 `class` 属性、`preStyle` 入 `style` 属性前统一转义（至少 `"` `<` `&`）— [highlighter.go#L210](file:///D:/code/personal/coder-mate-core/highlighter.go#L210)、[L226](file:///D:/code/personal/coder-mate-core/highlighter.go#L226)
- [x] 验收：含引号的 `preStyle` / `LineClassPrefix` 不破出属性边界（新增测试）

### 2.3 gofmt 治理 + CI 门禁

- [x] 根级 `.gitattributes`（`*.go text eol=lf`）→ `git add --renormalize .` → 一次性 `gofmt -w .` 提交（实测 `gofmt -l` 列出全部 123 个 .go 文件：index LF / 工作区 CRLF / 无 attributes，不先治行尾任何 gofmt CI 门禁直接全红）
- [x] GitHub Actions 三件套门禁：`go vet` + `gofmt -l`（零输出）+ `go test ./...`
- [ ] 后续 job：`-race` + TS 一致性 golden（此处 golden = 阶段 3 的跨仓库 TS 对照基建，依赖阶段 3；与 2.4 的内联期望 token 流无关）（注：-race 已随首版 CI 一并加入，golden 待阶段 3）
- [x] 验收：`gofmt -l` 零输出；故意引入格式错误 / 测试失败可被 CI 拦截

### 2.4 语法规则缺口修正（每项先写断言再修正则；断言用**内联期望 token 流**——表驱动测试里的 `expected []core.Token`，不依赖 TS 对照 golden 基建，本阶段即可独立完成）

- [x] JavaScript：补关键字 `typeof/new/this/super/yield/delete`（javascript 包全缺）— [javascript/rules.go](file:///D:/code/personal/coder-mate-core/language/javascript/rules.go)
- [x] JavaScript：数字正则支持 `e` 后符号（`1e+10` 现在只匹配出 `1e`）— [javascript/rules.go#L105](file:///D:/code/personal/coder-mate-core/language/javascript/rules.go#L105)；运算符表补 `**` — [javascript/rules.go#L99](file:///D:/code/personal/coder-mate-core/language/javascript/rules.go#L99)
- [x] Python：补关键字 `assert/del/global/nonlocal` — [python/rules.go#L10](file:///D:/code/personal/coder-mate-core/language/python/rules.go#L10)；数字支持 `0x/0o/0b` 与 `1_000` 分隔（`0x1F` 现被切成 `0` + `x1F`）— [python/rules.go#L49](file:///D:/code/personal/coder-mate-core/language/python/rules.go#L49)
- [x] Java：Type 规则先于 Function 导致 `Foo(` 恒判类型名，调整顺序 — [java/rules.go#L31-L32](file:///D:/code/personal/coder-mate-core/language/java/rules.go#L31-L32)
- [x] YAML：numeric 先于 `-` 导致 `-5` 整体吃成数字，调整顺序 — [yaml/rules.go#L26-L27](file:///D:/code/personal/coder-mate-core/language/yaml/rules.go#L26-L27)
- [x] 验收：内联期望 token 断言先红后绿，规则修正不引入既有用例回归

### 2.5 文档漂移校准

- [x] agents.md：`ts/` 引用删改或注明"参考实现已随 `896d07e` 移除"；"待进入实现阶段"改为现状；Go 版本统一为 1.23+（现三方打架：agents.md `1.21+` / go.mod `1.23.5` / README `1.23+`）— [agents.md#L14-L17](file:///D:/code/personal/coder-mate-core/agents.md#L14-L17)、[L36](file:///D:/code/personal/coder-mate-core/agents.md#L36)
- [x] specs 01/03/06 中 `ts/lib/` 路径引用同步标注
- [x] README 缓存卖点改"避免重复解析（re-parse）"（实测 CacheHit 165,755ns ≈ LargeJS 全流程 165,096ns，"return instantly"文案与自家基准矛盾）— [README.md#L17](file:///D:/code/personal/coder-mate-core/README.md#L17)
- [x] README `UpdateTheme` 示例补 error 处理

---

## 阶段 3：测试补强（依赖 1.3 确定性修复完成）

- [x] **TS 一致性 golden 验证落地**（注意与 2.4 区分：2.4 是内联期望 token 流，此处是跨仓库 TS 对照的 golden 文件基建）：spec 08 §2 承诺的"每语言 3-5 片段 JSON diff"目前零落地（仓库无 testdata/golden 文件），而"与 TS 版输出一致"是 agents.md 第一条核心约束 — [.agents/docs/specs/08-testing-strategy.md](file:///D:/code/personal/coder-mate-core/.agents/docs/specs/08-testing-strategy.md)（2026-10-09 完成：testdata/golden/ 51 基准 + dump.ts 生成脚本 + knownDivergences 分歧清单，24 例严格相等 / 27 例已归因分歧；顺带修正 TS catch-all 与 extends 两处移植走样）
- [x] 断言升级：16 个语言包的 `checkScopeExists`"存在即过"式弱断言 → 关键用例逐 token 校验（对齐 javascript 包 33 个精确断言的标准），覆盖误分类检测（行为锁定法，共 254 token）
- [x] 淘汰形式主义测试：如 [core_test.go#L559-L581](file:///D:/code/personal/coder-mate-core/core/core_test.go#L559-L581)（建 spec 后 `_ = spec`、对输出零断言）（改为注入式 adapter 验证分发；TestThemeType 已删）
- [x] panic 路径用例：无效 theme、空语言 ID（`ensureValidID` panic 面）当前零覆盖（无效 theme 回退已在阶段 1 覆盖；新增空 ID panic 用例）
- [x] 测试卫生：`ClearRegistry`/`ResetBuiltins`/`RegisterBuiltInThemes` 在 init 与各测试反复操作全局注册表导致无法并行——梳理后启用 `t.Parallel()`（17 包 202 个测试启用；注册表操作类测试保持串行）
- [x] 验收：golden 进 CI；原弱断言包各至少 1 条精确断言；测试可并行执行

---

## 阶段 4：P3 择机

- [ ] 抽象去重：`core.NewSimpleAdapter(id, aliases, spec)` 构造器 + 共享规则工厂（注释块收尾/字符串四件套），17 个语言包只留规则数据 — [language/*/adapter.go](file:///D:/code/personal/coder-mate-core/language/javascript/adapter.go)（17 份同构）
- [ ] TS 规则插入解耦：去 `insertPos := 3` 魔法数与硬编码 `"global"`，按锚定规则定位插入点并加断言锁 JS 规则顺序 — [typescript/rules.go#L60-L69](file:///D:/code/personal/coder-mate-core/language/typescript/rules.go#L60-L69)
- [ ] 死代码清理：`Parse` 不可达 `token == nil` 分支与恒真守卫；`MatchToken` 改调 `PushState`/`PopState` 消除双实现（pop 到底语义需统一：建议静默，与 TS 版对齐）— [tokenizer.go#L211-L215](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L211-L215)、[L17-L35](file:///D:/code/personal/coder-mate-core/core/tokenizer.go#L17-L35)
- [ ] API 细节：`interface{}` → `any`；`NewHighlighter` 直接收 `HighlighterOptions` 并连同 1.1 备选方案 B 的 `(*Highlighter, error)` 返回值**一并做签名变更**（影响面清单见 1.1，变更时逐项同步）；去冗余 `ResolveAlias` 双重解析；`ensureValidID` panic 改返回 error — [api.go#L16](file:///D:/code/personal/coder-mate-core/api.go#L16)、[L44](file:///D:/code/personal/coder-mate-core/api.go#L44)、[core/registry.go#L21-L27](file:///D:/code/personal/coder-mate-core/core/registry.go#L21-L27)
- [ ] 命名惯例：`GRAMMAR_RULES`/`THEME_ALIAS_MAP`/`copy_` 改 Go 风格；`ResolveScopeStyle` 参数名 `theme` 遮蔽包名；spec 导出模式统一（16 包可变导出 map vs TS `GetSpec()` vs JS 深拷贝三套并存，防篡改防护对齐）— [javascript/rules.go#L11-L20](file:///D:/code/personal/coder-mate-core/language/javascript/rules.go#L11-L20)
- [ ] 性能：缓存 LRU 上限或超长代码不缓存（现整段源码做 key 无淘汰）— [highlighter.go#L33-L34](file:///D:/code/personal/coder-mate-core/highlighter.go#L33-L34)；`fmt.Sprintf("%d")` → `strconv.Itoa`；空 style 跳过 `<span style="">`；`escapeHtml` 冗余转义（`` ` ``/`$`）清理
- [ ] （可选）缓存真收益：同 key 命中缓存渲染后 HTML 字符串而非 token 流（当前缓存只省 parse，clone + 重渲染占主导）
- [ ] 评审 nit（2026-10-09 复审补充）：① `CodeToHtml` 的 theme 快照与 `CodeToTokens` styled rows 间有极小 TOCTOU 窗口（UpdateTheme 恰好插入两步之间时混用旧样式 + 新 preStyle），纯外观非 race，可在快照后加版本号校验或接受现状；② Python 数字 `\d[\d_]*` 会过匹配尾下划线（`1_` 不合法）→ 改 `\d(?:[\d_]*\d)?`；③ Java `EntityFunction` 要求 `(` 紧跟标识符（`Foo (` 欠匹配仍判类型）→ 允许可选空白

---

## 验收总口径

- [ ] 探针 1/2/3 全部转绿：README 用法不 panic、中文/emoji 不乱码、同输入输出唯一确定
- [ ] `go vet` / `gofmt -l` / `go test ./...` / `go test -race ./...` 全绿，且前三件 + TS 一致性 golden（阶段 3 落地后）进 CI 门禁
- [ ] 断言覆盖：JS/Python 规则修正项与多字节用例（内联期望 token 流）、TS 一致性片段（golden 文件）
- [ ] 红线（静默吞错）清零：全库不再有丢弃 error 不上抛的调用点

> 阶段 1 全部为独立小改，合计改动量小但消除全部 P1 与红线；阶段 2.3（gofmt + CI）必须先于任何多人协作或长期维护动作，否则将持续产生伪 diff 且无回归拦截。
