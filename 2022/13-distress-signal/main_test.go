package main

import (
	"fmt"
	"testing"
)

var example = `[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]`

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example 1",
			input: example,
			want:  13,
		},
		{
			name:  "actual",
			input: input,
			want:  6656,
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
			want:  140,
		},
		{
			name:  "actual",
			input: input,
			want:  19716,
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

func Test_compare(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  int
	}{
		{"", "1", "1", 0},
		{"", "[1,1,1]", "[1,1,1]", 0},
		{"", "[1,1,3,1,1]", "[1,1,5,1,1]", 1},
		{"", "[1,1,6]", "[1,1,5]", -1},
		{"", "[[1],[2,3,4]]", "[[1],4]", 1},
		{"", "[9]", "[[8,7,6]]", -1},
		{"", "[[4,4],4,4]", "[[4,4],4,4,4]", 1},
		{"", "[7,7,7,7]", "[7,7,7]", -1},
		{"", "[]", "[3]", 1},
		{"", "[[[]]]", "[[]]", -1},
		{"", "[1,[2,[3,[4,[5,6,7]]]],8,9]", "[1,[2,[3,[4,[5,6,0]]]],8,9]", -1},
	}
	for i, tt := range tests {
		tt.name = fmt.Sprintf("Dataset #%d", i)
		t.Run(tt.name, func(t *testing.T) {
			if got := compare(parseUnknownJson(tt.left), parseUnknownJson(tt.right)); got != tt.want {
				t.Errorf("compare(%v, %v) = %v, want %v", tt.left, tt.right, got, tt.want)
			}
		})
	}
}
