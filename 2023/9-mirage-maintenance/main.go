package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"regexp"
	"slices"
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

	for _, item := range parsed {
		tower := [][]int{item}
		for !areAllZeroes(last(tower)) {
			tower = append(tower, calculateDiffs(last(tower)))
		}

		slices.Reverse(tower)
		for i := 0; i < len(tower); i++ {
			if i == 0 {
				tower[i] = append(tower[i], 0)
			} else {
				tower[i] = append(tower[i], last(tower[i])+last(tower[i-1]))
			}
		}

		ans += last(last(tower))
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	for _, item := range parsed {
		tower := [][]int{item}
		for !areAllZeroes(last(tower)) {
			tower = append(tower, calculateDiffs(last(tower)))
		}

		slices.Reverse(tower)
		for i := 0; i < len(tower); i++ {
			slices.Reverse(tower[i])
			if i == 0 {
				tower[i] = append(tower[i], 0)
			} else {
				tower[i] = append(tower[i], last(tower[i])-last(tower[i-1]))
			}
		}

		ans += last(last(tower))
	}

	return ans
}

func parseInput(input string) [][]int {
	lines := strings.Split(input, "\n")

	ans := make([][]int, len(lines))
	for i, line := range lines {
		ans[i] = parseNumbers(line)
	}

	return ans
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

func last[E any](s []E) E {
	return s[len(s)-1]
}

func calculateDiffs(input []int) []int {
	ans := make([]int, len(input)-1)
	for i := 0; i < len(input)-1; i++ {
		ans[i] = input[i+1] - input[i]
	}
	return ans
}

func areAllZeroes(input []int) bool {
	for _, v := range input {
		if v != 0 {
			return false
		}
	}
	return true
}
