package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"math/big"
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
		ans := part1(input, 200000000000000, 400000000000000)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string, min, max int) (ans int) {
	parsed := parseInput(input)

	minBig := big.NewFloat(float64(min))
	maxBig := big.NewFloat(float64(max))

	for i := 0; i < len(parsed)-1; i++ {
		for j := i + 1; j < len(parsed); j++ {
			c1 := new(HailstoneBig).FromHailstone(parsed[i])
			c2 := new(HailstoneBig).FromHailstone(parsed[j])

			p1 := c1.pos
			p2 := c1.AtSec(30)
			p3 := c2.pos
			p4 := c2.AtSec(30)

			px, err := findIntersection(p1, p2, p3, p4)

			//fmt.Printf("Hailstone A: {%v, %v}, {%v, %v}\n", c1.pos.x, c1.pos.y, c1.speed.x, c1.speed.y)
			//fmt.Printf("Hailstone B: {%v, %v}, {%v, %v}\n", c2.pos.x, c2.pos.y, c2.speed.x, c2.speed.y)
			//fmt.Printf("Intersection: {%v, %v}, %v\n", px.x, px.y, err)

			if err != nil {
				continue
			}
			if px.x.Cmp(minBig) == -1 || px.x.Cmp(maxBig) == 1 {
				continue
			}
			if px.y.Cmp(minBig) == -1 || px.y.Cmp(maxBig) == 1 {
				continue
			}

			//t1 := (px.X - float64(c1.pos.x)) / float64(c1.speed.x)
			//t2 := (px.Y - float64(c1.pos.y)) / float64(c1.speed.y)
			//t3 := (px.X - float64(c2.pos.x)) / float64(c2.speed.x)
			//t4 := (px.Y - float64(c2.pos.y)) / float64(c2.speed.y)
			t1 := new(big.Float).Quo(new(big.Float).Sub(px.x, c1.pos.x), c1.speed.x)
			t3 := new(big.Float).Quo(new(big.Float).Sub(px.x, c2.pos.x), c2.speed.x)
			//fmt.Printf("%v, %v, %v, %v \n", t1, t2, t3, t4)

			zero := big.NewFloat(float64(0))
			if t1.Cmp(zero) == -1 || t3.Cmp(zero) == -1 {
				continue
			}

			ans += 1
		}
	}

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

type Hailstone struct {
	pos, speed Coord
}

func (h *Hailstone) AtSec(t int) Coord {
	x := h.pos.x + h.speed.x*t
	y := h.pos.y + h.speed.y*t
	z := h.pos.z + h.speed.z*t
	return Coord{x, y, z}
}

func parseInput(input string) []Hailstone {
	lines := strings.Split(input, "\n")
	ans := make([]Hailstone, len(lines))
	for i, line := range lines {
		nums := parseNumbers(line)
		pos := Coord{nums[0], nums[1], nums[2]}
		speed := Coord{nums[3], nums[4], nums[5]}
		ans[i] = Hailstone{pos, speed}
	}
	return ans
}

func parseNumbers(input string) []int {
	re := regexp.MustCompile(`[-0-9]+`)
	matches := re.FindAllString(input, -1)
	nums := make([]int, len(matches))
	for i, match := range matches {
		nums[i], _ = strconv.Atoi(match)
	}
	return nums
}

type CoordBig struct {
	x, y, z *big.Float
}

func (cb *CoordBig) FromCoord(c Coord) *CoordBig {
	x := big.NewFloat(float64(c.x))
	y := big.NewFloat(float64(c.y))
	z := big.NewFloat(float64(c.z))
	return &CoordBig{x, y, z}
}

type HailstoneBig struct {
	pos   *CoordBig
	speed *CoordBig
}

func (hb *HailstoneBig) FromHailstone(h Hailstone) *HailstoneBig {
	pos := new(CoordBig).FromCoord(h.pos)
	speed := new(CoordBig).FromCoord(h.speed)
	return &HailstoneBig{pos, speed}
}

func (hb *HailstoneBig) AtSec(t int) *CoordBig {
	tb := big.NewFloat(float64(t))

	// x := h.pos.x + h.speed.x*t
	x := new(big.Float).Add(hb.pos.x, new(big.Float).Mul(hb.speed.x, tb))
	y := new(big.Float).Add(hb.pos.y, new(big.Float).Mul(hb.speed.y, tb))
	z := new(big.Float).Add(hb.pos.z, new(big.Float).Mul(hb.speed.z, tb))

	return &CoordBig{x, y, z}
}

func findIntersection(p1, p2, p3, p4 *CoordBig) (*CoordBig, error) {
	// Check if the lines are parallel
	// denominator := (p4.y-p3.y)*(p2.x-p1.x) - (p4.x-p3.x)*(p2.y-p1.y)
	denominator := new(big.Float).Sub(
		new(big.Float).Mul(new(big.Float).Sub(p4.y, p3.y), new(big.Float).Sub(p2.x, p1.x)),
		new(big.Float).Mul(new(big.Float).Sub(p4.x, p3.x), new(big.Float).Sub(p2.y, p1.y)),
	)

	// if denominator == 0 {
	if denominator.Cmp(new(big.Float).SetFloat64(0)) == 0 {
		return &CoordBig{}, errors.New("lines are parallel")
	}

	// Calculate the intersection point
	// ua := ((p4.x-p3.x)*(p1.y-p3.y) - (p4.y-p3.y)*(p1.x-p3.x)) / denominator
	// ub := ((p2.X-p1.X)*(p1.Y-p3.Y) - (p2.Y-p1.Y)*(p1.X-p3.X)) / denominator
	ua := new(big.Float).Quo(
		new(big.Float).Sub(
			new(big.Float).Mul(new(big.Float).Sub(p4.x, p3.x), new(big.Float).Sub(p1.y, p3.y)),
			new(big.Float).Mul(new(big.Float).Sub(p4.y, p3.y), new(big.Float).Sub(p1.x, p3.x)),
		),
		denominator,
	)

	//	intersectionX := p1.x + ua*(p2.x-p1.x)
	//	intersectionY := p1.y + ua*(p2.y-p1.y)
	intersectionX := new(big.Float).Add(p1.x, new(big.Float).Mul(ua, new(big.Float).Sub(p2.x, p1.x)))
	intersectionY := new(big.Float).Add(p1.y, new(big.Float).Mul(ua, new(big.Float).Sub(p2.y, p1.y)))

	return &CoordBig{intersectionX, intersectionY, &big.Float{}}, nil
}
