package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
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
	heightGrid, stepsGrid := parseInput(input)

	start := findSymbol(heightGrid, 'S')
	end := findSymbol(heightGrid, 'E')

	lookup(&heightGrid, &stepsGrid, start, 1)

	// printSteps(stepsGrid)

	return stepsGrid[end.y][end.x] - 1
}

func part2(input string) (ans int) {
	heightGrid, stepsGrid := parseInput(input)

	// hack: use only first vertical line as a start points
	// - I saw that my input has 'a' only there
	// - Example has few 'a' on second line, but we know that best start is still on the first line
	for y := 0; y < len(heightGrid); y++ {
		start := Coordinate{0, y}
		lookup(&heightGrid, &stepsGrid, start, 1)
	}

	end := findSymbol(heightGrid, 'E')
	return stepsGrid[end.y][end.x] - 1
}

func findSymbol(heightGrid [][]rune, symbol rune) Coordinate {
	for y, line := range heightGrid {
		for x, char := range line {
			if char == symbol {
				return Coordinate{x, y}
			}
		}
	}

	panic(fmt.Sprintf("Rune %v was not found on the grid.", symbol))
}

func lookup(heightGrid *[][]rune, stepsGrid *[][]int, cell Coordinate, step int) {
	val := (*stepsGrid)[cell.y][cell.x]
	if val == 0 || val > step {
		(*stepsGrid)[cell.y][cell.x] = step
	} else {
		return
	}

	var next Coordinate

	next = Coordinate{cell.x, cell.y - 1}
	if canStep(heightGrid, cell, next) {
		lookup(heightGrid, stepsGrid, next, step+1)
	}

	next = Coordinate{cell.x + 1, cell.y}
	if canStep(heightGrid, cell, next) {
		lookup(heightGrid, stepsGrid, next, step+1)
	}

	next = Coordinate{cell.x, cell.y + 1}
	if canStep(heightGrid, cell, next) {
		lookup(heightGrid, stepsGrid, next, step+1)
	}

	next = Coordinate{cell.x - 1, cell.y}
	if canStep(heightGrid, cell, next) {
		lookup(heightGrid, stepsGrid, next, step+1)
	}
}

func canStep(heightGrid *[][]rune, from Coordinate, to Coordinate) bool {
	if to.x < 0 || to.y < 0 {
		return false
	}

	maxX := len((*heightGrid)[0]) - 1
	maxY := len(*heightGrid) - 1
	if to.x > maxX || to.y > maxY {
		return false
	}

	valFrom := (*heightGrid)[from.y][from.x]
	valTo := (*heightGrid)[to.y][to.x]

	if valFrom == 'S' {
		return valTo == 'a'
	}
	if valTo == 'E' {
		return valFrom == 'z'
	}

	return int(valTo)-int(valFrom) <= 1
}

func parseInput(input string) ([][]rune, [][]int) {
	lines := strings.Split(input, "\n")

	heightGrid := make([][]rune, len(lines))
	stepsGrid := make([][]int, len(lines))

	for i, line := range lines {
		heightGrid[i] = []rune(line)
		stepsGrid[i] = make([]int, len(heightGrid[i]))
	}

	return heightGrid, stepsGrid
}

type Coordinate struct {
	x int
	y int
}

func printSteps(grid *[][]int) {
	for _, line := range *grid {
		for _, step := range line {
			fmt.Printf("%3d", step)
		}
		fmt.Println()
	}
}
