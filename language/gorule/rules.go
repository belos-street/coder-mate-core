package gorule

import (
	"regexp"

	"code-mate-core/core"
)

var goKeywordsDecl = regexp.MustCompile(`^(package|import|func|var|const|type|struct|interface|map|chan)\b`)
var goKeywordsControl = regexp.MustCompile(`^(if|else|switch|case|default|for|range|select|go|defer|return|break|continue|fallthrough|goto)\b`)
var goBuiltinTypes = regexp.MustCompile(`^(any|comparable|error|string|bool|byte|rune|int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64|uintptr|float32|float64|complex64|complex128)\b`)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^//.*`), Scope: ScopeCommentLine},
		{Regex: regexp.MustCompile(`^/\*`), Scope: ScopeCommentBlock, PushState: GrammarStateCommentBlock},
		{Regex: regexp.MustCompile("^`"), Scope: ScopeStringRaw, PushState: GrammarStateStringRaw},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PushState: GrammarStateStringSingle},
		{Regex: goKeywordsDecl, Scope: ScopeKeywordDeclaration},
		{Regex: goKeywordsControl, Scope: ScopeKeywordControl},
		{Regex: goBuiltinTypes, Scope: ScopeSupportBuiltin},
		{Regex: regexp.MustCompile(`^(true|false)\b`), Scope: ScopeConstantBoolean},
		{Regex: regexp.MustCompile(`^nil\b`), Scope: ScopeConstantNull},
		{Regex: regexp.MustCompile(`^iota\b`), Scope: ScopeConstantIota},
		{Regex: regexp.MustCompile(`^[a-z_][A-Za-z0-9_]*\.`), Scope: ScopeNamespace},
		{Regex: regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*\[`), Scope: ScopeEntityType},
		{Regex: regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*\(`), Scope: ScopeEntityFunction},
		{Regex: regexp.MustCompile(`^-?(?:0[xX][0-9A-Fa-f_]+|0[bB][01_]+|0[oO][0-7_]+|(?:\d[\d_]*)(?:\.\d[\d_]*)?(?:[eE][+-]?\d[\d_]*)?(?:i)?)`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(?::=|==|!=|<=|>=|&&|\|\||<-|\+\+|--|\.\.\.|[-+*/%=&|^~!?<>:.;,()[\]{}])`), Scope: ScopeOperator},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeDefault},
	},
	GrammarStateCommentBlock: {
		{Regex: regexp.MustCompile(`^[\s\S]*?\*/`), Scope: ScopeCommentBlock, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]+`), Scope: ScopeCommentBlock},
	},
	GrammarStateStringDouble: {
		{Regex: regexp.MustCompile(`^[^"\\\n]+`), Scope: ScopeStringDouble},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringDouble},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringDouble, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringDouble},
	},
	GrammarStateStringSingle: {
		{Regex: regexp.MustCompile(`^[^'\\\n]+`), Scope: ScopeStringSingle},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringSingle},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringSingle},
	},
	GrammarStateStringRaw: {
		{Regex: regexp.MustCompile("^[^`]+"), Scope: ScopeStringRaw},
		{Regex: regexp.MustCompile("^`"), Scope: ScopeStringRaw, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringRaw},
	},
}
