package parser

import (
	"fmt"
	"os"
	"regex-engine/internals/token"
)

type ParseContext struct {
	tokens []token.Token

	pos int
}

func (p *ParseContext) GetTokens() []token.Token {
	return p.tokens
}

func Parse(pattern string) *ParseContext {
	context := &ParseContext{
		pos:    0,
		tokens: []token.Token{},
	}

	for context.pos < len(pattern) {
		parsePattern(pattern, context)
		context.pos++
	}

	return context
}

func parsePattern(pattern string, context *ParseContext) {
	curChar := rune(pattern[context.pos])

	switch curChar {
	case '(': // (abc)
		groupContext := &ParseContext{
			pos:    0,
			tokens: []token.Token{},
		}

		parseGroup(pattern, groupContext)

		context.tokens = append(context.tokens, groupContext.tokens...)
        context.pos = groupContext.pos

	case '[': // [abc]
	case '|': //[a-e]
	case '{': // {3, } {3, 4} {,10}
	case '*', '?', '+': // a*, a?, a+
	default:
		// literal
		context.tokens = append(context.tokens, token.Token{
			Type:  token.LITERAL,
			Value: curChar,
		})
	}
}

func parseGroup(pattern string, groupContext *ParseContext) {
	groupContext.pos++ //skip (

	for groupContext.pos < len(pattern) && pattern[groupContext.pos] != ')' {
		parsePattern(pattern, groupContext)
		groupContext.pos++
	}

	if pattern[groupContext.pos] != ')' {
		fmt.Fprintf(os.Stderr, "[ERROR] Unclosed group bracket: %s", pattern)
		os.Exit(1)
	}

	groupContext.pos++
}
