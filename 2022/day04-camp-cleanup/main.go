package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/alexchao26/advent-of-code-go/cast"
	"github.com/alexchao26/advent-of-code-go/util"
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

	for _, group := range parsed {
		from1, to1, from2, to2 := group[0], group[1], group[2], group[3]
		if from1 <= from2 && to2 <= to1 {
			ans++
		} else if from2 <= from1 && to1 <= to2 {
			ans++
		}
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	for _, group := range parsed {
		from1, to1, from2, to2 := group[0], group[1], group[2], group[3]

		if from1 <= to2 && from2 <= to1 {
			ans++
		} else if from2 <= to1 && from1 <= to2 {
			ans++
		}
	}

	return ans
}

func parseInput(input string) (ans [][4]int) {
	for _, line := range strings.Split(input, "\n") {
		elves := strings.Split(line, ",")
		elf1 := strings.Split(elves[0], "-")
		elf2 := strings.Split(elves[1], "-")

		ans = append(ans, [4]int{
			cast.ToInt(elf1[0]),
			cast.ToInt(elf1[1]),
			cast.ToInt(elf2[0]),
			cast.ToInt(elf2[1]),
		})
	}
	return ans
}
