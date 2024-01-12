package utils

import (
	"errors"
	"fmt"
	"strings"
)

func RequireOne(m map[string]string, opts ...string) (string, error) {
	var foundKeys []string

	for _, key := range opts {
		if value, ok := m[key]; ok && value != "" {
			foundKeys = append(foundKeys, key)
		}
	}

	if len(foundKeys) == 0 {
		err := errors.New("Couldn't find any opts from map")
		return "", err
	}

	if len(foundKeys) > 1 {
		err := fmt.Errorf("Multiple opts set: %s", strings.Join(foundKeys, ", "))
		return "", err
	}

	return foundKeys[0], nil
}
