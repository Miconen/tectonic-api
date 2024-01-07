package handlers

import (
	"net/http"
	"tectonic-api/models"
)

// @Summary Get a guilds leaderboard by ID
// @Description Get guilds leaderboard details by unique guild Snowflake (ID)
// @Tags Leaderboard
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Success 200 {object} Users
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/leaderboard [GET]
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
	}

	h := func(r *http.Request) (interface{}, error) {

		users := models.Users{}
		for i := 0; i < 10; i++ {
			user := models.User{
				UserId:  "Hello World",
				GuildId: p["guild_id"],
				Points:  789,
			}

			users.Users = append(users.Users, user)
		}

		return users, nil
	}

	httpHandler(w, r, h, p)
}
