package php

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicSyntax(t *testing.T) {
	t.Parallel()

	code := `<?php

namespace App\Controller;

use App\Service\UserService;

class UserController {
  private UserService $service;

  public function __construct(UserService $service) {
    $this->service = $service;
  }

  public function show(int $id): void {
    $count = 42;
    $result = $this->service->find($id);
    echo "User: " . $result;
  }
}`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "meta.tag.php")
	checkScopeExists(t, tokens, "keyword.declaration.php")
	checkScopeExists(t, tokens, "keyword.modifier.php")
	checkScopeExists(t, tokens, "keyword.control.php")
	checkScopeExists(t, tokens, "support.type.builtin.php")
	checkScopeExists(t, tokens, "entity.name.type.php")
	checkScopeExists(t, tokens, "entity.name.function.php")
	checkScopeExists(t, tokens, "variable.language.php")
	checkScopeExists(t, tokens, "variable.other.php")
	checkScopeExists(t, tokens, "constant.numeric.php")
	checkScopeExists(t, tokens, "string.quoted.double.php")
	checkScopeExists(t, tokens, "operator.php")
}

func TestParse_Comments(t *testing.T) {
	t.Parallel()

	code := `<?php
// line comment
# hash comment
/* block */`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "comment.line.double-slash.php")
	checkScopeExists(t, tokens, "comment.line.number-sign.php")
	checkScopeExists(t, tokens, "comment.block.php")
}

func TestParse_Attribute(t *testing.T) {
	t.Parallel()

	code := `<?php
#[Route('/api', methods: ['GET'])]
class ApiController {}`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "meta.attribute.php")
}

func TestParse_BooleanNull(t *testing.T) {
	t.Parallel()

	code := `<?php
$x = true;
$y = false;
$z = null;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "constant.language.boolean.php")
	checkScopeExists(t, tokens, "constant.language.null.php")
}

func TestParse_MagicConstants(t *testing.T) {
	t.Parallel()

	code := `<?php
echo __FILE__;
echo __CLASS__;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "constant.language.php")
}

func TestParse_BacktickString(t *testing.T) {
	t.Parallel()

	code := "<?php\n$out = `ls -la`;"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "string.quoted.backtick.php")
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

	if Language.ID() != "php" {
		t.Errorf("expected ID 'php', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	t.Parallel()

	aliases := Language.Aliases()
	expected := []string{"phtml", "php8"}
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
		{"representative", `function f($x = 1) { return "s"; }`, []seqTok{
			{"function", "keyword.declaration.php"},
			{" ", "default"},
			{"f(", "entity.name.function.php"},
			{"$x", "variable.other.php"},
			{" ", "default"},
			{"=", "operator.php"},
			{" ", "default"},
			{"1", "constant.numeric.php"},
			{")", "operator.php"},
			{" ", "default"},
			{"{", "operator.php"},
			{" ", "default"},
			{"return", "keyword.control.php"},
			{" ", "default"},
			{"\"", "string.quoted.double.php"},
			{"s", "string.quoted.double.php"},
			{"\"", "string.quoted.double.php"},
			{";", "operator.php"},
			{" ", "default"},
			{"}", "operator.php"},
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
