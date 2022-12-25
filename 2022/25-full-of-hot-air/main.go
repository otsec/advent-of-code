package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"log"
	"math"
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

func part1(input string) string {
	lines := parseInput(input)

	ans := 0
	for _, line := range lines {
		ans += SnafuToHex(line)
	}

	return HexToSnafu(ans)
}

func part2(input string) (ans int) {
	parsed := parseInput(input)
	_ = parsed

	return ans
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}

func SnafuToHex(snafu string) (hex int) {
	values := map[rune]int{
		'=': -2,
		'-': -1,
		'0': 0,
		'1': 1,
		'2': 2,
	}

	for i := 0; i < len(snafu); i++ {
		symbol := rune(snafu[len(snafu)-1-i])
		val, found := values[symbol]
		if !found {
			log.Panicf("Unknown symbol %v in %v", symbol, snafu)
		}

		multiplier := PowInt(5, i)
		hex += val * multiplier
	}

	return hex
}

func HexToSnafu(hex int) (snafu string) {
	acc := hex
	digits := []int{}
	index := 0
	for acc > 0 {
		for index >= len(digits) {
			digits = append(digits, 0)
		}

		digits[index] += acc % 5
		if digits[index] > 2 {
			digits[index] -= 5
			digits = append(digits, 1)
		}

		acc = acc / 5
		index++
	}

	symbols := map[int]string{
		-2: "=",
		-1: "-",
		0:  "0",
		1:  "1",
		2:  "2",
	}

	for _, digit := range digits {
		val, found := symbols[digit]
		if !found {
			log.Panicf("Unknown digit %v in digits %v in hex %v", digit, digits, input)
		}

		snafu = val + snafu
	}

	return snafu
}

func PowInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
