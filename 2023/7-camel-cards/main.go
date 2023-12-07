package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"github.com/alexchao26/advent-of-code-go/util"
	"sort"
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

	const cardPriority = "23456789TJQKA"
	sort.Slice(parsed, func(i, j int) bool {
		hand1 := parsed[i].hand
		hand2 := parsed[j].hand

		res1 := parseGameResult(hand1)
		res2 := parseGameResult(hand2)
		if res1 == res2 {
			return isHandLowerByCards(cardPriority, hand1, hand2)
		} else {
			return res1 < res2
		}
	})

	for i, il := range parsed {
		ans += (i + 1) * il.bid
	}

	return ans
}

func part2(input string) (ans int) {
	parsed := parseInput(input)

	const cardPriority = "J23456789TQKA"
	sort.Slice(parsed, func(i, j int) bool {
		hand1 := parsed[i].hand
		hand2 := parsed[j].hand

		res1 := parseGameResultWithJokers(hand1)
		res2 := parseGameResultWithJokers(hand2)
		if res1 == res2 {
			return isHandLowerByCards(cardPriority, hand1, hand2)
		} else {
			return res1 < res2
		}
	})

	for i, il := range parsed {
		ans += (i + 1) * il.bid
	}

	return ans
}

type InputLine struct {
	hand string
	bid  int
}

const allCards = "23456789TJQKA"

type GameResult int

const (
	HighCard GameResult = iota
	OnePair
	TwoPairs
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

func parseInput(input string) []InputLine {
	lines := strings.Split(input, "\n")

	var parsed []InputLine
	for _, line := range lines {
		segments := strings.Split(line, " ")
		hand := segments[0]
		bid, _ := strconv.Atoi(segments[1])
		parsed = append(parsed, InputLine{hand, bid})
	}

	return parsed
}

func isHandLowerByCards(cardPriority, hand1, hand2 string) bool {
	for i := 0; i < 5; i++ {
		rank1, _ := getCardRank(cardPriority, hand1[i])
		rank2, _ := getCardRank(cardPriority, hand2[i])
		if rank1 != rank2 {
			return rank1 < rank2
		}
	}

	return false
}

func getCardRank(cardPriority string, card byte) (int, error) {
	for i := 0; i < len(cardPriority); i++ {
		if card == cardPriority[i] {
			return i, nil
		}
	}

	return 0, errors.New(fmt.Sprintf("Card %v not found", card))
}

func parseGameResult(hand string) GameResult {
	cardsMap := map[rune]int{}
	for _, card := range allCards {
		cardsMap[card] = 0
	}

	for _, card := range hand {
		cardsMap[card]++
	}

	pairs := 0
	thirds := 0
	for _, amount := range cardsMap {
		if amount == 5 {
			return FiveOfAKind
		}
		if amount == 4 {
			return FourOfAKind
		}
		if amount == 3 {
			thirds++
		}
		if amount == 2 {
			pairs++
		}
	}
	if thirds == 1 && pairs == 1 {
		return FullHouse
	} else if thirds == 1 {
		return ThreeOfAKind
	} else if pairs == 2 {
		return TwoPairs
	} else if pairs == 1 {
		return OnePair
	}

	return HighCard
}

func parseGameResultWithJokers(hand string) GameResult {
	cardsMap := map[rune]int{}
	for _, card := range allCards {
		cardsMap[card] = 0
	}

	for _, card := range hand {
		cardsMap[card]++
	}

	jokers := cardsMap['J']
	cardsMap['J'] = 0
	if jokers >= 4 {
		return FiveOfAKind
	}

	pairs := 0
	thirds := 0
	for _, amount := range cardsMap {
		if amount == 5 || amount+jokers == 5 {
			return FiveOfAKind
		}
		if amount == 4 || amount+jokers == 4 {
			return FourOfAKind
		}
		if amount == 3 {
			thirds++
		}
		if amount == 2 {
			pairs++
		}
	}

	if (thirds == 1 && pairs == 1) || (pairs == 2 && jokers == 1) {
		return FullHouse
	} else if (thirds == 1) || (pairs == 1 && jokers == 1) || (jokers == 2) {
		return ThreeOfAKind
	} else if (pairs == 2) || (pairs == 1 && jokers == 1) {
		return TwoPairs
	} else if (pairs == 1) || (jokers == 1) {
		return OnePair
	}

	return HighCard
}
