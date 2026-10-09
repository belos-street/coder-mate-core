package theme

import (
	"fmt"

	"code-mate-core/theme/darkplus"
)

const defaultFontFamily = "'Consolas', 'Monaco', monospace"

func buildPreStyle(background, foreground, border, fontFamily string) string {
	if fontFamily == "" {
		fontFamily = defaultFontFamily
	}
	borderStyle := ""
	if border != "" {
		borderStyle = "border: 1px solid " + border + "; "
	}
	return "background: " + background + "; color: " + foreground + "; " +
		borderStyle +
		"padding: 16px; border-radius: 8px; font-family: " +
		fontFamily + "; font-size: 14px; line-height: 1.5; white-space: pre;"
}

func mapFromDarkPlusPalette(palette VariantPalette) map[string]string {
	return map[string]string{
		"#D4D4D4": palette.Foreground,
		"#6A9955": palette.Comment,
		"#569CD6": palette.Keyword,
		"#C586C0": palette.Accent,
		"#4EC9B0": palette.TypeName,
		"#DCDCAA": palette.Callable,
		"#D16969": palette.Warning,
		"#B5CEA8": palette.Number,
		"#D7BA7D": palette.Symbol,
		"#9CDCFE": palette.Variable,
		"#CE9178": palette.String,
		"#808080": palette.Muted,
	}
}

type VariantPalette struct {
	Foreground string
	Comment    string
	Keyword    string
	Accent     string
	TypeName   string
	Callable   string
	Warning    string
	Number     string
	Symbol     string
	Variable   string
	String     string
	Muted      string
}

type VariantDefinition struct {
	ID          string
	DisplayName string
	Palette     VariantPalette
	Background  string
	Border      string
	FontFamily  string
}

func createPresetTheme(base *HighlightTheme, definition VariantDefinition) *HighlightTheme {
	colorMap := mapFromDarkPlusPalette(definition.Palette)

	return CreateThemeVariant(base, VariantOptions{
		ID:           definition.ID,
		DisplayName:  definition.DisplayName,
		DefaultStyle: "color: " + definition.Palette.Foreground + ";",
		PreStyle: buildPreStyle(
			definition.Background,
			definition.Palette.Foreground,
			definition.Border,
			definition.FontFamily,
		),
		ColorMap: colorMap,
	})
}

