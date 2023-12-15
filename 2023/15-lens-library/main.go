package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"slices"
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
	items := strings.Split(input, ",")

	for _, item := range items {
		ans += makeHash(item)
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	hashmap := make([][]InputItem, 256)
	for _, item := range parsed {
		box := makeHash(item.label)

		index := slices.IndexFunc(hashmap[box], func(saved InputItem) bool {
			return saved.label == item.label
		})

		if index == -1 {
			if item.op == "=" {
				hashmap[box] = append(hashmap[box], item)
			}
		} else {
			if item.op == "=" {
				hashmap[box][index] = item
			} else {
				hashmap[box] = slices.Delete(hashmap[box], index, index+1)
			}
		}
	}

	for box, _ := range hashmap {
		for slot, item := range hashmap[box] {
			ans += (box + 1) * (slot + 1) * item.focal
		}
	}

	return ans
}

type InputItem struct {
	raw   string
	label string
	op    string
	focal int
}

func parseInput(input string) []InputItem {
	items := strings.Split(input, ",")

	inputs := make([]InputItem, len(items))
	for i, item := range items {
		if strings.Contains(item, "-") {
			segments := strings.Split(item, "-")
			inputs[i] = InputItem{item, segments[0], "-", 0}
		} else {
			segments := strings.Split(item, "=")
			val, _ := strconv.Atoi(segments[1])
			inputs[i] = InputItem{item, segments[0], "=", val}
		}
	}

	return inputs
}

func makeHash(input string) (ans int) {
	for i := 0; i < len(input); i++ {
		ans += int(input[i])
		ans = ans * 17 % 256
	}

	return ans
}
