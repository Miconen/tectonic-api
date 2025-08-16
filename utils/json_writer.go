package utils

import (
	"encoding/json"
	"net/http"
	"tectonic-api/logging"
	"tectonic-api/models"
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
	logging.Get().Debug("writing http response", "body", body, "status", jw.statusCode)

	if body == nil || body == http.NoBody {
		body = struct{}{}
	}

	if jw.statusCode != 204 {
		jw.w.Header().Set("Content-Type", "application/json")
		jw.w.WriteHeader(jw.statusCode)
		enc := json.NewEncoder(jw.w)
		err := enc.Encode(body)

		if err != nil {
			logging.Get().Error("failed to write response", "error", err)
			http.Error(jw.w, "Internal server error", http.StatusInternalServerError)
		}
	} else {
		jw.w.WriteHeader(http.StatusNoContent)
	}
}

func (jw *JsonWriter) WriteError(code models.APIV1Error) {
	jw.SetStatus(code.Status())
	jw.WriteResponse(code.ToErrorResponse())
}
