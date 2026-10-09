package core

import (
	"regexp"
	"strings"
	"testing"
)

func TestNewParserContext(t *testing.T) {
	ctx := NewParserContext("global")

	if len(ctx.StateStack) != 1 {
		t.Fatalf("expected StateStack len 1, got %d", len(ctx.StateStack))
	}
	if ctx.StateStack[0] != "global" {
		t.Errorf("expected StateStack[0]='global', got %q", ctx.StateStack[0])
	}
	if ctx.Line != 1 {
		t.Errorf("expected Line=1, got %d", ctx.Line)
	}
	if ctx.Col != 1 {
		t.Errorf("expected Col=1, got %d", ctx.Col)
	}
}

func TestPushState_Immutability(t *testing.T) {
	ctx := NewParserContext("global")
	newCtx := PushState(ctx, "string")

	if len(ctx.StateStack) != 1 {
		t.Errorf("original should have 1 element, got %d", len(ctx.StateStack))
	}
	if len(newCtx.StateStack) != 2 {
		t.Errorf("new should have 2 elements, got %d", len(newCtx.StateStack))
	}
	if newCtx.StateStack[1] != "string" {
		t.Errorf("expected StateStack[1]='string', got %q", newCtx.StateStack[1])
	}

	ctx.StateStack[0] = "modified"
	if newCtx.StateStack[0] == "modified" {
		t.Error("newCtx.StateStack should be independent copy")
	}
}

func TestPopState(t *testing.T) {
	ctx := NewParserContext("global")
	ctx = PushState(ctx, "string")
	newCtx, err := PopState(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ctx.StateStack) != 2 {
		t.Errorf("original should still have 2 elements, got %d", len(ctx.StateStack))
	}
	if len(newCtx.StateStack) != 1 {
		t.Errorf("new should have 1 element, got %d", len(newCtx.StateStack))
	}
	if newCtx.StateStack[0] != "global" {
		t.Errorf("expected StateStack[0]='global', got %q", newCtx.StateStack[0])
	}
}

func TestPopState_SingleElementError(t *testing.T) {
	ctx := NewParserContext("global")
	_, err := PopState(ctx)
	if err == nil {
		t.Fatal("expected error when popping from single-element stack")
	}
}

func TestCurrentState(t *testing.T) {
	ctx := NewParserContext("global")
	if CurrentState(ctx) != "global" {
		t.Errorf("expected 'global', got %q", CurrentState(ctx))
	}

	ctx = PushState(ctx, "inner")
	if CurrentState(ctx) != "inner" {
		t.Errorf("expected 'inner', got %q", CurrentState(ctx))
	}
}

func TestSplitTokenByLineBreak_NoBreak(t *testing.T) {
	tokens := SplitTokenByLineBreak("hello", "scope.x", "fallback", 1, 1)
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].Text != "hello" {
		t.Errorf("expected Text='hello', got %q", tokens[0].Text)
	}
	if tokens[0].Scope != "scope.x" {
		t.Errorf("expected Scope='scope.x', got %q", tokens[0].Scope)
	}
	if tokens[0].Line != 1 {
		t.Errorf("expected Line=1, got %d", tokens[0].Line)
	}
	expectedCol := [2]int{1, 5}
	if tokens[0].Col != expectedCol {
		t.Errorf("expected Col=%v, got %v", expectedCol, tokens[0].Col)
	}
}

func TestSplitTokenByLineBreak_LF(t *testing.T) {
	tokens := SplitTokenByLineBreak("hello\nworld", "scope.x", "fallback", 1, 1)
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}
	// token 0: "hello"
	if tokens[0].Text != "hello" || tokens[0].Line != 1 {
		t.Errorf("token[0]: expected 'hello' line 1, got %q line %d", tokens[0].Text, tokens[0].Line)
	}
	// token 1: "\n" with fallback scope
	if tokens[1].Text != "\n" || tokens[1].Scope != "fallback" {
		t.Errorf("token[1]: expected '\\n' scope 'fallback', got %q scope %q", tokens[1].Text, tokens[1].Scope)
	}
	// token 2: "world"
	if tokens[2].Text != "world" || tokens[2].Line != 2 {
		t.Errorf("token[2]: expected 'world' line 2, got %q line %d", tokens[2].Text, tokens[2].Line)
	}
}

