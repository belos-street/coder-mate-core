package core

import "regexp"

type TokenStyle map[string]string

type Token struct {
	Text  string
	Scope string
	Line  int
	// Col 为 1 起始、双端闭合的列区间，按 UTF-8 字节偏移计（非 rune / UTF-16 计数），
	// 多字节字符按其字节数推进
	Col   [2]int
	Style TokenStyle
}

type GrammarRule struct {
	Regex     *regexp.Regexp
	Scope     string
	PushState string
	PopState  bool
	Skip      bool
}

type ParserContext struct {
	StateStack []string
	Line       int
	Col        int
}

type TokenizerSpec struct {
	InitialState  string
	Rules         map[string][]GrammarRule
	FallbackScope string
}

type TokenStream [][]Token
