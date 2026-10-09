package html

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicTags(t *testing.T) {
	t.Parallel()

	code := `<div class="container" id="app">
  <p>Hello world</p>
  <br />
</div>`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "entity.name.tag.html")
	checkScopeExists(t, tokens, "entity.other.attribute-name.html")
	checkScopeExists(t, tokens, "string.quoted.double.html")
	checkScopeExists(t, tokens, "punctuation.definition.tag.begin.html")
	checkScopeExists(t, tokens, "punctuation.definition.tag.end.html")
	checkScopeExists(t, tokens, "text.plain.html")
}

func TestParse_DoctypeAndComment(t *testing.T) {
	t.Parallel()

	code := `<!DOCTYPE html>
<!-- this is a comment -->
<html></html>`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "keyword.control.doctype.html")
	checkScopeExists(t, tokens, "comment.block.html")
}

func TestParse_SelfClosingTag(t *testing.T) {
	t.Parallel()

	code := `<img src="photo.jpg" alt="Photo" />`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "entity.name.tag.html")
	checkScopeExists(t, tokens, "entity.other.attribute-name.html")
	checkScopeExists(t, tokens, "string.quoted.double.html")
}

func TestParse_SingleQuotedAttributes(t *testing.T) {
	t.Parallel()

	code := `<input type='text' value='hello'>`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.quoted.single.html")
}

func TestParse_UnquotedAttributes(t *testing.T) {
	t.Parallel()

	code := `<div class=container>`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.unquoted.html")
}

func TestParse_EntityEncoding(t *testing.T) {
	t.Parallel()

	code := `<p>&amp; &lt; &gt; &#65;</p>`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "constant.character.entity.html")
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

	code := `<!-- not closed`
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "comment.block.html" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected comment tokens")
	}
}

func TestAdapter_ID(t *testing.T) {
	t.Parallel()

	if Language.ID() != "html" {
		t.Errorf("expected ID 'html', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	found := false
	for _, a := range aliases {
		if a == "htm" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected alias 'htm', got %v", aliases)
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
		{"representative", `<!-- c --><p class="a">hi</p>`, []seqTok{
			{"<!--", "comment.block.html"},
			{" c -->", "comment.block.html"},
			{"<p", "entity.name.tag.html"},
			{" ", "default"},
			{"class", "entity.other.attribute-name.html"},
			{"=", "punctuation.separator.key-value.html"},
			{"\"a\"", "string.quoted.double.html"},
			{">", "punctuation.definition.tag.end.html"},
			{"hi", "text.plain.html"},
			{"</", "punctuation.definition.tag.begin.html"},
			{"p", "entity.name.tag.html"},
			{">", "punctuation.definition.tag.end.html"},
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