func TestSplitTokenByLineBreak_CRLF(t *testing.T) {
	tokens := SplitTokenByLineBreak("hello\r\nworld", "scope.x", "fallback", 1, 1)
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}
	if tokens[1].Text != "\r\n" {
		t.Errorf("expected '\\r\\n', got %q", tokens[1].Text)
	}
}

func TestSplitTokenByLineBreak_MultipleLines(t *testing.T) {
	tokens := SplitTokenByLineBreak("a\nb\nc", "scope.x", "fallback", 1, 1)
	if len(tokens) != 5 {
		t.Fatalf("expected 5 tokens, got %d", len(tokens))
	}
	if tokens[0].Line != 1 || tokens[2].Line != 2 || tokens[4].Line != 3 {
		t.Errorf("line numbers: got %d, %d, %d", tokens[0].Line, tokens[2].Line, tokens[4].Line)
	}
}

func TestSplitTokenByLineBreak_EmptyString(t *testing.T) {
	tokens := SplitTokenByLineBreak("", "scope.x", "fallback", 1, 1)
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}
}

func TestSplitTokenByLineBreak_ConsecutiveNewlines(t *testing.T) {
	tokens := SplitTokenByLineBreak("a\n\nb", "scope.x", "fallback", 1, 1)
	if len(tokens) != 4 {
		t.Fatalf("expected 4 tokens, got %d", len(tokens))
	}
	if tokens[0].Text != "a" || tokens[0].Line != 1 {
		t.Errorf("token[0]: expected 'a' line 1, got %q line %d", tokens[0].Text, tokens[0].Line)
	}
	if tokens[1].Text != "\n" {
		t.Errorf("token[1]: expected '\\n', got %q", tokens[1].Text)
	}
	// Empty gap between \n and \n produces no token
	if tokens[2].Text != "\n" {
		t.Errorf("token[2]: expected '\\n', got %q", tokens[2].Text)
	}
	if tokens[3].Text != "b" || tokens[3].Line != 3 {
		t.Errorf("token[3]: expected 'b' line 3, got %q line %d", tokens[3].Text, tokens[3].Line)
	}
}

func TestSplitTokenByLineBreak_StartsWithNewline(t *testing.T) {
	tokens := SplitTokenByLineBreak("\nhello", "scope.x", "fallback", 1, 1)
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	// Empty before first \n produces no token
	if tokens[0].Text != "\n" {
		t.Errorf("token[0]: expected '\\n', got %q", tokens[0].Text)
	}
	if tokens[1].Text != "hello" || tokens[1].Line != 2 {
		t.Errorf("token[1]: expected 'hello' line 2, got %q line %d", tokens[1].Text, tokens[1].Line)
	}
}

func TestMatchToken_MatchesFirstRule(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules: map[string][]GrammarRule{
			"global": {
				{Regex: regexp.MustCompile(`^\/\/.*`), Scope: "comment.line.js"},
				{Regex: regexp.MustCompile(`^\w+`), Scope: "identifier.js"},
			},
		},
	}
	ctx := NewParserContext("global")
	token, newCtx := MatchToken("// comment", ctx, spec)
	if token == nil {
		t.Fatal("expected token, got nil")
	}
	if token.Text != "// comment" {
		t.Errorf("expected Text='// comment', got %q", token.Text)
	}
	if token.Scope != "comment.line.js" {
		t.Errorf("expected Scope='comment.line.js', got %q", token.Scope)
	}
	if len(newCtx.StateStack) != 1 {
		t.Errorf("expected StateStack len 1, got %d", len(newCtx.StateStack))
	}
}

