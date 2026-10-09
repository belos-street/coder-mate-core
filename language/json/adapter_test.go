package json

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_ObjectAndArray(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

	code := `{}`
	rows := Language.Parse(code)
	if len(rows) == 0 {
		t.Fatal("expected at least 1 row")
	}
}

func TestParse_NestedArrays(t *testing.T) {
	t.Parallel()

	code := `[[1, 2], [3, 4]]`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.structure.array.json")
	checkScopeExists(t, tokens, "constant.numeric.json")
}

func TestParse_EmptyString(t *testing.T) {
	t.Parallel()

	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty input, got %d", len(rows))
	}
}

func TestParse_UnclosedString(t *testing.T) {
	t.Parallel()

	code := `{"key": "unclosed`
	rows := Language.Parse(code)
	if len(rows) == 0 {
		t.Fatal("expected at least 1 row")
	}
}

func TestParse_BooleanValues(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

	code := `{"msg": "hello\nworld\t\"escaped\""}`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.quoted.double.json")
	checkScopeExists(t, tokens, "punctuation.separator.key-value.json")
}

func TestParse_LineTracking(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	if Language.ID() != "json" {
		t.Errorf("expected ID 'json', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	if len(aliases) != 0 {
		t.Errorf("expected empty aliases, got %v", aliases)
	}
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
		{"representative", `{"k": [1, true, null]}`, []seqTok{
			{"{", "meta.structure.dictionary.json"},
			{"\"", "string.quoted.double.json"},
			{"k\"", "string.quoted.double.json"},
			{":", "punctuation.separator.key-value.json"},
			{" ", "default"},
			{"[", "meta.structure.array.json"},
			{"1", "constant.numeric.json"},
			{",", "punctuation.separator.value.json"},
			{" ", "default"},
			{"true", "constant.language.boolean.json"},
			{",", "punctuation.separator.value.json"},
			{" ", "default"},
			{"null", "constant.language.null.json"},
			{"]", "meta.structure.array.json"},
			{"}", "meta.structure.dictionary.json"},
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
