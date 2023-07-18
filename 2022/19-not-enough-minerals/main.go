package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"github.com/alexchao26/advent-of-code-go/util"
	"github.com/emirpasic/gods/maps/hashmap"
	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	"log"
	"regexp"
	"strings"
	"time"
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
	blueprints := parseInput(input)

	//var wg sync.WaitGroup
	//for index, blueprint := range blueprints {
	//	wg.Add(1)
	//
	//	go func(number int, blueprint Blueprint) {
	//		defer wg.Done()
	//
	//		maxGeodes := findBestAlgorythm(blueprint, 24)
	//		fmt.Println(number, maxGeodes)
	//		ans += number * maxGeodes
	//	}(index+1, blueprint)
	//}
	//wg.Wait()

	for _, blueprint := range blueprints {
		maxGeodes := findBestAlgorythm(blueprint, 24)
		fmt.Printf("Blueprint %v. Max geodes: %v.\n", blueprint.number, maxGeodes)
		ans += blueprint.number * maxGeodes
	}

	//{
	//	number := 2
	//	blueprint := blueprints[1]
	//	maxGeodes := findBestAlgorythm(blueprint, 24)
	//	fmt.Printf("Blueprint %v. Max geodes: %v.\n", number, maxGeodes)
	//	ans += number * maxGeodes
	//}

	return ans
}

func part2(input string) (ans int) {
	blueprints := parseInput(input)

	ans = 1
	for _, blueprint := range blueprints {
		if blueprint.number > 3 {
			break
		}

		maxGeodes := findBestAlgorythm(blueprint, 32)
		fmt.Printf("Blueprint %v. Max geodes: %v.\n", blueprint.number, maxGeodes)
		ans *= maxGeodes
	}

	return ans
}

type ResourceType int
type RobotType ResourceType

const (
	RESOURCE_ORE ResourceType = iota
	RESOURCE_CLAY
	RESOURCE_OBSIDIAN
	RESOURCE_GEODE
)

type ProductionCost = [4]int

type Blueprint struct {
	number int
	robots map[RobotType]ProductionCost
}

type Context struct {
	minute     int
	blueprint  Blueprint
	robots     [4]int
	storage    [4]int
	production RobotType // or -1 if nothing is building right now
}

func NewContext(blueprint Blueprint) *Context {
	return &Context{
		blueprint:  blueprint,
		robots:     [4]int{1, 0, 0, 0},
		storage:    [4]int{0, 0, 0, 0},
		production: -1,
	}
}

func (ctx *Context) Copy() *Context {
	newCtx := NewContext(ctx.blueprint)
	newCtx.robots = ctx.robots
	newCtx.storage = ctx.storage
	return newCtx
}

// CanBuild checks if there are enough resources to StartProduction a certain robot.
func (ctx *Context) CanBuild(robot RobotType) bool {
	for resource, cnt := range ctx.blueprint.robots[robot] {
		if ctx.storage[resource] < cnt {
			return false
		}
	}
	return true
}

// StartProduction collects resources from storage and saves info that robot production started.
func (ctx *Context) StartProduction(robot RobotType) {
	for resource, cnt := range ctx.blueprint.robots[robot] {
		if ctx.storage[resource] < cnt {
			log.Panicf("Cannot StartProduction robot %d. Resource %d is not enough. Need %d got %d.", robot, resource, cnt, ctx.storage[resource])
		}
		ctx.storage[resource] -= cnt
	}

	ctx.production = robot
}

// FinishProduction mark robots in production as completed
func (ctx *Context) FinishProduction() {
	if ctx.production != -1 {
		robot := ctx.production
		ctx.robots[robot] += 1
		ctx.production = -1
	}
}

// Harvest resources with existing robots and collect them to storage.
func (ctx *Context) Harvest() {
	for robot, count := range ctx.robots {
		ctx.storage[robot] += count
	}
}

func (ctx *Context) Stringify() string {
	return fmt.Sprintf(
		"r [%d %d %d %d] s [%d %d %d %d]",
		ctx.robots[RESOURCE_ORE],
		ctx.robots[RESOURCE_CLAY],
		ctx.robots[RESOURCE_OBSIDIAN],
		ctx.robots[RESOURCE_GEODE],
		ctx.storage[RESOURCE_ORE],
		ctx.storage[RESOURCE_CLAY],
		ctx.storage[RESOURCE_OBSIDIAN],
		ctx.storage[RESOURCE_GEODE],
	)
}

