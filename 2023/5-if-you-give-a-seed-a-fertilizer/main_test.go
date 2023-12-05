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
			name:  "example 1",
			input: example,
			want:  35,
		},
		{
			name:  "actual",
			input: input,
			want:  331445006,
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
			want:  46,
		},
		{
			name:  "actual",
			input: input,
			want:  6472060,
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

func Test_convertRange(t *testing.T) {
	tests := []struct {
		name string
		sr   SeedRange
		cg   ConvertGroup
		want []SeedRange
	}{
		{
			"range before convert map",
			SeedRange{50, 60},
			ConvertGroup{"", []ConvertMap{{100, 110, -100}}},
			[]SeedRange{{50, 60}},
		},
		{
			"range after convert map",
			SeedRange{250, 260},
			ConvertGroup{"", []ConvertMap{{100, 110, -100}}},
			[]SeedRange{{250, 260}},
		},
		{
			"range inside convert map",
			SeedRange{105, 110},
			ConvertGroup{"", []ConvertMap{{100, 110, -100}}},
			[]SeedRange{{5, 10}},
		},
		{
			"range covert first part of convert map",
			SeedRange{90, 105},
			ConvertGroup{"", []ConvertMap{{100, 110, -100}}},
			[]SeedRange{{0, 5}, {90, 100}},
		},
		{
			"range covert second part of convert map",
			SeedRange{105, 115},
			ConvertGroup{"", []ConvertMap{{100, 110, -100}}},
			[]SeedRange{{5, 10}, {110, 115}},
		},
		{
			"range covers convert map",
			SeedRange{50, 150},
			ConvertGroup{"", []ConvertMap{{100, 110, -100}}},
			[]SeedRange{{0, 10}, {50, 100}, {110, 150}},
		},

		{
			"real example",
			SeedRange{55, 68},
			ConvertGroup{"", []ConvertMap{
				{53, 61, -4},
				{11, 53, -11},
				{0, 7, 42},
				{7, 11, 50}}},
			[]SeedRange{{51, 57}, {61, 68}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cg.convertRange(tt.sr)
			if len(got) != len(tt.want) {
				t.Errorf("ConvertGroup(%v).split(%v) = %v, want %v", tt.cg, tt.sr, got, tt.want)
				return
			}

			for i, _ := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ConvertGroup(%v).split(%v) = %v, want %v", tt.cg, tt.sr, got, tt.want)
				}
			}
		})
	}
}
