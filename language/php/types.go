package php

const (
	ScopeCommentLineDS   = "comment.line.double-slash.php"
	ScopeCommentLineHash = "comment.line.number-sign.php"
	ScopeCommentBlock    = "comment.block.php"
	ScopeMetaTag         = "meta.tag.php"
	ScopeMetaAttribute   = "meta.attribute.php"
	ScopeKeywordControl  = "keyword.control.php"
	ScopeKeywordDecl     = "keyword.declaration.php"
	ScopeKeywordModifier = "keyword.modifier.php"
	ScopeSupportBuiltin  = "support.type.builtin.php"
	ScopeEntityType      = "entity.name.type.php"
	ScopeEntityFunction  = "entity.name.function.php"
	ScopeNamespace       = "entity.name.namespace.php"
	ScopeVarLang         = "variable.language.php"
	ScopeVarOther        = "variable.other.php"
	ScopeConstantBoolean = "constant.language.boolean.php"
	ScopeConstantNull    = "constant.language.null.php"
	ScopeConstantMagic   = "constant.language.php"
	ScopeConstantNumeric = "constant.numeric.php"
	ScopeStringDouble    = "string.quoted.double.php"
	ScopeStringSingle    = "string.quoted.single.php"
	ScopeStringBacktick  = "string.quoted.backtick.php"
	ScopeOperator        = "operator.php"
	ScopeDefault         = "default"
)

const (
	GrammarStateGlobal         = "global"
	GrammarStateCommentBlock   = "comment-block"
	GrammarStateStringDouble   = "string-double"
	GrammarStateStringSingle   = "string-single"
	GrammarStateStringBacktick = "string-backtick"
)
