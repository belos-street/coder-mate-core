package python

import (
	"regexp"

	"code-mate-core/core"
)

var pythonBuiltins = regexp.MustCompile(`^(print|len|range|str|int|float|dict|list|set|tuple|sum|min|max|abs|enumerate|zip|map|filter|any|all|open)\b`)
var pythonKeywords = regexp.MustCompile(`^(if|elif|else|while|try|except|finally|with|pass|break|continue|return|raise|in|is|and|or|not|lambda|yield|await|async|from|import|match|case|assert|del|global|nonlocal)\b`)
var pythonAnnotationTypes = regexp.MustCompile(`^(int|str|float|bool|list|dict|set|tuple|bytes|Any|Optional|Union|Literal|TypedDict|Callable|Iterator|Iterable|Sequence|Mapping|Self|(?:[A-Z][A-Za-z0-9_]*))\b`)

var GrammarRules = map[string][]core.GrammarRule{
	StateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^#.*`), Scope: ScopeCommentLine},

		{Regex: regexp.MustCompile(`^@[A-Za-z_][A-Za-z0-9_.]*`), Scope: ScopeDecorator},

		// f-string triple
		{Regex: regexp.MustCompile(`^(?:[fF][rR]?|[rR][fF])"""`), Scope: ScopeStringInterp, PushState: StateFStringTripleD},
		{Regex: regexp.MustCompile(`^(?:[fF][rR]?|[rR][fF])'''`), Scope: ScopeStringInterp, PushState: StateFStringTripleS},

		// triple strings
		{Regex: regexp.MustCompile(`^(?:[rRuUbB])?"""`), Scope: ScopeStringTripleD, PushState: StateStringTripleD},
		{Regex: regexp.MustCompile(`^(?:[rRuUbB])?'''`), Scope: ScopeStringTripleS, PushState: StateStringTripleS},

		// f-string single/double
		{Regex: regexp.MustCompile(`^(?:[fF][rR]?|[rR][fF])"`), Scope: ScopeStringInterp, PushState: StateFStringDouble},
		{Regex: regexp.MustCompile(`^(?:[fF][rR]?|[rR][fF])'`), Scope: ScopeStringInterp, PushState: StateFStringSingle},

		// single/double strings
		{Regex: regexp.MustCompile(`^(?:[rRuUbB])?"`), Scope: ScopeStringDouble, PushState: StateStringDouble},
		{Regex: regexp.MustCompile(`^(?:[rRuUbB])?'`), Scope: ScopeStringSingle, PushState: StateStringSingle},

		// type annotations (-> return type, : parameter type)
		{Regex: regexp.MustCompile(`^->\s*(?:int|str|float|bool|list|dict|set|tuple|bytes|Any|Optional|Union|Literal|TypedDict|Callable|Iterator|Iterable|Sequence|Mapping|Self)`), Scope: ScopeSupportAnnotation},
		{Regex: regexp.MustCompile(`^:\s*(?:int|str|float|bool|list|dict|set|tuple|bytes|Any|Optional|Union|Literal|TypedDict|Callable|Iterator|Iterable|Sequence|Mapping|Self)`), Scope: ScopeSupportAnnotation},

		{Regex: regexp.MustCompile(`^(def)\b`), Scope: ScopeKeywordDecl, PushState: StateExpectFuncName},
		{Regex: regexp.MustCompile(`^(class)\b`), Scope: ScopeKeywordDecl, PushState: StateExpectClassName},
		{Regex: regexp.MustCompile(`^(for)\b`), Scope: ScopeKeywordControl, PushState: StateExpectComprehension},
		{Regex: regexp.MustCompile(`^(as)\b`), Scope: ScopeKeywordControl, PushState: StateExpectAsAlias},
		{Regex: pythonKeywords, Scope: ScopeKeywordControl},
		{Regex: pythonBuiltins, Scope: ScopeSupportBuiltin},
		{Regex: regexp.MustCompile(`^(True|False)\b`), Scope: ScopeConstantBoolean},
		{Regex: regexp.MustCompile(`^None\b`), Scope: ScopeConstantNone},
		{Regex: pythonAnnotationTypes, Scope: ScopeSupportAnnotation},
		{Regex: regexp.MustCompile(`^-?(?:0[xX][0-9a-fA-F_]+|0[oO][0-7_]+|0[bB][01_]+|\d[\d_]*(?:\.\d[\d_]*)?(?:[eE][+-]?\d+)?)[jJ]?`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(==|!=|<=|>=|\*\*|\/\/|->|:=|[+\-*/%=<>{}\[\]():.,|])`), Scope: ScopeOperator},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeVarIdentifier},
	},

	StateExpectFuncName: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeEntityFunction, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
	},
	StateExpectClassName: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeEntityClass, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
	},
	StateExpectComprehension: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^(in)\b`), Scope: ScopeKeywordControl, PopState: true},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeVarComprehension},
		{Regex: regexp.MustCompile(`^(,|\(|\)|\[|\])`), Scope: ScopeOperator},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
	},
	StateExpectAsAlias: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeVarAlias, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
	},

	StateStringSingle: {
		{Regex: regexp.MustCompile(`^[^\\']*(?:\\.[^\\']*)*'`), Scope: ScopeStringSingle, PopState: true},
		{Regex: regexp.MustCompile(`^.*`), Scope: ScopeStringSingle},
	},
	StateStringDouble: {
		{Regex: regexp.MustCompile(`^[^\\"]*(?:\\.[^\\"]*)*"`), Scope: ScopeStringDouble, PopState: true},
		{Regex: regexp.MustCompile(`^.*`), Scope: ScopeStringDouble},
	},
	StateStringTripleS: {
		{Regex: regexp.MustCompile(`^[\s\S]*?'''`), Scope: ScopeStringTripleS, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]*`), Scope: ScopeStringTripleS},
	},
	StateStringTripleD: {
		{Regex: regexp.MustCompile(`^[\s\S]*?"""`), Scope: ScopeStringTripleD, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]*`), Scope: ScopeStringTripleD},
	},

	StateFStringSingle: {
		{Regex: regexp.MustCompile(`^{{|^}}`), Scope: ScopeStringInterp},
		{Regex: regexp.MustCompile(`^\{`), Scope: ScopeInterpBegin, PushState: StateFStringInterpolation},
		{Regex: regexp.MustCompile(`^'`), Scope: ScopeStringInterp, PopState: true},
		{Regex: regexp.MustCompile(`^[^\\{'\n]+`), Scope: ScopeStringInterp},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringInterp},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringInterp},
	},
	StateFStringDouble: {
		{Regex: regexp.MustCompile(`^{{|^}}`), Scope: ScopeStringInterp},
		{Regex: regexp.MustCompile(`^\{`), Scope: ScopeInterpBegin, PushState: StateFStringInterpolation},
		{Regex: regexp.MustCompile(`^"`), Scope: ScopeStringInterp, PopState: true},
		{Regex: regexp.MustCompile(`^[^\\{"\n]+`), Scope: ScopeStringInterp},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringInterp},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringInterp},
	},
	StateFStringTripleS: {
		{Regex: regexp.MustCompile(`^{{|^}}`), Scope: ScopeStringInterp},
		{Regex: regexp.MustCompile(`^\{`), Scope: ScopeInterpBegin, PushState: StateFStringInterpolation},
		{Regex: regexp.MustCompile(`^'''`), Scope: ScopeStringInterp, PopState: true},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringInterp},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringInterp},
	},
	StateFStringTripleD: {
		{Regex: regexp.MustCompile(`^{{|^}}`), Scope: ScopeStringInterp},
		{Regex: regexp.MustCompile(`^\{`), Scope: ScopeInterpBegin, PushState: StateFStringInterpolation},
		{Regex: regexp.MustCompile(`^"""`), Scope: ScopeStringInterp, PopState: true},
		{Regex: regexp.MustCompile(`^\\.`), Scope: ScopeStringInterp},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringInterp},
	},

	StateFStringInterpolation: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^\{`), Scope: ScopeInterpBegin, PushState: StateFStringInterpolation},
		{Regex: regexp.MustCompile(`^}`), Scope: ScopeInterpEnd, PopState: true},
		{Regex: regexp.MustCompile(`^![sra]`), Scope: ScopeInterpFormat},
		{Regex: regexp.MustCompile(`^:`), Scope: ScopeInterpFormat},
		{Regex: regexp.MustCompile(`^(def|class)\b`), Scope: ScopeKeywordDecl},
		{Regex: regexp.MustCompile(`^(for)\b`), Scope: ScopeKeywordControl, PushState: StateExpectComprehension},
		{Regex: pythonKeywords, Scope: ScopeKeywordControl},
		{Regex: regexp.MustCompile(`^(True|False|None)\b`), Scope: ScopeConstantBoolean},
		{Regex: regexp.MustCompile(`^-?(?:0[xX][0-9a-fA-F_]+|0[oO][0-7_]+|0[bB][01_]+|\d[\d_]*(?:\.\d[\d_]*)?(?:[eE][+-]?\d+)?)[jJ]?`), Scope: ScopeConstantNumeric},
		{Regex: regexp.MustCompile(`^(==|!=|<=|>=|[+\-*/%=<>{}\[\]():.,|])`), Scope: ScopeOperator},
		{Regex: regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*`), Scope: ScopeVarIdentifier},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault},
	},
}
