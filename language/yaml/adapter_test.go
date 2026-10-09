package yaml

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_NegativeNumberSequenceItem(t *testing.T) {
	t.Parallel()

	// 2.4 回归：此前 numeric 规则先于 '-'，-5 被整体吃成数字
	tokens := flatten(Language.Parse("-5"))
	if len(tokens) < 2 {
		t.Fatalf("expected at least 2 tokens, got %v", tokens)
	}
	if tokens[0].Text != "-" || tokens[0].Scope != ScopePunctuationSequenceItem {
		t.Errorf("expected '-' as sequence item punctuation, got %q (%q)", tokens[0].Text, tokens[0].Scope)
	}
	if tokens[1].Text != "5" || tokens[1].Scope != ScopeConstantNumeric {
		t.Errorf("expected '5' numeric, got %q (%q)", tokens[1].Text, tokens[1].Scope)
	}
}

func TestParse_BasicMappingAndSequence(t *testing.T) {
	t.Parallel()

	code := `name: coder-mate
version: 1
enabled: true
description: "hello"
items:
  - web
  - api
meta: null`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "support.type.property-name.yaml")
	checkScopeExists(t, tokens, "punctuation.definition.sequence.item.yaml")
	checkScopeExists(t, tokens, "constant.numeric.yaml")
	checkScopeExists(t, tokens, "constant.language.boolean.yaml")
	checkScopeExists(t, tokens, "constant.language.null.yaml")
	checkScopeExists(t, tokens, "string.quoted.double.yaml")
	checkScopeExists(t, tokens, "string.unquoted.yaml")
}

func TestParse_CommentsAndDocuments(t *testing.T) {
	t.Parallel()

	code := `---
# app config
defaults: &defaults
  role: !user admin
profile:
  ref: *defaults
...
`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "keyword.control.document.begin.yaml")
	checkScopeExists(t, tokens, "comment.line.number-sign.yaml")
	checkScopeExists(t, tokens, "entity.name.tag.anchor.yaml")
	checkScopeExists(t, tokens, "variable.other.alias.yaml")
	checkScopeExists(t, tokens, "entity.name.tag.yaml")
	checkScopeExists(t, tokens, "keyword.control.document.end.yaml")
}

func TestParse_UnclosedStrings(t *testing.T) {
	t.Parallel()

	code := `name: "coder
title: 'mate`
	tokens := flatten(Language.Parse(code))

	if len(tokens) == 0 {
		t.Fatal("expected tokens for unclosed strings")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "string.quoted.double.yaml" ||
			tok.Scope == "string.quoted.single.yaml" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected quoted string tokens")
	}
}

func TestParse_EmptyString(t *testing.T) {
	t.Parallel()

	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty input, got %d", len(rows))
	}
}

func TestParse_LineTracking(t *testing.T) {
	t.Parallel()

	code := "a: 1\nb: 2\nc: 3"
	rows := Language.Parse(code)
	if len(rows) < 3 {
		t.Fatalf("expected at least 3 rows, got %d", len(rows))
	}
	for _, row := range rows {
		for _, tok := range row {
			if tok.Line < 0 || tok.Line > len(rows) {
				t.Errorf("line number %d out of range", tok.Line)
			}
		}
	}
}

func TestParse_SingleQuotedString(t *testing.T) {
	t.Parallel()

	code := `title: 'hello world'`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.quoted.single.yaml")
}

func TestParse_NestedCollections(t *testing.T) {
	t.Parallel()

	code := `list:
  - [1, 2]
  - {x: 1, y: 2}`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "punctuation.section.flow.yaml")
	checkScopeExists(t, tokens, "constant.numeric.yaml")
}

func TestParse_NullVariants(t *testing.T) {
	t.Parallel()

	code := `a: null
b: ~`
	tokens := flatten(Language.Parse(code))

	nullCount := 0
	for _, tok := range tokens {
		if tok.Scope == "constant.language.null.yaml" {
			nullCount++
		}
	}
	if nullCount == 0 {
		t.Error("expected null tokens")
	}
}

func TestAdapter_ID(t *testing.T) {
	t.Parallel()

	if Language.ID() != "yaml" {
		t.Errorf("expected ID 'yaml', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	found := false
	for _, a := range aliases {
		if a == "yml" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected alias 'yml', got %v", aliases)
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
		{"representative", `k: [1, true] # c`, []seqTok{
			{"k:", "support.type.property-name.yaml"},
			{" ", "default"},
			{"[", "punctuation.section.flow.yaml"},
			{"1", "constant.numeric.yaml"},
			{",", "punctuation.section.flow.yaml"},
			{" ", "default"},
			{"true", "constant.language.boolean.yaml"},
			{"]", "punctuation.section.flow.yaml"},
			{" ", "default"},
			{"# c", "comment.line.number-sign.yaml"},
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
