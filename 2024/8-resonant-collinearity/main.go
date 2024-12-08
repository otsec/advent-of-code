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
	nodesMap := createNodesMap(field)

	for _, coords := range nodesMap {
		for i := 0; i < len(coords); i++ {
			for j := i + 1; j < len(coords); j++ {
				placeAntiNodesP1(field, coords[i], coords[j], '#')
			}
		}
	}

	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			if field.At(Coord{x, y}) == '#' {
				ans++
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	field := NewCharField(input)
	nodesMap := createNodesMap(field)

	for _, coords := range nodesMap {
		for i := 0; i < len(coords); i++ {
			for j := i + 1; j < len(coords); j++ {
				placeAntiNodesP2(field, coords[i], coords[j], '#')
			}
		}
	}

	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			if field.At(Coord{x, y}) == '#' {
				ans++
			}
		}
	}

	return ans
}

func createNodesMap(field *CharField) map[byte][]Coord {
	nodesMap := map[byte][]Coord{}
	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			symbol := field.At(Coord{x, y})
			if symbol != '.' {
				if _, ok := nodesMap[symbol]; !ok {
					nodesMap[symbol] = []Coord{{x, y}}
				} else {
					nodesMap[symbol] = append(nodesMap[symbol], Coord{x, y})
				}
			}
		}
	}
	return nodesMap
}

func placeAntiNodesP1(f *CharField, antenna1, antenna2 Coord, ch byte) {
	x1 := antenna1.x - (antenna2.x - antenna1.x)
	x2 := antenna2.x + (antenna2.x - antenna1.x)

	y1 := antenna1.y - (antenna2.y - antenna1.y)
	y2 := antenna2.y + (antenna2.y - antenna1.y)

	antiNode1 := Coord{x1, y1}
	if f.Within(antiNode1) {
		f.Set(antiNode1, ch)
	}

	antiNode2 := Coord{x2, y2}
	if f.Within(antiNode2) {
		f.Set(antiNode2, ch)
	}
}

func placeAntiNodesP2(f *CharField, antenna1, antenna2 Coord, ch byte) {
	dx := antenna2.x - antenna1.x
	dy := antenna2.y - antenna1.y

	for x, y := antenna1.x, antenna1.y; x >= 0 && y >= 0; x, y = x-dx, y-dy {
		antiNode := Coord{x, y}
		if f.Within(antiNode) {
			f.Set(antiNode, ch)
		}
	}

	for x, y := antenna2.x, antenna2.y; x < f.Width() && y < f.Height(); x, y = x+dx, y+dy {
		antiNode := Coord{x, y}
		if f.Within(antiNode) {
			f.Set(antiNode, ch)
		}
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
