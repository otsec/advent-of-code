package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"slices"
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
	field := NewCharField(input)

	currPos := findStart(field)
	currDir := DIR_TOP

	for {
		field.Set(currPos, 'X')
		nextPos := takeAStep(currPos, currDir)
		if !field.Within(nextPos) {
			break
		} else if field.At(nextPos) == '#' {
			currDir = makeATurn(currDir)
			continue
		} else {
			currPos = nextPos
		}
	}

	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			if field.At(&Coord{x, y}) == 'X' {
				ans += 1
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	field := NewCharField(input)
	start := findStart(field)

	var variants []*Coord
	{
		currPos := &Coord{x: start.x, y: start.y}
		currDir := DIR_TOP
		for {
			nextPos := takeAStep(currPos, currDir)
			if !field.Within(nextPos) {
				break
			} else if field.At(nextPos) == '#' {
				currDir = makeATurn(currDir)
				continue
			} else {
				variants = append(variants, nextPos)
				currPos = nextPos
			}
		}
	}

	var obstacles []Coord
	for _, newObstacle := range variants {
		if slices.Contains(obstacles, *newObstacle) {
			continue
		} else {
			obstacles = append(obstacles, *newObstacle)
		}

		field = NewCharField(input)
		field.Set(newObstacle, '#')

		history := History{}
		currPos := &Coord{x: start.x, y: start.y}
		currDir := DIR_TOP
		for field.Within(currPos) {
			nextPos := takeAStep(currPos, currDir)
			if history.Contains(nextPos, currDir) {
				ans += 1
				break
			} else if !field.Within(nextPos) {
				break
			} else if field.At(nextPos) == '#' {
				currDir = makeATurn(currDir)
				continue
			} else {
				history.Add(nextPos, currDir)
				currPos = nextPos
			}
		}
	}

	return ans
}

func findStart(f *CharField) *Coord {
	for x := 0; x < f.Width(); x++ {
		for y := 0; y < f.Height(); y++ {
			c := &Coord{x, y}
			if f.At(c) == '^' {
				return c
			}
		}
	}

	panic("start position not found")
}

func takeAStep(c *Coord, dir Direction) *Coord {
	switch dir {
	case DIR_TOP:
		return c.Top()
	case DIR_RIGHT:
		return c.Right()
	case DIR_BOTTOM:
		return c.Bottom()
	case DIR_LEFT:
		return c.Left()
	}

	panic("unknown dir")
}

func makeATurn(dir Direction) Direction {
	switch dir {
	case DIR_TOP:
		return DIR_RIGHT
	case DIR_RIGHT:
		return DIR_BOTTOM
	case DIR_BOTTOM:
		return DIR_LEFT
	case DIR_LEFT:
		return DIR_TOP
	}

	panic("unknown dir")
}

type HistoryItem struct {
	x, y, dir int
}

type History struct {
	items []HistoryItem
}

func (h *History) Add(c *Coord, dir Direction) {
	h.items = append(h.items, HistoryItem{c.x, c.y, int(dir)})
}

func (h *History) Contains(c *Coord, dir Direction) bool {
	item := HistoryItem{c.x, c.y, int(dir)}
	return slices.Contains(h.items, item)
}

type Direction int

const (
	DIR_TOP Direction = iota
	DIR_RIGHT
	DIR_BOTTOM
	DIR_LEFT
)

type Coord struct {
	x, y int
}

func (c *Coord) Top() *Coord {
	return &Coord{c.x, c.y - 1}
}

func (c *Coord) Right() *Coord {
	return &Coord{c.x + 1, c.y}
}

func (c *Coord) Bottom() *Coord {
	return &Coord{c.x, c.y + 1}
}

func (c *Coord) Left() *Coord {
	return &Coord{c.x - 1, c.y}
}

type CharField struct {
	lines []string
}

func (f *CharField) Width() int {
	return len(f.lines[0])
}

func (f *CharField) Height() int {
	return len(f.lines)
}

func (f *CharField) Within(c *Coord) bool {
	if c.y < 0 || c.y >= f.Height() {
		return false
	}
	if c.x < 0 || c.x >= f.Width() {
		return false
	}
	return true
}

func (f *CharField) At(c *Coord) byte {
	return f.lines[c.y][c.x]
}

func (f *CharField) Set(c *Coord, v byte) {
	line := []byte(f.lines[c.y])
	line[c.x] = v
	f.lines[c.y] = string(line)
}

func (f *CharField) ToString() string {
	return strings.Join(f.lines, "\n")
}

func (f *CharField) Print() {
	fmt.Println(f.ToString())
}

func NewCharField(input string) *CharField {
	lines := strings.Split(input, "\n")
	return &CharField{lines}
}
