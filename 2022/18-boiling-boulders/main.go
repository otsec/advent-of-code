package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"github.com/alexchao26/advent-of-code-go/util"
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
	pixels := parseInput(input)

	mapping := map[Pixel3D]bool{}
	for _, pixel := range pixels {
		mapping[pixel] = true
	}

	for _, pixel := range pixels {
		for _, connected := range getConnectedPixels(pixel) {
			if _, found := mapping[connected]; !found {
				ans++
			}
		}
	}

	return ans
}

func part2(input string) (ans int) {
	pixels := parseInput(input)

	pixelsMap := map[Pixel3D]bool{}
	for _, pixel := range pixels {
		pixelsMap[pixel] = true
	}

	grid := Grid3D{pixels[0], pixels[0]}
	for _, pixel := range pixels {
		grid.min.x = min(pixel.x-1, grid.min.x)
		grid.min.y = min(pixel.y-1, grid.min.y)
		grid.min.z = min(pixel.z-1, grid.min.z)
		grid.max.x = max(pixel.x+1, grid.max.x)
		grid.max.y = max(pixel.y+1, grid.max.y)
		grid.max.z = max(pixel.z+1, grid.max.z)
	}

	airMap := map[Pixel3D]bool{}
	extendAir(grid.min, grid, &pixelsMap, &airMap)

	for _, pixel := range pixels {
		for _, connected := range getConnectedPixels(pixel) {
			if _, found := airMap[connected]; found {
				ans++
			}
		}
	}

	return ans
}

type Pixel3D struct {
	x, y, z int
}

type Grid3D struct {
	min Pixel3D
	max Pixel3D
}

func getConnectedPixels(pixel Pixel3D) []Pixel3D {
	return []Pixel3D{
		{pixel.x + 1, pixel.y, pixel.z},
		{pixel.x - 1, pixel.y, pixel.z},
		{pixel.x, pixel.y + 1, pixel.z},
		{pixel.x, pixel.y - 1, pixel.z},
		{pixel.x, pixel.y, pixel.z + 1},
		{pixel.x, pixel.y, pixel.z - 1},
	}
}

func extendAir(pixel Pixel3D, grid Grid3D, pixelsMap *map[Pixel3D]bool, airMap *map[Pixel3D]bool) {
	if pixel.x < grid.min.x || pixel.x > grid.max.x {
		return
	}
	if pixel.y < grid.min.y || pixel.y > grid.max.y {
		return
	}
	if pixel.z < grid.min.z || pixel.z > grid.max.z {
		return
	}

	if _, found := (*airMap)[pixel]; found {
		return
	}
	if _, found := (*pixelsMap)[pixel]; found {
		return
	}

	(*airMap)[pixel] = true

	for _, connected := range getConnectedPixels(pixel) {
		extendAir(connected, grid, pixelsMap, airMap)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parseInput(input string) []Pixel3D {
	lines := strings.Split(input, "\n")

	ans := make([]Pixel3D, len(lines))
	for i, line := range lines {
		numbers := strings.Split(line, ",")
		x, y, z := cast.ToInt(numbers[0]), cast.ToInt(numbers[1]), cast.ToInt(numbers[2])
		ans[i] = Pixel3D{x, y, z}
	}

	return ans
}
