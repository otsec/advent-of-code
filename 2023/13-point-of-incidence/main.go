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
	fields := parseInput(input)

	for _, f := range fields {
		ans += findHorizontalReflectionWithDiff(&f) * 100
		ans += findHorizontalReflectionWithDiff(f.rotate())
	}

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

func (f *Field) Print() {
	fmt.Println(strings.Join(f.lines, "\n"))
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
			numberLinesCanCheck := min(numberLinesBefore, numberLinesAfter)

			linesBefore := field.lines[i+1-numberLinesCanCheck : i+1]
			linesAfter := field.lines[i+1 : i+1+numberLinesCanCheck]

			linesAfterReversed := make([]string, len(linesAfter))
			copy(linesAfterReversed, linesAfter)
			slices.Reverse(linesAfterReversed)

			if reflect.DeepEqual(linesBefore, linesAfterReversed) {
				return i + 1
			}
		}
	}

	return 0
}

func findHorizontalReflectionWithDiff(field *Field) int {
	for i := 0; i < len(field.lines)-1; i++ {
		//diff := getDiff(field.lines[i], field.lines[i+1])
		//if diff > 1 {
		//	continue
		//}

		t, b := i, i+1
		diff := 0
		for {
			diff += getDiff(field.lines[t], field.lines[b])
			if diff > 1 {
				break
			}
			if t == 0 || b == len(field.lines)-1 {
				break
			} else {
				t--
				b++
			}
		}

		if diff == 1 {
			return i + 1
		}
	}

	return 0
}

func getDiff(line1, line2 string) int {
	ans := 0
	for i := 0; i < len(line1); i++ {
		if line1[i] != line2[i] {
			ans++
		}
	}
	return ans
}
