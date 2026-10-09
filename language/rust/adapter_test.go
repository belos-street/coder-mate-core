package rust

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
	code := `use std::collections::HashMap;

/// Doc comment
fn greet(name: &str) -> String {
    println!("Hello, {}!", name);
    String::from("ok")
}

fn main() {
    let x: i32 = 42;
    let y = true;
    let z: Option<i32> = None;
    let result = match z {
        Some(val) => val,
        None => 0,
    };
    let config: MyConfig = MyConfig {};
    println!("{}", result + x);
}`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "keyword.declaration.rust")
	checkScopeExists(t, tokens, "keyword.control.rust")
	checkScopeExists(t, tokens, "comment.line.double-slash.rust")
	checkScopeExists(t, tokens, "support.type.builtin.rust")
	checkScopeExists(t, tokens, "entity.name.function.rust")
	checkScopeExists(t, tokens, "constant.language.boolean.rust")
	checkScopeExists(t, tokens, "constant.language.null.rust")
	checkScopeExists(t, tokens, "constant.numeric.rust")
	checkScopeExists(t, tokens, "string.quoted.double.rust")
	checkScopeExists(t, tokens, "operator.rust")
	checkScopeExists(t, tokens, "entity.name.type.rust")
}

func TestParse_BlockComment(t *testing.T) {
	code := `/* block */ fn main() {}`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.block.rust")
}

func TestParse_Attribute(t *testing.T) {
	code := "#[derive(Debug)]\nstruct Foo {}"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.attribute.rust")
}

func TestParse_LanguageVars(t *testing.T) {
	code := "fn foo(self: &Self) { let _ = super::bar(); }"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "variable.language.rust")
}

func TestParse_EmptyString(t *testing.T) {
	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_UnclosedString(t *testing.T) {
	code := `let s = "unclosed`
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "string.quoted.double.rust" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected string tokens")
	}
}

func TestAdapter_ID(t *testing.T) {
	if Language.ID() != "rust" {
		t.Errorf("expected ID 'rust', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	aliases := Language.Aliases()
	found := false
	for _, a := range aliases {
		if a == "rs" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected alias 'rs', got %v", aliases)
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
