package handlers

import (
	"encoding/json"
	"net/http"
	"tectonic-api/models"
	// "tectonic-api/utils"
)

// Function types for method-specific logic
type Handler func(r *http.Request) (interface{}, error)

// Helper function to write a 500 Internal Server Error response with error message
func writeInternalServerError(w http.ResponseWriter, errMsg string) {
	errorResponse := models.Error{
		Content: errMsg,
		Error:   "Internal Server Error",
		Code:    http.StatusInternalServerError,
	}
	writeErrorResponse(w, errorResponse)
}

// Helper function to write an error response in JSON format
func writeErrorResponse(w http.ResponseWriter, err models.Error) {
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}

func httpHandler(w http.ResponseWriter, r *http.Request, h Handler, p map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	// Check if required parameters are missing based on the handler logic
	for _, v := range p {
		if v != "" {
			continue
		}

		errorResponse := models.Error{
			Content: "Missing required parameters",
			Error:   "Bad Request",
			Code:    http.StatusBadRequest,
		}
		writeErrorResponse(w, errorResponse)
		return
	}

	// Call the method-specific logic function and handle errors
	res, err := h(r)
	if err != nil {
		writeInternalServerError(w, err.Error())
		return
	}

	// err = utils.ValidateStruct(res)
	// if err != nil {
	// 	errorResponse := models.Error{
	// 		Content: "User not found",
	// 		Error:   "Not Found",
	// 		Code:    http.StatusNotFound,
	// 	}
	// 	writeErrorResponse(w, errorResponse)
	// 	return
	// }

	// Marshal response data into JSON
	userJSON, err := json.Marshal(res)
	if err != nil {
		writeInternalServerError(w, err.Error())
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}
