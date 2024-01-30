package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary Get a guilds leaderboard by ID
// @Description Get guilds leaderboard details by unique guild Snowflake (ID)
// @Tags Leaderboard
// @Produce json
// @Param guild_id path string false "Guild ID"
// @Success 200 {object} models.Users
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/leaderboard [GET]
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p := mux.Vars(r)

	leaderboard, err := database.SelectLeaderboard(r.Context(), p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(leaderboard).IntoHTTP(status)(w, r)
}
