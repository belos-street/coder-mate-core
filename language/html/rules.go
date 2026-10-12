package html

import (
	"regexp"

	"code-mate-core/core"
)

var htmlEntity = regexp.MustCompile(`^&(?:[a-zA-Z][a-zA-Z0-9]+|#[0-9]+|#x[0-9A-Fa-f]+);`)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^<!--`), Scope: ScopeCommentBlock, PushState: GrammarStateComment},
		{Regex: regexp.MustCompile(`(?i)^<!DOCTYPE\b[^>]*>`), Scope: ScopeDoctype},
		{Regex: regexp.MustCompile(`^</`), Scope: ScopeTagBegin, PushState: GrammarStateCloseTag},
		{Regex: regexp.MustCompile(`^<[A-Za-z][A-Za-z0-9:-]*`), Scope: ScopeEntityTag, PushState: GrammarStateOpenTag},
		{Regex: regexp.MustCompile(`^<`), Scope: ScopeTagBegin, PushState: GrammarStateOpenTag},
		{Regex: htmlEntity, Scope: ScopeEntityChar},
		{Regex: regexp.MustCompile(`^[^<&\r\n]+`), Scope: ScopeTextPlain},
	},
	GrammarStateComment: {
		{Regex: regexp.MustCompile(`^[\s\S]*?-->`), Scope: ScopeCommentBlock, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]+`), Scope: ScopeCommentBlock},
	},
	GrammarStateOpenTag: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^/>`), Scope: ScopeTagEnd, PopState: true},
		{Regex: regexp.MustCompile(`^>`), Scope: ScopeTagEnd, PopState: true},
		{Regex: regexp.MustCompile(`^[A-Za-z_:][A-Za-z0-9:._-]*`), Scope: ScopeAttrName},
		{Regex: regexp.MustCompile(`^=`), Scope: ScopeSeparatorKV, PushState: GrammarStateAttributeValue},
		{Regex: htmlEntity, Scope: ScopeEntityChar},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault},
	},
	GrammarStateCloseTag: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^[A-Za-z][A-Za-z0-9:-]*`), Scope: ScopeEntityTag},
		{Regex: regexp.MustCompile(`^>`), Scope: ScopeTagEnd, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
	},
	GrammarStateAttributeValue: {
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^"[^"\r\n]*"?`), Scope: ScopeStringQuotedDouble, PopState: true},
		{Regex: regexp.MustCompile(`^'[^'\r\n]*'?`), Scope: ScopeStringQuotedSingle, PopState: true},
		{Regex: regexp.MustCompile(`^(?:[^\s"'=<>` + "\x60" + `]+)`), Scope: ScopeStringUnquoted, PopState: true},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeDefault, PopState: true},
	},
}
