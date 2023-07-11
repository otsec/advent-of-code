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
			want:  33,
		},
		{
			name:  "actual",
			input: input,
			want:  1834,
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

func Test_findBestAlgorythm(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		blueprint int
		minutes   int
		want      int
	}{
		// test data for part 1
		{
			input:     example,
			blueprint: 1,
			minutes:   24,
			want:      9,
		},
		{
			input:     example,
			blueprint: 2,
			minutes:   24,
			want:      12,
		},

		// real values for part 1
		// 1, 2, 0, 0, 5, 30, 56, 16, 9, 40, 0, 48, 13, 168, 120, 0, 0, 90, 95, 0, 21, 110, 115, 0, 325, 0, 189, 84, 0, 210
		{
			input:     input,
			blueprint: 1,
			minutes:   24,
			want:      1,
		},
		{
			input:     input,
			blueprint: 2,
			minutes:   24,
			want:      1,
		},
		{
			input:     input,
			blueprint: 3,
			minutes:   24,
			want:      0,
		},
		{
			input:     input,
			blueprint: 4,
			minutes:   24,
			want:      0,
		},
		{
			input:     input,
			blueprint: 5,
			minutes:   24,
			want:      1,
		},
		{
			input:     input,
			blueprint: 6,
			minutes:   24,
			want:      5,
		},
		{
			input:     input,
			blueprint: 7,
			minutes:   24,
			want:      9,
		},
		{
			input:     input,
			blueprint: 13,
			minutes:   24,
			want:      1,
		},
		{
			input:     input,
			blueprint: 15,
			minutes:   24,
			want:      8,
		},
		{
			input:     input,
			blueprint: 16,
			minutes:   24,
			want:      0,
		},
		{
			input:     input,
			blueprint: 17,
			minutes:   24,
			want:      0,
		},
		{
			input:     input,
			blueprint: 21,
			minutes:   24,
			want:      1,
		},
		{
			input:     input,
			blueprint: 23,
			minutes:   24,
			want:      5,
		},
		{
			input:     input,
			blueprint: 24,
			minutes:   24,
			want:      0,
		},
		{
			input:     input,
			blueprint: 25,
			minutes:   24,
			want:      15,
		},
		{
			input:     input,
			blueprint: 26,
			minutes:   24,
			want:      0,
		},
		{
			input:     input,
			blueprint: 28,
			minutes:   24,
			want:      3,
		},
		{
			input:     input,
			blueprint: 29,
			minutes:   24,
			want:      0,
		},
		{
			input:     input,
			blueprint: 30,
			minutes:   24,
			want:      8,
		},

		// test data for part 2
		{
			input:     example,
			blueprint: 1,
			minutes:   32,
			want:      56,
		},
		{
			input:     example,
			blueprint: 2,
			minutes:   32,
			want:      62,
		},
	}
	for _, tt := range tests {
		tt.name = fmt.Sprintf("Bluepirnt %d. Max geodes is %d.", tt.blueprint, tt.want)
		t.Run(tt.name, func(t *testing.T) {
			blueprints := parseInput(tt.input)
			if got := findBestAlgorythm(blueprints[tt.blueprint-1], tt.minutes); got != tt.want {
				t.Errorf("findBestAlgorythm() = %v, want %v", got, tt.want)
			}
		})
	}
}
