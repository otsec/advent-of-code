package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/alexchao26/advent-of-code-go/util"
)

//go:embed input.txt
var input string

func main() {
	var part int
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.Parse()
	fmt.Println("Running part", part)

	var ans int
	if part == 1 {
		ans = part1(input)
	} else {
		ans = part2(input)
	}

	_ = util.CopyToClipboard(fmt.Sprintf("%v", ans))
	fmt.Println("Output:", ans)
}

func part1(input string) (result int) {
	rucksacks := strings.Split(input, "\n")

	for _, rucksack := range rucksacks {
		firstCompartment, secondCompartment := splitRucksack(rucksack)

		mapping := map[rune]bool{}
		for _, item := range firstCompartment {
			mapping[item] = true
		}

		for _, item := range secondCompartment {
			if _, exists := mapping[item]; exists {
				result += priorityOfItem(item)
				break
			}
		}
	}

	return result
}

func part2(input string) (result int) {
	rucksacks := strings.Split(input, "\n")

	for i := 0; i < len(rucksacks)/3; i++ {
		rucksack1 := rucksacks[i*3]
		mapping1 := map[rune]bool{}
		for _, item := range rucksack1 {
			mapping1[item] = true
		}

		rucksack2 := rucksacks[i*3+1]
		mapping2 := map[rune]bool{}
		for _, item := range rucksack2 {
			if _, exists := mapping1[item]; exists {
				mapping2[item] = true
			}
		}

		rucksack3 := rucksacks[i*3+2]
		for _, item := range rucksack3 {
			if _, exists := mapping2[item]; exists {
				result += priorityOfItem(item)
				break
			}
		}
	}

	return result
}

func splitRucksack(input string) (string, string) {
	return input[:len(input)/2], input[len(input)/2:]
}

func priorityOfItem(item rune) int {
	if item >= 'a' && item <= 'z' {
		return int(item) - 96
	}
	if item >= 'A' && item <= 'Z' {
		return int(item) - 38
	}
	return int(item)
}
