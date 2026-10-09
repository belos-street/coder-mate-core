package sql

import (
	"strings"
	"testing"

	"code-mate-core/core"
)

func TestParse_BasicQuery(t *testing.T) {
	code := `SELECT u.id, u.name, COUNT(o.id) AS order_count
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
WHERE u.active = TRUE AND o.total >= 99.5
GROUP BY u.id, u.name
ORDER BY order_count DESC
LIMIT 10;`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "keyword.control.sql")
	checkScopeExists(t, tokens, "support.function.builtin.sql")
	checkScopeExists(t, tokens, "entity.name.table.sql")
	checkScopeExists(t, tokens, "entity.name.column.sql")
	checkScopeExists(t, tokens, "constant.language.boolean.sql")
	checkScopeExists(t, tokens, "constant.numeric.sql")
	checkScopeExists(t, tokens, "keyword.operator.sql")
	checkScopeExists(t, tokens, "punctuation.separator.comma.sql")
	checkScopeExists(t, tokens, "punctuation.terminator.statement.sql")
}

func TestParse_CommentsAndStrings(t *testing.T) {
	code := `-- fetch one user
SELECT "display_name"
FROM "users"
WHERE id = $1
  AND email = :email
  AND status = @status
  AND nickname = ?
  AND deleted_at IS NULL
  AND remark = 'hello sql'
/* tail */`
	tokens := flatten(Language.Parse(code))

	checkScopeExists(t, tokens, "comment.line.double-dash.sql")
	checkScopeExists(t, tokens, "comment.block.sql")
	checkScopeExists(t, tokens, "string.quoted.single.sql")
	checkScopeExists(t, tokens, "string.quoted.double.sql")
	checkScopeExists(t, tokens, "constant.language.null.sql")

	paramCount := 0
	for _, tok := range tokens {
		if tok.Scope == "variable.parameter.sql" {
			paramCount++
		}
	}
	if paramCount == 0 {
		t.Error("expected at least one variable.parameter.sql token")
	}
}

func TestParse_UnclosedCommentAndString(t *testing.T) {
	code := `SELECT 'hello
/* not closed`
	tokens := flatten(Language.Parse(code))

	if len(tokens) == 0 {
		t.Fatal("expected tokens for unclosed comment/string")
	}
	found := false
	for _, tok := range tokens {
		if tok.Scope == "string.quoted.single.sql" ||
			tok.Scope == "comment.block.sql" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected string or comment tokens")
	}
}

func TestParse_LineTracking(t *testing.T) {
	code := `SELECT id
FROM users
WHERE id = 1;`
	rows := Language.Parse(code)
	if len(rows) < 3 {
		t.Fatalf("expected at least 3 rows, got %d", len(rows))
	}

	whereFound := false
	for _, row := range rows {
		for _, tok := range row {
			if tok.Scope == "keyword.control.sql" && tok.Text == "WHERE" {
				whereFound = true
				if tok.Line != 3 {
					t.Errorf("expected WHERE on line 3, got line %d", tok.Line)
				}
			}
		}
	}
	if !whereFound {
		t.Error("expected WHERE keyword")
	}
}

func TestParse_EmptyString(t *testing.T) {
	rows := Language.Parse("")
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty input, got %d", len(rows))
	}
}

func TestParse_TableNameQuoted(t *testing.T) {
	code := `SELECT * FROM "users" WHERE id = 1;`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "entity.name.table.sql")
}

func TestParse_JoinWithAlias(t *testing.T) {
	code := `SELECT * FROM users AS u JOIN orders AS o ON u.id = o.user_id`
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "keyword.control.sql")
	checkScopeExists(t, tokens, "entity.name.table.sql")
	checkScopeExists(t, tokens, "entity.name.column.sql")
}

func TestParse_BacktickIdentifiers(t *testing.T) {
	code := "SELECT `first_name` FROM `users`"
	tokens := flatten(Language.Parse(code))
	checkScopeExists(t, tokens, "entity.name.column.sql")
}

func TestAdapter_ID(t *testing.T) {
	if Language.ID() != "sql" {
		t.Errorf("expected ID 'sql', got %q", Language.ID())
	}
}

func TestAdapter_Aliases(t *testing.T) {
	aliases := Language.Aliases()
	expected := []string{"postgresql", "pgsql"}
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
