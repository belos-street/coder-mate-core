package gorule

const (
	ScopeCommentLine        = "comment.line.double-slash.go"
	ScopeCommentBlock       = "comment.block.go"
	ScopeKeywordControl     = "keyword.control.go"
	ScopeKeywordDeclaration = "keyword.declaration.go"
	ScopeSupportBuiltin     = "support.type.builtin.go"
	ScopeEntityFunction     = "entity.name.function.go"
	ScopeEntityType         = "entity.name.type.go"
	ScopeNamespace          = "entity.name.namespace.go"
	ScopeConstantBoolean    = "constant.language.boolean.go"
	ScopeConstantNull       = "constant.language.null.go"
	ScopeConstantIota       = "constant.language.iota.go"
	ScopeConstantNumeric    = "constant.numeric.go"
	ScopeStringDouble       = "string.quoted.double.go"
	ScopeStringSingle       = "string.quoted.single.go"
	ScopeStringRaw          = "string.quoted.raw.go"
	ScopeOperator           = "operator.go"
	ScopeDefault            = "default"
)

const (
	GrammarStateGlobal       = "global"
	GrammarStateCommentBlock = "comment-block"
	GrammarStateStringDouble = "string-double"
	GrammarStateStringSingle = "string-single"
	GrammarStateStringRaw    = "string-raw"
)
