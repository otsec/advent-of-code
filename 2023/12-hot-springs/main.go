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
	parsed := parseInput(input)

	for _, row := range parsed {
		ans += testPossibleVariants(row.symbols, row.groups)
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)
	_ = parsed

	return ans
}

type InputRow struct {
	symbols string
	groups  []int
}

func parseInput(input string) []InputRow {
	lines := strings.Split(input, "\n")

	parsed := make([]InputRow, len(lines))
	for i, line := range lines {
		segments := strings.Split(line, " ")
		parsed[i] = InputRow{segments[0], parseNumbers(segments[1])}
	}
	return parsed
}

func parseNumbers(input string) []int {
	re := regexp.MustCompile(`[-0-9]+`)
	matches := re.FindAllString(input, -1)
	nums := make([]int, len(matches))
	for i, match := range matches {
		nums[i], _ = strconv.Atoi(match)
	}
	return nums
}

func testPossibleVariants(input string, numbers []int) (ans int) {
	var indexes []int
	for i := 0; i < len(input); i++ {
		if input[i] == '?' {
			indexes = append(indexes, i)
		}
	}

	if len(indexes) == 0 {
		if isResultMatched(input, numbers) {
			return 1
		} else {
			return 0
		}
	}

	variant1 := strings.Replace(input, "?", "#", 1)
	ans += testPossibleVariants(variant1, numbers)

	variant2 := strings.Replace(input, "?", ".", 1)
	ans += testPossibleVariants(variant2, numbers)

	return ans
}

func isResultMatched(input string, numbers []int) bool {
	re := regexp.MustCompile(`#+`)

	matches := re.FindAllString(input, -1)
	if len(matches) != len(numbers) {
		return false
	}
	for i, m := range matches {
		if numbers[i] != len(m) {
			return false
		}
	}

	return true
}
