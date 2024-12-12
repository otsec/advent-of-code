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

	shapes := findAllShapes(field)
	for _, shape := range shapes {
		ans += shape.Area() * shape.Perimeter()
	}

	return ans
}

func part2(input string) (ans int) {
	field := NewCharField(input)

	shapes := findAllShapes(field)
	for _, shape := range shapes {
		ans += shape.Area() * shape.Sides()
	}

	return ans
}

func findAllShapes(field *CharField) []Shape {
	var shapes []Shape
	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			coord := Coord{x, y}
			if !isInAnotherShape(&shapes, coord) {
				shape := findShape(field, coord)
				shapes = append(shapes, *shape)
			}
		}
	}
	return shapes
}

func findShape(f *CharField, coord Coord) *Shape {
	shape := &Shape{
		symbol: f.At(coord),
		points: map[Coord]bool{},
	}
	shape.points[coord] = true

	queue := []Coord{coord}
	validated := map[Coord]bool{}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		variants := []Coord{curr.Top(), curr.Right(), curr.Bottom(), curr.Left()}
		for _, variant := range variants {
			if _, ok := validated[variant]; ok {
				continue
			} else {
				validated[variant] = true
			}

			if f.Within(variant) && f.At(variant) == shape.symbol {
				shape.Add(variant)
				queue = append(queue, variant)
			}
		}
	}
	return shape
}

func isInAnotherShape(shapes *[]Shape, coord Coord) bool {
	for _, shape := range *shapes {
		if shape.Has(coord) {
			return true
		}
	}
	return false
}

type Shape struct {
	symbol byte
	points map[Coord]bool
}

func (s *Shape) Add(coord Coord) {
	s.points[coord] = true
}

func (s *Shape) Has(coord Coord) bool {
	return s.points[coord]
}

func (s *Shape) Area() int {
	return len(s.points)
}

func (s *Shape) Perimeter() (ans int) {
	for point := range s.points {
		adjacent := []Coord{point.Top(), point.Right(), point.Bottom(), point.Left()}
		for _, coord := range adjacent {
			if !s.Has(coord) {
				ans += 1
			}
		}
	}
	return ans
}

func (s *Shape) Sides() int {
	memo := SideMemo{
		registered: map[string]bool{},
		calculated: map[string]bool{},
		ignored:    map[string]bool{},
	}

	for coord := range s.points {
		if !s.Has(coord.Top()) {
			memo.RegisterSegment(coord, Top)
		}
		if !s.Has(coord.Right()) {
			memo.RegisterSegment(coord, Right)
		}
		if !s.Has(coord.Bottom()) {
			memo.RegisterSegment(coord, Bottom)
		}
		if !s.Has(coord.Left()) {
			memo.RegisterSegment(coord, Left)
		}
	}

	for coord := range s.points {
		if !s.Has(coord.Top()) {
			memo.CalculateSide(coord, Top)
		}
		if !s.Has(coord.Right()) {
			memo.CalculateSide(coord, Right)
		}
		if !s.Has(coord.Bottom()) {
			memo.CalculateSide(coord, Bottom)
		}
		if !s.Has(coord.Left()) {
			memo.CalculateSide(coord, Left)
		}
	}

	return memo.Count()
}

type SideMemo struct {
	registered map[string]bool
	calculated map[string]bool
	ignored    map[string]bool
}

func (sm *SideMemo) RegisterSegment(curr Coord, side Direction) {
	code := sm.encode(curr, side)
	sm.registered[code] = true
}

func (sm *SideMemo) CalculateSide(curr Coord, side Direction) {
	var scan []Direction
	if side == Top || side == Bottom {
		scan = append(scan, Left, Right)
	} else if side == Left || side == Right {
		scan = append(scan, Top, Bottom)
	} else {
		panic("unknown side type")
	}

	code := sm.encode(curr, side)

	if sm.ignored[code] {
		return
	} else {
		sm.calculated[code] = true
	}

	for _, dir := range scan {
		next := curr.Next(dir)
		for {
			nextCode := sm.encode(next, side)
			if !sm.registered[nextCode] {
				break
			}
			sm.ignored[nextCode] = true
			next = next.Next(dir)
		}
	}
}

func (sm *SideMemo) Count() int {
	return len(sm.calculated)
}

func (sm *SideMemo) encode(coord Coord, edge Direction) string {
	return fmt.Sprintf("%v,%v,%v", coord.x, coord.y, edge)
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
