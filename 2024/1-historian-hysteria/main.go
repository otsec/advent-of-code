package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"sort"
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
	column1, column2 := parseInput(input)
	sort.Ints(column1)
	sort.Ints(column2)

	for i, val1 := range column1 {
		val2 := column2[i]

		if val2 > val1 {
			ans += val2 - val1
		} else {
			ans += val1 - val2
		}

	}

	return ans
}

func part2(input string) (ans int) {
	column1, column2 := parseInput(input)

	calcMap2 := map[int]int{}
	for _, val2 := range column2 {
		if _, ok := calcMap2[val2]; ok {
			calcMap2[val2] += 1
		} else {
			calcMap2[val2] = 1
		}
	}

	for _, val1 := range column1 {
		if count, ok := calcMap2[val1]; ok {
			ans += val1 * count
		}
	}

	_, _ = column1, column2

	return ans
}

func parseInput(input string) (column1 []int, column2 []int) {
	rows := strings.Split(input, "\n")
	for _, v := range rows {
		items := strings.Split(v, "   ")

		val1, _ := strconv.Atoi(items[0])
		column1 = append(column1, val1)

		val2, _ := strconv.Atoi(items[1])
		column2 = append(column2, val2)
	}

	return column1, column2
}
