package csharp

import (
	"regexp"

	"code-mate-core/core"
)

var csKeywordsDecl = regexp.MustCompile(`^(namespace|using|class|interface|struct|enum|record|delegate|event)\b`)
var csKeywordsModifier = regexp.MustCompile(`^(public|private|protected|internal|static|sealed|abstract|virtual|override|async|readonly|unsafe|partial|required|volatile|extern|new|ref|out|params)\b`)
var csKeywordsControl = regexp.MustCompile(`^(if|else|switch|case|default|for|foreach|while|do|break|continue|return|try|catch|finally|throw|lock|in|is|as|await|yield|from|where|select|group|join|let|orderby|into|when|get|set|init|add|remove|nameof|typeof|sizeof|checked|unchecked)\b`)
var csBuiltinTypes = regexp.MustCompile(`^(void|bool|byte|sbyte|short|ushort|int|uint|long|ulong|nint|nuint|float|double|decimal|char|string|object|dynamic|var|Task|ValueTask|DateTime|Guid|List|Dictionary|IEnumerable|IQueryable|Span|Memory|Func|Action|CancellationToken)\b`)
var csAttributePattern = regexp.MustCompile(`^\[[A-Za-z_][A-Za-z0-9_.]*(?:\s*:\s*[A-Za-z_][A-Za-z0-9_.]*)?(?:\s*\([^\)\r\n]*\))?\]`)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^//.*`), Scope: ScopeCommentLineDS},
		{Regex: regexp.MustCompile(`^/\*`), Scope: ScopeCommentBlock, PushState: GrammarStateCommentBlock},
		{Regex: regexp.MustCompile(`^#\s*(?:if|elif|else|endif|define|undef|region|endregion|pragma|warning|error|line)\b[^\r\n]*`), Scope: ScopePreprocessor},
		{Regex: csAttributePattern, Scope: ScopeMetaAttribute},
		{Regex: regexp.MustCompile(`^\$?"""`), Scope: ScopeStringRaw, PushState: GrammarStateStringRaw},
		{Regex: regexp.MustCompile(`^(?:@"|\$@"|@\$")`), Scope: ScopeStringVerbatim, PushState: GrammarStateStringVerbatim},
		{Regex: regexp.MustCompile(`^\$?"`), Scope: ScopeStringDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PushState: GrammarStateStringSingle},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)+`), Scope: ScopeNamespace},
		{Regex: csKeywordsDecl, Scope: ScopeKeywordDecl},
		{Regex: csKeywordsModifier, Scope: ScopeKeywordModifier},
		{Regex: csKeywordsControl, Scope: ScopeKeywordControl},
		{Regex: csBuiltinTypes, Scope: ScopeSupportBuiltin},
		{Regex: regexp.MustCompile(`^(true|false)\b`), Scope: ScopeConstantBool},
		{Regex: regexp.MustCompile(`^(null|default)\b`), Scope: ScopeConstantNull},
		{Regex: regexp.MustCompile(`^(this|base)\b`), Scope: ScopeVarLang},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*\(`), Scope: ScopeEntityFunction},
		{Regex: regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*(?:<[^>]*>)?`), Scope: ScopeEntityType},
		{Regex: regexp.MustCompile(`^-?(?:0[xX][0-9A-Fa-f_]+(?:[uUlLfFmMdD]{0,2})?|0[bB][01_]+(?:[uUlL]{0,2})?|(?:\d[\d_]*)(?:\.\d[\d_]*)?(?:[eE][+-]?\d[\d_]*)?(?:[uUlLfFmMdD]{0,2})?)`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(?:=>|\?\?|\?\.|::|==|!=|<=|>=|&&|\|\||<<|>>|\+\+|--|\+=|-=|\*=|\/=|%=|&=|\|=|\^=|[-+*/%=&|^~!?<>:.;,()[\]{}])`), Scope: ScopeOperator},
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
	GrammarStateStringVerbatim: {
		{Regex: regexp.MustCompile(`^[^"]+`), Scope: ScopeStringVerbatim},
		{Regex: regexp.MustCompile(`^""`), Scope: ScopeStringVerbatim},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringVerbatim, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringVerbatim},
	},
	GrammarStateStringRaw: {
		{Regex: regexp.MustCompile(`^[\s\S]*?"""`), Scope: ScopeStringRaw, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]+`), Scope: ScopeStringRaw},
	},
}