func TestMatchToken_FallbackWhenNoMatch(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules: map[string][]GrammarRule{
			"global": {
				{Regex: regexp.MustCompile(`^\/\/.*`), Scope: "comment.line.js"},
			},
		},
	}
	ctx := NewParserContext("global")
	token, _ := MatchToken("x", ctx, spec)
	if token == nil {
		t.Fatal("expected fallback token, got nil")
	}
	if token.Text != "x" {
		t.Errorf("expected Text='x', got %q", token.Text)
	}
	if token.Scope != "default" {
		t.Errorf("expected Scope='default', got %q", token.Scope)
	}
}

func TestMatchToken_PushState(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules: map[string][]GrammarRule{
			"global": {
				{Regex: regexp.MustCompile(`^"`), Scope: "string.start.js", PushState: "string-double"},
			},
		},
	}
	ctx := NewParserContext("global")
	_, newCtx := MatchToken(`"`, ctx, spec)
	if CurrentState(newCtx) != "string-double" {
		t.Errorf("expected state 'string-double', got %q", CurrentState(newCtx))
	}
}

func TestMatchToken_PopState(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules: map[string][]GrammarRule{
			"global": {
				{Regex: regexp.MustCompile(`^"`), Scope: "string.start.js", PushState: "string-double"},
			},
			"string-double": {
				{Regex: regexp.MustCompile(`^"`), Scope: "string.end.js", PopState: true},
			},
		},
	}
	ctx := NewParserContext("global")
	ctx = PushState(ctx, "string-double")
	_, newCtx := MatchToken(`"`, ctx, spec)
	if CurrentState(newCtx) != "global" {
		t.Errorf("expected state 'global', got %q", CurrentState(newCtx))
	}
}

func TestMatchToken_EmptyCode(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules:         map[string][]GrammarRule{},
	}
	ctx := NewParserContext("global")
	token, newCtx := MatchToken("", ctx, spec)
	if token != nil {
		t.Error("expected nil token for empty code")
	}
	if newCtx != ctx {
		t.Error("expected same context returned for empty code")
	}
}

func TestParse_EmptyInput(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules:         map[string][]GrammarRule{},
	}
	stream := Parse("", spec)
	if len(stream) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(stream))
	}
}

func TestParse_SingleLine(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules: map[string][]GrammarRule{
			"global": {
				{Regex: regexp.MustCompile(`^\/\/.*`), Scope: "comment.line.js"},
				{Regex: regexp.MustCompile(`^\w+`), Scope: "identifier.js"},
				{Regex: regexp.MustCompile(`^=`), Scope: "operator.js"},
				{Regex: regexp.MustCompile(`^\d+`), Scope: "number.js"},
				{Regex: regexp.MustCompile(`^\s+`), Scope: "whitespace"},
			},
		},
	}
	stream := Parse("x = 42", spec)
	if len(stream) != 1 {
		t.Fatalf("expected 1 row, got %d", len(stream))
	}
	if len(stream[0]) < 3 {
		t.Fatalf("expected at least 3 tokens, got %d", len(stream[0]))
	}
	if stream[0][0].Text != "x" || stream[0][0].Scope != "identifier.js" {
		t.Errorf("token[0]: expected 'x' scope 'identifier.js', got %q scope %q",
			stream[0][0].Text, stream[0][0].Scope)
	}
}

func TestParse_MultiLine(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules: map[string][]GrammarRule{
			"global": {
				{Regex: regexp.MustCompile(`^\w+`), Scope: "identifier.js"},
				{Regex: regexp.MustCompile(`^\s+`), Scope: "whitespace"},
			},
		},
	}
	stream := Parse("foo\nbar", spec)
	if len(stream) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(stream))
	}
}

