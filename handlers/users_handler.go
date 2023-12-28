package handlers

import "net/http"

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
// @Router /api/v1/users [GET]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
// @Router /api/v1/users [PUT]
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
