package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"regexp"
	"strconv"
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
	for _, pair := range extractMuls(parseInput(input)) {
		ans += pair[0] * pair[1]
	}

	return ans
}

func part2(input string) (ans int) {
	for _, pair := range extractMuls(removeDonts(parseInput(input))) {
		ans += pair[0] * pair[1]
	}

	return ans
}

func parseInput(input string) string {
	return strings.Replace(input, "\n", "", -1)
}

func extractMuls(input string) [][]int {
	re := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	matches := re.FindAllStringSubmatch(input, -1)

	res := make([][]int, len(matches))
	for i, m := range matches {
		v1, _ := strconv.Atoi(m[1])
		v2, _ := strconv.Atoi(m[2])
		res[i] = []int{v1, v2}
	}

	return res
}

func removeDonts(input string) string {
	dontIndex := strings.Index(input, "don't()")
	if dontIndex == -1 {
		return input
	}

	doIndex := strings.Index(input[dontIndex+7:], "do()")
	if doIndex == -1 {
		return input[:dontIndex]
	}

	return removeDonts(input[:dontIndex] + input[dontIndex+7+doIndex+4:])
}
