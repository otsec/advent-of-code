package main

import (
	_ "embed"
	"testing"
)

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name: "example 1",
			input: `AAAA
BBCD
BBCC
EEEC`,
			want: 140,
		},
		{
			name: "example 2",
			input: `OOOOO
OXOXO
OOOOO
OXOXO
OOOOO`,
			want: 772,
		},
		{
			name: "example 3",
			input: `RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`,
			want: 1930,
		},
		{
			name:  "actual",
			input: input,
			want:  1477762,
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
			name: "example 1",
			input: `AAAA
BBCD
BBCC
EEEC`,
			want: 80,
		},
		{
			name: "example 2",
			input: `OOOOO
OXOXO
OOOOO
OXOXO
OOOOO`,
			want: 436,
		},
		{
			name: "example 3",
			input: `EEEEE
EXXXX
EEEEE
EXXXX
EEEEE`,
			want: 236,
		},
		{
			name: "example 4",
			input: `AAAAAA
AAABBA
AAABBA
ABBAAA
ABBAAA
AAAAAA`,
			want: 368,
		},
		{
			name: "example 5",
			input: `RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`,
			want: 1206,
		},
		{
			name:  "actual",
			input: input,
			want:  923480,
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
