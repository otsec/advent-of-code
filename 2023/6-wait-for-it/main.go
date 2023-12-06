package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
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

	ans = 1
	for _, race := range parsed {
		ans *= beatNumbers(race)
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(strings.ReplaceAll(input, " ", ""))

	ans = 1
	for _, race := range parsed {
		ans *= beatNumbers(race)
	}

	return ans
}

type Race struct {
	time     int
	distance int
}

func parseInput(input string) []Race {
	lines := strings.Split(input, "\n")

	time := parseNumbers(lines[0])
	dist := parseNumbers(lines[1])

	var ans []Race
	for i, _ := range time {
		ans = append(ans, Race{time[i], dist[i]})
	}

	return ans
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

func beatNumbers(race Race) (ans int) {
	for i := 1; i < race.time; i++ {
		result := i * (race.time - i)
		if result > race.distance {
			ans += 1
		}
	}

	return ans
}
