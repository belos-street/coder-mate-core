package cpp

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
	t.Parallel()

	code := `#include <iostream>
#include "mylib.hpp"

#define VERSION 42
#define SQUARE(x) ((x)*(x))

#ifdef DEBUG
#pragma message "debug"
#endif

namespace app {
  class Runnable {
  public:
    virtual int run() = 0;
  };

  struct Config {
    int timeout;
    bool enabled;
  };

  template<typename T>
  T max(T a, T b) {
    return (a > b) ? a : b;
  }
}

int main() {
  constexpr int count = 42;
  auto name = "hello";
  char c = 'A';
  bool ok = true;
  int* ptr = nullptr;
  app::Config cfg;
  cfg.enabled = true;
  return 0;
}`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "keyword.control.directive.cpp")
	checkScopeExists(t, tokens, "string.quoted.angle.cpp")
	checkScopeExists(t, tokens, "string.quoted.double.cpp")
	checkScopeExists(t, tokens, "entity.name.macro.cpp")
	checkScopeExists(t, tokens, "keyword.declaration.cpp")
	checkScopeExists(t, tokens, "keyword.modifier.cpp")
	checkScopeExists(t, tokens, "keyword.control.cpp")
	checkScopeExists(t, tokens, "support.type.builtin.cpp")
	checkScopeExists(t, tokens, "entity.name.type.cpp")
	checkScopeExists(t, tokens, "entity.name.function.cpp")
	checkScopeExists(t, tokens, "entity.name.namespace.cpp")
	checkScopeExists(t, tokens, "constant.language.boolean.cpp")
	checkScopeExists(t, tokens, "constant.language.null.cpp")
	checkScopeExists(t, tokens, "constant.numeric.cpp")
	checkScopeExists(t, tokens, "operator.cpp")
}

func TestParse_LineComment(t *testing.T) {
	t.Parallel()

	code := "// comment\nint x = 1;"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.line.double-slash.cpp")
}

func TestParse_BlockComment(t *testing.T) {
	t.Parallel()

	code := `/* block */ int x = 1;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.block.cpp")
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

	code := `auto s = "unclosed`
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "string.quoted.double.cpp" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected string tokens")
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

	if Language.ID() != "cpp" {
		t.Errorf("expected ID 'cpp', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	expected := []string{"c++", "cc", "cxx", "hpp", "hxx", "hh"}
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
		{"representative", `class A { int x = 1; }; // c`, []seqTok{
			{"class", "keyword.declaration.cpp"},
			{" ", "default"},
			{"A", "entity.name.type.cpp"},
			{" ", "default"},
			{"{", "operator.cpp"},
			{" ", "default"},
			{"int", "support.type.builtin.cpp"},
			{" ", "default"},
			{"x", "default"},
			{" ", "default"},
			{"=", "operator.cpp"},
			{" ", "default"},
			{"1", "constant.numeric.cpp"},
			{";", "operator.cpp"},
			{" ", "default"},
			{"}", "operator.cpp"},
			{";", "operator.cpp"},
			{" ", "default"},
			{"// c", "comment.line.double-slash.cpp"},
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
