package core

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

func NewParserContext(initialState string) *ParserContext {
	return &ParserContext{
		StateStack: []string{initialState},
		Line:       1,
		Col:        1,
	}
}

func PushState(ctx *ParserContext, state string) *ParserContext {
	newCtx := *ctx
	newStack := make([]string, len(ctx.StateStack))
	copy(newStack, ctx.StateStack)
	newStack = append(newStack, state)
	newCtx.StateStack = newStack
	return &newCtx
}

func PopState(ctx *ParserContext) (*ParserContext, error) {
	if len(ctx.StateStack) <= 1 {
		return nil, errors.New("cannot pop the last state from stack")
	}
	newCtx := *ctx
	newStack := make([]string, len(ctx.StateStack)-1)
	copy(newStack, ctx.StateStack[:len(ctx.StateStack)-1])
	newCtx.StateStack = newStack
	return &newCtx, nil
}

func CurrentState(ctx *ParserContext) string {
	return ctx.StateStack[len(ctx.StateStack)-1]
}

func SplitTokenByLineBreak(text, tokenScope, fallbackScope string, startLine, startCol int) []Token {
	var tokens []Token
	remaining := text
	currentLine := startLine
	currentCol := startCol

	for len(remaining) > 0 {
		lineBreakIndex := strings.IndexByte(remaining, '\n')

		if lineBreakIndex == -1 {
			tokens = append(tokens, Token{
				Text:  remaining,
				Scope: tokenScope,
				Line:  currentLine,
				Col:   [2]int{currentCol, currentCol + len(remaining) - 1},
				Style: TokenStyle{},
			})
			break
		}

		isCRLF := lineBreakIndex > 0 && remaining[lineBreakIndex-1] == '\r'
		var lineBreakChar string
		if isCRLF {
			lineBreakChar = remaining[lineBreakIndex-1 : lineBreakIndex+1]
		} else {
			lineBreakChar = remaining[lineBreakIndex : lineBreakIndex+1]
		}

		var beforeLineBreak string
		if isCRLF {
			beforeLineBreak = remaining[:lineBreakIndex-1]
		} else {
			beforeLineBreak = remaining[:lineBreakIndex]
		}

		if beforeLineBreak != "" {
			beforeToken := Token{
				Text:  beforeLineBreak,
				Scope: tokenScope,
				Line:  currentLine,
				Col:   [2]int{currentCol, currentCol + len(beforeLineBreak) - 1},
				Style: TokenStyle{},
			}
			tokens = append(tokens, beforeToken)
			currentCol = beforeToken.Col[1] + 1
		}

		tokens = append(tokens, Token{
			Text:  lineBreakChar,
			Scope: fallbackScope,
			Line:  currentLine,
			Col:   [2]int{currentCol, currentCol + len(lineBreakChar) - 1},
			Style: TokenStyle{},
		})

		currentLine++
		currentCol = 1

		var sliceStart int
		if isCRLF {
			sliceStart = lineBreakIndex - 1
		} else {
			sliceStart = lineBreakIndex
		}
		remaining = remaining[sliceStart+len(lineBreakChar):]
	}

	return tokens
}

var lineBreakRegex = regexp.MustCompile(`^\r?\n`)

func MatchToken(code string, ctx *ParserContext, spec *TokenizerSpec) (*Token, *ParserContext) {
	if len(code) == 0 {
		return nil, ctx
	}

	currentState := CurrentState(ctx)
	rules, ok := spec.Rules[currentState]

	if ok {
		for _, rule := range rules {
			match := rule.Regex.FindStringIndex(code)
			if match == nil || match[0] != 0 {
				continue
			}

			matchedText := code[:match[1]]
			token := Token{
				Text:  matchedText,
				Scope: rule.Scope,
				Line:  ctx.Line,
				Col:   [2]int{ctx.Col, ctx.Col + len(matchedText) - 1},
				Style: TokenStyle{},
			}

			newCtx := *ctx
			if rule.PushState != "" {
				newCtx = *PushState(&newCtx, rule.PushState)
			}
			if rule.PopState {
				// 与 TS 版对齐：栈底 pop 静默忽略（PopState 返回错误时不回退）
				if popped, err := PopState(&newCtx); err == nil {
					newCtx = *popped
				}
			}
			newCtx.Col = token.Col[1] + 1

			return &token, &newCtx
		}
	}

	return matchFallbackToken(code, ctx, spec)
}

