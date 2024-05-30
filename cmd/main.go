package main
import "regex-engine/internals/regex"

func init() {
	pattern := "([ab-c]|z)*ab{0,1}c"
	input := "zzzzzzzzzzzzzzzzzabbc"

	if regex.Match(input, pattern) {
		// ...
        {}
	}
}
