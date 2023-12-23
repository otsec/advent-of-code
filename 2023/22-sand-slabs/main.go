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

func part1(input string) (ans int) {
	bricks := parseInput(input)

	moved := true
	for moved {
		bricks, moved = moveDown(bricks)
	}

	for i, _ := range bricks {
		var bricksWithoutOne []Brick
		for j, _ := range bricks {
			if i != j {
				bricksWithoutOne = append(bricksWithoutOne, bricks[j])
			}
		}

		_, movedWithoutOne := moveDown(bricksWithoutOne)
		if !movedWithoutOne {
			ans += 1
		}
	}

	//print(bricks)

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)
	_ = parsed

	return ans
}

type Coord struct {
	x, y, z int
}

type Brick struct {
	start, end Coord
}

func parseInput(input string) []Brick {
	lines := strings.Split(input, "\n")
	bricks := make([]Brick, len(lines))
	for i, line := range lines {
		nums := parseNumbers(line)
		bricks[i] = Brick{Coord{nums[0], nums[1], nums[2]}, Coord{nums[3], nums[4], nums[5]}}
	}
	return bricks
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

func moveDown(from []Brick) ([]Brick, bool) {
	var to []Brick
	var moved bool

	for _, brick := range from {
		// fmt.Println("currentState", to)

		if brick.start.z == 1 {
			to = append(to, brick)
			continue
		}

		nextStart := Coord{brick.start.x, brick.start.y, brick.start.z - 1}
		nextEnd := Coord{brick.end.x, brick.end.y, brick.end.z - 1}
		nextBrick := Brick{nextStart, nextEnd}

		if isInterfereAny(nextBrick, to) {
			to = append(to, brick)
		} else {
			to = append(to, nextBrick)
			moved = true
		}
	}

	return to, moved
}

func isInterfereAny(brick Brick, existed []Brick) bool {
	for _, anotherBrick := range existed {
		// fmt.Println("isInterfereOne", brick, anotherBrick, isInterfereOne(brick, anotherBrick))
		if isInterfereOne(brick, anotherBrick) {
			return true
		}
	}

	return false
}

func isInterfereOne(a Brick, b Brick) bool {
	if a.start.z > b.end.z || a.end.z < b.start.z {
		return false
	}

	if a.start.x == a.end.x && b.start.x == b.end.x {
		return a.start.x == b.start.x
	}

	pointsA := getAllPoints(a)
	for _, cb := range getAllPoints(b) {
		if slices.Contains(pointsA, cb) {
			return true
		}

		//for _, ca := range pointsA {
		//	if reflect.DeepEqual(ca, cb) {
		//		return true
		//	}
		//}
	}

	return false
}

func getAllPoints(brick Brick) []Coord {
	ans := []Coord{}
	for x := brick.start.x; x <= brick.end.x; x++ {
		ans = append(ans, Coord{x, brick.start.y, brick.start.z})
	}
	for y := brick.start.y; y <= brick.end.y; y++ {
		ans = append(ans, Coord{brick.start.x, y, brick.start.z})
	}
	return ans
}

func print(bricks []Brick) {
	var maxX, maxY, maxZ int
	for _, brick := range bricks {
		maxX = max(maxX, brick.start.x, brick.end.x)
		maxY = max(maxY, brick.start.y, brick.end.y)
		maxZ = max(maxZ, brick.start.z, brick.end.z)
	}

	for z := maxZ; z > 0; z-- {
		for x := 0; x <= maxX; x++ {
			symbol := "."

			for i, brick := range bricks {
				if brick.start.z <= z && z <= brick.end.z {
					if brick.start.x <= x && x <= brick.end.x {
						symbol = strconv.Itoa(i)
						break
					}
				}
			}

			fmt.Print(symbol)
		}

		fmt.Print("  ")

		for y := 0; y <= maxY; y++ {
			symbol := "."

			for i, brick := range bricks {
				if brick.start.z <= z && z <= brick.end.z {
					if brick.start.y <= y && y <= brick.end.y {
						symbol = strconv.Itoa(i)
						break
					}
				}
			}

			fmt.Print(symbol)
		}

		fmt.Println()
	}

	for x := 0; x <= maxX; x++ {
		fmt.Print("-")
	}
	fmt.Print("  ")
	for y := 0; y <= maxY; y++ {
		fmt.Print("-")
	}

	fmt.Println()
}
