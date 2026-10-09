package csharp

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
	t.Parallel()

	code := `using System;
using System.Collections.Generic;

namespace MyApp {
  [Serializable]
  public class Program {
    private string name;

    public Program(string name) {
      var x = this;
      this.name = name;
    }

    public string Greet() {
      var msg = "Hello, " + name + "!";
      return msg;
    }

    public async Task<int> ComputeAsync() {
      int result = 42;
      bool done = true;
      object data = null;
      var list = new List<int>();
      list.Add(result);
      return result;
    }
  }
}`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "keyword.declaration.csharp")
	checkScopeExists(t, tokens, "keyword.modifier.csharp")
	checkScopeExists(t, tokens, "keyword.control.csharp")
	checkScopeExists(t, tokens, "support.type.builtin.csharp")
	checkScopeExists(t, tokens, "entity.name.type.csharp")
	checkScopeExists(t, tokens, "entity.name.function.csharp")
	checkScopeExists(t, tokens, "variable.language.csharp")
	checkScopeExists(t, tokens, "constant.language.boolean.csharp")
	checkScopeExists(t, tokens, "constant.language.null.csharp")
	checkScopeExists(t, tokens, "constant.numeric.csharp")
	checkScopeExists(t, tokens, "string.quoted.double.csharp")
	checkScopeExists(t, tokens, "operator.csharp")
}

func TestParse_LineComment(t *testing.T) {
	t.Parallel()

	code := "// comment\nint x = 1;"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.line.double-slash.csharp")
}

func TestParse_BlockComment(t *testing.T) {
	t.Parallel()

	code := `/* block */ int x = 1;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.block.csharp")
}

func TestParse_Attribute(t *testing.T) {
	t.Parallel()

	code := `[Obsolete]
public void Old() {}`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.attribute.csharp")
}

func TestParse_Preprocessor(t *testing.T) {
	t.Parallel()

	code := `#if DEBUG
Console.WriteLine("debug");
#endif`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.preprocessor.csharp")
}

func TestParse_VerbatimString(t *testing.T) {
	t.Parallel()

	code := `var path = @"C:\Users\name";`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.quoted.verbatim.csharp")
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

	code := `var s = "unclosed`
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "string.quoted.double.csharp" {
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

	if Language.ID() != "csharp" {
		t.Errorf("expected ID 'csharp', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	expected := []string{"cs", "c#"}
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
		{"representative", `var s = "hi"; int x = 1; // c`, []seqTok{
			{"var", "support.type.builtin.csharp"},
			{" ", "default"},
			{"s", "default"},
			{" ", "default"},
			{"=", "operator.csharp"},
			{" ", "default"},
			{"\"", "string.quoted.double.csharp"},
			{"hi", "string.quoted.double.csharp"},
			{"\"", "string.quoted.double.csharp"},
			{";", "operator.csharp"},
			{" ", "default"},
			{"int", "support.type.builtin.csharp"},
			{" ", "default"},
			{"x", "default"},
			{" ", "default"},
			{"=", "operator.csharp"},
			{" ", "default"},
			{"1", "constant.numeric.csharp"},
			{";", "operator.csharp"},
			{" ", "default"},
			{"// c", "comment.line.double-slash.csharp"},
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
