package javascript

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func requireToken(t *testing.T, tokens core.TokenStream, row, idx int) core.Token {
	t.Helper()
	if row >= len(tokens) {
		t.Fatalf("expected at least %d rows, got %d", row+1, len(tokens))
	}
	if idx >= len(tokens[row]) {
		t.Fatalf("row %d: expected at least %d tokens, got %d", row, idx+1, len(tokens[row]))
	}
	return tokens[row][idx]
}

func findTokenByScope(t *testing.T, tokens core.TokenStream, scope string) core.Token {
	t.Helper()
	for _, row := range tokens {
		for _, tok := range row {
			if tok.Scope == scope {
				return tok
			}
		}
	}
	t.Fatalf("no token found with scope %q", scope)
	return core.Token{}
}

func findTokenByText(t *testing.T, tokens core.TokenStream, text string) core.Token {
	t.Helper()
	for _, row := range tokens {
		for _, tok := range row {
			if tok.Text == text {
				return tok
			}
		}
	}
	t.Fatalf("no token found with text %q", text)
	return core.Token{}
}

func countTokensByScope(tokens core.TokenStream, scope string) int {
	count := 0
	for _, row := range tokens {
		for _, tok := range row {
			if tok.Scope == scope {
				count++
			}
		}
	}
	return count
}

func TestParse_MissingKeywords(t *testing.T) {
	t.Parallel()

	// 2.4 回归：这些关键字此前落入 identifier（JS 关键字表缺口）
	for _, kw := range []string{"typeof", "new", "delete", "void", "yield", "instanceof"} {
		stream := Language.Parse(kw + " x")
		tok := requireToken(t, stream, 0, 0)
		if tok.Text != kw {
			t.Errorf("expected text=%q, got %q", kw, tok.Text)
		}
		if tok.Scope != ScopeKeywordControl {
			t.Errorf("%q: expected scope %q, got %q", kw, ScopeKeywordControl, tok.Scope)
		}
	}

	for _, kw := range []string{"this", "super"} {
		stream := Language.Parse(kw + ".x")
		tok := requireToken(t, stream, 0, 0)
		if tok.Text != kw {
			t.Errorf("expected text=%q, got %q", kw, tok.Text)
		}
		if tok.Scope != ScopeVariableThis {
			t.Errorf("%q: expected scope %q, got %q", kw, ScopeVariableThis, tok.Scope)
		}
	}
}

func TestParse_NumberExponentSign(t *testing.T) {
	t.Parallel()

	// 2.4 回归：此前 1e+10 只匹配出 1e
	stream := Language.Parse("1e+10")
	tok := requireToken(t, stream, 0, 0)
	if tok.Text != "1e+10" {
		t.Errorf("expected full scientific notation, got %q", tok.Text)
	}
	if tok.Scope != ScopeNumeric {
		t.Errorf("expected scope %q, got %q", ScopeNumeric, tok.Scope)
	}
}

func TestParse_ExponentOperator(t *testing.T) {
	t.Parallel()

	// 2.4 回归：此前运算符表无 **，被拆成两个 *
	stream := Language.Parse("2 ** 3")
	tok := requireToken(t, stream, 0, 2)
	if tok.Text != "**" {
		t.Errorf("expected '**' operator, got %q", tok.Text)
	}
}

func TestParse_EmptyInput(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("")
	if len(stream) != 0 {
		t.Fatalf("expected 0 rows for empty input, got %d", len(stream))
	}
}

func TestParse_OnlyNewlines(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("\n\n")
	if len(stream) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(stream))
	}
	if len(stream[0]) < 1 {
		t.Fatal("row 0 should have at least 1 token (newline)")
	}
	if len(stream[1]) < 1 {
		t.Fatal("row 1 should have at least 1 token (newline)")
	}
}

