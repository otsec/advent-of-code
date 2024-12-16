package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os/exec"
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
	field := NewCharField(input)

	startPos := Position{findSymbol(field, 'S'), Right, 0}

	queue := []Position{startPos}
	memo := make(map[Coord]int)
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]

		variants := []Position{
			{pos.coord.Next(pos.dir), pos.dir, pos.score + 1},
			{pos.coord.Next(pos.dir.TurnClockwise()), pos.dir.TurnClockwise(), pos.score + 1001},
			{pos.coord.Next(pos.dir.TurnCounterClockwise()), pos.dir.TurnCounterClockwise(), pos.score + 1001},
		}
		for _, variant := range variants {
			switch field.At(variant.coord) {
			case 'E':
				if ans == 0 || ans > variant.score {
					ans = variant.score
				}
			case '.':
				if memo[variant.coord] == 0 || memo[variant.coord] > variant.score {
					memo[variant.coord] = variant.score
					queue = append(queue, variant)
				}
			}
		}
	}

	return ans
}

func part2(input string) int {
	field := NewCharField(input)

	start := findSymbol(field, 'S')
	startPos := Position2{start, Right, 0, []Coord{start}}

	queue := []Position2{startPos}
	memo := make(map[Position]int)
	bestScore := 0
	winners := []Position2{}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]

		variants := []Position{
			{pos.coord.Next(pos.dir), pos.dir, pos.score + 1},
			{pos.coord.Next(pos.dir.TurnClockwise()), pos.dir.TurnClockwise(), pos.score + 1001},
			{pos.coord.Next(pos.dir.TurnCounterClockwise()), pos.dir.TurnCounterClockwise(), pos.score + 1001},
		}
		for _, variant := range variants {
			switch field.At(variant.coord) {
			case 'E':
				if bestScore == 0 || bestScore >= variant.score {
					bestScore = variant.score

					history := make([]Coord, len(pos.history)+1)
					copy(history, pos.history)
					history[len(history)-1] = variant.coord
					pos2 := Position2{variant.coord, variant.dir, variant.score, history}

					winners = append(winners, pos2)
				}
			case '.':
				memoPos := Position{variant.coord, pos.dir, 0}
				if memo[memoPos] == 0 || memo[memoPos] >= variant.score {
					memo[memoPos] = variant.score

					history := make([]Coord, len(pos.history)+1)
					copy(history, pos.history)
					history[len(history)-1] = variant.coord
					pos2 := Position2{variant.coord, variant.dir, variant.score, history}

					queue = append(queue, pos2)
				}
			}
		}
	}

	points := map[Coord]bool{}
	for _, pos := range winners {
		if pos.score == bestScore {
			for _, coord := range pos.history {
				points[coord] = true
			}
		}
	}

	return len(points)
}

func findSymbol(field *CharField, s byte) Coord {
	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			coord := Coord{x, y}
			if field.At(coord) == s {
				return coord
			}
		}
	}

	panic("symbol not found")
}

type Coord struct {
	x, y int
}

func (c Coord) Top() Coord {
	return Coord{c.x, c.y - 1}
}

func (c Coord) Right() Coord {
	return Coord{c.x + 1, c.y}
}

func (c Coord) Bottom() Coord {
	return Coord{c.x, c.y + 1}
}

func (c Coord) Left() Coord {
	return Coord{c.x - 1, c.y}
}

type Direction int

const (
	Top Direction = iota
	Right
	Bottom
	Left
)

func (d Direction) TurnClockwise() Direction {
	if d == 3 {
		return 0
	} else {
		return d + 1
	}
}

func (d Direction) TurnCounterClockwise() Direction {
	if d == 0 {
		return 3
	} else {
		return d - 1
	}
}

func (c Coord) Next(dir Direction) Coord {
	switch dir {
	case Top:
		return Coord{c.x, c.y - 1}
	case Right:
		return Coord{c.x + 1, c.y}
	case Bottom:
		return Coord{c.x, c.y + 1}
	case Left:
		return Coord{c.x - 1, c.y}
	default:
		panic("unknown direction type")
	}
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

func (f *CharField) Within(c Coord) bool {
	if c.y < 0 || c.y >= f.Height() {
		return false
	}
	if c.x < 0 || c.x >= f.Width() {
		return false
	}
	return true
}

func (f *CharField) At(c Coord) byte {
	return f.lines[c.y][c.x]
}

func (f *CharField) Set(c Coord, v byte) {
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

type Position struct {
	coord Coord
	dir   Direction
	score int
}

type Position2 struct {
	coord   Coord
	dir     Direction
	score   int
	history []Coord
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
