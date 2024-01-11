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

func (jsonWriter *jsonWriter) ServeHttp(w http.ResponseWriter, r *http.Request) {
	// We can pass the response writer directly because it won't write the
	// response if the marshalling had errors.
	enc := json.NewEncoder(w)
	err := enc.Encode(jsonWriter.data)

	if err != nil {
		// Can't write JSON, we have a big problem.
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
