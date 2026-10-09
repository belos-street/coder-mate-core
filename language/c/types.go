package c

func ScopeExists(st string) bool {
	switch st {
	case ScopeCommentLine, ScopeCommentBlock, ScopePreprocessor, ScopeDirective,
		ScopeKeywordControl, ScopeKeywordDeclaration, ScopeKeywordModifier,
		ScopeSupportBuiltin, ScopeEntityType, ScopeEntityFunction, ScopeEntityMacro,
		ScopeConstantBoolean, ScopeConstantNull, ScopeConstantNumeric,
		ScopeStringAngle, ScopeStringDouble, ScopeStringSingle, ScopeOperator, ScopeDefault:
		return true
	}
	return false
}

const (
	ScopeCommentLine        = "comment.line.double-slash.c"
	ScopeCommentBlock       = "comment.block.c"
	ScopePreprocessor       = "meta.preprocessor.c"
	ScopeDirective          = "keyword.control.directive.c"
	ScopeKeywordControl     = "keyword.control.c"
	ScopeKeywordDeclaration = "keyword.declaration.c"
	ScopeKeywordModifier    = "keyword.modifier.c"
	ScopeSupportBuiltin     = "support.type.builtin.c"
	ScopeEntityType         = "entity.name.type.c"
	ScopeEntityFunction     = "entity.name.function.c"
	ScopeEntityMacro        = "entity.name.macro.c"
	ScopeConstantBoolean    = "constant.language.boolean.c"
	ScopeConstantNull       = "constant.language.null.c"
	ScopeConstantNumeric    = "constant.numeric.c"
	ScopeStringAngle        = "string.quoted.angle.c"
	ScopeStringDouble       = "string.quoted.double.c"
	ScopeStringSingle       = "string.quoted.single.c"
	ScopeOperator           = "operator.c"
	ScopeDefault            = "default"
)

const (
	GrammarStateGlobal       = "global"
	GrammarStateCommentBlock = "comment-block"
	GrammarStateStringDouble = "string-double"
	GrammarStateStringSingle = "string-single"
)
