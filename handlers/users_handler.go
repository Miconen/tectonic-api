package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary Get multiple users
// @Description Get multiple users details by unique user Snowflakes (IDs)
// @Tags Users
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_ids query string false "User IDs"
// @Param wom_ids query string false "WOM IDs"
// @Param rsns query string false "RSNs"
// @Param user_ids query string false "User IDs"
// @Success 200 {object} models.Users
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users [GET]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	v := mux.Vars(r)
	p, err := utils.ParseParametersURL(r)
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	if len(p) > 0 {
		_, err = utils.RequireOne(p, "user_ids", "wom_ids", "rsns")
		if err != nil {
			status = http.StatusBadRequest
			utils.JsonWriter(err).IntoHTTP(status)(w, r)
			return
		}
	}

	users, err := database.SelectUsers(r.Context(), v["guild_id"], p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
		status = http.StatusInternalServerError
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(users).IntoHTTP(status)(w, r)
}