func TestParse_StringWithStateChange(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules: map[string][]GrammarRule{
			"global": {
				{Regex: regexp.MustCompile(`^"`), Scope: "string.start.js", PushState: "string-double"},
				{Regex: regexp.MustCompile(`^\w+`), Scope: "identifier.js"},
			},
			"string-double": {
				{Regex: regexp.MustCompile(`^[^"]+`), Scope: "string.content.js"},
				{Regex: regexp.MustCompile(`^"`), Scope: "string.end.js", PopState: true},
			},
		},
	}
	stream := Parse(`"hello"`, spec)
	if len(stream) != 1 {
		t.Fatalf("expected 1 row, got %d", len(stream))
	}
	if len(stream[0]) != 3 {
		t.Fatalf("expected 3 tokens (quote, content, quote), got %d", len(stream[0]))
	}
	if stream[0][0].Scope != "string.start.js" {
		t.Errorf("token[0] scope: expected 'string.start.js', got %q", stream[0][0].Scope)
	}
	if stream[0][1].Scope != "string.content.js" {
		t.Errorf("token[1] scope: expected 'string.content.js', got %q", stream[0][1].Scope)
	}
	if stream[0][2].Scope != "string.end.js" {
		t.Errorf("token[2] scope: expected 'string.end.js', got %q", stream[0][2].Scope)
	}
}

func TestParse_CrossLineToken(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules: map[string][]GrammarRule{
			"global": {
				{Regex: regexp.MustCompile(`^"`), Scope: "string.start.js", PushState: "string-double"},
			},
			"string-double": {
				{Regex: regexp.MustCompile(`^[^"]+`), Scope: "string.content.js"},
				{Regex: regexp.MustCompile(`^"`), Scope: "string.end.js", PopState: true},
			},
		},
	}
	stream := Parse("\"hello\nworld\"", spec)
	if len(stream) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(stream))
	}
	// Row 0: " hello
	if len(stream[0]) < 2 {
		t.Fatalf("expected at least 2 tokens in row 0, got %d", len(stream[0]))
	}
	if stream[0][0].Scope != "string.start.js" {
		t.Errorf("row0 token[0]: expected 'string.start.js', got %q", stream[0][0].Scope)
	}
	if stream[0][1].Scope != "string.content.js" {
		t.Errorf("row0 token[1]: expected 'string.content.js', got %q", stream[0][1].Scope)
	}
	// Row 1: world"
	if len(stream[1]) < 2 {
		t.Fatalf("expected at least 2 tokens in row 1, got %d", len(stream[1]))
	}
	if stream[1][0].Scope != "string.content.js" {
		t.Errorf("row1 token[0]: expected 'string.content.js', got %q", stream[1][0].Scope)
	}
	if stream[1][1].Scope != "string.end.js" {
		t.Errorf("row1 token[1]: expected 'string.end.js', got %q", stream[1][1].Scope)
	}
}

func TestParse_ConsecutiveNewlines(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "fallback",
		Rules: map[string][]GrammarRule{
			"global": {
				{Regex: regexp.MustCompile(`^\w+`), Scope: "ident"},
			},
		},
	}
	stream := Parse("a\n\nb", spec)
	if len(stream) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(stream))
	}
	// Row 1: a + \n
	if len(stream[0]) != 2 || stream[0][0].Text != "a" {
		t.Errorf("row 0: expected ['a', '\\n'], got %v", stream[0])
	}
	// Row 2: \n (empty line)
	if len(stream[1]) != 1 || stream[1][0].Scope != "fallback" {
		t.Errorf("row 1: expected [\\n fallback], got %v", stream[1])
	}
	// Row 3: b
	if len(stream[2]) != 1 || stream[2][0].Text != "b" {
		t.Errorf("row 2: expected ['b'], got %v", stream[2])
	}
}

