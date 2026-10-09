package markdown

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicFormatting(t *testing.T) {
	t.Parallel()

	code := `# Hello World

This is **bold** and *italic* and ~~strikethrough~~.

- item 1
- item 2

> blockquote
`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "markup.heading.markdown")
	checkScopeExists(t, tokens, "markup.bold.markdown")
	checkScopeExists(t, tokens, "markup.italic.markdown")
	checkScopeExists(t, tokens, "markup.strikethrough.markdown")
	checkScopeExists(t, tokens, "markup.list.markdown")
	checkScopeExists(t, tokens, "markup.quote.markdown")
}

func TestParse_CodeBlocks(t *testing.T) {
	t.Parallel()

	code := "```go\nfunc main() {}\n```"
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "markup.fenced-code.block.begin.markdown")
	checkScopeExists(t, tokens, "markup.fenced-code.block.end.markdown")
}

func TestParse_InlineCode(t *testing.T) {
	t.Parallel()

	code := "Use `code` inline"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "markup.inline.raw.markdown")
}

func TestParse_LinksAndImages(t *testing.T) {
	t.Parallel()

	code := `[link text](https://example.com)
![alt text](image.png)`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "markup.link.label.markdown")
	checkScopeExists(t, tokens, "markup.link.url.markdown")
	checkScopeExists(t, tokens, "markup.image.label.markdown")
}

func TestParse_HorizontalRule(t *testing.T) {
	t.Parallel()

	code := `---
***`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "markup.hr.markdown")
}

func TestParse_TaskList(t *testing.T) {
	t.Parallel()

	code := `- [x] done
- [ ] todo`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "markup.task.markdown")
}

func TestParse_Table(t *testing.T) {
	t.Parallel()

	code := `| A | B |
| --- | --- |
| 1 | 2 |`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "markup.table.row.markdown")
	checkScopeExists(t, tokens, "markup.table.separator.markdown")
}

func TestParse_EmptyString(t *testing.T) {
	t.Parallel()

	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_UnclosedFencedCode(t *testing.T) {
	t.Parallel()

	code := "```\nhello world"
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens for unclosed fenced code")
	}
}

func TestAdapter_ID(t *testing.T) {
	t.Parallel()

	if Language.ID() != "markdown" {
		t.Errorf("expected ID 'markdown', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	found := false
	for _, a := range aliases {
		if a == "md" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected alias 'md', got %v", aliases)
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
		{"representative", "*t* `c`", []seqTok{
			{"*", "markup.italic.markdown"},
			{"t", "markup.italic.markdown"},
			{"*", "markup.italic.markdown"},
			{" ", "default"},
			{"`", "markup.inline.raw.markdown"},
			{"c", "markup.inline.raw.markdown"},
			{"`", "markup.inline.raw.markdown"},
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
