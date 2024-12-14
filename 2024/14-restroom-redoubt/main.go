package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
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

func part1(input string) int {
	fieldSize := Coord{101, 103}
	seconds := 100
	return simulate(input, fieldSize, seconds)
}

func part2(input string) int {
	robots := parseInput(input)

	fieldSize := Coord{101, 103}

	for s := 1; s < 10000; s++ {
		resultPositions := make([]Coord, len(robots))
		for i, robot := range robots {
			resultPositions[i] = calculatePosition(robot.position, robot.speed, fieldSize, s)
		}

		if looksLikeTree(resultPositions) {
			field := NewEmptyCharField(fieldSize.x, fieldSize.y, '.')
			for _, pos := range resultPositions {
				field.Set(pos, '#')
			}
			field.Println()

			return s
		}
	}

	return -1
}

type Coord struct {
	x, y int
}

type ParsedInput struct {
	position, speed Coord
}

func looksLikeTree(positions []Coord) bool {
	pointsX := map[int][]int{}
	for _, pos := range positions {
		pointsX[pos.x] = append(pointsX[pos.x], pos.y)
	}

	for _, ys := range pointsX {
		sort.Ints(ys)

		consecutives := 0
		prev := ys[0]
		for i := 1; i < len(ys); i++ {
			if ys[i] == prev+1 {
				consecutives++
			} else {
				consecutives = 0
			}
			prev = ys[i]

			if consecutives > 5 {
				return true
			}
		}
	}

	return false
}

func simulate(rawInput string, fieldSize Coord, seconds int) int {
	robots := parseInput(rawInput)

	resultPositions := make([]Coord, len(robots))
	for i, robot := range robots {
		resultPositions[i] = calculatePosition(robot.position, robot.speed, fieldSize, seconds)
	}

	return calculateAns(resultPositions, fieldSize)
}

func calculatePosition(start, speed, size Coord, seconds int) Coord {
	finalX := (start.x + speed.x*seconds) % size.x
	if finalX < 0 {
		finalX += size.x
	}

	finalY := (start.y + speed.y*seconds) % size.y
	if finalY < 0 {
		finalY += size.y
	}

	return Coord{finalX, finalY}
}

func detectQuadrant(pos, fieldSize Coord) int {
	halfX := fieldSize.x / 2
	halfY := fieldSize.y / 2

	if pos.x < halfX && pos.y < halfY {
		return 1
	} else if pos.x > halfX && pos.y < halfY {
		return 2
	} else if pos.x < halfX && pos.y > halfY {
		return 3
	} else if pos.x > halfX && pos.y > halfY {
		return 4
	} else {
		return 0
	}
}

func calculateAns(robots []Coord, fieldSize Coord) (ans int) {
	memo := [5]int{}
	for _, robot := range robots {
		quadrant := detectQuadrant(robot, fieldSize)
		memo[quadrant]++
	}
	return memo[1] * memo[2] * memo[3] * memo[4]
}

func parseInput(input string) []ParsedInput {
	lines := strings.Split(input, "\n")
	inputs := make([]ParsedInput, len(lines))
	for i, line := range lines {
		nums := parseNumbers(line)
		position := Coord{nums[0], nums[1]}
		speed := Coord{nums[2], nums[3]}
		inputs[i] = ParsedInput{position, speed}
	}
	return inputs
}

func parseNumbers(input string) []int {
	re := regexp.MustCompile(`[-0-9]+`)
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

type CharField struct {
	width, height int
	symbols       []byte
}

func (f *CharField) Set(c Coord, v byte) {
	f.symbols[c.y*f.width+c.x] = v
}

func (f *CharField) Println() {
	for i := 0; i < f.height; i++ {
		fmt.Println(string(f.symbols[i*f.width : i*f.width+f.width]))
	}
	fmt.Println()
}

func NewEmptyCharField(width, height int, emptySymbol byte) *CharField {
	symbols := bytes.Repeat([]byte{emptySymbol}, width*height)
	return &CharField{width, height, symbols}
}
