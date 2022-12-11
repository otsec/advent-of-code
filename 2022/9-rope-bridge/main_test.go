package main

import (
	"fmt"
	"testing"
)

var example1 = `R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2`

var example2 = `R 5
U 8
L 8
D 3
R 17
D 10
L 25
U 20`

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example 1",
			input: example1,
			want:  13,
		},
		//{
		//	name:  "actual",
		//	input: input,
		//	want:  0,
		//},
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
			name:  "example 2",
			input: example2,
			want:  36,
		},
		{
			name:  "actual",
			input: input,
			want:  2717,
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

func Test_followMove(t *testing.T) {
	tests := []struct {
		inputHead Position
		inputTail Position
		wantTail  Position
	}{
		{
			inputHead: Position{0, 1},
			inputTail: Position{0, 0},
			wantTail:  Position{0, 0},
		},
		{
			inputHead: Position{0, 2},
			inputTail: Position{0, 0},
			wantTail:  Position{0, 1},
		},
		{
			inputHead: Position{1, 2},
			inputTail: Position{0, 0},
			wantTail:  Position{1, 1},
		},
		{
			inputHead: Position{2, 2},
			inputTail: Position{0, 0},
			wantTail:  Position{1, 1},
		},
		{
			inputHead: Position{2, 2},
			inputTail: Position{4, 4},
			wantTail:  Position{3, 3},
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("for head %v tail %v want %v", tt.inputHead, tt.inputTail, tt.wantTail)
		t.Run(name, func(t *testing.T) {
			if followMove(&tt.inputTail, &tt.inputHead); tt.inputTail != tt.wantTail {
				t.Errorf("followMove() = %v, want %v", tt.inputTail, tt.wantTail)
			}
		})
	}
}
