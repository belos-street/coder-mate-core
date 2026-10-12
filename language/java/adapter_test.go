package java

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_FunctionCallNotType(t *testing.T) {
	t.Parallel()

	// 2.4 回归：此前 EntityType 规则先于 EntityFunction，Foo( 恒被判为类型名
	tokens := flatten(Language.Parse("Foo(1);"))
	if len(tokens) == 0 || tokens[0].Text != "Foo(" {
		t.Fatalf("expected first token 'Foo(', got %v", tokens)
	}
	if tokens[0].Scope != ScopeEntityFunction {
		t.Errorf("expected scope %q, got %q", ScopeEntityFunction, tokens[0].Scope)
	}
}

func TestParse_FunctionCallWithSpace(t *testing.T) {
	// 阶段 4 nit：'Foo (' 形式的函数调用此前欠匹配、恒判类型名
	tokens := flatten(Language.Parse("Foo (1);"))
	found := false
	for _, tok := range tokens {
		if tok.Scope == ScopeEntityFunction {
			found = true
		}
	}
	if !found {
		t.Errorf("expected function scope for 'Foo (', got %+v", tokens)
	}
}

func TestParse_BasicSyntax(t *testing.T) {
	t.Parallel()

	code := `package com.example;

import java.util.List;

/**
 * Javadoc comment
 */
public class HelloWorld {
  private String name;

  public HelloWorld(String name) {
    this.name = name;
  }

  public String greet() {
    return "Hello, " + name + "!";
  }

  public static void main(String[] args) {
    int answer = 42;
    boolean done = false;
    Object obj = null;
    HelloWorld hw = new HelloWorld("World");
    System.out.println(hw.greet());
  }
}`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "keyword.declaration.java")
	checkScopeExists(t, tokens, "keyword.modifier.java")
	checkScopeExists(t, tokens, "keyword.control.java")
	checkScopeExists(t, tokens, "support.type.builtin.java")
	checkScopeExists(t, tokens, "entity.name.function.java")
	checkScopeExists(t, tokens, "constant.language.boolean.java")
	checkScopeExists(t, tokens, "constant.language.null.java")
	checkScopeExists(t, tokens, "constant.numeric.java")
	checkScopeExists(t, tokens, "string.quoted.double.java")
	checkScopeExists(t, tokens, "operator.java")
	checkScopeExists(t, tokens, "entity.name.type.java")
}

func TestParse_BlockComment(t *testing.T) {
	t.Parallel()

	code := `/* block */ int x = 1;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.block.java")
}

func TestParse_LineComment(t *testing.T) {
	t.Parallel()

	code := "// line comment\nint x = 1;"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.line.double-slash.java")
}

func TestParse_Annotation(t *testing.T) {
	t.Parallel()

	code := `@Override
public String toString() { return ""; }`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.annotation.java")
}

func TestParse_NullValues(t *testing.T) {
	t.Parallel()

	code := `Object o = null;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "constant.language.null.java")
}

func TestParse_ThisSuper(t *testing.T) {
	t.Parallel()

	code := `return this;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "variable.language.java")
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

	code := `String s = "unclosed`
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "string.quoted.double.java" {
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

	if Language.ID() != "java" {
		t.Errorf("expected ID 'java', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	found := false
	for _, a := range aliases {
		if a == "jav" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected alias 'jav', got %v", aliases)
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
		{"representative", `class A { int x = 1; } // c`, []seqTok{
			{"class", "keyword.declaration.java"},
			{" ", "default"},
			{"A", "entity.name.type.java"},
			{" ", "default"},
			{"{", "operator.java"},
			{" ", "default"},
			{"int", "support.type.builtin.java"},
			{" ", "default"},
			{"x", "default"},
			{" ", "default"},
			{"=", "operator.java"},
			{" ", "default"},
			{"1", "constant.numeric.java"},
			{";", "operator.java"},
			{" ", "default"},
			{"}", "operator.java"},
			{" ", "default"},
			{"// c", "comment.line.double-slash.java"},
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
