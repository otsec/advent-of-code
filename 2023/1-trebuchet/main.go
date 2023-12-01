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
	parsed := parseInput(input)

	re := regexp.MustCompile("[^0-9]")
	for _, line := range parsed {
		numbers := re.ReplaceAllString(line, "")
		n1, _ := strconv.Atoi(string(numbers[0]))
		n2, _ := strconv.Atoi(string(numbers[len(numbers)-1]))
		ans += n1*10 + n2
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	nums := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	words := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	for _, line := range parsed {
		numbers := ""
		for i, s := range line {
			for _, num := range nums {
				if string(s) == num {
					numbers += num
				}
			}

			for wi, word := range words {
				if i+len(word) > len(line) {
					continue
				}
				// fmt.Println(line, i, word, i+len(word), len(line), line[i:i+len(word)], word == line[i:i+len(word)])
				if line[i:i+len(word)] == word {
					numbers += nums[wi]
					break
				}
			}
		}

		// fmt.Println(line, numbers)
		n1, _ := strconv.Atoi(string(numbers[0]))
		n2, _ := strconv.Atoi(string(numbers[len(numbers)-1]))
		ans += n1*10 + n2
	}

	return ans
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}
