package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/cast"
	"github.com/alexchao26/advent-of-code-go/util"
	"log"
	"regexp"
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
	blueprints := parseInput(input)

	resources := []RESOURCE_TYPE{
		RESOURCE_ORE,
		RESOURCE_CLAY,
		RESOURCE_OBSIDIAN,
		RESOURCE_GEODE,
	}

	// 1, 2, 0, 0, 5, 30, 56, 16, 9, 40, 0, 48, 13, 168, 120, 0, 0, 90, 95, 0, 21, 110, 115, 0, 325, 0, 189, 84, 0, 210

	for bi, blueprint := range blueprints {
		variants := []Context{
			createContext(blueprint),
		}

		ansBlu := 0
		cache := make(map[string]bool)

		for m := 1; m <= 24; m++ {
			stats := struct {
				obsidianRobotsMin, obsidianRobotsMax int
				geodesMin, geodesMax                 int
			}{
				obsidianRobotsMin: variants[0].robots[RESOURCE_OBSIDIAN],
				obsidianRobotsMax: variants[0].robots[RESOURCE_OBSIDIAN],
				geodesMin:         variants[0].storage[RESOURCE_GEODE],
				geodesMax:         variants[0].storage[RESOURCE_GEODE],
			}

			// maybe build a robot
			for _, ctx := range variants {
				cacheKey := fmt.Sprint(m, ctx.robots, ctx.production, ctx.storage)
				cache[cacheKey] = true

				for _, resource := range resources {
					if canBuild(&ctx, resource) {
						next := copyContext(&ctx)
						startProduction(&next, resource)

						cacheKey = fmt.Sprint(m, next.robots, next.production, next.storage)
						if _, found := cache[cacheKey]; !found {
							variants = append(variants, next)
							cache[cacheKey] = true
						}
					}
				}
			}

			// exec all variants
			for _, ctx := range variants {
				harvest(&ctx)
				finishProduction(&ctx)

				// stats
				stats.obsidianRobotsMin = min(stats.obsidianRobotsMin, ctx.robots[RESOURCE_OBSIDIAN])
				stats.obsidianRobotsMax = max(stats.obsidianRobotsMax, ctx.robots[RESOURCE_OBSIDIAN])
				stats.geodesMin = min(stats.geodesMin, ctx.storage[RESOURCE_GEODE])
				stats.geodesMax = max(stats.geodesMax, ctx.storage[RESOURCE_GEODE])
			}

			// filter by obsidian robots
			if stats.obsidianRobotsMax > 3 && stats.obsidianRobotsMin == 0 {
				filtered := []Context{}
				for _, ctx := range variants {
					if ctx.robots[RESOURCE_OBSIDIAN] == 0 {
						continue
					}
					filtered = append(filtered, ctx)
				}
				variants = filtered
			}
			if stats.obsidianRobotsMax-stats.obsidianRobotsMin > 3 {
				filtered := []Context{}
				for _, ctx := range variants {
					if stats.obsidianRobotsMax-ctx.robots[RESOURCE_OBSIDIAN] > 3 {
						continue
					}
					filtered = append(filtered, ctx)
				}
				variants = filtered
			}

			// filter by geodes
			if stats.geodesMax-stats.geodesMin > 3 {
				filtered := []Context{}
				for _, ctx := range variants {
					if stats.geodesMax-ctx.storage[RESOURCE_GEODE] > 3 {
						continue
					}
					filtered = append(filtered, ctx)
				}
				variants = filtered
			}

			//if m == 20 && stats.geodesMax == 0 {
			//	break
			//}

			ansVar := stats.geodesMax * (bi + 1) // + len(blueprints)
			if ansVar > ansBlu {
				ansBlu = ansVar
			}

			fmt.Println("Blueprint", bi+1, "min", m, "cache", len(cache), "stats", stats)
		}

		fmt.Println("Blueprint", bi+1, "geodes", ansBlu)

		ans += ansBlu
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)
	_ = parsed

	return ans
}

type RESOURCE_TYPE int

const (
	RESOURCE_ORE RESOURCE_TYPE = iota
	RESOURCE_CLAY
	RESOURCE_OBSIDIAN
	RESOURCE_GEODE
)

