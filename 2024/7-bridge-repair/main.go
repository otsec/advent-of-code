package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
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
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) (ans int) {
	inputs := parseInput(input)

	for _, calibration := range inputs {
		if isPossibleP1(calibration) {
			ans += calibration.res
		}
	}

	return ans
}

func part2(input string) (ans int) {
	inputs := parseInput(input)

	for _, calibration := range inputs {
		if isPossibleP2(calibration) {
			ans += calibration.res
		}
	}

	return ans
}

func isPossibleP1(input Input) bool {
	memo := make([][]int, len(input.numbers))

	memo[0] = []int{input.numbers[0]}
	for i := 1; i < len(input.numbers); i++ {
		curr := input.numbers[i]
		isLast := i == len(input.numbers)-1
		for _, prev := range memo[i-1] {
			v1 := prev + curr
			v2 := prev * curr

			if isLast {
				if v1 == input.res || v2 == input.res {
					return true
				}
			} else {
				memo[i] = append(memo[i], v1, v2)
			}
		}
	}

	return false
}

func isPossibleP2(input Input) bool {
	memo := make([][]int, len(input.numbers))

	memo[0] = []int{input.numbers[0]}
	for i := 1; i < len(input.numbers); i++ {
		curr := input.numbers[i]
		isLast := i == len(input.numbers)-1
		for _, prev := range memo[i-1] {
			v1 := prev + curr
			v2 := prev * curr
			v3 := concat(prev, curr)

			if isLast {
				if v1 == input.res || v2 == input.res || v3 == input.res {
					return true
				}
			} else {
				memo[i] = append(memo[i], v1, v2, v3)
			}
		}
	}

	return false
}

func concat(n1, n2 int) int {
	s1 := strconv.Itoa(n1)
	s2 := strconv.Itoa(n2)
	ans, _ := strconv.Atoi(s1 + s2)
	return ans
}

type Input struct {
	res     int
	numbers []int
}

func parseInput(input string) []Input {
	lines := strings.Split(input, "\n")
	inputs := make([]Input, len(lines))
	for i, line := range lines {
		numbers := parseNumbers(line)
		inputs[i] = Input{numbers[0], numbers[1:]}
	}
	return inputs
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
