package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os/exec"
	"slices"
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

func part1(input string) (ans int) {
	field := NewIntField(input)

	memo := Memo{}

	for i := 0; i <= 9; i++ {
		for x := 0; x < field.Width(); x++ {
			for y := 0; y < field.Height(); y++ {
				coord := Coord{x, y}
				if field.At(coord) != i {
					continue
				}

				if i == 0 {
					memo.Set(coord, coord)
				}

				if i > 0 {
					scan := []Coord{coord.Top(), coord.Right(), coord.Bottom(), coord.Left()}
					for _, variant := range scan {
						if field.Within(variant) && field.At(variant) == i-1 {
							memo.Copy(variant, coord)
						}
					}
				}
			}
		}
	}

	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			coord := Coord{x, y}
			if field.At(coord) == 9 {
				ans += memo.Score(coord)
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	field := NewIntField(input)

	memoField := NewIntField(input)
	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			coord := Coord{x, y}
			if field.At(coord) == 0 {
				memoField.Set(coord, 1)
			} else {
				memoField.Set(coord, 0)
			}
		}
	}

	for i := 1; i <= 9; i++ {
		for x := 0; x < field.Width(); x++ {
			for y := 0; y < field.Height(); y++ {
				curr := Coord{x, y}
				if field.At(curr) != i {
					continue
				}

				variants := []Coord{curr.Top(), curr.Right(), curr.Bottom(), curr.Left()}
				for _, prev := range variants {
					if field.Within(prev) && field.At(prev) == i-1 {
						memoField.Set(curr, memoField.At(prev)+memoField.At(curr))
					}
				}
			}
		}
	}

	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			coord := Coord{x, y}
			if field.At(coord) == 9 {
				ans += memoField.At(coord)
			}
		}
	}

	return ans
}

type Memo map[Coord][]Coord

func (m *Memo) Set(at, start Coord) {
	if saved, ok := (*m)[at]; ok {
		saved = append(saved, start)
	} else {
		(*m)[at] = []Coord{start}
	}
}

func (m *Memo) Copy(from, to Coord) {
	if _, ok := (*m)[to]; ok {
		for _, v := range (*m)[from] {
			if !slices.Contains((*m)[to], v) {
				(*m)[to] = append((*m)[to], v)
			}
		}
	} else {
		(*m)[to] = []Coord{}
		(*m)[to] = append((*m)[to], (*m)[from]...)
	}
}

func (m *Memo) Score(coord Coord) int {
	if saved, ok := (*m)[coord]; ok {
		return len(saved)
	} else {
		return 0
	}
}

func (m *Memo) Print(width, height int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fmt.Print(m.Score(Coord{x, y}))
		}
		fmt.Println()
	}
}

type Direction int

const (
	Top Direction = iota
	Right
	Bottom
	Left
)

type Coord struct {
	x, y int
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
		panic(fmt.Sprintf("unknown dir %v", dir))
	}
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

type IntField struct {
	lines [][]int
}

func NewIntField(input string) *IntField {
	stringLines := strings.Split(input, "\n")
	numLines := make([][]int, len(stringLines))
	for i, line := range stringLines {
		numLines[i] = make([]int, len(line))
		for j, num := range line {
			atoi, _ := strconv.Atoi(string(num))
			numLines[i][j] = atoi
		}
	}
	return &IntField{numLines}
}

func (f *IntField) Width() int {
	return len(f.lines[0])
}

func (f *IntField) Height() int {
	return len(f.lines)
}

func (f *IntField) Within(c Coord) bool {
	if c.y < 0 || c.y >= f.Height() {
		return false
	}
	if c.x < 0 || c.x >= f.Width() {
		return false
	}
	return true
}

func (f *IntField) At(c Coord) int {
	return f.lines[c.y][c.x]
}

func (f *IntField) Set(c Coord, v int) {
	f.lines[c.y][c.x] = v
}

func (f *IntField) IsEmpty(c Coord) bool {
	return f.At(c) == '.'
}

func (f *IntField) ToString() string {
	stringLines := make([]string, len(f.lines))
	for i, numsLine := range f.lines {
		bytesLine := make([]byte, len(numsLine))
		for j, num := range numsLine {
			bytesLine[j] = strconv.Itoa(num)[0]
		}

		stringLines[i] = string(bytesLine)
	}

	return strings.Join(stringLines, "\n")
}

func (f *IntField) Print() {
	fmt.Println(f.ToString())
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
