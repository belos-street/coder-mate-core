package python

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
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
	code := `@app.route('/')
def index():
    return "ok"`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.decorator.python")
}

func TestParse_TripleStrings(t *testing.T) {
	code := `s = """hello
world"""
t = '''single
triple'''`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.quoted.double.triple.python")
	checkScopeExists(t, tokens, "string.quoted.single.triple.python")
}

func TestParse_FString(t *testing.T) {
	code := `name = "world"
msg = f"Hello, {name}!"`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.interpolated.python")
	checkScopeExists(t, tokens, "punctuation.definition.interpolation.begin.python")
	checkScopeExists(t, tokens, "punctuation.definition.interpolation.end.python")
}

func TestParse_EmptyString(t *testing.T) {
	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_UnclosedString(t *testing.T) {
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
	if Language.ID() != "python" {
		t.Errorf("expected ID 'python', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
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
