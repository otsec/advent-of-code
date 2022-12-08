package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"strings"

	"github.com/alexchao26/advent-of-code-go/util"
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
	treeMap := parseInput(input)

	//visibilityMap := createVisibilityMap(treeMap)
	//printMap(treeMap, visibilityMap)

	for i := 0; i < len(treeMap); i++ {
		for j := 0; j < len(treeMap[0]); j++ {
			if isVisible(treeMap, i, j) {
				ans++
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	treeMap := parseInput(input)

	//scoreMap := createScoreMap(treeMap)
	//printMap(treeMap, scoreMap)

	for i := 0; i < len(treeMap); i++ {
		for j := 0; j < len(treeMap[0]); j++ {
			score := scoreFor(treeMap, i, j)
			if score > ans {
				ans = score
			}
		}
	}

	return ans
}

func printMap(treeMap, extraMap [][]int) {
	for i := 0; i < len(treeMap); i++ {
		fmt.Print(treeMap[i], "  ", extraMap[i], "\n")
	}
}

func createVisibilityMap(treeMap [][]int) [][]int {
	visibleMap := make([][]int, len(treeMap))
	for i := 0; i < len(treeMap); i++ {
		visibleMap[i] = make([]int, len(treeMap[0]))

		for j := 0; j < len(treeMap[0]); j++ {
			if isVisible(treeMap, i, j) == false {
				visibleMap[i][j] = 1
			}
		}
	}

	return visibleMap
}

func createScoreMap(treeMap [][]int) [][]int {
	scoreMap := make([][]int, len(treeMap))
	for i := 0; i < len(treeMap); i++ {
		scoreMap[i] = make([]int, len(treeMap[0]))
		for j := 0; j < len(treeMap[0]); j++ {
			scoreMap[i][j] = scoreFor(treeMap, i, j)
		}
	}

	return scoreMap
}

func isVisible(treeMap [][]int, iTarget int, jTarget int) bool {
	visibleUp := true
	for i := 0; i < iTarget; i++ {
		if treeMap[i][jTarget] >= treeMap[iTarget][jTarget] {
			visibleUp = false
			break
		}
	}
	if visibleUp {
		return true
	}

	visibleRight := true
	for j := jTarget + 1; j < len(treeMap[0]); j++ {
		if treeMap[iTarget][j] >= treeMap[iTarget][jTarget] {
			visibleRight = false
			break
		}
	}
	if visibleRight {
		return true
	}

	visibleDown := true
	for i := iTarget + 1; i < len(treeMap); i++ {
		if treeMap[i][jTarget] >= treeMap[iTarget][jTarget] {
			visibleDown = false
			break
		}
	}
	if visibleDown {
		return true
	}

	visibleLeft := true
	for j := 0; j < jTarget; j++ {
		if treeMap[iTarget][j] >= treeMap[iTarget][jTarget] {
			visibleLeft = false
			break
		}
	}
	if visibleLeft {
		return true
	}

	return false
}

func scoreFor(treeMap [][]int, iTarget int, jTarget int) int {
	scoreUp := 0
	for i := iTarget - 1; i >= 0; i-- {
		scoreUp++
		if treeMap[i][jTarget] >= treeMap[iTarget][jTarget] {
			break
		}
	}

	scoreRight := 0
	for j := jTarget + 1; j < len(treeMap[0]); j++ {
		scoreRight++
		if treeMap[iTarget][j] >= treeMap[iTarget][jTarget] {
			break
		}
	}

	scoreDown := 0
	for i := iTarget + 1; i < len(treeMap); i++ {
		scoreDown++
		if treeMap[i][jTarget] >= treeMap[iTarget][jTarget] {
			break
		}
	}

	scoreLeft := 0
	for j := jTarget - 1; j >= 0; j-- {
		scoreLeft++
		if treeMap[iTarget][j] >= treeMap[iTarget][jTarget] {
			break
		}
	}

	return scoreUp * scoreRight * scoreDown * scoreLeft
}

func parseInput(input string) [][]int {
	lines := strings.Split(input, "\n")

	ans := [][]int{}
	for _, line := range lines {
		row := []int{}
		for _, char := range line {
			row = append(row, cast.ToInt(string(char)))
		}
		ans = append(ans, row)
	}

	return ans
}
