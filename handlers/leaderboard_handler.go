package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/utils"
)

// @Summary Get a guilds leaderboard by ID
// @Description Get guilds leaderboard details by unique guild Snowflake (ID)
// @Tags Leaderboard
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Success 200 {object} models.Users
// @Failure 400 {object} models.Empty
// @Failure 403 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/leaderboard [GET]
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p, err := utils.ParseParametersURL(r, "guild_id")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	leaderboard, err := database.SelectLeaderboard(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(leaderboard).IntoHTTP(status)(w, r)
}
