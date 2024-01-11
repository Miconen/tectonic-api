package utils

import "errors"

// Extract a value from a map, returning a new map
func ExtractByClone(p map[string]string, key string) (string, map[string]string, error) {
	val, ok := p[key]
	if !ok {
		return "", p, errors.New("No value found for key")
	}

	// Create a copy of the original map without provided key value pair
	c := make(map[string]string)
	for k, v := range p {
		if k != key {
			c[k] = v
		}
	}

	return val, c, nil
}

// Extract a value from a map, mutating the map given
func ExtractByMutate(p map[string]string, key string) (string, error) {
	val, ok := p[key]
	if !ok {
		return "", errors.New("No value found for key")
	}

	delete(p, key)

	return val, nil
}
