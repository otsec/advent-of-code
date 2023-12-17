package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/emirpasic/gods/sets/hashset"
	"math"
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
	field := NewCharField(input, '.')
	ans = math.MaxInt

	startCoord := Coord{0, 0}
	startHistory := NewHistory()
	startVariant := Variant{startCoord, field.At(startCoord), startHistory}

	finishCoord := Coord{field.Width() - 1, field.Height() - 1}

	queue := linkedlistqueue.New()
	queue.Enqueue(startVariant)
	for !queue.Empty() {
		item, _ := queue.Dequeue()
		variant := item.(Variant)

		if reflect.DeepEqual(variant.curr, finishCoord) {
			ans = min(ans, variant.heat)
			continue
		}

		candidates := []Coord{variant.curr.Top(), variant.curr.Right(), variant.curr.Bottom(), variant.curr.Left()}
		for _, c := range candidates {
			if !field.Within(c) {
				continue
			}
			if variant.Was(c) {
				continue
			}

			newHeat := variant.heat + field.At(c)
			newHistory := variant.history.Copy()
			newHistory.Add(c)
			newVariant := Variant{c, newHeat, newHistory}

			queue.Enqueue(newVariant)
		}
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)
	_ = parsed

	return ans
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

func NewCharField(input string, emptyChar byte) *IntField {
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
	curr    Coord
	heat    int
	history *History
}

func (v *Variant) Was(c Coord) bool {
	return v.history.Was(c)
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}
