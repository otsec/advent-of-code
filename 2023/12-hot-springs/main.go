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
	// do this in init (not main) so test file has same symbols
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		panic("empty symbols.txt file")
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
		ans += testPossibleVariants(row)
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	for i, row := range parsed {
		largeRow := makeItLarger(row)
		ans += testPossibleVariants(largeRow)
		fmt.Println("row", i)
	}

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

func testPossibleVariants(input InputRow) (ans int) {
	var indexes []int
	for i := 0; i < len(input.symbols); i++ {
		if input.symbols[i] == '?' {
			indexes = append(indexes, i)
		}
	}

	if len(indexes) == 0 {
		if isResultMatched(input) {
			return 1
		} else {
			return 0
		}
	}

	if !makesSense(input) {
		return 0
	}

	variant1 := strings.Replace(input.symbols, "?", "#", 1)
	ans += testPossibleVariants(InputRow{variant1, input.groups})

	variant2 := strings.Replace(input.symbols, "?", ".", 1)
	ans += testPossibleVariants(InputRow{variant2, input.groups})

	return ans
}

func makesSense(input InputRow) bool {
	stack := 0
	n := 0
	for i := 0; i < len(input.symbols); i++ {
		if input.symbols[i] == '?' {
			return true
		}

		if input.symbols[i] == '#' {
			stack++
			continue
		}

		if input.symbols[i] == '.' && stack > 0 {
			if n >= len(input.groups) {
				return false
			} else if stack != input.groups[n] {
				return false
			} else {
				stack = 0
				n += 1
			}
		}
	}

	return true
}

func shrink(input InputRow) InputRow {
	if input.symbols[0] == '?' {
		return input
	}

	if input.symbols[0] == '.' {
		newSymbols := strings.TrimLeft(input.symbols, ".")
		input = InputRow{newSymbols, input.groups}
	}

	return input
}

func isResultMatched(input InputRow) bool {
	re := regexp.MustCompile(`#+`)

	matches := re.FindAllString(input.symbols, -1)
	if len(matches) != len(input.groups) {
		return false
	}
	for i, m := range matches {
		if input.groups[i] != len(m) {
			return false
		}
	}

	return true
}

func makeItLarger(input InputRow) InputRow {
	largeSymbols := input.symbols
	largeGroups := input.groups
	for i := 0; i < 4; i++ {
		largeSymbols = largeSymbols + "?" + input.symbols
		largeGroups = append(largeGroups, input.groups...)
	}

	return InputRow{largeSymbols, largeGroups}
}
