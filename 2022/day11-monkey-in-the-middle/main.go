package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"github.com/alexchao26/advent-of-code-go/util"
	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
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
	monkeys := parseInput(input)

	//fmt.Println("Start.")
	//for _, monkey := range monkeys {
	//	fmt.Println("-", monkey)
	//}

	// This slice will count how many items every monkey have inspected.
	inspected := make([]int, len(monkeys))

	for round := 1; round <= 20; round++ {
		for _, monkey := range monkeys {
			for monkey.items.Size() > 0 {
				inspected[monkey.index]++

				val, _ := monkey.items.Dequeue()

				item := ToInt(val)
				item = performOperation(item, monkey.rawOperation)
				item /= 3

				if item%monkey.testValue == 0 {
					monkeys[monkey.ifTruePassTo].items.Enqueue(item)
				} else {
					monkeys[monkey.ifFalsePassTo].items.Enqueue(item)
				}
			}
		}

		//fmt.Printf("Round #%d.\n", round)
		//for _, monkey := range monkeys {
		//	fmt.Println("-", monkey)
		//}
	}

	// Find two max values and calculate answer.
	sort.Slice(inspected, func(i, j int) bool {
		return inspected[i] > inspected[j]
	})
	ans = inspected[0] * inspected[1]

	return ans
}

func part2(input string) (ans int) {
	monkeys := parseInput(input)

	// This slice will count how many items every monkey have inspected.
	inspected := make([]int, len(monkeys))

	for round := 1; round <= 10000; round++ {
		for _, monkey := range monkeys {
			for monkey.items.Size() > 0 {
				inspected[monkey.index]++

				val, _ := monkey.items.Dequeue()

				item := ToInt(val)
				item = performOperation(item, monkey.rawOperation)
				item = item % (2 * 3 * 5 * 7 * 11 * 13 * 17 * 19 * 23)

				if item%monkey.testValue == 0 {
					monkeys[monkey.ifTruePassTo].items.Enqueue(item)
				} else {
					monkeys[monkey.ifFalsePassTo].items.Enqueue(item)
				}
			}
		}

		//if round == 1 || round == 20 || round%1000 == 0 {
		//	fmt.Printf("== After round %d ==\n", round)
		//	for index, value := range inspected {
		//		fmt.Printf("Monkey %d inspected items %d times.\n", index, value)
		//	}
		//}
	}

	// Find two max values and calculate answer.
	sort.Slice(inspected, func(i, j int) bool {
		return inspected[i] > inspected[j]
	})
	ans = inspected[0] * inspected[1]

	return ans
}

type Monkey struct {
	index         int
	items         *llq.Queue
	rawOperation  string
	testValue     int
	ifTruePassTo  int
	ifFalsePassTo int
}

func (m Monkey) String() string {
	return fmt.Sprintf(
		"{Monkey #%d. Items: %v. Operation: %v. Test: /%d ? %d : %d}",
		m.index,
		m.items.Values(),
		m.rawOperation,
		m.testValue,
		m.ifTruePassTo,
		m.ifFalsePassTo,
	)
}

func parseInput(input string) []Monkey {
	template := `Monkey (\d+):\s+` +
		`Starting items: (.+?)\s+` +
		`Operation: new = old (.+?)\s+` +
		`Test: divisible by (\d+)\s+` +
		`If true: throw to monkey (\d+)\s+` +
		`If false: throw to monkey (\d+)`
	r := regexp.MustCompile(template)
	results := r.FindAllStringSubmatch(input, -1)

	monkeys := []Monkey{}
	for _, parsed := range results {
		monkeys = append(monkeys, Monkey{
			index:         cast.ToInt(parsed[1]),
			items:         parseStartingItems(parsed[2]),
			rawOperation:  parsed[3],
			testValue:     cast.ToInt(parsed[4]),
			ifTruePassTo:  cast.ToInt(parsed[5]),
			ifFalsePassTo: cast.ToInt(parsed[6]),
		})
	}

	return monkeys
}

func parseStartingItems(input string) *llq.Queue {
	q := llq.New()
	for _, item := range strings.Split(input, ", ") {
		q.Enqueue(cast.ToInt(item))
	}
	return q
}

func performOperation(item int, rawOperation string) int {
	if rawOperation == "* old" {
		return item * item
	}

	if strings.HasPrefix(rawOperation, "+ ") {
		return item + cast.ToInt(rawOperation[2:])
	}

	if strings.HasPrefix(rawOperation, "* ") {
		return item * cast.ToInt(rawOperation[2:])
	}

	panic("Unknown operation: " + rawOperation)
}

func ToInt(arg interface{}) int {
	var val int
	switch arg.(type) {
	case int:
		val = arg.(int)
	case string:
		val = cast.ToInt(arg.(string))
	default:
		panic(fmt.Sprintf("unhandled type for int casting %T", arg))
	}
	return val
}