type ContextQueue struct {
	queue *llq.Queue
}

func NewContextQueue() *ContextQueue {
	return &ContextQueue{
		queue: llq.New(),
	}
}

func (cq *ContextQueue) Size() int {
	return cq.queue.Size()
}

func (cq *ContextQueue) Enqueue(ctx *Context) {
	cq.queue.Enqueue(ctx)
}

func (cq *ContextQueue) Dequeue() (*Context, bool) {
	val, ok := cq.queue.Dequeue()
	return val.(*Context), ok
}

type Optimiser struct {
	maxMinutes int
	uniqueMap  *hashmap.Map
}

func NewOptimiser(maxMinutes int) *Optimiser {
	return &Optimiser{
		maxMinutes: maxMinutes,
		uniqueMap:  hashmap.New(),
	}
}

func (o *Optimiser) Enqueue(ctx *Context, cq *ContextQueue, stats *Stats) {
	// no reason to queue anything at last minute
	if stats.minute == o.maxMinutes {
		return
	}

	// what if we will create geode robot every next minute
	// will we be able to produce more geodes than already produced by better strategy
	if stats.maxGeodes > ctx.storage[RESOURCE_GEODE] {
		minutesLeft := o.maxMinutes - stats.minute
		maxPossibleGeodes := ctx.storage[RESOURCE_GEODE] + MaxPossibleGeodesAddition(minutesLeft, ctx.robots[RESOURCE_GEODE])
		if maxPossibleGeodes < stats.maxGeodes {
			return
		}
	}

	robotsKey := fmt.Sprint(ctx.robots)

	storageList, storageExists := o.uniqueMap.Get(robotsKey)
	if !storageExists {
		storageList = [][4]int{}
	}

	shouldAddCtx := true
	skipIndexes := map[int]bool{}
	for index, storageVariant := range storageList.([][4]int) {
		storageVariantBetter := true
		storageVariantWorse := true
		for resource := 0; resource < 4; resource++ {
			if storageVariant[resource] < ctx.storage[resource] {
				storageVariantBetter = false
			}
			if storageVariant[resource] > ctx.storage[resource] {
				storageVariantWorse = false
			}
		}

		if storageVariantBetter {
			shouldAddCtx = false
		}
		if storageVariantWorse {
			skipIndexes[index] = true
		}
	}

	newStorageList := [][4]int{}
	if shouldAddCtx {
		cq.Enqueue(ctx)
		newStorageList = append(newStorageList, ctx.storage)
	}
	if len(skipIndexes) == 0 {
		newStorageList = append(newStorageList, storageList.([][4]int)...)
	} else {
		for index, storageVariant := range storageList.([][4]int) {
			if _, exists := skipIndexes[index]; !exists {
				newStorageList = append(newStorageList, storageVariant)
			}
		}
	}

	o.uniqueMap.Put(robotsKey, newStorageList)

	// look through all saved storage variants
	// if current context has all 4 resource lower than any saved — no reason so queue ctx
	//for _, storageVariant := range storageList.([][4]int) {
	//	savedVariantBetter := true
	//	for resource := 0; resource < 4; resource++ {
	//		if ctx.storage[resource] > storageVariant[resource] {
	//			savedVariantBetter = false
	//		}
	//	}
	//	if savedVariantBetter {
	//		return
	//	}
	//}
	//
	//cq.Enqueue(ctx)
	//
	//newStorageList := append([][4]int{ctx.storage}, storageList.([][4]int)...)
	//o.uniqueMap.Put(robotsKey, newStorageList)

	//cacheKey := ctx.Stringify()
	//if _, exists := o.uniqueMap.Get(cacheKey); !exists {
	//	cq.Enqueue(ctx)
	//	o.uniqueMap.Put(cacheKey, nil)
	//}
}

type Stats struct {
	minute    int
	maxGeodes int
}

func NewStats() *Stats {
	return &Stats{1, 0}
}

func (s *Stats) SetMinute(minute int) {
	s.minute = minute
}

