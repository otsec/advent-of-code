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
	cards := parseInput(input)

	for _, card := range cards {
		ans += calcCardPoints(card)
	}

	return ans
}

func part2(input string) (ans int) {
	cards := parseInput(input)

	copiesMemo := make([]int, len(cards))
	for index, card := range cards {
		copiesMemo[index] += 1

		copiesWon := calcWonCopies(card)
		for i := index + 1; i <= index+copiesWon; i++ {
			copiesMemo[i] += copiesMemo[index]
		}
	}

	for index, _ := range cards {
		ans += copiesMemo[index]
	}

	return ans
}

type Card struct {
	id      int
	winning []int
	having  []int
}

func parseInput(input string) []*Card {
	lines := strings.Split(input, "\n")
	cards := make([]*Card, len(lines))

	for index, line := range lines {
		colonIndex := strings.Index(line, ":")
		numLine := line[colonIndex+2:]
		numGroups := strings.Split(numLine, "|")
		winning := parseNumbers(numGroups[0])
		having := parseNumbers(numGroups[1])

		cards[index] = &Card{index + 1, winning, having}
	}

	return cards
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

func calcCardPoints(card *Card) int {
	var points int

	for _, num := range card.winning {
		if findInts(card.having, num) {
			if points == 0 {
				points = 1
			} else {
				points *= 2
			}
		}
	}

	return points
}

func calcWonCopies(card *Card) int {
	var ans int
	for _, num := range card.winning {
		if findInts(card.having, num) {
			ans += 1
		}
	}
	return ans
}

func findInts(haystack []int, needle int) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
