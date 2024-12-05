package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"regexp"
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

type OrderingRule [2]int
type Page []int

func part1(input string) (ans int) {
	rules, pages := parseInput(input)

	for _, page := range pages {
		if isPageProperlyOrdered(page, rules) {
			middle := len(page) / 2
			ans += page[middle]
		}
	}

	return ans
}

func part2(input string) (ans int) {
	rules, pages := parseInput(input)

	for _, page := range pages {
		if !isPageProperlyOrdered(page, rules) {
			sortPage(page, rules)
			middle := len(page) / 2
			ans += page[middle]
		}
	}

	return ans
}

func isPageProperlyOrdered(page Page, rules []OrderingRule) bool {
	for _, rule := range rules {
		index1 := slices.Index(page, rule[0])
		index2 := slices.Index(page, rule[1])
		if index1 != -1 && index2 != -1 && index1 > index2 {
			return false
		}
	}

	return true
}

func sortPage(page Page, rules []OrderingRule) {
	slices.SortFunc(page, func(a, b int) int {
		for _, rule := range rules {
			if rule[0] == a && rule[1] == b {
				return -1
			}
			if rule[0] == b && rule[1] == a {
				return 1
			}
		}
		return 0
	})
}

func parseInput(input string) ([]OrderingRule, []Page) {
	segments := strings.Split(input, "\n\n")

	rulesLines := strings.Split(segments[0], "\n")
	rulesParsed := make([]OrderingRule, len(rulesLines))
	for i, line := range rulesLines {
		numsRaw := strings.Split(line, "|")
		v1, _ := strconv.Atoi(numsRaw[0])
		v2, _ := strconv.Atoi(numsRaw[1])
		rulesParsed[i] = [2]int{v1, v2}
	}

	pagesLines := strings.Split(segments[1], "\n")
	pagesParsed := make([]Page, len(pagesLines))
	for i, line := range pagesLines {
		pagesParsed[i] = parseNumbers(line)
	}

	return rulesParsed, pagesParsed
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