func (s *Stats) Analyze(ctx *Context) {
	if ctx.storage[RESOURCE_GEODE] > s.maxGeodes {
		s.maxGeodes = ctx.storage[RESOURCE_GEODE]
	}
}

func MaxPossibleGeodesAddition(minutes int, robots int) int {
	// todo use memorisation or dynamic programming
	if minutes == 0 {
		return 0
	} else if minutes == 1 {
		return robots
	} else {
		return robots + MaxPossibleGeodesAddition(minutes-1, robots+1)
	}
}

func findBestAlgorythm(blueprint Blueprint, maxMinutes int) int {
	startContext := NewContext(blueprint)

	stats := NewStats()
	optimiser := NewOptimiser(maxMinutes)

	startQueue := NewContextQueue()
	optimiser.Enqueue(startContext, startQueue, stats)

	thisMinuteQueue := startQueue
	for minute := 1; minute <= maxMinutes; minute++ {
		minuteStarted := time.Now()
		isLastMinute := minute == maxMinutes

		stats.SetMinute(minute)

		prevMinuteQueue := thisMinuteQueue
		thisMinuteQueue = NewContextQueue()
		fmt.Printf("Minute %v started. Queue size: %v.\n", minute, prevMinuteQueue.Size())

		for prevMinuteQueue.Size() > 0 {
			prevCtx, _ := prevMinuteQueue.Dequeue()

			if !isLastMinute {
				// if we can build geode robot — we do it first and do not look into another variants
				{
					robot := RobotType(RESOURCE_GEODE)
					if prevCtx.CanBuild(robot) {
						ctx := prevCtx
						ctx.StartProduction(robot)
						ctx.Harvest()
						ctx.FinishProduction()

						stats.Analyze(ctx)
						optimiser.Enqueue(ctx, thisMinuteQueue, stats)

						continue
					}
				}

				// skip GEODE, because we hardcoded it above
				for resource := 0; resource <= 2; resource++ {
					robot := RobotType(resource)
					if prevCtx.CanBuild(robot) {
						ctx := prevCtx.Copy()
						ctx.StartProduction(robot)
						ctx.Harvest()
						ctx.FinishProduction()

						stats.Analyze(ctx)
						optimiser.Enqueue(ctx, thisMinuteQueue, stats)
					}
				}
			}

			// save CPU for last minute loop
			if isLastMinute && prevCtx.robots[RESOURCE_GEODE] == 0 {
				continue
			}

			ctx := prevCtx
			ctx.Harvest()
			stats.Analyze(ctx)
			optimiser.Enqueue(ctx, thisMinuteQueue, stats)
		}

		minuteProcessTime := time.Since(minuteStarted)
		fmt.Printf("Minute %v processed in %v.\n", minute, minuteProcessTime)
	}

	//for thisMinuteQueue.Size() > 0 {
	//	ctx, _ := thisMinuteQueue.Dequeue()
	//	fmt.Println("r", ctx.robots, "s", ctx.storage)
	//}

	return stats.maxGeodes
}

func parseInput(input string) []Blueprint {
	template := `Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.`
	r := regexp.MustCompile(template)
	results := r.FindAllStringSubmatch(input, -1)

	blueprints := make([]Blueprint, len(results))
	for i, parsed := range results {
		//blueprintIndex := cast.ToInt(parsed[1])
		//if i != blueprintIndex-1 {
		//	log.Panicf("Parsing error. Index should be %d got %d. String: %s", blueprintIndex-1, i, parsed[0])
		//}

		blueprints[i] = Blueprint{
			number: cast.ToInt(parsed[1]),
			robots: map[RobotType]ProductionCost{
				RobotType(RESOURCE_ORE):      {cast.ToInt(parsed[2]), 0, 0, 0},
				RobotType(RESOURCE_CLAY):     {cast.ToInt(parsed[3]), 0, 0, 0},
				RobotType(RESOURCE_OBSIDIAN): {cast.ToInt(parsed[4]), cast.ToInt(parsed[5]), 0, 0},
				RobotType(RESOURCE_GEODE):    {cast.ToInt(parsed[6]), 0, cast.ToInt(parsed[7]), 0},
			},
		}
	}

	return blueprints
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
