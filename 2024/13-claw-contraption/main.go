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
	machines := parseInput(input)

	for _, m := range machines {
		res := findWin(m.a, m.b, m.win)
		if res != -1 {
			ans += res
		}
	}

	return ans
}

func part2(input string) (ans int) {
	machines := parseInput(input)
	for i, m := range machines {
		m.win.x += 10000000000000
		m.win.y += 10000000000000
		machines[i] = m
	}

	for _, m := range machines {
		res := findWinSmart(m.a, m.b, m.win)
		if res != -1 {
			ans += res
		}
	}

	return ans
}

func findWin(a, b, win Pair) int {
	maxA := min(win.x/a.x, win.y/a.y) + 1
	maxB := min(win.x/b.x, win.y/b.y) + 1

	ans := -1
	for i := 0; i < maxA; i++ {
		for j := 0; j < maxB; j++ {
			resX := a.x*i + b.x*j
			resY := a.y*i + b.y*j

			if resX != win.x || resY != win.y {
				continue
			}
			tokens := 3*i + j

			if ans == -1 {
				ans = tokens
			} else {
				ans = min(ans, tokens)
			}
		}
	}

	return ans
}

func findWinSmart(a, b, win Pair) int {
	// AX*A + BX*B = WX
	// AY*A + BY*B = WY
	//
	// B = (WX - AX*A) / BX
	//
	// AY*A + BY*((WX - AX*A) / BX) = WY
	// AY*A + BY*(WX - AX*A)/BX = WY
	// AY*A + (BY*WX - BY*AX*A)/BX = WY
	// AY*A + BY*WX/BX - BY*AX*A/BX - WY = 0
	// AY*A - A*BY*AX/BX = WY - BY*WX/BX
	// A(AY - BY*AX/BX) = WY - BY*WX/BX
	// A = (WY - BY*WX/BX) / (AY - BY*AX/BX)
	// A = (WY*BX - BY*WX) / (AY*BX - BY*AX)

	pressA := (win.y*b.x - b.y*win.x) / (a.y*b.x - b.y*a.x)
	pressB := (win.x - a.x*pressA) / b.x

	isXCorrect := (a.x*pressA + b.x*pressB) == win.x
	isYCorrect := (a.y*pressA + b.y*pressB) == win.y
	if isXCorrect && isYCorrect {
		return pressA*3 + pressB
	} else {
		return -1
	}
}

type Pair struct {
	x, y int
}

type Machine struct {
	a, b, win Pair
}

func parseInput(input string) []Machine {
	blocks := strings.Split(input, "\n\n")

	machines := make([]Machine, len(blocks))
	for i, block := range blocks {
		nums := parseNumbers(block)
		a := Pair{nums[0], nums[1]}
		b := Pair{nums[2], nums[3]}
		win := Pair{nums[4], nums[5]}
		machines[i] = Machine{a, b, win}
	}

	return machines
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
