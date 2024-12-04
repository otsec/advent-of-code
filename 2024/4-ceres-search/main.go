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
	field := NewCharField(input)

	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			start := &Coord{x, y}
			if field.At(start) != 'X' {
				continue
			}
			for d := 0; d <= 7; d++ {
				variant := createLineVariant(start, Direction(d))
				if isVariantCorrect(field, variant) {
					ans += 1
				}
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	field := NewCharField(input)

	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			start := &Coord{x, y}
			if field.At(start) != 'A' {
				continue
			}

			var11 := createLineVariant(start.Top().Left().Top().Left(), DIR_BOTTOM_RIGHT)
			var12 := createLineVariant(start.Bottom().Right().Bottom().Right(), DIR_TOP_LEFT)
			isVar1Correct := isVariantCorrect(field, var11) || isVariantCorrect(field, var12)

			var21 := createLineVariant(start.Top().Right().Top().Right(), DIR_BOTTOM_LEFT)
			var22 := createLineVariant(start.Bottom().Left().Bottom().Left(), DIR_TOP_RIGHT)
			isVar2Correct := isVariantCorrect(field, var21) || isVariantCorrect(field, var22)

			if isVar1Correct && isVar2Correct {
				ans += 1
			}
		}
	}

	return ans
}

type Direction int

const (
	DIR_TOP Direction = iota
	DIR_TOP_RIGHT
	DIR_RIGHT
	DIR_BOTTOM_RIGHT
	DIR_BOTTOM
	DIR_BOTTOM_LEFT
	DIR_LEFT
	DIR_TOP_LEFT
)

func createLineVariant(start *Coord, dir Direction) []*Coord {
	ans := make([]*Coord, 3)

	prev := start
	for i := 0; i < 3; i++ {
		switch dir {
		case DIR_TOP:
			ans[i] = prev.Top()
		case DIR_TOP_RIGHT:
			ans[i] = prev.Top().Right()
		case DIR_RIGHT:
			ans[i] = prev.Right()
		case DIR_BOTTOM_RIGHT:
			ans[i] = prev.Bottom().Right()
		case DIR_BOTTOM:
			ans[i] = prev.Bottom()
		case DIR_BOTTOM_LEFT:
			ans[i] = prev.Bottom().Left()
		case DIR_LEFT:
			ans[i] = prev.Left()
		case DIR_TOP_LEFT:
			ans[i] = prev.Top().Left()
		default:
			panic(fmt.Sprintf("unknown dir: %v", dir))
		}

		prev = ans[i]
	}

	return ans
}

func isVariantCorrect(field *CharField, variant []*Coord) bool {
	for i, coord := range variant {
		if !field.Within(coord) {
			return false
		}

		if i == 0 && field.At(coord) != 'M' {
			return false
		}
		if i == 1 && field.At(coord) != 'A' {
			return false
		}
		if i == 2 && field.At(coord) != 'S' {
			return false
		}
	}

	return true
}

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
