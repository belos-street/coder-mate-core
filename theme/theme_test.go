package theme

import (
	"testing"
)

func TestThemeType(t *testing.T) {
	theme := &HighlightTheme{
		ID:           "test",
		DisplayName:  "Test Theme",
		DefaultStyle: "color: #000;",
		PreStyle:     "background: #FFF;",
		Styles: ThemeStyleMap{
			"keyword": "color: #00F;",
		},
	}

	if theme.ID != "test" {
		t.Errorf("expected ID='test', got %q", theme.ID)
	}
	if theme.DefaultStyle != "color: #000;" {
		t.Errorf("expected DefaultStyle='color: #000;', got %q", theme.DefaultStyle)
	}
	if theme.Styles["keyword"] != "color: #00F;" {
		t.Errorf("expected keyword style='color: #00F;', got %q", theme.Styles["keyword"])
	}
}

func TestRegisterTheme(t *testing.T) {
	ClearRegistry()

	theme := &HighlightTheme{
		ID:           "my-theme",
		DisplayName:  "My Theme",
		DefaultStyle: "color: #000;",
		Styles:       ThemeStyleMap{"default": "color: #000;"},
	}

	err := RegisterTheme(theme)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := GetTheme("my-theme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.ID != "my-theme" {
		t.Errorf("expected ID='my-theme', got %q", found.ID)
	}
}

