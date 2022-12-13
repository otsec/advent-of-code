package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"reflect"
	"sort"
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

	for i, packet := range parsed {
		left, right := packet[0], packet[1]
		if compare(left, right) != -1 {
			ans += i + 1
		}
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	divider1 := parseUnknownJson("[[2]]")
	divider2 := parseUnknownJson("[[6]]")

	var all []interface{}
	all = append(all, divider1, divider2)
	for _, packet := range parsed {
		left, right := packet[0], packet[1]
		all = append(all, left, right)
	}

	sort.SliceStable(all, func(i, j int) bool {
		return compare(all[i], all[j]) == 1
	})

	ans = 1
	for i, packet := range all {
		if compare(packet, divider1) == 0 || compare(packet, divider2) == 0 {
			ans *= i + 1
		}
	}

	return ans
}

func parseInput(input string) [][]interface{} {
	packets := strings.Split(input, "\n\n")

	ans := make([][]interface{}, len(packets))
	for i, packet := range packets {
		parts := strings.Split(packet, "\n")

		ans[i] = make([]interface{}, 2)
		ans[i][0] = parseUnknownJson(parts[0])
		ans[i][1] = parseUnknownJson(parts[1])
	}

	return ans
}

func parseUnknownJson(input string) (ans interface{}) {
	_ = json.Unmarshal([]byte(input), &ans)
	return ans
}

func compare(rawLeft, rawRight interface{}) int {
	// Compare simple ints
	// json.Unmarshal converts ints to float64 by default
	if reflect.TypeOf(rawLeft).Kind() == reflect.Float64 && reflect.TypeOf(rawRight).Kind() == reflect.Float64 {
		res := int(rawRight.(float64)) - int(rawLeft.(float64))
		if res > 0 {
			return 1
		} else if res == 0 {
			return 0
		} else if res < 0 {
			return -1
		}
	}

	left, right := ensureSlice(rawLeft), ensureSlice(rawRight)
	for i := 0; i < len(left); i++ {
		if i >= len(right) {
			return -1
		}

		res := compare(left[i], right[i])
		if res == 0 {
			continue
		} else {
			return res
		}
	}

	if len(right) > len(left) {
		return 1
	}

	return 0
}

func ensureSlice(input interface{}) (ans []interface{}) {
	if reflect.TypeOf(input).Kind() == reflect.Slice {
		ans = input.([]interface{})
	} else {
		ans = append(ans, input)
	}
	return ans
}
