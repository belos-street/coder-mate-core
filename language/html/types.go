package html

const (
	ScopeCommentBlock       = "comment.block.html"
	ScopeDoctype            = "keyword.control.doctype.html"
	ScopeTagBegin           = "punctuation.definition.tag.begin.html"
	ScopeTagEnd             = "punctuation.definition.tag.end.html"
	ScopeEntityTag          = "entity.name.tag.html"
	ScopeAttrName           = "entity.other.attribute-name.html"
	ScopeSeparatorKV        = "punctuation.separator.key-value.html"
	ScopeStringQuotedDouble = "string.quoted.double.html"
	ScopeStringQuotedSingle = "string.quoted.single.html"
	ScopeStringUnquoted     = "string.unquoted.html"
	ScopeEntityChar         = "constant.character.entity.html"
	ScopeTextPlain          = "text.plain.html"
	ScopeDefault            = "default"
)

const (
	GrammarStateGlobal         = "global"
	GrammarStateComment        = "comment"
	GrammarStateOpenTag        = "open-tag"
	GrammarStateCloseTag       = "close-tag"
	GrammarStateAttributeValue = "attribute-value"
)
