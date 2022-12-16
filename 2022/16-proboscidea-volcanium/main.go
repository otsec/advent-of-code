package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"github.com/alexchao26/advent-of-code-go/util"
	"github.com/emirpasic/gods/sets/hashset"
	"regexp"
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
	pressures, paths := parseInput(input)

	// run few time to ensure that all paths calculated correctly
	for i := 0; i < 5; i++ {
		paths = completePathsGraph(paths)
	}

	ans = findBestPathPressure("AA", 30, 10, hashset.New(), &pressures, &paths)

	return ans
}

func part2(input string) (ans int) {
	pressures, paths := parseInput(input)

	// run few time to ensure that all paths calculated correctly
	for i := 0; i < 5; i++ {
		paths = completePathsGraph(paths)
	}

	me := Walker{"AA", 26}
	elephant := Walker{"AA", 26}
	ans = findBestPathTwoWalkers(me, elephant, 8, hashset.New(), &pressures, &paths)

	return ans
}

type PressureMap map[string]int

type PathMap map[string]map[string]int

func parseInput(input string) (PressureMap, PathMap) {
	template := `Valve (.+?) has flow rate=(.+?); tunnels? leads? to valves? (.+)`
	r := regexp.MustCompile(template)
	results := r.FindAllStringSubmatch(input, -1)

	pressures := make(PressureMap)
	paths := make(PathMap)

	for _, parsed := range results {
		name := parsed[1]
		pressure := cast.ToInt(parsed[2])
		tunnels := strings.Split(parsed[3], ", ")

		pressures[name] = pressure

		paths[name] = map[string]int{}
		for _, tunnel := range tunnels {
			paths[name][tunnel] = 1
		}
	}

	return pressures, paths
}

func completePathsGraph(paths PathMap) PathMap {
	for start, directions := range paths {
		for oldEnd, oldSteps := range directions {
			for newEnd, newSteps := range paths[oldEnd] {
				if newEnd == start {
					continue
				}

				newVal := oldSteps + newSteps

				if currVal, found := paths[start][newEnd]; !found || currVal > newVal {
					paths[start][newEnd] = newVal
				}

				if currVal, found := paths[newEnd][start]; !found || currVal > newVal {
					paths[start][newEnd] = newVal
				}
			}
		}
	}

	return paths
}

type Variant struct {
	valve         string
	pressureRate  int
	minutesToOpen int
}

func findBestVariants(currRoom string, minutesLeft int, openValves *hashset.Set, pressures *PressureMap, paths *PathMap) []Variant {
	variants := []Variant{}
	for valve, walkMinutes := range (*paths)[currRoom] {
		if openValves.Contains(valve) {
			continue
		}

		minutesToOpen := walkMinutes + 1
		if minutesToOpen > minutesLeft {
			continue
		}

		pressureRate := (*pressures)[valve] * (minutesLeft - minutesToOpen)
		if pressureRate == 0 {
			continue
		}

		variants = append(variants, Variant{valve, pressureRate, minutesToOpen})
	}

	sort.Slice(variants, func(i, j int) bool {
		return variants[i].pressureRate > variants[j].pressureRate
	})

	return variants
}

func findBestPathPressure(currRoom string, minutesLeft int, maxVariants int, openValves *hashset.Set, pressures *PressureMap, paths *PathMap) (maxPressureRate int) {
	variants := findBestVariants(currRoom, minutesLeft, openValves, pressures, paths)
	if len(variants) == 0 {
		return maxPressureRate
	}

	for i, variant := range variants {
		if i+1 > maxVariants {
			break
		}

		subOpenValves := hashset.New()
		subOpenValves.Add(openValves.Values()...)
		subOpenValves.Add(variant.valve)

		subPressureRate := findBestPathPressure(variant.valve, minutesLeft-variant.minutesToOpen, maxVariants, subOpenValves, pressures, paths)
		if subPressureRate+variant.pressureRate > maxPressureRate {
			maxPressureRate = subPressureRate + variant.pressureRate
		}
	}

	return maxPressureRate
}

type Walker struct {
	room    string
	minutes int
}

func findBestPathTwoWalkers(me Walker, elephant Walker, maxVariants int, openValves *hashset.Set, pressures *PressureMap, paths *PathMap) (maxPressureRate int) {
	var curr, other *Walker
	if me.minutes >= elephant.minutes {
		curr, other = &me, &elephant
	} else {
		curr, other = &elephant, &me
	}

	variants := findBestVariants(curr.room, curr.minutes, openValves, pressures, paths)
	if len(variants) == 0 {
		return maxPressureRate
	}

	for i, variant := range variants {
		if i+1 > maxVariants {
			break
		}

		subOpenValves := hashset.New()
		subOpenValves.Add(openValves.Values()...)
		subOpenValves.Add(variant.valve)

		newCurr := Walker{
			variant.valve,
			curr.minutes - variant.minutesToOpen,
		}

		subPressureRate := findBestPathTwoWalkers(newCurr, *other, maxVariants, subOpenValves, pressures, paths)
		if subPressureRate+variant.pressureRate > maxPressureRate {
			maxPressureRate = subPressureRate + variant.pressureRate
		}
	}

	return maxPressureRate
}
