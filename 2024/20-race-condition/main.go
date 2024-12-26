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

//go:embed input_test.txt
var input_test string

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

func part1(input string) int {
	field := NewCharField(input)
	return findFastCheats(field, 100)
}

func part2(input string) (ans int) {
	field := NewCharField(input)
	return findLongCheats(field, 20, 100)
}

func findFastCheats(field *CharField, minSavedSeconds int) (ans int) {
	pathMemo := findNormalPath(field)

	for coord, currVal := range pathMemo {
		variants := []Direction{Top, Right, Bottom, Left}
		for _, dir := range variants {
			oneStepCoord := coord.Next(dir)
			twoStepCoord := oneStepCoord.Next(dir)
			if !field.Within(oneStepCoord) || !field.Within(twoStepCoord) {
				continue
			}
			if field.At(oneStepCoord) != '#' {
				continue
			}

			if twoStepVal, ok := pathMemo[twoStepCoord]; ok {
				secondsSaved := currVal - twoStepVal - 2
				if secondsSaved >= minSavedSeconds {
					ans += 1
				}
			}
		}
	}

	return ans
}

func findLongCheats(field *CharField, maxCheatSeconds, minSavedSeconds int) (ans int) {
	pathMemo := findNormalPath(field)

	for coord, currVal := range pathMemo {
		cheatMemo := findCheatPaths(field, coord, maxCheatSeconds)

		for endCheatCoord, cheatSeconds := range cheatMemo {
			secondsSaved := currVal - pathMemo[endCheatCoord] - cheatSeconds
			if secondsSaved >= minSavedSeconds {
				ans += 1
			}
		}
	}

	return ans
}

func findNormalPath(field *CharField) map[Coord]int {
	memo := map[Coord]int{}

	end := findSymbol(field, 'E')

	memo[end] = 0
	queue := []Coord{end}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		variants := []Coord{curr.Top(), curr.Right(), curr.Bottom(), curr.Left()}
		for _, variant := range variants {
			if !field.Within(variant) {
				continue
			}
			if field.At(variant) == '#' {
				continue
			}
			if field.At(variant) == '.' || field.At(variant) == 'S' {
				newVal := memo[curr] + 1
				if oldVal, ok := memo[variant]; !ok || newVal < oldVal {
					memo[variant] = newVal
					queue = append(queue, variant)
				}
			}
		}
	}

	return memo
}

func findCheatPaths(field *CharField, start Coord, maxSeconds int) map[Coord]int {
	cheatMemo := map[Coord]int{}

	for x := start.x - maxSeconds; x <= start.x+maxSeconds; x++ {
		for y := start.y - maxSeconds; y <= start.y+maxSeconds; y++ {
			curr := Coord{x, y}
			dist := AbsInt(start.x-curr.x) + AbsInt(start.y-curr.y)
			if field.Within(curr) && dist <= maxSeconds {
				val := field.At(curr)
				if val == '.' || val == 'E' {
					cheatMemo[curr] = dist
				}
			}
		}
	}

	return cheatMemo
}

func findSymbol(field *CharField, s byte) Coord {
	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			coord := Coord{x, y}
			if field.At(coord) == s {
				return coord
			}
		}
	}

	panic("symbol not found")
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

func AbsInt(in int) int {
	if in < 0 {
		return -in
	}
	return in
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
