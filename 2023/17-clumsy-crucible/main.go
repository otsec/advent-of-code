package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/sets/hashset"
	"reflect"
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
	return findPath(input, 0, 3)
}

func part2(input string) (ans int) {
	return findPath(input, 4, 10)
}

func findPath(input string, minStraight, maxStraight int) int {
	field := NewIntField(input, '.')

	startCoord := Coord{0, 0}
	finishCoord := Coord{field.Width() - 1, field.Height() - 1}

	queue := priorityqueue.NewWith(func(a, b interface{}) int {
		varA := a.(Variant)
		varB := b.(Variant)
		return varA.heat - varB.heat
	})

	next := []Direction{Right, Bottom}
	for _, dir := range next {
		coord := startCoord.Next(dir)
		heat := field.At(coord)
		queue.Enqueue(Variant{coord, dir, 1, heat})
	}

	// key is Variant without hear
	// value is heat
	// for type safety we will recreate Variant with 0 heat
	memory := map[Variant]int{}

	for !queue.Empty() {
		item, _ := queue.Dequeue()
		variant := item.(Variant)

		memKey := Variant{variant.coord, variant.dir, variant.dirInRow, 0}
		if prevHeat, exists := memory[memKey]; exists {
			if prevHeat <= variant.heat {
				continue
			}
		}
		memory[memKey] = variant.heat

		if reflect.DeepEqual(variant.coord, finishCoord) {
			return variant.heat
		}

		candidates := []Direction{Top, Right, Bottom, Left}
		for _, dir := range candidates {
			if isOpposite(dir, variant.dir) {
				continue
			}

			coord := variant.coord.Next(dir)
			if !field.Within(coord) {
				continue
			}

			if dir != variant.dir && variant.dirInRow < minStraight {
				continue
			}

			dirInRow := 1
			if dir == variant.dir {
				dirInRow = variant.dirInRow + 1
			}
			if dirInRow > maxStraight {
				continue
			}

			heat := variant.heat + field.At(coord)
			queue.Enqueue(Variant{coord, dir, dirInRow, heat})
		}
	}

	panic("route not found")
}

type Direction int

const (
	Top Direction = iota
	Right
	Bottom
	Left
)

func isOpposite(d1, d2 Direction) bool {
	if min(d1, d2) == Top && max(d1, d2) == Bottom {
		return true
	}
	if min(d1, d2) == Right && max(d1, d2) == Left {
		return true
	}
	return false
}

type Coord struct {
	x, y int
}

func (c *Coord) Next(dir Direction) Coord {
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

type IntField struct {
	lines [][]int
}

func NewIntField(input string, emptyChar byte) *IntField {
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

type History struct {
	set *hashset.Set
}

func NewHistory() *History {
	return &History{hashset.New()}
}

func (h *History) Add(c Coord) {
	h.set.Add(c)
}

func (h *History) Was(c Coord) bool {
	return h.set.Contains(c)
}

func (h *History) Copy() *History {
	newHistory := NewHistory()
	newHistory.set.Add(h.set.Values()...)
	return newHistory
}

type Variant struct {
	coord    Coord
	dir      Direction
	dirInRow int
	heat     int
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}