// matchFallbackToken 当前状态无规则命中时的兜底：按单个 rune 消费，
// 保证多字节字符（CJK、emoji 等）不会被按字节切碎。code 须非空。
// Text 取原始字节切片（code[:size]）：非法 UTF-8 字节按原字节保留而非
// 替换成 U+FFFD，保证分词完全无损。
func matchFallbackToken(code string, ctx *ParserContext, spec *TokenizerSpec) (*Token, *ParserContext) {
	_, size := utf8.DecodeRuneInString(code)
	token := Token{
		Text:  code[:size],
		Scope: spec.FallbackScope,
		Line:  ctx.Line,
		Col:   [2]int{ctx.Col, ctx.Col + size - 1},
		Style: TokenStyle{},
	}
	return &token, &ParserContext{
		StateStack: newStateStack(ctx.StateStack),
		Line:       ctx.Line,
		Col:        ctx.Col + size,
	}
}

func newStateStack(original []string) []string {
	s := make([]string, len(original))
	copy(s, original)
	return s
}

func Parse(code string, spec *TokenizerSpec) TokenStream {
	var rows TokenStream
	var currentRowTokens []Token
	remainingCode := code
	ctx := NewParserContext(spec.InitialState)

	for len(remainingCode) > 0 {
		lineBreakMatch := lineBreakRegex.FindStringIndex(remainingCode)
		if lineBreakMatch != nil && lineBreakMatch[0] == 0 {
			lineBreakText := remainingCode[:lineBreakMatch[1]]
			currentRowTokens = append(currentRowTokens, Token{
				Text:  lineBreakText,
				Scope: spec.FallbackScope,
				Line:  ctx.Line,
				Col:   [2]int{ctx.Col, ctx.Col + len(lineBreakText) - 1},
				Style: TokenStyle{},
			})

			rows = append(rows, currentRowTokens)
			currentRowTokens = []Token{}
			ctx.Line++
			ctx.Col = 1

			remainingCode = remainingCode[lineBreakMatch[1]:]
			continue
		}

		token, newCtx := MatchToken(remainingCode, ctx, spec)

		// MatchToken 仅在空输入返回 nil，循环条件已保证非空；
		// 零宽匹配（token.Text 为空）由下方防护兜底
		if len(token.Text) == 0 {
			// 防护：零宽匹配的自定义规则会造成零长度消费死循环，按 1 rune 兜底
			token, newCtx = matchFallbackToken(remainingCode, ctx, spec)
		}

		if strings.Contains(token.Text, "\n") {
			splitTokens := SplitTokenByLineBreak(
				token.Text,
				token.Scope,
				spec.FallbackScope,
				token.Line,
				token.Col[0],
			)

			for _, splitToken := range splitTokens {
				if strings.Contains(splitToken.Text, "\n") {
					currentRowTokens = append(currentRowTokens, splitToken)
					rows = append(rows, currentRowTokens)
					currentRowTokens = []Token{}
					ctx.StateStack = newStateStack(newCtx.StateStack)
					ctx.Line = splitToken.Line + 1
					ctx.Col = 1
				} else {
					currentRowTokens = append(currentRowTokens, splitToken)
					ctx.StateStack = newStateStack(newCtx.StateStack)
					ctx.Line = splitToken.Line
					ctx.Col = splitToken.Col[1] + 1
				}
			}
		} else {
			currentRowTokens = append(currentRowTokens, *token)
			ctx.StateStack = newStateStack(newCtx.StateStack)
			ctx.Line = newCtx.Line
			ctx.Col = newCtx.Col
		}

		remainingCode = remainingCode[len(token.Text):]
	}

	if len(currentRowTokens) > 0 {
		rows = append(rows, currentRowTokens)
	}

	return rows
}
