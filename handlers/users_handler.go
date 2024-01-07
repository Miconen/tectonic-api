package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/models"
)

// @Summary Get multiple users
// @Description Get multiple users details by unique user Snowflakes (IDs)
// @Tags Users
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Param user_ids query string false "User IDs"
// @Success 200 {object} Users
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/users [GET]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"user_ids": r.URL.Query().Get("user_ids"),
	}

	h := func(r *http.Request) (interface{}, error) {

		users, err := database.FetchUsers(p["guild_id"])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
			return nil, err
		}

		return users, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Update multiple users information
// @Description Update multiple users by unique user Snowflakes (IDs)
// @Tags Users
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_ids path string true "User IDs"
// @Success 200 {object} Users
// @Success 201 {object} Users
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/users [PUT]
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"user_ids": r.URL.Query().Get("user_ids"),
	}

	h := func(r *http.Request) (interface{}, error) {

		users := models.Users{}
		for i := 0; i < 10; i++ {
			user := models.User{
				UserId:  p["user_ids"],
				GuildId: p["guild_id"],
				Points:  789,
			}

			users.Users = append(users.Users, user)
		}

		return users, nil
	}

	httpHandler(w, r, h, p)
}
