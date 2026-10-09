package c

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
	t.Parallel()

	code := `#include <stdio.h>

// comment
#define MAX 3

typedef struct User {
  int id;
  const char *name;
  bool active;
  uint32_t total;
} User;

int sum(int a, int b) {
  return a + b;
}

int main(void) {
  User u = {.id = 1, .name = "coder", .active = true};
  char flag = 'Y';
  User *ptr = NULL;
  return 0;
}`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "keyword.control.directive.c")
	checkScopeExists(t, tokens, "string.quoted.angle.c")
	checkScopeExists(t, tokens, "comment.line.double-slash.c")
	checkScopeExists(t, tokens, "keyword.declaration.c")
	checkScopeExists(t, tokens, "keyword.modifier.c")
	checkScopeExists(t, tokens, "keyword.control.c")
	checkScopeExists(t, tokens, "support.type.builtin.c")
	checkScopeExists(t, tokens, "entity.name.function.c")
	checkScopeExists(t, tokens, "entity.name.macro.c")
	checkScopeExists(t, tokens, "constant.language.boolean.c")
	checkScopeExists(t, tokens, "constant.language.null.c")
	checkScopeExists(t, tokens, "constant.numeric.c")
	checkScopeExists(t, tokens, "string.quoted.double.c")
	checkScopeExists(t, tokens, "string.quoted.single.c")
	checkScopeExists(t, tokens, "operator.c")

	typeFound := false
	for _, tok := range tokens {
		if tok.Scope == "entity.name.type.c" {
			typeFound = true
			break
		}
	}
	if !typeFound {
		t.Error("expected entity.name.type.c token")
	}
}

func TestParse_BlockComment(t *testing.T) {
	t.Parallel()

	code := `/* block comment */
int x = 1;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.block.c")
}

func TestParse_ConditionalCompilation(t *testing.T) {
	t.Parallel()

	code := `#ifndef APP_MODE
#define APP_MODE 1
#endif`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "keyword.control.directive.c")
}

func TestParse_EmptyString(t *testing.T) {
	t.Parallel()

	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_UnclosedComment(t *testing.T) {
	t.Parallel()

	code := `/* unclosed comment`
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "comment.block.c" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected comment tokens")
	}
}

func TestParse_LineTracking(t *testing.T) {
	t.Parallel()

	code := "int a = 1;\nint b = 2;"
	rows := Language.Parse(code)
	if len(rows) < 2 {
		t.Fatalf("expected at least 2 rows, got %d", len(rows))
	}
}

func TestAdapter_ID(t *testing.T) {
	t.Parallel()

	if Language.ID() != "c" {
		t.Errorf("expected ID 'c', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	expected := []string{"h", "c89", "c99", "c11"}
	for _, exp := range expected {
		found := false
		for _, a := range aliases {
			if a == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected alias %q, got %v", exp, aliases)
		}
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
		{"representative", `int main() { return 0; } // c`, []seqTok{
			{"int", "support.type.builtin.c"},
			{" ", "default"},
			{"main(", "entity.name.function.c"},
			{")", "operator.c"},
			{" ", "default"},
			{"{", "operator.c"},
			{" ", "default"},
			{"return", "keyword.control.c"},
			{" ", "default"},
			{"0", "constant.numeric.c"},
			{";", "operator.c"},
			{" ", "default"},
			{"}", "operator.c"},
			{" ", "default"},
			{"// c", "comment.line.double-slash.c"},
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
