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
			want:  21,
		},
		{
			name:  "actual",
			input: input,
			want:  6935,
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
		//{
		//	name:  "example 1",
		//	input: example,
		//	want:  525152,
		//},
		//{
		//	name:  "actual",
		//	symbols: symbols,
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

func Test_isResultMatched(t *testing.T) {
	tests := []struct {
		symbols string
		groups  []int
		want    bool
	}{
		{
			symbols: "###.###",
			groups:  []int{1, 1, 3},
			want:    false,
		},
		{
			symbols: "#.#.###",
			groups:  []int{1, 1, 3},
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := isResultMatched(InputRow{tt.symbols, tt.groups}); got != tt.want {
				t.Errorf("isResultMatched(%v, %v) = %v, want %v", tt.symbols, tt.groups, got, tt.want)
			}
		})
	}
}

func Test_makesSense(t *testing.T) {
	tests := []struct {
		symbols string
		groups  []int
		want    bool
	}{
		{symbols: "#.#", groups: []int{1, 1}, want: true},
		{symbols: ".#.#", groups: []int{1, 1}, want: true},
		{symbols: "...#..#...", groups: []int{1, 1}, want: true},
		{symbols: "...##..#...", groups: []int{1, 1}, want: false},
		{symbols: "???.###", groups: []int{1, 1}, want: true},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := makesSense(InputRow{tt.symbols, tt.groups}); got != tt.want {
				t.Errorf("makesSence(%v, %v) = %v, want %v", tt.symbols, tt.groups, got, tt.want)
			}
		})
	}
}

func Test_testPossibleVariants(t *testing.T) {
	tests := []struct {
		symbols string
		groups  []int
		want    int
	}{
		{
			symbols: "???.###",
			groups:  []int{1, 1, 3},
			want:    1,
		},
		{
			symbols: ".??..??...?##.",
			groups:  []int{1, 1, 3},
			want:    4,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := testPossibleVariants(InputRow{tt.symbols, tt.groups}); got != tt.want {
				t.Errorf("testPossibleVariants(%v, %v) = %v, want %v", tt.symbols, tt.groups, got, tt.want)
			}
		})
	}
}

func Test_makeItLarger(t *testing.T) {
	tests := []struct {
		symbols     string
		groups      []int
		wantSymbols string
		wantGroups  []int
	}{
		{
			symbols:     ".#",
			groups:      []int{1},
			wantSymbols: ".#?.#?.#?.#?.#",
			wantGroups:  []int{1, 1, 1, 1, 1},
		},
		{
			symbols:     "???.###",
			groups:      []int{1, 1, 3},
			wantSymbols: "???.###????.###????.###????.###????.###",
			wantGroups:  []int{1, 1, 3, 1, 1, 3, 1, 1, 3, 1, 1, 3, 1, 1, 3},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := makeItLarger(InputRow{tt.symbols, tt.groups})

			if got.symbols != tt.wantSymbols {
				t.Errorf("makeItLarger(%v, %v) got symbols %v, want %v", tt.symbols, tt.groups, got.symbols, tt.wantSymbols)
			}

			if !reflect.DeepEqual(got.groups, tt.wantGroups) {
				t.Errorf("makeItLarger(%v, %v) got groups %v, want %v", tt.symbols, tt.groups, got.groups, tt.wantGroups)
			}
		})
	}
}

func Test_shrink(t *testing.T) {
	tests := []struct {
		symbols     string
		groups      []int
		wantSymbols string
		wantGroups  []int
	}{
		{"?", []int{1}, "?", []int{1}},
		{"?..##", []int{1}, "?..##", []int{1}},
		{".?..##", []int{1}, "?..##", []int{1, 2}},
		{".#.?..##", []int{1}, "?..##", []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := shrink(InputRow{tt.symbols, tt.groups})

			if got.symbols != tt.wantSymbols {
				t.Errorf("shrink(%v, %v) got symbols %v, want %v", tt.symbols, tt.groups, got.symbols, tt.wantSymbols)
			}

			if !reflect.DeepEqual(got.groups, tt.wantGroups) {
				t.Errorf("shrink(%v, %v) got groups %v, want %v", tt.symbols, tt.groups, got.groups, tt.wantGroups)
			}
		})
	}
}
