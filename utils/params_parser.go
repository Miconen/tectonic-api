package utils

import (
	"errors"
	"fmt"
	"net/http"
)

// Helper function to check if all parameters were set by the client. It errors
// if one of more params are empty.
func ParseParametersURL(r *http.Request, params ...string) (map[string]string, error) {
	parsed := make(map[string]string)
	errs := make([]error, 0, len(params))

	for i := range params {
		value := r.URL.Query().Get(params[i])
		if value == "" {
			errs = append(errs, fmt.Errorf("parameter `%s` is empty", params[i]))
		}

		parsed[params[i]] = value
	}

	return parsed, errors.Join(errs...)
}
