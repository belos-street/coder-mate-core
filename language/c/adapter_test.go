package c

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
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
	code := `/* block comment */
int x = 1;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.block.c")
}

func TestParse_ConditionalCompilation(t *testing.T) {
	code := `#ifndef APP_MODE
#define APP_MODE 1
#endif`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "keyword.control.directive.c")
}

func TestParse_EmptyString(t *testing.T) {
	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_UnclosedComment(t *testing.T) {
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
	code := "int a = 1;\nint b = 2;"
	rows := Language.Parse(code)
	if len(rows) < 2 {
		t.Fatalf("expected at least 2 rows, got %d", len(rows))
	}
}

func TestAdapter_ID(t *testing.T) {
	if Language.ID() != "c" {
		t.Errorf("expected ID 'c', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
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
