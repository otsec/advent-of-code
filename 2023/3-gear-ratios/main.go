package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
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
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) (ans int) {
	field := parseInput(input)

	for y := 0; y < len(field.lines); y++ {
		for x := 0; x < len(field.lines[0]); x++ {
			c := &Coord{x, y}
			if field.IsNumber(c) {
				snum, inum := field.CutNumber(c)
				x += len(snum)
				// fmt.Println(c, snum, inum, field.IsNumberAdjacentToSymbol(c, len(snum)))

				if field.IsNumberAdjacentToSymbol(c, len(snum)) {
					ans += inum
				}
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	field := parseInput(input)

	for y := 0; y < len(field.lines); y++ {
		for x := 0; x < len(field.lines[0]); x++ {
			c := &Coord{x, y}
			if field.At(c) != '*' {
				continue
			}

			numCoords := field.FindAdjacentNumbers(c)
			if len(numCoords) == 2 {
				_, n1 := field.CutNumber(numCoords[0])
				_, n2 := field.CutNumber(numCoords[1])
				ans += n1 * n2
			}
		}
	}

	return ans
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

type Field struct {
	lines []string
}

func (f *Field) At(c *Coord) byte {
	if c.y < 0 || c.y >= len(f.lines) {
		return '.'
	}
	if c.x < 0 || c.x >= len(f.lines[0]) {
		return '.'
	}

	return f.lines[c.y][c.x]
}

func (f *Field) IsEmpty(c *Coord) bool {
	return f.At(c) == '.'
}

func (f *Field) IsNumber(c *Coord) bool {
	return f.At(c) >= '0' && f.At(c) <= '9'
}

func (f *Field) isSymbol(c *Coord) bool {
	return !f.IsNumber(c) && !f.IsEmpty(c)
}

func (f *Field) IsNumberAdjacentToSymbol(start *Coord, length int) bool {
	if f.isSymbol(start.Left().Top()) || f.isSymbol(start.Left()) || f.isSymbol(start.Left().Bottom()) {
		return true
	}

	c := start
	for f.IsNumber(c) {
		if f.isSymbol(c.Top()) || f.isSymbol(c.Bottom()) {
			return true
		}
		c = c.Right()
	}

	end := &Coord{x: start.x + length - 1, y: start.y}
	if f.isSymbol(end.Right().Top()) || f.isSymbol(end.Right()) || f.isSymbol(end.Right().Bottom()) {
		return true
	}

	return false
}

func (f *Field) FindAdjacentNumbers(c *Coord) []*Coord {
	var ans []*Coord

	if f.IsNumber(c.Top()) {
		ans = append(ans, c.Top())
	} else {
		if f.IsNumber(c.Top().Left()) {
			ans = append(ans, c.Top().Left())
		}
		if f.IsNumber(c.Top().Right()) {
			ans = append(ans, c.Top().Right())
		}
	}

	if f.IsNumber(c.Right()) {
		ans = append(ans, c.Right())
	}

	if f.IsNumber(c.Bottom()) {
		ans = append(ans, c.Bottom())
	} else {
		if f.IsNumber(c.Bottom().Left()) {
			ans = append(ans, c.Bottom().Left())
		}
		if f.IsNumber(c.Bottom().Right()) {
			ans = append(ans, c.Bottom().Right())
		}
	}

	if f.IsNumber(c.Left()) {
		ans = append(ans, c.Left())
	}

	return ans
}

func (f *Field) CutNumber(c *Coord) (string, int) {
	for f.IsNumber(c.Left()) {
		c = c.Left()
	}

	var digits []byte
	for f.IsNumber(c) {
		digits = append(digits, f.At(c))
		c = c.Right()
	}

	s := string(digits)
	i, _ := strconv.Atoi(s)

	return s, i
}

func parseInput(input string) *Field {
	return &Field{
		lines: strings.Split(input, "\n"),
	}
}
