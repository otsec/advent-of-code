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

	for _, r := range parsed {
		for i := r.start; i <= r.end; i++ {
			if isInvalidPart1(i) {
				ans += i
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	for _, r := range parsed {
		for i := r.start; i <= r.end; i++ {
			if isInvalidPart2(i) {
				ans += i
			}
		}
	}

	return ans
}

type IDRange struct {
	start, end int
}

func isInvalidPart1(id int) bool {
	str := strconv.Itoa(id)
	if len(str)%2 != 0 {
		return false
	}
	firstHalf := str[:len(str)/2]
	secondHalf := str[len(str)/2:]
	return firstHalf == secondHalf
}

func isInvalidPart2(id int) bool {
	str := strconv.Itoa(id)
mainLoop:
	for i := 1; i < len(str); i++ {
		if len(str)%i != 0 {
			continue
		}

		firstPart := str[:i]
		for j := i; j < len(str); j += i {
			nextPart := str[j : j+i]
			if firstPart != nextPart {
				continue mainLoop
			}
		}

		return true
	}
	return false
}

func parseInput(input string) []IDRange {
	rawRange := strings.Split(input, ",")
	ans := make([]IDRange, len(rawRange))
	for i, raw := range rawRange {
		segments := strings.Split(raw, "-")
		ans[i] = IDRange{parseInt(segments[0]), parseInt(segments[1])}
	}
	return ans
}

func parseInt(input string) int {
	val, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	return val
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
