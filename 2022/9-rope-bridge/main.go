package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"github.com/alexchao26/advent-of-code-go/mathy"
	"github.com/alexchao26/advent-of-code-go/util"
	"github.com/emirpasic/gods/sets/hashset"
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
	moves := parseInput(input)

	head := Position{0, 0}
	tail := Position{0, 0}
	history := hashset.New()

	for _, move := range moves {
		for move.steps > 0 {
			// fmt.Print("head: ", head, " -> ", move, " -> ")

			moveInDirection(&head, move.direction)
			followMove(&tail, &head)
			history.Add(tail)

			move.steps--

			// fmt.Print(head, " tail: ", tail, "\n")
		}
	}

	return history.Size()
}

func part2(input string) (ans int) {
	moves := parseInput(input)

	knots := make([]Position, 10)

	history := hashset.New()

	for _, move := range moves {
		//fmt.Println(move)

		for move.steps > 0 {
			moveInDirection(&knots[0], move.direction)
			for i := 1; i < len(knots); i++ {
				followMove(&knots[i], &knots[i-1])
			}

			history.Add(knots[9])

			move.steps--
		}

		//drawKnots(knots, []int{-11, -5}, []int{27, 20})
		//fmt.Println()
	}

	return history.Size()
}

type Move struct {
	direction string
	steps     int
}

func parseInput(input string) []Move {
	moves := []Move{}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		segments := strings.Split(line, " ")
		move := Move{segments[0], cast.ToInt(segments[1])}
		moves = append(moves, move)
	}

	return moves
}

type Position struct {
	x int
	y int
}

func moveInDirection(pos *Position, direction string) {
	switch direction {
	case "U":
		pos.y++
	case "R":
		pos.x++
	case "D":
		pos.y--
	case "L":
		pos.x--
	}
}

func followMove(tail *Position, head *Position) {
	if head.x-1 <= tail.x && tail.x <= head.x+1 {
		if head.y-1 <= tail.y && tail.y <= head.y+1 {
			return
		}
	}

	if tail.x != head.x && tail.y != head.y {
		diffX := mathy.AbsInt(tail.x - head.x)
		diffY := mathy.AbsInt(tail.y - head.y)

		if diffX == diffY {
			tail.x = (tail.x + head.x) / 2
			tail.y = (tail.y + head.y) / 2
		}

		if diffX > diffY {
			tail.y = head.y
		} else if diffY > diffX {
			tail.x = head.x
		}

		// if diffX != diffY, calculations are not finished
		// we need one of the next functions to calc other coordinate
	}

	if tail.x == head.x && tail.y != head.y {
		if tail.y > head.y {
			tail.y--
		} else {
			tail.y++
		}
	}

	if tail.y == head.y && tail.x != head.x {
		if tail.x > head.x {
			tail.x--
		} else {
			tail.x++
		}
	}
}

func drawKnots(knots []Position, start []int, size []int) {
	for y := size[1] - start[1] - 1; y >= start[1]; y-- {
		for x := start[0]; x < size[0]-start[0]; x++ {
			drawed := false

			for index, knot := range knots {
				if knot.x == x && knot.y == y {
					if index == 0 {
						fmt.Print("H")
					} else {
						fmt.Print(index)
					}

					drawed = true
					break
				}
			}

			if !drawed {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
