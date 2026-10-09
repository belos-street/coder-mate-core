# 06 — 公共 API

> **所属 Phase**：Phase 4
> **参考源码**：`ts/lib/api.ts`

> ⚠️ 历史注记：上述 `ts/lib/` 参考源码已随提交 `896d07e` 移除，仅作历史参考。

---

## 1. Highlighter 接口

实现文件：`api.go` + `highlighter.go`

```go
type Highlighter interface {
    CodeToTokens(code string, lang string) (core.TokenStream, error)
    CodeToHtml(code string, lang string, opts ...HtmlOption) (string, error)
    UpdateTheme(theme interface{}) error
}

type HighlighterOptions struct {
    Theme           interface{} // string 或 *HighlightTheme
    PreStyle        string
    LineClassPrefix string
}

func NewHighlighter(opts ...HighlighterOptions) Highlighter

// HtmlOption Functional Option 模式
type HtmlOption func(*htmlOptions)

type htmlOptions struct {
    PreStyle        string
    LineClassPrefix string
}

func WithPreStyle(style string) HtmlOption    { ... }
func WithLineClassPrefix(prefix string) HtmlOption { ... }
```

默认值：
- `PreStyle`：`"background: #1E1E1E; padding: 16px; border-radius: 8px; font-family: 'Consolas', 'Monaco', monospace; font-size: 14px; line-height: 1.5; white-space: pre;"`
- `LineClassPrefix`：`"line-"`

## 2. CodeToTokens 实现细节

1. 规范化语言 ID（trim + toLower）
2. 查找缓存：`key = lang + "\x00" + code`
3. 若命中缓存：深拷贝 styled token 流后返回
4. 若未命中：调用 `language.Tokenize(code, lang)` 获取基础 token 流
5. 调用 `applyThemeStyles()` 应用主题样式
6. 存入缓存（baseRows + styledRows）
7. 深拷贝后返回

### 样式查找逻辑（scope 解析）

1. 精确匹配：`theme.Styles[scope]`
2. 前缀回退：`a.b.c` → 查找 `a.b` → 查找 `a`
3. 兜底：`theme.Styles["default"]` → `theme.DefaultStyle`

### 样式解析/序列化工具函数（位于 `highlighter.go` 内部）

参考 [ts/lib/api.ts](../../../ts/lib/api.ts#L72-L97)：

**`parseInlineStyle(styleText string) TokenStyle`**
- 输入：`"color: #fff; font-weight: bold;"`
- 输出：`map[string]string{"color": "#fff", "font-weight": "bold"}`
- 按 `;` 分割，每段按第一个 `:` 分割为 property/value
- 忽略空段和无效段（没有 `:` 或 property/value 为空的段）
- 必须缓存解析结果（用 `map[string]TokenStyle`），同一 styleText 不重复解析

**`stringifyInlineStyle(style TokenStyle) string`**
- 输入：`map[string]string{"color": "#fff", "font-weight": "bold"}`
- 输出：`"color: #fff; font-weight: bold;"`
- 空 map 返回空字符串 `""`
- 格式：`"{property}: {value}; "` 用空格连接

**`applyThemeStyles` 流程**：
1. 遍历所有 token，对每个 token.scope 按 scope 解析逻辑找到对应 styleText
2. 调用 `parseInlineStyle(styleText)` 解析为 TokenStyle（带缓存）
3. 将解析后的 TokenStyle 深拷贝给 `token.Style`
4. 返回新的 TokenStream（不修改原始基础 token 流）

## 3. CodeToHtml 实现细节

1. 调用 `CodeToTokens` 获取带样式的 token 流
2. `escapeHtml()` 转义特殊字符，映射必须精确一致：

| 字符 | 转义为 |
|------|--------|
| `&` | `&amp;` |
| `<` | `&lt;` |
| `>` | `&gt;` |
| `"` | `&quot;` |
| `'` | `&#39;` |
| `` ` `` | `&#96;` |
| `$` | `&#36;` |
| `\t` | `&#9;` |

3. 每个 token 包裹 `<span style="...">`（style 值来自 `stringifyInlineStyle`）
4. 每行包裹 `<div class="code-line {lineClassPrefix}{lineNumber}">`
5. 外层包裹 `<pre style="{preStyle}"><code>...</code></pre>`

### HTML 输出格式

```html
<pre style="{preStyle}">
  <code>
    <div class="code-line line-1"><span style="...">token1</span><span style="...">token2</span></div>
    <div class="code-line line-2"><span style="...">token3</span></div>
  </code>
</pre>
```

## 4. 缓存机制

- **缓存 key**：`lang + "\x00" + code`
- **缓存内容**：基础 token 流（未应用主题样式）+ 已应用主题样式的 token 流
- **主题切换**：遍历缓存的所有 entry，只重算 styledRows，不重新解析代码
- **注意**：WASM 单线程环境，不需要 mutex 保护缓存

### 深拷贝（cloneTokenStream）

`CodeToTokens` 返回给调用方前，必须对缓存中的 styled token 流执行深拷贝，防止调用方修改返回的 token 对象污染缓存：

- `Token.Col` 数组（`[2]int`）—— 必须拷贝
- `Token.Style` map（`map[string]string`）—— 必须拷贝
- `Token.Text` 和 `Token.Scope` 是 string（Go 中不可变，无需单独拷贝）
