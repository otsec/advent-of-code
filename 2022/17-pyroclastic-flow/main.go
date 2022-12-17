package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"github.com/emirpasic/gods/maps/hashmap"
	"log"
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

func part1(input string) int {
	moves := parseInput(input)
	goal := 2022

	sim := NewTetrisSim(7)
	for sim.frozenFigures < goal {
		sim.beforeTick()
		sim.moveFigure(moves[sim.ticks%len(moves)])
		sim.tick()
	}

	return sim.maxY + 1
}

func part2(input string) int {
	moves := parseInput(input)
	goal := 1000000000000

	loopCheckMap := hashmap.New()

	sim := NewTetrisSim(7)
	fakeHeight := 0
	for sim.frozenFigures < goal {
		moveIndex := sim.ticks % len(moves)

		sim.beforeTick()
		sim.moveFigure(moves[moveIndex])
		sim.tick()

		if sim.figIsStuck {
			if fakeHeight != 0 {
				continue
			}

			minY := sim.maxYEveryX[0]
			for _, val := range sim.maxYEveryX {
				if val < minY {
					minY = val
				}
			}

			heightYEveryX := make([]int, sim.width)
			for i, _ := range sim.maxYEveryX {
				heightYEveryX[i] = sim.maxYEveryX[i] - minY
			}

			loopSignature := fmt.Sprintf("%v-%v-%v", moveIndex, sim.figType, heightYEveryX)
			if val, found := loopCheckMap.Get(loopSignature); !found {
				loopCheckMap.Put(loopSignature, []int{sim.maxY, sim.frozenFigures})
			} else {
				maxYPrev := val.([]int)[0]
				frozenFiguresPrev := val.([]int)[1]
				// fmt.Printf("Loop found. Y %v -> %v. Figures: %v -> %v.\n", maxYPrev, sim.maxY, frozenFiguresPrev, sim.frozenFigures)

				loopHeight := sim.maxY - maxYPrev
				loopFigures := sim.frozenFigures - frozenFiguresPrev
				//fmt.Println("Loop Height", loopHeight, "Loop figures", loopFigures)

				loopsRequired := (goal - sim.frozenFigures) / loopFigures
				fakeFigures := loopsRequired * loopFigures
				fakeHeight = loopsRequired * loopHeight
				//fmt.Println("Loops Required", loopsRequired, "Fake Height", fakeHeight, "Fake figures", fakeFigures)

				sim.frozenFigures += fakeFigures
			}
		}
	}

	return fakeHeight + sim.maxY + 1
}

type Move uint8

const (
	MoveDown Move = iota
	MoveLeft
	MoveRight
)

func parseInput(input string) []Move {
	moves := make([]Move, len(input))
	for i, sym := range input {
		if sym == '<' {
			moves[i] = MoveLeft
		} else if sym == '>' {
			moves[i] = MoveRight
		} else {
			log.Panicf("Unknown move: %s.", string(sym))
		}
	}
	return moves
}

type Pos struct {
	x, y int
}

type Figure []Pos

type FigureType uint8

const (
	FigureUnknown FigureType = iota
	FigureHorizontalPlank
	FigurePlus
	FigureReverseL
	FigureVerticalPlank
	FigureCube
)

func chooseNextFigureVariant(currFigure FigureType) FigureType {
	mapping := map[FigureType]FigureType{
		FigureHorizontalPlank: FigurePlus,
		FigurePlus:            FigureReverseL,
		FigureReverseL:        FigureVerticalPlank,
		FigureVerticalPlank:   FigureCube,
		FigureCube:            FigureHorizontalPlank,
	}

	return mapping[currFigure]
}

