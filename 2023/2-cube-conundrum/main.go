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

	for _, game := range parsed {
		if maxInts(game.red) <= 12 && maxInts(game.green) <= 13 && maxInts(game.blue) <= 14 {
			ans += game.id
		}
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	for _, game := range parsed {
		ans += maxInts(game.red) * maxInts(game.green) * maxInts(game.blue)
	}

	return ans
}

type Game struct {
	id    int
	red   []int
	green []int
	blue  []int
}

func parseInput(input string) []*Game {
	var games []*Game

	regexpRed := regexp.MustCompile(`(\d+) red`)
	regexpGreen := regexp.MustCompile(`(\d+) green`)
	regexpBlue := regexp.MustCompile(`(\d+) blue`)

	lines := strings.Split(input, "\n")
	for index, line := range lines {
		game := &Game{
			id:    index + 1,
			red:   parseColor(line, regexpRed),
			green: parseColor(line, regexpGreen),
			blue:  parseColor(line, regexpBlue),
		}

		games = append(games, game)
	}

	return games
}

func parseColor(record string, matcher *regexp.Regexp) []int {
	var cubes []int

	cubesSubmatch := matcher.FindAllStringSubmatch(record, -1)
	for _, match := range cubesSubmatch {
		val, _ := strconv.Atoi(match[1])
		cubes = append(cubes, val)
	}

	return cubes
}

func maxInts(items []int) (ans int) {
	for _, v := range items {
		if v > ans {
			ans = v
		}
	}
	return ans
}
