package c

import (
	"regexp"

	"code-mate-core/core"
)

var keywordsDeclaration = regexp.MustCompile(`^(typedef|struct|union|enum)\b`)
var keywordsModifier = regexp.MustCompile(`^(static|extern|register|auto|const|volatile|restrict|inline|signed|unsigned|short|long)\b`)
var keywordsControl = regexp.MustCompile(`^(if|else|switch|case|default|for|while|do|break|continue|return|goto|sizeof)\b`)
var directive = regexp.MustCompile(`^#\s*(?:include|define|undef|if|ifdef|ifndef|elif|else|endif|pragma|error|warning|line)\b`)
var builtinTypes = regexp.MustCompile(`^(void|char|int|float|double|_Bool|bool|size_t|ptrdiff_t|wchar_t|FILE|int8_t|int16_t|int32_t|int64_t|uint8_t|uint16_t|uint32_t|uint64_t)\b`)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^//.*`), Scope: ScopeCommentLine},
		{Regex: regexp.MustCompile(`^/\*`), Scope: ScopeCommentBlock, PushState: GrammarStateCommentBlock},
		{Regex: directive, Scope: ScopeDirective},
		{Regex: regexp.MustCompile(`^<[^>\r\n]+>`), Scope: ScopeStringAngle},
		{Regex: regexp.MustCompile(`^#[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopePreprocessor},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PushState: GrammarStateStringSingle},
		{Regex: keywordsDeclaration, Scope: ScopeKeywordDeclaration},
		{Regex: keywordsModifier, Scope: ScopeKeywordModifier},
		{Regex: keywordsControl, Scope: ScopeKeywordControl},
		{Regex: builtinTypes, Scope: ScopeSupportBuiltin},
		{Regex: regexp.MustCompile(`^(true|false)\b`), Scope: ScopeConstantBoolean},
		{Regex: regexp.MustCompile(`^NULL\b`), Scope: ScopeConstantNull},
		{Regex: regexp.MustCompile(`^(?:__[A-Z_][A-Z0-9_]*__|[A-Z][A-Z0-9_]+|[A-Z_][A-Z0-9_]+)`), Scope: ScopeEntityMacro},
		{Regex: regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*_t\b|^[A-Z][A-Za-z0-9_]*\b`), Scope: ScopeEntityType},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*\(`), Scope: ScopeEntityFunction},
		{Regex: regexp.MustCompile(`^(?:0[xX][0-9A-Fa-f]+(?:[uUlL]{0,2})?|0[bB][01]+(?:[uUlL]{0,2})?|(?:\d+\.\d*|\.\d+|\d+)(?:[eE][+-]?\d+)?[fFlLuU]{0,3})`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(?:->|==|!=|<=|>=|\+\+|--|&&|\|\||<<|>>|[-+*/%=&|^~!?<>:.;,()[\]{}])`), Scope: ScopeOperator},
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
}
