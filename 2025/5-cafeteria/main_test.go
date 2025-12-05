package main

import (
	_ "embed"
	"reflect"
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
			want:  3,
		},
		{
			name:  "actual",
			input: input,
			want:  607,
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
			want:  14,
		},
		{
			name:  "actual",
			input: input,
			want:  342433357244012,
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

func Test_mergeRanges(t *testing.T) {
	tests := []struct {
		name   string
		r1     FreshRange
		r2     FreshRange
		merged bool
		want   FreshRange
	}{
		{
			name:   "{1, 4} {5, 8} not merged",
			r1:     FreshRange{1, 4},
			r2:     FreshRange{5, 8},
			merged: false,
			want:   FreshRange{0, 0},
		},
		{
			name:   "{1, 5} {5, 8} merged to {1, 8}",
			r1:     FreshRange{1, 5},
			r2:     FreshRange{5, 8},
			merged: true,
			want:   FreshRange{1, 8},
		},
		{
			name:   "{1, 8} {4, 5} merged to {1, 8}",
			r1:     FreshRange{1, 8},
			r2:     FreshRange{4, 5},
			merged: true,
			want:   FreshRange{1, 8},
		},
		{
			name:   "{5, 8} {10, 15} not merged",
			r1:     FreshRange{5, 8},
			r2:     FreshRange{10, 15},
			merged: false,
			want:   FreshRange{0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, merged := mergeRanges(tt.r1, tt.r2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeRanges() got = %v, want %v", got, tt.want)
			}
			if merged != tt.merged {
				t.Errorf("mergeRanges() merged = %v, want %v", merged, tt.merged)
			}
		})
	}
}
