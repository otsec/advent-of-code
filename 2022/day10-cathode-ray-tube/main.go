package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"strings"

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
	operations := parseInput(input)

	cycle := 1
	x := 1
	for _, operation := range operations {
		for i := 1; i <= operation.cycles; i++ {
			cycle++

			if operation.command == "addx" && i == operation.cycles {
				x += operation.argument
			}

			if (cycle-20)%40 == 0 {
				ans += x * cycle
			}
		}
	}

	return ans
}

func part2(input string) string {
	operations := parseInput(input)

	crt := make([]rune, 240)
	cycle := 0
	x := 1
	for _, operation := range operations {
		for i := 1; i <= operation.cycles; i++ {
			cycle++

			pixel := '.'
			pixelIndex := cycle % 40
			if pixelIndex == 0 {
				pixelIndex = 40
			}
			if x <= pixelIndex && pixelIndex < x+3 {
				pixel = '#'
			}
			crt[cycle-1] = pixel

			if operation.command == "addx" && i == operation.cycles {
				x += operation.argument
			}
		}
	}

	//for i, pixel := range crt {
	//	fmt.Print(string(pixel))
	//
	//	cycle := i + 1
	//	if cycle%40 == 0 {
	//		fmt.Println()
	//	}
	//}
	//fmt.Println()

	return string(crt)
}

type Operation struct {
	command  string
	argument int
	cycles   int
}

func parseInput(input string) []Operation {
	lines := strings.Split(input, "\n")

	ans := make([]Operation, len(lines))
	for i, line := range lines {
		if strings.HasPrefix(line, "noop") {
			ans[i] = Operation{"noop", 0, 1}
		}
		if strings.HasPrefix(line, "addx") {
			argument := cast.ToInt(strings.Replace(line, "addx ", "", 1))
			ans[i] = Operation{"addx", argument, 2}
		}
	}

	return ans
}
