package codematecore

import (
	"strings"
	"sync"
	"testing"

	"code-mate-core/core"
	"code-mate-core/language"
	"code-mate-core/theme"
)

func init() {
	theme.ClearRegistry()
	core.ClearRegistry()
	language.ResetBuiltins()
	theme.RegisterBuiltInThemes()
}

func TestCodeToTokens_ParsesRows(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: "dark-plus"})

	rows, err := hl.CodeToTokens("const x = 42;", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if len(rows[0]) == 0 {
		t.Fatal("expected at least 1 token")
	}

	foundConst := false
	for _, tok := range rows[0] {
		if tok.Text == "const" {
			foundConst = true
			if tok.Scope != "keyword.declaration.js" {
				t.Errorf("expected scope 'keyword.declaration.js' for 'const', got %q", tok.Scope)
			}
		}
	}
	if !foundConst {
		t.Error("expected to find 'const' token")
	}
}

func TestCodeToTokens_IncludesStyle(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: "github-light"})

	rows, err := hl.CodeToTokens("// note", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rows) != 1 || len(rows[0]) == 0 {
		t.Fatal("expected tokens")
	}

	first := rows[0][0]
	if first.Scope != "comment.line.double-slash.js" {
		t.Errorf("expected scope 'comment.line.double-slash.js', got %q", first.Scope)
	}
	if first.Style == nil || len(first.Style) == 0 {
		t.Fatal("expected non-empty style")
	}
	if first.Style["color"] == "" {
		t.Error("expected color in style")
	}
}

func TestUpdateTheme_RefreshesStyles(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: "dark-plus"})

	darkRows, err := hl.CodeToTokens("const value = 1", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var darkConstToken *core.Token
	for _, row := range darkRows {
		for i := range row {
			if row[i].Text == "const" {
				darkConstToken = &row[i]
				break
			}
		}
	}
	if darkConstToken == nil {
		t.Fatal("expected to find 'const' token")
	}
	if darkConstToken.Scope != "keyword.declaration.js" {
		t.Errorf("expected scope 'keyword.declaration.js', got %q", darkConstToken.Scope)
	}
	if darkConstToken.Style["color"] != "#569CD6" {
		t.Errorf("expected dark color #569CD6, got %q", darkConstToken.Style["color"])
	}

	err = hl.UpdateTheme("github-light")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lightRows, err := hl.CodeToTokens("const value = 1", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var lightConstToken *core.Token
	for _, row := range lightRows {
		for i := range row {
			if row[i].Text == "const" {
				lightConstToken = &row[i]
				break
			}
		}
	}
	if lightConstToken == nil {
		t.Fatal("expected to find 'const' token")
	}
	if lightConstToken.Scope != "keyword.declaration.js" {
		t.Errorf("expected scope 'keyword.declaration.js', got %q", lightConstToken.Scope)
	}
	if lightConstToken.Style["color"] != "#D73A49" {
		t.Errorf("expected light color #D73A49, got %q", lightConstToken.Style["color"])
	}
}

func TestCodeToHtml_RendersHtml(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: "github-light"})

	html, err := hl.CodeToHtml("const x = 1", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, "<pre") {
		t.Error("expected <pre> element in HTML")
	}
	if !strings.Contains(html, "<code>") {
		t.Error("expected <code> element in HTML")
	}
	if !strings.Contains(html, "const") {
		t.Error("expected text 'const' in HTML")
	}
	if !strings.Contains(html, `class="code-line line-1"`) {
		t.Error("expected line class in HTML")
	}
}

func TestCodeToHtml_CustomPreStyle(t *testing.T) {
	hl := NewHighlighter()
	customStyle := "background: #000; color: #FFF;"

	html, err := hl.CodeToHtml("const x = 1", "javascript",
		WithPreStyle(customStyle),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, customStyle) {
		t.Errorf("expected custom pre style %q in HTML", customStyle)
	}
}

func TestCodeToHtml_CustomLineClassPrefix(t *testing.T) {
	hl := NewHighlighter()
	customPrefix := "ln-"

	html, err := hl.CodeToHtml("const x = 1", "javascript",
		WithLineClassPrefix(customPrefix),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, `class="code-line ln-1"`) {
		t.Error("expected custom line class prefix in HTML")
	}
}

func TestCodeToTokens_UnknownLanguage(t *testing.T) {
	hl := NewHighlighter()

	_, err := hl.CodeToTokens("hello", "unknown-lang")
	if err == nil {
		t.Fatal("expected error for unknown language")
	}

	expectedMsg := `language "unknown-lang" is not registered`
	if err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %q", expectedMsg, err.Error())
	}
}

