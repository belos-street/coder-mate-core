package bash

const (
	ScopeCommentLine     = "comment.line.number-sign.bash"
	ScopeKeywordControl  = "keyword.control.bash"
	ScopeKeywordDeclFunc = "keyword.declaration.function.bash"
	ScopeEntityFunction  = "entity.name.function.bash"
	ScopeEntityCommand   = "entity.name.command.bash"
	ScopeSupportBuiltin  = "support.function.builtin.bash"
	ScopeVarParameter    = "variable.parameter.bash"
	ScopeVarSpecial      = "variable.language.special.bash"
	ScopeVarReadWrite    = "variable.other.readwrite.bash"
	ScopeConstantNumeric = "constant.numeric.bash"
	ScopeStringDouble    = "string.quoted.double.bash"
	ScopeStringSingle    = "string.quoted.single.bash"
	ScopeOperator        = "operator.bash"
	ScopeDefault         = "default"
)

const (
	GrammarStateGlobal         = "global"
	GrammarStateExpectFuncName = "expect-function-name"
	GrammarStateStringDouble   = "string-double"
	GrammarStateStringSingle   = "string-single"
)
