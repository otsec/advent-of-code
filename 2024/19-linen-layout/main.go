package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
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
	towels, patterns := parseInput(input)

	for _, pattern := range patterns {
		if canMatch(pattern, &towels) {
			ans++
		}
	}

	return ans
}

func part2(input string) (ans int) {
	towels, patterns := parseInput(input)

	memo := map[string]int{}
	for _, pattern := range patterns {
		res := countMatches(pattern, &towels, &memo)
		ans += res
	}

	return ans
}

func canMatch(pattern string, towels *[]string) bool {
	for _, towel := range *towels {
		if pattern == towel {
			return true
		}
		if strings.HasPrefix(pattern, towel) {
			if canMatch(pattern[len(towel):], towels) {
				return true
			}
		}
	}
	return false
}

func countMatches(pattern string, towels *[]string, memo *map[string]int) (ans int) {
	for _, towel := range *towels {
		if len(pattern) < len(towel) {
			continue
		}

		if pattern == towel {
			ans += 1
			continue
		}

		if strings.HasPrefix(pattern, towel) {
			newPattern := pattern[len(towel):]
			if val, ok := (*memo)[newPattern]; ok {
				ans += val
			} else {
				val = countMatches(newPattern, towels, memo)
				(*memo)[newPattern] = val
				ans += val
			}
		}
	}

	return ans
}

func parseInput(input string) (towels []string, patterns []string) {
	segments := strings.Split(input, "\n\n")
	towels = strings.Split(segments[0], ", ")
	patterns = strings.Split(segments[1], "\n")
	return towels, patterns
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
