# Regex-engine

Implementing a Regex Engine in Golang while learning about Finite State Machines :)

### Usage

```go
	pattern := "([ab-c]|z)*ab{0,1}c"
	input := "zzzzzzzzzzzzzzzzzabbc"

	if regex.Match(input, pattern) {
		// ...
	}
```
