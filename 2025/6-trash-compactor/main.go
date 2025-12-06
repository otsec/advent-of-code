package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os/exec"
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
		CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) (ans int) {
	parsed := parseInput(input)

	for _, column := range parsed {
		curr := column.rawNums[0].Int()
		for i := 1; i < len(column.rawNums); i++ {
			if column.operator == '+' {
				curr += column.rawNums[i].Int()
			} else {
				curr *= column.rawNums[i].Int()
			}
		}
		ans += curr
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	for _, column := range parsed {
		nums := column.ReadRightLeft()
		curr := nums[0]
		for i := 1; i < len(nums); i++ {
			if column.operator == '+' {
				curr += nums[i]
			} else {
				curr *= nums[i]
			}
		}
		ans += curr
	}

	return ans
}

type RawNum string

func (r RawNum) Int() int {
	val := strings.TrimSpace(string(r))
	return toInt(val)
}

type RawNumColumn struct {
	operator rune
	rawNums  []RawNum
}

func (rc RawNumColumn) ReadRightLeft() []int {
	var ans []int

	size := len(rc.rawNums[0])
	for i := size - 1; i >= 0; i-- {
		curr := 0
		for _, rawNum := range rc.rawNums {
			if rawNum[i] == ' ' {
				continue
			}

			digit := int(rawNum[i] - '0')
			if curr == 0 {
				curr = digit
			} else {
				curr = curr*10 + digit
			}
		}
		ans = append(ans, curr)
	}

	return ans
}

func parseInput(input string) []RawNumColumn {
	lines := strings.Split(input, "\n")

	numLines := lines[:len(lines)-1]
	lastLine := lines[len(lines)-1]

	var ans []RawNumColumn
	currSymbol := lastLine[0]
	currSymbolIndex := 0
	for i := 1; i < len(lastLine); i++ {
		if lastLine[i] == ' ' {
			continue
		}

		numStartAt := currSymbolIndex
		numEndsAt := i - 1
		rawNums := make([]RawNum, len(numLines))
		for j, line := range numLines {
			rawNums[j] = RawNum(line[numStartAt:numEndsAt])
		}

		ans = append(ans, RawNumColumn{operator: rune(currSymbol), rawNums: rawNums})

		currSymbol = lastLine[i]
		currSymbolIndex = i
	}

	// process last column
	numStartAt := currSymbolIndex
	numEndsAt := len(lastLine)
	rawNums := make([]RawNum, len(numLines))
	for j, line := range numLines {
		rawNums[j] = RawNum(line[numStartAt:numEndsAt])
	}
	ans = append(ans, RawNumColumn{operator: rune(currSymbol), rawNums: rawNums})

	return ans
}

func toInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return num
}

// CopyToClipboard is for macOS
func CopyToClipboard(text string) error {
	command := exec.Command("pbcopy")
	command.Stdin = bytes.NewReader([]byte(text))

	if err := command.Start(); err != nil {
		return fmt.Errorf("error starting pbcopy command: %w", err)
	}

	err := command.Wait()
	if err != nil {
		return fmt.Errorf("error running pbcopy %w", err)
	}

	return nil
}
