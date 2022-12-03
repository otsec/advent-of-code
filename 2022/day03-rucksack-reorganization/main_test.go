package main

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "Part 1",
			input: input,
			want:  8298,
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
			name:  "Part 2",
			input: input,
			want:  2708,
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

func Test_splitRucksack(t *testing.T) {
	tests := []struct {
		input string
		want1 string
		want2 string
	}{
		{
			input: "GwrhJPDJCZFRcwfZWV",
			want1: "GwrhJPDJC",
			want2: "ZFRcwfZWV",
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("for %v want %v, %v", tt.input, tt.want1, tt.want2)
		t.Run(name, func(t *testing.T) {
			got, got1 := splitRucksack(tt.input)
			if got != tt.want1 {
				t.Errorf("splitRucksack() got1 = %v, want %v", got, tt.want1)
			}
			if got1 != tt.want2 {
				t.Errorf("splitRucksack() got2 = %v, want %v", got1, tt.want2)
			}
		})
	}
}

func Test_priorityOfItem(t *testing.T) {
	tests := []struct {
		input rune
		want  int
	}{
		{
			input: 'a',
			want:  1,
		},
		{
			input: 'b',
			want:  2,
		},
		{
			input: 'z',
			want:  26,
		},
		{
			input: 'A',
			want:  27,
		},
		{
			input: 'Z',
			want:  52,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("for %v want %v", tt.input, tt.want)
		t.Run(name, func(t *testing.T) {
			if got := priorityOfItem(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("priorityOfItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
