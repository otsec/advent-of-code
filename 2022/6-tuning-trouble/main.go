package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"strings"
)

//go:embed input.txt
var input string

func init() {
	// do this in init (not main) so test file has same input
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		panic("empty input.txt file")
	}
}

func main() {
	var part int
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.Parse()
	fmt.Println("Running part", part)

	if part == 1 {
		ans := part1(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) (ans int) {
	parsed := parseInput(input)

	for index, _ := range input {
		if isUnique(parsed, index, 4) {
			return index + 4
		}
	}

	return 0
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	for index, _ := range input {
		if isUnique(parsed, index, 14) {
			return index + 14
		}
	}

	return 0
}

func isUnique(input []rune, from int, len int) bool {
	mapping := map[rune]bool{}

	for i := from; i < from+len; i++ {
		if _, exists := mapping[input[i]]; exists {
			return false
		}

		mapping[input[i]] = true
	}

	return true
}

func parseInput(input string) []rune {
	return []rune(input)
}
