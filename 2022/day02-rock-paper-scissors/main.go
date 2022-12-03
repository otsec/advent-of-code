package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/alexchao26/advent-of-code-go/util"
)

func main() {
	var part int
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.Parse()
	fmt.Println("Running part", part)

	if part == 1 {
		ans := part1(util.ReadFile("./input.txt"))
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(util.ReadFile("./input.txt"))
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) int {
	lines := strings.Split(input, "\n")

	total := 0
	for _, line := range lines {
		items := strings.Split(line, " ")
		var1, var2 := items[0], items[1]
		total += playForWin(var1, var2)
	}

	return total
}

func part2(input string) int {
	lines := strings.Split(input, "\n")

	total := 0
	for _, line := range lines {
		items := strings.Split(line, " ")
		var1, var2 := items[0], items[1]
		total += playSecretStrategy(var1, var2)
	}

	return total
}

type GameMapping struct {
	score int
	wins  string
	draws string
	loses string
}

func playForWin(var1, var2 string) int {
	result := 0

	mapping := map[string]GameMapping{
		"X": {score: 1, wins: "C", draws: "A"},
		"Y": {score: 2, wins: "A", draws: "B"},
		"Z": {score: 3, wins: "B", draws: "C"},
	}

	this := mapping[var2]

	result += this.score
	if var1 == this.wins {
		result += 6
	}
	if var1 == this.draws {
		result += 3
	}

	return result
}

func playSecretStrategy(var1, var2 string) int {
	result := 0

	scores := map[string]int{
		"A": 1, // rock
		"B": 2, // paper
		"C": 3, // scissors
	}
	mapping := map[string]GameMapping{
		"A": {wins: "C", draws: "A", loses: "B"},
		"B": {wins: "A", draws: "B", loses: "C"},
		"C": {wins: "B", draws: "C", loses: "A"},
	}

	if var2 == "X" {
		myChoice := mapping[var1].wins
		result += scores[myChoice]
	}
	if var2 == "Y" {
		myChoice := mapping[var1].draws
		result += scores[myChoice] + 3
	}
	if var2 == "Z" {
		myChoice := mapping[var1].loses
		result += scores[myChoice] + 6
	}

	return result
}
