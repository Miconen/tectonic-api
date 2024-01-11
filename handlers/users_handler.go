package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/utils"
)

// @Summary Get multiple users
// @Description Get multiple users details by unique user Snowflakes (IDs)
// @Tags Users
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Param user_ids query string false "User IDs"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 429 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /v1/users [GET]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p, err := utils.ParseParametersURL(r, "guild_id", "user_ids")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	users, err := database.SelectUsers(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
		status = http.StatusInternalServerError
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(users).IntoHTTP(status)(w, r)
}
