package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
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
		if isSafeRow(row) {
			ans += 1
		}
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	for _, row := range parsed {
		if isSafeRow(row) || isErrorTolerate(row) {
			ans += 1
		}
	}

	return ans
}

func isSafeRow(row []int) bool {
	isIncr := row[1] > row[0]
	isDecr := row[1] < row[0]
	if !isIncr && !isDecr {
		return false
	}

	prev := row[0]
	for i := 1; i < len(row); i++ {
		curr := row[i]
		if isIncr {
			if curr <= prev || curr-prev > 3 {
				return false
			}
		}
		if isDecr {
			if curr >= prev || prev-curr > 3 {
				return false
			}
		}
		prev = curr
	}

	return true
}

func isErrorTolerate(row []int) bool {
	for i := 0; i < len(row); i++ {
		if isSafeWithExclude(row, i) {
			return true
		}
	}

	return false
}

func isSafeWithExclude(row []int, ignore int) bool {
	isIncr := row[1] > row[0]
	if ignore < 2 {
		isIncr = row[3] > row[2]
	}

	isDecr := row[1] < row[0]
	if ignore < 2 {
		isDecr = row[3] < row[2]
	}

	if !isIncr && !isDecr {
		return false
	}

	prev := row[0]
	if ignore == 0 {
		prev = row[1]
	}

	start := 1
	if ignore < 2 {
		start = 2
	}

	for i := start; i < len(row); i++ {
		if i == ignore {
			continue
		}

		curr := row[i]
		if isIncr {
			if curr <= prev || curr-prev > 3 {
				return false
			}
		}
		if isDecr {
			if curr >= prev || prev-curr > 3 {
				return false
			}
		}
		prev = curr
	}

	return true
}

func parseInput(input string) [][]int {
	rows := strings.Split(input, "\n")
	ans := make([][]int, len(rows))

	for rowIndex, row := range rows {
		ans[rowIndex] = []int{}
		for _, valStr := range strings.Split(row, " ") {
			valInt, _ := strconv.Atoi(valStr)
			ans[rowIndex] = append(ans[rowIndex], valInt)
		}
	}
	return ans
}
