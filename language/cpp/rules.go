package cpp

import (
	"regexp"

	"code-mate-core/core"
)

var cppKeywordsDecl = regexp.MustCompile(`^(class|struct|union|enum|template|typename|using|namespace|concept)\b`)
var cppKeywordsModifier = regexp.MustCompile(`^(static|extern|register|auto|const|volatile|mutable|inline|virtual|explicit|friend|thread_local|constexpr|consteval|constinit|signed|unsigned|short|long|final|override)\b`)
var cppKeywordsControl = regexp.MustCompile(`^(if|else|switch|case|default|for|while|do|break|continue|return|try|catch|throw|new|delete|this|sizeof|typeid|noexcept|static_assert|co_await|co_return|co_yield)\b`)
var cppIncludeDirective = regexp.MustCompile(`^#\s*include\b`)
var cppDefineDirective = regexp.MustCompile(`^#\s*define\b`)
var cppMacroRefDirective = regexp.MustCompile(`^#\s*(?:ifdef|ifndef|undef)\b`)
var cppOtherDirective = regexp.MustCompile(`^#\s*(?:if|elif|else|endif|pragma|error|warning|line)\b`)
var cppBuiltinTypes = regexp.MustCompile(`^(void|bool|char|char8_t|char16_t|char32_t|wchar_t|short|int|long|float|double|auto|size_t|ptrdiff_t|string|std::string|int8_t|int16_t|int32_t|int64_t|uint8_t|uint16_t|uint32_t|uint64_t)\b`)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^//.*`), Scope: ScopeCommentLineDS},
		{Regex: regexp.MustCompile(`^/\*`), Scope: ScopeCommentBlock, PushState: GrammarStateCommentBlock},
		{Regex: cppIncludeDirective, Scope: ScopeDirective, PushState: GrammarStatePreprocessorInclude},
		{Regex: cppDefineDirective, Scope: ScopeDirective, PushState: GrammarStatePreprocessorDefine},
		{Regex: cppMacroRefDirective, Scope: ScopeDirective, PushState: GrammarStatePreprocessorMacroRef},
		{Regex: cppOtherDirective, Scope: ScopeDirective},
		{Regex: regexp.MustCompile(`^#[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopePreprocessor},
		{Regex: regexp.MustCompile(`^R"[^\s()\\\r\n]{0,16}\(`), Scope: ScopeStringRaw, PushState: GrammarStateStringRaw},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PushState: GrammarStateStringSingle},
		{Regex: cppKeywordsDecl, Scope: ScopeKeywordDecl},
		{Regex: cppKeywordsModifier, Scope: ScopeKeywordModifier},
		{Regex: cppKeywordsControl, Scope: ScopeKeywordControl},
		{Regex: cppBuiltinTypes, Scope: ScopeSupportBuiltin},
		{Regex: regexp.MustCompile(`^(true|false)\b`), Scope: ScopeConstantBool},
		{Regex: regexp.MustCompile(`^(nullptr|NULL)\b`), Scope: ScopeConstantNull},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*::`), Scope: ScopeNamespace},
		{Regex: regexp.MustCompile(`^(?:__[A-Z_][A-Z0-9_]*__|[A-Z][A-Z0-9_]+|[A-Z_][A-Z0-9_]{2,})`), Scope: ScopeEntityMacro},
		{Regex: regexp.MustCompile(`^[A-Za-z_~][A-Za-z0-9_~]*\(`), Scope: ScopeEntityFunction},
		{Regex: regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*(?:<[^>]*>)?`), Scope: ScopeEntityType},
		{Regex: regexp.MustCompile(`^-?(?:0[xX][0-9A-Fa-f']+(?:[uUlLfFzZ]{0,4})?|0[bB][01']+(?:[uUlL]{0,3})?|(?:\d[\d']*)(?:\.(?:\d[\d']*)?)?(?:[eEpP][+-]?\d[\d']*)?[fFlLuUzZ]{0,4})`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(?:::|->\*|->|\.\*|\.\.\.|<=>|==|!=|<=|>=|<<=|>>=|\+\+|--|&&|\|\||<<|>>|\+=|-=|\*=|\/=|%=|&=|\|=|\^=|[-+*/%=&|^~!?<>:.;,()[\]{}])`), Scope: ScopeOperator},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeDefault},
	},
	GrammarStatePreprocessorInclude: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^<[^>\r\n]+>`), Scope: ScopeStringAngle, PopState: true},
		{Regex: regexp.MustCompile(`^"[^"\r\n]*"`), Scope: ScopeStringDouble, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
	},
	GrammarStatePreprocessorDefine: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeEntityMacro, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
	},
	GrammarStatePreprocessorMacroRef: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeEntityMacro, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
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
		{Regex: regexp.MustCompile(`^[\s\S]*?\)[^"\r\n\\]{0,16}"`), Scope: ScopeStringRaw, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]+`), Scope: ScopeStringRaw},
	},
}
