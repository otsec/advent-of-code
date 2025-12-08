package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"sort"
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
	junctionBoxes := parseInput(input)
	distancesMemo := calculateDistances(junctionBoxes)

	// 1k for main input and 10 for test input
	connections := 1000
	if len(junctionBoxes) == 20 {
		connections = 10
	}

	circuits := NewCircuitsMemo(junctionBoxes)
	for i := 0; i < connections; i++ {
		p1, p2 := distancesMemo[i].p1, distancesMemo[i].p2
		circuits.Connect(p1, p2)
	}

	top := circuits.TopCircuits()
	ans = top[0] * top[1] * top[2]

	return ans
}

func part2(input string) (ans int) {
	junctionBoxes := parseInput(input)
	distancesMemo := calculateDistances(junctionBoxes)

	circuits := NewCircuitsMemo(junctionBoxes)
	for i := 0; i < len(distancesMemo); i++ {
		p1, p2 := distancesMemo[i].p1, distancesMemo[i].p2
		circuits.Connect(p1, p2)
		if circuits.Count() == 1 {
			ans = p1.x * p2.x
			break
		}
	}

	return ans
}

type Point3D struct {
	x, y, z int
}

func distance(p1, p2 Point3D) float64 {
	dx := float64(p2.x - p1.x)
	dy := float64(p2.y - p1.y)
	dz := float64(p2.z - p1.z)

	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

type DistMemoItem struct {
	p1, p2   Point3D
	distance float64
}

func calculateDistances(junctionBoxes []Point3D) []DistMemoItem {
	memo := make([]DistMemoItem, 0, len(junctionBoxes)*len(junctionBoxes))
	for i := 0; i < len(junctionBoxes); i++ {
		for j := i + 1; j < len(junctionBoxes); j++ {
			memo = append(memo, DistMemoItem{junctionBoxes[i], junctionBoxes[j], distance(junctionBoxes[i], junctionBoxes[j])})
		}
	}
	sort.Slice(memo, func(i, j int) bool {
		return memo[i].distance < memo[j].distance
	})
	return memo
}

func NewCircuitsMemo(junctionBoxes []Point3D) *CircuitsMemo {
	return &CircuitsMemo{make(map[Point3D]int), 1, len(junctionBoxes)}
}

type CircuitsMemo struct {
	idMap         map[Point3D]int
	nextCircuitId int
	counter       int
}

func (cm *CircuitsMemo) Connect(p1, p2 Point3D) {
	c1 := cm.idMap[p1]
	c2 := cm.idMap[p2]

	if c1 == 0 && c2 == 0 {
		cm.idMap[p1] = cm.nextCircuitId
		cm.idMap[p2] = cm.nextCircuitId
		cm.nextCircuitId++
		cm.counter--
		return
	}

	if c1 == 0 && c2 != 0 {
		cm.idMap[p1] = c2
		cm.counter--
		return
	}

	if c1 != 0 && c2 == 0 {
		cm.idMap[p2] = c1
		cm.counter--
		return
	}

	if c1 != 0 && c2 != 0 && c1 != c2 {
		for k := range cm.idMap {
			if cm.idMap[k] == c2 {
				cm.idMap[k] = c1
			}
		}
		cm.counter--
	}
}

func (cm *CircuitsMemo) TopCircuits() []int {
	counter := make([]int, cm.nextCircuitId)
	for _, id := range cm.idMap {
		counter[id]++
	}
	sort.Slice(counter, func(i, j int) bool {
		return counter[i] > counter[j]
	})
	return counter
}

func (cm *CircuitsMemo) Count() int {
	return cm.counter
}

func parseInput(input string) []Point3D {
	lines := strings.Split(input, "\n")
	inputs := make([]Point3D, len(lines))
	for i, line := range lines {
		numbers := parseNumbers(line)
		inputs[i] = Point3D{numbers[0], numbers[1], numbers[2]}
	}
	return inputs
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
