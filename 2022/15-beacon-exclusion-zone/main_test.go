package main

import (
	"fmt"
	"testing"
)

var example = `Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3`

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		atY   int
		want  int
	}{
		{
			name:  "example 1",
			input: example,
			atY:   10,
			want:  26,
		},
		{
			name:  "actual",
			input: input,
			atY:   2000000,
			want:  4951427,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.input, tt.atY); got != tt.want {
				t.Errorf("part1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_part2(t *testing.T) {
	tests := []struct {
		name  string
		input string
		maxY  int
		want  int
	}{
		{
			name:  "example 1",
			input: example,
			maxY:  20,
			want:  56000011,
		},
		{
			name:  "actual",
			input: input,
			maxY:  4000000,
			want:  13029714573243,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.input, tt.maxY); got != tt.want {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calcManhattanDistance(t *testing.T) {
	tests := []struct {
		name   string
		sensor Coordinate
		beacon Coordinate
		want   int
	}{
		{
			sensor: Coordinate{8, 7},
			beacon: Coordinate{2, 10},
			want:   9,
		},
	}
	for i, tt := range tests {
		tt.name = fmt.Sprintf("Dataset #%d", i)
		t.Run(tt.name, func(t *testing.T) {
			if got := calcManhattanDistance(tt.sensor, tt.beacon); got != tt.want {
				t.Errorf("calcManhattanDistance(%v, %v) = %v, want %v", tt.sensor, tt.beacon, got, tt.want)
			}
		})
	}
}

func Test_mergeLines(t *testing.T) {
	tests := []struct {
		name string
		in   []Line
		new  Line
		want []Line
	}{
		{
			in:   []Line{{10, 19}},
			new:  Line{20, 30},
			want: []Line{{10, 30}},
		},
		{
			in:   []Line{{10, 20}},
			new:  Line{1, 9},
			want: []Line{{1, 20}},
		},
		{
			in:   []Line{{10, 20}, {25, 30}},
			new:  Line{40, 50},
			want: []Line{{10, 20}, {25, 30}, {40, 50}},
		},
		{
			in:   []Line{{25, 30}, {40, 50}},
			new:  Line{10, 20},
			want: []Line{{10, 20}, {25, 30}, {40, 50}},
		},
		{
			in:   []Line{{10, 20}, {40, 50}},
			new:  Line{25, 30},
			want: []Line{{10, 20}, {25, 30}, {40, 50}},
		},
		{
			in:   []Line{{10, 20}, {40, 50}},
			new:  Line{15, 30},
			want: []Line{{10, 30}, {40, 50}},
		},
		{
			in:   []Line{{10, 20}, {40, 50}},
			new:  Line{35, 55},
			want: []Line{{10, 20}, {35, 55}},
		},
		{
			in:   []Line{{10, 20}, {40, 50}, {60, 70}},
			new:  Line{15, 65},
			want: []Line{{10, 70}},
		},
		{
			in:   []Line{{10, 20}, {40, 50}, {60, 70}},
			new:  Line{5, 85},
			want: []Line{{5, 85}},
		},
	}
	for i, tt := range tests {
		tt.name = fmt.Sprintf("Dataset #%d", i)
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeLines(tt.in, tt.new); fmt.Sprint(got) != fmt.Sprint(tt.want) {
				t.Errorf("mergeLines() = %v, want %v", got, tt.want)
			}
		})
	}
}
