package main

import (
	"testing"
)

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example 1",
			input: "mjqjpqmgbljsphdztnvjfqwrcgsmlb",
			want:  7,
		},
		{
			name:  "example 2",
			input: "bvwbjplbgvbhsrlpgdmjqwftvncz",
			want:  5,
		},
		{
			name:  "example 3",
			input: "nppdvjthqldpwncqszvftbrmjlhg",
			want:  6,
		},
		{
			name:  "example 4",
			input: "nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg",
			want:  10,
		},
		{
			name:  "example 5",
			input: "zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw",
			want:  11,
		},
		{
			name:  "actual",
			input: input,
			want:  1909,
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
			input: "mjqjpqmgbljsphdztnvjfqwrcgsmlb",
			want:  19,
		},
		{
			name:  "example 2",
			input: "bvwbjplbgvbhsrlpgdmjqwftvncz",
			want:  23,
		},
		{
			name:  "example 3",
			input: "nppdvjthqldpwncqszvftbrmjlhg",
			want:  23,
		},
		{
			name:  "example 4",
			input: "nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg",
			want:  29,
		},
		{
			name:  "example 5",
			input: "zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw",
			want:  26,
		},
		{
			name:  "actual",
			input: input,
			want:  3380,
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
