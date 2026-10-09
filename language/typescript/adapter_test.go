package typescript

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
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
	code := `let x = value as string;
let y = <number>value;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "keyword.operator.assertion.typescript")
}

func TestParse_EmptyString(t *testing.T) {
	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_LineComment(t *testing.T) {
	code := "// comment\nlet x = 1;"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.line.double-slash.js")
}

func TestAdapter_ID(t *testing.T) {
	if Language.ID() != "typescript" {
		t.Errorf("expected ID 'typescript', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
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
