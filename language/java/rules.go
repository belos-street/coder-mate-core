package java

import (
	"regexp"

	"code-mate-core/core"
)

var javaKeywordsDecl = regexp.MustCompile(`^(package|import|class|interface|enum|record|extends|implements)\b`)
var javaKeywordsModifier = regexp.MustCompile(`^(public|protected|private|static|final|abstract|sealed|non-sealed|volatile|transient|native|strictfp|synchronized|default)\b`)
var javaKeywordsControl = regexp.MustCompile(`^(if|else|switch|case|for|while|do|break|continue|return|try|catch|finally|throw|throws|new|instanceof|assert|yield)\b`)
var javaBuiltinTypes = regexp.MustCompile(`^(void|boolean|byte|short|int|long|float|double|char|var|String)\b`)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^//.*`), Scope: ScopeCommentLine},
		{Regex: regexp.MustCompile(`^/\*`), Scope: ScopeCommentBlock, PushState: GrammarStateCommentBlock},
		{Regex: regexp.MustCompile(`^"""`), Scope: ScopeStringTriple, PushState: GrammarStateStringTriple},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringDouble, PushState: GrammarStateStringDouble},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringSingle, PushState: GrammarStateStringSingle},
		{Regex: regexp.MustCompile(`^@[A-Za-z_][A-Za-z0-9_.]*`), Scope: ScopeAnnotation},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)+`), Scope: ScopeNamespace},
		{Regex: javaKeywordsDecl, Scope: ScopeKeywordDeclaration},
		{Regex: javaKeywordsModifier, Scope: ScopeKeywordModifier},
		{Regex: javaKeywordsControl, Scope: ScopeKeywordControl},
		{Regex: regexp.MustCompile(`^(true|false)\b`), Scope: ScopeConstantBoolean},
		{Regex: regexp.MustCompile(`^null\b`), Scope: ScopeConstantNull},
		{Regex: regexp.MustCompile(`^(this|super)\b`), Scope: ScopeVariableLang},
		{Regex: javaBuiltinTypes, Scope: ScopeSupportBuiltin},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*[ \t]*\(`), Scope: ScopeEntityFunction},
		{Regex: regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*(?:<[^>]*>)?`), Scope: ScopeEntityType},
		{Regex: regexp.MustCompile(`^-?(?:0[xX][0-9A-Fa-f_]+|0[bB][01_]+|0[0-7_]+|(?:\d[\d_]*)(?:\.\d[\d_]*)?(?:[eE][+-]?\d[\d_]*)?)[fFdDlL]?`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(?:->|::|==|!=|<=|>=|\+\+|--|&&|\|\||<<|>>>?|[-+*/%=&|^~!?<>:.;,()[\]{}])`), Scope: ScopeOperator},
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
	GrammarStateStringTriple: {
		{Regex: regexp.MustCompile(`^[\s\S]*?"""`), Scope: ScopeStringTriple, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]+`), Scope: ScopeStringTriple},
	},
}
