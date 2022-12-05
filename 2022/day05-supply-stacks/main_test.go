package main

import (
	"reflect"
	"testing"
)

var example = `    [D]
[N] [C]
[Z] [M] [P]
 1   2   3

move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2`

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "example",
			input: example,
			want:  "CMZ",
		},
		{
			name:  "actual",
			input: input,
			want:  "VGBBJCRMN",
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
		want  string
	}{
		{
			name:  "example",
			input: example,
			want:  "MCD",
		},
		{
			name:  "actual",
			input: input,
			want:  "LBBVJBRMH",
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

func Test_parseInput(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantCrates map[int]CratesTower
		wantMoves  []Move
	}{
		{
			name:       "example",
			input:      example,
			wantCrates: map[int]CratesTower{1: {"Z", "N"}, 2: {"M", "C", "D"}, 3: {"P"}},
			wantMoves:  []Move{{2, 1, 1}, {1, 3, 3}, {2, 1, 2}, {1, 2, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: fix tests
			// crates, moves := parseInput(tt.input)
			// if !reflect.DeepEqual(crates, fmt.Sprintf("%v", tt.wantCrates)) {
			// 	 t.Errorf("parseInput() crates = %v, want %v", crates, tt.wantCrates)
			// }
			_, moves := parseInput(tt.input)
			if !reflect.DeepEqual(moves, tt.wantMoves) {
				t.Errorf("parseInput() moves = %v, want %v", moves, tt.wantMoves)
			}
		})
	}
}
