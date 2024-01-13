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

func part1(input string) int {
	field := NewCharField(input)

	for x := 0; x < field.Width(); x++ {
		line := &Line{Coord{x, 0}, Coord{x, field.Height() - 1}}
		slideLine(field, line)
	}

	return calcLoad(field)
}

func part2(input string) (ans int) {
	field := NewCharField(input)

	slideCycles(field, 1000000000)

	return calcLoad(field)
}

func slideCycles(f *CharField, cycles int) {
	maxX := f.Width() - 1
	maxY := f.Height() - 1

	testCycles := 1000
	mem := map[string]int{}

	for i := 1; i <= testCycles; i++ {
		// north
		for x := 0; x <= maxX; x++ {
			line := &Line{Coord{x, 0}, Coord{x, maxY}}
			slideLine(f, line)
		}

		// west
		for y := 0; y <= maxY; y++ {
			line := &Line{Coord{0, y}, Coord{maxX, y}}
			slideLine(f, line)
		}

		// south
		for x := 0; x <= maxX; x++ {
			line := &Line{Coord{x, maxY}, Coord{x, 0}}
			slideLine(f, line)
		}

		// east
		for y := 0; y <= maxY; y++ {
			line := &Line{Coord{maxX, y}, Coord{0, y}}
			slideLine(f, line)
		}

		//fmt.Println("Cycle", i, "load", calcLoad(f))
		//f.Print()
		//fmt.Println()

		hash := f.ToString()
		if _, exists := mem[hash]; exists {
			start := mem[hash]
			end := i
			loopLen := end - start
			cyclesToSkip := (cycles - i) / loopLen
			cyclesLeft := cycles - i - cyclesToSkip*loopLen
			i = testCycles - cyclesLeft
			mem = map[string]int{}

			// log.Fatalln("Loop found between cycles", start, "and", end, "skip", cyclesToSkip, "left", cyclesLeft)
		} else {
			mem[hash] = i
		}
	}
}

func slideLine(f *CharField, l *Line) {
	var coords []Coord

	curr := l.start
	coords = append(coords, curr)

	for curr != l.end {
		if curr.x < l.end.x {
			curr.x += 1
		} else if curr.x > l.end.x {
			curr.x -= 1
		} else if curr.y < l.end.y {
			curr.y += 1
		} else if curr.y > l.end.y {
			curr.y -= 1
		}

		coords = append(coords, curr)
	}

	bytesLine := make([]byte, len(coords))
	for i, coord := range coords {
		bytesLine[i] = f.At(coord)
	}

	slideBytes(bytesLine)

	for i, coord := range coords {
		f.Set(coord, bytesLine[i])
	}
}

func slideBytes(line []byte) {
	swapWith := 0
	for i := 0; i < len(line); i++ {
		if i == swapWith && line[i] == 'O' {
			swapWith = i + 1
			continue
		}
		if line[i] == '#' {
			swapWith = i + 1
			continue
		}
		if line[i] == 'O' {
			line[swapWith], line[i] = line[i], line[swapWith]
			swapWith++
		}
	}
}

func calcLoad(f *CharField) int {
	ans := 0

	for y := 0; y < f.Height(); y++ {
		for x := 0; x < f.Width(); x++ {
			if f.At(Coord{x, y}) == 'O' {
				ans += f.Height() - y
			}
		}
	}

	return ans
}

type Coord struct {
	x, y int
}

type Line struct {
	start, end Coord
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
