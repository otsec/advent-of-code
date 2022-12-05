package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"regexp"
	"strings"

	"github.com/alexchao26/advent-of-code-go/util"
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

func part1(input string) (ans string) {
	crates, moves := parseInput(input)

	for _, move := range moves {
		towerFrom := crates[move.from]
		towerTo := crates[move.to]

		for i := 0; i < move.count; i++ {
			towerTo.Push(towerFrom.Pop())
		}
	}

	for i := 1; i <= len(crates); i++ {
		ans += string(crates[i].Pop())
	}

	return ans
}

func part2(input string) (ans string) {
	crates, moves := parseInput(input)

	for _, move := range moves {
		towerFrom := crates[move.from]
		towerTo := crates[move.to]

		middle := CratesTower{}
		for i := 0; i < move.count; i++ {
			middle.Push(towerFrom.Pop())
		}
		for i := 0; i < move.count; i++ {
			towerTo.Push(middle.Pop())
		}
	}

	for i := 1; i <= len(crates); i++ {
		ans += string(crates[i].Pop())
	}

	return ans
}

type Crate string

type CratesTower []Crate

func (c *CratesTower) Push(item Crate) {
	*c = append(*c, item)
}

func (c *CratesTower) Pop() Crate {
	old := *c
	n := len(old)
	x := old[n-1]
	*c = old[0 : n-1]
	return x
}

//func (c *CratesTower) Last() Crate {
//	old := *c
//	n := len(old)
//	x := old[n-1]
//	return x
//}

type LandingDeck map[int]*CratesTower

type Move struct {
	from  int
	to    int
	count int
}

func parseInput(input string) (crates LandingDeck, moves []Move) {
	segments := strings.Split(input, "\n\n")

	inputCrates := strings.Split(segments[0], "\n")
	inputMoves := strings.Split(segments[1], "\n")

	// Parse Crates
	crates = make(LandingDeck)
	for col, numSymbol := range inputCrates[len(inputCrates)-1] {
		if numSymbol == ' ' {
			continue
		}

		num := cast.ToInt(cast.ToString(numSymbol))
		crates[num] = &CratesTower{}

		for row := len(inputCrates) - 2; row >= 0; row-- {
			if len(inputCrates[row]) < col {
				continue
			}

			crateSymbol := cast.ToString(inputCrates[row][col])
			if crateSymbol != " " {
				crates[num].Push(Crate(crateSymbol))
			}
		}
	}

	// Parse Moves
	for _, line := range inputMoves {
		re := regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)
		matches := re.FindAllStringSubmatch(line, 10)
		from, to, count := cast.ToInt(matches[0][2]), cast.ToInt(matches[0][3]), cast.ToInt(matches[0][1])
		moves = append(moves, Move{from, to, count})
	}

	return crates, moves
}
