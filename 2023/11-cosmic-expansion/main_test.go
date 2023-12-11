package main

import (
	_ "embed"
	"testing"
)

//go:embed input_test.txt
var example string

func Test_part1(t *testing.T) {
	tests := []struct {
		name       string
		emptySpace int
		input      string
		want       int
	}{
		{
			name:       "example 1",
			emptySpace: 2,
			input:      example,
			want:       374,
		},
		{
			name:       "actual",
			emptySpace: 2,
			input:      input,
			want:       9536038,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.input, tt.emptySpace); got != tt.want {
				t.Errorf("part1(..., %v) = %v, want %v", tt.emptySpace, got, tt.want)
			}
		})
	}
}

func Test_part2(t *testing.T) {
	tests := []struct {
		name       string
		emptySpace int
		input      string
		want       int
	}{
		{
			name:       "example 1 x10",
			emptySpace: 10,
			input:      example,
			want:       1030,
		},
		{
			name:       "example 1 x100",
			emptySpace: 100,
			input:      example,
			want:       8410,
		},
		{
			name:       "actual",
			emptySpace: 1000000,
			input:      input,
			want:       447744640566,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.input, tt.emptySpace); got != tt.want {
				t.Errorf("part1(..., %v) = %v, want %v", tt.emptySpace, got, tt.want)
			}
		})
	}
}
