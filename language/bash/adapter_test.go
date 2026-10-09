package bash

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicScript(t *testing.T) {
	t.Parallel()

	code := `#!/usr/bin/env bash
set -euo pipefail
name="coder"
count=3
echo "hi $name" | sed 's/hi/hello/'
if [ "$count" -gt 0 ]; then
  export PATH="$PATH:/tmp/bin"
fi`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "comment.line.number-sign.bash")
	checkScopeExists(t, tokens, "keyword.control.bash")
	checkScopeExists(t, tokens, "support.function.builtin.bash")
	checkScopeExists(t, tokens, "variable.other.readwrite.bash")
	checkScopeExists(t, tokens, "variable.parameter.bash")
	checkScopeExists(t, tokens, "constant.numeric.bash")
	checkScopeExists(t, tokens, "string.quoted.double.bash")
	checkScopeExists(t, tokens, "string.quoted.single.bash")
	checkScopeExists(t, tokens, "operator.bash")
}

func TestParse_FunctionDeclaration(t *testing.T) {
	t.Parallel()

	code := `function greet() {
  local user="${1:-guest}"
  printf '%s\n' "$user"
  echo $?
}`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "keyword.declaration.function.bash")
	checkScopeExists(t, tokens, "entity.name.function.bash")
	checkScopeExists(t, tokens, "variable.language.special.bash")
}

func TestParse_UnclosedStrings(t *testing.T) {
	t.Parallel()

	code := "echo \"hello\necho 'world"
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "string.quoted.double.bash" || tok.Scope == "string.quoted.single.bash" {
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

	code := "echo ok\ncount=1\nif [ \"$count\" -gt 0 ]; then\n  echo done\nfi"
	rows := Language.Parse(code)
	if len(rows) < 5 {
		t.Fatalf("expected at least 5 rows, got %d", len(rows))
	}
}

func TestParse_EmptyString(t *testing.T) {
	t.Parallel()

	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestAdapter_ID(t *testing.T) {
	t.Parallel()

	if Language.ID() != "bash" {
		t.Errorf("expected ID 'bash', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	expected := []string{"sh", "shell"}
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
		{"representative", `if [ $x -gt 2 ]; then echo "hi"; fi`, []seqTok{
			{"if", "keyword.control.bash"},
			{" ", "default"},
			{"[", "operator.bash"},
			{" ", "default"},
			{"$x", "variable.parameter.bash"},
			{" ", "default"},
			{"-gt", "entity.name.command.bash"},
			{" ", "default"},
			{"2", "constant.numeric.bash"},
			{" ", "default"},
			{"]", "operator.bash"},
			{";", "operator.bash"},
			{" ", "default"},
			{"then", "keyword.control.bash"},
			{" ", "default"},
			{"echo", "support.function.builtin.bash"},
			{" ", "default"},
			{"\"", "string.quoted.double.bash"},
			{"hi", "string.quoted.double.bash"},
			{"\"", "string.quoted.double.bash"},
			{";", "operator.bash"},
			{" ", "default"},
			{"fi", "keyword.control.bash"},
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
