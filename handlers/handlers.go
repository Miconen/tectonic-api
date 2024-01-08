package handlers

import (
	"encoding/json"
	"net/http"
	"tectonic-api/models"
	"tectonic-api/utils"
)

// Function types for method-specific logic
type Handler func(r *http.Request) (models.Body, int, error)

// Helper function to write an error response in JSON format
func writeErrorResponse(w http.ResponseWriter, e string, s int) {
	err := models.Body{Content: e}
	w.WriteHeader(s)
	json.NewEncoder(w).Encode(err)
}

// func get404() {
// 	err = utils.ValidateStruct(res)
// 	if err != nil {
// 		errorResponse := models.Response{
// 			Content: "User not found",
// 			Code:    http.StatusNotFound,
// 		}
// 		writeErrorResponse(w, errorResponse)
// 		return
// 	}
// }

func httpHandler(w http.ResponseWriter, r *http.Request, h Handler, p map[string]string) {
	// Check if required parameters are missing based on the handler logic
	for _, v := range p {
		if v != "" {
			continue
		}

		writeErrorResponse(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Call the method-specific logic function and handle errors
	res, code, err := h(r)
	if err != nil {
		writeErrorResponse(w, err.Error(), code)
		return
	}

	// Write JSON response
    utils.JsonWriter(res).ServeHttp(w, r)
}
