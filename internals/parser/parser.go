package parser

import (
	"fmt"
	"os"
	"regex-engine/internals/token"
	"strconv"
	"strings"
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
		parseOr(pattern, context)
	case '{': // {3, } {3, 4} {,10}
		parseRepeat(pattern, context)
	case '*', '?', '+': // a*, a?, a+
		parseRepeat(pattern, context)
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

func parseOr(pattern string, context *ParseContext) {
	context.pos++ // skipping |

	rightContext := &ParseContext{
		pos:    context.pos,
		tokens: []token.Token{},
	}

	for rightContext.pos < len(pattern) && pattern[rightContext.pos] != ')' {
		parsePattern(pattern, rightContext)
		rightContext.pos++
	}

	left := token.Token{
		Type:  token.OR,
		Value: context.tokens,
	}

	right := token.Token{
		Type:  token.OR,
		Value: rightContext.tokens,
	}

	context.pos = rightContext.pos

	context.tokens = []token.Token{left, right}
}

type RepeatValue struct {
	RepeatToken token.Token

	Min int
	Max int
}

const INFINITY = -1

func parseRepeat(pattern string, context *ParseContext) {
	switch pattern[context.pos] {
	case '{':
		parseRepeatBracket(pattern, context)
	case '*':
		makeRepeat(0, INFINITY, context)

	case '+':
		makeRepeat(1, INFINITY, context)

	case '?':
		makeRepeat(0, 1, context)

	default:
		fmt.Fprintf(os.Stderr, "Unreachable repeat character: %c", pattern[context.pos])
		os.Exit(1)
	}
}

func parseRepeatBracket(pattern string, context *ParseContext) {
	context.pos++ // skip {
	pos := context.pos

	for context.pos < len(pattern) && pattern[context.pos] != '}' {
		context.pos++
	}

	expr := pattern[pos:context.pos]

	split := strings.Split(expr, ",")

	rep := RepeatValue{}

	// TODO: What to do if both are empty? {,}

	if split[0] == "" {
		rep.Min = 0
	} else {
		val, err := strconv.Atoi(split[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not convert %s to int", split[0])
			os.Exit(1)
		}

		rep.Min = val
	}

	if split[len(split)-1] == "" {
		rep.Max = INFINITY
	} else {
		val, err := strconv.Atoi(split[len(split)-1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not convert %s to int", split[len(split)-1])
			os.Exit(1)
		}

		rep.Max = val
	}

	rep.RepeatToken = context.tokens[len(context.tokens)-1]

	context.tokens[len(context.tokens)-1] = token.Token{
		Type:  token.REPEAT,
		Value: rep,
	}
}

func makeRepeat(min, max int, context *ParseContext) {
	rep := RepeatValue{}
	rep.Min = min
	rep.Max = max

	rep.RepeatToken = context.tokens[len(context.tokens)-1]

	context.tokens[len(context.tokens)-1] = token.Token {
		Type:  token.REPEAT,
		Value: rep,
	}
}
