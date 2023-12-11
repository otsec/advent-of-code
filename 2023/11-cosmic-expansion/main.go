package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"math"
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
		ans := part1(input, 2)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part1(input, 1000000)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string, emptySpace int) (ans int) {
	field := NewCharField(input, '.')
	_ = field

	var galaxies []Coord
	for y := 0; y < field.Height(); y++ {
		for x := 0; x < field.Width(); x++ {
			c := Coord{x, y}
			if field.At(c) == '#' {
				galaxies = append(galaxies, c)
			}
		}
	}

	var emptyCols []int
	for x := 0; x < field.Width(); x++ {
		found := false
		for _, g := range galaxies {
			if g.x == x {
				found = true
				break
			}
		}

		if !found {
			emptyCols = append(emptyCols, x)
		}
	}

	var emptyRows []int
	for y := 0; y < field.Height(); y++ {
		found := false
		for _, g := range galaxies {
			if g.y == y {
				found = true
				break
			}
		}

		if !found {
			emptyRows = append(emptyRows, y)
		}
	}

	for i := 0; i < len(galaxies)-1; i++ {
		for j := i + 1; j < len(galaxies); j++ {
			c1, c2 := galaxies[i], galaxies[j]
			path := calcGalaxiesPath(c1, c2)

			emptyColsBetween := calcEmptyColsBetween(c1, c2, emptyCols)
			if emptyColsBetween > 0 {
				path += emptyColsBetween*emptySpace - emptyColsBetween
			}

			emptyRowsBetween := calcEmptyRowsBetween(c1, c2, emptyRows)
			if emptyRowsBetween > 0 {
				path += emptyRowsBetween*emptySpace - emptyRowsBetween
			}

			ans += path
		}
	}

	return ans
}

type Coord struct {
	x, y int
}

func (c *Coord) Top() Coord {
	return Coord{c.x, c.y - 1}
}

func (c *Coord) Right() Coord {
	return Coord{c.x + 1, c.y}
}

func (c *Coord) Bottom() Coord {
	return Coord{c.x, c.y + 1}
}

func (c *Coord) Left() Coord {
	return Coord{c.x - 1, c.y}
}

type CharField struct {
	lines     []string
	emptyChar byte
}

func (f *CharField) Width() int {
	return len(f.lines[0])
}

func (f *CharField) Height() int {
	return len(f.lines)
}

func (f *CharField) At(c Coord) byte {
	if c.y < 0 || c.y >= f.Height() {
		return '.'
	}
	if c.x < 0 || c.x >= f.Width() {
		return '.'
	}

	return f.lines[c.y][c.x]
}

func (f *CharField) IsEmpty(c Coord) bool {
	return f.At(c) == '.'
}

func NewCharField(input string, emptyChar byte) *CharField {
	lines := strings.Split(input, "\n")
	return &CharField{lines, emptyChar}
}

func calcGalaxiesPath(c1, c2 Coord) int {
	xPath := math.Abs(float64(c2.x - c1.x))
	yPath := math.Abs(float64(c2.y - c1.y))
	return int(xPath + yPath)
}

func calcEmptyColsBetween(c1, c2 Coord, emptyCols []int) int {
	var ans int

	minX := min(c1.x, c2.x)
	maxX := max(c1.x, c2.x)
	for _, row := range emptyCols {
		if minX < row && row < maxX {
			ans += 1
		}
	}

	return ans
}

func calcEmptyRowsBetween(c1, c2 Coord, emptyRows []int) int {
	var ans int

	minY := min(c1.y, c2.y)
	maxY := max(c1.y, c2.y)
	for _, row := range emptyRows {
		if minY < row && row < maxY {
			ans += 1
		}
	}

	return ans
}
