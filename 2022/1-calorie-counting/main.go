package main

import (
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var part int
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.Parse()
	fmt.Println("Running part", part)

	if part == 1 {
		ans := part1(util.ReadFile("./input.txt"))
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(util.ReadFile("./input.txt"))
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) int {
	lines := strings.Split(input, "\n")

	max := 0
	curr := 0
	for _, item := range lines {
		if item == "" {
			if curr > max {
				max = curr
			}
			curr = 0
		} else {
			if cal, err := strconv.Atoi(item); err == nil {
				curr += cal
			}
		}
	}

	return max
}

func part2(input string) int {
	lines := strings.Split(input, "\n")

	elves := []int{}
	curr := 0
	for _, item := range lines {
		if item == "" {
			elves = append(elves, curr)
			curr = 0
		} else {
			if cal, err := strconv.Atoi(item); err == nil {
				curr += cal
			}
		}
	}

	sort.Ints(elves)

	return elves[len(elves)-1] + elves[len(elves)-2] + elves[len(elves)-3]
}
