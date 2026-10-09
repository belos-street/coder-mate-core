package sql

const (
	ScopeCommentLine               = "comment.line.double-dash.sql"
	ScopeCommentBlock              = "comment.block.sql"
	ScopeKeywordControl            = "keyword.control.sql"
	ScopeKeywordOperator           = "keyword.operator.sql"
	ScopeSupportFunction           = "support.function.builtin.sql"
	ScopeConstantBoolean           = "constant.language.boolean.sql"
	ScopeConstantNull              = "constant.language.null.sql"
	ScopeConstantNumeric           = "constant.numeric.sql"
	ScopeStringQuotedSingle        = "string.quoted.single.sql"
	ScopeStringQuotedDouble        = "string.quoted.double.sql"
	ScopeVariableParameter         = "variable.parameter.sql"
	ScopeEntityTable               = "entity.name.table.sql"
	ScopeEntityColumn              = "entity.name.column.sql"
	ScopePunctuationSeparatorComma = "punctuation.separator.comma.sql"
	ScopePunctuationTerminator     = "punctuation.terminator.statement.sql"
	ScopePunctuationGroupBegin     = "punctuation.section.group.begin.sql"
	ScopePunctuationGroupEnd       = "punctuation.section.group.end.sql"
	ScopeDefault                   = "default"
)

const (
	GrammarStateGlobal          = "global"
	GrammarStateExpectTableName = "expect-table-name"
	GrammarStateStringSingle    = "string-single"
	GrammarStateStringDouble    = "string-double"
	GrammarStateCommentBlock    = "comment-block"
)
