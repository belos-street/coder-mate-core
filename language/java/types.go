package java

const (
	ScopeCommentLine        = "comment.line.double-slash.java"
	ScopeCommentBlock       = "comment.block.java"
	ScopeKeywordControl     = "keyword.control.java"
	ScopeKeywordDeclaration = "keyword.declaration.java"
	ScopeKeywordModifier    = "keyword.modifier.java"
	ScopeSupportBuiltin     = "support.type.builtin.java"
	ScopeEntityType         = "entity.name.type.java"
	ScopeEntityFunction     = "entity.name.function.java"
	ScopeNamespace          = "entity.name.namespace.java"
	ScopeAnnotation         = "meta.annotation.java"
	ScopeVariableLang       = "variable.language.java"
	ScopeConstantBoolean    = "constant.language.boolean.java"
	ScopeConstantNull       = "constant.language.null.java"
	ScopeConstantNumeric    = "constant.numeric.java"
	ScopeStringDouble       = "string.quoted.double.java"
	ScopeStringSingle       = "string.quoted.single.java"
	ScopeStringTriple       = "string.quoted.triple.java"
	ScopeOperator           = "operator.java"
	ScopeDefault            = "default"
)

const (
	GrammarStateGlobal       = "global"
	GrammarStateCommentBlock = "comment-block"
	GrammarStateStringDouble = "string-double"
	GrammarStateStringSingle = "string-single"
	GrammarStateStringTriple = "string-text-block"
)