func TestParse_Keywords(t *testing.T) {
	t.Parallel()

	ctrlKeywords := []string{"if", "else", "for", "while", "do", "switch", "case",
		"break", "continue", "return", "throw", "try", "catch", "finally"}
	for _, kw := range ctrlKeywords {
		stream := Language.Parse(kw + " x")
		tok := requireToken(t, stream, 0, 0)
		if tok.Text != kw {
			t.Errorf("expected text=%q, got %q", kw, tok.Text)
		}
		if tok.Scope != ScopeKeywordControl {
			t.Errorf("%q: expected scope=%q, got %q", kw, ScopeKeywordControl, tok.Scope)
		}
	}
}

func TestParse_AsyncAwait(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("async function() {}")
	tok := requireToken(t, stream, 0, 0)
	if tok.Text != "async" {
		t.Errorf("expected 'async', got %q", tok.Text)
	}
	if tok.Scope != ScopeKeywordAsync {
		t.Errorf("expected scope=%q, got %q", ScopeKeywordAsync, tok.Scope)
	}

	stream = Language.Parse("await promise")
	tok = requireToken(t, stream, 0, 0)
	if tok.Text != "await" {
		t.Errorf("expected 'await', got %q", tok.Text)
	}
	if tok.Scope != ScopeKeywordAsync {
		t.Errorf("expected scope=%q, got %q", ScopeKeywordAsync, tok.Scope)
	}
}

func TestParse_ClassKeywords(t *testing.T) {
	t.Parallel()

	classKws := []string{"class", "extends", "static", "constructor"}
	for _, kw := range classKws {
		stream := Language.Parse(kw + " X")
		tok := requireToken(t, stream, 0, 0)
		if tok.Text != kw {
			t.Errorf("%q: expected text=%q, got %q", kw, kw, tok.Text)
		}
		if tok.Scope != ScopeKeywordClass {
			t.Errorf("%q: expected scope=%q, got %q", kw, ScopeKeywordClass, tok.Scope)
		}
	}
}

func TestParse_DeclarationKeywords(t *testing.T) {
	t.Parallel()

	declKws := []string{"function", "var", "let", "const"}
	for _, kw := range declKws {
		stream := Language.Parse(kw + " x")
		tok := requireToken(t, stream, 0, 0)
		if tok.Text != kw {
			t.Errorf("%q: expected text=%q, got %q", kw, kw, tok.Text)
		}
		if tok.Scope != ScopeKeywordDeclaration {
			t.Errorf("%q: expected scope=%q, got %q", kw, ScopeKeywordDeclaration, tok.Scope)
		}
	}
}

func TestParse_BooleanLiterals(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("true")
	tok := requireToken(t, stream, 0, 0)
	if tok.Scope != ScopeBoolean {
		t.Errorf("expected scope=%q, got %q", ScopeBoolean, tok.Scope)
	}

	stream = Language.Parse("false")
	tok = requireToken(t, stream, 0, 0)
	if tok.Scope != ScopeBoolean {
		t.Errorf("expected scope=%q, got %q", ScopeBoolean, tok.Scope)
	}
}

func TestParse_NullUndefined(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("null")
	tok := requireToken(t, stream, 0, 0)
	if tok.Scope != ScopeNull {
		t.Errorf("expected scope=%q, got %q", ScopeNull, tok.Scope)
	}

	stream = Language.Parse("undefined")
	tok = requireToken(t, stream, 0, 0)
	if tok.Scope != ScopeNull {
		t.Errorf("expected scope=%q, got %q", ScopeNull, tok.Scope)
	}
}

func TestParse_LineComment(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("// this is a comment")
	tok := requireToken(t, stream, 0, 0)
	if tok.Text != "// this is a comment" {
		t.Errorf("expected '// this is a comment', got %q", tok.Text)
	}
	if tok.Scope != ScopeCommentLine {
		t.Errorf("expected scope=%q, got %q", ScopeCommentLine, tok.Scope)
	}
}

