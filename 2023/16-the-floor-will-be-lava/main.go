package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/emirpasic/gods/sets/hashset"
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

func part1(input string) int {
	field := NewCharField(input, '.')
	start := Beam{Coord{0, 0}, Right}
	return calcEnergy(field, start)
}

func part2(input string) (ans int) {
	field := NewCharField(input, '.')

	for x := 0; x < field.Width(); x++ {
		startTop := Beam{Coord{x, 0}, Bottom}
		ansTop := calcEnergy(field, startTop)

		startBottom := Beam{Coord{x, field.Height() - 1}, Top}
		ansBottom := calcEnergy(field, startBottom)

		ans = max(ans, ansTop, ansBottom)
	}

	for y := 0; y < field.Height(); y++ {
		startLeft := Beam{Coord{0, y}, Right}
		ansLeft := calcEnergy(field, startLeft)

		startRight := Beam{Coord{field.Width() - 1, y}, Left}
		ansRight := calcEnergy(field, startRight)

		ans = max(ans, ansLeft, ansRight)
	}

	return ans
}

func calcEnergy(field *CharField, start Beam) (ans int) {
	stepsField := NewCharField(input, '.')
	stepsSet := hashset.New()

	queue := linkedlistqueue.New()
	queue.Enqueue(start)
	for !queue.Empty() {
		v, _ := queue.Dequeue()
		beam := v.(Beam)

		if !field.Within(beam.coord) {
			continue
		}

		if stepsSet.Contains(beam) {
			continue
		} else {
			stepsSet.Add(beam)
			stepsField.Set(beam.coord, '#')
		}

		if field.At(beam.coord) == '.' {
			queue.Enqueue(beam.Next(beam.dir))
		}

		if field.At(beam.coord) == '|' && (beam.dir == Top || beam.dir == Bottom) {
			queue.Enqueue(beam.Next(beam.dir))
		}
		if field.At(beam.coord) == '|' && (beam.dir == Left || beam.dir == Right) {
			queue.Enqueue(beam.Next(Top))
			queue.Enqueue(beam.Next(Bottom))
		}

		if field.At(beam.coord) == '-' && (beam.dir == Top || beam.dir == Bottom) {
			queue.Enqueue(beam.Next(Left))
			queue.Enqueue(beam.Next(Right))
		}
		if field.At(beam.coord) == '-' && (beam.dir == Left || beam.dir == Right) {
			queue.Enqueue(beam.Next(beam.dir))
		}

		if field.At(beam.coord) == '/' {
			switch beam.dir {
			case Bottom:
				queue.Enqueue(beam.Next(Left))
			case Left:
				queue.Enqueue(beam.Next(Bottom))
			case Top:
				queue.Enqueue(beam.Next(Right))
			case Right:
				queue.Enqueue(beam.Next(Top))
			}
		}

		if field.At(beam.coord) == '\\' {
			switch beam.dir {
			case Bottom:
				queue.Enqueue(beam.Next(Right))
			case Left:
				queue.Enqueue(beam.Next(Top))
			case Top:
				queue.Enqueue(beam.Next(Left))
			case Right:
				queue.Enqueue(beam.Next(Bottom))
			}
		}
	}

	for y := 0; y < stepsField.Height(); y++ {
		for x := 0; x < stepsField.Width(); x++ {
			if stepsField.At(Coord{x, y}) == '#' {
				ans += 1
			}
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
	if f.Within(c) {
		return f.lines[c.y][c.x]
	} else {
		return f.emptyChar
	}
}

func (f *CharField) Set(c Coord, v byte) {
	line := []byte(f.lines[c.y])
	line[c.x] = v
	f.lines[c.y] = string(line)
}

func (f *CharField) IsEmpty(c Coord) bool {
	return f.At(c) == '.'
}

func (f *CharField) ToString() string {
	return strings.Join(f.lines, "\n")
}

func (f *CharField) Print() {
	fmt.Println(f.ToString())
}

type Direction int

const (
	Top Direction = iota
	Right
	Bottom
	Left
)

type Beam struct {
	coord Coord
	dir   Direction
}

func (b *Beam) Next(dir Direction) Beam {
	var next Coord

	switch dir {
	case Top:
		next = b.coord.Top()
	case Right:
		next = b.coord.Right()
	case Bottom:
		next = b.coord.Bottom()
	case Left:
		next = b.coord.Left()
	default:
		panic(fmt.Sprint("Unknown dir", dir))
	}

	return Beam{next, dir}
}

func NewCharField(input string, emptyChar byte) *CharField {
	lines := strings.Split(input, "\n")
	return &CharField{lines, emptyChar}
}
