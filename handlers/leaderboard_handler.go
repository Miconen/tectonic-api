package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/models"
)

// @Summary Get a guilds leaderboard by ID
// @Description Get guilds leaderboard details by unique guild Snowflake (ID)
// @Tags Leaderboard
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 429 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /v1/leaderboard [GET]
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		leaderboard, err := database.SelectLeaderboard(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
			return models.Body{}, http.StatusInternalServerError, err
		}

		return models.Body{Content: leaderboard}, http.StatusOK, nil
	}

	httpHandler(w, r, h, p)
}
