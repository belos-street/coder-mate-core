package json

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_ObjectAndArray(t *testing.T) {
	code := `{
  "name": "coder-mate",
  "count": 3,
  "active": true,
  "meta": null,
  "items": [1, 2]
}`
	rows := Language.Parse(code)
	tokens := flatten(rows)

	if len(rows) == 0 {
		t.Fatal("expected at least 1 row")
	}

	checkScopeExists(t, tokens, "meta.structure.dictionary.json")
	checkScopeExists(t, tokens, "meta.structure.array.json")
	checkScopeExists(t, tokens, "string.quoted.double.json")
	checkScopeExists(t, tokens, "constant.numeric.json")
	checkScopeExists(t, tokens, "constant.language.boolean.json")
	checkScopeExists(t, tokens, "constant.language.null.json")
	checkScopeExists(t, tokens, "punctuation.separator.key-value.json")
	checkScopeExists(t, tokens, "punctuation.separator.value.json")
}

func TestParse_NumericFormats(t *testing.T) {
	code := `{"a": -12.5e+2, "b": 0, "c": 1.23E-4}`
	tokens := flatten(Language.Parse(code))

	foundMinus := false
	foundZero := false
	foundExponent := false
	for _, tok := range tokens {
		if tok.Text == "-12.5e+2" {
			foundMinus = true
		}
		if tok.Text == "0" {
			foundZero = true
		}
		if tok.Text == "1.23E-4" {
			foundExponent = true
		}
	}
	if !foundMinus {
		t.Error("expected token with text '-12.5e+2'")
	}
	if !foundZero {
		t.Error("expected token with text '0'")
	}
	if !foundExponent {
		t.Error("expected token with text '1.23E-4'")
	}
}

func TestParse_EmptyObject(t *testing.T) {
	code := `{}`
	rows := Language.Parse(code)
	if len(rows) == 0 {
		t.Fatal("expected at least 1 row")
	}
}

func TestParse_NestedArrays(t *testing.T) {
	code := `[[1, 2], [3, 4]]`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.structure.array.json")
	checkScopeExists(t, tokens, "constant.numeric.json")
}

func TestParse_EmptyString(t *testing.T) {
	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty input, got %d", len(rows))
	}
}

func TestParse_UnclosedString(t *testing.T) {
	code := `{"key": "unclosed`
	rows := Language.Parse(code)
	if len(rows) == 0 {
		t.Fatal("expected at least 1 row")
	}
}

func TestParse_BooleanValues(t *testing.T) {
	code := `{"a": true, "b": false}`
	tokens := flatten(Language.Parse(code))

	trues := 0
	falses := 0
	for _, tok := range tokens {
		if tok.Text == "true" && tok.Scope == "constant.language.boolean.json" {
			trues++
		}
		if tok.Text == "false" && tok.Scope == "constant.language.boolean.json" {
			falses++
		}
	}
	if trues == 0 {
		t.Error("expected 'true' with boolean scope")
	}
	if falses == 0 {
		t.Error("expected 'false' with boolean scope")
	}
}

func TestParse_NullValue(t *testing.T) {
	code := `{"a": null}`
	tokens := flatten(Language.Parse(code))

	found := false
	for _, tok := range tokens {
		if tok.Text == "null" && tok.Scope == "constant.language.null.json" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'null' with null scope")
	}
}

func TestParse_StringWithEscapes(t *testing.T) {
	code := `{"msg": "hello\nworld\t\"escaped\""}`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.quoted.double.json")
	checkScopeExists(t, tokens, "punctuation.separator.key-value.json")
}

func TestParse_LineTracking(t *testing.T) {
	code := "{\n  \"key\": 42\n}"
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

func TestAdapter_ID(t *testing.T) {
	if Language.ID() != "json" {
		t.Errorf("expected ID 'json', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	aliases := Language.Aliases()
	if len(aliases) != 0 {
		t.Errorf("expected empty aliases, got %v", aliases)
	}
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
