package markdown

import (
	"regexp"

	"code-mate-core/core"
)

var grammarRules = map[string][]core.GrammarRule{
	GrammarStateGlobal: {
		{Regex: regexp.MustCompile(`^(?: {4}|\t)[^\r\n]*`), Scope: ScopeIndentedCode},
		{Regex: regexp.MustCompile(`^\s+`), Scope: ScopeDefault},
		{Regex: regexp.MustCompile(`^<!--`), Scope: ScopeComment, PushState: GrammarStateComment},
		{Regex: regexp.MustCompile(`^(?: {0,3})(?:` + "\x60{3,}" + `|~~~)[^\r\n]*`), Scope: ScopeFencedBegin, PushState: GrammarStateFencedCode},
		{Regex: regexp.MustCompile(`^(?: {0,3})#{1,6}\s+[^\r\n]*`), Scope: ScopeHeading},
		{Regex: regexp.MustCompile(`^(?: {0,3})(?:(?:\*[ \t]*){3,}|(?:-[ \t]*){3,}|(?:_[ \t]*){3,})`), Scope: ScopeHr},
		{Regex: regexp.MustCompile(`^(?: {0,3})(?:=+|-{2,})[ \t]*(?:\r?\n)?`), Scope: ScopeSetextHeading},
		{Regex: regexp.MustCompile(`^(?: {0,3})>\s?`), Scope: ScopeQuote},
		{Regex: regexp.MustCompile(`^(?: {0,3})(?:[-+*]|\d+[.)])\s+\[(?: |x|X)\]\s+`), Scope: ScopeTask},
		{Regex: regexp.MustCompile(`^(?: {0,3})(?:[-+*]|\d+[.)])\s+`), Scope: ScopeList},
		{Regex: regexp.MustCompile(`^(?: {0,3})\|?(?:[ \t]*:?-{3,}:?[ \t]*\|)+[ \t]*:?-{3,}:?[ \t]*\|?[ \t]*`), Scope: ScopeTableSep},
		{Regex: regexp.MustCompile(`^(?: {0,3})\|[^\r\n]*\|`), Scope: ScopeTableRow},
		{Regex: regexp.MustCompile(`^\[[^\]\r\n]+\]:`), Scope: ScopeRefLinkLabel},
		{Regex: regexp.MustCompile(`^(?:[ \t]*<https?:\/\/[^>\s]+>|[ \t]*https?:\/\/[^\s]+)`), Scope: ScopeRefLinkUrl},
		{Regex: regexp.MustCompile(`^!\[[^\]\r\n]*\]`), Scope: ScopeImageLabel},
		{Regex: regexp.MustCompile(`^\[[^\]\r\n]*\]`), Scope: ScopeLinkLabel},
		{Regex: regexp.MustCompile(`^\((?:[^()\r\n]|\\.)*\)`), Scope: ScopeLinkUrl},
		{Regex: regexp.MustCompile(`^[\[\]\(\)!]`), Scope: ScopeLinkPunct},
		{Regex: regexp.MustCompile("^\x60"), Scope: ScopeInlineRaw, PushState: GrammarStateInlineCode},
		{Regex: regexp.MustCompile(`^~~`), Scope: ScopeStrikethrough, PushState: GrammarStateStrikethrough},
		{Regex: regexp.MustCompile(`^(?:\*\*|__)`), Scope: ScopeBold, PushState: GrammarStateStrong},
		{Regex: regexp.MustCompile(`^(?:\*|_)`), Scope: ScopeItalic, PushState: GrammarStateEmphasis},
		{Regex: regexp.MustCompile(`^(?:<https?:\/\/[^>\s]+>|https?:\/\/[^\s)\]]+)`), Scope: ScopeLinkUrl},
		{Regex: regexp.MustCompile(`^[^\s` + "\x60" + `*_~#[\]()!<>|]+`), Scope: ScopeStringUnquoted},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringUnquoted},
	},
	GrammarStateComment: {
		{Regex: regexp.MustCompile(`^[\s\S]*?-->`), Scope: ScopeComment, PopState: true},
		{Regex: regexp.MustCompile(`^[\s\S]+`), Scope: ScopeComment},
	},
	GrammarStateFencedCode: {
		{Regex: regexp.MustCompile(`^(?: {0,3})(?:` + "\x60{3,}" + `|~~~)[ \t]*`), Scope: ScopeFencedEnd, PopState: true},
		{Regex: regexp.MustCompile(`^[^\r\n]+`), Scope: ScopeStringUnquoted},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStringUnquoted},
	},
	GrammarStateInlineCode: {
		{Regex: regexp.MustCompile("^\x60"), Scope: ScopeInlineRaw, PopState: true},
		{Regex: regexp.MustCompile(`^[^\x60\r\n]+`), Scope: ScopeInlineRaw},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeInlineRaw},
	},
	GrammarStateStrong: {
		{Regex: regexp.MustCompile(`^(?:\*\*|__)`), Scope: ScopeBold, PopState: true},
		{Regex: regexp.MustCompile(`^[^\r\n*]+?[ \t]*`), Scope: ScopeBold},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeBold},
	},
	GrammarStateEmphasis: {
		{Regex: regexp.MustCompile(`^(?:\*|_)`), Scope: ScopeItalic, PopState: true},
		{Regex: regexp.MustCompile(`^[^\r\n*_]+?[ \t]*`), Scope: ScopeItalic},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeItalic},
	},
	GrammarStateStrikethrough: {
		{Regex: regexp.MustCompile(`^~~`), Scope: ScopeStrikethrough, PopState: true},
		{Regex: regexp.MustCompile(`^[^\r\n~]+`), Scope: ScopeStrikethrough},
		{Regex: regexp.MustCompile(`^.`), Scope: ScopeStrikethrough},
	},
}
