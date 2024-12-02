package utils

import "testing"

func TestSortQueryString(t *testing.T) {
	tests := []struct {
		rawQuery string
		expected string
	}{
		{"b=2&a=1", "a=1&b=2"},
		{"a=1&b=2", "a=1&b=2"},
		{"b=2&b=1&a=3", "a=3&b=1&b=2"},
		{"", ""},
		{"z=last&a=first&c=middle", "a=first&c=middle&z=last"},
		{"c=3&b=2&a=1", "a=1&b=2&c=3"},
	}

	for _, test := range tests {
		result := SortQueryString(test.rawQuery)
		if result != test.expected {
			t.Errorf("SortQueryString(%q) = %q; expected %q", test.rawQuery, result, test.expected)
		}
	}
}