func createFigure(variant FigureType, atX, atY int) Figure {
	switch variant {
	case FigureHorizontalPlank:
		return Figure{
			Pos{atX + 0, atY},
			Pos{atX + 1, atY},
			Pos{atX + 2, atY},
			Pos{atX + 3, atY},
		}
	case FigurePlus:
		return Figure{
			Pos{atX + 1, atY},
			Pos{atX + 0, atY + 1},
			Pos{atX + 1, atY + 1},
			Pos{atX + 2, atY + 1},
			Pos{atX + 1, atY + 2},
		}
	case FigureReverseL:
		return Figure{
			Pos{atX + 0, atY},
			Pos{atX + 1, atY},
			Pos{atX + 2, atY},
			Pos{atX + 2, atY + 1},
			Pos{atX + 2, atY + 2},
		}
	case FigureVerticalPlank:
		return Figure{
			Pos{atX + 0, atY},
			Pos{atX + 0, atY + 1},
			Pos{atX + 0, atY + 2},
			Pos{atX + 0, atY + 3},
		}
	case FigureCube:
		return Figure{
			Pos{atX + 0, atY},
			Pos{atX + 1, atY},
			Pos{atX + 0, atY + 1},
			Pos{atX + 1, atY + 1},
		}
	default:
		panic(fmt.Sprintf("Unknown figure %d.", variant))
	}

}

type TetrisSimulator struct {
	width         int
	maxY          int
	maxYEveryX    []int
	figType       FigureType
	figPixels     Figure
	figIsStuck    bool
	frozenPixels  *hashmap.Map
	frozenFigures int
	ticks         int
}

func NewTetrisSim(width int) *TetrisSimulator {
	return &TetrisSimulator{
		width:        width,
		maxY:         -1,
		maxYEveryX:   make([]int, width),
		frozenPixels: hashmap.New(),
	}
}

func (s *TetrisSimulator) beforeTick() {
	if s.figType == FigureUnknown {
		s.figType = FigureHorizontalPlank
		s.figPixels = createFigure(s.figType, 2, s.maxY+4)
		//fmt.Println("S", s.figPixels)
	}
	if s.figIsStuck {
		s.figType = chooseNextFigureVariant(s.figType)
		s.figPixels = createFigure(s.figType, 2, s.maxY+4)
		s.figIsStuck = false
		//fmt.Println("S", s.figPixels)
	}
}

func (s *TetrisSimulator) tick() {
	if _, stuck := s.moveFigure(MoveDown); stuck {
		s.figIsStuck = true
		s.freezeFigure()
		s.frozenFigures++
	}

	s.ticks++
}

func (s *TetrisSimulator) moveFigure(move Move) (moved, stuck bool) {
	nextPixels := make(Figure, len(s.figPixels))
	for i, pixel := range s.figPixels {
		if move == MoveDown {
			nextPixels[i] = Pos{pixel.x, pixel.y - 1}
		}
		if move == MoveLeft {
			nextPixels[i] = Pos{pixel.x - 1, pixel.y}
		}
		if move == MoveRight {
			nextPixels[i] = Pos{pixel.x + 1, pixel.y}
		}
	}

	moved = true
	for _, pixel := range nextPixels {
		if pixel.x < 0 || pixel.x > s.width-1 {
			moved = false
		} else if pixel.y < 0 {
			moved = false
		} else if _, found := s.frozenPixels.Get(pixel); found {
			moved = false
		}

		if !moved {
			break
		}
	}

	if moved {
		s.figPixels = nextPixels
	} else {
		stuck = move == MoveDown
	}

	//fmt.Println(move, s.figPixels, stuck)

	return moved, stuck
}

func (s *TetrisSimulator) freezeFigure() {
	for _, pixel := range s.figPixels {
		s.frozenPixels.Put(pixel, true)
		s.maxY = max(pixel.y, s.maxY)
		s.maxYEveryX[pixel.x] = max(pixel.y, s.maxYEveryX[pixel.x])
	}
}

func (s *TetrisSimulator) draw() {
	maxY := s.maxY
	figMap := hashmap.New()
	for _, pixel := range s.figPixels {
		maxY = max(pixel.y, maxY)
		figMap.Put(pixel, true)
	}

	var pos Pos
	for y := maxY; y >= -1; y-- {
		fmt.Printf("%4d ", y)
		for x := -1; x <= s.width; x++ {
			if y == -1 || x == -1 || x == s.width {
				fmt.Print("#")
				continue
			}

			pos = Pos{x, y}
			if _, found := figMap.Get(pos); found {
				fmt.Print("@")
			} else if _, found = s.frozenPixels.Get(pos); found {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
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
