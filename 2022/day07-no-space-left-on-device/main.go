package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"sort"
	"strings"

	"github.com/alexchao26/advent-of-code-go/util"
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

type SizeMap map[string]int

func part1(input string) (ans int) {
	_, dirs := parseFs(input)

	for _, dirSize := range dirs {
		if dirSize <= 100000 {
			ans += dirSize
		}
	}

	return ans
}

func part2(input string) (ans int) {
	_, dirs := parseFs(input)

	used := dirs["/"]
	available := 70000000 - used
	shouldFree := 30000000 - available

	min := 70000000
	for _, dirSize := range dirs {
		if dirSize >= shouldFree && dirSize < min {
			min = dirSize
			ans = dirSize
		}
	}

	return ans
}

func parseFs(input string) (files SizeMap, dirs SizeMap) {
	files, dirs = make(SizeMap), make(SizeMap)

	lines := strings.Split(input, "\n")
	pathSegments := []string{}
	for _, line := range lines {
		if line == "$ cd /" {
			dirs["/"] = 0
			continue
		}
		if line == "$ cd .." {
			pathSegments = pathSegments[:len(pathSegments)-1]
			continue
		}
		if strings.HasPrefix(line, "$ cd ") {
			dirName := strings.Replace(line, "$ cd ", "", 1)
			pathSegments = append(pathSegments, dirName)

			dirPath := getDirName(pathSegments)
			dirs[dirPath] = 0

			continue
		}

		if strings.HasPrefix(line, "dir") {
			continue
		}
		if line == "$ ls" {
			continue
		}

		// parse file name and size
		{
			fileSegments := strings.Split(line, " ")
			fileSize := cast.ToInt(fileSegments[0])
			fileName := getDirName(pathSegments) + fileSegments[1]
			files[fileName] += fileSize
		}
	}

	for dirName, _ := range dirs {
		dirSize := getDirSize(files, dirName)
		dirs[dirName] = dirSize
	}

	return files, dirs
}

func printFs(fs map[string]int) {
	filenames := []string{}
	for name, _ := range fs {
		filenames = append(filenames, name)
	}

	sort.Strings(filenames)

	for _, name := range filenames {
		fmt.Println(name, fs[name])
	}
}

func getDirName(pathSegments []string) string {
	if len(pathSegments) > 0 {
		return "/" + strings.Join(pathSegments, "/") + "/"
	} else {
		return "/"
	}
}

func getDirSize(fs map[string]int, dir string) (ans int) {
	for name, size := range fs {
		if strings.HasPrefix(name, dir) {
			ans += size
		}
	}
	return ans
}