func TestParse_BlockComment(t *testing.T) {
	t.Parallel()

	code := `/* multi
line */`
	stream := Language.Parse(code)
	if len(stream) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(stream))
	}
	scopeCount := countTokensByScope(stream, ScopeCommentBlock)
	if scopeCount < 2 {
		t.Errorf("expected at least 2 comment block tokens, got %d", scopeCount)
	}
}

func TestParse_GlobalThis(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("globalThis")
	tok := requireToken(t, stream, 0, 0)
	if tok.Scope != ScopeGlobalThis {
		t.Errorf("expected scope=%q, got %q", ScopeGlobalThis, tok.Scope)
	}
}

func TestParse_PromiseAllSettled(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("Promise.allSettled")
	tok := requireToken(t, stream, 0, 0)
	if tok.Scope != ScopeSupportPromise {
		t.Errorf("expected scope=%q, got %q", ScopeSupportPromise, tok.Scope)
	}
}

func TestParse_OptionalChaining(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("?.property")
	tok := requireToken(t, stream, 0, 0)
	if tok.Text != "?." {
		t.Errorf("expected '?.', got %q", tok.Text)
	}
	if tok.Scope != ScopeOptionalChaining {
		t.Errorf("expected scope=%q, got %q", ScopeOptionalChaining, tok.Scope)
	}
}

func TestParse_NullishCoalescing(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("?? value")
	tok := requireToken(t, stream, 0, 0)
	if tok.Text != "??" {
		t.Errorf("expected '??', got %q", tok.Text)
	}
	if tok.Scope != ScopeNullishCoalescing {
		t.Errorf("expected scope=%q, got %q", ScopeNullishCoalescing, tok.Scope)
	}
}

func TestParse_ImportKeyword(t *testing.T) {
	t.Parallel()

	stream := Language.Parse(`import { x } from "./mod"`)
	tok := requireToken(t, stream, 0, 0)
	if tok.Text != "import" {
		t.Errorf("expected 'import', got %q", tok.Text)
	}
	if tok.Scope != ScopeKeywordModule {
		t.Errorf("expected scope=%q, got %q", ScopeKeywordModule, tok.Scope)
	}
}

func TestParse_ExportKeyword(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("export const x = 1")
	tok := requireToken(t, stream, 0, 0)
	if tok.Text != "export" {
		t.Errorf("expected 'export', got %q", tok.Text)
	}
	if tok.Scope != ScopeKeywordModule {
		t.Errorf("expected scope=%q, got %q", ScopeKeywordModule, tok.Scope)
	}
}

func TestParse_ExportStarAs(t *testing.T) {
	t.Parallel()

	stream := Language.Parse(`export * as ns from "./module.js"`)
	tok := requireToken(t, stream, 0, 0)
	if tok.Text != "export * as ns from" {
		t.Errorf("expected 'export * as ns from', got %q", tok.Text)
	}
	if tok.Scope != ScopeKeywordModule {
		t.Errorf("expected scope=%q, got %q", ScopeKeywordModule, tok.Scope)
	}
}

func TestParse_DynamicImport(t *testing.T) {
	t.Parallel()

	stream := Language.Parse(`import("./module.js")`)
	tok := requireToken(t, stream, 0, 0)
	if tok.Text != "import(" {
		t.Errorf("expected 'import(', got %q", tok.Text)
	}
	if tok.Scope != ScopeKeywordImport {
		t.Errorf("expected scope=%q, got %q", ScopeKeywordImport, tok.Scope)
	}
}

func TestParse_ArrowFunction(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("x => x + 1")
	arrowTok := findTokenByText(t, stream, "=>")
	if arrowTok.Scope != ScopeArrowFunction {
		t.Errorf("expected scope=%q, got %q", ScopeArrowFunction, arrowTok.Scope)
	}
}

