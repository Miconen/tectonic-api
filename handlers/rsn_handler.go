package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/utils"
)

// @Summary Get RSN related information by guild and user ID
// @Description Get RSN related details by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Param user_id query string false "User ID"
// @Success 200 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /api/v1/rsn [GET]
func GetRSN(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p, err := utils.ParseParametersURL(r, "guild_id", "user_id")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	rsn, err := database.SelectRsn(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting RSN: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(rsn).IntoHTTP(status)(w, r)
}

// @Summary Link an RSN to a user
// @Description Link an RSN to a guild and user in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Param rsn path string true "RSN"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 429 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/rsn [POST]
func CreateRSN(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	p, err := utils.ParseParametersURL(r, "guild_id", "user_id", "rsn")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	err = database.InsertRsn(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating RSN: %v\n", err)
		// TODO: Handle 404 Not Found errors
		status = http.StatusConflict
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Remove RSN from guild and user
// @Description Delete a RSN in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Param rsn path string true "RSN"
// @Success 200 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /api/v1/rsn [DELETE]
func RemoveRSN(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p, err := utils.ParseParametersURL(r, "guild_id", "user_id", "rsn")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	err = database.DeleteRsn(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting RSN: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
