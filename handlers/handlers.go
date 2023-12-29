package handlers

import (
	"encoding/json"
	"net/http"
	"tectonic-api/utils"
)

// Guild Model
// @Description Model of guild data
type Guild struct {
	GuildId     string `json:"guild_id"`
	Multiplier  int    `json:"multiplier"`
	PbChannelId string `json:"pb_channel_id"`
}

// User Model
// @Description Model of active guild member
type User struct {
	UserId  string `json:"user_id"`
	GuildId string `json:"guild_id"`
	Points  int    `json:"points"`
}

// Users Model
// @Description Model of active guild members
type Users struct {
	Users []User `json:"users"`
}

// Time Model
// @Description Model of a fetched time and the team
type Time struct {
	Time     int    `json:"time"`
	BossName string `json:"boss_name"`
	RunId    int    `json:"run_id"`
	Date     int    `json:"date"`
	Team     Users  `json:"team"`
}

// Error Model
// @Description HTTP Error model with content, error and code
type Error struct {
	Content string `json:"content"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

// Function types for method-specific logic
type Handler func(r *http.Request) (interface{}, error)

// Helper function to write a 500 Internal Server Error response with error message
func writeInternalServerError(w http.ResponseWriter, errMsg string) {
	errorResponse := Error{
		Content: errMsg,
		Error:   "Internal Server Error",
		Code:    http.StatusInternalServerError,
	}
	writeErrorResponse(w, errorResponse)
}

// Helper function to write an error response in JSON format
func writeErrorResponse(w http.ResponseWriter, err Error) {
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

		errorResponse := Error{
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

	err = utils.ValidateStruct(res)
	if err != nil {
		errorResponse := Error{
			Content: "User not found",
			Error:   "Not Found",
			Code:    http.StatusNotFound,
		}
		writeErrorResponse(w, errorResponse)
		return
	}

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
