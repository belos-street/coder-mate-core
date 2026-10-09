# 03 — 核心分词引擎

> **所属 Phase**：Phase 1
> **参考源码**：`ts/lib/core/types.ts`, `ts/lib/core/tokenizer.ts`, `ts/lib/core/registry.ts`

> ⚠️ 历史注记：上述 `ts/lib/` 参考源码已随提交 `896d07e` 移除，仅作历史参考。

---

## 1. 核心类型定义

实现文件：`core/types.go`

| TS 类型 | Go 类型 | 说明 |
|---------|---------|------|
| `Token<Scope>` | `Token` | 统一使用 `string` 类型的 Scope 字段 |
| `GrammarRule<State, Scope>` | `GrammarRule` | 包含 Regex, Scope, PushState, PopState, Skip |
| `ParserContext<State>` | `ParserContext` | 包含 StateStack ([]string), Line, Col |
| `TokenizerSpec<State, Scope>` | `TokenizerSpec` | 包含 InitialState, Rules (map[string][]GrammarRule), FallbackScope |
| `TokenStream<Scope>` | `TokenStream` | 即 `[][]Token` |
| `TokenStyle` | `TokenStyle` | 即 `map[string]string` |

**关键约束**：
- `Token.Col` 为 `[2]int`，表示 `[startCol, endCol]` 闭区间
- `GrammarRule.PushState` 为空字符串时表示不压栈
- `GrammarRule.PopState` 为 `true` 时在匹配后弹栈
- `GrammarRule.Skip` 字段保留用于未来扩展，当前实现不需要处理此字段（核心 tokenizer 不检查它）

### Go 类型定义

```go
package core

import "regexp"

type TokenStyle map[string]string

type Token struct {
    Text  string
    Scope string
    Line  int
    Col   [2]int // [startCol, endCol] 闭区间
    Style TokenStyle
}

type GrammarRule struct {
    Regex     *regexp.Regexp
    Scope     string
    PushState string // 空字符串表示不压栈
    PopState  bool
    Skip      bool   // 当前不处理，保留
}

type ParserContext struct {
    StateStack []string
    Line       int
    Col        int
}

type TokenizerSpec struct {
    InitialState  string
    Rules         map[string][]GrammarRule
    FallbackScope string
}

type TokenStream [][]Token
```

---

## 2. FSM 分词引擎

实现文件：`core/tokenizer.go`

算法逻辑必须与 [ts/lib/core/tokenizer.ts](../../../ts/lib/core/tokenizer.ts) 精确一致。

### 2.1 `NewParserContext(initialState string) *ParserContext`

创建初始上下文，StateStack 初始化为 `[initialState]`，Line=1, Col=1。

### 2.2 `PushState(ctx *ParserContext, state string) *ParserContext`

返回新上下文（不修改原对象），StateStack 追加 state。

### 2.3 `PopState(ctx *ParserContext) (*ParserContext, error)`

返回新上下文，StateStack 弹出最后一个元素。栈长度 ≤ 1 时返回 error。

### 2.4 `CurrentState(ctx *ParserContext) string`

返回 StateStack 最后一个元素。

### 2.5 `SplitTokenByLineBreak(text, tokenScope, fallbackScope string, startLine, startCol int) []Token`

- 处理 `\n` 和 `\r\n` 两种换行符
- 将跨行 token 拆分为多个单行 token
- 换行符 token 使用 fallbackScope 作为 scope
- 正确维护 line 和 col 计数

### 2.6 `MatchToken(code string, ctx *ParserContext, spec *TokenizerSpec) (*Token, *ParserContext)`

- 获取当前状态的规则列表
- **按顺序遍历规则**，用正则从代码头部匹配（`^` 锚定）
- 第一个匹配的规则胜出
- 根据规则的 PushState/PopState 修改上下文
- **所有规则都不匹配时，MatchToken 内部**消耗一个字符作为 fallback token（使用 `spec.FallbackScope`），不应该返回 nil/null
- 仅当 `code` 长度为 0（理论上不应发生）时才返回 nil/null

### 2.7 `Parse(code string, spec *TokenizerSpec) TokenStream`

主解析函数，维护 `remainingCode` 和 `context`。

每次迭代：
1. 先检查是否是换行符（`\n` 或 `\r\n`），是则推入当前行并换行
2. 如果不是换行，调用 MatchToken 尝试匹配
3. 匹配到的 token 如果包含 `\n`，调用 SplitTokenByLineBreak 拆分
4. 匹配不到则消耗一个字符作为 fallback token
5. 最后一行如果没有换行符结尾，也要推入结果

### Go 实现要点

- 使用**值拷贝**语义保持与 TS 一致的不可变性：`newCtx := *ctx`
- StateStack 的拷贝：`newStack := make([]string, len(ctx.StateStack)); copy(newStack, ctx.StateStack)`

---

## 3. 语言注册表

实现文件：`core/registry.go`

```go
type LanguageAdapter interface {
    ID() string
    Aliases() []string
    Parse(code string) TokenStream
}
```

必须实现：
- `RegisterLanguage(lang LanguageAdapter) error` — 注册语言（id + aliases 全部注册到 map），已注册则报错
- `GetLanguage(id string) (LanguageAdapter, bool)` — 查找语言（大小写不敏感）
- `ListLanguages() []LanguageAdapter` — 列出所有已注册语言（去重）
- `Tokenize(code, langID string) (TokenStream, error)` — 便捷函数

**语言 id 规范化**：`strings.TrimSpace(strings.ToLower(id))`
