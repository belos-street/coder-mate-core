package cpp

const (
	ScopeCommentLineDS   = "comment.line.double-slash.cpp"
	ScopeCommentBlock    = "comment.block.cpp"
	ScopePreprocessor    = "meta.preprocessor.cpp"
	ScopeDirective       = "keyword.control.directive.cpp"
	ScopeKeywordControl  = "keyword.control.cpp"
	ScopeKeywordDecl     = "keyword.declaration.cpp"
	ScopeKeywordModifier = "keyword.modifier.cpp"
	ScopeSupportBuiltin  = "support.type.builtin.cpp"
	ScopeEntityType      = "entity.name.type.cpp"
	ScopeEntityFunction  = "entity.name.function.cpp"
	ScopeNamespace       = "entity.name.namespace.cpp"
	ScopeEntityMacro     = "entity.name.macro.cpp"
	ScopeConstantBool    = "constant.language.boolean.cpp"
	ScopeConstantNull    = "constant.language.null.cpp"
	ScopeConstantNumeric = "constant.numeric.cpp"
	ScopeStringAngle     = "string.quoted.angle.cpp"
	ScopeStringDouble    = "string.quoted.double.cpp"
	ScopeStringSingle    = "string.quoted.single.cpp"
	ScopeStringRaw       = "string.quoted.raw.cpp"
	ScopeOperator        = "operator.cpp"
	ScopeDefault         = "default"
)

const (
	GrammarStateGlobal               = "global"
	GrammarStateCommentBlock         = "comment-block"
	GrammarStateStringDouble         = "string-double"
	GrammarStateStringSingle         = "string-single"
	GrammarStateStringRaw            = "string-raw"
	GrammarStatePreprocessorInclude  = "preprocessor-include"
	GrammarStatePreprocessorDefine   = "preprocessor-define"
	GrammarStatePreprocessorMacroRef = "preprocessor-macro-ref"
)
