package main

import (
	_ "embed"
	"testing"
)

//go:embed input_test.txt
var example string

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example 1",
			input: example,
			want:  6440,
		},
		{
			name:  "actual",
			input: input,
			want:  251058093,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.input); got != tt.want {
				t.Errorf("part1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_part2(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example 1",
			input: example,
			want:  5905,
		},
		{
			name:  "actual",
			input: input,
			want:  249781879,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.input); got != tt.want {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseGameResultWithJokers(t *testing.T) {
	tests := []struct {
		hand string
		want GameResult
	}{
		{"32T3K", OnePair},
		{"T55J5", FourOfAKind},
		{"KK677", TwoPairs},
		{"KTJJT", FourOfAKind},
		{"QQQJA", FourOfAKind},

		{"KKKKJ", FiveOfAKind},
		{"KKKJJ", FiveOfAKind},
		{"KKJJJ", FiveOfAKind},
		{"KJJJJ", FiveOfAKind},
		{"JJJJJ", FiveOfAKind},
		{"KKKJ2", FourOfAKind},
		{"KKJJ2", FourOfAKind},
		{"KJJJ2", FourOfAKind},
		{"KKJ22", FullHouse},
		{"KKJ12", ThreeOfAKind},
		{"K12JJ", ThreeOfAKind},
		{"K127J", OnePair},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := parseGameResultWithJokers(tt.hand); got != tt.want {
				t.Errorf("parseGameResultWithJokers(%v) = %v, want %v", tt.hand, got, tt.want)
			}
		})
	}
}