func TestParse_NumericLiterals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		code string
	}{
		{"42"},
		{"0xFF"},
		{"0b1010"},
		{"0o755"},
		{"3.14"},
		{"1e10"},
	}

	for _, tt := range tests {
		stream := Language.Parse(tt.code)
		tok := requireToken(t, stream, 0, 0)
		if tok.Scope != ScopeNumeric {
			t.Errorf("%q: expected scope=%q, got %q", tt.code, ScopeNumeric, tok.Scope)
		}
	}
}

func TestParse_DoubleQuotedString(t *testing.T) {
	t.Parallel()

	stream := Language.Parse(`"hello world"`)
	// 正则 ^[^\\"]*(?:\\.[^\\"]*)*" 将 "hello world" 整体匹配为一个 token
	// 第一个 token 是 " (PushState)，第二个 token 是 hello world" (PopState)
	if len(stream[0]) < 1 {
		t.Fatal("expected at least 1 token")
	}
	strCount := countTokensByScope(stream, ScopeStringDouble)
	if strCount < 1 {
		t.Errorf("expected at least 1 string.double token, got %d", strCount)
	}
}

func TestParse_SingleQuotedString(t *testing.T) {
	t.Parallel()

	stream := Language.Parse(`'hello world'`)
	strCount := countTokensByScope(stream, ScopeStringSingle)
	if strCount < 1 {
		t.Errorf("expected at least 1 string.single token, got %d", strCount)
	}
}

func TestParse_DoubleQuotedStringWithEscape(t *testing.T) {
	t.Parallel()

	stream := Language.Parse(`"hello \"world\""`)
	// The escaped string rule matches everything including the closing quote
	strCount := countTokensByScope(stream, ScopeStringDouble)
	if strCount < 1 {
		t.Errorf("expected at least 1 string token, got %d", strCount)
	}
}

func TestParse_SingleQuotedStringWithEscape(t *testing.T) {
	t.Parallel()

	stream := Language.Parse(`'hello \'world\''`)
	strCount := countTokensByScope(stream, ScopeStringSingle)
	if strCount < 1 {
		t.Errorf("expected at least 1 string token, got %d", strCount)
	}
}

func TestParse_TemplateString(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("`hello world`")
	templateCount := countTokensByScope(stream, ScopeStringBacktick)
	if templateCount < 2 {
		t.Errorf("expected at least 2 backtick tokens, got %d", templateCount)
	}
}

func TestParse_TemplateInterpolation(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("`hello ${name} world`")
	templateExprCount := countTokensByScope(stream, ScopeTemplateExpression)
	if templateExprCount < 2 {
		t.Errorf("expected at least 2 template expression tokens (${ and }), got %d", templateExprCount)
	}
}

func TestParse_MultiLineCode(t *testing.T) {
	t.Parallel()

	code := "const x = 1;\nlet y = 2;"
	stream := Language.Parse(code)
	if len(stream) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(stream))
	}
	firstTok := requireToken(t, stream, 0, 0)
	if firstTok.Text != "const" {
		t.Errorf("expected first token 'const', got %q", firstTok.Text)
	}
	if firstTok.Line != 1 {
		t.Errorf("expected line=1, got %d", firstTok.Line)
	}

	secondLineFirst := requireToken(t, stream, 1, 0)
	if secondLineFirst.Text != "let" {
		t.Errorf("expected second line first token 'let', got %q", secondLineFirst.Text)
	}
	if secondLineFirst.Line != 2 {
		t.Errorf("expected line=2, got %d", secondLineFirst.Line)
	}
}

func TestParse_CrossLineCommentLineTracking(t *testing.T) {
	t.Parallel()

	code := "/*a\nb*/\nconst x = 1;"
	stream := Language.Parse(code)

	if len(stream) < 3 {
		t.Fatalf("expected at least 3 rows, got %d", len(stream))
	}

	thirdLineFirst := requireToken(t, stream, 2, 0)
	if thirdLineFirst.Text != "const" {
		t.Errorf("expected 'const', got %q", thirdLineFirst.Text)
	}
	if thirdLineFirst.Line != 3 {
		t.Errorf("expected line=3, got %d", thirdLineFirst.Line)
	}
	if thirdLineFirst.Col[0] != 1 {
		t.Errorf("expected col start=1, got %d", thirdLineFirst.Col[0])
	}
}

