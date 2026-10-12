package bash

import (
	"regexp"

	"code-mate-core/core"
)

var bashKeywords = regexp.MustCompile(`^(if|then|else|elif|fi|for|while|until|do|done|case|esac|in|select|time|coproc)\b`)
var bashBuiltins = regexp.MustCompile(`^(echo|printf|read|cd|pwd|export|unset|alias|unalias|source|set|shift|test|true|false|trap|exec|exit|return|local|declare|typeset|readonly|mapfile|wait|kill|jobs|fg|bg)\b`)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^#.*`), Scope: ScopeCommentLine},
		{Regex: regexp.MustCompile(`^(function)\b`), Scope: ScopeKeywordDeclFunc, PushState: GrammarStateExpectFuncName},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*\s*\(`), Scope: ScopeEntityFunction},
		{Regex: bashKeywords, Scope: ScopeKeywordControl},
		{Regex: bashBuiltins, Scope: ScopeSupportBuiltin},
		{Regex: regexp.MustCompile(`^\$(?:[@*#?$!0-9]|-[A-Za-z?])`), Scope: ScopeVarSpecial},
		{Regex: regexp.MustCompile(`^\$\{[A-Za-z0-9_]+(?::[-=?+][^}]*)?\}`), Scope: ScopeVarParameter},
		{Regex: regexp.MustCompile(`^\$[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeVarParameter},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`), Scope: ScopeVarReadWrite},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PushState: GrammarStateStringSingle},
		{Regex: regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(?:\|\||&&|\|&|>>|<<-?|<&|>&|\||>|<|;;|;|&|\(|\)|\{|\}|\[|\]|=)`), Scope: ScopeOperator},
		{Regex: regexp.MustCompile(`^[A-Za-z_./-][A-Za-z0-9_./-]*`), Scope: ScopeEntityCommand},
	},
	GrammarStateExpectFuncName: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeEntityFunction, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
	},
	GrammarStateStringDouble: {
		{Regex: regexp.MustCompile(`^\$(?:[@*#?$!0-9]|-[A-Za-z?])`), Scope: ScopeVarSpecial},
		{Regex: regexp.MustCompile(`^\$\{[A-Za-z0-9_]+(?::[-=?+][^}]*)?\}`), Scope: ScopeVarParameter},
		{Regex: regexp.MustCompile(`^\$[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeVarParameter},
		{Regex: regexp.MustCompile(`^[^\\"$\n]+`), Scope: ScopeStringDouble},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringDouble},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringDouble, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringDouble},
	},
	GrammarStateStringSingle: {
		{Regex: regexp.MustCompile(`^[^'\n]+'`), Scope: ScopeStringSingle, PopState: true},
		{Regex: regexp.MustCompile(`^[^'\n]+`), Scope: ScopeStringSingle},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringSingle},
	},
}
