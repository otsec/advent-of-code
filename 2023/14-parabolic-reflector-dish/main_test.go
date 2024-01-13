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
			want:  136,
		},
		{
			name:  "actual",
			input: input,
			want:  109939,
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
			want:  64,
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

func Test_slideBytes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "....OO",
			want:  "OO....",
		},
		{
			input: "..#.OO",
			want:  "..#OO.",
		},
		{
			input: ".O#.OO",
			want:  "O.#OO.",
		},
		{
			input: "OO..O##..O",
			want:  "OOO..##O..",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			copiedInput := make([]byte, len(tt.input))
			copy(copiedInput, tt.input)
			slideBytes(copiedInput)
			got := string(copiedInput)

			if got != tt.want {
				t.Errorf("slideBytes(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
