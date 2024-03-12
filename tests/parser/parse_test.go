package parser_test

import (
	"fmt"
	"reflect"
	"regex-engine/internals/parser"
	"regex-engine/internals/token"
	"testing"
)

func TestParser(t *testing.T) {
	testcases := []struct {
		pattern string
		tokens  []token.Token
	}{
		// literals
		{
			pattern: "a",
			tokens: []token.Token{
				{Type: token.LITERAL, Value: byte('a')},
			},
		},

		// bracket
		{
			pattern: "(abc)",
			tokens: []token.Token{
				{Type: token.LITERAL, Value: byte('a')},
				{Type: token.LITERAL, Value: byte('b')},
				{Type: token.LITERAL, Value: byte('c')},
			},
		},

		{
			pattern: "[abc]",
			tokens: []token.Token{
				{Type: token.BRACKET, Value: map[byte]bool{
					byte('a'): true,
					byte('b'): true,
					byte('c'): true,
				}},
			},
		},
		{
			pattern: "[a-c]",
			tokens: []token.Token{
				{Type: token.BRACKET, Value: map[byte]bool{
					byte('a'): true,
					byte('b'): true,
					byte('c'): true,
				}},
			},
		},
		{
			pattern: "[ab-c]",
			tokens: []token.Token{
				{Type: token.BRACKET, Value: map[byte]bool{
					byte('a'): true,
					byte('b'): true,
					byte('c'): true,
				}},
			},
		},

		// Or
		{
			pattern: "a|b",
			tokens: []token.Token{
				{Type: token.OR, Value: []token.Token{
					{Type: token.LITERAL, Value: byte('a')},
				}},
				{Type: token.OR, Value: []token.Token{
					{Type: token.LITERAL, Value: byte('b')},
				}},
			},
		},
		{
			pattern: "[ab-c]|z",
			tokens: []token.Token{
				{Type: token.OR, Value: []token.Token{
					{Type: token.BRACKET, Value: map[byte]bool{
						byte('a'): true,
						byte('b'): true,
						byte('c'): true,
					}},
				}},
				{Type: token.OR, Value: []token.Token{
					{Type: token.LITERAL, Value: byte('z')},
				}},
			},
		},

		// Bracket repeat
		{
			pattern: "a{1,3}",
			tokens: []token.Token{
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 1,
					Max: 3,
				}},
			},
		},
		{
			pattern: "a{,3}c",
			tokens: []token.Token{
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 0,
					Max: 3,
				}},
				{Type: token.LITERAL, Value: byte('c')},
			},
		},
		{
			pattern: "ba{1,}",
			tokens: []token.Token{
				{Type: token.LITERAL, Value: byte('b')},
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 1,
					Max: parser.INFINITY,
				}},
			},
		},

		// *
		{
			pattern: "a*",
			tokens: []token.Token{
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 0,
					Max: parser.INFINITY,
				}},
			},
		},
		{
			pattern: "ba*",
			tokens: []token.Token{
				{Type: token.LITERAL, Value: byte('b')},
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 0,
					Max: parser.INFINITY,
				}},
			},
		},
		{
			pattern: "a*c",
			tokens: []token.Token{
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 0,
					Max: parser.INFINITY,
				}},
				{Type: token.LITERAL, Value: byte('c')},
			},
		},

		// +
		{
			pattern: "a+",
			tokens: []token.Token{
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 1,
					Max: parser.INFINITY,
				}},
			},
		},
		{
			pattern: "ba+",
			tokens: []token.Token{
				{Type: token.LITERAL, Value: byte('b')},
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 1,
					Max: parser.INFINITY,
				}},
			},
		},
		{
			pattern: "a+c",
			tokens: []token.Token{
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 1,
					Max: parser.INFINITY,
				}},
				{Type: token.LITERAL, Value: byte('c')},
			},
		},

		// ?
		{
			pattern: "a?",
			tokens: []token.Token{
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 0,
					Max: 1,
				}},
			},
		},
		{
			pattern: "ba?",
			tokens: []token.Token{
				{Type: token.LITERAL, Value: byte('b')},
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 0,
					Max: 1,
				}},
			},
		},
		{
			pattern: "a?c",
			tokens: []token.Token{
				{Type: token.REPEAT, Value: parser.RepeatValue{
					RepeatToken: token.Token{Type: token.LITERAL, Value: byte('a')},

					Min: 0,
					Max: 1,
				}},
				{Type: token.LITERAL, Value: byte('c')},
			},
		},
	}

	for _, test := range testcases {
		t.Run(fmt.Sprintf("Test for: %s", test.pattern), func(t *testing.T) {
			tokens := parser.Parse(test.pattern).GetTokens()

			if len(tokens) != len(test.tokens) {
				t.Logf("Expected %v, got %v", test.tokens, tokens)
				t.Fail()
			}

			for idx, token := range test.tokens {
				if !reflect.DeepEqual(token.Value, tokens[idx].Value) {
					t.Logf("Expected %v, got %v", test.tokens, tokens)
					t.Fail()
				}

				if token.Type != tokens[idx].Type {
					t.Logf("Expected %v, got %v", test.tokens, tokens)
					t.Fail()
				}
			}
		})
	}
}
