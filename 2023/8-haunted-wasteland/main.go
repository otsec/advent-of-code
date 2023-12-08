package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"regexp"
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
	parsed := parseInput(input)

	position := parsed.starts[0]
	dir := DirectionIterator{step: 0, direction: parsed.directions}
	for i := 1; i < 100000; i++ {
		position = parsed.makeStep(position, dir.Curr())
		if position == "ZZZ" {
			return i
		} else {
			dir.Next()
		}
	}

	panic("Max iterations reached for part 1")
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	var loops []int
	for _, p := range parsed.starts {
		if endsWith(p, 'A') {
			loops = append(loops, findLoop(parsed, p))
		}
	}

	ans = LCM(loops[0], loops[1])
	for i := 2; i < len(loops); i++ {
		ans = LCM(ans, loops[i])
	}

	return ans
}

type Step struct {
	Left, Right string
}

type Input struct {
	directions string
	starts     []string
	steps      map[string]Step
}

func (i *Input) makeStep(from string, direction byte) string {
	if direction == 'L' {
		return i.steps[from].Left
	} else {
		return i.steps[from].Right
	}
}

type DirectionIterator struct {
	step      int
	direction string
}

func (d *DirectionIterator) Curr() byte {
	return d.direction[d.step]
}

func (d *DirectionIterator) Next() {
	d.step += 1
	if d.step >= len(d.direction) {
		d.step = 0
	}
}

func parseInput(input string) Input {
	segments := strings.Split(input, "\n\n")

	starts := []string{}
	steps := map[string]Step{}
	re := regexp.MustCompile(`\w{3}`)
	for _, line := range strings.Split(segments[1], "\n") {
		matches := re.FindAllString(line, -1)
		starts = append(starts, matches[0])
		steps[matches[0]] = Step{matches[1], matches[2]}
	}

	return Input{segments[0], starts, steps}
}

func endsWith(position string, letter byte) bool {
	return position[len(position)-1] == letter
}

func findLoop(input Input, start string) int {
	position := start
	dir := DirectionIterator{step: 0, direction: input.directions}
	for i := 1; i < 100000; i++ {
		position = input.makeStep(position, dir.Curr())
		if endsWith(position, 'Z') {
			return i
		}

		dir.Next()
	}

	panic(fmt.Sprintf("Max iterations reached for %v", start))
}

// GCD finds greatest common divisor via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// LCM finds least common multiple via GCD
func LCM(a, b int) int {
	return a * b / GCD(a, b)
}
