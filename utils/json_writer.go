package utils

import (
	"encoding/json"
	"net/http"
)

type JsonWriter struct {
	statusCode int
	w          http.ResponseWriter
	r          *http.Request
}

func NewJsonWriter(w http.ResponseWriter, r *http.Request, statusCode int) *JsonWriter {
	return &JsonWriter{
		w:          w,
		r:          r,
		statusCode: statusCode,
	}
}

func (jw *JsonWriter) SetStatus(statusCode int) {
	jw.statusCode = statusCode
}

func (jw *JsonWriter) WriteResponse(body any) {
	log.Debug("writing http response", "body", body, "status", jw.statusCode)

	jw.w.Header().Set("Content-Type", "application/json")
	jw.w.WriteHeader(jw.statusCode)

	if body != nil && body != http.NoBody {
		enc := json.NewEncoder(jw.w)
		err := enc.Encode(body)

		if err != nil {
			log.Error("failed to write response", "error", err)
			http.Error(jw.w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
