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
// @Success 200 {object} User
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/user [GET]
func GetUser(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"user_id":  r.URL.Query().Get("user_id"),
	}

	h := func(r *http.Request) (interface{}, error) {

		user, err := database.FetchUser(p["guild_id"], p["user_id"])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching users: %v\n", err)
			return nil, err
		}

		return user, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Create / Initialize a puild
// @Description Initialize a guild in our backend by unique guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} User
// @Success 201 {object} User
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/user [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"user_id":  r.URL.Query().Get("user_id"),
	}

	h := func(r *http.Request) (interface{}, error) {

		err := database.InsertUser(p["guild_id"], p["user_id"])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error inserting user: %v\n", err)
			return nil, err
		}

		return models.User{GuildId: p["guild_id"], UserId: p["user_id"], Points: 0}, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Update a guilds information
// @Description Update a guild in our backend by unique guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} User
// @Success 201 {object} User
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/user [PUT]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"user_id":  r.URL.Query().Get("user_id"),
	}

	h := func(r *http.Request) (interface{}, error) {

		user := models.User{
			UserId:  p["user_id"],
			GuildId: p["guild_id"],
			Points:  789,
		}

		return user, nil
	}

	httpHandler(w, r, h, p)
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
// @Router /v1/user [DELETE]
func RemoveUser(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"user_id":  r.URL.Query().Get("user_id"),
	}

	h := func(r *http.Request) (interface{}, error) {

		user := models.User{
			UserId:  p["user_id"],
			GuildId: p["guild_id"],
			Points:  789,
		}

		return user, nil
	}

	httpHandler(w, r, h, p)
}
