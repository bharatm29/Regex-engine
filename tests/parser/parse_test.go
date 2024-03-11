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
		{
			pattern: "a",
			tokens: []token.Token{
				{Type: token.LITERAL, Value: 'a'},
			},
		},
		{
			pattern: "(abc)",
			tokens: []token.Token{
				{Type: token.LITERAL, Value: 'a'},
				{Type: token.LITERAL, Value: 'b'},
				{Type: token.LITERAL, Value: 'c'},
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
