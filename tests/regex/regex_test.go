package regex_test

import (
	"fmt"
	"regex-engine/internals/regex"
	"testing"
)

func TestRegex(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		match   bool
	}{
		{
			pattern: "a",
			input:   "a",
			match:   true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Test for: [%s] on [%s]", tt.pattern, tt.input), func(t *testing.T) {
			actual := regex.Match(tt.input, tt.pattern)

			if actual != tt.match {
                t.Logf("Expected %t, got %t: [%s] on [%s]", tt.match, actual, tt.pattern, tt.input)
				t.Fail()
			}
		})
	}
}
