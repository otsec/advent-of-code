package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"math"
	"os/exec"
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

	for _, line := range parsed {
		numLine := lineToNums(line)
		ans += calcLinePart1(numLine)
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	for _, line := range parsed {
		numLine := lineToNums(line)
		ans += calcLinePart2(numLine, 12)
	}

	return ans
}

func lineToNums(line string) []int {
	nums := make([]int, len(line))
	for i, c := range line {
		nums[i] = int(c - '0')
	}
	return nums
}

func calcLinePart1(line []int) int {
	firstMaxIndex := 0
	firstMax := line[firstMaxIndex]
	for i := 1; i < len(line)-1; i++ {
		if line[i] > firstMax {
			firstMax = line[i]
			firstMaxIndex = i
		}
	}

	lastMax := line[firstMaxIndex+1]
	for i := firstMaxIndex + 1; i < len(line); i++ {
		if line[i] > lastMax {
			lastMax = line[i]
		}
	}

	return firstMax*10 + lastMax
}

func calcLinePart2(line []int, useDigits int) int {
	maxDigitIndex := 0
	maxDigit := line[0]
	for i := maxDigitIndex + 1; i < len(line)-useDigits+1; i++ {
		if line[i] > maxDigit {
			maxDigit = line[i]
			maxDigitIndex = i
		}
	}

	if useDigits == 1 {
		return maxDigit
	}

	pow := int(math.Pow(10, float64(useDigits-1)))
	return maxDigit*pow + calcLinePart2(line[maxDigitIndex+1:], useDigits-1)
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
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
