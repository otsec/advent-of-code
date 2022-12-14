package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"github.com/alexchao26/advent-of-code-go/util"
	"log"
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
	lines := parseInput(input)

	leftTop, rightBottom := findCaveEdges(lines)
	leftTop.x -= 1
	rightBottom.x += 1
	rightBottom.y += 1
	cave := createCave(leftTop, rightBottom)
	cave.drawAt(Coordinate{500, 0}, "+")

	for _, line := range lines {
		cave.drawLine(line, "#")
	}

	for {
		c := dropAGrain(&cave)
		if c.y == cave.rightBottom.y {
			break
		}

		cave.drawAt(c, "o")
		ans++
	}

	//drawCave(&cave)

	return ans
}

func part2(input string) (ans int) {
	lines := parseInput(input)

	leftTop, rightBottom := findCaveEdges(lines)
	rightBottom.y += 2
	leftTop.x = 500 - rightBottom.y
	rightBottom.x = 500 + rightBottom.y
	cave := createCave(leftTop, rightBottom)
	cave.drawAt(Coordinate{500, 0}, "+")

	floor := Line{Coordinate{leftTop.x, rightBottom.y}, rightBottom}
	cave.drawLine(floor, "#")

	for _, line := range lines {
		cave.drawLine(line, "#")
	}

	for {
		c := dropAGrain(&cave)
		cave.drawAt(c, "o")
		ans++

		if c.x == 500 && c.y == 0 {
			break
		}
	}

	//drawCave(&cave)

	return ans
}

type Coordinate struct {
	x int
	y int
}

type Line struct {
	start Coordinate
	end   Coordinate
}

type Cave struct {
	leftTop     Coordinate
	rightBottom Coordinate
	cells       [][]string
}

func (c *Cave) at(x, y int) string {
	return c.cells[y][x-c.leftTop.x]
}

func (c *Cave) drawAt(crd Coordinate, symbol string) {
	c.cells[crd.y][crd.x-c.leftTop.x] = symbol
}

func (c *Cave) drawLine(line Line, symbol string) {
	if line.start.x == line.end.x {
		x := line.start.x
		minY := min(line.start.y, line.end.y)
		maxY := max(line.start.y, line.end.y)
		for y := minY; y <= maxY; y++ {
			c.drawAt(Coordinate{x, y}, symbol)
		}
	} else if line.start.y == line.end.y {
		y := line.start.y
		minX := min(line.start.x, line.end.x)
		maxX := max(line.start.x, line.end.x)
		for x := minX; x <= maxX; x++ {
			c.drawAt(Coordinate{x, y}, symbol)
		}
	} else {
		log.Fatalf("%v and %v are not on the same line.", line.start, line.end)
	}
}

func parseInput(input string) []Line {
	lines := []Line{}

	for _, line := range strings.Split(input, "\n") {
		var nextStart Coordinate
		for j, pair := range strings.Split(line, " -> ") {
			xy := strings.Split(pair, ",")
			thisEnd := Coordinate{cast.ToInt(xy[0]), cast.ToInt(xy[1])}
			if j != 0 {
				lines = append(lines, Line{nextStart, thisEnd})
			}
			nextStart = thisEnd
		}
	}

	return lines
}

func findCaveEdges(lines []Line) (Coordinate, Coordinate) {
	leftTop := Coordinate{lines[0].start.x, 0}
	rightBottom := lines[0].start

	for _, line := range lines {
		for _, c := range []Coordinate{line.start, line.end} {
			leftTop.x = min(leftTop.x, c.x)
			leftTop.y = min(leftTop.y, c.y)
			rightBottom.x = max(rightBottom.x, c.x)
			rightBottom.y = max(rightBottom.y, c.y)
		}
	}

	return leftTop, rightBottom
}

func createCave(leftTop, rightBottom Coordinate) Cave {
	height := rightBottom.y - leftTop.y + 1
	width := rightBottom.x - leftTop.x + 1

	cells := make([][]string, height)
	for y, _ := range cells {
		cells[y] = make([]string, width)
		for x, _ := range cells[y] {
			cells[y][x] = "."
		}
	}

	return Cave{leftTop, rightBottom, cells}
}

func dropAGrain(cave *Cave) Coordinate {
	now := Coordinate{500, 0}

	for {
		if now.y == cave.rightBottom.y {
			break
		}

		next := nextGrainCoordinate(cave, now)
		if next.x == now.x && next.y == now.y {
			break
		}

		now = next
	}

	return now
}

func nextGrainCoordinate(cave *Cave, c Coordinate) Coordinate {
	if cave.at(c.x, c.y+1) == "." {
		return Coordinate{c.x, c.y + 1}
	}
	if cave.at(c.x-1, c.y+1) == "." {
		return Coordinate{c.x - 1, c.y + 1}
	}
	if cave.at(c.x+1, c.y+1) == "." {
		return Coordinate{c.x + 1, c.y + 1}
	}

	return c
}

func drawCave(cave *Cave) {
	for y, row := range cave.cells {
		fmt.Printf("%3d %s \n", y, strings.Join(row, ""))
	}
	fmt.Println()
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
