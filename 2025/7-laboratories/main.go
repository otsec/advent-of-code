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

	for y := 1; y < field.Height(); y++ {
		for x := 0; x < field.Width(); x++ {
			curr := Coord{x, y}

			prevVal := field.At(curr.Top())
			if prevVal != '|' && prevVal != 'S' {
				continue
			}

			switch field.At(curr) {
			case '.':
				field.Set(curr, '|')
			case '^':
				field.Set(curr.Left(), '|')
				field.Set(curr.Right(), '|')
				ans++
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	field := NewCharField(input)

	for y := 1; y < field.Height(); y++ {
		for x := 0; x < field.Width(); x++ {
			curr := Coord{x, y}

			prevVal := field.At(curr.Top())
			if prevVal != '|' && prevVal != 'S' {
				continue
			}

			switch field.At(curr) {
			case '.':
				field.Set(curr, '|')
			case '^':
				field.Set(curr.Left(), '|')
				field.Set(curr.Right(), '|')
			}
		}
	}

	memo := make(map[Coord]int)

	for x := 0; x < field.Width(); x++ {
		curr := Coord{x, field.Height() - 1}
		if field.At(curr) == '|' {
			memo[curr] = 1
		}
	}

	for y := field.Height() - 2; y >= 0; y-- {
		for x := 0; x < field.Width(); x++ {
			curr := Coord{x, y}
			if field.At(curr) == '|' {
				memo[curr] = memo[curr.Bottom()]
			}
		}

		for x := 0; x < field.Width(); x++ {
			curr := Coord{x, y}
			if field.At(curr) == '^' {
				memo[curr] = memo[curr.Left()] + memo[curr.Right()]
			}
			if field.At(curr) == 'S' {
				return memo[curr.Bottom()]
			}
		}
	}

	return ans
}

type Queue []Coord

func (q *Queue) Len() int {
	return len(*q)
}

func (q *Queue) Push(c Coord) {
	*q = append(*q, c)
}

func (q *Queue) Pop() Coord {
	old := *q
	*q = old[1:]
	return old[0]
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

func NewCharField(input string) *CharField {
	lines := strings.Split(input, "\n")
	return &CharField{lines}
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