func TestParse_NewlineToken(t *testing.T) {
	t.Parallel()

	code := "x\ny"
	stream := Language.Parse(code)

	if len(stream) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(stream))
	}

	row0 := stream[0]
	foundNewline := false
	for _, tok := range row0 {
		if tok.Text == "\n" {
			foundNewline = true
			if tok.Line != 1 {
				t.Errorf("newline token line should be 1, got %d", tok.Line)
			}
		}
	}
	if !foundNewline {
		t.Error("expected newline token in first row")
	}

	row1 := stream[1]
	if len(row1) != 1 || row1[0].Text != "y" {
		t.Errorf("expected row 1 = ['y'], got %v", row1)
	}
	_ = requireToken(t, stream, 1, 0)
}

func TestParse_ComplexCode(t *testing.T) {
	t.Parallel()

	code := `const greeting = \x60Hello, ${name}!\x60;
// A comment
class Person {
  constructor(name) {
    this.name = name;
  }
  greet() {
    return greeting.replace('{name}', this.name);
  }
}
`
	stream := Language.Parse(code)
	if len(stream) < 5 {
		t.Fatalf("expected at least 5 rows, got %d", len(stream))
	}

	constToken := findTokenByText(t, stream, "const")
	if constToken.Scope != ScopeKeywordDeclaration {
		t.Errorf("expected 'const' to be declaration keyword, got %q", constToken.Scope)
	}

	commentToken := findTokenByText(t, stream, "// A comment")
	if commentToken.Scope != ScopeCommentLine {
		t.Errorf("expected comment scope, got %q", commentToken.Scope)
	}

	classToken := findTokenByText(t, stream, "class")
	if classToken.Scope != ScopeKeywordClass {
		t.Errorf("expected 'class' to be class keyword, got %q", classToken.Scope)
	}

	constructorToken := findTokenByText(t, stream, "constructor")
	if constructorToken.Scope != ScopeKeywordClass {
		t.Errorf("expected 'constructor' to be class keyword, got %q", constructorToken.Scope)
	}
}

func TestAdapter_ID(t *testing.T) {
	t.Parallel()

	if Language.ID() != "javascript" {
		t.Errorf("expected ID='javascript', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	if len(aliases) != 1 || aliases[0] != "js" {
		t.Errorf("expected aliases=['js'], got %v", aliases)
	}
}

func TestParse_Operators(t *testing.T) {
	t.Parallel()

	ops := []string{"===", "!==", "==", "!=", ">=", "<=", "&&", "||",
		"+", "-", "*", "/", "%", "=", ">", "<", "!", ".",
		",", ":", ";", "(", ")", "{", "}", "[", "]"}

	for _, op := range ops {
		stream := Language.Parse(op + " x")
		tok := requireToken(t, stream, 0, 0)
		if tok.Text != op {
			t.Errorf("%q: expected text=%q, got %q", op, op, tok.Text)
		}
		if tok.Scope != ScopeOperator {
			t.Errorf("%q: expected scope=%q, got %q", op, ScopeOperator, tok.Scope)
		}
	}
}

// 多字节字符不得被字节级切碎（分词无损性）
func TestParse_MultibyteLossless(t *testing.T) {
	t.Parallel()

	stream := Language.Parse("值")

	var b strings.Builder
	for _, row := range stream {
		for _, tok := range row {
			b.WriteString(tok.Text)
		}
	}
	if got := b.String(); got != "值" {
		t.Errorf("expected lossless tokenization of multibyte char, got %q", got)
	}
}
