package regex

import (
	"regex-engine/internals/fsm"
	"regex-engine/internals/parser"
)

// first parses then returns the nfa
func Match(input, pattern string) bool {
	ctx := parser.Parse(pattern)
	state, _ := fsm.ToNfa(ctx)

	return state.Check(input, 0)
}
