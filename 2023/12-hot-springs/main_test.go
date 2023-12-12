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
			want:  21,
		},
		{
			name:  "actual",
			input: input,
			want:  6935,
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
			want:  525152,
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

func Test_isResultMatched(t *testing.T) {
	tests := []struct {
		input   string
		numbers []int
		want    bool
	}{
		{
			input:   "###.###",
			numbers: []int{1, 1, 3},
			want:    false,
		},
		{
			input:   "#.#.###",
			numbers: []int{1, 1, 3},
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := isResultMatched(tt.input, tt.numbers); got != tt.want {
				t.Errorf("isResultMatched(%v, %v) = %v, want %v", tt.input, tt.numbers, got, tt.want)
			}
		})
	}
}

func Test_testPossibleVariants(t *testing.T) {
	tests := []struct {
		input   string
		numbers []int
		want    int
	}{
		{
			input:   "???.###",
			numbers: []int{1, 1, 3},
			want:    1,
		},
		{
			input:   ".??..??...?##.",
			numbers: []int{1, 1, 3},
			want:    4,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := testPossibleVariants(tt.input, tt.numbers); got != tt.want {
				t.Errorf("testPossibleVariants(%v, %v) = %v, want %v", tt.input, tt.numbers, got, tt.want)
			}
		})
	}
}
