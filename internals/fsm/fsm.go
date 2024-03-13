package fsm

import (
	"os"
	"regex-engine/internals/parser"
	"regex-engine/internals/token"
)

const (
	epsilonChar   = '\x00'
	START_OF_TEXT = 1
	END_OF_TEXT   = 2
)

type state struct {
	transition map[byte][]*state

	terminal bool
	start    bool
}

func (s *state) Check(input string, pos int) bool {
	ch := readChar(input, pos)

	if ch == END_OF_TEXT && s.terminal {
		return true
	}

	if states := s.transition[ch]; len(states) > 0 {
		nextState := states[0]
		if nextState.Check(input, pos+1) {
			return true
		}
	}

	for _, s := range s.transition[epsilonChar] {
		if s.Check(input, pos) {
			return true
		}

		if ch == START_OF_TEXT && s.Check(input, pos+1) {
			return true
		}
	}

	return false
}

func ToNfa(ctx *parser.ParseContext) (*state, *state) {
	startState := &state{
		start:      true,
		transition: map[byte][]*state{},
	}

	endState := &state{
		terminal:   true,
		transition: map[byte][]*state{},
	}

	tokens := ctx.GetTokens()

	if len(tokens) == 0 {
		startState.transition[epsilonChar] = append(startState.transition[epsilonChar], endState)
		return startState, endState
	}

	start, end := toNfaToken(tokens[0])

	for i := 1; i < len(tokens); i++ {
		s, e := toNfaToken(tokens[i])

		end.transition[epsilonChar] = append(end.transition[epsilonChar], s)
		end = e
	}

	startState.transition[epsilonChar] = append(startState.transition[epsilonChar], start)
	end.transition[epsilonChar] = append(end.transition[epsilonChar], endState)

	return startState, endState
}

func toNfaToken(tok token.Token) (start, end *state) {
	startState := &state{
		transition: map[byte][]*state{},
	}

	endState := &state{
		transition: map[byte][]*state{},
	}

	switch tok.Type {
	case token.LITERAL:
		ch := tok.Value.(byte)
		startState.transition[ch] = append(startState.transition[ch], endState)

	case token.GROUP, token.UNCAPTURE_GROUP:
		toks := tok.Value.([]token.Token)

		if len(toks) == 0 {
			startState.transition[epsilonChar] = append(startState.transition[epsilonChar], endState)
			return startState, endState
		}

		for i := 0; i < len(toks); i++ {
			s, e := toNfaToken(toks[i])

			startState.transition[epsilonChar] = append(startState.transition[epsilonChar], s)
			e.transition[epsilonChar] = append(e.transition[epsilonChar], endState)
		}

	case token.OR:
		tokens := tok.Value.([]token.Token)

		leftTok := tokens[0]
		rightTok := tokens[1]

		leftStart, leftEnd := toNfaToken(leftTok)
		rightStart, rightEnd := toNfaToken(rightTok)

		startState.transition[epsilonChar] = append(startState.transition[epsilonChar], leftStart, rightStart)

		leftEnd.transition[epsilonChar] = append(leftEnd.transition[epsilonChar], endState)
		rightEnd.transition[epsilonChar] = append(rightEnd.transition[epsilonChar], endState)

	default:
		os.Exit(1)
	}

	return startState, endState
}

func readChar(input string, pos int) byte {
	if pos >= len(input) {
		return END_OF_TEXT
	} else if pos < 0 {
		return START_OF_TEXT
	} else {
		return input[pos]
	}
}
