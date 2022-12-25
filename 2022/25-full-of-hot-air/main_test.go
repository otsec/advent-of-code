package main

import (
	_ "embed"
	"fmt"
	"testing"
)

//go:embed input_test.txt
var example string

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "example 1",
			input: example,
			want:  "2=-1=0",
		},
		{
			name:  "actual",
			input: input,
			want:  "2=--00--0220-0-21==1",
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
			want:  0,
		},
		//{
		//	name:  "actual",
		//	input: input,
		//	want:  0,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.input); got != tt.want {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_SnafuToHex_HexToSnafu(t *testing.T) {
	tests := []struct {
		snafu string
		hex   int
	}{
		{"1", 1},
		{"2", 2},
		{"1=", 3},
		{"1-", 4},
		{"10", 5},
		{"11", 6},
		{"12", 7},
		{"2=", 8},
		{"2-", 9},
		{"20", 10},
		{"1=0", 15},
		{"1-0", 20},
		{"1=11-2", 2022},
		{"1-0---0", 12345},
		{"1121-1110-1=0", 314159265},

		{"1=-0-2", 1747},
		{"12111", 906},
		{"2=0=", 198},
		{"21", 11},
		{"2=01", 201},
		{"111", 31},
		{"20012", 1257},
		{"112", 32},
		{"1=-1=", 353},
		{"1-12", 107},
		{"12", 7},
		{"1=", 3},
		{"122", 37},

		{"2=-1=0", 4890},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("snafu %v = hex %v", tt.snafu, tt.hex)
		t.Run(name, func(t *testing.T) {
			if got := SnafuToHex(tt.snafu); got != tt.hex {
				t.Errorf("SnafuToHex(%v) = %v, want %v", tt.snafu, got, tt.hex)
			}
		})
	}

	for _, tt := range tests {
		name := fmt.Sprintf("hex %v = snafu %v", tt.hex, tt.snafu)
		t.Run(name, func(t *testing.T) {
			if got := HexToSnafu(tt.hex); got != tt.snafu {
				t.Errorf("HexToSnafu(%v) = %v, want %v", tt.hex, got, tt.snafu)
			}
		})
	}
}
