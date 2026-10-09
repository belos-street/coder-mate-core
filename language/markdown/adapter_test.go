package markdown

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicFormatting(t *testing.T) {
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
	code := "```go\nfunc main() {}\n```"
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "markup.fenced-code.block.begin.markdown")
	checkScopeExists(t, tokens, "markup.fenced-code.block.end.markdown")
}

func TestParse_InlineCode(t *testing.T) {
	code := "Use `code` inline"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "markup.inline.raw.markdown")
}

func TestParse_LinksAndImages(t *testing.T) {
	code := `[link text](https://example.com)
![alt text](image.png)`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "markup.link.label.markdown")
	checkScopeExists(t, tokens, "markup.link.url.markdown")
	checkScopeExists(t, tokens, "markup.image.label.markdown")
}

func TestParse_HorizontalRule(t *testing.T) {
	code := `---
***`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "markup.hr.markdown")
}

func TestParse_TaskList(t *testing.T) {
	code := `- [x] done
- [ ] todo`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "markup.task.markdown")
}

func TestParse_Table(t *testing.T) {
	code := `| A | B |
| --- | --- |
| 1 | 2 |`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "markup.table.row.markdown")
	checkScopeExists(t, tokens, "markup.table.separator.markdown")
}

func TestParse_EmptyString(t *testing.T) {
	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_UnclosedFencedCode(t *testing.T) {
	code := "```\nhello world"
	tokens := flatten(Language.Parse(code))
	if len(tokens) == 0 {
		t.Fatal("expected tokens for unclosed fenced code")
	}
}

func TestAdapter_ID(t *testing.T) {
	if Language.ID() != "markdown" {
		t.Errorf("expected ID 'markdown', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
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
