package utils

func GetLastItems[T any](list []T, limit int) []T {
	if list == nil {
		return nil
	}

	if len(list) == 0 {
		return list
	}

	startIndex := Max(0, len(list)-limit)
	return list[startIndex:]
}

func GetFirstItems[T any](list []T, limit int) []T {
	if list == nil {
		return nil
	}

	if len(list) == 0 {
		return list
	}

	endIndex := Min(len(list), limit)
	return list[:endIndex]
}

func GetMidItems[T any](list []T, limit int, offset int) []T {
	if len(list) == 0 {
		return list
	}

	startIndex := min(offset, len(list))
	endIndex := min(offset+limit, len(list))
	return list[startIndex:endIndex]
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