func TestCodeToTokens_EmptyLanguage(t *testing.T) {
	hl := NewHighlighter()

	_, err := hl.CodeToTokens("hello", "   ")
	if err == nil {
		t.Fatal("expected error for empty language")
	}

	expectedMsg := `option "lang" cannot be empty`
	if err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %q", expectedMsg, err.Error())
	}
}

func TestCodeToHtml_ThemeAlias(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: "github"})

	html, err := hl.CodeToHtml("const x = 1", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, "<pre") {
		t.Error("expected HTML output with github alias")
	}
}

func TestCodeToHtml_DarkAlias(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: "dark"})

	html, err := hl.CodeToHtml("const x = 1", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, "<pre") {
		t.Error("expected HTML output with dark alias")
	}
}

func TestCache_Hit(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: "dark-plus"})

	rows1, err := hl.CodeToTokens("const x = 42;", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rows2, err := hl.CodeToTokens("const x = 42;", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rows1) != len(rows2) {
		t.Errorf("cache hit should return same number of rows: %d vs %d", len(rows1), len(rows2))
	}
	if len(rows1[0]) != len(rows2[0]) {
		t.Errorf("cache hit should return same number of tokens: %d vs %d", len(rows1[0]), len(rows2[0]))
	}
}

func TestCache_Clone(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: "dark-plus"})

	rows1, err := hl.CodeToTokens("const x = 42;", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rows1) > 0 && len(rows1[0]) > 0 {
		rows1[0][0].Style["color"] = "modified"
	}

	rows2, err := hl.CodeToTokens("const x = 42;", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rows2) > 0 && len(rows2[0]) > 0 {
		if rows2[0][0].Style["color"] == "modified" {
			t.Error("cache should be protected from mutation via clone")
		}
	}
}

func TestUpdateTheme_InvalidTheme(t *testing.T) {
	hl := NewHighlighter()

	err := hl.UpdateTheme("nonexistent-theme")
	if err == nil {
		t.Fatal("expected error for nonexistent theme")
	}
}

func TestCodeToHtml_HtmlEscaping(t *testing.T) {
	hl := NewHighlighter()

	html, err := hl.CodeToHtml(`const x = "<div>";`, "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(html, "<div>") && !strings.Contains(html, "&lt;div&gt;") {
		t.Error("expected HTML entities for < and > in string content")
	}
}

func TestCodeToHtml_MultiLineOutput(t *testing.T) {
	hl := NewHighlighter()

	code := "const x = 1;\nconst y = 2;"
	html, err := hl.CodeToHtml(code, "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, `class="code-line line-1"`) {
		t.Error("expected line-1 class")
	}
	if !strings.Contains(html, `class="code-line line-2"`) {
		t.Error("expected line-2 class")
	}
}

func TestNewHighlighter_NullTheme(t *testing.T) {
	hl := NewHighlighter()

	_, err := hl.CodeToTokens("const x = 42;", "javascript")
	if err != nil {
		t.Fatalf("unexpected error with no theme: %v", err)
	}
}

func TestCodeToHtml_WithThemePreStyle(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: "dark-plus"})

	html, err := hl.CodeToHtml("const x = 1", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, "background: #1E1E1E") {
		t.Error("expected dark-plus pre style with background")
	}
}

func TestCodeToHtml_WithHighlighterPreStyle(t *testing.T) {
	customStyle := "background: #000; padding: 10px;"
	hl := NewHighlighter(HighlighterOptions{
		Theme:    "dark-plus",
		PreStyle: customStyle,
	})

	html, err := hl.CodeToHtml("const x = 1", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, customStyle) {
		t.Errorf("expected highlighter preStyle in HTML, got: %s", html)
	}
}

