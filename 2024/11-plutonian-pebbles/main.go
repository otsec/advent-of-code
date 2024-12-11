package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"math"
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
	numbers := parseNumbers(input)
	memo := prepareMemo()
	return blinkSlice(memo, numbers, 25)
}

func part2(input string) (ans int) {
	numbers := parseNumbers(input)
	memo := prepareMemo()
	return blinkSlice(memo, numbers, 75)
}

func prepareMemo() *Memo {
	memo := &Memo{
		stones:    map[string][]int{},
		maxStones: 0,
		counts:    map[string]int{},
		maxCount:  0,
	}

	for n := 0; n < 10; n++ {
		nums := []int{n}
		for blinks := 1; blinks <= 5; blinks++ {
			nums = transformSlice(nums)
			memo.saveStones(n, blinks, nums)
		}
	}
	memo.maxStones = 5

	for blinks := 1; blinks <= 75; blinks++ {
		for n := 0; n < 10; n++ {
			res := blinkValue(memo, n, blinks)
			memo.saveCounts(n, blinks, res)
		}
		memo.maxCount = blinks
	}

	return memo
}

func blinkSlice(memo *Memo, numbers []int, times int) (ans int) {
	for _, n := range numbers {
		ans += blinkValue(memo, n, times)
	}

	return ans
}

func blinkValue(memo *Memo, num int, times int) (ans int) {
	if times < 0 {
		panic("times < 0")
	}

	if times == 0 {
		return 1
	}

	if times == 1 {
		if countOfDigits(num)%2 == 0 {
			return 2
		} else {
			return 1
		}
	}

	if num < 10 {
		if times <= memo.maxCount {
			return memo.readCounts(num, times)
		} else if times <= memo.maxStones {
			return len(memo.readStones(num, times))
		} else {
			nums := memo.readStones(num, memo.maxStones)
			return blinkSlice(memo, nums, times-memo.maxStones)
		}
	}

	return blinkSlice(memo, transformValue(num), times-1)
}

type Memo struct {
	stones    map[string][]int
	maxStones int
	counts    map[string]int
	maxCount  int
}

func (m *Memo) saveStones(num int, blinks int, stones []int) {
	key := fmt.Sprintf("%d,%d", num, blinks)
	m.stones[key] = stones
}

func (m *Memo) readStones(num int, blinks int) []int {
	key := fmt.Sprintf("%d,%d", num, blinks)
	if v, ok := m.stones[key]; ok {
		return v
	} else {
		panic(fmt.Sprintf("key %s not found", key))
	}
}

func (m *Memo) saveCounts(num int, blinks int, count int) {
	key := fmt.Sprintf("%d,%d", num, blinks)
	m.counts[key] = count
}

func (m *Memo) readCounts(num int, blinks int) int {
	key := fmt.Sprintf("%d,%d", num, blinks)
	if v, ok := m.counts[key]; ok {
		return v
	} else {
		panic(fmt.Sprintf("key %s not found", key))
	}
}

func transformSlice(nums []int) []int {
	var ans []int
	for _, n := range nums {
		ans = append(ans, transformValue(n)...)
	}
	return ans
}

func transformValue(initial int) []int {
	if initial == 0 {
		return []int{1}
	}

	digits := countOfDigits(initial)
	if digits%2 == 0 {
		mul := int(math.Pow10(digits / 2))
		n1 := initial / mul
		n2 := initial - n1*mul
		return []int{n1, n2}
	}

	return []int{initial * 2024}
}

func countOfDigits(n int) int {
	return int(math.Log10(float64(n))) + 1
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
