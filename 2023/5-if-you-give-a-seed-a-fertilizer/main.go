package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"math"
	"regexp"
	"slices"
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
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) (ans int) {
	seeds, conversions := parseInput(input)

	for i, _ := range seeds {
		for _, group := range conversions {
			seeds[i] = group.convert(seeds[i])
		}
	}

	return slices.Min(seeds)
}

func part2(input string) (ans int) {
	seeds, conversions := parseInput(input)
	seedRanges := parseSeedRanges(seeds)

	prevRanges := seedRanges
	for _, group := range conversions {
		nextRanges := []SeedRange{}
		for _, seedRange := range prevRanges {
			nextRanges = append(nextRanges, group.convertRange(seedRange)...)
		}
		prevRanges = nextRanges
	}

	ans = math.MaxInt
	for _, seedRange := range prevRanges {
		ans = min(ans, seedRange.start)
	}

	return ans
}

type ConvertMap struct {
	start, end, diff int
}

type SeedRange struct {
	start, end int
}

type ConvertGroup struct {
	name  string
	items []ConvertMap
}

func (gc *ConvertGroup) convert(seed int) int {
	for _, cm := range gc.items {
		if cm.start <= seed && seed < cm.end {
			return seed + cm.diff
		}
	}
	return seed
}

func (gc *ConvertGroup) convertRange(input SeedRange) []SeedRange {
	var ans []SeedRange

	iter := []SeedRange{input}
	for _, cm := range gc.items {
		nextIter := []SeedRange{}

		for _, sr := range iter {
			if sr.end < cm.start || cm.end < sr.start {
				nextIter = append(nextIter, sr)
				continue
			}

			if sr.start < cm.end && sr.end > cm.start {
				start := max(cm.start, sr.start)
				end := min(cm.end, sr.end)
				ans = append(ans, SeedRange{start + cm.diff, end + cm.diff})
			}

			if sr.start < cm.start {
				start := sr.start
				end := min(sr.end, cm.start)
				nextIter = append(nextIter, SeedRange{start, end})
			}

			if sr.end > cm.end {
				start := min(sr.end, cm.end)
				end := sr.end
				nextIter = append(nextIter, SeedRange{start, end})
			}
		}

		iter = nextIter
	}

	for _, sr := range iter {
		ans = append(ans, sr)
	}

	return ans
}

func parseInput(input string) ([]int, []ConvertGroup) {
	inputGroups := strings.Split(input, "\n\n")
	seeds := parseNumbers(inputGroups[0])

	convertGroups := make([]ConvertGroup, len(inputGroups)-1)
	for i := 1; i < len(inputGroups); i++ {
		lines := strings.Split(inputGroups[i], "\n")
		convertGroups[i-1] = ConvertGroup{name: lines[0], items: make([]ConvertMap, len(lines)-1)}
		for j := 1; j < len(lines); j++ {
			nums := parseNumbers(lines[j])
			convertGroups[i-1].items[j-1] = ConvertMap{nums[1], nums[1] + nums[2], nums[0] - nums[1]}
		}
	}

	return seeds, convertGroups
}

func parseNumbers(input string) []int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(input, -1)
	nums := make([]int, len(matches))
	for i, match := range matches {
		nums[i], _ = strconv.Atoi(match)
	}
	return nums
}

func parseSeedRanges(seeds []int) []SeedRange {
	var ans []SeedRange
	for i := 0; i < len(seeds); i += 2 {
		ans = append(ans, SeedRange{seeds[i], seeds[i] + seeds[i+1]})
	}
	return ans
}
