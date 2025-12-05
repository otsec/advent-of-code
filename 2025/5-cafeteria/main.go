package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os/exec"
	"regexp"
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
	freshRanges, productIds := parseInput(input)

	for _, id := range productIds {
		for _, fr := range freshRanges {
			if fr.Includes(id) {
				ans++
				break
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	freshRanges, _ := parseInput(input)

	var squashed bool
	for {
		freshRanges, squashed = squashRanges(freshRanges)
		if !squashed {
			break
		}
	}

	for _, fr := range freshRanges {
		ans += fr.high - fr.low + 1
	}

	return ans
}

func squashRanges(ranges []FreshRange) ([]FreshRange, bool) {
	newRanges := []FreshRange{}
	skipped := map[int]bool{}
	for i := 0; i < len(ranges); i++ {
		if skipped[i] {
			continue
		}

		curr := ranges[i]

		for j := i + 1; j < len(ranges); j++ {
			if skipped[j] {
				continue
			}

			merged, ok := mergeRanges(curr, ranges[j])
			if ok {
				curr = merged
				skipped[j] = true
			}
		}

		newRanges = append(newRanges, curr)
	}

	return newRanges, len(skipped) > 0
}

func mergeRanges(r1 FreshRange, r2 FreshRange) (FreshRange, bool) {
	if r1.high < r2.low || r2.high < r1.low {
		return FreshRange{}, false
	}
	return FreshRange{low: min(r1.low, r2.low), high: max(r1.high, r2.high)}, true
}

type FreshRange struct {
	low, high int
}

func (fr FreshRange) Includes(id int) bool {
	return fr.low <= id && id <= fr.high
}

func parseInput(input string) ([]FreshRange, []int) {
	blocks := strings.Split(input, "\n\n")
	if len(blocks) != 2 {
		panic("expected exactly 2 blocks")
	}

	return parseFreshRanges(blocks[0]), parseIds(blocks[1])
}

func parseFreshRanges(input string) []FreshRange {
	lines := strings.Split(input, "\n")
	ranges := make([]FreshRange, len(lines))

	re := regexp.MustCompile(`\d+`)
	for i, line := range lines {
		matches := re.FindAllString(line, -1)
		ranges[i] = FreshRange{
			low:  toInt(matches[0]),
			high: toInt(matches[1]),
		}
	}

	return ranges
}

func parseIds(input string) []int {
	lines := strings.Split(input, "\n")
	nums := make([]int, len(lines))

	for i, line := range lines {
		nums[i] = toInt(line)
	}
	return nums
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
