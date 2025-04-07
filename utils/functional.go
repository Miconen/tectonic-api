package utils

func MapField[T any, P any](slice []T, mapFunc func(obj T) (P)) []P {
	result := make([]P, len(slice))
	for i := range slice {
		result = append(result, mapFunc(slice[i]))
	}

	return result
}
