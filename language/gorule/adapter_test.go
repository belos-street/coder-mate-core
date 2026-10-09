package gorule

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
	t.Parallel()

	code := `package main

import "fmt"

// comment
const answer = 42

func greet(name string) string {
	return "Hello, " + name
}

func main() {
	var x int = 10
	y := 20
	ok := true
	if ok == false && x > 0 && y < 100 {
		fmt.Println(greet("World"))
	}
	ptr := &x
	nilCheck(ptr)
}

var data map[string]interface{}

func nilCheck(p *int) bool {
	return p != nil
}
`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "keyword.declaration.go")
	checkScopeExists(t, tokens, "keyword.control.go")
	checkScopeExists(t, tokens, "support.type.builtin.go")
	checkScopeExists(t, tokens, "entity.name.function.go")
	checkScopeExists(t, tokens, "constant.language.boolean.go")
	checkScopeExists(t, tokens, "constant.language.null.go")
	checkScopeExists(t, tokens, "constant.numeric.go")
	checkScopeExists(t, tokens, "string.quoted.double.go")
	checkScopeExists(t, tokens, "comment.line.double-slash.go")
	checkScopeExists(t, tokens, "operator.go")
}

func TestParse_BlockComment(t *testing.T) {
	t.Parallel()

	code := `/* block */ var x = 1`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.block.go")
}

func TestParse_RawString(t *testing.T) {
	t.Parallel()

	code := "var s = `raw\nstring`"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.quoted.raw.go")
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

	code := `s := "unclosed`
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "string.quoted.double.go" {
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

	if Language.ID() != "go" {
		t.Errorf("expected ID 'go', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	found := false
	for _, a := range aliases {
		if a == "golang" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected alias 'golang', got %v", aliases)
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
		{"representative", `func F(a int) string { return "s" }`, []seqTok{
			{"func", "keyword.declaration.go"},
			{" ", "default"},
			{"F(", "entity.name.function.go"},
			{"a", "default"},
			{" ", "default"},
			{"int", "support.type.builtin.go"},
			{")", "operator.go"},
			{" ", "default"},
			{"string", "support.type.builtin.go"},
			{" ", "default"},
			{"{", "operator.go"},
			{" ", "default"},
			{"return", "keyword.control.go"},
			{" ", "default"},
			{"\"", "string.quoted.double.go"},
			{"s", "string.quoted.double.go"},
			{"\"", "string.quoted.double.go"},
			{" ", "default"},
			{"}", "operator.go"},
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
