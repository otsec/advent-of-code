package main

import (
	_ "embed"
	"testing"
)

//go:embed input_test_1.txt
var example1 string

//go:embed input_test_2.txt
var example2 string

//go:embed input_test_3.txt
var example3 string

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example 1",
			input: example1,
			want:  2,
		},
		{
			name:  "example 2",
			input: example2,
			want:  6,
		},
		{
			name:  "actual",
			input: input,
			want:  13019,
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
			name:  "example 3",
			input: example3,
			want:  6,
		},
		{
			name:  "actual",
			input: input,
			want:  13524038372771,
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
