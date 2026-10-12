package codematecore

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"code-mate-core/core"
	"code-mate-core/language"
	"code-mate-core/theme"
)

const (
	// maxCacheEntries 缓存条目上限：超限时整体重置，保证服务端对任意
	// 用户输入调用时内存有界（简单有界策略，非严格 LRU）
	maxCacheEntries = 1024
	// maxCacheableCodeLen 超长代码不入缓存：大文件重复高亮的缓存价值低，
	// 却占用大量内存
	maxCacheableCodeLen = 64 * 1024
)

type cachedEntry struct {
	baseRows   core.TokenStream
	styledRows core.TokenStream
}

type highlighter struct {
	// mu 保护 cache / parsedStyleCache / currentTheme；defaultPreStyle 与
	// lineClassPrefix 构造后只读，无需加锁
	mu               sync.Mutex
	currentTheme     *theme.HighlightTheme
	defaultPreStyle  string
	lineClassPrefix  string
	cache            map[string]*cachedEntry
	parsedStyleCache map[string]core.TokenStyle
}

func normalizeLanguageID(lang string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(lang))
	if normalized == "" {
		return "", fmt.Errorf(`option "lang" cannot be empty`)
	}
	return normalized, nil
}

func createCacheKey(langID, code string) string {
	return langID + "\x00" + code
}

