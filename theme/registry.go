package theme

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

var DefaultThemeID = "dark-plus"

var THEME_ALIAS_MAP = map[string]string{
	"dark":   "dark-plus",
	"github": "github-light",
}

var (
	registryMu    sync.RWMutex
	themeRegistry = make(map[string]*HighlightTheme)
)

func normalizeThemeID(id string) string {
	return strings.TrimSpace(strings.ToLower(id))
}

func RegisterTheme(theme *HighlightTheme) error {
	normalizedID := normalizeThemeID(theme.ID)
	if normalizedID == "" {
		return fmt.Errorf("theme id cannot be empty")
	}

	registryMu.Lock()
	defer registryMu.Unlock()

	if existing, ok := themeRegistry[normalizedID]; ok && existing != theme {
		return fmt.Errorf("theme %q is already registered", theme.ID)
	}

	themeRegistry[normalizedID] = theme
	return nil
}

func GetTheme(themeID string) (*HighlightTheme, error) {
	if themeID == "" {
		themeID = DefaultThemeID
	}

	// 先检查别名
	if resolvedID, ok := THEME_ALIAS_MAP[normalizeThemeID(themeID)]; ok {
		themeID = resolvedID
	}

	normalizedID := normalizeThemeID(themeID)
	if normalizedID == "" {
		return nil, fmt.Errorf("theme id cannot be empty")
	}

	registryMu.RLock()
	defer registryMu.RUnlock()
	theme, ok := themeRegistry[normalizedID]
	if !ok {
		return nil, fmt.Errorf("theme %q is not registered", normalizedID)
	}

	return theme, nil
}

func ListThemes() []*HighlightTheme {
	registryMu.RLock()
	defer registryMu.RUnlock()

	themes := make([]*HighlightTheme, 0, len(themeRegistry))
	for _, theme := range themeRegistry {
		themes = append(themes, theme)
	}

	// 按 ID 排序，保证返回顺序确定（Go map 遍历顺序随机）
	sort.Slice(themes, func(i, j int) bool {
		return themes[i].ID < themes[j].ID
	})

	return themes
}

func ResolveAlias(theme interface{}) interface{} {
	if s, ok := theme.(string); ok {
		normalized := normalizeThemeID(s)
		if resolved, ok := THEME_ALIAS_MAP[normalized]; ok {
			return resolved
		}
		return normalized
	}
	return theme
}

func ResolveTheme(theme interface{}) (*HighlightTheme, error) {
	switch v := theme.(type) {
	case nil:
		return GetTheme("")
	case string:
		return GetTheme(v)
	case *HighlightTheme:
		return v, nil
	default:
		return nil, fmt.Errorf("invalid theme type: %T", theme)
	}
}

// ResolveThemeOrDefault 解析主题；未知主题名或无效类型时回退默认主题。
// 这是 NewHighlighter 的文档化行为：无效的主题配置不报错、不中断，
// 回退到 DefaultThemeID。调用前须已通过 EnsureBuiltinsRegistered 注册内置主题。
func ResolveThemeOrDefault(theme interface{}) *HighlightTheme {
	if resolved, err := ResolveTheme(theme); err == nil && resolved != nil {
		return resolved
	}

	resolved, err := ResolveTheme(nil)
	if err != nil {
		panic(fmt.Sprintf("theme: default theme %q is not registered: %v", DefaultThemeID, err))
	}
	return resolved
}

// ResolveScopeStyle 查找 scope 对应的样式：
// 1) 精确匹配
// 2) 前缀回退（a.b.c -> a.b -> a）
// 3) defaultStyle
func ResolveScopeStyle(scope string, theme *HighlightTheme) string {
	if style, ok := theme.Styles[scope]; ok {
		return style
	}

	candidate := scope
	for strings.Contains(candidate, ".") {
		idx := strings.LastIndex(candidate, ".")
		candidate = candidate[:idx]
		if style, ok := theme.Styles[candidate]; ok {
			return style
		}
	}

	if style, ok := theme.Styles["default"]; ok {
		return style
	}
	return theme.DefaultStyle
}

func ClearRegistry() {
	registryMu.Lock()
	defer registryMu.Unlock()
	themeRegistry = make(map[string]*HighlightTheme)
}
