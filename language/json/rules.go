package json

import (
	"regexp"

	"code-mate-core/core"
)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringQuotedDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(true|false)\b`), Scope: ScopeConstantBoolean},
		{Regex: regexp.MustCompile(`^null\b`), Scope: ScopeConstantNull},
		{Regex: regexp.MustCompile(`^:`), Scope: ScopePunctuationSeparatorKV},
		{Regex: regexp.MustCompile(`^,`), Scope: ScopePunctuationSeparatorVal},
		{Regex: regexp.MustCompile(`^[{}]`), Scope: ScopeMetaDictionary},
		{Regex: regexp.MustCompile(`^[\[\]]`), Scope: ScopeMetaArray},
	},
	GrammarStateStringDouble: {
		{Regex: regexp.MustCompile(`^[^\\"]*(?:\\.[^\\"]*)*"`), Scope: ScopeStringQuotedDouble, PopState: true},
		{Regex: regexp.MustCompile(`^.*`), Scope: ScopeStringQuotedDouble},
	},
}