func TestMatchToken_UnknownStateFallbackToGlobal(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules: map[string][]GrammarRule{
			"global": {
				{Regex: regexp.MustCompile(`^\w+`), Scope: "ident"},
			},
		},
	}
	ctx := NewParserContext("global")
	ctx = PushState(ctx, "nonexistent-state")
	token, newCtx := MatchToken("hello", ctx, spec)
	if token == nil {
		t.Fatal("expected token via fallback, got nil")
	}
	if token.Scope != "default" {
		t.Errorf("expected fallback scope 'default' for unknown state, got %q", token.Scope)
	}
	// Should still be in the unknown state
	if CurrentState(newCtx) != "nonexistent-state" {
		t.Errorf("expected state to remain 'nonexistent-state', got %q", CurrentState(newCtx))
	}
}

type mockAdapter struct {
	id      string
	aliases []string
}

func (m mockAdapter) ID() string                    { return m.id }
func (m mockAdapter) Aliases() []string             { return m.aliases }
func (m mockAdapter) Parse(code string) TokenStream { return nil }

func TestRegisterLanguage(t *testing.T) {
	ClearRegistry()

	lang := mockAdapter{id: "javascript", aliases: []string{"js"}}
	err := RegisterLanguage(lang)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, ok := GetLanguage("javascript")
	if !ok {
		t.Fatal("expected to find 'javascript'")
	}
	if found.ID() != "javascript" {
		t.Errorf("expected ID 'javascript', got %q", found.ID())
	}
}

func TestRegisterLanguage_Alias(t *testing.T) {
	ClearRegistry()

	lang := mockAdapter{id: "javascript", aliases: []string{"js"}}
	RegisterLanguage(lang)

	found, ok := GetLanguage("js")
	if !ok {
		t.Fatal("expected to find 'js' alias")
	}
	if found.ID() != "javascript" {
		t.Errorf("expected ID 'javascript' via alias, got %q", found.ID())
	}
}

func TestRegisterLanguage_Duplicate(t *testing.T) {
	ClearRegistry()

	RegisterLanguage(mockAdapter{id: "python", aliases: []string{"py"}})
	err := RegisterLanguage(mockAdapter{id: "python", aliases: []string{}})
	if err == nil {
		t.Fatal("expected error for duplicate register")
	}
}

func TestGetLanguage_CaseInsensitive(t *testing.T) {
	ClearRegistry()

	RegisterLanguage(mockAdapter{id: "javascript", aliases: []string{"js"}})
	_, ok := GetLanguage("JAVASCRIPT")
	if !ok {
		t.Fatal("expected case-insensitive lookup to succeed")
	}
	_, ok = GetLanguage("  javascript  ")
	if !ok {
		t.Fatal("expected trimmed lookup to succeed")
	}
}

func TestListLanguages(t *testing.T) {
	ClearRegistry()

	RegisterLanguage(mockAdapter{id: "javascript", aliases: []string{"js"}})
	RegisterLanguage(mockAdapter{id: "python", aliases: []string{"py"}})

	list := ListLanguages()
	if len(list) != 2 {
		t.Fatalf("expected 2 languages, got %d", len(list))
	}
}

func TestListLanguages_Deduplication(t *testing.T) {
	ClearRegistry()

	lang := mockAdapter{id: "go", aliases: []string{"golang"}}
	RegisterLanguage(lang)

	list := ListLanguages()
	if len(list) != 1 {
		t.Fatalf("expected 1 language (deduplicated), got %d", len(list))
	}
}

func TestTokenize_Success(t *testing.T) {
	ClearRegistry()

	// 注册一个返回真实 token 流的 adapter，验证 Tokenize 的分发行为
	RegisterLanguage(funcAdapter{fn: func(code string) TokenStream {
		return TokenStream{{{Text: code, Scope: "test.scope", Line: 1, Col: [2]int{1, len(code)}}}}
	}})

	rows, err := Tokenize("hello", "func")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 1 || len(rows[0]) != 1 {
		t.Fatalf("expected 1 row with 1 token, got %v", rows)
	}
	if rows[0][0].Text != "hello" || rows[0][0].Scope != "test.scope" {
		t.Errorf("unexpected token: %+v", rows[0][0])
	}
}

