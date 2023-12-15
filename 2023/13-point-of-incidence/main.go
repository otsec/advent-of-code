package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"reflect"
	"slices"
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
	fields := parseInput(input)

	for _, f := range fields {
		ans += findHorizontalReflection(&f) * 100
		ans += findHorizontalReflection(f.rotate())
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)
	_ = parsed

	return ans
}

type Field struct {
	lines []string
}

func (f *Field) rotate() *Field {
	newByteLines := make([][]byte, len(f.lines[0]))
	for i, _ := range newByteLines {
		newByteLines[i] = make([]byte, len(f.lines))
	}

	for i := 0; i < len(f.lines); i++ {
		for j := 0; j < len(f.lines[0]); j++ {
			newByteLines[j][i] = f.lines[i][j]
		}
	}

	newLines := make([]string, len(newByteLines))
	for i, line := range newByteLines {
		newLines[i] = string(line)
	}

	return &Field{newLines}
}

func parseInput(input string) []Field {
	groups := strings.Split(input, "\n\n")
	ans := make([]Field, len(groups))
	for i, group := range groups {
		lines := strings.Split(group, "\n")
		ans[i] = Field{lines}
	}
	return ans
}

func findHorizontalReflection(field *Field) int {
	for i := 0; i < len(field.lines)-1; i++ {
		if field.lines[i] == field.lines[i+1] {
			numberLinesBefore := i + 1
			numberLinesAfter := len(field.lines) - numberLinesBefore
			if numberLinesBefore < numberLinesAfter {
				continue
			}

			numberLinesCanCheck := min(numberLinesBefore, numberLinesAfter)
			//fmt.Println(i, numberLinesBefore, numberLinesAfter, numberLinesCanCheck)
			//fmt.Println(i+1-numberLinesCanCheck, i+1)
			//fmt.Println(field.lines[i+1-numberLinesCanCheck : i+1])
			//fmt.Println(i+1, i+1+numberLinesCanCheck)
			//fmt.Println(field.lines[i+1 : i+1+numberLinesCanCheck])

			linesBefore := field.lines[i+1-numberLinesCanCheck : i+1]
			linesAfter := field.lines[i+1 : i+1+numberLinesCanCheck]
			slices.Reverse(linesAfter)

			//fmt.Println(i, linesBefore, linesAfter)

			if reflect.DeepEqual(linesBefore, linesAfter) {
				return i + 1
			}
		}
	}

	return 0
}
