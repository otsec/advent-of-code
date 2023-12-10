package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	"slices"
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
	start := findStart(field)

	numbers := NewNumbersField(field.Width(), field.Height(), -1)
	numbers.Set(start, 0)

	queue := []Coord{start}
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		for _, coord := range findConnected(field, item) {
			if numbers.At(coord) == -1 {
				newVal := numbers.At(item) + 1
				numbers.Set(coord, newVal)
				ans = max(ans, newVal)
				queue = append(queue, coord)
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	field := parseInput(input)
	start := findStart(field)

	numbers := NewNumbersField(field.Width(), field.Height(), 0)
	numbers.Set(start, 1)

	// in part we use 0 as empty sybmol, not -1
	// and we fill loop borders with 1

	queue := llq.New()
	queue.Enqueue(start)
	for !queue.Empty() {
		item, _ := queue.Dequeue()

		for _, coord := range findConnected(field, item.(Coord)) {
			if numbers.At(coord) == 0 {
				numbers.Set(coord, 1)
				queue.Enqueue(coord)
			}
		}
	}

	// create new field with 1 extra tile around

	canvas := NewNumbersField(field.Width()+2, field.Height()+2, 0)
	for y := 0; y < numbers.Height(); y++ {
		for x := 0; x < numbers.Width(); x++ {
			canvas.Set(Coord{x + 1, y + 1}, numbers.At(Coord{x, y}))
		}
	}

	// fill everything outside the loop with -1

	queue = llq.New()
	queue.Enqueue(Coord{0, 0})
	for !queue.Empty() {
		item, _ := queue.Dequeue()
		coord := item.(Coord)

		if canvas.At(coord) == 0 {
			canvas.Set(coord, 2)
		} else {
			continue
		}

		next := []Coord{coord.Top(), coord.Right(), coord.Bottom(), coord.Left()}
		for _, n := range next {
			if canvas.IsInsideBorders(n) {
				queue.Enqueue(n)
			}
		}
	}

	for y := 0; y < canvas.Height(); y++ {
		for x := 0; x < canvas.Width(); x++ {
			if canvas.At(Coord{x, y}) == 0 {
				ans += 1
			}
		}
	}

	//for _, line := range canvas.arr {
	//	fmt.Println(line)
	//}
	//fmt.Println()

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

type NumbersField struct {
	arr [][]int
}

func (f *NumbersField) Width() int {
	return len(f.arr[0])
}

func (f *NumbersField) Height() int {
	return len(f.arr)
}

func (f *NumbersField) IsInsideBorders(c Coord) bool {
	return c.x >= 0 && c.x < f.Width() && c.y >= 0 && c.y < f.Height()
}

func (f *NumbersField) At(c Coord) int {
	return f.arr[c.y][c.x]
}

func (f *NumbersField) Set(c Coord, v int) {
	f.arr[c.y][c.x] = v
}

func NewNumbersField(width, height, emptyNum int) *NumbersField {
	numbers := make([][]int, height)
	for y, _ := range numbers {
		numbers[y] = make([]int, width)
		for x, _ := range numbers[y] {
			numbers[y][x] = emptyNum
		}
	}

	return &NumbersField{numbers}
}

func parseInput(input string) *CharField {
	input = strings.ReplaceAll(input, "\t", "")
	return NewCharField(input, '.')
}

func findStart(f *CharField) Coord {
	for y := 0; y < f.Height(); y++ {
		for x := 0; x < f.Width(); x++ {
			c := Coord{x, y}
			if f.At(c) == 'S' {
				return c
			}
		}
	}

	panic(fmt.Sprintf("Sybmol S not found on a field"))
}

type Direction int

const (
	Top Direction = iota
	Right
	Bottom
	Left
)

func canSymbolConnect(s byte, dir Direction) bool {
	symbols := map[byte][]Direction{
		'S': {Top, Right, Bottom, Left},
		'|': {Top, Bottom},
		'-': {Left, Right},
		'L': {Top, Right},
		'J': {Left, Top},
		'7': {Bottom, Left},
		'F': {Right, Bottom},
	}

	if directions, exists := symbols[s]; exists {
		return slices.Contains(directions, dir)
	}

	return false
}

func findConnected(f *CharField, c Coord) (ans []Coord) {
	symbol := f.At(c)

	if canSymbolConnect(symbol, Top) {
		if canSymbolConnect(f.At(c.Top()), Bottom) {
			ans = append(ans, c.Top())
		}
	}

	if canSymbolConnect(symbol, Right) {
		if canSymbolConnect(f.At(c.Right()), Left) {
			ans = append(ans, c.Right())
		}
	}

	if canSymbolConnect(symbol, Bottom) {
		if canSymbolConnect(f.At(c.Bottom()), Top) {
			ans = append(ans, c.Bottom())
		}
	}

	if canSymbolConnect(symbol, Left) {
		if canSymbolConnect(f.At(c.Left()), Right) {
			ans = append(ans, c.Left())
		}
	}

	return
}