func escapeHtml(text string) string {
	var b strings.Builder
	b.Grow(len(text) + 16)
	for _, r := range text {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&#39;")
		case '\t':
			b.WriteString("&#9;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// escapeHtmlAttr 转义 HTML 属性值：防止含引号的选项值（如 preStyle /
// lineClassPrefix）破出属性边界形成注入面。
func escapeHtmlAttr(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 8)
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&#39;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func parseInlineStyle(styleText string) core.TokenStyle {
	style := make(core.TokenStyle)
	declarations := strings.Split(styleText, ";")

	for _, declaration := range declarations {
		trimmed := strings.TrimSpace(declaration)
		if trimmed == "" {
			continue
		}

		separatorIndex := strings.Index(trimmed, ":")
		if separatorIndex == -1 {
			continue
		}

		property := strings.TrimSpace(trimmed[:separatorIndex])
		value := strings.TrimSpace(trimmed[separatorIndex+1:])

		if property == "" || value == "" {
			continue
		}
		style[property] = value
	}

	return style
}

func stringifyInlineStyle(style core.TokenStyle) string {
	if len(style) == 0 {
		return ""
	}

	// 按 key 排序序列化，保证相同 style 的输出逐字节确定（Go map 遍历顺序随机）
	properties := make([]string, 0, len(style))
	for property := range style {
		properties = append(properties, property)
	}
	sort.Strings(properties)

	var b strings.Builder
	for _, property := range properties {
		b.WriteString(property)
		b.WriteString(": ")
		b.WriteString(style[property])
		b.WriteString("; ")
	}
	result := b.String()
	if len(result) >= 2 {
		result = result[:len(result)-1]
	}
	return result
}

func cloneTokenStream(rows core.TokenStream) core.TokenStream {
	cloned := make(core.TokenStream, len(rows))
	for i, rowTokens := range rows {
		cloned[i] = make([]core.Token, len(rowTokens))
		for j, token := range rowTokens {
			newToken := token
			newToken.Col = token.Col
			newStyle := make(core.TokenStyle, len(token.Style))
			for k, v := range token.Style {
				newStyle[k] = v
			}
			newToken.Style = newStyle
			cloned[i][j] = newToken
		}
	}
	return cloned
}

func applyThemeStyles(rows core.TokenStream, t *theme.HighlightTheme, styleCache map[string]core.TokenStyle) core.TokenStream {
	result := make(core.TokenStream, len(rows))
	for i, rowTokens := range rows {
		result[i] = make([]core.Token, len(rowTokens))
		for j, token := range rowTokens {
			styleText := theme.ResolveScopeStyle(token.Scope, t)

			cached, ok := styleCache[styleText]
			if !ok {
				cached = parseInlineStyle(styleText)
				styleCache[styleText] = cached
			}

			newStyle := make(core.TokenStyle, len(cached))
			for k, v := range cached {
				newStyle[k] = v
			}

			newToken := token
			newToken.Style = newStyle
			result[i][j] = newToken
		}
	}
	return result
}

func (h *highlighter) CodeToTokens(code, lang string) (core.TokenStream, error) {
	langID, err := normalizeLanguageID(lang)
	if err != nil {
		return nil, err
	}

	cacheKey := createCacheKey(langID, code)

	h.mu.Lock()
	if cached, ok := h.cache[cacheKey]; ok {
		rows := cloneTokenStream(cached.styledRows)
		h.mu.Unlock()
		return rows, nil
	}
	h.mu.Unlock()

	// 解析在锁外进行：并发未命中同一 key 时可能重复解析，但结果幂等
	baseRows, err := language.Tokenize(code, langID)
	if err != nil {
		return nil, err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	styledRows := applyThemeStyles(baseRows, h.currentTheme, h.parsedStyleCache)

	// 有界缓存：超长代码不入缓存；条目超限时整体重置
	if len(code) <= maxCacheableCodeLen {
		if len(h.cache) >= maxCacheEntries {
			h.cache = make(map[string]*cachedEntry)
		}
		h.cache[cacheKey] = &cachedEntry{
			baseRows:   baseRows,
			styledRows: styledRows,
		}
	}

	return cloneTokenStream(styledRows), nil
}

func (h *highlighter) CodeToHtml(code, lang string, opts ...HtmlOption) (string, error) {
	rows, err := h.CodeToTokens(code, lang)
	if err != nil {
		return "", err
	}

	o := &htmlOptions{}
	for _, opt := range opts {
		opt(o)
	}

	// currentTheme 可能被并发的 UpdateTheme 修改，先在锁内取快照
	h.mu.Lock()
	currentTheme := h.currentTheme
	h.mu.Unlock()

	preStyle := h.defaultPreStyle
	if o.PreStyle != "" {
		preStyle = o.PreStyle
	} else if preStyle == "" && currentTheme != nil && currentTheme.PreStyle != "" {
		preStyle = currentTheme.PreStyle
	} else if preStyle == "" {
		preStyle = DefaultPreStyle
	}

	lineClassPrefix := h.lineClassPrefix
	if o.LineClassPrefix != "" {
		lineClassPrefix = o.LineClassPrefix
	} else if lineClassPrefix == "" {
		lineClassPrefix = DefaultLineClassPrefix
	}

	var rowsHtml strings.Builder
	for rowIndex, rowTokens := range rows {
		rowsHtml.WriteString(`<div class="code-line `)
		rowsHtml.WriteString(escapeHtmlAttr(lineClassPrefix))
		rowsHtml.WriteString(strconv.Itoa(rowIndex + 1))
		rowsHtml.WriteString(`">`)

		for _, token := range rowTokens {
			style := stringifyInlineStyle(token.Style)
			if style == "" {
				// 无样式 token 不输出空 span，减少 HTML 体积
				rowsHtml.WriteString(escapeHtml(token.Text))
				continue
			}
			rowsHtml.WriteString(`<span style="`)
			rowsHtml.WriteString(style)
			rowsHtml.WriteString(`">`)
			rowsHtml.WriteString(escapeHtml(token.Text))
			rowsHtml.WriteString(`</span>`)
		}

		rowsHtml.WriteString(`</div>`)
	}

	return `<pre style="` + escapeHtmlAttr(preStyle) + `"><code>` + rowsHtml.String() + `</code></pre>`, nil
}

func (h *highlighter) UpdateTheme(t any) error {
	newTheme, err := theme.ResolveTheme(t)
	if err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.currentTheme = newTheme

	for cacheKey, entry := range h.cache {
		h.cache[cacheKey] = &cachedEntry{
			baseRows:   entry.baseRows,
			styledRows: applyThemeStyles(entry.baseRows, newTheme, h.parsedStyleCache),
		}
	}

	return nil
}
