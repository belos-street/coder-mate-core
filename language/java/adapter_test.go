package java

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
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
	code := `/* block */ int x = 1;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.block.java")
}

func TestParse_LineComment(t *testing.T) {
	code := "// line comment\nint x = 1;"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.line.double-slash.java")
}

func TestParse_Annotation(t *testing.T) {
	code := `@Override
public String toString() { return ""; }`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.annotation.java")
}

func TestParse_NullValues(t *testing.T) {
	code := `Object o = null;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "constant.language.null.java")
}

func TestParse_ThisSuper(t *testing.T) {
	code := `return this;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "variable.language.java")
}

func TestParse_EmptyString(t *testing.T) {
	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_UnclosedString(t *testing.T) {
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
	if Language.ID() != "java" {
		t.Errorf("expected ID 'java', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
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
