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
	curChar := pattern[context.pos]

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
		parseBracket(pattern, context)
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
	groupContext.pos++ // skip (

	for groupContext.pos < len(pattern) && pattern[groupContext.pos] != ')' {
		parsePattern(pattern, groupContext)
		groupContext.pos++
	}

	if groupContext.pos >= len(pattern) || pattern[groupContext.pos] != ')' {
		fmt.Fprintf(os.Stderr, "[ERROR] Unclosed group bracket: %s", pattern)
		os.Exit(1)
	}

	groupContext.pos++
}

func parseBracket(pattern string, context *ParseContext) {
	context.pos++ // Skip [

	literals := []string{}

	for context.pos < len(pattern) && pattern[context.pos] != ']' {
		if pattern[context.pos] == '-' {
			literals = append(literals, fmt.Sprintf("%c%c", pattern[context.pos-1], pattern[context.pos+1]))
			context.pos++
		} else {
			literals = append(literals, string(pattern[context.pos]))
		}

		context.pos++
	}

	if context.pos >= len(pattern) || pattern[context.pos] != ']' {
		fmt.Fprintf(os.Stderr, "[ERROR] Unclosed bracket: %s", pattern)
		os.Exit(1)
	}

	literalSet := map[byte]bool{}

	for _, literal := range literals {
		for c := literal[0]; c <= literal[len(literal)-1]; c++ {
			literalSet[c] = true
		}
	}

	context.tokens = append(context.tokens, token.Token{
		Type:  token.BRACKET,
		Value: literalSet,
	})
}