func TestEnsureValidID_PanicsOnEmpty(t *testing.T) {
	// panic 路径：空语言 ID 在注册表入口 fail-fast（阶段 3 补齐的覆盖）
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty language id")
		}
	}()
	GetLanguage("   ")
}

func TestRegisterLanguage_PanicsOnEmptyID(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty adapter id")
		}
	}()
	RegisterLanguage(mockAdapter{id: "", aliases: nil})
}

// funcAdapter 按注入的函数返回解析结果
type funcAdapter struct{ fn func(code string) TokenStream }

func (f funcAdapter) ID() string                    { return "func" }
func (f funcAdapter) Aliases() []string             { return nil }
func (f funcAdapter) Parse(code string) TokenStream { return f.fn(code) }

func TestTokenize_UnknownLanguage(t *testing.T) {
	ClearRegistry()

	_, err := Tokenize("hello", "unknown")
	if err == nil {
		t.Fatal("expected error for unknown language")
	}
}

func TestTokenize_EmptyLangID(t *testing.T) {
	ClearRegistry()

	_, err := Tokenize("hello", "")
	if err == nil {
		t.Fatal("expected error for empty language ID")
	}
}

func TestMatchToken_FallbackMultibyte(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules:         map[string][]GrammarRule{},
	}
	ctx := NewParserContext("global")

	// 中文 = 3 字节
	token, newCtx := MatchToken("值", ctx, spec)
	if token == nil {
		t.Fatal("expected fallback token, got nil")
	}
	if token.Text != "值" {
		t.Errorf("expected complete multibyte char, got %q", token.Text)
	}
	if newCtx.Col != ctx.Col+3 { // Col 按字节推进
		t.Errorf("expected Col to advance by byte width (3), got %d", newCtx.Col)
	}

	// emoji = 4 字节
	token4, newCtx4 := MatchToken("🎨", ctx, spec)
	if token4 == nil {
		t.Fatal("expected fallback token, got nil")
	}
	if token4.Text != "🎨" {
		t.Errorf("expected complete emoji, got %q", token4.Text)
	}
	if newCtx4.Col != ctx.Col+4 {
		t.Errorf("expected Col to advance by byte width (4), got %d", newCtx4.Col)
	}
}

func TestMatchToken_FallbackInvalidUTF8(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules:         map[string][]GrammarRule{},
	}
	ctx := NewParserContext("global")

	token, newCtx := MatchToken("\xff", ctx, spec)
	if token == nil {
		t.Fatal("expected fallback token, got nil")
	}
	// 非法字节按原字节保留（不替换成 U+FFFD），分词完全无损
	if token.Text != "\xff" {
		t.Errorf("expected raw byte preserved, got %q", token.Text)
	}
	if newCtx.Col != ctx.Col+1 {
		t.Errorf("expected Col to advance by 1, got %d", newCtx.Col)
	}
}

func TestParse_ZeroWidthRuleGuard(t *testing.T) {
	spec := &TokenizerSpec{
		InitialState:  "global",
		FallbackScope: "default",
		Rules: map[string][]GrammarRule{
			"global": {
				// a* 可匹配空串：无零长度消费防护时 Parse 会死循环
				{Regex: regexp.MustCompile(`a*`), Scope: "zero"},
			},
		},
	}

	joinTokens := func(stream TokenStream) string {
		var b strings.Builder
		for _, row := range stream {
			for _, tok := range row {
				b.WriteString(tok.Text)
			}
		}
		return b.String()
	}

	// 输入不命中规则：零宽匹配触发防护，按 rune 兜底消费
	if got := joinTokens(Parse("bbb", spec)); got != "bbb" {
		t.Errorf("expected lossless consumption, got %q", got)
	}

	// 输入部分命中："a" 被规则正常消费（非空匹配），"b" 走兜底
	if got := joinTokens(Parse("ab", spec)); got != "ab" {
		t.Errorf("expected lossless consumption, got %q", got)
	}
}
