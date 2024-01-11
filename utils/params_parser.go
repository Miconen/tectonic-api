package utils

import (
	"errors"
	"fmt"
	"net/http"
)

// Helper function to check if required parameters were set by the client, and
// it will error if either the parameter is defined none or more than once.
//
// Note that error implements Unwrap() if it's not nill, meaning you can get
// all errors that ocurred inside the function and send back to the client.
func ParseParametersURL(r *http.Request, required ...string) (map[string]string, error) {
	parsed := make(map[string]string)
	errs := make([]error, 0, len(required))

	for k, v := range r.URL.Query() {
		if len(v) > 1 {
			errs = append(errs, fmt.Errorf("`%s` was defined more than once", k))
		}

		if v[0] != "" {
			parsed[k] = v[0]
		}
	}

	for i := range required {
		if _, ok := parsed[required[i]]; !ok {
			errs = append(errs, fmt.Errorf("`%s` not set", required[i]))
		}
	}

	return parsed, errors.Join(errs...)
}
