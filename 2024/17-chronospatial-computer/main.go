package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"math"
	"os/exec"
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
		ans := part1(input)
		CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) string {
	computer := parseInput(input)
	computer.PerformProgram()
	return intsToString(computer.out)
}

func part2(input string) (ans int) {
	computer := parseInput(input)

	// test input works only via bruteforce
	// for actual input it looks for longest match from the end

	stringProgram := intsToString(computer.program)
	longestMatch := 0
	for i := 0; i < 200000; i++ {
		computer.Reset()
		computer.registers.a = i
		computer.PerformProgram()

		if stringProgram == intsToString(computer.out) {
			return i
		}

		if strings.HasSuffix(stringProgram, intsToString(computer.out)) {
			if len(computer.out) > longestMatch {
				longestMatch = len(computer.out)
				ans = i
			}
		}
	}

	// my program: 2,4,1,1,7,5,0,3,1,4,4,5,5,5,3,0
	// tick 1: b = 2 (a % 8)
	// tick 2: b = 3 (b xor 1)
	// tick 3: c = a / 2^b
	// tick 4: a = a / 8
	// tick 5: b = 7 (b xor 4)
	// tick 6: b = b xor c
	// tick 7: out <- b % 8
	// tick 8: if a != 0 goto 0

	// reverse engineering
	// MUST: a % 8 = 2 (first number)
	// written number = (a / 8 % 8 xor 7)
	// every new number makes input number ~x8 (i guess because this processor is 3 bit)

	for longestMatch < len(computer.program) {
		output := intsToString(computer.program[len(computer.program)-longestMatch-1:])
		for i := ans*8 - 10; i < ans*8+10; i++ {
			if i < 0 {
				continue
			}

			computer.Reset()
			computer.registers.a = i
			computer.PerformProgram()

			if intsToString(computer.out) == output {
				longestMatch++
				ans = i
				break
			}
		}
	}

	return ans
}

type Registers struct {
	a, b, c int
}

type Computer struct {
	registers    Registers
	pos          int
	program, out []int
}

func (c *Computer) Reset() {
	c.pos = 0
	c.out = c.out[:0]
}

func (c *Computer) PerformProgram() {
	for c.pos < len(c.program) {
		c.Tick()
	}
}

func (c *Computer) Tick() {
	command := c.program[c.pos]
	operand := c.program[c.pos+1]

	switch command {
	case 0:
		// The adv instruction (opcode 0) performs division. The numerator is the value in the A register.
		// The denominator is found by raising 2 to the power of the instruction's combo operand.
		// (So, an operand of 2 would divide A by 4 (2^2); an operand of 5 would divide A by 2^B.)
		// The result of the division operation is truncated to an integer and then written to the A register.
		numerator := float64(c.registers.a)
		denominator := math.Pow(2.0, float64(c.Combo(operand)))
		c.registers.a = int(numerator / denominator)
	case 1:
		// The bxl instruction (opcode 1) calculates the bitwise XOR of register B
		// and the instruction's literal operand, then stores the result in register B.
		c.registers.b ^= operand
	case 2:
		// The bst instruction (opcode 2) calculates the value of its combo operand modulo 8
		// (thereby keeping only its lowest 3 bits), then writes that value to the B register.
		c.registers.b = c.Combo(operand) % 8
	case 3:
		// The jnz instruction (opcode 3) does nothing if the A register is 0.
		// However, if the A register is not zero, it jumps by setting the instruction
		// pointer to the value of its literal operand; if this instruction jumps,
		// the instruction pointer is not increased by 2 after this instruction.
		if c.registers.a != 0 {
			c.pos = operand
			return
		}
	case 4:
		// The bxc instruction (opcode 4) calculates the bitwise XOR of register B and register C,
		// then stores the result in register B. (For legacy reasons, this instruction reads an operand but ignores it.)
		c.registers.b ^= c.registers.c
	case 5:
		// The out instruction (opcode 5) calculates the value of its combo operand modulo 8, then outputs that value.
		// (If a program outputs multiple values, they are separated by commas.)
		val := c.Combo(operand) % 8
		c.out = append(c.out, val)
	case 6:
		// The bdv instruction (opcode 6) works exactly like the adv instruction
		// except that the result is stored in the B register. (The numerator is still read from the A register.)
		numerator := float64(c.registers.a)
		denominator := math.Pow(2.0, float64(c.Combo(operand)))
		c.registers.b = int(numerator / denominator)
	case 7:
		// The cdv instruction (opcode 7) works exactly like the adv instruction
		// except that the result is stored in the C register. (The numerator is still read from the A register.)
		numerator := float64(c.registers.a)
		denominator := math.Pow(2.0, float64(c.Combo(operand)))
		c.registers.c = int(numerator / denominator)
	}

	c.pos += 2
}

func (c *Computer) Combo(command int) int {
	switch command {
	case 0, 1, 2, 3:
		return command
	case 4:
		return c.registers.a
	case 5:
		return c.registers.b
	case 6:
		return c.registers.c
	default:
		panic("unsupported command")
	}
}

func intsToString(nums []int) string {
	outputStrings := make([]string, len(nums))
	for i, num := range nums {
		outputStrings[i] = strconv.Itoa(num)
	}
	return strings.Join(outputStrings, ",")
}

type Program []int

func parseInput(input string) Computer {
	nums := parseNumbers(input)
	registers := Registers{nums[0], nums[1], nums[2]}
	program := make(Program, len(nums)-3)
	copy(program, nums[3:])
	return Computer{registers, 0, program, []int{}}
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
