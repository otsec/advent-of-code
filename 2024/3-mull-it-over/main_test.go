package main

import (
	_ "embed"
	"fmt"
	"testing"
)

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example 1",
			input: "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))",
			want:  161,
		},
		{
			name:  "actual",
			input: input,
			want:  164730528,
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
			input: "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))",
			want:  48,
		},
		{
			name:  "actual",
			input: input,
			want:  70478672,
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

func Test_removeDonts(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))",
			want:  "xmul(2,4)&mul[3,7]!^?mul(8,5))",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			if got := removeDonts(tt.input); got != tt.want {
				t.Errorf("removeDonts() = %v", got)
			}
		})
	}
}
