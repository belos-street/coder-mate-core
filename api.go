package codematecore

import (
	"code-mate-core/core"
	"code-mate-core/language"
	"code-mate-core/theme"
)

const DefaultPreStyle = `background: #1E1E1E; padding: 16px; border-radius: 8px; font-family: 'Consolas', 'Monaco', monospace; font-size: 14px; line-height: 1.5; white-space: pre;`

const DefaultLineClassPrefix = "line-"

type Highlighter interface {
	CodeToTokens(code, lang string) (core.TokenStream, error)
	CodeToHtml(code, lang string, opts ...HtmlOption) (string, error)
	UpdateTheme(t interface{}) error
}

type HighlighterOptions struct {
	Theme           interface{}
	PreStyle        string
	LineClassPrefix string
}

type HtmlOption func(*htmlOptions)

type htmlOptions struct {
	PreStyle        string
	LineClassPrefix string
}

func WithPreStyle(style string) HtmlOption {
	return func(o *htmlOptions) {
		o.PreStyle = style
	}
}

func WithLineClassPrefix(prefix string) HtmlOption {
	return func(o *htmlOptions) {
		o.LineClassPrefix = prefix
	}
}

func NewHighlighter(opts ...HighlighterOptions) Highlighter {
	o := HighlighterOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	return newHighlighter(o)
}

func newHighlighter(opts HighlighterOptions) *highlighter {
	language.EnsureBuiltinsRegistered()
	theme.EnsureBuiltinsRegistered()

	// 未知主题名 / 无效类型回退默认主题（文档化行为，见 ResolveThemeOrDefault）
	currentTheme := theme.ResolveThemeOrDefault(theme.ResolveAlias(opts.Theme))

	return &highlighter{
		currentTheme:     currentTheme,
		defaultPreStyle:  opts.PreStyle,
		lineClassPrefix:  opts.LineClassPrefix,
		cache:            make(map[string]*cachedEntry),
		parsedStyleCache: make(map[string]core.TokenStyle),
	}
}
