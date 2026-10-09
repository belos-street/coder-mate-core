package css

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_SelectorsAndDeclarations(t *testing.T) {
	code := `.card:hover::before, #app[data-theme="dark"] {
  color: #1f2937;
  margin: 12px 8px;
  --brand-color: #2563eb;
}`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "entity.other.attribute-name.class.css")
	checkScopeExists(t, tokens, "entity.other.attribute-name.id.css")
	checkScopeExists(t, tokens, "entity.other.attribute-name.pseudo-class.css")
	checkScopeExists(t, tokens, "entity.other.attribute-name.pseudo-element.css")
	checkScopeExists(t, tokens, "punctuation.section.block.begin.css")
	checkScopeExists(t, tokens, "punctuation.section.block.end.css")
	checkScopeExists(t, tokens, "support.type.property-name.css")
	checkScopeExists(t, tokens, "variable.parameter.custom-property.css")
	checkScopeExists(t, tokens, "constant.other.color.hex.css")
	checkScopeExists(t, tokens, "constant.numeric.css")
}

func TestParse_CommentsAndFunctions(t *testing.T) {
	code := `/* reset */
@media screen and (min-width: 768px) {
  .box {
    width: calc(100% - 24px);
    background-image: linear-gradient(90deg, #fff, #000);
    color: rgba(0, 0, 0, 0.85) !important;
  }
}`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "comment.block.css")
	checkScopeExists(t, tokens, "keyword.control.at-rule.css")
	checkScopeExists(t, tokens, "support.function.css")
	checkScopeExists(t, tokens, "keyword.other.important.css")
}

func TestParse_UnclosedCommentAndString(t *testing.T) {
	code := `.title { content: "hello
/* not closed`
	tokens := flatten(Language.Parse(code))

	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "string.quoted.double.css" || tok.Scope == "comment.block.css" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected string or comment tokens")
	}
}

func TestParse_EmptyString(t *testing.T) {
	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestParse_TagAndPseudoSelectors(t *testing.T) {
	code := `div:hover { color: red; }`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "entity.name.tag.css")
	checkScopeExists(t, tokens, "entity.other.attribute-name.pseudo-class.css")
	checkScopeExists(t, tokens, "punctuation.section.block.begin.css")
	checkScopeExists(t, tokens, "punctuation.section.block.end.css")
	checkScopeExists(t, tokens, "support.type.property-name.css")
}

func TestParse_AttrSelectors(t *testing.T) {
	code := `input[type="text"] { }`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "entity.name.tag.css")
	checkScopeExists(t, tokens, "punctuation.definition.attribute.begin.css")
	checkScopeExists(t, tokens, "punctuation.definition.attribute.end.css")
	checkScopeExists(t, tokens, "string.quoted.double.css")
}

func TestParse_NestedAtRules(t *testing.T) {
	code := `@media screen {
  @supports (display: flex) {
    .box { display: flex; }
  }
}`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "keyword.control.at-rule.css")
	checkScopeExists(t, tokens, "punctuation.section.block.begin.css")
	checkScopeExists(t, tokens, "punctuation.section.block.end.css")
}

func TestParse_SingleQuotedString(t *testing.T) {
	code := `a[href='test'] { color: red; }`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.quoted.single.css")
}

func TestAdapter_ID(t *testing.T) {
	if Language.ID() != "css" {
		t.Errorf("expected ID 'css', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	aliases := Language.Aliases()
	if len(aliases) != 0 {
		t.Errorf("expected empty aliases, got %v", aliases)
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
