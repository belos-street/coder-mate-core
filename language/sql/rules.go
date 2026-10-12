package sql

import (
	"regexp"

	"code-mate-core/core"
)

var sqlKeywords = `(?i)^(select|from|where|group|by|order|having|limit|offset|insert|into|values|update|set|delete|join|left|right|full|outer|inner|on|as|distinct|union|all|create|alter|drop|table|view|index|truncate|case|when|then|else|end|with|returning|primary|key|foreign|references|constraint|and|or|not|in|like|between|exists|is|asc|desc)\b`

var sqlFunctions = `(?i)^(count|sum|avg|min|max|coalesce|concat|substring|length|lower|upper|now|current_date|current_timestamp|date_trunc|round|cast|json_extract|jsonb_extract_path_text)\b`

var sqlOperators = `(?i)^(?:<>|!=|<=|>=|==|:=|::|->>|->|&&|\|\||[-+*/%<>=])`

var sqlTableStart = `(?i)^(from|join|into|update|table|truncate)\b`

var sqlTableName = `^(?:` +
	`"[^"\r\n]+"|` +
	"\x60[^\x60\r\n]+\x60|" +
	`[A-Za-z_][A-Za-z0-9_$]*)` +
	`(?:\.(?:` +
	`"[^"\r\n]+"|` +
	"\x60[^\x60\r\n]+\x60|" +
	`[A-Za-z_][A-Za-z0-9_$]*))*`

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^--[^\r\n]*`), Scope: ScopeCommentLine},
		{Regex: regexp.MustCompile(`^/\*`), Scope: ScopeCommentBlock, PushState: GrammarStateCommentBlock},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringQuotedSingle, PushState: GrammarStateStringSingle},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringQuotedDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(sqlTableStart), Scope: ScopeKeywordControl, PushState: GrammarStateExpectTableName},
		{Regex: regexp.MustCompile(sqlKeywords), Scope: ScopeKeywordControl},
		{Regex: regexp.MustCompile(sqlFunctions), Scope: ScopeSupportFunction},
		{Regex: regexp.MustCompile(`(?i)^(?:true|false)\b`), Scope: ScopeConstantBoolean},
		{Regex: regexp.MustCompile(`(?i)^null\b`), Scope: ScopeConstantNull},
		{Regex: regexp.MustCompile(`^(?:\?[0-9]*|:[A-Za-z_][A-Za-z0-9_]*|@[A-Za-z_][A-Za-z0-9_]*|\$[0-9]+)`), Scope: ScopeVariableParameter},
		{Regex: regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^,`), Scope: ScopePunctuationSeparatorComma},
		{Regex: regexp.MustCompile(`^;`), Scope: ScopePunctuationTerminator},
		{Regex: regexp.MustCompile(`^\(`), Scope: ScopePunctuationGroupBegin},
		{Regex: regexp.MustCompile(`^\)`), Scope: ScopePunctuationGroupEnd},
		{Regex: regexp.MustCompile(sqlOperators), Scope: ScopeKeywordOperator},
		{Regex: regexp.MustCompile("^\x60[^\x60\r\n]+\x60"), Scope: ScopeEntityColumn},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_$]*`), Scope: ScopeEntityColumn},
	},
	GrammarStateExpectTableName: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(sqlTableName), Scope: ScopeEntityTable, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
	},
	GrammarStateStringSingle: {
		{Regex: regexp.MustCompile(`^[^'\\]*(?:\\.[^'\\]*)*'`), Scope: ScopeStringQuotedSingle, PopState: true},
		{Regex: regexp.MustCompile(`^.*`), Scope: ScopeStringQuotedSingle},
	},
	GrammarStateStringDouble: {
		{Regex: regexp.MustCompile(`^[^"\\]*(?:\\.[^"\\]*)*"`), Scope: ScopeStringQuotedDouble, PopState: true},
		{Regex: regexp.MustCompile(`^.*`), Scope: ScopeStringQuotedDouble},
	},
	GrammarStateCommentBlock: {
		{Regex: regexp.MustCompile(`^[\s\S]*?\*/`), Scope: ScopeCommentBlock, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]+`), Scope: ScopeCommentBlock},
	},
}
