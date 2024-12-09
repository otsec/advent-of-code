package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os/exec"
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
		_ = CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		_ = CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

func part1(input string) int {
	parsed := parseInput(input)
	disk := NewDisk(parsed)

	compressDiskP1(disk)

	return disk.checksum()
}

func part2(input string) int {
	parsed := parseInput(input)
	disk := NewDisk(parsed)

	compressDiskP2(disk)

	return disk.checksum()
}

type Disk struct {
	values []int
}

func NewDisk(parsed []int) *Disk {
	disk := Disk{}

	isFile := true
	fileId := 0
	for _, size := range parsed {
		if isFile {
			disk.append(fileId, size)
			isFile = false
		} else {
			disk.append(-1, size)
			isFile = true
			fileId += 1
		}
	}

	return &disk
}

func (d *Disk) append(value int, size int) {
	for i := 0; i < size; i++ {
		d.values = append(d.values, value)
	}
}

func (d *Disk) fill(value, start, size int) {
	for i := start; i < start+size; i++ {
		d.values[i] = value
	}
}

func (d *Disk) valSize(index int) int {
	i := index
	for i > 0 && d.values[i-1] == d.values[i] {
		i -= 1
	}

	j := index
	for j < len(d.values)-1 && d.values[j+1] == d.values[j] {
		j += 1
	}

	return j - i + 1
}

func (d *Disk) checksum() int {
	ans := 0
	for i, v := range d.values {
		if v != -1 {
			ans += i * v
		}
	}
	return ans
}

func compressDiskP1(disk *Disk) {
	i := 0
	j := len(disk.values) - 1
	for i < j {
		if disk.values[i] != -1 {
			i++
		} else if disk.values[j] == -1 {
			j--
		} else {
			disk.values[i], disk.values[j] = disk.values[j], disk.values[i]
			i++
			j--
		}
	}
}

func compressDiskP2(disk *Disk) {
	j := len(disk.values) - 1
	for j > 0 {
		rightValue := disk.values[j]
		if rightValue == -1 {
			j--
			continue
		}
		rightSize := disk.valSize(j)

		i := 0
		j = j - rightSize + 1
		for i < j {
			leftValue := disk.values[i]
			leftSize := disk.valSize(i)

			if leftValue != -1 {
				i += leftSize
				continue
			}
			if leftSize < rightSize {
				i += leftSize
				continue
			}

			disk.fill(rightValue, i, rightSize)
			disk.fill(-1, j, rightSize)

			break
		}

		j -= 1
	}
}

func parseInput(input string) []int {
	parsed := make([]int, len(input))
	for i, r := range input {
		parsed[i] = int(r - '0')
	}
	return parsed
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
