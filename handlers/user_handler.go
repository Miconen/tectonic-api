package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/utils"
)

// @Summary Get a user by ID
// @Description Get user details by unique user Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Param user_id query string false "User ID"
// @Success 200 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /v1/user [GET]
func GetUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p, err := utils.ParseParametersURL(r, "guild_id", "user_id")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	user, err := database.SelectUser(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching user: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(user).IntoHTTP(status)(w, r)
}

// @Summary Create / Initialize a new user
// @Description Initialize a guild in our backend by unique guild Snowflake (ID)
// @Tags User
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
// @Router /v1/user [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	p, err := utils.ParseParametersURL(r, "guild_id", "user_id", "rsn")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	wid, err := utils.GetWomId(p["rsn"])
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	err = database.InsertUser(p, wid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting user: %v\n", err)
		status = http.StatusConflict
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Delete a user from guild
// @Description Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 204 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /v1/user [DELETE]
func RemoveUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p, err := utils.ParseParametersURL(r, "guild_id", "user_id")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	err = database.DeleteUser(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting user: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
