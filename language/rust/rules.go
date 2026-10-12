package rust

import (
	"regexp"

	"code-mate-core/core"
)

var rustKeywordsDecl = regexp.MustCompile(`^(fn|let|const|static|struct|enum|trait|impl|type|mod|use|crate|extern)\b`)
var rustKeywordsModifier = regexp.MustCompile(`^(pub|mut|unsafe|async|move|ref|dyn|where|in)\b`)
var rustKeywordsControl = regexp.MustCompile(`^(if|else|match|for|while|loop|break|continue|return|await|as)\b`)
var rustBuiltinTypes = regexp.MustCompile(`^(bool|char|str|String|i8|i16|i32|i64|i128|isize|u8|u16|u32|u64|u128|usize|f32|f64|Option|Result|Vec|Box|HashMap|BTreeMap)\b`)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^//.*`), Scope: ScopeCommentLineDS},
		{Regex: regexp.MustCompile(`^/\*`), Scope: ScopeCommentBlock, PushState: GrammarStateCommentBlock},
		{Regex: regexp.MustCompile(`^#\s*!\[[^\]\r\n]*\]|^#\s*\[[^\]\r\n]*\]`), Scope: ScopeMetaAttribute},
		{Regex: regexp.MustCompile(`^r#*"`), Scope: ScopeStringRaw, PushState: GrammarStateStringRaw},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PushState: GrammarStateStringSingle},
		{Regex: rustKeywordsDecl, Scope: ScopeKeywordDecl},
		{Regex: rustKeywordsModifier, Scope: ScopeKeywordModifier},
		{Regex: rustKeywordsControl, Scope: ScopeKeywordControl},
		{Regex: rustBuiltinTypes, Scope: ScopeSupportBuiltin},
		{Regex: regexp.MustCompile(`^(true|false)\b`), Scope: ScopeConstantBoolean},
		{Regex: regexp.MustCompile(`^(None)\b`), Scope: ScopeConstantNull},
		{Regex: regexp.MustCompile(`^(self|Self|super|crate)\b`), Scope: ScopeVarLang},
		{Regex: regexp.MustCompile(`^[a-z_][A-Za-z0-9_]*!`), Scope: ScopeEntityMacro},
		{Regex: regexp.MustCompile(`^[a-z_][A-Za-z0-9_]*::`), Scope: ScopeNamespace},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(?:::<[^>]+>)?\(`), Scope: ScopeEntityFunction},
		{Regex: regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*`), Scope: ScopeEntityType},
		{Regex: regexp.MustCompile(`^-?(?:0[xX][0-9A-Fa-f_]+(?:[iu](?:8|16|32|64|128|size))?|0[bB][01_]+(?:[iu](?:8|16|32|64|128|size))?|0[oO][0-7_]+(?:[iu](?:8|16|32|64|128|size))?|(?:\d[\d_]*)(?:\.\d[\d_]*)?(?:[eE][+-]?\d[\d_]*)?(?:[iu](?:8|16|32|64|128|size)|f32|f64)?)`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(?:::|->|=>|==|!=|<=|>=|&&|\|\||<<|>>|\+=|-=|\*=|\/=|%=|&=|\|=|\^=|[-+*/%=&|^~!?<>:.;,()[\]{}])`), Scope: ScopeOperator},
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
		{Regex: regexp.MustCompile(`^[\s\S]*?"#*`), Scope: ScopeStringRaw, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]+`), Scope: ScopeStringRaw},
	},
}