func TestRegisterTheme_EmptyID(t *testing.T) {
	ClearRegistry()

	err := RegisterTheme(&HighlightTheme{ID: "", Styles: ThemeStyleMap{}})
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestRegisterTheme_Duplicate(t *testing.T) {
	ClearRegistry()

	theme1 := &HighlightTheme{ID: "dup", Styles: ThemeStyleMap{}}
	theme2 := &HighlightTheme{ID: "dup", Styles: ThemeStyleMap{}}

	RegisterTheme(theme1)
	err := RegisterTheme(theme2)
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestGetTheme_NotFound(t *testing.T) {
	ClearRegistry()

	_, err := GetTheme("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent theme")
	}
}

func TestGetTheme_Default(t *testing.T) {
	ClearRegistry()
	RegisterTheme(&HighlightTheme{
		ID:           DefaultThemeID,
		DisplayName:  "Test Default",
		DefaultStyle: "color: #FFF;",
		Styles:       ThemeStyleMap{"default": "color: #FFF;"},
	})

	theme, err := GetTheme("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if theme.ID != DefaultThemeID {
		t.Errorf("expected default theme ID=%q, got %q", DefaultThemeID, theme.ID)
	}
}

func TestGetTheme_Alias(t *testing.T) {
	ClearRegistry()
	RegisterTheme(&HighlightTheme{
		ID:           "dark-plus",
		DisplayName:  "Dark+",
		DefaultStyle: "color: #D4D4D4;",
		Styles:       ThemeStyleMap{"default": "color: #D4D4D4;"},
	})

	theme, err := GetTheme("dark")
	if err != nil {
		t.Fatalf("expected to find 'dark' alias: %v", err)
	}
	if theme.ID != "dark-plus" {
		t.Errorf("expected ID='dark-plus', got %q", theme.ID)
	}

	theme2, err2 := GetTheme("  Dark  ")
	if err2 != nil {
		t.Fatalf("expected to find 'Dark' with trim: %v", err2)
	}
	if theme2.ID != "dark-plus" {
		t.Errorf("expected ID='dark-plus', got %q", theme2.ID)
	}
}

func TestGetTheme_GitHubAlias(t *testing.T) {
	ClearRegistry()
	RegisterTheme(&HighlightTheme{
		ID:           "github-light",
		DisplayName:  "GitHub Light",
		DefaultStyle: "color: #24292F;",
		Styles:       ThemeStyleMap{"default": "color: #24292F;"},
	})

	theme, err := GetTheme("github")
	if err != nil {
		t.Fatalf("expected to find 'github' alias: %v", err)
	}
	if theme.ID != "github-light" {
		t.Errorf("expected ID='github-light', got %q", theme.ID)
	}
}

func TestListThemes(t *testing.T) {
	ClearRegistry()

	RegisterTheme(&HighlightTheme{ID: "theme-a", Styles: ThemeStyleMap{}})
	RegisterTheme(&HighlightTheme{ID: "theme-b", Styles: ThemeStyleMap{}})

	list := ListThemes()
	if len(list) != 2 {
		t.Fatalf("expected 2 themes, got %d", len(list))
	}
}

func TestResolveTheme_String(t *testing.T) {
	ClearRegistry()
	RegisterTheme(&HighlightTheme{
		ID:           "test",
		DisplayName:  "Test",
		DefaultStyle: "color: #000;",
		Styles:       ThemeStyleMap{"default": "color: #000;"},
	})

	resolved, err := ResolveTheme("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.ID != "test" {
		t.Errorf("expected ID='test', got %q", resolved.ID)
	}
}

func TestResolveTheme_Pointer(t *testing.T) {
	theme := &HighlightTheme{ID: "direct", Styles: ThemeStyleMap{}}
	resolved, err := ResolveTheme(theme)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != theme {
		t.Error("expected same pointer back")
	}
}

func TestResolveTheme_Nil(t *testing.T) {
	ClearRegistry()
	RegisterTheme(&HighlightTheme{
		ID:           DefaultThemeID,
		DisplayName:  "Default",
		DefaultStyle: "color: #FFF;",
		Styles:       ThemeStyleMap{"default": "color: #FFF;"},
	})

	resolved, err := ResolveTheme(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.ID != DefaultThemeID {
		t.Errorf("expected default theme, got %q", resolved.ID)
	}
}

func TestResolveTheme_InvalidType(t *testing.T) {
	_, err := ResolveTheme(42)
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
}

func TestResolveScopeStyle_Exact(t *testing.T) {
	theme := &HighlightTheme{
		ID:           "test",
		DisplayName:  "Test",
		DefaultStyle: "color: #000;",
		Styles: ThemeStyleMap{
			"comment.line.js":  "color: #6A9955;",
			"comment.block.js": "color: #6A9955; font-style: italic;",
		},
	}

	style := ResolveScopeStyle("comment.line.js", theme)
	if style != "color: #6A9955;" {
		t.Errorf("expected exact match, got %q", style)
	}
}

func TestResolveScopeStyle_PrefixFallback(t *testing.T) {
	theme := &HighlightTheme{
		ID:           "test",
		DisplayName:  "Test",
		DefaultStyle: "color: #000;",
		Styles: ThemeStyleMap{
			"keyword":         "color: #111;",
			"keyword.control": "color: #222;",
		},
	}

	style := ResolveScopeStyle("keyword.control.import.js", theme)
	if style != "color: #222;" {
		t.Errorf("expected keyword.control match, got %q", style)
	}

	style = ResolveScopeStyle("keyword.declaration.let.js", theme)
	if style != "color: #111;" {
		t.Errorf("expected keyword match, got %q", style)
	}
}

func TestResolveScopeStyle_FallbackToDefaultStyle(t *testing.T) {
	theme := &HighlightTheme{
		ID:           "test",
		DisplayName:  "Test",
		DefaultStyle: "color: #999;",
		Styles:       ThemeStyleMap{},
	}

	style := ResolveScopeStyle("completely.unknown.scope", theme)
	if style != "color: #999;" {
		t.Errorf("expected defaultStyle, got %q", style)
	}
}

func TestResolveScopeStyle_FallbackToDefaultDotted(t *testing.T) {
	theme := &HighlightTheme{
		ID:           "test",
		DisplayName:  "Test",
		DefaultStyle: "color: #999;",
		Styles: ThemeStyleMap{
			"default": "color: #666;",
		},
	}

	style := ResolveScopeStyle("completely.unknown", theme)
	if style != "color: #666;" {
		t.Errorf("expected 'default' scope style, got %q", style)
	}
}

func TestCreateThemeVariant_ColorReplacement(t *testing.T) {
	base := &HighlightTheme{
		ID:           "base",
		DisplayName:  "Base",
		DefaultStyle: "color: #D4D4D4;",
		Styles: ThemeStyleMap{
			"keyword": "color: #569CD6;",
			"comment": "color: #6A9955;",
			"default": "color: #D4D4D4;",
		},
	}

	variant := CreateThemeVariant(base, VariantOptions{
		ID:           "variant",
		DisplayName:  "Variant",
		DefaultStyle: "color: #FF0000;",
		ColorMap: map[string]string{
			"#569CD6": "#FF0000",
			"#6A9955": "#00FF00",
			"#D4D4D4": "#0000FF",
		},
	})

	if variant.ID != "variant" {
		t.Errorf("expected ID='variant', got %q", variant.ID)
	}
	if variant.Styles["keyword"] != "color: #FF0000;" {
		t.Errorf("expected keyword='color: #FF0000;', got %q", variant.Styles["keyword"])
	}
	if variant.Styles["comment"] != "color: #00FF00;" {
		t.Errorf("expected comment='color: #00FF00;', got %q", variant.Styles["comment"])
	}
	if variant.Styles["default"] != "color: #0000FF;" {
		t.Errorf("expected default='color: #0000FF;', got %q", variant.Styles["default"])
	}
}

func TestCreateThemeVariant_UnmappedColorPreserved(t *testing.T) {
	base := &HighlightTheme{
		ID: "base",
		Styles: ThemeStyleMap{
			"keyword": "color: #569CD6; font-weight: bold;",
		},
	}

	variant := CreateThemeVariant(base, VariantOptions{
		ID:           "variant",
		DisplayName:  "Variant",
		DefaultStyle: "color: #FFF;",
		ColorMap:     map[string]string{},
	})

	if variant.Styles["keyword"] != "color: #569CD6; font-weight: bold;" {
		t.Errorf("expected original style preserved, got %q", variant.Styles["keyword"])
	}
}

func TestCreateThemeVariant_CaseInsensitiveColor(t *testing.T) {
	base := &HighlightTheme{
		ID: "base",
		Styles: ThemeStyleMap{
			"keyword": "color: #569cd6;",
		},
	}

	variant := CreateThemeVariant(base, VariantOptions{
		ID:           "variant",
		DisplayName:  "Variant",
		DefaultStyle: "color: #FFF;",
		ColorMap: map[string]string{
			"#569CD6": "#FF0000",
		},
	})

	if variant.Styles["keyword"] != "color: #FF0000;" {
		t.Errorf("expected case-insensitive color match, got %q", variant.Styles["keyword"])
	}
}

func TestBuiltInThemes(t *testing.T) {
	ClearRegistry()
	RegisterBuiltInThemes()

	themeIDs := []string{
		"dark-plus", "github-light", "dracula", "one-dark-pro",
		"nord", "monokai", "material-ocean", "tokyo-night",
		"solarized-dark", "solarized-light",
	}

	for _, id := range themeIDs {
		theme, err := GetTheme(id)
		if err != nil {
			t.Errorf("expected theme %q to be registered, got error: %v", id, err)
			continue
		}
		if theme.ID != id {
			t.Errorf("theme %q has wrong ID: %q", id, theme.ID)
		}
		if theme.DefaultStyle == "" {
			t.Errorf("theme %q has empty DefaultStyle", id)
		}
		if len(theme.Styles) == 0 {
			t.Errorf("theme %q has empty Styles", id)
		}
	}

	list := ListThemes()
	if len(list) != 10 {
		t.Errorf("expected 10 themes, got %d", len(list))
	}
}

func TestDarkPlusTheme(t *testing.T) {
	ClearRegistry()
	RegisterBuiltInThemes()

	theme, err := GetTheme("dark-plus")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if theme.ID != "dark-plus" {
		t.Errorf("expected ID='dark-plus', got %q", theme.ID)
	}
	if theme.DisplayName != "Dark+" {
		t.Errorf("expected DisplayName='Dark+', got %q", theme.DisplayName)
	}

	if theme.Styles["comment.block.js"] != "color: #6A9955; font-style: italic;" {
		t.Errorf("unexpected comment.block.js style: %q", theme.Styles["comment.block.js"])
	}
	if theme.Styles["keyword.control.js"] != "color: #569CD6;" {
		t.Errorf("unexpected keyword.control.js style: %q", theme.Styles["keyword.control.js"])
	}
	if theme.Styles["string.quoted.double.js"] != "color: #CE9178;" {
		t.Errorf("unexpected string.quoted.double.js style: %q", theme.Styles["string.quoted.double.js"])
	}
	if theme.Styles["default"] != "color: #D4D4D4;" {
		t.Errorf("unexpected default style: %q", theme.Styles["default"])
	}
}

func TestGitHubLightThemeVariation(t *testing.T) {
	ClearRegistry()
	RegisterBuiltInThemes()

	theme, err := GetTheme("github-light")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if theme.Styles["keyword.control.js"] == "color: #569CD6;" {
		t.Error("github-light theme should have different colors from dark-plus")
	}
}

func TestDraculaThemeVariation(t *testing.T) {
	ClearRegistry()
	RegisterBuiltInThemes()

	theme, err := GetTheme("dracula")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if theme.ID != "dracula" {
		t.Errorf("expected ID='dracula', got %q", theme.ID)
	}
	if theme.Styles["comment.block.js"] != "color: #6272A4; font-style: italic;" {
		t.Errorf("unexpected dracula comment style: %q", theme.Styles["comment.block.js"])
	}
}
