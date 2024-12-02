package utils

import (
	"reflect"
	"testing"
)

func TestGetLastItems(t *testing.T) {
	tests := []struct {
		input    []int
		limit    int
		expected []int
	}{
		{input: []int{1, 2, 3, 4, 5}, limit: 3, expected: []int{3, 4, 5}},
		{input: []int{1, 2, 3, 4, 5}, limit: 10, expected: []int{1, 2, 3, 4, 5}},
		{input: []int{1, 2, 3, 4, 5}, limit: 0, expected: []int{}},
		{input: []int{}, limit: 3, expected: []int{}},
	}

	for _, test := range tests {
		output := GetLastItems(test.input, test.limit)
		if !reflect.DeepEqual(output, test.expected) {
			t.Errorf("GetLastItems(%v, %d) = %v; expected %v", test.input, test.limit, output, test.expected)
		}
	}
}

func TestGetFirstItems(t *testing.T) {
	tests := []struct {
		name     string
		list     []int
		limit    int
		expected []int
	}{
		{"Nil list", nil, 3, nil},
		{"Empty list", []int{}, 3, []int{}},
		{"Limit larger than list length", []int{1, 2, 3}, 5, []int{1, 2, 3}},
		{"Limit smaller than list length", []int{1, 2, 3, 4, 5}, 3, []int{1, 2, 3}},
		{"Limit equals list length", []int{1, 2, 3}, 3, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFirstItems(tt.list, tt.limit)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetFirstItems() = %v; want %v", result, tt.expected)
			}
		})
	}
}
