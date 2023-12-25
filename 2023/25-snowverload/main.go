package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"reflect"
	"regexp"
	"slices"
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
	parsed := parseInput(input)

	// dot -T svg graph_test.viz > graph_test.svg
	// dot -Tsvg -Kneato graph_part1.viz > graph_part1.svg

	slices.DeleteFunc(parsed, func(c Connection) bool {
		// test input
		if reflect.DeepEqual(c, Connection{"cmg", "bvb"}) {
			return true
		}
		if reflect.DeepEqual(c, Connection{"pzl", "hfx"}) {
			return true
		}
		if reflect.DeepEqual(c, Connection{"jqt", "nvd"}) {
			return true
		}

		// part1
		if reflect.DeepEqual(c, Connection{"znk", "mmr"}) {
			return true
		}
		if reflect.DeepEqual(c, Connection{"rnx", "ddj"}) {
			return true
		}
		if reflect.DeepEqual(c, Connection{"vcq", "lxb"}) {
			return true
		}

		return false
	})

	cmap := createConnectionMap(parsed)
	plugs := cmap.AllPlugs()

	var circuit Circuit
	circuit.FindAllPlugs(&cmap, plugs[0])
	if len(circuit) != len(plugs) {
		return len(circuit) * (len(plugs) - len(circuit))
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)
	_ = parsed

	return ans
}

type Connection [2]string

func parseInput(input string) []Connection {
	var connections []Connection

	lines := strings.Split(input, "\n")
	re := regexp.MustCompile(`\w+`)
	for _, line := range lines {
		items := re.FindAllString(line, -1)
		for i := 1; i < len(items); i++ {
			conn := Connection{items[0], items[i]}
			connections = append(connections, conn)
		}
	}

	return connections
}

type ConnectionMap map[string][]string

func (cmap *ConnectionMap) Connect(a, b string) {
	values, existsKey := (*cmap)[a]
	if !existsKey {
		(*cmap)[a] = []string{b}
		return
	}
	if !slices.Contains(values, b) {
		(*cmap)[a] = append(values, b)
	}
}

func (cmap *ConnectionMap) AllPlugs() []string {
	plugs := make([]string, 0, len(*cmap))
	for k := range *cmap {
		plugs = append(plugs, k)
	}
	return plugs
}

func createConnectionMap(conns []Connection) ConnectionMap {
	cmap := make(ConnectionMap)
	for _, conn := range conns {
		cmap.Connect(conn[0], conn[1])
		cmap.Connect(conn[1], conn[0])
	}
	return cmap
}

type Circuit []string

func (circuit *Circuit) FindAllPlugs(cmap *ConnectionMap, plug string) {
	if slices.Contains(*circuit, plug) {
		return
	}

	*circuit = append(*circuit, plug)

	if plugs, exists := (*cmap)[plug]; exists {
		for _, p := range plugs {
			circuit.FindAllPlugs(cmap, p)
		}
	}
}
