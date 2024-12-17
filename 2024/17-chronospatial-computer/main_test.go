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
		want  string
	}{
		{
			name: "example 1",
			input: `Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0`,
			want: "4,6,3,5,6,3,5,2,1,0",
		},
		{
			name:  "actual",
			input: input,
			want:  "4,0,4,7,1,2,7,1,6",
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
			input: `Register A: 2024
Register B: 0
Register C: 0

Program: 0,3,5,4,3,0`,
			want: 117440,
		},
		{
			name:  "actual",
			input: input,
			want:  202322348616234,
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
