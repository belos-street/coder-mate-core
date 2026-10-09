package rust

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	code := `/* block */ fn main() {}`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.block.rust")
}

func TestParse_Attribute(t *testing.T) {
	t.Parallel()

	code := "#[derive(Debug)]\nstruct Foo {}"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.attribute.rust")
}

func TestParse_LanguageVars(t *testing.T) {
	t.Parallel()

	code := "fn foo(self: &Self) { let _ = super::bar(); }"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "variable.language.rust")
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
	t.Parallel()

	if Language.ID() != "rust" {
		t.Errorf("expected ID 'rust', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

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
		{"representative", `fn main() { let x = 1; } // c`, []seqTok{
			{"fn", "keyword.declaration.rust"},
			{" ", "default"},
			{"main(", "entity.name.function.rust"},
			{")", "operator.rust"},
			{" ", "default"},
			{"{", "operator.rust"},
			{" ", "default"},
			{"let", "keyword.declaration.rust"},
			{" ", "default"},
			{"x", "default"},
			{" ", "default"},
			{"=", "operator.rust"},
			{" ", "default"},
			{"1", "constant.numeric.rust"},
			{";", "operator.rust"},
			{" ", "default"},
			{"}", "operator.rust"},
			{" ", "default"},
			{"// c", "comment.line.double-slash.rust"},
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
