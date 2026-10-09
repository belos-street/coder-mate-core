package yaml

import (
	"regexp"

	"code-mate-core/core"
)

var GrammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^#.*`), Scope: ScopeComment},
		{Regex: regexp.MustCompile(`^---`), Scope: ScopeKeywordDocumentBegin},
		{Regex: regexp.MustCompile(`^\.\.\.`), Scope: ScopeKeywordDocumentEnd},
		{Regex: regexp.MustCompile(`^"(?:\\.|[^"\\\r\n])*"\s*:`), Scope: ScopeSupportPropertyName},
		{Regex: regexp.MustCompile(`^'(?:[^'\r\n]|'')*'\s*:`), Scope: ScopeSupportPropertyName},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.-]*\s*:`), Scope: ScopeSupportPropertyName},
		{Regex: regexp.MustCompile(`^[\[\]\{\},]`), Scope: ScopePunctuationSectionFlow},
		{Regex: regexp.MustCompile(`^&[A-Za-z0-9_-]+`), Scope: ScopeEntityTagAnchor},
		{Regex: regexp.MustCompile(`^\*[A-Za-z0-9_-]+`), Scope: ScopeVariableAlias},
		{Regex: regexp.MustCompile(`^![A-Za-z0-9!:/._-]+`), Scope: ScopeEntityTag},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringQuotedDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringQuotedSingle, PushState: GrammarStateStringSingle},
		{Regex: regexp.MustCompile(`^(?:true|false|yes|no|on|off)\b`), Scope: ScopeConstantBoolean},
		{Regex: regexp.MustCompile(`^(?:null|~)\b`), Scope: ScopeConstantNull},
		{Regex: regexp.MustCompile(`^-`), Scope: ScopePunctuationSequenceItem},
		{Regex: regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^[^\s\[\]\{\},#][^#\r\n]*`), Scope: ScopeStringUnquoted},
	},
	GrammarStateStringDouble: {
		{Regex: regexp.MustCompile(`^[^"\\\n]+`), Scope: ScopeStringQuotedDouble},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringQuotedDouble},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringQuotedDouble, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringQuotedDouble},
	},
	GrammarStateStringSingle: {
		{Regex: regexp.MustCompile(`^[^'\n]+`), Scope: ScopeStringQuotedSingle},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringQuotedSingle, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringQuotedSingle},
	},
}
