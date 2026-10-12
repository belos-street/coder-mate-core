package css

import (
	"regexp"

	"code-mate-core/core"
)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^/\*`), Scope: ScopeCommentBlock, PushState: GrammarStateComment},
		{Regex: regexp.MustCompile(`^@[A-Za-z-]+\b`), Scope: ScopeKeywordAtRule},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringQuotedDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringQuotedSingle, PushState: GrammarStateStringSingle},
		{Regex: regexp.MustCompile(`^\.[A-Za-z_-][A-Za-z0-9_-]*`), Scope: ScopeEntityClass},
		{Regex: regexp.MustCompile(`^#[A-Za-z_-][A-Za-z0-9_-]*`), Scope: ScopeEntityId},
		{Regex: regexp.MustCompile(`^::[A-Za-z-]+(?:\([^()\r\n]*\))?`), Scope: ScopePseudoElement},
		{Regex: regexp.MustCompile(`^:[A-Za-z-]+(?:\([^()\r\n]*\))?`), Scope: ScopePseudoClass},
		{Regex: regexp.MustCompile(`^\[`), Scope: ScopeAttrBegin},
		{Regex: regexp.MustCompile(`^\]`), Scope: ScopeAttrEnd},
		{Regex: regexp.MustCompile(`^,`), Scope: ScopeSeparatorSelector},
		{Regex: regexp.MustCompile(`^\{`), Scope: ScopeBlockBegin, PushState: GrammarStateBlock},
		{Regex: regexp.MustCompile(`^;`), Scope: ScopeTerminatorRule},
		{Regex: regexp.MustCompile(`^:`), Scope: ScopeSeparatorKV},
		{Regex: regexp.MustCompile(`(?i)^!important\b`), Scope: ScopeKeywordImportant},
		{Regex: regexp.MustCompile(`^#[0-9A-Fa-f]{3,8}\b`), Scope: ScopeConstantColorHex},
		{Regex: regexp.MustCompile(`^-?(?:\d*\.\d+|\d+)(?:[A-Za-z%]+)?`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^-?[A-Za-z_][A-Za-z0-9_-]*\(`), Scope: ScopeSupportFunction},
		{Regex: regexp.MustCompile(`^(?:>|\+|~|\|\||\||\/|\*|=)`), Scope: ScopeOperator},
		{Regex: regexp.MustCompile(`^(?:\*|[A-Za-z][A-Za-z0-9-]*)`), Scope: ScopeEntityTag},
	},
	GrammarStateBlock: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^/\*`), Scope: ScopeCommentBlock, PushState: GrammarStateComment},
		{Regex: regexp.MustCompile(`^\}`), Scope: ScopeBlockEnd, PopState: true},
		{Regex: regexp.MustCompile(`^\{`), Scope: ScopeBlockBegin, PushState: GrammarStateBlock},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringQuotedDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringQuotedSingle, PushState: GrammarStateStringSingle},
		{Regex: regexp.MustCompile(`^@[A-Za-z-]+\b`), Scope: ScopeKeywordAtRule},
		{Regex: regexp.MustCompile(`(?i)^!important\b`), Scope: ScopeKeywordImportant},
		{Regex: regexp.MustCompile(`^--[A-Za-z_][A-Za-z0-9_-]*\s*:`), Scope: ScopeCustomProperty},
		{Regex: regexp.MustCompile(`^-?[A-Za-z_][A-Za-z0-9_-]*\s*:`), Scope: ScopeSupportPropertyName, Skip: false},
		{Regex: regexp.MustCompile(`^:`), Scope: ScopeSeparatorKV},
		{Regex: regexp.MustCompile(`^;`), Scope: ScopeTerminatorRule},
		{Regex: regexp.MustCompile(`^#[0-9A-Fa-f]{3,8}\b`), Scope: ScopeConstantColorHex},
		{Regex: regexp.MustCompile(`^\.[A-Za-z_-][A-Za-z0-9_-]*`), Scope: ScopeEntityClass},
		{Regex: regexp.MustCompile(`^#[A-Za-z_-][A-Za-z0-9_-]*`), Scope: ScopeEntityId},
		{Regex: regexp.MustCompile(`^::[A-Za-z-]+(?:\([^()\r\n]*\))?`), Scope: ScopePseudoElement},
		{Regex: regexp.MustCompile(`^:[A-Za-z-]+(?:\([^()\r\n]*\))?`), Scope: ScopePseudoClass},
		{Regex: regexp.MustCompile(`^\[`), Scope: ScopeAttrBegin},
		{Regex: regexp.MustCompile(`^\]`), Scope: ScopeAttrEnd},
		{Regex: regexp.MustCompile(`^,`), Scope: ScopeSeparatorSelector},
		{Regex: regexp.MustCompile(`^-?(?:\d*\.\d+|\d+)(?:[A-Za-z%]+)?`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^-?[A-Za-z_][A-Za-z0-9_-]*\(`), Scope: ScopeSupportFunction},
		{Regex: regexp.MustCompile(`^(?:>|\+|~|\|\||\||\/|\*|=|-)`), Scope: ScopeOperator},
	},
	GrammarStateComment: {
		{Regex: regexp.MustCompile(`^[\s\S]*?\*/`), Scope: ScopeCommentBlock, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]+`), Scope: ScopeCommentBlock},
	},
	GrammarStateStringDouble: {
		{Regex: regexp.MustCompile(`^[^\\"]*(?:\\.[^\\"]*)*"`), Scope: ScopeStringQuotedDouble, PopState: true},
		{Regex: regexp.MustCompile(`^.*`), Scope: ScopeStringQuotedDouble},
	},
	GrammarStateStringSingle: {
		{Regex: regexp.MustCompile(`^[^\\']*(?:\\.[^\\']*)*'`), Scope: ScopeStringQuotedSingle, PopState: true},
		{Regex: regexp.MustCompile(`^.*`), Scope: ScopeStringQuotedSingle},
	},
}
