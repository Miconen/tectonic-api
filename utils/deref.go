package utils

func DerefOr[T any](p *T, fallback T) T {
	if p != nil {
		return *p
	}
	return fallback
}
