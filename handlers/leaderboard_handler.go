package handlers

import (
	"net/http"
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
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	params := database.GetLeaderboardParams{
		GuildID:   p["guild_id"],
		UserLimit: 50,
	}

	log.DebugContext(r.Context(), "querying leaderboard from database", "guild_id", params.GuildID, "user_limit", params.UserLimit)
	rows, err := queries.GetLeaderboard(r.Context(), params)
	if err != nil {
		log.Error("Error fetching users", "error", err)
		jw.SetStatus(http.StatusNotFound)
		jw.WriteResponse(http.NoBody)
		return
	}

	leaderboard := database.NewLeaderboardFromRows(rows)
	jw.WriteResponse(leaderboard)
}
