package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
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
	monkeys := parseInput(input)

	ans = compute(monkeys, monkeys["root"])

	return ans
}

func part2(input string) (ans int) {
	monkeys := parseInput(input)

	root, _ := monkeys["root"]
	val1 := compute(monkeys, monkeys[root.parts[0]])
	val2 := compute(monkeys, monkeys[root.parts[1]])
	fmt.Println(monkeys["humn"], val1, val2)

	// I've checked manually that val1 is always changing and val2 always stays the same
	// fmt.Println(val1, val2)

	humMin := 0
	humMax := 40000000000000

	monkeys["humn"] = Job{_type: JOB_CONST, num: humMin}
	valHumMin := compute(monkeys, monkeys[root.parts[0]])

	monkeys["humn"] = Job{_type: JOB_CONST, num: humMax}
	valHumMax := compute(monkeys, monkeys[root.parts[0]])

	reqVal := val2
	if (valHumMin < reqVal && reqVal < valHumMax) || (valHumMax < reqVal && reqVal < valHumMin) {
		// ok
	} else {
		panic("Selected humMin and humMax are not wide enough to find reqVal.")
	}

	l := -40000000000000
	r := 40000000000000
	for l < r {
		m := (r + l) / 2

		monkeys["humn"] = Job{_type: JOB_CONST, num: m}
		val1 = compute(monkeys, monkeys[root.parts[0]])

		if val1 == val2 {
			ans = m
			break
		} else if (val1 < val2 && valHumMin < valHumMax) || (val1 > val2 && valHumMin > valHumMax) {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	return ans
}

func compute(monkeys map[string]Job, job Job) int {
	if job._type == JOB_CONST {
		return job.num
	}

	if job._type == JOB_EQUASION {
		val1 := compute(monkeys, monkeys[job.parts[0]])
		val2 := compute(monkeys, monkeys[job.parts[1]])
		switch job.op {
		case "+":
			return val1 + val2
		case "-":
			return val1 - val2
		case "*":
			return val1 * val2
		case "/":
			return val1 / val2
		}
	}

	panic(fmt.Sprintf("Cannot compute the job: %v.", job))
}

type JobType string

const (
	JOB_EQUASION JobType = "eq"
	JOB_CONST    JobType = "const"
)

type Job struct {
	_type JobType
	num   int
	op    string
	parts []string
}

func parseInput(input string) map[string]Job {
	lines := strings.Split(input, "\n")

	monkeys := make(map[string]Job)
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		name, eq := parts[0], parts[1]

		if val, err := strconv.Atoi(eq); err == nil {
			job := Job{_type: JOB_CONST, num: val}
			monkeys[name] = job
			continue
		}

		operations := []string{"+", "-", "*", "/"}
		for _, op := range operations {
			separator := fmt.Sprintf(" %s ", op)
			if strings.Contains(eq, separator) {
				parts2 := strings.Split(eq, separator)
				job := Job{_type: JOB_EQUASION, op: op, parts: parts2}
				monkeys[name] = job
				break
			}
		}

		if _, found := monkeys[name]; !found {
			panic(fmt.Sprintf("Monkey %v error. Cannot parse the equasion: %v.", name, eq))
		}
	}

	return monkeys
}
