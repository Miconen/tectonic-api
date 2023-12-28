package handlers

import "net/http"

// @Summary Get a user by ID
// @Description Get user details by unique user Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Param user_id query string false "User ID"
// @Success 200 {object} User
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /api/v1/user [GET]
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

// @Summary Create / Initialize a guild
// @Description Initialize a guild in our backend by unique guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} User
// @Success 201 {object} User
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /api/v1/user [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

// @Summary Update a guilds information
// @Description Update a guild in our backend by unique guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} User
// @Success 201 {object} User
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /api/v1/user [PUT]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

// @Summary Delete a user from guild
// @Description Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /api/v1/user [DELETE]
func RemoveUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
