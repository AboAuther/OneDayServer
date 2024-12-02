package utils

import (
	"reflect"
	"testing"
)

func TestGetMapValues(t *testing.T) {
	tests := []struct {
		input    map[string]int
		expected []int
	}{
		{input: map[string]int{"a": 1, "b": 2, "c": 3}, expected: []int{1, 2, 3}},
		{input: map[string]int{"x": 10, "y": 20}, expected: []int{10, 20}},
		{input: map[string]int{}, expected: []int{}},
	}

	for _, test := range tests {
		output := GetMapValues(test.input)
		if !isEquivalent(output, test.expected) {
			t.Errorf("GetMapValues(%v) = %v; expected %v", test.input, output, test.expected)
		}
	}
}

func isEquivalent(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	mapA := make(map[int]int)
	mapB := make(map[int]int)

	for _, v := range a {
		mapA[v]++
	}
	for _, v := range b {
		mapB[v]++
	}

	return reflect.DeepEqual(mapA, mapB)
}

func TestGetSortedMapValues(t *testing.T) {
	m := map[string]int{
		"one":   1,
		"three": 3,
		"two":   2,
		"four":  4,
	}
	expectedValues := []int{1, 2, 3, 4}

	compare := func(a, b int) bool {
		return a < b
	}

	sortedValues := GetSortedMapValues(m, compare)

	if !reflect.DeepEqual(sortedValues, expectedValues) {
		t.Errorf("GetSortedMapValues() = %v; want %v", sortedValues, expectedValues)
	}
}
