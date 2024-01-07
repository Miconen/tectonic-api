package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/models"
)

// @Summary Get a user by ID
// @Description Get user details by unique user Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Param user_id query string false "User ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 429 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /v1/user [GET]
func GetUser(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"user_id":  r.URL.Query().Get("user_id"),
	}

	h := func(r *http.Request) (models.Body, int, error) {
		user, err := database.SelectUser(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
			return models.Body{}, http.StatusInternalServerError, err
		}

		res := models.Body{Content: user}

		return res, http.StatusOK, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Create / Initialize a new user
// @Description Initialize a guild in our backend by unique guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 429 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /v1/user [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"user_id":  r.URL.Query().Get("user_id"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		user := models.User{
			UserId:  p["user_id"],
			GuildId: p["guild_id"],
		}

		err := database.InsertUser(user)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error inserting user: %v\n", err)
			return models.Body{}, http.StatusInternalServerError, err
		}

		return models.Body{}, http.StatusCreated, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Delete a user from guild
// @Description Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 204 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 429 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /v1/user [DELETE]
func RemoveUser(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"user_id":  r.URL.Query().Get("user_id"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		err := database.DeleteUser(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting user: %v\n", err)
			return models.Body{}, http.StatusInternalServerError, err
		}

		return models.Body{}, http.StatusNoContent, nil
	}

	httpHandler(w, r, h, p)
}
