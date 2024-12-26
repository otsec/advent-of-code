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
			name:  "actual",
			input: input,
			want:  1524,
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
			name:  "actual",
			input: input,
			want:  1033746,
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

func Test_findFastCheats(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		minSavedSeconds int
		want            int
	}{
		{
			name:            "example 1 (60 sec)",
			input:           example,
			minSavedSeconds: 60,
			want:            1,
		},
		{
			name:            "example 1 (20 sec)",
			input:           example,
			minSavedSeconds: 20,
			want:            5,
		},
		{
			name:            "example 1 (6 sec)",
			input:           example,
			minSavedSeconds: 6,
			want:            16,
		},
		{
			name:            "example 1 (4 sec)",
			input:           example,
			minSavedSeconds: 4,
			want:            30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := NewCharField(tt.input)
			if got := findFastCheats(field, tt.minSavedSeconds); got != tt.want {
				t.Errorf("Test_findFastCheats(%v, %v) = %v, want %v", tt.name, tt.minSavedSeconds, got, tt.want)
			}
		})
	}
}

func Test_findLongCheats(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		minSavedSeconds int
		want            int
	}{

		{
			name:            "example 1 (76 sec)",
			input:           example,
			minSavedSeconds: 76,
			want:            3,
		},
		{
			name:            "example 1 (74 sec)",
			input:           example,
			minSavedSeconds: 74,
			want:            7,
		},
		{
			name:            "example 1 (72 sec)",
			input:           example,
			minSavedSeconds: 72,
			want:            29,
		},
		{
			name:            "example 1 (70 sec)",
			input:           example,
			minSavedSeconds: 70,
			want:            41,
		},
		{
			name:            "example 1 (68 sec)",
			input:           example,
			minSavedSeconds: 68,
			want:            55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := NewCharField(tt.input)
			if got := findLongCheats(field, 20, tt.minSavedSeconds); got != tt.want {
				t.Errorf("Test_findLongCheats(%v, 20, %v) = %v, want %v", tt.name, tt.minSavedSeconds, got, tt.want)
			}
		})
	}
}
