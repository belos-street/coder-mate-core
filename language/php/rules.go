package php

import (
	"regexp"

	"code-mate-core/core"
)

var phpKeywordsDecl = regexp.MustCompile(`^(namespace|use|class|interface|trait|enum|function|const)\b`)
var phpKeywordsModifier = regexp.MustCompile(`^(public|private|protected|static|final|abstract|readonly)\b`)
var phpKeywordsControl = regexp.MustCompile(`^(if|else|elseif|switch|case|default|for|foreach|while|do|break|continue|return|try|catch|finally|throw|match|yield|as|new|instanceof|echo|include|require|include_once|require_once|fn)\b`)
var phpBuiltinTypes = regexp.MustCompile(`^(int|float|string|bool|array|object|mixed|void|null|callable|iterable|self|parent|never)\b`)

var GrammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^//.*`), Scope: ScopeCommentLineDS},
		{Regex: regexp.MustCompile(`^#\[[^\]\r\n]+\]`), Scope: ScopeMetaAttribute},
		{Regex: regexp.MustCompile(`^#.*`), Scope: ScopeCommentLineHash},
		{Regex: regexp.MustCompile(`^/\*`), Scope: ScopeCommentBlock, PushState: GrammarStateCommentBlock},
		{Regex: regexp.MustCompile(`^<\?(?:php|=)?|^\?>`), Scope: ScopeMetaTag},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PushState: GrammarStateStringSingle},
		{Regex: regexp.MustCompile("^`"), Scope: ScopeStringBacktick, PushState: GrammarStateStringBacktick},
		{Regex: regexp.MustCompile(`^[A-Za-z_\x80-\xff][A-Za-z0-9_\x80-\xff]*(?:\\[A-Za-z_\x80-\xff][A-Za-z0-9_\x80-\xff]*)+`), Scope: ScopeNamespace},
		{Regex: phpKeywordsDecl, Scope: ScopeKeywordDecl},
		{Regex: phpKeywordsModifier, Scope: ScopeKeywordModifier},
		{Regex: phpKeywordsControl, Scope: ScopeKeywordControl},
		{Regex: regexp.MustCompile(`(?i)^(true|false)\b`), Scope: ScopeConstantBoolean},
		{Regex: regexp.MustCompile(`(?i)^null\b`), Scope: ScopeConstantNull},
		{Regex: regexp.MustCompile(`^__(?:FILE|DIR|LINE|FUNCTION|CLASS|METHOD|NAMESPACE|TRAIT)__\b`), Scope: ScopeConstantMagic},
		{Regex: phpBuiltinTypes, Scope: ScopeSupportBuiltin},
		{Regex: regexp.MustCompile(`^\$this\b`), Scope: ScopeVarLang},
		{Regex: regexp.MustCompile(`^\$[A-Za-z_\x80-\xff][A-Za-z0-9_\x80-\xff]*`), Scope: ScopeVarOther},
		{Regex: regexp.MustCompile(`^[A-Za-z_\x80-\xff][A-Za-z0-9_\x80-\xff]*\(`), Scope: ScopeEntityFunction},
		{Regex: regexp.MustCompile(`^[A-Z][A-Za-z0-9_\x80-\xff]*`), Scope: ScopeEntityType},
		{Regex: regexp.MustCompile(`^-?(?:0[xX][0-9A-Fa-f_]+|0[bB][01_]+|0[oO][0-7_]+|(?:\d[\d_]*)(?:\.\d[\d_]*)?(?:[eE][+-]?\d[\d_]*)?)`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(?:\?\?=|\?\?|=>|->|::|===|!==|==|!=|<=|>=|<=>|&&|\|\||\+=|-=|\*=|\/=|%=|\.=|&=|\|=|\^=|<<|>>|[-+*/%=&|^~!?<>:.;,()[\]{}])`), Scope: ScopeOperator},
		{Regex: regexp.MustCompile(`^[A-Za-z_\x80-\xff][A-Za-z0-9_\x80-\xff]*`), Scope: ScopeDefault},
	},
	GrammarStateCommentBlock: {
		{Regex: regexp.MustCompile(`^[\s\S]*?\*/`), Scope: ScopeCommentBlock, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]+`), Scope: ScopeCommentBlock},
	},
	GrammarStateStringDouble: {
		{Regex: regexp.MustCompile(`^[^"\\$\n]+`), Scope: ScopeStringDouble},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringDouble},
		{Regex: regexp.MustCompile(`^\$[A-Za-z_\x80-\xff][A-Za-z0-9_\x80-\xff]*`), Scope: ScopeVarOther},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringDouble, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringDouble},
	},
	GrammarStateStringSingle: {
		{Regex: regexp.MustCompile(`^[^'\\\n]+`), Scope: ScopeStringSingle},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringSingle},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringSingle},
	},
	GrammarStateStringBacktick: {
		{Regex: regexp.MustCompile(`^[^` + "`" + `\\\n]+`), Scope: ScopeStringBacktick},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringBacktick},
		{Regex: regexp.MustCompile("^`"), Scope: ScopeStringBacktick, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringBacktick},
	},
}
