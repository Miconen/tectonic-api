package utils

import (
	"encoding/json"
	"net/http"
)

type jsonWriter struct {
	data any
}

func JsonWriter(data any) *jsonWriter {
	return &jsonWriter{
		data: data,
	}
}

func (jsonWriter *jsonWriter) IntoHTTP(status int) http.HandlerFunc {
    // TODO(robertoesteves13): Resolve json writes bug
    return func (w http.ResponseWriter, _ *http.Request) {
        w.WriteHeader(status)

        // We can pass the response writer directly because it won't write the
        // response if the marshalling had errors.
        enc := json.NewEncoder(w)
        err := enc.Encode(jsonWriter.data)

        if err != nil {
            // Can't write JSON, we have a big problem.
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
    }
}