func TestEscapeHtml(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"<div>", "&lt;div&gt;"},
		{"a & b", "a &amp; b"},
		{`"quote"`, "&quot;quote&quot;"},
		{"'s", "&#39;s"},
		{"`code`", "&#96;code&#96;"},
		{"$var", "&#36;var"},
		{"\t", "&#9;"},
	}

	for _, tt := range tests {
		result := escapeHtml(tt.input)
		if result != tt.expected {
			t.Errorf("escapeHtml(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestParseInlineStyle(t *testing.T) {
	style := parseInlineStyle("color: #569CD6; font-weight: bold;")
	if style["color"] != "#569CD6" {
		t.Errorf("expected color=#569CD6, got %q", style["color"])
	}
	if style["font-weight"] != "bold" {
		t.Errorf("expected font-weight=bold, got %q", style["font-weight"])
	}
}

func TestParseInlineStyle_Empty(t *testing.T) {
	style := parseInlineStyle("")
	if len(style) != 0 {
		t.Errorf("expected empty style, got %d entries", len(style))
	}
}

func TestStringifyInlineStyle(t *testing.T) {
	style := core.TokenStyle{
		"font-weight": "bold",
		"color":       "#569CD6",
	}
	result := stringifyInlineStyle(style)

	// 按 key 排序序列化：输出必须逐字节确定，与 map 内部顺序无关
	expected := "color: #569CD6; font-weight: bold;"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestStringifyInlineStyle_Empty(t *testing.T) {
	result := stringifyInlineStyle(core.TokenStyle{})
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestCloneTokenStream(t *testing.T) {
	original := core.TokenStream{
		{
			{
				Text:  "hello",
				Scope: "test",
				Line:  1,
				Col:   [2]int{1, 5},
				Style: core.TokenStyle{"color": "#FFF"},
			},
		},
	}

	cloned := cloneTokenStream(original)

	original[0][0].Style["color"] = "modified"
	original[0][0].Col = [2]int{99, 100}

	if cloned[0][0].Style["color"] != "#FFF" {
		t.Error("clone should have independent style map")
	}
	if cloned[0][0].Col != [2]int{1, 5} {
		t.Error("clone should have independent Col")
	}
}

func TestNormalizeLanguageID(t *testing.T) {
	id, err := normalizeLanguageID("  JavaScript  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "javascript" {
		t.Errorf("expected 'javascript', got %q", id)
	}
}

func TestNormalizeLanguageID_Empty(t *testing.T) {
	_, err := normalizeLanguageID("   ")
	if err == nil {
		t.Fatal("expected error for empty language ID")
	}
}

func TestCreateCacheKey(t *testing.T) {
	key := createCacheKey("javascript", "const x = 1;")
	if !strings.Contains(key, "\x00") {
		t.Error("cache key should contain null separator")
	}
}

func TestNewHighlighter_AutoRegistersBuiltins(t *testing.T) {
	// 模拟纯 Go 消费者：不手动注册语言/主题，直接按 README 用法调用
	theme.ClearRegistry()
	core.ClearRegistry()
	language.ResetBuiltins()

	hl := NewHighlighter()

	html, err := hl.CodeToHtml("const x = 1;", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, "const") {
		t.Error("expected highlighted output to contain 'const'")
	}
}

func TestNewHighlighter_UnknownThemeFallsBackToDefault(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: "no-such-theme"})
	defaultHl := NewHighlighter()

	code := "const x = 1; /* comment */"
	got, err := hl.CodeToHtml(code, "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want, err := defaultHl.CodeToHtml(code, "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Error("unknown theme should fall back to default theme output")
	}
}

func TestNewHighlighter_InvalidThemeTypeFallsBackToDefault(t *testing.T) {
	hl := NewHighlighter(HighlighterOptions{Theme: 42})
	defaultHl := NewHighlighter()

	code := "const x = 1;"
	got, err := hl.CodeToHtml(code, "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want, err := defaultHl.CodeToHtml(code, "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Error("invalid theme type should fall back to default theme output")
	}
}

func TestCodeToHtml_DeterministicOutput(t *testing.T) {
	hl := NewHighlighter()
	code := "/* block comment */\nconst x = 1; const s = \"hi\";"

	first, err := hl.CodeToHtml(code, "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 2000; i++ {
		out, err := hl.CodeToHtml(code, "javascript")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out != first {
			t.Fatalf("non-deterministic output at iteration %d", i)
		}
	}
}

func TestCodeToHtml_MultibyteNotMojibake(t *testing.T) {
	hl := NewHighlighter()

	html, err := hl.CodeToHtml("let 值 = 1", "javascript")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, "值") {
		t.Errorf("expected HTML to contain complete CJK char, got %q", html)
	}
}

func TestConcurrentHighlighting(t *testing.T) {
	hl := NewHighlighter()

	var wg sync.WaitGroup
	themes := []string{"dark-plus", "github-light", "dracula", "nord"}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				if _, err := hl.CodeToHtml("const x = 1; /* comment */", "javascript"); err != nil {
					t.Errorf("CodeToHtml error: %v", err)
					return
				}
				if _, err := hl.CodeToTokens("let 值 = 1", "javascript"); err != nil {
					t.Errorf("CodeToTokens error: %v", err)
					return
				}
				if err := hl.UpdateTheme(themes[j%len(themes)]); err != nil {
					t.Errorf("UpdateTheme error: %v", err)
					return
				}
			}
		}(i)
	}
	wg.Wait()
}
