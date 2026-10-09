package python

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_MissingKeywords(t *testing.T) {
	t.Parallel()

	// 2.4 回归：这些关键字此前落入 identifier（Python 关键字表缺口）
	for _, kw := range []string{"assert", "del", "global", "nonlocal"} {
		tokens := flatten(Language.Parse(kw + " x"))
		if len(tokens) == 0 || tokens[0].Text != kw {
			t.Errorf("%q: expected first token %q, got %v", kw, kw, tokens)
			continue
		}
		if tokens[0].Scope != ScopeKeywordControl {
			t.Errorf("%q: expected scope %q, got %q", kw, ScopeKeywordControl, tokens[0].Scope)
		}
	}
}

func TestParse_NumberLiterals(t *testing.T) {
	t.Parallel()

	// 2.4 回归：此前数字规则不支持 0x/0o/0b 前缀与下划线分隔
	for _, lit := range []string{"0x1F", "0o17", "0b101", "1_000", "1_000.5"} {
		tokens := flatten(Language.Parse(lit))
		if len(tokens) == 0 || tokens[0].Text != lit {
			t.Errorf("%q: expected single numeric token, got %v", lit, tokens)
			continue
		}
		if tokens[0].Scope != ScopeConstantNumeric {
			t.Errorf("%q: expected scope %q, got %q", lit, ScopeConstantNumeric, tokens[0].Scope)
		}
	}
}

func TestParse_BasicSyntax(t *testing.T) {
	t.Parallel()

	code := `# comment
def greet(name: str) -> str:
    if name:
        return "Hello, " + name
    return "World"

class MyClass:
    def __init__(self, count: int = 42):
        self.count = count

    def run(self):
        result = 0
        for i in range(10):
            result += i
        return result

x = True
y = None
z = 3.14
data = {"key": "value"}
`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "comment.line.number-sign.python")
	checkScopeExists(t, tokens, "keyword.declaration.python")
	checkScopeExists(t, tokens, "keyword.control.python")
	checkScopeExists(t, tokens, "entity.name.function.python")
	checkScopeExists(t, tokens, "entity.name.class.python")
	checkScopeExists(t, tokens, "support.function.builtin.python")
	checkScopeExists(t, tokens, "support.type.annotation.python")
	checkScopeExists(t, tokens, "constant.language.boolean.python")
	checkScopeExists(t, tokens, "constant.language.none.python")
	checkScopeExists(t, tokens, "constant.numeric.python")
	checkScopeExists(t, tokens, "string.quoted.double.python")
	checkScopeExists(t, tokens, "operator.python")
	checkScopeExists(t, tokens, "variable.identifier.python")
}

func TestParse_Decorator(t *testing.T) {
	t.Parallel()

	code := `@app.route('/')
def index():
    return "ok"`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.decorator.python")
}

func TestParse_TripleStrings(t *testing.T) {
	t.Parallel()

	code := `s = """hello
world"""
t = '''single
triple'''`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.quoted.double.triple.python")
	checkScopeExists(t, tokens, "string.quoted.single.triple.python")
}

func TestParse_FString(t *testing.T) {
	t.Parallel()

	code := `name = "world"
msg = f"Hello, {name}!"`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.interpolated.python")
	checkScopeExists(t, tokens, "punctuation.definition.interpolation.begin.python")
	checkScopeExists(t, tokens, "punctuation.definition.interpolation.end.python")
}

func TestParse_EmptyString(t *testing.T) {
	t.Parallel()

	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_UnclosedString(t *testing.T) {
	t.Parallel()

	code := "s = \"unclosed"
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "string.quoted.double.python" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected string tokens")
	}
}

func TestAdapter_ID(t *testing.T) {
	t.Parallel()

	if Language.ID() != "python" {
		t.Errorf("expected ID 'python', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	found := false
	for _, a := range aliases {
		if a == "py" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected alias 'py', got %v", aliases)
	}
}

func flatten(rows core.TokenStream) []core.Token {
	var result []core.Token
	for _, row := range rows {
		result = append(result, row...)
	}
	return result
}

func checkScopeExists(t *testing.T, tokens []core.Token, scope string) {
	t.Helper()
	for _, tok := range tokens {
		if tok.Scope == scope {
			return
		}
	}
	t.Errorf("expected token with scope %q not found", scope)
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

// 精确 token 序列（行为锁定）：任何规则变更都会在此显式失败
func TestParse_TokenSequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		want []seqTok
	}{
		{"representative", `x = 1  # note`, []seqTok{
			{"x", "variable.identifier.python"},
			{" ", "default"},
			{"=", "operator.python"},
			{" ", "default"},
			{"1", "constant.numeric.python"},
			{"  ", "default"},
			{"# note", "comment.line.number-sign.python"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got []seqTok
			for _, row := range Language.Parse(tt.code) {
				for _, tok := range row {
					got = append(got, seqTok{tok.Text, tok.Scope})
				}
			}
			if len(got) != len(tt.want) {
				t.Fatalf("token count = %d, want %d\ngot:  %+v", len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("token %d = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

type seqTok struct {
	Text  string
	Scope string
}
