package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os/exec"
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

	pos := 50
	for _, move := range parsed {
		if move.dir == 'L' {
			pos, _ = rotate(pos, -move.steps)
		} else {
			pos, _ = rotate(pos, move.steps)
		}
		if pos == 0 {
			ans++
		}
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	pos := 50
	for _, move := range parsed {
		var zeroPasses int
		if move.dir == 'L' {
			pos, zeroPasses = rotate(pos, -move.steps)
		} else {
			pos, zeroPasses = rotate(pos, move.steps)
		}

		if pos == 0 {
			ans += 1
		}
		ans += zeroPasses

		log.Printf("move=%v pos=%v zeroes=%v ans=%v", move, pos, zeroPasses, ans)
	}

	return ans
}

type Move struct {
	dir   rune
	steps int
}

func (m Move) String() string {
	return fmt.Sprintf("%c%v", m.dir, m.steps)
}

func rotate(from int, steps int) (res, zeroPasses int) {
	if steps > 0 {
		zeroPasses = steps / 100
	} else {
		zeroPasses = steps / -100
	}

	res = from + (steps % 100)

	if res < 0 {
		res += 100
		if from != 0 && res != 0 {
			zeroPasses += 1
		}
	}

	if res > 99 {
		res = res - 100
		if from != 0 && res != 0 {
			zeroPasses += 1
		}
	}

	return
}

func parseInput(input string) []Move {
	lines := strings.Split(input, "\n")

	moves := make([]Move, len(lines))
	for i, line := range lines {
		dir := rune(line[0])

		steps, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(err)
		}

		moves[i] = Move{dir, steps}
	}

	return moves
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
