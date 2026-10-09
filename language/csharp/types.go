package csharp

const (
	ScopeCommentLineDS   = "comment.line.double-slash.csharp"
	ScopeCommentBlock    = "comment.block.csharp"
	ScopePreprocessor    = "meta.preprocessor.csharp"
	ScopeMetaAttribute   = "meta.attribute.csharp"
	ScopeKeywordControl  = "keyword.control.csharp"
	ScopeKeywordDecl     = "keyword.declaration.csharp"
	ScopeKeywordModifier = "keyword.modifier.csharp"
	ScopeSupportBuiltin  = "support.type.builtin.csharp"
	ScopeEntityType      = "entity.name.type.csharp"
	ScopeEntityFunction  = "entity.name.function.csharp"
	ScopeNamespace       = "entity.name.namespace.csharp"
	ScopeVarLang         = "variable.language.csharp"
	ScopeConstantBool    = "constant.language.boolean.csharp"
	ScopeConstantNull    = "constant.language.null.csharp"
	ScopeConstantNumeric = "constant.numeric.csharp"
	ScopeStringDouble    = "string.quoted.double.csharp"
	ScopeStringSingle    = "string.quoted.single.csharp"
	ScopeStringVerbatim  = "string.quoted.verbatim.csharp"
	ScopeStringRaw       = "string.quoted.raw.csharp"
	ScopeOperator        = "operator.csharp"
	ScopeDefault         = "default"
)

const (
	GrammarStateGlobal         = "global"
	GrammarStateCommentBlock   = "comment-block"
	GrammarStateStringDouble   = "string-double"
	GrammarStateStringSingle   = "string-single"
	GrammarStateStringVerbatim = "string-verbatim"
	GrammarStateStringRaw      = "string-raw"
)