type Blueprint map[RESOURCE_TYPE]map[RESOURCE_TYPE]int

type Context struct {
	minute     int
	blueprint  Blueprint
	production map[RESOURCE_TYPE]int
	robots     map[RESOURCE_TYPE]int
	storage    map[RESOURCE_TYPE]int
}

func createContext(blueprint Blueprint) Context {
	return Context{
		blueprint: blueprint,
		production: map[RESOURCE_TYPE]int{
			RESOURCE_ORE:      0,
			RESOURCE_CLAY:     0,
			RESOURCE_OBSIDIAN: 0,
			RESOURCE_GEODE:    0,
		},
		robots: map[RESOURCE_TYPE]int{
			RESOURCE_ORE:      1,
			RESOURCE_CLAY:     0,
			RESOURCE_OBSIDIAN: 0,
			RESOURCE_GEODE:    0,
		},
		storage: map[RESOURCE_TYPE]int{
			RESOURCE_ORE:      0,
			RESOURCE_CLAY:     0,
			RESOURCE_OBSIDIAN: 0,
			RESOURCE_GEODE:    0,
		},
	}
}

func copyContext(from *Context) Context {
	ctx := createContext(from.blueprint)
	for resource, cnt := range from.production {
		ctx.production[resource] = cnt
	}
	for resource, cnt := range from.robots {
		ctx.robots[resource] = cnt
	}
	for resource, cnt := range from.storage {
		ctx.storage[resource] = cnt
	}
	return ctx
}

func isBadVariant(minute int, ctx *Context) bool {
	if minute >= 7 && ctx.robots[RESOURCE_CLAY] == 0 {
		return true
	}
	if minute >= 14 && ctx.robots[RESOURCE_OBSIDIAN] == 0 {
		return true
	}
	if minute >= 20 && ctx.robots[RESOURCE_GEODE] == 0 {
		return true
	}
	if minute >= 22 && ctx.robots[RESOURCE_GEODE] == 1 {
		return true
	}
	return false
}

func reduceVariants(minute int, variants []Context) []Context {
	if minute != 7 && minute != 14 && minute != 20 && minute != 22 {
		return variants
	}

	ans := []Context{}
	for _, ctx := range variants {
		if !isBadVariant(minute, &ctx) {
			ans = append(ans, ctx)
		}
	}
	return ans
}

// canBuild checks if resource are enough to startProduction a certain robot.
func canBuild(ctx *Context, robot RESOURCE_TYPE) bool {
	for resource, cnt := range ctx.blueprint[robot] {
		if ctx.storage[resource] < cnt {
			return false
		}
	}
	return true
}

// startProduction collects resources from storage and saves info that robot production started.
func startProduction(ctx *Context, robot RESOURCE_TYPE) {
	for resource, cnt := range ctx.blueprint[robot] {
		if ctx.storage[resource] < cnt {
			log.Panicf("Cannot startProduction robot %d. Resource %d is not enough. Need %d got %d.", robot, resource, cnt, ctx.storage[resource])
		}
		ctx.storage[resource] -= cnt
	}

	ctx.production[robot] += 1
}

// finishProduction mark robots in production as completed
func finishProduction(ctx *Context) {
	for robot, cnt := range ctx.production {
		ctx.robots[robot] += cnt
		ctx.production[robot] -= cnt
	}
}

// harvest resources with existing robots and collect them to storage.
func harvest(ctx *Context) {
	for robot, count := range ctx.robots {
		ctx.storage[robot] += count
	}
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
			RESOURCE_ORE: {
				RESOURCE_ORE: cast.ToInt(parsed[2]),
			},
			RESOURCE_CLAY: {
				RESOURCE_ORE: cast.ToInt(parsed[3]),
			},
			RESOURCE_OBSIDIAN: {
				RESOURCE_ORE:  cast.ToInt(parsed[4]),
				RESOURCE_CLAY: cast.ToInt(parsed[5]),
			},
			RESOURCE_GEODE: {
				RESOURCE_ORE:      cast.ToInt(parsed[6]),
				RESOURCE_OBSIDIAN: cast.ToInt(parsed[7]),
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
