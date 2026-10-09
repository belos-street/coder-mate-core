package typescript

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
	t.Parallel()

	code := `interface User {
  id: number;
  name: string;
}

type Status = "active" | "inactive";

enum Color { Red, Green, Blue }

function greet(name: string): void {
  if (name) {
    console.log("Hello, " + name);
  }
}

// a comment
const count: number = 42;
let ok: boolean = true;
let data: unknown = null;
`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "entity.name.type.typescript")
	checkScopeExists(t, tokens, "keyword.declaration.type.typescript")
	checkScopeExists(t, tokens, "support.type.builtin.typescript")
	checkScopeExists(t, tokens, "keyword.control.js")
	checkScopeExists(t, tokens, "keyword.declaration.js")
	checkScopeExists(t, tokens, "constant.language.boolean.js")
	checkScopeExists(t, tokens, "constant.language.null.js")
	checkScopeExists(t, tokens, "constant.numeric.js")
	checkScopeExists(t, tokens, "string.quoted.double.js")
	checkScopeExists(t, tokens, "comment.line.double-slash.js")
	checkScopeExists(t, tokens, "operator.js")
	checkScopeExists(t, tokens, "variable.identifier.js")
}

func TestParse_ClassDeclaration(t *testing.T) {
	t.Parallel()

	code := `export class MyClass<T> extends BaseClass implements IInterface {
  public name: string;
  private age: number;
  constructor(name: string) {
    super(name);
  }
}`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "entity.name.type.typescript")
	checkScopeExists(t, tokens, "keyword.modifier.access.typescript")
}

func TestParse_TypeAssertion(t *testing.T) {
	t.Parallel()

	code := `let x = value as string;
let y = <number>value;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "keyword.operator.assertion.typescript")
}

func TestParse_EmptyString(t *testing.T) {
	t.Parallel()

	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_LineComment(t *testing.T) {
	t.Parallel()

	code := "// comment\nlet x = 1;"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.line.double-slash.js")
}

func TestAdapter_ID(t *testing.T) {
	t.Parallel()

	if Language.ID() != "typescript" {
		t.Errorf("expected ID 'typescript', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	found := false
	for _, a := range aliases {
		if a == "ts" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected alias 'ts', got %v", aliases)
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
		{"representative", `const x: number = 1; // c`, []seqTok{
			{"const", "keyword.declaration.js"},
			{" ", "default"},
			{"x", "variable.identifier.js"},
			{":", "operator.js"},
			{" ", "default"},
			{"number", "support.type.builtin.typescript"},
			{" ", "default"},
			{"=", "operator.js"},
			{" ", "default"},
			{"1", "constant.numeric.js"},
			{";", "operator.js"},
			{" ", "default"},
			{"// c", "comment.line.double-slash.js"},
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
