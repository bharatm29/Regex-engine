package main

import (
	"fmt"
	"regex-engine/internals/regex"
)

func main() {
	for {
		fmt.Print("input< ")
		var input string
		fmt.Scanln(&input)

		fmt.Print("pattern< ")
		var pattern string
		fmt.Scanln(&pattern)

		if regex.Match(input, pattern) {
			fmt.Println("Match")
		} else {
			fmt.Println("Not Match")
		}
	}
}
