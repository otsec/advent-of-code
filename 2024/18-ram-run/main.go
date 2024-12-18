package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os/exec"
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
		CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) int {
	corrupted := parseInput(input)

	maze := NewEmptyCharField(71, 71, '.')
	if len(corrupted) < 50 {
		maze = NewEmptyCharField(7, 7, '.')
	}

	secondsPassed := 1024
	if len(corrupted) < 50 {
		secondsPassed = 12
	}
	fillCorruption(maze, corrupted[:secondsPassed])

	return findExit(maze)
}

func part2(input string) string {
	corrupted := parseInput(input)

	for i := 1; i < len(corrupted); i++ {
		maze := NewEmptyCharField(71, 71, '.')
		if len(corrupted) < 50 {
			maze = NewEmptyCharField(7, 7, '.')
		}
		fillCorruption(maze, corrupted[:i])

		exit := findExit(maze)
		if exit == -1 {
			coord := corrupted[i-1]
			return fmt.Sprintf("%v,%v", coord.x, coord.y)
		}
	}

	return ""
}

type Move struct {
	coord  Coord
	second int
}

func findExit(maze *CharField) int {
	start := Coord{0, 0}
	exit := Coord{maze.Width() - 1, maze.Height() - 1}

	memo := map[Coord]bool{start: true}
	queue := []Move{{start, 0}}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		variants := []Coord{curr.coord.Top(), curr.coord.Right(), curr.coord.Bottom(), curr.coord.Left()}
		for _, variant := range variants {
			if variant == exit {
				return curr.second + 1
			}

			if memo[variant] {
				continue
			}
			if !maze.Within(variant) {
				continue
			}
			if maze.At(variant) == '#' {
				continue
			}

			queue = append(queue, Move{variant, curr.second + 1})
			memo[variant] = true
		}
	}

	return -1
}

func fillCorruption(field *CharField, bytes []Coord) {
	for _, coord := range bytes {
		field.Set(coord, '#')
	}
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

func NewEmptyCharField(width, height int, emptySymbol byte) *CharField {
	lines := make([]string, height)
	for i, _ := range lines {
		lines[i] = string(bytes.Repeat([]byte{emptySymbol}, width))
	}

	return &CharField{lines}
}

func parseInput(input string) []Coord {
	lines := strings.Split(input, "\n")
	ans := make([]Coord, len(lines))
	for i, line := range lines {
		nums := parseNumbers(line)
		ans[i] = Coord{nums[0], nums[1]}
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
