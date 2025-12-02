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
		want  int
	}{
		{
			name:  "example 1",
			input: example,
			want:  1227775554,
		},
		{
			name:  "actual",
			input: input,
			want:  44854383294,
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
			want:  4174379265,
		},
		{
			name:  "actual",
			input: input,
			want:  55647141923,
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

func Test_isInvalidPart1(t *testing.T) {
	tests := []struct {
		id   int
		want bool
	}{
		{123, false},
		{1011, false},

		{11, true},
		{22, true},
		{1010, true},
		{1188511885, true},
		{222222, true},
		{446446, true},
		{38593859, true},

		{1111111, false},
		{123123123, false},
		{1212121212, false},

		{111, false},
		{111, false},
		{999, false},
		{565656, false},
		{824824824, false},
		{2121212121, false},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("for %v want %v", tt.id, tt.want)
		t.Run(name, func(t *testing.T) {
			if got := isInvalidPart1(tt.id); got != tt.want {
				t.Errorf("isInvalidPart1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isInvalidPart2(t *testing.T) {
	tests := []struct {
		id   int
		want bool
	}{
		{123, false},
		{1011, false},

		{11, true},
		{22, true},
		{1188511885, true},
		{222222, true},
		{446446, true},
		{38593859, true},

		{1111111, true},
		{123123123, true},
		{1212121212, true},

		{111, true},
		{999, true},
		{565656, true},
		{824824824, true},
		{2121212121, true},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("for %v want %v", tt.id, tt.want)
		t.Run(name, func(t *testing.T) {
			if got := isInvalidPart2(tt.id); got != tt.want {
				t.Errorf("isInvalidPart2() = %v, want %v", got, tt.want)
			}
		})
	}
}