var presetDefinitions = []VariantDefinition{
	{
		ID:          "github-light",
		DisplayName: "GitHub Light",
		Background:  "#FFFFFF",
		Border:      "#D0D7DE",
		FontFamily:  "'SFMono-Regular', 'Consolas', 'Monaco', monospace",
		Palette: VariantPalette{
			Foreground: "#24292F",
			Comment:    "#6A737D",
			Keyword:    "#D73A49",
			Accent:     "#6F42C1",
			TypeName:   "#005CC5",
			Callable:   "#6F42C1",
			Warning:    "#D73A49",
			Number:     "#005CC5",
			Symbol:     "#B08800",
			Variable:   "#24292F",
			String:     "#032F62",
			Muted:      "#24292F",
		},
	},
	{
		ID:          "dracula",
		DisplayName: "Dracula",
		Background:  "#282A36",
		Palette: VariantPalette{
			Foreground: "#F8F8F2",
			Comment:    "#6272A4",
			Keyword:    "#FF79C6",
			Accent:     "#BD93F9",
			TypeName:   "#8BE9FD",
			Callable:   "#F1FA8C",
			Warning:    "#FF5555",
			Number:     "#50FA7B",
			Symbol:     "#FFB86C",
			Variable:   "#8BE9FD",
			String:     "#F1FA8C",
			Muted:      "#6272A4",
		},
	},
	{
		ID:          "one-dark-pro",
		DisplayName: "One Dark Pro",
		Background:  "#282C34",
		Palette: VariantPalette{
			Foreground: "#ABB2BF",
			Comment:    "#5C6370",
			Keyword:    "#C678DD",
			Accent:     "#C678DD",
			TypeName:   "#56B6C2",
			Callable:   "#E5C07B",
			Warning:    "#E06C75",
			Number:     "#98C379",
			Symbol:     "#D19A66",
			Variable:   "#61AFEF",
			String:     "#98C379",
			Muted:      "#7F848E",
		},
	},
	{
		ID:          "nord",
		DisplayName: "Nord",
		Background:  "#2E3440",
		Palette: VariantPalette{
			Foreground: "#D8DEE9",
			Comment:    "#616E88",
			Keyword:    "#81A1C1",
			Accent:     "#B48EAD",
			TypeName:   "#88C0D0",
			Callable:   "#EBCB8B",
			Warning:    "#BF616A",
			Number:     "#A3BE8C",
			Symbol:     "#D08770",
			Variable:   "#8FBCBB",
			String:     "#A3BE8C",
			Muted:      "#4C566A",
		},
	},
	{
		ID:          "monokai",
		DisplayName: "Monokai",
		Background:  "#272822",
		Palette: VariantPalette{
			Foreground: "#F8F8F2",
			Comment:    "#75715E",
			Keyword:    "#F92672",
			Accent:     "#AE81FF",
			TypeName:   "#66D9EF",
			Callable:   "#E6DB74",
			Warning:    "#F92672",
			Number:     "#AE81FF",
			Symbol:     "#FD971F",
			Variable:   "#A6E22E",
			String:     "#E6DB74",
			Muted:      "#75715E",
		},
	},
	{
		ID:          "material-ocean",
		DisplayName: "Material Ocean",
		Background:  "#0F111A",
		Palette: VariantPalette{
			Foreground: "#A6ACCD",
			Comment:    "#546E7A",
			Keyword:    "#C792EA",
			Accent:     "#C792EA",
			TypeName:   "#89DDFF",
			Callable:   "#FFCB6B",
			Warning:    "#F07178",
			Number:     "#C3E88D",
			Symbol:     "#F78C6C",
			Variable:   "#82AAFF",
			String:     "#C3E88D",
			Muted:      "#717CB4",
		},
	},
	{
		ID:          "tokyo-night",
		DisplayName: "Tokyo Night",
		Background:  "#1A1B26",
		Palette: VariantPalette{
			Foreground: "#C0CAF5",
			Comment:    "#565F89",
			Keyword:    "#7AA2F7",
			Accent:     "#BB9AF7",
			TypeName:   "#7DCFFF",
			Callable:   "#E0AF68",
			Warning:    "#F7768E",
			Number:     "#9ECE6A",
			Symbol:     "#FF9E64",
			Variable:   "#73DACA",
			String:     "#9ECE6A",
			Muted:      "#414868",
		},
	},
	{
		ID:          "solarized-dark",
		DisplayName: "Solarized Dark",
		Background:  "#002B36",
		Palette: VariantPalette{
			Foreground: "#93A1A1",
			Comment:    "#586E75",
			Keyword:    "#268BD2",
			Accent:     "#6C71C4",
			TypeName:   "#2AA198",
			Callable:   "#B58900",
			Warning:    "#DC322F",
			Number:     "#859900",
			Symbol:     "#CB4B16",
			Variable:   "#268BD2",
			String:     "#859900",
			Muted:      "#657B83",
		},
	},
	{
		ID:          "solarized-light",
		DisplayName: "Solarized Light",
		Background:  "#FDF6E3",
		Border:      "#EEE8D5",
		Palette: VariantPalette{
			Foreground: "#657B83",
			Comment:    "#93A1A1",
			Keyword:    "#268BD2",
			Accent:     "#6C71C4",
			TypeName:   "#2AA198",
			Callable:   "#B58900",
			Warning:    "#DC322F",
			Number:     "#859900",
			Symbol:     "#CB4B16",
			Variable:   "#268BD2",
			String:     "#859900",
			Muted:      "#93A1A1",
		},
	},
}

var DarkPlusTheme = &HighlightTheme{
	ID:           darkplus.ID,
	DisplayName:  darkplus.DisplayName,
	DefaultStyle: darkplus.DefaultStyle,
	PreStyle:     darkplus.PreStyle,
	Styles:       darkplus.MergeStyles(),
}

var GithubLightTheme = createPresetTheme(DarkPlusTheme, presetDefinitions[0])
var DraculaTheme = createPresetTheme(DarkPlusTheme, presetDefinitions[1])
var OneDarkProTheme = createPresetTheme(DarkPlusTheme, presetDefinitions[2])
var NordTheme = createPresetTheme(DarkPlusTheme, presetDefinitions[3])
var MonokaiTheme = createPresetTheme(DarkPlusTheme, presetDefinitions[4])
var MaterialOceanTheme = createPresetTheme(DarkPlusTheme, presetDefinitions[5])
var TokyoNightTheme = createPresetTheme(DarkPlusTheme, presetDefinitions[6])
var SolarizedDarkTheme = createPresetTheme(DarkPlusTheme, presetDefinitions[7])
var SolarizedLightTheme = createPresetTheme(DarkPlusTheme, presetDefinitions[8])

var presetBuiltInThemes = []*HighlightTheme{
	GithubLightTheme,
	DraculaTheme,
	OneDarkProTheme,
	NordTheme,
	MonokaiTheme,
	MaterialOceanTheme,
	TokyoNightTheme,
	SolarizedDarkTheme,
	SolarizedLightTheme,
}

// EnsureBuiltinsRegistered 幂等注册内置主题：默认主题已存在时跳过；
// 注册表被 ClearRegistry 清空（如测试场景）后再次调用会自动重建。
func EnsureBuiltinsRegistered() {
	if _, ok := themeRegistry[DefaultThemeID]; !ok {
		RegisterBuiltInThemes()
	}
}

func RegisterBuiltInThemes() {
	mustRegisterTheme(DarkPlusTheme)
	for _, theme := range presetBuiltInThemes {
		mustRegisterTheme(theme)
	}
}

func mustRegisterTheme(theme *HighlightTheme) {
	if err := RegisterTheme(theme); err != nil {
		panic(fmt.Sprintf("theme: register built-in theme %q: %v", theme.ID, err))
	}
}
