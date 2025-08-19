package utils

func MapField[T any, P any](slice []T, mapFunc func(obj T) P) []P {
	result := make([]P, 0, len(slice))
	for _, v := range slice {
		result = append(result, mapFunc(v))
	}
	return result
}

func Filter[T any](slice []T, filterFunc func(elem T) bool) []T {
	result := make([]T, len(slice))
	for i := range slice {
		if filterFunc(slice[i]) {
			result = append(result, slice[i])
		}
	}

	return result
}
