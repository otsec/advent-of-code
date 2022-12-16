package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"github.com/alexchao26/advent-of-code-go/util"
	"math"
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
		ans := part1(input, 2000000)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input, 4000000)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string, magicY int) (ans int) {
	pairs := parseInput(input)

	scannedX := []Line{}
	for _, pair := range pairs {
		maxSignalDistance := calcManhattanDistance(pair.sensor, pair.beacon)
		line := calcCoverageAt(magicY, pair.sensor, maxSignalDistance)
		scannedX = mergeLines(scannedX, line)
	}

	for _, line := range scannedX {
		ans += line.end - line.start
	}

	return ans
}

func part2(input string, maxY int) (ans int) {
	pairs := parseInput(input)

	grid := make([][]Line, maxY+1)
	for _, pair := range pairs {
		maxSignalDistance := calcManhattanDistance(pair.sensor, pair.beacon)
		for y := 0; y <= maxY; y++ {
			line := calcCoverageAt(y, pair.sensor, maxSignalDistance)
			if line.start < 0 {
				line.start = 0
			}
			if line.end > maxY {
				line.end = maxY
			}
			if line.start == 0 && line.end == 0 {
				continue
			}

			grid[y] = mergeLines(grid[y], line)
		}
	}

	for y, lines := range grid {
		if len(lines) == 1 && lines[0].start == 0 && lines[0].end == maxY {
			continue
		}

		if len(lines) == 2 {
			x := lines[0].end + 1
			ans = x*4000000 + y
			break
		}
	}

	return ans
}

type Coordinate struct {
	x int
	y int
}

type Pair struct {
	sensor Coordinate
	beacon Coordinate
}

type Line struct {
	start int
	end   int
}

func parseInput(input string) []Pair {
	template := `Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`
	r := regexp.MustCompile(template)
	results := r.FindAllStringSubmatch(input, -1)

	ans := make([]Pair, len(results))
	for i, parsed := range results {
		sx, sy := cast.ToInt(parsed[1]), cast.ToInt(parsed[2])
		bx, by := cast.ToInt(parsed[3]), cast.ToInt(parsed[4])
		ans[i] = Pair{Coordinate{sx, sy}, Coordinate{bx, by}}
	}

	return ans
}

func calcManhattanDistance(sensor, beacon Coordinate) int {
	diffX := abs(sensor.x - beacon.x)
	diffY := abs(sensor.y - beacon.y)
	return diffX + diffY
}
func mergeLines(lines []Line, new Line) []Line {
	if len(lines) == 0 {
		return []Line{new}
	}

	if len(lines) == 1 {
		if lines[0].start <= new.start && new.end <= lines[0].end {
			return lines
		}
	}

	if new.end+1 < lines[0].start {
		return append([]Line{new}, lines...)
	}

	if lines[len(lines)-1].end+1 < new.start {
		return append(lines, new)
	}

	lines = append(lines, new)
	sort.Slice(lines, func(i, j int) bool {
		return lines[i].start < lines[j].start
	})

	merged := []Line{}
	for _, line := range lines {
		if len(merged) == 0 {
			merged = append(merged, line)
		} else if merged[len(merged)-1].end+1 < line.start {
			merged = append(merged, line)
		} else {
			merged[len(merged)-1].end = max(merged[len(merged)-1].end, line.end)
		}
	}

	return merged
}

func calcCoverageAt(y int, sensor Coordinate, maxSignalDistance int) Line {
	distanceToY := abs(sensor.y - y)

	if distanceToY > maxSignalDistance {
		return Line{0, 0}
	}

	startX := sensor.x - (maxSignalDistance - distanceToY)
	endX := sensor.x + (maxSignalDistance - distanceToY)

	return Line{startX, endX}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func abs(a int) int {
	return int(math.Abs(float64(a)))
}
