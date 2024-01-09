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
// @Success 200 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /v1/users [GET]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"user_ids": r.URL.Query().Get("user_ids"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		users, err := database.SelectUsers("users", p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
			return models.Body{}, http.StatusInternalServerError, err
		}

		return models.Body{Content: users}, http.StatusOK, nil
	}

	httpHandler(w, r, h, p)
}
