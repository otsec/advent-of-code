package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os/exec"
	"regexp"
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
		CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) (ans int) {
	parsed := parseInput(input)

	for _, v := range parsed {
		seq := countShortestSequence(v, 2)
		num := parseNumbers(v)[0]
		ans += seq * num
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	for _, v := range parsed {
		seq := countShortestSequence(v, 3)
		num := parseNumbers(v)[0]
		ans += seq * num
	}

	return ans
}

func countShortestSequence(code string, robots int) int {
	nkp := CreateNumericKeypad()
	curr := findSymbol(nkp, CharActivate)
	var paths []string
	for _, symbol := range []byte(code) {
		//fmt.Println(code, string(symbol), findMoves(nkp, curr, symbol))
		//fmt.Println(code, string(symbol), pathsToMoves(findMoves(nkp, curr, symbol)))
		paths = combinePaths(paths, findMoves(nkp, curr, symbol))
		curr = findSymbol(nkp, symbol)
	}

	queue := make([]int, 0)
	for _, path := range paths {
		queue.append({path, 25})
	}

	paths = []string{"^<<A", "<^<A"}

	dkp := CreateDirectionalKeypad()
	var minPathSize int

	type MoveMemo struct {
		startSymbol, endSymbol byte
		robots                 int
	}

	memo := map[MoveMemo][]string{}

	for _, path := range paths {
		startSymbol := byte('A')
		for _, endSymbol := range []byte(path) {
			start := findSymbol(dkp, startSymbol)
			symbolPaths := pathsToMoves(findMoves(dkp, start, endSymbol))
			memoKey := MoveMemo{startSymbol, endSymbol, 1}
			memo[memoKey] = symbolPaths
			startSymbol = endSymbol
		}
	}

	fmt.Println(memo)

	//for r := 0; r < robots; r++ {
	//	fmt.Println("robot", r+1)
	//
	//	var thisRobotPaths []string
	//	for _, path := range paths {
	//		curr = findSymbol(dkp, 'A')
	//		var thisGroupPaths []string
	//		for _, symbol := range []byte(path) {
	//			symbolPaths := findMoves(dkp, curr, symbol)
	//			thisGroupPaths = combinePaths(thisGroupPaths, symbolPaths)
	//			curr = findSymbol(dkp, symbol)
	//		}
	//		thisRobotPaths = append(thisRobotPaths, thisGroupPaths...)
	//	}
	//
	//	minPathSize = len(thisRobotPaths[0])
	//	for _, path := range thisRobotPaths {
	//		minPathSize = min(minPathSize, len(path))
	//	}
	//
	//	paths = slices.DeleteFunc(thisRobotPaths, func(s string) bool {
	//		return len(s) != minPathSize
	//	})
	//}

	return minPathSize
}

func combinePaths(paths []string, dirs [][]Direction) []string {
	var ans []string

	if len(paths) == 0 && len(dirs) == 0 {
		ans = append(ans, string(CharActivate))
	} else if len(dirs) == 0 {
		for _, path := range paths {
			ans = append(ans, path+string(CharActivate))
		}
	} else if len(paths) == 0 {
		for _, group := range dirs {
			ans = append(ans, stringifyDirections(group)+string(CharActivate))
		}
	} else {
		for _, path := range paths {
			for _, group := range dirs {
				ans = append(ans, path+stringifyDirections(group)+string(CharActivate))
			}
		}
	}

	return ans
}

func findShortestDirectionalPaths(dkp *CharField, memo *map[string][][]string, inputCode string) int {
	pathGroups := make([][]string, len(inputCode))

	startSymbol := byte('A')
	for i, endSymbol := range []byte(inputCode) {
		start := findSymbol(dkp, startSymbol)
		pathGroups[i] = pathsToMoves(findMoves(dkp, start, endSymbol))
		startSymbol = endSymbol
	}

	(*memo)[inputCode] = pathGroups

	return findShortestPath(pathGroups)
}

func findShortestPath(pathGroups [][]string) int {
	var ans int

	for _, group := range pathGroups {
		if len(group) > 0 {
			minPath := len(group[0])
			for _, path := range group {
				minPath = min(minPath, len(path))
			}
			ans += minPath
		}
	}

	return ans
}

func pathsToMoves(paths [][]Direction) []string {
	ans := make([]string, len(paths))
	for i, group := range paths {
		ans[i] = stringifyDirections(group) + "A"
	}
	return ans
}

func pathToMoves(path []Direction) string {
	return stringifyDirections(path) + "A"
}

func stringifyDirections(dirs []Direction) string {
	ans := make([]byte, len(dirs))
	for i, dir := range dirs {
		switch dir {
		case Top:
			ans[i] = CharTop
		case Right:
			ans[i] = CharRight
		case Bottom:
			ans[i] = CharBottom
		case Left:
			ans[i] = CharLeft
		default:
			panic("unknown direction")
		}
	}
	return string(ans)
}

const (
	CharTop      = '^'
	CharRight    = '>'
	CharBottom   = 'v'
	CharLeft     = '<'
	CharEmpty    = 'x'
	CharActivate = 'A'
)

func CreateNumericKeypad() *CharField {
	return &CharField{lines: []string{"789", "456", "123", "x0A"}}
}

func CreateDirectionalKeypad() *CharField {
	return &CharField{lines: []string{"x^A", "<v>"}}
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

func findMoves(field *CharField, start Coord, end byte) [][]Direction {
	if field.At(start) == end {
		return [][]Direction{}
	}

	type History struct {
		curr Coord
		path []Direction
	}

	var queue []History
	queue = append(queue, History{start, []Direction{}})

	var ans [][]Direction
	minAns := -1

	for len(queue) > 0 {
		curr := queue[0].curr
		path := queue[0].path

		queue = queue[1:]

		for _, direction := range []Direction{Top, Right, Bottom, Left} {
			next := curr.Next(direction)
			if !field.Within(next) || field.At(next) == CharEmpty {
				continue
			}

			newPath := make([]Direction, len(path)+1)
			copy(newPath, path)
			newPath[len(path)] = direction

			if field.At(next) == end {
				if minAns == -1 || minAns == len(newPath) {
					minAns = len(newPath)
					ans = append(ans, newPath)
				}
			} else {
				if minAns == -1 || len(newPath) < minAns {
					queue = append(queue, History{next, newPath})
				}
			}
		}
	}

	return ans
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

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}

func parseNumbers(input string) []int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(input, -1)
	nums := make([]int, len(matches))
	for i, match := range matches {
		nums[i], _ = strconv.Atoi(match)
	}
	return nums
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
