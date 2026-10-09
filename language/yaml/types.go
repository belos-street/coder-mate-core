package yaml

const (
	ScopeComment                 = "comment.line.number-sign.yaml"
	ScopeKeywordDocumentBegin    = "keyword.control.document.begin.yaml"
	ScopeKeywordDocumentEnd      = "keyword.control.document.end.yaml"
	ScopeSupportPropertyName     = "support.type.property-name.yaml"
	ScopePunctuationSeparatorKV  = "punctuation.separator.key-value.yaml"
	ScopePunctuationSequenceItem = "punctuation.definition.sequence.item.yaml"
	ScopePunctuationSectionFlow  = "punctuation.section.flow.yaml"
	ScopeStringQuotedDouble      = "string.quoted.double.yaml"
	ScopeStringQuotedSingle      = "string.quoted.single.yaml"
	ScopeStringUnquoted          = "string.unquoted.yaml"
	ScopeConstantNumeric         = "constant.numeric.yaml"
	ScopeConstantBoolean         = "constant.language.boolean.yaml"
	ScopeConstantNull            = "constant.language.null.yaml"
	ScopeEntityTagAnchor         = "entity.name.tag.anchor.yaml"
	ScopeVariableAlias           = "variable.other.alias.yaml"
	ScopeEntityTag               = "entity.name.tag.yaml"
	ScopeDefault                 = "default"
)

const (
	GrammarStateGlobal       = "global"
	GrammarStateStringDouble = "string-double"
	GrammarStateStringSingle = "string-single"
)
