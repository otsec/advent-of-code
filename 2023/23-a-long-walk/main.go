package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
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
	field := NewCharField(input)

	start := Coord{1, 0}
	finish := Coord{field.Width() - 2, field.Height() - 1}

	queue := linkedlistqueue.New()
	queue.Enqueue(Path{start, []Coord{}})
	for queue.Size() > 0 {
		item, _ := queue.Dequeue()
		p := item.(Path)

		if p.pos == finish {
			length := len(p.history)
			ans = max(ans, length)
			continue
		}

		if field.Within(p.pos.Top()) {
			if field.At(p.pos.Top()) == '.' || field.At(p.pos.Top()) == '^' {
				if !slices.Contains(p.history, p.pos.Top()) {
					queue.Enqueue(Path{p.pos.Top(), p.ToHistory()})
				}
			}
		}
		if field.At(p.pos.Right()) == '.' || field.At(p.pos.Right()) == '>' {
			if !slices.Contains(p.history, p.pos.Right()) {
				queue.Enqueue(Path{p.pos.Right(), p.ToHistory()})
			}
		}
		if field.At(p.pos.Bottom()) == '.' || field.At(p.pos.Bottom()) == 'v' {
			if !slices.Contains(p.history, p.pos.Bottom()) {
				queue.Enqueue(Path{p.pos.Bottom(), p.ToHistory()})
			}
		}
		if field.At(p.pos.Left()) == '.' || field.At(p.pos.Left()) == '<' {
			if !slices.Contains(p.history, p.pos.Left()) {
				queue.Enqueue(Path{p.pos.Left(), p.ToHistory()})
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	field := NewCharField(input)

	start := Coord{1, 0}
	finish := Coord{field.Width() - 2, field.Height() - 1}

	queue := linkedlistqueue.New()
	queue.Enqueue(Path{start, []Coord{}})
	for queue.Size() > 0 {
		item, _ := queue.Dequeue()
		p := item.(Path)

		if p.pos == finish {
			length := len(p.history)
			ans = max(ans, length)
			continue
		}

		candidates := []Coord{p.pos.Top(), p.pos.Right(), p.pos.Bottom(), p.pos.Left()}
		for _, c := range candidates {
			if field.Within(c) && field.At(c) != '#' && !p.WasAt(c) {
				queue.Enqueue(Path{c, p.ToHistory()})
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

type Path struct {
	pos     Coord
	history []Coord
}

func (p *Path) WasAt(c Coord) bool {
	return slices.Contains(p.history, c)
}

func (p *Path) ToHistory() []Coord {
	newHist := make([]Coord, len(p.history)+1)
	copy(newHist, p.history)
	newHist[len(newHist)-1] = p.pos
	return newHist
}

func (p *Path) Draw(f *CharField) {
	for _, c := range p.history {
		f.Set(c, 'O')
	}
	f.Set(p.pos, 'O')
}
