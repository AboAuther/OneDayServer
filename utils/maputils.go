package utils

import "sort"

func GetMapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func GetSortedMapValues[K comparable, V any](m map[K]V, compare func(a, b V) bool) []V {
	values := GetMapValues(m)
	sort.Slice(values, func(i, j int) bool {
		return compare(values[i], values[j])
	})
	return values
}
