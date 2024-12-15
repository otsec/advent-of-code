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

const (
	CharBox        = 'O'
	CharEmpty      = '.'
	CharRobot      = '@'
	CharWall       = '#'
	CharLargeBoxP1 = '['
	CharLargeBoxP2 = ']'
)

func part1(input string) (ans int) {
	field, moves := parseInput(input)

	pos := findStart(field)
	field.Set(pos, CharEmpty)
	for _, move := range moves {
		pos = makeMove(field, pos, move)
	}

	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			if field.At(Coord{x, y}) == CharBox {
				ans += x + 100*y
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	field, moves := parseInput(input)
	field = increaseField(field)

	pos := findStart(field)
	field.Set(pos, CharEmpty)
	for _, move := range moves {
		pos = makeMove(field, pos, move)
	}

	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			if field.At(Coord{x, y}) == CharLargeBoxP1 {
				ans += x + 100*y
			}
		}
	}

	return ans
}

func findStart(field *CharField) Coord {
	for x := 0; x < field.Width(); x++ {
		for y := 0; y < field.Height(); y++ {
			coord := Coord{x, y}
			if field.At(coord) == CharRobot {
				return coord
			}
		}
	}
	panic("start not found")
}

func makeMove(field *CharField, curr Coord, dir Direction) Coord {
	next := curr.Next(dir)
	nextSymbol := field.At(next)

	if nextSymbol == CharWall {
		return curr
	}

	if nextSymbol == CharBox {
		pushSmallBox(field, next, dir)

		nextSymbol = field.At(next)
		if nextSymbol == CharBox {
			return curr
		}
	}

	if nextSymbol == CharLargeBoxP1 || nextSymbol == CharLargeBoxP2 {
		if canLargeBoxMove(field, next, dir) {
			pushLargeBox(field, next, dir)
			return next
		} else {
			return curr
		}
	}

	return next
}

func pushSmallBox(field *CharField, pos Coord, dir Direction) {
	next := pos.Next(dir)

	if field.At(next) == CharBox {
		pushSmallBox(field, next, dir)
	}

	if field.At(next) == CharEmpty {
		field.Set(next, CharBox)
		field.Set(pos, CharEmpty)
	}
}

func canLargeBoxMove(field *CharField, pos Coord, dir Direction) bool {
	posSymbol := field.At(pos)

	var p1, p2 Coord
	if posSymbol == CharLargeBoxP1 {
		p1, p2 = pos, pos.Right()
	} else if posSymbol == CharLargeBoxP2 {
		p1, p2 = pos.Left(), pos
	} else {
		panic("wrong large box pos")
	}

	var nextCoords []Coord
	switch {
	case dir == Top:
		nextCoords = []Coord{p1.Top(), p2.Top()}
	case dir == Right:
		nextCoords = []Coord{p2.Right()}
	case dir == Bottom:
		nextCoords = []Coord{p1.Bottom(), p2.Bottom()}
	case dir == Left:
		nextCoords = []Coord{p1.Left()}
	default:
		panic("wrong push direction")
	}

	for _, next := range nextCoords {
		nextSymbol := field.At(next)
		if nextSymbol == CharWall {
			return false
		}
		if nextSymbol == CharLargeBoxP1 || nextSymbol == CharLargeBoxP2 {
			if !canLargeBoxMove(field, next, dir) {
				return false
			}
		}
	}

	return true
}

func pushLargeBox(field *CharField, pos Coord, dir Direction) {
	posSymbol := field.At(pos)

	var p1, p2 Coord
	if posSymbol == CharLargeBoxP1 {
		p1, p2 = pos, pos.Right()
	} else if posSymbol == CharLargeBoxP2 {
		p1, p2 = pos.Left(), pos
	} else {
		panic("wrong large box pos")
	}

	var nextCoords []Coord
	switch {
	case dir == Top:
		nextCoords = []Coord{p1.Top(), p2.Top()}
	case dir == Right:
		nextCoords = []Coord{p2.Right()}
	case dir == Bottom:
		nextCoords = []Coord{p1.Bottom(), p2.Bottom()}
	case dir == Left:
		nextCoords = []Coord{p1.Left()}
	default:
		panic("wrong push direction")
	}

	for _, next := range nextCoords {
		nextSymbol := field.At(next)
		if nextSymbol == CharLargeBoxP1 || nextSymbol == CharLargeBoxP2 {
			pushLargeBox(field, next, dir)
		}
	}

	field.Set(p1, CharEmpty)
	field.Set(p2, CharEmpty)
	field.Set(p1.Next(dir), CharLargeBoxP1)
	field.Set(p2.Next(dir), CharLargeBoxP2)
}

func increaseField(field *CharField) *CharField {
	newField := []byte{}
	for i, line := range field.lines {
		for _, symbol := range line {
			switch symbol {
			case CharBox:
				newField = append(newField, CharLargeBoxP1, CharLargeBoxP2)
			case CharEmpty:
				newField = append(newField, CharEmpty, CharEmpty)
			case CharRobot:
				newField = append(newField, CharRobot, CharEmpty)
			case CharWall:
				newField = append(newField, CharWall, CharWall)
			default:
				panic(fmt.Sprintf("unknown char to increase field %v", symbol))
			}
		}

		isLast := i == len(field.lines)-1
		if !isLast {
			newField = append(newField, '\n')
		}
	}
	return NewCharField(string(newField))
}

func parseInput(input string) (field *CharField, moves []Direction) {
	parts := strings.Split(input, "\n\n")
	field = NewCharField(parts[0])
	moves = parseMoves(parts[1])
	return field, moves
}

func parseMoves(input string) []Direction {
	moves := []Direction{}
	for _, char := range input {
		switch char {
		case '^':
			moves = append(moves, Top)
		case '>':
			moves = append(moves, Right)
		case 'v':
			moves = append(moves, Bottom)
		case '<':
			moves = append(moves, Left)
		case '\n':
			// nothing
		default:
			panic(fmt.Sprintf("unknown move %v", string(char)))
		}
	}
	return moves
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
